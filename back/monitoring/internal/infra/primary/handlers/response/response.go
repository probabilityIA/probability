package response

import "time"

type LoginResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ContainerResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Service   string         `json:"service"`
	Project   string         `json:"project"`
	Image     string         `json:"image"`
	State     string         `json:"state"`
	Status    string         `json:"status"`
	Health    string         `json:"health"`
	CreatedAt time.Time      `json:"created_at"`
	StartedAt time.Time      `json:"started_at"`
	Ports     []PortResponse `json:"ports"`
}

type PortResponse struct {
	HostPort      uint16 `json:"host_port"`
	ContainerPort uint16 `json:"container_port"`
	Protocol      string `json:"protocol"`
}

type StatsResponse struct {
	ContainerID   string  `json:"container_id"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsage   uint64  `json:"memory_usage"`
	MemoryLimit   uint64  `json:"memory_limit"`
	MemoryPercent float64 `json:"memory_percent"`
	NetworkRx     uint64  `json:"network_rx"`
	NetworkTx     uint64  `json:"network_tx"`
}

type LogEntryResponse struct {
	Timestamp string `json:"timestamp"`
	Stream    string `json:"stream"`
	Message   string `json:"message"`
}

type ComposeServiceResponse struct {
	Name        string         `json:"name"`
	ContainerID string         `json:"container_id"`
	State       string         `json:"state"`
	Health      string         `json:"health"`
	DependsOn   []string       `json:"depends_on"`
	Ports       []PortResponse `json:"ports"`
	Image       string         `json:"image"`
}

type ActionResponse struct {
	Message string `json:"message"`
}
