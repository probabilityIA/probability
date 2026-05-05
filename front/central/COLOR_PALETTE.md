# Color Palette Estándar - Probability Frontend

## Sistema de Colores Dinámicos (CSS Variables)

Todos los módulos deben usar estas variables CSS inyectadas por el ThemeProvider:

### Colores Primarios (Primary - Azul)
- `var(--color-primary-50)` → Fondo muy claro
- `var(--color-primary-100)` → Fondo claro
- `var(--color-primary-200)` → Bordes, divisores
- `var(--color-primary-500)` → Color activo, focus rings
- `var(--color-primary-600)` → Texto activo, iconos
- `var(--color-primary-900)` → Texto oscuro

**Casos de uso:**
- Inputs, selects (focus ring)
- Botones primarios
- Tabs activos
- Checkboxes (accentColor)
- Banners informativos
- Campos activos

---

### Colores Secundarios (Secondary - Púrpura)
- `var(--color-secondary-50)` → Fondo muy claro
- `var(--color-secondary-100)` → Fondo claro
- `var(--color-secondary-200)` → Bordes
- `var(--color-secondary-500)` → Color base
- `var(--color-secondary-900)` → Texto oscuro

**Casos de uso:**
- Headers de tabla (headerBackground)
- Badges especiales (scope: platform)
- Información de super admin
- Bordes decorativos

---

### Colores Terciarios (Tertiary - Púrpura claro)
- `var(--color-tertiary-50)` → Fondo muy claro
- `var(--color-tertiary-100)` → Fondo claro
- `var(--color-tertiary-200)` → Bordes
- `var(--color-tertiary-300)` → Hover states
- `var(--color-tertiary-500)` → Color base (botones, chips)

**Casos de uso:**
- Botones de acción (Asignar, Configurar)
- Quick amount chips (wallet)
- "Próximamente" modals
- Estados highlights

---

### Colores Cuaternarios (Quaternary - Rosa/Pink)
- `var(--color-quaternary-50)` → Fondo muy claro
- `var(--color-quaternary-100)` → Fondo claro
- `var(--color-quaternary-200)` → Bordes
- `var(--color-quaternary-900)` → Texto oscuro

**Casos de uso:**
- Badges especiales (integración Nequi, etc.)
- Información adicional
- Estados alternos

---

## Colores Semánticos (Hex Estáticos)

Estos colores NO cambian con la selección de tema:

### Success/Positivo (Verde)
- Fondo claro: `#dcfce7`
- Texto: `#166534`
- Color base: `#16a34a`

**Casos de uso:**
- Badges "Activo", "Completado"
- Checkmarks de éxito
- Botones de confirmación exitosa
- Texto de ingresos/balance positivos

---

### Error/Negativo (Rojo)
- Fondo claro: `#fee2e2`
- Texto: `#991b1b`
- Color base: `#dc2626`

**Casos de uso:**
- Badges "Inactivo", "Rechazado"
- X de error
- Botones de eliminación
- Texto de adeudos/balance negativo
- Borrar historial

---

### Warning/Advertencia (Amarillo)
- Fondo claro: `#fef3c7`
- Texto: `#92400e`
- Color base: `#fbbf24`

**Casos de uso:**
- Badges "Pendiente"
- Estados en proceso
- Alertas informativas

---

## Reglas de Consistencia

✅ **HACER:**
1. Usar `var(--color-primary-*)` para inputs, tabs, botones primarios
2. Usar `var(--color-secondary-*)` para headers de tabla
3. Usar `var(--color-tertiary-*)` para botones de acción (Asignar, Editar, etc.)
4. Usar `var(--color-quaternary-*)` para badges especiales
5. Usar hexadecimales fijos (#16a34a, #dc2626, #fbbf24) para semánticos
6. Usar inline styles `style={{}}` para aplicar variables CSS
7. Usar `onFocus`/`onBlur` para focus rings dinámicos

❌ **EVITAR:**
- Tailwind classes como `text-blue-600`, `bg-green-100`, etc.
- Múltiples tonalidades del mismo color en el mismo módulo
- Hardcoded hex colors que no sean semánticos
- Tailwind `focus:ring-blue-500` - usar boxShadow dinámico

---

## Ejemplo de Implementación

```tsx
// ✅ CORRECTO
<button
  className="px-4 py-2 rounded-lg text-white"
  style={{ backgroundColor: 'var(--color-tertiary-500)' }}
>
  Guardar
</button>

// ✅ CORRECTO - Checkboxes
<input
  type="checkbox"
  className="h-4 w-4 rounded border-gray-300"
  style={{ accentColor: 'var(--color-primary-600)' }}
/>

// ✅ CORRECTO - Focus states
<input
  type="text"
  className="border border-gray-300 rounded"
  onFocus={(e) => (e.target as HTMLInputElement).style.boxShadow = '0 0 0 3px var(--color-primary-500)'}
  onBlur={(e) => (e.target as HTMLInputElement).style.boxShadow = 'none'}
/>

// ✅ CORRECTO - Estados semánticos
<span style={{
  backgroundColor: '#dcfce7',
  color: '#166534'
}}>
  Activo
</span>

// ❌ INCORRECTO
<button className="bg-purple-500 hover:bg-purple-600">Save</button>

// ❌ INCORRECTO
<div className="text-blue-600">Active</div>
```

---

## Módulos Ya Refactorizados

- ✅ Permissions
- ✅ Resources  
- ✅ Roles
- ✅ Users
- ✅ Businesses
- ✅ Wallet (+ Financial Stats, Virtual Card, Bold Modal)

Última actualización: 2026-05-05
