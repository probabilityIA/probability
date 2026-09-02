package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/infra/primary/handlers/response"
)

func (h *Handlers) Get(c *gin.Context) {
	if _, ok := h.requireSuperAdmin(c); !ok {
		return
	}
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	sprint, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.FromSprint(sprint))
}
