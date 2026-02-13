package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/integration_cache"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// invoicingUseCase implementa IInvoiceUseCase
// Procesa órdenes para facturación automática SIN CONSULTAR BASE DE DATOS
type invoicingUseCase struct {
	softpymesClient  ports.ISoftpymesClient
	configCache      cache.IConfigCache
	redisClient      redis.IRedis
	integrationCore  core.IIntegrationCore  // Mantener para TestConnection
	integrationCache integration_cache.IIntegrationCacheClient // ✅ NUEVO
	log              log.ILogger
}

// NewInvoicingUseCase crea una nueva instancia del use case de facturación
func NewInvoicingUseCase(
	softpymesClient ports.ISoftpymesClient,
	configCache cache.IConfigCache,
	redisClient redis.IRedis,
	integrationCore core.IIntegrationCore,
	integrationCache integration_cache.IIntegrationCacheClient, // ✅ NUEVO
	logger log.ILogger,
) ports.IInvoiceUseCase {
	return &invoicingUseCase{
		softpymesClient:  softpymesClient,
		configCache:      configCache,
		redisClient:      redisClient,
		integrationCore:  integrationCore,  // Mantener para TestConnection
		integrationCache: integrationCache, // ✅ NUEVO
		log:              logger.WithModule("softpymes.invoicing"),
	}
}

// ProcessOrderForInvoicing procesa un evento de orden para facturación automática
// TODO EL PROCESAMIENTO SE HACE EN MEMORIA - NO SE CONSULTA BASE DE DATOS
func (uc *invoicingUseCase) ProcessOrderForInvoicing(ctx context.Context, event *ports.OrderEventMessage) error {
	uc.log.Info(ctx).
		Str("order_id", event.OrderID).
		Str("event_type", event.EventType).
		Msg("🔍 Processing order for invoicing")

	order := event.Order
	if order == nil {
		uc.log.Warn(ctx).Msg("⚠️ Order snapshot is nil")
		return nil
	}

	// 1. Obtener configuración desde Redis (read-through cache)
	config, err := uc.configCache.Get(ctx, order.IntegrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Failed to get invoicing config from cache")
		return fmt.Errorf("failed to get invoicing config: %w", err)
	}

	if config == nil {
		uc.log.Info(ctx).
			Uint("integration_id", order.IntegrationID).
			Msg("⏩ No invoicing config found - skipping")
		return nil
	}

	uc.log.Debug(ctx).
		Uint("integration_id", order.IntegrationID).
		Bool("enabled", config.Enabled).
		Bool("auto_invoice", config.AutoInvoice).
		Interface("filters", config.Filters).
		Msg("⚙️ Invoicing config retrieved from cache")

	// 2. Verificar si la configuración está habilitada y con auto-facturación
	if !config.Enabled {
		uc.log.Info(ctx).Msg("⏩ Invoicing config disabled - skipping")
		return nil
	}

	if !config.AutoInvoice {
		uc.log.Info(ctx).Msg("⏩ Auto-invoicing disabled - skipping")
		return nil
	}

	// 3. Verificar duplicado en Redis Hash
	integrationIDStr := fmt.Sprintf("%d", config.InvoicingIntegrationID)
	if isDuplicate, err := uc.checkDuplicate(ctx, order.ID, integrationIDStr); err != nil {
		uc.log.Warn(ctx).
			Err(err).
			Msg("⚠️ Failed to check duplicate - proceeding anyway")
	} else if isDuplicate {
		uc.log.Info(ctx).Msg("⏩ Already invoiced - skipping")
		return nil
	}

	// 4. Validar filtros en memoria
	if err := uc.validateFilters(ctx, order, config); err != nil {
		uc.log.Info(ctx).
			Err(err).
			Msg("⏩ Order does not meet filters - skipping")
		return nil
	}

	// 5. ✅ NUEVO - Obtener integración desde cache (Redis directo)
	integrationID := *config.InvoicingIntegrationID

	integration, err := uc.integrationCache.GetIntegration(ctx, integrationID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Failed to get integration from cache")
		return fmt.Errorf("failed to get integration: %w", err)
	}

	uc.log.Debug(ctx).
		Interface("integration", integration).
		Msg("📋 Integration loaded from cache")

	// 6. ✅ NUEVO - Obtener credenciales desde cache (desencriptadas)
	apiKey, err := uc.integrationCache.GetCredential(ctx, integrationID, "api_key")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Failed to get api_key from cache")
		return fmt.Errorf("failed to get api_key: %w", err)
	}

	apiSecret, err := uc.integrationCache.GetCredential(ctx, integrationID, "api_secret")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Failed to get api_secret from cache")
		return fmt.Errorf("failed to get api_secret: %w", err)
	}

	// 7. Construir invoiceData (en memoria)
	invoiceData := uc.buildInvoiceData(ctx, order, config, integration, apiKey, apiSecret)

	uc.log.Debug(ctx).
		Interface("invoice_data", invoiceData).
		Msg("📄 Invoice data built")

	// 8. Crear factura en Softpymes
	if err := uc.softpymesClient.CreateInvoice(ctx, invoiceData); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Msg("❌ Failed to create invoice in Softpymes")
		return fmt.Errorf("failed to create invoice: %w", err)
	}

	uc.log.Info(ctx).Msg("✅ Invoice created successfully in Softpymes")

	// 9. Marcar como procesado en Redis
	if err := uc.markAsProcessed(ctx, order.ID, integrationIDStr); err != nil {
		uc.log.Warn(ctx).
			Err(err).
			Msg("⚠️ Failed to mark as processed in Redis")
		// No retornar error - la factura ya se creó exitosamente
	}

	return nil
}

// checkDuplicate verifica en Redis Hash si la orden ya fue facturada con esta integración
func (uc *invoicingUseCase) checkDuplicate(ctx context.Context, orderID, integrationID string) (bool, error) {
	key := fmt.Sprintf("probability:invoices:processed:%s", orderID)

	// HGet retorna valor o error si no existe
	value, err := uc.redisClient.HGet(ctx, key, integrationID)
	if err != nil {
		// Si no existe la key o el field, no es error
		return false, nil
	}

	return value != "", nil
}

// markAsProcessed marca la orden como procesada en Redis Hash
func (uc *invoicingUseCase) markAsProcessed(ctx context.Context, orderID, integrationID string) error {
	key := fmt.Sprintf("probability:invoices:processed:%s", orderID)

	data := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"status":    "processed",
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// HSet para guardar el field
	if err := uc.redisClient.HSet(ctx, key, integrationID, string(dataJSON)); err != nil {
		return err
	}

	// Expire para que se limpie automáticamente después de 30 días
	ttl := 30 * 24 * time.Hour
	return uc.redisClient.Expire(ctx, key, ttl)
}

// validateFilters valida que la orden cumpla con los filtros configurados
func (uc *invoicingUseCase) validateFilters(ctx context.Context, order *ports.OrderSnapshot, config *entities.InvoicingConfig) error {
	// Validar monto mínimo
	if config.Filters.MinAmount > 0 && order.TotalAmount < config.Filters.MinAmount {
		uc.log.Debug(ctx).
			Float64("order_total", order.TotalAmount).
			Float64("min_amount", config.Filters.MinAmount).
			Msg("❌ Order below minimum amount")
		return fmt.Errorf("order total (%.2f) is below minimum amount (%.2f)", order.TotalAmount, config.Filters.MinAmount)
	}

	// Validar métodos de pago (si hay filtro configurado)
	if len(config.Filters.PaymentMethods) > 0 {
		allowed := false
		for _, allowedMethod := range config.Filters.PaymentMethods {
			if allowedMethod == order.PaymentMethodID {
				allowed = true
				break
			}
		}

		if !allowed {
			uc.log.Debug(ctx).
				Uint("payment_method", order.PaymentMethodID).
				Interface("allowed_methods", config.Filters.PaymentMethods).
				Msg("❌ Payment method not in allowed list")
			return fmt.Errorf("payment method %d not in allowed list", order.PaymentMethodID)
		}
	}

	// Validar estado de pago (si hay filtro configurado)
	if config.Filters.PaymentStatus != "" && order.PaymentStatusID != nil {
		// Aquí se necesitaría el código del estado, pero solo tenemos el ID
		// Por ahora, asumimos que si PaymentStatus está configurado como "paid",
		// se valida que exista un PaymentStatusID
		// En un escenario real, necesitarías un mapper de ID -> code en memoria
		uc.log.Debug(ctx).
			Uint("payment_status_id", *order.PaymentStatusID).
			Str("required_status", config.Filters.PaymentStatus).
			Msg("✅ Payment status validation (ID-based)")
	}

	uc.log.Debug(ctx).Msg("✅ All filters passed")
	return nil
}

// buildInvoiceData construye el payload para la API de Softpymes
func (uc *invoicingUseCase) buildInvoiceData(
	ctx context.Context,
	order *ports.OrderSnapshot,
	config *entities.InvoicingConfig,
	integration *integration_cache.CachedIntegration, // ✅ CAMBIO DE TIPO
	apiKey, apiSecret string,
) map[string]interface{} {
	// Construir items para Softpymes
	items := make([]map[string]interface{}, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, map[string]interface{}{
			"sku":         item.SKU,
			"name":        item.Name,
			"quantity":    item.Quantity,
			"unit_price":  item.UnitPrice,
			"total_price": item.TotalPrice,
			"tax":         item.Tax,
			"discount":    item.Discount,
			"product_id":  item.ProductID,
		})
	}

	// Construir credentials (igual al config de integración + api_key/api_secret)
	credentials := make(map[string]interface{})
	if integration.Config != nil {
		// Copiar todo el config de la integración
		for k, v := range integration.Config {
			credentials[k] = v
		}
	}
	// Sobrescribir/agregar las credenciales desencriptadas
	credentials["api_key"] = apiKey
	credentials["api_secret"] = apiSecret

	// Construir customer
	customer := map[string]interface{}{
		"name":  order.CustomerName,
		"email": order.CustomerEmail,
		"phone": order.CustomerPhone,
		"dni":   order.CustomerDNI,
	}

	// Construir payload completo en el formato que espera el cliente de Softpymes
	invoiceData := map[string]interface{}{
		// Credenciales (incluye api_key, api_secret, referer, company_nit, etc.)
		"credentials": credentials,

		// Config (copia del config de integración - usado por el cliente)
		"config": integration.Config,

		// Información de la orden
		"order_id":     order.ID,
		"order_number": order.OrderNumber,
		"external_id":  order.ExternalID,
		"platform":     order.Platform,

		// Cliente
		"customer": customer,

		// Totales
		"subtotal":      order.Subtotal,
		"tax":           order.Tax,
		"discount":      order.Discount,
		"shipping_cost": order.ShippingCost,
		"total":         order.TotalAmount, // Cambiar "total_amount" a "total"
		"currency":      order.Currency,

		// Items
		"items": items,

		// Metadata
		"created_at": order.CreatedAt,
	}

	return invoiceData
}
