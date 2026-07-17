package usecases

import (
	"encoding/json"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

const configEmparejada = `{
  "inventory_sync_enabled": true,
  "inventory_warehouse_mode": "mapped",
  "jumpseller_default_location_id": 327870,
  "jumpseller_location_mappings": [
    { "internal_warehouse_id": 7, "jumpseller_location_id": 327870 },
    { "internal_warehouse_id": 9, "jumpseller_location_id": 327870 },
    { "internal_warehouse_id": 11, "jumpseller_location_id": 400001 }
  ]
}`

const configUnaBodega = `{
  "inventory_sync_enabled": true,
  "inventory_warehouse_mode": "single",
  "inventory_single_warehouse_id": 7
}`

func parse(t *testing.T, raw string) domain.InventoryConfig {
	t.Helper()
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		t.Fatal(err)
	}
	return parseInventoryConfig(config)
}

func TestModoUnaBodega(t *testing.T) {
	cfg := parse(t, configUnaBodega)

	if cfg.Mode != domain.InventoryModeSingle {
		t.Fatalf("mode = %q", cfg.Mode)
	}
	ids := resolveWarehouseIDs(cfg)
	if len(ids) != 1 || ids[0] != 7 {
		t.Fatalf("una bodega debe consultar solo la 7, dio %v", ids)
	}
}

func TestModoEmparejadoAgrupaPorLocation(t *testing.T) {
	cfg := parse(t, configEmparejada)

	if cfg.Mode != domain.InventoryModeMapped {
		t.Fatalf("mode = %q", cfg.Mode)
	}
	if len(cfg.LocationMappings) != 3 {
		t.Fatalf("mappings = %d", len(cfg.LocationMappings))
	}

	groups := cfg.LocationGroups()
	if len(groups) != 2 {
		t.Fatalf("locations = %d, se esperaban 2 (327870 y 400001)", len(groups))
	}
	if len(groups[327870]) != 2 {
		t.Fatalf("la location 327870 debe recibir la suma de 2 bodegas, tiene %v", groups[327870])
	}
	if len(groups[400001]) != 1 || groups[400001][0] != 11 {
		t.Fatalf("la location 400001 debe recibir solo la bodega 11, tiene %v", groups[400001])
	}
}

func TestEmparejarDosBodegasAlMismoDestinoEsSumar(t *testing.T) {
	cfg := parse(t, `{
	  "inventory_warehouse_mode": "mapped",
	  "jumpseller_location_mappings": [
	    { "internal_warehouse_id": 7, "jumpseller_location_id": 327870 },
	    { "internal_warehouse_id": 9, "jumpseller_location_id": 327870 }
	  ]
	}`)

	groups := cfg.LocationGroups()
	if len(groups) != 1 {
		t.Fatalf("dos bodegas al mismo destino son UN grupo, dio %d", len(groups))
	}
	ids := resolveWarehouseIDs(cfg)
	if len(ids) != 2 {
		t.Fatalf("debe consultar el stock de ambas bodegas para sumarlas, dio %v", ids)
	}
}

func TestParseInventoryConfigIgnoraMapeosIncompletos(t *testing.T) {
	cfg := parse(t, `{
	  "jumpseller_location_mappings": [
	    { "internal_warehouse_id": 0, "jumpseller_location_id": 327870 },
	    { "internal_warehouse_id": 7, "jumpseller_location_id": 0 },
	    { "internal_warehouse_id": 7, "jumpseller_location_id": 327870 }
	  ]
	}`)

	if len(cfg.LocationMappings) != 1 {
		t.Fatalf("mappings = %d, solo el completo debe entrar", len(cfg.LocationMappings))
	}
}

func TestParseInventoryConfigVaciaNoRevienta(t *testing.T) {
	cfg := parseInventoryConfig(map[string]interface{}{})

	if cfg.Enabled {
		t.Fatal("sin config, el sync no debe quedar activo")
	}
	if cfg.Mode != domain.InventoryModeSingle {
		t.Fatalf("el modo por defecto debe ser una bodega, dio %q", cfg.Mode)
	}
	if len(cfg.LocationMappings) != 0 || cfg.DefaultLocationID != 0 {
		t.Fatal("sin config no se debe inventar mapeo")
	}
}
