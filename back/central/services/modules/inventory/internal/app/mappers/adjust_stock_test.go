package mappers

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func u(v uint) *uint {
	return &v
}

func TestAdjustStockDTOToTxParams_CopiaTodosLosCampos(t *testing.T) {
	dto := request.AdjustStockDTO{
		ProductID:   "prod-001",
		WarehouseID: 5,
		LocationID:  u(11),
		LotID:       u(22),
		StateID:     u(33),
		UomID:       u(44),
		BusinessID:  10,
		Quantity:    -15,
		Reason:      "ajuste por conteo",
		Notes:       "diferencia detectada en auditoria",
		CreatedByID: u(77),
	}

	got := AdjustStockDTOToTxParams(dto, 2, "manual_adjustment")

	assert.Equal(t, "prod-001", got.ProductID)
	assert.Equal(t, uint(5), got.WarehouseID)
	require.NotNil(t, got.LocationID)
	assert.Equal(t, uint(11), *got.LocationID)
	require.NotNil(t, got.LotID)
	assert.Equal(t, uint(22), *got.LotID)
	require.NotNil(t, got.StateID)
	assert.Equal(t, uint(33), *got.StateID)
	require.NotNil(t, got.UomID)
	assert.Equal(t, uint(44), *got.UomID)
	assert.Equal(t, uint(10), got.BusinessID)
	assert.Equal(t, -15, got.Quantity)
	assert.Equal(t, uint(2), got.MovementTypeID)
	assert.Equal(t, "ajuste por conteo", got.Reason)
	assert.Equal(t, "diferencia detectada en auditoria", got.Notes)
	assert.Equal(t, "manual_adjustment", got.ReferenceType)
	require.NotNil(t, got.CreatedByID)
	assert.Equal(t, uint(77), *got.CreatedByID)
}

func TestAdjustStockDTOToTxParams_CamposOpcionalesNulos(t *testing.T) {
	got := AdjustStockDTOToTxParams(request.AdjustStockDTO{
		ProductID:  "prod-001",
		BusinessID: 10,
		Quantity:   5,
	}, 1, "provider_sync")

	assert.Nil(t, got.LocationID)
	assert.Nil(t, got.LotID)
	assert.Nil(t, got.StateID)
	assert.Nil(t, got.UomID)
	assert.Nil(t, got.CreatedByID)
	assert.Equal(t, "provider_sync", got.ReferenceType)
	assert.Equal(t, uint(1), got.MovementTypeID)
}

func TestTransferStockDTOToTxParams_CopiaTodosLosCampos(t *testing.T) {
	dto := request.TransferStockDTO{
		ProductID:       "prod-001",
		FromWarehouseID: 5,
		ToWarehouseID:   9,
		FromLocationID:  u(11),
		ToLocationID:    u(12),
		LotID:           u(22),
		StateID:         u(33),
		UomID:           u(44),
		BusinessID:      10,
		Quantity:        7,
		Reason:          "traslado entre bodegas",
		Notes:           "reabastecimiento",
		CreatedByID:     u(77),
	}

	got := TransferStockDTOToTxParams(dto, 3, "transfer")

	assert.Equal(t, "prod-001", got.ProductID)
	assert.Equal(t, uint(5), got.FromWarehouseID)
	assert.Equal(t, uint(9), got.ToWarehouseID)
	require.NotNil(t, got.FromLocationID)
	assert.Equal(t, uint(11), *got.FromLocationID)
	require.NotNil(t, got.ToLocationID)
	assert.Equal(t, uint(12), *got.ToLocationID)
	assert.Equal(t, uint(10), got.BusinessID)
	assert.Equal(t, 7, got.Quantity)
	assert.Equal(t, uint(3), got.MovementTypeID)
	assert.Equal(t, "traslado entre bodegas", got.Reason)
	assert.Equal(t, "reabastecimiento", got.Notes)
	assert.Equal(t, "transfer", got.ReferenceType)
	require.NotNil(t, got.CreatedByID)
	assert.Equal(t, uint(77), *got.CreatedByID)
}

func TestTransferStockDTOToTxParams_NoMezclaOrigenYDestino(t *testing.T) {
	got := TransferStockDTOToTxParams(request.TransferStockDTO{
		FromWarehouseID: 1,
		ToWarehouseID:   2,
		FromLocationID:  u(100),
		ToLocationID:    u(200),
	}, 3, "transfer")

	assert.Equal(t, uint(1), got.FromWarehouseID)
	assert.Equal(t, uint(2), got.ToWarehouseID)
	assert.Equal(t, uint(100), *got.FromLocationID)
	assert.Equal(t, uint(200), *got.ToLocationID)
}
