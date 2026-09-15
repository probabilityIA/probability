package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

const (
	APIKey = "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"
	APIURL = "https://openrouter.ai/api/v1/chat/completions"
	Model  = "xiaomi/mimo-v2-flash:free"
)

type recommendationPayload struct {
	RecommendedCarrier string             `json:"recommended_carrier"`
	Reasoning          string             `json:"reasoning"`
	Alternatives       []string           `json:"alternatives"`
	Quotations         []quotationPayload `json:"quotations"`
}

type quotationPayload struct {
	Carrier               string  `json:"carrier"`
	EstimatedCost         float64 `json:"estimated_cost"`
	EstimatedDeliveryDays int     `json:"estimated_delivery_days"`
}

func (c *Client) loadTransportData() {
	cwd, _ := os.Getwd()

	positions := []string{
		filepath.Join(cwd, "services", "modules", "ai", "resources", "transportadoras.json"),
		filepath.Join(cwd, "..", "services", "modules", "ai", "resources", "transportadoras.json"),
		filepath.Join(cwd, "back", "central", "services", "modules", "ai", "resources", "transportadoras.json"),
		filepath.Join(cwd, "..", "..", "services", "modules", "ai", "resources", "transportadoras.json"),
	}

	var data []byte
	var err error
	for _, p := range positions {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to load transportadoras.json from any known path")
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		c.logger.Error().Err(err).Msg("Failed to parse transportadoras.json")
		return
	}
	c.transportData = result
}

func (c *Client) GetRecommendation(ctx context.Context, origin, destination string) (*entities.Recommendation, error) {
	if c.transportData == nil {
		c.loadTransportData()
	}

	transportDataJSON, _ := json.Marshal(c.transportData)
	prompt := fmt.Sprintf(recommendationPrompt, string(transportDataJSON), origin, destination)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	})

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, APIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("HTTP-Referer", "http://localhost:8000")
	httpReq.Header.Set("X-Title", "ProbabilityBackend")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var openRouterResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		c.logger.Error().Err(err).Str("body", string(body)).Msg("Failed to parse OpenRouter response")
		return nil, err
	}
	if len(openRouterResp.Choices) == 0 {
		c.logger.Error().Str("body", string(body)).Msg("No choices returned from AI")
		return nil, fmt.Errorf("no choices returned from AI: %s", string(body))
	}

	content := openRouterResp.Choices[0].Message.Content
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var payload recommendationPayload
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		return nil, err
	}
	return toRecommendation(payload), nil
}

func toRecommendation(p recommendationPayload) *entities.Recommendation {
	quotations := make([]entities.Quotation, 0, len(p.Quotations))
	for _, q := range p.Quotations {
		quotations = append(quotations, entities.Quotation{
			Carrier:               q.Carrier,
			EstimatedCost:         q.EstimatedCost,
			EstimatedDeliveryDays: q.EstimatedDeliveryDays,
		})
	}
	return &entities.Recommendation{
		RecommendedCarrier: p.RecommendedCarrier,
		Reasoning:          p.Reasoning,
		Alternatives:       p.Alternatives,
		Quotations:         quotations,
	}
}

const recommendationPrompt = "\n" +
	"    Act\u00faa como un experto en log\u00edstica y an\u00e1lisis de datos.\n" +
	"    Tienes la siguiente informaci\u00f3n hist\u00f3rica de transportadoras en formato JSON:\n" +
	"    %s\n" +
	"\n" +
	"    Basado en esta informaci\u00f3n, determina cu\u00e1l es la mejor transportadora para realizar un env\u00edo desde '%s' hasta '%s'.\n" +
	"\n" +
	"    Instrucciones de an\u00e1lisis:\n" +
	"    1.  **B\u00fasqueda Directa**: Busca si existe la ruta espec\u00edfica (origen -> destino) en la lista de 'coverage' de las transportadoras.\n" +
	"    2.  **Inferencia**: Si la ruta exacta no est\u00e1, usa tu conocimiento geogr\u00e1fico y el 'status_summary' para elegir la mejor opci\u00f3n.\n" +
	"    3.  **Rendimiento General**: Mira el 'status_summary' para justificar tu elecci\u00f3n.\n" +
	"\n" +
	"    IMPORTANTE:\n" +
	"    *   NO des cifras exactas de historial, pero S\u00cd estima costos y tiempos para las cotizaciones.\n" +
	"    *   Usa t\u00e9rminos cualitativos como \"alto volumen de entregas\", \"baja tasa de cancelaci\u00f3n\", \"amplia cobertura\".\n" +
	"    *   Genera 3 cotizaciones estimadas para las mejores opciones (incluyendo la recomendada). Costo en COP y d\u00edas h\u00e1biles.\n" +
	"\n" +
	"    Tu respuesta debe ser un JSON con el siguiente formato:\n" +
	"    {\n" +
	"        \"recommended_carrier\": \"Nombre exacto de la transportadora recomendada\",\n" +
	"        \"reasoning\": \"Explicaci\u00f3n detallada sin cifras exactas. Enf\u00f3cate en confiabilidad y cobertura.\",\n" +
	"        \"alternatives\": [\"Otras opciones viables\"],\n" +
	"        \"quotations\": [\n" +
	"            {\n" +
	"                \"carrier\": \"Nombre\",\n" +
	"                \"estimated_cost\": 15000,\n" +
	"                \"estimated_delivery_days\": 3\n" +
	"            }\n" +
	"        ]\n" +
	"    }\n" +
	"    Solo devuelve el JSON, nada m\u00e1s.\n" +
	"    "
