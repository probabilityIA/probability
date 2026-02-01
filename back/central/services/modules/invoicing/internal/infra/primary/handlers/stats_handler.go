package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ═══════════════════════════════════════════════════════════════
// HANDLERS DE ESTADÍSTICAS Y RESÚMENES
// ═══════════════════════════════════════════════════════════════

// GetSummary obtiene un resumen general de facturas con KPIs principales
// @Summary Obtener resumen de facturas
// @Description Retorna un resumen con KPIs principales de facturación
// @Tags invoicing-stats
// @Accept json
// @Produce json
// @Param business_id query uint true "ID del negocio"
// @Param period query string false "Período: today, week, month (default), year, all"
// @Success 200 {object} entities.InvoiceSummary
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /invoicing/summary [get]
func (h *handler) GetSummary(c *gin.Context) {
	// Obtener business_id del query param
	businessIDStr := c.Query("business_id")
	if businessIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id debe ser un número válido"})
		return
	}

	// Obtener período (default: month)
	period := c.DefaultQuery("period", "month")

	// Validar período
	validPeriods := []string{"today", "week", "month", "year", "all"}
	isValid := false
	for _, p := range validPeriods {
		if period == p {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "período inválido",
			"valid_periods": validPeriods,
		})
		return
	}

	// Ejecutar caso de uso
	summary, err := h.useCase.GetSummary(c.Request.Context(), uint(businessID), period)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to get invoice summary")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo resumen"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetStats obtiene estadísticas detalladas para dashboards
// @Summary Obtener estadísticas detalladas
// @Description Retorna estadísticas detalladas de facturación para dashboards
// @Tags invoicing-stats
// @Accept json
// @Produce json
// @Param business_id query uint true "ID del negocio"
// @Param start_date query string false "Fecha de inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha de fin (YYYY-MM-DD)"
// @Param integration_id query uint false "Filtrar por integración origen"
// @Param invoicing_integration_id query uint false "Filtrar por proveedor de facturación"
// @Success 200 {object} entities.DetailedStats
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /invoicing/stats [get]
func (h *handler) GetStats(c *gin.Context) {
	// Obtener business_id del query param
	businessIDStr := c.Query("business_id")
	if businessIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id debe ser un número válido"})
		return
	}

	// Construir filtros
	filters := make(map[string]interface{})

	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}
	if integrationIDStr := c.Query("integration_id"); integrationIDStr != "" {
		if integrationID, err := strconv.ParseUint(integrationIDStr, 10, 32); err == nil {
			filters["integration_id"] = uint(integrationID)
		}
	}
	if invoicingIntegrationIDStr := c.Query("invoicing_integration_id"); invoicingIntegrationIDStr != "" {
		if invoicingIntegrationID, err := strconv.ParseUint(invoicingIntegrationIDStr, 10, 32); err == nil {
			filters["invoicing_integration_id"] = uint(invoicingIntegrationID)
		}
	}

	// Ejecutar caso de uso
	stats, err := h.useCase.GetDetailedStats(c.Request.Context(), uint(businessID), filters)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to get detailed stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo estadísticas"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetTrends obtiene datos de tendencias temporales para gráficos
// @Summary Obtener tendencias de facturación
// @Description Retorna datos de tendencias temporales para visualización en gráficos
// @Tags invoicing-stats
// @Accept json
// @Produce json
// @Param business_id query uint true "ID del negocio"
// @Param start_date query string true "Fecha de inicio (YYYY-MM-DD)"
// @Param end_date query string true "Fecha de fin (YYYY-MM-DD)"
// @Param granularity query string false "Granularidad: day (default), week, month"
// @Param metric query string false "Métrica: count (default), amount, success_rate"
// @Success 200 {object} entities.TrendData
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /invoicing/trends [get]
func (h *handler) GetTrends(c *gin.Context) {
	// Obtener business_id del query param
	businessIDStr := c.Query("business_id")
	if businessIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id debe ser un número válido"})
		return
	}

	// Obtener fechas (requeridas)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date y end_date son requeridos (formato: YYYY-MM-DD)"})
		return
	}

	// Obtener parámetros opcionales con defaults
	granularity := c.DefaultQuery("granularity", "day")
	metric := c.DefaultQuery("metric", "count")

	// Ejecutar caso de uso
	trends, err := h.useCase.GetTrends(c.Request.Context(), uint(businessID), startDate, endDate, granularity, metric)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to get trends")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trends)
}
