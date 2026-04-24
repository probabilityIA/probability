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
)

func (h *handlers) CreatePutawayRule(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.CreatePutawayRuleBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isActive := true
	if body.IsActive != nil {
		isActive = *body.IsActive
	}
	rule, err := h.uc.CreatePutawayRule(c.Request.Context(), apprequest.CreatePutawayRuleDTO{
		BusinessID:   businessID,
		ProductID:    body.ProductID,
		CategoryID:   body.CategoryID,
		TargetZoneID: body.TargetZoneID,
		Priority:     body.Priority,
		Strategy:     body.Strategy,
		IsActive:     isActive,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

func (h *handlers) ListPutawayRules(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	activeOnly := c.Query("active_only") == "true"

	rules, total, err := h.uc.ListPutawayRules(c.Request.Context(), dtos.ListPutawayRulesParams{
		BusinessID: businessID,
		ActiveOnly: activeOnly,
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
	c.JSON(http.StatusOK, gin.H{
		"data":        rules,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func (h *handlers) UpdatePutawayRule(c *gin.Context) {
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
	var body request.UpdatePutawayRuleBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule, err := h.uc.UpdatePutawayRule(c.Request.Context(), apprequest.UpdatePutawayRuleDTO{
		ID:           uint(id),
		BusinessID:   businessID,
		ProductID:    body.ProductID,
		CategoryID:   body.CategoryID,
		TargetZoneID: body.TargetZoneID,
		Priority:     body.Priority,
		Strategy:     body.Strategy,
		IsActive:     body.IsActive,
	})
	if err != nil {
		if errors.Is(err, domainerrors.ErrPutawayRuleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (h *handlers) DeletePutawayRule(c *gin.Context) {
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
	if err := h.uc.DeletePutawayRule(c.Request.Context(), businessID, uint(id)); err != nil {
		if errors.Is(err, domainerrors.ErrPutawayRuleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *handlers) SuggestPutaway(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.SuggestPutawayBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items := make([]apprequest.PutawaySuggestItem, 0, len(body.Items))
	for _, i := range body.Items {
		items = append(items, apprequest.PutawaySuggestItem{ProductID: i.ProductID, Quantity: i.Quantity})
	}
	result, err := h.uc.SuggestPutaway(c.Request.Context(), apprequest.PutawaySuggestDTO{
		BusinessID: businessID,
		Items:      items,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *handlers) ConfirmPutaway(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid suggestion id"})
		return
	}
	var body request.ConfirmPutawayBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("user_id")
	var user *uint
	if userID > 0 {
		user = &userID
	}
	sug, err := h.uc.ConfirmPutaway(c.Request.Context(), apprequest.ConfirmPutawayDTO{
		BusinessID:       businessID,
		SuggestionID:     uint(id),
		ActualLocationID: body.ActualLocationID,
		UserID:           user,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrPutawaySuggestionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrPutawayAlreadyConfirmed):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, sug)
}

func (h *handlers) ListPutawaySuggestions(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	sugs, total, err := h.uc.ListPutawaySuggestions(c.Request.Context(), dtos.ListPutawaySuggestionsParams{
		BusinessID: businessID,
		Status:     c.Query("status"),
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
	c.JSON(http.StatusOK, gin.H{
		"data":        sugs,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func (h *handlers) CreateReplenishmentTask(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.CreateReplenishmentTaskBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := h.uc.CreateReplenishmentTask(c.Request.Context(), apprequest.CreateReplenishmentTaskDTO{
		BusinessID:     businessID,
		ProductID:      body.ProductID,
		WarehouseID:    body.WarehouseID,
		FromLocationID: body.FromLocationID,
		ToLocationID:   body.ToLocationID,
		Quantity:       body.Quantity,
		TriggeredBy:    body.TriggeredBy,
		Notes:          body.Notes,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *handlers) ListReplenishmentTasks(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	tasks, total, err := h.uc.ListReplenishmentTasks(c.Request.Context(), dtos.ListReplenishmentTasksParams{
		BusinessID:  businessID,
		WarehouseID: parseOptionalUintQuery(c, "warehouse_id"),
		Status:      c.Query("status"),
		AssignedTo:  parseOptionalUintQuery(c, "assigned_to"),
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(http.StatusOK, gin.H{
		"data":        tasks,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func (h *handlers) AssignReplenishment(c *gin.Context) {
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
	var body request.AssignReplenishmentBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := h.uc.AssignReplenishment(c.Request.Context(), apprequest.AssignReplenishmentDTO{
		BusinessID: businessID,
		TaskID:     uint(id),
		UserID:     body.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *handlers) CompleteReplenishment(c *gin.Context) {
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
	var body request.CompleteReplenishmentBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := h.uc.CompleteReplenishment(c.Request.Context(), apprequest.CompleteReplenishmentDTO{
		BusinessID: businessID,
		TaskID:     uint(id),
		Notes:      body.Notes,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *handlers) CancelReplenishment(c *gin.Context) {
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
	var body request.CancelReplenishmentBody
	_ = c.ShouldBindJSON(&body)
	task, err := h.uc.CancelReplenishment(c.Request.Context(), businessID, uint(id), body.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *handlers) DetectReplenishment(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	result, err := h.uc.DetectReplenishmentNeeds(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *handlers) CreateCrossDockLink(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.CreateCrossDockLinkBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	link, err := h.uc.CreateCrossDockLink(c.Request.Context(), apprequest.CreateCrossDockLinkDTO{
		BusinessID:        businessID,
		InboundShipmentID: body.InboundShipmentID,
		OutboundOrderID:   body.OutboundOrderID,
		ProductID:         body.ProductID,
		Quantity:          body.Quantity,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, link)
}

func (h *handlers) ListCrossDockLinks(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	links, total, err := h.uc.ListCrossDockLinks(c.Request.Context(), dtos.ListCrossDockLinksParams{
		BusinessID:      businessID,
		OutboundOrderID: c.Query("outbound_order_id"),
		Status:          c.Query("status"),
		Page:            page,
		PageSize:        pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(http.StatusOK, gin.H{
		"data":        links,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func (h *handlers) ExecuteCrossDock(c *gin.Context) {
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
	link, err := h.uc.ExecuteCrossDock(c.Request.Context(), apprequest.ExecuteCrossDockDTO{
		BusinessID: businessID,
		LinkID:     uint(id),
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrCrossDockNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrCrossDockExecuted):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, link)
}

func (h *handlers) RunSlotting(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.RunSlottingBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.uc.RunSlotting(c.Request.Context(), apprequest.RunSlottingDTO{
		BusinessID:  businessID,
		WarehouseID: body.WarehouseID,
		Period:      body.Period,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *handlers) ListVelocities(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	whID, _ := strconv.ParseUint(c.Query("warehouse_id"), 10, 64)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	velocities, err := h.uc.ListVelocities(c.Request.Context(), dtos.ListVelocityParams{
		BusinessID:  businessID,
		WarehouseID: uint(whID),
		Period:      c.Query("period"),
		Rank:        c.Query("rank"),
		Limit:       limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": velocities})
}
