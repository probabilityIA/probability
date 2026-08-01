package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type fakeClient struct {
	token     string
	tokenErr  error
	shops     []domain.ShopInfo
	queryErr  error
	lastShop  string
	lastLimit int
}

func (f *fakeClient) GetClientToken(_ context.Context, _, _, _ string) (string, error) {
	return f.token, f.tokenErr
}

func (f *fakeClient) QueryShopInfo(_ context.Context, _ domain.Credential, shopName string, limit int) ([]domain.ShopInfo, error) {
	f.lastShop = shopName
	f.lastLimit = limit
	return f.shops, f.queryErr
}

type fakeService struct {
	integration *domain.Integration
	creds       map[string]string
}

func (f *fakeService) GetIntegrationByID(_ context.Context, _ string) (*domain.Integration, error) {
	return f.integration, nil
}

func (f *fakeService) DecryptCredential(_ context.Context, _ string, field string) (string, error) {
	return f.creds[field], nil
}

type fakePublisher struct {
	published []*domain.ShopSnapshot
	err       error
}

func (f *fakePublisher) PublishSnapshot(_ context.Context, s *domain.ShopSnapshot) error {
	if f.err != nil {
		return f.err
	}
	f.published = append(f.published, s)
	return nil
}

type fakeRepo struct {
	saved []*domain.ShopSnapshot
	err   error
}

func (f *fakeRepo) Save(_ context.Context, s *domain.ShopSnapshot) error {
	if f.err != nil {
		return f.err
	}
	s.ID = uint(len(f.saved) + 1)
	f.saved = append(f.saved, s)
	return nil
}

func (f *fakeRepo) ListByIntegration(_ context.Context, _, _ uint, _, _ int) ([]domain.ShopSnapshot, int64, error) {
	out := make([]domain.ShopSnapshot, 0, len(f.saved))
	for _, s := range f.saved {
		out = append(out, *s)
	}
	return out, int64(len(out)), nil
}

func businessID(id uint) *uint {
	return &id
}

func newUseCase(c *fakeClient, s *fakeService, p *fakePublisher, r *fakeRepo) usecases.ITikTokUseCase {
	return usecases.New(c, s, p, r, log.New())
}

func TestTestConnectionSuccess(t *testing.T) {
	c := &fakeClient{token: "tok", shops: []domain.ShopInfo{{ShopName: "mitienda", ShopID: 1}}}
	uc := newUseCase(c, &fakeService{}, &fakePublisher{}, &fakeRepo{})

	config := map[string]interface{}{"base_url": "http://mock", "shop_name": "mitienda"}
	creds := map[string]interface{}{"client_key": "k", "client_secret": "s"}

	if err := uc.TestConnection(context.Background(), config, creds); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.lastShop != "mitienda" {
		t.Fatalf("expected shop query for mitienda, got %q", c.lastShop)
	}
}

func TestTestConnectionMissingCredentials(t *testing.T) {
	uc := newUseCase(&fakeClient{token: "tok"}, &fakeService{}, &fakePublisher{}, &fakeRepo{})

	config := map[string]interface{}{"base_url": "http://mock"}

	err := uc.TestConnection(context.Background(), config, map[string]interface{}{"client_secret": "s"})
	if !errors.Is(err, domain.ErrMissingClientKey) {
		t.Fatalf("expected ErrMissingClientKey, got %v", err)
	}

	err = uc.TestConnection(context.Background(), config, map[string]interface{}{"client_key": "k"})
	if !errors.Is(err, domain.ErrMissingClientSecret) {
		t.Fatalf("expected ErrMissingClientSecret, got %v", err)
	}
}

func TestTestConnectionTestingModeUsesTestURL(t *testing.T) {
	uc := newUseCase(&fakeClient{token: "tok"}, &fakeService{}, &fakePublisher{}, &fakeRepo{})

	config := map[string]interface{}{"is_testing": true}
	creds := map[string]interface{}{"client_key": "k", "client_secret": "s"}

	err := uc.TestConnection(context.Background(), config, creds)
	if !errors.Is(err, domain.ErrMissingBaseURLTest) {
		t.Fatalf("expected ErrMissingBaseURLTest, got %v", err)
	}
}

func TestTestConnectionShopNotFound(t *testing.T) {
	c := &fakeClient{token: "tok", shops: []domain.ShopInfo{}}
	uc := newUseCase(c, &fakeService{}, &fakePublisher{}, &fakeRepo{})

	config := map[string]interface{}{"base_url": "http://mock", "shop_name": "no-existe"}
	creds := map[string]interface{}{"client_key": "k", "client_secret": "s"}

	err := uc.TestConnection(context.Background(), config, creds)
	if !errors.Is(err, domain.ErrShopNotFound) {
		t.Fatalf("expected ErrShopNotFound, got %v", err)
	}
}

func TestSyncShopInfoPublishesSnapshot(t *testing.T) {
	c := &fakeClient{token: "tok", shops: []domain.ShopInfo{{
		ShopID:               99,
		ShopName:             "mitienda",
		ShopRating:           "4.8",
		ShopReviewCount:      10,
		ItemSoldCount:        20,
		ShopPerformanceValue: 70,
	}}}
	s := &fakeService{
		integration: &domain.Integration{
			ID:         5,
			BusinessID: businessID(7),
			Config:     map[string]interface{}{"shop_name": "mitienda"},
			BaseURL:    "http://mock",
		},
		creds: map[string]string{"client_key": "k", "client_secret": "s"},
	}
	p := &fakePublisher{}
	uc := newUseCase(c, s, p, &fakeRepo{})

	snapshot, err := uc.SyncShopInfo(context.Background(), "5", 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snapshot.ShopID != 99 || snapshot.BusinessID != 7 || snapshot.IntegrationID != 5 {
		t.Fatalf("unexpected snapshot: %+v", snapshot)
	}
	if len(p.published) != 1 {
		t.Fatalf("expected 1 published snapshot, got %d", len(p.published))
	}
}

func TestSyncShopInfoRejectsOtherBusiness(t *testing.T) {
	s := &fakeService{
		integration: &domain.Integration{
			ID:         5,
			BusinessID: businessID(7),
			Config:     map[string]interface{}{"shop_name": "mitienda"},
			BaseURL:    "http://mock",
		},
		creds: map[string]string{"client_key": "k", "client_secret": "s"},
	}
	uc := newUseCase(&fakeClient{token: "tok"}, s, &fakePublisher{}, &fakeRepo{})

	_, err := uc.SyncShopInfo(context.Background(), "5", 99)
	if !errors.Is(err, domain.ErrIntegrationNotFound) {
		t.Fatalf("expected ErrIntegrationNotFound for foreign business, got %v", err)
	}
}

func TestSaveSnapshotValidation(t *testing.T) {
	r := &fakeRepo{}
	uc := newUseCase(&fakeClient{}, &fakeService{}, &fakePublisher{}, r)

	if err := uc.SaveSnapshot(context.Background(), &domain.ShopSnapshot{}); err == nil {
		t.Fatal("expected error for snapshot without business/integration")
	}

	snap := &domain.ShopSnapshot{BusinessID: 1, IntegrationID: 2, ShopName: "mitienda"}
	if err := uc.SaveSnapshot(context.Background(), snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.saved) != 1 {
		t.Fatalf("expected 1 saved snapshot, got %d", len(r.saved))
	}
}
