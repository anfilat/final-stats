// +build linux

package cpu

import (
	"context"

	"github.com/anfilat/final-stats/internal/symo"
)

type cpuData struct {
	total  float64
	user   float64
	system float64
	idle   float64
}

var prevData cpuData

func Read(_ context.Context, init bool) (*symo.CPUData, error) {
	data, err := getCPU()
	if err != nil {
		return nil, err
	}

	if init {
		prevData = *data
		return nil, nil
	}

	total := data.total - prevData.total
	result := &symo.CPUData{
		User:   (data.user - prevData.user) / total * 100,
		System: (data.system - prevData.system) / total * 100,
		Idle:   (data.idle - prevData.idle) / total * 100,
	}
	prevData = *data
	return result, nil
}
