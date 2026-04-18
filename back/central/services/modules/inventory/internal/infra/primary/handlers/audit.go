package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	apprequest "github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/request"
)

func (h *handlers) CreateCountPlan(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.CreateCountPlanBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isActive := true
	if body.IsActive != nil {
		isActive = *body.IsActive
	}
	plan, err := h.uc.CreateCountPlan(c.Request.Context(), apprequest.CreateCountPlanDTO{
		BusinessID:    businessID,
		WarehouseID:   body.WarehouseID,
		Name:          body.Name,
		Strategy:      body.Strategy,
		FrequencyDays: body.FrequencyDays,
		NextRunAt:     body.NextRunAt,
		IsActive:      isActive,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, plan)
}

func (h *handlers) ListCountPlans(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	plans, total, err := h.uc.ListCountPlans(c.Request.Context(), dtos.ListCycleCountPlansParams{
		BusinessID:  businessID,
		WarehouseID: parseOptionalUintQuery(c, "warehouse_id"),
		ActiveOnly:  c.Query("active_only") == "true",
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
		"data":        plans,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func (h *handlers) UpdateCountPlan(c *gin.Context) {
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
	var body request.UpdateCountPlanBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plan, err := h.uc.UpdateCountPlan(c.Request.Context(), apprequest.UpdateCountPlanDTO{
		ID:            uint(id),
		BusinessID:    businessID,
		WarehouseID:   body.WarehouseID,
		Name:          body.Name,
		Strategy:      body.Strategy,
		FrequencyDays: body.FrequencyDays,
		NextRunAt:     body.NextRunAt,
		IsActive:      body.IsActive,
	})
	if err != nil {
		if errors.Is(err, domainerrors.ErrCountPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plan)
}

func (h *handlers) DeleteCountPlan(c *gin.Context) {
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
	if err := h.uc.DeleteCountPlan(c.Request.Context(), businessID, uint(id)); err != nil {
		if errors.Is(err, domainerrors.ErrCountPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *handlers) GenerateCountTask(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	var body request.GenerateCountTaskBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.uc.GenerateCountTask(c.Request.Context(), apprequest.GenerateCountTaskDTO{
		BusinessID: businessID,
		PlanID:     body.PlanID,
		ScopeType:  body.ScopeType,
		ScopeID:    body.ScopeID,
	})
	if err != nil {
		if errors.Is(err, domainerrors.ErrCountPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *handlers) ListCountTasks(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	tasks, total, err := h.uc.ListCountTasks(c.Request.Context(), dtos.ListCycleCountTasksParams{
		BusinessID:  businessID,
		WarehouseID: parseOptionalUintQuery(c, "warehouse_id"),
		PlanID:      parseOptionalUintQuery(c, "plan_id"),
		Status:      c.Query("status"),
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

func (h *handlers) StartCountTask(c *gin.Context) {
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
	var body request.StartCountTaskBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := h.uc.StartCountTask(c.Request.Context(), apprequest.StartCountTaskDTO{
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

func (h *handlers) FinishCountTask(c *gin.Context) {
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
	task, err := h.uc.FinishCountTask(c.Request.Context(), businessID, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *handlers) ListCountLines(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil || taskID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	lines, total, err := h.uc.ListCountLines(c.Request.Context(), dtos.ListCycleCountLinesParams{
		BusinessID: businessID,
		TaskID:     uint(taskID),
		Status:     c.Query("status"),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(http.StatusOK, gin.H{
		"data":        lines,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func (h *handlers) SubmitCountLine(c *gin.Context) {
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
	var body request.SubmitCountLineBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetUint("user_id")
	var uid *uint
	if userID > 0 {
		uid = &userID
	}
	result, err := h.uc.SubmitCountLine(c.Request.Context(), apprequest.SubmitCountLineDTO{
		BusinessID: businessID,
		LineID:     uint(id),
		CountedQty: body.CountedQty,
		UserID:     uid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *handlers) ListDiscrepancies(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	list, total, err := h.uc.ListDiscrepancies(c.Request.Context(), dtos.ListDiscrepanciesParams{
		BusinessID: businessID,
		TaskID:     parseOptionalUintQuery(c, "task_id"),
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
		"data":        list,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func (h *handlers) ApproveDiscrepancy(c *gin.Context) {
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
	var body request.ApproveDiscrepancyBody
	_ = c.ShouldBindJSON(&body)
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context required"})
		return
	}
	disc, err := h.uc.ApproveDiscrepancy(c.Request.Context(), apprequest.ApproveDiscrepancyDTO{
		BusinessID:    businessID,
		DiscrepancyID: uint(id),
		ReviewerID:    userID,
		Notes:         body.Notes,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, disc)
}

func (h *handlers) RejectDiscrepancy(c *gin.Context) {
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
	var body request.RejectDiscrepancyBody
	_ = c.ShouldBindJSON(&body)
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context required"})
		return
	}
	disc, err := h.uc.RejectDiscrepancy(c.Request.Context(), apprequest.RejectDiscrepancyDTO{
		BusinessID:    businessID,
		DiscrepancyID: uint(id),
		ReviewerID:    userID,
		Reason:        body.Reason,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, disc)
}

func (h *handlers) ExportKardex(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}
	productID := c.Query("product_id")
	whID, _ := strconv.ParseUint(c.Query("warehouse_id"), 10, 64)
	if productID == "" || whID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id and warehouse_id are required"})
		return
	}
	var from, to *time.Time
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			from = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			to = &t
		}
	}
	result, err := h.uc.ExportKardex(c.Request.Context(), apprequest.KardexExportDTO{
		BusinessID:  businessID,
		ProductID:   productID,
		WarehouseID: uint(whID),
		From:        from,
		To:          to,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
