package symo

import (
	"context"
	"errors"
	"time"

	"github.com/anfilat/final-stats/internal/pb"
)

// данные хранятся не более 10 минут.
const MaxSeconds = 10 * 60
const MaxOldPoints = MaxSeconds * time.Second

// сервис, запускающий каждую секунду сбор статистики и ее отправку клиентам.
type Heart interface {
	Start(ctx context.Context, config MetricConf, readers MetricReaders,
		toHeartChan <-chan HeartCommand, toClientsChan chan<- ClientsBeat)
	Stop(ctx context.Context)
}

// сервис, хранящий всех подключенных клиентов и отсылающий им статистику.
type Clients interface {
	Start(ctx context.Context, toHeartChan chan<- HeartCommand, toClientsChan <-chan ClientsBeat)
	Stop(ctx context.Context)
	// канал для получения отсылаемых данных и ф-ия отключения клиента
	NewClient(client NewClient) (<-chan *pb.Stats, func(), error)
}

// ошибка, возвращаемая grpc запросу, если приложение останавливается
var ErrStopped = errors.New("service is stopped")

// канал для управления сервисом Heart из сервиса Clients. Если клиентов нет, статистику собирать не нужно.
type ClientsToHeartChan chan HeartCommand

type HeartCommand int

const (
	Start HeartCommand = iota
	Stop
)

// канал для посекундной передачи накопленных данных сервису клиентов.
type HeartToClientsChan chan ClientsBeat

// сервису клиентов отсылается текущая секунда и копия всех собранных данных.
// Отсылается копия, чтобы сервис мог ее обрабатывать, не блокируя мьютекст с собираемыми данными.
type ClientsBeat struct {
	Time   time.Time
	Points Points
}

// информация, отправляемая горутинам, ответственным за сбор какой-то метрики.
type Beat struct {
	Time  time.Time // за какую секунду метрика
	Point *Point    // структура, в которую складываются метрики
}

// хранилище собранных посекундных наборов метрик.
type Points map[time.Time]*Point

// набор метрик (снапшот). За секунду или усредненный.
type Point struct {
	LoadAvg *LoadAvgData
	CPU     *CPUData
}

// набор функций, возвращающих свои метрики. Передается сервису Heart при его создании.
type MetricReaders struct {
	LoadAvg LoadAvg
	CPU     CPU
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
type CPU func(ctx context.Context, init bool) (*CPUData, error)

// средняя загрузка cpu. В процентах.
type CPUData struct {
	User   float64
	System float64
	Idle   float64
}

type GRPCServer interface {
	Start(addr string, clients Clients) error
	Stop(ctx context.Context)
}

// информация, передаваемая из grpc запроса сервису клиентов.
type NewClient struct {
	Ctx context.Context // контекст запроса, закрывается при его окончании. Нужен для определения отключения клиента
	N   int             // информация отправляется каждые N секунд
	M   int             // информация усредняется за M секунд
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}
