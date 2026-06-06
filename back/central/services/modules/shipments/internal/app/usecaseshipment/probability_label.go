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
		drawProbSender(pdf, tr, c, scale)
		drawProbRecipient(pdf, tr, c, scale)
		drawProbDetailsBox(pdf, tr, c, scale)
		drawProbBarcode(pdf, tr, c, wCm*10-6, scale)
		drawProbProofOfDelivery(pdf, tr, c, scale)
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
	pdf.SetTextColor(20, 40, 90)
	pdf.CellFormat(leftW-3, 2.5*scale, tr("REMITE (DESDE):"), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 7*scale)
	warehouse := strings.ToUpper(strings.TrimSpace(c.WarehouseCompany))
	if warehouse != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 2.8*scale, tr(warehouse), "", 1, "L", false, 0, "")
	}
	pdf.SetFont("Helvetica", "", 6*scale)
	if c.WarehouseAddress != "" {
		pdf.SetX(3)
		pdf.MultiCell(leftW-3, 2.5*scale, tr(c.WarehouseAddress), "", "L", false)
	}
	if c.WarehouseCity != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 2.5*scale, tr(c.WarehouseCity), "", 1, "L", false, 0, "")
	}
	if c.WarehousePhone != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 2.5*scale, tr("Tel: "+c.WarehousePhone), "", 1, "L", false, 0, "")
	}

	pdf.Ln(1*scale)
	pdf.SetX(3)
	pdf.SetFont("Helvetica", "B", 6*scale)
	pdf.SetTextColor(20, 40, 90)
	pdf.CellFormat(leftW-3, 2.5*scale, tr("PARA: (CLIENTE)"), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 7.5*scale)
	pdf.SetX(3)
	pdf.CellFormat(leftW-3, 3*scale, tr(c.CustomerName), "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 6*scale)
	if c.CustomerDNI != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 2.5*scale, tr("CC: "+c.CustomerDNI), "", 1, "L", false, 0, "")
	}
	if c.CustomerPhone != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 2.5*scale, tr("Tel: "+c.CustomerPhone), "", 1, "L", false, 0, "")
	}
	pdf.SetFont("Helvetica", "B", 6.5*scale)
	pdf.SetX(3)
	pdf.MultiCell(leftW-3, 2.8*scale, tr(c.DestinationAddress), "", "L", false)
	cityState := joinNonEmptyProb(", ", c.DestinationCity, c.DestinationState)
	if cityState != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 2.5*scale, tr(cityState), "", 1, "L", false, 0, "")
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

	pdf.SetY(hMm - 22*scale)
	pdf.SetDrawColor(20, 40, 90)
	pdf.SetLineWidth(0.3)
	pdf.Line(3, pdf.GetY(), wMm-3, pdf.GetY())
	pdf.Ln(0.5)

	pdf.SetFont("Helvetica", "B", 5*scale)
	pdf.SetTextColor(20, 40, 90)
	pdf.CellFormat(wMm-6, 2.2*scale, tr("COLVANES S.A.S. - ENVIA"), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 4.5*scale)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(wMm-6, 1.8*scale, tr("NIT: 800.185.306-4  |  Carrera 88 # 17B-10 Bogota"), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 4*scale)
	pdf.SetTextColor(60, 60, 60)
	pdf.CellFormat(wMm-6, 1.6*scale, tr("Tel: (1) 7943670 | www.envia.co | administradorpqr1@enviacolvanes.com.co"), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 3.5*scale)
	pdf.SetTextColor(80, 80, 80)
	lic1 := "Lic Min Transporte 0080 (14/3/2000) | Lic MinTIC 001368 (4/8/2020) | CIIU 5320 Mensajeria Express"
	pdf.MultiCell(wMm-6, 1.3*scale, tr(lic1), "", "L", false)

	pdf.SetFont("Helvetica", "", 3.5*scale)
	pdf.SetTextColor(80, 80, 80)
	lic2 := "Autorretenedores Res 4327 (Jul/97) | Grandes Contribuyentes Res 9061 (Dic/20) | Agente Retenedor de IVA"
	pdf.MultiCell(wMm-6, 1.3*scale, tr(lic2), "", "L", false)

	pdf.Ln(0.3*scale)
	pdf.SetFont("Helvetica", "", 3.5*scale)
	pdf.SetTextColor(80, 80, 80)
	weight := "N/D"
	if c.Weight > 0 {
		weight = fmt.Sprintf("%.1f kg", c.Weight)
	}
	dims := "N/D"
	if c.Length > 0 || c.Width > 0 || c.Height > 0 {
		dims = fmt.Sprintf("%.0fx%.0fx%.0f cm", c.Length, c.Width, c.Height)
	}
	details := joinNonEmptyProb("  |  ",
		"Peso: "+weight,
		"Dim: "+dims,
		"Vlr: $"+formatMoneyProb(c.DeclaredValue),
		time.Now().Format("2006-01-02 15:04"),
	)
	pdf.MultiCell(wMm-6, 1.3*scale, tr(details), "", "L", false)
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

	pdf.Ln(0.5)
	pdf.SetDrawColor(20, 40, 90)
	pdf.SetLineWidth(0.3)
	pdf.Line(3, pdf.GetY(), wMm-3, pdf.GetY())
	pdf.Ln(0.5)

	pdf.SetFont("Helvetica", "B", 4.5*scale)
	pdf.SetTextColor(20, 40, 90)
	pdf.CellFormat(0, 1.8*scale, tr("REMITE: "+strings.ToUpper(strings.TrimSpace(c.WarehouseCompany))), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 3.5*scale)
	pdf.SetTextColor(0, 0, 0)
	if c.WarehouseAddress != "" {
		pdf.CellFormat(0, 1.5*scale, tr(c.WarehouseAddress), "", 1, "L", false, 0, "")
	}

	pdf.SetFont("Helvetica", "B", 4*scale)
	pdf.SetTextColor(20, 40, 90)
	pdf.CellFormat(0, 1.6*scale, tr("COLVANES S.A.S. - ENVIA | NIT: 800.185.306-4"), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 2.8*scale)
	pdf.SetTextColor(80, 80, 80)
	weight := "N/D"
	if c.Weight > 0 {
		weight = fmt.Sprintf("%.1f kg", c.Weight)
	}
	dims := "N/D"
	if c.Length > 0 || c.Width > 0 || c.Height > 0 {
		dims = fmt.Sprintf("%.0fx%.0fx%.0f cm", c.Length, c.Width, c.Height)
	}
	details := joinNonEmptyProb(" | ", "Peso: "+weight, "Dim: "+dims)
	pdf.CellFormat(0, 1.4*scale, tr(details), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
}

func drawProbHeader(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, scale float64) {
	pageW := pageWidth(pdf)
	yStart := pdf.GetY()

	logoH := 7.0 * scale
	logoW := logoH * 4.7
	if logoW > pageW*0.45 {
		logoW = pageW * 0.45
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

	drawCarrierLogo(pdf, tr, c.Carrier, pageW-3, yStart, scale)

	pdf.SetY(yStart + logoH + 1)
	pdf.SetTextColor(20, 40, 90)
	pdf.SetFont("Helvetica", "B", 10*scale)
	pdf.CellFormat(0, 4.5*scale, tr("GUIA DE TRANSPORTE"), "", 1, "C", false, 0, "")

	pdf.SetTextColor(0, 0, 0)
	business := strings.ToUpper(strings.TrimSpace(c.BusinessName))
	if business == "" {
		business = strings.ToUpper(strings.TrimSpace(c.WarehouseCompany))
	}
	if business != "" {
		pdf.SetFont("Helvetica", "B", 8*scale)
		pdf.CellFormat(0, 3.5*scale, tr(business), "", 1, "L", false, 0, "")
	}

	line2 := []string{}
	if c.OrderNumber != "" {
		line2 = append(line2, "PEDIDO: "+c.OrderNumber)
	}
	if c.TrackingNumber != "" {
		line2 = append(line2, "GUIA: "+c.TrackingNumber)
	}
	if len(line2) > 0 {
		pdf.SetFont("Helvetica", "", 7*scale)
		pdf.SetTextColor(80, 80, 80)
		pdf.CellFormat(0, 3*scale, tr(strings.Join(line2, "   |   ")), "", 1, "L", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
	}

	pdf.SetDrawColor(20, 40, 90)
	pdf.SetLineWidth(0.4)
	pdf.Line(3, pdf.GetY()+0.5, pageW-3, pdf.GetY()+0.5)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.2)
	pdf.Ln(1.5 * scale)
}

func drawCarrierLogo(pdf *gofpdf.Fpdf, tr func(string) string, carrier string, rightX, y, scale float64) {
	logo := getCarrierLogoBytes(carrier)
	if len(logo) == 0 {
		drawCarrierBadge(pdf, tr, carrier, rightX, y, scale)
		return
	}
	opts := gofpdf.ImageOptions{ImageType: "PNG"}
	name := "carrier-logo-" + carrierLogoKey(carrier)
	pdf.RegisterImageOptionsReader(name, opts, bytes.NewReader(logo))
	info := pdf.GetImageInfo(name)
	boxH := 9.0 * scale
	boxW := 22.0 * scale
	dispW := boxW
	dispH := boxH
	if info != nil && info.Width() > 0 && info.Height() > 0 {
		ratio := info.Width() / info.Height()
		dispH = boxH
		dispW = dispH * ratio
		if dispW > boxW {
			dispW = boxW
			dispH = dispW / ratio
		}
	}
	x := rightX - dispW
	yImg := y + (boxH-dispH)/2
	pdf.ImageOptions(name, x, yImg, dispW, dispH, false, opts, 0, "")
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

func drawProbSender(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, scale float64) {
	sender := strings.ToUpper(strings.TrimSpace(c.WarehouseCompany))
	if sender == "" {
		sender = strings.ToUpper(strings.TrimSpace(c.BusinessName))
	}
	lines := []string{}
	if v := strings.TrimSpace(c.WarehouseContact); v != "" {
		lines = append(lines, "Contacto: "+v)
	}
	if v := strings.TrimSpace(c.WarehouseAddress); v != "" {
		lines = append(lines, v)
	}
	if cityState := joinNonEmptyProb(", ", c.WarehouseCity, c.WarehouseState); cityState != "" {
		lines = append(lines, cityState)
	}
	if v := strings.TrimSpace(c.WarehousePhone); v != "" {
		lines = append(lines, "Tel: "+v)
	}
	if sender == "" && len(lines) == 0 {
		return
	}

	pageW := pageWidth(pdf)
	pdf.SetFont("Helvetica", "B", 7*scale)
	pdf.SetFillColor(20, 40, 90)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(pageW-6, 4*scale, tr("REMITENTE"), "", 1, "L", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(0, 0, 0)

	if sender != "" {
		pdf.SetFont("Helvetica", "B", 8*scale)
		pdf.CellFormat(0, 3.5*scale, tr(sender), "", 1, "L", false, 0, "")
	}
	pdf.SetFont("Helvetica", "", 7*scale)
	for _, l := range lines {
		pdf.CellFormat(0, 3*scale, tr(l), "", 1, "L", false, 0, "")
	}
	pdf.Ln(1 * scale)
}

func drawProbRecipient(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, scale float64) {
	pageW := pageWidth(pdf)
	pdf.SetFont("Helvetica", "B", 7*scale)
	pdf.SetFillColor(20, 40, 90)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(pageW-6, 4*scale, tr("DESTINATARIO"), "", 1, "L", true, 0, "")
	pdf.SetFillColor(255, 255, 255)
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
	pdf.CellFormat(0, 3*scale, tr("CODIGO DE GUIA - "+strings.ToUpper(c.Carrier)), "", 1, "C", false, 0, "")
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

func drawProbDetailsBox(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, scale float64) {
	pageW := pageWidth(pdf)
	x := 3.0
	boxW := pageW - 6
	colW := boxW / 4

	type cell struct {
		label string
		value string
	}
	currency := c.Currency
	if currency == "" {
		currency = "COP"
	}
	weight := "N/D"
	if c.Weight > 0 {
		weight = fmt.Sprintf("%.1f kg", c.Weight)
	}
	dims := "N/D"
	if c.Length > 0 || c.Width > 0 || c.Height > 0 {
		dims = fmt.Sprintf("%.0fx%.0fx%.0f", c.Length, c.Width, c.Height)
	}
	declared := "N/D"
	if c.DeclaredValue > 0 {
		declared = "$" + formatMoneyProb(c.DeclaredValue)
	}
	cells := []cell{
		{"PESO", weight},
		{"DIM (cm)", dims},
		{"VLR DECLARADO", declared},
		{"FECHA", time.Now().Format("2006-01-02")},
	}
	if c.CodTotal > 0 {
		cells[2] = cell{"CONTRA ENTREGA", "$" + formatMoneyProb(c.CodTotal)}
	}

	y := pdf.GetY()
	labelH := 3.2 * scale
	valueH := 4.0 * scale

	pdf.SetDrawColor(120, 120, 120)
	pdf.SetLineWidth(0.2)
	pdf.SetFont("Helvetica", "B", 5.5*scale)
	pdf.SetFillColor(235, 238, 245)
	for i, cl := range cells {
		cx := x + float64(i)*colW
		pdf.Rect(cx, y, colW, labelH, "FD")
		pdf.SetXY(cx, y)
		pdf.SetTextColor(60, 60, 60)
		pdf.CellFormat(colW, labelH, tr(cl.label), "", 0, "C", false, 0, "")
	}
	pdf.SetFont("Helvetica", "B", 7.5*scale)
	for i, cl := range cells {
		cx := x + float64(i)*colW
		pdf.Rect(cx, y+labelH, colW, valueH, "D")
		pdf.SetXY(cx, y+labelH)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(colW, valueH, tr(cl.value), "", 0, "C", false, 0, "")
	}
	pdf.SetFillColor(255, 255, 255)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetY(y + labelH + valueH + 1.5*scale)
}

func drawProbProofOfDelivery(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, scale float64) {
	pageW := pageWidth(pdf)
	x := 3.0
	boxW := pageW - 6
	y := pdf.GetY()
	boxH := 14.0 * scale

	pdf.SetDrawColor(120, 120, 120)
	pdf.SetLineWidth(0.2)
	pdf.Rect(x, y, boxW, boxH, "D")

	pdf.SetXY(x+1.5, y+1)
	pdf.SetFont("Helvetica", "B", 6.5*scale)
	pdf.SetTextColor(60, 60, 60)
	pdf.CellFormat(boxW-3, 3*scale, tr("PRUEBA DE ENTREGA"), "", 1, "L", false, 0, "")

	lineY := y + boxH - 4.5*scale
	half := (boxW - 6) / 2
	pdf.SetDrawColor(0, 0, 0)
	pdf.Line(x+2, lineY, x+2+half, lineY)
	pdf.Line(x+boxW-2-half, lineY, x+boxW-2, lineY)

	pdf.SetXY(x+2, lineY+0.3)
	pdf.SetFont("Helvetica", "", 5.5*scale)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(half, 3*scale, tr("RECIBIDO POR (NOMBRE / CC)"), "", 0, "L", false, 0, "")
	pdf.SetX(x + boxW - 2 - half)
	pdf.CellFormat(half, 3*scale, tr("FIRMA / FECHA"), "", 0, "L", false, 0, "")

	pdf.SetDrawColor(0, 0, 0)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetY(y + boxH + 1*scale)
}

func drawProbFooter(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, scale float64) {
	_, pageH := pdf.GetPageSize()
	pageW := pageWidth(pdf)

	pdf.SetY(pageH - 24*scale)
	pdf.SetDrawColor(20, 40, 90)
	pdf.SetLineWidth(0.3)
	pdf.Line(3, pdf.GetY(), pageW-3, pdf.GetY())
	pdf.Ln(0.5)

	pdf.SetFont("Helvetica", "B", 5.5*scale)
	pdf.SetTextColor(20, 40, 90)
	pdf.CellFormat(pageW-6, 2.5*scale, tr("COLVANES S.A.S. - ENVIA"), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 4.5*scale)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(pageW-6, 2*scale, tr("NIT: 800.185.306-4  |  Carrera 88 # 17B-10 Bogota"), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 3.5*scale)
	pdf.SetTextColor(60, 60, 60)
	pdf.CellFormat(pageW-6, 1.6*scale, tr("Tel: (1) 7943670 | www.envia.co | administradorpqr1@enviacolvanes.com.co"), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 3*scale)
	pdf.SetTextColor(80, 80, 80)
	lic1 := "Lic Min Transporte 0080 (14/3/2000) | Lic MinTIC 001368 (4/8/2020) | CIIU 5320 Mensajeria Express"
	pdf.MultiCell(pageW-6, 1.4*scale, tr(lic1), "", "L", false)

	pdf.SetFont("Helvetica", "", 3*scale)
	pdf.SetTextColor(80, 80, 80)
	lic2 := "Autorretenedores Res 4327 (Jul/97) | Grandes Contribuyentes Res 9061 (Dic/20) | Agente Retenedor de IVA"
	pdf.MultiCell(pageW-6, 1.4*scale, tr(lic2), "", "L", false)

	pdf.Ln(0.3 * scale)
	legal := "ESTE CONTRATO DE TRANSPORTE SE RIGE POR EL DECRETO 229 DE 1995 Y NORMAS QUE LO MODIFIQUEN, Y POR LOS ARTICULOS 981 Y SIGUIENTES DEL CODIGO DE COMERCIO. EL REMITENTE DECLARA QUE LA INFORMACION DE ESTA GUIA ES VERIDICA Y QUE LA MERCANCIA NO CONTIENE ARTICULOS PROHIBIDOS O DE TENENCIA RESTRINGIDA. EL VALOR DECLARADO DETERMINA EL LIMITE DE RESPONSABILIDAD DEL TRANSPORTADOR. ESTE ES UN SERVICIO DE MENSAJERIA EXPRESA. RECLAMACIONES DENTRO DE LOS TERMINOS LEGALES."
	pdf.SetFont("Helvetica", "", 3*scale)
	pdf.SetTextColor(100, 100, 100)
	pdf.MultiCell(pageW-6, 1.4*scale, tr(legal), "", "J", false)

	pdf.SetFont("Helvetica", "B", 3.5*scale)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(0, 1.8*scale, tr("Generada por Probability "+time.Now().Format("2006-01-02 15:04")), "", 1, "L", false, 0, "")
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
