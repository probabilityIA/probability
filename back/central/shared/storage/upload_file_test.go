package storage

import "testing"

func TestContentTypePorExtension(t *testing.T) {
	casos := map[string]string{
		"mapas-direccion/abc.png": "image/png",
		"foto.jpg":                "image/jpeg",
		"doc.pdf":                 "application/pdf",
		"sin-extension":           "",
	}

	for archivo, esperado := range casos {
		got := contentTypePorExtension(archivo)
		if esperado == "" {
			if got != "" {
				t.Errorf("%q devolvio %q, se esperaba vacio", archivo, got)
			}
			continue
		}
		if got == "" {
			t.Errorf("%q no resolvio content type; S3 lo serviria como octet-stream y Meta rechaza la imagen", archivo)
			continue
		}
		if got[:len(esperado)] != esperado {
			t.Errorf("%q devolvio %q, se esperaba %q", archivo, got, esperado)
		}
	}
}
