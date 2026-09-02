package request

type CreateSprintRequest struct {
	Name      string `json:"name" binding:"required"`
	Goal      string `json:"goal"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Status    string `json:"status"`
}

type UpdateSprintRequest struct {
	Name      *string `json:"name"`
	Goal      *string `json:"goal"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	Status    *string `json:"status"`
}

type ChangeSprintStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
