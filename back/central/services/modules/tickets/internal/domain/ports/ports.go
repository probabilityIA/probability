package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
)

type IRepository interface {
	NextCode(ctx context.Context) (string, error)
	Create(ctx context.Context, ticket *entities.Ticket) (*entities.Ticket, error)
	GetByID(ctx context.Context, id uint) (*entities.Ticket, error)
	List(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error)
	Update(ctx context.Context, id uint, updates map[string]any) (*entities.Ticket, error)
	Delete(ctx context.Context, id uint) error

	UserExists(ctx context.Context, userID uint) (bool, error)

	AddComment(ctx context.Context, dto dtos.CreateCommentDTO) (*entities.TicketComment, error)
	ListComments(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error)

	AddAttachment(ctx context.Context, dto dtos.CreateAttachmentDTO) (*entities.TicketAttachment, error)
	GetAttachment(ctx context.Context, id uint) (*entities.TicketAttachment, error)
	DeleteAttachment(ctx context.Context, id uint) error
	ListAttachments(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error)

	AddHistory(ctx context.Context, ticketID uint, fromStatus, toStatus string, changedByID uint, note string) error
	ListHistory(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error)
}

type IStorageService interface {
	UploadFile(ctx context.Context, folder, filename string, data []byte, contentType string) (string, error)
	DeleteFile(ctx context.Context, fileURL string) error
}
