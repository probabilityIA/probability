# Quick Start: Agregar Nueva Integración

Guía rápida de 5 pasos para agregar una nueva integración al sistema.

---

## 🚀 Ejemplo: Agregar "Alegra" (Facturación)

### Paso 1: Backend - Crear IntegrationType

```sql
-- 1.1 Verificar que existe la categoría
SELECT * FROM integration_categories WHERE code = 'invoicing';
-- Anotar el ID (ej: 2)

-- 1.2 Crear el tipo de integración
INSERT INTO integration_types (
    name,
    code,
    category_id,
    description,
    icon,
    is_active,
    config_schema,
    credentials_schema
)
VALUES (
    'Alegra',
    'alegra',
    2, -- ID de categoría "invoicing"
    'Sistema de facturación electrónica Alegra',
    'DocumentTextIcon',
    true,
    '{
        "company_id": {"type": "string", "required": true, "label": "ID de Empresa"},
        "environment": {"type": "select", "required": true, "label": "Ambiente", "options": ["test", "production"]}
    }'::jsonb,
    '{
        "api_token": {"type": "password", "required": true, "label": "API Token"},
        "api_secret": {"type": "password", "required": true, "label": "API Secret"}
    }'::jsonb
);

-- Anotar el ID generado (ej: 9)
```

---

### Paso 2: Frontend - Crear Estructura de Carpetas

```bash
cd front/central/src/services/integrations/invoicing
mkdir -p alegra/{domain,ui/components}
```

---

### Paso 3: Definir Tipos de Dominio

**Archivo:** `alegra/domain/types.ts`

```typescript
export interface AlegraConfig {
    company_id: string;
    environment: 'test' | 'production';
}

export interface AlegraCredentials {
    api_token: string;
    api_secret: string;
}
```

---

### Paso 4: Crear Formulario de Configuración

**Archivo:** `alegra/ui/components/AlegraConfigForm.tsx`

```typescript
'use client';

import { useState, FormEvent } from 'react';
import { Button, Input, Select } from '@/shared/ui';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';

export function AlegraConfigForm({ onSuccess, onCancel }) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [formData, setFormData] = useState({
        name: '',
        company_id: '',
        environment: 'production',
        api_token: '',
        api_secret: ''
    });

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            await createIntegrationAction({
                name: formData.name,
                code: `alegra_${Date.now()}`,
                integration_type_id: 9, // ⚠️ Usar el ID del Paso 1
                category: 'invoicing',
                config: {
                    company_id: formData.company_id,
                    environment: formData.environment
                },
                credentials: {
                    api_token: formData.api_token,
                    api_secret: formData.api_secret
                },
                is_active: true
            });

            showToast('Integración Alegra creada exitosamente', 'success');
            onSuccess?.();
        } catch (error: any) {
            showToast('Error al crear integración', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6 p-6">
            <div>
                <h3 className="text-lg font-semibold mb-4">Configurar Alegra</h3>
            </div>

            <Input
                label="Nombre de la integración *"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="Alegra - Producción"
                required
            />

            <Input
                label="ID de Empresa *"
                value={formData.company_id}
                onChange={(e) => setFormData({ ...formData, company_id: e.target.value })}
                placeholder="12345"
                required
            />

            <Select
                label="Ambiente"
                value={formData.environment}
                onChange={(e) => setFormData({ ...formData, environment: e.target.value as any })}
            >
                <option value="production">Producción</option>
                <option value="test">Pruebas</option>
            </Select>

            <Input
                type="password"
                label="API Token *"
                value={formData.api_token}
                onChange={(e) => setFormData({ ...formData, api_token: e.target.value })}
                required
            />

            <Input
                type="password"
                label="API Secret *"
                value={formData.api_secret}
                onChange={(e) => setFormData({ ...formData, api_secret: e.target.value })}
                required
            />

            <div className="flex justify-end gap-3 pt-4 border-t">
                {onCancel && (
                    <Button variant="outline" onClick={onCancel} disabled={loading}>
                        Cancelar
                    </Button>
                )}
                <Button type="submit" variant="primary" disabled={loading}>
                    {loading ? 'Conectando...' : 'Conectar Alegra'}
                </Button>
            </div>
        </form>
    );
}
```

---

### Paso 5: Exportar Componentes

**Archivo:** `alegra/ui/components/index.ts`

```typescript
export { AlegraConfigForm } from './AlegraConfigForm';
```

---

## ✅ ¡Listo!

Tu nueva integración ya está disponible en el sistema:

1. Ve a `/integrations`
2. Click en "Nueva Integración"
3. Selecciona categoría "Facturación"
4. Aparecerá "Alegra" en la lista de proveedores
5. Completa el formulario y conecta

---

## 📋 Checklist de Validación

- [ ] Backend: IntegrationType creado con `category_id` correcto
- [ ] Frontend: Carpeta creada en categoría correspondiente
- [ ] `domain/types.ts` con interfaces TypeScript
- [ ] `ui/components/[Provider]ConfigForm.tsx` implementado
- [ ] `ui/components/index.ts` exporta el componente
- [ ] Compilación exitosa: `pnpm build`
- [ ] Probar flujo: Categoría → Proveedor → Configurar
- [ ] Integración aparece en lista de integraciones

---

## 🔧 Personalización Avanzada

### Agregar Vista Personalizada

**Archivo:** `alegra/ui/components/AlegraIntegrationView.tsx`

```typescript
'use client';

import { Integration } from '@/services/integrations/core/domain/types';
import { Badge, Button } from '@/shared/ui';

export function AlegraIntegrationView({ integration, onEdit }) {
    const config = integration.config as any;

    return (
        <div className="bg-white border rounded-lg p-6">
            <div className="flex items-start justify-between">
                <div>
                    <h3 className="text-lg font-semibold">{integration.name}</h3>
                    <p className="text-sm text-gray-500">Alegra - Facturación</p>
                </div>
                {integration.is_active ? (
                    <Badge type="success">Activo</Badge>
                ) : (
                    <Badge type="error">Inactivo</Badge>
                )}
            </div>

            <div className="mt-4 space-y-2 text-sm">
                <div className="flex justify-between">
                    <span className="text-gray-600">Empresa:</span>
                    <span className="font-medium">{config?.company_id}</span>
                </div>
                <div className="flex justify-between">
                    <span className="text-gray-600">Ambiente:</span>
                    <Badge type={config?.environment === 'production' ? 'primary' : 'warning'}>
                        {config?.environment === 'production' ? 'Producción' : 'Pruebas'}
                    </Badge>
                </div>
            </div>

            <div className="mt-4 flex gap-2 pt-4 border-t">
                <Button size="sm" variant="outline" onClick={onEdit}>
                    Editar
                </Button>
            </div>
        </div>
    );
}
```

### Agregar Lógica de Prueba de Conexión

**Archivo:** `alegra/infra/actions/index.ts`

```typescript
'use server';

export async function testAlegraConnection(config: any, credentials: any) {
    try {
        const response = await fetch('https://api.alegra.com/api/v1/companies', {
            headers: {
                'Authorization': `Bearer ${credentials.api_token}`,
                'Content-Type': 'application/json'
            }
        });

        if (response.ok) {
            return { success: true, message: 'Conexión exitosa' };
        } else {
            return { success: false, message: 'Credenciales inválidas' };
        }
    } catch (error: any) {
        return { success: false, message: error.message };
    }
}
```

Luego, agregar botón "Probar Conexión" en el formulario:

```typescript
import { testAlegraConnection } from '../../infra/actions';

// En el componente
const handleTest = async () => {
    const result = await testAlegraConnection(
        { company_id: formData.company_id, environment: formData.environment },
        { api_token: formData.api_token, api_secret: formData.api_secret }
    );

    showToast(result.message, result.success ? 'success' : 'error');
};

// En el JSX
<Button type="button" variant="outline" onClick={handleTest}>
    Probar Conexión
</Button>
```

---

## 🎯 Próximos Pasos

1. **Agregar Tests:** Crear tests unitarios para el formulario
2. **Validaciones:** Agregar validaciones específicas (regex, longitud, etc.)
3. **Webhooks:** Implementar configuración de webhooks si el proveedor lo soporta
4. **Documentación:** Agregar `README.md` específico en la carpeta del proveedor
5. **Localización:** Agregar traducciones si el sistema es multiidioma

---

**Última actualización:** 2026-01-31
