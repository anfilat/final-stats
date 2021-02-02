package heart

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

type heart struct {
	ctx           context.Context
	pointsMutex   *sync.Mutex
	points        symo.Points
	workerChains  []chan<- symo.Beat
	isClients     bool
	config        symo.MetricConf
	readers       symo.MetricReaders
	toHeartChan   <-chan symo.HeartCommand
	toClientsChan chan<- symo.ClientsBeat
	log           symo.Logger
}

func NewHeart(ctx context.Context, log symo.Logger, config symo.MetricConf, readers symo.MetricReaders,
	toHeartChan <-chan symo.HeartCommand, toClientsChan chan<- symo.ClientsBeat) symo.Heart {
	return &heart{
		ctx:           ctx,
		pointsMutex:   &sync.Mutex{},
		points:        make(symo.Points),
		isClients:     false,
		config:        config,
		readers:       readers,
		toHeartChan:   toHeartChan,
		toClientsChan: toClientsChan,
		log:           log,
	}
}

func (h *heart) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	h.mountMetrics(h.config, h.readers)

	go func() {
		defer func() {
			close(h.toClientsChan)
			h.unmountMetrics()
			wg.Done()
		}()

		for {
			if h.ctx.Err() != nil {
				return
			}
			if h.isClients {
				h.work()
			} else {
				h.waitClients()
			}
		}
	}()
}

func (h *heart) mountMetrics(config symo.MetricConf, readers symo.MetricReaders) {
	if config.Loadavg {
		h.loadavg(h.newWorkerChan(), readers.LoadAvg)
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

func (h *heart) waitClients() {
	select {
	case <-h.ctx.Done():
	case mess := <-h.toHeartChan:
		if mess == symo.Start {
			h.isClients = true
		}
	}
}

func (h *heart) work() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case mess := <-h.toHeartChan:
			if mess == symo.Stop {
				h.isClients = false
				return
			}
		case now := <-ticker.C:
			h.workPoint(now.Truncate(time.Second))
		}
	}
}

func (h *heart) workPoint(now time.Time) {
	// добавляется новая точка для статистики за эту секунду
	point := h.newPoint(now)

	// точка отправляется всем горутинам, ответственным за получение части статистики для заполнения
	for _, ch := range h.workerChains {
		ch <- symo.Beat{
			Time:  now,
			Point: point,
		}
	}

	// устаревшие точки удаляются
	h.cleanPoints(now)

	// накопленная статистика отправляются клиентам
	h.toClientsChan <- symo.ClientsBeat{
		Time:   now,
		Points: h.CloneOldPoints(now),
	}
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

func (h *heart) CloneOldPoints(now time.Time) symo.Points {
	h.pointsMutex.Lock()
	defer h.pointsMutex.Unlock()

	result := make(symo.Points, len(h.points))
	for key, point := range h.points {
		newPoint := *point
		result[key] = &newPoint
	}
	delete(result, now)

	return result
}

func (h *heart) loadavg(ch <-chan symo.Beat, reader symo.LoadAvg) {
	go func() {
		for {
			select {
			case <-h.ctx.Done():
				return
			case beat, ok := <-ch:
				if !ok {
					return
				}
				func() {
					ctx, cancel := context.WithTimeout(h.ctx, 950*time.Millisecond)
					defer cancel()

					loadAvg, err := reader(ctx)
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
