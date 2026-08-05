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

func TestCreateUser_ErrorGenericoAlAsignarBusinesses_NoSeConfundeConInexistentes(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) { return 12, nil },
		AssignBusinessesToUserFn: func(ctx context.Context, userID uint, businessIDs []uint) error {
			return errors.New("deadlock detected")
		},
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{
		Email:       "a@b.com",
		BusinessIDs: []uint{3},
	})

	require.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrBusinessesNotFound)
	assert.Contains(t, err.Error(), "error al asignar businesses")
}

func TestCreateUser_ErrorAlBuscarRolOperator_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout en roles")
	asignado := false
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) { return 20, nil },
		GetRoleIDByNameAndScopeFn: func(ctx context.Context, roleName string, scopeID uint) (uint, error) {
			return 0, dbErr
		},
		AssignRolesToUserFn: func(ctx context.Context, userID uint, roleIDs []uint) error {
			asignado = true
			return nil
		},
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{
		Email:   "op@b.com",
		ScopeID: uintPtr(3),
	})

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, asignado)
}

func TestCreateUser_ErrorAlAsignarRolOperator_SePropaga(t *testing.T) {
	repoErr := errors.New("violacion de fk en user_roles")
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) { return 20, nil },
		GetRoleIDByNameAndScopeFn: func(ctx context.Context, roleName string, scopeID uint) (uint, error) {
			return 77, nil
		},
		AssignRolesToUserFn: func(ctx context.Context, userID uint, roleIDs []uint) error { return repoErr },
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{
		Email:   "op@b.com",
		ScopeID: uintPtr(3),
	})

	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "error al asignar rol operator")
}

func TestGetUserByID_ErrorAlTraerRoles_NoAbortaElDetalle(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, Name: "Ana"}, nil
		},
		GetUserRolesFn: func(ctx context.Context, userID uint) ([]domain.Role, error) {
			return nil, errors.New("join fallido")
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUserByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, "Ana", got.Name)
	assert.Empty(t, got.Roles)
	assert.False(t, got.IsSuperUser)
}

func TestGetUserByID_ErrorAlTraerBusinesses_DejaListaVaciaSinFallar(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, Name: "Ana"}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return nil, errors.New("timeout")
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUserByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Empty(t, got.Businesses)
	assert.Empty(t, got.BusinessRoleAssignments)
}

func TestGetUserByID_NavbarRelativoSeCompletaConElDominioS3(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1}, nil
		},
		GetUserBusinessesFn: func(ctx context.Context, userID uint) ([]domain.BusinessInfoEntity, error) {
			return []domain.BusinessInfoEntity{
				{ID: 26, Name: "Demo", NavbarImageURL: "banners/demo.png"},
				{ID: 27, Name: "Otro", NavbarImageURL: "https://cdn.externo.com/x.png"},
			}, nil
		},
	}

	got, err := newUseCase(repo, nil, s3Config()).GetUserByID(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, got.Businesses, 2)
	assert.Equal(t, s3Base+"/banners/demo.png", got.Businesses[0].NavbarImageURL)
	assert.Equal(t, "https://cdn.externo.com/x.png", got.Businesses[1].NavbarImageURL)
}

func TestGetUsers_ErrorAlTraerRelacionesStaff_NoAbortaElListado(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUsersFn: func(ctx context.Context, f domain.UserFilters) ([]domain.UserQueryDTO, int64, error) {
			return []domain.UserQueryDTO{{ID: 1, Name: "Ana"}}, 1, nil
		},
		GetBusinessStaffRelationshipsFn: func(ctx context.Context, userID uint) ([]domain.BusinessRoleAssignmentDetailed, error) {
			return nil, errors.New("timeout en business_staff")
		},
	}

	got, err := newUseCase(repo, nil, nil).GetUsers(context.Background(), domain.UserFilters{})

	require.NoError(t, err)
	require.Len(t, got.Users, 1)
	assert.Equal(t, "Ana", got.Users[0].Name)
	assert.Empty(t, got.Users[0].BusinessRoleAssignments)
}

func TestUpdateUser_ErrorAlBuscarUsuario_SePropaga(t *testing.T) {
	dbErr := errors.New("connection reset")
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) { return nil, dbErr },
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{})

	assert.ErrorIs(t, err, dbErr)
	assert.Contains(t, err.Error(), "error al buscar usuario")
}

func TestUpdateUser_ErrorAlVerificarEmailDuplicado_SePropaga(t *testing.T) {
	dbErr := errors.New("statement timeout")
	actualizado := false
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, Email: "viejo@x.com"}, nil
		},
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) { return nil, dbErr },
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
			actualizado = true
			return "", nil
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{Email: "nuevo@x.com"})

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, actualizado)
}

func TestUpdateUser_FallaSubidaDeAvatar_NoActualiza(t *testing.T) {
	actualizado := false
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "", errors.New("s3 timeout")
		},
	}
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
			actualizado = true
			return "", nil
		},
	}

	_, err := newUseCase(repo, s3, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{
		AvatarFile: &multipart.FileHeader{Filename: "x.jpg"},
	})

	assert.ErrorIs(t, err, domain.ErrUserAvatarUploadFailed)
	assert.False(t, actualizado)
}

func TestUpdateUser_RemoveAvatarConUrlExterna_NoLlamaAS3PeroLimpiaElCampo(t *testing.T) {
	s3 := &mocks.S3ServiceMock{}
	var guardado domain.UsersEntity
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, AvatarURL: "https://lh3.googleusercontent.com/a/foto.jpg"}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
			guardado = user
			return "ok", nil
		},
	}

	_, err := newUseCase(repo, s3, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{RemoveAvatar: true})

	require.NoError(t, err)
	assert.Empty(t, s3.DeleteImageCall)
	assert.Empty(t, guardado.AvatarURL)
}

func TestUpdateUser_FalloAlBorrarAvatarEnRemove_NoRompeLaActualizacion(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
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

	got, err := newUseCase(repo, s3, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{RemoveAvatar: true})

	require.NoError(t, err)
	assert.Equal(t, "ok", got)
}

func TestUpdateUser_ErrorGenericoDelRepo_NoSeTraduceADuplicado(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) {
			return "", errors.New("deadlock detected")
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{})

	require.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrUserEmailExists)
	assert.Contains(t, err.Error(), "error al actualizar usuario")
}

func TestUpdateUser_ErrorGenericoAlAsignarBusinesses_NoSeConfundeConInexistentes(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) { return "ok", nil },
		AssignBusinessesToUserFn: func(ctx context.Context, userID uint, businessIDs []uint) error {
			return errors.New("deadlock detected")
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{BusinessIDs: []uint{3}})

	require.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrBusinessesNotFound)
	assert.Contains(t, err.Error(), "error al asignar businesses")
}

func TestUpdateUser_BusinessIDsSeReenvianTalCual(t *testing.T) {
	var recibidos []uint
	repo := &mocks.UserRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1}, nil
		},
		UpdateUserFn: func(ctx context.Context, id uint, user domain.UsersEntity) (string, error) { return "ok", nil },
		AssignBusinessesToUserFn: func(ctx context.Context, userID uint, businessIDs []uint) error {
			recibidos = businessIDs
			return nil
		},
	}

	_, err := newUseCase(repo, nil, nil).UpdateUser(context.Background(), 1, domain.UpdateUserDTO{BusinessIDs: []uint{26, 27}})

	require.NoError(t, err)
	assert.Equal(t, []uint{26, 27}, recibidos)
}
