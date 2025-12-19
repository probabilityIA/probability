-- Script para corregir la secuencia de la tabla integrations
-- Ejecutar este script en PostgreSQL para solucionar el error de clave primaria duplicada

-- 1. Verificar el estado actual de la secuencia
SELECT 
    pg_get_serial_sequence('integrations', 'id') as sequence_name,
    currval(pg_get_serial_sequence('integrations', 'id')) as current_value,
    (SELECT MAX(id) FROM integrations) as max_id;

-- 2. Corregir la secuencia para que sea mayor que el máximo ID existente
-- Esto asegura que los próximos IDs generados no entren en conflicto
SELECT setval(
    pg_get_serial_sequence('integrations', 'id'),
    COALESCE((SELECT MAX(id) FROM integrations), 0) + 1,
    false
);

-- 3. Verificar que se corrigió correctamente
SELECT 
    pg_get_serial_sequence('integrations', 'id') as sequence_name,
    currval(pg_get_serial_sequence('integrations', 'id')) as current_value,
    (SELECT MAX(id) FROM integrations) as max_id;




