package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apprequest "github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/response"
)

func (h *handlers) CreateLPN(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.CreateLPNBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lpn, err := h.uc.CreateLPN(c.Request.Context(), apprequest.CreateLPNDTO{
		BusinessID: businessID,
		Code:       body.Code,
		LpnType:    body.LpnType,
		LocationID: body.LocationID,
	})
	if err != nil {
		if errors.Is(err, domainerrors.ErrDuplicateLPNCode) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.LicensePlateFromEntity(lpn))
}

func (h *handlers) GetLPN(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	lpn, err := h.uc.GetLPN(c.Request.Context(), businessID, uint(id))
	if err != nil {
		if errors.Is(err, domainerrors.ErrLPNNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.LicensePlateFromEntity(lpn))
}

func (h *handlers) ListLPNs(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	lpns, total, err := h.uc.ListLPNs(c.Request.Context(), dtos.ListLPNParams{
		BusinessID: businessID,
		LpnType:    c.Query("lpn_type"),
		Status:     c.Query("status"),
		LocationID: parseOptionalUintQuery(c, "location_id"),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	data := make([]response.LicensePlateResponse, len(lpns))
	for i := range lpns {
		data[i] = response.LicensePlateFromEntity(&lpns[i])
	}
	c.JSON(http.StatusOK, response.LicensePlateListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *handlers) UpdateLPN(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body request.UpdateLPNBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lpn, err := h.uc.UpdateLPN(c.Request.Context(), apprequest.UpdateLPNDTO{
		ID:         uint(id),
		BusinessID: businessID,
		Code:       body.Code,
		LpnType:    body.LpnType,
		LocationID: body.LocationID,
		Status:     body.Status,
	})
	if err != nil {
		if errors.Is(err, domainerrors.ErrLPNNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.LicensePlateFromEntity(lpn))
}

func (h *handlers) DeleteLPN(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.uc.DeleteLPN(c.Request.Context(), businessID, uint(id)); err != nil {
		if errors.Is(err, domainerrors.ErrLPNNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *handlers) AddToLPN(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body request.AddToLPNBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	line, err := h.uc.AddToLPN(c.Request.Context(), apprequest.AddToLPNDTO{
		BusinessID: businessID,
		LpnID:      uint(id),
		ProductID:  body.ProductID,
		LotID:      body.LotID,
		SerialID:   body.SerialID,
		Qty:        body.Qty,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.LicensePlateLineFromEntity(*line))
}

func (h *handlers) MoveLPN(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body request.MoveLPNBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lpn, err := h.uc.MoveLPN(c.Request.Context(), apprequest.MoveLPNDTO{
		BusinessID:    businessID,
		LpnID:         uint(id),
		NewLocationID: body.NewLocationID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.LicensePlateFromEntity(lpn))
}

func (h *handlers) DissolveLPN(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.uc.DissolveLPN(c.Request.Context(), apprequest.DissolveLPNDTO{
		BusinessID: businessID,
		LpnID:      uint(id),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "dissolved"})
}

func (h *handlers) MergeLPN(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body request.MergeLPNBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lpn, err := h.uc.MergeLPN(c.Request.Context(), apprequest.MergeLPNDTO{
		BusinessID:  businessID,
		SourceLpnID: uint(id),
		TargetLpnID: body.TargetLpnID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.LicensePlateFromEntity(lpn))
}

func (h *handlers) Scan(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.ScanBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("user_id")
	var uid *uint
	if userID > 0 {
		uid = &userID
	}
	result, err := h.uc.Scan(c.Request.Context(), apprequest.ScanDTO{
		BusinessID: businessID,
		Code:       body.Code,
		DeviceID:   body.DeviceID,
		UserID:     uid,
		Action:     body.Action,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.ScanResponse{
		Resolved:   result.Resolved,
		Resolution: response.ScanResolutionFromEntity(result.Resolution),
		Event:      response.ScanEventFromEntity(result.Event),
	})
}

func (h *handlers) InboundSync(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	integrationID, err := strconv.ParseUint(c.Param("integrationId"), 10, 64)
	if err != nil || integrationID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid integration id"})
		return
	}
	var body request.InboundSyncBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.uc.InboundSync(c.Request.Context(), apprequest.InboundSyncDTO{
		BusinessID:    businessID,
		IntegrationID: uint(integrationID),
		Payload:       body.Payload,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.InboundSyncResultResponse{
		Log:       response.SyncLogFromEntity(result.Log),
		Duplicate: result.Duplicate,
	})
}

func (h *handlers) ListSyncLogs(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, total, err := h.uc.ListSyncLogs(c.Request.Context(), dtos.ListSyncLogsParams{
		BusinessID:    businessID,
		IntegrationID: parseOptionalUintQuery(c, "integration_id"),
		Direction:     c.Query("direction"),
		Status:        c.Query("status"),
		Page:          page,
		PageSize:      pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	data := make([]response.InventorySyncLogResponse, len(logs))
	for i := range logs {
		if log := response.SyncLogFromEntity(&logs[i]); log != nil {
			data[i] = *log
		}
	}
	c.JSON(http.StatusOK, response.InventorySyncLogListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}
