package usecaseshipment

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractCoordinadoraMetadata(ctx context.Context, labelURL string) (map[string]interface{}, error) {
	if labelURL == "" {
		return nil, fmt.Errorf("label URL is empty")
	}

	pdfBytes, err := downloadPDF(ctx, labelURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}

	fullText, err := extractTextFromPDF(pdfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}

	metadata := extractCoordinadoraFields(fullText)
	return metadata, nil
}

func extractTextFromPDF(pdfBytes []byte) (string, error) {
	reader := bytes.NewReader(pdfBytes)
	pdfReader, err := pdf.NewReader(reader, int64(len(pdfBytes)))
	if err != nil {
		return "", fmt.Errorf("error reading PDF: %w", err)
	}

	var fullText strings.Builder
	for pageNum := 1; pageNum <= pdfReader.NumPage(); pageNum++ {
		page := pdfReader.Page(pageNum)

		text, err := page.GetPlainText(make(map[string]*pdf.Font))
		if err != nil {
			continue
		}
		fullText.WriteString(text + "\n")
	}

	return fullText.String(), nil
}

func extractCoordinadoraFields(pdfText string) map[string]interface{} {
	result := make(map[string]interface{})

	// El extractor de texto (github.com/ledongthuc/pdf) no respeta el orden
	// visual del PDF, sino el orden en que el generador de COORDINADORA
	// dibujo cada elemento. Por eso el valor de "Destino" aparece ANTES
	// de su etiqueta, mientras que "Zona Hub" y "Equipo Reparto" aparecen
	// DESPUES de la suya. Estos patrones fueron calibrados contra el texto
	// real extraido de una guia de produccion, no adivinados.
	patterns := map[string]string{
		"origen":         `Origen\s*\n+\s*(\d+\s*\n+\s*[A-Z]{2,4})\b`,
		"as_code":        `AS\s+(\d+)`,
		"paq":            `PAQ\s+([0-9\-]+)`,
		"unidad":         `UNIDAD:\s+([0-9/]+)`,
		"destino":        `([A-Za-z0-9]{1,8})\s*\n+\s*Destino\b`,
		"zona_hub":       `Zona Hub\s*[\n\r]+\s*(\d+)`,
		"equipo_reparto": `Equipo Reparto\s*[\n\r]+\s*(\d+)`,
		"ref":            `(?i)ref:\s*\n+\s*(?:ORDEN\s*\n+\s*)?([A-Z0-9\-]+)`,
		"guia":           `GUIA:\s+([0-9\.]+)`,
		"postal_origen":  `Postal:\s+([0-9\-]+)`,
		"observaciones":  `(?s)Observaciones Cliente:\s*\n+\s*(.+?)\n\n`,
	}

	for key, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(pdfText)
		if len(matches) > 1 {
			result[key] = normalizeExtractedField(matches[1])
		}
	}

	return result
}

// normalizeExtractedField colapsa saltos de linea/espacios internos (que
// quedan en capturas multi-linea como "origen" u "observaciones") a un
// solo espacio.
func normalizeExtractedField(s string) string {
	return strings.Join(strings.Fields(s), " ")
}


