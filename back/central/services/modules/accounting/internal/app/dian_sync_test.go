package app

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func credencialesValidas() dtos.DianConfigDTO {
	return dtos.DianConfigDTO{
		ClientID: "cid", ClientSecret: "csecret", Username: "user", Password: "pass",
	}
}

func TestConfigString_ToleraNilYTiposDistintos(t *testing.T) {
	assert.Empty(t, configString(nil, "x"))
	assert.Empty(t, configString(map[string]interface{}{}, "x"))
	assert.Empty(t, configString(map[string]interface{}{"x": 123}, "x"))
	assert.Equal(t, "v", configString(map[string]interface{}{"x": "v"}, "x"))
}

func TestConfigInt_AceptaFloatYIntYToleraElResto(t *testing.T) {
	assert.Zero(t, configInt(nil, "x"))
	assert.Zero(t, configInt(map[string]interface{}{}, "x"))
	assert.Zero(t, configInt(map[string]interface{}{"x": "8"}, "x"))
	assert.Equal(t, 8, configInt(map[string]interface{}{"x": 8}, "x"))
	assert.Equal(t, 8, configInt(map[string]interface{}{"x": float64(8)}, "x"),
		"el JSON deserializado trae numeros como float64")
}

func TestGetDianConfig_SinIntegracion_ReportaNoConfigurado(t *testing.T) {
	integration := &mocks.DianIntegrationMock{
		GetGlobalIntegrationIDFn: func(ctx context.Context) (uint, error) {
			return 0, stderrors.New("no existe")
		},
	}

	got, err := newAccountingUseCase(nil, nil, nil, integration).GetDianConfig(context.Background())

	require.NoError(t, err, "no tener facturacion electronica no es un error, es un estado")
	assert.False(t, got.Configured)
	assert.Zero(t, got.IntegrationID)
}

func TestGetDianConfig_ErrorAlLeerConfig_ReportaConfiguradoSinDetalle(t *testing.T) {
	integration := &mocks.DianIntegrationMock{
		GetConfigFn: func(ctx context.Context, integrationID uint) (map[string]interface{}, error) {
			return nil, stderrors.New("timeout")
		},
	}

	got, err := newAccountingUseCase(nil, nil, nil, integration).GetDianConfig(context.Background())

	require.NoError(t, err)
	assert.True(t, got.Configured)
	assert.Empty(t, got.MunicipalityID)
}

func TestGetDianConfig_MapeaLaConfiguracion(t *testing.T) {
	integration := &mocks.DianIntegrationMock{
		GetGlobalIntegrationIDFn: func(ctx context.Context) (uint, error) { return 7, nil },
		GetConfigFn: func(ctx context.Context, integrationID uint) (map[string]interface{}, error) {
			return map[string]interface{}{
				"numbering_range_id":         float64(8),
				"municipality_id":            "11001",
				"identification_document_id": "3",
				"legal_organization_id":      "1",
				"api_url_display":            "https://api.factus.com",
				"username_display":           "user@x.com",
				"payment_method_code":        "10",
			}, nil
		},
	}

	got, err := newAccountingUseCase(nil, nil, nil, integration).GetDianConfig(context.Background())

	require.NoError(t, err)
	assert.True(t, got.Configured)
	assert.Equal(t, uint(7), got.IntegrationID)
	assert.Equal(t, 8, got.NumberingRangeID)
	assert.Equal(t, "11001", got.MunicipalityID)
	assert.Equal(t, "3", got.DocumentIDType)
	assert.Equal(t, "https://api.factus.com", got.ApiURL)
	assert.Equal(t, "user@x.com", got.Username)
	assert.Equal(t, "10", got.PaymentMethodCode)
}

func TestGetDianConfig_NoExponeElSecreto(t *testing.T) {
	integration := &mocks.DianIntegrationMock{
		GetConfigFn: func(ctx context.Context, integrationID uint) (map[string]interface{}, error) {
			return map[string]interface{}{
				"username_display": "user@x.com",
				"client_secret":    "SUPER_SECRETO",
				"password":         "SUPER_SECRETO",
			}, nil
		},
	}

	got, err := newAccountingUseCase(nil, nil, nil, integration).GetDianConfig(context.Background())

	require.NoError(t, err)
	assert.NotContains(t, got.Username, "SUPER_SECRETO")
	assert.NotContains(t, got.ApiURL, "SUPER_SECRETO")
	assert.NotContains(t, got.MunicipalityID, "SUPER_SECRETO")
}

func TestSaveDianConfig_CredencialesIncompletas_Rechaza(t *testing.T) {
	casos := []struct {
		nombre string
		ajuste func(*dtos.DianConfigDTO)
	}{
		{"sin client_id", func(d *dtos.DianConfigDTO) { d.ClientID = "" }},
		{"sin client_secret", func(d *dtos.DianConfigDTO) { d.ClientSecret = "" }},
		{"sin username", func(d *dtos.DianConfigDTO) { d.Username = "" }},
		{"sin password", func(d *dtos.DianConfigDTO) { d.Password = "" }},
		{"client_id en blanco", func(d *dtos.DianConfigDTO) { d.ClientID = "   " }},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			integration := &mocks.DianIntegrationMock{}
			dto := credencialesValidas()
			tc.ajuste(&dto)

			_, err := newAccountingUseCase(nil, nil, nil, integration).
				SaveDianConfig(context.Background(), dto)

			assert.ErrorIs(t, err, domainerrors.ErrDianCredentials)
			assert.Empty(t, integration.TestedCredentials, "no se prueba una conexion incompleta")
			assert.Empty(t, integration.EnsuredCredentials)
		})
	}
}

func TestSaveDianConfig_ProbamosLaConexionAntesDeGuardar(t *testing.T) {
	connErr := stderrors.New("401 unauthorized")
	integration := &mocks.DianIntegrationMock{
		TestConnectionFn: func(ctx context.Context, config, credentials map[string]interface{}) error {
			return connErr
		},
	}

	_, err := newAccountingUseCase(nil, nil, nil, integration).
		SaveDianConfig(context.Background(), credencialesValidas())

	assert.ErrorIs(t, err, connErr)
	assert.Empty(t, integration.EnsuredCredentials,
		"credenciales que no funcionan no deben quedar guardadas")
}

func TestSaveDianConfig_SeparaCredencialesDeConfiguracion(t *testing.T) {
	integration := &mocks.DianIntegrationMock{}

	_, err := newAccountingUseCase(nil, nil, nil, integration).SaveDianConfig(context.Background(),
		dtos.DianConfigDTO{
			ClientID: "cid", ClientSecret: "csecret", Username: "user", Password: "pass",
			ApiURL: "https://api.factus.com", NumberingRangeID: 8, MunicipalityID: "11001",
		})

	require.NoError(t, err)
	require.Len(t, integration.EnsuredCredentials, 1)
	require.Len(t, integration.EnsuredConfigs, 1)

	creds := integration.EnsuredCredentials[0]
	assert.Equal(t, "cid", creds["client_id"])
	assert.Equal(t, "csecret", creds["client_secret"])
	assert.Equal(t, "pass", creds["password"])

	config := integration.EnsuredConfigs[0]
	assert.NotContains(t, config, "client_secret",
		"el secreto va cifrado en credentials, nunca en el config visible")
	assert.NotContains(t, config, "password")
	assert.Equal(t, "user", config["username_display"], "solo se expone el usuario, para mostrar")
	assert.Equal(t, 8, config["numbering_range_id"])
	assert.Equal(t, "11001", config["municipality_id"])
}

func TestSaveDianConfig_RecortaEspaciosEnLasCredenciales(t *testing.T) {
	integration := &mocks.DianIntegrationMock{}

	_, err := newAccountingUseCase(nil, nil, nil, integration).SaveDianConfig(context.Background(),
		dtos.DianConfigDTO{
			ClientID: "  cid  ", ClientSecret: "  csecret  ",
			Username: "  user  ", Password: "  pass  ",
		})

	require.NoError(t, err)
	creds := integration.EnsuredCredentials[0]
	assert.Equal(t, "cid", creds["client_id"])
	assert.Equal(t, "user", creds["username"])
	assert.Equal(t, "pass", creds["password"])
}

func TestSaveDianConfig_CamposOpcionalesVaciosNoSeGuardan(t *testing.T) {
	integration := &mocks.DianIntegrationMock{}

	_, err := newAccountingUseCase(nil, nil, nil, integration).
		SaveDianConfig(context.Background(), credencialesValidas())

	require.NoError(t, err)
	config := integration.EnsuredConfigs[0]
	assert.NotContains(t, config, "numbering_range_id")
	assert.NotContains(t, config, "municipality_id")
	assert.NotContains(t, config, "identification_document_id")
	assert.NotContains(t, config, "legal_organization_id")
	assert.NotContains(t, config, "payment_method_code")
	assert.NotContains(t, integration.EnsuredCredentials[0], "api_url")
}

func TestSaveDianConfig_CamposOpcionalesLlenosSiSeGuardan(t *testing.T) {
	integration := &mocks.DianIntegrationMock{}

	_, err := newAccountingUseCase(nil, nil, nil, integration).SaveDianConfig(context.Background(),
		dtos.DianConfigDTO{
			ClientID: "cid", ClientSecret: "cs", Username: "u", Password: "p",
			ApiURL: "https://api.factus.com", NumberingRangeID: 8,
			MunicipalityID: "11001", DocumentIDType: "3", LegalOrgID: "1", PaymentMethodCode: "10",
		})

	require.NoError(t, err)
	config := integration.EnsuredConfigs[0]
	assert.Equal(t, "3", config["identification_document_id"])
	assert.Equal(t, "1", config["legal_organization_id"])
	assert.Equal(t, "10", config["payment_method_code"])
	assert.Equal(t, "https://api.factus.com", integration.EnsuredCredentials[0]["api_url"])
}

func TestSaveDianConfig_ErrorAlGuardar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("no se pudo cifrar")
	integration := &mocks.DianIntegrationMock{
		EnsureFn: func(ctx context.Context, name string, config, credentials map[string]interface{}) (uint, error) {
			return 0, dbErr
		},
	}

	_, err := newAccountingUseCase(nil, nil, nil, integration).
		SaveDianConfig(context.Background(), credencialesValidas())

	assert.ErrorIs(t, err, dbErr)
}

func TestEmitInvoiceDian_FacturaCancelada_NoEmite(t *testing.T) {
	emitter := &mocks.DianEmitterMock{}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusCancelled}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, emitter, nil).
		EmitInvoiceDian(context.Background(), dtos.EmitInvoiceDianDTO{ID: 9})

	assert.ErrorIs(t, err, domainerrors.ErrInvoiceCancelled)
	assert.Empty(t, emitter.EmitConfigs)
}

func TestEmitInvoiceDian_YaEmitida_NoReemite(t *testing.T) {
	casos := []struct {
		nombre  string
		invoice *entities.Invoice
	}{
		{"tiene CUFE", &entities.Invoice{ID: 9, Status: entities.InvoiceStatusSent, Cufe: "abc123"}},
		{"esta validada", &entities.Invoice{ID: 9, Status: entities.InvoiceStatusSent, DianStatus: entities.DianStatusValidated}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			emitter := &mocks.DianEmitterMock{}
			repo := &mocks.RepositoryMock{
				GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
					return tc.invoice, nil
				},
			}

			_, err := newAccountingUseCase(repo, nil, emitter, nil).
				EmitInvoiceDian(context.Background(), dtos.EmitInvoiceDianDTO{ID: 9})

			assert.ErrorIs(t, err, domainerrors.ErrDianAlreadyEmitted,
				"reemitir generaria una factura duplicada ante la DIAN")
			assert.Empty(t, emitter.EmitConfigs)
		})
	}
}

func TestEmitInvoiceDian_SinDocumentoDelCliente_Rechaza(t *testing.T) {
	emitter := &mocks.DianEmitterMock{}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusSent}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, emitter, nil).
		EmitInvoiceDian(context.Background(), dtos.EmitInvoiceDianDTO{ID: 9})

	assert.ErrorIs(t, err, domainerrors.ErrDianDocumentRequired)
	assert.Empty(t, emitter.EmitConfigs)
}

func TestEmitInvoiceDian_ElDocumentoDelRequestGanaSobreElDeLaFactura(t *testing.T) {
	var guardadoDoc, guardadoTel, guardadaDir string
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{
				ID: id, Status: entities.InvoiceStatusSent,
				CustomerDocument: "111", CustomerPhone: "300", CustomerAddress: "Cra 1",
			}, nil
		},
		SetInvoiceCustomerFn: func(ctx context.Context, id uint, document, phone, address string) error {
			guardadoDoc, guardadoTel, guardadaDir = document, phone, address
			return nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).EmitInvoiceDian(context.Background(),
		dtos.EmitInvoiceDianDTO{ID: 9, CustomerDocument: "  900123456  ", CustomerPhone: "  3001234567  "})

	require.NoError(t, err)
	assert.Equal(t, "900123456", guardadoDoc)
	assert.Equal(t, "3001234567", guardadoTel)
	assert.Equal(t, "Cra 1", guardadaDir, "lo que no viene en el request se conserva de la factura")
}

func TestEmitInvoiceDian_SinCambiosEnElCliente_NoReescribe(t *testing.T) {
	reescrito := false
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{
				ID: id, Status: entities.InvoiceStatusSent,
				CustomerDocument: "900123456", CustomerPhone: "300", CustomerAddress: "Cra 1",
			}, nil
		},
		SetInvoiceCustomerFn: func(ctx context.Context, id uint, document, phone, address string) error {
			reescrito = true
			return nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).EmitInvoiceDian(context.Background(),
		dtos.EmitInvoiceDianDTO{ID: 9})

	require.NoError(t, err)
	assert.False(t, reescrito)
}

func TestEmitInvoiceDian_ErrorAlGuardarElCliente_NoEmite(t *testing.T) {
	dbErr := stderrors.New("db caida")
	emitter := &mocks.DianEmitterMock{}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusSent}, nil
		},
		SetInvoiceCustomerFn: func(ctx context.Context, id uint, document, phone, address string) error {
			return dbErr
		},
	}

	_, err := newAccountingUseCase(repo, nil, emitter, nil).EmitInvoiceDian(context.Background(),
		dtos.EmitInvoiceDianDTO{ID: 9, CustomerDocument: "900123456"})

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, emitter.EmitConfigs)
}

func TestEmitInvoiceDian_SinIntegracionConfigurada_Rechaza(t *testing.T) {
	emitter := &mocks.DianEmitterMock{}
	integration := &mocks.DianIntegrationMock{
		GetGlobalIntegrationIDFn: func(ctx context.Context) (uint, error) {
			return 0, stderrors.New("no existe")
		},
	}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusSent, CustomerDocument: "900"}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, emitter, integration).
		EmitInvoiceDian(context.Background(), dtos.EmitInvoiceDianDTO{ID: 9})

	assert.ErrorIs(t, err, domainerrors.ErrDianNotConfigured)
	assert.Empty(t, emitter.EmitConfigs)
}

func TestEmitInvoiceDian_TraduceElPerfilFiscalAlFormatoDeFactus(t *testing.T) {
	casos := []struct {
		nombre         string
		personType     string
		documentType   string
		wantLegalOrg   string
		wantDocumentID string
	}{
		{"juridica con NIT", "JURIDICA", "NIT", "1", "3"},
		{"natural con CC", "NATURAL", "CC", "2", "1"},
		{"natural con NIT", "NATURAL", "NIT", "2", "3"},
		{"juridica con CC", "JURIDICA", "CC", "1", "1"},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			emitter := &mocks.DianEmitterMock{}
			repo := &mocks.RepositoryMock{
				GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
					return &entities.Invoice{
						ID: id, Number: "FV-1", Status: entities.InvoiceStatusSent, CustomerDocument: "900",
					}, nil
				},
				GetClientProfileFn: func(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
					return &dtos.ClientProfileDTO{
						Configured: true, PersonType: tc.personType, DocumentType: tc.documentType,
						DV: "7", MunicipalityID: "11001",
					}, nil
				},
			}

			_, err := newAccountingUseCase(repo, nil, emitter, nil).
				EmitInvoiceDian(context.Background(), dtos.EmitInvoiceDianDTO{ID: 9})

			require.NoError(t, err)
			require.Len(t, emitter.EmitConfigs, 1)
			config := emitter.EmitConfigs[0]
			assert.Equal(t, tc.wantLegalOrg, config["legal_organization_id"])
			assert.Equal(t, tc.wantDocumentID, config["identification_document_id"])
			assert.Equal(t, "7", config["customer_dv"])
			assert.Equal(t, "11001", config["municipality_id"])
			assert.Equal(t, "FV-1", config["reference_code"])
		})
	}
}

func TestEmitInvoiceDian_PerfilSinConfigurar_EmiteSoloConLaReferencia(t *testing.T) {
	emitter := &mocks.DianEmitterMock{}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{
				ID: id, Number: "FV-1", Status: entities.InvoiceStatusSent, CustomerDocument: "900",
			}, nil
		},
		GetClientProfileFn: func(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
			return &dtos.ClientProfileDTO{Configured: false}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, emitter, nil).
		EmitInvoiceDian(context.Background(), dtos.EmitInvoiceDianDTO{ID: 9})

	require.NoError(t, err)
	require.Len(t, emitter.EmitConfigs, 1)
	assert.Equal(t, "FV-1", emitter.EmitConfigs[0]["reference_code"])
	assert.NotContains(t, emitter.EmitConfigs[0], "legal_organization_id")
}

func TestEmitInvoiceDian_ErrorAlLeerElPerfil_NoAborta(t *testing.T) {
	emitter := &mocks.DianEmitterMock{}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{
				ID: id, Number: "FV-1", Status: entities.InvoiceStatusSent, CustomerDocument: "900",
			}, nil
		},
		GetClientProfileFn: func(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
			return nil, stderrors.New("timeout")
		},
	}

	_, err := newAccountingUseCase(repo, nil, emitter, nil).
		EmitInvoiceDian(context.Background(), dtos.EmitInvoiceDianDTO{ID: 9})

	require.NoError(t, err, "sin perfil fiscal se emite igual con los datos minimos")
	assert.Len(t, emitter.EmitConfigs, 1)
}

func TestEmitInvoiceDian_FallaLaEmision_NoGuardaResultado(t *testing.T) {
	emitErr := stderrors.New("factus rechazo la factura")
	guardado := false
	emitter := &mocks.DianEmitterMock{
		EmitFn: func(ctx context.Context, integrationID uint, inv *entities.Invoice, config map[string]interface{}) (*entities.DianEmitResult, error) {
			return nil, emitErr
		},
	}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusSent, CustomerDocument: "900"}, nil
		},
		SetInvoiceDianFn: func(ctx context.Context, id uint, result *entities.DianEmitResult) error {
			guardado = true
			return nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, emitter, nil).
		EmitInvoiceDian(context.Background(), dtos.EmitInvoiceDianDTO{ID: 9})

	assert.ErrorIs(t, err, emitErr)
	assert.False(t, guardado)
}

func TestEmitInvoiceDian_Exitoso_GuardaElCufe(t *testing.T) {
	var guardado *entities.DianEmitResult
	emitter := &mocks.DianEmitterMock{
		EmitFn: func(ctx context.Context, integrationID uint, inv *entities.Invoice, config map[string]interface{}) (*entities.DianEmitResult, error) {
			return &entities.DianEmitResult{Number: "SETP-990", CUFE: "abc123", QRCode: "qr"}, nil
		},
	}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusSent, CustomerDocument: "900"}, nil
		},
		SetInvoiceDianFn: func(ctx context.Context, id uint, result *entities.DianEmitResult) error {
			guardado = result
			return nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, emitter, nil).
		EmitInvoiceDian(context.Background(), dtos.EmitInvoiceDianDTO{ID: 9})

	require.NoError(t, err)
	require.NotNil(t, guardado)
	assert.Equal(t, "abc123", guardado.CUFE)
	assert.Equal(t, "SETP-990", guardado.Number)
}

func TestEmitInvoiceDian_ErrorAlGuardarElResultado_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusSent, CustomerDocument: "900"}, nil
		},
		SetInvoiceDianFn: func(ctx context.Context, id uint, result *entities.DianEmitResult) error {
			return dbErr
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).
		EmitInvoiceDian(context.Background(), dtos.EmitInvoiceDianDTO{ID: 9})

	assert.ErrorIs(t, err, dbErr)
}

func TestGetClientProfile_SinConfigurar_NoSugiereImpuestos(t *testing.T) {
	consultados := false
	repo := &mocks.RepositoryMock{
		GetClientProfileFn: func(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
			return &dtos.ClientProfileDTO{BusinessID: businessID, Configured: false}, nil
		},
		ListTaxesFn: func(ctx context.Context) ([]entities.Tax, error) {
			consultados = true
			return nil, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).GetClientProfile(context.Background(), 26)

	require.NoError(t, err)
	assert.False(t, got.Configured)
	assert.Empty(t, got.SuggestedTaxIDs)
	assert.False(t, consultados)
}

func TestGetClientProfile_SugiereSoloLosImpuestosQueAplica(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClientProfileFn: func(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
			return &dtos.ClientProfileDTO{
				BusinessID: businessID, Configured: true,
				AppliesIVA: true, AppliesReteFuente: false, AppliesReteICA: true,
			}, nil
		},
		ListTaxesFn: func(ctx context.Context) ([]entities.Tax, error) {
			return []entities.Tax{
				{ID: 1, Code: "IVA", IsActive: true},
				{ID: 2, Code: "RETEFUENTE", IsActive: true},
				{ID: 3, Code: "RETEICA", IsActive: true},
				{ID: 4, Code: "OTRO", IsActive: true},
			}, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).GetClientProfile(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, []uint{1, 3}, got.SuggestedTaxIDs,
		"solo IVA y ReteICA porque ReteFuente no aplica y OTRO no es un codigo conocido")
}

func TestGetClientProfile_IgnoraImpuestosInactivos(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClientProfileFn: func(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
			return &dtos.ClientProfileDTO{Configured: true, AppliesIVA: true}, nil
		},
		ListTaxesFn: func(ctx context.Context) ([]entities.Tax, error) {
			return []entities.Tax{{ID: 1, Code: "IVA", IsActive: false}}, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).GetClientProfile(context.Background(), 26)

	require.NoError(t, err)
	assert.Empty(t, got.SuggestedTaxIDs)
}

func TestGetClientProfile_ErrorAlListarImpuestos_DevuelveElPerfilSinSugerencias(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetClientProfileFn: func(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
			return &dtos.ClientProfileDTO{BusinessID: businessID, Configured: true, AppliesIVA: true}, nil
		},
		ListTaxesFn: func(ctx context.Context) ([]entities.Tax, error) {
			return nil, stderrors.New("timeout")
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).GetClientProfile(context.Background(), 26)

	require.NoError(t, err, "las sugerencias son un extra, no deben tumbar la consulta del perfil")
	require.NotNil(t, got)
	assert.True(t, got.Configured)
	assert.Empty(t, got.SuggestedTaxIDs)
}

func TestGetClientProfile_ErrorAlLeerElPerfil_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetClientProfileFn: func(ctx context.Context, businessID uint) (*dtos.ClientProfileDTO, error) {
			return nil, dbErr
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).GetClientProfile(context.Background(), 26)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestSync_IgnoraConceptosNoAutomaticosInactivosYManuales(t *testing.T) {
	var consultados []string
	repo := &mocks.RepositoryMock{
		ListConceptsFn: func(ctx context.Context) ([]entities.Concept, error) {
			return []entities.Concept{
				{ID: 1, SourceType: entities.SourceGuideMargin, IsAutomatic: false, IsActive: true},
				{ID: 2, SourceType: entities.SourceSubscription, IsAutomatic: true, IsActive: false},
				{ID: 3, SourceType: entities.SourceManual, IsAutomatic: true, IsActive: true},
				{ID: 4, SourceType: entities.SourceCodPayout, IsAutomatic: true, IsActive: true},
			}, nil
		},
		FindSyncCandidatesFn: func(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error) {
			consultados = append(consultados, sourceType)
			return nil, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).Sync(context.Background())

	require.NoError(t, err)
	assert.Equal(t, []string{entities.SourceCodPayout}, consultados,
		"solo se sincroniza el concepto automatico, activo y con origen distinto de MANUAL")
	assert.Zero(t, got.Created)
}

func TestSync_CreaLosMovimientosDeCadaCandidato(t *testing.T) {
	fecha := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	businessID := uint(26)
	llamadas := 0
	repo := &mocks.RepositoryMock{
		ListConceptsFn: func(ctx context.Context) ([]entities.Concept, error) {
			return []entities.Concept{
				{ID: 1, Kind: entities.KindIncome, SourceType: entities.SourceGuideMargin, IsAutomatic: true, IsActive: true},
			}, nil
		},
		FindSyncCandidatesFn: func(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error) {
			llamadas++
			if llamadas > 1 {
				return nil, nil
			}
			return []dtos.SyncCandidate{
				{SourceType: sourceType, SourceID: "G-1", BusinessID: &businessID, Amount: 2000.456, EntryDate: fecha, Description: "guia 1"},
				{SourceType: sourceType, SourceID: "G-2", BusinessID: &businessID, Amount: 3000, EntryDate: fecha, Description: "guia 2"},
			}, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).Sync(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 2, got.Created)
	assert.Equal(t, 2, got.BySource[entities.SourceGuideMargin])
	require.Len(t, repo.CreatedEntries, 2)
	assert.InDelta(t, 2000.46, repo.CreatedEntries[0].Amount, 0.001, "el monto se redondea a dos decimales")
	assert.Equal(t, "G-1", repo.CreatedEntries[0].SourceID)
	assert.True(t, repo.CreatedEntries[0].IsAutomatic)
	assert.Equal(t, entities.KindIncome, repo.CreatedEntries[0].Kind)
}

func TestSync_CandidatoYaAsentado_NoCuentaComoCreado(t *testing.T) {
	llamadas := 0
	repo := &mocks.RepositoryMock{
		ListConceptsFn: func(ctx context.Context) ([]entities.Concept, error) {
			return []entities.Concept{
				{ID: 1, SourceType: entities.SourceGuideMargin, IsAutomatic: true, IsActive: true},
			}, nil
		},
		FindSyncCandidatesFn: func(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error) {
			llamadas++
			if llamadas > 1 {
				return nil, nil
			}
			return []dtos.SyncCandidate{{SourceType: sourceType, SourceID: "G-1"}}, nil
		},
		CreateEntryFn: func(ctx context.Context, e *entities.Entry) (bool, error) {
			return false, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).Sync(context.Background())

	require.NoError(t, err)
	assert.Zero(t, got.Created,
		"el sync es idempotente: un movimiento que ya existia no se cuenta como nuevo")
	assert.NotContains(t, got.BySource, entities.SourceGuideMargin)
}

func TestSync_FalloEnUnCandidato_NoAbortaLosDemas(t *testing.T) {
	llamadas := 0
	repo := &mocks.RepositoryMock{
		ListConceptsFn: func(ctx context.Context) ([]entities.Concept, error) {
			return []entities.Concept{
				{ID: 1, SourceType: entities.SourceGuideMargin, IsAutomatic: true, IsActive: true},
			}, nil
		},
		FindSyncCandidatesFn: func(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error) {
			llamadas++
			if llamadas > 1 {
				return nil, nil
			}
			return []dtos.SyncCandidate{
				{SourceType: sourceType, SourceID: "G-1"},
				{SourceType: sourceType, SourceID: "G-2"},
			}, nil
		},
		CreateEntryFn: func(ctx context.Context, e *entities.Entry) (bool, error) {
			if e.SourceID == "G-1" {
				return false, stderrors.New("constraint violation")
			}
			return true, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).Sync(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, got.Created, "una guia que falla no debe dejar sin asentar a las demas")
}

func TestSync_FalloDeUnaFuente_NoAbortaLasOtras(t *testing.T) {
	llamadasPorFuente := map[string]int{}
	repo := &mocks.RepositoryMock{
		ListConceptsFn: func(ctx context.Context) ([]entities.Concept, error) {
			return []entities.Concept{
				{ID: 1, SourceType: entities.SourceGuideMargin, IsAutomatic: true, IsActive: true},
				{ID: 2, SourceType: entities.SourceSubscription, IsAutomatic: true, IsActive: true},
			}, nil
		},
		FindSyncCandidatesFn: func(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error) {
			llamadasPorFuente[sourceType]++
			if sourceType == entities.SourceGuideMargin {
				return nil, stderrors.New("query fallida")
			}
			if llamadasPorFuente[sourceType] > 1 {
				return nil, nil
			}
			return []dtos.SyncCandidate{{SourceType: sourceType, SourceID: "S-1"}}, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).Sync(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, got.Created)
	assert.Equal(t, 1, got.BySource[entities.SourceSubscription])
	assert.NotContains(t, got.BySource, entities.SourceGuideMargin)
}

func TestSync_AplicaLosImpuestosDelConcepto(t *testing.T) {
	llamadas := 0
	repo := &mocks.RepositoryMock{
		ListConceptsFn: func(ctx context.Context) ([]entities.Concept, error) {
			return []entities.Concept{
				{
					ID: 1, SourceType: entities.SourceGuideMargin, IsAutomatic: true, IsActive: true,
					Taxes: []entities.ConceptTax{{TaxCode: "IVA", RatePercent: 19, IsActive: true}},
				},
			}, nil
		},
		FindSyncCandidatesFn: func(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error) {
			llamadas++
			if llamadas > 1 {
				return nil, nil
			}
			return []dtos.SyncCandidate{{SourceType: sourceType, SourceID: "G-1", Amount: 100000}}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).Sync(context.Background())

	require.NoError(t, err)
	require.Len(t, repo.CreatedEntries, 1)
	assert.InDelta(t, 19000.0, repo.CreatedEntries[0].TaxTotal, 0.001)
}

func TestSync_SeDetieneCuandoUnaRondaNoInsertaNada(t *testing.T) {
	llamadas := 0
	repo := &mocks.RepositoryMock{
		ListConceptsFn: func(ctx context.Context) ([]entities.Concept, error) {
			return []entities.Concept{
				{ID: 1, SourceType: entities.SourceGuideMargin, IsAutomatic: true, IsActive: true},
			}, nil
		},
		FindSyncCandidatesFn: func(ctx context.Context, sourceType string, limit int) ([]dtos.SyncCandidate, error) {
			llamadas++
			candidatos := make([]dtos.SyncCandidate, syncBatchSize)
			for i := range candidatos {
				candidatos[i] = dtos.SyncCandidate{SourceType: sourceType, SourceID: "G"}
			}
			return candidatos, nil
		},
		CreateEntryFn: func(ctx context.Context, e *entities.Entry) (bool, error) {
			return false, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).Sync(context.Background())

	require.NoError(t, err)
	assert.Zero(t, got.Created)
	assert.Equal(t, 1, llamadas,
		"si un lote lleno no inserta nada nuevo se corta, no se cicla infinitamente")
}

func TestSync_ErrorAlListarConceptos_Aborta(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListConceptsFn: func(ctx context.Context) ([]entities.Concept, error) { return nil, dbErr },
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).Sync(context.Background())

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestSync_SinConceptos_ResultadoVacioNoNil(t *testing.T) {
	got, err := newAccountingUseCase(nil, nil, nil, nil).Sync(context.Background())

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Zero(t, got.Created)
	assert.NotNil(t, got.BySource)
}

func TestFormatCOP_SeparaMilesConPunto(t *testing.T) {
	casos := []struct {
		entrada float64
		want    string
	}{
		{0, "$ 0"},
		{100, "$ 100"},
		{1000, "$ 1.000"},
		{100000, "$ 100.000"},
		{1234567, "$ 1.234.567"},
		{-1000, "-$ 1.000"},
	}

	for _, tc := range casos {
		assert.Equal(t, tc.want, formatCOP(tc.entrada), "entrada=%v", tc.entrada)
	}
}

func TestFormatCOP_MuestraCentavosSoloSiLosHay(t *testing.T) {
	assert.Equal(t, "$ 1.000", formatCOP(1000))
	assert.Equal(t, "$ 1.000,50", formatCOP(1000.50))
	assert.Contains(t, formatCOP(1000.25), ",25")
}

func TestRenderInvoiceHTML_EscapaElContenidoDelCliente(t *testing.T) {
	inv := &entities.Invoice{
		Number: `<script>alert(1)</script>`,
		Notes:  `<img src=x onerror=alert(1)>`,
		Items: []entities.InvoiceItem{
			{Description: `"><script>alert(2)</script>`, Quantity: 1, UnitPrice: 100, Amount: 100},
		},
	}

	got := renderInvoiceHTML(inv)

	assert.NotContains(t, got, "<script>alert(1)</script>",
		"el numero de factura va escapado: el correo no debe poder inyectar HTML")
	assert.NotContains(t, got, "<script>alert(2)</script>")
	assert.NotContains(t, got, "<img src=x",
		"el < de la etiqueta va escapado, asi que el onerror queda como texto inerte")
	assert.Contains(t, got, "&lt;script&gt;")
	assert.Contains(t, got, "&lt;img src=x onerror=alert(1)&gt;")
}

func TestRenderInvoiceHTML_IncluyeLosTotales(t *testing.T) {
	inv := &entities.Invoice{
		Number:   "FV-1",
		Subtotal: 100000,
		TaxTotal: 19000,
		Total:    119000,
		Items: []entities.InvoiceItem{
			{Description: "Servicio", Quantity: 2, UnitPrice: 50000, Amount: 100000},
		},
		TaxDetail: []entities.TaxLine{{TaxCode: "IVA", TaxName: "IVA", RatePercent: 19, Amount: 19000}},
	}

	got := renderInvoiceHTML(inv)

	assert.Contains(t, got, "FV-1")
	assert.Contains(t, got, "Servicio")
	assert.Contains(t, got, formatCOP(119000))
	assert.True(t, strings.Contains(got, "IVA"))
}
