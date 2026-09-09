package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/push/internal/domain/errors"
)

func (u *UseCase) NotifyBusiness(ctx context.Context, event dtos.PushEventDTO) error {
	if event.BusinessID == 0 {
		return domainerrors.ErrBusinessRequired
	}

	if !u.sender.IsConfigured() {
		u.log.Warn(ctx).
			Uint("business_id", event.BusinessID).
			Str("event_type", event.EventType).
			Msg("envio de push sin configurar, evento descartado")
		return domainerrors.ErrFCMNotConfigured
	}

	message, ok := buildMessage(event)
	if !ok {
		u.log.Debug(ctx).
			Str("event_type", event.EventType).
			Str("status_code", event.StatusCode).
			Msg("evento sin mensaje push definido, se ignora")
		return nil
	}

	devices, err := u.repo.ListActiveTokensByBusiness(ctx, event.BusinessID)
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		u.log.Debug(ctx).
			Uint("business_id", event.BusinessID).
			Msg("el negocio no tiene dispositivos registrados")
		return nil
	}

	tokens := make([]string, 0, len(devices))
	for _, d := range devices {
		tokens = append(tokens, d.Token)
	}

	result, err := u.sender.Send(ctx, tokens, message)

	if len(result.InvalidTokens) > 0 {
		if deactivateErr := u.repo.DeactivateTokens(ctx, result.InvalidTokens); deactivateErr != nil {
			u.log.Warn(ctx).Err(deactivateErr).Msg("no se pudieron desactivar los tokens invalidos")
		}
	}

	if err != nil {
		return err
	}

	u.log.Info(ctx).
		Uint("business_id", event.BusinessID).
		Str("event_type", event.EventType).
		Int("enviados", result.Sent).
		Int("invalidos", len(result.InvalidTokens)).
		Msg("push enviado")
	return nil
}

func buildMessage(event dtos.PushEventDTO) (entities.PushMessage, bool) {
	data := map[string]string{
		"event_type": event.EventType,
	}
	if event.OrderID != "" {
		data["order_id"] = event.OrderID
	}
	if event.OrderNumber != "" {
		data["order_number"] = event.OrderNumber
	}
	if event.TrackingNo != "" {
		data["tracking_number"] = event.TrackingNo
	}

	orden := event.OrderNumber
	if orden == "" {
		orden = "una orden"
	} else {
		orden = "la orden " + orden
	}

	switch event.EventType {
	case "shipment.guide_generated", "order.guide_notification_requested":
		data["route"] = "/orders"
		carrier := event.Carrier
		if carrier == "" {
			carrier = "la transportadora"
		}
		body := fmt.Sprintf("%s ya tiene gu\u00eda con %s", capitalizar(orden), carrier)
		if event.TrackingNo != "" {
			body = fmt.Sprintf("%s. Gu\u00eda %s", body, event.TrackingNo)
		}
		return entities.PushMessage{Title: "Gu\u00eda generada", Body: body, Data: data}, true

	case "wallet.low_balance":
		data["route"] = "/wallet"
		return entities.PushMessage{
			Title: "Saldo bajo",
			Body:  "El saldo no alcanza para seguir generando gu\u00edas. Rec\u00e1rgalo para no frenar los env\u00edos",
			Data:  data,
		}, true

	case "order.delivered":
		data["route"] = "/orders"
		return entities.PushMessage{
			Title: "Pedido entregado",
			Body:  fmt.Sprintf("%s fue entregada", capitalizar(orden)),
			Data:  data,
		}, true

	case "order.status_changed":
		return buildStatusMessage(event, data, orden)
	}

	return entities.PushMessage{}, false
}

func buildStatusMessage(event dtos.PushEventDTO, data map[string]string, orden string) (entities.PushMessage, bool) {
	data["route"] = "/orders"
	data["status_code"] = event.StatusCode

	switch event.StatusCode {
	case "delivery_novelty":
		return entities.PushMessage{
			Title: "Novedad en la entrega",
			Body:  fmt.Sprintf("%s tiene una novedad. Atenderla hoy evita la devoluci\u00f3n", capitalizar(orden)),
			Data:  data,
		}, true

	case "delivery_failed":
		return entities.PushMessage{
			Title: "Entrega fallida",
			Body:  fmt.Sprintf("No se pudo entregar %s", orden),
			Data:  data,
		}, true

	case "delivered":
		return entities.PushMessage{
			Title: "Pedido entregado",
			Body:  fmt.Sprintf("%s fue entregada", capitalizar(orden)),
			Data:  data,
		}, true

	case "returned", "return_in_transit":
		return entities.PushMessage{
			Title: "Pedido devuelto",
			Body:  fmt.Sprintf("%s viene de regreso", capitalizar(orden)),
			Data:  data,
		}, true
	}

	return entities.PushMessage{}, false
}

func capitalizar(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
