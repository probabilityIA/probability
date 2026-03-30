package dtos

// ChangeStatusRequest es el DTO de dominio para cambio de estado de una orden
type ChangeStatusRequest struct {
	Status   string
	Metadata map[string]interface{}
	UserID   *uint
	UserName string
}
