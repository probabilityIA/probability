package usecasebusiness

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const s3Base = "https://cdn.probability.co"

func s3Config() map[string]string {
	return map[string]string{"URL_BASE_DOMAIN_S3": s3Base}
}

func newBusinessUseCase(repo *mocks.BusinessRepositoryMock, s3 *mocks.S3ServiceMock, cfg map[string]string) IUseCaseBusiness {
	if repo == nil {
		repo = &mocks.BusinessRepositoryMock{}
	}
	if s3 == nil {
		s3 = &mocks.S3ServiceMock{}
	}
	return New(repo, mocks.NewSilentLogger(), s3, mocks.NewConfigMock(cfg))
}

func rawUseCase(repo *mocks.BusinessRepositoryMock) *BusinessUseCase {
	if repo == nil {
		repo = &mocks.BusinessRepositoryMock{}
	}
	return &BusinessUseCase{
		repository: repo,
		log:        mocks.NewSilentLogger(),
		s3:         &mocks.S3ServiceMock{},
		env:        mocks.NewConfigMock(nil),
	}
}


func TestGenerateCodeFromName_NormalizaAMinusculasYSustituyeSeparadores(t *testing.T) {
	casos := []struct {
		entrada    string
		wantPrefix string
	}{
		{"Mi Negocio", "mi_negocio_"},
		{"MI NEGOCIO", "mi_negocio_"},
		{"Mi-Negocio", "mi_negocio_"},
		{"Mi_Negocio", "mi_negocio_"},
		{"Tienda 24/7", "tienda_247_"},
		{"Cafe & Bar", "cafe__bar_"},
		{"Tienda#1", "tienda1_"},
	}

	for _, tc := range casos {
		t.Run(tc.entrada, func(t *testing.T) {
			got := generateCodeFromName(tc.entrada)
			assert.True(t, strings.HasPrefix(got, tc.wantPrefix),
				"esperaba prefijo %q, obtuve %q", tc.wantPrefix, got)
		})
	}
}

func TestGenerateCodeFromName_SoloEspacioGuionYUnderscoreSonSeparadores(t *testing.T) {
	casos := []struct {
		entrada    string
		wantPrefix string
	}{
		{"a b", "a_b_"},
		{"a-b", "a_b_"},
		{"a_b", "a_b_"},
		{"a/b", "ab_"},
		{"a.b", "ab_"},
		{"a&b", "ab_"},
		{"a@b", "ab_"},
	}

	for _, tc := range casos {
		t.Run(tc.entrada, func(t *testing.T) {
			got := generateCodeFromName(tc.entrada)
			assert.True(t, strings.HasPrefix(got, tc.wantPrefix),
				"esperaba prefijo %q, obtuve %q", tc.wantPrefix, got)
		})
	}
}

func TestGenerateCodeFromName_SeparadoresConsecutivosNoSeColapsan(t *testing.T) {
	got := generateCodeFromName("Cafe   Bar")

	assert.True(t, strings.HasPrefix(got, "cafe___bar_"),
		"tres espacios producen tres underscores, no se colapsan: %q", got)
}

func TestGenerateCodeFromName_TruncaLaBaseAVeinteCaracteres(t *testing.T) {
	largo := "NegocioConUnNombreExtremadamenteLargoQueNoCabe"

	got := generateCodeFromName(largo)

	base := got[:strings.LastIndex(got, "_")]
	assert.LessOrEqual(t, len(base), 20, "la base del codigo se trunca a 20: %q", got)
}

func TestGenerateCodeFromName_NoRepiteCodigosParaElMismoNombre(t *testing.T) {
	vistos := make(map[string]bool, 200)
	for i := 0; i < 200; i++ {
		got := generateCodeFromName("Mi Negocio")
		require.False(t, vistos[got], "codigo repetido en la iteracion %d: %q", i, got)
		vistos[got] = true
	}
}

func TestGenerateCodeFromName_NoDejaGuionesNiSlashDelBase64(t *testing.T) {
	for i := 0; i < 200; i++ {
		got := generateCodeFromName("Mi Negocio")
		sufijo := got[strings.LastIndex(got, "_")+1:]
		assert.NotContains(t, sufijo, "-", "el sufijo aleatorio no debe traer guiones: %q", got)
	}
}

func TestGenerateCodeFromName_NombreSinAlfanumericos_ProduceCodigoUsable(t *testing.T) {
	got := generateCodeFromName("!!!")

	assert.True(t, strings.HasPrefix(got, "_"), "queda solo el separador y el sufijo: %q", got)
	assert.Greater(t, len(got), 1)
}

func TestResolveOrderPrefix_TomaLasTresPrimerasLetras(t *testing.T) {
	casos := []struct {
		entrada string
		want    string
	}{
		{"Probability", "PRO"},
		{"mi negocio", "MIN"},
		{"Cafe Bar", "CAF"},
		{"A B C D", "ABC"},
		{"123 Tienda", "TIE"},
	}

	uc := rawUseCase(nil)
	for _, tc := range casos {
		t.Run(tc.entrada, func(t *testing.T) {
			got := uc.resolveOrderPrefix(context.Background(), tc.entrada)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestResolveOrderPrefix_RellenaConXCuandoFaltanLetras(t *testing.T) {
	uc := rawUseCase(nil)

	assert.Equal(t, "AXX", uc.resolveOrderPrefix(context.Background(), "A"))
	assert.Equal(t, "ABX", uc.resolveOrderPrefix(context.Background(), "AB"))
}

func TestResolveOrderPrefix_SinLetras_CaeAlValorPorDefecto(t *testing.T) {
	uc := rawUseCase(nil)

	for _, entrada := range []string{"", "123", "!!!", "   "} {
		assert.Equal(t, "BIZ", uc.resolveOrderPrefix(context.Background(), entrada),
			"entrada %q", entrada)
	}
}

func TestResolveOrderPrefix_EvitaColisionesConLosYaUsados(t *testing.T) {
	casos := []struct {
		nombre   string
		ocupados []string
		want     string
	}{
		{"libre", []string{}, "PRO"},
		{"uno ocupado", []string{"PRO"}, "PRO2"},
		{"dos ocupados", []string{"PRO", "PRO2"}, "PRO3"},
		{"tres ocupados", []string{"PRO", "PRO2", "PRO3"}, "PRO4"},
		{"hueco intermedio no se reusa", []string{"PRO", "PRO3"}, "PRO2"},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := &mocks.BusinessRepositoryMock{
				GetExistingOrderPrefixesFn: func(ctx context.Context) ([]string, error) {
					return tc.ocupados, nil
				},
			}

			got := rawUseCase(repo).resolveOrderPrefix(context.Background(), "Probability")

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestResolveOrderPrefix_ComparaOcupadosSinImportarMayusculas(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetExistingOrderPrefixesFn: func(ctx context.Context) ([]string, error) {
			return []string{"pro", "Pro2"}, nil
		},
	}

	got := rawUseCase(repo).resolveOrderPrefix(context.Background(), "Probability")

	assert.Equal(t, "PRO3", got, "los prefijos ocupados se comparan en mayusculas")
}

func TestResolveOrderPrefix_ErrorAlLeerOcupados_NoBloquea(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetExistingOrderPrefixesFn: func(ctx context.Context) ([]string, error) {
			return nil, errors.New("timeout")
		},
	}

	got := rawUseCase(repo).resolveOrderPrefix(context.Background(), "Probability")

	assert.Equal(t, "PRO", got)
}

func TestCreateBusiness_SinTipoDeNegocio_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = true
			return 1, nil
		},
	}

	got, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "Mi Negocio",
	})

	assert.ErrorIs(t, err, domain.ErrBusinessTypeIDRequired)
	assert.Nil(t, got)
	assert.False(t, creado)
}

func TestCreateBusiness_TipoDeNegocioInexistente_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) {
			return nil, nil
		},
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = true
			return 1, nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 99,
	})

	assert.ErrorIs(t, err, domain.ErrBusinessTypeIDInvalid)
	assert.False(t, creado)
}

func TestCreateBusiness_ErrorAlVerificarTipo_Rechaza(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) {
			return nil, errors.New("timeout")
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
	})

	assert.ErrorIs(t, err, domain.ErrBusinessTypeIDInvalid)
}

func TestCreateBusiness_CodigoVacio_SeGeneraDesdeElNombre(t *testing.T) {
	var creado domain.Business
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = business
			return 5, nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "Mi Negocio", BusinessTypeID: 1,
	})

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(creado.Code, "mi_negocio_"),
		"sin codigo explicito se genera uno: %q", creado.Code)
}

func TestCreateBusiness_CodigoExplicito_SeRespeta(t *testing.T) {
	var creado domain.Business
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = business
			return 5, nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "Mi Negocio", Code: "MINEGOCIO", BusinessTypeID: 1,
	})

	require.NoError(t, err)
	assert.Equal(t, "MINEGOCIO", creado.Code)
}

func TestCreateBusiness_CodigoYaEnUso_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByCodeFn: func(ctx context.Context, code string) (*domain.Business, error) {
			return &domain.Business{ID: 1, Code: code}, nil
		},
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = true
			return 1, nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", Code: "DUP", BusinessTypeID: 1,
	})

	assert.ErrorIs(t, err, domain.ErrBusinessCodeAlreadyExists)
	assert.False(t, creado)
}

func TestCreateBusiness_ErrBusinessNotFoundAlBuscarCodigo_NoEsBloqueante(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByCodeFn: func(ctx context.Context, code string) (*domain.Business, error) {
			return nil, domain.ErrBusinessNotFound
		},
	}

	got, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", Code: "LIBRE", BusinessTypeID: 1,
	})

	require.NoError(t, err, "que el codigo no exista es justamente lo que queremos")
	assert.NotNil(t, got)
}

func TestCreateBusiness_ErrorRealAlBuscarCodigo_Aborta(t *testing.T) {
	dbErr := errors.New("connection refused")
	creado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByCodeFn: func(ctx context.Context, code string) (*domain.Business, error) {
			return nil, dbErr
		},
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = true
			return 1, nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
	})

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, creado)
}

func TestCreateBusiness_DominioPersonalizadoEnUso_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByCustomDomainFn: func(ctx context.Context, d string) (*domain.Business, error) {
			return &domain.Business{ID: 2, CustomDomain: d}, nil
		},
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = true
			return 1, nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1, CustomDomain: "tienda.com",
	})

	assert.ErrorIs(t, err, domain.ErrBusinessDomainAlreadyExists)
	assert.False(t, creado)
}

func TestCreateBusiness_SinDominioPersonalizado_NoLoConsulta(t *testing.T) {
	consultado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByCustomDomainFn: func(ctx context.Context, d string) (*domain.Business, error) {
			consultado = true
			return nil, nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
	})

	require.NoError(t, err)
	assert.False(t, consultado)
}

func TestCreateBusiness_ErrorRealAlBuscarDominio_Aborta(t *testing.T) {
	dbErr := errors.New("connection refused")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByCustomDomainFn: func(ctx context.Context, d string) (*domain.Business, error) {
			return nil, dbErr
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1, CustomDomain: "tienda.com",
	})

	assert.ErrorIs(t, err, dbErr)
}

func TestCreateBusiness_ErrBusinessNotFoundAlBuscarDominio_NoEsBloqueante(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByCustomDomainFn: func(ctx context.Context, d string) (*domain.Business, error) {
			return nil, domain.ErrBusinessNotFound
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1, CustomDomain: "libre.com",
	})

	require.NoError(t, err)
}

func TestCreateBusiness_SubeLogoYNavbarACarpetasDistintas(t *testing.T) {
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return folder + "/" + file.Filename, nil
		},
	}
	var creado domain.Business
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = business
			return 5, nil
		},
	}

	_, err := newBusinessUseCase(repo, s3, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
		LogoFile:        &multipart.FileHeader{Filename: "logo.png"},
		NavbarImageFile: &multipart.FileHeader{Filename: "navbar.png"},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"businessLogo", "navbar"}, s3.UploadImageCall)
	assert.Equal(t, "businessLogo/logo.png", creado.LogoURL)
	assert.Equal(t, "navbar/navbar.png", creado.NavbarImageURL)
}

func TestCreateBusiness_FallaSubidaDeLogo_NoCrea(t *testing.T) {
	creado := false
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			return "", errors.New("s3 timeout")
		},
	}
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = true
			return 1, nil
		},
	}

	_, err := newBusinessUseCase(repo, s3, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
		LogoFile: &multipart.FileHeader{Filename: "logo.png"},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error al subir el logo")
	assert.False(t, creado)
}

func TestCreateBusiness_FallaSubidaDeNavbar_NoCrea(t *testing.T) {
	creado := false
	s3 := &mocks.S3ServiceMock{
		UploadImageFn: func(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
			if folder == "navbar" {
				return "", errors.New("s3 timeout")
			}
			return "businessLogo/logo.png", nil
		},
	}
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = true
			return 1, nil
		},
	}

	_, err := newBusinessUseCase(repo, s3, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
		LogoFile:        &multipart.FileHeader{Filename: "logo.png"},
		NavbarImageFile: &multipart.FileHeader{Filename: "navbar.png"},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "imagen de navbar")
	assert.False(t, creado)
}

func TestCreateBusiness_AsignaOrderPrefixDerivadoDelNombre(t *testing.T) {
	var creado domain.Business
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = business
			return 5, nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "Probability", BusinessTypeID: 1,
	})

	require.NoError(t, err)
	assert.Equal(t, "PRO", creado.OrderPrefix)
}

func TestCreateBusiness_FalloAlCrearIntegracionDePlataforma_NoAbortaLaCreacion(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) { return 5, nil },
		CreatePlatformIntegrationFn: func(ctx context.Context, businessID uint) error {
			return errors.New("cola caida")
		},
	}

	got, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
	})

	require.NoError(t, err, "la integracion de plataforma es best-effort, no bloquea")
	require.NotNil(t, got)
	assert.Equal(t, uint(5), got.ID)
}

func TestCreateBusiness_CreaLaIntegracionDePlataformaConElIDNuevo(t *testing.T) {
	var integracionPara uint
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) { return 77, nil },
		CreatePlatformIntegrationFn: func(ctx context.Context, businessID uint) error {
			integracionPara = businessID
			return nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(77), integracionPara)
}

func TestCreateBusiness_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := errors.New("constraint violation")
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) { return 0, dbErr },
	}

	got, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateBusiness_ErrorAlRecargarElCreado_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) { return nil, dbErr },
	}

	got, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateBusiness_CompletaUrlsRelativasConElDominioS3(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return &domain.Business{
				ID: id, LogoURL: "businessLogo/logo.png", NavbarImageURL: "navbar/n.png",
			}, nil
		},
	}

	got, err := newBusinessUseCase(repo, nil, s3Config()).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
	})

	require.NoError(t, err)
	assert.Equal(t, s3Base+"/businessLogo/logo.png", got.LogoURL)
	assert.Equal(t, s3Base+"/navbar/n.png", got.NavbarImageURL)
}

func TestCreateBusiness_UrlsAbsolutasNoSeTocan(t *testing.T) {
	externa := "https://cdn.externo.com/logo.png"
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return &domain.Business{ID: id, LogoURL: externa, NavbarImageURL: externa}, nil
		},
	}

	got, err := newBusinessUseCase(repo, nil, s3Config()).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
	})

	require.NoError(t, err)
	assert.Equal(t, externa, got.LogoURL)
	assert.Equal(t, externa, got.NavbarImageURL)
}

func TestCreateBusiness_SinDominioS3_DejaElPathRelativo(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return &domain.Business{ID: id, LogoURL: "businessLogo/logo.png"}, nil
		},
	}

	got, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
	})

	require.NoError(t, err)
	assert.Equal(t, "businessLogo/logo.png", got.LogoURL)
}

func TestCreateBusiness_SinBusinessTypeCargado_DevuelveSoloElID(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return &domain.Business{ID: id, BusinessTypeID: 4, BusinessType: nil}, nil
		},
	}

	got, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 4,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(4), got.BusinessType.ID)
	assert.Empty(t, got.BusinessType.Name)
}

func TestCreateBusiness_ConBusinessTypeCargado_MapeaTodo(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessByIDFn: func(ctx context.Context, id uint) (*domain.Business, error) {
			return &domain.Business{
				ID: id, BusinessTypeID: 4,
				BusinessType: &domain.BusinessType{
					ID: 4, Name: "Retail", Code: "retail", Description: "d", Icon: "i", IsActive: true,
				},
			}, nil
		},
	}

	got, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 4,
	})

	require.NoError(t, err)
	assert.Equal(t, "Retail", got.BusinessType.Name)
	assert.Equal(t, "retail", got.BusinessType.Code)
	assert.True(t, got.BusinessType.IsActive)
}

func TestCreateBusiness_PropagaLosFlagsDeFuncionalidad(t *testing.T) {
	var creado domain.Business
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessFn: func(ctx context.Context, business domain.Business) (uint, error) {
			creado = business
			return 5, nil
		},
	}

	_, err := newBusinessUseCase(repo, nil, nil).CreateBusiness(context.Background(), domain.BusinessRequest{
		Name: "X", BusinessTypeID: 1,
		Timezone: "America/Bogota", Address: "Cra 1", Description: "d",
		PrimaryColor: "#111", SecondaryColor: "#222", TertiaryColor: "#333",
		QuaternaryColor: "#444", SidebarColor: "#555", TopbarColor: "#666",
		IsActive: true, EnableDelivery: true, EnablePickup: true, EnableReservations: true,
	})

	require.NoError(t, err)
	assert.Equal(t, "America/Bogota", creado.Timezone)
	assert.Equal(t, "#555", creado.SidebarColor)
	assert.Equal(t, "#666", creado.TopbarColor)
	assert.True(t, creado.EnableDelivery)
	assert.True(t, creado.EnablePickup)
	assert.True(t, creado.EnableReservations)
	assert.True(t, creado.IsActive)
}
