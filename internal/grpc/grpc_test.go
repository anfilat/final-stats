package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/anfilat/final-stats/internal/mocks"
	"github.com/anfilat/final-stats/internal/symo"
)

func TestGRPCStartStop(t *testing.T) {
	defer goleak.VerifyNone(t)

	config, _ := symo.NewConfig("")

	log := new(mocks.Logger)
	log.On("Debug", "grpc server is stopped")
	log.On("Debug", "starting grpc server on ", mock.Anything)

	clientsService := new(mocks.NewClienter)

	grpcServer := NewServer(log, config)
	go func() {
		err := grpcServer.Start(":"+config.Server.Port, clientsService)
		require.NoError(t, err)
	}()

	time.Sleep(50 * time.Millisecond)

	stopCtx := context.Background()
	grpcServer.Stop(stopCtx)

	log.AssertExpectations(t)
}
