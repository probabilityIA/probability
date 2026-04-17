package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecaseshipment"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type syncShipmentsRequest struct {
	Provider   string   `json:"provider"`
	DateFrom   string   `json:"date_from"`
	DateTo     string   `json:"date_to"`
	Statuses   []string `json:"statuses"`
	BusinessID uint     `json:"business_id"`
}

func (h *Handlers) SyncShipmentStatus(c *gin.Context) {
	ctx := c.Request.Context()

	var req syncShipmentsRequest
	_ = c.ShouldBindJSON(&req)

	if req.Provider == "" {
		req.Provider = domain.SyncProviderEnvioclick
	}

	businessID, err := h.resolveBusinessIDForSync(c, req.BusinessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	filter := domain.SyncShipmentsFilter{
		BusinessID: businessID,
		Provider:   req.Provider,
		Statuses:   req.Statuses,
	}

	if req.DateFrom != "" {
		if t, err := time.Parse("2006-01-02", req.DateFrom); err == nil {
			filter.DateFrom = &t
		}
	}
	if req.DateTo != "" {
		if t, err := time.Parse("2006-01-02", req.DateTo); err == nil {
			end := t.Add(24*time.Hour - time.Second)
			filter.DateTo = &end
		}
	}
	if filter.DateFrom == nil {
		from := time.Now().AddDate(0, 0, -30)
		filter.DateFrom = &from
	}
	if filter.DateTo == nil {
		to := time.Now()
		filter.DateTo = &to
	}

	syncUC := usecaseshipment.NewSyncShipments(h.uc.Repo(), h.transportPub)
	result, err := syncUC.SyncShipments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	if result.TotalShipments == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"total_shipments": 0,
			"message":         "No hay envios para sincronizar en el rango indicado",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success":                    true,
		"correlation_id":             result.CorrelationID,
		"total_shipments":            result.TotalShipments,
		"batches":                    result.Batches,
		"batch_size":                 result.BatchSize,
		"estimated_duration_seconds": result.EstimatedDurationSeconds,
		"message":                    "Sincronizacion iniciada. Los envios se actualizaran progresivamente.",
	})
}

func (h *Handlers) resolveBusinessIDForSync(c *gin.Context, bodyBusinessID uint) (uint, error) {
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		return 0, errors.New("no se pudo identificar la empresa")
	}
	if !middleware.IsSuperAdmin(c) {
		return businessID, nil
	}
	if bodyBusinessID == 0 {
		return 0, errors.New("super admin: business_id es requerido en el body")
	}
	return bodyBusinessID, nil
}
