package app

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newStorefrontUseCase(repo *mocks.RepositoryMock, pub *mocks.PublisherMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	if pub == nil {
		pub = &mocks.PublisherMock{}
	}
	return New(repo, mocks.NewSilentLogger(), pub)
}

func strPtr(v string) *string { return &v }

func sinClientePrevio(repo *mocks.RepositoryMock) *mocks.RepositoryMock {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	if repo.GetClientByUserIDFn == nil {
		repo.GetClientByUserIDFn = func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
			return nil, stderrors.New("no encontrado")
		}
	}
	return repo
}

func repoTiendaApagada() *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		IsIntegrationActiveOrMissingFn: func(ctx context.Context, businessID, typeID uint) (bool, error) {
			return false, nil
		},
	}
}

func ordenValida() *dtos.StorefrontCreateOrderDTO {
	return &dtos.StorefrontCreateOrderDTO{
		Items: []dtos.StorefrontOrderItemDTO{{ProductID: "P-1", Quantity: 2}},
	}
}

func decodificarOrden(t *testing.T, raw []byte) map[string]interface{} {
	t.Helper()
	var out map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &out))
	return out
}

func TestTiendaApagada_BloqueaTodasLasOperacionesPublicas(t *testing.T) {
	t.Run("catalogo", func(t *testing.T) {
		_, _, err := newStorefrontUseCase(repoTiendaApagada(), nil).
			ListCatalog(context.Background(), 26, dtos.CatalogFilters{})
		assert.ErrorIs(t, err, domainerrors.ErrStorefrontNotActive)
	})

	t.Run("detalle de producto", func(t *testing.T) {
		_, err := newStorefrontUseCase(repoTiendaApagada(), nil).
			GetProduct(context.Background(), 26, "P-1")
		assert.ErrorIs(t, err, domainerrors.ErrStorefrontNotActive)
	})

	t.Run("crear orden", func(t *testing.T) {
		pub := &mocks.PublisherMock{}
		err := newStorefrontUseCase(repoTiendaApagada(), pub).
			CreateOrder(context.Background(), 26, 42, ordenValida())
		assert.ErrorIs(t, err, domainerrors.ErrStorefrontNotActive)
		assert.Empty(t, pub.Published, "una tienda apagada no debe poder emitir ordenes")
	})

	t.Run("historial de ordenes", func(t *testing.T) {
		_, _, err := newStorefrontUseCase(repoTiendaApagada(), nil).
			ListMyOrders(context.Background(), 26, 42, 1, 10)
		assert.ErrorIs(t, err, domainerrors.ErrStorefrontNotActive)
	})

	t.Run("detalle de orden", func(t *testing.T) {
		_, err := newStorefrontUseCase(repoTiendaApagada(), nil).
			GetMyOrder(context.Background(), "ORD-1", 26, 42)
		assert.ErrorIs(t, err, domainerrors.ErrStorefrontNotActive)
	})
}

func TestElGateConsultaElTipoDeIntegracionDeTienda(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	uc := newStorefrontUseCase(repo, nil)

	_, _, _ = uc.ListCatalog(context.Background(), 26, dtos.CatalogFilters{})
	_, _ = uc.GetProduct(context.Background(), 26, "P-1")
	_, _, _ = uc.ListMyOrders(context.Background(), 26, 42, 1, 10)
	_, _ = uc.GetMyOrder(context.Background(), "ORD-1", 26, 42)
	_ = uc.CreateOrder(context.Background(), 26, 42, ordenValida())

	require.Len(t, repo.GateChecks, 5)
	for i, tipo := range repo.GateChecks {
		assert.Equal(t, uint(tiendaIntegrationTypeID), tipo, "consulta %d", i)
	}
}

func TestElGate_ErrorSePropagaYNoConcedeAcceso(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		IsIntegrationActiveOrMissingFn: func(ctx context.Context, businessID, typeID uint) (bool, error) {
			return false, dbErr
		},
	}
	pub := &mocks.PublisherMock{}
	uc := newStorefrontUseCase(repo, pub)

	_, _, err := uc.ListCatalog(context.Background(), 26, dtos.CatalogFilters{})
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.GetProduct(context.Background(), 26, "P-1")
	assert.ErrorIs(t, err, dbErr)

	err = uc.CreateOrder(context.Background(), 26, 42, ordenValida())
	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, pub.Published)

	_, _, err = uc.ListMyOrders(context.Background(), 26, 42, 1, 10)
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.GetMyOrder(context.Background(), "ORD-1", 26, 42)
	assert.ErrorIs(t, err, dbErr)
}

func TestCatalogFilters_Normalize(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 12, 1, 12},
		{-5, 12, 1, 12},
		{1, 0, 1, 12},
		{1, -1, 1, 12},
		{1, 101, 1, 12},
		{1, 100, 1, 100},
		{3, 24, 3, 24},
	}

	for _, tc := range casos {
		f := dtos.CatalogFilters{Page: tc.page, PageSize: tc.pageSize}
		f.Normalize()
		assert.Equal(t, tc.wantPage, f.Page, "page=%d", tc.page)
		assert.Equal(t, tc.wantSize, f.PageSize, "pageSize=%d", tc.pageSize)
	}
}

func TestCatalogFilters_Offset(t *testing.T) {
	assert.Equal(t, 0, dtos.CatalogFilters{Page: 1, PageSize: 12}.Offset())
	assert.Equal(t, 12, dtos.CatalogFilters{Page: 2, PageSize: 12}.Offset())
	assert.Equal(t, 24, dtos.CatalogFilters{Page: 3, PageSize: 12}.Offset())
	assert.Equal(t, 0, dtos.CatalogFilters{Page: 0, PageSize: 12}.Offset(),
		"una pagina invalida no debe producir un offset negativo")
	assert.Equal(t, 0, dtos.CatalogFilters{Page: -3, PageSize: 12}.Offset())
}

func TestListCatalog_NormalizaAntesDeConsultar(t *testing.T) {
	var visto dtos.CatalogFilters
	repo := &mocks.RepositoryMock{
		ListActiveProductsFn: func(ctx context.Context, businessID uint, f dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error) {
			visto = f
			return nil, 0, nil
		},
	}

	_, _, err := newStorefrontUseCase(repo, nil).ListCatalog(context.Background(), 26,
		dtos.CatalogFilters{Page: 0, PageSize: 500, Search: "camisa", Category: "ropa"})

	require.NoError(t, err)
	assert.Equal(t, 1, visto.Page)
	assert.Equal(t, 12, visto.PageSize)
	assert.Equal(t, "camisa", visto.Search)
	assert.Equal(t, "ropa", visto.Category)
}

func TestListCatalog_FiltraPorNegocio(t *testing.T) {
	var visto uint
	repo := &mocks.RepositoryMock{
		ListActiveProductsFn: func(ctx context.Context, businessID uint, f dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error) {
			visto = businessID
			return []entities.StorefrontProduct{{ID: "P-1"}}, 1, nil
		},
	}

	got, total, err := newStorefrontUseCase(repo, nil).
		ListCatalog(context.Background(), 26, dtos.CatalogFilters{})

	require.NoError(t, err)
	assert.Equal(t, uint(26), visto,
		"el catalogo es por negocio: una tienda no muestra productos de otra")
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
}

func TestListCatalog_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListActiveProductsFn: func(ctx context.Context, businessID uint, f dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error) {
			return nil, 0, dbErr
		},
	}

	_, _, err := newStorefrontUseCase(repo, nil).ListCatalog(context.Background(), 26, dtos.CatalogFilters{})

	assert.ErrorIs(t, err, dbErr)
}

func TestGetProduct_FiltraPorNegocio(t *testing.T) {
	var vistoNegocio uint
	var vistoProducto string
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error) {
			vistoNegocio, vistoProducto = businessID, productID
			return &entities.StorefrontProduct{ID: productID}, nil
		},
	}

	got, err := newStorefrontUseCase(repo, nil).GetProduct(context.Background(), 26, "P-1")

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio,
		"un producto de otra tienda no debe ser accesible adivinando el ID")
	assert.Equal(t, "P-1", vistoProducto)
	assert.Equal(t, "P-1", got.ID)
}

func TestGetProduct_ErrorDelRepo_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error) {
			return nil, domainerrors.ErrProductNotFound
		},
	}

	_, err := newStorefrontUseCase(repo, nil).GetProduct(context.Background(), 26, "P-404")

	assert.ErrorIs(t, err, domainerrors.ErrProductNotFound)
}

func TestListMyOrders_NormalizaPaginacionYFiltraPorUsuario(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 10, 1, 10},
		{-5, 10, 1, 10},
		{1, 0, 1, 10},
		{1, 101, 1, 10},
		{1, 100, 1, 100},
		{2, 50, 2, 50},
	}

	for _, tc := range casos {
		var vistoPage, vistoSize int
		var vistoNegocio, vistoUsuario uint
		repo := &mocks.RepositoryMock{
			ListOrdersByUserIDFn: func(ctx context.Context, businessID, userID uint, page, pageSize int) ([]entities.StorefrontOrder, int64, error) {
				vistoNegocio, vistoUsuario = businessID, userID
				vistoPage, vistoSize = page, pageSize
				return nil, 0, nil
			},
		}

		_, _, err := newStorefrontUseCase(repo, nil).
			ListMyOrders(context.Background(), 26, 42, tc.page, tc.pageSize)

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, vistoPage)
		assert.Equal(t, tc.wantSize, vistoSize)
		assert.Equal(t, uint(26), vistoNegocio)
		assert.Equal(t, uint(42), vistoUsuario,
			"el historial es del usuario autenticado, no de todos")
	}
}

func TestGetMyOrder_FiltraPorNegocioYUsuario(t *testing.T) {
	var vistoOrden string
	var vistoNegocio, vistoUsuario uint
	repo := &mocks.RepositoryMock{
		GetOrderByIDAndUserIDFn: func(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error) {
			vistoOrden, vistoNegocio, vistoUsuario = orderID, businessID, userID
			return &entities.StorefrontOrder{}, nil
		},
	}

	_, err := newStorefrontUseCase(repo, nil).GetMyOrder(context.Background(), "ORD-9", 26, 42)

	require.NoError(t, err)
	assert.Equal(t, "ORD-9", vistoOrden)
	assert.Equal(t, uint(26), vistoNegocio)
	assert.Equal(t, uint(42), vistoUsuario,
		"un cliente no puede ver la orden de otro aunque conozca el ID")
}

func TestGetMyOrder_ErrorDelRepo_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetOrderByIDAndUserIDFn: func(ctx context.Context, orderID string, businessID, userID uint) (*entities.StorefrontOrder, error) {
			return nil, domainerrors.ErrOrderNotFound
		},
	}

	_, err := newStorefrontUseCase(repo, nil).GetMyOrder(context.Background(), "ORD-404", 26, 42)

	assert.ErrorIs(t, err, domainerrors.ErrOrderNotFound)
}

func TestCreateOrder_SinItems_Rechaza(t *testing.T) {
	pub := &mocks.PublisherMock{}

	err := newStorefrontUseCase(nil, pub).CreateOrder(context.Background(), 26, 42,
		&dtos.StorefrontCreateOrderDTO{Items: nil})

	assert.ErrorIs(t, err, domainerrors.ErrNoItems)
	assert.Empty(t, pub.Published)
}

func TestCreateOrder_CantidadInvalida_Rechaza(t *testing.T) {
	for _, cantidad := range []int{0, -1, -100} {
		pub := &mocks.PublisherMock{}
		repo := &mocks.RepositoryMock{}

		err := newStorefrontUseCase(repo, pub).CreateOrder(context.Background(), 26, 42,
			&dtos.StorefrontCreateOrderDTO{
				Items: []dtos.StorefrontOrderItemDTO{{ProductID: "P-1", Quantity: cantidad}},
			})

		assert.ErrorIs(t, err, domainerrors.ErrInvalidQuantity, "cantidad=%d", cantidad)
		assert.Empty(t, pub.Published)
		assert.Empty(t, repo.ProductLookups, "se valida antes de consultar productos")
	}
}

func TestCreateOrder_UnItemInvalidoEntreVarios_RechazaTodo(t *testing.T) {
	pub := &mocks.PublisherMock{}

	err := newStorefrontUseCase(nil, pub).CreateOrder(context.Background(), 26, 42,
		&dtos.StorefrontCreateOrderDTO{
			Items: []dtos.StorefrontOrderItemDTO{
				{ProductID: "P-1", Quantity: 2},
				{ProductID: "P-2", Quantity: 0},
			},
		})

	assert.ErrorIs(t, err, domainerrors.ErrInvalidQuantity)
	assert.Empty(t, pub.Published, "no se publica una orden parcialmente valida")
}

func TestCreateOrder_ElPrecioSaleDelCatalogoNoDelCliente(t *testing.T) {
	pub := &mocks.PublisherMock{}
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error) {
			return &entities.StorefrontProduct{
				ID: productID, SKU: "SKU-1", Name: "Camisa",
				Price: 50000, Currency: "COP", ImageURL: "img.png",
			}, nil
		},
	}

	err := newStorefrontUseCase(repo, pub).CreateOrder(context.Background(), 26, 42,
		&dtos.StorefrontCreateOrderDTO{
			Items: []dtos.StorefrontOrderItemDTO{{ProductID: "P-1", Quantity: 3}},
		})

	require.NoError(t, err)
	require.Len(t, pub.Published, 1)
	orden := decodificarOrden(t, pub.Published[0])
	assert.InDelta(t, 150000.0, orden["total_amount"], 0.001,
		"el total lo calcula el backend con el precio del catalogo, el cliente no lo envia")
	assert.InDelta(t, 150000.0, orden["subtotal"], 0.001)

	items := orden["order_items"].([]interface{})
	require.Len(t, items, 1)
	item := items[0].(map[string]interface{})
	assert.InDelta(t, 50000.0, item["unit_price"], 0.001)
	assert.InDelta(t, 150000.0, item["total_price"], 0.001)
	assert.Equal(t, "SKU-1", item["product_sku"])
	assert.Equal(t, "Camisa", item["product_name"])
}

func TestCreateOrder_SumaVariosItems(t *testing.T) {
	precios := map[string]float64{"P-1": 50000, "P-2": 1000}
	pub := &mocks.PublisherMock{}
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error) {
			return &entities.StorefrontProduct{ID: productID, Price: precios[productID]}, nil
		},
	}

	err := newStorefrontUseCase(repo, pub).CreateOrder(context.Background(), 26, 42,
		&dtos.StorefrontCreateOrderDTO{
			Items: []dtos.StorefrontOrderItemDTO{
				{ProductID: "P-1", Quantity: 2},
				{ProductID: "P-2", Quantity: 5},
			},
		})

	require.NoError(t, err)
	orden := decodificarOrden(t, pub.Published[0])
	assert.InDelta(t, 105000.0, orden["total_amount"], 0.001)
	assert.Len(t, orden["order_items"], 2)
}

func TestCreateOrder_ProductoInexistente_NoPublica(t *testing.T) {
	pub := &mocks.PublisherMock{}
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error) {
			return nil, domainerrors.ErrProductNotFound
		},
	}

	err := newStorefrontUseCase(repo, pub).CreateOrder(context.Background(), 26, 42, ordenValida())

	require.Error(t, err)
	assert.ErrorIs(t, err, domainerrors.ErrProductNotFound)
	assert.Empty(t, pub.Published)
}

func TestCreateOrder_ElProductoSeBuscaEnElNegocioDelUsuario(t *testing.T) {
	var vistoNegocio uint
	pub := &mocks.PublisherMock{}
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, productID string) (*entities.StorefrontProduct, error) {
			vistoNegocio = businessID
			return &entities.StorefrontProduct{ID: productID, Price: 1000}, nil
		},
	}

	err := newStorefrontUseCase(repo, pub).CreateOrder(context.Background(), 26, 42, ordenValida())

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio,
		"no se puede comprar un producto de otra tienda pasando su ID")
}

func TestCreateOrder_UsaLosDatosDelClienteAutenticado(t *testing.T) {
	pub := &mocks.PublisherMock{}
	email := "ana@x.com"
	dni := "1020304050"
	repo := &mocks.RepositoryMock{
		GetClientByUserIDFn: func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
			return &entities.StorefrontClient{
				ID: 1, Name: "Ana Perez", Email: &email, Phone: "3001234567", Dni: &dni,
			}, nil
		},
	}

	err := newStorefrontUseCase(repo, pub).CreateOrder(context.Background(), 26, 42, ordenValida())

	require.NoError(t, err)
	orden := decodificarOrden(t, pub.Published[0])
	assert.Equal(t, "Ana Perez", orden["customer_name"])
	assert.Equal(t, "ana@x.com", orden["customer_email"])
	assert.Equal(t, "3001234567", orden["customer_phone"])
	assert.Equal(t, "1020304050", orden["customer_dni"])
	assert.EqualValues(t, 42, orden["user_id"])
}

func TestCreateOrder_ClienteSinEmailODni_NoRompe(t *testing.T) {
	pub := &mocks.PublisherMock{}
	repo := &mocks.RepositoryMock{
		GetClientByUserIDFn: func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
			return &entities.StorefrontClient{ID: 1, Name: "Ana", Email: nil, Dni: nil}, nil
		},
	}

	err := newStorefrontUseCase(repo, pub).CreateOrder(context.Background(), 26, 42, ordenValida())

	require.NoError(t, err)
	orden := decodificarOrden(t, pub.Published[0])
	assert.Equal(t, "", orden["customer_email"])
	assert.NotContains(t, orden, "customer_dni", "sin dni el campo no se envia")
}

func TestCreateOrder_ClienteInexistente_NoPublica(t *testing.T) {
	pub := &mocks.PublisherMock{}
	repo := &mocks.RepositoryMock{
		GetClientByUserIDFn: func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
			return nil, domainerrors.ErrClientNotFound
		},
	}

	err := newStorefrontUseCase(repo, pub).CreateOrder(context.Background(), 26, 42, ordenValida())

	require.Error(t, err)
	assert.Empty(t, pub.Published)
}

func TestCreateOrder_SinIntegracionDePlataforma_NoPublica(t *testing.T) {
	pub := &mocks.PublisherMock{}
	repo := &mocks.RepositoryMock{
		GetPlatformIntegrationIDFn: func(ctx context.Context, businessID uint) (uint, error) {
			return 0, domainerrors.ErrIntegrationNotFound
		},
	}

	err := newStorefrontUseCase(repo, pub).CreateOrder(context.Background(), 26, 42, ordenValida())

	require.Error(t, err)
	assert.Empty(t, pub.Published)
}

func TestCreateOrder_LaOrdenNaceEnPendienteYNoFacturable(t *testing.T) {
	pub := &mocks.PublisherMock{}

	err := newStorefrontUseCase(nil, pub).CreateOrder(context.Background(), 26, 42, ordenValida())

	require.NoError(t, err)
	orden := decodificarOrden(t, pub.Published[0])
	assert.Equal(t, "pending", orden["status"])
	assert.Equal(t, "pending", orden["original_status"])
	assert.Equal(t, false, orden["invoiceable"],
		"una orden de tienda nace no facturable, se factura despues del pago")
	assert.Equal(t, "COP", orden["currency"])
	assert.EqualValues(t, 0, orden["tax"])
	assert.EqualValues(t, 0, orden["discount"])
	assert.EqualValues(t, 0, orden["shipping_cost"])
}

func TestCreateOrder_MarcaElOrigenComoStorefront(t *testing.T) {
	pub := &mocks.PublisherMock{}
	repo := &mocks.RepositoryMock{
		GetPlatformIntegrationIDFn: func(ctx context.Context, businessID uint) (uint, error) { return 7, nil },
	}

	err := newStorefrontUseCase(repo, pub).CreateOrder(context.Background(), 26, 42, ordenValida())

	require.NoError(t, err)
	orden := decodificarOrden(t, pub.Published[0])
	assert.Equal(t, "platform", orden["integration_type"])
	assert.Equal(t, "storefront", orden["platform"])
	assert.EqualValues(t, 7, orden["integration_id"])
	assert.EqualValues(t, 26, orden["business_id"])
}

func TestCreateOrder_GeneraIdentificadoresUnicos(t *testing.T) {
	pub := &mocks.PublisherMock{}
	uc := newStorefrontUseCase(nil, pub)

	for i := 0; i < 25; i++ {
		require.NoError(t, uc.CreateOrder(context.Background(), 26, 42, ordenValida()))
	}

	externos := map[string]bool{}
	numeros := map[string]bool{}
	for _, raw := range pub.Published {
		orden := decodificarOrden(t, raw)
		ext := orden["external_id"].(string)
		num := orden["order_number"].(string)
		require.False(t, externos[ext], "external_id repetido: %q", ext)
		externos[ext] = true
		numeros[num] = true
		assert.True(t, strings.HasPrefix(ext, "sf-"))
		assert.True(t, strings.HasPrefix(num, "SF-"))
	}
	assert.Len(t, externos, 25)
}

func TestCreateOrder_SinDireccion_NoEnviaDirecciones(t *testing.T) {
	pub := &mocks.PublisherMock{}

	err := newStorefrontUseCase(nil, pub).CreateOrder(context.Background(), 26, 42, ordenValida())

	require.NoError(t, err)
	orden := decodificarOrden(t, pub.Published[0])
	assert.Nil(t, orden["addresses"])
}

func TestCreateOrder_ConDireccion_LaMarcaComoEnvio(t *testing.T) {
	pub := &mocks.PublisherMock{}
	dto := ordenValida()
	dto.Address = &dtos.StorefrontAddressDTO{
		FirstName: "Ana", LastName: "Perez", Phone: "3001234567",
		Street: "Cra 1 #2-3", City: "Bogota", State: "Cundinamarca",
		Country: "CO", PostalCode: "110111", Instructions: strPtr("timbre 2"),
	}

	err := newStorefrontUseCase(nil, pub).CreateOrder(context.Background(), 26, 42, dto)

	require.NoError(t, err)
	orden := decodificarOrden(t, pub.Published[0])
	direcciones := orden["addresses"].([]interface{})
	require.Len(t, direcciones, 1)
	d := direcciones[0].(map[string]interface{})
	assert.Equal(t, "shipping", d["type"])
	assert.Equal(t, "Cra 1 #2-3", d["street"])
	assert.Equal(t, "Bogota", d["city"])
	assert.Equal(t, "timbre 2", d["instructions"])
}

func TestCreateOrder_FallaLaPublicacion_SePropaga(t *testing.T) {
	pubErr := stderrors.New("rabbit caido")
	pub := &mocks.PublisherMock{
		PublishOrderFn: func(ctx context.Context, order []byte) error { return pubErr },
	}

	err := newStorefrontUseCase(nil, pub).CreateOrder(context.Background(), 26, 42, ordenValida())

	assert.ErrorIs(t, err, pubErr,
		"si la cola no acepta la orden el cliente debe enterarse, no quedar creyendo que compro")
}

func TestRegister_NegocioInexistente_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessByCodeFn: func(ctx context.Context, code string) (*entities.StorefrontBusiness, error) {
			return nil, stderrors.New("no existe")
		},
	}

	err := newStorefrontUseCase(repo, nil).Register(context.Background(),
		&dtos.RegisterDTO{BusinessCode: "inexistente", Email: "a@b.com"})

	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)
	assert.Empty(t, repo.CreatedUsers)
}

func TestRegister_EmailNuevo_CreaUsuarioYCliente(t *testing.T) {
	repo := sinClientePrevio(nil)

	err := newStorefrontUseCase(repo, nil).Register(context.Background(), &dtos.RegisterDTO{
		BusinessCode: "demo", Name: "Ana", Email: "ana@x.com",
		Password: "secreta", Phone: "3001234567", Dni: strPtr("1020304050"),
	})

	require.NoError(t, err)
	require.Len(t, repo.CreatedUsers, 1)
	assert.Equal(t, "ana@x.com", repo.CreatedUsers[0].Email)
	require.Len(t, repo.CreatedClients, 1)
	assert.Equal(t, "Ana", repo.CreatedClients[0].Name)
	assert.Equal(t, uint(26), repo.CreatedClients[0].BusinessID)
	require.NotNil(t, repo.CreatedClients[0].UserID)
	assert.Equal(t, uint(42), *repo.CreatedClients[0].UserID)
}

func TestRegister_EmailExistenteConPasswordCorrecta_ReusaElUsuario(t *testing.T) {
	repo := sinClientePrevio(&mocks.RepositoryMock{
		GetUserAuthByEmailFn: func(ctx context.Context, email string) (uint, string, error) {
			return 77, "$2a$hash", nil
		},
		VerifyPasswordFn: func(hash, password string) bool { return true },
	})

	err := newStorefrontUseCase(repo, nil).Register(context.Background(), &dtos.RegisterDTO{
		BusinessCode: "demo", Email: "ana@x.com", Password: "correcta",
	})

	require.NoError(t, err)
	assert.Empty(t, repo.CreatedUsers, "no se crea otro usuario para el mismo email")
	require.Len(t, repo.CreatedClients, 1)
	assert.Equal(t, uint(77), *repo.CreatedClients[0].UserID)
}

func TestRegister_EmailExistenteConPasswordIncorrecta_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetUserAuthByEmailFn: func(ctx context.Context, email string) (uint, string, error) {
			return 77, "$2a$hash", nil
		},
		VerifyPasswordFn: func(hash, password string) bool { return false },
	}

	err := newStorefrontUseCase(repo, nil).Register(context.Background(), &dtos.RegisterDTO{
		BusinessCode: "demo", Email: "ana@x.com", Password: "adivinada",
	})

	assert.ErrorIs(t, err, domainerrors.ErrPasswordMismatch,
		"registrarse con un email ajeno no debe dar acceso a esa cuenta")
	assert.Empty(t, repo.CreatedClients)
	assert.Empty(t, repo.StaffCalls)
}

func TestRegister_SinStaff_LoCreaConElRolClienteFinal(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClienteFinalRoleIDFn: func(ctx context.Context) (uint, error) { return 9, nil },
		StaffExistsFn:           func(ctx context.Context, userID, businessID uint) (bool, error) { return false, nil },
	}

	err := newStorefrontUseCase(repo, nil).Register(context.Background(),
		&dtos.RegisterDTO{BusinessCode: "demo", Email: "ana@x.com"})

	require.NoError(t, err)
	require.Len(t, repo.StaffCalls, 1)
	assert.Equal(t, uint(42), repo.StaffCalls[0].UserID)
	assert.Equal(t, uint(26), repo.StaffCalls[0].BusinessID)
	assert.Equal(t, uint(9), repo.StaffCalls[0].RoleID,
		"un registro publico solo puede quedar como cliente_final, nunca con otro rol")
}

func TestRegister_ConStaffExistente_NoLoDuplica(t *testing.T) {
	repo := &mocks.RepositoryMock{
		StaffExistsFn: func(ctx context.Context, userID, businessID uint) (bool, error) { return true, nil },
	}

	err := newStorefrontUseCase(repo, nil).Register(context.Background(),
		&dtos.RegisterDTO{BusinessCode: "demo", Email: "ana@x.com"})

	require.NoError(t, err)
	assert.Empty(t, repo.StaffCalls)
}

func TestRegister_SinRolClienteFinal_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClienteFinalRoleIDFn: func(ctx context.Context) (uint, error) {
			return 0, stderrors.New("no existe")
		},
	}

	err := newStorefrontUseCase(repo, nil).Register(context.Background(),
		&dtos.RegisterDTO{BusinessCode: "demo", Email: "ana@x.com"})

	assert.ErrorIs(t, err, domainerrors.ErrRoleNotFound)
	assert.Empty(t, repo.StaffCalls)
}

func TestRegister_YaEsClienteDelNegocio_EsIdempotente(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClientByUserIDFn: func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
			return &entities.StorefrontClient{ID: 5, BusinessID: businessID}, nil
		},
	}

	err := newStorefrontUseCase(repo, nil).Register(context.Background(),
		&dtos.RegisterDTO{BusinessCode: "demo", Email: "ana@x.com"})

	require.NoError(t, err, "registrarse dos veces no debe fallar ni duplicar el cliente")
	assert.Empty(t, repo.CreatedClients)
}

func TestRegister_ClientePorEmailSinUsuario_LoVincula(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClientByUserIDFn: func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
			return nil, stderrors.New("no encontrado")
		},
		GetClientByBusinessAndEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.StorefrontClient, error) {
			return &entities.StorefrontClient{ID: 5, BusinessID: businessID, UserID: nil}, nil
		},
	}

	err := newStorefrontUseCase(repo, nil).Register(context.Background(),
		&dtos.RegisterDTO{BusinessCode: "demo", Email: "ana@x.com"})

	require.NoError(t, err)
	require.Len(t, repo.LinkCalls, 1)
	assert.Equal(t, uint(5), repo.LinkCalls[0].ClientID)
	assert.Equal(t, uint(42), repo.LinkCalls[0].UserID,
		"un cliente cargado por el negocio se vincula a la cuenta, no se duplica")
	assert.Empty(t, repo.CreatedClients)
}

func TestRegister_ClientePorEmailDeOtroUsuario_Rechaza(t *testing.T) {
	otro := uint(99)
	repo := &mocks.RepositoryMock{
		GetClientByUserIDFn: func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
			return nil, stderrors.New("no encontrado")
		},
		GetClientByBusinessAndEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.StorefrontClient, error) {
			return &entities.StorefrontClient{ID: 5, UserID: &otro}, nil
		},
	}

	err := newStorefrontUseCase(repo, nil).Register(context.Background(),
		&dtos.RegisterDTO{BusinessCode: "demo", Email: "ana@x.com"})

	assert.ErrorIs(t, err, domainerrors.ErrEmailAlreadyExists,
		"ese email ya pertenece a otra cuenta en este negocio")
	assert.Empty(t, repo.LinkCalls)
	assert.Empty(t, repo.CreatedClients)
}

func TestRegister_ClientePorEmailYaVinculadoAlMismoUsuario_NoRevincula(t *testing.T) {
	mismo := uint(42)
	repo := &mocks.RepositoryMock{
		GetClientByUserIDFn: func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
			return nil, stderrors.New("no encontrado")
		},
		GetClientByBusinessAndEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.StorefrontClient, error) {
			return &entities.StorefrontClient{ID: 5, UserID: &mismo}, nil
		},
	}

	err := newStorefrontUseCase(repo, nil).Register(context.Background(),
		&dtos.RegisterDTO{BusinessCode: "demo", Email: "ana@x.com"})

	require.NoError(t, err)
	assert.Empty(t, repo.LinkCalls)
	assert.Empty(t, repo.CreatedClients)
}

func TestRegister_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al consultar el email", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetUserAuthByEmailFn: func(ctx context.Context, email string) (uint, string, error) {
				return 0, "", dbErr
			},
		}
		err := newStorefrontUseCase(repo, nil).Register(context.Background(),
			&dtos.RegisterDTO{BusinessCode: "demo", Email: "a@b.com"})
		assert.ErrorIs(t, err, dbErr)
		assert.Empty(t, repo.CreatedUsers)
	})

	t.Run("al crear el usuario", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			CreateUserFn: func(ctx context.Context, u *entities.NewUser) (uint, error) { return 0, dbErr },
		}
		err := newStorefrontUseCase(repo, nil).Register(context.Background(),
			&dtos.RegisterDTO{BusinessCode: "demo", Email: "a@b.com"})
		assert.ErrorIs(t, err, dbErr)
		assert.Empty(t, repo.CreatedClients)
	})

	t.Run("al verificar staff", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			StaffExistsFn: func(ctx context.Context, userID, businessID uint) (bool, error) {
				return false, dbErr
			},
		}
		err := newStorefrontUseCase(repo, nil).Register(context.Background(),
			&dtos.RegisterDTO{BusinessCode: "demo", Email: "a@b.com"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al crear staff", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			CreateBusinessStaffFn: func(ctx context.Context, userID, businessID, roleID uint) error {
				return dbErr
			},
		}
		err := newStorefrontUseCase(repo, nil).Register(context.Background(),
			&dtos.RegisterDTO{BusinessCode: "demo", Email: "a@b.com"})
		assert.ErrorIs(t, err, dbErr)
		assert.Empty(t, repo.CreatedClients)
	})

	t.Run("al buscar cliente por email", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetClientByUserIDFn: func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
				return nil, stderrors.New("no encontrado")
			},
			GetClientByBusinessAndEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.StorefrontClient, error) {
				return nil, dbErr
			},
		}
		err := newStorefrontUseCase(repo, nil).Register(context.Background(),
			&dtos.RegisterDTO{BusinessCode: "demo", Email: "a@b.com"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al vincular cliente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetClientByUserIDFn: func(ctx context.Context, businessID, userID uint) (*entities.StorefrontClient, error) {
				return nil, stderrors.New("no encontrado")
			},
			GetClientByBusinessAndEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.StorefrontClient, error) {
				return &entities.StorefrontClient{ID: 5, UserID: nil}, nil
			},
			LinkClientUserFn: func(ctx context.Context, clientID, userID uint) error { return dbErr },
		}
		err := newStorefrontUseCase(repo, nil).Register(context.Background(),
			&dtos.RegisterDTO{BusinessCode: "demo", Email: "a@b.com"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al crear el cliente", func(t *testing.T) {
		repo := sinClientePrevio(&mocks.RepositoryMock{
			CreateClientFn: func(ctx context.Context, c *entities.StorefrontClient) error { return dbErr },
		})
		err := newStorefrontUseCase(repo, nil).Register(context.Background(),
			&dtos.RegisterDTO{BusinessCode: "demo", Email: "a@b.com"})
		assert.ErrorIs(t, err, dbErr)
	})
}
