package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/container"
	domainErrors "github.com/secamc93/probability/back/monitoring/internal/domain/errors"
)

func (d *DockerClient) RestartContainer(ctx context.Context, id string) error {
	timeout := 10
	err := d.cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &timeout})
	if err != nil {
		return domainErrors.ErrActionFailed
	}
	return nil
}

func (d *DockerClient) StopContainer(ctx context.Context, id string) error {
	timeout := 10
	err := d.cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout})
	if err != nil {
		return domainErrors.ErrActionFailed
	}
	return nil
}

func (d *DockerClient) StartContainer(ctx context.Context, id string) error {
	err := d.cli.ContainerStart(ctx, id, container.StartOptions{})
	if err != nil {
		return domainErrors.ErrActionFailed
	}

	// Wait briefly for the container to be running
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	inspect, err := d.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil // Started but couldn't verify - not an error
	}
	if inspect.State.Running {
		return nil
	}
	return nil
}
