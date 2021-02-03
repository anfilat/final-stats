//go:generate protoc -I "/usr/local/include/" --proto_path=. --go_out=. --go-grpc_out=. ./symo.proto
package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/anfilat/final-stats/internal/symo"
)

type grpcServer struct {
	srv     *grpc.Server
	clients symo.Clients
	log     symo.Logger
}

func NewServer(log symo.Logger, clients symo.Clients) symo.GRPCServer {
	return &grpcServer{
		clients: clients,
		log:     log,
	}
}

func (g *grpcServer) Start(addr string) error {
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	g.srv = grpc.NewServer()
	RegisterSymoServer(g.srv, NewService(g.log, g.clients))

	g.log.Info("starting grpc server on ", addr)
	return g.srv.Serve(lsn)
}

func (g *grpcServer) Stop(_ context.Context) error {
	// TODO доработать
	g.srv.GracefulStop()
	g.log.Debug("grpc server is stopped")
	return nil
}
