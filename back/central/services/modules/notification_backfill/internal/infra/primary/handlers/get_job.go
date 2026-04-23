package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
)

func (h *Handlers) GetJob(c *gin.Context) {
	jobID := c.Param("job_id")
	job, ok := h.uc.GetJob(jobID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "job no encontrado"})
		return
	}

	resp := dtos.JobResponse{
		ID:            job.ID,
		EventCode:     job.EventCode,
		BusinessID:    job.BusinessID,
		Status:        string(job.Status),
		TotalEligible: job.TotalEligible,
		Sent:          job.Sent,
		Skipped:       job.Skipped,
		Failed:        job.Failed,
		StartedAt:     job.StartedAt,
		FinishedAt:    job.FinishedAt,
		ErrorMessage:  job.ErrorMessage,
	}
	c.JSON(http.StatusOK, resp)
}
