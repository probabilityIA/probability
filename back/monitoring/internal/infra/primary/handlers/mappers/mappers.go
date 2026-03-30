package mappers

import (
	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
	"github.com/secamc93/probability/back/monitoring/internal/infra/primary/handlers/response"
)

func ContainerToResponse(c *entities.Container) response.ContainerResponse {
	ports := make([]response.PortResponse, len(c.Ports))
	for i, p := range c.Ports {
		ports[i] = response.PortResponse{
			HostPort:      p.HostPort,
			ContainerPort: p.ContainerPort,
			Protocol:      p.Protocol,
		}
	}
	return response.ContainerResponse{
		ID:        c.ID,
		Name:      c.Name,
		Service:   c.Service,
		Project:   c.Project,
		Image:     c.Image,
		State:     c.State,
		Status:    c.Status,
		Health:    c.Health,
		CreatedAt: c.CreatedAt,
		StartedAt: c.StartedAt,
		Ports:     ports,
	}
}

func ContainersToResponse(containers []entities.Container) []response.ContainerResponse {
	result := make([]response.ContainerResponse, len(containers))
	for i := range containers {
		result[i] = ContainerToResponse(&containers[i])
	}
	return result
}

func StatsToResponse(s *entities.ContainerStats) response.StatsResponse {
	return response.StatsResponse{
		ContainerID:   s.ContainerID,
		CPUPercent:    s.CPUPercent,
		MemoryUsage:   s.MemoryUsage,
		MemoryLimit:   s.MemoryLimit,
		MemoryPercent: s.MemoryPercent,
		NetworkRx:     s.NetworkRx,
		NetworkTx:     s.NetworkTx,
	}
}

func LogsToResponse(logs []entities.LogEntry) []response.LogEntryResponse {
	result := make([]response.LogEntryResponse, len(logs))
	for i, l := range logs {
		result[i] = response.LogEntryResponse{
			Timestamp: l.Timestamp,
			Stream:    l.Stream,
			Message:   l.Message,
		}
	}
	return result
}

func ComposeServicesToResponse(services []entities.ComposeService) []response.ComposeServiceResponse {
	result := make([]response.ComposeServiceResponse, len(services))
	for i, s := range services {
		ports := make([]response.PortResponse, len(s.Ports))
		for j, p := range s.Ports {
			ports[j] = response.PortResponse{
				HostPort:      p.HostPort,
				ContainerPort: p.ContainerPort,
				Protocol:      p.Protocol,
			}
		}
		result[i] = response.ComposeServiceResponse{
			Name:        s.Name,
			ContainerID: s.ContainerID,
			State:       s.State,
			Health:      s.Health,
			DependsOn:   s.DependsOn,
			Ports:       ports,
			Image:       s.Image,
		}
	}
	return result
}
