package request

// ListMessageAudit representa los query params para listar logs de auditoría
type ListMessageAudit struct {
	BusinessID   uint    `form:"business_id"`
	Status       *string `form:"status"`
	Direction    *string `form:"direction"`
	TemplateName *string `form:"template_name"`
	DateFrom     *string `form:"date_from"`
	DateTo       *string `form:"date_to"`
	Page         int     `form:"page"`
	PageSize     int     `form:"page_size"`
}

// StatsMessageAudit representa los query params para estadísticas
type StatsMessageAudit struct {
	BusinessID uint    `form:"business_id"`
	DateFrom   *string `form:"date_from"`
	DateTo     *string `form:"date_to"`
}
