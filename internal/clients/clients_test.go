package clients

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/anfilat/final-stats/internal/mocks"
	"github.com/anfilat/final-stats/internal/symo"
)

func TestClientsStartStop(t *testing.T) {
	defer goleak.VerifyNone(t)

	log := new(mocks.Logger)
	log.On("Debug", "clients is stopped")

	toCollectorCh := make(symo.ClientsToCollectorCh, 1)
	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	clientsService := NewClients(log, clock.NewMock())
	clientsService.Start(startCtx, toCollectorCh, toClientsCh)

	stopCtx := context.Background()
	clientsService.Stop(stopCtx)

	log.AssertExpectations(t)
	require.Len(t, toCollectorCh, 0)
}

func TestClientsStartWithCanceledContext(t *testing.T) {
	defer goleak.VerifyNone(t)

	log := new(mocks.Logger)
	log.On("Debug", "clients is stopped")

	toCollectorCh := make(symo.ClientsToCollectorCh, 1)
	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx, cancel := context.WithCancel(context.Background())
	cancel()
	clientsService := NewClients(log, clock.NewMock())
	clientsService.Start(startCtx, toCollectorCh, toClientsCh)

	stopCtx := context.Background()
	clientsService.Stop(stopCtx)

	log.AssertExpectations(t)
	require.Len(t, toCollectorCh, 0)
}

func TestCommandsToCollector(t *testing.T) {
	log := new(mocks.Logger)
	log.On("Debug", "clients is stopped")
	log.On("Debug", "stats sent in ", mock.Anything)

	toCollectorCh := make(symo.ClientsToCollectorCh, 1)
	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	clientsService := NewClients(log, clock.NewMock())
	clientsService.Start(startCtx, toCollectorCh, toClientsCh)

	// при добавлении первого клиента коллектору отправляется сообщение о том, что есть кому отправлять метрики
	_, del1, err := clientsService.NewClient(symo.ClientData{N: 1, M: 1})
	require.NoError(t, err)
	cmd := <-toCollectorCh
	require.Equal(t, symo.Start, cmd)

	// при добавлении дополнительных клиентов коллектору ничего не отправляется
	_, del2, err := clientsService.NewClient(symo.ClientData{N: 1, M: 1})
	require.NoError(t, err)
	require.Len(t, toCollectorCh, 0)

	// при удалении всех клиентов коллектору отправляется сообщение о том, что метрики отправлять не надо
	del1()
	del2()
	toClientsCh <- symo.MetricsData{
		Time:   time.Now(),
		Points: nil,
	}
	cmd = <-toCollectorCh
	require.Equal(t, symo.Stop, cmd)

	// если снова появляется клиент, то все повторяется
	_, _, err = clientsService.NewClient(symo.ClientData{N: 1, M: 1})
	require.NoError(t, err)
	cmd = <-toCollectorCh
	require.Equal(t, symo.Start, cmd)

	stopCtx := context.Background()
	clientsService.Stop(stopCtx)

	log.AssertExpectations(t)
}

func TestClosingChannelsOnClose(t *testing.T) {
	log := new(mocks.Logger)
	log.On("Debug", "clients is stopped")

	toCollectorCh := make(symo.ClientsToCollectorCh, 1)
	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	clientsService := NewClients(log, clock.NewMock())
	clientsService.Start(startCtx, toCollectorCh, toClientsCh)

	ch1, _, err := clientsService.NewClient(symo.ClientData{N: 1, M: 1})
	require.NoError(t, err)
	ch2, _, err := clientsService.NewClient(symo.ClientData{N: 1, M: 1})
	require.NoError(t, err)
	ch3, _, err := clientsService.NewClient(symo.ClientData{N: 1, M: 1})
	require.NoError(t, err)

	stopCtx := context.Background()
	clientsService.Stop(stopCtx)

	// все каналы grpc клиентов закрыты
	_, ok := <-ch1
	require.False(t, ok)
	_, ok = <-ch2
	require.False(t, ok)
	_, ok = <-ch3
	require.False(t, ok)

	log.AssertExpectations(t)
}

func TestFailNewClientAfterClose(t *testing.T) {
	log := new(mocks.Logger)
	log.On("Debug", "clients is stopped")

	toCollectorCh := make(symo.ClientsToCollectorCh, 1)
	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	clientsService := NewClients(log, clock.NewMock())
	clientsService.Start(startCtx, toCollectorCh, toClientsCh)

	stopCtx := context.Background()
	clientsService.Stop(stopCtx)

	// нельзя добавлять клиентов после закрытия
	_, _, err := clientsService.NewClient(symo.ClientData{N: 1, M: 1})
	require.ErrorIs(t, err, symo.ErrStopped)

	log.AssertExpectations(t)
}

//nolint:gocognit, funlen
func TestSend(t *testing.T) {
	type grpcClient struct {
		n     int
		m     int
		start int
		stop  int
		ticks []bool
	}
	tests := []struct {
		name  string
		time  int
		grpcs []grpcClient
	}{
		{
			name:  "no clients",
			time:  3,
			grpcs: nil,
		},
		{
			name: "clients n=1, m=1",
			time: 5,
			grpcs: []grpcClient{
				{
					n:     1,
					m:     1,
					start: 0,
					stop:  0,
					ticks: []bool{false, true, true, true, true},
				},
				{
					n:     1,
					m:     1,
					start: 1,
					stop:  4,
					ticks: []bool{false, false, true, true, false},
				},
			},
		},
		{
			name: "clients n=3, m=6",
			time: 10,
			grpcs: []grpcClient{
				{
					n:     3,
					m:     6,
					start: 0,
					stop:  0,
					ticks: []bool{false, false, false, false, false, false, true, false, false, true},
				},
				{
					n:     3,
					m:     6,
					start: 0,
					stop:  9,
					ticks: []bool{false, false, false, false, false, false, true, false, false, false},
				},
			},
		},
		{
			name: "clients n=1, m=1, n=2, m=2, n=3, m=6",
			time: 10,
			grpcs: []grpcClient{
				{
					n:     1,
					m:     1,
					start: 0,
					stop:  0,
					ticks: []bool{false, true, true, true, true, true, true, true, true, true},
				},
				{
					n:     2,
					m:     2,
					start: 0,
					stop:  0,
					ticks: []bool{false, false, true, false, true, false, true, false, true, false},
				},
				{
					n:     3,
					m:     6,
					start: 0,
					stop:  0,
					ticks: []bool{false, false, false, false, false, false, true, false, false, true},
				},
			},
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			log := new(mocks.Logger)
			log.On("Debug", mock.Anything)
			log.On("Debug", mock.Anything, mock.Anything)
			toCollectorCh := make(symo.ClientsToCollectorCh, 1)
			toClientsCh := make(symo.CollectorToClientsCh, 1)

			startCtx := context.Background()
			mockedClock := clock.NewMock()
			clientsService := NewClients(log, mockedClock)
			clientsService.Start(startCtx, toCollectorCh, toClientsCh)
			defer func() {
				stopCtx := context.Background()
				clientsService.Stop(stopCtx)
			}()

			type grpcClient struct {
				ch  <-chan *symo.Stats
				del func()
			}
			clients := make([]grpcClient, len(tt.grpcs))
			for tick := 0; tick < tt.time; tick++ {
				// старт новых клиентов
				for i, grpc := range tt.grpcs {
					if grpc.start == tick {
						ch, del, err := clientsService.NewClient(symo.ClientData{N: grpc.n, M: grpc.m})
						require.NoErrorf(t, err, "tick %d, client %d", tick, i)
						clients[i].ch = ch
						clients[i].del = del
					}
				}

				// остановка клиентов
				for i, grpc := range tt.grpcs {
					if grpc.stop == tick {
						require.NotNilf(t, clients[i].del, "tick %d, client %d", tick, i)
						clients[i].del()
						clients[i].ch = nil
					}
				}

				// секундный тик
				now := mockedClock.Now()
				toClientsCh <- symo.MetricsData{
					Time:   now.Truncate(time.Second),
					Points: nil,
				}

				time.Sleep(100 * time.Millisecond)

				// проверка клиентских каналов
				for i, grpc := range tt.grpcs {
					if clients[i].ch == nil {
						continue
					}
					if grpc.ticks[tick] {
						require.Lenf(t, clients[i].ch, 1, "tick %d, client %d", tick, i)
						<-clients[i].ch
					} else {
						require.Lenf(t, clients[i].ch, 0, "tick %d, client %d", tick, i)
					}
				}

				mockedClock.Add(time.Second)
			}
		})
	}
}
