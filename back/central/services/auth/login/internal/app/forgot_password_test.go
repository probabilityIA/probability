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

const respuestaGenerica = "Si el correo esta registrado, recibiras un enlace para restablecer tu contrasena."

func buildForgotUseCase(
	repo *mocks.AuthRepositoryMock,
	email *mocks.EmailSenderMock,
	otp *mocks.OTPEventPublisherMock,
	valores map[string]string,
) *AuthUseCase {
	if email == nil {
		email = &mocks.EmailSenderMock{}
	}
	if otp == nil {
		otp = &mocks.OTPEventPublisherMock{}
	}
	if valores == nil {
		valores = map[string]string{}
	}
	return &AuthUseCase{
		repository:   repo,
		emailSender:  email,
		otpPublisher: otp,
		log:          mocks.NewSilentLogger(),
		env:          &configStub{valores: valores},
	}
}

func TestFrontendBaseURL_PorDefectoEsLocalhost(t *testing.T) {
	uc := buildForgotUseCase(&mocks.AuthRepositoryMock{}, nil, nil, nil)

	assert.Equal(t, "http://localhost:3000", uc.frontendBaseURL())
}

func TestFrontendBaseURL_RecortaLaBarraFinal(t *testing.T) {
	uc := buildForgotUseCase(&mocks.AuthRepositoryMock{}, nil, nil,
		map[string]string{"FRONTEND_BASE_URL": "https://app.test.com/"})

	assert.Equal(t, "https://app.test.com", uc.frontendBaseURL())
}

func TestForgotPassword_CorreoVacio_RespuestaGenerica(t *testing.T) {
	consultado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			consultado = true
			return nil, nil
		},
	}
	uc := buildForgotUseCase(repo, nil, nil, nil)

	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{Email: "  "})

	require.NoError(t, err)
	assert.True(t, got.Success)
	assert.Equal(t, respuestaGenerica, got.Message)
	assert.False(t, consultado)
}

func TestForgotPassword_UsuarioInexistente_NoRevelaNada(t *testing.T) {
	enviado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return nil, nil
		},
	}
	email := &mocks.EmailSenderMock{
		SendHTMLFn: func(ctx context.Context, to, subject, html string) error {
			enviado = true
			return nil
		},
	}
	uc := buildForgotUseCase(repo, email, nil, nil)

	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email: "noexiste@test.com",
	})

	require.NoError(t, err)
	assert.True(t, got.Success)
	assert.Equal(t, respuestaGenerica, got.Message)
	assert.False(t, enviado, "no se envia correo a una cuenta que no existe")
}

func TestForgotPassword_UsuarioInactivo_NoEnviaNada(t *testing.T) {
	enviado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(1, "x")
			u.IsActive = false
			return u, nil
		},
	}
	email := &mocks.EmailSenderMock{
		SendHTMLFn: func(ctx context.Context, to, subject, html string) error {
			enviado = true
			return nil
		},
	}
	uc := buildForgotUseCase(repo, email, nil, nil)

	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email: "ana@test.com",
	})

	require.NoError(t, err)
	assert.Equal(t, respuestaGenerica, got.Message)
	assert.False(t, enviado)
}

func TestForgotPassword_FallaConsultarUsuario_RespuestaGenerica(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return nil, errors.New("db down")
		},
	}
	uc := buildForgotUseCase(repo, nil, nil, nil)

	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email: "ana@test.com",
	})

	require.NoError(t, err, "un fallo de base de datos se ve igual que cualquier otro caso")
	assert.Equal(t, respuestaGenerica, got.Message)
}

func TestForgotPassword_NormalizaElCorreo(t *testing.T) {
	var consultado string
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			consultado = email
			return nil, nil
		},
	}
	uc := buildForgotUseCase(repo, nil, nil, nil)

	_, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email: "  ANA@Test.COM  ",
	})

	require.NoError(t, err)
	assert.Equal(t, "ana@test.com", consultado)
}

func TestForgotPassword_InvalidaTokensPrevios(t *testing.T) {
	var invalidadoPara uint
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(42, "x"), nil
		},
		InvalidateUserPasswordResetTokensFn: func(ctx context.Context, userID uint) error {
			invalidadoPara = userID
			return nil
		},
	}
	uc := buildForgotUseCase(repo, nil, nil, nil)

	_, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email: "ana@test.com",
	})

	require.NoError(t, err)
	assert.Equal(t, uint(42), invalidadoPara,
		"pedir un nuevo enlace invalida los anteriores")
}

func TestForgotPassword_PorCorreo_GuardaHashYEnviaEnlace(t *testing.T) {
	var guardadoHash, guardadoCanal string
	var expiracion time.Time
	var htmlEnviado, destinatario string

	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(42, "x")
			u.Name = "Ana"
			return u, nil
		},
		CreatePasswordResetTokenFn: func(ctx context.Context, userID uint, tokenHash, channel string, expiresAt time.Time) error {
			guardadoHash, guardadoCanal, expiracion = tokenHash, channel, expiresAt
			return nil
		},
	}
	email := &mocks.EmailSenderMock{
		SendHTMLFn: func(ctx context.Context, to, subject, html string) error {
			destinatario, htmlEnviado = to, html
			return nil
		},
	}
	uc := buildForgotUseCase(repo, email, nil, map[string]string{"FRONTEND_BASE_URL": "https://app.test.com"})

	antes := time.Now()
	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email: "ana@test.com",
	})

	require.NoError(t, err)
	assert.Equal(t, respuestaGenerica, got.Message)
	assert.Equal(t, "email", guardadoCanal)
	assert.Len(t, guardadoHash, 64)
	assert.WithinDuration(t, antes.Add(passwordResetTokenTTL), expiracion, 5*time.Second)
	assert.Equal(t, "ana@test.com", destinatario)
	assert.Contains(t, htmlEnviado, "https://app.test.com/reset-password?token=")
	assert.Contains(t, htmlEnviado, "Hola Ana")
}

func TestForgotPassword_PorCorreo_ElEnlaceLlevaElTokenNoElHash(t *testing.T) {
	var guardadoHash, htmlEnviado string
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(42, "x"), nil
		},
		CreatePasswordResetTokenFn: func(ctx context.Context, userID uint, tokenHash, channel string, expiresAt time.Time) error {
			guardadoHash = tokenHash
			return nil
		},
	}
	email := &mocks.EmailSenderMock{
		SendHTMLFn: func(ctx context.Context, to, subject, html string) error {
			htmlEnviado = html
			return nil
		},
	}
	uc := buildForgotUseCase(repo, email, nil, nil)

	_, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{Email: "ana@test.com"})

	require.NoError(t, err)
	assert.NotContains(t, htmlEnviado, guardadoHash,
		"el correo lleva el token en claro, la base de datos guarda solo su hash")

	inicio := strings.Index(htmlEnviado, "reset-password?token=") + len("reset-password?token=")
	fin := strings.Index(htmlEnviado[inicio:], `"`) + inicio
	tokenEnviado := htmlEnviado[inicio:fin]
	assert.Equal(t, guardadoHash, hashResetToken(tokenEnviado))
}

func TestForgotPassword_FallaGuardarToken_NoEnviaCorreo(t *testing.T) {
	enviado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(42, "x"), nil
		},
		CreatePasswordResetTokenFn: func(ctx context.Context, userID uint, tokenHash, channel string, expiresAt time.Time) error {
			return errors.New("db down")
		},
	}
	email := &mocks.EmailSenderMock{
		SendHTMLFn: func(ctx context.Context, to, subject, html string) error {
			enviado = true
			return nil
		},
	}
	uc := buildForgotUseCase(repo, email, nil, nil)

	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{Email: "ana@test.com"})

	require.NoError(t, err)
	assert.Equal(t, respuestaGenerica, got.Message)
	assert.False(t, enviado, "sin token guardado no se envia un enlace que no funcionaria")
}

func TestForgotPassword_FallaEnviarCorreo_RespuestaGenerica(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			return usuarioActivo(42, "x"), nil
		},
	}
	email := &mocks.EmailSenderMock{
		SendHTMLFn: func(ctx context.Context, to, subject, html string) error {
			return errors.New("smtp caido")
		},
	}
	uc := buildForgotUseCase(repo, email, nil, nil)

	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{Email: "ana@test.com"})

	require.NoError(t, err)
	assert.True(t, got.Success)
	assert.Equal(t, respuestaGenerica, got.Message)
}

func TestForgotPassword_PorWhatsApp_PublicaElCodigo(t *testing.T) {
	var guardadoHash, guardadoCanal string
	var expiracion time.Time
	var evento domain.PasswordResetOTPEvent

	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(42, "x")
			u.Name = "Ana"
			u.Phone = "3001234567"
			return u, nil
		},
		CreatePasswordResetTokenFn: func(ctx context.Context, userID uint, tokenHash, channel string, expiresAt time.Time) error {
			guardadoHash, guardadoCanal, expiracion = tokenHash, channel, expiresAt
			return nil
		},
	}
	otp := &mocks.OTPEventPublisherMock{
		PublishPasswordResetOTPFn: func(ctx context.Context, e domain.PasswordResetOTPEvent) error {
			evento = e
			return nil
		},
	}
	uc := buildForgotUseCase(repo, nil, otp, nil)

	antes := time.Now()
	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email:   "ana@test.com",
		Channel: "whatsapp",
	})

	require.NoError(t, err)
	assert.Equal(t, respuestaGenerica, got.Message)
	assert.Equal(t, "whatsapp", guardadoCanal)
	assert.WithinDuration(t, antes.Add(otpTTL), expiracion, 5*time.Second)
	assert.Equal(t, "3001234567", evento.Phone)
	assert.Len(t, evento.Code, 6)
	assert.Equal(t, "Ana", evento.UserName)
	assert.Equal(t, 10, evento.ExpiresMinutes)
	assert.Equal(t, hashOTPCode(42, evento.Code), guardadoHash,
		"se guarda el hash del codigo ligado al usuario, no el codigo en claro")
}

func TestForgotPassword_PorWhatsApp_CanalEsInsensibleAMayusculas(t *testing.T) {
	publicado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(42, "x")
			u.Phone = "3001234567"
			return u, nil
		},
	}
	otp := &mocks.OTPEventPublisherMock{
		PublishPasswordResetOTPFn: func(ctx context.Context, e domain.PasswordResetOTPEvent) error {
			publicado = true
			return nil
		},
	}
	uc := buildForgotUseCase(repo, nil, otp, nil)

	_, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email:   "ana@test.com",
		Channel: "  WhatsApp  ",
	})

	require.NoError(t, err)
	assert.True(t, publicado)
}

func TestForgotPassword_CanalVacio_UsaCorreo(t *testing.T) {
	correoEnviado := false
	otpPublicado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(42, "x")
			u.Phone = "3001234567"
			return u, nil
		},
	}
	email := &mocks.EmailSenderMock{
		SendHTMLFn: func(ctx context.Context, to, subject, html string) error {
			correoEnviado = true
			return nil
		},
	}
	otp := &mocks.OTPEventPublisherMock{
		PublishPasswordResetOTPFn: func(ctx context.Context, e domain.PasswordResetOTPEvent) error {
			otpPublicado = true
			return nil
		},
	}
	uc := buildForgotUseCase(repo, email, otp, nil)

	_, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{Email: "ana@test.com"})

	require.NoError(t, err)
	assert.True(t, correoEnviado)
	assert.False(t, otpPublicado)
}

func TestForgotPassword_WhatsAppSinTelefono_NoPublica(t *testing.T) {
	publicado := false
	tokenGuardado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(42, "x")
			u.Phone = "   "
			return u, nil
		},
		CreatePasswordResetTokenFn: func(ctx context.Context, userID uint, tokenHash, channel string, expiresAt time.Time) error {
			tokenGuardado = true
			return nil
		},
	}
	otp := &mocks.OTPEventPublisherMock{
		PublishPasswordResetOTPFn: func(ctx context.Context, e domain.PasswordResetOTPEvent) error {
			publicado = true
			return nil
		},
	}
	uc := buildForgotUseCase(repo, nil, otp, nil)

	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email:   "ana@test.com",
		Channel: "whatsapp",
	})

	require.NoError(t, err)
	assert.Equal(t, respuestaGenerica, got.Message)
	assert.False(t, publicado)
	assert.False(t, tokenGuardado, "sin telefono no se genera codigo alguno")
}

func TestForgotPassword_WhatsAppFallaGuardar_NoPublica(t *testing.T) {
	publicado := false
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(42, "x")
			u.Phone = "3001234567"
			return u, nil
		},
		CreatePasswordResetTokenFn: func(ctx context.Context, userID uint, tokenHash, channel string, expiresAt time.Time) error {
			return errors.New("db down")
		},
	}
	otp := &mocks.OTPEventPublisherMock{
		PublishPasswordResetOTPFn: func(ctx context.Context, e domain.PasswordResetOTPEvent) error {
			publicado = true
			return nil
		},
	}
	uc := buildForgotUseCase(repo, nil, otp, nil)

	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email:   "ana@test.com",
		Channel: "whatsapp",
	})

	require.NoError(t, err)
	assert.Equal(t, respuestaGenerica, got.Message)
	assert.False(t, publicado, "sin codigo guardado no se envia un OTP que nunca validaria")
}

func TestForgotPassword_WhatsAppFallaPublicar_RespuestaGenerica(t *testing.T) {
	repo := &mocks.AuthRepositoryMock{
		GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
			u := usuarioActivo(42, "x")
			u.Phone = "3001234567"
			return u, nil
		},
	}
	otp := &mocks.OTPEventPublisherMock{
		PublishPasswordResetOTPFn: func(ctx context.Context, e domain.PasswordResetOTPEvent) error {
			return errors.New("rabbitmq caido")
		},
	}
	uc := buildForgotUseCase(repo, nil, otp, nil)

	got, err := uc.ForgotPassword(context.Background(), domain.ForgotPasswordRequest{
		Email:   "ana@test.com",
		Channel: "whatsapp",
	})

	require.NoError(t, err)
	assert.True(t, got.Success)
	assert.Equal(t, respuestaGenerica, got.Message)
}

func TestForgotPassword_TodasLasRutasDevuelvenElMismoMensaje(t *testing.T) {
	casos := []struct {
		nombre string
		repo   *mocks.AuthRepositoryMock
		req    domain.ForgotPasswordRequest
	}{
		{
			nombre: "correo vacio",
			repo:   &mocks.AuthRepositoryMock{},
			req:    domain.ForgotPasswordRequest{Email: ""},
		},
		{
			nombre: "usuario inexistente",
			repo: &mocks.AuthRepositoryMock{
				GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
					return nil, nil
				},
			},
			req: domain.ForgotPasswordRequest{Email: "x@test.com"},
		},
		{
			nombre: "usuario existente",
			repo: &mocks.AuthRepositoryMock{
				GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
					return usuarioActivo(1, "x"), nil
				},
			},
			req: domain.ForgotPasswordRequest{Email: "ana@test.com"},
		},
		{
			nombre: "fallo de base de datos",
			repo: &mocks.AuthRepositoryMock{
				GetUserByEmailFn: func(ctx context.Context, email string) (*domain.UserAuthInfo, error) {
					return nil, errors.New("db down")
				},
			},
			req: domain.ForgotPasswordRequest{Email: "ana@test.com"},
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			uc := buildForgotUseCase(tc.repo, nil, nil, nil)

			got, err := uc.ForgotPassword(context.Background(), tc.req)

			require.NoError(t, err)
			assert.True(t, got.Success)
			assert.Equal(t, respuestaGenerica, got.Message,
				"ninguna ruta permite distinguir si el correo esta registrado")
		})
	}
}
