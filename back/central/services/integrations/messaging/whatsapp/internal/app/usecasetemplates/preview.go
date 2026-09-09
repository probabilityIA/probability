package usecasetemplates

import (
	"context"
	"fmt"
	"strings"
)

type TemplatePreview struct {
	Name       string   `json:"name"`
	Language   string   `json:"language"`
	Status     string   `json:"status"`
	Category   string   `json:"category"`
	MetaID     string   `json:"meta_id"`
	Header     string   `json:"header"`
	Body       string   `json:"body"`
	Footer     string   `json:"footer"`
	Buttons    []string `json:"buttons"`
	Variables  int      `json:"variables"`
	Condition  string   `json:"condition"`
	FoundInWA  bool     `json:"found_in_whatsapp"`
	RejectedBy string   `json:"rejected_reason"`
}

type templateCandidate struct {
	name      string
	condition string
}

func candidatesForEvent(eventCode string) []templateCandidate {
	switch strings.TrimSpace(eventCode) {
	case "order.created":
		return []templateCandidate{
			{"confirmacion_pedido", "Pedido pagado por adelantado"},
			{"confirmacion_pedido_contraentrega", "Contra entrega con valor a recaudar ya definido"},
			{"confirmacion_pedido_contraentrega_sin_valor", "Contra entrega cuando el flete todavia no se conoce"},
		}
	case "order.created_with_map":
		return []templateCandidate{
			{"confirmacion_pedido_contraentrega_mapa", "Contra entrega, con mapa de la direccion en el encabezado"},
		}
	case "order.shipped":
		return []templateCandidate{
			{"pedido_en_reparto", "Pedido pagado por adelantado"},
			{"pedido_en_reparto_cod", "Contra entrega"},
		}
	case "order.delivered":
		return []templateCandidate{
			{"pedido_entregado", "Pedido pagado por adelantado"},
			{"pedido_entregado_cod", "Contra entrega"},
		}
	case "shipment.guide_generated":
		return []templateCandidate{
			{"guia_envio_generada", "Pedido pagado por adelantado"},
			{"guia_envio_generada_cod", "Contra entrega"},
		}
	case "wallet.low_balance":
		return []templateCandidate{
			{"reporte_saldo_billetera_v2", "Saldo de la billetera por debajo del minimo"},
		}
	case "order.canceled":
		return []templateCandidate{
			{"pedido_cancelado", "Pedido cancelado"},
		}
	default:
		return nil
	}
}

func (u *usecase) PreviewByEvent(ctx context.Context, businessID uint, eventCode string) ([]TemplatePreview, error) {
	candidates := candidatesForEvent(eventCode)
	if len(candidates) == 0 {
		return []TemplatePreview{}, nil
	}

	wabaID, token, baseURL, err := u.resolveSubmissionTarget(ctx, businessID)
	if err != nil {
		return nil, err
	}

	remote, err := u.apiFactory(baseURL).ListTemplates(ctx, wabaID, token)
	if err != nil {
		return nil, fmt.Errorf("error consultando las plantillas en Meta: %w", err)
	}

	byName := make(map[string]int, len(remote))
	for i, tpl := range remote {
		byName[strings.ToLower(strings.TrimSpace(tpl.Name))] = i
	}

	previews := make([]TemplatePreview, 0, len(candidates))

	for _, candidate := range candidates {
		preview := TemplatePreview{
			Name:      candidate.name,
			Condition: candidate.condition,
		}

		index, found := byName[candidate.name]
		if !found {
			previews = append(previews, preview)
			continue
		}

		tpl := remote[index]
		preview.FoundInWA = true
		preview.Language = tpl.Language
		preview.Status = tpl.Status
		preview.Category = tpl.Category
		preview.MetaID = tpl.ID
		preview.RejectedBy = tpl.RejectedReason

		for _, component := range tpl.Components {
			switch strings.ToUpper(stringField(component, "type")) {
			case "HEADER":
				preview.Header = stringField(component, "text")
			case "BODY":
				preview.Body = stringField(component, "text")
			case "FOOTER":
				preview.Footer = stringField(component, "text")
			case "BUTTONS":
				preview.Buttons = buttonLabels(component)
			}
		}

		preview.Variables = countPlaceholders(preview.Body)

		previews = append(previews, preview)
	}

	return previews, nil
}

func buttonLabels(component map[string]any) []string {
	raw, ok := component["buttons"].([]any)
	if !ok {
		return nil
	}

	labels := make([]string, 0, len(raw))
	for _, item := range raw {
		button, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if text := stringField(button, "text"); text != "" {
			labels = append(labels, text)
		}
	}

	return labels
}

func countPlaceholders(body string) int {
	seen := make(map[string]bool)
	remaining := body

	for {
		start := strings.Index(remaining, "{{")
		if start < 0 {
			break
		}
		end := strings.Index(remaining[start:], "}}")
		if end < 0 {
			break
		}
		token := remaining[start+2 : start+end]
		seen[strings.TrimSpace(token)] = true
		remaining = remaining[start+end+2:]
	}

	return len(seen)
}
