package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/app"
)

// IHandlers define la interfaz de handlers del módulo clients
type IHandlers interface {
	ListClients(c *gin.Context)
	GetClient(c *gin.Context)
	CreateClient(c *gin.Context)
	UpdateClient(c *gin.Context)
	DeleteClient(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

// Handlers contiene el use case
type Handlers struct {
	uc app.IUseCase
}

// New crea una nueva instancia de los handlers
func New(uc app.IUseCase) IHandlers {
	return &Handlers{uc: uc}
}
