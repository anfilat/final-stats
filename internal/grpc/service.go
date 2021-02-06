package grpc

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anfilat/final-stats/internal/pb"
	"github.com/anfilat/final-stats/internal/symo"
)

type Service struct {
	pb.UnimplementedSymoServer

	clients symo.Clients
	config  symo.Config
	log     symo.Logger
}

func newService(log symo.Logger, config symo.Config, clients symo.Clients) *Service {
	return &Service{
		clients: clients,
		config:  config,
		log:     log,
	}
}

func (s *Service) GetStats(req *pb.StatsRequest, srv pb.Symo_GetStatsServer) error {
	s.log.Debug("new client. Every ", req.N, " for ", req.M)

	n := int(req.N)
	m := int(req.M)

	MaxSeconds := s.config.App.MaxSeconds
	if n > MaxSeconds {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("N must be less than %v seconds", MaxSeconds))
	}
	if m > MaxSeconds {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("M must be less than %v seconds", MaxSeconds))
	}

	ch, del, err := s.clients.NewClient(symo.NewClient{N: n, M: m})
	if err != nil {
		return status.Error(codes.InvalidArgument, "service is closing")
	}
	defer del()

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

			if err := srv.Send(data); err != nil {
				s.log.Debug(fmt.Errorf("unable to send message: %w", err))
				break L
			}
		}
	}

	return nil
}
