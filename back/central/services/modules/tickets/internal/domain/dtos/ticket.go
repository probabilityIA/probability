package dtos

type ListTicketsParams struct {
	Page     int
	PageSize int

	BusinessID    *uint
	CreatedByID   *uint
	AssignedToID  *uint
	Status        []string
	Priority      []string
	Type          []string
	Area          []string
	Source        string
	EscalatedOnly bool
	Search        string

	SortBy    string // created_at | updated_at | priority | status | code
	SortOrder string // asc | desc

	OnlyMine     bool
	UserID       uint
	IsSuperAdmin bool
}

type CreateTicketDTO struct {
	BusinessID  *uint
	CreatedByID uint
	Title       string
	Description string
	Type        string
	Category    string
	Priority    string
	Severity    string
	Source      string
	Area        string
	AssignedToID *uint
	DueDate      *string
}

type UpdateTicketDTO struct {
	ID           uint
	Title        *string
	Description  *string
	Type         *string
	Category     *string
	Priority     *string
	Severity     *string
	Area         *string
	AssignedToID *uint
	DueDate      *string
	ClearDueDate bool
}

type ChangeAreaDTO struct {
	TicketID    uint
	NewArea     string
	Note        string
	ChangedByID uint
}

type ChangeStatusDTO struct {
	TicketID    uint
	NewStatus   string
	Note        string
	ChangedByID uint
}

type AssignTicketDTO struct {
	TicketID     uint
	AssignedToID *uint
	ChangedByID  uint
}

type EscalateTicketDTO struct {
	TicketID    uint
	Note        string
	ChangedByID uint
}

type CreateCommentDTO struct {
	TicketID   uint
	UserID     uint
	Body       string
	IsInternal bool
}

type CreateAttachmentDTO struct {
	TicketID     uint
	CommentID    *uint
	UploadedByID uint
	FileURL      string
	FileName     string
	MimeType     string
	Size         int64
}
