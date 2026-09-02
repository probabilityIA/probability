package dtos

type ListSprintsParams struct {
	Page     int
	PageSize int
	Status   string
}

type CreateSprintDTO struct {
	Name        string
	Goal        string
	StartDate   string
	EndDate     string
	Status      string
	CreatedByID uint
}

type UpdateSprintDTO struct {
	ID        uint
	Name      *string
	Goal      *string
	StartDate *string
	EndDate   *string
	Status    *string
}

type ChangeSprintStatusDTO struct {
	SprintID uint
	Status   string
}
