package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
	"github.com/secamc93/probability/back/migration/shared/models"
	"golang.org/x/crypto/bcrypt"
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

// seedNotificationTypesAndEvents inserta los datos iniciales de notification_types y notification_event_types
func (r *Repository) seedNotificationTypesAndEvents(ctx context.Context) error {
	db := r.db.Conn(ctx)

	// 1. Insertar notification_types si no existen
	notificationTypes := []models.NotificationType{
		{
			Model:       gorm.Model{ID: 1},
			Name:        "SSE",
			Code:        "sse",
			Description: "Server-Sent Events para notificaciones en tiempo real",
			IsActive:    true,
		},
		{
			Model:       gorm.Model{ID: 2},
			Name:        "WhatsApp",
			Code:        "whatsapp",
			Description: "Mensajes de WhatsApp Business",
			IsActive:    true,
		},
		{
			Model:       gorm.Model{ID: 3},
			Name:        "Email",
			Code:        "email",
			Description: "Notificaciones por correo electrónico",
			IsActive:    true,
		},
		{
			Model:       gorm.Model{ID: 4},
			Name:        "SMS",
			Code:        "sms",
			Description: "Mensajes de texto SMS",
			IsActive:    false,
		},
	}

	for _, nt := range notificationTypes {
		var existing models.NotificationType
		err := db.Where("id = ?", nt.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// No existe, crear
			if err := db.Create(&nt).Error; err != nil {
				return fmt.Errorf("failed to create notification_type %s: %w", nt.Code, err)
			}
		}
	}

	// 2. Insertar notification_event_types si no existen
	notificationEventTypes := []models.NotificationEventType{
		// Eventos para SSE
		{
			Model:              gorm.Model{ID: 1},
			NotificationTypeID: 1,
			EventCode:          "order.created",
			EventName:          "Nueva Orden",
			IsActive:           true,
		},
		{
			Model:              gorm.Model{ID: 2},
			NotificationTypeID: 1,
			EventCode:          "order.status_changed",
			EventName:          "Cambio de Estado",
			IsActive:           true,
		},
		// Eventos para WhatsApp
		{
			Model:              gorm.Model{ID: 3},
			NotificationTypeID: 2,
			EventCode:          "order.created",
			EventName:          "Confirmación de Pedido",
			IsActive:           true,
		},
		{
			Model:              gorm.Model{ID: 4},
			NotificationTypeID: 2,
			EventCode:          "order.shipped",
			EventName:          "Pedido Enviado",
			IsActive:           true,
		},
		{
			Model:              gorm.Model{ID: 5},
			NotificationTypeID: 2,
			EventCode:          "order.delivered",
			EventName:          "Pedido Entregado",
			IsActive:           true,
		},
		{
			Model:              gorm.Model{ID: 6},
			NotificationTypeID: 2,
			EventCode:          "order.canceled",
			EventName:          "Pedido Cancelado",
			IsActive:           true,
		},
		{
			Model:              gorm.Model{ID: 7},
			NotificationTypeID: 2,
			EventCode:          "invoice.created",
			EventName:          "Factura Generada",
			IsActive:           true,
		},
	}

	for _, net := range notificationEventTypes {
		var existing models.NotificationEventType
		err := db.Where("id = ?", net.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// No existe, crear
			if err := db.Create(&net).Error; err != nil {
				return fmt.Errorf("failed to create notification_event_type %s: %w", net.EventCode, err)
			}
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
		{1, 1},  // pending
		{2, 2},  // processing
		{9, 3},  // on_hold
		{3, 4},  // shipped
		{4, 5},  // delivered
		{5, 6},  // completed
		{7, 7},  // refunded
		{6, 8},  // cancelled
		{8, 9},  // failed
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
		1: {1, 2},       // SSE order.created → pending, processing
		2: {},            // SSE order.status_changed → todos (vacío)
		3: {1, 2},       // WA order.created → pending, processing
		4: {3, 4},       // WA order.shipped → shipped, delivered
		5: {4, 5},       // WA order.delivered → delivered, completed
		6: {6, 7},       // WA order.canceled → cancelled, refunded
		7: {},            // WA invoice.created → todos (vacío)
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

// ptrUint es un helper para crear punteros a uint
func ptrUint(v uint) *uint {
	return &v
}
