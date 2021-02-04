package heart

import (
	"context"
	"sync"
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

const timeToGetMetric = 950 * time.Millisecond

type heart struct {
	ctx           context.Context // контекст приложения, сервис завершается по закрытию контекста
	pointsMutex   *sync.Mutex
	points        symo.Points        // собираемые данные
	workerChans   []chan<- symo.Beat // каналы горутин, ответственных за сбор конкретных метрик
	isClients     bool               // сервис работает, только если есть клиенты
	config        symo.MetricConf    // какие метрики собираются
	readers       symo.MetricReaders // функции возвращающие конкретные метрики
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
	h.mountMetrics()

	go func() {
		defer func() {
			close(h.toClientsChan)
			h.unmountMetrics()
			h.log.Debug("heart is stopped")
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

func (h *heart) mountMetrics() {
	if h.config.Loadavg {
		go loadavg(h, h.newWorkerChan(), h.readers.LoadAvg)
	}
	if h.config.CPU {
		go cpu(h, h.newWorkerChan(), h.readers.CPU)
	}
}

func (h *heart) initMetrics() {
	if h.config.CPU {
		initCPU(h)
	}
}

func (h *heart) newWorkerChan() chan symo.Beat {
	ch := make(chan symo.Beat, 1)
	h.workerChans = append(h.workerChans, ch)

	return ch
}

func (h *heart) unmountMetrics() {
	for _, ch := range h.workerChans {
		close(ch)
	}
}

func (h *heart) waitClients() {
	select {
	case <-h.ctx.Done():
	case mess, ok := <-h.toHeartChan:
		if !ok {
			return
		}
		if mess == symo.Start {
			h.isClients = true
			h.log.Debug("start heart work")
		}
	}
}

func (h *heart) work() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	h.initMetrics()

	for {
		select {
		case <-h.ctx.Done():
			return
		case mess, ok := <-h.toHeartChan:
			if !ok {
				return
			}
			if mess == symo.Stop {
				h.isClients = false
				h.log.Debug("pause heart work")
				return
			}
		case now := <-ticker.C:
			h.processTick(now.Truncate(time.Second))
		}
	}
}

func (h *heart) processTick(now time.Time) {
	h.log.Debug("tick ", now)

	// добавляется новая точка для статистики за эту секунду
	point := h.addPoint(now)

	// точка отправляется всем горутинам, ответственным за получение части статистики для заполнения
	for _, ch := range h.workerChans {
		beat := symo.Beat{
			Time:  now,
			Point: point,
		}
		select {
		case ch <- beat:
		default:
		}
	}

	// устаревшие точки удаляются
	h.cleanPoints(now)

	// накопленная статистика отправляются клиентам

	data := symo.ClientsBeat{
		Time:   now,
		Points: h.cloneOldPoints(now),
	}

	select {
	case h.toClientsChan <- data:
	default:
	}
}

func (h *heart) addPoint(now time.Time) *symo.Point {
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

func (h *heart) cloneOldPoints(now time.Time) symo.Points {
	h.pointsMutex.Lock()
	defer h.pointsMutex.Unlock()

	result := make(symo.Points, len(h.points))
	for key, point := range h.points {
		newPoint := *point
		result[key] = &newPoint
	}
	// последняя, еще пустая точка удаляется
	delete(result, now)

	return result
}
