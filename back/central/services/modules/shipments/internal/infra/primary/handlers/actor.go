package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

type actor struct {
	ID   *uint
	Name string
}

func (h *Handlers) resolveActor(c *gin.Context) actor {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		return actor{}
	}

	name := h.uc.Repo().GetUserDisplayName(c.Request.Context(), userID)
	if name == "" {
		if email, ok := middleware.GetUserEmail(c); ok {
			name = email
		}
	}

	id := userID
	return actor{ID: &id, Name: name}
}
