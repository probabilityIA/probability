-- ============================================
-- Tabla de Estados de Fulfillment de Probability
-- ============================================
-- Este archivo crea la tabla fulfillment_statuses e inserta los estados iniciales
-- del sistema Probability.
--
-- Uso: Ejecutar este script en PostgreSQL para crear la tabla y poblar
-- con los estados básicos del sistema.
-- ============================================

-- Crear tabla fulfillment_statuses
CREATE TABLE IF NOT EXISTS fulfillment_statuses (
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
CREATE INDEX IF NOT EXISTS idx_fulfillment_statuses_code ON fulfillment_statuses(code);
CREATE INDEX IF NOT EXISTS idx_fulfillment_statuses_category ON fulfillment_statuses(category);
CREATE INDEX IF NOT EXISTS idx_fulfillment_statuses_is_active ON fulfillment_statuses(is_active);
CREATE INDEX IF NOT EXISTS idx_fulfillment_statuses_deleted_at ON fulfillment_statuses(deleted_at);

-- Insertar estados iniciales de Probability
INSERT INTO fulfillment_statuses (code, name, description, category, is_active, created_at, updated_at)
VALUES 
    ('unfulfilled', 'No Cumplida', 'Orden aún no ha sido cumplida', 'pending', true, NOW(), NOW()),
    ('partial', 'Parcial', 'Solo parte de la orden ha sido cumplida', 'in_progress', true, NOW(), NOW()),
    ('fulfilled', 'Cumplida', 'Orden cumplida completamente', 'completed', true, NOW(), NOW()),
    ('shipped', 'Enviada', 'Orden enviada al cliente', 'in_progress', true, NOW(), NOW()),
    ('in_transit', 'En Tránsito', 'Orden en camino al cliente', 'in_progress', true, NOW(), NOW()),
    ('delivered', 'Entregada', 'Orden entregada al cliente', 'completed', true, NOW(), NOW()),
    ('cancelled', 'Cancelada', 'Fulfillment cancelado', 'cancelled', true, NOW(), NOW()),
    ('failed', 'Fallido', 'Fulfillment fallido', 'failed', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- Notas:
-- ============================================
-- 1. Los códigos deben coincidir con los estados de las plataformas externas
--    (ej: Shopify fulfillment_status)
--
-- 2. Las categorías ayudan a agrupar estados:
--    - pending: Estados pendientes de procesamiento
--    - in_progress: Estados en proceso
--    - completed: Estados completados exitosamente
--    - cancelled: Estados cancelados
--    - failed: Estados de fallo
--
-- 3. Los campos icon, color y metadata pueden ser utilizados por el frontend
--    para personalizar la visualización de los estados
-- ============================================
