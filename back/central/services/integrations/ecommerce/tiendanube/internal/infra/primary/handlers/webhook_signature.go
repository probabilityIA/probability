package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const hmacHeaderName = "x-linkedstore-hmac-sha256"

func firmaValida(body []byte, secret, firma string) bool {
	firma = strings.TrimSpace(firma)
	if firma == "" || strings.TrimSpace(secret) == "" {
		return false
	}
	recibida, err := hex.DecodeString(firma)
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hmac.Equal(mac.Sum(nil), recibida)
}

func storeIDDeIntegracion(storeID string, config map[string]interface{}) string {
	if s := strings.TrimSpace(storeID); s != "" {
		return s
	}
	if config != nil {
		if v, ok := config["store_id"]; ok && v != nil {
			return strings.TrimSpace(fmt.Sprintf("%v", v))
		}
	}
	return ""
}

func (h *tiendanubeHandler) firmaPrivacidadValida(ctx context.Context, body []byte, firma string) bool {
	if strings.TrimSpace(firma) == "" {
		return false
	}
	if firmaValida(body, h.platformConfig(ctx, "client_secret", false), firma) {
		return true
	}
	return firmaValida(body, h.platformConfig(ctx, "client_secret", true), firma)
}

func (h *tiendanubeHandler) verificarWebhook(ctx context.Context, body []byte, firma, integrationID, storeIDPayload string) error {
	if h.coreIntegration == nil {
		return fmt.Errorf("core de integraciones no disponible")
	}

	integracion, err := h.coreIntegration.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("consultando la integracion: %w", err)
	}
	if integracion == nil {
		return fmt.Errorf("la integracion no existe")
	}

	secret := h.platformConfig(ctx, "client_secret", integracion.IsTesting)
	if strings.TrimSpace(secret) == "" {
		return fmt.Errorf("client_secret ausente en las credenciales de plataforma de Tiendanube")
	}

	if !firmaValida(body, secret, firma) {
		return fmt.Errorf("la firma no corresponde al cuerpo recibido")
	}

	esperado := storeIDDeIntegracion(integracion.StoreID, integracion.Config)
	storeIDPayload = strings.TrimSpace(storeIDPayload)
	if esperado != "" && storeIDPayload != "" && esperado != storeIDPayload {
		return fmt.Errorf("el store_id del webhook (%s) no corresponde al de la integracion (%s)", storeIDPayload, esperado)
	}

	return nil
}
