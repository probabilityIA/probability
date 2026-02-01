-- Actualizar schemas de Softpymes para el formulario dinámico
UPDATE integration_types
SET
    credentials_schema = '{
        "type": "object",
        "properties": {
            "api_key": {
                "type": "string",
                "title": "API Key",
                "description": "Clave de API proporcionada por Softpymes",
                "required": true,
                "order": 1,
                "placeholder": "Ingresa tu API Key de Softpymes",
                "error_message": "La API Key es requerida"
            },
            "api_secret": {
                "type": "string",
                "title": "API Secret",
                "description": "Secreto de API proporcionado por Softpymes",
                "required": true,
                "order": 2,
                "placeholder": "Ingresa tu API Secret de Softpymes",
                "error_message": "El API Secret es requerido",
                "format": "password"
            }
        },
        "required": ["api_key", "api_secret"]
    }'::jsonb,
    config_schema = '{
        "type": "object",
        "properties": {}
    }'::jsonb,
    setup_instructions = '# Configuración de Softpymes

## Paso 1: Obtener credenciales
1. Inicia sesión en tu cuenta de Softpymes
2. Ve a Configuración → Integraciones → API
3. Copia tu **API Key** y **API Secret**

## Paso 2: Configurar integración
1. Pega tu API Key en el campo correspondiente
2. Pega tu API Secret en el campo correspondiente
3. Haz clic en "Probar Conexión" para validar las credenciales
4. Si la prueba es exitosa, guarda la integración

## Notas importantes
- Mantén tus credenciales seguras y no las compartas
- Las credenciales se almacenan de forma encriptada
- Puedes tener múltiples configuraciones de Softpymes para diferentes negocios'
WHERE code = 'softpymes';

-- Verificar el cambio
SELECT
    id,
    name,
    code,
    credentials_schema->'properties'->>'api_key' as has_api_key,
    credentials_schema->'properties'->>'api_secret' as has_api_secret,
    LEFT(setup_instructions, 50) as instructions_preview
FROM integration_types
WHERE code = 'softpymes';
