package app

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const s3Base = "https://cdn.probability.co"

func s3Config() map[string]string {
	return map[string]string{"URL_BASE_DOMAIN_S3": s3Base}
}

func TestGetUsers_NormalizaPaginacionFueraDeRango(t *testing.T) {
	casos := []struct {
		nombre       string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{"page cero sube a uno", 0, 10, 1, 10},
		{"page negativo sube a uno", -5, 10, 1, 10},
		{"pageSize cero cae al default", 1, 0, 1, 10},
		{"pageSize negativo cae al default", 1, -3, 1, 10},
		{"pageSize sobre el maximo se topa en 100", 1, 5000, 1, 100},
		{"valores validos se respetan", 4, 25, 4, 25},
		{"el maximo exacto se respeta", 1, 100, 1, 100},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			var visto domain.UserFilters
			repo := &mocks.UserRepositoryMock{
				GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
					visto = f
					return nil, 0, nil
				},
			}

			got, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{
				Page:     tc.page,
				PageSize: tc.pageSize,
			})

			require.NoError(t, err)
			assert.Equal(t, tc.wantPage, visto.Page)
			assert.Equal(t, tc.wantPageSize, visto.PageSize)
			assert.Equal(t, tc.wantPage, got.Page)
			assert.Equal(t, tc.wantPageSize, got.PageSize)
		})
	}
}

func TestGetUsers_SortByNoPermitidoCaeACreatedAt(t *testing.T) {
	casos := []struct {
		entrada string
		want    string
	}{
		{"id", "id"},
		{"name", "name"},
		{"email", "email"},
		{"phone", "phone"},
		{"is_active", "is_active"},
		{"created_at", "created_at"},
		{"updated_at", "updated_at"},
		{"password", "created_at"},
		{"users.password", "created_at"},
		{"id; DROP TABLE users", "created_at"},
		{"(SELECT 1)", "created_at"},
	}

	for _, tc := range casos {
		t.Run(tc.entrada, func(t *testing.T) {
			var visto domain.UserFilters
			repo := &mocks.UserRepositoryMock{
				GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
					visto = f
					return nil, 0, nil
				},
			}

			_, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{SortBy: tc.entrada})

			require.NoError(t, err)
			assert.Equal(t, tc.want, visto.SortBy)
		})
	}
}

func TestGetUsers_SortOrderInvalidoCaeADesc(t *testing.T) {
	casos := []struct {
		entrada string
		want    string
	}{
		{"asc", "asc"},
		{"desc", "desc"},
		{"ASC", "desc"},
		{"ascending", "desc"},
		{"random()", "desc"},
	}

	for _, tc := range casos {
		t.Run(tc.entrada, func(t *testing.T) {
			var visto domain.UserFilters
			repo := &mocks.UserRepositoryMock{
				GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
					visto = f
					return nil, 0, nil
				},
			}

			_, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{SortOrder: tc.entrada})

			require.NoError(t, err)
			assert.Equal(t, tc.want, visto.SortOrder)
		})
	}
}

func TestGetUsers_CalculaTotalPagesConRedondeoHaciaArriba(t *testing.T) {
	casos := []struct {
		total    int64
		pageSize int
		want     int
	}{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{99, 10, 10},
		{100, 10, 10},
		{101, 10, 11},
	}

	for _, tc := range casos {
		repo := &mocks.UserRepositoryMock{
			GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
				return nil, tc.total, nil
			},
		}

		got, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.want, got.TotalPages, "total=%d pageSize=%d", tc.total, tc.pageSize)
	}
}

func TestGetUsers_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("statement timeout")
	repo := &mocks.UserRepositoryMock{
		GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
			return nil, 0, dbErr
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetUsers_AvatarRelativoSeCompletaConElDominioS3(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
			return []domain.UserQueryDTO{{ID: 1, AvatarURL: "avatars/foto.jpg"}}, 1, nil
		},
	}

	got, err := newUseCase(repo, nil, s3Config()).GetUsers(context.Background(), domain.UserFilters{})

	require.NoError(t, err)
	require.Len(t, got.Users, 1)
	assert.Equal(t, s3Base+"/avatars/foto.jpg", got.Users[0].AvatarURL)
}

func TestGetUsers_AvatarAbsolutoNoSeToca(t *testing.T) {
	externa := "https://lh3.googleusercontent.com/a/foto.jpg"
	repo := &mocks.UserRepositoryMock{
		GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
			return []domain.UserQueryDTO{{ID: 1, AvatarURL: externa}}, 1, nil
		},
	}

	got, err := newUseCase(repo, nil, s3Config()).GetUsers(context.Background(), domain.UserFilters{})

	require.NoError(t, err)
	assert.Equal(t, externa, got.Users[0].AvatarURL)
}

func TestGetUsers_SinDominioS3ConfiguradoDejaElPathRelativo(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
			return []domain.UserQueryDTO{{ID: 1, AvatarURL: "avatars/foto.jpg"}}, 1, nil
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{})

	require.NoError(t, err)
	assert.Equal(t, "avatars/foto.jpg", got.Users[0].AvatarURL)
}

func TestGetUsers_NoDuplicaBarrasAlUnirDominioYPath(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
			return []domain.UserQueryDTO{{ID: 1, AvatarURL: "/avatars/foto.jpg"}}, 1, nil
		},
	}
	cfg := map[string]string{"URL_BASE_DOMAIN_S3": s3Base + "/"}

	got, err := newUseCase(repo, nil, cfg).GetUsers(context.Background(), domain.UserFilters{})

	require.NoError(t, err)
	assert.Equal(t, s3Base+"/avatars/foto.jpg", got.Users[0].AvatarURL)
	assert.NotContains(t, got.Users[0].AvatarURL, "//avatars")
}

func TestGetUsers_RolConScopePlatform_MarcaSuperUsuario(t *testing.T) {
	casos := []struct {
		nombre    string
		rol       domain.Role
		wantSuper bool
	}{
		{"scope_id 1", domain.Role{ID: 1, Name: "Super Admin", ScopeID: 1}, true},
		{"scope_code platform", domain.Role{ID: 2, Name: "Super Admin", ScopeCode: "platform"}, true},
		{"scope business", domain.Role{ID: 3, Name: "Admin", ScopeID: 2, ScopeCode: "business"}, false},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := &mocks.UserRepositoryMock{
				GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
					return []domain.UserQueryDTO{{ID: 1}}, 1, nil
				},
				GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
					return []domain.Role{tc.rol}, nil
				},
			}

			got, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{})

			require.NoError(t, err)
			assert.Equal(t, tc.wantSuper, got.Users[0].IsSuperUser)
		})
	}
}

func TestGetUsers_SuperUsuario_RecibeAssignmentConBusinessCero(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
			return []domain.UserQueryDTO{{ID: 1}}, 1, nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{{ID: 9, Name: "Super Admin", ScopeCode: "platform"}}, nil
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{})

	require.NoError(t, err)
	require.NotEmpty(t, got.Users[0].BusinessRoleAssignments)
	primero := got.Users[0].BusinessRoleAssignments[0]
	assert.Equal(t, uint(0), primero.BusinessID)
	assert.Equal(t, uint(9), primero.RoleID)
	assert.Equal(t, "Super Admin", primero.RoleName)
	assert.Empty(t, primero.BusinessName)
}

func TestGetUsers_RelacionConBusinessCero_SeDescarta(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
			return []domain.UserQueryDTO{{ID: 1}}, 1, nil
		},
		GetBusinessStaffRelationshipsFn: func(ctx context.Context, userID uint) ([]domain.BusinessRoleAssignmentDetailed, error) {
			return []domain.BusinessRoleAssignmentDetailed{
				{BusinessID: 0, RoleID: 5},
				{BusinessID: 26, RoleID: 6, BusinessName: "Demo"},
			}, nil
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{})

	require.NoError(t, err)
	require.Len(t, got.Users[0].BusinessRoleAssignments, 1)
	assert.Equal(t, uint(26), got.Users[0].BusinessRoleAssignments[0].BusinessID)
}

func TestGetUsers_ErrorAlTraerRolesNoAbortaElListado(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
			return []domain.UserQueryDTO{{ID: 1, Name: "Ana"}}, 1, nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return nil, errors.New("join fallido")
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{})

	require.NoError(t, err)
	require.Len(t, got.Users, 1)
	assert.Equal(t, "Ana", got.Users[0].Name)
	assert.Empty(t, got.Users[0].Roles)
}

func TestGetUsers_BusinessDelUsuarioTraeElRolDeLaRelacion(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
			return []domain.UserQueryDTO{{ID: 1}}, 1, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return []domain.BusinessInfoEntity{{ID: 26, Name: "Demo", NavbarImageURL: "banners/x.png"}}, nil
		},
		GetBusinessStaffRelationshipsFn: func(ctx context.Context, userID uint) ([]domain.BusinessRoleAssignmentDetailed, error) {
			return []domain.BusinessRoleAssignmentDetailed{{BusinessID: 26, RoleID: 6, RoleName: "Admin"}}, nil
		},
		GetRoleByIDFn: func(ctx context.Context, id uint) (*domain.Role, error) {
			return &domain.Role{ID: 6, Name: "Admin", Level: 2}, nil
		},
	}

	got, err := newUseCase(repo, nil, s3Config()).GetUsers(context.Background(), domain.UserFilters{})

	require.NoError(t, err)
	require.Len(t, got.Users[0].Businesses, 1)
	b := got.Users[0].Businesses[0]
	assert.Equal(t, s3Base+"/banners/x.png", b.NavbarImageURL)
	require.NotNil(t, b.Role)
	assert.Equal(t, uint(6), b.Role.ID)
	assert.Equal(t, 2, b.Role.Level)
}

func TestGetUserByID_UsuarioInexistente_RetornaError(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) { return nil, nil },
	}

	got, err := newUseCase(repo, nil, nil).GetUserByID(context.Background(), 404)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "usuario no encontrado")
}

func TestGetUserByID_ErrorDelRepo_NoFiltraElDetalleInterno(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return nil, errors.New("pq: relation users does not exist")
		},
	}

	_, err := newUseCase(repo, nil, nil).GetUserByID(context.Background(), 1)

	require.Error(t, err)
	assert.Equal(t, "usuario no encontrado", err.Error())
	assert.NotContains(t, err.Error(), "pq:")
}

func TestGetUserByID_MapeaDatosBasicosYCompletaAvatar(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{
				ID: 42, Name: "Ana", Email: "ana@x.com", Phone: "3001234567",
				AvatarURL: "avatars/ana.jpg", IsActive: true,
			}, nil
		},
	}

	got, err := newUseCase(repo, nil, s3Config()).GetUserByID(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, uint(42), got.ID)
	assert.Equal(t, "Ana", got.Name)
	assert.Equal(t, "3001234567", got.Phone)
	assert.True(t, got.IsActive)
	assert.Equal(t, s3Base+"/avatars/ana.jpg", got.AvatarURL)
}

func TestGetUserByID_NoExponeElHashDeContrasena(t *testing.T) {
	const hash = "$2a$10$abcdefghijklmnopqrstuv"
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, Email: "a@b.com", Password: hash}, nil
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUserByID(context.Background(), 1)

	require.NoError(t, err)
	assert.NotContains(t, got.Email, hash)
	assert.NotContains(t, got.Name, hash)
	assert.NotContains(t, got.AvatarURL, hash)
}

func TestGetUserByID_ConstruyeAssignmentsDesdeSusBusinesses(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return []domain.BusinessInfoEntity{
				{ID: 26, Name: "Demo"},
				{ID: 27, Name: "Otro"},
			}, nil
		},
		GetUserRoleByBusinessFn: func(ctx context.Context, userID uint, businessID uint) (*domain.Role, error) {
			if businessID == 26 {
				return &domain.Role{ID: 6, Name: "Admin"}, nil
			}
			return nil, errors.New("sin rol")
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUserByID(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, got.BusinessRoleAssignments, 2)
	assert.Equal(t, uint(26), got.BusinessRoleAssignments[0].BusinessID)
	assert.Equal(t, uint(6), got.BusinessRoleAssignments[0].RoleID)
	assert.Equal(t, "Admin", got.BusinessRoleAssignments[0].RoleName)
	assert.Equal(t, uint(27), got.BusinessRoleAssignments[1].BusinessID)
	assert.Equal(t, uint(0), got.BusinessRoleAssignments[1].RoleID)
	assert.Empty(t, got.BusinessRoleAssignments[1].RoleName)
}

func TestGetUserByID_SuperUsuario_PoneElAssignmentCeroAlInicio(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1}, nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return []domain.Role{{ID: 1, Name: "Super Admin", ScopeID: 1}}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return []domain.BusinessInfoEntity{{ID: 26, Name: "Demo"}}, nil
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUserByID(context.Background(), 1)

	require.NoError(t, err)
	assert.True(t, got.IsSuperUser)
	require.Len(t, got.BusinessRoleAssignments, 2)
	assert.Equal(t, uint(0), got.BusinessRoleAssignments[0].BusinessID)
	assert.Equal(t, uint(26), got.BusinessRoleAssignments[1].BusinessID)
}

func TestGetUserBusinesses_DelegaAlRepositorio(t *testing.T) {
	esperado := []domain.BusinessInfoEntity{{ID: 26, Name: "Demo"}}
	var vistoUserID uint
	repo := &mocks.UserRepositoryMock{
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			vistoUserID = userID
			return esperado, nil
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUserBusinesses(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, uint(42), vistoUserID)
	assert.Equal(t, esperado, got)
}

func TestGetUserBusinesses_PropagaError(t *testing.T) {
	dbErr := errors.New("boom")
	repo := &mocks.UserRepositoryMock{
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return nil, dbErr
		},
	}

	_, err := newUseCase(repo, nil, nil).GetUserBusinesses(context.Background(), 1)

	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateUser_UsuarioInexistente_NoActualiza(t *testing.T) {
	actualizado := false
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) { return nil, nil },
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
			actualizado = true
			return "", nil
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 404, domain.UpdateUserDTO{})

	assert.ErrorIs(t, err, domain.ErrUserNotFound)
	assert.False(t, actualizado)
}

func TestUpdateUser_EmailDeOtroUsuario_RetornaEmailExists(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, Email: "yo@x.com"}, nil
		},
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 2, Email: email}, nil
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{Email: "otro@x.com"})

	assert.ErrorIs(t, err, domain.ErrUserEmailExists)
}

func TestUpdateUser_MismoEmailDelPropioUsuario_NoConsultaDuplicados(t *testing.T) {
	consultado := false
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, Email: "Yo@X.com"}, nil
		},
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			consultado = true
			return nil, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
			return "ok", nil
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{Email: "  YO@x.COM "})

	require.NoError(t, err)
	assert.False(t, consultado, "el email normalizado coincide con el propio, no hay que buscar duplicados")
}

func TestUpdateUser_EmailPropioConIDDistintoNoBloquea(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, Email: "viejo@x.com"}, nil
		},
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, Email: email}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) { return "ok", nil },
	}

	got, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{Email: "nuevo@x.com"})

	require.NoError(t, err)
	assert.Equal(t, "ok", got)
}

func TestUpdateUser_NuevoAvatar_BorraElAnteriorRelativo(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "avatars/nuevo.jpg", nil
		},
	}
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, AvatarURL: "avatars/viejo.jpg"}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) { return "ok", nil },
	}

	_, err := newUseCase(repo, s3, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{
		AvatarFile: &multipart.FileHeader{Filename: "nuevo.jpg"},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"avatars/viejo.jpg"}, s3.DeleteImageCall)
}

func TestUpdateUser_NuevoAvatar_NoIntentaBorrarUrlExterna(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "avatars/nuevo.jpg", nil
		},
	}
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, AvatarURL: "https://lh3.googleusercontent.com/a/foto.jpg"}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) { return "ok", nil },
	}

	_, err := newUseCase(repo, s3, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{
		AvatarFile: &multipart.FileHeader{Filename: "nuevo.jpg"},
	})

	require.NoError(t, err)
	assert.Empty(t, s3.DeleteImageCall, "una URL externa no es nuestra, no se borra de S3")
}

func TestUpdateUser_FalloAlBorrarAvatarAnterior_NoRompeLaActualizacion(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "avatars/nuevo.jpg", nil
		},
		DeleteImageFn: func(ctx context.Context, filename string) error {
			return errors.New("s3 access denied")
		},
	}
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, AvatarURL: "avatars/viejo.jpg"}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) { return "ok", nil },
	}

	got, err := newUseCase(repo, s3, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{
		AvatarFile: &multipart.FileHeader{Filename: "nuevo.jpg"},
	})

	require.NoError(t, err)
	assert.Equal(t, "ok", got)
}

func TestUpdateUser_RemoveAvatar_LimpiaLaUrlYBorraDeS3(t *testing.T) {
	s3 := &mocks.S3ServiceMock{}
	var guardado domain.UsersEntity
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, AvatarURL: "avatars/viejo.jpg"}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
			guardado = user
			return "ok", nil
		},
	}

	_, err := newUseCase(repo, s3, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{
		AvatarURL:    "avatars/viejo.jpg",
		RemoveAvatar: true,
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"avatars/viejo.jpg"}, s3.DeleteImageCall)
	assert.Empty(t, guardado.AvatarURL)
}

func TestUpdateUser_SinAvatarFileNiRemove_NoTocaS3(t *testing.T) {
	s3 := &mocks.S3ServiceMock{}
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, AvatarURL: "avatars/viejo.jpg"}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) { return "ok", nil },
	}

	_, err := newUseCase(repo, s3, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{Name: "Nuevo Nombre"})

	require.NoError(t, err)
	assert.Empty(t, s3.DeleteImageCall)
}

func TestUpdateUser_NoEnviaContrasenaAlRepositorio(t *testing.T) {
	var guardado domain.UsersEntity
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, Password: "$2a$10$hashviejo"}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
			guardado = user
			return "ok", nil
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{Name: "Ana"})

	require.NoError(t, err)
	assert.Empty(t, guardado.Password, "este endpoint no debe poder cambiar la contrasena")
}

func TestUpdateUser_ErrorDuplicadoDelRepo_SeTraduce(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
			return "", errors.New(`duplicate key value violates unique constraint "users_email_key"`)
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{})

	assert.ErrorIs(t, err, domain.ErrUserEmailExists)
}

func TestUpdateUser_BusinessIDsVacio_NoReasignaRelaciones(t *testing.T) {
	reasignado := false
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) { return "ok", nil },
		AssignBusinessesToUserFn: func(ctx context.Context, userID uint, businessIDs []uint) error {
			reasignado = true
			return nil
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{Name: "Ana"})

	require.NoError(t, err)
	assert.False(t, reasignado, "una lista vacia no debe borrar las relaciones existentes")
}

func TestUpdateUser_BusinessInexistente_SeTraduceAErrorDeDominio(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) { return "ok", nil },
		AssignBusinessesToUserFn: func(ctx context.Context, userID uint, businessIDs []uint) error {
			return errors.New("algunos businesses no existen")
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{BusinessIDs: []uint{99}})

	assert.ErrorIs(t, err, domain.ErrBusinessesNotFound)
}

func TestDeleteUser_UsuarioInexistente_NoLlamaAlRepo(t *testing.T) {
	borrado := false
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) { return nil, nil },
		DeleteUserFn: func(ctx context.Context, id uint) (string, error) {
			borrado = true
			return "", nil
		},
	}

	_, err := newUseCase(repo, nil, nil).DeleteUser(context.Background(), 404)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "usuario no encontrado")
	assert.False(t, borrado)
}

func TestDeleteUser_Exitoso_RetornaMensajeDelRepo(t *testing.T) {
	var borradoID uint
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: userID}, nil
		},
		DeleteUserFn: func(ctx context.Context, id uint) (string, error) {
			borradoID = id
			return "usuario eliminado", nil
		},
	}

	got, err := newUseCase(repo, nil, nil).DeleteUser(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, uint(7), borradoID)
	assert.Equal(t, "usuario eliminado", got)
}

func TestDeleteUser_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("foreign key violation")
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: userID}, nil
		},
		DeleteUserFn: func(ctx context.Context, id uint) (string, error) { return "", dbErr },
	}

	_, err := newUseCase(repo, nil, nil).DeleteUser(context.Background(), 7)

	assert.ErrorIs(t, err, dbErr)
}

func TestAssignRoleToUserBusiness_SinAsignaciones_Rechaza(t *testing.T) {
	llamado := false
	repo := &mocks.UserRepositoryMock{
		AssignRoleToUserBusinessFn: func(ctx context.Context, userID uint, a []domain.BusinessRoleAssignment) error {
			llamado = true
			return nil
		},
	}

	for _, entrada := range [][]domain.BusinessRoleAssignment{nil, {}} {
		err := newUseCase(repo, nil, nil).AssignRoleToUserBusiness(context.Background(), 1, entrada)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no se proporcionaron asignaciones")
	}
	assert.False(t, llamado)
}

func TestAssignRoleToUserBusiness_PasaLasAsignacionesTalCual(t *testing.T) {
	esperadas := []domain.BusinessRoleAssignment{{BusinessID: 26, RoleID: 6}, {BusinessID: 27, RoleID: 7}}
	var recibidas []domain.BusinessRoleAssignment
	var recibidoUserID uint
	repo := &mocks.UserRepositoryMock{
		AssignRoleToUserBusinessFn: func(ctx context.Context, userID uint, a []domain.BusinessRoleAssignment) error {
			recibidoUserID = userID
			recibidas = a
			return nil
		},
	}

	err := newUseCase(repo, nil, nil).AssignRoleToUserBusiness(context.Background(), 42, esperadas)

	require.NoError(t, err)
	assert.Equal(t, uint(42), recibidoUserID)
	assert.Equal(t, esperadas, recibidas)
}

func TestAssignRoleToUserBusiness_ErrorDelRepo_SeEnvuelve(t *testing.T) {
	repoErr := errors.New("el rol no pertenece al tipo de business")
	repo := &mocks.UserRepositoryMock{
		AssignRoleToUserBusinessFn: func(ctx context.Context, userID uint, a []domain.BusinessRoleAssignment) error {
			return repoErr
		},
	}

	err := newUseCase(repo, nil, nil).AssignRoleToUserBusiness(context.Background(), 1,
		[]domain.BusinessRoleAssignment{{BusinessID: 26, RoleID: 6}})

	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "error al asignar roles")
}
