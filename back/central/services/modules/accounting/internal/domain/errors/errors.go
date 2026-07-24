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

	ErrInvoiceNotFound    = errors.New("factura no encontrada")
	ErrInvoiceNotDraft    = errors.New("solo se pueden modificar facturas en borrador")
	ErrInvoiceNoItems     = errors.New("la factura debe tener al menos un item")
	ErrInvoiceInvalidItem = errors.New("cada item requiere descripcion, cantidad > 0 y precio > 0")
	ErrInvoiceAlreadyPaid = errors.New("la factura ya esta pagada")
	ErrInvoiceCancelled   = errors.New("la factura esta cancelada")
	ErrInvoiceNoEmail     = errors.New("no hay correo destinatario: envia email_to o registra un usuario en el negocio")
	ErrInvoiceConceptKind = errors.New("el concepto de la factura debe ser de tipo INCOME")
	ErrBusinessNotFound   = errors.New("negocio no encontrado")
	ErrInvoiceNotPayable  = errors.New("solo se puede marcar pagada una factura enviada o en borrador")
)
