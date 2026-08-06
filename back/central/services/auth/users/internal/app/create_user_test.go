package app

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"
	"testing"
	"unicode"

	"github.com/secamc93/probability/back/central/services/auth/users/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newUseCase(repo *mocks.UserRepositoryMock, s3 *mocks.S3ServiceMock, cfg map[string]string) Iapp {
	if repo == nil {
		repo = &mocks.UserRepositoryMock{}
	}
	if s3 == nil {
		s3 = &mocks.S3ServiceMock{}
	}
	return New(repo, mocks.NewSilentLogger(), s3, mocks.NewConfigMock(cfg))
}

func uintPtr(v uint) *uint { return &v }

func TestCreateUser_EmailSeNormalizaAMinusculasYSinEspacios(t *testing.T) {
	var lookedUp string
	var created domain.UsersEntity
	repo := &mocks.UserRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			lookedUp = email
			return nil, nil
		},
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) {
			created = user
			return 55, nil
		},
	}

	email, password, message, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{
		Name:  "Ana",
		Email: "  ANA.Perez@Example.COM  ",
	})

	require.NoError(t, err)
	assert.Equal(t, "ana.perez@example.com", lookedUp)
	assert.Equal(t, "ana.perez@example.com", created.Email)
	assert.Equal(t, "ana.perez@example.com", email)
	assert.NotEmpty(t, password)
	assert.Contains(t, message, "55")
}

func TestCreateUser_EmailYaExiste_RetornaErrorSinCrear(t *testing.T) {
	creado := false
	repo := &mocks.UserRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return &domain.UserAuthInfo{ID: 1, Email: email}, nil
		},
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) {
			creado = true
			return 1, nil
		},
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{Email: "dup@x.com"})

	assert.ErrorIs(t, err, domain.ErrUserEmailExists)
	assert.False(t, creado, "no debe crear el usuario si el email ya existe")
}

func TestCreateUser_ErrorAlVerificarEmail_NoCrea(t *testing.T) {
	dbErr := errors.New("connection refused")
	creado := false
	repo := &mocks.UserRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return nil, dbErr
		},
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) {
			creado = true
			return 1, nil
		},
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{Email: "x@x.com"})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.False(t, creado)
}

func TestCreateUser_DuplicateKeyDelRepo_SeTraduceAErrorDeDominio(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) {
			return 0, errors.New(`ERROR: duplicate key value violates unique constraint "users_email_key"`)
		},
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{Email: "x@x.com"})

	assert.ErrorIs(t, err, domain.ErrUserEmailExists)
}

func TestCreateUser_ErrorGenericoDelRepo_NoSeTraduceADuplicado(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) {
			return 0, errors.New("deadlock detected")
		},
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{Email: "x@x.com"})

	require.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrUserEmailExists)
	assert.Contains(t, err.Error(), "deadlock")
}

func TestCreateUser_ConAvatarFile_SubeAS3YGuardaPathRelativo(t *testing.T) {
	var folder string
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, f string) (string, error) {
			folder = f
			return "avatars/123_foto.jpg", nil
		},
	}
	var created domain.UsersEntity
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) {
			created = user
			return 7, nil
		},
	}

	_, _, _, err := newUseCase(repo, s3, nil).CreateUser(context.Background(), domain.CreateUserDTO{
		Email:      "a@b.com",
		AvatarFile: &multipart.FileHeader{Filename: "foto.jpg"},
	})

	require.NoError(t, err)
	assert.Equal(t, "avatars", folder)
	assert.Equal(t, "avatars/123_foto.jpg", created.AvatarURL)
}

func TestCreateUser_FallaSubidaAvatar_AbortaSinCrearUsuario(t *testing.T) {
	creado := false
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "", errors.New("s3 timeout")
		},
	}
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) {
			creado = true
			return 1, nil
		},
	}

	_, _, _, err := newUseCase(repo, s3, nil).CreateUser(context.Background(), domain.CreateUserDTO{
		Email:      "a@b.com",
		AvatarFile: &multipart.FileHeader{Filename: "foto.jpg"},
	})

	assert.ErrorIs(t, err, domain.ErrUserAvatarUploadFailed)
	assert.False(t, creado)
}

func TestCreateUser_ConBusinessIDs_LosAsigna(t *testing.T) {
	var asignados []uint
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) { return 12, nil },
		AssignBusinessesToUserFn: func(ctx context.Context, userID uint, businessIDs []uint) error {
			asignados = businessIDs
			return nil
		},
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{
		Email:       "a@b.com",
		BusinessIDs: []uint{3, 9},
	})

	require.NoError(t, err)
	assert.Equal(t, []uint{3, 9}, asignados)
}

func TestCreateUser_BusinessInexistente_SeTraduceAErrorDeDominio(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) { return 12, nil },
		AssignBusinessesToUserFn: func(ctx context.Context, userID uint, businessIDs []uint) error {
			return errors.New("algunos businesses no existen: [99]")
		},
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{
		Email:       "a@b.com",
		BusinessIDs: []uint{99},
	})

	assert.ErrorIs(t, err, domain.ErrBusinessesNotFound)
}

func TestCreateUser_ScopeOperator_AsignaRolOperatorAdminAutomaticamente(t *testing.T) {
	var pedidoRol string
	var pedidoScope uint
	var rolesAsignados []uint
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) { return 20, nil },
		GetRoleIDByNameAndScopeFn: func(ctx context.Context, roleName string, scopeID uint) (uint, error) {
			pedidoRol = roleName
			pedidoScope = scopeID
			return 77, nil
		},
		AssignRolesToUserFn: func(ctx context.Context, userID uint, roleIDs []uint) error {
			rolesAsignados = roleIDs
			return nil
		},
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{
		Email:   "op@b.com",
		ScopeID: uintPtr(3),
	})

	require.NoError(t, err)
	assert.Equal(t, "Operator Admin", pedidoRol)
	assert.Equal(t, uint(3), pedidoScope)
	assert.Equal(t, []uint{77}, rolesAsignados)
}

func TestCreateUser_ScopeNoOperator_NoBuscaRolAutomatico(t *testing.T) {
	buscado := false
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) { return 20, nil },
		GetRoleIDByNameAndScopeFn: func(ctx context.Context, roleName string, scopeID uint) (uint, error) {
			buscado = true
			return 1, nil
		},
	}

	for _, scope := range []*uint{nil, uintPtr(1), uintPtr(2)} {
		_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{
			Email:   "x@b.com",
			ScopeID: scope,
		})
		require.NoError(t, err)
	}

	assert.False(t, buscado, "solo el scope 3 dispara la asignacion automatica de rol")
}

func TestCreateUser_ScopeOperatorSinRolEnBD_Falla(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		CreateUserFn: func(ctx context.Context, user domain.UsersEntity) (uint, error) { return 20, nil },
		GetRoleIDByNameAndScopeFn: func(ctx context.Context, roleName string, scopeID uint) (uint, error) {
			return 0, nil
		},
	}

	_, _, _, err := newUseCase(repo, nil, nil).CreateUser(context.Background(), domain.CreateUserDTO{
		Email:   "op@b.com",
		ScopeID: uintPtr(3),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "operator admin no encontrado")
}

func TestGenerateRandomPassword_LongitudMinimaDeOcho(t *testing.T) {
	for _, pedido := range []int{-5, 0, 1, 7} {
		got, err := generateRandomPassword(pedido)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(got), 8, "pedido %d debe elevarse al minimo de 8", pedido)
	}
}

func TestGenerateRandomPassword_RespetaLongitudPedida(t *testing.T) {
	for _, pedido := range []int{8, 12, 32, 64} {
		got, err := generateRandomPassword(pedido)
		require.NoError(t, err)
		assert.Equal(t, pedido, len(got))
	}
}

func TestGenerateRandomPassword_SiempreIncluyeLasCuatroClases(t *testing.T) {
	const specials = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	for i := 0; i < 200; i++ {
		got, err := generateRandomPassword(12)
		require.NoError(t, err)

		var tieneMinuscula, tieneMayuscula, tieneDigito, tieneEspecial bool
		for _, r := range got {
			switch {
			case unicode.IsLower(r):
				tieneMinuscula = true
			case unicode.IsUpper(r):
				tieneMayuscula = true
			case unicode.IsDigit(r):
				tieneDigito = true
			case strings.ContainsRune(specials, r):
				tieneEspecial = true
			}
		}

		require.True(t, tieneMinuscula, "iteracion %d sin minuscula: %q", i, got)
		require.True(t, tieneMayuscula, "iteracion %d sin mayuscula: %q", i, got)
		require.True(t, tieneDigito, "iteracion %d sin digito: %q", i, got)
		require.True(t, tieneEspecial, "iteracion %d sin especial: %q", i, got)
	}
}

func TestGenerateRandomPassword_NoRepiteValores(t *testing.T) {
	vistas := make(map[string]bool, 300)
	for i := 0; i < 300; i++ {
		got, err := generateRandomPassword(12)
		require.NoError(t, err)
		require.False(t, vistas[got], "contrasena repetida en la iteracion %d: %q", i, got)
		vistas[got] = true
	}
}

func TestGenerateRandomPassword_NoDejaLasClasesEnPosicionFija(t *testing.T) {
	const specials = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	posicionesConEspecial := make(map[int]bool)
	for i := 0; i < 300; i++ {
		got, err := generateRandomPassword(12)
		require.NoError(t, err)
		for pos, r := range got {
			if strings.ContainsRune(specials, r) {
				posicionesConEspecial[pos] = true
			}
		}
	}
	assert.Greater(t, len(posicionesConEspecial), 1,
		"si el caracter especial cayera siempre en la misma posicion la mezcla no estaria funcionando")
}
