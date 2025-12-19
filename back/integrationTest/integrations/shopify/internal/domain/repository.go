package domain

import (
	"sync"
)

// OrderRepository almacena órdenes en memoria
type OrderRepository struct {
	orders map[string]*Order
	mu     sync.RWMutex
}

// NewOrderRepository crea un nuevo repositorio de órdenes
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]*Order),
	}
}

// Save guarda una orden (por OrderNumber como clave)
func (r *OrderRepository) Save(order *Order) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.Name] = order
}

// Get obtiene una orden por su nombre (OrderNumber)
func (r *OrderRepository) Get(orderNumber string) (*Order, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, exists := r.orders[orderNumber]
	return order, exists
}

// GetAll obtiene todas las órdenes
func (r *OrderRepository) GetAll() []*Order {
	r.mu.RLock()
	defer r.mu.RUnlock()
	orders := make([]*Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}
	return orders
}

// Delete elimina una orden
func (r *OrderRepository) Delete(orderNumber string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.orders, orderNumber)
}

// Exists verifica si existe una orden
func (r *OrderRepository) Exists(orderNumber string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.orders[orderNumber]
	return exists
}

// Count retorna el número de órdenes almacenadas
func (r *OrderRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.orders)
}



