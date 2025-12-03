package tester

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/shared/log"
)

// WhatsAppTester implementa ITestIntegration para testear conexiones de WhatsApp
type WhatsAppTester struct {
	log log.ILogger
}

// NewWhatsAppTester crea una nueva instancia del tester de WhatsApp
func NewWhatsAppTester(logger log.ILogger) core.ITestIntegration {
	return &WhatsAppTester{
		log: logger,
	}
}

// TestConnection prueba la conexión con WhatsApp usando las credenciales y configuración proporcionadas
func (t *WhatsAppTester) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	// Validar credenciales requeridas
	accessToken, ok := credentials["access_token"].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("access_token requerido en las credenciales")
	}

	// Validar configuración requerida
	phoneNumberIDStr, ok := config["phone_number_id"].(string)
	if !ok {
		// Intentar como número
		phoneNumberIDNum, ok := config["phone_number_id"].(float64)
		if !ok {
			return fmt.Errorf("phone_number_id requerido en la configuración")
		}
		phoneNumberIDStr = fmt.Sprintf("%.0f", phoneNumberIDNum)
	}

	// Validar formato de phone_number_id
	if phoneNumberIDStr == "" {
		return fmt.Errorf("phone_number_id no puede estar vacío")
	}

	// Validación básica: verificar que el access_token tenga formato válido
	// Los tokens de WhatsApp suelen tener un formato específico
	if len(accessToken) < 10 {
		return fmt.Errorf("access_token parece inválido (muy corto)")
	}

	t.log.Info(ctx).
		Str("phone_number_id", phoneNumberIDStr).
		Msg("Test de conexión WhatsApp: credenciales y configuración válidas")

	// TODO: En el futuro, hacer un request real a la API de WhatsApp
	// Por ejemplo: GET /v18.0/{phone-number-id} para verificar que el token funciona
	// Por ahora, solo validamos el formato

	return nil
}
