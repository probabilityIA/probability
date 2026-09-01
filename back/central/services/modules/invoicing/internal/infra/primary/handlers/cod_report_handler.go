package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
)

// @Summary Obtener reporte de facturas contra entrega
// @Description Retorna cuantas facturas se emitieron, a que cuenta bancaria
// @Description se registro cada recibo de caja y el desglose contra entrega vs. contado
// @Tags invoicing-stats
// @Accept json
// @Produce json
// @Param business_id query uint true "ID del negocio"
// @Param start_date query string false "Fecha de inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha de fin (YYYY-MM-DD)"
// @Success 200 {object} entities.CODAccountReport
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /invoicing/invoices/cod-report [get]
func (h *handler) GetCODReport(c *gin.Context) {
	businessIDStr := c.Query("business_id")
	if businessIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id debe ser un numero valido"})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	report, err := h.useCase.GetCODAccountReport(c.Request.Context(), uint(businessID), startDate, endDate)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to get COD account report")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo reporte de contra entrega"})
		return
	}

	c.JSON(http.StatusOK, mappers.CODReportToResponse(report))
}

// @Summary Listar facturas de una cuenta bancaria del reporte contra entrega
// @Description Detalle paginado de facturas para una fila del desglose por cuenta
// @Tags invoicing-stats
// @Accept json
// @Produce json
// @Param business_id query uint true "ID del negocio"
// @Param account_number query string true "Numero de cuenta bancaria"
// @Param is_cod query bool true "Contra entrega (true) o pagada por adelantado (false)"
// @Param start_date query string false "Fecha de inicio (YYYY-MM-DD)"
// @Param end_date query string false "Fecha de fin (YYYY-MM-DD)"
// @Param page query int false "Pagina (default 1)"
// @Param page_size query int false "Tamano de pagina (default 10)"
// @Success 200 {object} response.InvoiceList
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /invoicing/invoices/cod-report/invoices [get]
func (h *handler) GetCODReportInvoices(c *gin.Context) {
	businessIDStr := c.Query("business_id")
	if businessIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id debe ser un numero valido"})
		return
	}

	accountNumber := c.Query("account_number")
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_number es requerido"})
		return
	}

	isCOD, err := strconv.ParseBool(c.Query("is_cod"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "is_cod debe ser true o false"})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	page := 1
	if p := c.Query("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}
	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		pageSize, _ = strconv.Atoi(ps)
	}
	if pageSize > 100 {
		pageSize = 100
	}

	invoices, total, err := h.useCase.GetCODAccountReportInvoices(c.Request.Context(), uint(businessID), startDate, endDate, accountNumber, isCOD, page, pageSize)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to get COD account report invoices")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo facturas del reporte de contra entrega"})
		return
	}

	baseURL, bucket := h.getS3Config()
	c.JSON(http.StatusOK, mappers.InvoicesToResponse(invoices, total, page, pageSize, baseURL, bucket))
}
