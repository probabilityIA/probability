-- =====================================================
-- Migración: Agregar campos document_json y document_metadata a invoices
-- Propósito: Almacenar JSON completo del documento retornado por Softpymes
-- Fecha: 2026-02-09
-- =====================================================

-- PROBLEMA:
-- La API de Softpymes NO retorna URLs de PDF/XML en ningún endpoint.
-- Solo retorna la estructura del documento con metadata (cliente, sucursal, totales, etc.)

-- SOLUCIÓN:
-- Almacenar el JSON completo del documento de Softpymes para:
-- 1. Tener información detallada del documento
-- 2. Mostrar metadata en el frontend (cliente, sucursal, totales)
-- 3. Tener un registro completo para auditoría
-- 4. Permitir reconstrucción de información si se necesita

BEGIN;

-- Agregar columnas JSONB para almacenar documento completo
ALTER TABLE invoices
ADD COLUMN document_json JSONB,
ADD COLUMN document_metadata JSONB;

-- Comentarios descriptivos
COMMENT ON COLUMN invoices.document_json IS
'JSON completo del documento retornado por Softpymes endpoint /search/documents/.
Contiene: branchCode, branchName, customerIdentification, customerName, documentDate,
documentNumber, total, totalIva, seller, details, etc.
Se obtiene después de crear la factura mediante consulta con documentNumber.';

COMMENT ON COLUMN invoices.document_metadata IS
'Metadata adicional extraída del document_json para acceso rápido.
Contiene: branch_code, branch_name, customer_id, customer_name, seller_nit,
seller_name, document_name, document_date, total, total_iva.
Permite consultas rápidas sin parsear el JSON completo.';

-- Índice para búsquedas rápidas por número de documento en el JSON
CREATE INDEX idx_invoices_document_json_number
ON invoices ((document_json->>'documentNumber'));

-- Índice para búsquedas por identificación del cliente
CREATE INDEX idx_invoices_document_json_customer_id
ON invoices ((document_json->>'customerIdentification'));

COMMIT;
