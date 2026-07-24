package errors

import "errors"

var (
	ErrConceptNotFound      = errors.New("concepto contable no encontrado")
	ErrTaxNotFound          = errors.New("impuesto no encontrado")
	ErrEntryNotFound        = errors.New("movimiento contable no encontrado")
	ErrDuplicateCode        = errors.New("ya existe un registro con ese codigo")
	ErrInvalidKind          = errors.New("kind debe ser INCOME o EXPENSE")
	ErrInvalidAmount        = errors.New("el monto debe ser mayor a cero")
	ErrInvalidRate          = errors.New("la tarifa debe estar entre 0 y 100")
	ErrEntryNotManual       = errors.New("solo se pueden eliminar movimientos manuales")
	ErrConceptInactive      = errors.New("el concepto esta inactivo")
	ErrInvalidPeriod        = errors.New("rango de fechas invalido")
	ErrCodeRequired         = errors.New("el codigo es requerido")
	ErrNameRequired         = errors.New("el nombre es requerido")
	ErrConceptTaxNotAllowed = errors.New("no se puede asociar un impuesto inactivo")
)
