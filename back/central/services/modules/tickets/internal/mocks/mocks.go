package mocks

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type HistoryCall struct {
	TicketID    uint
	From        string
	To          string
	ChangedByID uint
	Note        string
}

type RepositoryMock struct {
	NextCodeFn       func(ctx context.Context) (string, error)
	CreateFn         func(ctx context.Context, ticket *entities.Ticket) (*entities.Ticket, error)
	GetByIDFn        func(ctx context.Context, id uint) (*entities.Ticket, error)
	ListFn           func(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error)
	UpdateFn         func(ctx context.Context, id uint, updates map[string]any) (*entities.Ticket, error)
	DeleteFn         func(ctx context.Context, id uint) error
	UserExistsFn     func(ctx context.Context, userID uint) (bool, error)
	FindSprintNameFn func(ctx context.Context, sprintID uint) (string, bool, error)

	AddCommentFn   func(ctx context.Context, dto dtos.CreateCommentDTO) (*entities.TicketComment, error)
	ListCommentsFn func(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error)

	AddAttachmentFn    func(ctx context.Context, dto dtos.CreateAttachmentDTO) (*entities.TicketAttachment, error)
	GetAttachmentFn    func(ctx context.Context, id uint) (*entities.TicketAttachment, error)
	DeleteAttachmentFn func(ctx context.Context, id uint) error
	ListAttachmentsFn  func(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error)

	AddHistoryFn     func(ctx context.Context, ticketID uint, fromStatus, toStatus string, changedByID uint, note string) error
	AddAreaHistoryFn func(ctx context.Context, ticketID uint, fromArea, toArea string, changedByID uint, note string) error
	ListHistoryFn    func(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error)

	CreatedTicket   *entities.Ticket
	Updates         []map[string]any
	History         []HistoryCall
	AreaHistory     []HistoryCall
	AddedComments   []dtos.CreateCommentDTO
	AddedAttachment *dtos.CreateAttachmentDTO
	DeletedTickets  []uint
}

var _ ports.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) NextCode(ctx context.Context) (string, error) {
	if m.NextCodeFn != nil {
		return m.NextCodeFn(ctx)
	}
	return "TKT-001", nil
}

func (m *RepositoryMock) Create(ctx context.Context, ticket *entities.Ticket) (*entities.Ticket, error) {
	m.CreatedTicket = ticket
	if m.CreateFn != nil {
		return m.CreateFn(ctx, ticket)
	}
	ticket.ID = 1
	return ticket, nil
}

func (m *RepositoryMock) GetByID(ctx context.Context, id uint) (*entities.Ticket, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return &entities.Ticket{ID: id, Status: "open", Area: "soporte"}, nil
}

func (m *RepositoryMock) List(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) Update(ctx context.Context, id uint, updates map[string]any) (*entities.Ticket, error) {
	m.Updates = append(m.Updates, updates)
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, updates)
	}
	return &entities.Ticket{ID: id}, nil
}

func (m *RepositoryMock) Delete(ctx context.Context, id uint) error {
	m.DeletedTickets = append(m.DeletedTickets, id)
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) UserExists(ctx context.Context, userID uint) (bool, error) {
	if m.UserExistsFn != nil {
		return m.UserExistsFn(ctx, userID)
	}
	return true, nil
}

func (m *RepositoryMock) FindSprintName(ctx context.Context, sprintID uint) (string, bool, error) {
	if m.FindSprintNameFn != nil {
		return m.FindSprintNameFn(ctx, sprintID)
	}
	return "Sprint mock", true, nil
}

func (m *RepositoryMock) AddComment(ctx context.Context, dto dtos.CreateCommentDTO) (*entities.TicketComment, error) {
	m.AddedComments = append(m.AddedComments, dto)
	if m.AddCommentFn != nil {
		return m.AddCommentFn(ctx, dto)
	}
	return &entities.TicketComment{ID: 1, Body: dto.Body}, nil
}

func (m *RepositoryMock) ListComments(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error) {
	if m.ListCommentsFn != nil {
		return m.ListCommentsFn(ctx, ticketID, includeInternal)
	}
	return nil, nil
}

func (m *RepositoryMock) AddAttachment(ctx context.Context, dto dtos.CreateAttachmentDTO) (*entities.TicketAttachment, error) {
	m.AddedAttachment = &dto
	if m.AddAttachmentFn != nil {
		return m.AddAttachmentFn(ctx, dto)
	}
	return &entities.TicketAttachment{ID: 1, FileURL: dto.FileURL}, nil
}

func (m *RepositoryMock) GetAttachment(ctx context.Context, id uint) (*entities.TicketAttachment, error) {
	if m.GetAttachmentFn != nil {
		return m.GetAttachmentFn(ctx, id)
	}
	return &entities.TicketAttachment{ID: id}, nil
}

func (m *RepositoryMock) DeleteAttachment(ctx context.Context, id uint) error {
	if m.DeleteAttachmentFn != nil {
		return m.DeleteAttachmentFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) ListAttachments(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error) {
	if m.ListAttachmentsFn != nil {
		return m.ListAttachmentsFn(ctx, ticketID)
	}
	return nil, nil
}

func (m *RepositoryMock) AddHistory(ctx context.Context, ticketID uint, fromStatus, toStatus string, changedByID uint, note string) error {
	m.History = append(m.History, HistoryCall{ticketID, fromStatus, toStatus, changedByID, note})
	if m.AddHistoryFn != nil {
		return m.AddHistoryFn(ctx, ticketID, fromStatus, toStatus, changedByID, note)
	}
	return nil
}

func (m *RepositoryMock) AddAreaHistory(ctx context.Context, ticketID uint, fromArea, toArea string, changedByID uint, note string) error {
	m.AreaHistory = append(m.AreaHistory, HistoryCall{ticketID, fromArea, toArea, changedByID, note})
	if m.AddAreaHistoryFn != nil {
		return m.AddAreaHistoryFn(ctx, ticketID, fromArea, toArea, changedByID, note)
	}
	return nil
}

func (m *RepositoryMock) ListHistory(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error) {
	if m.ListHistoryFn != nil {
		return m.ListHistoryFn(ctx, ticketID)
	}
	return nil, nil
}

type StorageServiceMock struct {
	UploadFileFn func(ctx context.Context, folder, filename string, data []byte, contentType string) (string, error)
	DeleteFileFn func(ctx context.Context, fileURL string) error

	UploadedFolders []string
	DeletedFiles    []string
}

var _ ports.IStorageService = (*StorageServiceMock)(nil)

func (m *StorageServiceMock) UploadFile(ctx context.Context, folder, filename string, data []byte, contentType string) (string, error) {
	m.UploadedFolders = append(m.UploadedFolders, folder)
	if m.UploadFileFn != nil {
		return m.UploadFileFn(ctx, folder, filename, data, contentType)
	}
	return "s3://" + folder + "/" + filename, nil
}

func (m *StorageServiceMock) DeleteFile(ctx context.Context, fileURL string) error {
	m.DeletedFiles = append(m.DeletedFiles, fileURL)
	if m.DeleteFileFn != nil {
		return m.DeleteFileFn(ctx, fileURL)
	}
	return nil
}

type SilentLogger struct{}

func NewSilentLogger() log.ILogger { return &SilentLogger{} }

func (l *SilentLogger) nop() zerolog.Logger { return zerolog.Nop() }

func (l *SilentLogger) Info(ctx ...context.Context) *zerolog.Event  { n := l.nop(); return n.Info() }
func (l *SilentLogger) Error(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Error() }
func (l *SilentLogger) Warn(ctx ...context.Context) *zerolog.Event  { n := l.nop(); return n.Warn() }
func (l *SilentLogger) Debug(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Debug() }
func (l *SilentLogger) Fatal(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Fatal() }
func (l *SilentLogger) Panic(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Panic() }
func (l *SilentLogger) With() zerolog.Context                       { n := l.nop(); return n.With() }
func (l *SilentLogger) WithService(service string) log.ILogger      { return l }
func (l *SilentLogger) WithModule(module string) log.ILogger        { return l }
func (l *SilentLogger) WithBusinessID(businessID uint) log.ILogger  { return l }
