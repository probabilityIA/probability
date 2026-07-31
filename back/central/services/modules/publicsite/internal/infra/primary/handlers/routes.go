package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	pub := router.Group("/public/tienda")
	{
		pub.GET("/:slug", h.GetBusinessPage)
		pub.GET("/:slug/catalog", h.ListCatalog)
		pub.GET("/:slug/categories", h.GetCategories)
		pub.GET("/:slug/product/:id", h.GetProduct)
		pub.POST("/:slug/contact", h.SubmitContact)
		pub.POST("/:slug/checkout/bold/signature", h.CreateCheckoutSession)
		pub.GET("/:slug/checkout/:reference/status", h.GetCheckoutStatus)
		pub.GET("/:slug/session", middleware.JWT(), h.GetSession)
		pub.GET("/:slug/my-prices", middleware.JWT(), h.GetMyPrices)
		pub.GET("/:slug/my-orders", middleware.JWT(), h.GetMyOrders)
		pub.PUT("/:slug/profile", middleware.JWT(), h.UpdateProfile)
		pub.POST("/:slug/profile/avatar", middleware.JWT(), h.UploadAvatar)
		pub.GET("/:slug/addresses", middleware.JWT(), h.ListAddresses)
		pub.POST("/:slug/addresses", middleware.JWT(), h.AddAddress)
		pub.PUT("/:slug/addresses/:id/primary", middleware.JWT(), h.SetPrimaryAddress)
		pub.DELETE("/:slug/addresses/:id", middleware.JWT(), h.DeleteAddress)
	}
}
