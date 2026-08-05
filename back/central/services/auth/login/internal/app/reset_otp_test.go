package app

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tokenValido(userID uint, hash string) *domain.PasswordResetTokenInfo {
	return &domain.PasswordResetTokenInfo{
		ID:        1,
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
}

func TestHashResetToken_EsDeterministaYSHA256(t *testing.T) {
	a := hashResetToken("token-abc")
	b := hashResetToken("token-abc")
	c := hashResetToken("token-abd")

	assert.Equal(t, a, b)
	assert.NotEqual(t, a, c)
	assert.Len(t, a, 64, "sha256 en hexadecimal son 64 caracteres")
	assert.NotContains(t, a, "token-abc", "el hash no contiene el token en claro")
}

func TestGenerateResetToken_TokenYHashCorresponden(t *testing.T) {
	raw, hash, err := generateResetToken()

	require.NoError(t, err)
	assert.Len(t, raw, 64, "32 bytes en hexadecimal")
	assert.Len(t, hash, 64)
	assert.Equal(t, hash, hashResetToken(raw),
		"el hash guardado es el sha256 del token entregado al usuario")
	assert.NotEqual(t, raw, hash, "nunca se guarda el token en claro")
}

func TestGenerateResetToken_NoSeRepite(t *testing.T) {
	vistos := make(map[string]bool)

	for i := 0; i < 100; i++ {
		raw, _, err := generateResetToken()
		require.NoError(t, err)
		assert.False(t, vistos[raw], "token de reseteo repetido")
		vistos[raw] = true
	}
}

func TestGenerateOTPCode_SeisDigitos(t *testing.T) {
	for i := 0; i < 100; i++ {
		got, err := generateOTPCode()

		require.NoError(t, err)
		assert.Len(t, got, 6, "el codigo OTP siempre tiene seis digitos")
		assert.Regexp(t, `^\d{6}$`, got, "solo digitos, con ceros a la izquierda si hace falta")
	}
}

func TestGenerateOTPCode_Varia(t *testing.T) {
	vistos := make(map[string]bool)

	for i := 0; i < 60; i++ {
		got, err := generateOTPCode()
		require.NoError(t, err)
		vistos[got] = true
	}

	assert.Greater(t, len(vistos), 30, "los codigos generados deben variar")
}

func TestHashOTPCode_IncluyeElUsuario(t *testing.T) {
	mismoCodigoUsuarioA := hashOTPCode(1, "123456")
	mismoCodigoUsuarioB := hashOTPCode(2, "123456")

	assert.NotEqual(t, mismoCodigoUsuarioA, mismoCodigoUsuarioB,
		"el mismo codigo para otro usuario produce otro hash, no es reutilizable entre cuentas")
	assert.Equal(t, mismoCodigoUsuarioA, hashOTPCode(1, "123456"))
	assert.Len(t, mismoCodigoUsuarioA, 64)
}

func TestMaskPhone(t *testing.T) {
	cases := map[string]string{
		"3001234567": "***567",
		"123":        "***",
		"12":         "***",
		"":           "***",
		"  3009999":  "***999",
	}

	for entrada, esperado := range cases {
		t.Run(entrada, func(t *testing.T) {
			assert.Equal(t, esperado, maskPhone(entrada))
		})
	}
}

func TestMaskPhone_NoRevelaElNumeroCompleto(t *testing.T) {
	got := maskPhone("3001234567")

	assert.NotContains(t, got, "3001234")
	assert.True(t, strings.HasPrefix(got, "***"))
}

func TestBuildResetPasswordEmail_IncluyeElEnlace(t *testing.T) {
	got := buildResetPasswordEmail("Ana", "https://app/reset?token=abc")

	assert.Contains(t, got, "https://app/reset?token=abc")
	assert.Contains(t, got, "Hola Ana")
}

func TestBuildResetPasswordEmail_SinNombre_UsaSaludoGenerico(t *testing.T) {
	got := buildResetPasswordEmail("   ", "https://app/reset")

	assert.Contains(t, got, "Hola,")
	assert.NotContains(t, got, "Hola ,")
}

func TestResetPassword_TokenVacio_Rechaza(t *testing.T) {
	consultado := false
	repo := &mocks.AuthRepositoryMock{
		GetValidPasswordResetTokenFn: func(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
			consultado = true
			return nil, nil
		},
	}
	uc := buildUseCase(repo)

	for _, token := range []string{"", "   "} {
		got, err := uc.ResetPassword(context.Background(), domain.ResetPasswordRequest{
			Token:       token,
			NewPassword: "nueva-segura",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token invalido")
		assert.Nil(t, got)
	}
	assert.False(t, consultado)
}

func TestResetPassword_ContrasenaCorta_Rechaza(t *testing.T) {
	cambiada := false
	repo := &mocks.AuthRepositoryMock{
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			cambiada = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	for _, pass := range []string{"", "a", "12345"} {
		got, err := uc.ResetPassword(context.Background(), domain.ResetPasswordRequest{
			Token:       "token-valido",
			NewPassword: pass,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "al menos 6 caracteres")
		assert.Nil(t, got)
	}
	assert.False(t, cambiada)
}

func TestResetPassword_TokenInexistente_Rechaza(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetValidPasswordResetTokenFn: func(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
			return nil, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ResetPassword(context.Background(), domain.ResetPasswordRequest{
		Token:       "inventado",
		NewPassword: "nueva-segura",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "token invalido o expirado")
	assert.Nil(t, got)
}

func TestResetPassword_TokenYaUsado_Rechaza(t *testing.T) {
	usado := time.Now().Add(-time.Hour)
	cambiada := false
	repo := &mocks.AuthRepositoryMock{
		GetValidPasswordResetTokenFn: func(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
			info := tokenValido(1, tokenHash)
			info.UsedAt = &usado
			return info, nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			cambiada = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ResetPassword(context.Background(), domain.ResetPasswordRequest{
		Token:       "token-valido",
		NewPassword: "nueva-segura",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "token invalido o expirado")
	assert.Nil(t, got)
	assert.False(t, cambiada, "un token de reseteo es de un solo uso")
}

func TestResetPassword_TokenExpirado_Rechaza(t *testing.T) {
	cambiada := false
	repo := &mocks.AuthRepositoryMock{
		GetValidPasswordResetTokenFn: func(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
			info := tokenValido(1, tokenHash)
			info.ExpiresAt = time.Now().Add(-time.Minute)
			return info, nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			cambiada = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ResetPassword(context.Background(), domain.ResetPasswordRequest{
		Token:       "token-valido",
		NewPassword: "nueva-segura",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "expirado")
	assert.Nil(t, got)
	assert.False(t, cambiada)
}

func TestResetPassword_BuscaPorHashNoPorTokenEnClaro(t *testing.T) {
	var hashConsultado string
	repo := &mocks.AuthRepositoryMock{
		GetValidPasswordResetTokenFn: func(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
			hashConsultado = tokenHash
			return tokenValido(1, tokenHash), nil
		},
	}
	uc := buildUseCase(repo)

	_, err := uc.ResetPassword(context.Background(), domain.ResetPasswordRequest{
		Token:       "mi-token-secreto",
		NewPassword: "nueva-segura",
	})

	require.NoError(t, err)
	assert.NotEqual(t, "mi-token-secreto", hashConsultado)
	assert.Equal(t, hashResetToken("mi-token-secreto"), hashConsultado)
}

func TestResetPassword_Exitoso_MarcaElTokenComoUsado(t *testing.T) {
	var marcado uint
	var cambiadaPara uint
	repo := &mocks.AuthRepositoryMock{
		GetValidPasswordResetTokenFn: func(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
			return tokenValido(42, tokenHash), nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			cambiadaPara = userID
			return nil
		},
		MarkPasswordResetTokenUsedFn: func(ctx context.Context, tokenID uint) error {
			marcado = tokenID
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ResetPassword(context.Background(), domain.ResetPasswordRequest{
		Token:       "token-valido",
		NewPassword: "nueva-segura",
	})

	require.NoError(t, err)
	assert.True(t, got.Success)
	assert.Equal(t, uint(42), cambiadaPara)
	assert.Equal(t, uint(1), marcado)
}

func TestResetPassword_FallaMarcarUsado_NoAbortaElReseteo(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetValidPasswordResetTokenFn: func(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
			return tokenValido(1, tokenHash), nil
		},
		MarkPasswordResetTokenUsedFn: func(ctx context.Context, tokenID uint) error {
			return errors.New("db down")
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ResetPassword(context.Background(), domain.ResetPasswordRequest{
		Token:       "token-valido",
		NewPassword: "nueva-segura",
	})

	require.NoError(t, err)
	assert.True(t, got.Success)
}

func TestResetPassword_FallaConsultarToken_NoFiltraElDetalle(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetValidPasswordResetTokenFn: func(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
			return nil, errors.New("connection refused to postgres at 10.0.0.5")
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ResetPassword(context.Background(), domain.ResetPasswordRequest{
		Token:       "token-valido",
		NewPassword: "nueva-segura",
	})

	require.Error(t, err)
	assert.Equal(t, "error interno del servidor", err.Error())
	assert.Nil(t, got)
}

func TestResetPassword_FallaCambiarContrasena_RetornaError(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetValidPasswordResetTokenFn: func(ctx context.Context, tokenHash string) (*domain.PasswordResetTokenInfo, error) {
			return tokenValido(1, tokenHash), nil
		},
		ChangePasswordFn: func(ctx context.Context, userID uint, newPassword string) error {
			return errors.New("db down")
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.ResetPassword(context.Background(), domain.ResetPasswordRequest{
		Token:       "token-valido",
		NewPassword: "nueva-segura",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error al actualizar")
	assert.Nil(t, got)
}

func TestVerifyOTP_DatosVacios_RespuestaGenerica(t *testing.T) {
	consultado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			consultado = true
			return nil, nil
		},
	}
	uc := buildUseCase(repo)

	casos := []domain.VerifyOTPRequest{
		{Email: "", Code: "123456"},
		{Email: "ana@test.com", Code: ""},
		{Email: "  ", Code: "  "},
	}

	for _, req := range casos {
		got, err := uc.VerifyOTP(context.Background(), req)

		require.NoError(t, err)
		assert.False(t, got.Success)
		assert.Equal(t, "Codigo invalido o expirado", got.Message)
	}
	assert.False(t, consultado)
}

func TestVerifyOTP_UsuarioInexistente_NoRevelaSiExiste(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return nil, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.VerifyOTP(context.Background(), domain.VerifyOTPRequest{
		Email: "noexiste@test.com",
		Code:  "123456",
	})

	require.NoError(t, err, "no se devuelve error para no revelar si la cuenta existe")
	assert.False(t, got.Success)
	assert.Equal(t, "Codigo invalido o expirado", got.Message)
}

func TestVerifyOTP_UsuarioInactivo_RespuestaGenerica(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(1, "x")
			u.IsActive = false
			return u, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.VerifyOTP(context.Background(), domain.VerifyOTPRequest{
		Email: "ana@test.com",
		Code:  "123456",
	})

	require.NoError(t, err)
	assert.False(t, got.Success)
}

func TestVerifyOTP_NormalizaElCorreo(t *testing.T) {
	var consultado string
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			consultado = email
			return nil, nil
		},
	}
	uc := buildUseCase(repo)

	_, err := uc.VerifyOTP(context.Background(), domain.VerifyOTPRequest{
		Email: "  ANA@Test.COM  ",
		Code:  "123456",
	})

	require.NoError(t, err)
	assert.Equal(t, "ana@test.com", consultado)
}

func TestVerifyOTP_SinTokenActivo_RespuestaGenerica(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(1, "x"), nil
		},
		GetActiveOTPTokenFn: func(ctx context.Context, userID uint) (*domain.PasswordResetTokenInfo, error) {
			return nil, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.VerifyOTP(context.Background(), domain.VerifyOTPRequest{
		Email: "ana@test.com",
		Code:  "123456",
	})

	require.NoError(t, err)
	assert.False(t, got.Success)
}

func TestVerifyOTP_TokenExpirado_RespuestaGenerica(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(1, "x"), nil
		},
		GetActiveOTPTokenFn: func(ctx context.Context, userID uint) (*domain.PasswordResetTokenInfo, error) {
			info := tokenValido(userID, hashOTPCode(userID, "123456"))
			info.ExpiresAt = time.Now().Add(-time.Minute)
			return info, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.VerifyOTP(context.Background(), domain.VerifyOTPRequest{
		Email: "ana@test.com",
		Code:  "123456",
	})

	require.NoError(t, err)
	assert.False(t, got.Success)
}

func TestVerifyOTP_ExcesoDeIntentos_InvalidaElToken(t *testing.T) {
	invalidado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(1, "x"), nil
		},
		GetActiveOTPTokenFn: func(ctx context.Context, userID uint) (*domain.PasswordResetTokenInfo, error) {
			info := tokenValido(userID, hashOTPCode(userID, "123456"))
			info.Attempts = maxOTPAttempts
			return info, nil
		},
		MarkPasswordResetTokenUsedFn: func(ctx context.Context, tokenID uint) error {
			invalidado = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.VerifyOTP(context.Background(), domain.VerifyOTPRequest{
		Email: "ana@test.com",
		Code:  "123456",
	})

	require.NoError(t, err)
	assert.False(t, got.Success)
	assert.True(t, invalidado,
		"al agotar los intentos el codigo se invalida aunque el codigo enviado sea correcto")
}

func TestVerifyOTP_CodigoIncorrecto_IncrementaIntentos(t *testing.T) {
	var incrementado uint
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(1, "x"), nil
		},
		GetActiveOTPTokenFn: func(ctx context.Context, userID uint) (*domain.PasswordResetTokenInfo, error) {
			return tokenValido(userID, hashOTPCode(userID, "111111")), nil
		},
		IncrementPasswordResetTokenAttemptsFn: func(ctx context.Context, tokenID uint) error {
			incrementado = tokenID
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.VerifyOTP(context.Background(), domain.VerifyOTPRequest{
		Email: "ana@test.com",
		Code:  "999999",
	})

	require.NoError(t, err)
	assert.False(t, got.Success)
	assert.Equal(t, uint(1), incrementado)
}

func TestVerifyOTP_CodigoCorrecto_EntregaTokenDeReseteo(t *testing.T) {
	var guardadoHash string
	var guardadoCanal string
	marcadoUsado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(42, "x"), nil
		},
		GetActiveOTPTokenFn: func(ctx context.Context, userID uint) (*domain.PasswordResetTokenInfo, error) {
			return tokenValido(userID, hashOTPCode(userID, "123456")), nil
		},
		MarkPasswordResetTokenUsedFn: func(ctx context.Context, tokenID uint) error {
			marcadoUsado = true
			return nil
		},
		CreatePasswordResetTokenFn: func(ctx context.Context, userID uint, tokenHash, channel string, expiresAt time.Time) error {
			guardadoHash, guardadoCanal = tokenHash, channel
			return nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.VerifyOTP(context.Background(), domain.VerifyOTPRequest{
		Email: "ana@test.com",
		Code:  "123456",
	})

	require.NoError(t, err)
	assert.True(t, got.Success)
	assert.NotEmpty(t, got.Token)
	assert.True(t, marcadoUsado, "el OTP se consume al verificarse")
	assert.Equal(t, "email", guardadoCanal)
	assert.Equal(t, hashResetToken(got.Token), guardadoHash,
		"se guarda el hash del token, no el token entregado")
}

func TestVerifyOTP_FallaGuardarTokenNuevo_RetornaErrorInterno(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(1, "x"), nil
		},
		GetActiveOTPTokenFn: func(ctx context.Context, userID uint) (*domain.PasswordResetTokenInfo, error) {
			return tokenValido(userID, hashOTPCode(userID, "123456")), nil
		},
		CreatePasswordResetTokenFn: func(ctx context.Context, userID uint, tokenHash, channel string, expiresAt time.Time) error {
			return errors.New("db down")
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.VerifyOTP(context.Background(), domain.VerifyOTPRequest{
		Email: "ana@test.com",
		Code:  "123456",
	})

	require.NoError(t, err)
	assert.False(t, got.Success)
	assert.Empty(t, got.Token)
}

func TestVerifyOTP_FallaConsultarUsuario_RespuestaGenerica(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return nil, errors.New("db down")
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.VerifyOTP(context.Background(), domain.VerifyOTPRequest{
		Email: "ana@test.com",
		Code:  "123456",
	})

	require.NoError(t, err)
	assert.False(t, got.Success)
	assert.Equal(t, "Codigo invalido o expirado", got.Message,
		"un fallo de base de datos se ve igual que un codigo invalido")
}
