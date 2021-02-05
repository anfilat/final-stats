package engine

import (
	"context"
	"sync"
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

const timeToGetMetric = 950 * time.Millisecond

type engine struct {
	ctx           context.Context // управление остановкой сервиса
	ctxCancel     context.CancelFunc
	closedCh      chan interface{}
	mutex         *sync.Mutex
	points        symo.Points        // собираемые данные
	workerChans   []chan<- timePoint // каналы горутин, ответственных за сбор конкретных метрик
	isClients     bool               // сервис работает, только если есть клиенты
	config        symo.MetricConf    // какие метрики собираются
	readers       symo.MetricReaders // функции возвращающие конкретные метрики
	toEngineChan  <-chan symo.EngineCommand
	toClientsChan chan<- symo.MetricsData
	log           symo.Logger
}

// информация, отправляемая горутинам, ответственным за сбор конкретных метрик.
type timePoint struct {
	time  time.Time   // за какую секунду метрика
	point *symo.Point // структура, в которую складываются метрики
}

func NewEngine(log symo.Logger) symo.Engine {
	return &engine{
		log: log,
	}
}

func (h *engine) Start(ctx context.Context, config symo.MetricConf, readers symo.MetricReaders,
	toEngineChan <-chan symo.EngineCommand, toClientsChan chan<- symo.MetricsData) {
	h.config = config
	h.readers = readers
	h.toEngineChan = toEngineChan
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

func (h *engine) Stop(ctx context.Context) {
	h.ctxCancel()

	select {
	case <-ctx.Done():
	case <-h.closedCh:
	}

	close(h.toClientsChan)
	h.unmountMetrics()
	h.log.Debug("engine is stopped")
}

func (h *engine) mountMetrics() {
	if h.config.Loadavg {
		go loadavg(h.ctx, h.newWorkerChan(), h.readers.LoadAvg, h.log)
	}
	if h.config.CPU {
		go cpu(h.ctx, h.newWorkerChan(), h.readers.CPU, h.log)
	}
}

func (h *engine) initMetrics() {
	if h.config.CPU {
		initCPU(h.ctx, h.readers.CPU, h.log)
	}
}

func (h *engine) newWorkerChan() chan timePoint {
	ch := make(chan timePoint, 1)
	h.workerChans = append(h.workerChans, ch)

	return ch
}

func (h *engine) unmountMetrics() {
	for _, ch := range h.workerChans {
		close(ch)
	}
}

func (h *engine) work() {
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

func (h *engine) waitClients() {
	select {
	case <-h.ctx.Done():
	case mess, ok := <-h.toEngineChan:
		if !ok {
			return
		}
		if mess == symo.Start {
			h.isClients = true
			h.log.Debug("start engine work")
		}
	}
}

func (h *engine) process() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	h.initMetrics()

	for {
		select {
		case <-h.ctx.Done():
			return
		case mess, ok := <-h.toEngineChan:
			if !ok {
				return
			}
			if mess == symo.Stop {
				h.isClients = false
				h.log.Debug("pause engine work")
				return
			}
		case now := <-ticker.C:
			h.processTick(now.Truncate(time.Second))
		}
	}
}

func (h *engine) processTick(now time.Time) {
	h.log.Debug("tick ", now)

	// добавляется новая точка для статистики за эту секунду
	point := h.addPoint(now)

	// точка отправляется всем горутинам, ответственным за получение части статистики для заполнения
	for _, ch := range h.workerChans {
		tp := timePoint{
			time:  now,
			point: point,
		}
		select {
		case ch <- tp:
		default:
		}
	}

	// устаревшие точки удаляются
	h.cleanPoints(now)

	// накопленная статистика отправляются клиентам

	data := symo.MetricsData{
		Time:   now,
		Points: h.cloneOldPoints(now),
	}

	select {
	case h.toClientsChan <- data:
	default:
	}
}

func (h *engine) addPoint(now time.Time) *symo.Point {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	point := &symo.Point{}
	h.points[now] = point

	return point
}

func (h *engine) cleanPoints(now time.Time) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	limit := now.Add(-symo.MaxOldPoints)

	for key := range h.points {
		if key.Before(limit) {
			delete(h.points, key)
		}
	}
}

func (h *engine) cloneOldPoints(now time.Time) symo.Points {
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
