package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	invoicingErrors "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// handleDomainError mapea errores de dominio a respuestas HTTP con mensajes en español.
// Llama c.JSON internamente; el handler debe hacer return después de invocarla.
func handleDomainError(c *gin.Context, err error, code string) {
	status, message := resolveInvoicingError(err)
	c.JSON(status, response.Error{
		Error:   code,
		Message: message,
	})
}

// resolveInvoicingError determina el HTTP status code y el mensaje en español
// apropiado para cada error de dominio del módulo de facturación.
func resolveInvoicingError(err error) (int, string) {
	switch {

	// ── 404 Not Found ─────────────────────────────────────────────────────────
	case errors.Is(err, invoicingErrors.ErrInvoiceNotFound):
		return http.StatusNotFound, "Factura no encontrada"
	case errors.Is(err, invoicingErrors.ErrConfigNotFound):
		return http.StatusNotFound, "Configuración de facturación no encontrada"
	case errors.Is(err, invoicingErrors.ErrProviderNotFound):
		return http.StatusNotFound, "Proveedor de facturación no encontrado"
	case errors.Is(err, invoicingErrors.ErrProviderTypeNotFound):
		return http.StatusNotFound, "Tipo de proveedor de facturación no encontrado"
	case errors.Is(err, invoicingErrors.ErrSyncLogNotFound):
		return http.StatusNotFound, "Registro de sincronización no encontrado"
	case errors.Is(err, invoicingErrors.ErrCreditNoteNotFound):
		return http.StatusNotFound, "Nota de crédito no encontrada"

	// ── 409 Conflict ──────────────────────────────────────────────────────────
	case errors.Is(err, invoicingErrors.ErrConfigAlreadyExists):
		return http.StatusConflict, "Ya existe una configuración de facturación para esta integración. Edítala para cambiar el proveedor, o elimínala para crear una nueva."
	case errors.Is(err, invoicingErrors.ErrActiveInvoicingConfigExists):
		return http.StatusConflict, "Ya existe una configuración de facturación activa para este negocio. Por favor desactívela antes de activar otra."
	case errors.Is(err, invoicingErrors.ErrInvoiceAlreadyExists):
		return http.StatusConflict, "Ya existe una factura para esta orden"
	case errors.Is(err, invoicingErrors.ErrInvoiceAlreadyIssued):
		return http.StatusConflict, "La factura ya fue emitida"
	case errors.Is(err, invoicingErrors.ErrInvoiceAlreadyCancelled):
		return http.StatusConflict, "La factura ya fue cancelada"
	case errors.Is(err, invoicingErrors.ErrOrderAlreadyInvoiced):
		return http.StatusConflict, "La orden ya tiene una factura asociada"
	case errors.Is(err, invoicingErrors.ErrCreditNoteAlreadyIssued):
		return http.StatusConflict, "La nota de crédito ya fue emitida"
	case errors.Is(err, invoicingErrors.ErrSyncInProgress):
		return http.StatusConflict, "Ya hay una sincronización en progreso"

	// ── 422 Unprocessable Entity (reglas de negocio) ───────────────────────────
	case errors.Is(err, invoicingErrors.ErrInvoiceCannotBeCancelled):
		return http.StatusUnprocessableEntity, "La factura no puede ser cancelada en su estado actual"
	case errors.Is(err, invoicingErrors.ErrCancelNotImplemented):
		return http.StatusUnprocessableEntity, "La cancelación de facturas no está implementada para este proveedor"
	case errors.Is(err, invoicingErrors.ErrOrderNotInvoiceable):
		return http.StatusUnprocessableEntity, "La orden no es facturable"
	case errors.Is(err, invoicingErrors.ErrConfigNotEnabled):
		return http.StatusUnprocessableEntity, "La configuración de facturación no está habilitada"
	case errors.Is(err, invoicingErrors.ErrAutoInvoiceNotEnabled):
		return http.StatusUnprocessableEntity, "La facturación automática no está habilitada en esta configuración"
	case errors.Is(err, invoicingErrors.ErrProviderNotActive):
		return http.StatusUnprocessableEntity, "El proveedor de facturación no está activo"
	case errors.Is(err, invoicingErrors.ErrProviderNotConfigured):
		return http.StatusUnprocessableEntity, "El proveedor de facturación no está configurado para esta integración"
	case errors.Is(err, invoicingErrors.ErrRetryNotAllowed):
		return http.StatusUnprocessableEntity, "No se permite reintentar esta factura en su estado actual"
	case errors.Is(err, invoicingErrors.ErrNoRetriesToCancel):
		return http.StatusUnprocessableEntity, "No hay reintentos pendientes para cancelar"
	case errors.Is(err, invoicingErrors.ErrMaxRetriesExceeded):
		return http.StatusUnprocessableEntity, "Se ha superado el número máximo de reintentos permitidos"
	case errors.Is(err, invoicingErrors.ErrCreditNoteAmountExceeds):
		return http.StatusUnprocessableEntity, "El monto de la nota de crédito supera el total de la factura"
	case errors.Is(err, invoicingErrors.ErrInvoiceNotIssued):
		return http.StatusUnprocessableEntity, "La factura debe estar emitida para crear una nota de crédito"
	// Filtros de facturación automática
	case errors.Is(err, invoicingErrors.ErrOrderBelowMinAmount):
		return http.StatusUnprocessableEntity, "El monto de la orden está por debajo del mínimo configurado"
	case errors.Is(err, invoicingErrors.ErrOrderAboveMaxAmount):
		return http.StatusUnprocessableEntity, "El monto de la orden supera el máximo configurado"
	case errors.Is(err, invoicingErrors.ErrOrderNotPaid):
		return http.StatusUnprocessableEntity, "La orden no está pagada"
	case errors.Is(err, invoicingErrors.ErrPaymentMethodNotAllowed):
		return http.StatusUnprocessableEntity, "El método de pago de la orden no está permitido para facturación"
	case errors.Is(err, invoicingErrors.ErrOrderTypeNotAllowed):
		return http.StatusUnprocessableEntity, "El tipo de orden no está permitido para facturación"
	case errors.Is(err, invoicingErrors.ErrOrderStatusExcluded):
		return http.StatusUnprocessableEntity, "El estado de la orden está excluido de la facturación"
	case errors.Is(err, invoicingErrors.ErrProductExcluded):
		return http.StatusUnprocessableEntity, "La orden contiene productos excluidos de la facturación"
	case errors.Is(err, invoicingErrors.ErrProductNotAllowed):
		return http.StatusUnprocessableEntity, "La orden contiene productos que no están en la lista permitida"
	case errors.Is(err, invoicingErrors.ErrMinItemsNotMet):
		return http.StatusUnprocessableEntity, "La orden no cumple el mínimo de ítems requeridos"
	case errors.Is(err, invoicingErrors.ErrMaxItemsExceeded):
		return http.StatusUnprocessableEntity, "La orden supera el máximo de ítems permitidos"
	case errors.Is(err, invoicingErrors.ErrCustomerTypeNotAllowed):
		return http.StatusUnprocessableEntity, "El tipo de cliente no está permitido para facturación"
	case errors.Is(err, invoicingErrors.ErrCustomerExcluded):
		return http.StatusUnprocessableEntity, "El cliente está excluido de la facturación"
	case errors.Is(err, invoicingErrors.ErrShippingRegionNotAllowed):
		return http.StatusUnprocessableEntity, "La región de envío no está permitida para facturación"
	case errors.Is(err, invoicingErrors.ErrOrderOutsideDateRange):
		return http.StatusUnprocessableEntity, "La orden está fuera del rango de fechas permitido para facturación"

	// ── 400 Bad Request (datos inválidos) ─────────────────────────────────────
	case errors.Is(err, invoicingErrors.ErrInvalidInvoiceData):
		return http.StatusBadRequest, "Datos de factura inválidos"
	case errors.Is(err, invoicingErrors.ErrInvalidProviderConfig):
		return http.StatusBadRequest, "Configuración de proveedor inválida"
	case errors.Is(err, invoicingErrors.ErrInvalidCredentials):
		return http.StatusBadRequest, "Credenciales inválidas"
	case errors.Is(err, invoicingErrors.ErrInvalidAmount):
		return http.StatusBadRequest, "Monto inválido"
	case errors.Is(err, invoicingErrors.ErrInvalidCurrency):
		return http.StatusBadRequest, "Moneda inválida"
	case errors.Is(err, invoicingErrors.ErrInvalidCustomerData):
		return http.StatusBadRequest, "Datos del cliente inválidos"
	case errors.Is(err, invoicingErrors.ErrInvalidOrderData):
		return http.StatusBadRequest, "Datos de la orden inválidos"
	case errors.Is(err, invoicingErrors.ErrMissingRequiredField):
		return http.StatusBadRequest, "Falta un campo requerido"
	case errors.Is(err, invoicingErrors.ErrInvalidFilterConfig):
		return http.StatusBadRequest, "Configuración de filtros inválida"

	// ── 401 Unauthorized ──────────────────────────────────────────────────────
	case errors.Is(err, invoicingErrors.ErrProviderUnauthorized):
		return http.StatusUnauthorized, "No autorizado con el proveedor de facturación"
	case errors.Is(err, invoicingErrors.ErrAuthenticationFailed):
		return http.StatusUnauthorized, "Autenticación fallida con el proveedor de facturación"
	case errors.Is(err, invoicingErrors.ErrTokenExpired):
		return http.StatusUnauthorized, "El token del proveedor de facturación ha expirado"
	case errors.Is(err, invoicingErrors.ErrTokenRefreshFailed):
		return http.StatusUnauthorized, "No se pudo renovar el token del proveedor de facturación"

	// ── 429 Too Many Requests ─────────────────────────────────────────────────
	case errors.Is(err, invoicingErrors.ErrProviderRateLimitExceeded):
		return http.StatusTooManyRequests, "Límite de solicitudes del proveedor de facturación excedido"

	// ── 500 Internal Server Error (default) ───────────────────────────────────
	default:
		return http.StatusInternalServerError, "Error interno del servidor de facturación"
	}
}
