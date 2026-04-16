package response

// BusinessResourceConfiguredResponse representa un recurso configurado del negocio
//
//	@Description	Información detallada de un recurso del negocio
type BusinessResourceConfiguredResponse struct {
	ResourceID   uint   `json:"resource_id" example:"1" description:"ID único del recurso"`
	ResourceName string `json:"resource_name" example:"users" description:"Nombre del recurso"`
	IsActive     bool   `json:"is_active" example:"true" description:"Indica si el recurso está activo para el negocio"`
}

// BusinessResourcesResponse representa la respuesta completa de recursos del negocio
//
//	@Description	Lista completa de recursos del negocio con estadísticas
type BusinessResourcesResponse struct {
	BusinessID uint                                 `json:"business_id" example:"123" description:"ID del negocio"`
	Resources  []BusinessResourceConfiguredResponse `json:"resources" description:"Lista de recursos del negocio"`
	Total      int                                  `json:"total" example:"5" description:"Total de recursos disponibles"`
	Active     int                                  `json:"active" example:"3" description:"Número de recursos activos"`
	Inactive   int                                  `json:"inactive" example:"2" description:"Número de recursos inactivos"`
}

// GetBusinessResourcesResponse representa la respuesta para obtener recursos del negocio
//
//	@Description	Respuesta exitosa al obtener recursos del negocio
type GetBusinessResourcesResponse struct {
	Success bool                      `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string                    `json:"message" example:"Recursos del negocio obtenidos exitosamente" description:"Mensaje descriptivo de la operación"`
	Data    BusinessResourcesResponse `json:"data" description:"Datos de los recursos del negocio"`
}

// GetBusinessResourceStatusResponse representa la respuesta para obtener estado de un recurso
//
//	@Description	Respuesta exitosa al obtener el estado de un recurso específico
type GetBusinessResourceStatusResponse struct {
	Success bool                               `json:"success" example:"true" description:"Indica si la operación fue exitosa"`
	Message string                             `json:"message" example:"Estado del recurso obtenido exitosamente" description:"Mensaje descriptivo de la operación"`
	Data    BusinessResourceConfiguredResponse `json:"data" description:"Datos del estado del recurso"`
}
