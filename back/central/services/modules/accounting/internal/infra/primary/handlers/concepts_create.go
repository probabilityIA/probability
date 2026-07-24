package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func (h *Handlers) CreateConcept(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	var req request.CreateConceptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	isRealIncome := true
	if req.IsRealIncome != nil {
		isRealIncome = *req.IsRealIncome
	}
	concept, err := h.uc.CreateConcept(c.Request.Context(), dtos.CreateConceptDTO{
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		Kind:         req.Kind,
		IsRealIncome: isRealIncome,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrDuplicateCode):
			c.JSON(http.StatusConflict, gin.H{"success": false, "error": err.Error()})
		case errors.Is(err, domainerrors.ErrInvalidKind), errors.Is(err, domainerrors.ErrCodeRequired), errors.Is(err, domainerrors.ErrNameRequired):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": response.FromConcept(concept)})
}
