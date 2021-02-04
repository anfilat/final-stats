package clients

import (
	"context"
	"sync"
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

type clients struct {
	ctx           context.Context // контекст приложения, сервис завершается по закрытию контекста
	clientsMutex  *sync.Mutex
	clients       clientsList // список клиентов
	isDead        bool
	toHeartChan   chan<- symo.HeartCommand
	toClientsChan <-chan symo.ClientsBeat
	log           symo.Logger
}

func NewClients(ctx context.Context, log symo.Logger,
	toHeartChan chan<- symo.HeartCommand, toClientsChan <-chan symo.ClientsBeat) symo.Clients {
	return &clients{
		ctx:           ctx,
		clientsMutex:  &sync.Mutex{},
		isDead:        false,
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
func (c *clients) NewClient(cl symo.NewClient) (<-chan *symo.Stat, func()) {
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
		c.isDead = true
	}
}

func (c *clients) sendStat(data *symo.ClientsBeat) {
	res := c.filterReadyClients(data.Time)

	for m, list := range res {
		c.sendToClients(list, makeSnapshot(data, m))
	}

	c.delDeadClients()
}

// возвращает клиентов, которым надо отправить данные. Клиенты сгруппированы по M.
// У возвращенных клиентов устанавливается время следующей отправки.
func (c *clients) filterReadyClients(now time.Time) map[int]clientsList {
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	result := make(map[int]clientsList, len(c.clients))
	for _, client := range c.clients {
		if client.isReady(now) {
			result[client.m] = append(result[client.m], client)
			client.setNextReady(now)
		}
	}
	return result
}

func (c *clients) sendToClients(list clientsList, stat *symo.Stat) {
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	for _, client := range list {
		if client.dead {
			continue
		}

		select {
		case client.ch <- stat:
		default:
		}
	}
}

func (c *clients) delDeadClients() {
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	if !c.isDead {
		return
	}

	clients := make(clientsList, 0, len(c.clients))
	for _, client := range c.clients {
		if client.dead {
			continue
		}
		clients = append(clients, client)
	}
	c.clients = clients

	if len(c.clients) == 0 {
		select {
		case c.toHeartChan <- symo.Stop:
		default:
		}
	}
}
