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

func TestUsedFS(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan timePoint, 1)

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

	point := &symo.Point{}
	go func() {
		ch <- timePoint{
			time:  time.Now().Truncate(time.Second),
			point: point,
		}
		close(ch)
	}()

	usedFS(ctx, ch, collector, log)

	log.AssertExpectations(t)
	require.Equal(t, ufData, point.UsedFS)
}

func TestUsedFSError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan timePoint, 1)

	log := new(mocks.Logger)
	log.On("Debug", mock.Anything)

	collector := func(_ context.Context, _ symo.MetricCommand) (symo.UsedFSData, error) {
		return nil, fmt.Errorf("cannot parse df line")
	}

	point := &symo.Point{}
	go func() {
		ch <- timePoint{
			time:  time.Now().Truncate(time.Second),
			point: point,
		}
		close(ch)
	}()

	usedFS(ctx, ch, collector, log)

	log.AssertExpectations(t)
	require.Nil(t, point.UsedFS)
}

func TestUsedFSCloseByContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan timePoint, 1)

	go func() {
		cancel()
	}()

	usedFS(ctx, ch, nil, nil)
}
