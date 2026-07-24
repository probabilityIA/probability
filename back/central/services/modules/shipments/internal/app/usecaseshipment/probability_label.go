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

// buildCoordinadoraLabel replica el layout real de las guias de COORDINADORA
// (logo + GUIA/UNIDAD arriba, codigo de barras, DE:/PARA:, Observaciones,
// Ref, banner "Recaudos Contra Entrega" solo si el envio es contra entrega,
// y la fila inferior Origen/QR/Destino/Zona Hub/Equipo Reparto), usando los
// datos que ExtractCoordinadoraMetadata ya extrae del PDF real del carrier.
func buildCoordinadoraLabel(c *domain.GuidePDFContext, format *domain.GuideFormat) ([]byte, error) {
	wCm := format.WidthCm
	hCm := format.HeightCm
	if wCm < 4 || hCm < 4 {
		wCm = 10
		hCm = 15
	}
	wMm := wCm * 10.0
	hMm := hCm * 10.0

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: wMm, Ht: hMm},
		OrientationStr: "P",
	})
	margin := 3.0
	pdf.SetMargins(margin, margin, margin)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	scale := wMm / 100.0
	if scale < 0.6 {
		scale = 0.6
	} else if scale > 1.6 {
		scale = 1.6
	}

	usableW := wMm - margin*2
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	gray := color.RGBA{R: 90, G: 90, B: 90, A: 255}

	pdf.SetDrawColor(0, 0, 0)
	pdf.SetTextColor(0, 0, 0)

	// ── Header: [logo / barcode] | [GUIA, caja alta] | [UNIDAD / CC / PAQ] ─
	gap := 1.5
	topH := 24.0 * scale
	logoBarcodeW := usableW * 0.45
	guiaW := usableW * 0.27
	rightColW := usableW - logoBarcodeW - guiaW - gap*2

	x1 := margin
	x2 := x1 + logoBarcodeW + gap
	x3 := x2 + guiaW + gap

	logoH := topH*0.4 - 1
	barcodeY := margin + logoH + 1.5*scale
	barcodeH := topH - logoH - 1.5*scale

	// El logo fuente de COORDINADORA ya trae el icono + el nombre en una
	// sola imagen (wordmark completo). Se dibuja respetando su proporcion
	// real (contain, sin deformar), del alto disponible o del ancho
	// disponible, lo que resulte mas chico.
	logoBytes := getCarrierLogoBytes(c.Carrier)
	if len(logoBytes) > 0 {
		logoW := logoBarcodeW
		drawH := logoH
		if cfg, _, err := image.DecodeConfig(bytes.NewReader(logoBytes)); err == nil && cfg.Width > 0 && cfg.Height > 0 {
			aspect := float64(cfg.Width) / float64(cfg.Height)
			drawH = logoH
			logoW = drawH * aspect
			if logoW > logoBarcodeW {
				logoW = logoBarcodeW
				drawH = logoW / aspect
			}
		}
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("coord-logo", opts, bytes.NewReader(logoBytes))
		pdf.ImageOptions("coord-logo", x1, margin+(logoH-drawH)/2, logoW, drawH, false, opts, 0, "")
	} else {
		pdf.SetXY(x1, margin)
		pdf.SetFont("Helvetica", "B", 11*scale)
		pdf.SetTextColor(int(black.R), int(black.G), int(black.B))
		pdf.CellFormat(logoBarcodeW, logoH, tr("COORDINADORA"), "", 0, "L", false, 0, "")
	}

	barcodePNG := buildCode128PNGProb(strings.TrimSpace(c.TrackingNumber), int(logoBarcodeW*8), int(barcodeH*4))
	if len(barcodePNG) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		key := fmt.Sprintf("coord-barcode-%d", time.Now().UnixNano())
		pdf.RegisterImageOptionsReader(key, opts, bytes.NewReader(barcodePNG))
		pdf.ImageOptions(key, x1, barcodeY, logoBarcodeW, barcodeH, false, opts, 0, "")
	}

	guiaVal := strings.TrimSpace(c.Guia)
	if guiaVal == "" {
		guiaVal = strings.TrimSpace(c.TrackingNumber)
	}
	drawCoordBox(pdf, tr, "GUIA:", guiaVal, x2, margin, guiaW, topH, scale, false)

	unidadVal := strings.TrimSpace(c.Unidad)
	if unidadVal == "" {
		unidadVal = "1/1"
	}
	paqVal := strings.TrimSpace(c.Paq)
	paqTxt := "PAQ"
	if paqVal != "" {
		paqTxt = "PAQ " + paqVal
	}
	cellH := topH / 3
	drawCoordBadge(pdf, tr, unidadVal, x3, margin, rightColW, cellH, scale, false)
	// "CC" es un badge fijo que COORDINADORA imprime en sus guias (tipo de
	// documento); no proviene de datos extraidos, se replica tal cual.
	drawCoordBadge(pdf, tr, "CC", x3, margin+cellH, rightColW, cellH, scale, false)
	drawCoordBadge(pdf, tr, paqTxt, x3, margin+cellH*2, rightColW, cellH, scale, true)

	y := margin + topH + 1.5*scale
	pdf.SetDrawColor(int(black.R), int(black.G), int(black.B))
	pdf.SetLineWidth(0.3)
	pdf.Line(margin, y, margin+usableW, y)
	y += 1.5 * scale

	// ── DE: ─────────────────────────────────────────────────────────────
	sender := strings.TrimSpace(c.WarehouseCompany)
	if sender == "" {
		sender = strings.TrimSpace(c.BusinessName)
	}
	senderAddress := strings.TrimSpace(c.WarehouseAddress)
	if senderAddress == "" {
		// Sin bodega ni direccion de origen configurada: usar la direccion
		// general del negocio como ultimo recurso, mejor que dejarlo vacio.
		senderAddress = strings.TrimSpace(c.BusinessAddress)
	}
	senderCity := joinNonEmptyProb(", ", strings.TrimSpace(c.WarehouseCity), strings.TrimSpace(c.WarehouseState))
	senderLine2 := joinNonEmptyProb("  ", cityTelLine(senderCity, c.WarehousePhone), postalLine(c.WarehousePostal))
	y = drawCoordAddressBlock(pdf, tr, "DE:", sender, senderAddress, senderLine2, margin, y, usableW, scale, black, gray)

	// ── PARA: ───────────────────────────────────────────────────────────
	recipient := strings.TrimSpace(c.CustomerName)
	destCity := joinNonEmptyProb(", ", strings.TrimSpace(c.DestinationCity), strings.TrimSpace(c.DestinationState))
	destLine2 := cityTelLine(destCity, c.CustomerPhone)
	y = drawCoordAddressBlock(pdf, tr, "PARA:", recipient, strings.TrimSpace(c.DestinationAddress), destLine2, margin, y, usableW, scale, black, gray)

	// ── Observaciones Cliente: ─────────────────────────────────────────
	if obs := strings.TrimSpace(c.Observaciones); obs != "" {
		pdf.SetXY(margin, y)
		pdf.SetFont("Helvetica", "B", 7.5*scale)
		pdf.CellFormat(usableW, 3.2*scale, tr("Observaciones Cliente:"), "", 1, "L", false, 0, "")
		pdf.SetX(margin)
		pdf.SetFont("Helvetica", "", 7*scale)
		pdf.MultiCell(usableW, 3.2*scale, tr(obs), "", "L", false)
		y = pdf.GetY() + 1.2
	}

	// ── Ref: ────────────────────────────────────────────────────────────
	refVal := strings.TrimSpace(c.Ref)
	if refVal == "" {
		refVal = strings.TrimSpace(c.OrderNumber)
	}
	if refVal != "" {
		refLabelW := 10 * scale
		pdf.SetXY(margin, y)
		pdf.SetFont("Helvetica", "BU", 7.5*scale)
		pdf.CellFormat(refLabelW, 3.2*scale, tr("Ref:"), "", 0, "L", false, 0, "")
		pdf.SetFont("Helvetica", "", 7*scale)
		pdf.CellFormat(usableW-refLabelW, 3.2*scale, tr(refVal), "", 1, "L", false, 0, "")
		y = pdf.GetY() + 1.2
	}

	// ── Banner "Recaudos Contra Entrega" (solo COD) ────────────────────
	if c.CodTotal > 0 {
		bannerH := 7 * scale
		pdf.SetFillColor(0, 0, 0)
		pdf.Rect(margin, y, usableW, bannerH, "F")
		pdf.SetTextColor(255, 255, 255)
		pdf.SetXY(margin, y+1.3*scale)
		pdf.SetFont("Helvetica", "B", 8*scale)
		pdf.CellFormat(usableW, 3.8*scale, tr("Recaudos Contra Entrega"), "", 1, "C", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
		y += bannerH + 2
	}

	// Linea divisoria entre la info de la orden (DE:/PARA:/Observaciones/Ref)
	// y la info de origen/destino de la guia. Si hay banner de "Recaudos
	// Contra Entrega" ya cumple ese rol, no hace falta la linea extra.
	if c.CodTotal <= 0 {
		pdf.SetDrawColor(int(black.R), int(black.G), int(black.B))
		pdf.SetLineWidth(0.5)
		pdf.Line(margin, y, margin+usableW, y)
		y += 2
	}

	// ── Fila inferior: Origen | QR | Destino / Zona Hub / Equipo Reparto ─
	footerH := hMm - y - margin
	minFooterH := 22 * scale
	maxFooterH := 30 * scale
	if footerH < minFooterH {
		footerH = minFooterH
		y = hMm - margin - footerH
	} else if footerH > maxFooterH {
		// Sobra espacio vertical (formatos altos como 10x15): no estirar las
		// cajas, dejar el resto de la pagina en blanco como en la guia real.
		footerH = maxFooterH
	}
	drawCoordFooterRow(pdf, tr, c, margin, y, usableW, footerH, scale, black)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// drawCoordBox dibuja una caja con borde, etiqueta pequena arriba y el
// valor en negrita debajo — el estilo de "GUIA:" / "UNIDAD" / "Origen" /
// "Destino" / "Zona Hub" / "Equipo Reparto" que usa COORDINADORA.
func drawCoordBox(pdf *gofpdf.Fpdf, tr func(string) string, label, value string, x, y, w, h, scale float64, centered bool) {
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)
	pdf.Rect(x, y, w, h, "D")

	align := "L"
	pad := 1.2 * scale
	if centered {
		align = "C"
		pad = 0
	}

	pdf.SetXY(x+pad, y+0.8*scale)
	pdf.SetFont("Helvetica", "", 5.5*scale)
	pdf.CellFormat(w-pad*2, 2.4*scale, tr(label), "", 1, align, false, 0, "")

	pdf.SetX(x + pad)
	pdf.SetFont("Helvetica", "B", 8.5*scale)
	pdf.CellFormat(w-pad*2, h-4.5*scale, tr(value), "", 0, align, false, 0, "")
}

// drawCoordBadge dibuja una caja pequena de una sola linea, centrada,
// opcionalmente con fondo negro y texto blanco (como el badge "PAQ 1-2").
func drawCoordBadge(pdf *gofpdf.Fpdf, tr func(string) string, text string, x, y, w, h, scale float64, filled bool) {
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)
	if filled {
		pdf.SetFillColor(0, 0, 0)
		pdf.Rect(x, y, w, h, "FD")
		pdf.SetTextColor(255, 255, 255)
	} else {
		pdf.Rect(x, y, w, h, "D")
	}
	pdf.SetFont("Helvetica", "B", 7.5*scale)
	lineH := 3.6 * scale
	pdf.SetXY(x, y+(h-lineH)/2)
	pdf.CellFormat(w, lineH, tr(text), "", 0, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
}

// drawCoordAddressBlock dibuja un bloque "DE:"/"PARA:" con nombre en
// negrita, direccion y una linea de ciudad/telefono/postal. Retorna el Y
// donde termino de escribir, para encadenar el siguiente bloque.
func drawCoordAddressBlock(pdf *gofpdf.Fpdf, tr func(string) string, label, name, address, line2 string, x, y, w, scale float64, textCol, mutedCol color.RGBA) float64 {
	pdf.SetXY(x, y)
	pdf.SetFont("Helvetica", "B", 8*scale)
	pdf.SetTextColor(int(textCol.R), int(textCol.G), int(textCol.B))
	labelTxt := label
	if name != "" {
		labelTxt = label + " " + name
	}
	pdf.MultiCell(w, 3.6*scale, tr(labelTxt), "", "L", false)
	y = pdf.GetY()

	if address != "" {
		pdf.SetX(x)
		pdf.SetFont("Helvetica", "", 7*scale)
		pdf.SetTextColor(int(textCol.R), int(textCol.G), int(textCol.B))
		pdf.MultiCell(w, 3.2*scale, tr(address), "", "L", false)
		y = pdf.GetY()
	}

	if line2 != "" {
		pdf.SetX(x)
		pdf.SetFont("Helvetica", "", 6.5*scale)
		pdf.SetTextColor(int(mutedCol.R), int(mutedCol.G), int(mutedCol.B))
		pdf.CellFormat(w, 3*scale, tr(line2), "", 1, "L", false, 0, "")
		y = pdf.GetY()
	}

	pdf.SetTextColor(0, 0, 0)
	return y + 1*scale
}

// drawCoordFooterRow dibuja la fila inferior de la guia: caja Origen (numero
// + codigo de hub en dos lineas), QR de tracking con la fecha debajo, caja
// Destino con el valor grande, y una sola caja Zona Hub/Equipo Reparto
// dividida en dos mitades — igual que en la guia real de COORDINADORA.
func drawCoordFooterRow(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, x, y, w, h, scale float64, black color.RGBA) {
	gap := 2.0
	origenW := w * 0.17
	qrColW := w * 0.22
	destinoW := w * 0.20
	statW := w - origenW - qrColW - destinoW - gap*3

	boxH := h

	origenParts := strings.Fields(strings.TrimSpace(c.Origen))
	origenLine1, origenLine2 := "-", ""
	if len(origenParts) > 0 {
		origenLine1 = origenParts[0]
	}
	if len(origenParts) > 1 {
		origenLine2 = strings.Join(origenParts[1:], " ")
	}
	drawCoordBoxTwoLine(pdf, tr, "Origen", origenLine1, origenLine2, x, y, origenW, boxH, scale)

	qrX := x + origenW + gap
	qrSize := qrColW
	dateH := 5.5 * scale
	if boxH-dateH-1 < qrSize {
		qrSize = boxH - dateH - 1
	}
	qrData := strings.TrimSpace(c.TrackingNumber)
	qrPNG := buildQRPNGProb(qrData)
	if len(qrPNG) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		key := fmt.Sprintf("coord-qr-%d", time.Now().UnixNano())
		pdf.RegisterImageOptionsReader(key, opts, bytes.NewReader(qrPNG))
		pdf.ImageOptions(key, qrX+(qrColW-qrSize)/2, y, qrSize, qrSize, false, opts, 0, "")
	}

	dateStr := ""
	if c.CreatedAt != nil {
		dateStr = c.CreatedAt.Format("2006-01-02")
	}
	pdf.SetDrawColor(int(black.R), int(black.G), int(black.B))
	pdf.SetLineWidth(0.3)
	pdf.Rect(qrX, y+boxH-dateH, qrColW, dateH, "D")
	pdf.SetXY(qrX, y+boxH-dateH+(dateH-2.8*scale)/2)
	pdf.SetFont("Helvetica", "", 5.5*scale)
	pdf.CellFormat(qrColW, 2.8*scale, tr(dateStr), "", 0, "C", false, 0, "")

	destinoX := qrX + qrColW + gap
	destinoVal := strings.TrimSpace(c.Destino)
	if destinoVal == "" {
		destinoVal = "-"
	}
	drawCoordBoxBig(pdf, tr, "Destino", destinoVal, destinoX, y, destinoW, boxH, scale)

	statX := destinoX + destinoW + gap
	zonaVal := strings.TrimSpace(c.ZonaHub)
	if zonaVal == "" {
		zonaVal = "-"
	}
	equipoVal := strings.TrimSpace(c.EquipoReparto)
	if equipoVal == "" {
		equipoVal = "-"
	}
	drawCoordSplitBox(pdf, tr, "Zona Hub", zonaVal, "Equipo Reparto", equipoVal, statX, y, statW, boxH, scale)
}

// drawCoordBoxTwoLine es como drawCoordBox pero el valor va en dos lineas
// grandes centradas (el numero de Origen y su codigo de hub, ej. "28"/"FLA").
func drawCoordBoxTwoLine(pdf *gofpdf.Fpdf, tr func(string) string, label, line1, line2 string, x, y, w, h, scale float64) {
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)
	pdf.Rect(x, y, w, h, "D")

	pdf.SetXY(x, y+1*scale)
	pdf.SetFont("Helvetica", "", 5.5*scale)
	pdf.CellFormat(w, 2.4*scale, tr(label), "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "B", 11.5*scale)
	midY := y + h/2
	if line2 != "" {
		pdf.SetXY(x, midY-5*scale)
		pdf.CellFormat(w, 5*scale, tr(line1), "", 1, "C", false, 0, "")
		pdf.SetX(x)
		pdf.CellFormat(w, 5*scale, tr(line2), "", 0, "C", false, 0, "")
	} else {
		pdf.SetXY(x, midY-2.5*scale)
		pdf.CellFormat(w, 5*scale, tr(line1), "", 0, "C", false, 0, "")
	}
}

// drawCoordBoxBig es como drawCoordBox pero con el valor en un tamano de
// fuente mucho mayor (el "Destino" real se imprime muy grande y bold).
func drawCoordBoxBig(pdf *gofpdf.Fpdf, tr func(string) string, label, value string, x, y, w, h, scale float64) {
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)
	pdf.Rect(x, y, w, h, "D")

	pdf.SetXY(x, y+1*scale)
	pdf.SetFont("Helvetica", "B", 6*scale)
	pdf.CellFormat(w, 2.6*scale, tr(label), "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "B", 16.5*scale)
	pdf.SetXY(x, y+h/2-4.5*scale)
	pdf.CellFormat(w, 9*scale, tr(value), "", 0, "C", false, 0, "")
}

// drawCoordSplitBox dibuja una sola caja partida por una linea horizontal
// en dos mitades, cada una con su propia etiqueta y valor — el estilo real
// de la caja combinada "Zona Hub" / "Equipo Reparto".
func drawCoordSplitBox(pdf *gofpdf.Fpdf, tr func(string) string, label1, val1, label2, val2 string, x, y, w, h, scale float64) {
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.3)
	pdf.Rect(x, y, w, h, "D")
	halfH := h / 2
	pdf.Line(x, y+halfH, x+w, y+halfH)

	drawCoordHalfCell(pdf, tr, label1, val1, x, y, w, halfH, scale)
	drawCoordHalfCell(pdf, tr, label2, val2, x, y+halfH, w, halfH, scale)
}

func drawCoordHalfCell(pdf *gofpdf.Fpdf, tr func(string) string, label, value string, x, y, w, h, scale float64) {
	pdf.SetXY(x, y+0.8*scale)
	pdf.SetFont("Helvetica", "", 5.5*scale)
	pdf.CellFormat(w, 2.4*scale, tr(label), "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "B", 12*scale)
	pdf.SetXY(x, y+h-7*scale)
	pdf.CellFormat(w, 5*scale, tr(value), "", 0, "C", false, 0, "")
}

// cityTelLine junta "ciudad" y "Tel: numero" con un separador visual,
// omitiendo cualquiera de los dos que venga vacio.
func cityTelLine(city, phone string) string {
	phone = strings.TrimSpace(phone)
	if phone != "" {
		phone = "TEL " + phone
	}
	return joinNonEmptyProb("  ", city, phone)
}

// postalLine antepone "Z.Postal:" solo si hay un valor que mostrar.
func postalLine(postal string) string {
	postal = strings.TrimSpace(postal)
	if postal == "" {
		return ""
	}
	return "Z.Postal: " + postal
}

func buildGenericCarrierLabel(c *domain.GuidePDFContext, format *domain.GuideFormat) ([]byte, error) {
	wCm := format.WidthCm
	hCm := format.HeightCm
	if wCm < 4 || hCm < 4 {
		wCm = 10
		hCm = 15
	}

	wMm := wCm * 10.0
	hMm := hCm * 10.0

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:        "mm",
		Size:           gofpdf.SizeType{Wd: wMm, Ht: hMm},
		OrientationStr: "P",
	})
	pdf.SetMargins(0, 0, 0)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	scale := 1.0
	if wMm < 80 {
		scale = 0.65
	} else if wMm >= 210 {
		scale = 1.3
	}

	margin := 4.0

	drawProfessionalHeaderNewDesign(pdf, tr, c, wMm, hMm, scale)
	drawThreeColumnsSection(pdf, tr, c, wMm, hMm, margin, scale)
	drawTrackingSectionNewDesign(pdf, tr, c, wMm, hMm, margin, scale)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func drawProfessionalHeaderNewDesign(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, wMm, hMm, scale float64) {
	headerH := hMm * 0.25
	darkBlue := color.RGBA{R: 22, G: 54, B: 92, A: 255}
	darkerBlue := color.RGBA{R: 15, G: 39, B: 66, A: 255}
	lightGray := color.RGBA{R: 159, G: 180, B: 206, A: 255}

	pdf.SetFillColor(int(darkBlue.R), int(darkBlue.G), int(darkBlue.B))
	pdf.Rect(0, 0, wMm, headerH, "F")

	pdf.SetXY(0, 1.5)
	pdf.SetFont("Helvetica", "B", 12*scale)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(wMm, 3*scale, "PROBABILITY", "", 0, "C", false, 0, "")

	pdf.SetXY(wMm-22, 1.5)
	pdf.SetFont("Helvetica", "B", 6.5*scale)
	pdf.SetTextColor(255, 255, 255)
	createdDate := ""
	if c.CreatedAt != nil {
		createdDate = c.CreatedAt.Format("02/01/2006")
	}
	pdf.CellFormat(20, 3*scale, tr(createdDate), "", 0, "R", false, 0, "")

	pdf.SetXY(0, 4.5*scale)
	pdf.SetFont("Helvetica", "", 5.5*scale)
	pdf.SetTextColor(int(lightGray.R), int(lightGray.G), int(lightGray.B))
	pdf.CellFormat(wMm, 1.8*scale, "GESTION DE ENVIOS", "", 1, "C", false, 0, "")

	pdf.SetFillColor(int(darkerBlue.R), int(darkerBlue.G), int(darkerBlue.B))
	pdf.Rect(0, headerH*0.65, wMm, headerH*0.35, "F")

	pdf.SetXY(0, headerH*0.68)
	pdf.SetFont("Helvetica", "B", 15*scale)
	pdf.SetTextColor(255, 255, 255)
	carrier := strings.ToUpper(strings.TrimSpace(c.Carrier))
	pdf.CellFormat(wMm, 4*scale, tr(carrier), "", 1, "C", false, 0, "")

	pdf.SetDrawColor(int(lightGray.R), int(lightGray.G), int(lightGray.B))
	pdf.SetLineWidth(0.5)
	pdf.Line(0, headerH, wMm, headerH)
}

func drawThreeColumnsSection(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, wMm, hMm, margin, scale float64) {
	headerH := hMm * 0.25
	contentH := hMm * 0.50
	usableW := wMm - (margin * 2)
	colW := usableW / 3.0
	startY := headerH

	darkBlue := color.RGBA{R: 22, G: 54, B: 92, A: 255}
	darkGray := color.RGBA{R: 17, G: 17, B: 17, A: 255}
	lightBgGray := color.RGBA{R: 250, G: 251, B: 252, A: 255}
	borderGray := color.RGBA{R: 184, G: 188, B: 194, A: 255}

	pdf.SetDrawColor(int(borderGray.R), int(borderGray.G), int(borderGray.B))
	pdf.SetLineWidth(0.4)

	xCol1 := margin
	xCol2 := margin + colW
	xCol3 := margin + (colW * 2)

	drawColumnRemite(pdf, tr, c, xCol1, startY, colW, contentH, darkBlue, darkGray, borderGray, scale)
	pdf.Line(xCol2-0.5, startY, xCol2-0.5, startY+contentH)

	drawColumnDestino(pdf, tr, c, xCol2, startY, colW, contentH, lightBgGray, darkBlue, darkGray, borderGray, scale)
	pdf.Line(xCol3-0.5, startY, xCol3-0.5, startY+contentH)

	drawColumnPaquete(pdf, tr, c, xCol3, startY, colW, contentH, darkBlue, darkGray, borderGray, scale)

	pdf.SetDrawColor(int(darkBlue.R), int(darkBlue.G), int(darkBlue.B))
	pdf.SetLineWidth(0.6)
	pdf.Line(0, startY+contentH, wMm, startY+contentH)
}

func drawColumnRemite(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, x, y, w, h float64, headerCol, textCol, borderCol color.RGBA, scale float64) {
	margin := 2.5
	currentY := y + 1.5

	pdf.SetXY(x+margin, currentY)
	pdf.SetFont("Helvetica", "B", 7*scale)
	pdf.SetTextColor(int(headerCol.R), int(headerCol.G), int(headerCol.B))
	pdf.CellFormat(w-margin*2, 2*scale, "REMITE", "", 1, "L", false, 0, "")
	currentY += 3.5

	pdf.SetXY(x+margin, currentY)
	pdf.SetFont("Helvetica", "B", 6.5*scale)
	pdf.SetTextColor(int(textCol.R), int(textCol.G), int(textCol.B))
	sender := strings.ToUpper(strings.TrimSpace(c.WarehouseCompany))
	if sender == "" {
		sender = strings.ToUpper(strings.TrimSpace(c.BusinessName))
	}
	pdf.MultiCell(w-margin*2, 1.8*scale, tr(sender), "", "L", false)
	currentY = pdf.GetY() + 0.5

	pdf.SetFont("Helvetica", "", 5*scale)
	pdf.SetXY(x+margin, currentY)
	pdf.MultiCell(w-margin*2, 1.4*scale, tr(strings.TrimSpace(c.WarehouseAddress)), "", "L", false)
	currentY = pdf.GetY() + 0.3

	pdf.SetFont("Helvetica", "B", 5.5*scale)
	pdf.SetXY(x+margin, currentY)
	city := strings.TrimSpace(c.WarehouseCity)
	state := strings.TrimSpace(c.WarehouseState)
	if state != "" {
		city += ", " + state
	}
	pdf.CellFormat(w-margin*2, 1.6*scale, tr(city), "", 1, "L", false, 0, "")
	currentY = pdf.GetY() + 0.3

	pdf.SetFont("Helvetica", "", 4.5*scale)
	pdf.SetTextColor(100, 100, 100)
	if c.WarehousePhone != "" {
		pdf.SetXY(x+margin, currentY)
		pdf.CellFormat(w-margin*2, 1.3*scale, tr("Tel: "+c.WarehousePhone), "", 1, "L", false, 0, "")
	}
}

func drawColumnDestino(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, x, y, w, h float64, bgCol, headerCol, textCol, borderCol color.RGBA, scale float64) {
	pdf.SetFillColor(int(bgCol.R), int(bgCol.G), int(bgCol.B))
	pdf.Rect(x, y, w, h, "F")

	margin := 2.5
	currentY := y + 1.5

	pdf.SetXY(x+margin, currentY)
	pdf.SetFont("Helvetica", "B", 7*scale)
	pdf.SetTextColor(int(headerCol.R), int(headerCol.G), int(headerCol.B))
	pdf.CellFormat(w-margin*2, 2*scale, "DESTINO", "", 1, "L", false, 0, "")
	currentY += 3.5

	pdf.SetXY(x+margin, currentY)
	pdf.SetFont("Helvetica", "B", 6.5*scale)
	pdf.SetTextColor(int(textCol.R), int(textCol.G), int(textCol.B))
	pdf.MultiCell(w-margin*2, 1.8*scale, tr(strings.ToUpper(strings.TrimSpace(c.CustomerName))), "", "L", false)
	currentY = pdf.GetY() + 0.5

	pdf.SetFont("Helvetica", "", 5*scale)
	pdf.SetXY(x+margin, currentY)
	pdf.MultiCell(w-margin*2, 1.4*scale, tr(strings.TrimSpace(c.DestinationAddress)), "", "L", false)
	currentY = pdf.GetY() + 0.3

	pdf.SetFont("Helvetica", "", 4.8*scale)
	pdf.SetXY(x+margin, currentY)
	suburb := strings.TrimSpace(c.DestinationSuburb)
	if suburb != "" {
		pdf.CellFormat(w-margin*2, 1.3*scale, tr(suburb), "", 1, "L", false, 0, "")
		currentY = pdf.GetY() + 0.2
	}

	pdf.SetFont("Helvetica", "B", 5.5*scale)
	pdf.SetTextColor(int(textCol.R), int(textCol.G), int(textCol.B))
	pdf.SetXY(x+margin, currentY)
	dest := strings.TrimSpace(c.DestinationCity)
	state := strings.TrimSpace(c.DestinationState)
	if state != "" {
		dest += ", " + state
	}
	pdf.CellFormat(w-margin*2, 1.6*scale, tr(dest), "", 1, "L", false, 0, "")
	currentY = pdf.GetY() + 0.3

	pdf.SetFont("Helvetica", "", 4.5*scale)
	pdf.SetTextColor(100, 100, 100)
	if c.CustomerPhone != "" {
		pdf.SetXY(x+margin, currentY)
		pdf.CellFormat(w-margin*2, 1.3*scale, tr("Tel: "+c.CustomerPhone), "", 1, "L", false, 0, "")
	}
}

func drawColumnPaquete(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, x, y, w, h float64, headerCol, textCol, borderCol color.RGBA, scale float64) {
	margin := 2.5
	currentY := y + 1.5

	pdf.SetXY(x+margin, currentY)
	pdf.SetFont("Helvetica", "B", 7*scale)
	pdf.SetTextColor(int(headerCol.R), int(headerCol.G), int(headerCol.B))
	pdf.CellFormat(w-margin*2, 2*scale, "PAQUETE", "", 1, "L", false, 0, "")
	currentY += 3.5

	pdf.SetFont("Helvetica", "", 4.8*scale)
	pdf.SetTextColor(int(textCol.R), int(textCol.G), int(textCol.B))
	pdf.SetXY(x+margin, currentY)

	if len(c.OrderItems) > 0 {
		for i, item := range c.OrderItems {
			if i >= 2 {
				break
			}
			line := fmt.Sprintf("%dx %s", item.Quantity, item.ProductName)
			if item.SKU != "" {
				line = fmt.Sprintf("%dx SKU %s", item.Quantity, item.SKU)
			}
			pdf.MultiCell(w-margin*2, 1.3*scale, tr(line), "", "L", false)
		}
	} else {
		pdf.CellFormat(w-margin*2, 1.3*scale, tr("Items no especificados"), "", 1, "L", false, 0, "")
	}
	currentY = pdf.GetY() + 2

	pdf.SetFont("Helvetica", "B", 4.8*scale)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(x+margin, currentY)
	pdf.CellFormat(w-margin*2, 1.5*scale, "Peso", "", 1, "L", false, 0, "")
	currentY = pdf.GetY()

	pdf.SetFont("Helvetica", "B", 6*scale)
	pdf.SetTextColor(int(textCol.R), int(textCol.G), int(textCol.B))
	weightStr := "—"
	if c.Weight > 0 {
		weightStr = fmt.Sprintf("%.1f kg", c.Weight)
	}
	pdf.SetXY(x+margin, currentY)
	pdf.CellFormat(w-margin*2, 1.8*scale, tr(weightStr), "", 1, "L", false, 0, "")
	currentY = pdf.GetY() + 1

	pdf.SetFont("Helvetica", "B", 4.8*scale)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetXY(x+margin, currentY)
	pdf.CellFormat(w-margin*2, 1.5*scale, "Dimensiones", "", 1, "L", false, 0, "")
	currentY = pdf.GetY()

	pdf.SetFont("Helvetica", "B", 6*scale)
	pdf.SetTextColor(int(textCol.R), int(textCol.G), int(textCol.B))
	dimsStr := "—"
	if c.Width > 0 && c.Height > 0 && c.Length > 0 {
		dimsStr = fmt.Sprintf("%.0fx%.0fx%.0f cm", c.Width, c.Height, c.Length)
	}
	pdf.SetXY(x+margin, currentY)
	pdf.CellFormat(w-margin*2, 1.8*scale, tr(dimsStr), "", 1, "L", false, 0, "")

	if c.CodTotal > 0 {
		pdf.SetFillColor(200, 16, 46)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Helvetica", "B", 5*scale)
		currency := c.Currency
		if currency == "" {
			currency = "COP"
		}
		totalCliente := c.CodTotal + c.CodCarrierFee
		codTxt := fmt.Sprintf("CONTRA ENTREGA: $%s", formatMoneyProb(totalCliente))
		yPos := y + h - 8
		pdf.Rect(x+margin, yPos, w-margin*2, 6.5*scale, "F")
		pdf.SetXY(x+margin+0.5, yPos+0.8)
		pdf.MultiCell(w-margin*2-1, 2*scale, tr(codTxt), "", "C", false)
	}
}

func drawTrackingSectionNewDesign(pdf *gofpdf.Fpdf, tr func(string) string, c *domain.GuidePDFContext, wMm, hMm, margin, scale float64) {
	headerH := hMm * 0.25
	contentH := hMm * 0.50
	trackingY := headerH + contentH
	trackingH := hMm * 0.25

	darkBlue := color.RGBA{R: 22, G: 54, B: 92, A: 255}
	darkGray := color.RGBA{R: 17, G: 17, B: 17, A: 255}
	lightGray := color.RGBA{R: 100, G: 100, B: 100, A: 255}

	usableW := wMm - (margin * 2)

	pdf.SetXY(margin, trackingY+2)
	pdf.SetFont("Helvetica", "", 5*scale)
	pdf.SetTextColor(int(lightGray.R), int(lightGray.G), int(lightGray.B))
	pdf.CellFormat(usableW*0.3, 1.8*scale, "TRACKING", "", 0, "L", false, 0, "")

	pdf.SetXY(margin+usableW*0.5, trackingY+2)
	pdf.SetFont("Helvetica", "", 5*scale)
	pdf.CellFormat(usableW*0.5, 1.8*scale, "ORDEN", "", 1, "R", false, 0, "")

	pdf.SetXY(margin, trackingY+4.5)
	pdf.SetFont("Helvetica", "B", 12*scale)
	pdf.SetTextColor(int(darkGray.R), int(darkGray.G), int(darkGray.B))
	pdf.CellFormat(usableW*0.5, 4*scale, tr(strings.TrimSpace(c.TrackingNumber)), "", 0, "L", false, 0, "")

	pdf.SetXY(margin+usableW*0.5, trackingY+4.5)
	pdf.SetFont("Helvetica", "B", 9*scale)
	pdf.CellFormat(usableW*0.5, 4*scale, tr(strings.TrimSpace(c.OrderNumber)), "", 1, "R", false, 0, "")

	pdf.SetDrawColor(int(lightGray.R), int(lightGray.G), int(lightGray.B))
	pdf.SetLineWidth(0.3)
	pdf.Line(margin, trackingY+8.5, margin+usableW, trackingY+8.5)

	barcodeH := trackingH - 9.5
	barcodeY := trackingY + 9.0

	barcodePNG := buildCode128PNGProb(strings.TrimSpace(c.TrackingNumber), int(usableW*7), int(barcodeH*2))
	if barcodePNG != nil && len(barcodePNG) > 0 {
		opts := gofpdf.ImageOptions{ImageType: "PNG"}
		barcodeKey := fmt.Sprintf("barcode_%d", time.Now().UnixNano())
		pdf.RegisterImageOptionsReader(barcodeKey, opts, bytes.NewReader(barcodePNG))

		barcodeX := margin + (usableW - (usableW * 0.8)) / 2
		pdf.ImageOptions(barcodeKey, barcodeX, barcodeY, usableW*0.8, barcodeH*0.8, false, opts, 0, "")
	}

	pdf.SetXY(margin, trackingY+trackingH-2)
	pdf.SetFont("Helvetica", "", 4*scale)
	pdf.SetTextColor(int(darkBlue.R), int(darkBlue.G), int(darkBlue.B))
	barcodeTxt := fmt.Sprintf("CODE 128 · *%s*", strings.TrimSpace(c.TrackingNumber))
	pdf.CellFormat(usableW*0.5, 1.5*scale, tr(barcodeTxt), "", 0, "L", false, 0, "")

	pdf.SetXY(margin+usableW*0.5, trackingY+trackingH-2)
	pdf.SetFont("Helvetica", "", 4*scale)
	pdf.SetTextColor(int(lightGray.R), int(lightGray.G), int(lightGray.B))
	guiaTxt := fmt.Sprintf("Guia: %s", strings.TrimSpace(c.Guia))
	pdf.CellFormat(usableW*0.5, 1.5*scale, tr(guiaTxt), "", 1, "R", false, 0, "")
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
	width := b.Dx()
	height := b.Dy()

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			origX := b.Min.X + x
			origY := b.Min.Y + y
			r, _, _, a := scaled.At(origX, origY).RGBA()

			if a == 0 {
				rgba.SetRGBA(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else if r>>8 > 127 {
				rgba.SetRGBA(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else {
				rgba.SetRGBA(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, rgba); err != nil {
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

func pageWidth(pdf *gofpdf.Fpdf) float64 {
	w, _ := pdf.GetPageSize()
	return w
}
