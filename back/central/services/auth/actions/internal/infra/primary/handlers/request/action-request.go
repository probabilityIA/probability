package request

// CreateActionRequest representa la solicitud para crear un action
//
//	@Description	Solicitud para crear un nuevo action en el sistema
type CreateActionRequest struct {
	Name        string `json:"name" binding:"required" example:"create" description:"Nombre único del action"`
	Description string `json:"description" example:"Permite crear nuevos registros" description:"Descripción del action"`
}

// UpdateActionRequest representa la solicitud para actualizar un action
//
//	@Description	Solicitud para actualizar un action existente
type UpdateActionRequest struct {
	Name        string `json:"name" binding:"required" example:"create" description:"Nombre único del action"`
	Description string `json:"description" example:"Permite crear nuevos registros" description:"Descripción del action"`
}

// GetActionsRequest representa los parámetros de consulta para obtener actions
//
//	@Description	Parámetros para filtrar y paginar la lista de actions
type GetActionsRequest struct {
	Page     int    `form:"page" example:"1" description:"Número de página (por defecto: 1)"`
	PageSize int    `form:"page_size" example:"10" description:"Tamaño de página (por defecto: 10)"`
	Name     string `form:"name" example:"create" description:"Filtrar por nombre (búsqueda parcial)"`
}
