package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

type dataSummaryCellResponse struct {
	IntegrationID   uint   `json:"integration_id"`
	IntegrationName string `json:"integration_name"`
	ChannelCode     string `json:"channel_code"`
	CanFill         int    `json:"can_fill"`
	CanOverwrite    int    `json:"can_overwrite"`
}

type dataSummaryRowResponse struct {
	Field string                    `json:"field"`
	Label string                    `json:"label"`
	Note  string                    `json:"note,omitempty"`
	Cells []dataSummaryCellResponse `json:"cells"`
}

func (h *handler) DataSummary(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c, queryBusinessID(c))
	if !ok {
		return
	}

	summary, err := h.useCase.DataSummary(c.Request.Context(), businessID)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Uint("business_id", businessID).
			Msg("Error al resumir los datos que se pueden traer de los canales")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "no se pudo calcular el resumen"})
		return
	}

	rows := make([]dataSummaryRowResponse, 0, len(summary.Rows))
	for _, row := range summary.Rows {
		cells := make([]dataSummaryCellResponse, 0, len(row.Cells))
		for _, cell := range row.Cells {
			cells = append(cells, dataSummaryCellResponse{
				IntegrationID:   cell.IntegrationID,
				IntegrationName: cell.IntegrationName,
				ChannelCode:     cell.ChannelCode,
				CanFill:         cell.CanFill,
				CanOverwrite:    cell.CanOverwrite,
			})
		}
		rows = append(rows, dataSummaryRowResponse{Field: row.Field, Label: row.Label, Note: row.Note, Cells: cells})
	}

	var snapshotAt *time.Time
	if summary.Snapshotal != nil {
		snapshotAt = summary.Snapshotal
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows, "snapshot_at": snapshotAt})
}

type dataPreviewRequest struct {
	BusinessID    *uint    `json:"business_id"`
	IntegrationID uint     `json:"integration_id" binding:"required"`
	Field         string   `json:"field" binding:"required"`
	Mode          string   `json:"mode"`
	SKUs          []string `json:"skus"`
}

type dataConflictResponse struct {
	Source          string     `json:"source"`
	IntegrationID   *uint      `json:"integration_id,omitempty"`
	IntegrationName string     `json:"integration_name,omitempty"`
	Count           int        `json:"count"`
	LastChangeAt    *time.Time `json:"last_change_at,omitempty"`
}

type dataSampleResponse struct {
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Name      string `json:"name,omitempty"`
	Current   string `json:"current,omitempty"`
	Incoming  string `json:"incoming,omitempty"`
}

func (h *handler) DataPreview(c *gin.Context) {
	var req dataPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id y field son requeridos"})
		return
	}

	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	preview, err := h.useCase.DataPreview(c.Request.Context(), domain.DataPreviewQuery{
		BusinessID:    businessID,
		IntegrationID: req.IntegrationID,
		Field:         strings.TrimSpace(req.Field),
		Mode:          strings.TrimSpace(req.Mode),
		SKUs:          req.SKUs,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidDataField) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrIntegrationNotOwned) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": err.Error()})
			return
		}
		h.logger.Error(c.Request.Context()).Err(err).Uint("business_id", businessID).
			Str("field", req.Field).Msg("Error al calcular la vista previa de datos")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "no se pudo calcular la vista previa"})
		return
	}

	conflicts := make([]dataConflictResponse, 0, len(preview.Conflicts))
	for _, conflict := range preview.Conflicts {
		conflicts = append(conflicts, dataConflictResponse{
			Source:          conflict.Source,
			IntegrationID:   conflict.IntegrationID,
			IntegrationName: conflict.IntegrationName,
			Count:           conflict.Count,
			LastChangeAt:    conflict.LastChangeAt,
		})
	}

	samples := make([]dataSampleResponse, 0, len(preview.Samples))
	for _, sample := range preview.Samples {
		samples = append(samples, dataSampleResponse{
			ProductID: sample.ProductID,
			SKU:       sample.SKU,
			Name:      sample.Name,
			Current:   sample.Current,
			Incoming:  sample.Incoming,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"field":         preview.Field,
		"mode":          preview.Mode,
		"would_fill":    preview.WouldFill,
		"would_replace": preview.WouldReplace,
		"conflicts":     conflicts,
		"samples":       samples,
	})
}

func queryBusinessID(c *gin.Context) *uint {
	raw := c.Query("business_id")
	if raw == "" {
		return nil
	}
	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return nil
	}
	value := uint(parsed)
	return &value
}
