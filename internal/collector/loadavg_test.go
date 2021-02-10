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

func TestLoadAvg(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan timePoint, 1)

	log := new(mocks.Logger)

	loadAvg := symo.LoadAvgData{
		Load1:  0.1,
		Load5:  0.2,
		Load15: 0.3,
	}

	reader := func(_ context.Context) (*symo.LoadAvgData, error) {
		return &loadAvg, nil
	}

	point := &symo.Point{}
	go func() {
		ch <- timePoint{
			time:  time.Now().Truncate(time.Second),
			point: point,
		}
		close(ch)
	}()

	loadavg(ctx, ch, reader, log)

	log.AssertExpectations(t)
	require.Equal(t, &loadAvg, point.LoadAvg)
}

func TestLoadAvgError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan timePoint, 1)

	log := new(mocks.Logger)
	log.On("Debug", mock.Anything)

	reader := func(_ context.Context) (*symo.LoadAvgData, error) {
		return nil, fmt.Errorf("cannot read the loadavg file")
	}

	point := &symo.Point{}
	go func() {
		ch <- timePoint{
			time:  time.Now().Truncate(time.Second),
			point: point,
		}
		close(ch)
	}()

	loadavg(ctx, ch, reader, log)

	log.AssertExpectations(t)
	require.Nil(t, point.LoadAvg)
}

func TestLoadAvgCloseByContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan timePoint, 1)

	go func() {
		cancel()
	}()

	loadavg(ctx, ch, nil, nil)
}
