package app

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAccountingUseCase(
	repo *mocks.RepositoryMock,
	email *mocks.EmailSenderMock,
	emitter *mocks.DianEmitterMock,
	integration *mocks.DianIntegrationMock,
) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	if email == nil {
		email = &mocks.EmailSenderMock{}
	}
	if emitter == nil {
		emitter = &mocks.DianEmitterMock{}
	}
	if integration == nil {
		integration = &mocks.DianIntegrationMock{}
	}
	return New(repo, email, emitter, integration, mocks.NewSilentLogger())
}

func itemsValidos() []dtos.InvoiceItemDTO {
	return []dtos.InvoiceItemDTO{
		{Description: "Servicio de envios", Quantity: 2, UnitPrice: 50000},
	}
}

func TestRound2_RedondeaADosDecimales(t *testing.T) {
	casos := []struct {
		entrada float64
		want    float64
	}{
		{0, 0},
		{1.004, 1.0},
		{1.006, 1.01},
		{19.999, 20.0},
		{123456.789, 123456.79},
		{100000, 100000},
	}

	for _, tc := range casos {
		assert.InDelta(t, tc.want, round2(tc.entrada), 0.0001, "entrada=%v", tc.entrada)
	}
}

func TestRound2_ElMedioCentavoNoSiempreSubePorLaRepresentacionFloat(t *testing.T) {
	assert.InDelta(t, 1.0, round2(1.005), 0.0001,
		"1.005 se guarda como 1.00499... en float64, asi que round2 lo baja a 1.00 en vez de subirlo a 1.01")
	assert.InDelta(t, 2.01, round2(2.005), 0.0001,
		"2.005 si se guarda por encima y sube: el sesgo del medio centavo depende del valor, no es uniforme")
}

func TestBuildInvoice_SinItems_Rechaza(t *testing.T) {
	uc := newAccountingUseCase(nil, nil, nil, nil)

	got, err := uc.CreateInvoice(context.Background(), dtos.CreateInvoiceDTO{BusinessID: 26, Items: nil})

	assert.ErrorIs(t, err, domainerrors.ErrInvoiceNoItems)
	assert.Nil(t, got)
}

func TestBuildInvoice_ItemsInvalidos_Rechaza(t *testing.T) {
	casos := []struct {
		nombre string
		item   dtos.InvoiceItemDTO
	}{
		{"sin descripcion", dtos.InvoiceItemDTO{Description: "", Quantity: 1, UnitPrice: 100}},
		{"descripcion en blanco", dtos.InvoiceItemDTO{Description: "   ", Quantity: 1, UnitPrice: 100}},
		{"cantidad cero", dtos.InvoiceItemDTO{Description: "X", Quantity: 0, UnitPrice: 100}},
		{"cantidad negativa", dtos.InvoiceItemDTO{Description: "X", Quantity: -1, UnitPrice: 100}},
		{"precio cero", dtos.InvoiceItemDTO{Description: "X", Quantity: 1, UnitPrice: 0}},
		{"precio negativo", dtos.InvoiceItemDTO{Description: "X", Quantity: 1, UnitPrice: -100}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			got, err := newAccountingUseCase(nil, nil, nil, nil).CreateInvoice(context.Background(),
				dtos.CreateInvoiceDTO{BusinessID: 26, Items: []dtos.InvoiceItemDTO{tc.item}})

			assert.ErrorIs(t, err, domainerrors.ErrInvoiceInvalidItem)
			assert.Nil(t, got)
		})
	}
}

func TestBuildInvoice_ConceptoDeEgreso_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetConceptByIDFn: func(ctx context.Context, id uint) (*entities.Concept, error) {
			return &entities.Concept{ID: id, Kind: entities.KindExpense}, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{BusinessID: 26, ConceptID: 1, Items: itemsValidos()})

	assert.ErrorIs(t, err, domainerrors.ErrInvoiceConceptKind,
		"una factura no puede colgar de un concepto de egreso")
	assert.Nil(t, got)
}

func TestBuildInvoice_CalculaSubtotalComoCantidadPorPrecio(t *testing.T) {
	var creada *entities.Invoice
	repo := &mocks.RepositoryMock{
		CreateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			creada = inv
			inv.ID = 9
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{
			BusinessID: 26, ConceptID: 1,
			Items: []dtos.InvoiceItemDTO{
				{Description: "A", Quantity: 2, UnitPrice: 50000},
				{Description: "B", Quantity: 3, UnitPrice: 1000},
			},
		})

	require.NoError(t, err)
	require.NotNil(t, creada)
	require.Len(t, creada.Items, 2)
	assert.InDelta(t, 100000.0, creada.Items[0].Amount, 0.001)
	assert.InDelta(t, 3000.0, creada.Items[1].Amount, 0.001)
	assert.InDelta(t, 103000.0, creada.Subtotal, 0.001)
}

func TestBuildInvoice_RedondeaCadaItemADosDecimales(t *testing.T) {
	var creada *entities.Invoice
	repo := &mocks.RepositoryMock{
		CreateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			creada = inv
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{
			BusinessID: 26, ConceptID: 1,
			Items: []dtos.InvoiceItemDTO{{Description: "A", Quantity: 3, UnitPrice: 33.333}},
		})

	require.NoError(t, err)
	assert.InDelta(t, 100.0, creada.Items[0].Amount, 0.001)
	assert.InDelta(t, 100.0, creada.Subtotal, 0.001)
}

func TestBuildInvoice_ImpuestoDeCargo_SumaAlTotal(t *testing.T) {
	var creada *entities.Invoice
	repo := &mocks.RepositoryMock{
		GetTaxesByIDsFn: func(ctx context.Context, ids []uint) ([]entities.Tax, error) {
			return []entities.Tax{
				{ID: 1, Code: "IVA", Name: "IVA", RatePercent: 19, Kind: entities.TaxKindCharge, IsActive: true},
			}, nil
		},
		CreateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			creada = inv
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{BusinessID: 26, ConceptID: 1, TaxIDs: []uint{1}, Items: itemsValidos()})

	require.NoError(t, err)
	assert.InDelta(t, 100000.0, creada.Subtotal, 0.001)
	assert.InDelta(t, 19000.0, creada.TaxTotal, 0.001)
	assert.InDelta(t, 119000.0, creada.Total, 0.001, "el total es subtotal + impuestos de cargo")
	require.Len(t, creada.TaxDetail, 1)
	assert.Equal(t, "IVA", creada.TaxDetail[0].TaxCode)
	assert.Empty(t, creada.WithholdingDetail)
}

func TestBuildInvoice_RetencionNoRestaDelTotal(t *testing.T) {
	var creada *entities.Invoice
	repo := &mocks.RepositoryMock{
		GetTaxesByIDsFn: func(ctx context.Context, ids []uint) ([]entities.Tax, error) {
			return []entities.Tax{
				{ID: 2, Code: "RETEFUENTE", Name: "ReteFuente", RatePercent: 11, Kind: entities.TaxKindWithholding, IsActive: true},
			}, nil
		},
		CreateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			creada = inv
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{BusinessID: 26, ConceptID: 1, TaxIDs: []uint{2}, Items: itemsValidos()})

	require.NoError(t, err)
	assert.InDelta(t, 11000.0, creada.WithholdingTotal, 0.001)
	assert.Zero(t, creada.TaxTotal, "una retencion no es un impuesto de cargo")
	assert.InDelta(t, 100000.0, creada.Total, 0.001,
		"la retencion se informa aparte, no cambia el total a pagar")
	require.Len(t, creada.WithholdingDetail, 1)
	assert.Empty(t, creada.TaxDetail)
}

func TestBuildInvoice_CargoYRetencionJuntos_VanASusRespectivosBloques(t *testing.T) {
	var creada *entities.Invoice
	repo := &mocks.RepositoryMock{
		GetTaxesByIDsFn: func(ctx context.Context, ids []uint) ([]entities.Tax, error) {
			return []entities.Tax{
				{ID: 1, Code: "IVA", RatePercent: 19, Kind: entities.TaxKindCharge, IsActive: true},
				{ID: 2, Code: "RETEFUENTE", RatePercent: 11, Kind: entities.TaxKindWithholding, IsActive: true},
			}, nil
		},
		CreateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			creada = inv
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{BusinessID: 26, ConceptID: 1, TaxIDs: []uint{1, 2}, Items: itemsValidos()})

	require.NoError(t, err)
	assert.Len(t, creada.TaxDetail, 1)
	assert.Len(t, creada.WithholdingDetail, 1)
	assert.InDelta(t, 19000.0, creada.TaxTotal, 0.001)
	assert.InDelta(t, 11000.0, creada.WithholdingTotal, 0.001)
	assert.InDelta(t, 119000.0, creada.Total, 0.001)
}

func TestBuildInvoice_ImpuestosInactivosYDeTipoOther_SeIgnoran(t *testing.T) {
	var creada *entities.Invoice
	repo := &mocks.RepositoryMock{
		GetTaxesByIDsFn: func(ctx context.Context, ids []uint) ([]entities.Tax, error) {
			return []entities.Tax{
				{ID: 1, Code: "IVA_VIEJO", RatePercent: 16, Kind: entities.TaxKindCharge, IsActive: false},
				{ID: 2, Code: "OTRO", RatePercent: 5, Kind: entities.TaxKindOther, IsActive: true},
			}, nil
		},
		CreateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			creada = inv
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{BusinessID: 26, ConceptID: 1, TaxIDs: []uint{1, 2}, Items: itemsValidos()})

	require.NoError(t, err)
	assert.Empty(t, creada.TaxDetail, "un impuesto inactivo o de tipo OTHER no se cobra")
	assert.Zero(t, creada.TaxTotal)
	assert.InDelta(t, 100000.0, creada.Total, 0.001)
}

func TestBuildInvoice_SinImpuestos_DetallesVaciosNoNil(t *testing.T) {
	var creada *entities.Invoice
	repo := &mocks.RepositoryMock{
		CreateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			creada = inv
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{BusinessID: 26, ConceptID: 1, Items: itemsValidos()})

	require.NoError(t, err)
	assert.NotNil(t, creada.TaxDetail)
	assert.NotNil(t, creada.WithholdingDetail)
	assert.Empty(t, creada.TaxDetail)
}

func TestBuildInvoice_ErrorAlLeerImpuestos_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetTaxesByIDsFn: func(ctx context.Context, ids []uint) ([]entities.Tax, error) {
			return nil, dbErr
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{BusinessID: 26, ConceptID: 1, TaxIDs: []uint{1}, Items: itemsValidos()})

	assert.ErrorIs(t, err, dbErr)
}

func TestBuildInvoice_ErrorAlLeerConcepto_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetConceptByIDFn: func(ctx context.Context, id uint) (*entities.Concept, error) {
			return nil, domainerrors.ErrConceptNotFound
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{BusinessID: 26, ConceptID: 404, Items: itemsValidos()})

	assert.ErrorIs(t, err, domainerrors.ErrConceptNotFound)
}

func TestCreateInvoice_NegocioInexistente_NoCrea(t *testing.T) {
	creada := false
	repo := &mocks.RepositoryMock{
		GetBusinessNameFn: func(ctx context.Context, id uint) (string, error) {
			return "", domainerrors.ErrBusinessNotFound
		},
		CreateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			creada = true
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{BusinessID: 404, ConceptID: 1, Items: itemsValidos()})

	assert.ErrorIs(t, err, domainerrors.ErrBusinessNotFound)
	assert.False(t, creada)
}

func TestCreateInvoice_NaceEnBorradorYRecortaLosTextos(t *testing.T) {
	var creada *entities.Invoice
	repo := &mocks.RepositoryMock{
		CreateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			creada = inv
			inv.ID = 9
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{
			BusinessID: 26, ConceptID: 1, Items: itemsValidos(),
			Notes: "  nota  ", EmailTo: "  a@b.com  ",
			CustomerDocument: "  900123456  ", CustomerPhone: "  3001234567  ",
			CustomerAddress: "  Cra 1  ", IsTest: true,
		})

	require.NoError(t, err)
	assert.Equal(t, entities.InvoiceStatusDraft, creada.Status)
	assert.Equal(t, "nota", creada.Notes)
	assert.Equal(t, "a@b.com", creada.EmailTo)
	assert.Equal(t, "900123456", creada.CustomerDocument)
	assert.Equal(t, "3001234567", creada.CustomerPhone)
	assert.Equal(t, "Cra 1", creada.CustomerAddress)
	assert.True(t, creada.IsTest)
	assert.Equal(t, uint(26), creada.BusinessID)
}

func TestCreateInvoice_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("constraint violation")
	repo := &mocks.RepositoryMock{
		CreateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			return nil, dbErr
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).CreateInvoice(context.Background(),
		dtos.CreateInvoiceDTO{BusinessID: 26, ConceptID: 1, Items: itemsValidos()})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateInvoice_SoloEnBorrador(t *testing.T) {
	for _, estado := range []string{
		entities.InvoiceStatusSent,
		entities.InvoiceStatusPaid,
		entities.InvoiceStatusCancelled,
	} {
		t.Run(estado, func(t *testing.T) {
			actualizada := false
			repo := &mocks.RepositoryMock{
				GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
					return &entities.Invoice{ID: id, Status: estado}, nil
				},
				UpdateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
					actualizada = true
					return inv, nil
				},
			}

			_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateInvoice(context.Background(),
				dtos.UpdateInvoiceDTO{ID: 9, ConceptID: 1, Items: itemsValidos()})

			assert.ErrorIs(t, err, domainerrors.ErrInvoiceNotDraft)
			assert.False(t, actualizada, "una factura ya enviada o pagada no se puede modificar")
		})
	}
}

func TestUpdateInvoice_ConservaNegocioYEstado(t *testing.T) {
	var guardada *entities.Invoice
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, BusinessID: 26, Status: entities.InvoiceStatusDraft}, nil
		},
		UpdateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			guardada = inv
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateInvoice(context.Background(),
		dtos.UpdateInvoiceDTO{ID: 9, ConceptID: 1, Items: itemsValidos()})

	require.NoError(t, err)
	require.NotNil(t, guardada)
	assert.Equal(t, uint(9), guardada.ID)
	assert.Equal(t, uint(26), guardada.BusinessID, "el negocio dueno no se cambia por update")
	assert.Equal(t, entities.InvoiceStatusDraft, guardada.Status)
}

func TestUpdateInvoice_RecalculaLosTotales(t *testing.T) {
	var guardada *entities.Invoice
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusDraft, Subtotal: 999, Total: 999}, nil
		},
		UpdateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			guardada = inv
			return inv, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateInvoice(context.Background(),
		dtos.UpdateInvoiceDTO{
			ID: 9, ConceptID: 1,
			Items: []dtos.InvoiceItemDTO{{Description: "A", Quantity: 1, UnitPrice: 7000}},
		})

	require.NoError(t, err)
	assert.InDelta(t, 7000.0, guardada.Subtotal, 0.001, "los totales se recalculan, no se arrastran")
	assert.InDelta(t, 7000.0, guardada.Total, 0.001)
}

func TestUpdateInvoice_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("deadlock")
	repo := &mocks.RepositoryMock{
		UpdateInvoiceFn: func(ctx context.Context, inv *entities.Invoice) (*entities.Invoice, error) {
			return nil, dbErr
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).UpdateInvoice(context.Background(),
		dtos.UpdateInvoiceDTO{ID: 9, ConceptID: 1, Items: itemsValidos()})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListInvoices_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 10, 1, 10},
		{-5, 10, 1, 10},
		{1, 0, 1, 10},
		{1, -1, 1, 10},
		{1, 101, 1, 100},
		{1, 100, 1, 100},
		{3, 50, 3, 50},
	}

	for _, tc := range casos {
		var visto dtos.ListInvoicesParams
		repo := &mocks.RepositoryMock{
			ListInvoicesFn: func(ctx context.Context, params dtos.ListInvoicesParams) ([]entities.Invoice, int64, error) {
				visto = params
				return nil, 0, nil
			},
		}

		_, _, err := newAccountingUseCase(repo, nil, nil, nil).ListInvoices(context.Background(),
			dtos.ListInvoicesParams{Page: tc.page, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, visto.Page)
		assert.Equal(t, tc.wantSize, visto.PageSize)
	}
}

func TestDeleteInvoice_SoloBorradorOCancelada(t *testing.T) {
	casos := []struct {
		estado    string
		wantBorra bool
	}{
		{entities.InvoiceStatusDraft, true},
		{entities.InvoiceStatusCancelled, true},
		{entities.InvoiceStatusSent, false},
		{entities.InvoiceStatusPaid, false},
	}

	for _, tc := range casos {
		t.Run(tc.estado, func(t *testing.T) {
			borrada := false
			repo := &mocks.RepositoryMock{
				GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
					return &entities.Invoice{ID: id, Status: tc.estado}, nil
				},
				DeleteInvoiceFn: func(ctx context.Context, id uint) error {
					borrada = true
					return nil
				},
			}

			err := newAccountingUseCase(repo, nil, nil, nil).DeleteInvoice(context.Background(), 9)

			if tc.wantBorra {
				require.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, domainerrors.ErrInvoiceNotDraft,
					"borrar una factura enviada o pagada dejaria un hueco en la numeracion")
			}
			assert.Equal(t, tc.wantBorra, borrada)
		})
	}
}

func TestSendInvoice_Cancelada_NoEnvia(t *testing.T) {
	email := &mocks.EmailSenderMock{}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusCancelled}, nil
		},
	}

	_, err := newAccountingUseCase(repo, email, nil, nil).SendInvoice(context.Background(),
		dtos.SendInvoiceDTO{ID: 9})

	assert.ErrorIs(t, err, domainerrors.ErrInvoiceCancelled)
	assert.Empty(t, email.Sent)
}

func TestSendInvoice_PrioridadDelDestinatario(t *testing.T) {
	casos := []struct {
		nombre      string
		dtoEmail    string
		invoiceMail string
		defaultMail string
		want        string
	}{
		{"el del request gana", "req@x.com", "inv@x.com", "def@x.com", "req@x.com"},
		{"sin request usa el de la factura", "", "inv@x.com", "def@x.com", "inv@x.com"},
		{"sin ninguno usa el del negocio", "", "", "def@x.com", "def@x.com"},
		{"el request se recorta", "  req@x.com  ", "inv@x.com", "", "req@x.com"},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			email := &mocks.EmailSenderMock{}
			repo := &mocks.RepositoryMock{
				GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
					return &entities.Invoice{
						ID: id, Number: "FV-1", Status: entities.InvoiceStatusDraft, EmailTo: tc.invoiceMail,
					}, nil
				},
				GetBusinessDefEmailFn: func(ctx context.Context, id uint) (string, error) {
					return tc.defaultMail, nil
				},
			}

			_, err := newAccountingUseCase(repo, email, nil, nil).SendInvoice(context.Background(),
				dtos.SendInvoiceDTO{ID: 9, EmailTo: tc.dtoEmail})

			require.NoError(t, err)
			require.Len(t, email.Sent, 1)
			assert.Equal(t, tc.want, email.Sent[0].To)
		})
	}
}

func TestSendInvoice_SinNingunCorreo_Rechaza(t *testing.T) {
	email := &mocks.EmailSenderMock{}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusDraft}, nil
		},
	}

	_, err := newAccountingUseCase(repo, email, nil, nil).SendInvoice(context.Background(),
		dtos.SendInvoiceDTO{ID: 9})

	assert.ErrorIs(t, err, domainerrors.ErrInvoiceNoEmail)
	assert.Empty(t, email.Sent)
}

func TestSendInvoice_ErrorAlBuscarCorreoDelNegocio_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusDraft}, nil
		},
		GetBusinessDefEmailFn: func(ctx context.Context, id uint) (string, error) { return "", dbErr },
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).SendInvoice(context.Background(),
		dtos.SendInvoiceDTO{ID: 9})

	assert.ErrorIs(t, err, dbErr)
}

func TestSendInvoice_AsuntoIncluyeElNumero(t *testing.T) {
	email := &mocks.EmailSenderMock{}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{
				ID: id, Number: "FV-00123", Status: entities.InvoiceStatusDraft, EmailTo: "a@b.com",
			}, nil
		},
	}

	_, err := newAccountingUseCase(repo, email, nil, nil).SendInvoice(context.Background(),
		dtos.SendInvoiceDTO{ID: 9})

	require.NoError(t, err)
	require.Len(t, email.Sent, 1)
	assert.Contains(t, email.Sent[0].Subject, "FV-00123")
	assert.NotEmpty(t, email.Sent[0].HTML)
}

func TestSendInvoice_FallaElCorreo_NoCambiaElEstado(t *testing.T) {
	mailErr := stderrors.New("smtp caido")
	email := &mocks.EmailSenderMock{
		SendHTMLFn: func(ctx context.Context, to, subject, html string) error { return mailErr },
	}
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusDraft, EmailTo: "a@b.com"}, nil
		},
	}

	_, err := newAccountingUseCase(repo, email, nil, nil).SendInvoice(context.Background(),
		dtos.SendInvoiceDTO{ID: 9})

	assert.ErrorIs(t, err, mailErr)
	assert.Empty(t, repo.StatusCalls, "si el correo no salio, la factura no queda marcada como enviada")
}

func TestSendInvoice_BorradorPasaAEnviada(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusDraft, EmailTo: "a@b.com"}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).SendInvoice(context.Background(),
		dtos.SendInvoiceDTO{ID: 9})

	require.NoError(t, err)
	require.Len(t, repo.StatusCalls, 1)
	assert.Equal(t, entities.InvoiceStatusSent, repo.StatusCalls[0].Status)
	require.NotNil(t, repo.StatusCalls[0].SentAt)
}

func TestSendInvoice_ReenviarUnaPagadaNoLaDevuelveAEnviada(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusPaid, EmailTo: "a@b.com"}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).SendInvoice(context.Background(),
		dtos.SendInvoiceDTO{ID: 9})

	require.NoError(t, err)
	require.Len(t, repo.StatusCalls, 1)
	assert.Equal(t, entities.InvoiceStatusPaid, repo.StatusCalls[0].Status,
		"reenviar el correo de una factura pagada no debe revertir su estado")
}

func TestSendInvoice_ErrorAlGuardarElEstado_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusDraft, EmailTo: "a@b.com"}, nil
		},
		SetInvoiceStatusFn: func(ctx context.Context, id uint, status string, sentAt, paidAt *time.Time, emailTo string) error {
			return dbErr
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).SendInvoice(context.Background(),
		dtos.SendInvoiceDTO{ID: 9})

	assert.ErrorIs(t, err, dbErr)
}

func TestMarkInvoicePaid_EstadosNoPagables(t *testing.T) {
	casos := []struct {
		estado  string
		wantErr error
	}{
		{entities.InvoiceStatusPaid, domainerrors.ErrInvoiceAlreadyPaid},
		{entities.InvoiceStatusCancelled, domainerrors.ErrInvoiceCancelled},
		{"ESTADO_RARO", domainerrors.ErrInvoiceNotPayable},
	}

	for _, tc := range casos {
		t.Run(tc.estado, func(t *testing.T) {
			repo := &mocks.RepositoryMock{
				GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
					return &entities.Invoice{ID: id, Status: tc.estado}, nil
				},
			}

			_, err := newAccountingUseCase(repo, nil, nil, nil).MarkInvoicePaid(context.Background(), 9)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Empty(t, repo.CreatedEntries, "no se asienta nada si la factura no es pagable")
		})
	}
}

func TestMarkInvoicePaid_CreaElMovimientoContablePorElSubtotal(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{
				ID: id, BusinessID: 26, ConceptID: 3, Number: "FV-1",
				Status:   entities.InvoiceStatusSent,
				Subtotal: 100000, TaxTotal: 19000, Total: 119000,
				TaxDetail: []entities.TaxLine{{TaxCode: "IVA", Amount: 19000}},
			}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).MarkInvoicePaid(context.Background(), 9)

	require.NoError(t, err)
	require.Len(t, repo.CreatedEntries, 1)
	entry := repo.CreatedEntries[0]
	assert.InDelta(t, 100000.0, entry.Amount, 0.001,
		"el ingreso contable es el subtotal, no el total con IVA")
	assert.InDelta(t, 19000.0, entry.TaxTotal, 0.001)
	assert.Equal(t, entities.KindIncome, entry.Kind)
	assert.Equal(t, entities.SourceInvoice, entry.SourceType)
	assert.Equal(t, "FV-1", entry.SourceID)
	assert.True(t, entry.IsAutomatic)
	require.NotNil(t, entry.BusinessID)
	assert.Equal(t, uint(26), *entry.BusinessID)
	assert.Equal(t, uint(3), entry.ConceptID)
}

func TestMarkInvoicePaid_FacturaDePrueba_NoAsientaNada(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{
				ID: id, Status: entities.InvoiceStatusSent, IsTest: true, Subtotal: 100000,
			}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).MarkInvoicePaid(context.Background(), 9)

	require.NoError(t, err)
	assert.Empty(t, repo.CreatedEntries,
		"una factura de prueba no debe ensuciar la contabilidad")
	require.Len(t, repo.StatusCalls, 1)
	assert.Equal(t, entities.InvoiceStatusPaid, repo.StatusCalls[0].Status)
	require.NotNil(t, repo.StatusCalls[0].PaidAt)
}

func TestMarkInvoicePaid_FallaElAsiento_NoMarcaPagada(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusSent, Number: "FV-1"}, nil
		},
		CreateEntryFn: func(ctx context.Context, e *entities.Entry) (bool, error) { return false, dbErr },
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).MarkInvoicePaid(context.Background(), 9)

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, repo.StatusCalls,
		"si no se pudo asentar el ingreso la factura no queda pagada")
}

func TestMarkInvoicePaid_DesdeBorradorTambienSePuede(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusDraft, Number: "FV-1"}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).MarkInvoicePaid(context.Background(), 9)

	require.NoError(t, err, "se puede cobrar una factura que nunca se envio por correo")
	assert.Len(t, repo.CreatedEntries, 1)
}

func TestCancelInvoice_YaCancelada_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusCancelled}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CancelInvoice(context.Background(), 9)

	assert.ErrorIs(t, err, domainerrors.ErrInvoiceCancelled)
	assert.Empty(t, repo.StatusCalls)
}

func TestCancelInvoice_SiEstabaPagada_ReversaElMovimientoContable(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusPaid, Number: "FV-1"}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CancelInvoice(context.Background(), 9)

	require.NoError(t, err)
	assert.Equal(t, []string{"INVOICE:FV-1"}, repo.DeletedEntrySources,
		"anular una factura pagada debe borrar su ingreso de la contabilidad")
	require.Len(t, repo.StatusCalls, 1)
	assert.Equal(t, entities.InvoiceStatusCancelled, repo.StatusCalls[0].Status)
	assert.Nil(t, repo.StatusCalls[0].PaidAt, "al anular se limpia la fecha de pago")
}

func TestCancelInvoice_SiNoEstabaPagada_NoTocaLaContabilidad(t *testing.T) {
	for _, estado := range []string{entities.InvoiceStatusDraft, entities.InvoiceStatusSent} {
		t.Run(estado, func(t *testing.T) {
			repo := &mocks.RepositoryMock{
				GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
					return &entities.Invoice{ID: id, Status: estado, Number: "FV-1"}, nil
				},
			}

			_, err := newAccountingUseCase(repo, nil, nil, nil).CancelInvoice(context.Background(), 9)

			require.NoError(t, err)
			assert.Empty(t, repo.DeletedEntrySources)
		})
	}
}

func TestCancelInvoice_FallaLaReversa_NoCancela(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			return &entities.Invoice{ID: id, Status: entities.InvoiceStatusPaid, Number: "FV-1"}, nil
		},
		DeleteEntryBySourceFn: func(ctx context.Context, sourceType, sourceID string) error { return dbErr },
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CancelInvoice(context.Background(), 9)

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, repo.StatusCalls,
		"si no se pudo reversar el ingreso la factura no queda cancelada")
}

func TestGetInvoice_DelegaAlRepositorio(t *testing.T) {
	var visto uint
	repo := &mocks.RepositoryMock{
		GetInvoiceByIDFn: func(ctx context.Context, id uint) (*entities.Invoice, error) {
			visto = id
			return &entities.Invoice{ID: id, Number: "FV-1"}, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).GetInvoice(context.Background(), 9)

	require.NoError(t, err)
	assert.Equal(t, uint(9), visto)
	assert.Equal(t, "FV-1", got.Number)
}
