package symo

import (
	"context"
	"errors"
	"time"
)

// сервис, запускающий каждую секунду сбор статистики и ее отправку клиентам.
type Collector interface {
	Start(context.Context, MetricCollectors, chan<- MetricsData)
	Stop(context.Context)
}

type NewClienter interface {
	// канал для получения отсылаемых данных и ф-ия отключения клиента
	NewClient(ClientData) (<-chan *Stats, func(), error)
}

// сервис, хранящий всех подключенных клиентов и отсылающий им статистику.
type Clients interface {
	Start(context.Context, <-chan MetricsData)
	Stop(context.Context)
	NewClienter
}

// ErrStopped ошибка, возвращаемая grpc запросу, если приложение останавливается.
var ErrStopped = errors.New("service is stopped")

// канал для посекундной передачи накопленных данных сервису клиентов.
type CollectorToClientsCh chan MetricsData

// сервису клиентов отсылается текущая секунда и копия всех собранных данных.
// Отсылается копия, чтобы сервис мог ее обрабатывать, не блокируя мьютекс с собираемыми данными.
type MetricsData struct {
	Time   time.Time
	Points Points
}

// данные для клиента.
type Stats struct {
	Time      time.Time
	LoadAvg   *LoadAvgData
	CPU       *CPUData
	LoadDisks LoadDisksData
	UsedFS    UsedFSData
}

// хранилище собранных посекундных наборов метрик.
type Points map[time.Time]*Point

// набор метрик (снапшот). За секунду или усредненный.
type Point struct {
	LoadAvg   *LoadAvgData
	CPU       *CPUData
	LoadDisks LoadDisksData
	UsedFS    UsedFSData
}

type MetricCommand int

const (
	StartMetric MetricCommand = iota
	StopMetric
	GetMetric
)

// набор функций, возвращающих свои метрики. Передается сервису Collector при его создании.
type MetricCollectors struct {
	LoadAvg   LoadAvg
	CPU       CPU
	LoadDisks LoadDisks
	UsedFS    UsedFS
}

// функция возвращающая среднюю загрузку системы.
type LoadAvg func(ctx context.Context) (*LoadAvgData, error)

// средняя загрузка системы.
type LoadAvgData struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

// функция возвращающая среднюю загрузку cpu.
type CPU func(ctx context.Context, action MetricCommand) (*CPUData, error)

// средняя загрузка cpu. В процентах.
type CPUData struct {
	User   float64
	System float64
	Idle   float64
}

// функция возвращающая загрузку дисков.
type LoadDisks func(ctx context.Context, action MetricCommand) (LoadDisksData, error)

type LoadDisksData []DiskData

// загрузка дисков.
type DiskData struct {
	Name    string
	Tps     float64
	KBRead  float64
	KBWrite float64
}

// функция возвращающая использование файловых систем.
type UsedFS func(ctx context.Context, action MetricCommand) (UsedFSData, error)

type UsedFSData []FSData

// использовано в каждой файловой системе.
type FSData struct {
	Path      string
	UsedSpace float64
	UsedInode float64
}

type GRPCServer interface {
	Start(addr string, clients NewClienter) error
	Stop(ctx context.Context)
}

// информация, передаваемая из grpc запроса сервису клиентов.
type ClientData struct {
	N int // информация отправляется каждые N секунд
	M int // информация усредняется за M секунд
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}
