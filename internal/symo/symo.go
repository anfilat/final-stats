package symo

import (
	"context"
	"sync"
	"time"
)

// данные хранятся не более 10 минут.
const MaxSeconds = 10 * 60
const MaxOldPoints = MaxSeconds * time.Second

type Heart interface {
	Start(wg *sync.WaitGroup)
}

type Beat struct {
	Time  time.Time
	Point *Point
}

type MetricReaders struct {
	LoadAvg LoadAvg
}

type Points map[time.Time]*Point

type Point struct {
	LoadAvg *LoadAvgData
}

type LoadAvg func(ctx context.Context) (*LoadAvgData, error)

type LoadAvgData struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

type Clients interface {
	Start(wg *sync.WaitGroup)
	NewClient(client NewClient) <-chan Stat
}

type ClientsToHeartChan chan HeartCommand

type HeartToClientsChan chan ClientsBeat

type HeartCommand int

const (
	Start HeartCommand = iota
	Stop
)

type ClientsBeat struct {
	Time   time.Time
	Points Points
}

type GRPCServer interface {
	Start(addr string) error
	Stop(ctx context.Context) error
}

type NewClient struct {
	Ctx context.Context
	N   int
	M   int
}

type Stat struct {
	Time time.Time
	Stat *Point
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}
