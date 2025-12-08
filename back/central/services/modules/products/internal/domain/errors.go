package domain

import "errors"

var (
	// ErrProductNotFound se retorna cuando un producto no existe
	ErrProductNotFound = errors.New("product not found")

	// ErrProductAlreadyExists se retorna cuando un producto con el mismo SKU ya existe
	ErrProductAlreadyExists = errors.New("product with this SKU already exists for this business")

	// ErrInvalidProductData se retorna cuando los datos del producto son inválidos
	ErrInvalidProductData = errors.New("invalid product data")
)

