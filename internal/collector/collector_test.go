package collector

import (
	"context"
	"errors"
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

	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toClientsCh)

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

	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx, cancel := context.WithCancel(context.Background())
	cancel()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toClientsCh)

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

	laData := &symo.LoadAvgData{
		Load1:  1,
		Load5:  2,
		Load15: 3,
	}
	LoadAvg := new(mocks.LoadAvg)
	LoadAvg.On("Execute", mock.Anything).Return(laData, nil)

	cpuData := &symo.CPUData{
		User:   0.1,
		System: 0.2,
		Idle:   0.3,
	}
	CPU := new(mocks.CPU)
	CPU.On("Execute", mock.Anything, mock.Anything).Return(cpuData, nil)

	ldData := symo.LoadDisksData{
		{
			Name:    "sda",
			Tps:     7,
			KBRead:  8,
			KBWrite: 9,
		},
	}
	LoadDisks := new(mocks.LoadDisks)
	LoadDisks.On("Execute", mock.Anything, mock.Anything).Return(ldData, nil)

	fsData := symo.UsedFSData{
		{
			Path:      "/",
			UsedSpace: 13,
			UsedInode: 31,
		},
	}
	UsedFS := new(mocks.UsedFS)
	UsedFS.On("Execute", mock.Anything, mock.Anything).Return(fsData, nil)

	collectors := symo.MetricCollectors{
		LoadAvg:   LoadAvg.Execute,
		CPU:       CPU.Execute,
		LoadDisks: LoadDisks.Execute,
		UsedFS:    UsedFS.Execute,
	}

	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toClientsCh)

	time.Sleep(50 * time.Millisecond)

	time.Sleep(time.Second)
	data := <-toClientsCh
	// текущая секунда еще не заполнена, предыдущих нет - статистика должна быть пустой
	require.Len(t, data.Points, 0)

	time.Sleep(time.Second)
	data = <-toClientsCh
	require.Len(t, data.Points, 1)
	for _, point := range data.Points {
		require.Equal(t, laData, point.LoadAvg)
		require.Equal(t, cpuData, point.CPU)
		require.Equal(t, ldData, point.LoadDisks)
		require.Equal(t, fsData, point.UsedFS)
	}

	stopCtx := context.Background()
	collectorService.Stop(stopCtx)

	log.AssertExpectations(t)
}

func TestCollectorTickWithErrors(t *testing.T) {
	defer goleak.VerifyNone(t)

	config, err := symo.NewConfig("")
	require.NoError(t, err)

	log := new(mocks.Logger)
	log.On("Debug", mock.Anything)
	log.On("Debug", "tick ", mock.Anything)

	laErr := errors.New("LoadAvg Error")
	LoadAvg := new(mocks.LoadAvg)
	LoadAvg.On("Execute", mock.Anything).Return(nil, laErr)

	cpuErr := errors.New("CPU Error")
	CPU := new(mocks.CPU)
	CPU.On("Execute", mock.Anything, symo.StartMetric).Return(nil, nil)
	CPU.On("Execute", mock.Anything, symo.StopMetric).Return(nil, nil)
	CPU.On("Execute", mock.Anything, symo.GetMetric).Return(nil, cpuErr)

	ldErr := errors.New("LoadDisks Error")
	LoadDisks := new(mocks.LoadDisks)
	LoadDisks.On("Execute", mock.Anything, symo.StartMetric).Return(nil, nil)
	LoadDisks.On("Execute", mock.Anything, symo.StopMetric).Return(nil, nil)
	LoadDisks.On("Execute", mock.Anything, symo.GetMetric).Return(nil, ldErr)

	fsErr := errors.New("UsedFS Error")
	UsedFS := new(mocks.UsedFS)
	UsedFS.On("Execute", mock.Anything, symo.StartMetric).Return(nil, nil)
	UsedFS.On("Execute", mock.Anything, symo.StopMetric).Return(nil, nil)
	UsedFS.On("Execute", mock.Anything, symo.GetMetric).Return(nil, fsErr)

	collectors := symo.MetricCollectors{
		LoadAvg:   LoadAvg.Execute,
		CPU:       CPU.Execute,
		LoadDisks: LoadDisks.Execute,
		UsedFS:    UsedFS.Execute,
	}

	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	collectorService := NewCollector(log, config)
	collectorService.Start(startCtx, collectors, toClientsCh)

	time.Sleep(50 * time.Millisecond)

	time.Sleep(time.Second)
	data := <-toClientsCh
	// текущая секунда еще не заполнена, предыдущих нет - статистика должна быть пустой
	require.Len(t, data.Points, 0)

	// все коллекторы вернули ошибки, данные должны быть пустые
	time.Sleep(time.Second)
	data = <-toClientsCh
	require.Len(t, data.Points, 1)
	for _, point := range data.Points {
		require.Nil(t, point.LoadAvg)
		require.Nil(t, point.CPU)
		require.Nil(t, point.LoadDisks)
		require.Nil(t, point.UsedFS)
	}

	stopCtx := context.Background()
	collectorService.Stop(stopCtx)

	log.AssertExpectations(t)
}
