package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.Health)

	mock := router.Group("/_mock")
	{
		mock.GET("/sent", h.ListSent)
		mock.GET("/templates", h.ListTemplates)
		mock.GET("/script", h.GetScript)
		mock.POST("/script", h.SetScript)
		mock.POST("/reply", h.ManualReply)
		mock.POST("/approve", h.ApproveByName)
		mock.POST("/reset", h.Reset)
		mock.GET("/media", h.ListMedia)
		mock.GET("/media/:id", h.DownloadMedia)
		mock.POST("/inbound-media", h.InboundMedia)
	}

	router.POST("/:id/messages", h.SendMessage)
	router.POST("/:id/message_templates", h.CreateTemplate)
	router.GET("/:id/message_templates", h.ListMetaTemplates)
	router.DELETE("/:id/message_templates", h.DeleteMetaTemplate)
	router.POST("/:id/uploads", h.StartUpload)
	router.POST("/:id/media", h.UploadMedia)
	router.GET("/:id", h.GetMediaInfo)
	router.POST("/:id", h.FinishUpload)
}
