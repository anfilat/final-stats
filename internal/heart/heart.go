package heart

import (
	"context"
	"sync"
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

const timeToGetMetric = 950 * time.Millisecond

type heart struct {
	ctx           context.Context // управление остановкой сервиса
	ctxCancel     context.CancelFunc
	closedCh      chan interface{}
	mutex         *sync.Mutex
	points        symo.Points        // собираемые данные
	workerChans   []chan<- symo.Beat // каналы горутин, ответственных за сбор конкретных метрик
	isClients     bool               // сервис работает, только если есть клиенты
	config        symo.MetricConf    // какие метрики собираются
	readers       symo.MetricReaders // функции возвращающие конкретные метрики
	toHeartChan   <-chan symo.HeartCommand
	toClientsChan chan<- symo.ClientsBeat
	log           symo.Logger
}

func NewHeart(log symo.Logger) symo.Heart {
	return &heart{
		log: log,
	}
}

func (h *heart) Start(ctx context.Context, config symo.MetricConf, readers symo.MetricReaders,
	toHeartChan <-chan symo.HeartCommand, toClientsChan chan<- symo.ClientsBeat) {
	h.config = config
	h.readers = readers
	h.toHeartChan = toHeartChan
	h.toClientsChan = toClientsChan

	h.ctx, h.ctxCancel = context.WithCancel(ctx)
	h.closedCh = make(chan interface{})
	h.mutex = &sync.Mutex{}
	h.points = make(symo.Points)
	h.workerChans = nil
	h.isClients = false

	h.mountMetrics()

	go h.work()
}

func (h *heart) Stop(ctx context.Context) {
	h.ctxCancel()

	select {
	case <-ctx.Done():
	case <-h.closedCh:
	}

	close(h.toClientsChan)
	h.unmountMetrics()
	h.log.Debug("heart is stopped")
}

func (h *heart) mountMetrics() {
	if h.config.Loadavg {
		go loadavg(h.ctx, h.newWorkerChan(), h.readers.LoadAvg, h.log)
	}
	if h.config.CPU {
		go cpu(h.ctx, h.newWorkerChan(), h.readers.CPU, h.log)
	}
}

func (h *heart) initMetrics() {
	if h.config.CPU {
		initCPU(h.ctx, h.readers.CPU, h.log)
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

func (h *heart) work() {
	for {
		select {
		case <-h.ctx.Done():
			close(h.closedCh)
			return
		default:
		}

		if h.isClients {
			h.process()
		} else {
			h.waitClients()
		}
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

func (h *heart) process() {
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
	h.mutex.Lock()
	defer h.mutex.Unlock()

	point := &symo.Point{}
	h.points[now] = point

	return point
}

func (h *heart) cleanPoints(now time.Time) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	limit := now.Add(-symo.MaxOldPoints)

	for key := range h.points {
		if key.Before(limit) {
			delete(h.points, key)
		}
	}
}

func (h *heart) cloneOldPoints(now time.Time) symo.Points {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	result := make(symo.Points, len(h.points))
	for key, point := range h.points {
		newPoint := *point
		result[key] = &newPoint
	}
	// последняя, еще пустая точка удаляется
	delete(result, now)

	return result
}
