package docker

import (
	"github.com/docker/docker/client"
	"github.com/secamc93/probability/back/monitoring/internal/domain/ports"
)

type DockerClient struct {
	cli     *client.Client
	project string
}

func New(projectName string) (ports.IDockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerClient{
		cli:     cli,
		project: projectName,
	}, nil
}
