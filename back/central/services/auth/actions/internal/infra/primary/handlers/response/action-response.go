package response

import "time"

// ActionResponse representa la respuesta de un action
//
//	@Description	Información completa de un action
type ActionResponse struct {
	ID          uint      `json:"id" example:"1" description:"ID único del action"`
	Name        string    `json:"name" example:"create" description:"Nombre del action"`
	Description string    `json:"description" example:"Permite crear nuevos registros" description:"Descripción del action"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z" description:"Fecha de creación"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z" description:"Fecha de última actualización"`
}

// ActionListResponse representa la respuesta de lista paginada de actions
//
//	@Description	Lista paginada de actions con metadatos de paginación
type ActionListResponse struct {
	Actions    []ActionResponse `json:"actions" description:"Lista de actions"`
	Total      int64            `json:"total" example:"50" description:"Total de actions disponibles"`
	Page       int              `json:"page" example:"1" description:"Página actual"`
	PageSize   int              `json:"page_size" example:"10" description:"Tamaño de página"`
	TotalPages int              `json:"total_pages" example:"5" description:"Total de páginas"`
}

// GetActionsResponse representa la respuesta para obtener actions
//
//	@Description	Respuesta exitosa al obtener lista de actions
type GetActionsResponse struct {
	Success bool               `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string             `json:"message" example:"Actions obtenidos exitosamente" description:"Mensaje descriptivo"`
	Data    ActionListResponse `json:"data" description:"Datos de la lista de actions"`
}

// GetActionByIDResponse representa la respuesta para obtener un action por ID
//
//	@Description	Respuesta exitosa al obtener un action específico
type GetActionByIDResponse struct {
	Success bool           `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string         `json:"message" example:"Action obtenido exitosamente" description:"Mensaje descriptivo"`
	Data    ActionResponse `json:"data" description:"Datos del action"`
}

// CreateActionResponse representa la respuesta para crear un action
//
//	@Description	Respuesta exitosa al crear un nuevo action
type CreateActionResponse struct {
	Success bool           `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string         `json:"message" example:"Action creado exitosamente" description:"Mensaje descriptivo"`
	Data    ActionResponse `json:"data" description:"Datos del action creado"`
}

// UpdateActionResponse representa la respuesta para actualizar un action
//
//	@Description	Respuesta exitosa al actualizar un action
type UpdateActionResponse struct {
	Success bool           `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string         `json:"message" example:"Action actualizado exitosamente" description:"Mensaje descriptivo"`
	Data    ActionResponse `json:"data" description:"Datos del action actualizado"`
}

// DeleteActionResponse representa la respuesta para eliminar un action
//
//	@Description	Respuesta exitosa al eliminar un action
type DeleteActionResponse struct {
	Success bool   `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string `json:"message" example:"Action eliminado exitosamente" description:"Mensaje descriptivo"`
}

// ErrorResponse representa una respuesta de error
//
//	@Description	Respuesta de error estándar
type ErrorResponse struct {
	Success bool   `json:"success" example:"false" description:"Indica que la operación falló"`
	Message string `json:"message" example:"Error al procesar la solicitud" description:"Mensaje de error"`
	Error   string `json:"error,omitempty" example:"Detalles específicos del error" description:"Detalles del error"`
}
