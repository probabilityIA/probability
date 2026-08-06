package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/roles/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRolesUseCase(repo *mocks.RoleRepositoryMock) IUseCaseRole {
	if repo == nil {
		repo = &mocks.RoleRepositoryMock{}
	}
	return New(repo, mocks.NewSilentLogger())
}

func strPtr(v string) *string { return &v }
func intPtr(v int) *int       { return &v }
func uintPtr(v uint) *uint    { return &v }
func boolPtr(v bool) *bool    { return &v }

func catalogoDeRoles() []domain.Role {
	return []domain.Role{
		{ID: 1, Name: "Super Admin", Level: 1, IsSystem: true, ScopeID: 1, BusinessTypeID: 0},
		{ID: 2, Name: "Business Admin", Level: 2, IsSystem: true, ScopeID: 2, BusinessTypeID: 1},
		{ID: 3, Name: "Vendedor", Level: 3, IsSystem: false, ScopeID: 2, BusinessTypeID: 1},
		{ID: 4, Name: "Bodeguero", Level: 3, IsSystem: false, ScopeID: 2, BusinessTypeID: 2},
		{ID: 5, Name: "Operator Admin", Level: 2, IsSystem: true, ScopeID: 3, BusinessTypeID: 0},
	}
}

func repoConCatalogo() *mocks.RoleRepositoryMock {
	return &mocks.RoleRepositoryMock{
		GetRolesFn: func(ctx context.Context) ([]domain.Role, error) { return catalogoDeRoles(), nil },
	}
}

func idsDe(roles []domain.RoleDTO) []uint {
	out := make([]uint, len(roles))
	for i, r := range roles {
		out[i] = r.ID
	}
	return out
}

func TestGetRoles_SinFiltros_DevuelveTodo(t *testing.T) {
	got, err := newRolesUseCase(repoConCatalogo()).GetRoles(context.Background(), domain.RoleFilters{})

	require.NoError(t, err)
	assert.Equal(t, []uint{1, 2, 3, 4, 5}, idsDe(got))
}

func TestGetRoles_FiltraPorScopeID(t *testing.T) {
	got, err := newRolesUseCase(repoConCatalogo()).GetRoles(context.Background(), domain.RoleFilters{
		ScopeID: uintPtr(2),
	})

	require.NoError(t, err)
	assert.Equal(t, []uint{2, 3, 4}, idsDe(got))
}

func TestGetRoles_FiltraPorBusinessTypeID(t *testing.T) {
	got, err := newRolesUseCase(repoConCatalogo()).GetRoles(context.Background(), domain.RoleFilters{
		BusinessTypeID: uintPtr(1),
	})

	require.NoError(t, err)
	assert.Equal(t, []uint{2, 3}, idsDe(got))
}

func TestGetRoles_FiltraPorIsSystemEnAmbosSentidos(t *testing.T) {
	uc := newRolesUseCase(repoConCatalogo())

	soloSistema, err := uc.GetRoles(context.Background(), domain.RoleFilters{IsSystem: boolPtr(true)})
	require.NoError(t, err)
	assert.Equal(t, []uint{1, 2, 5}, idsDe(soloSistema))

	soloCustom, err := uc.GetRoles(context.Background(), domain.RoleFilters{IsSystem: boolPtr(false)})
	require.NoError(t, err)
	assert.Equal(t, []uint{3, 4}, idsDe(soloCustom))
}

func TestGetRoles_FiltraPorLevel(t *testing.T) {
	got, err := newRolesUseCase(repoConCatalogo()).GetRoles(context.Background(), domain.RoleFilters{
		Level: intPtr(3),
	})

	require.NoError(t, err)
	assert.Equal(t, []uint{3, 4}, idsDe(got))
}

func TestGetRoles_FiltraPorNombreParcialSinImportarMayusculas(t *testing.T) {
	casos := []struct {
		entrada string
		want    []uint
	}{
		{"admin", []uint{1, 2, 5}},
		{"ADMIN", []uint{1, 2, 5}},
		{"AdMiN", []uint{1, 2, 5}},
		{"vendedor", []uint{3}},
		{"super", []uint{1}},
		{"inexistente", []uint{}},
	}

	uc := newRolesUseCase(repoConCatalogo())
	for _, tc := range casos {
		t.Run(tc.entrada, func(t *testing.T) {
			got, err := uc.GetRoles(context.Background(), domain.RoleFilters{Name: strPtr(tc.entrada)})
			require.NoError(t, err)
			assert.Equal(t, tc.want, idsDe(got))
		})
	}
}

func TestGetRoles_NombreVacioNoFiltra(t *testing.T) {
	got, err := newRolesUseCase(repoConCatalogo()).GetRoles(context.Background(), domain.RoleFilters{
		Name: strPtr(""),
	})

	require.NoError(t, err)
	assert.Len(t, got, 5, "un nombre vacio no debe descartar nada")
}

func TestGetRoles_FiltrosCombinadosSeAplicanEnConjuncion(t *testing.T) {
	got, err := newRolesUseCase(repoConCatalogo()).GetRoles(context.Background(), domain.RoleFilters{
		ScopeID:        uintPtr(2),
		BusinessTypeID: uintPtr(1),
		IsSystem:       boolPtr(false),
	})

	require.NoError(t, err)
	assert.Equal(t, []uint{3}, idsDe(got))
}

func TestGetRoles_FiltrosContradictorios_DevuelveVacioNoNil(t *testing.T) {
	got, err := newRolesUseCase(repoConCatalogo()).GetRoles(context.Background(), domain.RoleFilters{
		ScopeID: uintPtr(1),
		Level:   intPtr(99),
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Empty(t, got)
}

func TestGetRoles_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.RoleRepositoryMock{
		GetRolesFn: func(ctx context.Context) ([]domain.Role, error) { return nil, dbErr },
	}

	got, err := newRolesUseCase(repo).GetRoles(context.Background(), domain.RoleFilters{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetRoles_MapeaTodosLosCamposAlDTO(t *testing.T) {
	origen := domain.Role{
		ID: 7, Name: "Vendedor", Description: "vende", Level: 3, IsSystem: false,
		ScopeID: 2, ScopeName: "Negocio", ScopeCode: "business",
		BusinessTypeID: 4, BusinessTypeName: "Retail",
	}
	repo := &mocks.RoleRepositoryMock{
		GetRolesFn: func(ctx context.Context) ([]domain.Role, error) { return []domain.Role{origen}, nil },
	}

	got, err := newRolesUseCase(repo).GetRoles(context.Background(), domain.RoleFilters{})

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, domain.RoleDTO{
		ID: 7, Name: "Vendedor", Description: "vende", Level: 3, IsSystem: false,
		ScopeID: 2, ScopeName: "Negocio", ScopeCode: "business",
		BusinessTypeID: 4, BusinessTypeName: "Retail",
	}, got[0])
}

func TestGetRolesByLevel_SinLevelUsaCero(t *testing.T) {
	var visto int
	repo := &mocks.RoleRepositoryMock{
		GetRolesByLevelFn: func(ctx context.Context, level int) ([]domain.Role, error) {
			visto = level
			return nil, nil
		},
	}

	_, err := newRolesUseCase(repo).GetRolesByLevel(context.Background(), domain.RoleFilters{})

	require.NoError(t, err)
	assert.Equal(t, 0, visto)
}

func TestGetRolesByLevel_PasaElLevelPedido(t *testing.T) {
	var visto int
	repo := &mocks.RoleRepositoryMock{
		GetRolesByLevelFn: func(ctx context.Context, level int) ([]domain.Role, error) {
			visto = level
			return []domain.Role{{ID: 3, Level: level}}, nil
		},
	}

	got, err := newRolesUseCase(repo).GetRolesByLevel(context.Background(), domain.RoleFilters{Level: intPtr(3)})

	require.NoError(t, err)
	assert.Equal(t, 3, visto)
	assert.Equal(t, []uint{3}, idsDe(got))
}

func TestGetRolesByLevel_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("boom")
	repo := &mocks.RoleRepositoryMock{
		GetRolesByLevelFn: func(ctx context.Context, level int) ([]domain.Role, error) { return nil, dbErr },
	}

	got, err := newRolesUseCase(repo).GetRolesByLevel(context.Background(), domain.RoleFilters{Level: intPtr(1)})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetRolesByScopeID_PasaElScopeYMapea(t *testing.T) {
	var visto uint
	repo := &mocks.RoleRepositoryMock{
		GetRolesByScopeIDFn: func(ctx context.Context, scopeID uint) ([]domain.Role, error) {
			visto = scopeID
			return []domain.Role{{ID: 2, ScopeID: scopeID, ScopeCode: "business"}}, nil
		},
	}

	got, err := newRolesUseCase(repo).GetRolesByScopeID(context.Background(), 2)

	require.NoError(t, err)
	assert.Equal(t, uint(2), visto)
	require.Len(t, got, 1)
	assert.Equal(t, "business", got[0].ScopeCode)
}

func TestGetRolesByScopeID_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("boom")
	repo := &mocks.RoleRepositoryMock{
		GetRolesByScopeIDFn: func(ctx context.Context, scopeID uint) ([]domain.Role, error) { return nil, dbErr },
	}

	_, err := newRolesUseCase(repo).GetRolesByScopeID(context.Background(), 1)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetSystemRoles_MapeaLoQueDevuelveElRepo(t *testing.T) {
	repo := &mocks.RoleRepositoryMock{
		GetSystemRolesFn: func(ctx context.Context) ([]domain.Role, error) {
			return []domain.Role{{ID: 1, Name: "Super Admin", IsSystem: true}}, nil
		},
	}

	got, err := newRolesUseCase(repo).GetSystemRoles(context.Background())

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.True(t, got[0].IsSystem)
}

func TestGetSystemRoles_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("boom")
	repo := &mocks.RoleRepositoryMock{
		GetSystemRolesFn: func(ctx context.Context) ([]domain.Role, error) { return nil, dbErr },
	}

	_, err := newRolesUseCase(repo).GetSystemRoles(context.Background())

	assert.ErrorIs(t, err, dbErr)
}

func TestGetRoleByID_RolInexistente_RetornaError(t *testing.T) {
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) { return nil, nil },
	}

	got, err := newRolesUseCase(repo).GetRoleByID(context.Background(), 404)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Equal(t, "rol no encontrado", err.Error())
}

func TestGetRoleByID_ErrorDelRepo_NoFiltraElDetalleInterno(t *testing.T) {
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return nil, errors.New("pq: relation roles does not exist")
		},
	}

	_, err := newRolesUseCase(repo).GetRoleByID(context.Background(), 1)

	require.Error(t, err)
	assert.Equal(t, "rol no encontrado", err.Error())
	assert.NotContains(t, err.Error(), "pq:")
}

func TestGetRoleByID_Exitoso_MapeaADTO(t *testing.T) {
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, Name: "Vendedor", Level: 3, ScopeCode: "business"}, nil
		},
	}

	got, err := newRolesUseCase(repo).GetRoleByID(context.Background(), 3)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uint(3), got.ID)
	assert.Equal(t, "Vendedor", got.Name)
	assert.Equal(t, "business", got.ScopeCode)
}

func TestCreateRole_DelegaAlRepositorio(t *testing.T) {
	var recibido domain.CreateRoleDTO
	repo := &mocks.RoleRepositoryMock{
		CreateRoleFn: func(ctx context.Context, role domain.CreateRoleDTO) (*domain.Role, error) {
			recibido = role
			return &domain.Role{ID: 10, Name: role.Name}, nil
		},
	}
	entrada := domain.CreateRoleDTO{Name: "Auditor", Description: "solo lectura", Level: 4, ScopeID: 2, BusinessTypeID: 1}

	got, err := newRolesUseCase(repo).CreateRole(context.Background(), entrada)

	require.NoError(t, err)
	assert.Equal(t, entrada, recibido)
	assert.Equal(t, uint(10), got.ID)
}

func TestCreateRole_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("duplicate key")
	repo := &mocks.RoleRepositoryMock{
		CreateRoleFn: func(ctx context.Context, role domain.CreateRoleDTO) (*domain.Role, error) {
			return nil, dbErr
		},
	}

	got, err := newRolesUseCase(repo).CreateRole(context.Background(), domain.CreateRoleDTO{Name: "X"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateRole_RolInexistente_NoActualiza(t *testing.T) {
	actualizado := false
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) { return nil, nil },
		UpdateRoleFn: func(ctx context.Context, id uint, role domain.UpdateRoleDTO) (*domain.Role, error) {
			actualizado = true
			return nil, nil
		},
	}

	_, err := newRolesUseCase(repo).UpdateRole(context.Background(), 404, domain.UpdateRoleDTO{})

	require.Error(t, err)
	assert.Equal(t, "rol no encontrado", err.Error())
	assert.False(t, actualizado)
}

func TestUpdateRole_ErrorAlBuscarRol_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) { return nil, dbErr },
	}

	_, err := newRolesUseCase(repo).UpdateRole(context.Background(), 1, domain.UpdateRoleDTO{})

	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateRole_NombreNuevoDuplicado_Rechaza(t *testing.T) {
	actualizado := false
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, Name: "Viejo"}, nil
		},
		RoleExistsByNameFn: func(ctx context.Context, name string, excludeID *uint) (bool, error) {
			return true, nil
		},
		UpdateRoleFn: func(ctx context.Context, id uint, role domain.UpdateRoleDTO) (*domain.Role, error) {
			actualizado = true
			return nil, nil
		},
	}

	_, err := newRolesUseCase(repo).UpdateRole(context.Background(), 1, domain.UpdateRoleDTO{Name: strPtr("Nuevo")})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ya existe un rol con el nombre 'Nuevo'")
	assert.False(t, actualizado)
}

func TestUpdateRole_ExcluyeElPropioIDAlBuscarDuplicados(t *testing.T) {
	var excluido *uint
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, Name: "Viejo"}, nil
		},
		RoleExistsByNameFn: func(ctx context.Context, name string, excludeID *uint) (bool, error) {
			excluido = excludeID
			return false, nil
		},
	}

	_, err := newRolesUseCase(repo).UpdateRole(context.Background(), 42, domain.UpdateRoleDTO{Name: strPtr("Nuevo")})

	require.NoError(t, err)
	require.NotNil(t, excluido)
	assert.Equal(t, uint(42), *excluido)
}

func TestUpdateRole_MismoNombreNoDisparaChequeoDeDuplicados(t *testing.T) {
	chequeado := false
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, Name: "Igual"}, nil
		},
		RoleExistsByNameFn: func(ctx context.Context, name string, excludeID *uint) (bool, error) {
			chequeado = true
			return false, nil
		},
	}

	_, err := newRolesUseCase(repo).UpdateRole(context.Background(), 1, domain.UpdateRoleDTO{Name: strPtr("Igual")})

	require.NoError(t, err)
	assert.False(t, chequeado)
}

func TestUpdateRole_SinNombreNoDisparaChequeoDeDuplicados(t *testing.T) {
	chequeado := false
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, Name: "Igual"}, nil
		},
		RoleExistsByNameFn: func(ctx context.Context, name string, excludeID *uint) (bool, error) {
			chequeado = true
			return false, nil
		},
	}

	_, err := newRolesUseCase(repo).UpdateRole(context.Background(), 1, domain.UpdateRoleDTO{
		Description: strPtr("otra descripcion"),
	})

	require.NoError(t, err)
	assert.False(t, chequeado)
}

func TestUpdateRole_ErrorAlChequearDuplicados_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, Name: "Viejo"}, nil
		},
		RoleExistsByNameFn: func(ctx context.Context, name string, excludeID *uint) (bool, error) {
			return false, dbErr
		},
	}

	_, err := newRolesUseCase(repo).UpdateRole(context.Background(), 1, domain.UpdateRoleDTO{Name: strPtr("Nuevo")})

	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "error verificando existencia de rol")
}

func TestUpdateRole_ErrorDelRepoAlActualizar_SePropaga(t *testing.T) {
	dbErr := errors.New("constraint violation")
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id}, nil
		},
		UpdateRoleFn: func(ctx context.Context, id uint, role domain.UpdateRoleDTO) (*domain.Role, error) {
			return nil, dbErr
		},
	}

	got, err := newRolesUseCase(repo).UpdateRole(context.Background(), 1, domain.UpdateRoleDTO{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateRole_Exitoso_ReenviaElDTOCompleto(t *testing.T) {
	var recibido domain.UpdateRoleDTO
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: id, Name: "Viejo"}, nil
		},
		UpdateRoleFn: func(ctx context.Context, id uint, role domain.UpdateRoleDTO) (*domain.Role, error) {
			recibido = role
			return &domain.Role{ID: id, Name: *role.Name}, nil
		},
	}
	entrada := domain.UpdateRoleDTO{Name: strPtr("Nuevo"), Level: intPtr(5), IsSystem: boolPtr(true)}

	got, err := newRolesUseCase(repo).UpdateRole(context.Background(), 1, entrada)

	require.NoError(t, err)
	assert.Equal(t, entrada, recibido)
	assert.Equal(t, "Nuevo", got.Name)
}

func TestAssignPermissionsToRole_RolInexistente_NoAsigna(t *testing.T) {
	asignado := false
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) { return nil, nil },
		AssignPermissionsToRoleFn: func(ctx context.Context, roleID uint, permissionIDs []uint) error {
			asignado = true
			return nil
		},
	}

	err := newRolesUseCase(repo).AssignPermissionsToRole(context.Background(), 404, []uint{1, 2})

	require.Error(t, err)
	assert.Equal(t, "rol no encontrado", err.Error())
	assert.False(t, asignado)
}

func TestAssignPermissionsToRole_ErrorAlBuscarRol_NoAsigna(t *testing.T) {
	asignado := false
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return nil, errors.New("timeout")
		},
		AssignPermissionsToRoleFn: func(ctx context.Context, roleID uint, permissionIDs []uint) error {
			asignado = true
			return nil
		},
	}

	err := newRolesUseCase(repo).AssignPermissionsToRole(context.Background(), 1, []uint{1})

	require.Error(t, err)
	assert.Equal(t, "rol no encontrado", err.Error())
	assert.False(t, asignado)
}

func TestAssignPermissionsToRole_ListaVaciaLlegaAlRepositorio(t *testing.T) {
	var recibidos []uint
	llamado := false
	repo := &mocks.RoleRepositoryMock{
		AssignPermissionsToRoleFn: func(ctx context.Context, roleID uint, permissionIDs []uint) error {
			llamado = true
			recibidos = permissionIDs
			return nil
		},
	}

	err := newRolesUseCase(repo).AssignPermissionsToRole(context.Background(), 1, []uint{})

	require.NoError(t, err)
	assert.True(t, llamado, "una lista vacia es un reemplazo valido: quitar todos los permisos")
	assert.Empty(t, recibidos)
}

func TestAssignPermissionsToRole_PasaLosIDsTalCual(t *testing.T) {
	var recibidoRol uint
	var recibidos []uint
	repo := &mocks.RoleRepositoryMock{
		AssignPermissionsToRoleFn: func(ctx context.Context, roleID uint, permissionIDs []uint) error {
			recibidoRol = roleID
			recibidos = permissionIDs
			return nil
		},
	}

	err := newRolesUseCase(repo).AssignPermissionsToRole(context.Background(), 6, []uint{10, 20, 30})

	require.NoError(t, err)
	assert.Equal(t, uint(6), recibidoRol)
	assert.Equal(t, []uint{10, 20, 30}, recibidos)
}

func TestAssignPermissionsToRole_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("fk violation")
	repo := &mocks.RoleRepositoryMock{
		AssignPermissionsToRoleFn: func(ctx context.Context, roleID uint, permissionIDs []uint) error {
			return dbErr
		},
	}

	err := newRolesUseCase(repo).AssignPermissionsToRole(context.Background(), 1, []uint{1})

	assert.ErrorIs(t, err, dbErr)
}

func TestRemovePermissionFromRole_RolInexistente_NoElimina(t *testing.T) {
	eliminado := false
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) { return nil, nil },
		RemovePermissionFromRoleFn: func(ctx context.Context, roleID uint, permissionID uint) error {
			eliminado = true
			return nil
		},
	}

	err := newRolesUseCase(repo).RemovePermissionFromRole(context.Background(), 404, 1)

	require.Error(t, err)
	assert.Equal(t, "rol no encontrado", err.Error())
	assert.False(t, eliminado)
}

func TestRemovePermissionFromRole_ErrorAlBuscarRol_NoElimina(t *testing.T) {
	eliminado := false
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return nil, errors.New("timeout")
		},
		RemovePermissionFromRoleFn: func(ctx context.Context, roleID uint, permissionID uint) error {
			eliminado = true
			return nil
		},
	}

	err := newRolesUseCase(repo).RemovePermissionFromRole(context.Background(), 1, 1)

	require.Error(t, err)
	assert.False(t, eliminado)
}

func TestRemovePermissionFromRole_PasaAmbosIDs(t *testing.T) {
	var rol, permiso uint
	repo := &mocks.RoleRepositoryMock{
		RemovePermissionFromRoleFn: func(ctx context.Context, roleID uint, permissionID uint) error {
			rol, permiso = roleID, permissionID
			return nil
		},
	}

	err := newRolesUseCase(repo).RemovePermissionFromRole(context.Background(), 6, 99)

	require.NoError(t, err)
	assert.Equal(t, uint(6), rol)
	assert.Equal(t, uint(99), permiso)
}

func TestRemovePermissionFromRole_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("no existe la relacion")
	repo := &mocks.RoleRepositoryMock{
		RemovePermissionFromRoleFn: func(ctx context.Context, roleID uint, permissionID uint) error {
			return dbErr
		},
	}

	err := newRolesUseCase(repo).RemovePermissionFromRole(context.Background(), 1, 1)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetRolePermissions_RolInexistente_NoConsultaPermisos(t *testing.T) {
	consultado := false
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) { return nil, nil },
		GetRolePermissionsFn: func(ctx context.Context, roleID uint) ([]domain.Permission, error) {
			consultado = true
			return nil, nil
		},
	}

	got, err := newRolesUseCase(repo).GetRolePermissions(context.Background(), 404)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.False(t, consultado)
}

func TestGetRolePermissions_ErrorAlBuscarRol_NoConsultaPermisos(t *testing.T) {
	consultado := false
	repo := &mocks.RoleRepositoryMock{
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return nil, errors.New("timeout")
		},
		GetRolePermissionsFn: func(ctx context.Context, roleID uint) ([]domain.Permission, error) {
			consultado = true
			return nil, nil
		},
	}

	_, err := newRolesUseCase(repo).GetRolePermissions(context.Background(), 1)

	require.Error(t, err)
	assert.False(t, consultado)
}

func TestGetRolePermissions_DevuelveLoQueEntregaElRepo(t *testing.T) {
	esperados := []domain.Permission{
		{ID: 1, Name: "orders.read", Resource: "orders", Action: "read"},
		{ID: 2, Name: "orders.write", Resource: "orders", Action: "write"},
	}
	repo := &mocks.RoleRepositoryMock{
		GetRolePermissionsFn: func(ctx context.Context, roleID uint) ([]domain.Permission, error) {
			return esperados, nil
		},
	}

	got, err := newRolesUseCase(repo).GetRolePermissions(context.Background(), 6)

	require.NoError(t, err)
	assert.Equal(t, esperados, got)
}

func TestGetRolePermissions_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("join fallido")
	repo := &mocks.RoleRepositoryMock{
		GetRolePermissionsFn: func(ctx context.Context, roleID uint) ([]domain.Permission, error) {
			return nil, dbErr
		},
	}

	got, err := newRolesUseCase(repo).GetRolePermissions(context.Background(), 1)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}
