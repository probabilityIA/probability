package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBulkJobStatus obtiene el estado de un job de facturación masiva
// GET /api/v1/invoicing/bulk-jobs/:id
func (h *handler) GetBulkJobStatus(c *gin.Context) {
	ctx := c.Request.Context()
	jobID := c.Param("id")

	// Obtener job status
	job, err := h.useCase.GetBulkJobStatus(ctx, jobID)
	if err != nil {
		h.log.Error(ctx).Err(err).Str("job_id", jobID).Msg("Failed to get bulk job status")
		c.JSON(404, gin.H{"error": "Job not found: " + err.Error()})
		return
	}

	// Construir response
	response := gin.H{
		"job_id":       job.ID,
		"status":       job.Status,
		"total_orders": job.TotalOrders,
		"processed":    job.Processed,
		"successful":   job.Successful,
		"failed":       job.Failed,
		"progress":     job.GetProgress(),
		"success_rate": job.GetSuccessRate(),
		"started_at":   job.StartedAt,
		"completed_at": job.CompletedAt,
		"created_at":   job.CreatedAt,
	}

	// Incluir items si se solicita
	if c.Query("include_items") == "true" {
		items, err := h.useCase.GetBulkJobItems(ctx, jobID)
		if err != nil {
			h.log.Warn(ctx).Err(err).Str("job_id", jobID).Msg("Failed to get bulk job items")
		} else {
			response["items"] = items
		}
	}

	c.JSON(200, response)
}

// ListBulkJobs lista los jobs de facturación masiva
// GET /api/v1/invoicing/bulk-jobs
func (h *handler) ListBulkJobs(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener business_id del contexto
	businessID, ok := ctx.Value("business_id").(uint)
	if !ok || businessID == 0 {
		h.log.Error(ctx).Msg("business_id not found in context")
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Parsear parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Validar límites
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Listar jobs
	jobs, total, err := h.useCase.ListBulkJobs(ctx, businessID, page, pageSize)
	if err != nil {
		h.log.Error(ctx).
			Err(err).
			Uint("business_id", businessID).
			Msg("Failed to list bulk jobs")
		c.JSON(500, gin.H{"error": "Failed to list jobs: " + err.Error()})
		return
	}

	// Calcular total de páginas
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	// Retornar respuesta paginada
	c.JSON(200, gin.H{
		"data":        jobs,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}
