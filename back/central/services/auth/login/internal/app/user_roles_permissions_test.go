package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uintPtr(v uint) *uint {
	return &v
}

func jwtQueValida(userID uint, suscripcion string) *mocks.JWTServiceMock {
	return &mocks.JWTServiceMock{
		ValidateTokenFn: func(tokenString string) (*domain.JWTClaims, error) {
			return &domain.JWTClaims{UserID: userID, SubscriptionStatus: suscripcion}, nil
		},
	}
}

func relacionConNegocio(businessID uint, roleID *uint) *domain.BusinessStaffRelation {
	return &domain.BusinessStaffRelation{
		BusinessID: &businessID,
		RoleID:     roleID,
		Business: &domain.BusinessInfoEntity{
			ID:               businessID,
			Name:             "Mi Tienda",
			Code:             "MT",
			BusinessTypeID:   3,
			BusinessTypeName: "Retail",
			BusinessTypeCode: "retail",
		},
	}
}

func permiso(id, resourceID uint, recurso, accion string) domain.Permission {
	return domain.Permission{
		ID:         id,
		Name:       recurso + "." + accion,
		Resource:   recurso,
		Action:     accion,
		ResourceID: resourceID,
	}
}

func TestGetUserRolesPermissions_SinToken_Rechaza(t *testing.T) {
	uc := buildLoginUseCase(&mocks.AuthRepositoryMock{}, nil, nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 10, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "token")
	assert.Nil(t, got)
}

func TestGetUserRolesPermissions_TokenInvalido_Rechaza(t *testing.T) {
	consultado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			consultado = true
			return nil, nil
		},
	}
	jwt := &mocks.JWTServiceMock{
		ValidateTokenFn: func(tokenString string) (*domain.JWTClaims, error) {
			return nil, errors.New("firma invalida")
		},
	}
	uc := buildLoginUseCase(repo, jwt, nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 10, "token-falso")

	require.Error(t, err)
	assert.Nil(t, got)
	assert.False(t, consultado, "con token invalido no se consulta la base de datos")
}

func TestGetUserRolesPermissions_TokenDeOtroUsuario_Rechaza(t *testing.T) {
	consultado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			consultado = true
			return usuarioActivo(userID, "x"), nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(99, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 10, "token-de-otro")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "acceso denegado")
	assert.Nil(t, got)
	assert.False(t, consultado,
		"un token valido de otro usuario no da acceso a sus permisos")
}

func TestGetUserRolesPermissions_UsuarioInexistente_Rechaza(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return nil, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 10, "token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "usuario no encontrado")
	assert.Nil(t, got)
}

func TestGetUserRolesPermissions_UsuarioInactivo_Rechaza(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(userID, "x")
			u.IsActive = false
			return u, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 10, "token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "usuario inactivo")
	assert.Nil(t, got)
}

func TestGetUserRolesPermissions_SinRelacionConElNegocio_Rechaza(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return nil, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 10, "token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "usuario-business no encontrada")
	assert.Nil(t, got)
}

func TestGetUserRolesPermissions_PasaElNegocioSoloSiNoEsCero(t *testing.T) {
	var recibido *uint
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			recibido = businessID
			return nil, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	_, _ = uc.GetUserRolesPermissions(context.Background(), 1, 0, "token")
	assert.Nil(t, recibido, "business_id cero se envia como nulo")

	_, _ = uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")
	require.NotNil(t, recibido)
	assert.Equal(t, uint(36), *recibido)
}

func TestGetUserRolesPermissions_SuperAdmin_SinNegocioEnLaRelacion(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return &domain.BusinessStaffRelation{BusinessID: nil, RoleID: nil}, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 0, "token")

	require.NoError(t, err)
	assert.True(t, got.IsSuper,
		"business_id nulo en business_staff identifica al super admin")
	assert.Equal(t, uint(0), got.BusinessID)
}

func TestGetUserRolesPermissions_RolDePlataforma_MarcaSuper(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return relacionConNegocio(36, uintPtr(5)), nil
		},
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, Name: "super", ScopeCode: "platform"}, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")

	require.NoError(t, err)
	assert.True(t, got.IsSuper, "un rol con scope platform eleva a super aunque tenga negocio")
}

func TestGetUserRolesPermissions_UsuarioDeNegocio_NoEsSuper(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return relacionConNegocio(36, uintPtr(5)), nil
		},
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, Name: "vendedor", ScopeCode: "business", ScopeName: "Negocio"}, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")

	require.NoError(t, err)
	assert.False(t, got.IsSuper)
	assert.Equal(t, uint(36), got.BusinessID)
	assert.Equal(t, "Mi Tienda", got.BusinessName)
	assert.Equal(t, uint(3), got.BusinessTypeID)
	assert.Equal(t, "Retail", got.BusinessTypeName)
	assert.Equal(t, "vendedor", got.Role.Name)
}

func TestGetUserRolesPermissions_NegocioNuloEnLaRelacion_Rechaza(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			bid := uint(36)
			return &domain.BusinessStaffRelation{BusinessID: &bid, RoleID: uintPtr(5), Business: nil}, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "business no encontrado")
	assert.Nil(t, got)
}

func TestGetUserRolesPermissions_FiltraPermisosPorRecursosActivos(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return relacionConNegocio(36, uintPtr(5)), nil
		},
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, Name: "vendedor", ScopeCode: "business"}, nil
		},
		GetBusinessConfiguredResourcesIDsFn: func(ctx context.Context, businessID uint) ([]uint, error) {
			return []uint{1}, nil
		},
		GetRolePermissionsFn: func(ctx context.Context, roleID uint) ([]domain.Permission, error) {
			return []domain.Permission{
				permiso(1, 1, "orders", "read"),
				permiso(2, 2, "inventory", "write"),
			}, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")

	require.NoError(t, err)
	require.Len(t, got.Permissions, 1,
		"solo se devuelven permisos de recursos configurados para el negocio")
	assert.Equal(t, "orders", got.Permissions[0].Resource)
	assert.True(t, got.Permissions[0].Active)
}

func TestGetUserRolesPermissions_SinRecursosConfigurados_DevuelveTodos(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return relacionConNegocio(36, uintPtr(5)), nil
		},
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, ScopeCode: "business"}, nil
		},
		GetBusinessConfiguredResourcesIDsFn: func(ctx context.Context, businessID uint) ([]uint, error) {
			return nil, nil
		},
		GetRolePermissionsFn: func(ctx context.Context, roleID uint) ([]domain.Permission, error) {
			return []domain.Permission{
				permiso(1, 1, "orders", "read"),
				permiso(2, 2, "inventory", "write"),
			}, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")

	require.NoError(t, err)
	assert.Len(t, got.Permissions, 2,
		"un negocio sin recursos configurados no filtra nada")
}

func TestGetUserRolesPermissions_SuperAdmin_TodosLosPermisosActivos(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return &domain.BusinessStaffRelation{BusinessID: nil, RoleID: uintPtr(1)}, nil
		},
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, ScopeCode: "platform"}, nil
		},
		GetRolePermissionsFn: func(ctx context.Context, roleID uint) ([]domain.Permission, error) {
			return []domain.Permission{
				permiso(1, 1, "orders", "read"),
				permiso(2, 99, "admin", "write"),
			}, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 0, "token")

	require.NoError(t, err)
	assert.True(t, got.IsSuper)
	require.Len(t, got.Permissions, 2)
	for _, p := range got.Permissions {
		assert.True(t, p.Active, "el super admin ve todos sus permisos activos")
	}
}

func TestGetUserRolesPermissions_EliminaPermisosDuplicados(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return &domain.BusinessStaffRelation{BusinessID: nil, RoleID: uintPtr(1)}, nil
		},
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, ScopeCode: "platform"}, nil
		},
		GetRolePermissionsFn: func(ctx context.Context, roleID uint) ([]domain.Permission, error) {
			return []domain.Permission{
				permiso(1, 1, "orders", "read"),
				permiso(2, 1, "orders", "read"),
				permiso(3, 1, "orders", "write"),
			}, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 0, "token")

	require.NoError(t, err)
	assert.Len(t, got.Permissions, 2,
		"la combinacion recurso:accion se deduplica")
}

func TestGetUserRolesPermissions_SinRol_SinPermisos(t *testing.T) {
	consultadoPermisos := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return relacionConNegocio(36, nil), nil
		},
		GetRolePermissionsFn: func(ctx context.Context, roleID uint) ([]domain.Permission, error) {
			consultadoPermisos = true
			return nil, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")

	require.NoError(t, err)
	assert.Empty(t, got.Permissions)
	assert.False(t, consultadoPermisos, "sin rol asignado no hay permisos que cargar")
}

func TestGetUserRolesPermissions_PropagaLaSuscripcionDelNegocio(t *testing.T) {
	relacion := relacionConNegocio(36, uintPtr(5))
	relacion.Business.SubscriptionStatus = "cancelled"

	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return relacion, nil
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")

	require.NoError(t, err)
	assert.Equal(t, "cancelled", got.SubscriptionStatus,
		"el estado de suscripcion sale de business.subscription_status en vivo, no del JWT viejo")
}

func TestGetUserRolesPermissions_FallaRelacion_ErrorInterno(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return nil, errors.New("connection refused to postgres at 10.0.0.5")
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")

	require.Error(t, err)
	assert.Equal(t, "error interno del servidor", err.Error())
	assert.NotContains(t, err.Error(), "postgres")
	assert.Nil(t, got)
}

func TestGetUserRolesPermissions_FallaRecursos_ErrorInterno(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return relacionConNegocio(36, uintPtr(5)), nil
		},
		GetBusinessConfiguredResourcesIDsFn: func(ctx context.Context, businessID uint) ([]uint, error) {
			return nil, errors.New("db down")
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")

	require.Error(t, err)
	assert.Equal(t, "error interno del servidor", err.Error())
	assert.Nil(t, got)
}

func TestGetUserRolesPermissions_FallaCargarPermisos_DevuelveRespuestaSinPermisos(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "x"), nil
		},
		GetBusinessStaffRelationFn: func(ctx context.Context, userID uint, businessID *uint) (*domain.BusinessStaffRelation, error) {
			return relacionConNegocio(36, uintPtr(5)), nil
		},
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, ScopeCode: "business"}, nil
		},
		GetRolePermissionsFn: func(ctx context.Context, roleID uint) ([]domain.Permission, error) {
			return nil, errors.New("db down")
		},
	}
	uc := buildLoginUseCase(repo, jwtQueValida(1, "active"), nil)

	got, err := uc.GetUserRolesPermissions(context.Background(), 1, 36, "token")

	require.NoError(t, err)
	assert.Empty(t, got.Permissions,
		"un fallo cargando permisos deja al usuario sin permisos, no con todos")
}
