package repository

import (
	"sync"

	"service/internal/domain"
)

type DBGetAllOrders interface {
	GetAllOrders() ([]domain.Order, error)
}

type cache struct {
	storage map[string]*domain.Order
	mtx     sync.RWMutex
}

func NewCache(db DBGetAllOrders) (*cache, error) {
	orders, err := db.GetAllOrders()
	if err != nil {
		return nil, err
	}

	c := cache{
		storage: make(map[string]*domain.Order),
	}

	for _, order := range orders {
		c.storage[order.UID] = &order
	}

	return &c, nil
}

func (c *cache) Stop() {}

func (c *cache) GetOrder(id string) *domain.Order {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.storage[id]
}

func (c *cache) WriteOrder(order *domain.Order) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.storage[order.UID] = order
}
