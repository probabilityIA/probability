package usecases

import (
	"fmt"
	"testing"
)

func TestFailedSKUsAcumula(t *testing.T) {
	f := &failedSKUs{}
	f.add("A-1")
	f.add("A-2")

	if f.count() != 2 {
		t.Fatalf("count = %d, se esperaban 2", f.count())
	}
	if len(f.list()) != 2 || f.list()[0] != "A-1" {
		t.Fatalf("list = %v", f.list())
	}
	if f.truncated() != 0 {
		t.Fatalf("truncated = %d, no se oculto nada", f.truncated())
	}
}

func TestFailedSKUsSinFallosDevuelveListaVaciaNoNil(t *testing.T) {
	f := &failedSKUs{}
	if f.list() == nil {
		t.Fatal("list() nunca debe ser nil: serializaria como null en el evento")
	}
	if len(f.list()) != 0 || f.count() != 0 {
		t.Fatal("sin fallos la lista debe ir vacia")
	}
}

func TestFailedSKUsCuentaAunSinSKU(t *testing.T) {
	f := &failedSKUs{}
	f.add("")
	f.add("A-1")

	if f.count() != 2 {
		t.Fatalf("count = %d: un fallo sin SKU igual debe contarse", f.count())
	}
	if len(f.list()) != 1 {
		t.Fatalf("list = %v: no debe incluir SKUs vacios", f.list())
	}
}

func TestFailedSKUsToleraTopeSinMentirElTotal(t *testing.T) {
	f := &failedSKUs{}
	for i := 0; i < maxReportedFailedSKUs+30; i++ {
		f.add(fmt.Sprintf("SKU-%d", i))
	}

	if f.count() != maxReportedFailedSKUs+30 {
		t.Fatalf("count = %d: el total debe reflejar TODOS los fallos", f.count())
	}
	if len(f.list()) != maxReportedFailedSKUs {
		t.Fatalf("list = %d, debe toparse en %d", len(f.list()), maxReportedFailedSKUs)
	}
	if f.truncated() != 30 {
		t.Fatalf("truncated = %d, se esperaban 30 ocultos: el tope no puede ser silencioso", f.truncated())
	}
}
