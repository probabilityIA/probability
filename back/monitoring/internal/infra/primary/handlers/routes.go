package handlers

import "github.com/gin-gonic/gin"

func (h *handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.Health)

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", h.Login)
		auth.GET("/verify", h.jwtMiddleware(), h.Verify)
	}

	api := router.Group("/api/v1", h.jwtMiddleware())
	{
		containers := api.Group("/containers")
		{
			containers.GET("", h.ListContainers)
			containers.GET("/:id", h.GetContainer)
			containers.GET("/:id/stats", h.GetStats)
			containers.GET("/:id/logs", h.GetLogs)
			containers.GET("/:id/logs/stream", h.StreamLogs)
			containers.POST("/:id/restart", h.ContainerAction("restart"))
			containers.POST("/:id/stop", h.ContainerAction("stop"))
			containers.POST("/:id/start", h.ContainerAction("start"))
		}

		api.GET("/compose/services", h.GetComposeServices)
	}
}
