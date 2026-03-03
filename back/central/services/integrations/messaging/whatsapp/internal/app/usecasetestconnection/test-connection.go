package usecasetestconnection

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// ITestConnectionUseCase define la interfaz para el caso de uso de prueba de conexión
type ITestConnectionUseCase interface {
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}, clientFactory func(string, log.ILogger) ports.IWhatsApp) error
}

type testConnectionUseCase struct {
	config env.IConfig
	logger log.ILogger
}

// New crea una nueva instancia del caso de uso de prueba de conexión
func New(config env.IConfig, logger log.ILogger) ITestConnectionUseCase {
	return &testConnectionUseCase{
		config: config,
		logger: logger,
	}
}

// TestConnection prueba la conexión enviando un mensaje de prueba con credenciales dinámicas
// Si test_phone_number está presente en config, envía mensaje hello_world.
// Si no está presente, solo valida credenciales básicas (para creación sin test_phone_number).
func (u *testConnectionUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}, clientFactory func(string, log.ILogger) ports.IWhatsApp) error {
	// 1. Extraer y validar parámetros básicos
	accessToken, ok := credentials["access_token"].(string)
	if !ok || accessToken == "" {
		u.logger.Error().Msg("access_token no encontrado o vacío")
		return fmt.Errorf("access_token es requerido")
	}

	phoneNumberIDStr, ok := config["phone_number_id"].(string)
	if !ok || phoneNumberIDStr == "" {
		if num, ok := config["phone_number_id"].(float64); ok {
			phoneNumberIDStr = fmt.Sprintf("%.0f", num)
		} else {
			u.logger.Error().Msg("phone_number_id no encontrado o vacío")
			return fmt.Errorf("phone_number_id es requerido")
		}
	}

	// 2. Verificar si hay test_phone_number para enviar mensaje
	testPhone, ok := config["test_phone_number"].(string)
	if !ok || testPhone == "" {
		// Si no hay test_phone_number, solo validar credenciales básicas
		u.logger.Info().
			Str("phone_number_id", phoneNumberIDStr).
			Msg("test_phone_number no encontrado, solo validando credenciales básicas (access_token y phone_number_id)")
		// Validación básica: access_token y phone_number_id ya validados arriba
		return nil
	}

	u.logger.Info().
		Str("phone_number_id", phoneNumberIDStr).
		Str("test_phone", testPhone).
		Msg("Usando test_phone_number desde la configuración de la integración")

	u.logger.Info().
		Str("phone_number_id", phoneNumberIDStr).
		Str("test_phone", testPhone).
		Msg("Parámetros extraídos correctamente")

	// 2. Obtener whatsapp_url del config map (enviado desde frontend), con fallback a .env
	whatsappURL, _ := config["whatsapp_url"].(string)
	if whatsappURL == "" {
		whatsappURL = u.config.Get("WHATSAPP_URL")
	}

	if whatsappURL == "" {
		u.logger.Error().Msg("whatsapp_url no encontrado en config ni en variables de entorno")
		return fmt.Errorf("whatsapp_url es requerido")
	}

	// 3. Crear cliente usando la factory con la URL dinámica
	waClient := clientFactory(whatsappURL, u.logger)

	// 4. Convertir ID
	pID, err := strconv.ParseUint(phoneNumberIDStr, 10, 64)
	if err != nil {
		u.logger.Error().Err(err).Str("phone_number_id", phoneNumberIDStr).Msg("Error al convertir phone_number_id")
		return fmt.Errorf("phone_number_id inválido: %w", err)
	}

	// 5. Construir mensaje
	msg := entities.TemplateMessage{
		MessagingProduct: "whatsapp",
		To:               testPhone,
		Type:             "template",
		Template: entities.TemplateData{
			Name: "hello_world",
			Language: entities.TemplateLanguage{
				Code: "en_US",
			},
		},
	}

	u.logger.Info().Msg("Enviando mensaje de prueba...")

	// 6. Enviar mensaje
	resp, err := waClient.SendMessage(ctx, uint(pID), msg, accessToken)
	if err != nil {
		u.logger.Error().Err(err).Msg("Error al enviar mensaje de prueba")
		return fmt.Errorf("error al enviar mensaje de prueba: %w", err)
	}

	u.logger.Info().Str("response", resp).Msg("Mensaje de prueba enviado exitosamente")

	return nil
}

