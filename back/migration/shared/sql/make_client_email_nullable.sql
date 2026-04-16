-- ============================================
-- Migración: Hacer email nullable en tabla client
-- Fecha: 2026-03-01
-- Contexto: Permitir crear clientes sin email.
--           El unique index actual impide múltiples clientes sin email
--           porque Go guarda "" en vez de NULL.
-- ============================================

-- 1. Convertir emails vacíos a NULL
UPDATE client SET email = NULL WHERE email = '';

-- 2. Permitir NULL en la columna
ALTER TABLE client ALTER COLUMN email DROP NOT NULL;

-- 3. Recrear unique index como partial (ignora NULLs)
DROP INDEX IF EXISTS idx_business_client_email;
CREATE UNIQUE INDEX idx_business_client_email ON client(business_id, email) WHERE email IS NOT NULL AND deleted_at IS NULL;
