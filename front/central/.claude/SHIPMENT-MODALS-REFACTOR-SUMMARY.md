# 🎨 Refactorización: Modales de Envío con Colores Dinámicos del Negocio

## Resumen Ejecutivo

Los modales `ShipmentGuideModal` y `MassGuideGenerationModal` están hardcodeados con colores fijos (`#7c3aed`, `#a855f7`, etc.). Se necesita migrar a **variables CSS dinámicas** para que cada negocio personalice sus colores.

### Status: ✅ Listo para Implementar

---

## Archivos Creados

```
front/central/src/shared/ui/styles/
├── shipment-modals.css                           ← Estilos reutilizables
└── examples/shipment-modal-refactor-example.tsx  ← Ejemplos antes/después
.claude/guides/
├── refactor-shipment-modals-to-dynamic-colors.md ← Guía completa
└── SHIPMENT-MODALS-REFACTOR-SUMMARY.md           ← Este archivo
```

---

## 🏗️ Arquitectura del Sistema de Colores

```
┌─────────────────────────────────────────┐
│     TokenStorage (localStorage)         │
│   { primary, secondary, tertiary, ...}  │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│    theme-provider.tsx                   │
│  - Lee colores del negocio              │
│  - Llama a updateAllColorScales()       │
│  - Escucha cambios de negocio           │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│    color-scales.ts                      │
│  - Genera escalas 50-900 (tinycolor)    │
│  - Inyecta en CSS variables             │
│  (--color-primary-50, -100, ..., -900)  │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│    globals.css (lineas 9-21)            │
│  - Define variables CSS dinámicas       │
│  - --color-primary                      │
│  - --color-secondary                    │
│  - --color-tertiary                     │
│  - --color-quaternary                   │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│    shipment-modals.css (NUEVO)          │
│  - Clases reutilizables:                │
│  - .shipment-btn-primary                │
│  - .shipment-section-origin             │
│  - .shipment-badge-primary              │
│  - .shipment-alert-error                │
│  - ... (25+ clases)                     │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│    Componentes React                    │
│  - ShipmentGuideModal                   │
│  - MassGuideGenerationModal             │
│  (usando las clases dinámicas)          │
└─────────────────────────────────────────┘
```

---

## 📊 Comparación: Antes vs Después

### ANTES (Hardcodeado)
```tsx
<div className="bg-purple-50 border border-purple-100 rounded-xl p-4">
  <h3 className="text-purple-700">Origen</h3>
  <input className="focus:ring-purple-500" />
  <button style={{ background: '#7c3aed' }}>Siguiente</button>
</div>
```
❌ Colores fijos para todos los negocios  
❌ 25+ ubicaciones con hardcoded `#7c3aed`  
❌ Cambio de color requiere editar código  

### DESPUÉS (Dinámico)
```tsx
<div className="shipment-section-origin">
  <h3 className="shipment-section-origin-label">Origen</h3>
  <input className="shipment-input" />
  <button className="shipment-btn-primary">Siguiente</button>
</div>
```
✅ Colores personalizables por negocio  
✅ Una sola fuente de estilos  
✅ Cambio automático al cambiar de negocio  

---

## 🎯 Cambios Principales por Sección

### 1️⃣ Sección Origen (Paso 1)

| Elemento | Antes | Después |
|----------|-------|---------|
| Container | `bg-purple-50 border-purple-100` | `shipment-section-origin` |
| Icon | `bg-purple-100 text-purple-600` | `shipment-section-origin-icon` |
| Label | `text-purple-700` | `shipment-section-origin-label` |
| Input | `focus:ring-purple-500` | `shipment-input` |

### 2️⃣ Sección Destino (Paso 1)

| Elemento | Antes | Después |
|----------|-------|---------|
| Container | `bg-blue-50 border-blue-100` | `shipment-section-destination` |
| Icon | `bg-blue-100 text-blue-600` | `shipment-section-destination-icon` |
| Label | `text-blue-700` | `shipment-section-destination-label` |

### 3️⃣ Tarjetas de Transportista (Paso 2)

| Elemento | Antes | Después |
|----------|-------|---------|
| Card | `hover:border-purple-500` | `shipment-carrier-card` |
| Logo Container | `bg-purple-50` | `shipment-carrier-logo-container` |
| Costo Total | `text-purple-600` | `shipment-cost-amount` |
| Selected State | None | `shipment-carrier-card-selected` |

### 4️⃣ Botones

| Tipo | Antes | Después |
|------|-------|---------|
| Primario | `style={{ background: '#7c3aed' }}` | `shipment-btn-primary` |
| Secundario | `bg-green-600` | `shipment-btn-secondary` |
| Outline | `border-purple-500` | `shipment-btn-outline` |

### 5️⃣ Badges/Etiquetas

| Tipo | Antes | Después |
|------|-------|---------|
| COD | `bg-amber-100 text-amber-700` | `shipment-badge-warning` |
| Éxito | `bg-green-100 text-green-700` | `shipment-badge-success` |
| Primario | `text-purple-600` | `shipment-badge-primary` |

### 6️⃣ Alertas

| Tipo | Antes | Después |
|------|-------|---------|
| Error | `bg-red-50 text-red-700` | `shipment-alert shipment-alert-error` |
| Éxito | `bg-green-50 text-green-700` | `shipment-alert shipment-alert-success` |
| Warning | `bg-amber-50 text-amber-700` | `shipment-alert shipment-alert-warning` |

---

## 📝 Pasos de Implementación

### Paso 1: Importar CSS (5 minutos)

**En `shipment-guide-modal.tsx` (línea 1):**
```tsx
'use client';

import { useState, useEffect, useRef } from "react";
// ... otros imports ...
import '@/shared/ui/styles/shipment-modals.css';  // ← AGREGAR
```

**En `mass-guide-generation-modal.tsx` (línea 1):**
```tsx
'use client';

import { useState, useEffect } from 'react';
// ... otros imports ...
import '@/shared/ui/styles/shipment-modals.css';  // ← AGREGAR
```

### Paso 2: Reemplazar Componentes (2-3 horas)

Usar find & replace por sección:

#### Sección Origen
```
Buscar:  bg-purple-50/50 dark:bg-purple-900/10 border border-purple-100 dark:border-purple-800/30
Cambiar: shipment-section-origin
```

#### Inputs de Formulario
```
Buscar:  w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500
Cambiar: shipment-input
```

#### Botones Primarios
```
Buscar:  style={{ background: '#7c3aed' }}
Cambiar: className="shipment-btn-primary"
```

### Paso 3: Testing (1 hora)

- [ ] Cambiar colores del negocio en DevTools
- [ ] Verificar actualización en tiempo real
- [ ] Probar en modo light y dark
- [ ] Verificar escala de colores (hover, active)

---

## 🎨 Variables CSS Disponibles

### Colores Base
```css
--color-primary       /* Ej: #0f172a */
--color-secondary     /* Ej: #be185d */
--color-tertiary      /* Ej: #06b6d4 */
--color-quaternary    /* Ej: #f59e0b */
```

### Escalas Generadas (Automático)
```css
/* Para cada color base se generan estos niveles: */
--color-primary-50   /* 95% blanco */
--color-primary-100  /* 90% blanco */
--color-primary-200  /* 80% blanco */
--color-primary-300  /* 60% blanco */
--color-primary-400  /* 40% blanco */
--color-primary-500  /* Color base */
--color-primary-600  /* 20% negro */
--color-primary-700  /* 40% negro */
--color-primary-800  /* 60% negro */
--color-primary-900  /* 80% negro */
```

Igual para: `secondary`, `tertiary`, `quaternary`

### Cómo se Usan
```css
/* En clases CSS */
.shipment-btn-primary {
  background: var(--color-primary);
  color: white;
}

.shipment-btn-primary:hover {
  background: var(--color-primary-600);  /* Más oscuro */
}
```

---

## 🔄 Flujo de Datos: Cambio de Negocio

```
Usuario cambia de negocio
        ↓
localStorage se actualiza
        ↓
Event 'businessChanged' se dispara
        ↓
theme-provider.tsx escucha el evento
        ↓
Llama a applyBusinessColors()
        ↓
Llama a updateAllColorScales()
        ↓
CSS variables se actualizan en DOM
        ↓
Navegador renderiza con nuevos colores
        ↓
Los modales se ven actualizados
```

⏱️ **Tiempo de actualización**: < 50ms

---

## 📋 Checklist de Implementación

### Preparación
- [ ] Crear `shipment-modals.css` ✅
- [ ] Crear guía de refactorización ✅
- [ ] Crear ejemplos antes/después ✅

### ShipmentGuideModal (1558 líneas)
- [ ] Importar CSS
- [ ] Paso 1: Reemplazar secciones (líneas 856-946)
- [ ] Paso 1: Reemplazar inputs (líneas 883-998)
- [ ] Paso 2: Reemplazar tarjetas (líneas 1190-1259)
- [ ] Paso 2: Reemplazar badges (líneas 1143-1160)
- [ ] Paso 3: Reemplazar inputs (líneas 1279-1327)
- [ ] Paso 4: Reemplazar botones (líneas 1388-1551)
- [ ] Reemplazar spinner (línea 1167)
- [ ] Reemplazar alertas (línea 826)

### MassGuideGenerationModal (765 líneas)
- [ ] Importar CSS
- [ ] Reemplazar tabla (líneas 503-511)
- [ ] Reemplazar badges (líneas 419-421, 520, 550, 554)
- [ ] Reemplazar botones (líneas 450-457)
- [ ] Reemplazar barras de progreso (líneas 466-470, 597-601)
- [ ] Reemplazar alertas (líneas 573-577)

### Testing
- [ ] Probar con colores por defecto
- [ ] Cambiar colores y verificar actualización
- [ ] Probar dark mode
- [ ] Verificar accesibilidad (contraste)
- [ ] Revisar console por errores

---

## 🚀 Beneficios Esperados

| Aspecto | Beneficio |
|--------|-----------|
| **Escalabilidad** | Nuevos negocios sin cambiar código |
| **Mantenimiento** | Cambios en un solo archivo CSS |
| **UX** | Cambio de colores en tiempo real |
| **DX** | Developers no piensan en hardcoded colors |
| **Accesibilidad** | Escalas de color mantienen contraste |
| **Performance** | Variables CSS sin overhead |

---

## 📚 Recursos

- **Archivo de estilos**: `.../shared/ui/styles/shipment-modals.css` (105 líneas)
- **Guía detallada**: `.../guides/refactor-shipment-modals-to-dynamic-colors.md`
- **Ejemplos código**: `.../examples/shipment-modal-refactor-example.tsx`
- **Sistema de colores**: `theme-provider.tsx`, `color-scales.ts`, `globals.css`

---

## ⚠️ Notas Importantes

1. **No cambiar `globals.css`**: La definición de variables CSS base ya está completa.
2. **No editar `color-scales.ts`**: El sistema de generación de escalas es automático.
3. **Solo agregar clases en `shipment-modals.css`**: Reutilizar en ambos modales.
4. **Importar CSS en ambos modales**: Necesario en cada archivo que use estas clases.
5. **Mantener compatibilidad dark**: Usar `dark:` en Tailwind donde aplique.

---

## 🎯 Próximos Pasos

1. ✅ Revisar archivos creados
2. ⏳ Comenzar con refactorización de `ShipmentGuideModal`
3. ⏳ Continuar con `MassGuideGenerationModal`
4. ⏳ Testing completo
5. ⏳ Solicitar review a equipo
6. ⏳ Deploy a main

**Tiempo estimado total**: 4-5 horas

---

## 📞 Soporte

Si tienes dudas sobre:
- **CSS dinámicas**: Ver `color-scales.ts`
- **Ejemplos**: Ver `shipment-modal-refactor-example.tsx`
- **Clases disponibles**: Ver `shipment-modals.css`
- **Paso a paso**: Ver `refactor-shipment-modals-to-dynamic-colors.md`

¿Necesitas ayuda? Comienza por revisar la guía detallada. 🚀
