package clients

import (
	"context"
	"sync"
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

type clients struct {
	ctx           context.Context // управление остановкой сервиса
	ctxCancel     context.CancelFunc
	closedCh      chan interface{}
	mutex         *sync.Mutex
	clients       clientsList // список клиентов
	toCollectorCh chan<- symo.CollectorCommand
	toClientsCh   <-chan symo.MetricsData
	log           symo.Logger
}

func NewClients(log symo.Logger) symo.Clients {
	return &clients{
		log: log,
	}
}

func (c *clients) Start(_ context.Context, toCollectorCh chan<- symo.CollectorCommand, toClientsCh <-chan symo.MetricsData) {
	c.toCollectorCh = toCollectorCh
	c.toClientsCh = toClientsCh

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())
	c.closedCh = make(chan interface{})
	c.mutex = &sync.Mutex{}
	c.clients = nil

	go c.work()
}

func (c *clients) Stop(ctx context.Context) {
	c.ctxCancel()

	select {
	case <-ctx.Done():
		return
	case <-c.closedCh:
	}

	c.closeClients()

	c.log.Debug("clients is stopped")
}

func (c *clients) closeClients() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, client := range c.clients {
		client.close()
	}
}

func (c *clients) work() {
	defer func() {
		c.mutex.Lock()
		close(c.closedCh)
		c.mutex.Unlock()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		select {
		case <-c.ctx.Done():
			return
		case data := <-c.toClientsCh:
			c.sendStat(&data)
		}
	}
}

// подключение нового клиента.
func (c *clients) NewClient(cl symo.ClientData) (<-chan *symo.Stats, func(), error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	select {
	case <-c.closedCh:
		return nil, nil, symo.ErrStopped
	default:
	}

	client := newClient(cl)

	c.clients = append(c.clients, client)

	if len(c.clients) == 1 {
		select {
		case c.toCollectorCh <- symo.Start:
		default:
		}
	}

	delClient := func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		client.dead = true
	}

	return client.ch, delClient, nil
}

func (c *clients) sendStat(data *symo.MetricsData) {
	from := time.Now()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := data.Time
	results := make(map[int]*symo.Stats)

	clients := make(clientsList, 0, len(c.clients))
	for _, client := range c.clients {
		if client.dead {
			client.close()
			continue
		}
		clients = append(clients, client)

		if !client.isReady(now) {
			continue
		}
		client.setNextReady(now)

		stats, ok := results[client.m]
		if !ok {
			stats = makeSnapshot(data, client.m)
			results[client.m] = stats
		}

		select {
		case client.ch <- stats:
		default:
		}
	}
	c.clients = clients

	if len(c.clients) == 0 {
		select {
		case c.toCollectorCh <- symo.Stop:
		default:
		}
	}
	c.log.Debug("stats sent in ", time.Since(from))
}
