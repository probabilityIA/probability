package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
)

type TicketResponse struct {
	ID             uint       `json:"id"`
	Code           string     `json:"code"`
	BusinessID     *uint      `json:"business_id"`
	BusinessName   string     `json:"business_name,omitempty"`
	CreatedByID    uint       `json:"created_by_id"`
	CreatedByName  string     `json:"created_by_name,omitempty"`
	AssignedToID   *uint      `json:"assigned_to_id"`
	AssignedToName string     `json:"assigned_to_name,omitempty"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Type           string     `json:"type"`
	Category       string     `json:"category,omitempty"`
	Priority       string     `json:"priority"`
	Status         string     `json:"status"`
	Source         string     `json:"source"`
	Severity       string     `json:"severity,omitempty"`
	EscalatedToDev bool       `json:"escalated_to_dev"`
	EscalatedAt    *time.Time `json:"escalated_at,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	ClosedAt       *time.Time `json:"closed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	CommentsCount    int64 `json:"comments_count"`
	AttachmentsCount int64 `json:"attachments_count"`
}

type CommentResponse struct {
	ID         uint                  `json:"id"`
	TicketID   uint                  `json:"ticket_id"`
	UserID     uint                  `json:"user_id"`
	UserName   string                `json:"user_name"`
	Body       string                `json:"body"`
	IsInternal bool                  `json:"is_internal"`
	CreatedAt  time.Time             `json:"created_at"`
	Attachments []AttachmentResponse `json:"attachments,omitempty"`
}

type AttachmentResponse struct {
	ID             uint      `json:"id"`
	TicketID       uint      `json:"ticket_id"`
	CommentID      *uint     `json:"comment_id,omitempty"`
	UploadedByID   uint      `json:"uploaded_by_id"`
	UploadedByName string    `json:"uploaded_by_name"`
	FileURL        string    `json:"file_url"`
	FileName       string    `json:"file_name"`
	MimeType       string    `json:"mime_type"`
	Size           int64     `json:"size"`
	CreatedAt      time.Time `json:"created_at"`
}

type HistoryResponse struct {
	ID            uint      `json:"id"`
	TicketID      uint      `json:"ticket_id"`
	FromStatus    string    `json:"from_status"`
	ToStatus      string    `json:"to_status"`
	ChangedByID   uint      `json:"changed_by_id"`
	ChangedByName string    `json:"changed_by_name"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}

func FromTicket(t *entities.Ticket) TicketResponse {
	return TicketResponse{
		ID:               t.ID,
		Code:             t.Code,
		BusinessID:       t.BusinessID,
		BusinessName:     t.BusinessName,
		CreatedByID:      t.CreatedByID,
		CreatedByName:    t.CreatedByName,
		AssignedToID:     t.AssignedToID,
		AssignedToName:   t.AssignedToName,
		Title:            t.Title,
		Description:      t.Description,
		Type:             t.Type,
		Category:         t.Category,
		Priority:         t.Priority,
		Status:           t.Status,
		Source:           t.Source,
		Severity:         t.Severity,
		EscalatedToDev:   t.EscalatedToDev,
		EscalatedAt:      t.EscalatedAt,
		DueDate:          t.DueDate,
		ResolvedAt:       t.ResolvedAt,
		ClosedAt:         t.ClosedAt,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
		CommentsCount:    t.CommentsCount,
		AttachmentsCount: t.AttachmentsCount,
	}
}

func FromComment(c *entities.TicketComment) CommentResponse {
	resp := CommentResponse{
		ID:         c.ID,
		TicketID:   c.TicketID,
		UserID:     c.UserID,
		UserName:   c.UserName,
		Body:       c.Body,
		IsInternal: c.IsInternal,
		CreatedAt:  c.CreatedAt,
	}
	for i := range c.Attachments {
		resp.Attachments = append(resp.Attachments, FromAttachment(&c.Attachments[i]))
	}
	return resp
}

func FromAttachment(a *entities.TicketAttachment) AttachmentResponse {
	return AttachmentResponse{
		ID:             a.ID,
		TicketID:       a.TicketID,
		CommentID:      a.CommentID,
		UploadedByID:   a.UploadedByID,
		UploadedByName: a.UploadedByName,
		FileURL:        a.FileURL,
		FileName:       a.FileName,
		MimeType:       a.MimeType,
		Size:           a.Size,
		CreatedAt:      a.CreatedAt,
	}
}

func FromHistory(h *entities.TicketStatusHistory) HistoryResponse {
	return HistoryResponse{
		ID:            h.ID,
		TicketID:      h.TicketID,
		FromStatus:    h.FromStatus,
		ToStatus:      h.ToStatus,
		ChangedByID:   h.ChangedByID,
		ChangedByName: h.ChangedByName,
		Note:          h.Note,
		CreatedAt:     h.CreatedAt,
	}
}
