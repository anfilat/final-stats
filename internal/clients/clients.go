package clients

import (
	"context"
	"sync"

	"github.com/anfilat/final-stats/internal/symo"
)

// длина очереди на отправку данных клиенту. На случай временного замедления сети.
const MaxQueueLen = 100

type clients struct {
	ctx           context.Context // контекст приложения, сервис завершается по закрытию контекста
	clientsMutex  *sync.Mutex
	clients       []*grpcClient // список клиентов
	toHeartChan   chan<- symo.HeartCommand
	toClientsChan <-chan symo.ClientsBeat
	log           symo.Logger
}

// данные клиента.
type grpcClient struct {
	ctx  context.Context  // контекст клиента
	n    int              // информация отправляется каждые N секунд
	m    int              // информация усредняется за M секунд
	ch   chan<- symo.Stat // переданный клиенту канал
	dead bool             // контекст клиента закрыт, нужно удалить этого клиента из списка
}

func NewClients(ctx context.Context, log symo.Logger,
	toHeartChan chan<- symo.HeartCommand, toClientsChan <-chan symo.ClientsBeat) symo.Clients {
	return &clients{
		ctx:           ctx,
		clientsMutex:  &sync.Mutex{},
		toHeartChan:   toHeartChan,
		toClientsChan: toClientsChan,
		log:           log,
	}
}

func (c *clients) Start(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer func() {
			close(c.toHeartChan)
			c.log.Debug("clients is stopped")
			wg.Done()
		}()

		for {
			select {
			case <-c.ctx.Done():
				return
			case data, ok := <-c.toClientsChan:
				if !ok {
					return
				}
				c.sendStat(&data)
			}
		}
	}()
}

// подключение нового клиента.
func (c *clients) NewClient(client symo.NewClient) <-chan symo.Stat {
	ch := make(chan symo.Stat, MaxQueueLen)
	cl := &grpcClient{
		ctx:  client.Ctx,
		n:    client.N,
		m:    client.M,
		ch:   ch,
		dead: false,
	}

	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	c.clients = append(c.clients, cl)

	if len(c.clients) == 1 {
		select {
		case <-c.ctx.Done():
		case c.toHeartChan <- symo.Start:
		default:
		}
	}

	return ch
}

func (c *clients) sendStat(data *symo.ClientsBeat) {
	// TODO заменить на реальный код
	stat := &symo.Stat{
		Time: data.Time,
		Stat: &symo.Point{
			LoadAvg: &symo.LoadAvgData{
				Load1:  0,
				Load5:  0,
				Load15: 0,
			},
		},
	}

	isDead := c.sendToClients(stat)
	if isDead {
		c.log.Debug("isDead")
		c.delDeadClients()
	}
}

func (c *clients) sendToClients(stat *symo.Stat) (isDead bool) {
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	for _, client := range c.clients {
		select {
		case <-client.ctx.Done():
			client.dead = true
			isDead = true
			break
		default:
		}

		select {
		case <-c.ctx.Done():
			return
		case <-client.ctx.Done():
			client.dead = true
			isDead = true
		case client.ch <- *stat:
		default:
		}
	}
	return
}

func (c *clients) delDeadClients() {
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	clients := make([]*grpcClient, 0, len(c.clients))
	for _, client := range c.clients {
		if client.dead {
			continue
		}
		clients = append(clients, client)
	}
	c.clients = clients

	if len(c.clients) == 0 {
		select {
		case c.toHeartChan <- symo.Stop:
		default:
		}
	}
}
