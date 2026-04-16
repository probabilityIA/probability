package docker

import (
	"context"
	"sort"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/monitoring/internal/domain/errors"
)

func (d *DockerClient) ListContainers(ctx context.Context) ([]entities.Container, error) {
	opts := container.ListOptions{
		All: true,
	}

	if d.project != "" {
		opts.Filters = filters.NewArgs(
			filters.Arg("label", "com.docker.compose.project="+d.project),
		)
	}

	containers, err := d.cli.ContainerList(ctx, opts)
	if err != nil {
		return nil, err
	}

	result := make([]entities.Container, 0, len(containers))
	for _, c := range containers {
		result = append(result, mapContainer(c))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Service < result[j].Service
	})

	return result, nil
}

func (d *DockerClient) GetContainer(ctx context.Context, id string) (*entities.Container, error) {
	inspect, err := d.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, domainErrors.ErrContainerNotFound
	}

	ports := make([]entities.PortMapping, 0)
	for containerPort, bindings := range inspect.NetworkSettings.Ports {
		for _, binding := range bindings {
			var hostPort uint16
			if binding.HostPort != "" {
				for _, ch := range binding.HostPort {
					hostPort = hostPort*10 + uint16(ch-'0')
				}
			}
			ports = append(ports, entities.PortMapping{
				HostPort:      hostPort,
				ContainerPort: uint16(containerPort.Int()),
				Protocol:      containerPort.Proto(),
			})
		}
	}

	var startedAt time.Time
	if inspect.State.StartedAt != "" {
		startedAt, _ = time.Parse(time.RFC3339Nano, inspect.State.StartedAt)
	}

	health := ""
	if inspect.State.Health != nil {
		health = inspect.State.Health.Status
	}

	c := &entities.Container{
		ID:        inspect.ID[:12],
		Name:      inspect.Name[1:], // remove leading /
		Service:   inspect.Config.Labels["com.docker.compose.service"],
		Project:   inspect.Config.Labels["com.docker.compose.project"],
		Image:     inspect.Config.Image,
		State:     inspect.State.Status,
		Status:    inspect.State.Status,
		Health:    health,
		CreatedAt: parseTime(inspect.Created),
		StartedAt: startedAt,
		Ports:     ports,
	}

	return c, nil
}

func (d *DockerClient) GetComposeServices(ctx context.Context) ([]entities.ComposeService, error) {
	containers, err := d.ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	services := make([]entities.ComposeService, 0, len(containers))
	for _, c := range containers {
		services = append(services, entities.ComposeService{
			Name:        c.Service,
			ContainerID: c.ID,
			State:       c.State,
			Health:      c.Health,
			Ports:       c.Ports,
			Image:       c.Image,
		})
	}

	return services, nil
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339Nano, s)
	return t
}

func mapContainer(c types.Container) entities.Container {
	ports := make([]entities.PortMapping, 0, len(c.Ports))
	for _, p := range c.Ports {
		ports = append(ports, entities.PortMapping{
			HostPort:      p.PublicPort,
			ContainerPort: p.PrivatePort,
			Protocol:      p.Type,
		})
	}

	name := ""
	if len(c.Names) > 0 {
		name = c.Names[0][1:] // remove leading /
	}

	health := ""
	if c.State == "running" {
		health = "healthy"
	}

	return entities.Container{
		ID:      c.ID[:12],
		Name:    name,
		Service: c.Labels["com.docker.compose.service"],
		Project: c.Labels["com.docker.compose.project"],
		Image:   c.Image,
		State:   c.State,
		Status:  c.Status,
		Health:  health,
		Ports:   ports,
	}
}
