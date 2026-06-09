package usecaseshipment

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func buildProbabilityLabel(c *domain.GuidePDFContext, format *domain.GuideFormat) ([]byte, error) {
	carrier := strings.ToUpper(strings.TrimSpace(c.Carrier))

	if carrier == "ENVIA" {
		return buildEnviaLabel(c, format)
	}

	if carrier == "COORDINADORA" {
		return buildCoordinadoraLabel(c, format)
	}

	return buildGenericCarrierLabel(c, format)
}

func buildEnviaLabel(c *domain.GuidePDFContext, format *domain.GuideFormat) ([]byte, error) {
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

func buildGenericCarrierLabel(c *domain.GuidePDFContext, format *domain.GuideFormat) ([]byte, error) {
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
	pageW := wCm*10 - 6
	scale := 1.0
	if wCm < 8 {
		scale = 0.75
	} else if wCm >= 20 {
		scale = 1.5
	}

	drawGenericCarrierHeader(pdf, tr, c, pageW, scale)
	drawGenericSenderRecipient(pdf, tr, c, pageW, scale)
	drawGenericQRAndBarcode(pdf, tr, c, pageW, scale)

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
	pdf.SetFont("Helvetica", "B", 5.5*scale)
	pdf.SetTextColor(20, 40, 90)
	pdf.CellFormat(leftW-3, 2*scale, tr("REMITE:"), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 6*scale)
	warehouse := strings.ToUpper(strings.TrimSpace(c.WarehouseCompany))
	if warehouse != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 2.2*scale, tr(warehouse), "", 1, "L", false, 0, "")
	}
	pdf.SetFont("Helvetica", "", 5*scale)
	if c.WarehouseCity != "" && c.WarehousePhone != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 1.8*scale, tr(c.WarehouseCity+" | Tel: "+c.WarehousePhone), "", 1, "L", false, 0, "")
	} else if c.WarehouseCity != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 1.8*scale, tr(c.WarehouseCity), "", 1, "L", false, 0, "")
	}

	pdf.Ln(0.5*scale)
	pdf.SetX(3)
	pdf.SetFont("Helvetica", "B", 5.5*scale)
	pdf.SetTextColor(20, 40, 90)
	pdf.CellFormat(leftW-3, 2*scale, tr("PARA:"), "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "B", 6.5*scale)
	pdf.SetX(3)
	pdf.CellFormat(leftW-3, 2.5*scale, tr(c.CustomerName), "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 5*scale)
	if c.CustomerPhone != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 1.8*scale, tr("Tel: "+c.CustomerPhone), "", 1, "L", false, 0, "")
	}
	pdf.SetFont("Helvetica", "B", 5.5*scale)
	pdf.SetX(3)
	pdf.MultiCell(leftW-3, 2.2*scale, tr(c.DestinationAddress), "", "L", false)
	cityState := joinNonEmptyProb(", ", c.DestinationCity, c.DestinationState)
	if cityState != "" {
		pdf.SetX(3)
		pdf.CellFormat(leftW-3, 1.8*scale, tr(cityState), "", 1, "L", false, 0, "")
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
	addr := strings.TrimSpace(c.WarehouseAddress)
	if addr == "" {
		addr = strings.TrimSpace(c.BusinessAddress)
	}
	if addr != "" {
		lines = append(lines, addr)
	}
	city := strings.TrimSpace(c.WarehouseCity)
	state := strings.TrimSpace(c.WarehouseState)
	if cityState := joinNonEmptyProb(", ", city, state); cityState != "" {
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

	if len(lines) > 0 {
		pdf.SetFont("Helvetica", "", 7*scale)
		for _, l := range lines {
			pdf.CellFormat(0, 3*scale, tr(l), "", 1, "L", false, 0, "")
		}
	} else if sender != "" {
		pdf.SetFont("Helvetica", "", 7*scale)
		pdf.CellFormat(0, 3*scale, tr("(Sin información adicional en BD)"), "", 1, "L", false, 0, "")
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
	carrierName := strings.ToUpper(strings.TrimSpace(c.Carrier))
	if carrierName == "" {
		carrierName = "PROBABILITY"
	}
	pdf.CellFormat(pageW-6, 2.5*scale, tr(carrierName), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 4.5*scale)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(pageW-6, 2*scale, tr("COLVANES S.A.S."), "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 4*scale)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(pageW-6, 1.8*scale, tr("NIT: 800.185.306-4  |  Carrera 88 # 17B-10 Bogota"), "", 1, "L", false, 0, "")

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

func drawGenericCarrierHeader(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, pageW float64, scale float64) {
	y := 1.5
	logoH := 7 * scale
	carrierLogoW := logoH * 2.5
	probLogoW := logoH * 4.2

	carrier := strings.ToUpper(strings.TrimSpace(c.Carrier))

	logoImg := getCarrierLogo(carrier)
	if logoImg != nil && len(logoImg) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("carrier_logo.png", opts, bytes.NewReader(logoImg))
		pdf.ImageOptions("carrier_logo.png", 5, y, carrierLogoW, logoH, false, opts, 0, "")
	}

	probLogoImg := readLocalAsset("probability-logo.png")
	if probLogoImg != nil && len(probLogoImg) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("prob_logo.png", opts, bytes.NewReader(probLogoImg))
		rightLogoX := 5 + pageW - probLogoW - 0.5
		pdf.ImageOptions("prob_logo.png", rightLogoX, y, probLogoW, logoH, false, opts, 0, "")
	}

	y += logoH + 1.5

	pdf.SetFont("Helvetica", "B", 7*scale)
	if c.TrackingNumber != "" {
		pdf.SetXY(5, y)
		pdf.CellFormat(pageW, 3*scale, tr("Tracking: "+c.TrackingNumber), "1", 1, "C", false, 0, "")
		y += 3.5*scale
	}
}

func drawGenericSenderRecipient(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, pageW float64, scale float64) {
	y := pdf.GetY() + 2*scale
	colW := (pageW - 2) / 2

	pdf.SetFont("Helvetica", "B", 7*scale)
	pdf.SetXY(5, y)
	pdf.CellFormat(colW, 3*scale, tr("REMITE"), "1", 0, "C", false, 0, "")
	pdf.SetX(5 + colW + 2)
	pdf.CellFormat(colW, 3*scale, tr("DESTINATARIO"), "1", 0, "C", false, 0, "")

	y += 3.5*scale
	pdf.SetFont("Helvetica", "", 5.5*scale)

	senderText := fmt.Sprintf("%s\n%s\n%s\n%s %s\nTel: %s",
		c.WarehouseCompany,
		c.WarehouseAddress,
		c.WarehouseCity,
		c.WarehouseContact,
		c.WarehousePhone,
		c.WarehousePhone,
	)

	recipientText := fmt.Sprintf("%s\n%s\n%s, %s\nTel: %s\nDNI: %s",
		c.CustomerName,
		c.DestinationAddress,
		c.DestinationCity,
		c.DestinationState,
		c.CustomerPhone,
		c.CustomerDNI,
	)

	pdf.SetXY(5, y)
	pdf.MultiCell(colW, 2.2*scale, tr(senderText), "1", "L", false)

	senderH := pdf.GetY() - y

	pdf.SetXY(5+colW+2, y)
	pdf.MultiCell(colW, 2.2*scale, tr(recipientText), "1", "L", false)

	if pdf.GetY() < y+senderH {
		pdf.SetY(y + senderH)
	}

	y = pdf.GetY() + 1.5*scale
	pdf.SetFont("Helvetica", "B", 6*scale)
	pdf.SetXY(5, y)

	createdDate := ""
	if c.CreatedAt != nil {
		createdDate = c.CreatedAt.Format("02/01/2006")
	}
	currency := c.Currency
	if currency == "" {
		currency = "COP"
	}

	infoText := fmt.Sprintf("Orden: %s | Fecha: %s | Valor: %.0f %s", c.OrderNumber, createdDate, c.DeclaredValue, currency)
	pdf.CellFormat(pageW, 3*scale, tr(infoText), "1", 1, "C", false, 0, "")
}

func drawGenericQRAndBarcode(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, pageW float64, scale float64) {
	y := pdf.GetY() + 1*scale
	boxMargin := 0.5 * scale

	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)

	leftX := 5.0
	bcHeight := 10.0 * scale

	if c.TrackingNumber != "" {
		barcodePNG := buildCode128PNGProb(c.TrackingNumber, int(pageW*8), int(bcHeight*8))
		if barcodePNG != nil {
			pdf.Rect(leftX-boxMargin, y-boxMargin, pageW+2*boxMargin, bcHeight+2*boxMargin, "D")

			opts := gofpdf.ImageOptions{ImageType: "PNG"}
			pdf.RegisterImageOptionsReader("bc.png", opts, bytes.NewReader(barcodePNG))
			pdf.ImageOptions("bc.png", leftX, y, pageW, bcHeight, false, opts, 0, "")
		}
	}

	y += bcHeight + 1.5*scale

	pdf.SetFont("Helvetica", "B", 7.5*scale)
	pdf.SetXY(leftX, y)

	dimText := fmt.Sprintf("Peso: %.1f kg | Dim: %.0f x %.0f x %.0f cm",
		c.Weight, c.Length, c.Width, c.Height)
	pdf.CellFormat(pageW, 3.5*scale, tr(dimText), "1", 1, "C", false, 0, "")

	y = pdf.GetY() + 1*scale
	colW := (pageW - 2) / 2
	leftColW := colW
	rightColW := colW
	rightX := leftX + leftColW + 2
	boxHeight := 35*scale

	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.2)
	pdf.Rect(leftX-boxMargin, y-boxMargin, leftColW+2*boxMargin, boxHeight+2*boxMargin, "D")
	pdf.Rect(rightX-boxMargin, y-boxMargin, rightColW+2*boxMargin, boxHeight+2*boxMargin, "D")

	pdf.SetFont("Helvetica", "", 4.5*scale)

	qrSize := rightColW * 0.65

	if qrPNG := buildQRPNGProb(c.TrackingNumber); qrPNG != nil {
		qrCenterX := rightX + (rightColW-qrSize)/2
		qrCenterY := y + (boxHeight-qrSize)/2

		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("qr.png", opts, bytes.NewReader(qrPNG))
		pdf.ImageOptions("qr.png", qrCenterX, qrCenterY, qrSize, qrSize, false, opts, 0, "")
	}

	pdf.SetY(y + boxHeight)
}

func buildCoordinadoraLabel(c *domain.GuidePDFContext, format *domain.GuideFormat) ([]byte, error) {
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
	pdf.SetCellMargin(1.0)

	tr := pdf.UnicodeTranslatorFromDescriptor("")
	scale := 1.0
	pageW := wCm * 10 - 6
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)
	pdf.SetTextColor(0, 0, 0)

	y := 3.0

	logoH := 5.5 * scale

	probLogo := readLocalAsset("probability-logo.png")
	if len(probLogo) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("prob_logo_coord.png", opts, bytes.NewReader(probLogo))
		pdf.ImageOptions("prob_logo_coord.png", 3, y, 0, logoH, true, opts, 0, "")
	}

	y = y + logoH + 1.5

	pdf.SetXY(3, y)
	pdf.SetFont("Helvetica", "B", 6.5*scale)
	colW1 := pageW / 3
	colW2 := pageW / 3
	colW3 := pageW / 3

	pdf.CellFormat(colW1, 7*scale, "Origin\n1\nBOG", "1", 0, "C", false, 0, "")
	pdf.SetX(3 + colW1)
	pdf.CellFormat(colW2, 7*scale, "AS\nPAQ\n1-2", "1", 0, "C", false, 0, "")
	pdf.SetX(3 + colW1 + colW2)
	pdf.CellFormat(colW3, 7*scale, tr("UNIDAD:\n1/1"), "1", 1, "C", false, 0, "")
	y = pdf.GetY() + 2.0

	pdf.SetXY(3, y)
	pdf.SetFont("Helvetica", "B", 5.5*scale)
	pdf.SetFillColor(240, 240, 240)
	colRef := pageW / 2
	colObs := pageW / 2

	pdf.CellFormat(colRef-0.2, 4*scale, tr("Ref:"), "1", 0, "L", true, 0, "")
	pdf.SetX(3 + colRef + 0.2)
	pdf.CellFormat(colObs-0.2, 4*scale, tr("Observaciones Cliente:"), "1", 1, "L", true, 0, "")

	pdf.SetFont("Helvetica", "", 4.5*scale)
	pdf.SetFillColor(255, 255, 255)
	pdf.SetXY(3.5, pdf.GetY())
	refText := "ORDEN\nORD-" + c.OrderNumber
	pdf.MultiCell(colRef-0.5, 2.5*scale, tr(refText), "1", "L", false)

	obsStartY := pdf.GetY()
	pdf.SetXY(3+colRef+0.7, pdf.GetY()-5*scale)
	obsText := tr("CASA 126 Doc.\nORDEN " + c.OrderNumber)
	pdf.MultiCell(colObs-0.5, 2.5*scale, obsText, "1", "L", false)

	maxObsY := pdf.GetY()
	if obsStartY > maxObsY {
		maxObsY = obsStartY
	}
	pdf.SetY(maxObsY)
	y = pdf.GetY() + 1.5

	pdf.SetXY(3, y)
	pdf.SetFont("Helvetica", "B", 5.5*scale)
	colDestino := pageW / 3
	colZona := pageW / 3
	colEquipo := pageW / 3

	pdf.CellFormat(colDestino, 3.5*scale, "Destino\n1", "1", 0, "C", false, 0, "")
	pdf.SetX(3 + colDestino)
	pdf.CellFormat(colZona, 3.5*scale, "Zona Hub", "1", 0, "C", false, 0, "")
	pdf.SetX(3 + colDestino + colZona)
	pdf.CellFormat(colEquipo, 3.5*scale, "Equipo\nReparto", "1", 1, "C", false, 0, "")
	y = pdf.GetY()

	pdf.SetXY(3, y)
	pdf.SetFont("Helvetica", "B", 9*scale)
	pdf.CellFormat(pageW, 3.5*scale, time.Now().Format("2006-01-02"), "1", 1, "C", false, 0, "")
	y = pdf.GetY() + 2.0

	pdf.SetXY(3, y)
	pdf.SetFont("Helvetica", "B", 5.5*scale)
	colRemDest := pageW / 2
	pdf.CellFormat(colRemDest, 3*scale, tr("REMITENTE"), "1", 0, "C", false, 0, "")
	pdf.SetX(3 + colRemDest)
	pdf.CellFormat(colRemDest, 3*scale, tr("DESTINATARIO"), "1", 1, "C", false, 0, "")
	y = pdf.GetY()

	warehouse := c.WarehouseCompany
	if warehouse == "" {
		warehouse = c.BusinessName
	}

	pdf.SetXY(3.5, y)
	pdf.SetFont("Helvetica", "", 3.8*scale)
	remText := tr(warehouse + "\n" + c.WarehouseAddress + "\n" + c.WarehouseCity + "\nTel: " + c.WarehousePhone)
	pdf.MultiCell(colRemDest-0.8, 2.2*scale, remText, "1", "L", false)
	remEndY := pdf.GetY()

	pdf.SetXY(3+colRemDest+0.7, y)
	pdf.SetFont("Helvetica", "", 3.8*scale)
	destText := tr(c.CustomerName + "\n" + c.DestinationAddress + "\n" + c.DestinationCity + "\nTel: " + c.CustomerPhone)
	pdf.MultiCell(colRemDest-0.8, 2.2*scale, destText, "1", "L", false)
	destEndY := pdf.GetY()

	y = remEndY
	if destEndY > y {
		y = destEndY
	}
	y = y + 3.0

	pdf.SetXY(3, y)
	pdf.SetFont("Helvetica", "B", 4.5*scale)
	pdf.CellFormat(pageW, 2.5*scale, "CODIGO DE GUIA", "1", 1, "C", false, 0, "")
	y = pdf.GetY()

	bcImg := buildCode128PNGProb(c.TrackingNumber, int(pageW*8), int(9*scale*8))
	if bcImg != nil {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("coord_bc.png", opts, bytes.NewReader(bcImg))
		pdf.ImageOptions("coord_bc.png", 3, y, pageW, 9*scale, false, opts, 0, "")
		pdf.SetY(y + 9*scale)
	}
	y = pdf.GetY()

	pdf.SetFont("Courier", "B", 9*scale)
	pdf.SetXY(3, y)
	pdf.CellFormat(pageW, 3*scale, c.TrackingNumber, "1", 1, "C", false, 0, "")
	y = pdf.GetY() + 1.5

	pdf.SetXY(3, y)
	pdf.SetFont("Helvetica", "", 4.5*scale)
	dimText := fmt.Sprintf("Peso: %.1f kg | Dim: %.0f x %.0f x %.0f cm", c.Weight, c.Length, c.Width, c.Height)
	pdf.CellFormat(pageW, 2.5*scale, dimText, "1", 1, "C", false, 0, "")
	y = pdf.GetY() + 1.5

	colLogoBox := pageW * 0.3
	colQRBox := pageW * 0.7

	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)

	coordLogo := getCarrierLogo("COORDINADORA")
	logoBoxH := 9.0 * scale

	pdf.Rect(3, y, colLogoBox, logoBoxH, "")
	if len(coordLogo) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("coord_logo_box.png", opts, bytes.NewReader(coordLogo))
		logoImgW := colLogoBox - 1.2
		logoImgH := logoBoxH - 1.2
		logoX := 3 + (colLogoBox-logoImgW)/2
		logoY := y + 0.6
		pdf.ImageOptions("coord_logo_box.png", logoX, logoY, logoImgW, logoImgH, true, opts, 0, "")
	}

	pdf.Rect(3+colLogoBox, y, colQRBox, logoBoxH, "")
	qrImg := buildQRPNGProb(c.TrackingNumber)
	if qrImg != nil {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("qr_coord.png", opts, bytes.NewReader(qrImg))
		qrSize := colQRBox - 1.5
		qrX := 3 + colLogoBox + (colQRBox-qrSize)/2
		qrY := y + logoBoxH - qrSize - 0.8
		pdf.ImageOptions("qr_coord.png", qrX, qrY, qrSize, qrSize, false, opts, 0, "")
	}

	y = y + logoBoxH + 1.0

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getCarrierLogo(carrier string) []byte {
	switch carrier {
	case "COORDINADORA":
		return downloadLogoFromS3("https://images-cam93.s3.us-east-1.amazonaws.com/imagen_coordinadora.png")
	case "INTERRAPIDISIMO":
		return downloadLogoFromS3("https://images-cam93.s3.us-east-1.amazonaws.com/imagen_inerapidisimo.png")
	case "SERVIENTREGA":
		return downloadLogoFromS3("https://images-cam93.s3.us-east-1.amazonaws.com/imagen_servientrega.png")
	default:
		return nil
	}
}

func readLocalAsset(filename string) []byte {
	paths := []string{
		fmt.Sprintf("./services/modules/shipments/internal/app/usecaseshipment/assets/%s", filename),
		fmt.Sprintf("services/modules/shipments/internal/app/usecaseshipment/assets/%s", filename),
		fmt.Sprintf("back/central/services/modules/shipments/internal/app/usecaseshipment/assets/%s", filename),
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil && len(data) > 0 {
			return data
		}
	}
	return nil
}

func downloadLogoFromS3(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return data
}
