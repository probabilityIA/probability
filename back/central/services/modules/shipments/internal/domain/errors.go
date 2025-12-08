package domain

import "errors"

var (
	// ErrShipmentNotFound se retorna cuando un envío no existe
	ErrShipmentNotFound = errors.New("shipment not found")

	// ErrShipmentAlreadyExists se retorna cuando un envío con el mismo tracking number ya existe para una orden
	ErrShipmentAlreadyExists = errors.New("shipment with this tracking number already exists for this order")

	// ErrInvalidShipmentData se retorna cuando los datos del envío son inválidos
	ErrInvalidShipmentData = errors.New("invalid shipment data")

	// ErrOrderIDRequired se retorna cuando el order_id es requerido
	ErrOrderIDRequired = errors.New("order_id is required")
)

