package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Iapp interface {
	Login(ctx context.Context, request domain.LoginRequest) (*domain.LoginResponse, error)
	GetUserRolesPermissions(ctx context.Context, userID uint, businessID uint, token string) (*domain.UserRolesPermissionsResponse, error)
	ChangePassword(ctx context.Context, request domain.ChangePasswordRequest) (*domain.ChangePasswordResponse, error)
	GeneratePassword(ctx context.Context, request domain.GeneratePasswordRequest) (*domain.GeneratePasswordResponse, error)
	GenerateBusinessToken(ctx context.Context, userID uint, businessID uint) (string, error)
	// GenerateAPIKey(ctx context.Context, request domain.GenerateAPIKeyRequest) (*domain.GenerateAPIKeyResponse, error)
	// ValidateAPIKey(ctx context.Context, request domain.ValidateAPIKeyRequest) (*domain.ValidateAPIKeyResponse, error)
}

type IAuthUseCase interface {
	ValidateAPIKey(ctx context.Context, request domain.ValidateAPIKeyRequest) (*domain.ValidateAPIKeyResponse, error)
}

type AuthUseCase struct {
	repository domain.IAuthRepository
	jwtService domain.IJWTService
	log        log.ILogger
	env        env.IConfig
}

func New(repository domain.IAuthRepository, jwtService domain.IJWTService, log log.ILogger, env env.IConfig) Iapp {
	return &AuthUseCase{
		repository: repository,
		jwtService: jwtService,
		log:        log,
		env:        env,
	}
}
