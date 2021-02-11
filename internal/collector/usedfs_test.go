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

func TestUsedFS(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.Logger)

	ufData := symo.UsedFSData{
		{
			Path:      "/",
			UsedSpace: 12.3,
			UsedInode: 7.77,
		},
	}

	collector := func(_ context.Context, _ symo.MetricCommand) (symo.UsedFSData, error) {
		return ufData, nil
	}

	usedFSCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Equal(t, ufData, point.UsedFS)
}

func TestUsedFSError(t *testing.T) {
	ctx, mutex, ch, point := testCollector()

	log := new(mocks.Logger)
	log.On("Debug", mock.Anything)

	collector := func(_ context.Context, _ symo.MetricCommand) (symo.UsedFSData, error) {
		return nil, fmt.Errorf("cannot parse df line")
	}

	usedFSCollect(ctx, mutex, ch, collector, log)

	log.AssertExpectations(t)
	require.Nil(t, point.UsedFS)
}
