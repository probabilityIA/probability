package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type configStub struct {
	valores map[string]string
}

func (c *configStub) Get(key string) string {
	return c.valores[key]
}

func buildLoginUseCase(repo *mocks.AuthRepositoryMock, jwt *mocks.JWTServiceMock, valores map[string]string) *AuthUseCase {
	if jwt == nil {
		jwt = &mocks.JWTServiceMock{}
	}
	if valores == nil {
		valores = map[string]string{}
	}
	return &AuthUseCase{
		repository: repo,
		jwtService: jwt,
		log:        mocks.NewSilentLogger(),
		env:        &configStub{valores: valores},
	}
}

func rolPlataforma() domain.Role {
	return domain.Role{ID: 1, Name: "super_admin", ScopeCode: "platform"}
}

func rolNegocio() domain.Role {
	return domain.Role{ID: 2, Name: "admin", ScopeCode: "business"}
}

func negocio(id uint) domain.BusinessInfoEntity {
	return domain.BusinessInfoEntity{
		ID:                 id,
		Name:               "Mi Tienda",
		Code:               "MT",
		BusinessTypeID:     3,
		SubscriptionStatus: "active",
	}
}

func TestIsSuperAdmin(t *testing.T) {
	assert.True(t, isSuperAdmin([]domain.Role{rolPlataforma()}))
	assert.True(t, isSuperAdmin([]domain.Role{rolNegocio(), rolPlataforma()}))
	assert.False(t, isSuperAdmin([]domain.Role{rolNegocio()}))
	assert.False(t, isSuperAdmin(nil))
	assert.False(t, isSuperAdmin([]domain.Role{}))
}

func TestIsSuperAdmin_SoloElScopePlatformCuenta(t *testing.T) {
	assert.False(t, isSuperAdmin([]domain.Role{{Name: "super_admin", ScopeCode: "business"}}),
		"el nombre del rol no define super admin, solo el scope_code")
	assert.False(t, isSuperAdmin([]domain.Role{{ScopeCode: "Platform"}}),
		"el scope es sensible a mayusculas")
}

func TestLogin_CredencialesVacias_Rechaza(t *testing.T) {
	consultado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			consultado = true
			return nil, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	casos := []domain.LoginRequest{
		{Email: "", Password: "x"},
		{Email: "ana@test.com", Password: ""},
		{Email: "   ", Password: "x"},
	}

	for _, req := range casos {
		got, err := uc.Login(context.Background(), req)

		assert.ErrorIs(t, err, domain.ErrEmailPasswordRequired)
		assert.Nil(t, got)
	}
	assert.False(t, consultado)
}

func TestLogin_NormalizaElCorreo(t *testing.T) {
	var consultado string
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			consultado = email
			return nil, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	_, _ = uc.Login(context.Background(), domain.LoginRequest{
		Email:    "  ANA@Test.COM  ",
		Password: "x",
	})

	assert.Equal(t, "ana@test.com", consultado)
}

func TestLogin_UsuarioInexistente_RetornaNotFound(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return nil, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "noexiste@test.com",
		Password: "x",
	})

	assert.ErrorIs(t, err, domain.ErrUserNotFound)
	assert.Nil(t, got)
}

func TestLogin_ContrasenaInvalida_Rechaza(t *testing.T) {
	tokenGenerado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(1, "la-correcta"), nil
		},
	}
	jwt := &mocks.JWTServiceMock{
		GenerateTokenFn: func(userID, businessID, businessTypeID, roleID uint, subscriptionStatus string) (string, error) {
			tokenGenerado = true
			return "token", nil
		},
	}
	uc := buildLoginUseCase(repo, jwt, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "la-equivocada",
	})

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	assert.Nil(t, got)
	assert.False(t, tokenGenerado, "sin contrasena valida no se emite token")
}

func TestLogin_UsuarioInactivoSinVerificacion_RetornaInactivo(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(1, "correcta")
			u.IsActive = false
			return u, nil
		},
		HasPendingEmailVerificationFn: func(ctx context.Context, userID uint) (bool, error) {
			return false, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	assert.ErrorIs(t, err, domain.ErrUserInactive)
	assert.Nil(t, got)
}

func TestLogin_UsuarioPendienteDeVerificacion_ErrorEspecifico(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(1, "correcta")
			u.IsActive = false
			return u, nil
		},
		HasPendingEmailVerificationFn: func(ctx context.Context, userID uint) (bool, error) {
			return true, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	assert.ErrorIs(t, err, domain.ErrUserPendingVerification)
	assert.Nil(t, got)
}

func TestLogin_FallaConsultarVerificacion_CaeAInactivo(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(1, "correcta")
			u.IsActive = false
			return u, nil
		},
		HasPendingEmailVerificationFn: func(ctx context.Context, userID uint) (bool, error) {
			return false, errors.New("db down")
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	assert.ErrorIs(t, err, domain.ErrUserInactive,
		"ante un fallo se rechaza el acceso, no se concede")
	assert.Nil(t, got)
}

func TestLogin_SuperAdmin_TokenConBusinessCero(t *testing.T) {
	var tokenBusinessID, tokenRoleID uint
	var tokenSubscripcion string
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(42, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolPlataforma()}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return []domain.BusinessInfoEntity{negocio(36)}, nil
		},
	}
	jwt := &mocks.JWTServiceMock{
		GenerateTokenFn: func(userID, businessID, businessTypeID, roleID uint, subscriptionStatus string) (string, error) {
			tokenBusinessID, tokenRoleID, tokenSubscripcion = businessID, roleID, subscriptionStatus
			return "token-super", nil
		},
	}
	uc := buildLoginUseCase(repo, jwt, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.IsSuperAdmin)
	assert.Equal(t, "platform", got.Scope)
	assert.Equal(t, uint(0), tokenBusinessID,
		"un super admin lleva business_id = 0 en el token aunque tenga negocios asociados")
	assert.Equal(t, uint(1), tokenRoleID)
	assert.Equal(t, "active", tokenSubscripcion)
}

func TestLogin_UsuarioDeNegocio_TokenConSuPrimerNegocio(t *testing.T) {
	var tokenBusinessID, tokenTipoNegocio, tokenRoleID uint
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(7, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return []domain.BusinessInfoEntity{negocio(36), negocio(99)}, nil
		},
		GetUserRoleByBusinessFn: func(ctx context.Context, userID, businessID uint) (*domain.Role, error) {
			return &domain.Role{ID: 5, Name: "vendedor"}, nil
		},
	}
	jwt := &mocks.JWTServiceMock{
		GenerateTokenFn: func(userID, businessID, businessTypeID, roleID uint, subscriptionStatus string) (string, error) {
			tokenBusinessID, tokenTipoNegocio, tokenRoleID = businessID, businessTypeID, roleID
			return "token", nil
		},
	}
	uc := buildLoginUseCase(repo, jwt, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.False(t, got.IsSuperAdmin)
	assert.Equal(t, "business", got.Scope)
	assert.Equal(t, uint(36), tokenBusinessID, "se usa el primer negocio del usuario")
	assert.Equal(t, uint(3), tokenTipoNegocio)
	assert.Equal(t, uint(5), tokenRoleID, "se usa el rol especifico en ese negocio")
	assert.Len(t, got.Businesses, 2)
}

func TestLogin_SinRolEnElNegocio_CaeAlPrimerRol(t *testing.T) {
	var tokenRoleID uint
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(7, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{{ID: 9, Name: "operario", ScopeCode: "business"}}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return []domain.BusinessInfoEntity{negocio(36)}, nil
		},
		GetUserRoleByBusinessFn: func(ctx context.Context, userID, businessID uint) (*domain.Role, error) {
			return nil, errors.New("sin rol asignado")
		},
	}
	jwt := &mocks.JWTServiceMock{
		GenerateTokenFn: func(userID, businessID, businessTypeID, roleID uint, subscriptionStatus string) (string, error) {
			tokenRoleID = roleID
			return "token", nil
		},
	}
	uc := buildLoginUseCase(repo, jwt, nil)

	_, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.Equal(t, uint(9), tokenRoleID)
}

func TestLogin_UsuarioSinNegocios_TokenConBusinessCero(t *testing.T) {
	var tokenBusinessID uint
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(7, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return nil, nil
		},
	}
	jwt := &mocks.JWTServiceMock{
		GenerateTokenFn: func(userID, businessID, businessTypeID, roleID uint, subscriptionStatus string) (string, error) {
			tokenBusinessID = businessID
			return "token", nil
		},
	}
	uc := buildLoginUseCase(repo, jwt, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.Equal(t, uint(0), tokenBusinessID)
	assert.Empty(t, got.Businesses)
}

func TestLogin_SuscripcionVacia_CaeAActive(t *testing.T) {
	var tokenSubscripcion string
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(7, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			b := negocio(36)
			b.SubscriptionStatus = ""
			return []domain.BusinessInfoEntity{b}, nil
		},
	}
	jwt := &mocks.JWTServiceMock{
		GenerateTokenFn: func(userID, businessID, businessTypeID, roleID uint, subscriptionStatus string) (string, error) {
			tokenSubscripcion = subscriptionStatus
			return "token", nil
		},
	}
	uc := buildLoginUseCase(repo, jwt, nil)

	_, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.Equal(t, "active", tokenSubscripcion)
}

func TestLogin_SuscripcionVencida_SePropagaAlToken(t *testing.T) {
	var tokenSubscripcion string
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(7, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			b := negocio(36)
			b.SubscriptionStatus = "expired"
			return []domain.BusinessInfoEntity{b}, nil
		},
	}
	jwt := &mocks.JWTServiceMock{
		GenerateTokenFn: func(userID, businessID, businessTypeID, roleID uint, subscriptionStatus string) (string, error) {
			tokenSubscripcion = subscriptionStatus
			return "token", nil
		},
	}
	uc := buildLoginUseCase(repo, jwt, nil)

	_, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.Equal(t, "expired", tokenSubscripcion)
}

func TestLogin_PrimerLogin_ExigeCambioDeContrasena(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(7, "correcta")
			u.LastLoginAt = nil
			return u, nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.True(t, got.RequirePasswordChange)
}

func TestLogin_LoginPosterior_NoExigeCambio(t *testing.T) {
	anterior := time.Now().Add(-24 * time.Hour)
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(7, "correcta")
			u.LastLoginAt = &anterior
			return u, nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.False(t, got.RequirePasswordChange)
}

func TestLogin_CompletaLaURLDelAvatar(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(7, "correcta")
			u.AvatarURL = "avatars/ana.jpg"
			return u, nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, map[string]string{"URL_BASE_DOMAIN_S3": "https://cdn.test.com/"})

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.Equal(t, "https://cdn.test.com/avatars/ana.jpg", got.User.AvatarURL)
}

func TestLogin_AvatarQueYaEsURL_NoSeToca(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(7, "correcta")
			u.AvatarURL = "https://otro-cdn.com/ana.jpg"
			return u, nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, map[string]string{"URL_BASE_DOMAIN_S3": "https://cdn.test.com"})

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.Equal(t, "https://otro-cdn.com/ana.jpg", got.User.AvatarURL)
}

func TestLogin_SinBaseDeS3_DejaLaRutaRelativa(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(7, "correcta")
			u.AvatarURL = "avatars/ana.jpg"
			return u, nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.Equal(t, "avatars/ana.jpg", got.User.AvatarURL)
}

func TestLogin_FallaObtenerRoles_RetornaErrorInterno(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(7, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return nil, errors.New("connection refused to postgres at 10.0.0.5")
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.Error(t, err)
	assert.Equal(t, "error interno del servidor", err.Error())
	assert.NotContains(t, err.Error(), "postgres")
	assert.Nil(t, got)
}

func TestLogin_FallaObtenerNegocios_NoAbortaElLogin(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(7, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return nil, errors.New("db down")
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err, "un fallo al listar negocios no impide iniciar sesion")
	assert.Empty(t, got.Businesses)
}

func TestLogin_FallaGenerarToken_RetornaErrorInterno(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(7, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
	}
	jwt := &mocks.JWTServiceMock{
		GenerateTokenFn: func(userID, businessID, businessTypeID, roleID uint, subscriptionStatus string) (string, error) {
			return "", errors.New("clave jwt invalida")
		},
	}
	uc := buildLoginUseCase(repo, jwt, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.Error(t, err)
	assert.Equal(t, "error interno del servidor", err.Error())
	assert.Nil(t, got)
}

func TestLogin_FallaActualizarUltimoLogin_NoAbortaElLogin(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(7, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
		UpdateLastLoginFn: func(ctx context.Context, userID uint) error {
			return errors.New("db down")
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, got.Token)
}

func TestLogin_CompletaLasImagenesDeLosNegocios(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(7, "correcta"), nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			b := negocio(36)
			b.LogoURL = "logos/mt.png"
			b.NavbarImageURL = "navbars/mt.png"
			return []domain.BusinessInfoEntity{b}, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, map[string]string{"URL_BASE_DOMAIN_S3": "https://cdn.test.com"})

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	require.Len(t, got.Businesses, 1)
	assert.Equal(t, "https://cdn.test.com/logos/mt.png", got.Businesses[0].LogoURL)
	assert.Equal(t, "https://cdn.test.com/navbars/mt.png", got.Businesses[0].NavbarImageURL)
}

func TestLogin_DevuelveLosDatosDelUsuario(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(7, "correcta")
			u.Name = "Ana Gomez"
			u.Phone = "3001234567"
			return u, nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{rolNegocio()}, nil
		},
	}
	uc := buildLoginUseCase(repo, nil, nil)

	got, err := uc.Login(context.Background(), domain.LoginRequest{
		Email:    "ana@test.com",
		Password: "correcta",
	})

	require.NoError(t, err)
	assert.Equal(t, uint(7), got.User.ID)
	assert.Equal(t, "Ana Gomez", got.User.Name)
	assert.Equal(t, "ana@test.com", got.User.Email)
	assert.Equal(t, "3001234567", got.User.Phone)
	assert.True(t, got.User.IsActive)
	assert.NotEmpty(t, got.Token)
}
