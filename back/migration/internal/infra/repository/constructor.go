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
	if err := r.db.Conn(ctx).AutoMigrate(
		&models.BusinessType{},
		&models.Scope{},
		&models.Business{},
		&models.BusinessResourceConfigured{},
		&models.Resource{},
		&models.Role{},
		&models.Permission{},
		&models.User{},
		&models.BusinessStaff{},
		&models.Client{},
		&models.Action{},
		&models.APIKey{},
		&models.IntegrationCategory{},
		&models.IntegrationType{},
		&models.Integration{},

		// Integration Notification Configs (debe ir después de Integration)
		&models.IntegrationNotificationConfig{},

		// Payment Methods
		&models.PaymentMethod{},
		&models.PaymentMethodMapping{},
		&models.OrderStatusMapping{},
		&models.OrderStatus{},
		&models.PaymentStatus{},
		&models.FulfillmentStatus{},

		// Notification Types (nueva arquitectura)
		&models.NotificationType{},
		&models.NotificationEventType{},

		// Business Notification Config (debe ir después de Integration, NotificationType, NotificationEventType)
		&models.BusinessNotificationConfig{},

		// Business Notification Config Order Status (tabla intermedia)
		// Debe ir después de BusinessNotificationConfig y OrderStatus
		&models.BusinessNotificationConfigOrderStatus{},

		&models.Product{},

		// Orders
		&models.Order{},
		&models.OrderHistory{},
		&models.OrderError{},

		// Order Channel Metadata
		&models.OrderChannelMetadata{},

		// Order Items
		&models.OrderItem{},

		// Addresses
		&models.Address{},

		// Payments
		&models.Payment{},

		// Shipments
		&models.Shipment{},

		// WhatsApp Integration
		&models.WhatsAppConversation{},
		&models.WhatsAppMessageLog{},

		// Invoicing System
		&models.InvoicingProviderType{},
		&models.InvoicingProvider{},
		&models.InvoicingConfig{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.InvoiceSyncLog{},
		&models.CreditNote{},

		// Bulk Invoice Jobs (Async Bulk Invoicing)
		&models.BulkInvoiceJob{},
		&models.BulkInvoiceJobItem{},

		// Origin Addresses
		&models.OriginAddress{},

		// Payment Transactions (pasarelas de pago externas)
		&models.PaymentTransaction{},
		&models.PaymentSyncLog{},
	); err != nil {
		return err
	}

	// Insertar datos iniciales de integration_categories
	if err := r.seedIntegrationCategories(ctx); err != nil {
		return fmt.Errorf("failed to seed integration categories: %w", err)
	}

	// Migrar integration_types existentes a las nuevas categorías
	if err := r.migrateIntegrationTypesToCategories(ctx); err != nil {
		return fmt.Errorf("failed to migrate integration types to categories: %w", err)
	}

	// Crear tipo de integración para Softpymes (facturación)
	if err := r.seedSoftpymesIntegrationType(ctx); err != nil {
		return fmt.Errorf("failed to seed softpymes integration type: %w", err)
	}

	// Crear tipo de integración para Plataforma (órdenes manuales)
	if err := r.seedPlatformIntegrationType(ctx); err != nil {
		return fmt.Errorf("failed to seed platform integration type: %w", err)
	}

	// Crear integraciones de plataforma para negocios existentes
	if err := r.seedPlatformIntegrationsForBusinesses(ctx); err != nil {
		return fmt.Errorf("failed to seed platform integrations for businesses: %w", err)
	}

	// Migrar datos de invoicing_providers a integrations
	if err := r.migrateInvoicingProvidersToIntegrations(ctx); err != nil {
		return fmt.Errorf("failed to migrate invoicing providers: %w", err)
	}

	// Insertar datos iniciales de notification_types y notification_event_types
	if err := r.seedNotificationTypesAndEvents(ctx); err != nil {
		return fmt.Errorf("failed to seed notification types and events: %w", err)
	}

	// Migrar datos de business_notification_configs después de crear columnas
	if err := r.migrateBusinessNotificationConfigData(ctx); err != nil {
		return fmt.Errorf("failed to migrate business_notification_configs data: %w", err)
	}

	// Crear tipo de integración para EnvioClick (transporte)
	if err := r.seedEnvioClickIntegrationType(ctx); err != nil {
		return fmt.Errorf("failed to seed envioclick integration type: %w", err)
	}

	// Seed base_url de EnvioClick si no está configurada
	if err := r.seedEnvioClickBaseURL(ctx); err != nil {
		return fmt.Errorf("failed to seed envioclick base url: %w", err)
	}

	// Marcar integration types en desarrollo (IDs 13-21)
	if err := r.markInDevelopmentIntegrationTypes(ctx); err != nil {
		return fmt.Errorf("failed to mark in-development integration types: %w", err)
	}

	// Agregar columnas de modo testing
	if err := r.addIsTestingColumns(ctx); err != nil {
		return fmt.Errorf("failed to add is_testing columns: %w", err)
	}

	// Agregar columna de credenciales de plataforma a integration_types
	if err := r.addPlatformCredentialsToIntegrationTypes(ctx); err != nil {
		return fmt.Errorf("failed to add platform_credentials_encrypted to integration_types: %w", err)
	}

	return r.createDefaultUserIfNotExists(ctx)
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

// ptrUint es un helper para crear punteros a uint
func ptrUint(v uint) *uint {
	return &v
}
