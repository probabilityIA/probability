-- =====================================================
-- Rollback: Eliminar campos document_json y document_metadata de invoices
-- Propósito: Revertir migración 028
-- Fecha: 2026-02-09
-- =====================================================

BEGIN;

-- Eliminar índices
DROP INDEX IF EXISTS idx_invoices_document_json_number;
DROP INDEX IF EXISTS idx_invoices_document_json_customer_id;

-- Eliminar columnas
ALTER TABLE invoices
DROP COLUMN IF EXISTS document_json,
DROP COLUMN IF EXISTS document_metadata;

COMMIT;
