package clients

import (
	"context"
	"sync"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/anfilat/final-stats/internal/symo"
)

type clients struct {
	ctx         context.Context // управление остановкой сервиса
	ctxCancel   context.CancelFunc
	closedCh    chan interface{}
	mutex       *sync.Mutex
	clients     clientsList // список клиентов
	toClientsCh <-chan symo.MetricsData
	log         symo.Logger
	clock       clock.Clock
}

// NewClients возвращает сервис клиентов.
func NewClients(log symo.Logger, clock clock.Clock) symo.Clients {
	return &clients{
		log:   log,
		clock: clock,
	}
}

func (c *clients) Start(_ context.Context, toClientsCh <-chan symo.MetricsData) {
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

	now := c.clock.Now().Truncate(time.Second)
	client := newClient(cl, now)

	c.clients = append(c.clients, client)

	delClient := func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		client.dead = true
	}

	return client.ch, delClient, nil
}

func (c *clients) sendStat(data *symo.MetricsData) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.clients) == 0 {
		return
	}

	from := time.Now()
	defer c.log.Debug("stats sent in ", time.Since(from))

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
}
