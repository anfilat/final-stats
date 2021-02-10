package collector

import (
	"context"
	"sync"
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

const timeToGetMetric = 950 * time.Millisecond

type collector struct {
	ctx           context.Context // управление остановкой сервиса
	ctxCancel     context.CancelFunc
	closedCh      chan interface{}
	mutex         *sync.Mutex
	points        symo.Points        // собираемые данные
	workerChans   []chan<- timePoint // каналы горутин, ответственных за сбор конкретных метрик
	isClients     bool               // для отключения сервиса при отсутствии клиентов
	config        symo.Config
	readers       symo.MetricReaders // функции возвращающие конкретные метрики
	toCollectorCh <-chan symo.CollectorCommand
	toClientsCh   chan<- symo.MetricsData
	log           symo.Logger
}

// информация, отправляемая горутинам, ответственным за сбор конкретных метрик.
type timePoint struct {
	time  time.Time   // за какую секунду метрика
	point *symo.Point // структура, в которую складываются метрики
}

func NewCollector(log symo.Logger, config symo.Config) symo.Collector {
	return &collector{
		config: config,
		log:    log,
	}
}

func (c *collector) Start(_ context.Context, readers symo.MetricReaders,
	toCollectorCh <-chan symo.CollectorCommand, toClientsCh chan<- symo.MetricsData) {
	c.readers = readers
	c.toCollectorCh = toCollectorCh
	c.toClientsCh = toClientsCh

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())
	c.closedCh = make(chan interface{})
	c.mutex = &sync.Mutex{}
	c.points = make(symo.Points)
	c.workerChans = nil
	c.isClients = false

	c.mountMetrics()

	go c.work()
}

func (c *collector) Stop(ctx context.Context) {
	c.ctxCancel()

	select {
	case <-ctx.Done():
	case <-c.closedCh:
	}

	close(c.toClientsCh)
	c.stopMetrics()
	c.unmountMetrics()
	c.log.Debug("collector is stopped")
}

func (c *collector) mountMetrics() {
	if c.config.Metric.Loadavg {
		go loadavg(c.ctx, c.newWorkerChan(), c.readers.LoadAvg, c.log)
	}
	if c.config.Metric.CPU {
		go startCPU(c.ctx, c.readers.CPU, c.log)
		go cpu(c.ctx, c.newWorkerChan(), c.readers.CPU, c.log)
	}
	if c.config.Metric.Loaddisks {
		go startLoadDisks(c.ctx, c.readers.LoadDisks, c.log)
		go loadDisks(c.ctx, c.newWorkerChan(), c.readers.LoadDisks, c.log)
	}
	if c.config.Metric.UsedFS {
		go startUsedFS(c.ctx, c.readers.UsedFS, c.log)
		go usedFS(c.ctx, c.newWorkerChan(), c.readers.UsedFS, c.log)
	}
}

func (c *collector) stopMetrics() {
	ctx := context.Background()
	if c.config.Metric.Loaddisks {
		go stopLoadDisks(ctx, c.readers.LoadDisks, c.log)
	}
	if c.config.Metric.UsedFS {
		go stopUsedFS(ctx, c.readers.UsedFS, c.log)
	}
}

func (c *collector) newWorkerChan() chan timePoint {
	ch := make(chan timePoint, 1)
	c.workerChans = append(c.workerChans, ch)

	return ch
}

func (c *collector) unmountMetrics() {
	for _, ch := range c.workerChans {
		close(ch)
	}
}

func (c *collector) work() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			close(c.closedCh)
			return
		default:
		}

		select {
		case <-c.ctx.Done():
			close(c.closedCh)
			return
		case mess, ok := <-c.toCollectorCh:
			if !ok {
				return
			}
			if mess == symo.Start {
				c.isClients = true
			} else if mess == symo.Stop {
				c.isClients = false
			}
		case now := <-ticker.C:
			c.processTick(now.Truncate(time.Second))
		}
	}
}

func (c *collector) processTick(now time.Time) {
	c.log.Debug("tick ", now)

	// добавляется новая точка для статистики за эту секунду
	point := c.addPoint(now)

	// точка отправляется всем горутинам, ответственным за получение части статистики для заполнения
	for _, ch := range c.workerChans {
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
	c.cleanPoints(now)

	if !c.isClients {
		return
	}
	// накопленная статистика отправляются клиентам

	data := symo.MetricsData{
		Time:   now,
		Points: c.cloneOldPoints(now),
	}

	select {
	case c.toClientsCh <- data:
	default:
	}
}

func (c *collector) addPoint(now time.Time) *symo.Point {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	point := &symo.Point{}
	c.points[now] = point

	return point
}

func (c *collector) cleanPoints(now time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	limit := now.Add(-time.Duration(c.config.App.MaxSeconds) * time.Second)

	for key := range c.points {
		if key.Before(limit) {
			delete(c.points, key)
		}
	}
}

func (c *collector) cloneOldPoints(now time.Time) symo.Points {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	result := make(symo.Points, len(c.points))
	for key, point := range c.points {
		newPoint := *point
		result[key] = &newPoint
	}
	// последняя, еще пустая точка удаляется
	delete(result, now)

	return result
}
