package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
)

func TestUpdateCarrierFee_ActualizaYDevuelveLaComisionAnterior(t *testing.T) {
	repo := &repoMock{feePrevia: 6236}
	uc := New(repo, log.New())

	res, err := uc.UpdateCarrierFee(context.Background(), dtos.UpdateCarrierFeeDTO{
		BusinessID: 46,
		ShipmentID: 43001,
		Fee:        6835,
	})

	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if repo.feeLlamadas != 1 {
		t.Errorf("se esperaba 1 actualizacion, se obtuvieron: %d", repo.feeLlamadas)
	}
	if repo.feeBusinessID != 46 {
		t.Errorf("se esperaba business 46, se obtuvo: %d", repo.feeBusinessID)
	}
	if repo.feeShipmentID != 43001 {
		t.Errorf("se esperaba envio 43001, se obtuvo: %d", repo.feeShipmentID)
	}
	if repo.feeAplicada != 6835 {
		t.Errorf("se esperaba comision 6835, se obtuvo: %v", repo.feeAplicada)
	}
	if res.PreviousFee != 6236 || res.NewFee != 6835 {
		t.Errorf("resultado inesperado: anterior=%v nueva=%v", res.PreviousFee, res.NewFee)
	}
	if res.OrderNumber != "VIG-0083" {
		t.Errorf("se esperaba VIG-0083, se obtuvo: %s", res.OrderNumber)
	}
}

func TestUpdateCarrierFee_RechazaEntradasInvalidas(t *testing.T) {
	casos := []struct {
		nombre string
		dto    dtos.UpdateCarrierFeeDTO
	}{
		{"sin negocio", dtos.UpdateCarrierFeeDTO{ShipmentID: 1, Fee: 100}},
		{"sin envio", dtos.UpdateCarrierFeeDTO{BusinessID: 46, Fee: 100}},
		{"comision negativa", dtos.UpdateCarrierFeeDTO{BusinessID: 46, ShipmentID: 1, Fee: -1}},
		{"comision desbordada", dtos.UpdateCarrierFeeDTO{BusinessID: 46, ShipmentID: 1, Fee: 100000001}},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			repo := &repoMock{}
			uc := New(repo, log.New())

			if _, err := uc.UpdateCarrierFee(context.Background(), caso.dto); err == nil {
				t.Error("se esperaba error")
			}
			if repo.feeLlamadas != 0 {
				t.Errorf("no se debio tocar el repositorio, llamadas: %d", repo.feeLlamadas)
			}
		})
	}
}

func TestUpdateCarrierFee_ComisionCeroEsValida(t *testing.T) {
	repo := &repoMock{feePrevia: 6236}
	uc := New(repo, log.New())

	if _, err := uc.UpdateCarrierFee(context.Background(), dtos.UpdateCarrierFeeDTO{
		BusinessID: 46,
		ShipmentID: 43001,
		Fee:        0,
	}); err != nil {
		t.Fatalf("cero debe ser valido, se obtuvo: %v", err)
	}
	if repo.feeAplicada != 0 || repo.feeLlamadas != 1 {
		t.Errorf("se esperaba aplicar 0 una vez, se obtuvo %v en %d llamadas", repo.feeAplicada, repo.feeLlamadas)
	}
}

func TestUpdateCarrierFee_PropagaElErrorDelRepositorio(t *testing.T) {
	repo := &repoMock{feeErr: errors.New("el envio no existe o no pertenece al negocio")}
	uc := New(repo, log.New())

	if _, err := uc.UpdateCarrierFee(context.Background(), dtos.UpdateCarrierFeeDTO{
		BusinessID: 46,
		ShipmentID: 99999,
		Fee:        6835,
	}); err == nil {
		t.Error("se esperaba el error del repositorio")
	}
}
