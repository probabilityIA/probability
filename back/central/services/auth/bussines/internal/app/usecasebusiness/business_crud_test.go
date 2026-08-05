package usecasebusiness

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strPtr(v string) *string { return &v }
func uintPtr(v uint) *uint    { return &v }
func boolPtr(v bool) *bool    { return &v }

func negocioExistente() *domain.Business {
	return &domain.Business{
		ID: 26, Name: "Demo", Code: "demo", BusinessTypeID: 1,
		Timezone: "America/Bogota", Address: "Cra 1", Description: "desc",
		LogoURL: "businessLogo/viejo.png", NavbarImageURL: "navbar/viejo.png",
		PrimaryColor: "#111", SecondaryColor: "#222", TertiaryColor: "#333",
		QuaternaryColor: "#444", SidebarColor: "#555", TopbarColor: "#666",
		CustomDomain: "demo.com", IsActive: true,
		EnableDelivery: true, EnablePickup: false, EnableReservations: true,
	}
}

func repoConNegocio() *mocks.BusinessRepositoryMock {
	return &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return negocioExistente(), nil
		},
	}
}

func TestGetBusinesses_PasaTodosLosFiltrosAlRepositorio(t *testing.T) {
	var vistoPage, vistoPerPage int
	var vistoName string
	var vistoTipo *uint
	var vistoActivo *bool
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessesFn: func(ctx context.Context, page, perPage int, name string, businessTypeID *uint, isActive *bool) ([]domain.Business, int64, error) {
			vistoPage, vistoPerPage, vistoName = page, perPage, name
			vistoTipo, vistoActivo = businessTypeID, isActive
			return nil, 0, nil
		},
	}

	_, _, err := newBusinessUseCase(repo, nil, nil).GetBusinesses(
		context.Background(), 3, 25, "demo", uintPtr(1), boolPtr(true))

	require.NoError(t, err)
	assert.Equal(t, 3, vistoPage)
	assert.Equal(t, 25, vistoPerPage)
	assert.Equal(t, "demo", vistoName)
	require.NotNil(t, vistoTipo)
	assert.Equal(t, uint(1), *vistoTipo)
	require.NotNil(t, vistoActivo)
	assert.True(t, *vistoActivo)
}

func TestGetBusinesses_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessesFn: func(ctx context.Context, page, perPage int, name string, businessTypeID *uint, isActive *bool) ([]domain.Business, int64, error) {
			return nil, 0, dbErr
		},
	}

	got, total, err := newBusinessUseCase(repo, nil, nil).GetBusinesses(context.Background(), 1, 10, "", nil, nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
	assert.Zero(t, total)
}

func TestGetBusinesses_CompletaUrlsRelativasYRespetaAbsolutas(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessesFn: func(ctx context.Context, page, perPage int, name string, businessTypeID *uint, isActive *bool) ([]domain.Business, int64, error) {
			return []domain.Business{
				{ID: 1, LogoURL: "businessLogo/a.png", NavbarImageURL: "navbar/a.png"},
				{ID: 2, LogoURL: "https://cdn.externo.com/b.png", NavbarImageURL: "https://cdn.externo.com/c.png"},
				{ID: 3, LogoURL: "", NavbarImageURL: ""},
			}, 3, nil
		},
	}

	got, total, err := newBusinessUseCase(repo, nil, s3Config()).GetBusinesses(context.Background(), 1, 10, "", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, s3Base+"/businessLogo/a.png", got[0].LogoURL)
	assert.Equal(t, s3Base+"/navbar/a.png", got[0].NavbarImageURL)
	assert.Equal(t, "https://cdn.externo.com/b.png", got[1].LogoURL)
	assert.Empty(t, got[2].LogoURL, "una url vacia se queda vacia, no se le antepone el dominio")
}

func TestGetBusinesses_SinDominioS3_DejaLosPathsRelativos(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessesFn: func(ctx context.Context, page, perPage int, name string, businessTypeID *uint, isActive *bool) ([]domain.Business, int64, error) {
			return []domain.Business{{ID: 1, LogoURL: "businessLogo/a.png", NavbarImageURL: "navbar/a.png"}}, 1, nil
		},
	}

	got, _, err := newBusinessUseCase(repo, nil, nil).GetBusinesses(context.Background(), 1, 10, "", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, "businessLogo/a.png", got[0].LogoURL)
	assert.Equal(t, "navbar/a.png", got[0].NavbarImageURL)
}

func TestGetBusinesses_MapeaBusinessTypeCargadoYNoCargado(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessesFn: func(ctx context.Context, page, perPage int, name string, businessTypeID *uint, isActive *bool) ([]domain.Business, int64, error) {
			return []domain.Business{
				{ID: 1, BusinessTypeID: 4, BusinessType: nil},
				{ID: 2, BusinessTypeID: 4, BusinessType: &domain.BusinessType{ID: 4, Name: "Retail", Code: "retail", IsActive: true}},
			}, 2, nil
		},
	}

	got, _, err := newBusinessUseCase(repo, nil, nil).GetBusinesses(context.Background(), 1, 10, "", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, uint(4), got[0].BusinessType.ID)
	assert.Empty(t, got[0].BusinessType.Name)
	assert.Equal(t, "Retail", got[1].BusinessType.Name)
	assert.True(t, got[1].BusinessType.IsActive)
}

func TestGetBusinesses_SinResultados_DevuelveVacioNoNil(t *testing.T) {
	got, total, err := newBusinessUseCase(nil, nil, nil).GetBusinesses(context.Background(), 1, 10, "", nil, nil)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Empty(t, got)
	assert.Zero(t, total)
}

func TestGetBusinessByID_NegocioInexistente_RetornaError(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) { return nil, nil },
	}

	got, err := newBusinessUseCase(repo, nil, nil).GetBusinessByID(context.Background(), 404)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "negocio no encontrado")
}

func TestGetBusinessByID_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) { return nil, dbErr },
	}

	_, err := newBusinessUseCase(repo, nil, nil).GetBusinessByID(context.Background(), 1)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetBusinessByID_MapeaTodosLosCampos(t *testing.T) {
	got, err := newBusinessUseCase(repoConNegocio(), nil, s3Config()).GetBusinessByID(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, uint(26), got.ID)
	assert.Equal(t, "Demo", got.Name)
	assert.Equal(t, "demo", got.Code)
	assert.Equal(t, "America/Bogota", got.Timezone)
	assert.Equal(t, s3Base+"/businessLogo/viejo.png", got.LogoURL)
	assert.Equal(t, s3Base+"/navbar/viejo.png", got.NavbarImageURL)
	assert.Equal(t, "#555", got.SidebarColor)
	assert.Equal(t, "#666", got.TopbarColor)
	assert.Equal(t, "demo.com", got.CustomDomain)
	assert.True(t, got.EnableDelivery)
	assert.False(t, got.EnablePickup)
	assert.True(t, got.EnableReservations)
}

func TestGetBusinessByID_ConBusinessTypeCargado_LoMapea(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return &domain.Business{
				ID: id, BusinessTypeID: 4,
				BusinessType: &domain.BusinessType{ID: 4, Name: "Retail", Code: "retail", Description: "d", Icon: "i", IsActive: true},
			}, nil
		},
	}

	got, err := newBusinessUseCase(repo, nil, nil).GetBusinessByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, "Retail", got.BusinessType.Name)
	assert.Equal(t, "i", got.BusinessType.Icon)
}

func TestUpdateBusiness_NegocioInexistente_NoActualiza(t *testing.T) {
	actualizado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) { return nil, nil },
		UpdateBusinessFn: func(ctx context.Context, id uint, business domain.Business) (string, error) {
			actualizado = true
			return "", nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 404, domain.UpdateBusinessRequest{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "negocio no encontrado")
	assert.False(t, actualizado)
}

func TestUpdateBusiness_ErrorAlBuscar_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) { return nil, dbErr },
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 1, domain.UpdateBusinessRequest{})

	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateBusiness_CamposNilNoModificanNada(t *testing.T) {
	original := negocioExistente()
	var guardado domain.Business
	repo := repoConNegocio()
	repo.UpdateBusinessFn = func(ctx context.Context, id uint, business domain.Business) (string, error) {
		guardado = business
		return "ok", nil
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{})

	require.NoError(t, err)
	assert.Equal(t, original.Name, guardado.Name)
	assert.Equal(t, original.Code, guardado.Code)
	assert.Equal(t, original.BusinessTypeID, guardado.BusinessTypeID)
	assert.Equal(t, original.Timezone, guardado.Timezone)
	assert.Equal(t, original.Address, guardado.Address)
	assert.Equal(t, original.Description, guardado.Description)
	assert.Equal(t, original.LogoURL, guardado.LogoURL)
	assert.Equal(t, original.PrimaryColor, guardado.PrimaryColor)
	assert.Equal(t, original.SidebarColor, guardado.SidebarColor)
	assert.Equal(t, original.TopbarColor, guardado.TopbarColor)
	assert.Equal(t, original.NavbarImageURL, guardado.NavbarImageURL)
	assert.Equal(t, original.CustomDomain, guardado.CustomDomain)
	assert.Equal(t, original.IsActive, guardado.IsActive)
	assert.Equal(t, original.EnableDelivery, guardado.EnableDelivery)
	assert.Equal(t, original.EnablePickup, guardado.EnablePickup)
	assert.Equal(t, original.EnableReservations, guardado.EnableReservations)
}

func TestUpdateBusiness_CadaCampoSeAplicaPorSeparado(t *testing.T) {
	casos := []struct {
		nombre  string
		request domain.UpdateBusinessRequest
		check   func(t *testing.T, b domain.Business)
	}{
		{"nombre", domain.UpdateBusinessRequest{Name: strPtr("Nuevo")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "Nuevo", b.Name) }},
		{"codigo", domain.UpdateBusinessRequest{Code: strPtr("nuevo")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "nuevo", b.Code) }},
		{"tipo", domain.UpdateBusinessRequest{BusinessTypeID: uintPtr(9)},
			func(t *testing.T, b domain.Business) { assert.Equal(t, uint(9), b.BusinessTypeID) }},
		{"timezone", domain.UpdateBusinessRequest{Timezone: strPtr("UTC")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "UTC", b.Timezone) }},
		{"direccion", domain.UpdateBusinessRequest{Address: strPtr("Cra 9")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "Cra 9", b.Address) }},
		{"descripcion", domain.UpdateBusinessRequest{Description: strPtr("nueva")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "nueva", b.Description) }},
		{"primario", domain.UpdateBusinessRequest{PrimaryColor: strPtr("#AAA")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "#AAA", b.PrimaryColor) }},
		{"secundario", domain.UpdateBusinessRequest{SecondaryColor: strPtr("#BBB")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "#BBB", b.SecondaryColor) }},
		{"terciario", domain.UpdateBusinessRequest{TertiaryColor: strPtr("#CCC")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "#CCC", b.TertiaryColor) }},
		{"cuaternario", domain.UpdateBusinessRequest{QuaternaryColor: strPtr("#DDD")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "#DDD", b.QuaternaryColor) }},
		{"sidebar", domain.UpdateBusinessRequest{SidebarColor: strPtr("#EEE")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "#EEE", b.SidebarColor) }},
		{"topbar", domain.UpdateBusinessRequest{TopbarColor: strPtr("#FFF")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "#FFF", b.TopbarColor) }},
		{"dominio", domain.UpdateBusinessRequest{CustomDomain: strPtr("otro.com")},
			func(t *testing.T, b domain.Business) { assert.Equal(t, "otro.com", b.CustomDomain) }},
		{"activo a false", domain.UpdateBusinessRequest{IsActive: boolPtr(false)},
			func(t *testing.T, b domain.Business) { assert.False(t, b.IsActive) }},
		{"delivery a false", domain.UpdateBusinessRequest{EnableDelivery: boolPtr(false)},
			func(t *testing.T, b domain.Business) { assert.False(t, b.EnableDelivery) }},
		{"pickup a true", domain.UpdateBusinessRequest{EnablePickup: boolPtr(true)},
			func(t *testing.T, b domain.Business) { assert.True(t, b.EnablePickup) }},
		{"reservas a false", domain.UpdateBusinessRequest{EnableReservations: boolPtr(false)},
			func(t *testing.T, b domain.Business) { assert.False(t, b.EnableReservations) }},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			var guardado domain.Business
			repo := repoConNegocio()
			repo.UpdateBusinessFn = func(ctx context.Context, id uint, business domain.Business) (string, error) {
				guardado = business
				return "ok", nil
			}

			_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, tc.request)

			require.NoError(t, err)
			tc.check(t, guardado)
		})
	}
}

func TestUpdateBusiness_FalseExplicitoSiApaga(t *testing.T) {
	var guardado domain.Business
	repo := repoConNegocio()
	repo.UpdateBusinessFn = func(ctx context.Context, id uint, business domain.Business) (string, error) {
		guardado = business
		return "ok", nil
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		EnableDelivery:     boolPtr(false),
		EnableReservations: boolPtr(false),
		IsActive:           boolPtr(false),
	})

	require.NoError(t, err)
	assert.False(t, guardado.EnableDelivery, "un false explicito debe apagar, no confundirse con ausente")
	assert.False(t, guardado.EnableReservations)
	assert.False(t, guardado.IsActive)
}

func TestUpdateBusiness_CodigoNuevoDuplicado_Rechaza(t *testing.T) {
	actualizado := false
	repo := repoConNegocio()
	repo.GetBusinessByCodeFn = func(ctx context.Context, code string) (*domain.Business, error) {
		return &domain.Business{ID: 99, Code: code}, nil
	}
	repo.UpdateBusinessFn = func(ctx context.Context, id uint, business domain.Business) (string, error) {
		actualizado = true
		return "", nil
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		Code: strPtr("ocupado"),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "el código 'ocupado' ya existe")
	assert.False(t, actualizado)
}

func TestUpdateBusiness_MismoCodigoNoDisparaChequeo(t *testing.T) {
	chequeado := false
	repo := repoConNegocio()
	repo.GetBusinessByCodeFn = func(ctx context.Context, code string) (*domain.Business, error) {
		chequeado = true
		return nil, nil
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		Code: strPtr("demo"),
	})

	require.NoError(t, err)
	assert.False(t, chequeado, "si el codigo no cambia no hay que buscar duplicados")
}

func TestUpdateBusiness_ErrorRealAlChequearCodigo_Aborta(t *testing.T) {
	dbErr := errors.New("connection refused")
	repo := repoConNegocio()
	repo.GetBusinessByCodeFn = func(ctx context.Context, code string) (*domain.Business, error) {
		return nil, dbErr
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		Code: strPtr("nuevo"),
	})

	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateBusiness_NoEncontradoAlChequearCodigo_NoEsBloqueante(t *testing.T) {
	repo := repoConNegocio()
	repo.GetBusinessByCodeFn = func(ctx context.Context, code string) (*domain.Business, error) {
		return nil, errors.New("negocio no encontrado")
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		Code: strPtr("libre"),
	})

	require.NoError(t, err)
}

func TestUpdateBusiness_DominioNuevoDuplicado_Rechaza(t *testing.T) {
	repo := repoConNegocio()
	repo.GetBusinessByCustomDomainFn = func(ctx context.Context, d string) (*domain.Business, error) {
		return &domain.Business{ID: 99, CustomDomain: d}, nil
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		CustomDomain: strPtr("ocupado.com"),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "el dominio 'ocupado.com' ya existe")
}

func TestUpdateBusiness_MismoDominioNoDisparaChequeo(t *testing.T) {
	chequeado := false
	repo := repoConNegocio()
	repo.GetBusinessByCustomDomainFn = func(ctx context.Context, d string) (*domain.Business, error) {
		chequeado = true
		return nil, nil
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		CustomDomain: strPtr("demo.com"),
	})

	require.NoError(t, err)
	assert.False(t, chequeado)
}

func TestUpdateBusiness_ErrorRealAlChequearDominio_Aborta(t *testing.T) {
	dbErr := errors.New("connection refused")
	repo := repoConNegocio()
	repo.GetBusinessByCustomDomainFn = func(ctx context.Context, d string) (*domain.Business, error) {
		return nil, dbErr
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		CustomDomain: strPtr("nuevo.com"),
	})

	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateBusiness_NoEncontradoAlChequearDominio_NoEsBloqueante(t *testing.T) {
	repo := repoConNegocio()
	repo.GetBusinessByCustomDomainFn = func(ctx context.Context, d string) (*domain.Business, error) {
		return nil, errors.New("negocio no encontrado")
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		CustomDomain: strPtr("libre.com"),
	})

	require.NoError(t, err)
}

func TestUpdateBusiness_NuevoLogo_BorraElAnteriorRelativo(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "businessLogo/nuevo.png", nil
		},
	}
	repo := repoConNegocio()

	_, err := newBusinessUseCase(repo, s3, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		LogoFile: &multipart.FileHeader{Filename: "nuevo.png"},
	})

	require.NoError(t, err)
	assert.Contains(t, s3.DeleteImageCall, "businessLogo/viejo.png")
}

func TestUpdateBusiness_NuevoLogo_NoBorraUrlExterna(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "businessLogo/nuevo.png", nil
		},
	}
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return &domain.Business{ID: id, LogoURL: "https://cdn.externo.com/logo.png"}, nil
		},
	}

	_, err := newBusinessUseCase(repo, s3, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		LogoFile: &multipart.FileHeader{Filename: "nuevo.png"},
	})

	require.NoError(t, err)
	assert.Empty(t, s3.DeleteImageCall, "una url externa no es nuestra, no se borra de S3")
}

func TestUpdateBusiness_MismoPathDeLogo_NoBorra(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "businessLogo/viejo.png", nil
		},
	}
	repo := repoConNegocio()

	_, err := newBusinessUseCase(repo, s3, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		LogoFile: &multipart.FileHeader{Filename: "viejo.png"},
	})

	require.NoError(t, err)
	assert.Empty(t, s3.DeleteImageCall, "si el path resultante es el mismo no hay que borrar nada")
}

func TestUpdateBusiness_FalloAlBorrarLogoAnterior_NoRompeLaActualizacion(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "businessLogo/nuevo.png", nil
		},
		DeleteImageFn: func(ctx context.Context, filename string) error {
			return errors.New("s3 access denied")
		},
	}

	got, err := newBusinessUseCase(repoConNegocio(), s3, nil).UpdateBusiness(context.Background(), 26,
		domain.UpdateBusinessRequest{LogoFile: &multipart.FileHeader{Filename: "nuevo.png"}})

	require.NoError(t, err)
	assert.NotNil(t, got)
}

func TestUpdateBusiness_FallaSubidaDeLogo_NoActualiza(t *testing.T) {
	actualizado := false
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "", errors.New("s3 timeout")
		},
	}
	repo := repoConNegocio()
	repo.UpdateBusinessFn = func(ctx context.Context, id uint, business domain.Business) (string, error) {
		actualizado = true
		return "", nil
	}

	_, err := newBusinessUseCase(repo, s3, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		LogoFile: &multipart.FileHeader{Filename: "nuevo.png"},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error al subir logo")
	assert.False(t, actualizado)
}

func TestUpdateBusiness_NuevaImagenDeNavbar_BorraLaAnteriorRelativa(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "navbar/nuevo.png", nil
		},
	}

	_, err := newBusinessUseCase(repoConNegocio(), s3, nil).UpdateBusiness(context.Background(), 26,
		domain.UpdateBusinessRequest{NavbarImageFile: &multipart.FileHeader{Filename: "nuevo.png"}})

	require.NoError(t, err)
	assert.Contains(t, s3.DeleteImageCall, "navbar/viejo.png")
}

func TestUpdateBusiness_FallaSubidaDeNavbar_NoActualiza(t *testing.T) {
	actualizado := false
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "", errors.New("s3 timeout")
		},
	}
	repo := repoConNegocio()
	repo.UpdateBusinessFn = func(ctx context.Context, id uint, business domain.Business) (string, error) {
		actualizado = true
		return "", nil
	}

	_, err := newBusinessUseCase(repo, s3, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{
		NavbarImageFile: &multipart.FileHeader{Filename: "nuevo.png"},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "imagen de navbar")
	assert.False(t, actualizado)
}

func TestUpdateBusiness_FalloAlBorrarNavbarAnterior_NoRompeLaActualizacion(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "navbar/nuevo.png", nil
		},
		DeleteImageFn: func(ctx context.Context, filename string) error {
			return errors.New("s3 access denied")
		},
	}

	got, err := newBusinessUseCase(repoConNegocio(), s3, nil).UpdateBusiness(context.Background(), 26,
		domain.UpdateBusinessRequest{NavbarImageFile: &multipart.FileHeader{Filename: "nuevo.png"}})

	require.NoError(t, err)
	assert.NotNil(t, got)
}

func TestUpdateBusiness_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := errors.New("constraint violation")
	repo := repoConNegocio()
	repo.UpdateBusinessFn = func(ctx context.Context, id uint, business domain.Business) (string, error) {
		return "", dbErr
	}

	got, err := newBusinessUseCase(repo, nil, nil).UpdateBusiness(context.Background(), 26, domain.UpdateBusinessRequest{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestDeleteBusiness_NegocioInexistente_NoBorra(t *testing.T) {
	borrado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) { return nil, nil },
		DeleteBusinessFn: func(ctx context.Context, id uint) (string, error) {
			borrado = true
			return "", nil
		},
	}

	err := newBusinessUseCase(repo, nil, nil).DeleteBusiness(context.Background(), 404)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "negocio no encontrado")
	assert.False(t, borrado)
}

func TestDeleteBusiness_ErrorAlBuscar_NoBorra(t *testing.T) {
	dbErr := errors.New("timeout")
	borrado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) { return nil, dbErr },
		DeleteBusinessFn: func(ctx context.Context, id uint) (string, error) {
			borrado = true
			return "", nil
		},
	}

	err := newBusinessUseCase(repo, nil, nil).DeleteBusiness(context.Background(), 1)

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, borrado)
}

func TestDeleteBusiness_Exitoso(t *testing.T) {
	var borradoID uint
	repo := repoConNegocio()
	repo.DeleteBusinessFn = func(ctx context.Context, id uint) (string, error) {
		borradoID = id
		return "ok", nil
	}

	err := newBusinessUseCase(repo, nil, nil).DeleteBusiness(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, uint(26), borradoID)
}

func TestDeleteBusiness_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("fk violation")
	repo := repoConNegocio()
	repo.DeleteBusinessFn = func(ctx context.Context, id uint) (string, error) { return "", dbErr }

	err := newBusinessUseCase(repo, nil, nil).DeleteBusiness(context.Background(), 26)

	assert.ErrorIs(t, err, dbErr)
}

func TestToggleBusinessActive_ErrorAlBuscar_NoAplicaElCambio(t *testing.T) {
	aplicado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return nil, errors.New("timeout")
		},
		ToggleBusinessActiveFn: func(ctx context.Context, businessID uint, active bool) error {
			aplicado = true
			return nil
		},
	}

	err := newBusinessUseCase(repo, nil, nil).ToggleBusinessActive(context.Background(), 26, false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "business no encontrado")
	assert.False(t, aplicado)
}

func TestToggleBusinessActive_PasaElEstadoPedido(t *testing.T) {
	for _, activo := range []bool{true, false} {
		var vistoID uint
		var vistoActivo bool
		repo := repoConNegocio()
		repo.ToggleBusinessActiveFn = func(ctx context.Context, businessID uint, active bool) error {
			vistoID, vistoActivo = businessID, active
			return nil
		}

		err := newBusinessUseCase(repo, nil, nil).ToggleBusinessActive(context.Background(), 26, activo)

		require.NoError(t, err)
		assert.Equal(t, uint(26), vistoID)
		assert.Equal(t, activo, vistoActivo)
	}
}

func TestToggleBusinessActive_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("deadlock")
	repo := repoConNegocio()
	repo.ToggleBusinessActiveFn = func(ctx context.Context, businessID uint, active bool) error { return dbErr }

	err := newBusinessUseCase(repo, nil, nil).ToggleBusinessActive(context.Background(), 26, true)

	assert.ErrorIs(t, err, dbErr)
}

func TestToggleBusinessResourceActive_RecursoNoPermitidoParaElTipo_Rechaza(t *testing.T) {
	aplicado := false
	repo := repoConNegocio()
	repo.GetResourceByIDFn = func(ctx context.Context, resourceID uint) (*domain.Resource, error) {
		return &domain.Resource{ID: resourceID, Name: "shipments"}, nil
	}
	repo.GetBusinessTypeResourcesPermittedFn = func(ctx context.Context, businessTypeID uint) ([]domain.BusinessTypeResourcePermitted, error) {
		return []domain.BusinessTypeResourcePermitted{{ResourceID: 1}, {ResourceID: 2}}, nil
	}
	repo.ToggleBusinessResourceActiveFn = func(ctx context.Context, businessID uint, resourceID uint, active bool) error {
		aplicado = true
		return nil
	}

	err := newBusinessUseCase(repo, nil, nil).ToggleBusinessResourceActive(context.Background(), 26, 99, true)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no está permitido para este tipo de business")
	assert.False(t, aplicado, "activar un recurso fuera del tipo de negocio abriria funcionalidad no contratada")
}

func TestToggleBusinessResourceActive_RecursoPermitido_Aplica(t *testing.T) {
	var vistoBusiness, vistoRecurso uint
	var vistoActivo bool
	repo := repoConNegocio()
	repo.GetBusinessTypeResourcesPermittedFn = func(ctx context.Context, businessTypeID uint) ([]domain.BusinessTypeResourcePermitted, error) {
		return []domain.BusinessTypeResourcePermitted{{ResourceID: 5}}, nil
	}
	repo.ToggleBusinessResourceActiveFn = func(ctx context.Context, businessID uint, resourceID uint, active bool) error {
		vistoBusiness, vistoRecurso, vistoActivo = businessID, resourceID, active
		return nil
	}

	err := newBusinessUseCase(repo, nil, nil).ToggleBusinessResourceActive(context.Background(), 26, 5, true)

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoBusiness)
	assert.Equal(t, uint(5), vistoRecurso)
	assert.True(t, vistoActivo)
}

func TestToggleBusinessResourceActive_ConsultaLosPermitidosDelTipoDelNegocio(t *testing.T) {
	var vistoTipo uint
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return &domain.Business{ID: id, BusinessTypeID: 7}, nil
		},
		GetBusinessTypeResourcesPermittedFn: func(ctx context.Context, businessTypeID uint) ([]domain.BusinessTypeResourcePermitted, error) {
			vistoTipo = businessTypeID
			return []domain.BusinessTypeResourcePermitted{{ResourceID: 5}}, nil
		},
	}

	err := newBusinessUseCase(repo, nil, nil).ToggleBusinessResourceActive(context.Background(), 26, 5, true)

	require.NoError(t, err)
	assert.Equal(t, uint(7), vistoTipo, "los recursos permitidos salen del tipo del negocio, no de otro")
}

func TestToggleBusinessResourceActive_ErrorAlBuscarNegocio_Rechaza(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return nil, errors.New("timeout")
		},
	}

	err := newBusinessUseCase(repo, nil, nil).ToggleBusinessResourceActive(context.Background(), 26, 5, true)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "business no encontrado")
}

func TestToggleBusinessResourceActive_ErrorAlBuscarRecurso_Rechaza(t *testing.T) {
	repo := repoConNegocio()
	repo.GetResourceByIDFn = func(ctx context.Context, resourceID uint) (*domain.Resource, error) {
		return nil, errors.New("no existe")
	}

	err := newBusinessUseCase(repo, nil, nil).ToggleBusinessResourceActive(context.Background(), 26, 5, true)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "recurso no encontrado")
}

func TestToggleBusinessResourceActive_ErrorAlLeerPermitidos_Rechaza(t *testing.T) {
	aplicado := false
	repo := repoConNegocio()
	repo.GetBusinessTypeResourcesPermittedFn = func(ctx context.Context, businessTypeID uint) ([]domain.BusinessTypeResourcePermitted, error) {
		return nil, errors.New("timeout")
	}
	repo.ToggleBusinessResourceActiveFn = func(ctx context.Context, businessID uint, resourceID uint, active bool) error {
		aplicado = true
		return nil
	}

	err := newBusinessUseCase(repo, nil, nil).ToggleBusinessResourceActive(context.Background(), 26, 5, true)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error al verificar recursos permitidos")
	assert.False(t, aplicado, "si no podemos verificar el permiso, no se aplica")
}

func TestToggleBusinessResourceActive_ErrorAlAplicar_SePropaga(t *testing.T) {
	dbErr := errors.New("deadlock")
	repo := repoConNegocio()
	repo.GetBusinessTypeResourcesPermittedFn = func(ctx context.Context, businessTypeID uint) ([]domain.BusinessTypeResourcePermitted, error) {
		return []domain.BusinessTypeResourcePermitted{{ResourceID: 5}}, nil
	}
	repo.ToggleBusinessResourceActiveFn = func(ctx context.Context, businessID uint, resourceID uint, active bool) error {
		return dbErr
	}

	err := newBusinessUseCase(repo, nil, nil).ToggleBusinessResourceActive(context.Background(), 26, 5, true)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetBusinessesConfiguredResources_PasaFiltrosYPropagaError(t *testing.T) {
	var vistoPage, vistoPerPage int
	var vistoBusiness, vistoTipo *uint
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessesWithConfiguredResourcesPaginatedFn: func(ctx context.Context, page, perPage int, businessID *uint, businessTypeID *uint) ([]domain.BusinessWithConfiguredResourcesResponse, int64, error) {
			vistoPage, vistoPerPage = page, perPage
			vistoBusiness, vistoTipo = businessID, businessTypeID
			return []domain.BusinessWithConfiguredResourcesResponse{{ID: 26}}, 1, nil
		},
	}

	got, total, err := newBusinessUseCase(repo, nil, nil).GetBusinessesConfiguredResources(
		context.Background(), 2, 50, uintPtr(26), uintPtr(1))

	require.NoError(t, err)
	assert.Equal(t, 2, vistoPage)
	assert.Equal(t, 50, vistoPerPage)
	require.NotNil(t, vistoBusiness)
	assert.Equal(t, uint(26), *vistoBusiness)
	require.NotNil(t, vistoTipo)
	assert.Equal(t, uint(1), *vistoTipo)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
}

func TestGetBusinessesConfiguredResources_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessesWithConfiguredResourcesPaginatedFn: func(ctx context.Context, page, perPage int, businessID *uint, businessTypeID *uint) ([]domain.BusinessWithConfiguredResourcesResponse, int64, error) {
			return nil, 0, dbErr
		},
	}

	got, total, err := newBusinessUseCase(repo, nil, nil).GetBusinessesConfiguredResources(context.Background(), 1, 10, nil, nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
	assert.Zero(t, total)
}

func TestGetBusinessConfiguredResourcesByID_DelegaYPropaga(t *testing.T) {
	var visto uint
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDWithConfiguredResourcesFn: func(ctx context.Context, businessID uint) (*domain.BusinessWithConfiguredResourcesResponse, error) {
			visto = businessID
			return &domain.BusinessWithConfiguredResourcesResponse{ID: businessID, Name: "Demo"}, nil
		},
	}

	got, err := newBusinessUseCase(repo, nil, nil).GetBusinessConfiguredResourcesByID(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, uint(26), visto)
	assert.Equal(t, "Demo", got.Name)
}

func TestGetBusinessConfiguredResourcesByID_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("no encontrado")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDWithConfiguredResourcesFn: func(ctx context.Context, businessID uint) (*domain.BusinessWithConfiguredResourcesResponse, error) {
			return nil, dbErr
		},
	}

	got, err := newBusinessUseCase(repo, nil, nil).GetBusinessConfiguredResourcesByID(context.Background(), 26)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetFiscalProfile_SinPerfil_DevuelveDefaultsColombianos(t *testing.T) {
	repo := repoConNegocio()
	repo.GetFiscalProfileFn = func(ctx context.Context, businessID uint) (*domain.FiscalProfile, error) {
		return nil, nil
	}

	got, err := newBusinessUseCase(repo, nil, nil).GetFiscalProfile(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, uint(26), got.BusinessID)
	assert.Equal(t, "NIT", got.DocumentType)
	assert.Equal(t, "JURIDICA", got.PersonType)
	assert.InDelta(t, 19.0, got.IVARate, 0.001)
	assert.InDelta(t, 11.0, got.ReteFuenteRate, 0.001)
	assert.InDelta(t, 0.966, got.ReteICARate, 0.0001)
}

func TestGetFiscalProfile_ConPerfil_LoDevuelveTalCual(t *testing.T) {
	guardado := &domain.FiscalProfile{
		ID: 1, BusinessID: 26, DocumentType: "CC", PersonType: "NATURAL", IVARate: 5,
	}
	repo := repoConNegocio()
	repo.GetFiscalProfileFn = func(ctx context.Context, businessID uint) (*domain.FiscalProfile, error) {
		return guardado, nil
	}

	got, err := newBusinessUseCase(repo, nil, nil).GetFiscalProfile(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, guardado, got)
}

func TestGetFiscalProfile_NegocioInexistente_SePropaga(t *testing.T) {
	dbErr := errors.New("no encontrado")
	consultado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) { return nil, dbErr },
		GetFiscalProfileFn: func(ctx context.Context, businessID uint) (*domain.FiscalProfile, error) {
			consultado = true
			return nil, nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).GetFiscalProfile(context.Background(), 404)

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, consultado, "no se consulta el perfil fiscal de un negocio que no podemos verificar")
}

func TestGetFiscalProfile_ErrorAlLeerPerfil_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := repoConNegocio()
	repo.GetFiscalProfileFn = func(ctx context.Context, businessID uint) (*domain.FiscalProfile, error) {
		return nil, dbErr
	}

	_, err := newBusinessUseCase(repo, nil, nil).GetFiscalProfile(context.Background(), 26)

	assert.ErrorIs(t, err, dbErr)
}

func TestUpsertFiscalProfile_NormalizaDocumentTypeAMayusculas(t *testing.T) {
	var guardado *domain.FiscalProfile
	repo := repoConNegocio()
	repo.UpsertFiscalProfileFn = func(ctx context.Context, profile *domain.FiscalProfile) error {
		guardado = profile
		return nil
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpsertFiscalProfile(context.Background(), &domain.FiscalProfile{
		BusinessID: 26, DocumentType: "  cc  ", PersonType: "natural",
	})

	require.NoError(t, err)
	assert.Equal(t, "CC", guardado.DocumentType)
	assert.Equal(t, "NATURAL", guardado.PersonType)
}

func TestUpsertFiscalProfile_DocumentTypeVacio_CaeANIT(t *testing.T) {
	var guardado *domain.FiscalProfile
	repo := repoConNegocio()
	repo.UpsertFiscalProfileFn = func(ctx context.Context, profile *domain.FiscalProfile) error {
		guardado = profile
		return nil
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpsertFiscalProfile(context.Background(), &domain.FiscalProfile{
		BusinessID: 26, DocumentType: "   ", PersonType: "JURIDICA",
	})

	require.NoError(t, err)
	assert.Equal(t, "NIT", guardado.DocumentType)
}

func TestUpsertFiscalProfile_PersonTypeInvalido_Rechaza(t *testing.T) {
	casos := []string{"", "EMPRESA", "juridico", "MIXTA", "N/A"}

	for _, entrada := range casos {
		t.Run(entrada, func(t *testing.T) {
			guardado := false
			repo := repoConNegocio()
			repo.UpsertFiscalProfileFn = func(ctx context.Context, profile *domain.FiscalProfile) error {
				guardado = true
				return nil
			}

			_, err := newBusinessUseCase(repo, nil, nil).UpsertFiscalProfile(context.Background(), &domain.FiscalProfile{
				BusinessID: 26, PersonType: entrada,
			})

			require.Error(t, err)
			assert.Contains(t, err.Error(), "JURIDICA o NATURAL")
			assert.False(t, guardado)
		})
	}
}

func TestUpsertFiscalProfile_TarifasFueraDeRango_Rechaza(t *testing.T) {
	casos := []struct {
		nombre string
		perfil domain.FiscalProfile
	}{
		{"IVA negativo", domain.FiscalProfile{PersonType: "JURIDICA", IVARate: -1}},
		{"IVA sobre cien", domain.FiscalProfile{PersonType: "JURIDICA", IVARate: 101}},
		{"ReteFuente negativo", domain.FiscalProfile{PersonType: "JURIDICA", ReteFuenteRate: -0.5}},
		{"ReteFuente sobre cien", domain.FiscalProfile{PersonType: "JURIDICA", ReteFuenteRate: 150}},
		{"ReteICA negativo", domain.FiscalProfile{PersonType: "JURIDICA", ReteICARate: -2}},
		{"ReteICA sobre cien", domain.FiscalProfile{PersonType: "JURIDICA", ReteICARate: 100.1}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			guardado := false
			repo := repoConNegocio()
			repo.UpsertFiscalProfileFn = func(ctx context.Context, profile *domain.FiscalProfile) error {
				guardado = true
				return nil
			}
			perfil := tc.perfil
			perfil.BusinessID = 26

			_, err := newBusinessUseCase(repo, nil, nil).UpsertFiscalProfile(context.Background(), &perfil)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "entre 0 y 100")
			assert.False(t, guardado)
		})
	}
}

func TestUpsertFiscalProfile_LosLimitesExactosSeAceptan(t *testing.T) {
	casos := []domain.FiscalProfile{
		{PersonType: "JURIDICA", IVARate: 0, ReteFuenteRate: 0, ReteICARate: 0},
		{PersonType: "JURIDICA", IVARate: 100, ReteFuenteRate: 100, ReteICARate: 100},
		{PersonType: "NATURAL", IVARate: 19, ReteFuenteRate: 11, ReteICARate: 0.966},
	}

	for i, perfil := range casos {
		repo := repoConNegocio()
		perfil.BusinessID = 26

		_, err := newBusinessUseCase(repo, nil, nil).UpsertFiscalProfile(context.Background(), &perfil)

		require.NoError(t, err, "caso %d", i)
	}
}

func TestUpsertFiscalProfile_NegocioInexistente_NoGuarda(t *testing.T) {
	dbErr := errors.New("no encontrado")
	guardado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) { return nil, dbErr },
		UpsertFiscalProfileFn: func(ctx context.Context, profile *domain.FiscalProfile) error {
			guardado = true
			return nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).UpsertFiscalProfile(context.Background(), &domain.FiscalProfile{
		BusinessID: 404, PersonType: "JURIDICA",
	})

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, guardado)
}

func TestUpsertFiscalProfile_ErrorAlGuardar_SePropaga(t *testing.T) {
	dbErr := errors.New("constraint violation")
	repo := repoConNegocio()
	repo.UpsertFiscalProfileFn = func(ctx context.Context, profile *domain.FiscalProfile) error { return dbErr }

	got, err := newBusinessUseCase(repo, nil, nil).UpsertFiscalProfile(context.Background(), &domain.FiscalProfile{
		BusinessID: 26, PersonType: "JURIDICA",
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpsertFiscalProfile_Exitoso_ReleeElPerfilGuardado(t *testing.T) {
	releido := &domain.FiscalProfile{ID: 9, BusinessID: 26, DocumentType: "NIT", PersonType: "JURIDICA"}
	repo := repoConNegocio()
	repo.GetFiscalProfileFn = func(ctx context.Context, businessID uint) (*domain.FiscalProfile, error) {
		return releido, nil
	}

	got, err := newBusinessUseCase(repo, nil, nil).UpsertFiscalProfile(context.Background(), &domain.FiscalProfile{
		BusinessID: 26, PersonType: "JURIDICA",
	})

	require.NoError(t, err)
	assert.Equal(t, releido, got, "la respuesta viene de la relectura, no del input")
}
