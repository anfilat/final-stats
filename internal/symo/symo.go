package symo

import (
	"context"
	"errors"
	"time"
)

// Collector представляет сервис, запускающий каждую секунду сбор статистики и ее отправку клиентам.
type Collector interface {
	Start(context.Context, MetricCollectors, chan<- MetricsData)
	Stop(context.Context)
}

// NewClienter представляет интерфейс для подключения новых клиентов.
type NewClienter interface {
	// возвращает канал для получения отсылаемых данных и ф-ия отключения клиента
	NewClient(ClientData) (<-chan *Stats, func(), error)
}

// Clients представляет сервис, хранящий всех подключенных клиентов и отсылающий им статистику.
type Clients interface {
	Start(context.Context, <-chan MetricsData)
	Stop(context.Context)
	NewClienter
}

// ErrStopped ошибка, возвращаемая grpc запросу, если приложение останавливается.
var ErrStopped = errors.New("service is stopped")

// CollectorToClientsCh - канал для посекундной передачи накопленных данных сервису клиентов.
type CollectorToClientsCh chan MetricsData

// MetricsData содержит данные, отсылаемые сервису клиентов.
// Отсылается текущая секунда и копия всех собранных данных.
// Отсылается копия, чтобы сервис мог ее обрабатывать, не блокируя мьютекс с собираемыми данными.
type MetricsData struct {
	Time   time.Time
	Points Points
}

// Stats содержит данные, отсылаемые каждому клиенту.
type Stats struct {
	Time      time.Time
	LoadAvg   *LoadAvgData
	CPU       *CPUData
	LoadDisks LoadDisksData
	UsedFS    UsedFSData
}

// Points хранит собранные посекундные наборы метрик.
type Points map[time.Time]*Point

// Point содержит набор метрик (снапшот). За секунду или усредненный.
type Point struct {
	LoadAvg   *LoadAvgData
	CPU       *CPUData
	LoadDisks LoadDisksData
	UsedFS    UsedFSData
}

// MetricCommand - команды для взаимодействия сервиса метрик и коллекторами, собирающими метрики.
type MetricCommand int

const (
	// StartMetric - начать собирать метрики.
	StartMetric MetricCommand = iota
	// StopMetric - остановить сбор метрик.
	StopMetric
	// GetMetric - получить метрики.
	GetMetric
)

// MetricCollectors - набор функций, возвращающих свои метрики. Передается сервису Collector при его создании.
type MetricCollectors struct {
	LoadAvg   LoadAvg
	CPU       CPU
	LoadDisks LoadDisks
	UsedFS    UsedFS
}

// LoadAvg - функция возвращающая среднюю загрузку системы.
type LoadAvg func(ctx context.Context) (*LoadAvgData, error)

// LoadAvgData содержит метрики средней загрузки системы.
type LoadAvgData struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

// CPU - функция возвращающая среднюю загрузку cpu.
type CPU func(ctx context.Context, action MetricCommand) (*CPUData, error)

// CPUData содержит метрики средней загрузки cpu. В процентах.
type CPUData struct {
	User   float64
	System float64
	Idle   float64
}

// LoadDisks - функция возвращающая загрузку дисков.
type LoadDisks func(ctx context.Context, action MetricCommand) (LoadDisksData, error)

// LoadDisksData - слайс информации о загрузке дисков.
type LoadDisksData []DiskData

// DiskData содержит метрики загрузки дисков.
type DiskData struct {
	Name    string
	Tps     float64
	KBRead  float64
	KBWrite float64
}

// UsedFS - функция возвращающая использование файловых систем.
type UsedFS func(ctx context.Context, action MetricCommand) (UsedFSData, error)

// UsedFSData - слайс информации об использовании файловых систем.
type UsedFSData []FSData

// FSData содержит информацию об использовании файловой системы.
type FSData struct {
	Path      string
	UsedSpace float64
	UsedInode float64
}

// GRPCServer представляет gRPC сервер.
type GRPCServer interface {
	Start(addr string, clients NewClienter) error
	Stop(ctx context.Context)
}

// ClientData - информация, передаваемая из grpc запроса сервису клиентов.
type ClientData struct {
	N int // информация отправляется каждые N секунд
	M int // информация усредняется за M секунд
}

// Logger представляет логгер.
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}
