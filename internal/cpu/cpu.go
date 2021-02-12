package cpu

import (
	"context"
	"sync"

	"github.com/anfilat/final-stats/internal/symo"
)

type cpuData struct {
	total  float64
	user   float64
	system float64
	idle   float64
}

var (
	mutex    sync.Mutex
	prevData *cpuData
)

func Collect(_ context.Context, action symo.MetricCommand) (*symo.CPUData, error) {
	switch action {
	case symo.StartMetric:
		return nil, start()
	case symo.StopMetric:
		return nil, nil
	default:
		return get()
	}
}

func start() error {
	data, err := getCPU()
	if err != nil {
		return err
	}

	mutex.Lock()
	defer mutex.Unlock()

	prevData = data

	return nil
}

func get() (*symo.CPUData, error) {
	data, err := getCPU()
	if err != nil {
		return nil, err
	}

	mutex.Lock()
	defer mutex.Unlock()

	if prevData == nil {
		prevData = data
		return nil, nil
	}

	total := data.total - prevData.total
	if total == 0 {
		return &symo.CPUData{}, nil
	}

	result := &symo.CPUData{
		User:   (data.user - prevData.user) / total * 100,
		System: (data.system - prevData.system) / total * 100,
		Idle:   (data.idle - prevData.idle) / total * 100,
	}
	prevData = data
	return result, nil
}
