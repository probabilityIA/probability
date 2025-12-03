package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/payments/app/usecases"
)

// ═══════════════════════════════════════════
// INTERFACE
// ═══════════════════════════════════════════

// IHandlers define la interfaz de handlers del módulo payments
type IHandlers interface {
	// Payment Methods
	ListPaymentMethods(c *gin.Context)
	GetPaymentMethod(c *gin.Context)
	CreatePaymentMethod(c *gin.Context)
	UpdatePaymentMethod(c *gin.Context)
	DeletePaymentMethod(c *gin.Context)
	TogglePaymentMethod(c *gin.Context)

	// Payment Mappings
	ListPaymentMappings(c *gin.Context)
	GetPaymentMapping(c *gin.Context)
	GetPaymentMappingsByIntegration(c *gin.Context)
	CreatePaymentMapping(c *gin.Context)
	UpdatePaymentMapping(c *gin.Context)
	DeletePaymentMapping(c *gin.Context)
	TogglePaymentMapping(c *gin.Context)

	// Routes
	RegisterRoutes(router *gin.RouterGroup)
}

// ═══════════════════════════════════════════
// CONSTRUCTOR
// ═══════════════════════════════════════════

type PaymentHandlers struct {
	uc usecases.IUseCase
}

// New crea una nueva instancia de los handlers
func New(uc usecases.IUseCase) IHandlers {
	return &PaymentHandlers{
		uc: uc,
	}
}
