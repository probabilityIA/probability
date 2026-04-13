package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type handler struct {
	uc  app.IUseCase
	log log.ILogger
}

func New(useCase app.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		uc:  useCase,
		log: logger.WithModule("announcements-handler"),
	}
}

func (h *handler) parseIDParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id parameter"})
		return 0, false
	}
	return uint(id), true
}
