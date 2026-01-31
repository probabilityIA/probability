-- ==========================================
-- SEED: Tipos de Proveedores de Facturación
-- Autor: Claude Code
-- Fecha: 2026-01-31
-- Descripción: Inserta los tipos de proveedores de facturación iniciales
-- ==========================================

-- Limpiar datos existentes (solo en desarrollo)
-- TRUNCATE TABLE invoicing_provider_types RESTART IDENTITY CASCADE;

-- ==========================================
-- 1. SOFTPYMES (Colombia)
-- ==========================================
INSERT INTO invoicing_provider_types (
    name,
    code,
    description,
    icon,
    image_url,
    api_base_url,
    documentation_url,
    is_active,
    supported_countries,
    created_at,
    updated_at
) VALUES (
    'Softpymes',
    'softpymes',
    'Proveedor de facturación electrónica para Colombia. Permite generar facturas, notas de crédito y otros documentos tributarios de acuerdo con la DIAN.',
    'receipt',
    '/providers/softpymes-logo.png',
    'https://api-integracion.softpymes.com.co/app/integration/',
    'https://docs.softpymes.com.co/api/',
    TRUE,
    'CO',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    api_base_url = EXCLUDED.api_base_url,
    documentation_url = EXCLUDED.documentation_url,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

-- ==========================================
-- 2. SIIGO (Colombia, México, Chile)
-- ==========================================
INSERT INTO invoicing_provider_types (
    name,
    code,
    description,
    icon,
    image_url,
    api_base_url,
    documentation_url,
    is_active,
    supported_countries,
    created_at,
    updated_at
) VALUES (
    'Siigo',
    'siigo',
    'Software de facturación electrónica y gestión empresarial para Colombia, México y Chile. Cumple con normativas fiscales locales.',
    'receipt',
    '/providers/siigo-logo.png',
    'https://api.siigo.com/v1/',
    'https://developers.siigo.com/',
    FALSE,
    'CO,MX,CL',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    api_base_url = EXCLUDED.api_base_url,
    documentation_url = EXCLUDED.documentation_url,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

-- ==========================================
-- 3. FACTURAMA (México)
-- ==========================================
INSERT INTO invoicing_provider_types (
    name,
    code,
    description,
    icon,
    image_url,
    api_base_url,
    documentation_url,
    is_active,
    supported_countries,
    created_at,
    updated_at
) VALUES (
    'Facturama',
    'facturama',
    'Proveedor de facturación electrónica para México. Genera CFDIs (Comprobantes Fiscales Digitales por Internet) certificados por el SAT.',
    'receipt',
    '/providers/facturama-logo.png',
    'https://api.facturama.mx/',
    'https://www.facturama.mx/api/docs',
    FALSE,
    'MX',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    api_base_url = EXCLUDED.api_base_url,
    documentation_url = EXCLUDED.documentation_url,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

-- ==========================================
-- 4. ALEGRA (Multi-país: CO, MX, PE, AR, CL, CR, PA)
-- ==========================================
INSERT INTO invoicing_provider_types (
    name,
    code,
    description,
    icon,
    image_url,
    api_base_url,
    documentation_url,
    is_active,
    supported_countries,
    created_at,
    updated_at
) VALUES (
    'Alegra',
    'alegra',
    'Plataforma de facturación electrónica y contabilidad para múltiples países de Latinoamérica.',
    'receipt',
    '/providers/alegra-logo.png',
    'https://api.alegra.com/api/v1/',
    'https://developer.alegra.com/',
    FALSE,
    'CO,MX,PE,AR,CL,CR,PA',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    api_base_url = EXCLUDED.api_base_url,
    documentation_url = EXCLUDED.documentation_url,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

-- ==========================================
-- 5. NUBEFACT (Perú)
-- ==========================================
INSERT INTO invoicing_provider_types (
    name,
    code,
    description,
    icon,
    image_url,
    api_base_url,
    documentation_url,
    is_active,
    supported_countries,
    created_at,
    updated_at
) VALUES (
    'NubeFact',
    'nubefact',
    'Sistema de facturación electrónica para Perú. Genera facturas, boletas y otros comprobantes de pago electrónicos certificados por SUNAT.',
    'receipt',
    '/providers/nubefact-logo.png',
    'https://api.nubefact.com/api/v1/',
    'https://nubefact.com/api',
    FALSE,
    'PE',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    api_base_url = EXCLUDED.api_base_url,
    documentation_url = EXCLUDED.documentation_url,
    is_active = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

-- ==========================================
-- VERIFICACIÓN
-- ==========================================
-- Mostrar los proveedores insertados
SELECT
    id,
    name,
    code,
    supported_countries,
    is_active,
    created_at
FROM invoicing_provider_types
ORDER BY id;

-- ==========================================
-- NOTAS
-- ==========================================
-- - Solo Softpymes está activo inicialmente (is_active = TRUE)
-- - Los demás proveedores están preparados para futuras integraciones
-- - Ejecutar después de las migraciones de tablas
-- - El seeder usa ON CONFLICT para evitar duplicados
--
-- Para ejecutar:
-- psql -h localhost -p 5433 -U postgres -d probability -f seed_invoicing_providers.sql
