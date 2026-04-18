package request

type CreateStockMovementTypeDTO struct {
	Code        string
	Name        string
	Description string
	Direction   string
}

type UpdateStockMovementTypeDTO struct {
	ID          uint
	Name        string
	Description string
	IsActive    *bool
	Direction   string
}
