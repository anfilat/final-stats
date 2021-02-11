package collector

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/anfilat/final-stats/internal/cpu"
	"github.com/anfilat/final-stats/internal/loadavg"
	"github.com/anfilat/final-stats/internal/loaddisks"
	"github.com/anfilat/final-stats/internal/mocks"
	"github.com/anfilat/final-stats/internal/symo"
	"github.com/anfilat/final-stats/internal/usedfs"
)

func TestCollectorStartStop(t *testing.T) {
	defer goleak.VerifyNone(t)

	config, err := symo.NewConfig("")
	require.NoError(t, err)

	log := new(mocks.Logger)
	log.On("Debug", "collector is stopped")

	collectors := symo.MetricCollectors{
		LoadAvg:   loadavg.Collect,
		CPU:       cpu.Collect,
		LoadDisks: loaddisks.Collect,
		UsedFS:    usedfs.Collect,
	}

	toCollectorCh := make(symo.ClientsToCollectorCh, 1)
	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toCollectorCh, toClientsCh)

	stopCtx := context.Background()
	collectorService.Stop(stopCtx)

	log.AssertExpectations(t)
	require.Len(t, toClientsCh, 0)
}

func TestCollectorStartWithCanceledContext(t *testing.T) {
	defer goleak.VerifyNone(t)

	config, err := symo.NewConfig("")
	require.NoError(t, err)

	log := new(mocks.Logger)
	log.On("Debug", mock.Anything)

	collectors := symo.MetricCollectors{
		LoadAvg:   loadavg.Collect,
		CPU:       cpu.Collect,
		LoadDisks: loaddisks.Collect,
		UsedFS:    usedfs.Collect,
	}

	toCollectorCh := make(symo.ClientsToCollectorCh, 1)
	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx, cancel := context.WithCancel(context.Background())
	cancel()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toCollectorCh, toClientsCh)

	stopCtx := context.Background()
	collectorService.Stop(stopCtx)

	log.AssertExpectations(t)
	require.Len(t, toClientsCh, 0)
}

func TestCollectorTick(t *testing.T) {
	defer goleak.VerifyNone(t)

	config, err := symo.NewConfig("")
	require.NoError(t, err)

	log := new(mocks.Logger)
	log.On("Debug", mock.Anything)
	log.On("Debug", "tick ", mock.Anything)

	collectors := symo.MetricCollectors{
		LoadAvg:   loadavg.Collect,
		CPU:       cpu.Collect,
		LoadDisks: loaddisks.Collect,
		UsedFS:    usedfs.Collect,
	}

	toCollectorCh := make(symo.ClientsToCollectorCh, 1)
	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toCollectorCh, toClientsCh)

	// тик с клиентами
	toCollectorCh <- symo.Start
	time.Sleep(1050 * time.Millisecond)
	require.Len(t, toClientsCh, 1)
	data := <-toClientsCh
	// текущая секунда еще не заполнена, предыдущих нет - статистика должна быть пустой
	require.Len(t, data.Points, 0)

	time.Sleep(1050 * time.Millisecond)
	require.Len(t, toClientsCh, 1)
	data = <-toClientsCh
	require.Len(t, data.Points, 1)

	// тик без клиентов
	toCollectorCh <- symo.Stop
	time.Sleep(1050 * time.Millisecond)
	require.Len(t, toClientsCh, 0)

	stopCtx := context.Background()
	collectorService.Stop(stopCtx)

	log.AssertExpectations(t)
}
