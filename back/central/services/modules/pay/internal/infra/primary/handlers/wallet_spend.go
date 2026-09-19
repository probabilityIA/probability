package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/mappers"
)

func spendDateRange(c *gin.Context) (string, string) {
	start := c.Query("start_date")
	end := c.Query("end_date")
	if start != "" && end != "" {
		return start, end
	}
	now := time.Now()
	return now.AddDate(0, 0, -29).Format("2006-01-02"), now.Format("2006-01-02")
}

func spendPagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	if pageSize < 1 {
		pageSize = 15
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// GetSpendSummary maneja GET /pay/wallet/spend-summary: cuanto ha gastado el
// negocio de su propia billetera (guias, membresia, etc.), sin margenes.
func (h *walletHandler) GetSpendSummary(c *gin.Context) {
	businessID, ok := resolveBusinessID(c)
	if !ok {
		return
	}
	start, end := spendDateRange(c)

	summaryDTO := &dtos.SpendSummaryDTO{
		BusinessID: businessID,
		StartDate:  start,
		EndDate:    end,
	}

	totals, err := h.walletUC.GetSpendSummary(c.Request.Context(), summaryDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	guideStatusBreakdown, err := h.walletUC.GetGuideStatusSummary(c.Request.Context(), summaryDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":                true,
		"start_date":             start,
		"end_date":               end,
		"by_concept":             totals,
		"guide_status_breakdown": guideStatusBreakdown,
	})
}

// ListSpendTransactions maneja GET /pay/wallet/spend-transactions: detalle
// paginado de los debitos de billetera del negocio.
func (h *walletHandler) ListSpendTransactions(c *gin.Context) {
	businessID, ok := resolveBusinessID(c)
	if !ok {
		return
	}
	start, end := spendDateRange(c)
	page, pageSize := spendPagination(c)

	txs, total, err := h.walletUC.ListSpendTransactions(c.Request.Context(), &dtos.TransactionFilterDTO{
		BusinessID: businessID,
		StartDate:  start,
		EndDate:    end,
		Concept:    c.Query("concept"),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        mappers.WalletTxListToResponse(txs),
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}
