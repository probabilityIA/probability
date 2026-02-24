package domain

import (
	"fmt"
	"sync"
)

// ShipmentRepository almacena los envios simulados en memoria
type ShipmentRepository struct {
	mu              sync.RWMutex
	shipmentsByID   map[string]*StoredShipment
	shipmentsByTrack map[string]*StoredShipment
	shipmentSeq     int
}

// NewShipmentRepository crea una nueva instancia del repositorio
func NewShipmentRepository() *ShipmentRepository {
	return &ShipmentRepository{
		shipmentsByID:    make(map[string]*StoredShipment),
		shipmentsByTrack: make(map[string]*StoredShipment),
		shipmentSeq:      5000,
	}
}

// SaveShipment guarda un envio con doble indice (por ID y por tracking)
func (r *ShipmentRepository) SaveShipment(shipment *StoredShipment) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shipmentsByID[shipment.ID] = shipment
	r.shipmentsByTrack[shipment.TrackingNumber] = shipment
}

// GetByID obtiene un envio por su ID
func (r *ShipmentRepository) GetByID(id string) (*StoredShipment, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, exists := r.shipmentsByID[id]
	return s, exists
}

// GetByTracking obtiene un envio por su tracking number
func (r *ShipmentRepository) GetByTracking(trackingNumber string) (*StoredShipment, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, exists := r.shipmentsByTrack[trackingNumber]
	return s, exists
}

// GetAll retorna todos los envios almacenados
func (r *ShipmentRepository) GetAll() []*StoredShipment {
	r.mu.RLock()
	defer r.mu.RUnlock()
	shipments := make([]*StoredShipment, 0, len(r.shipmentsByID))
	for _, s := range r.shipmentsByID {
		shipments = append(shipments, s)
	}
	return shipments
}

// MarkCancelled marca un envio como cancelado
func (r *ShipmentRepository) MarkCancelled(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, exists := r.shipmentsByID[id]
	if !exists {
		return false
	}
	s.Status = "cancelled"
	return true
}

// GenerateShipmentID genera un ID de envio secuencial
func (r *ShipmentRepository) GenerateShipmentID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shipmentSeq++
	return fmt.Sprintf("EC-%06d", r.shipmentSeq)
}
