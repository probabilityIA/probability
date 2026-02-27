package handlers

import (
	"strconv"

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

// resolveBusinessID obtiene el business_id efectivo.
// Para usuarios normales usa el del JWT.
// Para super admins (business_id=0 en JWT) lee el query param ?business_id=X.
func (h *Handlers) resolveBusinessID(c *gin.Context) (uint, bool) {
	businessID := c.GetUint("business_id")
	if businessID > 0 {
		return businessID, true
	}
	// Super admin: leer de query param
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			return uint(id), true
		}
	}
	return 0, false
}
