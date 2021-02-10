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

func TestLoadDisks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan timePoint, 1)

	log := new(mocks.Logger)

	ldData := symo.LoadDisksData{
		{
			Name:    "sda",
			Tps:     5,
			KBRead:  7,
			KBWrite: 12,
		},
	}

	reader := func(_ context.Context, _ symo.MetricCommand) (symo.LoadDisksData, error) {
		return ldData, nil
	}

	point := &symo.Point{}
	go func() {
		ch <- timePoint{
			time:  time.Now().Truncate(time.Second),
			point: point,
		}
		close(ch)
	}()

	loadDisks(ctx, ch, reader, log)

	log.AssertExpectations(t)
	require.Equal(t, ldData, point.LoadDisks)
}

func TestLoadDisksError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan timePoint, 1)

	log := new(mocks.Logger)
	log.On("Debug", mock.Anything)

	reader := func(_ context.Context, _ symo.MetricCommand) (symo.LoadDisksData, error) {
		return nil, fmt.Errorf("cannot parse iostat line")
	}

	point := &symo.Point{}
	go func() {
		ch <- timePoint{
			time:  time.Now().Truncate(time.Second),
			point: point,
		}
		close(ch)
	}()

	loadDisks(ctx, ch, reader, log)

	log.AssertExpectations(t)
	require.Nil(t, point.LoadDisks)
}

func TestLoadDisksCloseByContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan timePoint, 1)

	go func() {
		cancel()
	}()

	loadDisks(ctx, ch, nil, nil)
}
