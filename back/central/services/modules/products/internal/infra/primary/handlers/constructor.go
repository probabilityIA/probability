package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

// Handlers contiene todos los handlers del módulo products
type Handlers struct {
	uc  *usecases.UseCases
	log log.ILogger
	s3  storage.IS3Service
	env env.IConfig
}

// New crea una nueva instancia de Handlers
func New(uc *usecases.UseCases, logger log.ILogger, s3 storage.IS3Service, environment env.IConfig) *Handlers {
	return &Handlers{
		uc:  uc,
		log: logger,
		s3:  s3,
		env: environment,
	}
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

// respondBusinessIDRequired retorna un 400 cuando no se puede resolver el business_id
func (h *Handlers) respondBusinessIDRequired(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": "business_id es requerido",
		"error":   "business_id is required",
	})
}
