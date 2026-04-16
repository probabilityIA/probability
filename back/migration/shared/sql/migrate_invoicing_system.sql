-- ==========================================
-- MIGRACIÓN: Sistema de Facturación
-- Autor: Claude Code
-- Fecha: 2026-01-31
-- Descripción: Crea las tablas necesarias para el módulo de facturación electrónica
-- ==========================================

-- NOTA: GORM ya crea las tablas con AutoMigrate, pero este archivo
-- sirve como documentación y backup de la estructura

-- ==========================================
-- 1. INVOICING PROVIDER TYPES
-- ==========================================
-- Tipos de proveedores de facturación disponibles
-- Ejemplo: Softpymes, Siigo, Facturama, Alegra
CREATE TABLE IF NOT EXISTS invoicing_provider_types (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(500),
    icon VARCHAR(100),
    image_url VARCHAR(500),

    api_base_url VARCHAR(255),
    documentation_url VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    supported_countries VARCHAR(500)
);

CREATE INDEX idx_invoicing_provider_types_deleted_at ON invoicing_provider_types(deleted_at);
CREATE INDEX idx_invoicing_provider_types_code ON invoicing_provider_types(code);
CREATE INDEX idx_invoicing_provider_types_is_active ON invoicing_provider_types(is_active);

-- ==========================================
-- 2. INVOICING PROVIDERS
-- ==========================================
-- Instancias configuradas de proveedores por negocio
CREATE TABLE IF NOT EXISTS invoicing_providers (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    business_id INTEGER NOT NULL REFERENCES businesses(id) ON UPDATE CASCADE ON DELETE CASCADE,
    provider_type_id INTEGER NOT NULL REFERENCES invoicing_provider_types(id) ON UPDATE CASCADE ON DELETE RESTRICT,

    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),

    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,

    config JSONB,
    credentials JSONB,

    created_by_id INTEGER NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    updated_by_id INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,

    UNIQUE (business_id, provider_type_id)
);

CREATE INDEX idx_invoicing_providers_deleted_at ON invoicing_providers(deleted_at);
CREATE INDEX idx_invoicing_providers_business_id ON invoicing_providers(business_id);
CREATE INDEX idx_invoicing_providers_provider_type_id ON invoicing_providers(provider_type_id);
CREATE INDEX idx_invoicing_providers_is_active ON invoicing_providers(is_active);
CREATE INDEX idx_invoicing_providers_is_default ON invoicing_providers(is_default);
CREATE INDEX idx_invoicing_providers_created_by_id ON invoicing_providers(created_by_id);
CREATE INDEX idx_invoicing_providers_updated_by_id ON invoicing_providers(updated_by_id);

-- ==========================================
-- 3. INVOICING CONFIGS
-- ==========================================
-- Configuración de facturación por integración
CREATE TABLE IF NOT EXISTS invoicing_configs (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    business_id INTEGER NOT NULL REFERENCES businesses(id) ON UPDATE CASCADE ON DELETE CASCADE,
    integration_id INTEGER NOT NULL REFERENCES integrations(id) ON UPDATE CASCADE ON DELETE CASCADE,
    invoicing_provider_id INTEGER NOT NULL REFERENCES invoicing_providers(id) ON UPDATE CASCADE ON DELETE CASCADE,

    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    auto_invoice BOOLEAN NOT NULL DEFAULT FALSE,

    filters JSONB,
    invoice_config JSONB,

    description VARCHAR(500),
    created_by_id INTEGER NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    updated_by_id INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,

    UNIQUE (business_id, integration_id)
);

CREATE INDEX idx_invoicing_configs_deleted_at ON invoicing_configs(deleted_at);
CREATE INDEX idx_invoicing_configs_business_id ON invoicing_configs(business_id);
CREATE INDEX idx_invoicing_configs_integration_id ON invoicing_configs(integration_id);
CREATE INDEX idx_invoicing_configs_invoicing_provider_id ON invoicing_configs(invoicing_provider_id);
CREATE INDEX idx_invoicing_configs_enabled ON invoicing_configs(enabled);
CREATE INDEX idx_invoicing_configs_auto_invoice ON invoicing_configs(auto_invoice);
CREATE INDEX idx_invoicing_configs_created_by_id ON invoicing_configs(created_by_id);
CREATE INDEX idx_invoicing_configs_updated_by_id ON invoicing_configs(updated_by_id);

-- ==========================================
-- 4. INVOICES
-- ==========================================
-- Facturas generadas
CREATE TABLE IF NOT EXISTS invoices (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    order_id VARCHAR(36) NOT NULL REFERENCES orders(id) ON UPDATE CASCADE ON DELETE CASCADE,
    business_id INTEGER NOT NULL REFERENCES businesses(id) ON UPDATE CASCADE ON DELETE CASCADE,
    invoicing_provider_id INTEGER NOT NULL REFERENCES invoicing_providers(id) ON UPDATE CASCADE ON DELETE RESTRICT,

    invoice_number VARCHAR(128) NOT NULL,
    external_id VARCHAR(255),
    internal_number VARCHAR(128) NOT NULL UNIQUE,

    subtotal DECIMAL(12,2) NOT NULL,
    tax DECIMAL(12,2) NOT NULL,
    discount DECIMAL(12,2) NOT NULL,
    shipping_cost DECIMAL(12,2) NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'COP',

    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255),
    customer_phone VARCHAR(32),
    customer_dni VARCHAR(64),

    status VARCHAR(64) NOT NULL DEFAULT 'pending',

    issued_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    expires_at TIMESTAMP,

    invoice_url VARCHAR(512),
    pdf_url VARCHAR(512),
    xml_url VARCHAR(512),
    cufe VARCHAR(255),

    notes TEXT,
    metadata JSONB,
    provider_response JSONB,

    UNIQUE (order_id, invoicing_provider_id)
);

CREATE INDEX idx_invoices_deleted_at ON invoices(deleted_at);
CREATE INDEX idx_invoices_order_id ON invoices(order_id);
CREATE INDEX idx_invoices_business_id ON invoices(business_id);
CREATE INDEX idx_invoices_invoicing_provider_id ON invoices(invoicing_provider_id);
CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX idx_invoices_external_id ON invoices(external_id);
CREATE INDEX idx_invoices_internal_number ON invoices(internal_number);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_issued_at ON invoices(issued_at);

-- ==========================================
-- 5. INVOICE ITEMS
-- ==========================================
-- Items/líneas de facturas
CREATE TABLE IF NOT EXISTS invoice_items (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    invoice_id INTEGER NOT NULL REFERENCES invoices(id) ON UPDATE CASCADE ON DELETE CASCADE,
    product_id VARCHAR(64) REFERENCES products(id) ON UPDATE CASCADE ON DELETE SET NULL,

    sku VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(12,2) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'COP',

    tax DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_rate DECIMAL(5,4),
    discount DECIMAL(12,2) NOT NULL DEFAULT 0,

    provider_item_id VARCHAR(255),
    metadata JSONB
);

CREATE INDEX idx_invoice_items_deleted_at ON invoice_items(deleted_at);
CREATE INDEX idx_invoice_items_invoice_id ON invoice_items(invoice_id);
CREATE INDEX idx_invoice_items_product_id ON invoice_items(product_id);

-- ==========================================
-- 6. INVOICE SYNC LOGS
-- ==========================================
-- Logs de sincronización con proveedores
CREATE TABLE IF NOT EXISTS invoice_sync_logs (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    invoice_id INTEGER NOT NULL REFERENCES invoices(id) ON UPDATE CASCADE ON DELETE CASCADE,

    operation_type VARCHAR(64) NOT NULL,
    status VARCHAR(64) NOT NULL,

    request_payload JSONB,
    request_headers JSONB,
    request_url VARCHAR(512),

    response_status INTEGER,
    response_body JSONB,
    response_headers JSONB,

    error_message TEXT,
    error_code VARCHAR(64),
    error_details JSONB,

    retry_count INTEGER NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMP,
    max_retries INTEGER NOT NULL DEFAULT 3,
    retried_at TIMESTAMP,

    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    duration INTEGER,

    triggered_by VARCHAR(64),
    user_id INTEGER REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE INDEX idx_invoice_sync_logs_deleted_at ON invoice_sync_logs(deleted_at);
CREATE INDEX idx_invoice_sync_logs_invoice_id ON invoice_sync_logs(invoice_id);
CREATE INDEX idx_invoice_sync_logs_operation_type ON invoice_sync_logs(operation_type);
CREATE INDEX idx_invoice_sync_logs_status ON invoice_sync_logs(status);
CREATE INDEX idx_invoice_sync_logs_retry_count ON invoice_sync_logs(retry_count);
CREATE INDEX idx_invoice_sync_logs_next_retry_at ON invoice_sync_logs(next_retry_at);
CREATE INDEX idx_invoice_sync_logs_started_at ON invoice_sync_logs(started_at);
CREATE INDEX idx_invoice_sync_logs_completed_at ON invoice_sync_logs(completed_at);
CREATE INDEX idx_invoice_sync_logs_response_status ON invoice_sync_logs(response_status);
CREATE INDEX idx_invoice_sync_logs_user_id ON invoice_sync_logs(user_id);

-- ==========================================
-- 7. CREDIT NOTES
-- ==========================================
-- Notas de crédito
CREATE TABLE IF NOT EXISTS credit_notes (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    invoice_id INTEGER NOT NULL REFERENCES invoices(id) ON UPDATE CASCADE ON DELETE CASCADE,
    business_id INTEGER NOT NULL REFERENCES businesses(id) ON UPDATE CASCADE ON DELETE CASCADE,

    credit_note_number VARCHAR(128) NOT NULL,
    external_id VARCHAR(255),
    internal_number VARCHAR(128) NOT NULL UNIQUE,

    note_type VARCHAR(64) NOT NULL,

    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'COP',

    reason VARCHAR(255) NOT NULL,
    description TEXT,

    status VARCHAR(64) NOT NULL DEFAULT 'pending',

    issued_at TIMESTAMP,
    cancelled_at TIMESTAMP,

    note_url VARCHAR(512),
    pdf_url VARCHAR(512),
    xml_url VARCHAR(512),
    cufe VARCHAR(255),

    metadata JSONB,
    provider_response JSONB,

    created_by_id INTEGER NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX idx_credit_notes_deleted_at ON credit_notes(deleted_at);
CREATE INDEX idx_credit_notes_invoice_id ON credit_notes(invoice_id);
CREATE INDEX idx_credit_notes_business_id ON credit_notes(business_id);
CREATE INDEX idx_credit_notes_credit_note_number ON credit_notes(credit_note_number);
CREATE INDEX idx_credit_notes_external_id ON credit_notes(external_id);
CREATE INDEX idx_credit_notes_internal_number ON credit_notes(internal_number);
CREATE INDEX idx_credit_notes_note_type ON credit_notes(note_type);
CREATE INDEX idx_credit_notes_status ON credit_notes(status);
CREATE INDEX idx_credit_notes_issued_at ON credit_notes(issued_at);
CREATE INDEX idx_credit_notes_created_by_id ON credit_notes(created_by_id);

-- ==========================================
-- COMENTARIOS SOBRE LA MIGRACIÓN
-- ==========================================
-- GORM AutoMigrate ya crea estas tablas automáticamente cuando se ejecuta el servicio.
-- Este archivo SQL es un backup y documentación de la estructura esperada.
--
-- Para ejecutar manualmente (solo si es necesario):
-- psql -h localhost -p 5433 -U postgres -d probability -f migrate_invoicing_system.sql
