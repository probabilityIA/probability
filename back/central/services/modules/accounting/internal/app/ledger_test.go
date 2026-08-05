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

func TestCreateConcept_NormalizaCodigoAMayusculasYRecorta(t *testing.T) {
	var vistoExists string
	var creado dtos.CreateConceptDTO
	repo := &mocks.RepositoryMock{
		ConceptCodeExistsFn: func(ctx context.Context, code string, excludeID *uint) (bool, error) {
			vistoExists = code
			return false, nil
		},
		CreateConceptFn: func(ctx context.Context, dto dtos.CreateConceptDTO) (*entities.Concept, error) {
			creado = dto
			return &entities.Concept{ID: 1, Code: dto.Code}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateConcept(context.Background(),
		dtos.CreateConceptDTO{Code: "  margen_guia  ", Name: "  Margen de guia  ", Kind: entities.KindIncome})

	require.NoError(t, err)
	assert.Equal(t, "MARGEN_GUIA", vistoExists)
	assert.Equal(t, "MARGEN_GUIA", creado.Code)
	assert.Equal(t, "Margen de guia", creado.Name)
}

func TestCreateConcept_CamposObligatorios(t *testing.T) {
	casos := []struct {
		nombre  string
		dto     dtos.CreateConceptDTO
		wantErr error
	}{
		{"sin codigo", dtos.CreateConceptDTO{Name: "X", Kind: entities.KindIncome}, domainerrors.ErrCodeRequired},
		{"codigo en blanco", dtos.CreateConceptDTO{Code: "   ", Name: "X", Kind: entities.KindIncome}, domainerrors.ErrCodeRequired},
		{"sin nombre", dtos.CreateConceptDTO{Code: "A", Kind: entities.KindIncome}, domainerrors.ErrNameRequired},
		{"nombre en blanco", dtos.CreateConceptDTO{Code: "A", Name: "  ", Kind: entities.KindIncome}, domainerrors.ErrNameRequired},
		{"kind vacio", dtos.CreateConceptDTO{Code: "A", Name: "X"}, domainerrors.ErrInvalidKind},
		{"kind invalido", dtos.CreateConceptDTO{Code: "A", Name: "X", Kind: "OTRO"}, domainerrors.ErrInvalidKind},
		{"kind en minusculas", dtos.CreateConceptDTO{Code: "A", Name: "X", Kind: "income"}, domainerrors.ErrInvalidKind},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			creado := false
			repo := &mocks.RepositoryMock{
				CreateConceptFn: func(ctx context.Context, dto dtos.CreateConceptDTO) (*entities.Concept, error) {
					creado = true
					return nil, nil
				},
			}

			_, err := newAccountingUseCase(repo, nil, nil, nil).CreateConcept(context.Background(), tc.dto)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.False(t, creado)
		})
	}
}

func TestCreateConcept_AceptaAmbosKinds(t *testing.T) {
	for _, kind := range []string{entities.KindIncome, entities.KindExpense} {
		got, err := newAccountingUseCase(nil, nil, nil, nil).CreateConcept(context.Background(),
			dtos.CreateConceptDTO{Code: "A", Name: "X", Kind: kind})
		require.NoError(t, err, "kind=%s", kind)
		assert.NotNil(t, got)
	}
}

func TestCreateConcept_CodigoDuplicado_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.RepositoryMock{
		ConceptCodeExistsFn: func(ctx context.Context, code string, excludeID *uint) (bool, error) {
			return true, nil
		},
		CreateConceptFn: func(ctx context.Context, dto dtos.CreateConceptDTO) (*entities.Concept, error) {
			creado = true
			return nil, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateConcept(context.Background(),
		dtos.CreateConceptDTO{Code: "A", Name: "X", Kind: entities.KindIncome})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateCode)
	assert.False(t, creado)
}

func TestCreateConcept_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("timeout")

	t.Run("al verificar duplicado", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			ConceptCodeExistsFn: func(ctx context.Context, code string, excludeID *uint) (bool, error) {
				return false, dbErr
			},
		}
		_, err := newAccountingUseCase(repo, nil, nil, nil).CreateConcept(context.Background(),
			dtos.CreateConceptDTO{Code: "A", Name: "X", Kind: entities.KindIncome})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al crear", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			CreateConceptFn: func(ctx context.Context, dto dtos.CreateConceptDTO) (*entities.Concept, error) {
				return nil, dbErr
			},
		}
		_, err := newAccountingUseCase(repo, nil, nil, nil).CreateConcept(context.Background(),
			dtos.CreateConceptDTO{Code: "A", Name: "X", Kind: entities.KindIncome})
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestUpdateConcept_NombreRequerido(t *testing.T) {
	for _, name := range []string{"", "   "} {
		actualizado := false
		repo := &mocks.RepositoryMock{
			UpdateConceptFn: func(ctx context.Context, dto dtos.UpdateConceptDTO) (*entities.Concept, error) {
				actualizado = true
				return nil, nil
			},
		}

		_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateConcept(context.Background(),
			dtos.UpdateConceptDTO{ID: 1, Name: name})

		assert.ErrorIs(t, err, domainerrors.ErrNameRequired)
		assert.False(t, actualizado)
	}
}

func TestUpdateConcept_ConceptoInexistente_NoActualiza(t *testing.T) {
	actualizado := false
	repo := &mocks.RepositoryMock{
		GetConceptByIDFn: func(ctx context.Context, id uint) (*entities.Concept, error) {
			return nil, domainerrors.ErrConceptNotFound
		},
		UpdateConceptFn: func(ctx context.Context, dto dtos.UpdateConceptDTO) (*entities.Concept, error) {
			actualizado = true
			return nil, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateConcept(context.Background(),
		dtos.UpdateConceptDTO{ID: 404, Name: "X"})

	assert.ErrorIs(t, err, domainerrors.ErrConceptNotFound)
	assert.False(t, actualizado)
}

func TestUpdateConcept_NoPermiteCambiarElCodigo(t *testing.T) {
	var recibido dtos.UpdateConceptDTO
	repo := &mocks.RepositoryMock{
		UpdateConceptFn: func(ctx context.Context, dto dtos.UpdateConceptDTO) (*entities.Concept, error) {
			recibido = dto
			return &entities.Concept{ID: dto.ID}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateConcept(context.Background(),
		dtos.UpdateConceptDTO{ID: 1, Name: "  Nuevo  "})

	require.NoError(t, err)
	assert.Equal(t, "Nuevo", recibido.Name)
	assert.Equal(t, uint(1), recibido.ID)
}

func TestListConcepts_DelegaYPropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListConceptsFn: func(ctx context.Context) ([]entities.Concept, error) {
			return []entities.Concept{{ID: 1, Code: "A"}}, nil
		},
	}
	got, err := newAccountingUseCase(repo, nil, nil, nil).ListConcepts(context.Background())
	require.NoError(t, err)
	assert.Len(t, got, 1)

	dbErr := stderrors.New("timeout")
	repoErr := &mocks.RepositoryMock{
		ListConceptsFn: func(ctx context.Context) ([]entities.Concept, error) { return nil, dbErr },
	}
	_, err = newAccountingUseCase(repoErr, nil, nil, nil).ListConcepts(context.Background())
	assert.ErrorIs(t, err, dbErr)
}

func TestSetConceptTax_ImpuestoInactivo_NoSePuedeActivar(t *testing.T) {
	asociado := false
	repo := &mocks.RepositoryMock{
		GetTaxByIDFn: func(ctx context.Context, id uint) (*entities.Tax, error) {
			return &entities.Tax{ID: id, IsActive: false}, nil
		},
		SetConceptTaxFn: func(ctx context.Context, dto dtos.SetConceptTaxDTO) error {
			asociado = true
			return nil
		},
	}

	err := newAccountingUseCase(repo, nil, nil, nil).SetConceptTax(context.Background(),
		dtos.SetConceptTaxDTO{ConceptID: 1, TaxID: 2, IsActive: true})

	assert.ErrorIs(t, err, domainerrors.ErrConceptTaxNotAllowed)
	assert.False(t, asociado)
}

func TestSetConceptTax_DesactivarUnImpuestoInactivoSiSePuede(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetTaxByIDFn: func(ctx context.Context, id uint) (*entities.Tax, error) {
			return &entities.Tax{ID: id, IsActive: false}, nil
		},
	}

	err := newAccountingUseCase(repo, nil, nil, nil).SetConceptTax(context.Background(),
		dtos.SetConceptTaxDTO{ConceptID: 1, TaxID: 2, IsActive: false})

	require.NoError(t, err, "quitar la asociacion de un impuesto ya inactivo debe ser posible")
}

func TestSetConceptTax_ErroresDeLectura_SePropagan(t *testing.T) {
	t.Run("concepto inexistente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetConceptByIDFn: func(ctx context.Context, id uint) (*entities.Concept, error) {
				return nil, domainerrors.ErrConceptNotFound
			},
		}
		err := newAccountingUseCase(repo, nil, nil, nil).SetConceptTax(context.Background(),
			dtos.SetConceptTaxDTO{ConceptID: 404, TaxID: 1})
		assert.ErrorIs(t, err, domainerrors.ErrConceptNotFound)
	})

	t.Run("impuesto inexistente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetTaxByIDFn: func(ctx context.Context, id uint) (*entities.Tax, error) {
				return nil, domainerrors.ErrTaxNotFound
			},
		}
		err := newAccountingUseCase(repo, nil, nil, nil).SetConceptTax(context.Background(),
			dtos.SetConceptTaxDTO{ConceptID: 1, TaxID: 404})
		assert.ErrorIs(t, err, domainerrors.ErrTaxNotFound)
	})
}

func TestCreateTax_ValidacionesDeEntrada(t *testing.T) {
	casos := []struct {
		nombre  string
		dto     dtos.CreateTaxDTO
		wantErr error
	}{
		{"sin codigo", dtos.CreateTaxDTO{Name: "X", RatePercent: 19}, domainerrors.ErrCodeRequired},
		{"sin nombre", dtos.CreateTaxDTO{Code: "IVA", RatePercent: 19}, domainerrors.ErrNameRequired},
		{"tarifa negativa", dtos.CreateTaxDTO{Code: "IVA", Name: "X", RatePercent: -1}, domainerrors.ErrInvalidRate},
		{"tarifa sobre cien", dtos.CreateTaxDTO{Code: "IVA", Name: "X", RatePercent: 100.1}, domainerrors.ErrInvalidRate},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			_, err := newAccountingUseCase(nil, nil, nil, nil).CreateTax(context.Background(), tc.dto)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestCreateTax_LimitesExactosSeAceptan(t *testing.T) {
	for _, rate := range []float64{0, 100} {
		got, err := newAccountingUseCase(nil, nil, nil, nil).CreateTax(context.Background(),
			dtos.CreateTaxDTO{Code: "X", Name: "X", RatePercent: rate})
		require.NoError(t, err, "rate=%v", rate)
		assert.NotNil(t, got)
	}
}

func TestCreateTax_NormalizaElCodigo(t *testing.T) {
	var visto string
	repo := &mocks.RepositoryMock{
		TaxCodeExistsFn: func(ctx context.Context, code string, excludeID *uint) (bool, error) {
			visto = code
			return false, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateTax(context.Background(),
		dtos.CreateTaxDTO{Code: "  iva  ", Name: "IVA", RatePercent: 19})

	require.NoError(t, err)
	assert.Equal(t, "IVA", visto)
}

func TestCreateTax_CodigoDuplicado_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		TaxCodeExistsFn: func(ctx context.Context, code string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateTax(context.Background(),
		dtos.CreateTaxDTO{Code: "IVA", Name: "IVA", RatePercent: 19})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateCode)
}

func TestCreateTax_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("timeout")

	repo := &mocks.RepositoryMock{
		TaxCodeExistsFn: func(ctx context.Context, code string, excludeID *uint) (bool, error) {
			return false, dbErr
		},
	}
	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateTax(context.Background(),
		dtos.CreateTaxDTO{Code: "IVA", Name: "IVA", RatePercent: 19})
	assert.ErrorIs(t, err, dbErr)

	repo2 := &mocks.RepositoryMock{
		CreateTaxFn: func(ctx context.Context, dto dtos.CreateTaxDTO) (*entities.Tax, error) { return nil, dbErr },
	}
	_, err = newAccountingUseCase(repo2, nil, nil, nil).CreateTax(context.Background(),
		dtos.CreateTaxDTO{Code: "IVA", Name: "IVA", RatePercent: 19})
	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateTax_ValidacionesYPropagacion(t *testing.T) {
	t.Run("sin nombre", func(t *testing.T) {
		_, err := newAccountingUseCase(nil, nil, nil, nil).UpdateTax(context.Background(),
			dtos.UpdateTaxDTO{ID: 1, Name: "  ", RatePercent: 19})
		assert.ErrorIs(t, err, domainerrors.ErrNameRequired)
	})

	t.Run("tarifa fuera de rango", func(t *testing.T) {
		for _, rate := range []float64{-0.1, 100.1} {
			_, err := newAccountingUseCase(nil, nil, nil, nil).UpdateTax(context.Background(),
				dtos.UpdateTaxDTO{ID: 1, Name: "X", RatePercent: rate})
			assert.ErrorIs(t, err, domainerrors.ErrInvalidRate)
		}
	})

	t.Run("impuesto inexistente", func(t *testing.T) {
		actualizado := false
		repo := &mocks.RepositoryMock{
			GetTaxByIDFn: func(ctx context.Context, id uint) (*entities.Tax, error) {
				return nil, domainerrors.ErrTaxNotFound
			},
			UpdateTaxFn: func(ctx context.Context, dto dtos.UpdateTaxDTO) (*entities.Tax, error) {
				actualizado = true
				return nil, nil
			},
		}
		_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateTax(context.Background(),
			dtos.UpdateTaxDTO{ID: 404, Name: "X", RatePercent: 19})
		assert.ErrorIs(t, err, domainerrors.ErrTaxNotFound)
		assert.False(t, actualizado)
	})
}

func TestUpdateTax_Exitoso_RecortaElNombre(t *testing.T) {
	var recibido dtos.UpdateTaxDTO
	repo := &mocks.RepositoryMock{
		UpdateTaxFn: func(ctx context.Context, dto dtos.UpdateTaxDTO) (*entities.Tax, error) {
			recibido = dto
			return &entities.Tax{ID: dto.ID}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateTax(context.Background(),
		dtos.UpdateTaxDTO{ID: 1, Name: "  IVA  ", RatePercent: 19})

	require.NoError(t, err)
	assert.Equal(t, "IVA", recibido.Name)
}

func TestListTaxes_DelegaYPropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListTaxesFn: func(ctx context.Context) ([]entities.Tax, error) {
			return []entities.Tax{{ID: 1, Code: "IVA"}}, nil
		},
	}
	got, err := newAccountingUseCase(repo, nil, nil, nil).ListTaxes(context.Background())
	require.NoError(t, err)
	assert.Len(t, got, 1)

	dbErr := stderrors.New("timeout")
	repoErr := &mocks.RepositoryMock{
		ListTaxesFn: func(ctx context.Context) ([]entities.Tax, error) { return nil, dbErr },
	}
	_, err = newAccountingUseCase(repoErr, nil, nil, nil).ListTaxes(context.Background())
	assert.ErrorIs(t, err, dbErr)
}

func TestApplyTaxes_CalculaCadaImpuestoSobreElMonto(t *testing.T) {
	entry := &entities.Entry{Amount: 100000}
	concept := &entities.Concept{
		Taxes: []entities.ConceptTax{
			{TaxCode: "IVA", TaxName: "IVA", RatePercent: 19, IsActive: true},
			{TaxCode: "RETEICA", TaxName: "ReteICA", RatePercent: 0.966, IsActive: true},
		},
	}

	applyTaxes(entry, concept)

	require.Len(t, entry.TaxDetail, 2)
	assert.InDelta(t, 19000.0, entry.TaxDetail[0].Amount, 0.001)
	assert.InDelta(t, 966.0, entry.TaxDetail[1].Amount, 0.001)
	assert.InDelta(t, 19966.0, entry.TaxTotal, 0.001)
}

func TestApplyTaxes_IgnoraLosInactivos(t *testing.T) {
	entry := &entities.Entry{Amount: 100000}
	concept := &entities.Concept{
		Taxes: []entities.ConceptTax{
			{TaxCode: "IVA", RatePercent: 19, IsActive: true},
			{TaxCode: "VIEJO", RatePercent: 16, IsActive: false},
		},
	}

	applyTaxes(entry, concept)

	require.Len(t, entry.TaxDetail, 1)
	assert.Equal(t, "IVA", entry.TaxDetail[0].TaxCode)
	assert.InDelta(t, 19000.0, entry.TaxTotal, 0.001)
}

func TestApplyTaxes_SinImpuestos_DejaDetalleVacioNoNil(t *testing.T) {
	entry := &entities.Entry{Amount: 100000}

	applyTaxes(entry, &entities.Concept{})

	assert.NotNil(t, entry.TaxDetail)
	assert.Empty(t, entry.TaxDetail)
	assert.Zero(t, entry.TaxTotal)
}

func TestApplyTaxes_EsIdempotenteSobreElMismoEntry(t *testing.T) {
	entry := &entities.Entry{Amount: 100000}
	concept := &entities.Concept{
		Taxes: []entities.ConceptTax{{TaxCode: "IVA", RatePercent: 19, IsActive: true}},
	}

	applyTaxes(entry, concept)
	applyTaxes(entry, concept)

	require.Len(t, entry.TaxDetail, 1, "aplicar dos veces no debe duplicar las lineas")
	assert.InDelta(t, 19000.0, entry.TaxTotal, 0.001, "ni acumular el total")
}

func TestCreateManualEntry_MontoInvalido_Rechaza(t *testing.T) {
	for _, monto := range []float64{0, -1, -100000} {
		repo := &mocks.RepositoryMock{}

		_, err := newAccountingUseCase(repo, nil, nil, nil).CreateManualEntry(context.Background(),
			dtos.CreateEntryDTO{ConceptID: 1, Amount: monto})

		assert.ErrorIs(t, err, domainerrors.ErrInvalidAmount, "monto=%v", monto)
		assert.Empty(t, repo.CreatedEntries)
	}
}

func TestCreateManualEntry_ConceptoInactivo_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetConceptByIDFn: func(ctx context.Context, id uint) (*entities.Concept, error) {
			return &entities.Concept{ID: id, IsActive: false}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateManualEntry(context.Background(),
		dtos.CreateEntryDTO{ConceptID: 1, Amount: 1000})

	assert.ErrorIs(t, err, domainerrors.ErrConceptInactive)
	assert.Empty(t, repo.CreatedEntries)
}

func TestCreateManualEntry_HeredaElKindDelConcepto(t *testing.T) {
	for _, kind := range []string{entities.KindIncome, entities.KindExpense} {
		repo := &mocks.RepositoryMock{
			GetConceptByIDFn: func(ctx context.Context, id uint) (*entities.Concept, error) {
				return &entities.Concept{ID: id, Kind: kind, IsActive: true}, nil
			},
		}

		_, err := newAccountingUseCase(repo, nil, nil, nil).CreateManualEntry(context.Background(),
			dtos.CreateEntryDTO{ConceptID: 1, Amount: 1000})

		require.NoError(t, err)
		require.Len(t, repo.CreatedEntries, 1)
		assert.Equal(t, kind, repo.CreatedEntries[0].Kind,
			"el kind no lo elige el usuario, lo define el concepto")
	}
}

func TestCreateManualEntry_QuedaMarcadoComoManualYNoAutomatico(t *testing.T) {
	businessID := uint(26)
	fecha := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	repo := &mocks.RepositoryMock{}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateManualEntry(context.Background(),
		dtos.CreateEntryDTO{
			ConceptID: 1, Amount: 1000, BusinessID: &businessID,
			EntryDate: fecha, Description: "  ajuste  ",
		})

	require.NoError(t, err)
	require.Len(t, repo.CreatedEntries, 1)
	entry := repo.CreatedEntries[0]
	assert.Equal(t, entities.SourceManual, entry.SourceType)
	assert.Empty(t, entry.SourceID)
	assert.False(t, entry.IsAutomatic)
	assert.Equal(t, "ajuste", entry.Description)
	assert.Equal(t, fecha, entry.EntryDate)
	require.NotNil(t, entry.BusinessID)
	assert.Equal(t, uint(26), *entry.BusinessID)
}

func TestCreateManualEntry_AplicaLosImpuestosDelConcepto(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetConceptByIDFn: func(ctx context.Context, id uint) (*entities.Concept, error) {
			return &entities.Concept{
				ID: id, Kind: entities.KindIncome, IsActive: true,
				Taxes: []entities.ConceptTax{{TaxCode: "IVA", RatePercent: 19, IsActive: true}},
			}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateManualEntry(context.Background(),
		dtos.CreateEntryDTO{ConceptID: 1, Amount: 100000})

	require.NoError(t, err)
	require.Len(t, repo.CreatedEntries, 1)
	assert.InDelta(t, 19000.0, repo.CreatedEntries[0].TaxTotal, 0.001)
}

func TestCreateManualEntry_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("concepto inexistente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetConceptByIDFn: func(ctx context.Context, id uint) (*entities.Concept, error) {
				return nil, domainerrors.ErrConceptNotFound
			},
		}
		_, err := newAccountingUseCase(repo, nil, nil, nil).CreateManualEntry(context.Background(),
			dtos.CreateEntryDTO{ConceptID: 404, Amount: 1000})
		assert.ErrorIs(t, err, domainerrors.ErrConceptNotFound)
	})

	t.Run("al crear", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			CreateEntryFn: func(ctx context.Context, e *entities.Entry) (bool, error) { return false, dbErr },
		}
		_, err := newAccountingUseCase(repo, nil, nil, nil).CreateManualEntry(context.Background(),
			dtos.CreateEntryDTO{ConceptID: 1, Amount: 1000})
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestDeleteManualEntry_MovimientoAutomatico_NoSeBorra(t *testing.T) {
	casos := []struct {
		nombre string
		entry  *entities.Entry
	}{
		{"marcado como automatico", &entities.Entry{ID: 1, IsAutomatic: true, SourceType: entities.SourceManual}},
		{"origen guia", &entities.Entry{ID: 1, SourceType: entities.SourceGuideMargin}},
		{"origen suscripcion", &entities.Entry{ID: 1, SourceType: entities.SourceSubscription}},
		{"origen factura", &entities.Entry{ID: 1, SourceType: entities.SourceInvoice}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			borrado := false
			repo := &mocks.RepositoryMock{
				GetEntryByIDFn: func(ctx context.Context, id uint) (*entities.Entry, error) {
					return tc.entry, nil
				},
				DeleteEntryFn: func(ctx context.Context, id uint) error {
					borrado = true
					return nil
				},
			}

			err := newAccountingUseCase(repo, nil, nil, nil).DeleteManualEntry(context.Background(), 1)

			assert.ErrorIs(t, err, domainerrors.ErrEntryNotManual,
				"un movimiento generado por el sistema no se borra a mano: se recrearia en el proximo sync")
			assert.False(t, borrado)
		})
	}
}

func TestDeleteManualEntry_MovimientoManual_SeBorra(t *testing.T) {
	var borradoID uint
	repo := &mocks.RepositoryMock{
		GetEntryByIDFn: func(ctx context.Context, id uint) (*entities.Entry, error) {
			return &entities.Entry{ID: id, SourceType: entities.SourceManual, IsAutomatic: false}, nil
		},
		DeleteEntryFn: func(ctx context.Context, id uint) error {
			borradoID = id
			return nil
		},
	}

	err := newAccountingUseCase(repo, nil, nil, nil).DeleteManualEntry(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, uint(7), borradoID)
}

func TestDeleteManualEntry_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	repo := &mocks.RepositoryMock{
		GetEntryByIDFn: func(ctx context.Context, id uint) (*entities.Entry, error) { return nil, dbErr },
	}
	assert.ErrorIs(t, newAccountingUseCase(repo, nil, nil, nil).
		DeleteManualEntry(context.Background(), 1), dbErr)

	repo2 := &mocks.RepositoryMock{
		DeleteEntryFn: func(ctx context.Context, id uint) error { return dbErr },
	}
	assert.ErrorIs(t, newAccountingUseCase(repo2, nil, nil, nil).
		DeleteManualEntry(context.Background(), 1), dbErr)
}

func TestListEntries_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 10, 1, 10},
		{-5, 10, 1, 10},
		{1, 0, 1, 10},
		{1, 101, 1, 100},
		{1, 100, 1, 100},
		{2, 50, 2, 50},
	}

	for _, tc := range casos {
		var visto dtos.ListEntriesParams
		repo := &mocks.RepositoryMock{
			ListEntriesFn: func(ctx context.Context, params dtos.ListEntriesParams) ([]entities.Entry, int64, error) {
				visto = params
				return nil, 0, nil
			},
		}

		_, _, err := newAccountingUseCase(repo, nil, nil, nil).ListEntries(context.Background(),
			dtos.ListEntriesParams{Page: tc.page, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, visto.Page)
		assert.Equal(t, tc.wantSize, visto.PageSize)
	}
}

func TestListEntries_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListEntriesFn: func(ctx context.Context, params dtos.ListEntriesParams) ([]entities.Entry, int64, error) {
			return nil, 0, dbErr
		},
	}

	_, _, err := newAccountingUseCase(repo, nil, nil, nil).ListEntries(context.Background(), dtos.ListEntriesParams{})

	assert.ErrorIs(t, err, dbErr)
}

func TestReport_PeriodoInvalido_Rechaza(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	casos := []struct {
		nombre string
		params dtos.ReportParams
	}{
		{"sin desde", dtos.ReportParams{To: base}},
		{"sin hasta", dtos.ReportParams{From: base}},
		{"sin ninguna", dtos.ReportParams{}},
		{"hasta antes de desde", dtos.ReportParams{From: base, To: base.AddDate(0, -1, 0)}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			_, err := newAccountingUseCase(nil, nil, nil, nil).Report(context.Background(), tc.params)
			assert.ErrorIs(t, err, domainerrors.ErrInvalidPeriod)
		})
	}
}

func TestReport_MismoDiaEsValido(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	got, err := newAccountingUseCase(nil, nil, nil, nil).Report(context.Background(),
		dtos.ReportParams{From: base, To: base})

	require.NoError(t, err, "un reporte de un solo dia es valido")
	assert.Equal(t, "2026-06-01", got.From)
	assert.Equal(t, "2026-06-01", got.To)
}

func TestReport_SeparaCajaDeIngresoReal(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	repo := &mocks.RepositoryMock{
		ReportFn: func(ctx context.Context, params dtos.ReportParams) ([]dtos.ReportConceptRow, []dtos.ReportTaxRow, error) {
			return []dtos.ReportConceptRow{
				{Kind: entities.KindIncome, Amount: 100000, IsRealIncome: true, TaxTotal: 19000},
				{Kind: entities.KindIncome, Amount: 500000, IsRealIncome: false},
				{Kind: entities.KindExpense, Amount: 30000, IsRealIncome: true},
				{Kind: entities.KindExpense, Amount: 200000, IsRealIncome: false},
			}, nil, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).Report(context.Background(),
		dtos.ReportParams{From: base, To: base.AddDate(0, 1, 0)})

	require.NoError(t, err)
	assert.InDelta(t, 600000.0, got.Totals.CashIn, 0.001, "la caja entra toda")
	assert.InDelta(t, 100000.0, got.Totals.RealIncome, 0.001,
		"solo lo marcado como ingreso real cuenta como ingreso del negocio")
	assert.InDelta(t, 230000.0, got.Totals.CashOut, 0.001)
	assert.InDelta(t, 30000.0, got.Totals.RealExpense, 0.001)
	assert.InDelta(t, 19000.0, got.Totals.TaxTotal, 0.001)
}

func TestReport_NetoEsIngresoRealMenosGastoRealMenosImpuestos(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	repo := &mocks.RepositoryMock{
		ReportFn: func(ctx context.Context, params dtos.ReportParams) ([]dtos.ReportConceptRow, []dtos.ReportTaxRow, error) {
			return []dtos.ReportConceptRow{
				{Kind: entities.KindIncome, Amount: 100000, IsRealIncome: true, TaxTotal: 19000},
				{Kind: entities.KindExpense, Amount: 30000, IsRealIncome: true},
			}, nil, nil
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).Report(context.Background(),
		dtos.ReportParams{From: base, To: base.AddDate(0, 1, 0)})

	require.NoError(t, err)
	assert.InDelta(t, 51000.0, got.Totals.Net, 0.001, "100000 - 30000 - 19000")
}

func TestReport_SinMovimientos_TotalesEnCero(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	got, err := newAccountingUseCase(nil, nil, nil, nil).Report(context.Background(),
		dtos.ReportParams{From: base, To: base.AddDate(0, 1, 0)})

	require.NoError(t, err)
	assert.Zero(t, got.Totals.CashIn)
	assert.Zero(t, got.Totals.Net)
}

func TestReport_ErrorDelRepo_SePropaga(t *testing.T) {
	base := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ReportFn: func(ctx context.Context, params dtos.ReportParams) ([]dtos.ReportConceptRow, []dtos.ReportTaxRow, error) {
			return nil, nil, dbErr
		},
	}

	got, err := newAccountingUseCase(repo, nil, nil, nil).Report(context.Background(),
		dtos.ReportParams{From: base, To: base.AddDate(0, 1, 0)})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateService_ValidacionesYNormalizacion(t *testing.T) {
	t.Run("sin codigo", func(t *testing.T) {
		_, err := newAccountingUseCase(nil, nil, nil, nil).CreateService(context.Background(),
			dtos.SaveServiceDTO{Name: "X"})
		assert.ErrorIs(t, err, domainerrors.ErrCodeRequired)
	})

	t.Run("sin nombre", func(t *testing.T) {
		_, err := newAccountingUseCase(nil, nil, nil, nil).CreateService(context.Background(),
			dtos.SaveServiceDTO{Code: "A"})
		assert.ErrorIs(t, err, domainerrors.ErrNameRequired)
	})

	t.Run("precio negativo", func(t *testing.T) {
		_, err := newAccountingUseCase(nil, nil, nil, nil).CreateService(context.Background(),
			dtos.SaveServiceDTO{Code: "A", Name: "X", UnitPrice: -1})
		assert.ErrorIs(t, err, domainerrors.ErrInvalidAmount)
	})

	t.Run("precio cero es valido", func(t *testing.T) {
		_, err := newAccountingUseCase(nil, nil, nil, nil).CreateService(context.Background(),
			dtos.SaveServiceDTO{Code: "A", Name: "X", UnitPrice: 0})
		require.NoError(t, err)
	})
}

func TestCreateService_NaceActivoYNormalizaElCodigo(t *testing.T) {
	var creado dtos.SaveServiceDTO
	repo := &mocks.RepositoryMock{
		CreateServiceFn: func(ctx context.Context, dto dtos.SaveServiceDTO) (*entities.Service, error) {
			creado = dto
			return &entities.Service{ID: 1}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateService(context.Background(),
		dtos.SaveServiceDTO{Code: "  envio  ", Name: "  Envio  ", UnspscCode: "  78102200  ", IsActive: false})

	require.NoError(t, err)
	assert.Equal(t, "ENVIO", creado.Code)
	assert.Equal(t, "Envio", creado.Name)
	assert.Equal(t, "78102200", creado.UnspscCode)
	assert.True(t, creado.IsActive, "un servicio nuevo nace activo, el flag de entrada se ignora")
}

func TestCreateService_CodigoDuplicado_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ServiceCodeExistsFn: func(ctx context.Context, code string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).CreateService(context.Background(),
		dtos.SaveServiceDTO{Code: "A", Name: "X"})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateCode)
}

func TestUpdateService_ExcluyeElPropioIDAlBuscarDuplicados(t *testing.T) {
	var excluido *uint
	repo := &mocks.RepositoryMock{
		ServiceCodeExistsFn: func(ctx context.Context, code string, excludeID *uint) (bool, error) {
			excluido = excludeID
			return false, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateService(context.Background(),
		dtos.SaveServiceDTO{ID: 7, Code: "A", Name: "X"})

	require.NoError(t, err)
	require.NotNil(t, excluido)
	assert.Equal(t, uint(7), *excluido,
		"al actualizar, su propio codigo no debe contar como duplicado")
}

func TestUpdateService_ServicioInexistente_NoActualiza(t *testing.T) {
	actualizado := false
	repo := &mocks.RepositoryMock{
		GetServiceByIDFn: func(ctx context.Context, id uint) (*entities.Service, error) {
			return nil, domainerrors.ErrServiceNotFound
		},
		UpdateServiceFn: func(ctx context.Context, dto dtos.SaveServiceDTO) (*entities.Service, error) {
			actualizado = true
			return nil, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateService(context.Background(),
		dtos.SaveServiceDTO{ID: 404, Code: "A", Name: "X"})

	assert.ErrorIs(t, err, domainerrors.ErrServiceNotFound)
	assert.False(t, actualizado)
}

func TestUpdateService_CodigoDeOtroServicio_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ServiceCodeExistsFn: func(ctx context.Context, code string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateService(context.Background(),
		dtos.SaveServiceDTO{ID: 7, Code: "A", Name: "X"})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateCode)
}

func TestUpdateService_ValidaAntesDeLeer(t *testing.T) {
	leido := false
	repo := &mocks.RepositoryMock{
		GetServiceByIDFn: func(ctx context.Context, id uint) (*entities.Service, error) {
			leido = true
			return &entities.Service{ID: id}, nil
		},
	}

	_, err := newAccountingUseCase(repo, nil, nil, nil).UpdateService(context.Background(),
		dtos.SaveServiceDTO{ID: 7, Code: "", Name: "X"})

	assert.ErrorIs(t, err, domainerrors.ErrCodeRequired)
	assert.False(t, leido)
}

func TestDeleteService_InexistenteNoBorra(t *testing.T) {
	borrado := false
	repo := &mocks.RepositoryMock{
		GetServiceByIDFn: func(ctx context.Context, id uint) (*entities.Service, error) {
			return nil, domainerrors.ErrServiceNotFound
		},
		DeleteServiceFn: func(ctx context.Context, id uint) error {
			borrado = true
			return nil
		},
	}

	err := newAccountingUseCase(repo, nil, nil, nil).DeleteService(context.Background(), 404)

	assert.ErrorIs(t, err, domainerrors.ErrServiceNotFound)
	assert.False(t, borrado)
}

func TestDeleteService_Exitoso(t *testing.T) {
	var borradoID uint
	repo := &mocks.RepositoryMock{
		DeleteServiceFn: func(ctx context.Context, id uint) error {
			borradoID = id
			return nil
		},
	}

	require.NoError(t, newAccountingUseCase(repo, nil, nil, nil).DeleteService(context.Background(), 7))
	assert.Equal(t, uint(7), borradoID)
}

func TestListServices_DelegaYPropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListServicesFn: func(ctx context.Context) ([]entities.Service, error) {
			return []entities.Service{{ID: 1, Code: "ENVIO"}}, nil
		},
	}
	got, err := newAccountingUseCase(repo, nil, nil, nil).ListServices(context.Background())
	require.NoError(t, err)
	assert.Len(t, got, 1)

	dbErr := stderrors.New("timeout")
	repoErr := &mocks.RepositoryMock{
		ListServicesFn: func(ctx context.Context) ([]entities.Service, error) { return nil, dbErr },
	}
	_, err = newAccountingUseCase(repoErr, nil, nil, nil).ListServices(context.Background())
	assert.ErrorIs(t, err, dbErr)
}
