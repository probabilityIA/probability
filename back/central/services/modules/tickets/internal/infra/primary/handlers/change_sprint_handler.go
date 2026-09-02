package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/infra/primary/handlers/response"
)

func (h *Handlers) ChangeSprint(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	var req request.ChangeSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _, _ := h.requesterContext(c)
	t, err := h.uc.ChangeSprint(c.Request.Context(), dtos.ChangeSprintDTO{
		TicketID:    id,
		SprintID:    req.SprintID,
		ChangedByID: userID,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.FromTicket(t))
}
