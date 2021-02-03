package grpc

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/anfilat/final-stats/internal/symo"
)

type Service struct {
	UnimplementedSymoServer

	clients symo.Clients
	log     symo.Logger
}

func NewService(log symo.Logger, clients symo.Clients) *Service {
	return &Service{
		clients: clients,
		log:     log,
	}
}

func (s *Service) GetStats(req *StatsRequest, srv Symo_GetStatsServer) error {
	s.log.Debug("new client. Every ", req.N, " for ", req.M)

	ctx, cancel := context.WithCancel(srv.Context())
	defer cancel()

	ch := s.clients.NewClient(symo.NewClient{
		Ctx: ctx,
		N:   int(req.N),
		M:   int(req.M),
	})

L:
	for {
		select {
		case <-srv.Context().Done():
			s.log.Debug("client disconnected")
			break L
		case data, ok := <-ch:
			if !ok {
				break L
			}

			message := dataToGRPC(&data)
			if err := srv.Send(message); err != nil {
				s.log.Debug(fmt.Errorf("unable to send message: %w", err))
				break L
			}
		}
	}

	return nil
}

func dataToGRPC(data *symo.Stat) *Stats {
	stat := data.Stat

	return &Stats{
		Time: timestamppb.New(data.Time),
		LoadAvg: &LoadAvg{
			Load1:  stat.LoadAvg.Load1,
			Load5:  stat.LoadAvg.Load5,
			Load15: stat.LoadAvg.Load15,
		},
	}
}
