package clients

import (
	"context"
	"sync"
	"time"

	"github.com/anfilat/final-stats/internal/pb"
	"github.com/anfilat/final-stats/internal/symo"
)

type clients struct {
	ctx           context.Context // контекст приложения, сервис завершается по закрытию контекста
	clientsMutex  *sync.Mutex
	clients       clientsList // список клиентов
	toHeartChan   chan<- symo.HeartCommand
	toClientsChan <-chan symo.ClientsBeat
	log           symo.Logger
}

func NewClients(ctx context.Context, log symo.Logger,
	toHeartChan chan<- symo.HeartCommand, toClientsChan <-chan symo.ClientsBeat) symo.Clients {
	return &clients{
		ctx:           ctx,
		clientsMutex:  &sync.Mutex{},
		toHeartChan:   toHeartChan,
		toClientsChan: toClientsChan,
		log:           log,
	}
}

func (c *clients) Start(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer func() {
			close(c.toHeartChan)
			c.log.Debug("clients is stopped")
			wg.Done()
		}()

		for {
			select {
			case <-c.ctx.Done():
				return
			case data, ok := <-c.toClientsChan:
				if !ok {
					return
				}
				c.sendStat(&data)
			}
		}
	}()
}

// подключение нового клиента.
func (c *clients) NewClient(cl symo.NewClient) (<-chan *pb.Stats, func()) {
	client := newClient(cl)

	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	c.clients = append(c.clients, client)

	if len(c.clients) == 1 {
		select {
		case c.toHeartChan <- symo.Start:
		default:
		}
	}

	return client.ch, func() {
		c.clientsMutex.Lock()
		defer c.clientsMutex.Unlock()

		client.dead = true
	}
}

func (c *clients) sendStat(data *symo.ClientsBeat) {
	from := time.Now()
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	now := data.Time
	results := make(map[int]*pb.Stats)

	clients := make(clientsList, 0, len(c.clients))
	for _, client := range c.clients {
		if client.dead {
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
		case c.toHeartChan <- symo.Stop:
		default:
		}
	}
	c.log.Info(time.Since(from))
}
