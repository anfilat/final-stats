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
	Start(wg *sync.WaitGroup, config MetricConf, readers MetricReaders)
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

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}
