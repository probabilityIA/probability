package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) ports.IRepository {
	return &Repository{db: database}
}

func (r *Repository) NextCode(ctx context.Context) (string, error) {
	var maxID uint
	row := r.db.Conn(ctx).Unscoped().Model(&models.Ticket{}).Select("COALESCE(MAX(id), 0)").Row()
	if err := row.Scan(&maxID); err != nil {
		return "", err
	}
	return fmt.Sprintf("TKT-%06d", maxID+1), nil
}

func (r *Repository) UserExists(ctx context.Context, userID uint) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).Table("user").Where("id = ? AND deleted_at IS NULL", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) Create(ctx context.Context, t *entities.Ticket) (*entities.Ticket, error) {
	m := entityToModel(t)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	t.ID = m.ID
	t.CreatedAt = m.CreatedAt
	t.UpdatedAt = m.UpdatedAt
	return t, nil
}

func (r *Repository) GetByID(ctx context.Context, id uint) (*entities.Ticket, error) {
	var m models.Ticket
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("CreatedBy").
		Preload("AssignedTo").
		Where("id = ?", id).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dom.ErrTicketNotFound
		}
		return nil, err
	}
	out := modelToEntity(&m)

	var counts struct {
		Comments    int64
		Attachments int64
	}
	r.db.Conn(ctx).Model(&models.TicketComment{}).Where("ticket_id = ?", id).Count(&counts.Comments)
	r.db.Conn(ctx).Model(&models.TicketAttachment{}).Where("ticket_id = ?", id).Count(&counts.Attachments)
	out.CommentsCount = counts.Comments
	out.AttachmentsCount = counts.Attachments
	return out, nil
}

func (r *Repository) List(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error) {
	q := r.db.Conn(ctx).Model(&models.Ticket{}).
		Preload("Business").
		Preload("CreatedBy").
		Preload("AssignedTo")

	if !params.IsSuperAdmin {
		if params.BusinessID != nil {
			q = q.Where("business_id = ?", *params.BusinessID)
		} else {
			q = q.Where("1 = 0")
		}
	} else if params.BusinessID != nil {
		q = q.Where("business_id = ?", *params.BusinessID)
	}

	if params.CreatedByID != nil {
		q = q.Where("created_by_id = ?", *params.CreatedByID)
	}
	if params.AssignedToID != nil {
		q = q.Where("assigned_to_id = ?", *params.AssignedToID)
	}
	if params.OnlyMine && params.UserID > 0 {
		q = q.Where("created_by_id = ? OR assigned_to_id = ?", params.UserID, params.UserID)
	}
	if len(params.Status) > 0 {
		q = q.Where("status IN ?", params.Status)
	}
	if len(params.Priority) > 0 {
		q = q.Where("priority IN ?", params.Priority)
	}
	if len(params.Type) > 0 {
		q = q.Where("type IN ?", params.Type)
	}
	if len(params.Area) > 0 {
		q = q.Where("area IN ?", params.Area)
	}
	if params.Source != "" {
		q = q.Where("source = ?", params.Source)
	}
	if params.EscalatedOnly {
		q = q.Where("escalated_to_dev = ?", true)
	}
	if s := strings.TrimSpace(params.Search); s != "" {
		like := "%" + s + "%"
		q = q.Where("title ILIKE ? OR description ILIKE ? OR code ILIKE ? OR category ILIKE ?", like, like, like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSort := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"priority":   "priority",
		"status":     "status",
		"area":       "area",
		"code":       "code",
		"due_date":   "due_date",
	}
	sortCol, ok := allowedSort[strings.ToLower(strings.TrimSpace(params.SortBy))]
	if !ok {
		sortCol = "created_at"
	}
	sortDir := strings.ToLower(strings.TrimSpace(params.SortOrder))
	if sortDir != "asc" {
		sortDir = "desc"
	}

	var ms []models.Ticket
	offset := (params.Page - 1) * params.PageSize
	if err := q.Order(sortCol + " " + sortDir).Limit(params.PageSize).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	out := make([]entities.Ticket, 0, len(ms))
	for i := range ms {
		out = append(out, *modelToEntity(&ms[i]))
	}
	return out, total, nil
}

func (r *Repository) Update(ctx context.Context, id uint, updates map[string]any) (*entities.Ticket, error) {
	res := r.db.Conn(ctx).Model(&models.Ticket{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		var m models.Ticket
		if err := r.db.Conn(ctx).Where("id = ?", id).First(&m).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, dom.ErrTicketNotFound
			}
			return nil, err
		}
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Unscoped().Delete(&models.Ticket{}, id).Error
}

func (r *Repository) AddComment(ctx context.Context, dto dtos.CreateCommentDTO) (*entities.TicketComment, error) {
	m := models.TicketComment{
		TicketID:   dto.TicketID,
		UserID:     dto.UserID,
		Body:       dto.Body,
		IsInternal: dto.IsInternal,
	}
	if err := r.db.Conn(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	var loaded models.TicketComment
	if err := r.db.Conn(ctx).Preload("User").First(&loaded, m.ID).Error; err != nil {
		return nil, err
	}
	return commentToEntity(&loaded), nil
}

func (r *Repository) ListComments(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error) {
	q := r.db.Conn(ctx).Preload("User").Preload("Attachments").Where("ticket_id = ?", ticketID)
	if !includeInternal {
		q = q.Where("is_internal = ?", false)
	}
	var ms []models.TicketComment
	if err := q.Order("created_at ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]entities.TicketComment, 0, len(ms))
	for i := range ms {
		out = append(out, *commentToEntity(&ms[i]))
	}
	return out, nil
}

func (r *Repository) AddAttachment(ctx context.Context, dto dtos.CreateAttachmentDTO) (*entities.TicketAttachment, error) {
	m := models.TicketAttachment{
		TicketID:     dto.TicketID,
		CommentID:    dto.CommentID,
		UploadedByID: dto.UploadedByID,
		FileURL:      dto.FileURL,
		FileName:     dto.FileName,
		MimeType:     dto.MimeType,
		Size:         dto.Size,
	}
	if err := r.db.Conn(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	var loaded models.TicketAttachment
	if err := r.db.Conn(ctx).Preload("UploadedBy").First(&loaded, m.ID).Error; err != nil {
		return nil, err
	}
	return attachmentToEntity(&loaded), nil
}

func (r *Repository) GetAttachment(ctx context.Context, id uint) (*entities.TicketAttachment, error) {
	var m models.TicketAttachment
	err := r.db.Conn(ctx).Preload("UploadedBy").Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dom.ErrAttachmentNotFound
		}
		return nil, err
	}
	return attachmentToEntity(&m), nil
}

func (r *Repository) DeleteAttachment(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Unscoped().Delete(&models.TicketAttachment{}, id).Error
}

func (r *Repository) ListAttachments(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error) {
	var ms []models.TicketAttachment
	err := r.db.Conn(ctx).Preload("UploadedBy").Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&ms).Error
	if err != nil {
		return nil, err
	}
	out := make([]entities.TicketAttachment, 0, len(ms))
	for i := range ms {
		out = append(out, *attachmentToEntity(&ms[i]))
	}
	return out, nil
}

func (r *Repository) AddHistory(ctx context.Context, ticketID uint, fromStatus, toStatus string, changedByID uint, note string) error {
	m := models.TicketStatusHistory{
		TicketID:    ticketID,
		ChangeType:  "status",
		FromStatus:  fromStatus,
		ToStatus:    toStatus,
		ChangedByID: changedByID,
		Note:        note,
	}
	return r.db.Conn(ctx).Create(&m).Error
}

func (r *Repository) AddAreaHistory(ctx context.Context, ticketID uint, fromArea, toArea string, changedByID uint, note string) error {
	m := models.TicketStatusHistory{
		TicketID:    ticketID,
		ChangeType:  "area",
		FromArea:    fromArea,
		ToArea:      toArea,
		ChangedByID: changedByID,
		Note:        note,
	}
	return r.db.Conn(ctx).Create(&m).Error
}

func (r *Repository) ListHistory(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error) {
	var ms []models.TicketStatusHistory
	err := r.db.Conn(ctx).Preload("ChangedBy").Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&ms).Error
	if err != nil {
		return nil, err
	}
	out := make([]entities.TicketStatusHistory, 0, len(ms))
	for i := range ms {
		out = append(out, *historyToEntity(&ms[i]))
	}
	return out, nil
}
