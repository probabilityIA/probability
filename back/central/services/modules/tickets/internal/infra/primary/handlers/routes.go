package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/tickets", middleware.JWT())
	{
		g.GET("", h.List)
		g.POST("", h.Create)
		g.GET(":id", h.Get)
		g.PUT(":id", h.Update)
		g.DELETE(":id", h.Delete)

		g.PATCH(":id/status", h.ChangeStatus)
		g.PATCH(":id/assign", h.Assign)
		g.PATCH(":id/area", h.ChangeArea)
		g.PATCH(":id/escalate", h.Escalate)

		g.GET(":id/comments", h.ListComments)
		g.POST(":id/comments", h.AddComment)

		g.GET(":id/attachments", h.ListAttachments)
		g.POST(":id/attachments", h.UploadAttachment)
		g.DELETE("attachments/:attachment_id", h.DeleteAttachment)

		g.GET(":id/history", h.ListHistory)
	}
}
