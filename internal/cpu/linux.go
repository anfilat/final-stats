// +build linux

package cpu

import (
	"context"

	"github.com/anfilat/final-stats/internal/symo"
)

func Read(ctx context.Context, init bool) (*symo.CPUData, error) {
	return &symo.CPUData{
		User:   0,
		System: 0,
		Idle:   0,
	}, nil
}
