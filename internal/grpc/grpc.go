package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/anfilat/final-stats/internal/pb"
	"github.com/anfilat/final-stats/internal/symo"
)

type grpcServer struct {
	srv *grpc.Server
	log symo.Logger
}

func NewServer(log symo.Logger) symo.GRPCServer {
	return &grpcServer{
		log: log,
	}
}

func (g *grpcServer) Start(addr string, clients symo.Clients) error {
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	g.srv = grpc.NewServer()
	pb.RegisterSymoServer(g.srv, newService(g.log, clients))

	g.log.Debug("starting grpc server on ", addr)
	return g.srv.Serve(lsn)
}

func (g *grpcServer) Stop(ctx context.Context) {
	stopped := make(chan interface{})
	go func() {
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
