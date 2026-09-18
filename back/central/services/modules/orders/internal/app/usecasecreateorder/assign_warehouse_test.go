package usecasecreateorder

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func runCreateWithWarehouse(t *testing.T, dto *dtos.ProbabilityOrderDTO, resolve func(ctx context.Context, businessID, integrationID uint) (*entities.WarehouseRef, error)) (*entities.ProbabilityOrder, int) {
	t.Helper()
	var saved *entities.ProbabilityOrder
	calls := 0
	repo := &mockRepository{
		OrderExistsFn: func(ctx context.Context, externalID string, integrationID uint) (bool, error) {
			return false, nil
		},
		GetClientByEmailFn: func(ctx context.Context, businessID uint, email string) (*entities.Client, error) {
			return &entities.Client{ID: 99, Email: &email}, nil
		},
		CreateOrderFn: func(ctx context.Context, order *entities.ProbabilityOrder) error {
			order.ID = "ORDER-WH"
			saved = order
			return nil
		},
		ResolveOrderWarehouseFn: func(ctx context.Context, businessID, integrationID uint) (*entities.WarehouseRef, error) {
			calls++
			return resolve(ctx, businessID, integrationID)
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)
	if _, err := uc.MapAndSaveOrder(context.Background(), dto); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if saved == nil {
		t.Fatal("la orden no se guardo")
	}
	return saved, calls
}

func TestMapAndSaveOrder_SinBodega_AsignaLaResuelta(t *testing.T) {
	dto := newMinimalDTO("EXT-WH-1", 228, 46)
	order, calls := runCreateWithWarehouse(t, dto, func(ctx context.Context, businessID, integrationID uint) (*entities.WarehouseRef, error) {
		if businessID != 46 || integrationID != 228 {
			t.Errorf("parametros inesperados: business %d integracion %d", businessID, integrationID)
		}
		return &entities.WarehouseRef{ID: 22, Name: "Bodega Principal"}, nil
	})
	if calls != 1 {
		t.Errorf("se esperaba 1 llamada a ResolveOrderWarehouse, hubo %d", calls)
	}
	if order.WarehouseID == nil || *order.WarehouseID != 22 {
		t.Fatalf("bodega esperada 22, obtenida %v", order.WarehouseID)
	}
	if order.WarehouseName != "Bodega Principal" {
		t.Errorf("nombre esperado %q, obtenido %q", "Bodega Principal", order.WarehouseName)
	}
}

func TestMapAndSaveOrder_ConBodega_NoLaCambia(t *testing.T) {
	dto := newMinimalDTO("EXT-WH-2", 228, 46)
	wh := uint(7)
	dto.Shipments = []dtos.ProbabilityShipmentDTO{{WarehouseID: &wh, WarehouseName: "Bodega Sur"}}
	order, calls := runCreateWithWarehouse(t, dto, func(ctx context.Context, businessID, integrationID uint) (*entities.WarehouseRef, error) {
		return &entities.WarehouseRef{ID: 22, Name: "Bodega Principal"}, nil
	})
	if calls != 0 {
		t.Errorf("no se esperaba resolver bodega, hubo %d llamadas", calls)
	}
	if order.WarehouseID == nil || *order.WarehouseID != 7 {
		t.Fatalf("bodega esperada 7, obtenida %v", order.WarehouseID)
	}
}

func TestMapAndSaveOrder_NegocioSinBodega_QuedaSinBodega(t *testing.T) {
	dto := newMinimalDTO("EXT-WH-3", 228, 46)
	order, _ := runCreateWithWarehouse(t, dto, func(ctx context.Context, businessID, integrationID uint) (*entities.WarehouseRef, error) {
		return nil, nil
	})
	if order.WarehouseID != nil {
		t.Errorf("se esperaba bodega nil, obtenida %d", *order.WarehouseID)
	}
}
