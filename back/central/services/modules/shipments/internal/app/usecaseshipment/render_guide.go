package usecaseshipment

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpulib "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	pdftypes "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type RenderedGuide struct {
	PDF      []byte
	Format   *domain.GuideFormat
	Filename string
}

func (uc *UseCaseShipment) RenderGuide(ctx context.Context, shipmentID uint, formatCode string) (*RenderedGuide, error) {
	shipment, err := uc.repo.GetShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("get shipment: %w", err)
	}
	if shipment == nil {
		return nil, fmt.Errorf("shipment %d no encontrado", shipmentID)
	}

	carrier := ""
	if shipment.Carrier != nil {
		carrier = strings.ToUpper(strings.TrimSpace(*shipment.Carrier))
	}

	var format *domain.GuideFormat
	if formatCode != "" {
		format, err = uc.repo.GetGuideFormatByCode(ctx, formatCode)
		if err != nil {
			return nil, fmt.Errorf("get format: %w", err)
		}
		if format == nil {
			return nil, fmt.Errorf("formato %q no encontrado", formatCode)
		}
		fmtCarrier := strings.ToUpper(strings.TrimSpace(format.Carrier))
		if fmtCarrier != domain.CarrierUniversal && fmtCarrier != carrier {
			return nil, fmt.Errorf("formato %q no aplica al carrier %s", formatCode, carrier)
		}
	} else {
		format, err = uc.repo.GetDefaultGuideFormat(ctx, carrier)
		if err != nil {
			return nil, fmt.Errorf("get default format: %w", err)
		}
	}

	if format != nil && format.Strategy == domain.GuideStrategyRebuild {
		pdfCtx, err := uc.repo.GetGuidePDFContext(ctx, shipmentID)
		if err != nil {
			return nil, fmt.Errorf("get pdf context: %w", err)
		}
		if pdfCtx == nil {
			return nil, fmt.Errorf("contexto del shipment %d no disponible", shipmentID)
		}

		if pdfCtx.Destino == "" && pdfCtx.ZonaHub == "" && pdfCtx.EquipoReparto == "" && shipment.GuideURL != nil && *shipment.GuideURL != "" {
			metadata, err := ExtractCoordinadoraMetadata(ctx, *shipment.GuideURL)
			if err == nil && len(metadata) > 0 {
				pdfCtx.Destino = getMetaStr(metadata, "destino")
				pdfCtx.ZonaHub = getMetaStr(metadata, "zona_hub")
				pdfCtx.EquipoReparto = getMetaStr(metadata, "equipo_reparto")
				pdfCtx.Origen = getMetaStr(metadata, "origen")
				pdfCtx.AsCode = getMetaStr(metadata, "as_code")
				pdfCtx.Paq = getMetaStr(metadata, "paq")
				pdfCtx.Unidad = getMetaStr(metadata, "unidad")
				pdfCtx.WarehousePostal = getMetaStr(metadata, "postal_origen")
				pdfCtx.Ref = getMetaStr(metadata, "ref")
				pdfCtx.Guia = getMetaStr(metadata, "guia")
				pdfCtx.Observaciones = getMetaStr(metadata, "observaciones")
			}
		}

		pdfBytes, err := buildProbabilityLabel(pdfCtx, format)
		if err != nil {
			return nil, fmt.Errorf("build probability label: %w", err)
		}
		filename := fmt.Sprintf("probability-%d-%s.pdf", shipment.ID, format.Code)
		return &RenderedGuide{PDF: pdfBytes, Format: format, Filename: filename}, nil
	}

	if shipment.GuideURL == nil || *shipment.GuideURL == "" {
		return nil, fmt.Errorf("shipment %d sin guide_url del carrier", shipmentID)
	}

	originalPDF, err := downloadPDF(ctx, *shipment.GuideURL)
	if err != nil {
		return nil, fmt.Errorf("download original pdf: %w", err)
	}

	pdfBytes, err := applyGuideFormat(originalPDF, format)
	if err != nil {
		return nil, fmt.Errorf("apply format: %w", err)
	}

	formatCodeForName := ""
	if format != nil {
		formatCodeForName = "-" + format.Code
	}
	filename := fmt.Sprintf("guia-%d%s.pdf", shipment.ID, formatCodeForName)

	return &RenderedGuide{PDF: pdfBytes, Format: format, Filename: filename}, nil
}

func downloadPDF(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func applyGuideFormat(pdfBytes []byte, format *domain.GuideFormat) ([]byte, error) {
	if format == nil || format.Strategy == domain.GuideStrategyPassthrough {
		return pdfBytes, nil
	}

	page := format.SourcePage
	if page <= 0 {
		page = 1
	}
	selPages := []string{fmt.Sprintf("%d", page)}

	cropped, err := cropPDF(pdfBytes, format, selPages)
	if err != nil {
		return nil, err
	}

	if format.Strategy == domain.GuideStrategyCrop {
		return cropped, nil
	}

	if format.Strategy == domain.GuideStrategyResize {
		return resizePDF(cropped, format, selPages)
	}

	return cropped, nil
}

func cropPDF(pdfBytes []byte, format *domain.GuideFormat, selPages []string) ([]byte, error) {
	conf := model.NewDefaultConfiguration()

	if format.CropLLxFrac == 0 && format.CropLLyFrac == 0 && format.CropURxFrac == 1 && format.CropURyFrac == 1 {
		return singlePageOnly(pdfBytes, selPages, conf)
	}

	dims, err := api.PageDims(bytes.NewReader(pdfBytes), conf)
	if err != nil {
		return nil, fmt.Errorf("read dims: %w", err)
	}
	if len(dims) == 0 {
		return pdfBytes, nil
	}
	pageIdx := 0
	if format.SourcePage > 0 && format.SourcePage-1 < len(dims) {
		pageIdx = format.SourcePage - 1
	}
	pageW := dims[pageIdx].Width
	pageH := dims[pageIdx].Height

	cropBox := &model.Box{
		Rect: pdftypes.NewRectangle(
			format.CropLLxFrac*pageW,
			format.CropLLyFrac*pageH,
			format.CropURxFrac*pageW,
			format.CropURyFrac*pageH,
		),
	}
	pb := &model.PageBoundaries{Crop: cropBox, Media: cropBox}

	var withBoxes bytes.Buffer
	if err := api.AddBoxes(bytes.NewReader(pdfBytes), &withBoxes, selPages, pb, conf); err != nil {
		return nil, fmt.Errorf("add boxes: %w", err)
	}

	var trimmed bytes.Buffer
	if err := api.Trim(bytes.NewReader(withBoxes.Bytes()), &trimmed, selPages, conf); err != nil {
		return withBoxes.Bytes(), nil
	}
	return trimmed.Bytes(), nil
}

func singlePageOnly(pdfBytes []byte, selPages []string, conf *model.Configuration) ([]byte, error) {
	dims, err := api.PageDims(bytes.NewReader(pdfBytes), conf)
	if err != nil || len(dims) <= 1 {
		return pdfBytes, nil
	}
	var trimmed bytes.Buffer
	if err := api.Trim(bytes.NewReader(pdfBytes), &trimmed, selPages, conf); err != nil {
		return pdfBytes, nil
	}
	return trimmed.Bytes(), nil
}

func resizePDF(pdfBytes []byte, format *domain.GuideFormat, selPages []string) ([]byte, error) {
	conf := model.NewDefaultConfiguration()

	targetW := format.WidthCm
	targetH := format.HeightCm
	dims, _ := api.PageDims(bytes.NewReader(pdfBytes), conf)
	if len(dims) > 0 {
		srcLandscape := dims[0].Width > dims[0].Height
		targetLandscape := targetW > targetH
		if srcLandscape != targetLandscape {
			targetW, targetH = targetH, targetW
		}
	}

	wPt := cmToPoints(targetW)
	hPt := cmToPoints(targetH)
	resDesc := fmt.Sprintf("dim:%.2f %.2f", wPt, hPt)
	res, err := pdfcpulib.ParseResizeConfig(resDesc, pdftypes.POINTS)
	if err != nil {
		return nil, fmt.Errorf("parse resize: %w", err)
	}

	var out bytes.Buffer
	if err := api.Resize(bytes.NewReader(pdfBytes), &out, selPages, res, conf); err != nil {
		return nil, fmt.Errorf("resize: %w", err)
	}
	return out.Bytes(), nil
}

func cmToPoints(cm float64) float64 {
	return cm * 28.3464567
}

func getMetaStr(metadata map[string]interface{}, key string) string {
	if metadata == nil {
		return ""
	}
	if v, ok := metadata[key]; ok {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}
