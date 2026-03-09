package handlers

import "github.com/gin-gonic/gin"

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	pub := router.Group("/public/tienda")
	{
		pub.GET("/:slug", h.GetBusinessPage)
		pub.GET("/:slug/catalog", h.ListCatalog)
		pub.GET("/:slug/product/:id", h.GetProduct)
		pub.POST("/:slug/contact", h.SubmitContact)
	}
}
