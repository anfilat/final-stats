package clients

import (
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

const maxQueueLen = 100

type clientsList []*grpcClient

// данные клиента.
type grpcClient struct {
	n     int              // информация отправляется каждые N секунд
	m     int              // информация усредняется за M секунд
	ch    chan *symo.Stats // переданный клиенту канал
	after time.Time        // когда отправлять следующий пакет данных
	dead  bool             // контекст клиента закрыт, нужно удалить этого клиента из списка
}

func newClient(cl symo.ClientData, now time.Time) *grpcClient {
	ch := make(chan *symo.Stats, maxQueueLen)
	client := &grpcClient{
		n:    cl.N,
		m:    cl.M,
		ch:   ch,
		dead: false,
	}
	client.after = now.Add(time.Duration(client.m-1) * time.Second)
	return client
}

func (g *grpcClient) close() {
	close(g.ch)
}

func (g *grpcClient) isReady(now time.Time) bool {
	return now.After(g.after)
}

func (g *grpcClient) setNextReady(now time.Time) {
	g.after = now.Add(time.Duration(g.n-1) * time.Second)
}
