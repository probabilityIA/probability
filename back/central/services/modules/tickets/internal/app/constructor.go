package app

import (
	"context"
	"mime/multipart"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	Create(ctx context.Context, dto dtos.CreateTicketDTO) (*entities.Ticket, error)
	Get(ctx context.Context, id uint, requesterUserID uint, requesterBusinessID *uint, isSuperAdmin bool) (*entities.Ticket, error)
	List(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error)
	Update(ctx context.Context, dto dtos.UpdateTicketDTO) (*entities.Ticket, error)
	Delete(ctx context.Context, id uint) error
	ChangeStatus(ctx context.Context, dto dtos.ChangeStatusDTO) (*entities.Ticket, error)
	ChangeArea(ctx context.Context, dto dtos.ChangeAreaDTO) (*entities.Ticket, error)
	Assign(ctx context.Context, dto dtos.AssignTicketDTO) (*entities.Ticket, error)
	Escalate(ctx context.Context, dto dtos.EscalateTicketDTO) (*entities.Ticket, error)

	AddComment(ctx context.Context, dto dtos.CreateCommentDTO) (*entities.TicketComment, error)
	ListComments(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error)

	UploadAttachment(ctx context.Context, ticketID uint, commentID *uint, uploaderID uint, file *multipart.FileHeader) (*entities.TicketAttachment, error)
	ListAttachments(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID uint, requesterID uint, isSuperAdmin bool) error

	ListHistory(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error)
}

type UseCase struct {
	repo    ports.IRepository
	storage ports.IStorageService
	log     log.ILogger
}

func New(repo ports.IRepository, storage ports.IStorageService, logger log.ILogger) IUseCase {
	return &UseCase{repo: repo, storage: storage, log: logger}
}
