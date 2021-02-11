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

func TestLoadDisks(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.Logger)

	ldData := symo.LoadDisksData{
		{
			Name:    "sda",
			Tps:     5,
			KBRead:  7,
			KBWrite: 12,
		},
	}

	collector := func(_ context.Context, _ symo.MetricCommand) (symo.LoadDisksData, error) {
		return ldData, nil
	}

	loadDisksCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Equal(t, ldData, point.LoadDisks)
}

func TestLoadDisksError(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.Logger)
	log.On("Debug", mock.Anything)

	collector := func(_ context.Context, _ symo.MetricCommand) (symo.LoadDisksData, error) {
		return nil, fmt.Errorf("cannot parse iostat line")
	}

	loadDisksCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Nil(t, point.LoadDisks)
}
