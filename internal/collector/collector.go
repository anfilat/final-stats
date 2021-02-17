package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

const timeToGetMetric = 950 * time.Millisecond

type collector struct {
	ctx         context.Context // управление остановкой сервиса
	ctxCancel   context.CancelFunc
	stoppedCh   chan interface{}
	mutex       *sync.Mutex
	points      symo.Points        // собираемые данные
	workerChans []chan<- timePoint // каналы горутин, ответственных за сбор конкретных метрик
	config      symo.Config
	collectors  symo.MetricCollectors // функции возвращающие конкретные метрики
	toClientsCh chan<- symo.MetricsData
	log         symo.Logger
}

// информация, отправляемая горутинам, ответственным за сбор конкретных метрик.
type timePoint struct {
	time  time.Time   // за какую секунду метрика
	point *symo.Point // структура, в которую складываются метрики
}

// NewCollector возвращает сервис сбора метрик.
func NewCollector(log symo.Logger, config symo.Config) symo.Collector {
	return &collector{
		config: config,
		log:    log,
	}
}

func (c *collector) Start(ctx context.Context, collectors symo.MetricCollectors, toClientsCh chan<- symo.MetricsData) {
	c.collectors = collectors
	c.toClientsCh = toClientsCh

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())
	c.stoppedCh = make(chan interface{})
	c.mutex = &sync.Mutex{}
	c.points = make(symo.Points)
	c.workerChans = nil

	mountedCh := make(chan interface{})
	go c.mountMetrics(ctx, mountedCh)

	select {
	case <-ctx.Done():
		// стартануть не успели, помечаем сервис как остановленный
		close(c.stoppedCh)
		return
	case <-mountedCh:
		go c.work()
	}
}

func (c *collector) Stop(ctx context.Context) {
	c.ctxCancel()

	select {
	case <-ctx.Done():
		return
	case <-c.stoppedCh:
	}

	unmountedCh := make(chan interface{})
	go c.unmountMetrics(ctx, unmountedCh)

	select {
	case <-ctx.Done():
		return
	case <-unmountedCh:
		c.log.Debug("collector is stopped")
	}
}

func (c *collector) mountMetrics(startCtx context.Context, mountedCh chan interface{}) {
	wg := &sync.WaitGroup{}

	if c.config.Metric.Loadavg {
		go loadavgCollect(c.ctx, c.mutex, c.newWorkerChan(), c.collectors.LoadAvg, c.log)
	}
	if c.config.Metric.CPU {
		wg.Add(1)
		go c.mountCPU(startCtx, wg)
	}
	if c.config.Metric.Loaddisks {
		wg.Add(1)
		go c.mountLoadDisks(startCtx, wg)
	}
	if c.config.Metric.UsedFS {
		wg.Add(1)
		go c.mountUsedFS(startCtx, wg)
	}

	wg.Wait()
	close(mountedCh)
}

func (c *collector) mountCPU(startCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := c.collectors.CPU(startCtx, symo.StartMetric)
	if err != nil {
		c.log.Debug(fmt.Errorf("cannot start cpu: %w", err))
		return
	}
	go cpuCollect(c.ctx, c.mutex, c.newWorkerChan(), c.collectors.CPU, c.log)
}

func (c *collector) mountLoadDisks(startCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := c.collectors.LoadDisks(startCtx, symo.StartMetric)
	if err != nil {
		c.log.Debug(fmt.Errorf("cannot start load disks: %w", err))
	}
	go loadDisksCollect(c.ctx, c.mutex, c.newWorkerChan(), c.collectors.LoadDisks, c.log)
}

func (c *collector) mountUsedFS(startCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := c.collectors.UsedFS(startCtx, symo.StartMetric)
	if err != nil {
		c.log.Debug(fmt.Errorf("cannot start used fs: %w", err))
	}
	go usedFSCollect(c.ctx, c.mutex, c.newWorkerChan(), c.collectors.UsedFS, c.log)
}

func (c *collector) unmountMetrics(stopCtx context.Context, unmountedCh chan interface{}) {
	wg := &sync.WaitGroup{}

	if c.config.Metric.Loaddisks {
		wg.Add(1)
		go c.unmountLoadDisks(stopCtx, wg)
	}
	if c.config.Metric.UsedFS {
		wg.Add(1)
		go c.unmountUsedFS(stopCtx, wg)
	}

	wg.Wait()
	close(unmountedCh)
}

func (c *collector) unmountLoadDisks(stopCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := c.collectors.LoadDisks(stopCtx, symo.StopMetric)
	if err != nil {
		c.log.Debug(fmt.Errorf("cannot stop load disks: %w", err))
	}
}

func (c *collector) unmountUsedFS(stopCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := c.collectors.UsedFS(stopCtx, symo.StopMetric)
	if err != nil {
		c.log.Debug(fmt.Errorf("cannot stop used fs: %w", err))
	}
}

func (c *collector) newWorkerChan() chan timePoint {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ch := make(chan timePoint, 1)
	c.workerChans = append(c.workerChans, ch)

	return ch
}

func (c *collector) work() {
	defer close(c.stoppedCh)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		select {
		case <-c.ctx.Done():
			return
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
		if key.Equal(now) {
			continue
		}
		newPoint := *point
		result[key] = &newPoint
	}

	return result
}
