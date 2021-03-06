//go:generate protoc -I "/usr/local/include/" --proto_path=. --go_out=. --go-grpc_out=. ./symo.proto

package grpc

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/anfilat/final-stats/internal/symo"
)

type grpcServer struct {
	mutex  *sync.Mutex
	srv    *grpc.Server
	config symo.Config
	log    symo.Logger
}

// NewServer возвращает gRPC сервер.
func NewServer(log symo.Logger, config symo.Config) symo.GRPCServer {
	return &grpcServer{
		mutex:  &sync.Mutex{},
		config: config,
		log:    log,
	}
}

func (g *grpcServer) Start(addr string, clients symo.NewClienter) error {
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	g.mutex.Lock()
	g.srv = grpc.NewServer()
	g.mutex.Unlock()

	RegisterSymoServer(g.srv, newService(g.log, g.config, clients))

	g.log.Debug("starting grpc server on ", addr)
	return g.srv.Serve(lsn)
}

func (g *grpcServer) Stop(ctx context.Context) {
	stopped := make(chan interface{})
	go func() {
		g.mutex.Lock()
		defer g.mutex.Unlock()

		g.srv.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		g.srv.Stop()
	case <-stopped:
	}

	g.log.Debug("grpc server is stopped")
}
