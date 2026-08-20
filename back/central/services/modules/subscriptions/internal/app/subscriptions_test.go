package app

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSubsUseCase(repo *mocks.RepositoryMock, wallet *mocks.WalletDebiterMock, ann *mocks.AnnouncementsGatewayMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	if wallet == nil {
		wallet = &mocks.WalletDebiterMock{}
	}
	if ann == nil {
		ann = &mocks.AnnouncementsGatewayMock{}
	}
	return New(repo, wallet, ann, mocks.NewSilentLogger())
}

func uintPtr(v uint) *uint           { return &v }
func timePtr(v time.Time) *time.Time { return &v }

func planActivo(price float64) *entities.SubscriptionType {
	return &entities.SubscriptionType{
		ID: 3, Name: "Pro", Code: "pro", Price: price, Active: true,
		ModuleCodes: []string{"orders", "shipments"}, MaxEcommerceChannels: 5,
	}
}

func TestPurchaseSubscription_MesesInvalidos_NoTocaLaBilletera(t *testing.T) {
	for _, meses := range []int{0, -1, -12} {
		wallet := &mocks.WalletDebiterMock{}

		got, err := newSubsUseCase(nil, wallet, nil).PurchaseSubscription(context.Background(),
			dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: meses})

		assert.ErrorIs(t, err, errs.ErrInvalidMonths, "meses=%d", meses)
		assert.Nil(t, got)
		assert.Empty(t, wallet.DebitCalls)
	}
}

func TestPurchaseSubscription_PlanInexistente_NoDebita(t *testing.T) {
	wallet := &mocks.WalletDebiterMock{}
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return nil, nil
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 404, Months: 1})

	assert.ErrorIs(t, err, errs.ErrSubscriptionTypeNotFound)
	assert.Empty(t, wallet.DebitCalls)
}

func TestPurchaseSubscription_PlanInactivo_NoDebita(t *testing.T) {
	wallet := &mocks.WalletDebiterMock{}
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, Active: false, Price: 100}, nil
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	assert.ErrorIs(t, err, errs.ErrSubscriptionTypeInactive)
	assert.Empty(t, wallet.DebitCalls)
}

func TestPurchaseSubscription_SaldoInsuficiente_NoDebitaNiCrea(t *testing.T) {
	creada := false
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 99999, nil },
	}
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return planActivo(50000), nil
		},
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			creada = true
			return nil
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 2})

	assert.ErrorIs(t, err, errs.ErrInsufficientBalance,
		"2 meses a 50000 son 100000 y el saldo es 99999")
	assert.Empty(t, wallet.DebitCalls)
	assert.False(t, creada)
}

func TestPurchaseSubscription_SaldoExacto_Alcanza(t *testing.T) {
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 100000, nil },
	}
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return planActivo(50000), nil
		},
	}

	got, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 2})

	require.NoError(t, err)
	assert.InDelta(t, 100000.0, got.Amount, 0.001)
	require.Len(t, wallet.DebitCalls, 1)
	assert.InDelta(t, 100000.0, wallet.DebitCalls[0].Amount, 0.001)
}

func TestPurchaseSubscription_ErrorAlLeerSaldo_NoDebita(t *testing.T) {
	walletErr := stderrors.New("wallet caida")
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 0, walletErr },
	}

	_, err := newSubsUseCase(nil, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	assert.ErrorIs(t, err, walletErr)
	assert.Empty(t, wallet.DebitCalls)
}

func TestPurchaseSubscription_MontoEsPrecioPorMeses(t *testing.T) {
	casos := []struct {
		precio float64
		meses  int
		want   float64
	}{
		{50000, 1, 50000},
		{50000, 3, 150000},
		{50000, 12, 600000},
		{33333.33, 3, 99999.99},
	}

	for _, tc := range casos {
		wallet := &mocks.WalletDebiterMock{
			GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
		}
		repo := &mocks.RepositoryMock{
			GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
				return planActivo(tc.precio), nil
			},
		}

		got, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
			dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: tc.meses})

		require.NoError(t, err)
		assert.InDelta(t, tc.want, got.Amount, 0.01, "precio=%v meses=%d", tc.precio, tc.meses)
	}
}

func TestPurchaseSubscription_ReferenciaIncluyeNegocioCodigoYMeses(t *testing.T) {
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return planActivo(1000), nil
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 6, UserID: 42})

	require.NoError(t, err)
	require.Len(t, wallet.DebitCalls, 1)
	assert.Equal(t, "SUB-26-pro-6M", wallet.DebitCalls[0].Reference)
	assert.Equal(t, "SUBSCRIPTION", wallet.DebitCalls[0].Concept)
	assert.Equal(t, uint(42), wallet.DebitCalls[0].UserID)
}

func TestPurchaseSubscription_FallaElDebito_NoCreaSuscripcion(t *testing.T) {
	creada := false
	debitErr := stderrors.New("saldo bloqueado")
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
		DebitFn: func(ctx context.Context, businessID uint, amount float64, reference, concept string, userID uint) error {
			return debitErr
		},
	}
	repo := &mocks.RepositoryMock{
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			creada = true
			return nil
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	assert.ErrorIs(t, err, debitErr)
	assert.False(t, creada, "si no se pudo cobrar no debe quedar suscripcion")
}

func TestPurchaseSubscription_DebitoOKPeroFallaElRegistro_QuedaPlataCobradaSinSuscripcion(t *testing.T) {
	dbErr := stderrors.New("db caida")
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}
	repo := &mocks.RepositoryMock{
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			return dbErr
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	assert.ErrorIs(t, err, dbErr)
	assert.Len(t, wallet.DebitCalls, 1,
		"el debito ya ocurrio y no se revierte: queda plata cobrada sin suscripcion, requiere conciliacion manual")
}

func TestPurchaseSubscription_NaceEnEstadoPagada(t *testing.T) {
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}

	got, err := newSubsUseCase(nil, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	require.NoError(t, err)
	assert.Equal(t, entities.SubscriptionStatusPaid, got.Status)
}

func TestPurchaseSubscription_MarcaElNegocioComoActivo(t *testing.T) {
	var vistoTipo uint
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return planActivo(1000), nil
		},
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			vistoTipo = subscriptionTypeID
			return nil
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	require.NoError(t, err)
	assert.Equal(t, uint(3), vistoTipo)
}

func TestPurchaseSubscription_FallaMarcarElNegocio_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}
	repo := &mocks.RepositoryMock{
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			return dbErr
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	assert.ErrorIs(t, err, dbErr)
}

func TestPurchaseSubscription_ErrorAlLeerElPlan_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	wallet := &mocks.WalletDebiterMock{}
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return nil, dbErr
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, wallet.DebitCalls)
}

func TestPurchaseSubscription_SinSuscripcionVigente_ArrancaDesdeAhora(t *testing.T) {
	antes := time.Now()
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}

	got, err := newSubsUseCase(nil, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 3})

	require.NoError(t, err)
	assert.False(t, got.StartDate.Before(antes))
	assert.Equal(t, got.StartDate.AddDate(0, 3, 0), got.EndDate)
}

func TestPurchaseSubscription_ConSuscripcionVigente_SeEncadenaAlFinal(t *testing.T) {
	finVigente := time.Now().AddDate(0, 2, 0)
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return &entities.BusinessSubscription{ID: 1, EndDate: finVigente}, nil
		},
	}

	got, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	require.NoError(t, err)
	assert.Equal(t, finVigente, got.StartDate,
		"renovar antes de vencer no debe regalar ni quitar dias: se encadena al final")
	assert.Equal(t, finVigente.AddDate(0, 1, 0), got.EndDate)
}

func TestPurchaseSubscription_ConSuscripcionYaVencida_ArrancaDesdeAhora(t *testing.T) {
	finVencido := time.Now().AddDate(0, -2, 0)
	antes := time.Now()
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return &entities.BusinessSubscription{ID: 1, EndDate: finVencido}, nil
		},
	}

	got, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	require.NoError(t, err)
	assert.False(t, got.StartDate.Before(antes),
		"una suscripcion vencida no encadena hacia atras")
}

func TestPurchaseSubscription_ErrorAlLeerLaVigente_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return nil, dbErr
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	assert.ErrorIs(t, err, dbErr,
		"un error real de BD al leer la vigente no debe acortar silenciosamente el periodo pagado")
	assert.Empty(t, wallet.DebitCalls,
		"si no se puede calcular la ventana no debe cobrarse")
}

func TestPurchaseSubscription_DesactivaLosAvisosDeVencimiento(t *testing.T) {
	ann := &mocks.AnnouncementsGatewayMock{
		FindActiveBusinessAlertFn: func(ctx context.Context, businessID uint, title string) (*uint, error) {
			if title == "Tu suscripcion vencio" {
				return uintPtr(88), nil
			}
			return uintPtr(77), nil
		},
	}
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}

	_, err := newSubsUseCase(nil, wallet, ann).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	require.NoError(t, err)
	assert.ElementsMatch(t, []uint{77, 88}, ann.DeactivatedIDs,
		"al pagar se bajan los dos avisos: el de por vencer y el de vencido")
}

func TestPurchaseSubscription_SinAvisosActivos_NoDesactivaNada(t *testing.T) {
	ann := &mocks.AnnouncementsGatewayMock{}
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}

	_, err := newSubsUseCase(nil, wallet, ann).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	require.NoError(t, err)
	assert.Empty(t, ann.DeactivatedIDs)
}

func TestPurchaseSubscription_FalloAlDesactivarAviso_NoRompeLaCompra(t *testing.T) {
	ann := &mocks.AnnouncementsGatewayMock{
		FindActiveBusinessAlertFn: func(ctx context.Context, businessID uint, title string) (*uint, error) {
			return uintPtr(77), nil
		},
		DeactivateAnnouncementFn: func(ctx context.Context, id uint) error {
			return stderrors.New("anuncios caido")
		},
	}
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1e9, nil },
	}

	got, err := newSubsUseCase(nil, wallet, ann).PurchaseSubscription(context.Background(),
		dtos.PurchaseSubscriptionDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1})

	require.NoError(t, err, "el aviso es cosmetico, no debe tumbar una compra ya cobrada")
	assert.NotNil(t, got)
}

func TestRegisterPayment_MesesInvalidos_Rechaza(t *testing.T) {
	for _, meses := range []int{0, -3} {
		_, err := newSubsUseCase(nil, nil, nil).RegisterPayment(context.Background(),
			dtos.RegisterPaymentDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: meses}, 1)
		assert.ErrorIs(t, err, errs.ErrInvalidMonths)
	}
}

func TestRegisterPayment_PlanInexistente_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return nil, nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).RegisterPayment(context.Background(),
		dtos.RegisterPaymentDTO{BusinessID: 26, SubscriptionTypeID: 404, Months: 1}, 1)
	assert.ErrorIs(t, err, errs.ErrSubscriptionTypeNotFound)
}

func TestRegisterPayment_NoTocaLaBilletera(t *testing.T) {
	wallet := &mocks.WalletDebiterMock{}

	_, err := newSubsUseCase(nil, wallet, nil).RegisterPayment(context.Background(),
		dtos.RegisterPaymentDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1}, 1)
	require.NoError(t, err)
	assert.Empty(t, wallet.DebitCalls,
		"registrar un pago externo no debe descontar de la billetera")
}

func TestRegisterPayment_AceptaPlanInactivo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, Code: "legacy", Price: 1000, Active: false}, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).RegisterPayment(context.Background(),
		dtos.RegisterPaymentDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1}, 1)
	require.NoError(t, err, "un pago manual sobre un plan viejo debe poder registrarse")
	assert.NotNil(t, got)
}

func TestRegisterPayment_ConFechaDeInicioManual_LaRespeta(t *testing.T) {
	inicio := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	consultada := false
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			consultada = true
			return &entities.BusinessSubscription{EndDate: time.Now().AddDate(1, 0, 0)}, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).RegisterPayment(context.Background(),
		dtos.RegisterPaymentDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 3, StartDate: timePtr(inicio)}, 1)
	require.NoError(t, err)
	assert.Equal(t, inicio, got.StartDate)
	assert.Equal(t, inicio.AddDate(0, 3, 0), got.EndDate)
	assert.False(t, consultada, "con fecha manual no se consulta la vigente")
}

func TestRegisterPayment_GuardaReferenciaYNotas(t *testing.T) {
	ref := "TRANSF-12345"
	notas := "pago por consignacion"
	var creada *entities.BusinessSubscription
	repo := &mocks.RepositoryMock{
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			creada = s
			return nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).RegisterPayment(context.Background(),
		dtos.RegisterPaymentDTO{
			BusinessID: 26, SubscriptionTypeID: 3, Months: 1,
			PaymentReference: &ref, Notes: &notas,
		}, 1)
	require.NoError(t, err)
	require.NotNil(t, creada)
	require.NotNil(t, creada.PaymentReference)
	assert.Equal(t, ref, *creada.PaymentReference)
	require.NotNil(t, creada.Notes)
	assert.Equal(t, notas, *creada.Notes)
	assert.Equal(t, entities.SubscriptionStatusPaid, got.Status)
}

func TestRegisterPayment_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			return dbErr
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).RegisterPayment(context.Background(),
		dtos.RegisterPaymentDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1}, 1)
	assert.ErrorIs(t, err, dbErr)
}

func TestRegisterPayment_ErrorAlMarcarElNegocio_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			return dbErr
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).RegisterPayment(context.Background(),
		dtos.RegisterPaymentDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1}, 1)
	assert.ErrorIs(t, err, dbErr)
}

func TestRegisterPayment_ErrorAlLeerElPlan_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return nil, dbErr
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).RegisterPayment(context.Background(),
		dtos.RegisterPaymentDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1}, 1)
	assert.ErrorIs(t, err, dbErr)
}

func TestEditSubscriptionDates_RangoInvalido_Rechaza(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	casos := []struct {
		nombre string
		inicio time.Time
		fin    time.Time
	}{
		{"fin antes del inicio", base, base.AddDate(0, -1, 0)},
		{"fin igual al inicio", base, base},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			actualizado := false
			repo := &mocks.RepositoryMock{
				UpdateSubscriptionDatesFn: func(ctx context.Context, id uint, s, e time.Time) error {
					actualizado = true
					return nil
				},
			}

			_, err := newSubsUseCase(repo, nil, nil).EditSubscriptionDates(context.Background(),
				dtos.EditSubscriptionDatesDTO{BusinessID: 26, StartDate: tc.inicio, EndDate: tc.fin}, 1)
			assert.ErrorIs(t, err, errs.ErrInvalidDateRange)
			assert.False(t, actualizado)
		})
	}
}

func TestEditSubscriptionDates_SinSuscripcion_Rechaza(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return nil, nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).EditSubscriptionDates(context.Background(),
		dtos.EditSubscriptionDatesDTO{BusinessID: 26, StartDate: base, EndDate: base.AddDate(0, 1, 0)}, 1)
	assert.ErrorIs(t, err, errs.ErrSubscriptionNotFound)
}

func TestEditSubscriptionDates_ActualizaAmbasTablasYDevuelveLasNuevas(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	fin := base.AddDate(0, 6, 0)
	var vistoID uint
	var vistoInicio, vistoFin, vistoFinNegocio time.Time
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return &entities.BusinessSubscription{ID: 9, BusinessID: businessID}, nil
		},
		UpdateSubscriptionDatesFn: func(ctx context.Context, id uint, s, e time.Time) error {
			vistoID, vistoInicio, vistoFin = id, s, e
			return nil
		},
		UpdateBusinessSubscriptionEndDateFn: func(ctx context.Context, businessID uint, endDate time.Time) error {
			vistoFinNegocio = endDate
			return nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).EditSubscriptionDates(context.Background(),
		dtos.EditSubscriptionDatesDTO{BusinessID: 26, StartDate: base, EndDate: fin}, 1)
	require.NoError(t, err)
	assert.Equal(t, uint(9), vistoID)
	assert.Equal(t, base, vistoInicio)
	assert.Equal(t, fin, vistoFin)
	assert.Equal(t, fin, vistoFinNegocio, "el fin del negocio queda alineado con el de la suscripcion")
	assert.Equal(t, base, got.StartDate)
	assert.Equal(t, fin, got.EndDate)
}

func TestEditSubscriptionDates_NoCreaRegistroDePago(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	creada := false
	wallet := &mocks.WalletDebiterMock{}
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return &entities.BusinessSubscription{ID: 9}, nil
		},
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			creada = true
			return nil
		},
	}

	_, err := newSubsUseCase(repo, wallet, nil).EditSubscriptionDates(context.Background(),
		dtos.EditSubscriptionDatesDTO{BusinessID: 26, StartDate: base, EndDate: base.AddDate(0, 1, 0)}, 1)
	require.NoError(t, err)
	assert.False(t, creada, "corregir fechas no es un pago nuevo")
	assert.Empty(t, wallet.DebitCalls)
}

func TestEditSubscriptionDates_ErroresDelRepo_SePropagan(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	dbErr := stderrors.New("db caida")

	t.Run("al leer la vigente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
				return nil, dbErr
			},
		}
		_, err := newSubsUseCase(repo, nil, nil).EditSubscriptionDates(context.Background(),
			dtos.EditSubscriptionDatesDTO{BusinessID: 26, StartDate: base, EndDate: base.AddDate(0, 1, 0)}, 1)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al actualizar la suscripcion", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
				return &entities.BusinessSubscription{ID: 9}, nil
			},
			UpdateSubscriptionDatesFn: func(ctx context.Context, id uint, s, e time.Time) error { return dbErr },
		}
		_, err := newSubsUseCase(repo, nil, nil).EditSubscriptionDates(context.Background(),
			dtos.EditSubscriptionDatesDTO{BusinessID: 26, StartDate: base, EndDate: base.AddDate(0, 1, 0)}, 1)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al actualizar el negocio", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
				return &entities.BusinessSubscription{ID: 9}, nil
			},
			UpdateBusinessSubscriptionEndDateFn: func(ctx context.Context, businessID uint, endDate time.Time) error {
				return dbErr
			},
		}
		_, err := newSubsUseCase(repo, nil, nil).EditSubscriptionDates(context.Background(),
			dtos.EditSubscriptionDatesDTO{BusinessID: 26, StartDate: base, EndDate: base.AddDate(0, 1, 0)}, 1)
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestDisableSubscription_MarcaCanceladoConFechaDeCorte(t *testing.T) {
	antes := time.Now()
	var vistoID uint
	var vistoStatus string
	var vistoFin *time.Time
	repo := &mocks.RepositoryMock{
		UpdateBusinessSubscriptionStatusFn: func(ctx context.Context, businessID uint, status string, endDate *time.Time) error {
			vistoID, vistoStatus, vistoFin = businessID, status, endDate
			return nil
		},
	}

	err := newSubsUseCase(repo, nil, nil).DisableSubscription(context.Background(), 26, 1)
	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoID)
	assert.Equal(t, entities.BusinessStatusCancelled, vistoStatus)
	require.NotNil(t, vistoFin)
	assert.False(t, vistoFin.Before(antes), "el corte es ahora, no una fecha futura")
}

func TestDisableSubscription_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		UpdateBusinessSubscriptionStatusFn: func(ctx context.Context, businessID uint, status string, endDate *time.Time) error {
			return dbErr
		},
	}

	err := newSubsUseCase(repo, nil, nil).DisableSubscription(context.Background(), 26, 1)
	assert.ErrorIs(t, err, dbErr)
}

func TestGetBusinessSubscription_DelegaAlRepositorio(t *testing.T) {
	esperada := &entities.BusinessSubscription{ID: 9, BusinessID: 26}
	var visto uint
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			visto = businessID
			return esperada, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).GetBusinessSubscription(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, uint(26), visto)
	assert.Equal(t, esperada, got)
}
