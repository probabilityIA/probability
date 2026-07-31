package errors

import "errors"

var (
	ErrBusinessNotFound    = errors.New("negocio no encontrado")
	ErrProductNotFound     = errors.New("producto no encontrado")
	ErrInvalidContact      = errors.New("nombre y mensaje son requeridos")
	ErrPublicSiteNotActive = errors.New("el sitio web público no está activo para este negocio")
	ErrEmptyCart           = errors.New("el carrito esta vacio")
	ErrInvalidCartItem     = errors.New("cantidad invalida en un item del carrito")
	ErrCheckoutNotFound    = errors.New("checkout no encontrado")
	ErrOnlinePayNotReady   = errors.New("el pago en linea no esta disponible para esta tienda")
)
