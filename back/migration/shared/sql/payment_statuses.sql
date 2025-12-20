-- ============================================
-- Tabla de Estados de Pago de Probability
-- ============================================
-- Este archivo crea la tabla payment_statuses e inserta los estados iniciales
-- del sistema Probability.
--
-- Uso: Ejecutar este script en PostgreSQL para crear la tabla y poblar
-- con los estados básicos del sistema.
-- ============================================

-- Crear tabla payment_statuses
CREATE TABLE IF NOT EXISTS payment_statuses (
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
CREATE INDEX IF NOT EXISTS idx_payment_statuses_code ON payment_statuses(code);
CREATE INDEX IF NOT EXISTS idx_payment_statuses_category ON payment_statuses(category);
CREATE INDEX IF NOT EXISTS idx_payment_statuses_is_active ON payment_statuses(is_active);
CREATE INDEX IF NOT EXISTS idx_payment_statuses_deleted_at ON payment_statuses(deleted_at);

-- Insertar estados iniciales de Probability
INSERT INTO payment_statuses (code, name, description, category, is_active, created_at, updated_at)
VALUES 
    ('pending', 'Pendiente', 'Pago pendiente de procesamiento', 'pending', true, NOW(), NOW()),
    ('authorized', 'Autorizado', 'Pago autorizado pero aún no capturado', 'pending', true, NOW(), NOW()),
    ('paid', 'Pagado', 'Pago completado exitosamente', 'completed', true, NOW(), NOW()),
    ('partially_paid', 'Parcialmente Pagado', 'Solo se ha pagado una parte del total', 'pending', true, NOW(), NOW()),
    ('refunded', 'Reembolsado', 'Pago reembolsado completamente', 'refunded', true, NOW(), NOW()),
    ('partially_refunded', 'Parcialmente Reembolsado', 'Solo se ha reembolsado una parte', 'refunded', true, NOW(), NOW()),
    ('voided', 'Anulado', 'Pago anulado antes de completarse', 'failed', true, NOW(), NOW()),
    ('failed', 'Fallido', 'Pago fallido durante el procesamiento', 'failed', true, NOW(), NOW()),
    ('unpaid', 'No Pagado', 'Orden no pagada', 'pending', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- Notas:
-- ============================================
-- 1. Los códigos deben coincidir con los estados de las plataformas externas
--    (ej: Shopify financial_status)
--
-- 2. Las categorías ayudan a agrupar estados:
--    - pending: Estados pendientes de procesamiento
--    - completed: Estados completados exitosamente
--    - refunded: Estados relacionados con reembolsos
--    - failed: Estados de fallo o anulación
--
-- 3. Los campos icon, color y metadata pueden ser utilizados por el frontend
--    para personalizar la visualización de los estados
-- ============================================
