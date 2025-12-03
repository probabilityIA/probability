package response

import "time"

// ResourceResponse representa la respuesta de un recurso
//
//	@Description	Información completa de un recurso
type ResourceResponse struct {
	ID               uint      `json:"id" example:"1" description:"ID único del recurso"`
	Name             string    `json:"name" example:"users" description:"Nombre del recurso"`
	Description      string    `json:"description" example:"Gestión de usuarios del sistema" description:"Descripción del recurso"`
	BusinessTypeID   uint      `json:"business_type_id" example:"11" description:"ID del tipo de business"`
	BusinessTypeName string    `json:"business_type_name" example:"Propiedad Horizontal" description:"Nombre del tipo de business"`
	CreatedAt        time.Time `json:"created_at" example:"2024-01-01T00:00:00Z" description:"Fecha de creación"`
	UpdatedAt        time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z" description:"Fecha de última actualización"`
}

// ResourceListResponse representa la respuesta de lista paginada de recursos
//
//	@Description	Lista paginada de recursos con metadatos de paginación
type ResourceListResponse struct {
	Resources  []ResourceResponse `json:"resources" description:"Lista de recursos"`
	Total      int64              `json:"total" example:"50" description:"Total de recursos disponibles"`
	Page       int                `json:"page" example:"1" description:"Página actual"`
	PageSize   int                `json:"page_size" example:"10" description:"Tamaño de página"`
	TotalPages int                `json:"total_pages" example:"5" description:"Total de páginas"`
}

// GetResourcesResponse representa la respuesta para obtener recursos
//
//	@Description	Respuesta exitosa al obtener lista de recursos
type GetResourcesResponse struct {
	Success bool                 `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string               `json:"message" example:"Recursos obtenidos exitosamente" description:"Mensaje descriptivo"`
	Data    ResourceListResponse `json:"data" description:"Datos de la lista de recursos"`
}

// GetResourceByIDResponse representa la respuesta para obtener un recurso por ID
//
//	@Description	Respuesta exitosa al obtener un recurso específico
type GetResourceByIDResponse struct {
	Success bool             `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string           `json:"message" example:"Recurso obtenido exitosamente" description:"Mensaje descriptivo"`
	Data    ResourceResponse `json:"data" description:"Datos del recurso"`
}

// CreateResourceResponse representa la respuesta para crear un recurso
//
//	@Description	Respuesta exitosa al crear un nuevo recurso
type CreateResourceResponse struct {
	Success bool             `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string           `json:"message" example:"Recurso creado exitosamente" description:"Mensaje descriptivo"`
	Data    ResourceResponse `json:"data" description:"Datos del recurso creado"`
}

// UpdateResourceResponse representa la respuesta para actualizar un recurso
//
//	@Description	Respuesta exitosa al actualizar un recurso
type UpdateResourceResponse struct {
	Success bool             `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string           `json:"message" example:"Recurso actualizado exitosamente" description:"Mensaje descriptivo"`
	Data    ResourceResponse `json:"data" description:"Datos del recurso actualizado"`
}

// DeleteResourceResponse representa la respuesta para eliminar un recurso
//
//	@Description	Respuesta exitosa al eliminar un recurso
type DeleteResourceResponse struct {
	Success bool   `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string `json:"message" example:"Recurso eliminado exitosamente" description:"Mensaje descriptivo"`
}

// ErrorResponse representa una respuesta de error
//
//	@Description	Respuesta de error estándar
type ErrorResponse struct {
	Success bool   `json:"success" example:"false" description:"Indica que la operación falló"`
	Message string `json:"message" example:"Error al procesar la solicitud" description:"Mensaje de error"`
	Error   string `json:"error,omitempty" example:"Detalles específicos del error" description:"Detalles del error"`
}
