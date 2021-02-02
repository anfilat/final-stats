package clients

import (
	"context"
	"sync"

	"github.com/anfilat/final-stats/internal/symo"
)

type clients struct {
	ctx           context.Context
	toHeartChan   chan<- symo.HeartCommand
	toClientsChan <-chan symo.ClientsBeat
	clientsMutex  *sync.Mutex
	clients       []*grpcClient
	log           symo.Logger
}

type grpcClient struct {
	ctx  context.Context
	n    int
	m    int
	ch   chan<- symo.Stat
	dead bool
}

func NewClients(ctx context.Context, log symo.Logger,
	toHeartChan chan<- symo.HeartCommand, toClientsChan <-chan symo.ClientsBeat) symo.Clients {
	return &clients{
		ctx:           ctx,
		toHeartChan:   toHeartChan,
		toClientsChan: toClientsChan,
		clientsMutex:  &sync.Mutex{},
		log:           log,
	}
}

func (c *clients) Start(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer func() {
			c.log.Debug("clients stooped")
			close(c.toHeartChan)
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

func (c *clients) NewClient(client symo.NewClient) <-chan symo.Stat {
	ch := make(chan symo.Stat, 1)
	cl := &grpcClient{
		ctx:  client.Ctx,
		n:    client.N,
		m:    client.M,
		ch:   ch,
		dead: false,
	}

	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	c.clients = append(c.clients, cl)

	if len(c.clients) == 1 {
		c.toHeartChan <- symo.Start
	}

	return ch
}

func (c *clients) sendStat(data *symo.ClientsBeat) {
	// TODO заменить на реальный код
	stat := &symo.Stat{
		Time: data.Time,
		Stat: &symo.Point{
			LoadAvg: &symo.LoadAvgData{
				Load1:  0,
				Load5:  0,
				Load15: 0,
			},
		},
	}

	c.sendToClients(stat)
	c.delDeadClients()
}

func (c *clients) sendToClients(stat *symo.Stat) {
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	for _, client := range c.clients {
		select {
		case <-c.ctx.Done():
			return
		case <-client.ctx.Done():
			client.dead = true
		case client.ch <- *stat:
		}
	}
}

func (c *clients) delDeadClients() {
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	clients := make([]*grpcClient, 0, len(c.clients))
	for _, client := range c.clients {
		if client.dead {
			continue
		}
		clients = append(clients, client)
	}
	c.clients = clients

	if len(c.clients) == 0 {
		c.toHeartChan <- symo.Stop
	}
}
