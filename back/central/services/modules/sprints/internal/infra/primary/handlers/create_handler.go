package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/infra/primary/handlers/response"
)

func (h *Handlers) Create(c *gin.Context) {
	userID, ok := h.requireSuperAdmin(c)
	if !ok {
		return
	}
	var req request.CreateSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sprint, err := h.uc.Create(c.Request.Context(), dtos.CreateSprintDTO{
		Name:        req.Name,
		Goal:        req.Goal,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      req.Status,
		CreatedByID: userID,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.FromSprint(sprint))
}
