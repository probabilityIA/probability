package usecases

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

type accionCanal struct {
	Fulfillment   string
	TrackingEvent string
	Descripcion   string
	Cancelar      bool
}

var accionesPorEstado = map[string]accionCanal{
	"picking":            {Fulfillment: domain.FulfillmentPacked, Descripcion: "El pedido esta siendo preparado"},
	"packing":            {Fulfillment: domain.FulfillmentPacked, Descripcion: "El pedido fue empacado"},
	"ready_to_ship":      {Fulfillment: domain.FulfillmentPacked, Descripcion: "El pedido esta listo para despacho"},
	"assigned_to_driver": {Fulfillment: domain.FulfillmentDispatched, TrackingEvent: domain.TrackingDispatched, Descripcion: "El pedido fue asignado al transportador"},
	"picked_up":          {Fulfillment: domain.FulfillmentDispatched, TrackingEvent: domain.TrackingDispatched, Descripcion: "El transportador recogio el pedido"},
	"shipped":            {Fulfillment: domain.FulfillmentDispatched, TrackingEvent: domain.TrackingDispatched, Descripcion: "El pedido fue despachado"},
	"in_transit":         {Fulfillment: domain.FulfillmentDispatched, TrackingEvent: domain.TrackingInTransit, Descripcion: "El pedido va en camino"},
	"out_for_delivery":   {Fulfillment: domain.FulfillmentDispatched, TrackingEvent: domain.TrackingOutForDelivery, Descripcion: "El pedido salio a reparto"},
	"delivered":          {Fulfillment: domain.FulfillmentDelivered, TrackingEvent: domain.TrackingDelivered, Descripcion: "El pedido fue entregado"},
	"delivery_novelty":   {Fulfillment: domain.FulfillmentDispatched, TrackingEvent: domain.TrackingAttemptFailed, Descripcion: "Novedad en la entrega"},
	"delivery_failed":    {Fulfillment: domain.FulfillmentDispatched, TrackingEvent: domain.TrackingAttemptFailed, Descripcion: "La entrega no pudo completarse"},
	"rejected":           {Fulfillment: domain.FulfillmentDispatched, TrackingEvent: domain.TrackingReturned, Descripcion: "El cliente rechazo el pedido"},
	"return_in_transit":  {Fulfillment: domain.FulfillmentDispatched, TrackingEvent: domain.TrackingReturned, Descripcion: "El pedido va de regreso"},
	"returned":           {Fulfillment: domain.FulfillmentDispatched, TrackingEvent: domain.TrackingReturned, Descripcion: "El pedido fue devuelto"},
	"cancelled":          {Cancelar: true},
}

var ordenFulfillment = map[string]int{
	domain.FulfillmentUnpacked:   0,
	domain.FulfillmentPacked:     1,
	domain.FulfillmentDispatched: 2,
	domain.FulfillmentDelivered:  3,
}

func pasosHasta(actual, objetivo string) []string {
	desde, ok := ordenFulfillment[strings.ToUpper(strings.TrimSpace(actual))]
	if !ok {
		desde = 0
	}
	hasta, ok := ordenFulfillment[objetivo]
	if !ok || hasta <= desde {
		return nil
	}

	secuencia := []string{domain.FulfillmentUnpacked, domain.FulfillmentPacked, domain.FulfillmentDispatched, domain.FulfillmentDelivered}
	return secuencia[desde+1 : hasta+1]
}

func (uc *tiendanubeUseCase) UpdateOrderStatus(ctx context.Context, integrationID string, update domain.OrderStatusUpdate) error {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return err
	}

	if habilitado, _ := integration.Config[domain.ConfigStatusSyncEnabled].(bool); !habilitado {
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Msg("Sync de estados desactivado para la integracion Tiendanube, actualizacion omitida")
		return nil
	}

	accion, ok := accionesPorEstado[strings.ToLower(strings.TrimSpace(update.Status))]
	if !ok {
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Str("status", update.Status).
			Msg("Estado sin homologacion a Tiendanube, actualizacion omitida")
		return nil
	}

	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return err
	}

	orderID := strings.TrimSpace(update.ExternalOrderID)
	if orderID == "" {
		return domain.ErrMissingExternalOrder
	}

	if accion.Cancelar {
		if err := uc.client.CancelOrder(ctx, cred, orderID); err != nil {
			uc.logger.Error(ctx).Err(err).
				Str("integration_id", integrationID).
				Str("order_id", orderID).
				Msg("Error al cancelar la orden en Tiendanube")
			return err
		}
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Str("order_id", orderID).
			Msg("Orden cancelada en Tiendanube")
		return nil
	}

	fulfillments, err := uc.client.ListFulfillmentOrders(ctx, cred, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			uc.logger.Warn(ctx).
				Str("integration_id", integrationID).
				Str("order_id", orderID).
				Msg("Tiendanube no expone fulfillment orders para esta orden, actualizacion omitida")
			return nil
		}
		return err
	}
	if len(fulfillments) == 0 {
		uc.logger.Warn(ctx).
			Str("integration_id", integrationID).
			Str("order_id", orderID).
			Msg("La orden de Tiendanube no tiene fulfillment orders, no hay que actualizar")
		return nil
	}

	tracking := construirTracking(update)
	momento := update.HappenedAt
	if momento.IsZero() {
		momento = time.Now()
	}

	for _, fulfillment := range fulfillments {
		if err := uc.avanzarFulfillment(ctx, cred, orderID, fulfillment, accion, tracking, momento, integrationID, update); err != nil {
			return err
		}
	}

	return nil
}

func construirTracking(update domain.OrderStatusUpdate) *domain.TrackingInfo {
	if strings.TrimSpace(update.TrackingNumber) == "" {
		return nil
	}
	return &domain.TrackingInfo{
		Code:           update.TrackingNumber,
		URL:            update.TrackingURL,
		NotifyCustomer: true,
	}
}

func (uc *tiendanubeUseCase) avanzarFulfillment(
	ctx context.Context,
	cred domain.Credential,
	orderID string,
	fulfillment domain.FulfillmentOrder,
	accion accionCanal,
	tracking *domain.TrackingInfo,
	momento time.Time,
	integrationID string,
	update domain.OrderStatusUpdate,
) error {
	for _, paso := range pasosHasta(fulfillment.Status, accion.Fulfillment) {
		var info *domain.TrackingInfo
		if paso == domain.FulfillmentDispatched {
			info = tracking
		}
		if err := uc.client.UpdateFulfillmentOrder(ctx, cred, orderID, fulfillment.ID, paso, info); err != nil {
			uc.logger.Error(ctx).Err(err).
				Str("integration_id", integrationID).
				Str("order_id", orderID).
				Str("fulfillment_order_id", fulfillment.ID).
				Str("paso", paso).
				Msg("Error al avanzar el fulfillment order en Tiendanube")
			return err
		}
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Str("order_id", orderID).
			Str("fulfillment_order_id", fulfillment.ID).
			Str("paso", paso).
			Msg("Fulfillment order actualizado en Tiendanube")
	}

	if accion.TrackingEvent == "" {
		return nil
	}

	evento := domain.TrackingEvent{
		Status:      accion.TrackingEvent,
		Description: accion.Descripcion,
		Address:     strings.TrimSpace(update.City),
		HappenedAt:  momento,
	}
	if err := uc.client.CreateTrackingEvent(ctx, cred, orderID, fulfillment.ID, evento); err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Str("order_id", orderID).
			Str("fulfillment_order_id", fulfillment.ID).
			Str("evento", accion.TrackingEvent).
			Msg("Error al registrar el evento de seguimiento en Tiendanube")
		return err
	}

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("order_id", orderID).
		Str("fulfillment_order_id", fulfillment.ID).
		Str("evento", accion.TrackingEvent).
		Msg("Evento de seguimiento registrado en Tiendanube")

	return nil
}
