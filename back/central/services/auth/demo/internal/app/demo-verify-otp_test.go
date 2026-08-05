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

func TestDemoVerifyOTP_ActivaLaCuentaYAprovisiona(t *testing.T) {
	e := nuevoEntorno(t)
	conToken(e, tokenVigente())

	resp, err := e.uc.DemoVerifyOTP(context.Background(), domain.DemoVerifyOTPRequest{
		Email: "demo@tienda.com", Code: "123456",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Len(t, e.repo.Activaciones, 1)
	assert.Len(t, e.repo.Aprovisionamientos, 1)
}

func TestDemoVerifyOTP_BuscaPorElHashDeCorreoMasCodigo(t *testing.T) {
	e := nuevoEntorno(t)
	recibido := conToken(e, tokenVigente())

	_, err := e.uc.DemoVerifyOTP(context.Background(), domain.DemoVerifyOTPRequest{
		Email: "  DEMO@Tienda.COM ", Code: "  123456  ",
	})

	require.NoError(t, err)
	assert.Equal(t, hashOTP("demo@tienda.com", "123456"), *recibido,
		"correo y codigo se normalizan igual que en el registro, si no el hash nunca coincidiria")
}

func TestDemoVerifyOTP_TodoFalloDevuelveElMismoMensajeGenerico(t *testing.T) {
	usado := time.Now()
	casos := map[string]func(*entorno) domain.DemoVerifyOTPRequest{
		"sin correo": func(e *entorno) domain.DemoVerifyOTPRequest {
			return domain.DemoVerifyOTPRequest{Email: "  ", Code: "123456"}
		},
		"sin codigo": func(e *entorno) domain.DemoVerifyOTPRequest {
			return domain.DemoVerifyOTPRequest{Email: "demo@tienda.com", Code: " "}
		},
		"codigo que no existe": func(e *entorno) domain.DemoVerifyOTPRequest {
			conToken(e, nil)
			return domain.DemoVerifyOTPRequest{Email: "demo@tienda.com", Code: "000000"}
		},
		"codigo vencido": func(e *entorno) domain.DemoVerifyOTPRequest {
			conToken(e, &domain.EmailVerificationTokenInfo{ID: 1, UserID: 7, ExpiresAt: time.Now().Add(-time.Second)})
			return domain.DemoVerifyOTPRequest{Email: "demo@tienda.com", Code: "123456"}
		},
		"falla la consulta": func(e *entorno) domain.DemoVerifyOTPRequest {
			e.repo.GetValidEmailVerificationTokenFn = func(ctx context.Context, tokenHash string) (*domain.EmailVerificationTokenInfo, error) {
				return nil, errors.New("timeout")
			}
			return domain.DemoVerifyOTPRequest{Email: "demo@tienda.com", Code: "123456"}
		},
	}

	for nombre, preparar := range casos {
		t.Run(nombre, func(t *testing.T) {
			e := nuevoEntorno(t)
			req := preparar(e)

			resp, err := e.uc.DemoVerifyOTP(context.Background(), req)

			require.NoError(t, err,
				"nunca devuelve error de transporte: el handler responde 200 con success=false")
			assert.False(t, resp.Success)
			assert.Equal(t, "Codigo invalido o expirado", resp.Message,
				"un mensaje unico evita que se pueda distinguir codigo errado de codigo vencido a fuerza de intentos")
			assert.Empty(t, e.repo.Activaciones)
		})
	}
	_ = usado
}

func TestDemoVerifyOTP_NoHayLimiteDeIntentosContraFuerzaBruta(t *testing.T) {
	e := nuevoEntorno(t)
	consultas := 0
	e.repo.GetValidEmailVerificationTokenFn = func(ctx context.Context, tokenHash string) (*domain.EmailVerificationTokenInfo, error) {
		consultas++
		return nil, nil
	}

	for i := 0; i < 50; i++ {
		resp, err := e.uc.DemoVerifyOTP(context.Background(), domain.DemoVerifyOTPRequest{
			Email: "demo@tienda.com", Code: "000000",
		})
		require.NoError(t, err)
		require.False(t, resp.Success)
	}

	assert.Equal(t, 50, consultas,
		"el caso de uso no bloquea ni cuenta intentos: un OTP de 6 digitos con ventana de 15 min solo esta protegido por el rate limit del gateway")
}

func TestDemoVerifyOTP_UnCodigoYaUsadoResponderExitoSinReactivar(t *testing.T) {
	usado := time.Now().Add(-time.Hour)
	e := nuevoEntorno(t)
	conToken(e, &domain.EmailVerificationTokenInfo{
		ID: 11, UserID: 7, ExpiresAt: time.Now().Add(-time.Hour), UsedAt: &usado,
	})

	resp, err := e.uc.DemoVerifyOTP(context.Background(), domain.DemoVerifyOTPRequest{
		Email: "demo@tienda.com", Code: "123456",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Empty(t, e.repo.Activaciones)
	assert.Len(t, e.repo.Aprovisionamientos, 1)
}

func TestDemoVerifyOTP_SiFallaLaActivacionAvisaConMensajePropio(t *testing.T) {
	e := nuevoEntorno(t)
	conToken(e, tokenVigente())
	e.repo.ActivateUserAndConsumeTokenFn = func(ctx context.Context, tokenID, userID uint) error {
		return errors.New("deadlock")
	}

	resp, err := e.uc.DemoVerifyOTP(context.Background(), domain.DemoVerifyOTPRequest{
		Email: "demo@tienda.com", Code: "123456",
	})

	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "No se pudo verificar la cuenta", resp.Message,
		"un fallo de base si se distingue del codigo errado: el usuario tiene que saber que puede reintentar")
	assert.Empty(t, e.repo.Aprovisionamientos)
}
