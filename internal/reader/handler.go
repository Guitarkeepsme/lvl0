package reader

import (
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"

	"service/internal/domain"
)

func (r *reader) Handler(msg *stan.Msg) {
	order := &domain.Order{}

	if err := json.Unmarshal(msg.Data, order); err != nil {
		log.Printf("reader handler error: %v", err)
	}

	select {
	case <-r.done:
	case r.orders <- order:
	}
}

func (r *reader) DBWriter() {
	r.wg.Add(1)
	defer r.wg.Done()

	for order := range r.orders {
		err := r.db.WriteOrder(order)
		if err != nil {
			log.Printf("write db order error: %v", err)
		}

		r.cache.WriteOrder(order)
	}
}
