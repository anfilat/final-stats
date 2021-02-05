package main

import (
	"context"
	"sync"
	"time"
)

type serviceStop func(ctx context.Context)

type serviceStopper struct {
	list []serviceStop
}

func newServiceStopper() *serviceStopper {
	return &serviceStopper{}
}

func (s *serviceStopper) add(closeFn serviceStop) {
	s.list = append(s.list, closeFn)
}

func (s *serviceStopper) stop() {
	softCtx, softCancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer softCancel()
	hardCtx, hardCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer hardCancel()

	wg := &sync.WaitGroup{}

	for _, closeFn := range s.list {
		wg.Add(1)
		go func(closeFn serviceStop) {
			defer wg.Done()

			closeFn(softCtx)
		}(closeFn)
	}

	stopped := make(chan interface{})
	go func() {
		wg.Wait()
		close(stopped)
	}()

	select {
	case <-hardCtx.Done():
	case <-stopped:
	}
}
