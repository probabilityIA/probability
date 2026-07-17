package usecases

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type fakeWebhookService struct {
	integration *domain.Integration
}

func (f *fakeWebhookService) GetIntegrationByID(ctx context.Context, integrationID string) (*domain.Integration, error) {
	return f.integration, nil
}

func (f *fakeWebhookService) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	return "valor-" + fieldName, nil
}

func (f *fakeWebhookService) UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error {
	return nil
}

type fakeWebhookClient struct {
	domain.IJumpsellerClient
	hooks        []domain.WebhookItem
	listErr      error
	createErr    map[string]error
	createdOrder []string
}

func (f *fakeWebhookClient) GetStoreInfo(ctx context.Context, cred domain.Credential) (*domain.StoreInfo, error) {
	return &domain.StoreInfo{Code: "tienda-test"}, nil
}

func (f *fakeWebhookClient) ListHooks(ctx context.Context, cred domain.Credential) ([]domain.WebhookItem, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.hooks, nil
}

func (f *fakeWebhookClient) CreateHook(ctx context.Context, cred domain.Credential, event, url string) (string, error) {
	if err, bad := f.createErr[event]; bad {
		return "", err
	}
	f.createdOrder = append(f.createdOrder, event)
	return "hook-" + event, nil
}

func newWebhookUseCase(client domain.IJumpsellerClient) *jumpsellerUseCase {
	businessID := uint(7)
	return &jumpsellerUseCase{
		client: client,
		service: &fakeWebhookService{
			integration: &domain.Integration{
				ID:         232,
				BusinessID: &businessID,
				BaseURL:    "https://api.jumpseller.com/v1",
			},
		},
		logger: log.New(),
	}
}

func TestSoloRegistraEventosValidosDeJumpseller(t *testing.T) {
	eventosValidos := map[string]bool{
		domain.EventOrderPaid:           true,
		domain.EventOrderPendingPayment: true,
		domain.EventOrderShipped:        true,
		domain.EventOrderCanceled:       true,
		domain.EventOrderAbandoned:      true,
		domain.EventOrderUpdated:        true,
	}

	for _, event := range domain.WebhookOrderEvents {
		if !eventosValidos[event] {
			t.Errorf("WebhookOrderEvents contiene %q, que la API de Jumpseller rechaza con 'Evento no valido'", event)
		}
	}
}

func TestFallaFuerteCuandoJumpsellerRechazaUnWebhook(t *testing.T) {
	rechazo := errors.New("No se pudo conectar a la URL")
	client := &fakeWebhookClient{
		createErr: map[string]error{domain.EventOrderCanceled: rechazo},
	}
	uc := newWebhookUseCase(client)

	err := uc.CreateWebhooks(context.Background(), "232", "https://probabilityia.com.co")

	if err == nil {
		t.Fatal("se esperaba error cuando Jumpseller rechaza un webhook, no un exito silencioso")
	}
	if !errors.Is(err, domain.ErrWebhookCreationFailed) {
		t.Errorf("el error debe envolver ErrWebhookCreationFailed, se obtuvo: %v", err)
	}
	if !strings.Contains(err.Error(), domain.EventOrderCanceled) {
		t.Errorf("el error debe nombrar el evento que fallo, se obtuvo: %v", err)
	}
	if !errors.Is(err, rechazo) {
		t.Errorf("el error debe conservar la causa de Jumpseller, se obtuvo: %v", err)
	}
}

func TestFallaCuandoNoSePuedenListarLosWebhooksExistentes(t *testing.T) {
	client := &fakeWebhookClient{listErr: errors.New("timeout")}
	uc := newWebhookUseCase(client)

	err := uc.CreateWebhooks(context.Background(), "232", "https://probabilityia.com.co")

	if err == nil {
		t.Fatal("sin poder listar no se puede deduplicar: debe fallar en vez de crear webhooks duplicados")
	}
	if len(client.createdOrder) != 0 {
		t.Errorf("no debe crear ningun webhook si el listado fallo, creo: %v", client.createdOrder)
	}
}

func TestNoDuplicaWebhooksYaRegistrados(t *testing.T) {
	deliveryURL := WebhookDeliveryURL("https://probabilityia.com.co", 232)
	client := &fakeWebhookClient{
		hooks: []domain.WebhookItem{
			{ID: "1", Address: deliveryURL, Topic: domain.EventOrderPaid},
		},
	}
	uc := newWebhookUseCase(client)

	if err := uc.CreateWebhooks(context.Background(), "232", "https://probabilityia.com.co"); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	for _, event := range client.createdOrder {
		if event == domain.EventOrderPaid {
			t.Error("order_paid ya estaba registrado en esa URL y no debia recrearse")
		}
	}
	if len(client.createdOrder) != len(domain.WebhookOrderEvents)-1 {
		t.Errorf("se esperaban %d webhooks nuevos, se crearon %v", len(domain.WebhookOrderEvents)-1, client.createdOrder)
	}
}
