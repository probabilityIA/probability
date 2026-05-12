# Guía: Refactorizar Modales de Envío a Colores Dinámicos del Negocio

## Visión General

Los modales de generación de guías de envío (`ShipmentGuideModal` y `MassGuideGenerationModal`) actualmente tienen colores hardcodeados (ej: `#7c3aed`, `#a855f7`). El objetivo es usar las **variables CSS dinámicas** del negocio para que cada negocio pueda personalizar los colores.

### Sistema Actual
- **Colores dinámicos del negocio** se definen en: `globals.css` (líneas 9-21)
- Variables CSS: `--color-primary`, `--color-secondary`, `--color-tertiary`, `--color-quaternary`
- Se inyectan dinámicamente en `theme-provider.tsx` y se actualizan con `color-scales.ts`
- Se generan escalas de colores automáticamente (50-900)

---

## Cambios Necesarios

### 1. Importar CSS de Estilos Reutilizables

**Archivo:** `shipment-guide-modal.tsx` y `mass-guide-generation-modal.tsx`

```typescript
// Al inicio del archivo, después de los otros imports
import '@/shared/ui/styles/shipment-modals.css';
```

---

### 2. Reemplazar Colores Hardcodeados

#### ANTES: Modal Individual - Paso 1 (Línea 856)
```tsx
<div className="bg-purple-50/50 dark:bg-purple-900/10 border border-purple-100 dark:border-purple-800/30 rounded-xl p-4 space-y-2">
    <div className="flex items-center gap-2">
        <div className="w-8 h-8 rounded-lg bg-purple-100 dark:bg-purple-800/40 flex items-center justify-center text-purple-600 dark:text-purple-400 text-sm font-bold">A</div>
        <h3 className="font-semibold text-base text-purple-700 dark:text-purple-400">Origen</h3>
    </div>
```

#### DESPUÉS: Usando Variable Dinámica
```tsx
<div className="shipment-section-origin">
    <div className="flex items-center gap-2">
        <div className="shipment-section-origin-icon w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold">A</div>
        <h3 className="shipment-section-origin-label">Origen</h3>
    </div>
```

---

### 3. Formulario - Sección Destino

#### ANTES (Línea 948)
```tsx
<div className="bg-blue-50/50 dark:bg-blue-900/10 border border-blue-100 dark:border-blue-800/30 rounded-xl p-4 space-y-2">
    <div className="flex items-center gap-2">
        <div className="w-8 h-8 rounded-lg bg-blue-100 dark:bg-blue-800/40 flex items-center justify-center text-blue-600 dark:text-blue-400 text-sm font-bold">B</div>
        <h3 className="font-semibold text-base text-blue-700 dark:text-blue-400">Destino</h3>
    </div>
```

#### DESPUÉS
```tsx
<div className="shipment-section-destination">
    <div className="flex items-center gap-2">
        <div className="shipment-section-destination-icon w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold">B</div>
        <h3 className="shipment-section-destination-label">Destino</h3>
    </div>
```

---

### 4. Inputs del Formulario

#### ANTES (Línea 883)
```tsx
<input
    type="text"
    className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white ${step1Form.formState.errors.originDaneCode ? "border-red-500 bg-red-50 dark:bg-red-900/20" : "border-gray-300 dark:border-gray-600"}`}
/>
```

#### DESPUÉS
```tsx
<input
    type="text"
    className={`shipment-input ${step1Form.formState.errors.originDaneCode ? "shipment-input-error" : ""}`}
/>
```

---

### 5. Botón Principal (Paso 1 - Siguiente)

#### ANTES (Línea 1503)
```tsx
<Button
    variant="primary"
    onClick={...}
    disabled={loading}
    style={{ background: '#7c3aed' }}
>
    {loading ? "Cotizando..." : "Siguiente"}
</Button>
```

#### DESPUÉS
```tsx
<Button
    className="shipment-btn-primary"
    onClick={...}
    disabled={loading}
>
    {loading ? "Cotizando..." : "Siguiente"}
</Button>
```

---

### 6. Tarjetas de Transportista (Paso 2)

#### ANTES (Línea 1199)
```tsx
<div
    key={rate.idRate}
    onClick={() => handleRateSelection(rate)}
    className="border border-gray-200 dark:border-gray-600 rounded-lg p-3 hover:border-purple-500 hover:shadow-md cursor-pointer transition-all bg-white dark:bg-gray-800"
>
    <div className="grid grid-cols-3 gap-3 h-full">
        <div className="col-span-1 flex flex-col items-center justify-center">
            <div className={`${getCarrierLogoSize(rate.carrier).container} bg-purple-50 rounded-lg flex items-center justify-center overflow-hidden`}>
```

#### DESPUÉS
```tsx
<div
    key={rate.idRate}
    onClick={() => handleRateSelection(rate)}
    className={`shipment-carrier-card ${selectedRate?.idRate === rate.idRate ? 'shipment-carrier-card-selected' : ''}`}
>
    <div className="grid grid-cols-3 gap-3 h-full">
        <div className="col-span-1 flex flex-col items-center justify-center">
            <div className={`${getCarrierLogoSize(rate.carrier).container} shipment-carrier-logo-container rounded-lg flex items-center justify-center overflow-hidden`}>
```

---

### 7. Badges de Información (COD, Contra Entrega)

#### ANTES (Línea 1145, Paso 2)
```tsx
<span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-100 text-amber-700 border border-amber-300 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-600">
    Contra Entrega - Solo opciones contra entrega
</span>
```

#### DESPUÉS
```tsx
<span className="shipment-badge-warning">
    Contra Entrega - Solo opciones contra entrega
</span>
```

#### Variantes de Badges
```tsx
// Primario (color del negocio)
<span className="shipment-badge-primary">Badge Text</span>

// Secundario
<span className="shipment-badge-secondary">Badge Text</span>

// Éxito
<span className="shipment-badge-success">✓ Cotizada</span>

// Error
<span className="shipment-badge-error">✗ Error</span>
```

---

### 8. Alertas de Error

#### ANTES (Línea 826)
```tsx
{error && (
    <div className="mb-3 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg text-red-700 dark:text-red-400 text-sm">
        {error.includes('\n') ? (...) : error}
    </div>
)}
```

#### DESPUÉS
```tsx
{error && (
    <div className="shipment-alert shipment-alert-error text-sm">
        {error.includes('\n') ? (...) : error}
    </div>
)}
```

#### Variantes de Alertas
```tsx
<div className="shipment-alert shipment-alert-primary">Información primaria</div>
<div className="shipment-alert shipment-alert-secondary">Información secundaria</div>
<div className="shipment-alert shipment-alert-success">Éxito</div>
<div className="shipment-alert shipment-alert-warning">Advertencia</div>
<div className="shipment-alert shipment-alert-error">Error</div>
```

---

### 9. Spinner/Cargador (Paso 2)

#### ANTES (Línea 1167)
```tsx
<div style={{ width: 28, height: 28, border: '3px solid #a855f7', borderTopColor: 'transparent', borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} />
```

#### DESPUÉS
```tsx
<div className="shipment-spinner" style={{ width: 28, height: 28 }} />
```

---

### 10. Botones Generales

#### ANTES (Múltiples ubicaciones)
```tsx
// Botón primario con color hardcodeado
<Button variant="primary" style={{ background: '#7c3aed' }}>Generar</Button>

// Botón con hover personalizado
<button className="... bg-green-600 hover:bg-green-700">...</button>
```

#### DESPUÉS
```tsx
// Botón primario (usa color del negocio)
<button className="shipment-btn-primary">Generar</button>

// Botón secundario (usa color secundario del negocio)
<button className="shipment-btn-secondary">Cancelar</button>

// Botón outline (usa color primario)
<button className="shipment-btn-outline">Volver</button>
```

---

### 11. Checkbox y Radio (Paso 1 y 3)

#### ANTES
```tsx
<input
    type="checkbox"
    {...step1Form.register("insurance")}
    className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
/>
```

#### DESPUÉS
```tsx
<input
    type="checkbox"
    {...step1Form.register("insurance")}
    className="shipment-checkbox"
/>
```

---

### 12. Sección de Éxito (Paso 4 - Guía Generada)

#### ANTES (Línea 1401)
```tsx
<div className="rounded-xl border-2 border-emerald-200 bg-emerald-50 p-4 flex flex-col items-center gap-3">
    <div className="w-12 h-12 rounded-full bg-emerald-100 flex items-center justify-center flex-shrink-0">
        <svg className="w-6 h-6 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
        </svg>
    </div>
    <div className="text-center min-w-0">
        <p className="font-bold text-emerald-800 text-sm">¡Guía generada exitosamente!</p>
```

#### DESPUÉS
```tsx
<div className="shipment-success-container">
    <div className="shipment-success-icon">
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
        </svg>
    </div>
    <div className="text-center min-w-0">
        <p className="shipment-success-text">¡Guía generada exitosamente!</p>
```

---

### 13. Modal Masivo - Tabla de Cotizaciones

#### ANTES (Línea 503)
```tsx
<thead className="bg-gray-50 sticky top-0">
    <tr className="border-b">
        <th className="text-left p-3 font-semibold">Orden</th>
```

#### DESPUÉS
```tsx
<thead className="sticky top-0">
    <tr className="border-b shipment-table-header">
        <th className="text-left p-3">Orden</th>
```

---

### 14. Barra de Progreso

#### ANTES (Línea 466, 597)
```tsx
<div className="w-full bg-gray-200 rounded-full h-4">
    <div
        className="bg-orange-500 h-4 rounded-full transition-all duration-300"
        style={{ width: `${quotingProgress}%` }}
    />
</div>
```

#### DESPUÉS
```tsx
<div className="shipment-progress-bar rounded-full h-4">
    <div
        className="shipment-progress-fill"
        style={{ width: `${quotingProgress}%` }}
    />
</div>
```

---

## Plan de Implementación

### Fase 1: Preparación (30 min)
- [ ] Crear archivo `shipment-modals.css` ✅ (ya hecho)
- [ ] Importar CSS en ambos modales

### Fase 2: Refactorizar ShipmentGuideModal (1.5 horas)
- [ ] Paso 1: Reemplazar secciones origen/destino
- [ ] Paso 1: Reemplazar inputs
- [ ] Paso 2: Reemplazar tarjetas de transportista
- [ ] Paso 3: Reemplazar inputs
- [ ] Paso 4: Reemplazar botones y secciones
- [ ] Reemplazar todos los colores hardcodeados

### Fase 3: Refactorizar MassGuideGenerationModal (1 hora)
- [ ] Selector de bodega
- [ ] Lista de órdenes
- [ ] Tabla de cotizaciones
- [ ] Botones

### Fase 4: Testing (1 hora)
- [ ] Probar con negocio default (colores por defecto)
- [ ] Cambiar colores del negocio y verificar actualización
- [ ] Probar en modo dark
- [ ] Verificar accesibilidad

---

## Variables CSS Disponibles

### Colores Base
```css
--color-primary      /* Color principal del negocio */
--color-secondary    /* Color secundario */
--color-tertiary     /* Color terciario (cyan/turquesa) */
--color-quaternary   /* Color cuaternario */
```

### Escalas de Color (Generadas Automáticamente)
```css
--color-primary-50, --color-primary-100, ..., --color-primary-900
--color-secondary-50, --color-secondary-100, ..., --color-secondary-900
--color-tertiary-50, --color-tertiary-100, ..., --color-tertiary-900
--color-quaternary-50, --color-quaternary-100, ..., --color-quaternary-900
```

### Estados Fijos (No Dinámicos)
```css
--success: #10b981
--error: #ef4444
--warning: #7c3aed   /* Este se puede cambiar si deseas que sea dinámico */
--info: #3b82f6
```

---

## Ejemplo: Aplicar Colores Dinámicos en Componente Personalizado

Si necesitas crear un componente personalizado dentro del modal:

```tsx
// En tu componente
<div className="shipment-section-origin">
    <h3 className="shipment-section-origin-label">Mi Sección</h3>
    <input className="shipment-input" />
    <button className="shipment-btn-primary">Enviar</button>
</div>

// O usar variables CSS directamente
<div style={{
    background: 'color-mix(in oklab, var(--color-primary) 10%, transparent)',
    borderColor: 'var(--color-primary)',
}}>
    Contenido personalizado
</div>
```

---

## Beneficios de Esta Implementación

✅ **Escalabilidad**: Los negocios pueden cambiar colores sin modificar código  
✅ **Consistencia**: Un solo origen de verdad para colores  
✅ **Mantenimiento**: Cambios en CSS se aplican automáticamente  
✅ **Accesibilidad**: Escalas de color generadas mantienen contraste  
✅ **Rendimiento**: Variables CSS se aplican instantáneamente  
✅ **DX**: Developers no necesitan pensar en colores hardcodeados  

---

## Checklist de Validación

Después de refactorizar, verifica:

- [ ] Los colores cambian cuando se cambia de negocio
- [ ] No hay colores `#7c3aed` directos en los modales
- [ ] Las escalas de color se generan correctamente (hover, active)
- [ ] Funciona en modo light y dark
- [ ] Los badges y alertas usan colores del negocio
- [ ] Los botones tienen el color correcto
- [ ] La barra de progreso usa color primario
- [ ] Las tarjetas de transportista se destacan correctamente
- [ ] El spinner/loader usa color primario
- [ ] No hay errores de TypeScript/ESLint
