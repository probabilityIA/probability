package campaign

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

type lifecycleAction func(ctx *gin.Context, id, businessID uint) (*entities.Campaign, error)

func (h *handler) runLifecycle(c *gin.Context, action lifecycleAction) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	id, ok := parseID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}

	campaign, err := action(c, id, businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": campaign})
}

func (h *handler) Launch(c *gin.Context) {
	h.runLifecycle(c, func(ctx *gin.Context, id, businessID uint) (*entities.Campaign, error) {
		return h.useCase.Launch(ctx.Request.Context(), id, businessID)
	})
}

func (h *handler) Pause(c *gin.Context) {
	h.runLifecycle(c, func(ctx *gin.Context, id, businessID uint) (*entities.Campaign, error) {
		return h.useCase.Pause(ctx.Request.Context(), id, businessID)
	})
}

func (h *handler) Resume(c *gin.Context) {
	h.runLifecycle(c, func(ctx *gin.Context, id, businessID uint) (*entities.Campaign, error) {
		return h.useCase.Resume(ctx.Request.Context(), id, businessID)
	})
}

func (h *handler) Cancel(c *gin.Context) {
	h.runLifecycle(c, func(ctx *gin.Context, id, businessID uint) (*entities.Campaign, error) {
		return h.useCase.Cancel(ctx.Request.Context(), id, businessID)
	})
}
