package domain

import "testing"

func TestIsMatrixSearchBy(t *testing.T) {
	validos := []string{MatrixSearchAll, MatrixSearchSKU, MatrixSearchName, MatrixSearchBarcode}
	for _, v := range validos {
		if !IsMatrixSearchBy(v) {
			t.Fatalf("%q deberia ser valido", v)
		}
	}

	invalidos := []string{"", "SKU", "ean", "p.sku ILIKE ?", "todo"}
	for _, v := range invalidos {
		if IsMatrixSearchBy(v) {
			t.Fatalf("%q no deberia ser valido", v)
		}
	}
}

func TestMatrixQueryNormalizeSearchBy(t *testing.T) {
	q := MatrixQuery{}
	q.Normalize()
	if q.SearchBy != MatrixSearchAll {
		t.Fatalf("vacio deberia quedar en %q, quedo en %q", MatrixSearchAll, q.SearchBy)
	}

	q = MatrixQuery{SearchBy: MatrixSearchSKU}
	q.Normalize()
	if q.SearchBy != MatrixSearchSKU {
		t.Fatalf("un valor valido no se debe cambiar, quedo en %q", q.SearchBy)
	}

	q = MatrixQuery{SearchBy: "inventado"}
	q.Normalize()
	if q.SearchBy != "inventado" {
		t.Fatalf("Normalize no debe tapar un valor invalido: el handler lo rechaza")
	}
}
