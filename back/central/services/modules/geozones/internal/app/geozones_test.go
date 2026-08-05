package app

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newGeozonesUseCase(repo *mocks.RepositoryMock, cache *mocks.DisplayCacheMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	if cache == nil {
		return New(repo, nil)
	}
	return New(repo, cache)
}

func newProbUseCase(
	repo *mocks.ProbabilityRepositoryMock,
	mainRepo *mocks.RepositoryMock,
	resolver *mocks.ResolverMock,
	cache *mocks.ProbabilityCacheMock,
) interface {
	GetProbability(ctx context.Context, req dtos.ProbabilityRequest) (*dtos.ProbabilityResult, error)
	GetOrderZone(ctx context.Context, orderID string, businessID uint) (*entities.Geozone, error)
	GetProbabilityByCarrier(ctx context.Context, orderID string, businessID uint) ([]dtos.ProbabilityResult, error)
} {
	if repo == nil {
		repo = &mocks.ProbabilityRepositoryMock{}
	}
	if mainRepo == nil {
		mainRepo = &mocks.RepositoryMock{}
	}
	if resolver == nil {
		resolver = &mocks.ResolverMock{}
	}
	if cache == nil {
		return NewProbability(repo, mainRepo, resolver, nil)
	}
	return NewProbability(repo, mainRepo, resolver, cache)
}

func f64Ptr(v float64) *float64 { return &v }
func uintPtr(v uint) *uint      { return &v }
func strPtr(v string) *string   { return &v }

func geometriaValida() json.RawMessage {
	return json.RawMessage(`{"type":"Polygon","coordinates":[[[0,0],[1,0],[1,1],[0,0]]]}`)
}

func TestValidType_AceptaSoloLosTiposConocidos(t *testing.T) {
	validos := []string{
		dtos.TypeCountry, dtos.TypeState, dtos.TypeCity, dtos.TypeAdminDistrict,
		dtos.TypeLocality, dtos.TypeNeighborhood, dtos.TypeBarrio, dtos.TypeCustom,
	}
	for _, tipo := range validos {
		assert.True(t, validType(tipo), "tipo %q deberia ser valido", tipo)
	}

	invalidos := []string{"", "COUNTRY", "pais", "departamento", "municipio", "zona", "  city  "}
	for _, tipo := range invalidos {
		assert.False(t, validType(tipo), "tipo %q no deberia ser valido", tipo)
	}
}

func TestParentTypeFor_DefineLaJerarquia(t *testing.T) {
	casos := []struct {
		hijo  string
		padre string
	}{
		{dtos.TypeState, dtos.TypeCountry},
		{dtos.TypeCity, dtos.TypeState},
		{dtos.TypeLocality, dtos.TypeCity},
		{dtos.TypeNeighborhood, dtos.TypeLocality},
		{dtos.TypeCountry, ""},
		{dtos.TypeBarrio, ""},
		{dtos.TypeCustom, ""},
		{dtos.TypeAdminDistrict, ""},
	}

	for _, tc := range casos {
		assert.Equal(t, tc.padre, parentTypeFor(tc.hijo), "padre de %q", tc.hijo)
	}
}

func TestCreate_TipoInvalido_NoLlegaAlRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	got, err := newGeozonesUseCase(repo, nil).Create(context.Background(),
		dtos.CreateGeozoneDTO{Type: "municipio", Name: "Bogota", Geometry: geometriaValida()})

	assert.ErrorIs(t, err, domainerrors.ErrInvalidType)
	assert.Nil(t, got)
	assert.Empty(t, repo.CreatedDTOs)
}

func TestCreate_GeometriaVacia_Rechaza(t *testing.T) {
	for _, geom := range []json.RawMessage{nil, {}} {
		repo := &mocks.RepositoryMock{}

		got, err := newGeozonesUseCase(repo, nil).Create(context.Background(),
			dtos.CreateGeozoneDTO{Type: dtos.TypeCity, Name: "Bogota", Geometry: geom})

		assert.ErrorIs(t, err, domainerrors.ErrInvalidGeometry)
		assert.Nil(t, got)
		assert.Empty(t, repo.CreatedDTOs)
	}
}

func TestCreate_Valido_PasaElDTOTalCual(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	dto := dtos.CreateGeozoneDTO{
		BusinessID: 26, Type: dtos.TypeCity, Code: strPtr("11001"),
		Name: "Bogota", Geometry: geometriaValida(),
	}

	got, err := newGeozonesUseCase(repo, nil).Create(context.Background(), dto)

	require.NoError(t, err)
	require.Len(t, repo.CreatedDTOs, 1)
	assert.Equal(t, dto, repo.CreatedDTOs[0])
	assert.Equal(t, "Bogota", got.Name)
}

func TestCreate_ErrorDelRepo_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, dto dtos.CreateGeozoneDTO) (*entities.Geozone, error) {
			return nil, domainerrors.ErrDuplicateGeozone
		},
	}

	_, err := newGeozonesUseCase(repo, nil).Create(context.Background(),
		dtos.CreateGeozoneDTO{Type: dtos.TypeCity, Name: "X", Geometry: geometriaValida()})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateGeozone)
}

func TestBulkImport_CuentaCreadosYSaltados(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	got, err := newGeozonesUseCase(repo, nil).BulkImport(context.Background(), dtos.BulkImportDTO{
		BusinessID: 26,
		Features: []dtos.BulkImportFeature{
			{Type: dtos.TypeCity, Name: "Bogota", Geometry: geometriaValida()},
			{Type: "municipio", Name: "Malo", Geometry: geometriaValida()},
			{Type: dtos.TypeCity, Name: "Sin geometria"},
			{Type: dtos.TypeCity, Name: "Medellin", Geometry: geometriaValida()},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, 2, got.Created)
	assert.Equal(t, 2, got.Skipped)
	require.Len(t, got.Errors, 2)
	assert.Contains(t, got.Errors[0], "tipo invalido")
	assert.Contains(t, got.Errors[1], "geometria vacia")
}

func TestBulkImport_UnaFeatureQueFallaNoAbortaElResto(t *testing.T) {
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, dto dtos.CreateGeozoneDTO) (*entities.Geozone, error) {
			if dto.Name == "Malo" {
				return nil, stderrors.New("constraint violation")
			}
			return &entities.Geozone{ID: 1, Name: dto.Name}, nil
		},
	}

	got, err := newGeozonesUseCase(repo, nil).BulkImport(context.Background(), dtos.BulkImportDTO{
		Features: []dtos.BulkImportFeature{
			{Type: dtos.TypeCity, Name: "Bogota", Geometry: geometriaValida()},
			{Type: dtos.TypeCity, Name: "Malo", Geometry: geometriaValida()},
			{Type: dtos.TypeCity, Name: "Cali", Geometry: geometriaValida()},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, 2, got.Created)
	assert.Equal(t, 1, got.Skipped)
	assert.Contains(t, got.Errors[0], "Malo")
}

func TestBulkImport_ResuelveElPadreDentroDelMismoLote(t *testing.T) {
	consultasAlRepo := 0
	repo := &mocks.RepositoryMock{
		GetByCodeFn: func(ctx context.Context, businessID uint, geozoneType, code string) (*entities.Geozone, error) {
			consultasAlRepo++
			return nil, nil
		},
	}

	got, err := newGeozonesUseCase(repo, nil).BulkImport(context.Background(), dtos.BulkImportDTO{
		BusinessID: 26,
		Features: []dtos.BulkImportFeature{
			{Type: dtos.TypeState, Code: strPtr("11"), Name: "Cundinamarca", Geometry: geometriaValida()},
			{Type: dtos.TypeCity, Code: strPtr("11001"), ParentCode: strPtr("11"), Name: "Bogota", Geometry: geometriaValida()},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, 2, got.Created)
	assert.Zero(t, consultasAlRepo,
		"el padre ya venia en el mismo lote, no hace falta ir a la base")

	require.Len(t, repo.CreatedDTOs, 2)
	require.NotNil(t, repo.CreatedDTOs[1].ParentID)
	assert.Equal(t, uint(1), *repo.CreatedDTOs[1].ParentID,
		"la ciudad queda colgando del departamento recien creado")
}

func TestBulkImport_BuscaElPadreEnLaBaseSiNoEstaEnElLote(t *testing.T) {
	var vistoTipo, vistoCode string
	var vistoNegocio uint
	repo := &mocks.RepositoryMock{
		GetByCodeFn: func(ctx context.Context, businessID uint, geozoneType, code string) (*entities.Geozone, error) {
			vistoNegocio, vistoTipo, vistoCode = businessID, geozoneType, code
			return &entities.Geozone{ID: 77}, nil
		},
	}

	_, err := newGeozonesUseCase(repo, nil).BulkImport(context.Background(), dtos.BulkImportDTO{
		BusinessID: 26,
		Features: []dtos.BulkImportFeature{
			{Type: dtos.TypeCity, ParentCode: strPtr("11"), Name: "Bogota", Geometry: geometriaValida()},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio)
	assert.Equal(t, dtos.TypeState, vistoTipo, "el padre de una ciudad se busca como departamento")
	assert.Equal(t, "11", vistoCode)
	require.NotNil(t, repo.CreatedDTOs[0].ParentID)
	assert.Equal(t, uint(77), *repo.CreatedDTOs[0].ParentID)
}

func TestBulkImport_PadreQueNoExiste_CreaSinPadre(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByCodeFn: func(ctx context.Context, businessID uint, geozoneType, code string) (*entities.Geozone, error) {
			return nil, stderrors.New("no encontrado")
		},
	}

	got, err := newGeozonesUseCase(repo, nil).BulkImport(context.Background(), dtos.BulkImportDTO{
		Features: []dtos.BulkImportFeature{
			{Type: dtos.TypeCity, ParentCode: strPtr("99"), Name: "Huerfana", Geometry: geometriaValida()},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, 1, got.Created)
	assert.Nil(t, repo.CreatedDTOs[0].ParentID,
		"si el padre no existe la geozona se crea igual, sin colgar de nadie")
}

func TestBulkImport_ParentCodeVacio_NoBuscaPadre(t *testing.T) {
	buscado := false
	repo := &mocks.RepositoryMock{
		GetByCodeFn: func(ctx context.Context, businessID uint, geozoneType, code string) (*entities.Geozone, error) {
			buscado = true
			return nil, nil
		},
	}

	_, err := newGeozonesUseCase(repo, nil).BulkImport(context.Background(), dtos.BulkImportDTO{
		Features: []dtos.BulkImportFeature{
			{Type: dtos.TypeCity, ParentCode: strPtr(""), Name: "X", Geometry: geometriaValida()},
		},
	})

	require.NoError(t, err)
	assert.False(t, buscado)
}

func TestBulkImport_LoteVacio_ResultadoEnCero(t *testing.T) {
	got, err := newGeozonesUseCase(nil, nil).BulkImport(context.Background(), dtos.BulkImportDTO{})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Zero(t, got.Created)
	assert.Zero(t, got.Skipped)
	assert.Empty(t, got.Errors)
}

func TestList_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 20, 1, 20},
		{-5, 20, 1, 20},
		{1, 0, 1, 20},
		{1, -1, 1, 20},
		{1, 101, 1, 100},
		{1, 100, 1, 100},
		{3, 50, 3, 50},
	}

	for _, tc := range casos {
		var visto dtos.ListGeozonesParams
		repo := &mocks.RepositoryMock{
			ListFn: func(ctx context.Context, params dtos.ListGeozonesParams) ([]entities.Geozone, int64, error) {
				visto = params
				return nil, 0, nil
			},
		}

		_, _, err := newGeozonesUseCase(repo, nil).List(context.Background(),
			dtos.ListGeozonesParams{Page: tc.page, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, visto.Page)
		assert.Equal(t, tc.wantSize, visto.PageSize)
	}
}

func TestListParams_Offset(t *testing.T) {
	assert.Equal(t, 0, dtos.ListGeozonesParams{Page: 1, PageSize: 20}.Offset())
	assert.Equal(t, 20, dtos.ListGeozonesParams{Page: 2, PageSize: 20}.Offset())
	assert.Equal(t, 200, dtos.ListGeozonesParams{Page: 3, PageSize: 100}.Offset())
}

func TestGetLookupYDelete_DeleganAlRepositorio(t *testing.T) {
	var vistoID uint
	var vistoGeom bool
	var vistoLookup dtos.LookupParams
	var borradoID uint
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint, includeGeom bool) (*entities.Geozone, error) {
			vistoID, vistoGeom = id, includeGeom
			return &entities.Geozone{ID: id}, nil
		},
		LookupContainingFn: func(ctx context.Context, params dtos.LookupParams) ([]entities.Geozone, error) {
			vistoLookup = params
			return []entities.Geozone{{ID: 1}}, nil
		},
		DeleteFn: func(ctx context.Context, id uint) error {
			borradoID = id
			return nil
		},
	}
	uc := newGeozonesUseCase(repo, nil)

	_, err := uc.Get(context.Background(), 5, true)
	require.NoError(t, err)
	assert.Equal(t, uint(5), vistoID)
	assert.True(t, vistoGeom)

	params := dtos.LookupParams{BusinessID: 26, Lat: 4.6, Lng: -74.1, Type: dtos.TypeCity}
	got, err := uc.Lookup(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, params, vistoLookup)
	assert.Len(t, got, 1)

	require.NoError(t, uc.Delete(context.Background(), 9))
	assert.Equal(t, uint(9), borradoID)
}

func TestToleranceForZoom_BajaDetalleEnZoomLejano(t *testing.T) {
	casos := []struct {
		zoom int
		want float64
	}{
		{0, 0.0005},
		{5, 0.0005},
		{7, 0.0005},
		{8, 0.0001},
		{9, 0.0001},
		{10, 0},
		{18, 0},
	}

	for _, tc := range casos {
		assert.InDelta(t, tc.want, toleranceForZoom(tc.zoom), 1e-9, "zoom=%d", tc.zoom)
	}
}

func TestZoomBucket_AgrupaLosZoomsEnTresNiveles(t *testing.T) {
	casos := []struct {
		zoom int
		want string
	}{
		{0, "z7"}, {7, "z7"},
		{8, "z9"}, {9, "z9"},
		{10, "zfull"}, {18, "zfull"},
	}

	for _, tc := range casos {
		assert.Equal(t, tc.want, zoomBucket(tc.zoom), "zoom=%d", tc.zoom)
	}
}

func TestZoomBucket_LaToleranciaYElBucketVanDeLaMano(t *testing.T) {
	for zoom := 0; zoom <= 20; zoom++ {
		bucket := zoomBucket(zoom)
		tol := toleranceForZoom(zoom)
		switch bucket {
		case "z7":
			assert.InDelta(t, 0.0005, tol, 1e-9, "zoom=%d", zoom)
		case "z9":
			assert.InDelta(t, 0.0001, tol, 1e-9, "zoom=%d", zoom)
		default:
			assert.Zero(t, tol, "zoom=%d", zoom)
		}
	}
}

func TestGetForDisplay_SinTipo_UsaAllEnLaClave(t *testing.T) {
	cache := &mocks.DisplayCacheMock{}

	_, etag, err := newGeozonesUseCase(nil, cache).GetForDisplay(context.Background(), "", 5, nil, nil)

	require.NoError(t, err)
	require.Len(t, cache.GetKeys, 1)
	assert.True(t, strings.HasPrefix(cache.GetKeys[0], "all:p0:z7:"), "clave: %q", cache.GetKeys[0])
	assert.Contains(t, etag, "all-p0-z7")
}

func TestGetForDisplay_LaClaveDistingueTipoPadreYZoom(t *testing.T) {
	casos := []struct {
		nombre   string
		tipo     string
		zoom     int
		parentID *uint
		want     string
	}{
		{"ciudad sin padre z5", dtos.TypeCity, 5, nil, "city:p0:z7:"},
		{"ciudad con padre z5", dtos.TypeCity, 5, uintPtr(11), "city:p11:z7:"},
		{"ciudad sin padre z9", dtos.TypeCity, 9, nil, "city:p0:z9:"},
		{"ciudad sin padre z14", dtos.TypeCity, 14, nil, "city:p0:zfull:"},
		{"barrio con padre z14", dtos.TypeBarrio, 14, uintPtr(3), "barrio:p3:zfull:"},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			cache := &mocks.DisplayCacheMock{}

			_, _, err := newGeozonesUseCase(nil, cache).
				GetForDisplay(context.Background(), tc.tipo, tc.zoom, nil, tc.parentID)

			require.NoError(t, err)
			require.Len(t, cache.GetKeys, 1)
			assert.True(t, strings.HasPrefix(cache.GetKeys[0], tc.want),
				"esperaba prefijo %q, obtuve %q", tc.want, cache.GetKeys[0])
		})
	}
}

func TestGetForDisplay_ClavesDistintasNoColisionan(t *testing.T) {
	vistas := map[string]bool{}
	for _, tipo := range []string{dtos.TypeCity, dtos.TypeBarrio} {
		for _, zoom := range []int{5, 9, 14} {
			for _, parent := range []*uint{nil, uintPtr(1), uintPtr(2)} {
				cache := &mocks.DisplayCacheMock{}
				_, _, err := newGeozonesUseCase(nil, cache).
					GetForDisplay(context.Background(), tipo, zoom, nil, parent)
				require.NoError(t, err)
				key := cache.GetKeys[0]
				require.False(t, vistas[key], "clave repetida para combinacion distinta: %q", key)
				vistas[key] = true
			}
		}
	}
	assert.Len(t, vistas, 18)
}

func TestGetForDisplay_ConBboxNoUsaCache(t *testing.T) {
	cache := &mocks.DisplayCacheMock{}
	bbox := &dtos.Bbox{}

	_, _, err := newGeozonesUseCase(nil, cache).
		GetForDisplay(context.Background(), dtos.TypeCity, 5, bbox, nil)

	require.NoError(t, err)
	assert.Empty(t, cache.GetKeys, "un recorte por bbox es especifico del viewport, no se cachea")
	assert.Empty(t, cache.SetKeys)
}

func TestGetForDisplay_ConCacheCaliente_NoVaAlRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	cache := &mocks.DisplayCacheMock{
		GetFn: func(ctx context.Context, key string) ([]byte, bool) {
			return []byte(`{"cached":true}`), true
		},
	}

	got, etag, err := newGeozonesUseCase(repo, cache).
		GetForDisplay(context.Background(), dtos.TypeCity, 5, nil, nil)

	require.NoError(t, err)
	assert.JSONEq(t, `{"cached":true}`, string(got))
	assert.NotEmpty(t, etag)
	assert.Empty(t, repo.DisplayParams, "con cache caliente no se toca PostGIS")
}

func TestGetForDisplay_CacheFrio_ConsultaYGuarda(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	cache := &mocks.DisplayCacheMock{}

	got, _, err := newGeozonesUseCase(repo, cache).
		GetForDisplay(context.Background(), dtos.TypeCity, 5, nil, nil)

	require.NoError(t, err)
	require.Len(t, repo.DisplayParams, 1)
	assert.Equal(t, dtos.TypeCity, repo.DisplayParams[0].Type)
	assert.InDelta(t, 0.0005, repo.DisplayParams[0].Tolerance, 1e-9)
	require.Len(t, cache.SetKeys, 1)
	assert.Equal(t, cache.GetKeys[0], cache.SetKeys[0], "se guarda bajo la misma clave que se consulto")
	assert.Contains(t, string(got), "FeatureCollection")
}

func TestGetForDisplay_SinResultados_DevuelveColeccionValida(t *testing.T) {
	got, _, err := newGeozonesUseCase(nil, nil).
		GetForDisplay(context.Background(), dtos.TypeCity, 5, nil, nil)

	require.NoError(t, err)
	var fc map[string]interface{}
	require.NoError(t, json.Unmarshal(got, &fc))
	assert.Equal(t, "FeatureCollection", fc["type"],
		"el mapa espera un GeoJSON valido aunque no haya features")
}

func TestGetForDisplay_ErrorDelRepo_NoGuardaEnCache(t *testing.T) {
	dbErr := stderrors.New("postgis caido")
	cache := &mocks.DisplayCacheMock{}
	repo := &mocks.RepositoryMock{
		GetForDisplayFn: func(ctx context.Context, params dtos.DisplayParams) ([]dtos.DisplayFeature, error) {
			return nil, dbErr
		},
	}

	got, etag, err := newGeozonesUseCase(repo, cache).
		GetForDisplay(context.Background(), dtos.TypeCity, 5, nil, nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
	assert.Empty(t, etag)
	assert.Empty(t, cache.SetKeys)
}

func TestGetForDisplay_SinCacheConfigurado_Funciona(t *testing.T) {
	got, etag, err := newGeozonesUseCase(nil, nil).
		GetForDisplay(context.Background(), dtos.TypeCity, 5, nil, nil)

	require.NoError(t, err)
	assert.NotEmpty(t, got)
	assert.NotEmpty(t, etag)
}

func TestFlushDisplayCache_SinCache_NoFalla(t *testing.T) {
	require.NoError(t, newGeozonesUseCase(nil, nil).FlushDisplayCache(context.Background()))
}

func TestFlushDisplayCache_ConCache_Delega(t *testing.T) {
	cache := &mocks.DisplayCacheMock{}

	require.NoError(t, newGeozonesUseCase(nil, cache).FlushDisplayCache(context.Background()))
	assert.Equal(t, 1, cache.Flushed)

	cacheErr := &mocks.DisplayCacheMock{
		FlushAllFn: func(ctx context.Context) error { return stderrors.New("redis caido") },
	}
	assert.Error(t, newGeozonesUseCase(nil, cacheErr).FlushDisplayCache(context.Background()))
}

func TestBaselineForCarrier_ConoceLasTransportadorasPrincipales(t *testing.T) {
	casos := []struct {
		carrier string
		want    float64
	}{
		{"SERVIENTREGA", 0.92},
		{"servientrega", 0.92},
		{"Servientrega", 0.92},
		{"Inter Rapidisimo", 0.89},
		{"inter-rapidisimo", 0.89},
		{"99 minutos", 0.91},
		{"TCC", 0.87},
		{"FedEx", 0.93},
	}

	for _, tc := range casos {
		assert.InDelta(t, tc.want, baselineForCarrier(tc.carrier), 1e-9, "carrier=%q", tc.carrier)
	}
}

func TestBaselineForCarrier_DesconocidaCaeAlDefault(t *testing.T) {
	for _, carrier := range []string{"", "TRANSPORTADORA_NUEVA", "xyz", "123"} {
		assert.InDelta(t, defaultBaseline, baselineForCarrier(carrier), 1e-9, "carrier=%q", carrier)
	}
}

func TestBaselineForCarrier_TodasLasBaselinesSonProbabilidadesValidas(t *testing.T) {
	for carrier, rate := range carrierBaselines {
		assert.Greater(t, rate, 0.0, "%s tiene baseline no positiva", carrier)
		assert.LessOrEqual(t, rate, 1.0, "%s tiene baseline sobre 1", carrier)
	}
	assert.Greater(t, defaultBaseline, 0.0)
	assert.LessOrEqual(t, defaultBaseline, 1.0)
}

func TestApplyBaselineCascade_TasaRealSuficiente_NoSeToca(t *testing.T) {
	results := []dtos.ProbabilityResult{
		{Found: true, Carrier: "SERVIENTREGA", DeliveryRate: f64Ptr(0.75)},
	}

	applyBaselineCascade(results)

	require.NotNil(t, results[0].DeliveryRate)
	assert.InDelta(t, 0.75, *results[0].DeliveryRate, 1e-9)
	assert.False(t, results[0].IsEstimated, "un dato real no se marca como estimado")
	assert.Empty(t, results[0].EstimateSource)
}

func TestApplyBaselineCascade_SinTasaLocal_CaeALaGlobalDelTransportador(t *testing.T) {
	results := []dtos.ProbabilityResult{
		{Carrier: "SERVIENTREGA", DeliveryRate: nil, GlobalRate: f64Ptr(0.81), GlobalTotal: 500},
	}

	applyBaselineCascade(results)

	require.NotNil(t, results[0].DeliveryRate)
	assert.InDelta(t, 0.81, *results[0].DeliveryRate, 1e-9)
	assert.True(t, results[0].IsEstimated)
	assert.Equal(t, "global_carrier", results[0].EstimateSource)
}

func TestApplyBaselineCascade_SinGlobalUtil_CaeAlBaselineDelTransportador(t *testing.T) {
	casos := []struct {
		nombre string
		result dtos.ProbabilityResult
	}{
		{"sin global", dtos.ProbabilityResult{Carrier: "SERVIENTREGA"}},
		{"global sin muestras", dtos.ProbabilityResult{Carrier: "SERVIENTREGA", GlobalRate: f64Ptr(0.81), GlobalTotal: 0}},
		{"global practicamente cero", dtos.ProbabilityResult{Carrier: "SERVIENTREGA", GlobalRate: f64Ptr(0.005), GlobalTotal: 100}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			results := []dtos.ProbabilityResult{tc.result}

			applyBaselineCascade(results)

			require.NotNil(t, results[0].DeliveryRate)
			assert.InDelta(t, 0.92, *results[0].DeliveryRate, 1e-9)
			assert.True(t, results[0].IsEstimated)
			assert.Equal(t, "carrier_baseline", results[0].EstimateSource)
		})
	}
}

func TestApplyBaselineCascade_TasaLocalPracticamenteCero_SeConsideraSinDato(t *testing.T) {
	results := []dtos.ProbabilityResult{
		{Carrier: "SERVIENTREGA", DeliveryRate: f64Ptr(0.005), GlobalRate: f64Ptr(0.81), GlobalTotal: 500},
	}

	applyBaselineCascade(results)

	require.NotNil(t, results[0].DeliveryRate)
	assert.InDelta(t, 0.81, *results[0].DeliveryRate, 1e-9,
		"una tasa de 0.5% no es un dato util, es ruido: se cae a la global")
	assert.True(t, results[0].IsEstimated)
}

func TestApplyBaselineCascade_ElUmbralEsExcluyente(t *testing.T) {
	justoEnElUmbral := []dtos.ProbabilityResult{{Carrier: "SERVIENTREGA", DeliveryRate: f64Ptr(0.01)}}
	applyBaselineCascade(justoEnElUmbral)
	assert.True(t, justoEnElUmbral[0].IsEstimated, "0.01 exacto NO supera el umbral, se estima")

	apenasEncima := []dtos.ProbabilityResult{{Carrier: "SERVIENTREGA", DeliveryRate: f64Ptr(0.0101)}}
	applyBaselineCascade(apenasEncima)
	assert.False(t, apenasEncima[0].IsEstimated, "apenas por encima del umbral si se respeta")
}

func TestApplyBaselineCascade_NuncaDejaTasaNil(t *testing.T) {
	results := []dtos.ProbabilityResult{
		{Carrier: "SERVIENTREGA"},
		{Carrier: "DESCONOCIDA"},
		{Carrier: ""},
		{Carrier: "TCC", GlobalRate: f64Ptr(0.7), GlobalTotal: 10},
	}

	applyBaselineCascade(results)

	for i := range results {
		require.NotNil(t, results[i].DeliveryRate,
			"resultado %d quedo sin tasa: la UI mostraria un dato vacio", i)
	}
}

func TestApplyBaselineCascade_ListaVacia_NoRompe(t *testing.T) {
	assert.NotPanics(t, func() { applyBaselineCascade(nil) })
	assert.NotPanics(t, func() { applyBaselineCascade([]dtos.ProbabilityResult{}) })
}

func TestGetProbabilityByCarrier_ConCacheCaliente_NoConsultaLaBase(t *testing.T) {
	repoLlamado := false
	repo := &mocks.ProbabilityRepositoryMock{
		AncestorsByOrderIDFn: func(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
			repoLlamado = true
			return nil, nil
		},
	}
	cache := &mocks.ProbabilityCacheMock{
		GetByOrderFn: func(ctx context.Context, businessID uint, orderID string) ([]dtos.ProbabilityResult, bool) {
			return []dtos.ProbabilityResult{{Carrier: "SERVIENTREGA", Found: true}}, true
		},
	}

	got, err := newProbUseCase(repo, nil, nil, cache).
		GetProbabilityByCarrier(context.Background(), "ORD-1", 26)

	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.False(t, repoLlamado)
}

func TestGetProbabilityByCarrier_SinZonaResuelta_DevuelveVacioNoNil(t *testing.T) {
	repo := &mocks.ProbabilityRepositoryMock{
		AncestorsByOrderIDFn: func(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
			return nil, nil
		},
	}

	got, err := newProbUseCase(repo, nil, nil, nil).
		GetProbabilityByCarrier(context.Background(), "ORD-1", 26)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Empty(t, got)
}

func TestGetProbabilityByCarrier_AplicaLaCascadaYGuardaEnCache(t *testing.T) {
	cache := &mocks.ProbabilityCacheMock{}
	repo := &mocks.ProbabilityRepositoryMock{
		ProbabilityByOrderFn: func(ctx context.Context, a *entities.GeozoneAncestors) ([]dtos.ProbabilityResult, error) {
			return []dtos.ProbabilityResult{{Carrier: "SERVIENTREGA"}}, nil
		},
	}

	got, err := newProbUseCase(repo, nil, nil, cache).
		GetProbabilityByCarrier(context.Background(), "ORD-1", 26)

	require.NoError(t, err)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].DeliveryRate)
	assert.True(t, got[0].IsEstimated)
	assert.Equal(t, 1, cache.SetCalls)
}

func TestGetProbabilityByCarrier_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("timeout")

	t.Run("al resolver ancestros", func(t *testing.T) {
		repo := &mocks.ProbabilityRepositoryMock{
			AncestorsByOrderIDFn: func(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
				return nil, dbErr
			},
		}
		_, err := newProbUseCase(repo, nil, nil, nil).
			GetProbabilityByCarrier(context.Background(), "ORD-1", 26)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al calcular probabilidad", func(t *testing.T) {
		repo := &mocks.ProbabilityRepositoryMock{
			ProbabilityByOrderFn: func(ctx context.Context, a *entities.GeozoneAncestors) ([]dtos.ProbabilityResult, error) {
				return nil, dbErr
			},
		}
		_, err := newProbUseCase(repo, nil, nil, nil).
			GetProbabilityByCarrier(context.Background(), "ORD-1", 26)
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestGetOrderZone_SinZonaProfunda_DevuelveNilSinError(t *testing.T) {
	casos := []struct {
		nombre    string
		ancestors *entities.GeozoneAncestors
	}{
		{"sin ancestros", nil},
		{"deepest nil", &entities.GeozoneAncestors{DeepestID: nil}},
		{"deepest cero", &entities.GeozoneAncestors{DeepestID: uintPtr(0)}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			consultado := false
			mainRepo := &mocks.RepositoryMock{
				GetByIDFn: func(ctx context.Context, id uint, includeGeom bool) (*entities.Geozone, error) {
					consultado = true
					return nil, nil
				},
			}
			repo := &mocks.ProbabilityRepositoryMock{
				AncestorsByOrderIDFn: func(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
					return tc.ancestors, nil
				},
			}

			got, err := newProbUseCase(repo, mainRepo, nil, nil).
				GetOrderZone(context.Background(), "ORD-1", 26)

			require.NoError(t, err)
			assert.Nil(t, got)
			assert.False(t, consultado)
		})
	}
}

func TestGetOrderZone_ConZona_TraeLaGeometria(t *testing.T) {
	var vistoID uint
	var vistoGeom bool
	mainRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint, includeGeom bool) (*entities.Geozone, error) {
			vistoID, vistoGeom = id, includeGeom
			return &entities.Geozone{ID: id, Name: "Chapinero"}, nil
		},
	}
	repo := &mocks.ProbabilityRepositoryMock{
		AncestorsByOrderIDFn: func(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
			return &entities.GeozoneAncestors{DeepestID: uintPtr(42)}, nil
		},
	}

	got, err := newProbUseCase(repo, mainRepo, nil, nil).GetOrderZone(context.Background(), "ORD-1", 26)

	require.NoError(t, err)
	assert.Equal(t, uint(42), vistoID)
	assert.True(t, vistoGeom, "la zona se pinta en el mapa, necesita geometria")
	assert.Equal(t, "Chapinero", got.Name)
}

func TestGetOrderZone_ErrorAlResolver_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.ProbabilityRepositoryMock{
		AncestorsByOrderIDFn: func(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
			return nil, dbErr
		},
	}

	_, err := newProbUseCase(repo, nil, nil, nil).GetOrderZone(context.Background(), "ORD-1", 26)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetProbability_SinOrdenNiCoordenadas_Rechaza(t *testing.T) {
	got, err := newProbUseCase(nil, nil, nil, nil).
		GetProbability(context.Background(), dtos.ProbabilityRequest{BusinessID: 26})

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "order_id or lat/lng")
}

func TestGetProbability_ConOrdenNoUsaElResolver(t *testing.T) {
	resolver := &mocks.ResolverMock{}

	_, err := newProbUseCase(nil, nil, resolver, nil).GetProbability(context.Background(),
		dtos.ProbabilityRequest{BusinessID: 26, OrderID: "ORD-1"})

	require.NoError(t, err)
	assert.Zero(t, resolver.Calls)
}

func TestGetProbability_ConCoordenadasUsaElResolver(t *testing.T) {
	var vistoLat, vistoLng float64
	var vistoNegocio uint
	resolver := &mocks.ResolverMock{
		ResolveFn: func(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error) {
			vistoLat, vistoLng, vistoNegocio = lat, lng, businessID
			return &entities.GeozoneAncestors{}, nil
		},
	}

	_, err := newProbUseCase(nil, nil, resolver, nil).GetProbability(context.Background(),
		dtos.ProbabilityRequest{BusinessID: 26, Lat: f64Ptr(4.65), Lng: f64Ptr(-74.05)})

	require.NoError(t, err)
	assert.Equal(t, 1, resolver.Calls)
	assert.InDelta(t, 4.65, vistoLat, 1e-9)
	assert.InDelta(t, -74.05, vistoLng, 1e-9)
	assert.Equal(t, uint(26), vistoNegocio)
}

func TestGetProbability_SoloLatSinLng_Rechaza(t *testing.T) {
	casos := []dtos.ProbabilityRequest{
		{BusinessID: 26, Lat: f64Ptr(4.65)},
		{BusinessID: 26, Lng: f64Ptr(-74.05)},
	}

	for _, req := range casos {
		_, err := newProbUseCase(nil, nil, nil, nil).GetProbability(context.Background(), req)
		require.Error(t, err, "una coordenada a medias no sirve para ubicar nada")
	}
}

func TestGetProbability_SinZonaResuelta_DevuelveNoEncontrado(t *testing.T) {
	repo := &mocks.ProbabilityRepositoryMock{
		AncestorsByOrderIDFn: func(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
			return nil, nil
		},
	}

	got, err := newProbUseCase(repo, nil, nil, nil).GetProbability(context.Background(),
		dtos.ProbabilityRequest{BusinessID: 26, OrderID: "ORD-1", Carrier: "Servientrega"})

	require.NoError(t, err)
	assert.False(t, got.Found)
	assert.Equal(t, "Servientrega", got.Carrier, "se devuelve el transportador pedido aunque no haya zona")
	assert.Nil(t, got.DeliveryRate, "sin zona no se inventa una tasa")
}

func TestGetProbability_SinTransportador_PrefiereElPrimeroEncontrado(t *testing.T) {
	repo := &mocks.ProbabilityRepositoryMock{
		ProbabilityByOrderFn: func(ctx context.Context, a *entities.GeozoneAncestors) ([]dtos.ProbabilityResult, error) {
			return []dtos.ProbabilityResult{
				{Carrier: "A", Found: false},
				{Carrier: "B", Found: true, DeliveryRate: f64Ptr(0.9)},
				{Carrier: "C", Found: true, DeliveryRate: f64Ptr(0.8)},
			}, nil
		},
	}

	got, err := newProbUseCase(repo, nil, nil, nil).GetProbability(context.Background(),
		dtos.ProbabilityRequest{BusinessID: 26, OrderID: "ORD-1"})

	require.NoError(t, err)
	assert.Equal(t, "B", got.Carrier, "se prefiere el primero con dato real, no el primero de la lista")
}

func TestGetProbability_SinTransportadorNiEncontrados_DevuelveElPrimero(t *testing.T) {
	repo := &mocks.ProbabilityRepositoryMock{
		ProbabilityByOrderFn: func(ctx context.Context, a *entities.GeozoneAncestors) ([]dtos.ProbabilityResult, error) {
			return []dtos.ProbabilityResult{{Carrier: "A", Found: false}, {Carrier: "B", Found: false}}, nil
		},
	}

	got, err := newProbUseCase(repo, nil, nil, nil).GetProbability(context.Background(),
		dtos.ProbabilityRequest{BusinessID: 26, OrderID: "ORD-1"})

	require.NoError(t, err)
	assert.Equal(t, "A", got.Carrier)
}

func TestGetProbability_SinTransportadorNiResultados_DevuelveNoEncontrado(t *testing.T) {
	repo := &mocks.ProbabilityRepositoryMock{
		ProbabilityByOrderFn: func(ctx context.Context, a *entities.GeozoneAncestors) ([]dtos.ProbabilityResult, error) {
			return nil, nil
		},
	}

	got, err := newProbUseCase(repo, nil, nil, nil).GetProbability(context.Background(),
		dtos.ProbabilityRequest{BusinessID: 26, OrderID: "ORD-1"})

	require.NoError(t, err)
	assert.False(t, got.Found)
}

func TestGetProbability_NormalizaElNombreDelTransportador(t *testing.T) {
	casos := []struct {
		entrada string
		want    string
	}{
		{"Servientrega", "SERVIENTREGA"},
		{"inter rapidisimo", "INTERRAPIDISIMO"},
		{"Inter-Rapidisimo", "INTERRAPIDISIMO"},
		{"99 minutos", "99MINUTOS"},
	}

	for _, tc := range casos {
		t.Run(tc.entrada, func(t *testing.T) {
			repo := &mocks.ProbabilityRepositoryMock{}

			_, err := newProbUseCase(repo, nil, nil, nil).GetProbability(context.Background(),
				dtos.ProbabilityRequest{BusinessID: 26, OrderID: "ORD-1", Carrier: tc.entrada})

			require.NoError(t, err)
			require.Len(t, repo.CarrierKeys, 1)
			assert.Equal(t, tc.want, repo.CarrierKeys[0])
		})
	}
}

func TestGetProbability_ConTransportadorSinDato_AplicaLaCascada(t *testing.T) {
	repo := &mocks.ProbabilityRepositoryMock{
		ProbabilityForCarrierFn: func(ctx context.Context, a *entities.GeozoneAncestors, key string) (*dtos.ProbabilityResult, error) {
			return nil, nil
		},
	}

	got, err := newProbUseCase(repo, nil, nil, nil).GetProbability(context.Background(),
		dtos.ProbabilityRequest{BusinessID: 26, OrderID: "ORD-1", Carrier: "Servientrega"})

	require.NoError(t, err)
	assert.False(t, got.Found)
	assert.Equal(t, "Servientrega", got.Carrier, "se conserva el nombre que pidio el usuario")
	require.NotNil(t, got.DeliveryRate)
	assert.InDelta(t, 0.92, *got.DeliveryRate, 1e-9)
	assert.Equal(t, "carrier_baseline", got.EstimateSource)
}

func TestGetProbability_ConTransportadorYDatoReal_NoLoPisa(t *testing.T) {
	repo := &mocks.ProbabilityRepositoryMock{
		ProbabilityForCarrierFn: func(ctx context.Context, a *entities.GeozoneAncestors, key string) (*dtos.ProbabilityResult, error) {
			return &dtos.ProbabilityResult{Found: true, Carrier: "SERVIENTREGA", DeliveryRate: f64Ptr(0.73)}, nil
		},
	}

	got, err := newProbUseCase(repo, nil, nil, nil).GetProbability(context.Background(),
		dtos.ProbabilityRequest{BusinessID: 26, OrderID: "ORD-1", Carrier: "Servientrega"})

	require.NoError(t, err)
	assert.InDelta(t, 0.73, *got.DeliveryRate, 1e-9)
	assert.False(t, got.IsEstimated)
}

func TestGetProbability_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("timeout")

	t.Run("al resolver por orden", func(t *testing.T) {
		repo := &mocks.ProbabilityRepositoryMock{
			AncestorsByOrderIDFn: func(ctx context.Context, orderID string, businessID uint) (*entities.GeozoneAncestors, error) {
				return nil, dbErr
			},
		}
		_, err := newProbUseCase(repo, nil, nil, nil).GetProbability(context.Background(),
			dtos.ProbabilityRequest{OrderID: "ORD-1"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al resolver por coordenadas", func(t *testing.T) {
		resolver := &mocks.ResolverMock{
			ResolveFn: func(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error) {
				return nil, dbErr
			},
		}
		_, err := newProbUseCase(nil, nil, resolver, nil).GetProbability(context.Background(),
			dtos.ProbabilityRequest{Lat: f64Ptr(4.6), Lng: f64Ptr(-74.1)})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al calcular sin transportador", func(t *testing.T) {
		repo := &mocks.ProbabilityRepositoryMock{
			ProbabilityByOrderFn: func(ctx context.Context, a *entities.GeozoneAncestors) ([]dtos.ProbabilityResult, error) {
				return nil, dbErr
			},
		}
		_, err := newProbUseCase(repo, nil, nil, nil).GetProbability(context.Background(),
			dtos.ProbabilityRequest{OrderID: "ORD-1"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al calcular con transportador", func(t *testing.T) {
		repo := &mocks.ProbabilityRepositoryMock{
			ProbabilityForCarrierFn: func(ctx context.Context, a *entities.GeozoneAncestors, key string) (*dtos.ProbabilityResult, error) {
				return nil, dbErr
			},
		}
		_, err := newProbUseCase(repo, nil, nil, nil).GetProbability(context.Background(),
			dtos.ProbabilityRequest{OrderID: "ORD-1", Carrier: "Servientrega"})
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestResolver_DelegaAlRepositorio(t *testing.T) {
	var vistoLat, vistoLng float64
	var vistoNegocio uint
	repo := &mocks.RepositoryMock{
		ResolveAncestorsFn: func(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error) {
			vistoLat, vistoLng, vistoNegocio = lat, lng, businessID
			return &entities.GeozoneAncestors{DeepestID: uintPtr(9)}, nil
		},
	}

	got, err := NewResolver(repo).Resolve(context.Background(), 4.65, -74.05, 26)

	require.NoError(t, err)
	assert.InDelta(t, 4.65, vistoLat, 1e-9)
	assert.InDelta(t, -74.05, vistoLng, 1e-9)
	assert.Equal(t, uint(26), vistoNegocio)
	require.NotNil(t, got.DeepestID)
	assert.Equal(t, uint(9), *got.DeepestID)
}

func TestResolver_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("postgis caido")
	repo := &mocks.RepositoryMock{
		ResolveAncestorsFn: func(ctx context.Context, lat, lng float64, businessID uint) (*entities.GeozoneAncestors, error) {
			return nil, dbErr
		},
	}

	_, err := NewResolver(repo).Resolve(context.Background(), 4.65, -74.05, 26)

	assert.ErrorIs(t, err, dbErr)
}
