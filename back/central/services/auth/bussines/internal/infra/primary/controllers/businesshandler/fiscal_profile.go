package businesshandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

type fiscalProfileRequest struct {
	DocumentType      string  `json:"document_type"`
	DocumentNumber    string  `json:"document_number"`
	DV                string  `json:"dv"`
	PersonType        string  `json:"person_type"`
	TaxRegime         string  `json:"tax_regime"`
	MunicipalityID    string  `json:"municipality_id"`
	Address           string  `json:"address"`
	Phone             string  `json:"phone"`
	BillingEmail      string  `json:"billing_email"`
	AppliesIVA        bool    `json:"applies_iva"`
	IVARate           float64 `json:"iva_rate"`
	AppliesReteFuente bool    `json:"applies_retefuente"`
	ReteFuenteRate    float64 `json:"retefuente_rate"`
	AppliesReteICA    bool    `json:"applies_reteica"`
	ReteICARate       float64 `json:"reteica_rate"`
	Notes             string  `json:"notes"`
}

func (h *BusinessHandler) fiscalProfileAccess(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id de negocio invalido"})
		return 0, false
	}
	if !middleware.IsSuperAdmin(c) {
		tokenBusinessID, ok := middleware.GetBusinessID(c)
		if !ok || tokenBusinessID == 0 || tokenBusinessID != uint(id) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "no tienes permisos sobre este negocio"})
			return 0, false
		}
	}
	return uint(id), true
}

func (h *BusinessHandler) GetFiscalProfileHandler(c *gin.Context) {
	businessID, ok := h.fiscalProfileAccess(c)
	if !ok {
		return
	}
	profile, err := h.usecase.GetFiscalProfile(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}

func (h *BusinessHandler) UpsertFiscalProfileHandler(c *gin.Context) {
	businessID, ok := h.fiscalProfileAccess(c)
	if !ok {
		return
	}
	var req fiscalProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	profile, err := h.usecase.UpsertFiscalProfile(c.Request.Context(), &domain.FiscalProfile{
		BusinessID:        businessID,
		DocumentType:      req.DocumentType,
		DocumentNumber:    req.DocumentNumber,
		DV:                req.DV,
		PersonType:        req.PersonType,
		TaxRegime:         req.TaxRegime,
		MunicipalityID:    req.MunicipalityID,
		Address:           req.Address,
		Phone:             req.Phone,
		BillingEmail:      req.BillingEmail,
		AppliesIVA:        req.AppliesIVA,
		IVARate:           req.IVARate,
		AppliesReteFuente: req.AppliesReteFuente,
		ReteFuenteRate:    req.ReteFuenteRate,
		AppliesReteICA:    req.AppliesReteICA,
		ReteICARate:       req.ReteICARate,
		Notes:             req.Notes,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}
