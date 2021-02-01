package symo

import (
	"sync"
	"time"
)

const MaxOldPoints = 10 * time.Minute

type Heart interface {
	Start(wg *sync.WaitGroup, config MetricConf)
}

type Beat struct {
	Time  time.Time
	Point *Point
}

type Points map[time.Time]*Point

type Point struct {
	LoadAvg *LoadAvgData
}

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
