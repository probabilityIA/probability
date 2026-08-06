package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
)

const mensajeGenerico = "Si existe una cuenta pendiente de verificacion, te enviamos un nuevo enlace."

func pendiente(e *entorno, ajustar func(*domain.PendingDemoUser)) {
	e.repo.GetDemoUserByEmailFn = func(ctx context.Context, email string) (*domain.PendingDemoUser, error) {
		u := &domain.PendingDemoUser{
			UserID:       7,
			FullName:     "Sebastian",
			BusinessName: "Tienda Probability",
			IsActive:     false,
		}
		if ajustar != nil {
			ajustar(u)
		}
		return u, nil
	}
}

func TestDemoResend_ReenviaElCorreoConUnTokenNuevo(t *testing.T) {
	e := nuevoEntorno(t)
	pendiente(e, nil)

	resp, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{Email: "  DEMO@Tienda.com "})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	require.Len(t, e.repo.TokensInvalidados, 1)
	assert.EqualValues(t, 7, e.repo.TokensInvalidados[0],
		"los tokens viejos se invalidan primero para que solo el ultimo enlace enviado funcione")
	require.Len(t, e.repo.TokensCreados, 1)
	require.Len(t, e.email.Enviados, 1)
	assert.Equal(t, "demo@tienda.com", e.email.Enviados[0].To)
}

func TestDemoResend_ElTokenNuevoDeCorreoVuelveACaducarEn24Horas(t *testing.T) {
	e := nuevoEntorno(t)
	pendiente(e, nil)
	antes := time.Now()

	_, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{Email: "demo@tienda.com"})

	require.NoError(t, err)
	assert.WithinDuration(t, antes.Add(24*time.Hour), e.repo.TokensCreados[0].ExpiresAt, time.Minute)
}

func TestDemoResend_SiempreDevuelveElMismoMensajeSinRevelarSiLaCuentaExiste(t *testing.T) {
	casos := map[string]func(*entorno){
		"correo vacio": func(e *entorno) {},
		"el usuario no existe": func(e *entorno) {
			e.repo.GetDemoUserByEmailFn = func(ctx context.Context, email string) (*domain.PendingDemoUser, error) {
				return nil, nil
			}
		},
		"la cuenta ya esta activa": func(e *entorno) {
			pendiente(e, func(u *domain.PendingDemoUser) { u.IsActive = true })
		},
		"falla la consulta": func(e *entorno) {
			e.repo.GetDemoUserByEmailFn = func(ctx context.Context, email string) (*domain.PendingDemoUser, error) {
				return nil, errors.New("timeout")
			}
		},
	}

	for nombre, preparar := range casos {
		t.Run(nombre, func(t *testing.T) {
			e := nuevoEntorno(t)
			preparar(e)
			correo := "demo@tienda.com"
			if nombre == "correo vacio" {
				correo = "   "
			}

			resp, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{Email: correo})

			require.NoError(t, err)
			assert.True(t, resp.Success)
			assert.Equal(t, mensajeGenerico, resp.Message,
				"la respuesta no cambia nunca: el endpoint es publico y no debe servir para enumerar correos registrados")
			assert.Empty(t, e.email.Enviados)
			assert.Empty(t, e.repo.TokensCreados)
		})
	}
}

func TestDemoResend_BloqueaUnSegundoReenvioAntesDelMinutoDeEspera(t *testing.T) {
	hace30s := time.Now().Add(-30 * time.Second)
	e := nuevoEntorno(t)
	pendiente(e, func(u *domain.PendingDemoUser) { u.LastTokenCreatedAt = &hace30s })

	resp, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{Email: "demo@tienda.com"})

	require.NoError(t, err)
	assert.Equal(t, mensajeGenerico, resp.Message,
		"el cooldown es invisible para el cliente: responde igual que un reenvio exitoso")
	assert.Empty(t, e.email.Enviados)
	assert.Empty(t, e.repo.TokensCreados)
}

func TestDemoResend_PermiteElReenvioPasadoElMinuto(t *testing.T) {
	hace2min := time.Now().Add(-2 * time.Minute)
	e := nuevoEntorno(t)
	pendiente(e, func(u *domain.PendingDemoUser) { u.LastTokenCreatedAt = &hace2min })

	_, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{Email: "demo@tienda.com"})

	require.NoError(t, err)
	assert.Len(t, e.email.Enviados, 1)
}

func TestDemoResend_SinTokenPrevioNoAplicaCooldown(t *testing.T) {
	e := nuevoEntorno(t)
	pendiente(e, func(u *domain.PendingDemoUser) { u.LastTokenCreatedAt = nil })

	_, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{Email: "demo@tienda.com"})

	require.NoError(t, err)
	assert.Len(t, e.email.Enviados, 1)
}

func TestDemoResend_PorWhatsAppPublicaUnOTPNuevoYNoMandaCorreo(t *testing.T) {
	e := nuevoEntorno(t)
	pendiente(e, func(u *domain.PendingDemoUser) { u.Phone = "573001112233" })

	_, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{
		Email: "demo@tienda.com", Channel: "WhatsApp", Phone: "573001112233",
	})

	require.NoError(t, err)
	require.Len(t, e.otp.Publicados, 1)
	assert.Regexp(t, `^\d{6}$`, e.otp.Publicados[0].Code)
	assert.Equal(t, "Sebastian", e.otp.Publicados[0].UserName)
	assert.Equal(t, 15, e.otp.Publicados[0].ExpiresMinutes)
	assert.Empty(t, e.email.Enviados)
	assert.WithinDuration(t, time.Now().Add(15*time.Minute), e.repo.TokensCreados[0].ExpiresAt, time.Minute)
}

func TestDemoResend_ActualizaElTelefonoSiElUsuarioLoCorrige(t *testing.T) {
	e := nuevoEntorno(t)
	pendiente(e, func(u *domain.PendingDemoUser) { u.Phone = "573009998877" })

	_, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{
		Email: "demo@tienda.com", Channel: "whatsapp", Phone: "573001112233",
	})

	require.NoError(t, err)
	require.Len(t, e.repo.TelefonosCambiados, 1)
	assert.Equal(t, "573001112233", e.repo.TelefonosCambiados[0].Phone,
		"si alguien se equivoco de numero al registrarse, el reenvio es el unico modo de corregirlo sin crear otra cuenta")
}

func TestDemoResend_NoTocaElTelefonoSiEsElMismo(t *testing.T) {
	e := nuevoEntorno(t)
	pendiente(e, func(u *domain.PendingDemoUser) { u.Phone = "573001112233" })

	_, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{
		Email: "demo@tienda.com", Channel: "whatsapp", Phone: "573001112233",
	})

	require.NoError(t, err)
	assert.Empty(t, e.repo.TelefonosCambiados)
}

func TestDemoResend_SiFallaGuardarElTelefonoIgualSeMandaElCodigo(t *testing.T) {
	e := nuevoEntorno(t)
	pendiente(e, func(u *domain.PendingDemoUser) { u.Phone = "otro" })
	e.repo.UpdateUserPhoneFn = func(ctx context.Context, userID uint, phone string) error {
		return errors.New("update fallido")
	}

	_, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{
		Email: "demo@tienda.com", Channel: "whatsapp", Phone: "573001112233",
	})

	require.NoError(t, err)
	assert.Len(t, e.otp.Publicados, 1,
		"el codigo va al telefono que el usuario acaba de escribir aunque la columna no se haya podido actualizar")
}

func TestDemoResend_WhatsAppSinTelefonoSiDevuelveError(t *testing.T) {
	e := nuevoEntorno(t)
	pendiente(e, nil)

	_, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{
		Email: "demo@tienda.com", Channel: "whatsapp", Phone: "  ",
	})

	require.Error(t, err,
		"es el unico caso que rompe el mensaje generico, pero no filtra nada de la cuenta: es un error de forma del request")
	assert.Empty(t, e.repo.TokensCreados)
}

func TestDemoResend_UnFalloAlEscribirElTokenNoSeLeCuentaAlUsuario(t *testing.T) {
	casos := map[string]func(*entorno){
		"falla invalidar los tokens previos": func(e *entorno) {
			e.repo.InvalidateEmailVerificationFn = func(ctx context.Context, userID uint) error {
				return errors.New("timeout")
			}
		},
		"falla crear el token nuevo": func(e *entorno) {
			e.repo.CreateEmailVerificationTokenFn = func(ctx context.Context, userID uint, tokenHash string, expiresAt time.Time) error {
				return errors.New("timeout")
			}
		},
		"falla el envio del correo": func(e *entorno) {
			e.email.SendHTMLFn = func(ctx context.Context, to, subject, html string) error {
				return errors.New("SMTP 550")
			}
		},
		"falla la publicacion del OTP": func(e *entorno) {
			e.otp.PublishDemoOTPFn = func(ctx context.Context, event domain.DemoOTPEvent) error {
				return errors.New("rabbitmq caido")
			}
		},
	}

	for nombre, romper := range casos {
		t.Run(nombre, func(t *testing.T) {
			e := nuevoEntorno(t)
			pendiente(e, func(u *domain.PendingDemoUser) { u.Phone = "573001112233" })
			romper(e)
			req := domain.DemoResendRequest{Email: "demo@tienda.com"}
			if nombre == "falla la publicacion del OTP" {
				req.Channel = "whatsapp"
				req.Phone = "573001112233"
			}

			resp, err := e.uc.DemoResend(context.Background(), req)

			require.NoError(t, err)
			assert.Equal(t, mensajeGenerico, resp.Message,
				"el cliente ve exito aunque el reenvio no haya llegado; solo el log del backend registra la falla")
		})
	}
}

func TestDemoResend_ElTokenDelCorreoNoSeGuardaEnClaro(t *testing.T) {
	e := nuevoEntorno(t)
	pendiente(e, nil)

	_, err := e.uc.DemoResend(context.Background(), domain.DemoResendRequest{Email: "demo@tienda.com"})

	require.NoError(t, err)
	assert.Len(t, e.repo.TokensCreados[0].TokenHash, 64)
	assert.NotContains(t, e.email.Enviados[0].HTML, e.repo.TokensCreados[0].TokenHash)
}
