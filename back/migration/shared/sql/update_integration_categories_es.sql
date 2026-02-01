-- Actualizar categorías de integración con nombres y descripciones en español
-- Este script configura los datos de las categorías de integraciones

-- Categoría: E-commerce
UPDATE integration_categories
SET
    name = 'E-commerce',
    description = 'Plataformas de comercio electrónico como Shopify, WooCommerce, Mercado Libre, etc.',
    icon = 'ShoppingCartIcon',
    color = '#3b82f6',
    display_order = 1
WHERE code = 'ecommerce';

-- Categoría: Facturación
UPDATE integration_categories
SET
    name = 'Facturación',
    description = 'Sistemas de facturación electrónica y contabilidad',
    icon = 'DocumentTextIcon',
    color = '#10b981',
    display_order = 2
WHERE code = 'invoicing';

-- Categoría: Mensajería
UPDATE integration_categories
SET
    name = 'Mensajería',
    description = 'Plataformas de mensajería y comunicación con clientes (WhatsApp, Telegram, etc.)',
    icon = 'ChatBubbleLeftRightIcon',
    color = '#8b5cf6',
    display_order = 3
WHERE code = 'messaging';

-- Categoría: Pagos
UPDATE integration_categories
SET
    name = 'Pagos',
    description = 'Pasarelas de pago y procesadores de transacciones',
    icon = 'BanknotesIcon',
    color = '#f59e0b',
    display_order = 4
WHERE code = 'payment';

-- Categoría: Envíos
UPDATE integration_categories
SET
    name = 'Envíos',
    description = 'Servicios de logística y transporte de mercancía',
    icon = 'TruckIcon',
    color = '#ef4444',
    display_order = 5
WHERE code = 'shipping';

-- Categoría: Sistema
UPDATE integration_categories
SET
    name = 'Sistema',
    description = 'Integraciones internas y servicios del sistema',
    icon = 'CogIcon',
    color = '#6b7280',
    display_order = 6
WHERE code = 'system';

-- Insertar categorías si no existen (para nueva instalación)
INSERT INTO integration_categories (code, name, description, icon, color, display_order, is_active, is_visible, created_at, updated_at)
SELECT 'ecommerce', 'E-commerce', 'Plataformas de comercio electrónico como Shopify, WooCommerce, Mercado Libre, etc.', 'ShoppingCartIcon', '#3b82f6', 1, true, true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM integration_categories WHERE code = 'ecommerce');

INSERT INTO integration_categories (code, name, description, icon, color, display_order, is_active, is_visible, created_at, updated_at)
SELECT 'invoicing', 'Facturación', 'Sistemas de facturación electrónica y contabilidad', 'DocumentTextIcon', '#10b981', 2, true, true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM integration_categories WHERE code = 'invoicing');

INSERT INTO integration_categories (code, name, description, icon, color, display_order, is_active, is_visible, created_at, updated_at)
SELECT 'messaging', 'Mensajería', 'Plataformas de mensajería y comunicación con clientes (WhatsApp, Telegram, etc.)', 'ChatBubbleLeftRightIcon', '#8b5cf6', 3, true, true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM integration_categories WHERE code = 'messaging');

INSERT INTO integration_categories (code, name, description, icon, color, display_order, is_active, is_visible, created_at, updated_at)
SELECT 'payment', 'Pagos', 'Pasarelas de pago y procesadores de transacciones', 'BanknotesIcon', '#f59e0b', 4, true, true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM integration_categories WHERE code = 'payment');

INSERT INTO integration_categories (code, name, description, icon, color, display_order, is_active, is_visible, created_at, updated_at)
SELECT 'shipping', 'Envíos', 'Servicios de logística y transporte de mercancía', 'TruckIcon', '#ef4444', 5, true, true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM integration_categories WHERE code = 'shipping');

INSERT INTO integration_categories (code, name, description, icon, color, display_order, is_active, is_visible, created_at, updated_at)
SELECT 'system', 'Sistema', 'Integraciones internas y servicios del sistema', 'CogIcon', '#6b7280', 6, true, true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM integration_categories WHERE code = 'system');
