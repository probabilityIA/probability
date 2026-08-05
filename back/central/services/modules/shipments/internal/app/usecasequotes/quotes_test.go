package usecasequotes

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func s(v string) *string {
	return &v
}

func i64(v int64) *int64 {
	return &v
}

func savedQuote(id, businessID uint) *domain.SavedQuote {
	return &domain.SavedQuote{
		ID:         id,
		BusinessID: businessID,
		Status:     domain.QuoteStatusCreated,
	}
}

func TestAssociate_NoExiste_RetornaNotFound(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) { return nil, nil },
	}
	uc := New(repo)

	resp, err := uc.Associate(context.Background(), domain.AssociateQuoteInput{QuoteID: 1})

	assert.ErrorIs(t, err, ErrQuoteNotFound)
	assert.Nil(t, resp)
}

func TestAssociate_CotizacionDeOtroNegocio_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) {
			return savedQuote(id, 36), nil
		},
	}
	uc := New(repo)

	resp, err := uc.Associate(context.Background(), domain.AssociateQuoteInput{
		QuoteID:    1,
		BusinessID: 99,
		OrderUUID:  "order-uuid",
	})

	assert.ErrorIs(t, err, ErrQuoteBusinessMatch)
	assert.Nil(t, resp)
}

func TestAssociate_SuperAdminSinBusinessID_NoValidaPertenencia(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) {
			return savedQuote(id, 36), nil
		},
		UpdateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error { return nil },
	}
	uc := New(repo)

	resp, err := uc.Associate(context.Background(), domain.AssociateQuoteInput{
		QuoteID:    1,
		BusinessID: 0,
		OrderUUID:  "order-uuid",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestAssociate_YaAsociadaAOtraOrden_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) {
			q := savedQuote(id, 36)
			q.OrderUUID = s("otra-orden")
			return q, nil
		},
	}
	uc := New(repo)

	resp, err := uc.Associate(context.Background(), domain.AssociateQuoteInput{
		QuoteID:    1,
		BusinessID: 36,
		OrderUUID:  "order-uuid",
	})

	assert.ErrorIs(t, err, ErrQuoteAlreadyLinked)
	assert.Nil(t, resp)
}

func TestAssociate_MismaOrden_EsIdempotente(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) {
			q := savedQuote(id, 36)
			q.OrderUUID = s("order-uuid")
			return q, nil
		},
		UpdateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error { return nil },
	}
	uc := New(repo)

	resp, err := uc.Associate(context.Background(), domain.AssociateQuoteInput{
		QuoteID:    1,
		BusinessID: 36,
		OrderUUID:  "order-uuid",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, domain.QuoteStatusAssociated, resp.Status)
}

func TestAssociate_OrderUUIDVacioEnLaCotizacion_PermiteAsociar(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) {
			q := savedQuote(id, 36)
			q.OrderUUID = s("")
			return q, nil
		},
		UpdateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error { return nil },
	}
	uc := New(repo)

	resp, err := uc.Associate(context.Background(), domain.AssociateQuoteInput{
		QuoteID:    1,
		BusinessID: 36,
		OrderUUID:  "order-uuid",
	})

	require.NoError(t, err)
	assert.Equal(t, "order-uuid", *resp.OrderUUID)
}

func TestAssociate_ConGuiaSolicitada_MarcaGuideGenerated(t *testing.T) {
	var saved *domain.SavedQuote
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) {
			return savedQuote(id, 36), nil
		},
		UpdateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error {
			saved = quote
			return nil
		},
	}
	uc := New(repo)

	resp, err := uc.Associate(context.Background(), domain.AssociateQuoteInput{
		QuoteID:         1,
		BusinessID:      36,
		OrderUUID:       "order-uuid",
		SelectedCarrier: "Servientrega",
		SelectedIDRate:  i64(4521),
		GuideRequested:  true,
	})

	require.NoError(t, err)
	require.NotNil(t, saved)
	assert.Equal(t, domain.QuoteStatusGuideGenerated, saved.Status)
	assert.Equal(t, "Servientrega", saved.SelectedCarrier)
	require.NotNil(t, saved.SelectedIDRate)
	assert.Equal(t, int64(4521), *saved.SelectedIDRate)
	assert.Equal(t, domain.QuoteStatusGuideGenerated, resp.Status)
}

func TestAssociate_SinCarrierNiRate_NoSobreescribeLosExistentes(t *testing.T) {
	var saved *domain.SavedQuote
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) {
			q := savedQuote(id, 36)
			q.SelectedCarrier = "Interrapidisimo"
			q.SelectedIDRate = i64(111)
			return q, nil
		},
		UpdateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error {
			saved = quote
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.Associate(context.Background(), domain.AssociateQuoteInput{
		QuoteID:    1,
		BusinessID: 36,
		OrderUUID:  "order-uuid",
	})

	require.NoError(t, err)
	assert.Equal(t, "Interrapidisimo", saved.SelectedCarrier)
	assert.Equal(t, int64(111), *saved.SelectedIDRate)
}

func TestAssociate_FallaLeer_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) { return nil, dbErr },
	}
	uc := New(repo)

	resp, err := uc.Associate(context.Background(), domain.AssociateQuoteInput{QuoteID: 1})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

func TestAssociate_FallaGuardar_PropagaError(t *testing.T) {
	dbErr := errors.New("update failed")
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) {
			return savedQuote(id, 36), nil
		},
		UpdateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error { return dbErr },
	}
	uc := New(repo)

	resp, err := uc.Associate(context.Background(), domain.AssociateQuoteInput{
		QuoteID:    1,
		BusinessID: 36,
		OrderUUID:  "order-uuid",
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

func TestSaveQuote_TTLPorDefectoEs24Horas(t *testing.T) {
	var saved *domain.SavedQuote
	repo := &mocks.RepositoryMock{
		CreateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error {
			saved = quote
			return nil
		},
	}
	uc := New(repo)

	antes := time.Now()
	got, err := uc.SaveQuote(context.Background(), domain.SaveQuoteInput{BusinessID: 36})

	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, saved.ExpiresAt)
	assert.WithinDuration(t, antes.Add(24*time.Hour), *saved.ExpiresAt, 5*time.Second)
	assert.Equal(t, domain.QuoteStatusCreated, saved.Status)
}

func TestSaveQuote_TTLExplicitoSeRespeta(t *testing.T) {
	var saved *domain.SavedQuote
	repo := &mocks.RepositoryMock{
		CreateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error {
			saved = quote
			return nil
		},
	}
	uc := New(repo)

	antes := time.Now()
	_, err := uc.SaveQuote(context.Background(), domain.SaveQuoteInput{
		BusinessID: 36,
		TTL:        2 * time.Hour,
	})

	require.NoError(t, err)
	assert.WithinDuration(t, antes.Add(2*time.Hour), *saved.ExpiresAt, 5*time.Second)
}

func TestSaveQuote_TTLNegativo_UsaElPorDefecto(t *testing.T) {
	var saved *domain.SavedQuote
	repo := &mocks.RepositoryMock{
		CreateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error {
			saved = quote
			return nil
		},
	}
	uc := New(repo)

	antes := time.Now()
	_, err := uc.SaveQuote(context.Background(), domain.SaveQuoteInput{TTL: -1 * time.Hour})

	require.NoError(t, err)
	assert.WithinDuration(t, antes.Add(24*time.Hour), *saved.ExpiresAt, 5*time.Second)
}

func TestSaveQuote_CopiaLosCamposDeEntrada(t *testing.T) {
	var saved *domain.SavedQuote
	repo := &mocks.RepositoryMock{
		CreateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error {
			saved = quote
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.SaveQuote(context.Background(), domain.SaveQuoteInput{
		BusinessID:       36,
		IntegrationID:    197,
		Source:           "checkout",
		CorrelationID:    "corr-1",
		OrderUUID:        s("order-uuid"),
		ExternalOrderRef: "EXT-1",
		Rates:            []map[string]interface{}{{"carrier": "Servientrega"}},
	})

	require.NoError(t, err)
	assert.Equal(t, uint(36), saved.BusinessID)
	assert.Equal(t, uint(197), saved.IntegrationID)
	assert.Equal(t, "checkout", saved.Source)
	assert.Equal(t, "corr-1", saved.CorrelationID)
	assert.Equal(t, "EXT-1", saved.ExternalOrderRef)
	assert.Len(t, saved.Rates, 1)
}

func TestSaveQuote_FallaCrear_PropagaError(t *testing.T) {
	dbErr := errors.New("insert failed")
	repo := &mocks.RepositoryMock{
		CreateSavedQuoteFn: func(ctx context.Context, quote *domain.SavedQuote) error { return dbErr },
	}
	uc := New(repo)

	got, err := uc.SaveQuote(context.Background(), domain.SaveQuoteInput{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestList_NormalizaPaginacionYCalculaPaginas(t *testing.T) {
	cases := []struct {
		name         string
		page         int
		pageSize     int
		total        int64
		wantPage     int
		wantPageSize int
		wantPages    int
	}{
		{name: "pagina cero", page: 0, pageSize: 10, total: 25, wantPage: 1, wantPageSize: 10, wantPages: 3},
		{name: "page size cero", page: 2, pageSize: 0, total: 5, wantPage: 2, wantPageSize: 10, wantPages: 1},
		{name: "sin resultados", page: 1, pageSize: 10, total: 0, wantPage: 1, wantPageSize: 10, wantPages: 0},
		{name: "division exacta", page: 1, pageSize: 10, total: 30, wantPage: 1, wantPageSize: 10, wantPages: 3},
		{name: "valores validos", page: 3, pageSize: 20, total: 41, wantPage: 3, wantPageSize: 20, wantPages: 3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mocks.RepositoryMock{
				ListSavedQuotesFn: func(ctx context.Context, filter domain.SavedQuoteFilter) ([]domain.SavedQuote, int64, error) {
					return []domain.SavedQuote{}, tc.total, nil
				},
			}
			uc := New(repo)

			resp, err := uc.List(context.Background(), domain.SavedQuoteFilter{
				Page:     tc.page,
				PageSize: tc.pageSize,
			})

			require.NoError(t, err)
			assert.Equal(t, tc.wantPage, resp.Page)
			assert.Equal(t, tc.wantPageSize, resp.PageSize)
			assert.Equal(t, tc.wantPages, resp.TotalPages)
			assert.Equal(t, tc.total, resp.Total)
		})
	}
}

func TestList_MapeaCadaCotizacion(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListSavedQuotesFn: func(ctx context.Context, filter domain.SavedQuoteFilter) ([]domain.SavedQuote, int64, error) {
			return []domain.SavedQuote{
				{ID: 1, BusinessID: 36, Status: domain.QuoteStatusCreated},
				{ID: 2, BusinessID: 36, Status: domain.QuoteStatusAssociated},
			}, 2, nil
		},
	}
	uc := New(repo)

	resp, err := uc.List(context.Background(), domain.SavedQuoteFilter{Page: 1, PageSize: 10})

	require.NoError(t, err)
	require.Len(t, resp.Data, 2)
	assert.Equal(t, uint(1), resp.Data[0].ID)
	assert.Equal(t, domain.QuoteStatusAssociated, resp.Data[1].Status)
}

func TestList_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("query failed")
	repo := &mocks.RepositoryMock{
		ListSavedQuotesFn: func(ctx context.Context, filter domain.SavedQuoteFilter) ([]domain.SavedQuote, int64, error) {
			return nil, 0, dbErr
		},
	}
	uc := New(repo)

	resp, err := uc.List(context.Background(), domain.SavedQuoteFilter{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

func TestGetByID_NoExiste_RetornaNilSinError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) { return nil, nil },
	}
	uc := New(repo)

	resp, err := uc.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Nil(t, resp)
}

func TestGetByID_Existe_RetornaResponse(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) {
			q := savedQuote(id, 36)
			q.SelectedCarrier = "Servientrega"
			return q, nil
		},
	}
	uc := New(repo)

	resp, err := uc.GetByID(context.Background(), 9)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint(9), resp.ID)
	assert.Equal(t, "Servientrega", resp.SelectedCarrier)
}

func TestGetByID_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetSavedQuoteByIDFn: func(ctx context.Context, id uint) (*domain.SavedQuote, error) { return nil, dbErr },
	}
	uc := New(repo)

	resp, err := uc.GetByID(context.Background(), 9)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}
