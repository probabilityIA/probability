package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/app"
)

type IHandlers interface {
	ListConcepts(c *gin.Context)
	CreateConcept(c *gin.Context)
	UpdateConcept(c *gin.Context)
	SetConceptTax(c *gin.Context)
	ListTaxes(c *gin.Context)
	CreateTax(c *gin.Context)
	UpdateTax(c *gin.Context)
	ListEntries(c *gin.Context)
	CreateEntry(c *gin.Context)
	DeleteEntry(c *gin.Context)
	Report(c *gin.Context)
	Sync(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc app.IUseCase
}

func New(uc app.IUseCase) IHandlers {
	return &Handlers{uc: uc}
}

func (h *Handlers) requireSuperAdmin(c *gin.Context) bool {
	if c.GetBool("is_super_admin") {
		return true
	}
	c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "solo el super admin puede acceder a contabilidad"})
	c.Abort()
	return false
}
