package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		families := products.Group("/families")
		{
			families.GET("", middleware.JWT(), h.ListProductFamilies)
			families.GET("/:family_id", middleware.JWT(), h.GetProductFamilyByID)
			families.GET("/:family_id/variants", middleware.JWT(), h.ListProductFamilyVariants)
			families.POST("", middleware.JWT(), h.CreateProductFamily)
			families.PUT("/:family_id", middleware.JWT(), h.UpdateProductFamily)
			families.DELETE("/:family_id", middleware.JWT(), h.DeleteProductFamily)
		}

		products.GET("/lookup-by-external", middleware.JWT(), h.LookupProductByExternalRef)

		products.GET("", middleware.JWT(), h.ListProducts)
		products.GET("/:id", middleware.JWT(), h.GetProductByID)
		products.POST("", middleware.JWT(), h.CreateProduct)
		products.PUT("/:id", middleware.JWT(), h.UpdateProduct)
		products.DELETE("/:id", middleware.JWT(), h.DeleteProduct)

		products.POST("/:id/image", middleware.JWT(), h.UploadProductImage)

		products.POST("/:id/integrations", middleware.JWT(), h.AddProductIntegration)
		products.GET("/:id/integrations", middleware.JWT(), h.GetProductIntegrations)
		products.PUT("/:id/integrations/:integration_id", middleware.JWT(), h.UpdateProductIntegration)
		products.DELETE("/:id/integrations/:integration_id", middleware.JWT(), h.RemoveProductIntegration)
	}
}
