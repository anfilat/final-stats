package heart

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/anfilat/final-stats/internal/loadavg"
	"github.com/anfilat/final-stats/internal/symo"
)

type heart struct {
	ctx          context.Context
	pointsMutex  *sync.Mutex
	points       symo.Points
	workerChains []chan<- symo.Beat
	log          symo.Logger
}

func NewHeart(ctx context.Context, log symo.Logger) symo.Heart {
	return &heart{
		ctx:         ctx,
		pointsMutex: &sync.Mutex{},
		points:      make(symo.Points),
		log:         log,
	}
}

func (h *heart) Start(wg *sync.WaitGroup, config symo.MetricConf) {
	wg.Add(1)
	h.mountMetrics(config)

	go func() {
		defer func() {
			h.unmountMetrics()
			wg.Done()
		}()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-h.ctx.Done():
				return
			case now := <-ticker.C:
				h.work(now)
			}
		}
	}()
}

func (h *heart) mountMetrics(config symo.MetricConf) {
	if config.Loadavg {
		h.loadavg(h.newWorkerChan())
	}
}

func (h *heart) unmountMetrics() {
	for _, ch := range h.workerChains {
		close(ch)
	}
}

func (h *heart) newWorkerChan() chan symo.Beat {
	ch := make(chan symo.Beat, 1)
	h.workerChains = append(h.workerChains, ch)

	return ch
}

func (h *heart) work(now time.Time) {
	point := h.newPoint(now)

	for _, ch := range h.workerChains {
		ch <- symo.Beat{
			Time:  now,
			Point: point,
		}
	}

	h.cleanPoints(now)
}

func (h *heart) newPoint(now time.Time) *symo.Point {
	h.pointsMutex.Lock()
	defer h.pointsMutex.Unlock()

	point := &symo.Point{}
	h.points[now] = point

	return point
}

func (h *heart) cleanPoints(now time.Time) {
	h.pointsMutex.Lock()
	defer h.pointsMutex.Unlock()

	limit := now.Add(-symo.MaxOldPoints)

	for key := range h.points {
		if key.Before(limit) {
			delete(h.points, key)
		}
	}
}

func (h *heart) loadavg(ch <-chan symo.Beat) {
	go func() {
		for {
			select {
			case <-h.ctx.Done():
				return
			case beat := <-ch:
				func() {
					ctx, cancel := context.WithTimeout(h.ctx, 950*time.Millisecond)
					defer cancel()

					loadAvg, err := loadavg.Avg(ctx)
					if err != nil {
						h.log.Debug(fmt.Errorf("cannot get load average: %w", err))
						return
					}

					h.pointsMutex.Lock()
					defer h.pointsMutex.Unlock()

					beat.Point.LoadAvg = loadAvg
				}()
			}
		}
	}()
}
