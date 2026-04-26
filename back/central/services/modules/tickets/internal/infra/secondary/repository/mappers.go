package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func entityToModel(t *entities.Ticket) *models.Ticket {
	return &models.Ticket{
		Code:           t.Code,
		BusinessID:     t.BusinessID,
		CreatedByID:    t.CreatedByID,
		AssignedToID:   t.AssignedToID,
		Title:          t.Title,
		Description:    t.Description,
		Type:           t.Type,
		Category:       t.Category,
		Priority:       t.Priority,
		Status:         t.Status,
		Source:         t.Source,
		Severity:       t.Severity,
		EscalatedToDev: t.EscalatedToDev,
		EscalatedAt:    t.EscalatedAt,
		DueDate:        t.DueDate,
		ResolvedAt:     t.ResolvedAt,
		ClosedAt:       t.ClosedAt,
	}
}

func modelToEntity(m *models.Ticket) *entities.Ticket {
	out := &entities.Ticket{
		ID:             m.ID,
		Code:           m.Code,
		BusinessID:     m.BusinessID,
		CreatedByID:    m.CreatedByID,
		AssignedToID:   m.AssignedToID,
		Title:          m.Title,
		Description:    m.Description,
		Type:           m.Type,
		Category:       m.Category,
		Priority:       m.Priority,
		Status:         m.Status,
		Source:         m.Source,
		Severity:       m.Severity,
		EscalatedToDev: m.EscalatedToDev,
		EscalatedAt:    m.EscalatedAt,
		DueDate:        m.DueDate,
		ResolvedAt:     m.ResolvedAt,
		ClosedAt:       m.ClosedAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	if m.Business != nil {
		out.BusinessName = m.Business.Name
	}
	out.CreatedByName = m.CreatedBy.Name
	if m.AssignedTo != nil {
		out.AssignedToName = m.AssignedTo.Name
	}
	return out
}

func commentToEntity(m *models.TicketComment) *entities.TicketComment {
	c := &entities.TicketComment{
		ID:         m.ID,
		TicketID:   m.TicketID,
		UserID:     m.UserID,
		UserName:   m.User.Name,
		Body:       m.Body,
		IsInternal: m.IsInternal,
		CreatedAt:  m.CreatedAt,
	}
	for i := range m.Attachments {
		c.Attachments = append(c.Attachments, *attachmentToEntity(&m.Attachments[i]))
	}
	return c
}

func attachmentToEntity(m *models.TicketAttachment) *entities.TicketAttachment {
	return &entities.TicketAttachment{
		ID:             m.ID,
		TicketID:       m.TicketID,
		CommentID:      m.CommentID,
		UploadedByID:   m.UploadedByID,
		UploadedByName: m.UploadedBy.Name,
		FileURL:        m.FileURL,
		FileName:       m.FileName,
		MimeType:       m.MimeType,
		Size:           m.Size,
		CreatedAt:      m.CreatedAt,
	}
}

func historyToEntity(m *models.TicketStatusHistory) *entities.TicketStatusHistory {
	return &entities.TicketStatusHistory{
		ID:            m.ID,
		TicketID:      m.TicketID,
		FromStatus:    m.FromStatus,
		ToStatus:      m.ToStatus,
		ChangedByID:   m.ChangedByID,
		ChangedByName: m.ChangedBy.Name,
		Note:          m.Note,
		CreatedAt:     m.CreatedAt,
	}
}
