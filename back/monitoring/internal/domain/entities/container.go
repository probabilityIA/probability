package entities

import "time"

type Container struct {
	ID        string
	Name      string
	Service   string
	Project   string
	Image     string
	State     string
	Status    string
	Health    string
	CreatedAt time.Time
	StartedAt time.Time
	Ports     []PortMapping
}

type PortMapping struct {
	HostPort      uint16
	ContainerPort uint16
	Protocol      string
}

type ContainerStats struct {
	ContainerID string
	CPUPercent  float64
	MemoryUsage uint64
	MemoryLimit uint64
	MemoryPercent float64
	NetworkRx   uint64
	NetworkTx   uint64
}

type LogEntry struct {
	Timestamp string
	Stream    string
	Message   string
}

type SystemStats struct {
	CPUPercent    float64
	CPUCores      int
	MemoryTotal   uint64
	MemoryUsed    uint64
	MemoryPercent float64
	DiskTotal     uint64
	DiskUsed      uint64
	DiskPercent   float64
}

type ComposeService struct {
	Name         string
	ContainerID  string
	State        string
	Health       string
	DependsOn    []string
	Ports        []PortMapping
	Image        string
}
