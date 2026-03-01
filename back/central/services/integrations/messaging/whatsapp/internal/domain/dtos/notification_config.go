package dtos

// NotificationConfigData contiene los datos de una configuración de notificación
// Este DTO es puro de dominio: sin tags json, sin dependencias de infraestructura
type NotificationConfigData struct {
	ID                  uint
	IntegrationID       uint
	NotificationType    string
	IsActive            bool
	TemplateName        string
	Language            string
	RecipientType       string
	Trigger             string
	Statuses            []string
	PaymentMethods      []uint
	SourceIntegrationID *uint
	Priority            int
	Description         string
}
