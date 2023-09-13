package reader

import (
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/nats-io/stan.go"

	"service/config"
	"service/internal/domain"
)

// Писатель в базу данных отправляет информацию, собственно, в базу
// и возвращает ошибку при её возникновении
type DBWriter interface {
	WriteOrder(order *domain.Order) error
}

type CacheWriter interface {
	WriteOrder(order *domain.Order)
}

type reader struct {
	db     DBWriter
	cache  CacheWriter
	sc     stan.Conn
	sub    stan.Subscription
	orders chan *domain.Order
	done   chan struct{} // Нужен для закрытия... ?
	wg     sync.WaitGroup
}

func NewReader(config config.Config, db DBWriter, cache CacheWriter) *reader {
	var err error

	r := &reader{
		db:     db,
		cache:  cache,
		orders: make(chan *domain.Order, config.OrdersBufSize),
		done:   make(chan struct{}),
	}

	r.sc, err = stan.Connect(config.NatsClusterID, "orders", stan.NatsURL(config.NatsUrl), stan.ConnectWait(5*time.Second))
	if err != nil {
		log.Fatalf("Failed to connect to NATS Streaming Server: %v", err)
	}

	r.sub, err = r.sc.Subscribe("order", r.Handler, stan.DurableName("my-durable"))
	if err != nil {
		log.Fatalf("Failed to subscribe to channel: %v", err)
	}

	go r.DBWriter()

	return r
}

func (r *reader) Stop() {
	r.sub.Close()
	r.sc.Close()
	close(r.done)
	close(r.orders)
	r.wg.Wait()
}
