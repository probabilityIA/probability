package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Iapp interface {
	GetUsers(ctx context.Context, filters domain.UserFilters) (*domain.UserListDTO, error)
	GetUserByID(ctx context.Context, id uint) (*domain.UserDTO, error)
	CreateUser(ctx context.Context, user domain.CreateUserDTO) (string, string, string, error)
	UpdateUser(ctx context.Context, id uint, user domain.UpdateUserDTO) (string, error)
	DeleteUser(ctx context.Context, id uint) (string, error)
	AssignRoleToUserBusiness(ctx context.Context, userID uint, assignments []domain.BusinessRoleAssignment) error
}

// UserUseCase implementa los casos de uso para usuarios
type UserUseCase struct {
	repository domain.IUserRepository
	log        log.ILogger
	s3         domain.IS3Service
	env        env.IConfig
}

// NewUserUseCase crea una nueva instancia del caso de uso de usuarios
func New(repository domain.IUserRepository, log log.ILogger, s3 domain.IS3Service, env env.IConfig) Iapp {
	return &UserUseCase{
		repository: repository,
		log:        log,
		s3:         s3,
		env:        env,
	}
}
