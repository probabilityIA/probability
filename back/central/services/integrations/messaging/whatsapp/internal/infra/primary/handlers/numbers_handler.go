package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasenumbers"
)

type addNumberRequest struct {
	CountryCode  string `json:"country_code"`
	PhoneNumber  string `json:"phone_number"`
	VerifiedName string `json:"verified_name"`
}

type requestCodeRequest struct {
	Method string `json:"method"`
}

type verifyCodeRequest struct {
	Code string `json:"code"`
}

func (h *handler) numbersReady(c *gin.Context) bool {
	if h.numbersUseCase == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "el alta de números de WhatsApp no está disponible",
		})
		return false
	}
	return true
}

func (h *handler) GetNumberState(c *gin.Context) {
	if !h.numbersReady(c) {
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	state, err := h.numbersUseCase.GetState(c.Request.Context(), businessID)
	h.respondNumber(c, state, err, "consultando el estado del número")
}

func (h *handler) AddNumber(c *gin.Context) {
	if !h.numbersReady(c) {
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	var req addNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "payload inválido"})
		return
	}

	state, err := h.numbersUseCase.AddNumber(c.Request.Context(), businessID, usecasenumbers.AddNumberInput{
		CountryCode:  req.CountryCode,
		PhoneNumber:  req.PhoneNumber,
		VerifiedName: req.VerifiedName,
	})
	h.respondNumber(c, state, err, "agregando el número")
}

func (h *handler) RequestNumberCode(c *gin.Context) {
	if !h.numbersReady(c) {
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	var req requestCodeRequest
	_ = c.ShouldBindJSON(&req)

	state, err := h.numbersUseCase.RequestCode(c.Request.Context(), businessID, req.Method)
	h.respondNumber(c, state, err, "pidiendo el código")
}

func (h *handler) VerifyNumberCode(c *gin.Context) {
	if !h.numbersReady(c) {
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	var req verifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "payload inválido"})
		return
	}

	state, err := h.numbersUseCase.VerifyCode(c.Request.Context(), businessID, req.Code)
	h.respondNumber(c, state, err, "verificando el código")
}

func (h *handler) RegisterNumber(c *gin.Context) {
	if !h.numbersReady(c) {
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	state, err := h.numbersUseCase.Register(c.Request.Context(), businessID)
	h.respondNumber(c, state, err, "registrando el número")
}

func (h *handler) RemoveNumber(c *gin.Context) {
	if !h.numbersReady(c) {
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		return
	}

	state, err := h.numbersUseCase.RemoveNumber(c.Request.Context(), businessID)
	h.respondNumber(c, state, err, "quitando el número")
}

func (h *handler) respondNumber(c *gin.Context, state *usecasenumbers.NumberState, err error, accion string) {
	ctx := c.Request.Context()

	if err != nil {
		h.log.Warn(ctx).Err(err).Msgf("[WhatsApp Números] - error %s", accion)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": mapNumberState(state)})
}

func mapNumberState(state *usecasenumbers.NumberState) gin.H {
	if state == nil {
		return gin.H{}
	}

	return gin.H{
		"integration_id":           state.IntegrationID,
		"business_id":              state.BusinessID,
		"phone_number_id":          state.PhoneNumberID,
		"status":                   state.Status,
		"display_phone_number":     state.DisplayPhoneNumber,
		"verified_name":            state.VerifiedName,
		"name_status":              state.NameStatus,
		"code_verification_status": state.CodeVerification,
		"quality_rating":           state.QualityRating,
		"hosted_by_platform":       state.HostedByPlatform,
		"active":                   state.Active,
		"pin":                      state.Pin,
	}
}
