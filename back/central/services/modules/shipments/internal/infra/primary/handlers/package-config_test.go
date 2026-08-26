package handlers

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
	"github.com/secamc93/probability/back/central/shared/shippingpkg"
)

func f64(v float64) *float64 { return &v }

func demoBoxConfig() *domain.PackageConfig {
	return &domain.PackageConfig{
		Strategy: shippingpkg.StrategyStandardBox,
		Boxes: []shippingpkg.Box{{
			Name: "Caja Mediana", Weight: f64(3), Length: f64(30), Width: f64(40), Height: f64(30), MaxItems: 8,
		}},
	}
}

func newPackageTestHandlers(repo *mocks.RepositoryMock) *Handlers {
	return &Handlers{uc: usecases.New(repo, nil, nil)}
}

func TestResolveWooPackageDimensions_UsaCajaEstandarDelNegocio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessPackageConfigFn: func(_ context.Context, businessID uint, warehouseID *uint) (*domain.PackageConfig, error) {
			if businessID != 26 {
				t.Fatalf("business inesperado: %d", businessID)
			}
			if warehouseID == nil || *warehouseID != 12 {
				t.Fatalf("debe buscar la config de la bodega de origen (12), llego %v", warehouseID)
			}
			return demoBoxConfig(), nil
		},
		GetProductDimensionsBySKUsFn: func(_ context.Context, _ uint, _ []string) (map[string]domain.ProductDimensions, error) {
			return map[string]domain.ProductDimensions{"SKU-A": {Weight: f64(0.4), Length: f64(10), Width: f64(10), Height: f64(5)}}, nil
		},
	}
	h := newPackageTestHandlers(repo)
	resolved := &wooResolved{OriginIsWarehouse: true, Origin: &domain.OriginAddress{ID: 12}}
	req := wooRateRequest{Contents: []wooRateItem{{Sku: "SKU-A", Quantity: 2, WeightGrams: 500}}}

	pkg := h.resolveWooPackageDimensions(context.Background(), 26, resolved, req)

	if pkg.Length != 30 || pkg.Width != 40 || pkg.Height != 30 {
		t.Fatalf("dimensiones deben ser las de la caja 30x40x30, llego %+v", pkg)
	}
	if pkg.Weight != 3 {
		t.Fatalf("peso debe ser el de la caja (3kg > 1kg del carrito), llego %v", pkg.Weight)
	}
}

func TestResolveWooPackageDimensions_SinConfigUsaProductoOCarrito(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessPackageConfigFn: func(_ context.Context, _ uint, _ *uint) (*domain.PackageConfig, error) {
			return nil, nil
		},
		GetProductDimensionsBySKUsFn: func(_ context.Context, _ uint, _ []string) (map[string]domain.ProductDimensions, error) {
			return map[string]domain.ProductDimensions{"SKU-A": {Length: f64(20), Width: f64(15), Height: f64(12)}}, nil
		},
	}
	h := newPackageTestHandlers(repo)
	req := wooRateRequest{Contents: []wooRateItem{{Sku: "SKU-A", Quantity: 3, WeightGrams: 700}}}

	pkg := h.resolveWooPackageDimensions(context.Background(), 26, &wooResolved{}, req)

	if pkg.Length != 20 || pkg.Width != 15 || pkg.Height != 12 {
		t.Fatalf("sin caja deben ir las dimensiones del producto, llego %+v", pkg)
	}
	if pkg.Weight != 2.1 {
		t.Fatalf("peso debe ser el del carrito (3 x 0.7kg), llego %v", pkg.Weight)
	}
}

func TestApplyPackageConfig_ConOrdenAplicaCaja(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetOrderPackageItemsFn: func(_ context.Context, orderID string) ([]shippingpkg.PackageItem, uint, error) {
			return []shippingpkg.PackageItem{{SKU: "SKU-A", Quantity: 2, Weight: f64(0.4)}}, 12, nil
		},
		GetBusinessPackageConfigFn: func(_ context.Context, _ uint, warehouseID *uint) (*domain.PackageConfig, error) {
			if warehouseID == nil || *warehouseID != 12 {
				t.Fatalf("debe buscar la config de la bodega de la orden (12), llego %v", warehouseID)
			}
			return demoBoxConfig(), nil
		},
	}
	h := newPackageTestHandlers(repo)
	raw := map[string]interface{}{"packages": []interface{}{map[string]interface{}{"weight": 1.0, "height": 10.0, "width": 10.0, "length": 10.0}}}

	h.applyPackageConfig(context.Background(), 26, "order-1", raw)

	pkg := raw["packages"].([]interface{})[0].(map[string]interface{})
	if pkg["length"] != 30.0 || pkg["width"] != 40.0 || pkg["height"] != 30.0 || pkg["weight"] != 3.0 {
		t.Fatalf("paquete debe ser la caja 30x40x30 de 3kg, llego %v", pkg)
	}
}

func TestApplyPackageConfig_SinOrdenAplicaCajaDelNegocio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetOrderPackageItemsFn: func(_ context.Context, _ string) ([]shippingpkg.PackageItem, uint, error) {
			t.Fatal("sin orden no debe consultar items")
			return nil, 0, nil
		},
		GetBusinessPackageConfigFn: func(_ context.Context, _ uint, warehouseID *uint) (*domain.PackageConfig, error) {
			if warehouseID != nil {
				t.Fatalf("sin orden debe usar la config del negocio, llego bodega %v", *warehouseID)
			}
			return demoBoxConfig(), nil
		},
	}
	h := newPackageTestHandlers(repo)
	raw := map[string]interface{}{"packages": []interface{}{map[string]interface{}{"weight": 1.0, "height": 10.0, "width": 10.0, "length": 10.0}}}

	h.applyPackageConfig(context.Background(), 26, "", raw)

	pkg := raw["packages"].([]interface{})[0].(map[string]interface{})
	if pkg["length"] != 30.0 || pkg["width"] != 40.0 || pkg["height"] != 30.0 || pkg["weight"] != 3.0 {
		t.Fatalf("cotizador expres debe usar la caja del negocio, llego %v", pkg)
	}
}

func TestApplyPackageConfig_RespetaPaqueteEditadoPorElUsuario(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessPackageConfigFn: func(_ context.Context, _ uint, _ *uint) (*domain.PackageConfig, error) {
			t.Fatal("si el usuario ya puso medidas no se debe consultar la config")
			return nil, nil
		},
	}
	h := newPackageTestHandlers(repo)
	raw := map[string]interface{}{"packages": []interface{}{map[string]interface{}{"weight": 5.0, "height": 30.0, "width": 30.0, "length": 30.0}}}

	h.applyPackageConfig(context.Background(), 26, "", raw)

	pkg := raw["packages"].([]interface{})[0].(map[string]interface{})
	if pkg["length"] != 30.0 || pkg["weight"] != 5.0 {
		t.Fatalf("paquete editado no debe cambiar, llego %v", pkg)
	}
}
