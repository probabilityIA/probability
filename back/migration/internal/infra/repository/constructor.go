package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Repository struct {
	db  db.IDatabase
	cfg env.IConfig
}

func New(db db.IDatabase, cfg env.IConfig) *Repository {
	return &Repository{
		db:  db,
		cfg: cfg,
	}
}

func (r *Repository) Migrate(ctx context.Context) error {
	// Clean up orphan wallet rows (business_id=0) that block FK constraint creation
	if err := r.cleanOrphanWalletRows(ctx); err != nil {
		return fmt.Errorf("failed to clean orphan wallet rows: %w", err)
	}

	// Solo migrar modelos nuevos/modificados (no todos)
	if err := r.db.Conn(ctx).AutoMigrate(
		// Wallet (agregar FK que faltaba)
		&models.Wallet{},
		&models.WalletTransaction{},

		// Warehouses & Locations (tablas nuevas)
		&models.Warehouse{},
		&models.WarehouseLocation{},

		// Stock Movement Types (debe ir antes de StockMovement por FK)
		&models.StockMovementType{},

		// Inventory (tablas nuevas)
		&models.InventoryLevel{},
		&models.StockMovement{},

		// Ampliar columnas varchar para soportar datos largos de Shopify
		// Client.Phone: 20→50, Client.Dni: 30→50
		&models.Client{},
		// Order.CustomerPhone: 32→50
		&models.Order{},
		// Invoice.CustomerPhone: 32→50
		&models.Invoice{},

		// Ultima milla (drivers, vehicles, routes)
		&models.Driver{},
		&models.Vehicle{},
		&models.Route{},
		&models.RouteStop{},

		// Subscriptions
		&models.BusinessSubscription{},
	); err != nil {
		return err
	}

	// Seed tipos de movimiento de inventario
	if err := r.seedStockMovementTypes(ctx); err != nil {
		return fmt.Errorf("failed to seed stock movement types: %w", err)
	}

	// Fix notification config unique index: el viejo idx_business_event_type
	// usaba (business_id, event_type) donde event_type es deprecated y siempre ''.
	// Reemplazar con partial unique index sobre las columnas nuevas, excluyendo soft-deleted.
	if err := r.fixNotificationConfigIndex(ctx); err != nil {
		return fmt.Errorf("failed to fix notification config index: %w", err)
	}

	// Migrar tabla pivote AllowedOrderStatuses para NotificationEventType
	if err := r.db.Conn(ctx).AutoMigrate(&models.NotificationEventType{}); err != nil {
		return fmt.Errorf("failed to auto-migrate notification_event_type (allowed statuses): %w", err)
	}

	// Seed allowed statuses por tipo de evento
	if err := r.seedAllowedOrderStatusesByEventType(ctx); err != nil {
		return fmt.Errorf("failed to seed allowed order statuses by event type: %w", err)
	}

	// Migrate email_logs table
	if err := r.db.Conn(ctx).AutoMigrate(&models.EmailLog{}); err != nil {
		return fmt.Errorf("failed to auto-migrate email_logs: %w", err)
	}

	// Seed Email integration type (id=29)
	if err := r.seedEmailIntegrationType(ctx); err != nil {
		return fmt.Errorf("failed to seed email integration type: %w", err)
	}

	// Migrate legacy order items from JSONB to order_items table
	if err := r.migrateOrderItemsFromJSONB(ctx); err != nil {
		return fmt.Errorf("failed to migrate order items from JSONB: %w", err)
	}

	// Drop the legacy items column (data is now in order_items table)
	if err := r.db.Conn(ctx).Exec("ALTER TABLE orders DROP COLUMN IF EXISTS items").Error; err != nil {
		return fmt.Errorf("failed to drop orders.items column: %w", err)
	}

	// Migrate InvoicingConfig: join table for multiple integration_ids
	if err := r.migrateInvoicingConfigIntegrations(ctx); err != nil {
		return fmt.Errorf("failed to migrate invoicing config integrations: %w", err)
	}

	// Fix client DNI unique index: was UNIQUE(dni) global, should be UNIQUE(business_id, dni)
	// to allow same DNI across different businesses (multi-tenant)
	if err := r.fixClientDniIndex(ctx); err != nil {
		return fmt.Errorf("failed to fix client DNI index: %w", err)
	}

	// Add discount_percent column to order_items
	if err := r.db.Conn(ctx).AutoMigrate(&models.OrderItem{}); err != nil {
		return fmt.Errorf("failed to auto-migrate order_items (discount_percent): %w", err)
	}

	// Seed new resources and permissions for sidebar modules
	if err := r.seedNewResourcesAndPermissions(ctx); err != nil {
		return fmt.Errorf("failed to seed new resources and permissions: %w", err)
	}

	// Add user_id column to clients table (storefront)
	if err := r.db.Conn(ctx).AutoMigrate(&models.Client{}); err != nil {
		return fmt.Errorf("failed to auto-migrate clients (user_id): %w", err)
	}

	// Seed cliente_final role, storefront resource and permissions
	if err := r.seedStorefrontRoleAndPermissions(ctx); err != nil {
		return fmt.Errorf("failed to seed storefront role and permissions: %w", err)
	}

	// Add unit_price_base columns to order_items and invoice_items
	if err := r.db.Conn(ctx).AutoMigrate(&models.OrderItem{}, &models.InvoiceItem{}); err != nil {
		return fmt.Errorf("failed to auto-migrate order_items/invoice_items (unit_price_base): %w", err)
	}

	// Update product defaults: is_active=true, status='active' and apply to all existing products
	if err := r.db.Conn(ctx).AutoMigrate(&models.Product{}); err != nil {
		return fmt.Errorf("failed to auto-migrate products (is_active default): %w", err)
	}
	if err := r.activateAllProducts(ctx); err != nil {
		return fmt.Errorf("failed to activate all products: %w", err)
	}

	// Public website config + contact submissions (tienda pública)
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.BusinessWebsiteConfig{},
		&models.ContactSubmission{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate public website tables: %w", err)
	}

	// Seed storefront integration category (id=7) and types (id=30, 31)
	if err := r.seedStorefrontIntegrationCategory(ctx); err != nil {
		return fmt.Errorf("failed to seed storefront integration category: %w", err)
	}
	if err := r.seedTiendaIntegrationType(ctx); err != nil {
		return fmt.Errorf("failed to seed tienda integration type: %w", err)
	}
	if err := r.seedTiendaWebIntegrationType(ctx); err != nil {
		return fmt.Errorf("failed to seed tienda web integration type: %w", err)
	}

	// Seed notification event type: shipment.guide_generated for WhatsApp (id=13)
	if err := r.seedShipmentGuideNotificationEventType(ctx); err != nil {
		return fmt.Errorf("failed to seed shipment guide notification event type: %w", err)
	}

	// Pricing: client pricing rules + quantity discounts
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.ClientPricingRule{},
		&models.QuantityDiscount{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate pricing tables: %w", err)
	}

	// Partial unique index for global client pricing rules (product_id IS NULL)
	if err := r.db.Conn(ctx).Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_client_pricing_rule_global
		ON client_pricing_rules (business_id, client_id)
		WHERE product_id IS NULL AND deleted_at IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to create partial index for global pricing rules: %w", err)
	}

	// Partial unique index for global quantity discounts (product_id IS NULL)
	if err := r.db.Conn(ctx).Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_qty_discount_global
		ON quantity_discounts (business_id, min_quantity)
		WHERE product_id IS NULL AND deleted_at IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to create partial index for global quantity discounts: %w", err)
	}

	// Order Statuses v2: nuevos estados de última milla + normalización
	if err := r.migrateOrderStatusesV2(ctx); err != nil {
		return fmt.Errorf("failed to migrate order statuses v2: %w", err)
	}

	// Add cash receipt audit columns to invoice_sync_logs
	if err := r.db.Conn(ctx).AutoMigrate(&models.InvoiceSyncLog{}); err != nil {
		return fmt.Errorf("failed to auto-migrate invoice_sync_logs (cash receipt audit): %w", err)
	}

	return nil
}

// migrateOrderStatusesV2 agrega los nuevos estados logísticos, normaliza códigos existentes,
// migra órdenes de estados deprecados y actualiza mappings.
// Idempotente: usa ON CONFLICT DO NOTHING y WHERE condicionales.
func (r *Repository) migrateOrderStatusesV2(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// 1. Insertar nuevos estados (idempotente via ON CONFLICT)
	if err := db.Exec(`
		INSERT INTO order_statuses (code, name, description, category, is_active, priority, color, created_at, updated_at)
		VALUES
			('picking', 'Seleccionando productos', 'Seleccionando productos del inventario', 'active', true, 10, '#3B82F6', NOW(), NOW()),
			('packing', 'Empacando', 'Empacando el pedido', 'active', true, 20, '#6366F1', NOW(), NOW()),
			('ready_to_ship', 'Listo para despacho', 'Listo para despacho', 'active', true, 30, '#8B5CF6', NOW(), NOW()),
			('assigned_to_driver', 'Asignado a piloto', 'Asignado a piloto/conductor', 'active', true, 40, '#A855F7', NOW(), NOW()),
			('picked_up', 'Recogido', 'Recogido por el piloto', 'active', true, 50, '#D946EF', NOW(), NOW()),
			('in_transit', 'En camino', 'En camino al destino', 'active', true, 60, '#EC4899', NOW(), NOW()),
			('out_for_delivery', 'En reparto final', 'En reparto final (última milla)', 'active', true, 70, '#F43F5E', NOW(), NOW()),
			('delivery_failed', 'Entrega fallida', 'Entrega fallida', 'issue', true, 76, '#EF4444', NOW(), NOW()),
			('rejected', 'Rechazado', 'Rechazado por el cliente', 'issue', true, 77, '#DC2626', NOW(), NOW()),
			('return_in_transit', 'Devolución en camino', 'Devolución en camino al almacén', 'return', true, 80, '#F59E0B', NOW(), NOW()),
			('inventory_issue', 'Novedad de inventario', 'Novedad de inventario (sin stock, producto dañado)', 'issue', true, 15, '#FB923C', NOW(), NOW())
		ON CONFLICT (code) DO NOTHING
	`).Error; err != nil {
		return fmt.Errorf("failed to insert new order statuses: %w", err)
	}

	// 2. Normalizar estados existentes: Novelty → delivery_novelty, Refund → returned
	db.Exec(`
		UPDATE order_statuses
		SET code = 'delivery_novelty', name = 'Novedad de entrega', description = 'Novedad de entrega',
			category = 'issue', priority = 75, color = '#F97316', updated_at = NOW()
		WHERE id = 10 AND code = 'Novelty'
	`)
	db.Exec(`
		UPDATE order_statuses
		SET code = 'returned', name = 'Devuelto', description = 'Devuelto al almacén',
			category = 'return', priority = 85, color = '#EAB308', updated_at = NOW()
		WHERE id = 11 AND code = 'Refund'
	`)

	// 3. Migrar órdenes de estados deprecados
	db.Exec(`
		UPDATE orders
		SET status = 'picking',
			status_id = (SELECT id FROM order_statuses WHERE code = 'picking' LIMIT 1),
			updated_at = NOW()
		WHERE status = 'processing' AND deleted_at IS NULL
	`)
	db.Exec(`
		UPDATE orders
		SET status = 'in_transit',
			status_id = (SELECT id FROM order_statuses WHERE code = 'in_transit' LIMIT 1),
			updated_at = NOW()
		WHERE status = 'shipped' AND deleted_at IS NULL
	`)

	// 4. Actualizar mappings que apuntaban a processing/shipped
	db.Exec(`
		UPDATE order_status_mappings
		SET order_status_id = (SELECT id FROM order_statuses WHERE code = 'picking' LIMIT 1),
			updated_at = NOW()
		WHERE order_status_id = (SELECT id FROM order_statuses WHERE code = 'processing' LIMIT 1)
		  AND deleted_at IS NULL
	`)
	db.Exec(`
		UPDATE order_status_mappings
		SET order_status_id = (SELECT id FROM order_statuses WHERE code = 'in_transit' LIMIT 1),
			updated_at = NOW()
		WHERE order_status_id = (SELECT id FROM order_statuses WHERE code = 'shipped' LIMIT 1)
		  AND deleted_at IS NULL
	`)

	// 5. Desactivar estados deprecados
	db.Exec(`
		UPDATE order_statuses
		SET is_active = false, updated_at = NOW()
		WHERE code IN ('processing', 'shipped')
	`)

	return nil
}

// migrateInvoicingConfigIntegrations crea la tabla pivote invoicing_config_integrations,
// migra los integration_id existentes a ella, y elimina la columna integration_id de invoicing_configs.
// Es idempotente: detecta si ya fue ejecutada verificando si la columna integration_id existe.
func (r *Repository) migrateInvoicingConfigIntegrations(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// 1. Crear la tabla join (idempotente via AutoMigrate)
	if err := db.AutoMigrate(&models.InvoicingConfigIntegration{}); err != nil {
		return fmt.Errorf("failed to auto-migrate invoicing_config_integrations: %w", err)
	}

	// 2. Si la columna legacy integration_id todavía existe, migrar y limpiar
	if db.Migrator().HasColumn("invoicing_configs", "integration_id") {
		// 3. Leer configs con integration_id legacy usando GORM (struct temporal para scan)
		type legacyRow struct {
			ID            uint
			IntegrationID uint
		}
		var legacyRows []legacyRow
		if err := db.Model(&models.InvoicingConfig{}).
			Select("id, integration_id").
			Where("integration_id IS NOT NULL AND integration_id > ?", 0).
			Scan(&legacyRows).Error; err != nil {
			return fmt.Errorf("failed to query legacy configs: %w", err)
		}

		// 4. Insertar en join table sólo las filas que no existen aún
		for _, row := range legacyRows {
			var count int64
			db.Model(&models.InvoicingConfigIntegration{}).
				Where("config_id = ? AND integration_id = ? AND deleted_at IS NULL",
					row.ID, row.IntegrationID).
				Count(&count)
			if count == 0 {
				entry := models.InvoicingConfigIntegration{
					ConfigID:      row.ID,
					IntegrationID: row.IntegrationID,
				}
				if err := db.Create(&entry).Error; err != nil {
					return fmt.Errorf("failed to create join entry for config %d: %w", row.ID, err)
				}
			}
		}

		// 5. Eliminar la columna legacy via Migrator (usa ALTER TABLE internamente, seguro)
		if err := db.Migrator().DropColumn("invoicing_configs", "integration_id"); err != nil {
			return fmt.Errorf("failed to drop integration_id column: %w", err)
		}
	}

	// 6. Eliminar el viejo unique index via Migrator
	if db.Migrator().HasIndex("invoicing_configs", "idx_business_integration_config") {
		if err := db.Migrator().DropIndex("invoicing_configs", "idx_business_integration_config"); err != nil {
			return fmt.Errorf("failed to drop idx_business_integration_config: %w", err)
		}
	}

	// 7. Deduplicar: por cada (business_id, invoicing_integration_id) con múltiples registros activos,
	// conservar el de menor id, reasignar join entries al superviviente, soft-delete el resto.
	type dupGroup struct {
		BusinessID             uint
		InvoicingIntegrationID uint
		MinID                  uint
	}
	var groups []dupGroup
	if err := db.Model(&models.InvoicingConfig{}).
		Select("business_id, invoicing_integration_id, MIN(id) AS min_id").
		Where("invoicing_integration_id IS NOT NULL AND deleted_at IS NULL").
		Group("business_id, invoicing_integration_id").
		Having("COUNT(*) > 1").
		Scan(&groups).Error; err != nil {
		return fmt.Errorf("failed to query duplicate configs: %w", err)
	}

	for _, g := range groups {
		// Obtener IDs de los configs duplicados (todos excepto el superviviente)
		var duplicateIDs []uint
		if err := db.Model(&models.InvoicingConfig{}).
			Where("business_id = ? AND invoicing_integration_id = ? AND id != ? AND deleted_at IS NULL",
				g.BusinessID, g.InvoicingIntegrationID, g.MinID).
			Pluck("id", &duplicateIDs).Error; err != nil {
			return fmt.Errorf("failed to find duplicate config ids: %w", err)
		}

		for _, dupID := range duplicateIDs {
			// Leer join entries del duplicado
			var joinEntries []models.InvoicingConfigIntegration
			if err := db.Where("config_id = ? AND deleted_at IS NULL", dupID).
				Find(&joinEntries).Error; err != nil {
				continue
			}
			// Reasignar cada entry al config superviviente si no existe ya
			for _, entry := range joinEntries {
				var exists int64
				db.Model(&models.InvoicingConfigIntegration{}).
					Where("config_id = ? AND integration_id = ? AND deleted_at IS NULL",
						g.MinID, entry.IntegrationID).
					Count(&exists)
				if exists == 0 {
					db.Model(&models.InvoicingConfigIntegration{}).
						Where("id = ?", entry.ID).
						Update("config_id", g.MinID)
				}
			}
			// Soft-delete el config duplicado via GORM (respeta gorm.Model.DeletedAt)
			db.Delete(&models.InvoicingConfig{}, dupID)
		}
	}

	// 8. Crear índice único parcial (DDL estático, sin valores de usuario — no hay riesgo de inyección)
	// GORM no soporta índices parciales via struct tags, por eso se usa Exec con SQL estático.
	if !db.Migrator().HasIndex("invoicing_configs", "idx_business_invoicing_integration") {
		if err := db.Exec(
			`CREATE UNIQUE INDEX idx_business_invoicing_integration
			 ON invoicing_configs (business_id, invoicing_integration_id)
			 WHERE invoicing_integration_id IS NOT NULL AND deleted_at IS NULL`,
		).Error; err != nil {
			return fmt.Errorf("failed to create idx_business_invoicing_integration: %w", err)
		}
	}

	return nil
}

// createDefaultUserIfNotExists crea el usuario principal solo si no existe
// Las migraciones solo deben crear la estructura de tablas, no datos adicionales
func (r *Repository) createDefaultUserIfNotExists(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// Verificar si el usuario principal ya existe
	var count int64
	if err := db.Model(&models.User{}).Where("id = ?", 1).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check for existing user: %w", err)
	}

	// Solo crear el usuario si no existe
	if count == 0 {
		email := r.cfg.Get("EMAIL_USER_DEFAULT")
		password := r.cfg.Get("USER_PASS_DEFAULT")

		if email == "" || password == "" {
			return fmt.Errorf("EMAIL_USER_DEFAULT or USER_PASS_DEFAULT not set")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		user := models.User{
			Model:    gorm.Model{ID: 1},
			Name:     "Admin",
			Email:    email,
			Password: string(hashedPassword),
			IsActive: true,
		}

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create default user: %w", err)
		}
	}

	return nil
}

// seedShipmentGuideNotificationEventType inserta el notification_event_type para shipment.guide_generated (WhatsApp)
func (r *Repository) seedShipmentGuideNotificationEventType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	net := models.NotificationEventType{
		Model:              gorm.Model{ID: 13},
		NotificationTypeID: 2,
		EventCode:          "shipment.guide_generated",
		EventName:          "Guía de Envío Generada",
		Description:        "Envia un mensaje de WhatsApp al cliente con el numero de guia y transportadora cuando se genera la guia de envio.",
		IsActive:           true,
	}

	var existing models.NotificationEventType
	err := db.Where("id = ?", net.ID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		if err := db.Create(&net).Error; err != nil {
			return fmt.Errorf("failed to create notification_event_type shipment.guide_generated: %w", err)
		}
	}

	return nil
}

// migrateBusinessNotificationConfigData migra los datos existentes de business_notification_configs
// a la nueva estructura con integration_id, notification_type_id, notification_event_type_id
func (r *Repository) migrateBusinessNotificationConfigData(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// 1. Poblar integration_id con la primera integración activa de cada business
	if err := db.Exec(`
		UPDATE business_notification_configs bnc
		SET integration_id = (
			SELECT i.id
			FROM integrations i
			WHERE i.business_id = bnc.business_id
			AND i.is_active = true
			LIMIT 1
		)
		WHERE integration_id IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to populate integration_id: %w", err)
	}

	// 2. Eliminar configs huérfanas (businesses sin integraciones)
	if err := db.Exec(`
		DELETE FROM business_notification_configs
		WHERE integration_id IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to delete orphaned configs: %w", err)
	}

	// 3. Poblar notification_type_id con 1 (SSE) por defecto
	if err := db.Exec(`
		UPDATE business_notification_configs
		SET notification_type_id = 1
		WHERE notification_type_id IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to populate notification_type_id: %w", err)
	}

	// 4. Poblar notification_event_type_id basándose en event_type
	if err := db.Exec(`
		UPDATE business_notification_configs bnc
		SET notification_event_type_id = (
			SELECT id
			FROM notification_event_types
			WHERE notification_type_id = bnc.notification_type_id
			AND event_code = bnc.event_type
			LIMIT 1
		)
		WHERE notification_event_type_id IS NULL
		AND event_type IS NOT NULL
		AND event_type != ''
	`).Error; err != nil {
		return fmt.Errorf("failed to map notification_event_type_id from event_type: %w", err)
	}

	// 5. Poblar notification_event_type_id con el primer evento disponible si aún es NULL
	if err := db.Exec(`
		UPDATE business_notification_configs bnc
		SET notification_event_type_id = (
			SELECT id
			FROM notification_event_types
			WHERE notification_type_id = bnc.notification_type_id
			LIMIT 1
		)
		WHERE notification_event_type_id IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to populate notification_event_type_id with default: %w", err)
	}

	// 6. Eliminar configs que no se pudieron migrar (no deberían existir si los datos son consistentes)
	if err := db.Exec(`
		DELETE FROM business_notification_configs
		WHERE notification_event_type_id IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to delete unmigrated configs: %w", err)
	}

	return nil
}

// seedIntegrationCategories inserta las categorías iniciales de integraciones
func (r *Repository) seedIntegrationCategories(ctx context.Context) error {
	db := r.db.Conn(ctx)

	categories := []models.IntegrationCategory{
		{
			Model:        gorm.Model{ID: 1},
			Code:         "ecommerce",
			Name:         "E-commerce",
			Description:  "Plataformas de venta online",
			Icon:         "shopping-cart",
			Color:        "#3B82F6",
			DisplayOrder: 1,
			IsActive:     true,
			IsVisible:    true,
		},
		{
			Model:        gorm.Model{ID: 2},
			Code:         "invoicing",
			Name:         "Facturación Electrónica",
			Description:  "Proveedores de facturación",
			Icon:         "receipt",
			Color:        "#10B981",
			DisplayOrder: 2,
			IsActive:     true,
			IsVisible:    true,
		},
		{
			Model:        gorm.Model{ID: 3},
			Code:         "messaging",
			Name:         "Mensajería",
			Description:  "Canales de comunicación",
			Icon:         "message-circle",
			Color:        "#8B5CF6",
			DisplayOrder: 3,
			IsActive:     true,
			IsVisible:    true,
		},
		{
			Model:        gorm.Model{ID: 4},
			Code:         "payment",
			Name:         "Pagos",
			Description:  "Pasarelas de pago",
			Icon:         "credit-card",
			Color:        "#F59E0B",
			DisplayOrder: 4,
			IsActive:     false,
			IsVisible:    true,
		},
		{
			Model:        gorm.Model{ID: 5},
			Code:         "shipping",
			Name:         "Logística",
			Description:  "Operadores logísticos",
			Icon:         "truck",
			Color:        "#EF4444",
			DisplayOrder: 5,
			IsActive:     true,
			IsVisible:    true,
		},
		{
			Model:        gorm.Model{ID: 6},
			Code:         "platform",
			Name:         "Plataforma",
			Description:  "Órdenes creadas directamente en la plataforma",
			Icon:         "squares-plus",
			Color:        "#6366F1",
			DisplayOrder: 0,
			IsActive:     true,
			IsVisible:    true,
		},
	}

	for _, category := range categories {
		var existing models.IntegrationCategory
		err := db.Where("id = ?", category.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// No existe, crear
			if err := db.Create(&category).Error; err != nil {
				return fmt.Errorf("failed to create integration_category %s: %w", category.Code, err)
			}
		}
	}

	return nil
}

// migrateIntegrationTypesToCategories actualiza los integration_types existentes con category_id
func (r *Repository) migrateIntegrationTypesToCategories(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// Mapeo de códigos de integration_types a category_id
	migrations := []struct {
		typeCode   string
		categoryID uint
	}{
		{"shopify", 1},      // ecommerce
		{"mercadolibre", 1}, // ecommerce
		{"amazon", 1},       // ecommerce
		{"whatsapp", 3},     // messaging
		{"whatsap", 3},      // messaging (typo histórico)
	}

	for _, m := range migrations {
		if err := db.Exec(`
			UPDATE integration_types
			SET category_id = ?
			WHERE code = ?
			AND category_id IS NULL
		`, m.categoryID, m.typeCode).Error; err != nil {
			return fmt.Errorf("failed to update category_id for %s: %w", m.typeCode, err)
		}
	}

	return nil
}

// seedSoftpymesIntegrationType crea el tipo de integración para Softpymes (facturación electrónica)
func (r *Repository) seedSoftpymesIntegrationType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	integrationType := models.IntegrationType{
		Model:       gorm.Model{ID: 5},
		Name:        "Softpymes Facturación",
		Code:        "softpymes",
		CategoryID:  ptrUint(2), // invoicing category
		IsActive:    true,
		Description: "Proveedor de facturación electrónica Softpymes",
		Icon:        "receipt",
	}

	var existing models.IntegrationType
	err := db.Where("id = ?", 5).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		// No existe, crear
		if err := db.Create(&integrationType).Error; err != nil {
			return fmt.Errorf("failed to create softpymes integration type: %w", err)
		}
	}

	return nil
}

// migrateInvoicingProvidersToIntegrations migra datos de invoicing_providers → integrations
func (r *Repository) migrateInvoicingProvidersToIntegrations(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// 1. Verificar si ya se migró
	var count int64
	if err := db.Model(&models.Integration{}).
		Where("integration_type_id = ? AND store_id LIKE 'softpymes-%'", 5).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing migrations: %w", err)
	}

	if count > 0 {
		// Ya se migró, salir
		return nil
	}

	// 2. Migrar invoicing_providers → integrations
	err := db.Exec(`
		INSERT INTO integrations (
			business_id, integration_type_id, name, description,
			store_id, credentials, config, is_active, is_default,
			created_at, updated_at, created_by_id
		)
		SELECT
			ip.business_id,
			5 AS integration_type_id,
			ip.name,
			ip.description,
			CONCAT('softpymes-', ip.id) AS store_id,
			ip.credentials,
			ip.config,
			ip.is_active,
			ip.is_default,
			ip.created_at,
			ip.updated_at,
			ip.created_by_id
		FROM invoicing_providers ip
		WHERE NOT EXISTS (
			SELECT 1 FROM integrations i
			WHERE i.store_id = CONCAT('softpymes-', ip.id)
		)
	`).Error

	if err != nil {
		return fmt.Errorf("failed to migrate invoicing_providers: %w", err)
	}

	// 3. Actualizar invoicing_configs con invoicing_integration_id
	err = db.Exec(`
		UPDATE invoicing_configs ic
		SET invoicing_integration_id = i.id
		FROM integrations i
		JOIN invoicing_providers ip ON i.store_id = CONCAT('softpymes-', ip.id)
		WHERE ic.invoicing_provider_id = ip.id
		AND ic.invoicing_integration_id IS NULL
	`).Error

	if err != nil {
		return fmt.Errorf("failed to update invoicing_configs: %w", err)
	}

	// 4. Actualizar invoices con invoicing_integration_id
	err = db.Exec(`
		UPDATE invoices inv
		SET invoicing_integration_id = i.id
		FROM integrations i
		JOIN invoicing_providers ip ON i.store_id = CONCAT('softpymes-', ip.id)
		WHERE inv.invoicing_provider_id = ip.id
		AND inv.invoicing_integration_id IS NULL
	`).Error

	if err != nil {
		return fmt.Errorf("failed to update invoices: %w", err)
	}

	return nil
}

// seedPlatformIntegrationType crea el tipo de integración para Plataforma (órdenes manuales)
func (r *Repository) seedPlatformIntegrationType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	integrationType := models.IntegrationType{
		Model:       gorm.Model{ID: 6},
		Name:        "Plataforma",
		Code:        "platform",
		CategoryID:  ptrUint(6), // platform category
		IsActive:    true,
		Description: "Órdenes creadas directamente en la plataforma",
		Icon:        "squares-plus",
	}

	var existing models.IntegrationType
	err := db.Where("id = ?", 6).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		if err := db.Create(&integrationType).Error; err != nil {
			return fmt.Errorf("failed to create platform integration type: %w", err)
		}
	}

	return nil
}

// seedPlatformIntegrationsForBusinesses crea una integración de plataforma por cada negocio activo que no tenga una
func (r *Repository) seedPlatformIntegrationsForBusinesses(ctx context.Context) error {
	db := r.db.Conn(ctx)

	err := db.Exec(`
		INSERT INTO integrations (name, code, category, integration_type_id, business_id, is_active, is_default, created_by_id, created_at, updated_at)
		SELECT 'Plataforma', CONCAT('platform_', b.id), 'platform', 6, b.id, true, false, 1, NOW(), NOW()
		FROM business b
		WHERE b.deleted_at IS NULL
		AND NOT EXISTS (
			SELECT 1 FROM integrations i
			WHERE i.business_id = b.id
			AND i.integration_type_id = 6
			AND i.deleted_at IS NULL
		)
	`).Error

	if err != nil {
		return fmt.Errorf("failed to seed platform integrations for businesses: %w", err)
	}

	return nil
}

// seedEnvioClickIntegrationType crea el tipo de integración para EnvioClick (transporte)
func (r *Repository) seedEnvioClickIntegrationType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	integrationType := models.IntegrationType{
		Model:       gorm.Model{ID: 12},
		Name:        "EnvioClick",
		Code:        "envioclick",
		CategoryID:  ptrUint(5), // shipping category
		IsActive:    true,
		Description: "Plataforma de envíos EnvioClick",
		Icon:        "truck",
	}

	var existing models.IntegrationType
	err := db.Where("id = ?", 12).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		if err := db.Create(&integrationType).Error; err != nil {
			return fmt.Errorf("failed to create envioclick integration type: %w", err)
		}
	}

	return nil
}

// seedEnvioClickBaseURL actualiza la base_url de producción para EnvioClick si no está configurada
func (r *Repository) seedEnvioClickBaseURL(ctx context.Context) error {
	return r.db.Conn(ctx).Exec(`
		UPDATE integration_types
		SET base_url = ?
		WHERE id = 12
		  AND (base_url IS NULL OR base_url = '')
		  AND deleted_at IS NULL
	`, "https://api.envioclickpro.com.co/api/v2").Error
}

// markInDevelopmentIntegrationTypes marca los integration types con IDs 13-21 como en desarrollo
func (r *Repository) markInDevelopmentIntegrationTypes(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// IDs 13-21 corresponden a los skeleton de integraciones nuevas:
	// enviame, tu, mipaquete, vtex, tiendanube, magento, amazon, falabella, exito
	if err := db.Exec(`
		UPDATE integration_types
		SET in_development = true
		WHERE id BETWEEN 13 AND 21
		AND deleted_at IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to mark in-development integration types: %w", err)
	}

	return nil
}

// addIsTestingColumns agrega las columnas de modo testing si no existen
func (r *Repository) addIsTestingColumns(ctx context.Context) error {
	db := r.db.Conn(ctx)
	if err := db.Exec(`ALTER TABLE integrations ADD COLUMN IF NOT EXISTS is_testing BOOLEAN NOT NULL DEFAULT FALSE`).Error; err != nil {
		return fmt.Errorf("failed to add is_testing to integrations: %w", err)
	}
	if err := db.Exec(`ALTER TABLE shipments ADD COLUMN IF NOT EXISTS is_test BOOLEAN NOT NULL DEFAULT FALSE`).Error; err != nil {
		return fmt.Errorf("failed to add is_test to shipments: %w", err)
	}
	return nil
}

// addPlatformCredentialsToIntegrationTypes agrega la columna de credenciales de plataforma
func (r *Repository) addPlatformCredentialsToIntegrationTypes(ctx context.Context) error {
	return r.db.Conn(ctx).Exec(`
		ALTER TABLE integration_types
		ADD COLUMN IF NOT EXISTS platform_credentials_encrypted BYTEA
	`).Error
}

// seedPaymentIntegrationTypes crea los tipos de integración para pasarelas de pago
func (r *Repository) seedPaymentIntegrationTypes(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// Activar la categoría de pagos (id=4) en caso de que esté desactivada
	if err := db.Exec(`
		UPDATE integration_categories
		SET is_active = true
		WHERE id = 4 AND deleted_at IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to activate payment category: %w", err)
	}

	paymentTypes := []models.IntegrationType{
		{
			Model:         gorm.Model{ID: 22},
			Name:          "Nequi",
			Code:          "nequi",
			CategoryID:    ptrUint(4),
			IsActive:      true,
			InDevelopment: false,
			Description:   "Pasarela de pago Nequi - Colombia",
			Icon:          "credit-card",
			BaseURL:       "https://api.sandbox.connect.nequi.com",
			BaseURLTest:   "https://api.sandbox.connect.nequi.com",
		},
		{
			Model:         gorm.Model{ID: 23},
			Name:          "Bold",
			Code:          "bold",
			CategoryID:    ptrUint(4),
			IsActive:      true,
			InDevelopment: false,
			Description:   "Pasarela de pago Bold - Colombia",
			Icon:          "credit-card",
			BaseURL:       "https://integrations.api.bold.co",
			BaseURLTest:   "https://integrations.api.bold.co",
		},
		{
			Model:         gorm.Model{ID: 24},
			Name:          "Wompi",
			Code:          "wompi",
			CategoryID:    ptrUint(4),
			IsActive:      true,
			InDevelopment: true,
			Description:   "Pasarela de pago Wompi - Colombia",
			Icon:          "credit-card",
			BaseURL:       "https://production.wompi.co/v1",
			BaseURLTest:   "https://sandbox.wompi.co/v1",
		},
		{
			Model:         gorm.Model{ID: 25},
			Name:          "Stripe",
			Code:          "stripe",
			CategoryID:    ptrUint(4),
			IsActive:      true,
			InDevelopment: true,
			Description:   "Pasarela de pago Stripe - Internacional",
			Icon:          "credit-card",
			BaseURL:       "https://api.stripe.com",
			BaseURLTest:   "https://api.stripe.com",
		},
		{
			Model:         gorm.Model{ID: 26},
			Name:          "PayU",
			Code:          "payu",
			CategoryID:    ptrUint(4),
			IsActive:      true,
			InDevelopment: true,
			Description:   "Pasarela de pago PayU - Latam",
			Icon:          "credit-card",
			BaseURL:       "https://api.payulatam.com",
			BaseURLTest:   "https://sandbox.api.payulatam.com",
		},
		{
			Model:         gorm.Model{ID: 27},
			Name:          "ePayco",
			Code:          "epayco",
			CategoryID:    ptrUint(4),
			IsActive:      true,
			InDevelopment: true,
			Description:   "Pasarela de pago ePayco - Colombia",
			Icon:          "credit-card",
			BaseURL:       "https://secure.epayco.io",
			BaseURLTest:   "https://secure.epayco.io",
		},
		{
			Model:         gorm.Model{ID: 28},
			Name:          "MercadoPago",
			Code:          "melipago",
			CategoryID:    ptrUint(4),
			IsActive:      true,
			InDevelopment: true,
			Description:   "Pasarela de pago MercadoPago - Latam",
			Icon:          "credit-card",
			BaseURL:       "https://api.mercadopago.com",
			BaseURLTest:   "https://api.mercadopago.com",
		},
	}

	for _, pt := range paymentTypes {
		var existing models.IntegrationType
		err := db.Where("id = ?", pt.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&pt).Error; err != nil {
				return fmt.Errorf("failed to create payment integration type %s: %w", pt.Code, err)
			}
		}
	}

	return nil
}

// seedChannelPaymentMethods inserta los métodos de pago nativos por canal de venta
func (r *Repository) seedChannelPaymentMethods(ctx context.Context) error {
	db := r.db.Conn(ctx)

	methods := []models.ChannelPaymentMethod{
		// Shopify
		{Model: gorm.Model{ID: 1}, IntegrationType: "shopify", Code: "shopify_payments", Name: "Shopify Payments", IsActive: true, DisplayOrder: 1},
		{Model: gorm.Model{ID: 2}, IntegrationType: "shopify", Code: "manual", Name: "Manual", IsActive: true, DisplayOrder: 2},
		{Model: gorm.Model{ID: 3}, IntegrationType: "shopify", Code: "cash_on_delivery", Name: "Contra entrega", IsActive: true, DisplayOrder: 3},
		{Model: gorm.Model{ID: 4}, IntegrationType: "shopify", Code: "bank_transfer", Name: "Transferencia bancaria", IsActive: true, DisplayOrder: 4},
		// MercadoLibre
		{Model: gorm.Model{ID: 5}, IntegrationType: "mercado_libre", Code: "account_money", Name: "Dinero en cuenta", IsActive: true, DisplayOrder: 1},
		{Model: gorm.Model{ID: 6}, IntegrationType: "mercado_libre", Code: "credit_card", Name: "Tarjeta de crédito", IsActive: true, DisplayOrder: 2},
		{Model: gorm.Model{ID: 7}, IntegrationType: "mercado_libre", Code: "debit_card", Name: "Tarjeta débito", IsActive: true, DisplayOrder: 3},
		{Model: gorm.Model{ID: 8}, IntegrationType: "mercado_libre", Code: "ticket", Name: "Efectivo (ticket)", IsActive: true, DisplayOrder: 4},
		// WhatsApp
		{Model: gorm.Model{ID: 9}, IntegrationType: "whatsapp", Code: "cash", Name: "Efectivo", IsActive: true, DisplayOrder: 1},
		{Model: gorm.Model{ID: 10}, IntegrationType: "whatsapp", Code: "bank_transfer", Name: "Transferencia bancaria", IsActive: true, DisplayOrder: 2},
		{Model: gorm.Model{ID: 11}, IntegrationType: "whatsapp", Code: "nequi", Name: "Nequi", IsActive: true, DisplayOrder: 3},
		{Model: gorm.Model{ID: 12}, IntegrationType: "whatsapp", Code: "daviplata", Name: "Daviplata", IsActive: true, DisplayOrder: 4},
	}

	for _, m := range methods {
		var existing models.ChannelPaymentMethod
		err := db.Where("id = ?", m.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&m).Error; err != nil {
				return fmt.Errorf("failed to create channel_payment_method %s/%s: %w", m.IntegrationType, m.Code, err)
			}
		}
	}
	return nil
}

// migrateOrderStatusPriority mueve el campo priority de order_status_mappings a order_statuses.
// El mapeo hereda la prioridad del estado de Probability al que apunta.
func (r *Repository) migrateOrderStatusPriority(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// 1. Agregar columna priority a order_statuses si no existe (AutoMigrate ya lo hace,
	//    pero el ALTER explícito garantiza ejecución ordenada)
	if err := db.Exec(`
		ALTER TABLE order_statuses ADD COLUMN IF NOT EXISTS priority INT NOT NULL DEFAULT 0
	`).Error; err != nil {
		return fmt.Errorf("failed to add priority to order_statuses: %w", err)
	}

	// 2. Seed de prioridades según el ciclo de vida de una orden
	//    Mayor número = estado más avanzado en el ciclo
	type prioritySeed struct {
		id       uint
		priority int
	}
	seeds := []prioritySeed{
		{1, 1}, // pending
		{2, 2}, // processing
		{9, 3}, // on_hold
		{3, 4}, // shipped
		{4, 5}, // delivered
		{5, 6}, // completed
		{7, 7}, // refunded
		{6, 8}, // cancelled
		{8, 9}, // failed
	}
	for _, s := range seeds {
		if err := db.Exec(`
			UPDATE order_statuses SET priority = ? WHERE id = ? AND priority = 0
		`, s.priority, s.id).Error; err != nil {
			return fmt.Errorf("failed to seed priority for order_status id=%d: %w", s.id, err)
		}
	}

	// 3. Eliminar columna priority de order_status_mappings si existe
	if err := db.Exec(`
		ALTER TABLE order_status_mappings DROP COLUMN IF EXISTS priority
	`).Error; err != nil {
		return fmt.Errorf("failed to drop priority from order_status_mappings: %w", err)
	}

	return nil
}

// seedIntegrationChannelStatuses inserta los estados nativos de los canales ecommerce.
// Se hace lookup por code (no hardcoded IDs) para que sea idempotente.
func (r *Repository) seedIntegrationChannelStatuses(ctx context.Context) error {
	db := r.db.Conn(ctx)

	type typeSeed struct {
		code     string
		statuses []models.IntegrationChannelStatus
	}

	channelData := []typeSeed{
		{
			code: "Shopify",
			statuses: []models.IntegrationChannelStatus{
				{Code: "pending", Name: "Pendiente", DisplayOrder: 1, IsActive: true},
				{Code: "authorized", Name: "Autorizada", DisplayOrder: 2, IsActive: true},
				{Code: "paid", Name: "Pagada", DisplayOrder: 3, IsActive: true},
				{Code: "partially_paid", Name: "Parcialmente pagada", DisplayOrder: 4, IsActive: true},
				{Code: "refunded", Name: "Reembolsada", DisplayOrder: 5, IsActive: true},
				{Code: "partially_refunded", Name: "Parcialmente reembolsada", DisplayOrder: 6, IsActive: true},
				{Code: "voided", Name: "Anulada", DisplayOrder: 7, IsActive: true},
				{Code: "cancelled", Name: "Cancelada", DisplayOrder: 8, IsActive: true},
			},
		},
		{
			code: "Mercado Libre",
			statuses: []models.IntegrationChannelStatus{
				{Code: "pending", Name: "Pendiente", DisplayOrder: 1, IsActive: true},
				{Code: "payment_required", Name: "Pago requerido", DisplayOrder: 2, IsActive: true},
				{Code: "payment_in_process", Name: "Pago en proceso", DisplayOrder: 3, IsActive: true},
				{Code: "partially_paid", Name: "Parcialmente pagado", DisplayOrder: 4, IsActive: true},
				{Code: "paid", Name: "Pagado", DisplayOrder: 5, IsActive: true},
				{Code: "money_returned", Name: "Dinero devuelto", DisplayOrder: 6, IsActive: true},
				{Code: "cancelled", Name: "Cancelado", DisplayOrder: 7, IsActive: true},
			},
		},
	}

	for _, ch := range channelData {
		var it models.IntegrationType
		if err := db.Where("code = ?", ch.code).First(&it).Error; err != nil {
			// Si el tipo no existe aún, saltamos
			continue
		}

		for _, s := range ch.statuses {
			s.IntegrationTypeID = it.ID
			var existing models.IntegrationChannelStatus
			err := db.Where("integration_type_id = ? AND code = ?", it.ID, s.Code).First(&existing).Error
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&s).Error; err != nil {
					return fmt.Errorf("failed to create channel status %s/%s: %w", ch.code, s.Code, err)
				}
			}
		}
	}

	return nil
}

// seedStockMovementTypes inserta los tipos de movimiento de inventario por defecto
func (r *Repository) seedStockMovementTypes(ctx context.Context) error {
	db := r.db.Conn(ctx)

	types := []models.StockMovementType{
		{Model: gorm.Model{ID: 1}, Code: "inbound", Name: "Entrada de mercancía", Description: "Ingreso de productos al inventario", Direction: "in", IsActive: true},
		{Model: gorm.Model{ID: 2}, Code: "outbound", Name: "Salida de mercancía", Description: "Salida de productos del inventario", Direction: "out", IsActive: true},
		{Model: gorm.Model{ID: 3}, Code: "adjustment", Name: "Ajuste de inventario", Description: "Corrección manual de cantidades de inventario", Direction: "neutral", IsActive: true},
		{Model: gorm.Model{ID: 4}, Code: "transfer", Name: "Transferencia entre bodegas", Description: "Movimiento de productos entre bodegas", Direction: "neutral", IsActive: true},
		{Model: gorm.Model{ID: 5}, Code: "return", Name: "Devolución", Description: "Reingreso de productos por devolución", Direction: "in", IsActive: true},
		{Model: gorm.Model{ID: 6}, Code: "sync", Name: "Sincronización desde canal", Description: "Actualización de inventario desde canal de venta", Direction: "neutral", IsActive: true},
		{Model: gorm.Model{ID: 7}, Code: "reserve", Name: "Reserva de stock", Description: "Stock reservado por orden pendiente", Direction: "neutral", IsActive: true},
		{Model: gorm.Model{ID: 8}, Code: "confirm_sale", Name: "Confirmación de venta", Description: "Stock confirmado como vendido", Direction: "out", IsActive: true},
		{Model: gorm.Model{ID: 9}, Code: "release", Name: "Liberación de reserva", Description: "Stock reservado liberado por cancelación", Direction: "neutral", IsActive: true},
	}

	for _, t := range types {
		var existing models.StockMovementType
		err := db.Where("id = ?", t.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&t).Error; err != nil {
				return fmt.Errorf("failed to create stock_movement_type %s: %w", t.Code, err)
			}
		}
	}

	return nil
}

// cleanOrphanWalletRows deletes wallet rows with business_id that don't exist in the business table.
// This prevents FK constraint creation from failing during AutoMigrate.
func (r *Repository) cleanOrphanWalletRows(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// Check if wallet table exists first
	var count int64
	if err := db.Raw(`
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'wallet'
	`).Scan(&count).Error; err != nil || count == 0 {
		return nil // Table doesn't exist yet, nothing to clean
	}

	// Delete transactions for orphan wallets first (FK dependency)
	if err := db.Exec(`
		DELETE FROM transaction
		WHERE wallet_id IN (
			SELECT w.id FROM wallet w
			LEFT JOIN business b ON w.business_id = b.id
			WHERE b.id IS NULL
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to delete orphan wallet transactions: %w", err)
	}

	// Delete orphan wallet rows
	if err := db.Exec(`
		DELETE FROM wallet w
		WHERE NOT EXISTS (
			SELECT 1 FROM business b WHERE b.id = w.business_id
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to delete orphan wallets: %w", err)
	}

	return nil
}

// fixNotificationConfigIndex reemplaza el índice único viejo idx_business_event_type
// (business_id, event_type) por uno parcial que usa las columnas nuevas y excluye soft-deleted.
func (r *Repository) fixNotificationConfigIndex(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// 1. Eliminar el índice viejo (sobre campo deprecated event_type)
	if err := db.Exec(`DROP INDEX IF EXISTS idx_business_event_type`).Error; err != nil {
		return fmt.Errorf("failed to drop old idx_business_event_type: %w", err)
	}

	// 2. Crear nuevo índice único parcial sobre las columnas correctas.
	// Una config es única por: business + integration + notification_type + event_type (solo activos).
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_notification_config_unique
		ON business_notification_configs (business_id, integration_id, notification_type_id, notification_event_type_id)
		WHERE deleted_at IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to create idx_notification_config_unique: %w", err)
	}

	return nil
}

// seedAllowedOrderStatusesByEventType inserta los estados de orden permitidos por tipo de evento
// Si vacío → significa "todos los estados permitidos"
func (r *Repository) seedAllowedOrderStatusesByEventType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// Mapeo: eventTypeID → []orderStatusID
	// IDs de order_statuses: pending=1, processing=2, shipped=3, delivered=4, completed=5, cancelled=6, refunded=7
	allowedMap := map[uint][]uint{
		1: {1, 2}, // SSE order.created → pending, processing
		2: {},     // SSE order.status_changed → todos (vacío)
		3: {1, 2}, // WA order.created → pending, processing
		4: {3, 4}, // WA order.shipped → shipped, delivered
		5: {4, 5}, // WA order.delivered → delivered, completed
		6: {6, 7}, // WA order.canceled → cancelled, refunded
		7: {},     // WA invoice.created → todos (vacío)
	}

	for eventTypeID, statusIDs := range allowedMap {
		if len(statusIDs) == 0 {
			continue // vacío = todos permitidos, no insertar nada
		}

		// Verificar que el event type existe
		var eventType models.NotificationEventType
		if err := db.Where("id = ?", eventTypeID).First(&eventType).Error; err != nil {
			continue // No existe, saltar
		}

		// Insertar solo si no existen ya
		for _, statusID := range statusIDs {
			if err := db.Exec(`
				INSERT INTO notification_event_type_allowed_statuses (notification_event_type_id, order_status_id)
				VALUES (?, ?)
				ON CONFLICT DO NOTHING
			`, eventTypeID, statusID).Error; err != nil {
				return fmt.Errorf("failed to seed allowed status %d for event type %d: %w", statusID, eventTypeID, err)
			}
		}
	}

	return nil
}

// seedEmailIntegrationType crea el tipo de integración para Email (notificaciones por correo)
func (r *Repository) seedEmailIntegrationType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	integrationType := models.IntegrationType{
		Model:         gorm.Model{ID: 29},
		Name:          "Email",
		Code:          "email",
		CategoryID:    ptrUint(3), // messaging category
		IsActive:      true,
		InDevelopment: false,
		Description:   "Notificaciones por correo electrónico via Amazon SES",
		Icon:          "mail",
	}

	var existing models.IntegrationType
	err := db.Where("id = ?", 29).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		if err := db.Create(&integrationType).Error; err != nil {
			return fmt.Errorf("failed to create email integration type: %w", err)
		}
	}

	return nil
}

// migrateOrderItemsFromJSONB backfills order_items from the legacy orders.items JSONB column.
// It is idempotent: only processes orders that have no rows in order_items yet.
// Must run BEFORE the DROP COLUMN for orders.items.
func (r *Repository) migrateOrderItemsFromJSONB(ctx context.Context) error {
	db := r.db.Conn(ctx)

	type orderRow struct {
		ID         string
		Items      []byte
		Currency   string
		BusinessID *uint
	}
	type legacyItem struct {
		ID       string  `json:"id"`
		SKU      string  `json:"sku"`
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}

	// Check if orders.items column still exists (idempotency)
	var colExists int64
	db.Raw(`
		SELECT COUNT(*) FROM information_schema.columns
		WHERE table_schema = CURRENT_SCHEMA()
		  AND table_name = 'orders'
		  AND column_name = 'items'
	`).Scan(&colExists)
	if colExists == 0 {
		// Column already dropped — migration was previously completed
		return nil
	}

	var rows []orderRow
	err := db.Raw(`
		SELECT id, items, currency, business_id
		FROM orders
		WHERE items IS NOT NULL
		  AND items::text NOT IN ('null', '[]', '')
		  AND deleted_at IS NULL
		  AND id NOT IN (
		      SELECT DISTINCT order_id FROM order_items WHERE deleted_at IS NULL
		  )
	`).Scan(&rows).Error
	if err != nil {
		return fmt.Errorf("failed to query orders for JSONB migration: %w", err)
	}

	if len(rows) == 0 {
		return nil
	}

	for _, row := range rows {
		var items []legacyItem
		if err := json.Unmarshal(row.Items, &items); err != nil {
			// Skip orders with unparseable items
			continue
		}

		for _, item := range items {
			quantity := item.Quantity
			if quantity <= 0 {
				quantity = 1
			}

			// Try to find product by SKU and business_id
			var productID *string
			if item.SKU != "" && row.BusinessID != nil {
				var product models.Product
				if err := db.Where("sku = ? AND business_id = ? AND deleted_at IS NULL", item.SKU, *row.BusinessID).First(&product).Error; err == nil {
					productID = &product.ID
				}
			}

			// Build metadata with original item data
			originalJSON, _ := json.Marshal(item)
			metadataJSON, _ := json.Marshal(map[string]interface{}{
				"_source":  "json_migration",
				"original": json.RawMessage(originalJSON),
			})

			currency := row.Currency
			if currency == "" {
				currency = "COP"
			}

			unitPrice := item.Price
			totalPrice := unitPrice * float64(quantity)

			orderItem := models.OrderItem{
				OrderID:               row.ID,
				ProductID:             productID,
				Quantity:              quantity,
				UnitPrice:             unitPrice,
				TotalPrice:            totalPrice,
				Currency:              currency,
				UnitPricePresentment:  unitPrice,
				TotalPricePresentment: totalPrice,
				Metadata:              datatypes.JSON(metadataJSON),
			}

			if err := db.Create(&orderItem).Error; err != nil {
				return fmt.Errorf("failed to create order_item for order %s: %w", row.ID, err)
			}
		}
	}

	return nil
}

// fixClientDniIndex corrige el índice único de DNI en la tabla client.
// El índice original era UNIQUE(dni) global, lo cual impedía que dos negocios
// tuvieran clientes con el mismo DNI. Se cambia a UNIQUE(business_id, dni)
// con filtro parcial para soportar multi-tenant correctamente.
func (r *Repository) fixClientDniIndex(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// 1. Eliminar el índice viejo global UNIQUE(dni)
	if err := db.Exec(`DROP INDEX IF EXISTS idx_business_client_dni`).Error; err != nil {
		return fmt.Errorf("failed to drop old idx_business_client_dni: %w", err)
	}

	// 2. Crear nuevo índice compuesto UNIQUE(business_id, dni) con filtro parcial
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_business_client_dni
		ON client (business_id, dni)
		WHERE dni IS NOT NULL AND deleted_at IS NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to create new idx_business_client_dni: %w", err)
	}

	return nil
}

// seedNewResourcesAndPermissions creates missing resources, their CRUD permissions,
// assigns them to the Administrador role, and adds them to business_resource_configured.
func (r *Repository) seedNewResourcesAndPermissions(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// =============================================
	// 1. New resources
	// =============================================
	newResources := []models.Resource{
		{Model: gorm.Model{ID: 19}, Name: "Clientes", Description: "Gestion de clientes"},
		{Model: gorm.Model{ID: 20}, Name: "Ultima Milla", Description: "Gestion de ultima milla (delivery)"},
		{Model: gorm.Model{ID: 21}, Name: "Billetera", Description: "Gestion de billetera y transacciones"},
		{Model: gorm.Model{ID: 22}, Name: "Inventario", Description: "Gestion de inventario"},
		{Model: gorm.Model{ID: 23}, Name: "Bodegas", Description: "Gestion de bodegas y ubicaciones"},
	}

	for _, res := range newResources {
		var count int64
		db.Model(&models.Resource{}).Where("id = ?", res.ID).Count(&count)
		if count == 0 {
			if err := db.Create(&res).Error; err != nil {
				return fmt.Errorf("failed to create resource %s: %w", res.Name, err)
			}
		}
	}

	// =============================================
	// 2. CRUD permissions for new resources + missing ones for existing resources
	// =============================================
	type permDef struct {
		ID         uint
		Name       string
		ResourceID uint
		ActionID   uint
	}

	businessTypeID := ptrUint(1)
	scopeID := uint(2) // Business scope

	allPerms := []permDef{
		// Clientes (resource 19)
		{54, "Create Clientes", 19, 1},
		{55, "Read Clientes", 19, 2},
		{56, "Update Clientes", 19, 3},
		{57, "Delete Clientes", 19, 4},
		// Ultima Milla (resource 20)
		{58, "Create Ultima Milla", 20, 1},
		{59, "Read Ultima Milla", 20, 2},
		{60, "Update Ultima Milla", 20, 3},
		{61, "Delete Ultima Milla", 20, 4},
		// Billetera (resource 21)
		{62, "Create Billetera", 21, 1},
		{63, "Read Billetera", 21, 2},
		{64, "Update Billetera", 21, 3},
		{65, "Delete Billetera", 21, 4},
		// Inventario (resource 22)
		{66, "Create Inventario", 22, 1},
		{67, "Read Inventario", 22, 2},
		{68, "Update Inventario", 22, 3},
		{69, "Delete Inventario", 22, 4},
		// Bodegas (resource 23)
		{70, "Create Bodegas", 23, 1},
		{71, "Read Bodegas", 23, 2},
		{72, "Update Bodegas", 23, 3},
		{73, "Delete Bodegas", 23, 4},
		// Missing CRUD for Permisos (resource 2)
		{74, "Create Permisos", 2, 1},
		{75, "Read Permisos", 2, 2},
		{76, "Update Permisos", 2, 3},
		{77, "Delete Permisos", 2, 4},
		// Missing CRUD for Roles (resource 3)
		{78, "Create Roles", 3, 1},
		{79, "Read Roles", 3, 2},
		{80, "Update Roles", 3, 3},
		{81, "Delete Roles", 3, 4},
		// Missing CRUD for Recursos (resource 4)
		{82, "Create Recursos", 4, 1},
		{83, "Read Recursos", 4, 2},
		{84, "Update Recursos", 4, 3},
		{85, "Delete Recursos", 4, 4},
	}

	newPermIDs := []uint{}
	for _, p := range allPerms {
		var count int64
		db.Model(&models.Permission{}).Where("id = ?", p.ID).Count(&count)
		if count == 0 {
			perm := models.Permission{
				Model:          gorm.Model{ID: p.ID},
				Name:           p.Name,
				Description:    p.Name,
				ResourceID:     p.ResourceID,
				ActionID:       p.ActionID,
				ScopeID:        scopeID,
				BusinessTypeID: businessTypeID,
			}
			if err := db.Create(&perm).Error; err != nil {
				return fmt.Errorf("failed to create permission %s: %w", p.Name, err)
			}
			newPermIDs = append(newPermIDs, p.ID)
		}
	}

	// =============================================
	// 3. Assign new permissions to Administrador role (ID 4)
	// =============================================
	adminRoleID := uint(4)
	for _, permID := range newPermIDs {
		var count int64
		db.Table("role_permissions").
			Where("role_id = ? AND permission_id = ?", adminRoleID, permID).
			Count(&count)
		if count == 0 {
			if err := db.Exec(
				"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
				adminRoleID, permID,
			).Error; err != nil {
				return fmt.Errorf("failed to assign permission %d to admin role: %w", permID, err)
			}
		}
	}

	// =============================================
	// 4. Add new resources + Notificaciones to business_resource_configured
	// =============================================
	resourceIDs := []uint{18, 19, 20, 21, 22, 23} // 18=Notificaciones + new ones

	var businessIDs []uint
	db.Table("business_resource_configured").
		Select("DISTINCT business_id").
		Where("deleted_at IS NULL").
		Scan(&businessIDs)

	for _, bizID := range businessIDs {
		for _, resID := range resourceIDs {
			var count int64
			db.Model(&models.BusinessResourceConfigured{}).
				Where("business_id = ? AND resource_id = ? AND deleted_at IS NULL", bizID, resID).
				Count(&count)
			if count == 0 {
				brc := models.BusinessResourceConfigured{
					BusinessID: bizID,
					ResourceID: resID,
					Active:     true,
				}
				if err := db.Create(&brc).Error; err != nil {
					return fmt.Errorf("failed to create business_resource_configured biz=%d res=%d: %w", bizID, resID, err)
				}
			}
		}
	}

	// Also assign Notificaciones permissions (50-53) to admin role if missing
	for _, permID := range []uint{50, 51, 52, 53} {
		var count int64
		db.Table("role_permissions").
			Where("role_id = ? AND permission_id = ?", adminRoleID, permID).
			Count(&count)
		if count == 0 {
			db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", adminRoleID, permID)
		}
	}

	return nil
}

// seedStorefrontRoleAndPermissions creates the cliente_final role, storefront resource and permissions
func (r *Repository) seedStorefrontRoleAndPermissions(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// 1. Create storefront resource (id=24)
	var resCount int64
	db.Model(&models.Resource{}).Where("id = ?", 24).Count(&resCount)
	if resCount == 0 {
		res := models.Resource{
			Model:       gorm.Model{ID: 24},
			Name:        "Storefront",
			Description: "Modulo de tienda para clientes finales",
		}
		if err := db.Create(&res).Error; err != nil {
			return fmt.Errorf("failed to create storefront resource: %w", err)
		}
	}

	// 2. Create storefront permissions (read=86, create=87)
	scopeID := uint(2) // Business scope
	storefrontPerms := []struct {
		ID       uint
		Name     string
		ActionID uint // 1=create, 2=read
	}{
		{86, "Read Storefront", 2},
		{87, "Create Storefront", 1},
	}

	permIDs := []uint{}
	for _, p := range storefrontPerms {
		var count int64
		db.Model(&models.Permission{}).Where("id = ?", p.ID).Count(&count)
		if count == 0 {
			perm := models.Permission{
				Model:       gorm.Model{ID: p.ID},
				Name:        p.Name,
				Description: p.Name,
				ResourceID:  24,
				ActionID:    p.ActionID,
				ScopeID:     scopeID,
			}
			if err := db.Create(&perm).Error; err != nil {
				return fmt.Errorf("failed to create permission %s: %w", p.Name, err)
			}
			permIDs = append(permIDs, p.ID)
		}
	}

	// 3. Create cliente_final role (level=5, scope=business, is_system=true)
	var roleCount int64
	db.Model(&models.Role{}).Where("name = ?", "cliente_final").Count(&roleCount)
	if roleCount == 0 {
		role := models.Role{
			Name:        "cliente_final",
			Description: "Rol para clientes finales del storefront",
			Level:       5,
			IsSystem:    true,
			ScopeID:     scopeID,
		}
		if err := db.Create(&role).Error; err != nil {
			return fmt.Errorf("failed to create cliente_final role: %w", err)
		}

		// Assign storefront permissions to the new role
		for _, permID := range []uint{86, 87} {
			if err := db.Exec(
				"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
				role.ID, permID,
			).Error; err != nil {
				return fmt.Errorf("failed to assign permission %d to cliente_final role: %w", permID, err)
			}
		}
	}

	// 4. Add storefront resource to all businesses' configured resources
	var businessIDs []uint
	db.Table("business_resource_configured").
		Select("DISTINCT business_id").
		Where("deleted_at IS NULL").
		Scan(&businessIDs)

	for _, bizID := range businessIDs {
		var count int64
		db.Model(&models.BusinessResourceConfigured{}).
			Where("business_id = ? AND resource_id = ? AND deleted_at IS NULL", bizID, 24).
			Count(&count)
		if count == 0 {
			brc := models.BusinessResourceConfigured{
				BusinessID: bizID,
				ResourceID: 24,
				Active:     true,
			}
			if err := db.Create(&brc).Error; err != nil {
				return fmt.Errorf("failed to create business_resource_configured biz=%d res=24: %w", bizID, err)
			}
		}
	}

	return nil
}

// activateAllProducts sets is_active=true and status='active' for all existing products
func (r *Repository) activateAllProducts(ctx context.Context) error {
	return r.db.Conn(ctx).Exec(
		"UPDATE products SET is_active = true, status = 'active' WHERE is_active = false OR status != 'active'",
	).Error
}

// seedStorefrontIntegrationCategory crea la categoría "Tu Tienda" (id=7) para integraciones de storefront
func (r *Repository) seedStorefrontIntegrationCategory(ctx context.Context) error {
	db := r.db.Conn(ctx)

	category := models.IntegrationCategory{
		Model:        gorm.Model{ID: 7},
		Code:         "storefront",
		Name:         "Tu Tienda",
		Description:  "Tu tienda propia y sitio web público",
		Icon:         "storefront",
		Color:        "#059669",
		DisplayOrder: 7,
		IsActive:     true,
		IsVisible:    true,
	}

	var existing models.IntegrationCategory
	err := db.Where("id = ?", 7).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		if err := db.Create(&category).Error; err != nil {
			return fmt.Errorf("failed to create storefront integration category: %w", err)
		}
	}

	return nil
}

// seedTiendaIntegrationType crea el tipo de integración "Tienda" (id=30)
func (r *Repository) seedTiendaIntegrationType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	integrationType := models.IntegrationType{
		Model:       gorm.Model{ID: 30},
		Name:        "Tienda",
		Code:        "tienda",
		CategoryID:  ptrUint(7),
		IsActive:    true,
		Description: "Tienda con login para clientes de tu negocio",
		Icon:        "shopping-bag",
	}

	var existing models.IntegrationType
	err := db.Where("id = ?", 30).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		if err := db.Create(&integrationType).Error; err != nil {
			return fmt.Errorf("failed to create tienda integration type: %w", err)
		}
	}

	return nil
}

// seedTiendaWebIntegrationType crea el tipo de integración "Tienda Web" (id=31)
func (r *Repository) seedTiendaWebIntegrationType(ctx context.Context) error {
	db := r.db.Conn(ctx)

	integrationType := models.IntegrationType{
		Model:       gorm.Model{ID: 31},
		Name:        "Sitio Web",
		Code:        "tienda_web",
		CategoryID:  ptrUint(7),
		IsActive:    true,
		Description: "Sitio web público para tu negocio",
		Icon:        "globe-alt",
	}

	var existing models.IntegrationType
	err := db.Where("id = ?", 31).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		if err := db.Create(&integrationType).Error; err != nil {
			return fmt.Errorf("failed to create tienda web integration type: %w", err)
		}
	}

	return nil
}

// ptrUint es un helper para crear punteros a uint
func ptrUint(v uint) *uint {
	return &v
}
