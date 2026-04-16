package usecases

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// GenerateMockGuidePDF generates a minimal valid PDF simulating a shipping guide.
// Uses no external dependencies — pure Go PDF construction.
func GenerateMockGuidePDF(trackingNumber, carrier, product, shipmentID, originCity, destCity string, flete float64) ([]byte, error) {
	now := time.Now().Format("02/01/2006 15:04:05")

	content := buildPDFContent(trackingNumber, carrier, product, shipmentID, originCity, destCity, flete, now)
	contentLen := len(content)

	var buf bytes.Buffer

	// PDF header
	buf.WriteString("%PDF-1.4\n")
	// Binary marker (4 bytes > 127 to signal binary content)
	buf.Write([]byte{'%', 0xe2, 0xe3, 0xcf, 0xd3, '\n'})

	// Track byte offsets for xref table (objects 1-5)
	offsets := make([]int64, 6)

	// Object 1: Catalog
	offsets[1] = int64(buf.Len())
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")

	// Object 2: Pages
	offsets[2] = int64(buf.Len())
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")

	// Object 3: Page (A4: 595 x 842 pts)
	offsets[3] = int64(buf.Len())
	buf.WriteString("3 0 obj\n")
	buf.WriteString("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842]\n")
	buf.WriteString("   /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>\n")
	buf.WriteString("endobj\n")

	// Object 4: Content stream
	offsets[4] = int64(buf.Len())
	buf.WriteString(fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n", contentLen))
	buf.WriteString(content)
	buf.WriteString("\nendstream\nendobj\n")

	// Object 5: Font (Helvetica — always available in PDF viewers)
	offsets[5] = int64(buf.Len())
	buf.WriteString("5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding >>\nendobj\n")

	// Cross-reference table
	xrefOffset := int64(buf.Len())
	buf.WriteString("xref\n")
	buf.WriteString("0 6\n")
	buf.WriteString("0000000000 65535 f \n")
	for i := 1; i <= 5; i++ {
		buf.WriteString(fmt.Sprintf("%010d 00000 n \n", offsets[i]))
	}

	// Trailer
	buf.WriteString("trailer\n<< /Size 6 /Root 1 0 R >>\n")
	buf.WriteString(fmt.Sprintf("startxref\n%d\n%%%%EOF\n", xrefOffset))

	return buf.Bytes(), nil
}

// buildPDFContent builds the page content stream for the mock guide PDF.
func buildPDFContent(trackingNumber, carrier, product, shipmentID, originCity, destCity string, flete float64, date string) string {
	var sb strings.Builder

	line := func(text string, yOffset int) {
		sb.WriteString(fmt.Sprintf("0 -%d Td\n(%s) Tj\n", yOffset, escPDF(text)))
	}

	sb.WriteString("BT\n")

	// Title
	sb.WriteString("/F1 16 Tf\n")
	sb.WriteString("50 800 Td\n")
	sb.WriteString("(EnvioClick Testing - Guia de Envio) Tj\n")

	// Separator line (approximated with underscores)
	sb.WriteString("/F1 10 Tf\n")
	line("------------------------------------------------------------", 14)
	line(fmt.Sprintf("DOCUMENTO DE PRUEBA - Generado: %s", date), 18)
	line("------------------------------------------------------------", 14)

	// Details
	sb.WriteString("/F1 12 Tf\n")
	line(fmt.Sprintf("ID de Envio:     %s", shipmentID), 28)
	line(fmt.Sprintf("Tracking:        %s", trackingNumber), 22)
	line(fmt.Sprintf("Transportadora:  %s", carrier), 22)
	line(fmt.Sprintf("Producto:        %s", product), 22)

	sb.WriteString("/F1 10 Tf\n")
	line("", 22)
	line("RUTA:", 18)
	sb.WriteString("/F1 12 Tf\n")
	line(fmt.Sprintf("  Origen:        %s", originCity), 20)
	line(fmt.Sprintf("  Destino:       %s", destCity), 20)

	sb.WriteString("/F1 10 Tf\n")
	line("", 22)
	line("COSTOS:", 18)
	sb.WriteString("/F1 12 Tf\n")
	line(fmt.Sprintf("  Flete:         $%.0f COP", flete), 20)

	// Footer
	sb.WriteString("/F1 8 Tf\n")
	line("", 80)
	line("------------------------------------------------------------", 14)
	line("AVISO: Este documento es generado por el simulador de pruebas de EnvioClick.", 16)
	line("No tiene validez legal ni operativa. Solo para uso en entornos de desarrollo.", 14)

	sb.WriteString("ET\n")
	return sb.String()
}

// escPDF escapes special PDF string literal characters.
func escPDF(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}
