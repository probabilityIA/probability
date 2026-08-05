package usecasemanifest

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCarriers_MapeaLasOpciones(t *testing.T) {
	var pedidoBusiness uint
	var pedidoChildren bool
	repo := &mocks.RepositoryMock{
		ListPendingCarriersFn: func(ctx context.Context, businessID uint, includeChildren bool) ([]domain.ManifestCarrierCount, error) {
			pedidoBusiness, pedidoChildren = businessID, includeChildren
			return []domain.ManifestCarrierCount{
				{Carrier: "Servientrega", Count: 12},
				{Carrier: "Interrapidisimo", Count: 5},
			}, nil
		},
	}
	uc := New(repo)

	got, err := uc.ListCarriers(context.Background(), 36, true)

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "Servientrega", got[0].Carrier)
	assert.Equal(t, int64(12), got[0].Count)
	assert.Equal(t, uint(36), pedidoBusiness)
	assert.True(t, pedidoChildren)
}

func TestListCarriers_SinResultados_RetornaListaVacia(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListPendingCarriersFn: func(ctx context.Context, businessID uint, includeChildren bool) ([]domain.ManifestCarrierCount, error) {
			return nil, nil
		},
	}
	uc := New(repo)

	got, err := uc.ListCarriers(context.Background(), 36, false)

	require.NoError(t, err)
	assert.NotNil(t, got)
	assert.Empty(t, got)
}

func TestListCarriers_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		ListPendingCarriersFn: func(ctx context.Context, businessID uint, includeChildren bool) ([]domain.ManifestCarrierCount, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.ListCarriers(context.Background(), 36, false)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListPending_NormalizaPaginacion(t *testing.T) {
	cases := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{name: "pagina cero", page: 0, pageSize: 25, wantPage: 1, wantPageSize: 25},
		{name: "pagina negativa", page: -3, pageSize: 25, wantPage: 1, wantPageSize: 25},
		{name: "page size cero usa 25", page: 1, pageSize: 0, wantPage: 1, wantPageSize: 25},
		{name: "page size negativo usa 25", page: 1, pageSize: -1, wantPage: 1, wantPageSize: 25},
		{name: "page size sobre el tope usa 200", page: 1, pageSize: 500, wantPage: 1, wantPageSize: 200},
		{name: "page size en el tope", page: 1, pageSize: 200, wantPage: 1, wantPageSize: 200},
		{name: "valores validos", page: 4, pageSize: 50, wantPage: 4, wantPageSize: 50},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var recibido domain.ManifestFilter
			repo := &mocks.RepositoryMock{
				ListPendingForManifestFn: func(ctx context.Context, f domain.ManifestFilter) ([]domain.ManifestShipmentRow, int64, error) {
					recibido = f
					return nil, 0, nil
				},
			}
			uc := New(repo)

			got, err := uc.ListPending(context.Background(), PendingFilter{
				BusinessID: 36,
				Page:       tc.page,
				PageSize:   tc.pageSize,
			})

			require.NoError(t, err)
			assert.Equal(t, tc.wantPage, recibido.Page)
			assert.Equal(t, tc.wantPageSize, recibido.PageSize)
			assert.Equal(t, tc.wantPage, got.Page)
			assert.Equal(t, tc.wantPageSize, got.PageSize)
		})
	}
}

func TestListPending_PropagaLosFiltros(t *testing.T) {
	var recibido domain.ManifestFilter
	repo := &mocks.RepositoryMock{
		ListPendingForManifestFn: func(ctx context.Context, f domain.ManifestFilter) ([]domain.ManifestShipmentRow, int64, error) {
			recibido = f
			return nil, 0, nil
		},
	}
	uc := New(repo)

	_, err := uc.ListPending(context.Background(), PendingFilter{
		BusinessID:      36,
		IncludeChildren: true,
		Carrier:         "Servientrega",
		OnlyWithGuide:   true,
		Page:            2,
		PageSize:        50,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(36), recibido.BusinessID)
	assert.True(t, recibido.IncludeChildren)
	assert.Equal(t, "Servientrega", recibido.Carrier)
	assert.True(t, recibido.OnlyWithGuide)
}

func TestListPending_TotalPaginasMinimoEsUno(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListPendingForManifestFn: func(ctx context.Context, f domain.ManifestFilter) ([]domain.ManifestShipmentRow, int64, error) {
			return nil, 0, nil
		},
	}
	uc := New(repo)

	got, err := uc.ListPending(context.Background(), PendingFilter{BusinessID: 36})

	require.NoError(t, err)
	assert.Equal(t, 1, got.TotalPages)
	assert.Equal(t, int64(0), got.Total)
	assert.Empty(t, got.Data)
}

func TestListPending_CalculaTotalDePaginas(t *testing.T) {
	cases := []struct {
		total     int64
		pageSize  int
		wantPages int
	}{
		{total: 1, pageSize: 25, wantPages: 1},
		{total: 25, pageSize: 25, wantPages: 1},
		{total: 26, pageSize: 25, wantPages: 2},
		{total: 100, pageSize: 25, wantPages: 4},
		{total: 101, pageSize: 25, wantPages: 5},
	}

	for _, tc := range cases {
		repo := &mocks.RepositoryMock{
			ListPendingForManifestFn: func(ctx context.Context, f domain.ManifestFilter) ([]domain.ManifestShipmentRow, int64, error) {
				return nil, tc.total, nil
			},
		}
		uc := New(repo)

		got, err := uc.ListPending(context.Background(), PendingFilter{BusinessID: 36, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPages, got.TotalPages)
	}
}

func TestListPending_MapeaLasFilas(t *testing.T) {
	orderID := "order-uuid"
	repo := &mocks.RepositoryMock{
		ListPendingForManifestFn: func(ctx context.Context, f domain.ManifestFilter) ([]domain.ManifestShipmentRow, int64, error) {
			return []domain.ManifestShipmentRow{
				{
					ShipmentID:         34678,
					OrderID:            &orderID,
					OrderNumber:        "MYS-0302",
					TrackingNumber:     "034057376067",
					Carrier:            "Servientrega",
					CarrierCode:        "servientrega",
					CustomerName:       "Ana Gomez",
					CustomerDocument:   "1020304050",
					DestinationAddress: "Calle 1 # 2-3",
					DestinationCity:    "Bogota",
					DestinationState:   "Cundinamarca",
					Weight:             1.5,
					DeclaredValue:      120000,
					CodTotal:           150000,
					BusinessID:         36,
					BusinessName:       "Mi Tienda",
					ShipmentStatus:     "pending",
					OrderStatus:        "confirmed",
				},
			}, 1, nil
		},
	}
	uc := New(repo)

	got, err := uc.ListPending(context.Background(), PendingFilter{BusinessID: 36, Carrier: "Servientrega"})

	require.NoError(t, err)
	require.Len(t, got.Data, 1)

	row := got.Data[0]
	assert.Equal(t, uint(34678), row.ShipmentID)
	assert.Equal(t, "order-uuid", *row.OrderID)
	assert.Equal(t, "MYS-0302", row.OrderNumber)
	assert.Equal(t, "034057376067", row.TrackingNumber)
	assert.Equal(t, "Servientrega", row.Carrier)
	assert.Equal(t, "Ana Gomez", row.CustomerName)
	assert.Equal(t, "1020304050", row.CustomerDocument)
	assert.Equal(t, "Bogota", row.DestinationCity)
	assert.InDelta(t, 1.5, row.Weight, 0.001)
	assert.InDelta(t, 150000.0, row.CodTotal, 0.001)
	assert.Equal(t, "Mi Tienda", row.BusinessName)
	assert.Equal(t, "Servientrega", got.Carrier)
}

func TestListPending_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("query failed")
	repo := &mocks.RepositoryMock{
		ListPendingForManifestFn: func(ctx context.Context, f domain.ManifestFilter) ([]domain.ManifestShipmentRow, int64, error) {
			return nil, 0, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.ListPending(context.Background(), PendingFilter{BusinessID: 36})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGeneratePDF_SinBusinessID_RetornaError(t *testing.T) {
	uc := New(&mocks.RepositoryMock{})

	got, err := uc.GeneratePDF(context.Background(), GeneratePDFInput{ShipmentIDs: []uint{1}})

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "business_id requerido")
}

func TestGeneratePDF_SinEnvios_RetornaError(t *testing.T) {
	uc := New(&mocks.RepositoryMock{})

	got, err := uc.GeneratePDF(context.Background(), GeneratePDFInput{BusinessID: 36})

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "no hay envios seleccionados")
}

func TestTruncate(t *testing.T) {
	cases := []struct {
		in   string
		n    int
		want string
	}{
		{in: "corto", n: 10, want: "corto"},
		{in: "exacto", n: 6, want: "exacto"},
		{in: "demasiado largo", n: 8, want: "demasia."},
		{in: "", n: 5, want: ""},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.want, truncate(tc.in, tc.n))
	}
}

func TestStrOr(t *testing.T) {
	assert.Equal(t, "primero", strOr("primero", "segundo"))
	assert.Equal(t, "segundo", strOr("", "segundo"))
	assert.Equal(t, "", strOr("", ""))
}

func TestCarrierKey(t *testing.T) {
	cases := map[string]string{
		"servientrega":     "SERVIENTREGA",
		"Inter Rapidisimo": "INTERRAPIDISIMO",
		"dhl-express":      "DHLEXPRESS",
		"mi_paquete":       "MIPAQUETE",
		"":                 "",
		"   ":              "",
	}

	for in, want := range cases {
		assert.Equal(t, want, carrierKey(in))
	}
}

func TestGetCarrierLogoPNG_CarrierDesconocido_RetornaNil(t *testing.T) {
	assert.Nil(t, getCarrierLogoPNG("inventada"))
	assert.Nil(t, getCarrierLogoPNG(""))
}
