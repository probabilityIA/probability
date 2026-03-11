package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/app/usecases"
)

type SubscriptionHandler struct {
	uc *usecases.SubscriptionUsecase
}

func NewSubscriptionHandler(uc *usecases.SubscriptionUsecase) *SubscriptionHandler {
	return &SubscriptionHandler{uc: uc}
}

// getBusinessIDFromContext extrae el businessId del JWT (seteado por el middleware de auth)
func getBusinessIDFromContext(c *gin.Context) uint {
	claim, exists := c.Get("businessId")
	if exists {
		switch v := claim.(type) {
		case uint:
			return v
		case float64:
			return uint(v)
		case string:
			id, _ := strconv.ParseUint(v, 10, 32)
			return uint(id)
		}
	}
	return 0
}

// GetCurrentSubscription retorna la suscripción del negocio autenticado.
// Super admins pueden pasar ?businessId=X para ver la suscripción de otro negocio.
func (h *SubscriptionHandler) GetCurrentSubscription(c *gin.Context) {
	businessID := getBusinessIDFromContext(c)

	// Si hay un ?businessId en la query (para superadmin), usarlo
	if qbid := c.Query("businessId"); qbid != "" {
		if bid, err := strconv.ParseUint(qbid, 10, 32); err == nil && bid > 0 {
			businessID = uint(bid)
		}
	}

	if businessID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "business ID not found"})
		return
	}

	sub, err := h.uc.GetBusinessSubscription(c.Request.Context(), businessID)
	if err != nil {
		log.Printf("[Error] GetCurrentSubscription: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sub})
}

// RegisterPayment registra un pago de suscripción y activa el negocio.
// Solo puede ser llamado por el Super Admin.
type RegisterPaymentReq struct {
	BusinessID       uint    `json:"businessId"`
	Amount           float64 `json:"amount"`
	MonthsToAdd      int     `json:"monthsToAdd"`
	PaymentReference *string `json:"paymentReference"`
	Notes            *string `json:"notes"`
}

func (h *SubscriptionHandler) RegisterPayment(c *gin.Context) {
	var req RegisterPaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido: " + err.Error()})
		return
	}

	if req.BusinessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "businessId es requerido"})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount debe ser mayor a 0"})
		return
	}
	if req.MonthsToAdd <= 0 {
		req.MonthsToAdd = 1
	}

	err := h.uc.RegisterSubscriptionPayment(
		c.Request.Context(),
		req.BusinessID,
		req.Amount,
		req.MonthsToAdd,
		req.PaymentReference,
		req.Notes,
	)
	if err != nil {
		log.Printf("[Error] RegisterPayment businessID=%d: %v\n", req.BusinessID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al registrar pago: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pago registrado. El negocio ahora está activo."})
}

// DisableSubscription suspende manualmente la suscripción de un negocio.
func (h *SubscriptionHandler) DisableSubscription(c *gin.Context) {
	businessIDStr := c.Query("businessId")
	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil || businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "businessId inválido"})
		return
	}

	err = h.uc.DisableBusinessSubscription(c.Request.Context(), uint(businessID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to disable subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Suscripción suspendida correctamente"})
}
