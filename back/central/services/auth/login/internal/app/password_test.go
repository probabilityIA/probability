package app

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const (
	minusculas = "abcdefghijklmnopqrstuvwxyz"
	mayusculas = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitos    = "0123456789"
	especiales = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

func buildUseCase(repo *mocks.AuthRepositoryMock) *AuthUseCase {
	return &AuthUseCase{
		repository: repo,
		log:        mocks.NewSilentLogger(),
	}
}

func usuarioActivo(id uint, passwordPlano string) *domain.UserAuthInfo {
	hash, _ := bcrypt.GenerateFromPassword([]byte(passwordPlano), bcrypt.MinCost)
	return &domain.UserAuthInfo{
		ID:       id,
		Email:    "ana@test.com",
		Password: string(hash),
		IsActive: true,
	}
}

func contieneAlgunoDe(s, conjunto string) bool {
	return strings.ContainsAny(s, conjunto)
}

func TestGenerateRandomPassword_LongitudMinimaEsOcho(t *testing.T) {
	for _, n := range []int{0, 1, 5, 7, -3} {
		got, err := generateRandomPassword(n)

		require.NoError(t, err)
		assert.Len(t, got, 8, "una longitud menor a 8 se eleva al minimo recomendado")
	}
}

func TestGenerateRandomPassword_RespetaLaLongitudPedida(t *testing.T) {
	for _, n := range []int{8, 12, 16, 32, 64} {
		got, err := generateRandomPassword(n)

		require.NoError(t, err)
		assert.Len(t, got, n)
	}
}

func TestGenerateRandomPassword_IncluyeLosCuatroTiposDeCaracter(t *testing.T) {
	for i := 0; i < 50; i++ {
		got, err := generateRandomPassword(12)
		require.NoError(t, err)

		assert.True(t, contieneAlgunoDe(got, minusculas), "falta minuscula en %q", got)
		assert.True(t, contieneAlgunoDe(got, mayusculas), "falta mayuscula en %q", got)
		assert.True(t, contieneAlgunoDe(got, digitos), "falta digito en %q", got)
		assert.True(t, contieneAlgunoDe(got, especiales), "falta caracter especial en %q", got)
	}
}

func TestGenerateRandomPassword_LongitudMinimaTambienCumpleLosCuatroTipos(t *testing.T) {
	for i := 0; i < 30; i++ {
		got, err := generateRandomPassword(8)
		require.NoError(t, err)

		assert.True(t, contieneAlgunoDe(got, minusculas))
		assert.True(t, contieneAlgunoDe(got, mayusculas))
		assert.True(t, contieneAlgunoDe(got, digitos))
		assert.True(t, contieneAlgunoDe(got, especiales))
	}
}

func TestGenerateRandomPassword_NoSeRepite(t *testing.T) {
	vistas := make(map[string]bool)

	for i := 0; i < 100; i++ {
		got, err := generateRandomPassword(12)
		require.NoError(t, err)
		assert.False(t, vistas[got], "se genero dos veces la misma contrasena: %q", got)
		vistas[got] = true
	}
}

func TestGenerateRandomPassword_NoDejaLosTiposEnPosicionFija(t *testing.T) {
	primerCaracterSiempreMinuscula := true

	for i := 0; i < 60; i++ {
		got, err := generateRandomPassword(12)
		require.NoError(t, err)

		if !strings.ContainsAny(string(got[0]), minusculas) {
			primerCaracterSiempreMinuscula = false
			break
		}
	}

	assert.False(t, primerCaracterSiempreMinuscula,
		"la mezcla debe evitar que el primer caracter sea siempre del primer conjunto")
}

func TestGeneratePassword_UsuarioInexistente_RetornaError(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return nil, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.GeneratePassword(context.Background(), domain.GeneratePasswordRequest{UserID: 1})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "usuario no encontrado")
	assert.Nil(t, got)
}

func TestGeneratePassword_UsuarioInactivo_Rechaza(t *testing.T) {
	cambiada := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(userID, "actual")
			u.IsActive = false
			return u, nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			cambiada = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.GeneratePassword(context.Background(), domain.GeneratePasswordRequest{UserID: 1})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "usuario inactivo")
	assert.Nil(t, got)
	assert.False(t, cambiada, "no se genera contrasena para un usuario desactivado")
}

func TestGeneratePassword_FallaConsultarUsuario_NoFiltraElDetalle(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return nil, errors.New("connection refused to postgres at 10.0.0.5")
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.GeneratePassword(context.Background(), domain.GeneratePasswordRequest{UserID: 1})

	require.Error(t, err)
	assert.Equal(t, "error interno del servidor", err.Error(),
		"el detalle de infraestructura no se expone al cliente")
	assert.NotContains(t, err.Error(), "postgres")
	assert.Nil(t, got)
}

func TestGeneratePassword_Exitoso_DevuelveLaContrasenaGenerada(t *testing.T) {
	var guardada string
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "actual"), nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			guardada = newPassword
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.GeneratePassword(context.Background(), domain.GeneratePasswordRequest{UserID: 1})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.Success)
	assert.Equal(t, "ana@test.com", got.Email)
	assert.Len(t, got.Password, 12)
	assert.Equal(t, guardada, got.Password,
		"la contrasena devuelta es la misma que se guardo")
}

func TestGeneratePassword_FallaGuardar_RetornaError(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "actual"), nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			return errors.New("db down")
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.GeneratePassword(context.Background(), domain.GeneratePasswordRequest{UserID: 1})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error al actualizar")
	assert.Nil(t, got)
}

func TestChangePassword_UsuarioInexistente_RetornaError(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return nil, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ChangePassword(context.Background(), domain.ChangePasswordRequest{UserID: 1})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "usuario no encontrado")
	assert.Nil(t, got)
}

func TestChangePassword_UsuarioInactivo_Rechaza(t *testing.T) {
	cambiada := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(userID, "actual")
			u.IsActive = false
			return u, nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			cambiada = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ChangePassword(context.Background(), domain.ChangePasswordRequest{
		UserID:          1,
		CurrentPassword: "actual",
		NewPassword:     "nueva",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "usuario inactivo")
	assert.Nil(t, got)
	assert.False(t, cambiada)
}

func TestChangePassword_ContrasenaActualIncorrecta_Rechaza(t *testing.T) {
	cambiada := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "la-correcta"), nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			cambiada = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ChangePassword(context.Background(), domain.ChangePasswordRequest{
		UserID:          1,
		CurrentPassword: "la-equivocada",
		NewPassword:     "nueva",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "actual incorrecta")
	assert.Nil(t, got)
	assert.False(t, cambiada, "sin la contrasena actual no se cambia nada")
}

func TestChangePassword_NuevaIgualALaActual_Rechaza(t *testing.T) {
	cambiada := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "la-misma"), nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			cambiada = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ChangePassword(context.Background(), domain.ChangePasswordRequest{
		UserID:          1,
		CurrentPassword: "la-misma",
		NewPassword:     "la-misma",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "diferente a la actual")
	assert.Nil(t, got)
	assert.False(t, cambiada)
}

func TestChangePassword_Exitoso(t *testing.T) {
	var guardada string
	loginActualizado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "actual"), nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			guardada = newPassword
			return nil
		},
		UpdateLastLoginFn: func(ctx context.Context, userID uint) error {
			loginActualizado = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ChangePassword(context.Background(), domain.ChangePasswordRequest{
		UserID:          1,
		CurrentPassword: "actual",
		NewPassword:     "nueva-segura",
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.Success)
	assert.Equal(t, "nueva-segura", guardada)
	assert.True(t, loginActualizado)
}

func TestChangePassword_FallaActualizarUltimoLogin_NoAbortaElCambio(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "actual"), nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			return nil
		},
		UpdateLastLoginFn: func(ctx context.Context, userID uint) error {
			return errors.New("db down")
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ChangePassword(context.Background(), domain.ChangePasswordRequest{
		UserID:          1,
		CurrentPassword: "actual",
		NewPassword:     "nueva",
	})

	require.NoError(t, err, "el cambio de contrasena ya fue exitoso, el ultimo login es secundario")
	assert.True(t, got.Success)
}

func TestChangePassword_FallaGuardar_RetornaError(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return usuarioActivo(userID, "actual"), nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			return errors.New("db down")
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ChangePassword(context.Background(), domain.ChangePasswordRequest{
		UserID:          1,
		CurrentPassword: "actual",
		NewPassword:     "nueva",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error al cambiar")
	assert.Nil(t, got)
}

func TestChangePassword_FallaConsultarUsuario_NoFiltraElDetalle(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByIDFn: func(ctx context.Context, userID uint) (*domain.UserAuthInfo, error) {
			return nil, errors.New("connection refused to postgres at 10.0.0.5")
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ChangePassword(context.Background(), domain.ChangePasswordRequest{UserID: 1})

	require.Error(t, err)
	assert.Equal(t, "error interno del servidor", err.Error())
	assert.Nil(t, got)
}
