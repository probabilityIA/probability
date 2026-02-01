# Integrations Service

## Descripción

Sistema de integración con plataformas externas organizado por categorías. Permite a los negocios conectar múltiples proveedores de e-commerce, facturación, mensajería y otros servicios.

## Arquitectura por Categorías

Las integraciones están organizadas en categorías para una mejor experiencia de usuario:

```
services/integrations/
├── core/                  # Hub central - Infraestructura compartida
│   ├── domain/           # Tipos, ports, entidades
│   ├── app/              # Use cases compartidos
│   ├── infra/            # Repository, actions
│   └── ui/               # Componentes genéricos
│       ├── components/   # Modales, tabs, forms
│       ├── hooks/        # useIntegrations, useCategories
│       └── index.ts      # Exports públicos
│
├── ecommerce/            # Categoría: E-commerce
│   └── shopify/          # (Existente - bajo core/ui/components)
│
├── messaging/            # Categoría: Mensajería
│   └── whatsapp/         # (Existente - bajo core/ui/components)
│
├── invoicing/            # Categoría: Facturación
│   └── softpymes/        # ✅ Ejemplo completo
│       ├── domain/
│       ├── ui/
│       └── infra/
│
└── notification-config/  # Configuración de notificaciones (separado)
```

---

## Categorías de Integración

### 1. **E-commerce** (`ecommerce`)
Plataformas de comercio electrónico.

**Proveedores:**
- Shopify
- MercadoLibre
- WooCommerce

### 2. **Facturación** (`invoicing`)
Sistemas de facturación electrónica.

**Proveedores:**
- Softpymes
- Siigo
- Factus

### 3. **Mensajería** (`messaging`)
Canales de comunicación con clientes.

**Proveedores:**
- WhatsApp
- Email
- SMS

### 4. **Sistema** (`system`)
Integraciones internas o sin categoría específica.

---

## Flujo de Creación de Integración (2 Pasos)

### UI/UX del Usuario

1. **Paso 1: Seleccionar Categoría**
   - Usuario hace click en "Nueva Integración"
   - Se muestra un grid de categorías (E-commerce, Facturación, Mensajería, Sistema)
   - Cada categoría muestra icono, nombre y descripción

2. **Paso 2: Seleccionar Proveedor**
   - Se filtran solo los proveedores de la categoría seleccionada
   - Grid con logos, nombres y descripciones de proveedores
   - Botón "← Volver a categorías" para regresar

3. **Paso 3: Configurar Credenciales**
   - Formulario específico del proveedor seleccionado
   - Campos dinámicos según `config_schema` y `credentials_schema`
   - Botón "Probar Conexión" (opcional)
   - Botón "← Volver a proveedores"

### Componentes Involucrados

```typescript
// Modal principal que orquesta los 3 pasos
<CreateIntegrationModal
    isOpen={showModal}
    onClose={handleClose}
    categories={categories}
    onSuccess={handleSuccess}
/>

// Internamente:
// Paso 1
<CategorySelector
    categories={categories}
    onSelect={(category) => setStep(2)}
/>

// Paso 2
<ProviderSelector
    category={selectedCategory}
    onSelect={(provider) => setStep(3)}
    onBack={() => setStep(1)}
/>

// Paso 3
<DynamicIntegrationForm
    integrationType={selectedProvider}
    onSubmit={handleCreate}
    onCancel={handleClose}
    onTest={handleTestConnection}
/>
```

---

## Navegación por Categorías (CategoryTabs)

### UI en la Página Principal

En `/integrations`, se muestran dos niveles de tabs:

**Nivel 1: Integraciones vs Tipos**
- Mis Integraciones
- Tipos de Integración

**Nivel 2: Categorías (solo en "Mis Integraciones")**
- Todas
- E-commerce
- Facturación
- Mensajería
- Sistema

### Implementación

```typescript
// Hooks
const { categories, loading } = useCategories();
const { setFilterCategory, refresh } = useIntegrations();

// Handler
const handleCategoryChange = (categoryCode: string | null) => {
    setActiveCategoryCode(categoryCode);
    setFilterCategory(categoryCode || '');
    refresh();
};

// Renderizado
<CategoryTabs
    categories={categories}
    activeCategory={activeCategoryCode}
    onSelectCategory={handleCategoryChange}
/>
```

---

## Cómo Agregar una Nueva Integración por Categoría

### Ejemplo: Agregar "Siigo" (Facturación)

#### 1. Backend: Crear IntegrationType

```sql
INSERT INTO integration_types (name, code, category_id, description, is_active)
VALUES (
    'Siigo',
    'siigo',
    2, -- ID de categoría "invoicing"
    'Sistema de facturación electrónica Siigo',
    true
);
```

#### 2. Frontend: Crear Módulo

```bash
mkdir -p src/services/integrations/invoicing/siigo/{domain,ui/components}
```

**Estructura:**
```
invoicing/siigo/
├── domain/
│   └── types.ts
├── ui/
│   └── components/
│       ├── SiigoConfigForm.tsx
│       ├── SiigoIntegrationView.tsx
│       └── index.ts
└── infra/ (opcional)
    └── actions/
        └── index.ts
```

#### 3. Definir Tipos de Dominio

```typescript
// domain/types.ts
export interface SiigoConfig {
    company_nit: string;
    company_name: string;
    api_key: string;
    environment: 'test' | 'production';
}

export interface SiigoCredentials {
    username: string;
    access_key: string;
}
```

#### 4. Crear Formulario de Configuración

```typescript
// ui/components/SiigoConfigForm.tsx
'use client';

import { useState } from 'react';
import { Button, Input } from '@/shared/ui';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';

export function SiigoConfigForm({ onSuccess, onCancel }) {
    const [formData, setFormData] = useState({
        name: '',
        company_nit: '',
        company_name: '',
        username: '',
        access_key: '',
        environment: 'production'
    });

    const handleSubmit = async (e) => {
        e.preventDefault();

        const config = {
            company_nit: formData.company_nit,
            company_name: formData.company_name,
            environment: formData.environment
        };

        const credentials = {
            username: formData.username,
            access_key: formData.access_key
        };

        await createIntegrationAction({
            name: formData.name,
            code: `siigo_${Date.now()}`,
            integration_type_id: 8, // ID de Siigo en backend
            category: 'invoicing',
            config,
            credentials,
            is_active: true
        });

        onSuccess?.();
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {/* Campos del formulario */}
            <Input
                label="Nombre de la integración"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                required
            />

            <Input
                label="NIT de la empresa"
                value={formData.company_nit}
                onChange={(e) => setFormData({ ...formData, company_nit: e.target.value })}
                required
            />

            {/* ... más campos ... */}

            <Button type="submit">Conectar Siigo</Button>
        </form>
    );
}
```

#### 5. Crear Componente de Vista

```typescript
// ui/components/SiigoIntegrationView.tsx
'use client';

import { Integration } from '@/services/integrations/core/domain/types';
import { Badge, Button } from '@/shared/ui';

export function SiigoIntegrationView({ integration, onEdit }) {
    return (
        <div className="bg-white border rounded-lg p-6">
            <h3>{integration.name}</h3>
            <p>Siigo - Facturación</p>

            {integration.is_active ? (
                <Badge type="success">Activo</Badge>
            ) : (
                <Badge type="error">Inactivo</Badge>
            )}

            <Button onClick={onEdit}>Editar</Button>
        </div>
    );
}
```

#### 6. Exportar Componentes

```typescript
// ui/components/index.ts
export { SiigoConfigForm } from './SiigoConfigForm';
export { SiigoIntegrationView } from './SiigoIntegrationView';
```

#### 7. Usar en DynamicIntegrationForm (Automático)

El `DynamicIntegrationForm` del core ya renderiza formularios dinámicos basados en el `config_schema` del IntegrationType. Si necesitas un formulario personalizado:

```typescript
// En core/ui/components/IntegrationForm.tsx
import { SiigoConfigForm } from '@/services/integrations/invoicing/siigo/ui/components';

// Dentro del switch/case por tipo de integración
case 'siigo':
    return <SiigoConfigForm onSuccess={onSuccess} onCancel={onCancel} />;
```

---

## Hooks Disponibles

### `useCategories()`

Obtiene y gestiona las categorías de integración.

```typescript
const { categories, loading, error, refresh } = useCategories();

// Retorna:
// - categories: IntegrationCategory[]
// - loading: boolean
// - error: string | null
// - refresh: () => void
```

### `useIntegrations()`

Gestiona integraciones con filtrado por categoría.

```typescript
const {
    integrations,
    loading,
    filterCategory,
    setFilterCategory,
    refresh
} = useIntegrations();

// Filtrar por categoría
setFilterCategory('invoicing'); // Solo integraciones de facturación
setFilterCategory(''); // Todas las integraciones
```

---

## Componentes Compartidos

### `<CategoryTabs />`

Navegación horizontal por categorías.

```typescript
<CategoryTabs
    categories={categories}
    activeCategory={activeCategoryCode}
    onSelectCategory={(code) => handleCategoryChange(code)}
/>
```

### `<CreateIntegrationModal />`

Modal de 3 pasos para crear integraciones.

```typescript
<CreateIntegrationModal
    isOpen={showModal}
    onClose={() => setShowModal(false)}
    categories={categories}
    onSuccess={() => {
        setShowModal(false);
        refresh();
    }}
/>
```

### `<CategorySelector />`

Paso 1: Selección de categoría (grid).

```typescript
<CategorySelector
    categories={categories}
    onSelect={(category) => console.log(category)}
/>
```

### `<ProviderSelector />`

Paso 2: Selección de proveedor filtrado por categoría.

```typescript
<ProviderSelector
    category={selectedCategory}
    onSelect={(provider) => console.log(provider)}
    onBack={() => goBackToCategories()}
/>
```

---

## Tipos y Interfaces

### `IntegrationCategory`

```typescript
interface IntegrationCategory {
    id: number;
    code: string;                    // 'ecommerce', 'invoicing', 'messaging'
    name: string;                    // 'E-commerce', 'Facturación', 'Mensajería'
    description?: string;
    icon?: string;                   // Icon name (heroicons)
    color?: string;                  // Tailwind color class
    display_order: number;
    parent_category_id?: number;
    is_active: boolean;
    is_visible: boolean;
    created_at: string;
    updated_at: string;
}
```

### `IntegrationType`

```typescript
interface IntegrationType {
    id: number;
    name: string;
    code: string;
    description?: string;
    icon?: string;
    image_url?: string;
    category: string;                           // Legacy field
    category_id?: number;                       // NEW - FK a IntegrationCategory
    integration_category?: IntegrationCategory; // NEW - Populated
    is_active: boolean;
    config_schema?: any;
    credentials_schema?: any;
    setup_instructions?: string;
    created_at: string;
    updated_at: string;
}
```

---

## Server Actions

### Categorías

```typescript
// Obtener todas las categorías
const response = await getIntegrationCategoriesAction();
// response.data: IntegrationCategory[]
```

### Integraciones

```typescript
// Obtener integraciones con filtro de categoría
const response = await getIntegrationsAction({
    category_id: 2, // ID de categoría "invoicing"
    page: 1,
    page_size: 10
});

// Crear integración
const response = await createIntegrationAction({
    name: 'Mi Integración',
    code: 'my_integration',
    integration_type_id: 7,
    category: 'invoicing',
    config: { /* ... */ },
    credentials: { /* ... */ },
    is_active: true
});
```

---

## Checklist para Nueva Integración

- [ ] Backend: Crear IntegrationType con `category_id` correcto
- [ ] Frontend: Crear carpeta en categoría correspondiente
- [ ] Crear `domain/types.ts` con interfaces de config y credentials
- [ ] Crear `ui/components/[Provider]ConfigForm.tsx`
- [ ] Crear `ui/components/[Provider]IntegrationView.tsx`
- [ ] Exportar en `ui/components/index.ts`
- [ ] (Opcional) Agregar caso específico en `IntegrationForm.tsx` si necesitas lógica custom
- [ ] Probar flujo completo: Categoría → Proveedor → Configurar
- [ ] Verificar filtrado por categoría en `/integrations`

---

## Convenciones

1. **Nombres de carpetas:** lowercase con guiones (e.g., `softpymes`, `mercado-libre`)
2. **Nombres de componentes:** PascalCase (e.g., `SoftpymesConfigForm`)
3. **Código de integración:** lowercase_snake_case (e.g., `softpymes_12345`)
4. **Category codes:** lowercase (e.g., `ecommerce`, `invoicing`, `messaging`, `system`)

---

## Troubleshooting

### Las categorías no aparecen en los tabs

**Causa:** Backend no tiene categorías seeded o `is_visible=false`

**Solución:**
```sql
-- Verificar categorías
SELECT * FROM integration_categories WHERE is_active = true AND is_visible = true;

-- Crear categorías faltantes
INSERT INTO integration_categories (code, name, description, display_order, is_active, is_visible)
VALUES
    ('ecommerce', 'E-commerce', 'Plataformas de comercio electrónico', 1, true, true),
    ('invoicing', 'Facturación', 'Sistemas de facturación electrónica', 2, true, true),
    ('messaging', 'Mensajería', 'Canales de comunicación', 3, true, true),
    ('system', 'Sistema', 'Integraciones internas', 4, true, true);
```

### No aparecen proveedores en una categoría

**Causa:** IntegrationTypes no tienen `category_id` configurado

**Solución:**
```sql
-- Verificar integration types
SELECT id, name, code, category_id FROM integration_types WHERE is_active = true;

-- Asignar categoría a un tipo
UPDATE integration_types SET category_id = 2 WHERE code = 'softpymes';
```

### Error al crear integración

**Causa:** `integration_type_id` incorrecto o credenciales inválidas

**Solución:**
- Verificar que el ID del tipo existe y está activo
- Revisar que los campos de config/credentials coincidan con los schemas del backend
- Verificar logs del servidor para más detalles

---

**Última actualización:** 2026-01-31
**Versión:** 1.0.0
