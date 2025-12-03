package request

// CreateResourceRequest representa la solicitud para crear un recurso
//
//	@Description	Solicitud para crear un nuevo recurso en el sistema
type CreateResourceRequest struct {
	Name           string `json:"name" binding:"required" example:"users" description:"Nombre único del recurso"`
	Description    string `json:"description" example:"Gestión de usuarios del sistema" description:"Descripción del recurso"`
	BusinessTypeID *uint  `json:"business_type_id" example:"11" description:"ID del tipo de business (opcional, nil = genérico)"`
}

// UpdateResourceRequest representa la solicitud para actualizar un recurso
//
//	@Description	Solicitud para actualizar un recurso existente
type UpdateResourceRequest struct {
	Name           string `json:"name" binding:"required" example:"users" description:"Nombre único del recurso"`
	Description    string `json:"description" example:"Gestión de usuarios del sistema" description:"Descripción del recurso"`
	BusinessTypeID *uint  `json:"business_type_id" example:"11" description:"ID del tipo de business (opcional, nil = genérico)"`
}

// GetResourcesRequest representa los parámetros de consulta para obtener recursos
//
//	@Description	Parámetros para filtrar y paginar la lista de recursos
type GetResourcesRequest struct {
	Page           int    `form:"page" example:"1" description:"Número de página (por defecto: 1)"`
	PageSize       int    `form:"page_size" example:"10" description:"Tamaño de página (por defecto: 10)"`
	Name           string `form:"name" example:"user" description:"Filtrar por nombre (búsqueda parcial)"`
	Description    string `form:"description" example:"gestión" description:"Filtrar por descripción (búsqueda parcial)"`
	BusinessTypeID uint   `form:"business_type_id" example:"11" description:"Filtrar por tipo de business (incluye genéricos)"`
	SortBy         string `form:"sort_by" example:"name" description:"Campo para ordenar: name, created_at, updated_at"`
	SortOrder      string `form:"sort_order" example:"asc" description:"Orden: asc o desc"`
}
