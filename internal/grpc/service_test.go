package grpc

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/anfilat/final-stats/internal/mocks"
	"github.com/anfilat/final-stats/internal/symo"
)

func TestGRPC(t *testing.T) {
	srv, listener, clientsService, log := startGRPCServer()
	defer stopGRPCServer(srv, listener)

	conn := getConnect(t, listener)
	defer conn.Close()

	log.On("Debug", "client disconnected")

	ch := make(chan *symo.Stats, 1)
	del := func() {}
	clientsService.On("NewClient", mock.Anything).Return((<-chan *symo.Stats)(ch), del, nil)

	client := NewSymoClient(conn)

	ctx := context.Background()
	req := &StatsRequest{
		N: 1,
		M: 1,
	}
	reqClient, err := client.GetStats(ctx, req)
	require.NoError(t, err)

	ch <- someStats()
	stats, err := reqClient.Recv()
	require.NoError(t, err)
	require.NotNil(t, stats)
	require.NotNil(t, stats.Time)
	require.NotNil(t, stats.LoadAvg)
	require.NotNil(t, stats.Cpu)
	require.NotNil(t, stats.LoadDisks)
	require.NotNil(t, stats.UsedFs)
}

func TestGRPCFails(t *testing.T) {
	tests := []struct {
		name    string
		n       int32
		m       int32
		message string
	}{
		{
			name:    "n = 0",
			n:       0,
			m:       1,
			message: "N must be greater than 0 seconds",
		},
		{
			name:    "n < 0",
			n:       -1,
			m:       1,
			message: "N must be greater than 0 seconds",
		},
		{
			name:    "m = 0",
			n:       1,
			m:       0,
			message: "M must be greater than 0 seconds",
		},
		{
			name:    "m < 0",
			n:       1,
			m:       -1,
			message: "M must be greater than 0 seconds",
		},
		{
			name:    "n > MaxSeconds",
			n:       symo.MaxSeconds + 1,
			m:       1,
			message: fmt.Sprintf("N must be less than %v seconds", symo.MaxSeconds),
		},
		{
			name:    "m > MaxSeconds",
			n:       1,
			m:       symo.MaxSeconds + 1,
			message: fmt.Sprintf("M must be less than %v seconds", symo.MaxSeconds),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			srv, listener, _, _ := startGRPCServer()
			defer stopGRPCServer(srv, listener)

			conn := getConnect(t, listener)
			defer conn.Close()

			client := NewSymoClient(conn)

			ctx := context.Background()
			req := &StatsRequest{
				N: tt.n,
				M: tt.m,
			}
			reqClient, err := client.GetStats(ctx, req)
			require.NoError(t, err)
			_, err = reqClient.Recv()
			require.NotNil(t, err)
			er, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, codes.InvalidArgument, er.Code())
			require.Equal(t, tt.message, er.Message())
		})
	}
}

func TestGRPCFailInClosingTime(t *testing.T) {
	srv, listener, clientsService, _ := startGRPCServer()
	defer stopGRPCServer(srv, listener)

	conn := getConnect(t, listener)
	defer conn.Close()

	client := NewSymoClient(conn)

	clientsService.On("NewClient", mock.Anything).Return(nil, nil, symo.ErrStopped)

	ctx := context.Background()
	req := &StatsRequest{
		N: 1,
		M: 1,
	}
	reqClient, err := client.GetStats(ctx, req)
	require.NoError(t, err)
	_, err = reqClient.Recv()
	require.NotNil(t, err)
	er, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unavailable, er.Code())
	require.Equal(t, "service is closing", er.Message())
}

func startGRPCServer() (*grpc.Server, *bufconn.Listener, *mocks.NewClienter, *mocks.Logger) {
	listener := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()

	log := new(mocks.Logger)
	log.On("Debug", "new client. Every ", mock.Anything, " for ", mock.Anything)

	config, _ := symo.NewConfig("")

	clientsService := new(mocks.NewClienter)

	RegisterSymoServer(srv, newService(log, config, clientsService))

	go func() {
		_ = srv.Serve(listener)
	}()

	return srv, listener, clientsService, log
}

func stopGRPCServer(srv *grpc.Server, listener io.Closer) {
	srv.GracefulStop()
	_ = listener.Close()
}

func dialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func getConnect(t *testing.T, listener *bufconn.Listener) *grpc.ClientConn {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(listener)))
	require.NoError(t, err)
	return conn
}

func someStats() *symo.Stats {
	return &symo.Stats{
		Time: time.Now().Truncate(time.Second),
		LoadAvg: &symo.LoadAvgData{
			Load1:  1,
			Load5:  1,
			Load15: 1,
		},
		CPU: &symo.CPUData{
			User:   1,
			System: 1,
			Idle:   1,
		},
		LoadDisks: symo.LoadDisksData{
			{
				Name:    "sda",
				Tps:     7,
				KBRead:  8,
				KBWrite: 9,
			},
		},
		UsedFS: symo.UsedFSData{
			{
				Path:      "/",
				UsedSpace: 12.3,
				UsedInode: 7.77,
			},
		},
	}
}
