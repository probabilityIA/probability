-- ============================================
-- Tabla de Estados de Órdenes de Probability
-- ============================================
-- Este archivo crea la tabla order_statuses e inserta los estados iniciales
-- del sistema Probability.
--
-- Uso: Ejecutar este script en PostgreSQL para crear la tabla y poblar
-- con los estados básicos del sistema.
-- ============================================

-- Crear tabla order_statuses
CREATE TABLE IF NOT EXISTS order_statuses (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    -- Identificación
    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    
    -- Categorización
    category VARCHAR(64),
    
    -- Configuración
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    -- UI/UX
    icon VARCHAR(255),
    color VARCHAR(32),
    metadata JSONB
);

-- Crear índices
CREATE INDEX IF NOT EXISTS idx_order_statuses_code ON order_statuses(code);
CREATE INDEX IF NOT EXISTS idx_order_statuses_category ON order_statuses(category);
CREATE INDEX IF NOT EXISTS idx_order_statuses_is_active ON order_statuses(is_active);
CREATE INDEX IF NOT EXISTS idx_order_statuses_deleted_at ON order_statuses(deleted_at);

-- Insertar estados iniciales de Probability
INSERT INTO order_statuses (code, name, description, category, is_active, created_at, updated_at)
VALUES 
    ('pending', 'Pendiente', 'Orden recibida, pendiente de procesamiento', 'active', true, NOW(), NOW()),
    ('processing', 'En Procesamiento', 'Orden en proceso de preparación', 'active', true, NOW(), NOW()),
    ('shipped', 'Enviada', 'Orden enviada al cliente', 'active', true, NOW(), NOW()),
    ('delivered', 'Entregada', 'Orden entregada al cliente', 'completed', true, NOW(), NOW()),
    ('completed', 'Completada', 'Orden completada exitosamente', 'completed', true, NOW(), NOW()),
    ('cancelled', 'Cancelada', 'Orden cancelada', 'cancelled', true, NOW(), NOW()),
    ('refunded', 'Reembolsada', 'Orden reembolsada al cliente', 'refunded', true, NOW(), NOW()),
    ('failed', 'Fallida', 'Orden fallida durante el procesamiento', 'cancelled', true, NOW(), NOW()),
    ('on_hold', 'En Espera', 'Orden en espera por alguna razón', 'active', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- Notas:
-- ============================================
-- 1. Los códigos deben coincidir con las constantes definidas en
--    orders/internal/domain/entities.go
--
-- 2. Las categorías ayudan a agrupar estados:
--    - active: Estados activos del proceso
--    - completed: Estados finales exitosos
--    - cancelled: Estados de cancelación
--    - refunded: Estados relacionados con reembolsos
--
-- 3. Los campos icon, color y metadata pueden ser utilizados por el frontend
--    para personalizar la visualización de los estados
-- ============================================
