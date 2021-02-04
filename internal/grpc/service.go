package grpc

import (
	"context"
	"fmt"

	"github.com/anfilat/final-stats/internal/pb"
	"github.com/anfilat/final-stats/internal/symo"
)

type Service struct {
	pb.UnimplementedSymoServer

	ctx     context.Context // контекст приложения, сервис завершается по закрытию контекста
	clients symo.Clients
	log     symo.Logger
}

func NewService(ctx context.Context, log symo.Logger, clients symo.Clients) *Service {
	return &Service{
		ctx:     ctx,
		clients: clients,
		log:     log,
	}
}

func (s *Service) GetStats(req *pb.StatsRequest, srv pb.Symo_GetStatsServer) error {
	s.log.Debug("new client. Every ", req.N, " for ", req.M)

	ch, del := s.clients.NewClient(symo.NewClient{
		N: int(req.N),
		M: int(req.M),
	})
	defer del()

L:
	for {
		select {
		case <-s.ctx.Done():
			break L
		case <-srv.Context().Done():
			s.log.Debug("client disconnected")
			break L
		case data, ok := <-ch:
			if !ok {
				break L
			}

			if err := srv.Send(data); err != nil {
				s.log.Debug(fmt.Errorf("unable to send message: %w", err))
				break L
			}
		}
	}

	return nil
}
