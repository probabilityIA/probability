package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/infra/primary/handlers/response"
)

func (h *Handlers) ChangeStatus(c *gin.Context) {
	if _, ok := h.requireSuperAdmin(c); !ok {
		return
	}
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	var req request.ChangeSprintStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sprint, err := h.uc.ChangeStatus(c.Request.Context(), dtos.ChangeSprintStatusDTO{
		SprintID: id,
		Status:   req.Status,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.FromSprint(sprint))
}
