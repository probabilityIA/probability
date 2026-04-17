package domain

// DashboardStatsResponse es la respuesta del endpoint de estadísticas
type DashboardStatsResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	ServerTime string         `json:"server_time"`  // Fecha/hora actual del servidor (RFC3339)
	Data       DashboardStats `json:"data"`
}
