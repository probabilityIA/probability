package entities

import "time"

type Ticket struct {
	ID          uint
	Code        string
	BusinessID  *uint
	CreatedByID uint
	AssignedToID *uint
	Title       string
	Description string
	Type        string
	Category    string
	Priority    string
	Status      string
	Source      string
	Severity    string
	EscalatedToDev bool
	EscalatedAt    *time.Time
	DueDate     *time.Time
	ResolvedAt  *time.Time
	ClosedAt    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	BusinessName  string
	CreatedByName string
	AssignedToName string

	CommentsCount    int64
	AttachmentsCount int64

	Comments    []TicketComment
	Attachments []TicketAttachment
	History     []TicketStatusHistory
}

type TicketComment struct {
	ID         uint
	TicketID   uint
	UserID     uint
	UserName   string
	Body       string
	IsInternal bool
	CreatedAt  time.Time

	Attachments []TicketAttachment
}

type TicketAttachment struct {
	ID           uint
	TicketID     uint
	CommentID    *uint
	UploadedByID uint
	UploadedByName string
	FileURL      string
	FileName     string
	MimeType     string
	Size         int64
	CreatedAt    time.Time
}

type TicketStatusHistory struct {
	ID            uint
	TicketID      uint
	FromStatus    string
	ToStatus      string
	ChangedByID   uint
	ChangedByName string
	Note          string
	CreatedAt     time.Time
}
