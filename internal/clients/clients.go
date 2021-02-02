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
	log           symo.Logger
}

func NewClients(ctx context.Context, log symo.Logger,
	toHeartChan chan<- symo.HeartCommand, toClientsChan <-chan symo.ClientsBeat) symo.Clients {
	return &clients{
		ctx:           ctx,
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
			wg.Done()
		}()

		for {
			select {
			case <-c.ctx.Done():
				return
			case _, ok := <-c.toClientsChan:
				if !ok {
					return
				}
				// TODO
			}
		}
	}()
}
