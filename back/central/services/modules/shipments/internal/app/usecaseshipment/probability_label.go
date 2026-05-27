package usecaseshipment

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func buildProbabilityLabel(c *domain.GuidePDFContext, format *domain.GuideFormat) ([]byte, error) {
	wCm := format.WidthCm
	hCm := format.HeightCm
	if wCm < 4 || hCm < 4 {
		wCm = 10
		hCm = 15
	}

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: wCm * 10, Ht: hCm * 10},
		OrientationStr: "P",
	})
	pdf.SetMargins(3, 3, 3)
	pdf.SetAutoPageBreak(false, 3)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	scale := 1.0
	if wCm < 8 {
		scale = 0.75
	} else if wCm >= 20 {
		scale = 1.5
	}

	isLandscape := wCm > hCm
	isSquare := wCm == hCm

	if isLandscape {
		drawProbLabelLandscape(pdf, tr, c, wCm*10, hCm*10, scale)
	} else if isSquare {
		drawProbLabelSquare(pdf, tr, c, wCm*10, hCm*10, scale)
	} else {
		drawProbHeader(pdf, tr, c, scale)
		drawProbCOD(pdf, tr, c, scale)
		drawProbRecipient(pdf, tr, c, scale)
		drawProbBarcode(pdf, tr, c, wCm*10-6, scale)
		drawProbFooter(pdf, tr, c, scale)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func drawProbLabelLandscape(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, wMm, hMm, scale float64) {
	leftW := wMm * 0.55
	rightX := leftW + 2
	rightW := wMm - rightX - 3

	logoH := 6.0 * scale
	logoW := logoH * 4.7
	if logoW > leftW-3 {
		logoW = leftW - 3
		logoH = logoW / 4.7
	}
	if len(probabilityLogoPNG) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("prob-logo-l", opts, bytes.NewReader(probabilityLogoPNG))
		pdf.ImageOptions("prob-logo-l", 3, 3, logoW, logoH, false, opts, 0, "")
	} else {
		pdf.SetXY(3, 3)
		pdf.SetFont("Helvetica", "B", 10*scale)
		pdf.SetTextColor(20, 40, 90)
		pdf.CellFormat(leftW-3, 4.5*scale, tr("PROBABILITY"), "", 1, "L", false, 0, "")
	}
	drawCarrierBadge(pdf, tr, c.Carrier, leftW, 3, scale)
	pdf.SetXY(3, 3+logoH+0.5)

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 7*scale)
	business := strings.ToUpper(strings.TrimSpace(c.BusinessName))
	if business != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 3*scale, tr(business), "", 1, "L", false, 0, "")
	}
	if c.OrderNumber != "" {
		pdf.SetX(3)
		pdf.SetFont("Helvetica", "", 6*scale)
		pdf.SetTextColor(80, 80, 80)
		pdf.CellFormat(leftW-3, 2.5*scale, tr("Pedido: "+c.OrderNumber), "", 1, "L", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
	}

	pdf.SetX(3)
	pdf.SetDrawColor(20, 40, 90)
	pdf.SetLineWidth(0.3)
	pdf.Line(3, pdf.GetY()+0.5, leftW, pdf.GetY()+0.5)
	pdf.SetDrawColor(0, 0, 0)
	pdf.Ln(1.5 * scale)

	if c.CodTotal > 0 {
		y := pdf.GetY()
		pdf.SetFillColor(220, 40, 40)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Helvetica", "B", 7*scale)
		currency := c.Currency
		if currency == "" {
			currency = "COP"
		}
		codTxt := fmt.Sprintf("COD: $%s %s", formatMoneyProb(c.CodTotal), currency)
		pdf.Rect(3, y, leftW-3, 4*scale, "F")
		pdf.SetXY(3, y+0.3*scale)
		pdf.CellFormat(leftW-3, 3.4*scale, tr(codTxt), "", 1, "C", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFillColor(255, 255, 255)
		pdf.Ln(0.5)
	}

	pdf.SetX(3)
	pdf.SetFont("Helvetica", "B", 6*scale)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(leftW-3, 2.5*scale, tr("PARA:"), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 8*scale)
	pdf.SetX(3)
	pdf.CellFormat(leftW-3, 3.5*scale, tr(c.CustomerName), "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 6.5*scale)
	if c.CustomerDNI != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 2.8*scale, tr("CC: "+c.CustomerDNI), "", 1, "L", false, 0, "")
	}
	if c.CustomerPhone != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 2.8*scale, tr("Tel: "+c.CustomerPhone), "", 1, "L", false, 0, "")
	}
	pdf.SetFont("Helvetica", "B", 7*scale)
	pdf.SetX(3)
	pdf.MultiCell(leftW-3, 3*scale, tr(c.DestinationAddress), "", "L", false)
	cityState := joinNonEmptyProb(", ", c.DestinationCity, c.DestinationState)
	if cityState != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 3*scale, tr(cityState), "", 1, "L", false, 0, "")
	}

	pdf.SetXY(rightX, 3)
	pdf.SetFont("Helvetica", "", 6*scale)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(rightW, 2.5*scale, tr("GUIA - "+strings.ToUpper(c.Carrier)), "", 1, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	bcImg := buildCode128PNGProb(c.TrackingNumber, int(rightW*8), int(20*scale*8))
	if bcImg != nil {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("bcL.png", opts, bytes.NewReader(bcImg))
		pdf.ImageOptions("bcL.png", rightX, pdf.GetY(), rightW, 20*scale, false, opts, 0, "")
		pdf.SetY(pdf.GetY() + 20*scale + 0.5)
	}
	pdf.SetX(rightX)
	pdf.SetFont("Courier", "B", 11*scale)
	pdf.CellFormat(rightW, 5*scale, tr(c.TrackingNumber), "", 1, "C", false, 0, "")

	pdf.SetY(hMm - 7*scale)
	pdf.SetDrawColor(20, 40, 90)
	pdf.SetLineWidth(0.3)
	pdf.Line(3, pdf.GetY(), wMm-3, pdf.GetY())
	pdf.Ln(0.5)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetFont("Helvetica", "", 5*scale)
	footer := joinNonEmptyProb("  |  ",
		"Desde: "+joinNonEmptyProb(" - ", c.WarehouseCompany, c.WarehouseCity),
		fmt.Sprintf("%.1f kg", c.Weight),
		fmt.Sprintf("Vlr $%s", formatMoneyProb(c.DeclaredValue)),
		"Probability "+time.Now().Format("2006-01-02 15:04"),
	)
	pdf.SetX(3)
	pdf.CellFormat(wMm-6, 2.5*scale, tr(footer), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
}

func drawProbLabelSquare(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, wMm, hMm, scale float64) {
	logoH := 6.0 * scale
	logoW := logoH * 4.7
	if logoW > wMm*0.55 {
		logoW = wMm * 0.55
		logoH = logoW / 4.7
	}
	if len(probabilityLogoPNG) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("prob-logo-s", opts, bytes.NewReader(probabilityLogoPNG))
		pdf.ImageOptions("prob-logo-s", 3, 3, logoW, logoH, false, opts, 0, "")
	} else {
		pdf.SetXY(3, 3)
		pdf.SetFont("Helvetica", "B", 10*scale)
		pdf.SetTextColor(20, 40, 90)
		pdf.CellFormat(0, 4.5*scale, tr("PROBABILITY"), "", 1, "L", false, 0, "")
	}
	drawCarrierBadge(pdf, tr, c.Carrier, wMm-3, 3, scale)
	pdf.SetXY(3, 3+logoH+0.5)

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 7*scale)
	business := strings.ToUpper(strings.TrimSpace(c.BusinessName))
	if business != "" {
		pdf.CellFormat(0, 3*scale, tr(business), "", 1, "L", false, 0, "")
	}
	if c.OrderNumber != "" {
		pdf.SetFont("Helvetica", "", 6*scale)
		pdf.SetTextColor(80, 80, 80)
		pdf.CellFormat(0, 2.5*scale, tr("Pedido: "+c.OrderNumber), "", 1, "L", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
	}

	pdf.SetDrawColor(20, 40, 90)
	pdf.SetLineWidth(0.3)
	pdf.Line(3, pdf.GetY()+0.3, wMm-3, pdf.GetY()+0.3)
	pdf.SetDrawColor(0, 0, 0)
	pdf.Ln(1 * scale)

	if c.CodTotal > 0 {
		y := pdf.GetY()
		pdf.SetFillColor(220, 40, 40)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Helvetica", "B", 7*scale)
		currency := c.Currency
		if currency == "" {
			currency = "COP"
		}
		codTxt := fmt.Sprintf("COD $%s %s", formatMoneyProb(c.CodTotal), currency)
		pdf.Rect(3, y, wMm-6, 4*scale, "F")
		pdf.SetXY(3, y+0.3*scale)
		pdf.CellFormat(wMm-6, 3.4*scale, tr(codTxt), "", 1, "C", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFillColor(255, 255, 255)
		pdf.Ln(0.3)
	}

	pdf.SetFont("Helvetica", "B", 8*scale)
	pdf.CellFormat(0, 3*scale, tr(c.CustomerName), "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 6.5*scale)
	if c.CustomerPhone != "" {
		pdf.CellFormat(0, 2.5*scale, tr("Tel: "+c.CustomerPhone), "", 1, "L", false, 0, "")
	}
	pdf.SetFont("Helvetica", "B", 7*scale)
	pdf.MultiCell(0, 2.8*scale, tr(c.DestinationAddress), "", "L", false)
	cityState := joinNonEmptyProb(", ", c.DestinationCity, c.DestinationState)
	if cityState != "" {
		pdf.CellFormat(0, 2.8*scale, tr(cityState), "", 1, "L", false, 0, "")
	}
	pdf.Ln(0.5)

	pdf.SetFont("Helvetica", "", 5*scale)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(0, 2*scale, tr("GUIA - "+strings.ToUpper(c.Carrier)), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	bcImg := buildCode128PNGProb(c.TrackingNumber, int((wMm-6)*8), int(12*scale*8))
	if bcImg != nil {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("bcS.png", opts, bytes.NewReader(bcImg))
		pdf.ImageOptions("bcS.png", 3, pdf.GetY(), wMm-6, 12*scale, false, opts, 0, "")
		pdf.SetY(pdf.GetY() + 12*scale + 0.3)
	}
	pdf.SetFont("Courier", "B", 9*scale)
	pdf.CellFormat(0, 3.5*scale, tr(c.TrackingNumber), "", 1, "C", false, 0, "")
}

func drawProbHeader(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, scale float64) {
	pageW := pageWidth(pdf)
	yStart := pdf.GetY()

	logoH := 7.0 * scale
	logoW := logoH * 4.7
	if logoW > pageW*0.5 {
		logoW = pageW * 0.5
		logoH = logoW / 4.7
	}
	if len(probabilityLogoPNG) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("prob-logo", opts, bytes.NewReader(probabilityLogoPNG))
		pdf.ImageOptions("prob-logo", 3, yStart, logoW, logoH, false, opts, 0, "")
	} else {
		pdf.SetFont("Helvetica", "B", 11*scale)
		pdf.SetTextColor(20, 40, 90)
		pdf.SetXY(3, yStart)
		pdf.CellFormat(logoW, logoH, tr("PROBABILITY"), "", 0, "L", false, 0, "")
	}

	drawCarrierBadge(pdf, tr, c.Carrier, pageW-3, yStart, scale)

	pdf.SetY(yStart + logoH + 1)
	pdf.SetTextColor(0, 0, 0)
	business := strings.ToUpper(strings.TrimSpace(c.BusinessName))
	if business == "" {
		business = strings.ToUpper(strings.TrimSpace(c.WarehouseCompany))
	}
	if business != "" {
		pdf.SetFont("Helvetica", "B", 8*scale)
		pdf.CellFormat(0, 3.5*scale, tr(business), "", 1, "L", false, 0, "")
	}

	if c.OrderNumber != "" {
		pdf.SetFont("Helvetica", "", 7*scale)
		pdf.SetTextColor(80, 80, 80)
		pdf.CellFormat(0, 3*scale, tr("Pedido: "+c.OrderNumber), "", 1, "L", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
	}

	pdf.SetDrawColor(20, 40, 90)
	pdf.SetLineWidth(0.4)
	pdf.Line(3, pdf.GetY()+0.5, pageW-3, pdf.GetY()+0.5)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.2)
	pdf.Ln(1.5 * scale)
}

func drawCarrierBadge(pdf *gofpdf.Fpdf, tr func(string) string, carrier string, rightX, y, scale float64) {
	carrier = strings.ToUpper(strings.TrimSpace(carrier))
	if carrier == "" {
		return
	}
	st := styleForCarrier(carrier)
	pdf.SetFont("Helvetica", "B", 8*scale)
	textW := pdf.GetStringWidth(carrier) + 4*scale
	h := 5.5 * scale
	x := rightX - textW
	pdf.SetFillColor(st.BgR, st.BgG, st.BgB)
	pdf.Rect(x, y, textW, h, "F")
	pdf.SetTextColor(st.TxtR, st.TxtG, st.TxtB)
	pdf.SetXY(x, y+0.5*scale)
	pdf.CellFormat(textW, h-1*scale, tr(carrier), "", 0, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
}

func drawProbCOD(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, scale float64) {
	if c.CodTotal <= 0 {
		return
	}
	pageW := pageWidth(pdf)
	y := pdf.GetY()
	pdf.SetFillColor(220, 40, 40)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 9*scale)
	currency := c.Currency
	if currency == "" {
		currency = "COP"
	}
	codTxt := fmt.Sprintf("CONTRA ENTREGA: $%s %s", formatMoneyProb(c.CodTotal), currency)
	pdf.Rect(3, y, pageW-6, 5.5*scale, "F")
	pdf.SetXY(3, y+0.5*scale)
	pdf.CellFormat(pageW-6, 4.5*scale, tr(codTxt), "", 1, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(255, 255, 255)
	pdf.Ln(1 * scale)
}

func drawProbRecipient(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, scale float64) {
	pdf.SetFont("Helvetica", "B", 7*scale)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(0, 3*scale, tr("PARA:"), "", 1, "L", false, 0, "")

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 10*scale)
	name := strings.TrimSpace(c.CustomerName)
	if name == "" {
		name = "(sin nombre)"
	}
	pdf.CellFormat(0, 5*scale, tr(name), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 8*scale)
	if c.CustomerDNI != "" {
		pdf.CellFormat(0, 3.5*scale, tr("CC: "+c.CustomerDNI), "", 1, "L", false, 0, "")
	}
	if c.CustomerPhone != "" {
		pdf.CellFormat(0, 3.5*scale, tr("Tel: "+c.CustomerPhone), "", 1, "L", false, 0, "")
	}

	pdf.SetFont("Helvetica", "B", 9*scale)
	if c.DestinationAddress != "" {
		pdf.MultiCell(0, 4*scale, tr(c.DestinationAddress), "", "L", false)
	}
	cityState := joinNonEmptyProb(", ", c.DestinationCity, c.DestinationState)
	if cityState != "" {
		pdf.CellFormat(0, 4*scale, tr(cityState), "", 1, "L", false, 0, "")
	}
	if c.DestinationSuburb != "" {
		pdf.SetFont("Helvetica", "", 7*scale)
		pdf.CellFormat(0, 3*scale, tr("Barrio: "+c.DestinationSuburb), "", 1, "L", false, 0, "")
	}
	pdf.Ln(1.5 * scale)
}

func drawProbBarcode(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, widthMm float64, scale float64) {
	if c.TrackingNumber == "" {
		return
	}

	pdf.SetFont("Helvetica", "", 7*scale)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(0, 3*scale, tr("GUIA - "+strings.ToUpper(c.Carrier)), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	bcImg := buildCode128PNGProb(c.TrackingNumber, int(widthMm*8), int(10*scale*8))
	if bcImg != nil {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("bc.png", opts, bytes.NewReader(bcImg))
		pdf.ImageOptions("bc.png", 3, pdf.GetY(), widthMm, 10*scale, false, opts, 0, "")
		pdf.SetY(pdf.GetY() + 10*scale + 0.5)
	}

	pdf.SetFont("Courier", "B", 12*scale)
	pdf.CellFormat(0, 5*scale, tr(c.TrackingNumber), "", 1, "C", false, 0, "")
	pdf.Ln(1 * scale)
}

func drawProbFooter(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, scale float64) {
	_, pageH := pdf.GetPageSize()
	pdf.SetY(pageH - 10*scale)
	pdf.SetDrawColor(20, 40, 90)
	pdf.SetLineWidth(0.3)
	pdf.Line(3, pdf.GetY(), pageWidth(pdf)-3, pdf.GetY())
	pdf.Ln(0.5)

	pdf.SetFont("Helvetica", "", 6*scale)
	pdf.SetTextColor(80, 80, 80)
	origen := joinNonEmptyProb(" - ", c.WarehouseCompany, c.WarehouseCity, c.WarehouseState)
	if origen != "" {
		pdf.CellFormat(0, 2.5*scale, tr("Desde: "+origen), "", 1, "L", false, 0, "")
	}

	details := []string{}
	if c.Weight > 0 {
		details = append(details, fmt.Sprintf("%.1f kg", c.Weight))
	}
	if c.Length > 0 || c.Width > 0 || c.Height > 0 {
		details = append(details, fmt.Sprintf("%.0fx%.0fx%.0f cm", c.Length, c.Width, c.Height))
	}
	if c.DeclaredValue > 0 {
		details = append(details, fmt.Sprintf("Vlr $%s", formatMoneyProb(c.DeclaredValue)))
	}
	if len(details) > 0 {
		pdf.CellFormat(0, 2.5*scale, tr(strings.Join(details, "  |  ")), "", 1, "L", false, 0, "")
	}
	pdf.CellFormat(0, 2.5*scale, tr("Generada por Probability "+time.Now().Format("2006-01-02 15:04")), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
}

func buildCode128PNGProb(data string, widthPx, heightPx int) []byte {
	if data == "" {
		return nil
	}
	bc, err := code128.Encode(data)
	if err != nil {
		return nil
	}
	if widthPx < 200 {
		widthPx = 200
	}
	if heightPx < 40 {
		heightPx = 40
	}
	scaled, err := barcode.Scale(bc, widthPx, heightPx)
	if err != nil {
		return nil
	}
	b := scaled.Bounds()
	gray := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, _, _, _ := scaled.At(x, y).RGBA()
			if r>>8 > 200 {
				gray.SetGray(x, y, color.Gray{Y: 255})
			} else {
				gray.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, gray); err != nil {
		return nil
	}
	return buf.Bytes()
}

func buildQRPNGProb(data string) []byte {
	if data == "" {
		return nil
	}
	pngBytes, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return nil
	}
	return pngBytes
}

func joinNonEmptyProb(sep string, parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, sep)
}

func formatMoneyProb(v float64) string {
	n := int64(v)
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	pre := len(s) % 3
	if pre > 0 {
		b.WriteString(s[:pre])
		if len(s) > pre {
			b.WriteString(".")
		}
	}
	for i := pre; i < len(s); i += 3 {
		b.WriteString(s[i : i+3])
		if i+3 < len(s) {
			b.WriteString(".")
		}
	}
	return b.String()
}

var _ = buildQRPNGProb

func pageWidth(pdf *gofpdf.Fpdf) float64 {
	w, _ := pdf.GetPageSize()
	return w
}
