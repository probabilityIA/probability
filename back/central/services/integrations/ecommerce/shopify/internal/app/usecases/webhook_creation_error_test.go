package usecases

import (
	"errors"
	"strings"
	"testing"
)

func TestWebhookCreationError_ProtectedCustomerData(t *testing.T) {
	shopifyErr := errors.New("acceso denegado al crear el webhook: {\"errors\":\"You do not have permission to create or update webhooks with orders/create topic. This topic contains protected customer data.\"}")

	err := webhookCreationError(shopifyErr)
	if err == nil {
		t.Fatal("se esperaba un error, se obtuvo nil")
	}

	msg := err.Error()
	if !strings.Contains(msg, "Protected customer data") {
		t.Errorf("el mensaje debe mencionar 'Protected customer data', se obtuvo: %s", msg)
	}
	if !strings.Contains(msg, "Partner Dashboard") {
		t.Errorf("el mensaje debe guiar al Partner Dashboard, se obtuvo: %s", msg)
	}
}

func TestWebhookCreationError_ProtectedDataTakesPriorityOverScope(t *testing.T) {
	shopifyErr := errors.New("permission to create or update webhooks; requires access scope")

	err := webhookCreationError(shopifyErr)
	if err == nil {
		t.Fatal("se esperaba un error, se obtuvo nil")
	}

	if !strings.Contains(err.Error(), "Partner Dashboard") {
		t.Errorf("el caso de datos protegidos debe tener prioridad, se obtuvo: %s", err.Error())
	}
}

func TestWebhookCreationError_TokenExpired(t *testing.T) {
	err := webhookCreationError(errors.New("token de acceso inválido o expirado"))
	if err == nil {
		t.Fatal("se esperaba un error, se obtuvo nil")
	}
	if !strings.Contains(err.Error(), "Reconecta la integracion") {
		t.Errorf("se esperaba guia de reconexion, se obtuvo: %s", err.Error())
	}
}
