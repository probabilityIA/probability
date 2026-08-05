package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func f64(v float64) *float64 {
	return &v
}

func repoParaCubing(dims *ports.ProductDimensions, cap *ports.LocationCapacityInfo, ocupado int) *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "SKU-1", "Producto", true, nil
		},
		GetProductDimensionsFn: func(ctx context.Context, productID string, businessID uint) (*ports.ProductDimensions, error) {
			return dims, nil
		},
		GetLocationCapacityFn: func(ctx context.Context, locationID uint) (*ports.LocationCapacityInfo, error) {
			return cap, nil
		},
		GetLocationOccupiedQtyFn: func(ctx context.Context, locationID uint) (int, error) {
			return ocupado, nil
		},
	}
}

func cubingDTO(qty int) request.ValidateCubingDTO {
	return request.ValidateCubingDTO{
		ProductID:  "prod-001",
		LocationID: 3,
		BusinessID: 10,
		Quantity:   qty,
	}
}

func TestConvertWeightToKg(t *testing.T) {
	cases := []struct {
		unidad string
		valor  float64
		want   float64
	}{
		{unidad: "kg", valor: 2.5, want: 2.5},
		{unidad: "", valor: 2.5, want: 2.5},
		{unidad: "g", valor: 2500, want: 2.5},
		{unidad: "lb", valor: 1, want: 0.453592},
		{unidad: "oz", valor: 1, want: 0.0283495},
		{unidad: "unidad-rara", valor: 7, want: 7},
		{unidad: "g", valor: 0, want: 0},
	}

	for _, tc := range cases {
		t.Run(tc.unidad, func(t *testing.T) {
			assert.InDelta(t, tc.want, convertWeightToKg(tc.valor, tc.unidad), 0.000001)
		})
	}
}

func TestConvertVolumeToCm3(t *testing.T) {
	cases := []struct {
		nombre string
		unidad string
		l      float64
		w      float64
		h      float64
		want   float64
	}{
		{nombre: "centimetros", unidad: "cm", l: 10, w: 10, h: 10, want: 1000},
		{nombre: "sin unidad usa cm", unidad: "", l: 10, w: 10, h: 10, want: 1000},
		{nombre: "milimetros", unidad: "mm", l: 100, w: 100, h: 100, want: 1000},
		{nombre: "metros", unidad: "m", l: 1, w: 1, h: 1, want: 1000000},
		{nombre: "pulgadas", unidad: "in", l: 1, w: 1, h: 1, want: 16.387064},
		{nombre: "unidad desconocida usa factor 1", unidad: "yardas", l: 2, w: 2, h: 2, want: 8},
		{nombre: "dimension cero", unidad: "cm", l: 10, w: 0, h: 10, want: 0},
	}

	for _, tc := range cases {
		t.Run(tc.nombre, func(t *testing.T) {
			assert.InDelta(t, tc.want, convertVolumeToCm3(tc.l, tc.w, tc.h, tc.unidad), 0.0001)
		})
	}
}

func TestValidateCubing_CantidadInvalida_RetornaError(t *testing.T) {
	uc := buildUseCase(&mocks.RepositoryMock{}, nil, nil)

	for _, qty := range []int{0, -1, -100} {
		got, err := uc.ValidateCubing(context.Background(), cubingDTO(qty))

		assert.ErrorIs(t, err, domainerrors.ErrInvalidQuantity)
		assert.Nil(t, got)
	}
}

func TestValidateCubing_ProductoInexistente_RetornaError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "", "", false, errors.New("not found")
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(1))

	assert.ErrorIs(t, err, domainerrors.ErrProductNotFound)
	assert.Nil(t, got)
}

func TestValidateCubing_CabeEnLaUbicacion(t *testing.T) {
	repo := repoParaCubing(
		&ports.ProductDimensions{Weight: 1, WeightU: "kg", Length: 10, Width: 10, Height: 10, DimU: "cm"},
		&ports.LocationCapacityInfo{MaxWeightKg: f64(100), MaxVolumeCm3: f64(100000)},
		7,
	)
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(5))

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.Fits)
	assert.Equal(t, "", got.Reason)
	assert.InDelta(t, 5.0, got.WeightNeededKg, 0.001)
	assert.InDelta(t, 5000.0, got.VolumeNeededCm3, 0.001)
	assert.Equal(t, 7, got.OccupiedQty)
	assert.InDelta(t, 100.0, got.WeightMaxKg, 0.001)
	assert.InDelta(t, 100000.0, got.VolumeMaxCm3, 0.001)
}

func TestValidateCubing_ExcedePeso(t *testing.T) {
	repo := repoParaCubing(
		&ports.ProductDimensions{Weight: 30, WeightU: "kg", Length: 1, Width: 1, Height: 1, DimU: "cm"},
		&ports.LocationCapacityInfo{MaxWeightKg: f64(100), MaxVolumeCm3: f64(100000)},
		0,
	)
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(5))

	require.NoError(t, err)
	assert.False(t, got.Fits)
	assert.Equal(t, "weight exceeds location capacity", got.Reason)
}

func TestValidateCubing_ExcedeVolumen(t *testing.T) {
	repo := repoParaCubing(
		&ports.ProductDimensions{Weight: 0.1, WeightU: "kg", Length: 50, Width: 50, Height: 50, DimU: "cm"},
		&ports.LocationCapacityInfo{MaxWeightKg: f64(100), MaxVolumeCm3: f64(100000)},
		0,
	)
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(5))

	require.NoError(t, err)
	assert.False(t, got.Fits)
	assert.Equal(t, "volume exceeds location capacity", got.Reason)
}

func TestValidateCubing_ExcedePesoYVolumen_MensajeCombinado(t *testing.T) {
	repo := repoParaCubing(
		&ports.ProductDimensions{Weight: 50, WeightU: "kg", Length: 50, Width: 50, Height: 50, DimU: "cm"},
		&ports.LocationCapacityInfo{MaxWeightKg: f64(100), MaxVolumeCm3: f64(100000)},
		0,
	)
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(5))

	require.NoError(t, err)
	assert.False(t, got.Fits)
	assert.Equal(t, "weight and volume exceed location capacity", got.Reason)
}

func TestValidateCubing_SinLimitesDeCapacidad_SiempreCabe(t *testing.T) {
	repo := repoParaCubing(
		&ports.ProductDimensions{Weight: 1000, WeightU: "kg", Length: 500, Width: 500, Height: 500, DimU: "cm"},
		&ports.LocationCapacityInfo{},
		0,
	)
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(100))

	require.NoError(t, err)
	assert.True(t, got.Fits, "sin limites configurados la ubicacion acepta cualquier carga")
	assert.InDelta(t, 0.0, got.WeightMaxKg, 0.001)
	assert.InDelta(t, 0.0, got.VolumeMaxCm3, 0.001)
}

func TestValidateCubing_SoloLimiteDePeso(t *testing.T) {
	repo := repoParaCubing(
		&ports.ProductDimensions{Weight: 1, WeightU: "kg", Length: 100, Width: 100, Height: 100, DimU: "cm"},
		&ports.LocationCapacityInfo{MaxWeightKg: f64(10)},
		0,
	)
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(5))

	require.NoError(t, err)
	assert.True(t, got.Fits)
	assert.InDelta(t, 5.0, got.WeightNeededKg, 0.001)
}

func TestValidateCubing_ConvierteUnidadesAntesDeComparar(t *testing.T) {
	repo := repoParaCubing(
		&ports.ProductDimensions{Weight: 500, WeightU: "g", Length: 100, Width: 100, Height: 100, DimU: "mm"},
		&ports.LocationCapacityInfo{MaxWeightKg: f64(10), MaxVolumeCm3: f64(10000)},
		0,
	)
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(4))

	require.NoError(t, err)
	assert.True(t, got.Fits)
	assert.InDelta(t, 2.0, got.WeightNeededKg, 0.001, "500 g por unidad son 0.5 kg")
	assert.InDelta(t, 4000.0, got.VolumeNeededCm3, 0.001, "100 mm de lado son 10 cm, o sea 1000 cm3")
}

func TestValidateCubing_FallaLeerDimensiones_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "SKU-1", "Producto", true, nil
		},
		GetProductDimensionsFn: func(ctx context.Context, productID string, businessID uint) (*ports.ProductDimensions, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(1))

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestValidateCubing_FallaLeerCapacidad_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "SKU-1", "Producto", true, nil
		},
		GetProductDimensionsFn: func(ctx context.Context, productID string, businessID uint) (*ports.ProductDimensions, error) {
			return &ports.ProductDimensions{}, nil
		},
		GetLocationCapacityFn: func(ctx context.Context, locationID uint) (*ports.LocationCapacityInfo, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(1))

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestValidateCubing_FallaLeerOcupacion_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "SKU-1", "Producto", true, nil
		},
		GetProductDimensionsFn: func(ctx context.Context, productID string, businessID uint) (*ports.ProductDimensions, error) {
			return &ports.ProductDimensions{}, nil
		},
		GetLocationCapacityFn: func(ctx context.Context, locationID uint) (*ports.LocationCapacityInfo, error) {
			return &ports.LocationCapacityInfo{}, nil
		},
		GetLocationOccupiedQtyFn: func(ctx context.Context, locationID uint) (int, error) {
			return 0, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ValidateCubing(context.Background(), cubingDTO(1))

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}
