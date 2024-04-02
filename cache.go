package golangwb1

import (
	"sync"
)

type Cache struct {
	mu     sync.RWMutex
	Orders map[string]Order
}

func NewCache() *Cache {
	return &Cache{
		Orders: make(map[string]Order),
	}
}
func (c *Cache) AddOrder(order Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Orders[order.OrderUID] = order
}

func (с *Cache) RestoreFromDatabase(orders []Order) {
	for _, order := range orders {
		с.AddOrder(order)
	}
}
func (c *Cache) GetOrderById(orderUID string) (Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, ok := c.Orders[orderUID]
	return order, ok
}
