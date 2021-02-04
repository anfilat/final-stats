package grpc

import (
	"context"
	"fmt"

	"github.com/anfilat/final-stats/internal/symo"
)

type Service struct {
	UnimplementedSymoServer

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

func (s *Service) GetStats(req *StatsRequest, srv Symo_GetStatsServer) error {
	s.log.Debug("new client. Every ", req.N, " for ", req.M)

	message := newMessage()
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

			dataToGRPC(data, message)
			if err := srv.Send(message); err != nil {
				s.log.Debug(fmt.Errorf("unable to send message: %w", err))
				break L
			}
		}
	}

	return nil
}

func newMessage() *Stats {
	return &Stats{
		LoadAvg: &LoadAvg{
			Load1:  0,
			Load5:  0,
			Load15: 0,
		},
	}
}

func dataToGRPC(data *symo.Stat, message *Stats) {
	stat := data.Stat

	message.Time = data.Time
	message.LoadAvg.Load1 = stat.LoadAvg.Load1
	message.LoadAvg.Load5 = stat.LoadAvg.Load5
	message.LoadAvg.Load15 = stat.LoadAvg.Load15
}
