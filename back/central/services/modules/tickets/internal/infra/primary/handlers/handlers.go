package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/infra/primary/handlers/response"
)

func (h *Handlers) Create(c *gin.Context) {
	var req request.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, businessID, isSuperAdmin := h.requesterContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	source := req.Source
	dtoBusinessID := req.BusinessID
	if !isSuperAdmin {
		source = "business"
		dtoBusinessID = businessID
	} else if source == "" {
		source = "internal"
	}

	dto := dtos.CreateTicketDTO{
		BusinessID:   dtoBusinessID,
		CreatedByID:  userID,
		Title:        req.Title,
		Description:  req.Description,
		Type:         req.Type,
		Category:     req.Category,
		Priority:     req.Priority,
		Severity:     req.Severity,
		Source:       source,
		Area:         req.Area,
		AssignedToID: req.AssignedToID,
		DueDate:      req.DueDate,
	}
	if !isSuperAdmin {
		dto.AssignedToID = nil
	}

	t, err := h.uc.Create(c.Request.Context(), dto)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.FromTicket(t))
}

func (h *Handlers) Get(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	userID, businessID, isSuperAdmin := h.requesterContext(c)
	t, err := h.uc.Get(c.Request.Context(), id, userID, businessID, isSuperAdmin)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.FromTicket(t))
}

func (h *Handlers) List(c *gin.Context) {
	userID, businessID, isSuperAdmin := h.requesterContext(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := dtos.ListTicketsParams{
		Page:          page,
		PageSize:      pageSize,
		Status:        h.splitCSV(c.Query("status")),
		Priority:      h.splitCSV(c.Query("priority")),
		Type:          h.splitCSV(c.Query("type")),
		Area:          h.splitCSV(c.Query("area")),
		Source:        c.Query("source"),
		EscalatedOnly: c.Query("escalated") == "true",
		Search:        c.Query("search"),
		SortBy:        c.Query("sort_by"),
		SortOrder:     c.Query("sort_order"),
		OnlyMine:      c.Query("only_mine") == "true",
		UserID:        userID,
		IsSuperAdmin:  isSuperAdmin,
	}
	if isSuperAdmin {
		if v := c.Query("business_id"); v != "" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil && id > 0 {
				u := uint(id)
				params.BusinessID = &u
			}
		}
		if v := c.Query("created_by_id"); v != "" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				u := uint(id)
				params.CreatedByID = &u
			}
		}
		if v := c.Query("assigned_to_id"); v != "" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				u := uint(id)
				params.AssignedToID = &u
			}
		}
	} else {
		params.BusinessID = businessID
	}

	items, total, err := h.uc.List(c.Request.Context(), params)
	if err != nil {
		h.handleError(c, err)
		return
	}
	out := make([]response.TicketResponse, 0, len(items))
	for i := range items {
		out = append(out, response.FromTicket(&items[i]))
	}
	totalPages := total / int64(params.PageSize)
	if total%int64(params.PageSize) > 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, gin.H{
		"data":        out,
		"total":       total,
		"page":        params.Page,
		"page_size":   params.PageSize,
		"total_pages": totalPages,
	})
}

func (h *Handlers) Update(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var req request.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := h.uc.Update(c.Request.Context(), dtos.UpdateTicketDTO{
		ID:           id,
		Title:        req.Title,
		Description:  req.Description,
		Type:         req.Type,
		Category:     req.Category,
		Priority:     req.Priority,
		Severity:     req.Severity,
		Area:         req.Area,
		AssignedToID: req.AssignedToID,
		DueDate:      req.DueDate,
		ClearDueDate: req.ClearDueDate,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.FromTicket(t))
}

func (h *Handlers) Delete(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handlers) ChangeStatus(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	var req request.ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _, _ := h.requesterContext(c)
	t, err := h.uc.ChangeStatus(c.Request.Context(), dtos.ChangeStatusDTO{
		TicketID:    id,
		NewStatus:   req.Status,
		Note:        req.Note,
		ChangedByID: userID,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.FromTicket(t))
}

func (h *Handlers) ChangeArea(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	var req request.ChangeAreaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _, _ := h.requesterContext(c)
	t, err := h.uc.ChangeArea(c.Request.Context(), dtos.ChangeAreaDTO{
		TicketID:    id,
		NewArea:     req.Area,
		Note:        req.Note,
		ChangedByID: userID,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.FromTicket(t))
}

func (h *Handlers) Assign(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	var req request.AssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _, _ := h.requesterContext(c)
	t, err := h.uc.Assign(c.Request.Context(), dtos.AssignTicketDTO{
		TicketID:     id,
		AssignedToID: req.AssignedToID,
		ChangedByID:  userID,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.FromTicket(t))
}

func (h *Handlers) Escalate(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	var req request.EscalateRequest
	_ = c.ShouldBindJSON(&req)
	userID, _, _ := h.requesterContext(c)
	t, err := h.uc.Escalate(c.Request.Context(), dtos.EscalateTicketDTO{
		TicketID:    id,
		Note:        req.Note,
		ChangedByID: userID,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.FromTicket(t))
}

func (h *Handlers) ListComments(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	isSuperAdmin := middleware.IsSuperAdmin(c)
	items, err := h.uc.ListComments(c.Request.Context(), id, isSuperAdmin)
	if err != nil {
		h.handleError(c, err)
		return
	}
	out := make([]response.CommentResponse, 0, len(items))
	for i := range items {
		out = append(out, response.FromComment(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *Handlers) AddComment(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	var req request.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _, isSuperAdmin := h.requesterContext(c)
	internal := req.IsInternal && isSuperAdmin
	cmt, err := h.uc.AddComment(c.Request.Context(), dtos.CreateCommentDTO{
		TicketID:   id,
		UserID:     userID,
		Body:       req.Body,
		IsInternal: internal,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.FromComment(cmt))
}

func (h *Handlers) ListAttachments(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	items, err := h.uc.ListAttachments(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	out := make([]response.AttachmentResponse, 0, len(items))
	for i := range items {
		out = append(out, response.FromAttachment(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *Handlers) UploadAttachment(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	var commentID *uint
	if v := c.PostForm("comment_id"); v != "" {
		if cid, err := strconv.ParseUint(v, 10, 64); err == nil && cid > 0 {
			u := uint(cid)
			commentID = &u
		}
	}
	userID, _, _ := h.requesterContext(c)
	att, err := h.uc.UploadAttachment(c.Request.Context(), id, commentID, userID, file)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.FromAttachment(att))
}

func (h *Handlers) DeleteAttachment(c *gin.Context) {
	id, ok := h.parseUintParam(c, "attachment_id")
	if !ok {
		return
	}
	userID, _, isSuperAdmin := h.requesterContext(c)
	if err := h.uc.DeleteAttachment(c.Request.Context(), id, userID, isSuperAdmin); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handlers) ListHistory(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	items, err := h.uc.ListHistory(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	out := make([]response.HistoryResponse, 0, len(items))
	for i := range items {
		out = append(out, response.FromHistory(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}
