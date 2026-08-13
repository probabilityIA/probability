package domain

import (
	"strings"
	"testing"
	"time"
)

func producto(name, description string) *Product {
	return &Product{
		ID:          "PRD_1",
		BusinessID:  46,
		Name:        name,
		Description: description,
		Price:       15000,
	}
}

func porCampo(changes []FieldChange) map[string]FieldChange {
	out := make(map[string]FieldChange, len(changes))
	for _, c := range changes {
		out[c.Field] = c
	}
	return out
}

func TestDiffProductSinCambiosNoRegistraNada(t *testing.T) {
	antes := producto("Torero", "algo")
	despues := producto("Torero", "algo")

	if changes := DiffProduct(antes, despues, ChannelOrigin(254), "", time.Now()); len(changes) != 0 {
		t.Fatalf("una sincronizacion que reescribe lo mismo no debe dejar historial, dejo %d", len(changes))
	}
}

func TestDiffProductGuardaElValorAnterior(t *testing.T) {
	antes := producto("Torero", "")
	despues := producto("Torero deportivo", "")

	changes := DiffProduct(antes, despues, ChannelOrigin(254), "lote-1", time.Now())
	if len(changes) != 1 {
		t.Fatalf("esperaba 1 cambio, hubo %d", len(changes))
	}

	change := changes[0]
	if change.Field != FieldName {
		t.Errorf("campo = %q, esperaba %q", change.Field, FieldName)
	}
	if change.OldValue != "Torero" {
		t.Errorf("el historial debe guardar el valor anterior, guardo %q", change.OldValue)
	}
	if change.BatchID != "lote-1" {
		t.Errorf("batch = %q, esperaba lote-1", change.BatchID)
	}
	if change.Origin.IntegrationID == nil || *change.Origin.IntegrationID != 254 {
		t.Error("el origen debe apuntar a la integracion que escribio")
	}
}

func TestDiffProductLlenarVacioNoOcupaEspacio(t *testing.T) {
	antes := producto("Torero", "")
	despues := producto("Torero", "Buzo con capota")

	changes := DiffProduct(antes, despues, ChannelOrigin(221), "", time.Now())
	if len(changes) != 1 {
		t.Fatalf("esperaba 1 cambio, hubo %d", len(changes))
	}
	if changes[0].OldValue != "" {
		t.Errorf("llenar un campo vacio no guarda valor anterior, guardo %q", changes[0].OldValue)
	}
}

func TestDiffProductRecortaLaDescripcionLarga(t *testing.T) {
	larga := strings.Repeat("x", MaxOldValueLen+500)
	antes := producto("Torero", larga)
	despues := producto("Torero", "corta")

	changes := porCampo(DiffProduct(antes, despues, ManualOrigin(7), "", time.Now()))
	change, ok := changes[FieldDescription]
	if !ok {
		t.Fatal("esperaba un cambio de descripcion")
	}
	if len(change.OldValue) != MaxOldValueLen {
		t.Errorf("largo = %d, esperaba el tope %d", len(change.OldValue), MaxOldValueLen)
	}
	if !change.Truncated {
		t.Error("un valor recortado debe quedar marcado como tal")
	}
}

func TestDiffProductIgnoraSoloEspacios(t *testing.T) {
	antes := producto("Torero", "algo")
	despues := producto("  Torero  ", "algo")

	if changes := DiffProduct(antes, despues, ChannelOrigin(254), "", time.Now()); len(changes) != 0 {
		t.Errorf("un cambio de espacios no es un cambio de dato, dejo %d filas", len(changes))
	}
}

func TestDiffProductManualQuedaConUsuario(t *testing.T) {
	antes := producto("Torero", "")
	despues := producto("Torero editado", "")

	changes := DiffProduct(antes, despues, ManualOrigin(9), "", time.Now())
	if len(changes) != 1 {
		t.Fatalf("esperaba 1 cambio, hubo %d", len(changes))
	}
	if changes[0].Origin.Source != SourceManual {
		t.Errorf("source = %q, esperaba %q", changes[0].Origin.Source, SourceManual)
	}
	if changes[0].Origin.UserID == nil || *changes[0].Origin.UserID != 9 {
		t.Error("una edicion manual debe quedar sellada con el usuario")
	}
	if changes[0].Origin.IntegrationID != nil {
		t.Error("una edicion manual no tiene integracion")
	}
}
