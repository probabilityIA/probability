package usecases

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/inventorycompare"
)

func filaPorSKU(filas []inventorycompare.Row, sku string) *inventorycompare.Row {
	for i := range filas {
		if filas[i].SKU == sku {
			return &filas[i]
		}
	}
	return nil
}

func TestElComparativoMarcaParaActualizarCuandoElCanalTieneOtroStock(t *testing.T) {
	repo := &fakeRepo{
		mapped: []domain.MappedItem{
			{ProductID: "p1", SKU: "SKU-1", Name: "Camiseta", ExternalItemID: "2001:9001", ExternalVariantID: "9001"},
		},
		stock: map[string]int{"p1": 10},
	}
	client := &fakeClient{
		channelStock: []domain.ChannelStock{
			{ExternalID: "2001:9001", Quantity: 3, ManageStock: true, Found: true},
		},
	}

	page, err := newTestUseCase(client, repo).CompareInventory(context.Background(), testIntegrationID, testBusinessID, 1, 100)
	if err != nil {
		t.Fatalf("CompareInventory fallo: %v", err)
	}

	fila := filaPorSKU(page.Rows, "SKU-1")
	if fila == nil {
		t.Fatal("el SKU mapeado no aparece en el comparativo")
	}
	if fila.Action != inventorycompare.ActionUpdate {
		t.Fatalf("Probability tiene 10 y el canal 3: debe quedar para actualizar, quedo %q", fila.Action)
	}
	if fila.ProbabilityQty == nil || *fila.ProbabilityQty != 10 {
		t.Fatalf("la cantidad de Probability no se leyo: %+v", fila.ProbabilityQty)
	}
	if fila.ChannelQty == nil || *fila.ChannelQty != 3 {
		t.Fatalf("la cantidad del canal no se leyo: %+v", fila.ChannelQty)
	}
	if page.Totals.ToUpdate != 1 {
		t.Fatalf("el resumen debe contar 1 por actualizar, conto %d", page.Totals.ToUpdate)
	}
}

func TestElComparativoNoTocaLoQueYaCoincide(t *testing.T) {
	repo := &fakeRepo{
		mapped: []domain.MappedItem{
			{ProductID: "p1", SKU: "SKU-1", ExternalItemID: "2001:9001"},
		},
		stock: map[string]int{"p1": 7},
	}
	client := &fakeClient{
		channelStock: []domain.ChannelStock{
			{ExternalID: "2001:9001", Quantity: 7, ManageStock: true, Found: true},
		},
	}

	page, err := newTestUseCase(client, repo).CompareInventory(context.Background(), testIntegrationID, testBusinessID, 1, 100)
	if err != nil {
		t.Fatalf("CompareInventory fallo: %v", err)
	}

	fila := filaPorSKU(page.Rows, "SKU-1")
	if fila == nil || fila.Action != inventorycompare.ActionUnchanged {
		t.Fatalf("con el mismo stock en ambos lados no debe proponerse ningun cambio: %+v", fila)
	}
}

func TestElComparativoSeSaltaLaPublicacionQueNoManejaStock(t *testing.T) {
	repo := &fakeRepo{
		mapped: []domain.MappedItem{
			{ProductID: "p1", SKU: "SKU-1", ExternalItemID: "2001:9001"},
		},
		stock: map[string]int{"p1": 5},
	}
	client := &fakeClient{
		channelStock: []domain.ChannelStock{
			{ExternalID: "2001:9001", Quantity: 0, ManageStock: false, Found: true},
		},
	}

	page, err := newTestUseCase(client, repo).CompareInventory(context.Background(), testIntegrationID, testBusinessID, 1, 100)
	if err != nil {
		t.Fatalf("CompareInventory fallo: %v", err)
	}

	fila := filaPorSKU(page.Rows, "SKU-1")
	if fila == nil || fila.Action != inventorycompare.ActionSkip {
		t.Fatalf("una variante con stock ilimitado en Tiendanube no se debe pisar: %+v", fila)
	}
}

func TestElComparativoGuardaLaFotoParaConsultarlaSinLlamarAlCanal(t *testing.T) {
	repo := &fakeRepo{
		mapped: []domain.MappedItem{
			{ProductID: "p1", SKU: "SKU-1", ExternalItemID: "2001:9001"},
		},
		stock: map[string]int{"p1": 4},
	}
	client := &fakeClient{
		channelStock: []domain.ChannelStock{
			{ExternalID: "2001:9001", Quantity: 1, ManageStock: true, Found: true},
		},
	}

	uc := newTestUseCase(client, repo)
	if _, err := uc.CompareInventory(context.Background(), testIntegrationID, testBusinessID, 1, 100); err != nil {
		t.Fatalf("CompareInventory fallo: %v", err)
	}
	if len(repo.savedRows) == 0 {
		t.Fatal("el comparativo no guardo la foto: la pestana de solo lectura quedaria vacia")
	}

	guardado, err := uc.LoadInventoryCompare(context.Background(), testIntegrationID, testBusinessID, inventorycompare.LoadOptions{Page: 1, PageSize: 100})
	if err != nil {
		t.Fatalf("LoadInventoryCompare fallo: %v", err)
	}
	if !guardado.FromCache || len(guardado.Rows) == 0 {
		t.Fatalf("la foto guardada no se pudo releer: %+v", guardado)
	}
}

func TestElComparativoFiltraPorLosSKUsPedidos(t *testing.T) {
	repo := &fakeRepo{
		mapped: []domain.MappedItem{
			{ProductID: "p1", SKU: "SKU-1", ExternalItemID: "2001:9001"},
			{ProductID: "p2", SKU: "SKU-2", ExternalItemID: "2002:9002"},
		},
		stock: map[string]int{"p1": 1, "p2": 2},
	}
	client := &fakeClient{
		channelStock: []domain.ChannelStock{
			{ExternalID: "2001:9001", Quantity: 9, ManageStock: true, Found: true},
			{ExternalID: "2002:9002", Quantity: 9, ManageStock: true, Found: true},
		},
	}

	page, err := newTestUseCase(client, repo).CompareInventory(context.Background(), testIntegrationID, testBusinessID, 1, 100, "sku-2")
	if err != nil {
		t.Fatalf("CompareInventory fallo: %v", err)
	}
	if len(page.Rows) != 1 || page.Rows[0].SKU != "SKU-2" {
		t.Fatalf("el filtro por SKU debe ignorar mayusculas y traer solo el pedido: %+v", page.Rows)
	}
}

func TestSinProductosMapeadosElComparativoNoLlamaAlCanal(t *testing.T) {
	repo := &fakeRepo{}
	client := &fakeClient{stockErr: errNoDebioLlamarse}

	page, err := newTestUseCase(client, repo).CompareInventory(context.Background(), testIntegrationID, testBusinessID, 1, 100)
	if err != nil {
		t.Fatalf("sin productos mapeados no debe fallar: %v", err)
	}
	if len(page.Rows) != 0 || page.Total != 0 {
		t.Fatalf("sin productos mapeados el comparativo debe venir vacio: %+v", page)
	}
}
