package request

type CreateTicketRequest struct {
	BusinessID   *uint   `json:"business_id"`
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Type         string  `json:"type"`
	Category     string  `json:"category"`
	Priority     string  `json:"priority"`
	Severity     string  `json:"severity"`
	Source       string  `json:"source"`
	AssignedToID *uint   `json:"assigned_to_id"`
	DueDate      *string `json:"due_date"`
}

type UpdateTicketRequest struct {
	Title        *string `json:"title"`
	Description  *string `json:"description"`
	Type         *string `json:"type"`
	Category     *string `json:"category"`
	Priority     *string `json:"priority"`
	Severity     *string `json:"severity"`
	AssignedToID *uint   `json:"assigned_to_id"`
	DueDate      *string `json:"due_date"`
	ClearDueDate bool    `json:"clear_due_date"`
}

type ChangeStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Note   string `json:"note"`
}

type AssignRequest struct {
	AssignedToID *uint `json:"assigned_to_id"`
}

type EscalateRequest struct {
	Note string `json:"note"`
}

type CreateCommentRequest struct {
	Body       string `json:"body" binding:"required"`
	IsInternal bool   `json:"is_internal"`
}
