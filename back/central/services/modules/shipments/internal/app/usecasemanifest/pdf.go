package usecasemanifest

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

const (
	pageW = 210.0
	pageH = 297.0
	mLeft = 8.0
	mTop  = 8.0
)

func buildManifestPDF(in domain.ManifestPDFInput) ([]byte, error) {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: pageW, Ht: pageH},
		OrientationStr: "P",
	})
	pdf.SetMargins(mLeft, mTop, mLeft)
	pdf.SetAutoPageBreak(true, 12)
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	rowsPerPage := 28
	total := len(in.Rows)
	pages := (total + rowsPerPage - 1) / rowsPerPage
	if pages == 0 {
		pages = 1
	}

	totalBultos := 0
	for range in.Rows {
		totalBultos++
	}

	for p := 0; p < pages; p++ {
		pdf.AddPage()
		drawHeader(pdf, tr, in, p+1, pages)
		drawInfoBox(pdf, tr, in)
		start := p * rowsPerPage
		end := start + rowsPerPage
		if end > total {
			end = total
		}
		drawTable(pdf, tr, in.Rows[start:end], start+1)
		if p == pages-1 {
			drawSignatures(pdf, tr, total, totalBultos)
		}
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func drawHeader(pdf *gofpdf.Fpdf, tr func(string) string, in domain.ManifestPDFInput, page, totalPages int) {
	if len(probabilityLogoPNG) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("manifest-logo", opts, bytes.NewReader(probabilityLogoPNG))
		pdf.ImageOptions("manifest-logo", mLeft, mTop, 38, 12, false, opts, 0, "")
	}

	pdf.SetXY(mLeft+45, mTop+2)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(110, 8, tr("RELACION DE DESPACHO"), "", 0, "C", false, 0, "")

	pdf.SetXY(pageW-mLeft-30, mTop)
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(30, 4, tr(fmt.Sprintf("Pagina  %d/%d", page, totalPages)), "", 0, "R", false, 0, "")
}

func drawInfoBox(pdf *gofpdf.Fpdf, tr func(string) string, in domain.ManifestPDFInput) {
	y := mTop + 16
	w := pageW - 2*mLeft
	h := 22.0
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)
	pdf.Rect(mLeft, y, w, h, "D")

	pdf.SetFont("Helvetica", "B", 8)
	leftCol := mLeft + 2
	midCol := mLeft + 70
	rightCol := mLeft + 140

	pdf.SetXY(leftCol, y+2)
	pdf.CellFormat(20, 4, tr("NIT"), "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(45, 4, tr(in.BusinessCode), "", 0, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetXY(midCol, y+2)
	pdf.CellFormat(35, 4, tr("POBLACION IMPOSICION:"), "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(30, 4, tr(strOr(in.OriginCity, "BOGOTA")), "", 0, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetXY(rightCol, y+2)
	pdf.CellFormat(28, 4, tr("N° MANIFIESTO"), "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(30, 4, tr(in.ManifestNo), "", 0, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetXY(leftCol, y+8)
	pdf.CellFormat(20, 4, tr("CUENTA"), "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(45, 4, tr(in.BusinessCode), "", 0, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetXY(midCol, y+8)
	pdf.CellFormat(35, 4, tr("FECHA DE CREACION:"), "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(30, 4, in.GeneratedAt.Format("02/01/2006"), "", 0, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetXY(rightCol, y+8)
	pdf.CellFormat(28, 4, tr("USUARIO:"), "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(30, 4, tr(in.GeneratedBy), "", 0, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetXY(leftCol, y+14)
	pdf.CellFormat(28, 4, tr("RAZON SOCIAL:"), "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(140, 4, tr(in.BusinessName), "", 0, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetXY(leftCol, y+18)
	pdf.CellFormat(28, 3, tr("TRANSPORTADORA:"), "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(60, 3, tr(in.Carrier), "", 0, "L", false, 0, "")

	if logo := getCarrierLogoPNG(in.Carrier); len(logo) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		key := "carrier-" + carrierKey(in.Carrier)
		pdf.RegisterImageOptionsReader(key, opts, bytes.NewReader(logo))
		pdf.ImageOptions(key, pageW-mLeft-26, y+2, 22, 18, false, opts, 0, "")
	}
}

func drawTable(pdf *gofpdf.Fpdf, tr func(string) string, rows []domain.ManifestShipmentRow, startNum int) {
	y := mTop + 42
	pdf.SetY(y)
	pdf.SetX(mLeft)

	cols := []struct {
		W     float64
		Label string
		Align string
	}{
		{10, "N°", "C"},
		{14, "Prod", "C"},
		{30, "Guia", "L"},
		{16, "N° Paquete", "C"},
		{24, "Documento", "L"},
		{58, "Destinatario", "L"},
		{42, "Ciudad", "L"},
	}

	pdf.SetFillColor(220, 220, 220)
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetDrawColor(0, 0, 0)
	for _, c := range cols {
		pdf.CellFormat(c.W, 6, tr(c.Label), "1", 0, c.Align, true, 0, "")
	}
	pdf.Ln(6)

	pdf.SetFont("Helvetica", "", 7.5)
	for i, r := range rows {
		num := fmt.Sprintf("%d", startNum+i)
		pdf.CellFormat(cols[0].W, 5, num, "1", 0, "C", false, 0, "")
		pdf.CellFormat(cols[1].W, 5, tr(truncate(r.CarrierCode, 6)), "1", 0, "C", false, 0, "")
		pdf.CellFormat(cols[2].W, 5, tr(truncate(r.TrackingNumber, 16)), "1", 0, "L", false, 0, "")
		pdf.CellFormat(cols[3].W, 5, fmt.Sprintf("%d", r.ShipmentID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(cols[4].W, 5, tr(truncate(r.CustomerDocument, 14)), "1", 0, "L", false, 0, "")
		pdf.CellFormat(cols[5].W, 5, tr(truncate(r.CustomerName, 40)), "1", 0, "L", false, 0, "")
		pdf.CellFormat(cols[6].W, 5, tr(truncate(r.DestinationCity, 28)), "1", 0, "L", false, 0, "")
		pdf.Ln(5)
	}
}

func drawSignatures(pdf *gofpdf.Fpdf, tr func(string) string, totalEnvios, totalBultos int) {
	y := pageH - 50
	w := (pageW - 2*mLeft) / 3
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)

	labels := []string{"CLIENTE", "TRANSPORTE", "SUCURSAL"}
	for i, lab := range labels {
		x := mLeft + float64(i)*w
		pdf.Rect(x, y, w, 28, "D")
		pdf.SetXY(x, y+1)
		pdf.SetFont("Helvetica", "B", 9)
		pdf.CellFormat(w, 5, tr(lab), "", 0, "C", false, 0, "")
		pdf.SetXY(x, y+22)
		pdf.SetFont("Helvetica", "", 7)
		pdf.CellFormat(w, 4, tr("NOMBRE, NIT, FIRMA y FECHA"), "T", 0, "C", false, 0, "")
	}

	pdf.SetXY(mLeft, y+30)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.CellFormat(50, 6, tr(fmt.Sprintf("TOTAL ENVIOS:  %d", totalEnvios)), "", 0, "L", false, 0, "")
	pdf.SetX(pageW - mLeft - 50)
	pdf.CellFormat(50, 6, tr(fmt.Sprintf("TOTAL BULTOS:  %d", totalBultos)), "", 0, "R", false, 0, "")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "."
}

func strOr(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
