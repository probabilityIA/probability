package app

import (
	"github.com/secamc93/probability/back/monitoring/internal/domain/ports"
)

type useCase struct {
	docker    ports.IDockerService
	userRepo  ports.IUserRepository
	jwtSecret string
}

func New(docker ports.IDockerService, userRepo ports.IUserRepository, jwtSecret string) ports.IUseCase {
	return &useCase{
		docker:    docker,
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}
