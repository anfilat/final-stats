package collector

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/anfilat/final-stats/internal/mocks"
	"github.com/anfilat/final-stats/internal/symo"
)

func TestLoadAvg(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.Logger)

	loadAvg := symo.LoadAvgData{
		Load1:  0.1,
		Load5:  0.2,
		Load15: 0.3,
	}

	collector := func(_ context.Context) (*symo.LoadAvgData, error) {
		return &loadAvg, nil
	}

	loadavgCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Equal(t, &loadAvg, point.LoadAvg)
}

func TestLoadAvgError(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.Logger)
	log.On("Debug", mock.Anything)

	collector := func(_ context.Context) (*symo.LoadAvgData, error) {
		return nil, fmt.Errorf("cannot read the loadavg file")
	}

	loadavgCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Nil(t, point.LoadAvg)
}
