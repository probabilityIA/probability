-- ═══════════════════════════════════════════════════════════════
-- Migration: Bulk Invoice Jobs
-- Description: Sistema de tracking para facturación masiva asíncrona
-- Author: Claude Code
-- Date: 2026-02-07
-- ═══════════════════════════════════════════════════════════════

-- Tabla principal de jobs de facturación masiva
CREATE TABLE IF NOT EXISTS bulk_invoice_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id INT NOT NULL,
    created_by_user_id INT,
    total_orders INT NOT NULL CHECK (total_orders > 0),
    processed INT NOT NULL DEFAULT 0 CHECK (processed >= 0),
    successful INT NOT NULL DEFAULT 0 CHECK (successful >= 0),
    failed INT NOT NULL DEFAULT 0 CHECK (failed >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT fk_bulk_invoice_jobs_business
        FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE,

    CONSTRAINT chk_processed_le_total
        CHECK (processed <= total_orders),

    CONSTRAINT chk_successful_le_processed
        CHECK (successful <= processed),

    CONSTRAINT chk_failed_le_processed
        CHECK (failed <= processed),

    CONSTRAINT chk_successful_plus_failed_le_processed
        CHECK (successful + failed <= processed)
);

-- Items individuales del job (una fila por orden a facturar)
CREATE TABLE IF NOT EXISTS bulk_invoice_job_items (
    id SERIAL PRIMARY KEY,
    job_id UUID NOT NULL,
    order_id VARCHAR(255) NOT NULL,
    invoice_id INT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'success', 'failed')),
    error_message TEXT,
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_bulk_job_items_job
        FOREIGN KEY (job_id) REFERENCES bulk_invoice_jobs(id) ON DELETE CASCADE,

    CONSTRAINT fk_bulk_job_items_invoice
        FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE SET NULL,

    CONSTRAINT uq_job_order
        UNIQUE (job_id, order_id)
);

-- Índices para optimizar consultas
CREATE INDEX IF NOT EXISTS idx_bulk_jobs_business
    ON bulk_invoice_jobs(business_id);

CREATE INDEX IF NOT EXISTS idx_bulk_jobs_status
    ON bulk_invoice_jobs(status)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_bulk_jobs_created_at
    ON bulk_invoice_jobs(created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_bulk_job_items_job
    ON bulk_invoice_job_items(job_id);

CREATE INDEX IF NOT EXISTS idx_bulk_job_items_status
    ON bulk_invoice_job_items(status);

CREATE INDEX IF NOT EXISTS idx_bulk_job_items_order
    ON bulk_invoice_job_items(order_id);

-- Función trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_bulk_invoice_jobs_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para actualizar updated_at en bulk_invoice_jobs
CREATE TRIGGER trigger_bulk_invoice_jobs_updated_at
    BEFORE UPDATE ON bulk_invoice_jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_bulk_invoice_jobs_updated_at();

-- Comentarios para documentación
COMMENT ON TABLE bulk_invoice_jobs IS 'Jobs de facturación masiva - tracking de procesamiento asíncrono';
COMMENT ON TABLE bulk_invoice_job_items IS 'Items individuales de cada job - una fila por orden';

COMMENT ON COLUMN bulk_invoice_jobs.status IS 'pending: creado, processing: al menos un item procesándose, completed: todos procesados, failed: error crítico';
COMMENT ON COLUMN bulk_invoice_job_items.status IS 'pending: publicado a queue, processing: consumer tomó mensaje, success: factura creada, failed: error al crear';
COMMENT ON COLUMN bulk_invoice_jobs.processed IS 'Cantidad de items procesados (success + failed)';
COMMENT ON COLUMN bulk_invoice_jobs.successful IS 'Cantidad de facturas creadas exitosamente';
COMMENT ON COLUMN bulk_invoice_jobs.failed IS 'Cantidad de facturas que fallaron';
