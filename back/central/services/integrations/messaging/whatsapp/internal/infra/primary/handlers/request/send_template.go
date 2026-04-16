package request

// SendTemplateRequest define la estructura de la petición para enviar plantillas
type SendTemplateRequest struct {
	TemplateName string            `json:"template_name" binding:"required"`
	PhoneNumber  string            `json:"phone_number" binding:"required"`
	Variables    map[string]string `json:"variables"`
	OrderNumber  string            `json:"order_number" binding:"required"`
	BusinessID   uint              `json:"business_id" binding:"required"`
}
