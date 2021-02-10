package grpc

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/anfilat/final-stats/internal/symo"
)

type Service struct {
	UnimplementedSymoServer

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

func (s *Service) GetStats(req *StatsRequest, srv Symo_GetStatsServer) error {
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

	ch, del, err := s.clients.NewClient(symo.ClientData{N: n, M: m})
	if err != nil {
		return status.Error(codes.Unavailable, "service is closing")
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

			if err := srv.Send(dataToGRPC(data)); err != nil {
				s.log.Debug(fmt.Errorf("unable to send message: %w", err))
				break L
			}
		}
	}

	return nil
}

func dataToGRPC(data *symo.Stats) *Stats {
	result := &Stats{}
	result.Time = timestamppb.New(data.Time)

	if data.LoadAvg != nil {
		result.LoadAvg = &LoadAvg{
			Load1:  data.LoadAvg.Load1,
			Load5:  data.LoadAvg.Load5,
			Load15: data.LoadAvg.Load15,
		}
	}
	if data.CPU != nil {
		result.Cpu = &CPU{
			User:   data.CPU.User,
			System: data.CPU.System,
			Idle:   data.CPU.Idle,
		}
	}
	if data.LoadDisks != nil {
		for _, diskData := range data.LoadDisks {
			result.LoadDisks = append(result.LoadDisks, &LoadDisk{
				Name:    diskData.Name,
				Tps:     diskData.Tps,
				KBRead:  diskData.KBRead,
				KBWrite: diskData.KBWrite,
			})
		}
	}
	if data.UsedFS != nil {
		for _, fsData := range data.UsedFS {
			result.UsedFs = append(result.UsedFs, &UsedFS{
				Path:      fsData.Path,
				UsedSpace: fsData.UsedSpace,
				UsedInode: fsData.UsedInode,
			})
		}
	}
	return result
}
