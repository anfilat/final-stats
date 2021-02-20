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

	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	clientsService := NewClients(log, clock.NewMock())
	clientsService.Start(startCtx, toClientsCh)

	stopCtx := context.Background()
	clientsService.Stop(stopCtx)

	log.AssertExpectations(t)
}

func TestClientsStartWithCanceledContext(t *testing.T) {
	defer goleak.VerifyNone(t)

	log := new(mocks.Logger)
	log.On("Debug", "clients is stopped")

	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx, cancel := context.WithCancel(context.Background())
	cancel()
	clientsService := NewClients(log, clock.NewMock())
	clientsService.Start(startCtx, toClientsCh)

	stopCtx := context.Background()
	clientsService.Stop(stopCtx)

	log.AssertExpectations(t)
}

func TestClosingChannelsOnClose(t *testing.T) {
	log := new(mocks.Logger)
	log.On("Debug", "clients is stopped")

	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	clientsService := NewClients(log, clock.NewMock())
	clientsService.Start(startCtx, toClientsCh)

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

	toClientsCh := make(symo.CollectorToClientsCh, 1)

	startCtx := context.Background()
	clientsService := NewClients(log, clock.NewMock())
	clientsService.Start(startCtx, toClientsCh)

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
			toClientsCh := make(symo.CollectorToClientsCh, 1)

			startCtx := context.Background()
			mockedClock := clock.NewMock()
			clientsService := NewClients(log, mockedClock)
			clientsService.Start(startCtx, toClientsCh)
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

				// проверка клиентских каналов
				slept := false
				for i, grpc := range tt.grpcs {
					ch := clients[i].ch
					if ch == nil {
						continue
					}
					if grpc.ticks[tick] {
						require.Eventuallyf(t, func() bool {
							return len(ch) == 1
						}, 100*time.Millisecond, time.Millisecond, "tick %d, client %d", tick, i)
						<-ch
					} else {
						if !slept {
							time.Sleep(100 * time.Millisecond)
							slept = true
						}
						require.Lenf(t, ch, 0, "tick %d, client %d", tick, i)
					}
				}

				mockedClock.Add(time.Second)
			}
		})
	}
}
