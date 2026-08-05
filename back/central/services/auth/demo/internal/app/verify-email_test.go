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

func tokenVigente() *domain.EmailVerificationTokenInfo {
	return &domain.EmailVerificationTokenInfo{ID: 11, UserID: 7, ExpiresAt: time.Now().Add(time.Hour)}
}

func conToken(e *entorno, info *domain.EmailVerificationTokenInfo) *string {
	var recibido string
	e.repo.GetValidEmailVerificationTokenFn = func(ctx context.Context, tokenHash string) (*domain.EmailVerificationTokenInfo, error) {
		recibido = tokenHash
		return info, nil
	}
	return &recibido
}

func TestVerifyEmail_ActivaLaCuentaYAprovisionaLasIntegracionesDemo(t *testing.T) {
	e := nuevoEntorno(t)
	conToken(e, tokenVigente())

	resp, err := e.uc.VerifyEmail(context.Background(), domain.VerifyEmailRequest{Token: "token-crudo"})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	require.Len(t, e.repo.Activaciones, 1)
	assert.EqualValues(t, 11, e.repo.Activaciones[0].TokenID)
	assert.EqualValues(t, 7, e.repo.Activaciones[0].UserID)
	require.Len(t, e.repo.Aprovisionamientos, 1,
		"al verificar se siembran las integraciones simuladas para que la demo tenga datos desde el primer login")
	assert.EqualValues(t, 100, e.repo.Aprovisionamientos[0].BusinessID)
}

func TestVerifyEmail_BuscaElTokenPorSuHashNoPorElValorCrudo(t *testing.T) {
	e := nuevoEntorno(t)
	recibido := conToken(e, tokenVigente())

	_, err := e.uc.VerifyEmail(context.Background(), domain.VerifyEmailRequest{Token: "  token-crudo  "})

	require.NoError(t, err)
	assert.Equal(t, hashToken("token-crudo"), *recibido,
		"el token de la URL se recorta y se hashea antes de tocar la base")
}

func TestVerifyEmail_RechazaUnTokenVacio(t *testing.T) {
	for _, entrada := range []string{"", "   ", "\t"} {
		e := nuevoEntorno(t)

		_, err := e.uc.VerifyEmail(context.Background(), domain.VerifyEmailRequest{Token: entrada})

		require.Error(t, err)
		assert.Empty(t, e.repo.Activaciones)
	}
}

func TestVerifyEmail_RechazaUnTokenQueNoExiste(t *testing.T) {
	e := nuevoEntorno(t)
	conToken(e, nil)

	_, err := e.uc.VerifyEmail(context.Background(), domain.VerifyEmailRequest{Token: "inventado"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalido o expirado")
	assert.Empty(t, e.repo.Activaciones)
}

func TestVerifyEmail_RechazaUnTokenVencidoQueNuncaSeUso(t *testing.T) {
	e := nuevoEntorno(t)
	conToken(e, &domain.EmailVerificationTokenInfo{ID: 11, UserID: 7, ExpiresAt: time.Now().Add(-time.Minute)})

	_, err := e.uc.VerifyEmail(context.Background(), domain.VerifyEmailRequest{Token: "viejo"})

	require.Error(t, err)
	assert.Empty(t, e.repo.Activaciones)
	assert.Empty(t, e.repo.Aprovisionamientos,
		"un token vencido no dispara aprovisionamiento")
}

func TestVerifyEmail_UnTokenYaUsadoEsIdempotenteYNoReactivaLaCuenta(t *testing.T) {
	usado := time.Now().Add(-time.Hour)
	e := nuevoEntorno(t)
	conToken(e, &domain.EmailVerificationTokenInfo{
		ID: 11, UserID: 7, ExpiresAt: time.Now().Add(-30 * time.Minute), UsedAt: &usado,
	})

	resp, err := e.uc.VerifyEmail(context.Background(), domain.VerifyEmailRequest{Token: "ya-usado"})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.Message, "ya estaba verificada")
	assert.Empty(t, e.repo.Activaciones,
		"si ya se uso no se vuelve a activar; ademas la expiracion se ignora, asi que un clic en el enlace viejo sigue respondiendo exito")
}

func TestVerifyEmail_UnTokenYaUsadoIgualReintentaElAprovisionamiento(t *testing.T) {
	usado := time.Now()
	e := nuevoEntorno(t)
	conToken(e, &domain.EmailVerificationTokenInfo{ID: 11, UserID: 7, ExpiresAt: time.Now(), UsedAt: &usado})

	_, err := e.uc.VerifyEmail(context.Background(), domain.VerifyEmailRequest{Token: "ya-usado"})

	require.NoError(t, err)
	assert.Len(t, e.repo.Aprovisionamientos, 1,
		"reabrir el enlace rescata una demo cuyo aprovisionamiento fallo la primera vez")
}

func TestVerifyEmail_SiFallaLaConsultaNoFiltraElErrorInterno(t *testing.T) {
	e := nuevoEntorno(t)
	e.repo.GetValidEmailVerificationTokenFn = func(ctx context.Context, tokenHash string) (*domain.EmailVerificationTokenInfo, error) {
		return nil, errors.New("relation does not exist")
	}

	_, err := e.uc.VerifyEmail(context.Background(), domain.VerifyEmailRequest{Token: "x"})

	require.Error(t, err)
	assert.Equal(t, "error interno del servidor", err.Error())
}

func TestVerifyEmail_SiFallaLaActivacionNoSeAprovisiona(t *testing.T) {
	e := nuevoEntorno(t)
	conToken(e, tokenVigente())
	e.repo.ActivateUserAndConsumeTokenFn = func(ctx context.Context, tokenID, userID uint) error {
		return errors.New("deadlock detected")
	}

	_, err := e.uc.VerifyEmail(context.Background(), domain.VerifyEmailRequest{Token: "x"})

	require.Error(t, err)
	assert.Empty(t, e.repo.Aprovisionamientos)
}

func TestVerifyEmail_SiFallaElAprovisionamientoLaCuentaIgualQuedaVerificada(t *testing.T) {
	casos := map[string]func(*entorno){
		"no se resuelve el negocio": func(e *entorno) {
			e.repo.GetBusinessIDByUserIDFn = func(ctx context.Context, userID uint) (uint, error) {
				return 0, errors.New("no encontrado")
			}
		},
		"el negocio devuelve cero": func(e *entorno) {
			e.repo.GetBusinessIDByUserIDFn = func(ctx context.Context, userID uint) (uint, error) {
				return 0, nil
			}
		},
		"falla la siembra de integraciones": func(e *entorno) {
			e.repo.ProvisionDemoIntegrationsFn = func(ctx context.Context, businessID, userID uint) error {
				return errors.New("insert fallido")
			}
		},
	}

	for nombre, romper := range casos {
		t.Run(nombre, func(t *testing.T) {
			e := nuevoEntorno(t)
			conToken(e, tokenVigente())
			romper(e)

			resp, err := e.uc.VerifyEmail(context.Background(), domain.VerifyEmailRequest{Token: "x"})

			require.NoError(t, err,
				"el aprovisionamiento es best-effort: la cuenta queda activa aunque la demo arranque sin datos de ejemplo")
			assert.True(t, resp.Success)
			assert.Len(t, e.repo.Activaciones, 1)
		})
	}
}

func TestBuildVerificationEmail_SaludaConElNombreYLlevaElEnlaceDosVeces(t *testing.T) {
	html := buildVerificationEmail("Sebastian", "Tienda Probability", "https://app.com/verify-email?token=abc")

	assert.Contains(t, html, "Hola Sebastian,")
	assert.Contains(t, html, "Tienda Probability")
	assert.Equal(t, 2, contarOcurrencias(html, "https://app.com/verify-email?token=abc"),
		"el enlace va en el boton y en texto plano por si el cliente de correo bloquea el boton")
}

func TestBuildVerificationEmail_SinNombreSaludaGenerico(t *testing.T) {
	html := buildVerificationEmail("   ", "Tienda", "https://app.com/x")

	assert.Contains(t, html, "Hola,")
	assert.NotContains(t, html, "Hola ,")
}

func contarOcurrencias(s, sub string) int {
	n := 0
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			n++
		}
	}
	return n
}
