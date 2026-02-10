-- =====================================================
-- Crear índices para document_json (ejecución manual)
-- Fecha: 2026-02-09
-- =====================================================
-- NOTA: Estos índices mejoran el rendimiento de búsquedas
-- pero NO son obligatorios para el funcionamiento básico

-- Índice para búsquedas por número de documento en el JSON
CREATE INDEX IF NOT EXISTS idx_invoices_document_json_number
ON invoices ((document_json->>'documentNumber'));

-- Índice para búsquedas por identificación del cliente
CREATE INDEX IF NOT EXISTS idx_invoices_document_json_customer_id
ON invoices ((document_json->>'customerIdentification'));

-- Verificar que se crearon
SELECT indexname, indexdef
FROM pg_indexes
WHERE tablename = 'invoices'
AND indexname LIKE '%document%';
