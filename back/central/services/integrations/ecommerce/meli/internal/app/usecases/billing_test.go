package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type mockBillingClient struct {
	domain.IMeliClient

	tokensUsados []string
	errores      []error
	info         *domain.MeliBillingInfo
}

func (m *mockBillingClient) WithBaseURL(string) domain.IMeliClient { return m }
func (m *mockBillingClient) ForAccount(string) domain.IMeliClient  { return m }

func (m *mockBillingClient) GetBillingInfo(_ context.Context, accessToken string, _ int64) (*domain.MeliBillingInfo, error) {
	idx := len(m.tokensUsados)
	m.tokensUsados = append(m.tokensUsados, accessToken)
	if idx < len(m.errores) && m.errores[idx] != nil {
		return nil, m.errores[idx]
	}
	return m.info, nil
}

func servicioConTokenVigente(token string) *mockService {
	return &mockService{
		integration: &domain.Integration{ID: 7, Config: futureToken()},
		token:       token,
	}
}

func TestTokenVencidoSeRefrescaYRecuperaElDocumento(t *testing.T) {
	cli := &mockBillingClient{
		errores: []error{domain.ErrTokenExpired},
		info:    &domain.MeliBillingInfo{DocNumber: "1020304050"},
	}
	svc := servicioConTokenVigente("token-nuevo")
	uc := newUseCase(cli, svc, nil)

	orden := &domain.MeliOrder{ID: 2000018097840680}
	uc.enrichBillingInfo(context.Background(), "7", "token-viejo", orden)

	if len(cli.tokensUsados) != 2 {
		t.Fatalf("esperaba un reintento tras refrescar, hubo %d llamadas", len(cli.tokensUsados))
	}
	if cli.tokensUsados[1] != "token-nuevo" {
		t.Fatalf("el reintento debia usar el token refrescado, uso %q", cli.tokensUsados[1])
	}
	if orden.Buyer.BillingInfo == nil || orden.Buyer.BillingInfo.DocNumber != "1020304050" {
		t.Fatal("la orden debia quedar con el documento del comprador, que es lo que se necesita para facturar")
	}
}

func TestSinTokenVencidoNoSeRefresca(t *testing.T) {
	cli := &mockBillingClient{info: &domain.MeliBillingInfo{DocNumber: "999"}}
	svc := servicioConTokenVigente("otro-token")
	uc := newUseCase(cli, svc, nil)

	uc.enrichBillingInfo(context.Background(), "7", "token-vigente", &domain.MeliOrder{ID: 1})

	if len(cli.tokensUsados) != 1 {
		t.Fatalf("esperaba una sola llamada, hubo %d", len(cli.tokensUsados))
	}
	if svc.decryptCalls != 0 {
		t.Fatalf("no debia tocar el token, lo desencripto %d veces", svc.decryptCalls)
	}
}

func TestOrdenQueYaTraeDocumentoNoLlamaAlCanal(t *testing.T) {
	cli := &mockBillingClient{}
	uc := newUseCase(cli, servicioConTokenVigente("t"), nil)

	orden := &domain.MeliOrder{ID: 1}
	orden.Buyer.BillingInfo = &domain.MeliBillingInfo{DocNumber: "ya-esta"}

	uc.enrichBillingInfo(context.Background(), "7", "token", orden)

	if len(cli.tokensUsados) != 0 {
		t.Fatalf("no debia consultar el canal, hizo %d llamadas", len(cli.tokensUsados))
	}
}

type mockLabelClient struct {
	domain.IMeliClient

	tokensUsados []string
	errores      []error
	label        *domain.ShipmentLabel
}

func (m *mockLabelClient) WithBaseURL(string) domain.IMeliClient { return m }
func (m *mockLabelClient) ForAccount(string) domain.IMeliClient  { return m }

func (m *mockLabelClient) GetShipmentLabel(_ context.Context, accessToken string, _ int64, _ string) (*domain.ShipmentLabel, error) {
	idx := len(m.tokensUsados)
	m.tokensUsados = append(m.tokensUsados, accessToken)
	if idx < len(m.errores) && m.errores[idx] != nil {
		return nil, m.errores[idx]
	}
	return m.label, nil
}

type mockLabelLookup struct {
	domain.IOrderLookupRepository
	ref *domain.MeliLabelRef
}

func (m *mockLabelLookup) GetMeliLabelRefByShipmentID(context.Context, uint) (*domain.MeliLabelRef, error) {
	return m.ref, nil
}

func usecaseDeEtiquetas(cli domain.IMeliClient, svc domain.IIntegrationService) *meliUseCase {
	negocio := uint(26)
	return &meliUseCase{
		client:  cli,
		service: svc,
		logger:  testLogger(),
		orderLookupRepo: &mockLabelLookup{ref: &domain.MeliLabelRef{
			BusinessID: negocio, IntegrationID: 7, MeliShipmentID: 47846252129,
		}},
	}
}

func TestEtiquetaConEnvioYaDespachadoNoEsErrorInterno(t *testing.T) {
	cli := &mockLabelClient{errores: []error{domain.ErrLabelAlreadyShipped}}
	svc := &mockService{
		integration: &domain.Integration{ID: 7, Config: futureToken(), BusinessID: uintPtr(26)},
		token:       "token",
	}

	_, err := usecaseDeEtiquetas(cli, svc).GetShipmentLabel(context.Background(), 1, 26, "pdf")

	if !errors.Is(err, domain.ErrLabelAlreadyShipped) {
		t.Fatalf("esperaba ErrLabelAlreadyShipped, obtuve %v", err)
	}
	if len(cli.tokensUsados) != 1 {
		t.Fatalf("no debia reintentar un envio ya despachado, hizo %d llamadas", len(cli.tokensUsados))
	}
}

func uintPtr(v uint) *uint { return &v }
