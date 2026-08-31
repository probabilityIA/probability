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
	"github.com/secamc93/probability/back/central/shared/moduleregistry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasModuleAccess_OverrideGanaSobreTodo(t *testing.T) {
	consultadoPlan := false
	repo := &mocks.RepositoryMock{
		ListOverridesByBusinessFn: func(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error) {
			return []entities.BusinessModuleOverride{{BusinessID: businessID, ModuleCode: "storefront"}}, nil
		},
		GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
			consultadoPlan = true
			return nil, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).HasModuleAccess(context.Background(), 26, "storefront")

	require.NoError(t, err)
	assert.True(t, got, "un override habilita incluso un modulo restringido por defecto")
	assert.False(t, consultadoPlan, "con override no hace falta mirar el plan")
}

func TestHasModuleAccess_ModulosRestringidosPorDefecto_SinOverrideNoSeVen(t *testing.T) {
	for _, code := range []string{"storefront", "tickets", "delivery"} {
		t.Run(code, func(t *testing.T) {
			repo := &mocks.RepositoryMock{
				GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
					return &entities.SubscriptionType{ID: id, ModuleCodes: []string{code}}, nil
				},
				GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
					return uintPtr(3), nil
				},
			}

			got, err := newSubsUseCase(repo, nil, nil).HasModuleAccess(context.Background(), 26, code)

			require.NoError(t, err)
			assert.False(t, got,
				"un modulo restringido por defecto no se abre ni aunque el plan lo liste: solo por override")
		})
	}
}

func TestHasModuleAccess_SinPlanAsignado_TodoLoNoRestringidoEsVisible(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
			return nil, nil
		},
	}
	uc := newSubsUseCase(repo, nil, nil)

	for _, code := range []string{"orders", "shipments", "inventory", "invoicing"} {
		got, err := uc.HasModuleAccess(context.Background(), 26, code)
		require.NoError(t, err)
		assert.True(t, got, "sin plan, %s queda abierto", code)
	}
}

func TestHasModuleAccess_PlanQueNoExiste_SeComportaComoSinPlan(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
			return uintPtr(999), nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return nil, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).HasModuleAccess(context.Background(), 26, "orders")

	require.NoError(t, err)
	assert.True(t, got)
}

func TestHasModuleAccess_ModuloEnElPlan_SeVe(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
			return uintPtr(3), nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, ModuleCodes: []string{"orders", "shipments"}}, nil
		},
	}
	uc := newSubsUseCase(repo, nil, nil)

	ok, err := uc.HasModuleAccess(context.Background(), 26, "orders")
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = uc.HasModuleAccess(context.Background(), 26, "inventory")
	require.NoError(t, err)
	assert.False(t, ok, "un modulo fuera del plan no se ve")
}

func TestHasModuleAccess_OverrideDeOtroModuloNoAbreEste(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListOverridesByBusinessFn: func(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error) {
			return []entities.BusinessModuleOverride{{ModuleCode: "tickets"}}, nil
		},
		GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
			return uintPtr(3), nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, ModuleCodes: []string{"orders"}}, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).HasModuleAccess(context.Background(), 26, "storefront")

	require.NoError(t, err)
	assert.False(t, got)
}

func TestHasModuleAccess_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al listar overrides", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			ListOverridesByBusinessFn: func(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error) {
				return nil, dbErr
			},
		}
		got, err := newSubsUseCase(repo, nil, nil).HasModuleAccess(context.Background(), 26, "orders")
		assert.ErrorIs(t, err, dbErr)
		assert.False(t, got, "ante un error no se concede acceso")
	})

	t.Run("al leer el plan del negocio", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
				return nil, dbErr
			},
		}
		got, err := newSubsUseCase(repo, nil, nil).HasModuleAccess(context.Background(), 26, "orders")
		assert.ErrorIs(t, err, dbErr)
		assert.False(t, got)
	})

	t.Run("al leer el plan", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
				return uintPtr(3), nil
			},
			GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
				return nil, dbErr
			},
		}
		got, err := newSubsUseCase(repo, nil, nil).HasModuleAccess(context.Background(), 26, "orders")
		assert.ErrorIs(t, err, dbErr)
		assert.False(t, got)
	})
}

func TestGetModuleCodes_DevuelveTodoElRegistro(t *testing.T) {
	got := newSubsUseCase(nil, nil, nil).GetModuleCodes()

	assert.Len(t, got, len(moduleregistry.All))
	for _, m := range moduleregistry.All {
		assert.Contains(t, got, string(m))
	}
}

func TestGetModuleCatalog_TraeCodigoYNombreVisible(t *testing.T) {
	got := newSubsUseCase(nil, nil, nil).GetModuleCatalog()

	require.Len(t, got, len(moduleregistry.All))
	porCodigo := map[string]string{}
	for _, m := range got {
		porCodigo[m.Code] = m.Name
		assert.NotEmpty(t, m.Name, "el modulo %s no tiene nombre visible", m.Code)
	}
	assert.Equal(t, "Ordenes", porCodigo["orders"])
	assert.Equal(t, "Facturacion", porCodigo["invoicing"])
}

func TestGetAccessibleModules_SinPlan_ExcluyeLosRestringidos(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
			return nil, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).GetAccessibleModules(context.Background(), 26)

	require.NoError(t, err)
	assert.NotContains(t, got, "storefront")
	assert.NotContains(t, got, "tickets")
	assert.NotContains(t, got, "delivery")
	assert.Contains(t, got, "orders")
	assert.Len(t, got, len(moduleregistry.All)-len(moduleregistry.RestrictedByDefault))
}

func TestGetAccessibleModules_ConPlanAcotado_SoloLosDelPlan(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
			return uintPtr(3), nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, ModuleCodes: []string{"orders", "inventory"}}, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).GetAccessibleModules(context.Background(), 26)

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"orders", "inventory"}, got)
}

func TestGetAccessibleModules_SumaLosOverrides(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListOverridesByBusinessFn: func(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error) {
			return []entities.BusinessModuleOverride{{ModuleCode: "tickets"}}, nil
		},
		GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
			return uintPtr(3), nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, ModuleCodes: []string{"orders"}}, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).GetAccessibleModules(context.Background(), 26)

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"orders", "tickets"}, got)
}

func TestGetAccessibleModules_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		ListOverridesByBusinessFn: func(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error) {
			return nil, dbErr
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).GetAccessibleModules(context.Background(), 26)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestEcommerceChannelLimit_SinPlan_EsCero(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
			return nil, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).EcommerceChannelLimit(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, 0, got, "sin plan no se puede conectar ningun canal")
}

func TestEcommerceChannelLimit_PlanInexistente_EsCero(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
			return uintPtr(999), nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return nil, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).EcommerceChannelLimit(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, 0, got)
}

func TestEcommerceChannelLimit_DevuelveElLimiteDelPlan(t *testing.T) {
	for _, limite := range []int{0, 1, 5, 100} {
		repo := &mocks.RepositoryMock{
			GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
				return uintPtr(3), nil
			},
			GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
				return &entities.SubscriptionType{ID: id, MaxEcommerceChannels: limite}, nil
			},
		}

		got, err := newSubsUseCase(repo, nil, nil).EcommerceChannelLimit(context.Background(), 26)

		require.NoError(t, err)
		assert.Equal(t, limite, got)
	}
}

func TestEcommerceChannelLimit_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al leer el plan del negocio", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
				return nil, dbErr
			},
		}
		got, err := newSubsUseCase(repo, nil, nil).EcommerceChannelLimit(context.Background(), 26)
		assert.ErrorIs(t, err, dbErr)
		assert.Zero(t, got)
	})

	t.Run("al leer el plan", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetBusinessCurrentSubscriptionTypeIDFn: func(ctx context.Context, businessID uint) (*uint, error) {
				return uintPtr(3), nil
			},
			GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
				return nil, dbErr
			},
		}
		got, err := newSubsUseCase(repo, nil, nil).EcommerceChannelLimit(context.Background(), 26)
		assert.ErrorIs(t, err, dbErr)
		assert.Zero(t, got)
	})
}

func TestGrantOverride_ModuloInvalido_Rechaza(t *testing.T) {
	casos := []string{"", "modulo_inventado", "ORDERS", "orders "}

	for _, code := range casos {
		t.Run(code, func(t *testing.T) {
			creado := false
			repo := &mocks.RepositoryMock{
				CreateOverrideFn: func(ctx context.Context, o *entities.BusinessModuleOverride) error {
					creado = true
					return nil
				},
			}

			err := newSubsUseCase(repo, nil, nil).GrantOverride(context.Background(),
				dtos.GrantOverrideDTO{BusinessID: 26, ModuleCode: code})

			assert.ErrorIs(t, err, errs.ErrInvalidModuleCode)
			assert.False(t, creado)
		})
	}
}

func TestGrantOverride_ModuloValido_GuardaQuienLoOtorgo(t *testing.T) {
	notas := "cliente piloto"
	var creado *entities.BusinessModuleOverride
	repo := &mocks.RepositoryMock{
		CreateOverrideFn: func(ctx context.Context, o *entities.BusinessModuleOverride) error {
			creado = o
			return nil
		},
	}

	err := newSubsUseCase(repo, nil, nil).GrantOverride(context.Background(), dtos.GrantOverrideDTO{
		BusinessID: 26, ModuleCode: "storefront", GrantedByUserID: 42, Notes: &notas,
	})

	require.NoError(t, err)
	require.NotNil(t, creado)
	assert.Equal(t, uint(26), creado.BusinessID)
	assert.Equal(t, "storefront", creado.ModuleCode)
	assert.Equal(t, uint(42), creado.GrantedByUserID, "queda trazabilidad de quien abrio el modulo")
	require.NotNil(t, creado.Notes)
	assert.Equal(t, notas, *creado.Notes)
}

func TestGrantOverride_AceptaTodosLosModulosDelRegistro(t *testing.T) {
	uc := newSubsUseCase(nil, nil, nil)

	for _, m := range moduleregistry.All {
		err := uc.GrantOverride(context.Background(), dtos.GrantOverrideDTO{
			BusinessID: 26, ModuleCode: string(m),
		})
		assert.NoError(t, err, "modulo %s del registro debe ser otorgable", m)
	}
}

func TestGrantOverride_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("duplicate key")
	repo := &mocks.RepositoryMock{
		CreateOverrideFn: func(ctx context.Context, o *entities.BusinessModuleOverride) error { return dbErr },
	}

	err := newSubsUseCase(repo, nil, nil).GrantOverride(context.Background(),
		dtos.GrantOverrideDTO{BusinessID: 26, ModuleCode: "orders"})

	assert.ErrorIs(t, err, dbErr)
}

func TestRevokeOverride_PasaNegocioYModulo(t *testing.T) {
	var vistoID uint
	var vistoCode string
	repo := &mocks.RepositoryMock{
		DeleteOverrideFn: func(ctx context.Context, businessID uint, moduleCode string) error {
			vistoID, vistoCode = businessID, moduleCode
			return nil
		},
	}

	err := newSubsUseCase(repo, nil, nil).RevokeOverride(context.Background(), 26, "storefront", 1)

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoID)
	assert.Equal(t, "storefront", vistoCode)
}

func TestRevokeOverride_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("no existe")
	repo := &mocks.RepositoryMock{
		DeleteOverrideFn: func(ctx context.Context, businessID uint, moduleCode string) error { return dbErr },
	}

	err := newSubsUseCase(repo, nil, nil).RevokeOverride(context.Background(), 26, "storefront", 1)

	assert.ErrorIs(t, err, dbErr)
}

func TestListOverrides_DelegaAlRepositorio(t *testing.T) {
	esperados := []entities.BusinessModuleOverride{{ID: 1, ModuleCode: "storefront"}}
	var visto uint
	repo := &mocks.RepositoryMock{
		ListOverridesByBusinessFn: func(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error) {
			visto = businessID
			return esperados, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).ListOverrides(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, uint(26), visto)
	assert.Equal(t, esperados, got)
}

func TestCreateSubscriptionType_CamposObligatorios(t *testing.T) {
	casos := []struct {
		nombre string
		dto    dtos.CreateSubscriptionTypeDTO
	}{
		{"sin nombre", dtos.CreateSubscriptionTypeDTO{Code: "pro", Price: 1000}},
		{"sin codigo", dtos.CreateSubscriptionTypeDTO{Name: "Pro", Price: 1000}},
		{"precio negativo", dtos.CreateSubscriptionTypeDTO{Name: "Pro", Code: "pro", Price: -1}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			creado := false
			repo := &mocks.RepositoryMock{
				CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
					creado = true
					return nil
				},
			}

			_, err := newSubsUseCase(repo, nil, nil).CreateSubscriptionType(context.Background(), tc.dto)

			assert.ErrorIs(t, err, errs.ErrInvalidSubscriptionType)
			assert.False(t, creado)
		})
	}
}

func TestCreateSubscriptionType_ModuloInvalido_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.RepositoryMock{
		CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			creado = true
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).CreateSubscriptionType(context.Background(),
		dtos.CreateSubscriptionTypeDTO{
			Name: "Pro", Code: "pro", Price: 1000,
			ModuleCodes: []string{"orders", "modulo_inventado"},
		})

	assert.ErrorIs(t, err, errs.ErrInvalidModuleCode)
	assert.False(t, creado, "un plan con un modulo inexistente dejaria acceso indefinido")
}

func TestCreateSubscriptionType_BillingPeriodVacio_CaeAMonthly(t *testing.T) {
	var creado *entities.SubscriptionType
	repo := &mocks.RepositoryMock{
		CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			creado = s
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).CreateSubscriptionType(context.Background(),
		dtos.CreateSubscriptionTypeDTO{Name: "Pro", Code: "pro", Price: 1000})

	require.NoError(t, err)
	assert.Equal(t, "monthly", creado.BillingPeriod)
}

func TestCreateSubscriptionType_BillingPeriodExplicito_SeRespeta(t *testing.T) {
	var creado *entities.SubscriptionType
	repo := &mocks.RepositoryMock{
		CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			creado = s
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).CreateSubscriptionType(context.Background(),
		dtos.CreateSubscriptionTypeDTO{Name: "Pro", Code: "pro", Price: 1000, BillingPeriod: "yearly"})

	require.NoError(t, err)
	assert.Equal(t, "yearly", creado.BillingPeriod)
}

func TestCreateSubscriptionType_NaceActivoYSinNegocioDueno(t *testing.T) {
	var creado *entities.SubscriptionType
	repo := &mocks.RepositoryMock{
		CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			creado = s
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).CreateSubscriptionType(context.Background(),
		dtos.CreateSubscriptionTypeDTO{
			Name: "Pro", Code: "pro", Price: 1000,
			ModuleCodes: []string{"orders"}, MaxEcommerceChannels: 5,
		})

	require.NoError(t, err)
	assert.True(t, creado.Active)
	assert.Nil(t, creado.BusinessID, "un plan del catalogo no pertenece a ningun negocio")
	assert.Equal(t, 5, creado.MaxEcommerceChannels)
}

func TestCreateSubscriptionType_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("duplicate code")
	repo := &mocks.RepositoryMock{
		CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error { return dbErr },
	}

	got, err := newSubsUseCase(repo, nil, nil).CreateSubscriptionType(context.Background(),
		dtos.CreateSubscriptionTypeDTO{Name: "Pro", Code: "pro", Price: 1000})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateCustomPlan_ValidacionesDeEntrada(t *testing.T) {
	casos := []struct {
		nombre  string
		dto     dtos.CreateCustomPlanDTO
		wantErr error
	}{
		{"sin negocio", dtos.CreateCustomPlanDTO{Months: 1, Name: "X", Code: "x"}, errs.ErrBusinessRequired},
		{"meses cero", dtos.CreateCustomPlanDTO{BusinessID: 26, Months: 0, Name: "X", Code: "x"}, errs.ErrInvalidMonths},
		{"meses negativos", dtos.CreateCustomPlanDTO{BusinessID: 26, Months: -1, Name: "X", Code: "x"}, errs.ErrInvalidMonths},
		{"sin nombre", dtos.CreateCustomPlanDTO{BusinessID: 26, Months: 1, Code: "x"}, errs.ErrInvalidSubscriptionType},
		{"sin codigo", dtos.CreateCustomPlanDTO{BusinessID: 26, Months: 1, Name: "X"}, errs.ErrInvalidSubscriptionType},
		{"precio negativo", dtos.CreateCustomPlanDTO{BusinessID: 26, Months: 1, Name: "X", Code: "x", Price: -1}, errs.ErrInvalidSubscriptionType},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			creado := false
			repo := &mocks.RepositoryMock{
				CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
					creado = true
					return nil
				},
			}

			_, err := newSubsUseCase(repo, nil, nil).CreateCustomPlan(context.Background(), tc.dto, 1)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.False(t, creado)
		})
	}
}

func TestCreateCustomPlan_PrecioCeroEsValido(t *testing.T) {
	repo := &mocks.RepositoryMock{
		CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			s.ID = 9
			return nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).CreateCustomPlan(context.Background(),
		dtos.CreateCustomPlanDTO{BusinessID: 26, Months: 1, Name: "Cortesia", Code: "cortesia", Price: 0}, 1)

	require.NoError(t, err, "un plan de cortesia a precio cero es valido, a diferencia del catalogo general")
	assert.NotNil(t, got)
}

func TestCreateCustomPlan_ModuloInvalido_Rechaza(t *testing.T) {
	_, err := newSubsUseCase(nil, nil, nil).CreateCustomPlan(context.Background(),
		dtos.CreateCustomPlanDTO{
			BusinessID: 26, Months: 1, Name: "X", Code: "x", Price: 100,
			ModuleCodes: []string{"modulo_inventado"},
		}, 1)

	assert.ErrorIs(t, err, errs.ErrInvalidModuleCode)
}

func TestCreateCustomPlan_QuedaAmarradoAlNegocio(t *testing.T) {
	var creado *entities.SubscriptionType
	repo := &mocks.RepositoryMock{
		CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			creado = s
			s.ID = 9
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).CreateCustomPlan(context.Background(),
		dtos.CreateCustomPlanDTO{BusinessID: 26, Months: 3, Name: "A Medida", Code: "medida", Price: 500}, 1)

	require.NoError(t, err)
	require.NotNil(t, creado)
	require.NotNil(t, creado.BusinessID)
	assert.Equal(t, uint(26), *creado.BusinessID,
		"un plan a medida pertenece a un solo negocio, no aparece en el catalogo general")
	assert.True(t, creado.Active)
}

func TestCreateCustomPlan_RegistraElPagoDelPlanRecienCreado(t *testing.T) {
	ref := "TRANSF-999"
	var suscripcion *entities.BusinessSubscription
	repo := &mocks.RepositoryMock{
		CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			s.ID = 77
			return nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, Code: "medida", Price: 500, Active: true}, nil
		},
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			suscripcion = s
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).CreateCustomPlan(context.Background(),
		dtos.CreateCustomPlanDTO{
			BusinessID: 26, Months: 3, Name: "A Medida", Code: "medida", Price: 500,
			PaymentReference: &ref,
		}, 1)

	require.NoError(t, err)
	require.NotNil(t, suscripcion)
	assert.Equal(t, uint(77), suscripcion.SubscriptionTypeID, "el pago apunta al plan recien creado")
	assert.Equal(t, 3, suscripcion.Months)
	assert.InDelta(t, 1500.0, suscripcion.Amount, 0.001)
	require.NotNil(t, suscripcion.PaymentReference)
	assert.Equal(t, ref, *suscripcion.PaymentReference)
}

func TestCreateCustomPlan_NoDebitaDeLaBilletera(t *testing.T) {
	wallet := &mocks.WalletDebiterMock{}

	_, err := newSubsUseCase(nil, wallet, nil).CreateCustomPlan(context.Background(),
		dtos.CreateCustomPlanDTO{BusinessID: 26, Months: 1, Name: "X", Code: "x", Price: 100}, 1)

	require.NoError(t, err)
	assert.Empty(t, wallet.DebitCalls)
}

func TestCreateCustomPlan_FallaElRegistroDelPago_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			return dbErr
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).CreateCustomPlan(context.Background(),
		dtos.CreateCustomPlanDTO{BusinessID: 26, Months: 1, Name: "X", Code: "x", Price: 100}, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got, "el plan quedo creado pero sin suscripcion, se reporta el fallo")
}

func TestCreateCustomPlan_FallaLaCreacionDelPlan_NoRegistraPago(t *testing.T) {
	dbErr := stderrors.New("duplicate code")
	pagoRegistrado := false
	repo := &mocks.RepositoryMock{
		CreateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error { return dbErr },
		CreateSubscriptionAndActivateFn: func(ctx context.Context, s *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			pagoRegistrado = true
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).CreateCustomPlan(context.Background(),
		dtos.CreateCustomPlanDTO{BusinessID: 26, Months: 1, Name: "X", Code: "x", Price: 100}, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, pagoRegistrado)
}

func TestListCustomPlans_PasaElFiltroDeNegocio(t *testing.T) {
	var visto *uint
	repo := &mocks.RepositoryMock{
		ListCustomPlansFn: func(ctx context.Context, businessID *uint) ([]entities.SubscriptionType, error) {
			visto = businessID
			return []entities.SubscriptionType{{ID: 1}}, nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).ListCustomPlans(context.Background(), uintPtr(26))

	require.NoError(t, err)
	require.NotNil(t, visto)
	assert.Equal(t, uint(26), *visto)
	assert.Len(t, got, 1)
}

func TestUpdateSubscriptionType_PlanInexistente_Rechaza(t *testing.T) {
	actualizado := false
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return nil, nil
		},
		UpdateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			actualizado = true
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).UpdateSubscriptionType(context.Background(),
		dtos.UpdateSubscriptionTypeDTO{ID: 404, Name: "X", Price: 100})

	assert.ErrorIs(t, err, errs.ErrSubscriptionTypeNotFound)
	assert.False(t, actualizado)
}

func TestUpdateSubscriptionType_CamposObligatorios(t *testing.T) {
	casos := []struct {
		nombre string
		dto    dtos.UpdateSubscriptionTypeDTO
	}{
		{"sin nombre", dtos.UpdateSubscriptionTypeDTO{ID: 3, Price: 100}},
		{"precio negativo", dtos.UpdateSubscriptionTypeDTO{ID: 3, Name: "X", Price: -5}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			_, err := newSubsUseCase(nil, nil, nil).UpdateSubscriptionType(context.Background(), tc.dto)
			assert.ErrorIs(t, err, errs.ErrInvalidSubscriptionType)
		})
	}
}

func TestUpdateSubscriptionType_ModuloInvalido_Rechaza(t *testing.T) {
	_, err := newSubsUseCase(nil, nil, nil).UpdateSubscriptionType(context.Background(),
		dtos.UpdateSubscriptionTypeDTO{ID: 3, Name: "X", Price: 100, ModuleCodes: []string{"inventado"}})

	assert.ErrorIs(t, err, errs.ErrInvalidModuleCode)
}

func TestUpdateSubscriptionType_NoPermiteCambiarElCodigo(t *testing.T) {
	var guardado *entities.SubscriptionType
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, Code: "pro", Name: "Pro", Price: 1000}, nil
		},
		UpdateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			guardado = s
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).UpdateSubscriptionType(context.Background(),
		dtos.UpdateSubscriptionTypeDTO{ID: 3, Name: "Pro Plus", Price: 2000})

	require.NoError(t, err)
	assert.Equal(t, "pro", guardado.Code, "el codigo es la llave del plan, no se toca en update")
	assert.Equal(t, "Pro Plus", guardado.Name)
	assert.InDelta(t, 2000.0, guardado.Price, 0.001)
}

func TestUpdateSubscriptionType_PuedeDesactivarElPlan(t *testing.T) {
	var guardado *entities.SubscriptionType
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, Code: "pro", Active: true}, nil
		},
		UpdateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			guardado = s
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).UpdateSubscriptionType(context.Background(),
		dtos.UpdateSubscriptionTypeDTO{ID: 3, Name: "Pro", Price: 1000, Active: false})

	require.NoError(t, err)
	assert.False(t, guardado.Active)
}

func TestUpdateSubscriptionType_BillingPeriodVacio_CaeAMonthly(t *testing.T) {
	var guardado *entities.SubscriptionType
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, BillingPeriod: "yearly"}, nil
		},
		UpdateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error {
			guardado = s
			return nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).UpdateSubscriptionType(context.Background(),
		dtos.UpdateSubscriptionTypeDTO{ID: 3, Name: "X", Price: 100})

	require.NoError(t, err)
	assert.Equal(t, "monthly", guardado.BillingPeriod)
}

func TestUpdateSubscriptionType_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al leer", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
				return nil, dbErr
			},
		}
		_, err := newSubsUseCase(repo, nil, nil).UpdateSubscriptionType(context.Background(),
			dtos.UpdateSubscriptionTypeDTO{ID: 3, Name: "X", Price: 100})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al guardar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			UpdateSubscriptionTypeFn: func(ctx context.Context, s *entities.SubscriptionType) error { return dbErr },
		}
		_, err := newSubsUseCase(repo, nil, nil).UpdateSubscriptionType(context.Background(),
			dtos.UpdateSubscriptionTypeDTO{ID: 3, Name: "X", Price: 100})
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestGetYListYDeleteSubscriptionType_DeleganAlRepositorio(t *testing.T) {
	var borradoID uint
	var vistoActiveOnly bool
	esperado := &entities.SubscriptionType{ID: 3, Name: "Pro"}
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return esperado, nil
		},
		ListSubscriptionTypesFn: func(ctx context.Context, activeOnly bool) ([]entities.SubscriptionType, error) {
			vistoActiveOnly = activeOnly
			return []entities.SubscriptionType{*esperado}, nil
		},
		DeleteSubscriptionTypeFn: func(ctx context.Context, id uint) error {
			borradoID = id
			return nil
		},
	}
	uc := newSubsUseCase(repo, nil, nil)

	got, err := uc.GetSubscriptionType(context.Background(), 3)
	require.NoError(t, err)
	assert.Equal(t, esperado, got)

	lista, err := uc.ListSubscriptionTypes(context.Background(), true)
	require.NoError(t, err)
	assert.True(t, vistoActiveOnly)
	assert.Len(t, lista, 1)

	require.NoError(t, uc.DeleteSubscriptionType(context.Background(), 3))
	assert.Equal(t, uint(3), borradoID)
}

func TestCheckExpiringSubscriptions_VencidosQuedanExpired(t *testing.T) {
	var marcados []uint
	repo := &mocks.RepositoryMock{
		ListBusinessesJustExpiredFn: func(ctx context.Context, before time.Time) ([]entities.ExpiringBusiness, error) {
			return []entities.ExpiringBusiness{{BusinessID: 26}}, nil
		},
		MarkExpiredIfStillActiveFn: func(ctx context.Context, businessID uint, before time.Time) error {
			marcados = append(marcados, businessID)
			return nil
		},
	}

	err := newSubsUseCase(repo, nil, nil).CheckExpiringSubscriptions(context.Background())

	require.NoError(t, err)
	assert.Equal(t, []uint{26}, marcados)
}

func TestCheckExpiringSubscriptions_NoCreaAnuncios(t *testing.T) {
	ann := &mocks.AnnouncementsGatewayMock{}
	repo := &mocks.RepositoryMock{
		ListBusinessesJustExpiredFn: func(ctx context.Context, before time.Time) ([]entities.ExpiringBusiness, error) {
			return []entities.ExpiringBusiness{{BusinessID: 26}}, nil
		},
	}

	err := newSubsUseCase(repo, nil, ann).CheckExpiringSubscriptions(context.Background())

	require.NoError(t, err)
	assert.Empty(t, ann.CreatedAlerts, "el aviso de vencimiento se elimino: nunca debe crear anuncios")
}

func TestCheckExpiringSubscriptions_FalloAlMarcarExpirado_NoAbortaElResto(t *testing.T) {
	intentos := 0
	repo := &mocks.RepositoryMock{
		ListBusinessesJustExpiredFn: func(ctx context.Context, before time.Time) ([]entities.ExpiringBusiness, error) {
			return []entities.ExpiringBusiness{{BusinessID: 26}, {BusinessID: 27}}, nil
		},
		MarkExpiredIfStillActiveFn: func(ctx context.Context, businessID uint, before time.Time) error {
			intentos++
			return stderrors.New("db caida")
		},
	}

	err := newSubsUseCase(repo, nil, nil).CheckExpiringSubscriptions(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 2, intentos)
}

func TestCheckExpiringSubscriptions_ErrorAlListarVencidos_AbortaElBarrido(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := &mocks.RepositoryMock{
		ListBusinessesJustExpiredFn: func(ctx context.Context, before time.Time) ([]entities.ExpiringBusiness, error) {
			return nil, dbErr
		},
	}
	err := newSubsUseCase(repo, nil, nil).CheckExpiringSubscriptions(context.Background())
	assert.ErrorIs(t, err, dbErr)
}

func TestCheckExpiringSubscriptions_SinVencidos_NoMarcaNada(t *testing.T) {
	err := newSubsUseCase(nil, nil, nil).CheckExpiringSubscriptions(context.Background())

	require.NoError(t, err)
}
