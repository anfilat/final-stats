package clients

import (
	"context"
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

// длина очереди на отправку данных клиенту. На случай временного замедления сети.
const MaxQueueLen = 100

type clientsList []*grpcClient

// данные клиента.
type grpcClient struct {
	ctx   context.Context // контекст клиента
	n     int             // информация отправляется каждые N секунд
	m     int             // информация усредняется за M секунд
	ch    chan *symo.Stat // переданный клиенту канал
	after time.Time       // когда отправлять следующий пакет данных
	dead  bool            // контекст клиента закрыт, нужно удалить этого клиента из списка
}

func newClient(cl symo.NewClient) *grpcClient {
	ch := make(chan *symo.Stat, MaxQueueLen)
	client := &grpcClient{
		ctx:  cl.Ctx,
		n:    cl.N,
		m:    cl.M,
		ch:   ch,
		dead: false,
	}
	now := time.Now().Truncate(time.Second)
	client.after = now.Add(time.Duration(client.m-1) * time.Second)
	return client
}

func (g *grpcClient) isReady(now time.Time) bool {
	return now.After(g.after)
}

func (g *grpcClient) setNextReady(now time.Time) {
	g.after = now.Add(time.Duration(g.n-1) * time.Second)
}