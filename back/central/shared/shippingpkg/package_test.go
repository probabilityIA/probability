package shippingpkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func f(v float64) *float64 { return &v }

func vigaBoxes() []Box {
	return []Box{
		{Name: "Caja Mediana", Weight: f(3), Length: f(30), Width: f(40), Height: f(30), MaxItems: 8},
		{Name: "Caja muy chica", Weight: f(1), Length: f(10), Width: f(10), Height: f(10), MaxItems: 3},
		{Name: "Caja Chica", Weight: f(1), Length: f(20), Width: f(20), Height: f(20), MaxItems: 5},
	}
}

func TestResolvePackage_OrdenViga14636(t *testing.T) {
	boxes := vigaBoxes()

	out := Resolve(StrategyStandardBox, boxes, PackageInput{
		TotalQuantity: 4,
		CartWeightKg:  3,
		Items: []PackageItem{
			{SKU: "BT16-5XL", Quantity: 1},
			{SKU: "BN16-5XL", Quantity: 1},
			{SKU: "CH113-5XL", Quantity: 1},
			{SKU: "CHB36-5XL", Quantity: 1},
		},
	})

	assert.Equal(t, PackageSourceStandardBox, out.Source)
	assert.Equal(t, "Caja Chica", out.BoxName)
	assert.Equal(t, 20.0, out.Length)
	assert.Equal(t, 3.0, out.Weight,
		"el peso real del carrito manda sobre el peso nominal de la caja")
}

func TestResolvePackage_SinCajaNiDimensiones(t *testing.T) {
	out := Resolve(StrategyProductDimensions, nil, PackageInput{TotalQuantity: 2})

	assert.Equal(t, PackageSourceDefault, out.Source)
	assert.Equal(t, DefaultDimCm, out.Length)
	assert.Equal(t, DefaultWeightKg, out.Weight)
}

func TestResolvePackage_PesoDelCatalogoCuandoElCarritoNoLoTrae(t *testing.T) {
	out := Resolve(StrategyProductDimensions, nil, PackageInput{
		TotalQuantity:  2,
		CartWeightKg:   0,
		CatalogWeights: map[string]float64{"SKU-A": 0.75},
		Items: []PackageItem{
			{SKU: "SKU-A", Quantity: 2},
		},
	})

	assert.Equal(t, 1.5, out.Weight, "0.75 kg por 2 unidades sale del catalogo")
}

func TestResolvePackage_DimensionesParcialesNoDescartanTodo(t *testing.T) {
	out := Resolve(StrategyProductDimensions, nil, PackageInput{
		TotalQuantity: 1,
		Items: []PackageItem{
			{SKU: "SKU-A", Quantity: 1, Length: f(35), Width: f(25)},
		},
	})

	assert.Equal(t, PackageSourceProduct, out.Source)
	assert.Equal(t, 35.0, out.Length)
	assert.Equal(t, 25.0, out.Width)
	assert.Equal(t, DefaultDimCm, out.Height, "el lado faltante cae al default, no descarta los otros")
}

func TestResolvePackage_PesoDeLaCajaCuandoSuperaAlContenido(t *testing.T) {
	boxes := vigaBoxes()

	out := Resolve(StrategyStandardBox, boxes, PackageInput{TotalQuantity: 2, CartWeightKg: 0.4})

	assert.Equal(t, PackageSourceStandardBox, out.Source)
	assert.Equal(t, 1.0, out.Weight, "la caja declara 1 kg y el contenido pesa menos")
}

func TestResolvePackage_ProductoGrandeDescartaCajaPequena(t *testing.T) {
	boxes := vigaBoxes()

	out := Resolve(StrategyStandardBox, boxes, PackageInput{
		TotalQuantity: 2,
		CartWeightKg:  2,
		Items: []PackageItem{
			{SKU: "GRANDE", Quantity: 1, Length: f(38), Width: f(28), Height: f(25)},
		},
	})

	assert.Equal(t, "Caja Mediana", out.BoxName,
		"solo la mediana admite un item de 38x28x25")
}

func TestSelectBox_EquivalenteAlSelectorAnterior(t *testing.T) {
	boxes := vigaBoxes()

	cases := []struct {
		quantity int
		want     string
	}{
		{1, "Caja muy chica"},
		{3, "Caja muy chica"},
		{4, "Caja Chica"},
		{5, "Caja Chica"},
		{6, "Caja Mediana"},
		{8, "Caja Mediana"},
	}

	for _, tc := range cases {
		box := SelectBox(StrategyStandardBox, boxes, tc.quantity, 0, 0, 0)
		if assert.NotNil(t, box, "cantidad %d debe encontrar caja", tc.quantity) {
			assert.Equal(t, tc.want, box.Name, "cantidad %d", tc.quantity)
		}
	}

	overflow := SelectBox(StrategyStandardBox, boxes, 9, 0, 0, 0)
	if assert.NotNil(t, overflow, "9 items supera la caja mas grande, pero se combina con otra") {
		assert.Equal(t, "Caja Mediana + Caja muy chica", overflow.Name)
	}
}

func TestSelectBox_CombinaCajasCuandoNingunaSolaAlcanza(t *testing.T) {
	boxes := []Box{
		{Name: "Caja Chica", Weight: f(1), Length: f(15), Width: f(10), Height: f(15), MaxItems: 2},
		{Name: "Caja Mediana", Weight: f(2), Length: f(20), Width: f(20), Height: f(20), MaxItems: 4},
		{Name: "Caja Grande", Weight: f(4), Length: f(20), Width: f(20), Height: f(20), MaxItems: 8},
	}

	combo10 := SelectBox(StrategyStandardBox, boxes, 10, 0, 0, 0)
	if assert.NotNil(t, combo10) {
		assert.Equal(t, "Caja Grande + Caja Chica", combo10.Name)
		assert.Equal(t, 20.0, *combo10.Length)
		assert.Equal(t, 20.0, *combo10.Width)
		assert.Equal(t, 35.0, *combo10.Height, "alto de grande + alto de chica, van apiladas")
		assert.Equal(t, 5.0, *combo10.Weight, "peso base de grande + chica")
	}

	combo12 := SelectBox(StrategyStandardBox, boxes, 12, 0, 0, 0)
	if assert.NotNil(t, combo12) {
		assert.Equal(t, "Caja Grande + Caja Mediana", combo12.Name)
	}
}
