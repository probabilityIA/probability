package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandlers interface {
	ListEvents(c *gin.Context)
	Preview(c *gin.Context)
	Run(c *gin.Context)
	GetJob(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc  ports.IUseCase
	log log.ILogger
}

func New(uc ports.IUseCase, logger log.ILogger) IHandlers {
	return &Handlers{uc: uc, log: logger.WithModule("notification_backfill.handler")}
}

func (h *Handlers) resolveScope(c *gin.Context) (*uint, error) {
	jwtBiz := c.GetUint("business_id")
	if jwtBiz > 0 {
		id := jwtBiz
		return &id, nil
	}

	if param := c.Query("business_id"); param != "" {
		id, err := strconv.ParseUint(param, 10, 64)
		if err != nil || id == 0 {
			return nil, http.ErrMissingFile
		}
		bizID := uint(id)
		return &bizID, nil
	}

	return nil, nil
}
