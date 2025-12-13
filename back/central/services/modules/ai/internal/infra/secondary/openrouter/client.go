package openrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	APIKey = "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"
	APIURL = "https://openrouter.ai/api/v1/chat/completions"
	Model  = "amazon/nova-2-lite-v1:free"
)

type Client struct {
	logger        log.ILogger
	transportData map[string]interface{}
}

func New(logger log.ILogger) *Client {
	client := &Client{
		logger: logger,
	}
	client.loadTransportData()
	return client
}

func (c *Client) loadTransportData() {
	// Try multiple paths to find the file
	// Get current working directory
	cwd, _ := os.Getwd()

	positions := []string{
		filepath.Join(cwd, "services", "modules", "ai", "resources", "transportadoras.json"),
		filepath.Join(cwd, "..", "services", "modules", "ai", "resources", "transportadoras.json"),
		filepath.Join(cwd, "back", "central", "services", "modules", "ai", "resources", "transportadoras.json"),
		// Hardcoded fallback for typical structure if running from cmd
		filepath.Join(cwd, "..", "..", "services", "modules", "ai", "resources", "transportadoras.json"),
	}

	var data []byte
	var err error
	var foundPath string

	for _, p := range positions {
		data, err = os.ReadFile(p)
		if err == nil {
			foundPath = p
			break
		}
	}

	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to load transportadoras.json from any known path")
		return
	}
	c.logger.Info().Str("path", foundPath).Msg("Loaded transportadoras.json")

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		c.logger.Error().Err(err).Msg("Failed to parse transportadoras.json")
		return
	}
	c.transportData = result
}

func (c *Client) GetRecommendation(req domain.RecommendationRequest) (*domain.AIRecommendation, error) {
	if c.transportData == nil {
		c.loadTransportData() // Retry loading
	}

	transportDataJSON, _ := json.Marshal(c.transportData)
	availableText := ""
	if len(req.Carriers) > 0 {
		availableText = fmt.Sprintf("IMPORTANTE: Solo puedes recomendar una de las siguientes transportadoras disponibles: %s. Si ninguna es ideal, elige la mejor de esta lista.", strings.Join(req.Carriers, ", "))
	}

	prompt := fmt.Sprintf(`
    Actúa como un experto en logística y análisis de datos.
    Tienes la siguiente información histórica de transportadoras en formato JSON:
    %s

    %s

    Basado en esta información, determina cuál es la mejor transportadora para realizar un envío desde '%s' hasta '%s'.
    
    Instrucciones de análisis:
    1.  **Filtro de Disponibilidad**: Si se proporcionó una lista de transportadoras disponibles, TU RECOMENDACIÓN DEBE SER UNA DE ELLAS.
    2.  **Búsqueda Directa**: Busca si existe la ruta específica (origen -> destino) en la lista de 'coverage' de las transportadoras disponibles.
    3.  **Inferencia**: Si la ruta exacta no está, usa tu conocimiento geográfico y el 'status_summary' para elegir la mejor opción DE LAS DISPONIBLES.
    4.  **Rendimiento General**: Mira el 'status_summary' para justificar tu elección.
    
    IMPORTANTE:
    *   NO des cifras exactas de historial, pero SÍ estima costos y tiempos para las cotizaciones.
    *   Usa términos cualitativos como "alto volumen de entregas", "baja tasa de cancelación", "amplia cobertura".
    *   Genera 3 cotizaciones estimadas para las mejores opciones (incluyendo la recomendada). Costo en COP y días hábiles.

    Tu respuesta debe ser un JSON con el siguiente formato:
    {
        "recommended_carrier": "Nombre exacto de la transportadora recomendada",
        "reasoning": "Explicación detallada sin cifras exactas. Enfócate en confiabilidad y cobertura.",
        "alternatives": ["Otras opciones viables"],
        "quotations": [
            {
                "carrier": "Nombre",
                "estimated_cost": 15000,
                "estimated_delivery_days": 3
            }
        ]
    }
    Solo devuelve el JSON, nada más.
    `, string(transportDataJSON), availableText, req.Origin, req.Destination)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	})

	httpReq, err := http.NewRequest("POST", APIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("HTTP-Referer", "http://localhost:8000")
	httpReq.Header.Set("X-Title", "ProbabilityBackend")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
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
		return nil, err
	}

	if len(openRouterResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from AI")
	}

	content := openRouterResp.Choices[0].Message.Content
	// Clean markdown
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var recommendation domain.AIRecommendation
	if err := json.Unmarshal([]byte(content), &recommendation); err != nil {
		return nil, err
	}

	return &recommendation, nil
}
