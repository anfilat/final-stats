package collector

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/anfilat/final-stats/internal/mocks"
	"github.com/anfilat/final-stats/internal/symo"
)

func TestCPU(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan timePoint, 1)

	log := new(mocks.Logger)

	cpuData := symo.CPUData{
		User:   0.1,
		System: 0.2,
		Idle:   0.3,
	}

	collector := func(_ context.Context, _ symo.MetricCommand) (*symo.CPUData, error) {
		return &cpuData, nil
	}

	point := &symo.Point{}
	go func() {
		ch <- timePoint{
			time:  time.Now().Truncate(time.Second),
			point: point,
		}
		close(ch)
	}()

	cpuCollect(ctx, ch, collector, log)

	log.AssertExpectations(t)
	require.Equal(t, &cpuData, point.CPU)
}

func TestCPUError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan timePoint, 1)

	log := new(mocks.Logger)
	log.On("Debug", mock.Anything)

	collector := func(_ context.Context, _ symo.MetricCommand) (*symo.CPUData, error) {
		return nil, fmt.Errorf("cannot read the stat file")
	}

	point := &symo.Point{}
	go func() {
		ch <- timePoint{
			time:  time.Now().Truncate(time.Second),
			point: point,
		}
		close(ch)
	}()

	cpuCollect(ctx, ch, collector, log)

	log.AssertExpectations(t)
	require.Nil(t, point.CPU)
}

func TestCPUCloseByContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan timePoint, 1)

	go func() {
		cancel()
	}()

	cpuCollect(ctx, ch, nil, nil)
}
