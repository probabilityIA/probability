package request

// ChangeStatus representa la petición HTTP para cambiar el estado de una orden
type ChangeStatus struct {
	Status   string                 `json:"status" binding:"required,max=64"`
	Metadata map[string]interface{} `json:"metadata"`
}
