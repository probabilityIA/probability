# Resumen Módulo: Vista de Saldos (Wallet)

## 📍 Ubicación
- **Path**: `/wallet/saldos`
- **Componente Principal**: `WalletSaldosPage`
- **Archivo**: `front/central/src/app/(auth)/wallet/saldos/page.tsx`
- **Vista Detallada**: `front/central/src/app/(auth)/wallet/wallet-views.tsx`

---

## 🏗️ Estructura General

El módulo tiene **DOS vistas principales** dependiendo del rol del usuario:

### 1. **Vista Admin (Super Admin)**
Para administradores que ven saldos de TODOS los negocios del sistema.

#### Componentes:
- **Tabla "Saldos de Negocios"**
  - Columnas: Negocio, Saldo, Acciones
  - Muestra: Nombre del negocio, saldo en $ (verde/rojo según disponibilidad)
  - Acciones por fila: "Agregar Saldo", "Borrar Historial"

- **Acordeón "En revisión"**
  - Estado: Cerrado por defecto
  - Contiene: Tabla de solicitudes de recarga PENDIENTES
  - Columnas: Fecha, Negocio, Monto, Acciones
  - Acciones: Botones "Aprobar" y "Rechazar"

- **Tabla "Aprobados"**
  - Columnas: Fecha, Negocio, Monto
  - **Filtros de fecha**: Desde/Hasta (a la derecha del título)
  - Búsqueda por rango de fechas de recargas completadas

- **Tabla "Rechazados"**
  - Columnas: Fecha, Negocio, Monto
  - **Filtros de fecha**: Desde/Hasta (igual que Aprobados)
  - Muestra solicitudes rechazadas

---

### 2. **Vista Business (Usuario Normal)**
Para usuarios normales que ven su propia billetera.

#### Componentes:
- **Tarjeta Virtual**
  - Diseño: Card con gradiente púrpura
  - Muestra: Saldo disponible en grande, nombre del negocio, últimas 4 dígitos de wallet

- **Panel "Recargar Saldo"**
  - Input de monto con helper text "Mínimo $15.000"
  - Botones de montos rápidos: $15k, $50k, $100k, $200k, $500k
  - Alert de advertencia cuando hay monto ingresado
  - Botón "Proceder al Pago"

- **Modal de Pago (QR Nequi)**
  - Muestra: QR de Nequi, monto a pagar, instrucciones
  - Estados: "Escanea para Pagar"
  - Botones: "Ya generé el pago", "Regresar"

- **Modal de Procesamiento Bold**
  - Estados: Waiting (spinner), Success, Failed, Timeout
  - Muestra: Monto, Order ID, Tiempo restante
  - Permite cerrar en segundo plano
  - Auto-refresh con polling cada 5s

- **Historial de Transacciones** (2 tablas)
  - **Transacciones Recientes/Procesadas**
    - Columnas: Fecha, Referencia, Método, Monto, Estado
    - Métodos: Nequi (badge), Bold (badge), Débito manual
  - **Transacciones Pendientes**
    - Misma estructura
    - Filtra por estado PENDING

---

## 🎨 Componentes Visuales Principales

### Colores & Estilos
- **Tarjeta Virtual**: Gradiente púrpura (primary color)
- **Botones CTA**: Tertiary color (naranja/rojo)
- **Alertas**: 
  - Success: Verde (#16a34a)
  - Warning: Amarillo (#fef3c7)
  - Error: Rojo (#dc2626)
- **Badges**: Color según método (Nequi=quaternary, Bold=tertiary, Manual=gris)

### Estados de Carga
- Spinner en tablas mientras cargan datos
- Loading en botones con spinner

---

## 📋 Modales Actuales

1. **Modal QR Nequi** - Muestra código QR para pago
2. **Modal Confirmación Pago** - "¡Pago Reportado!"
3. **Modal Selector de Métodos** - Elige entre Nequi, Bold, Otros (coming soon)
4. **Modal Coming Soon** - Para métodos no disponibles yet
5. **Modal Procesamiento Bold** - Muestra estado del pago con Bold

---

## 🔄 Flujos Principales

### Flujo de Recarga (Usuario Normal)
1. Usuario ingresa monto (mín $15k)
2. Elige método de pago (Nequi/Bold)
3. Si Nequi: Ve QR → Escanea → Paga → Confirma
4. Si Bold: Se abre checkout → Paga → Modal muestra progreso
5. Historial se actualiza automáticamente

### Flujo de Aprobación (Admin)
1. Admin ve solicitudes "En revisión" (acordeón)
2. Puede Aprobar o Rechazar cada una
3. Solicitudes se mueven a "Aprobados" o "Rechazados"
4. Puede filtrar por fechas

---

## 📊 Datos Mostrados

### Por Transacción:
- **CreatedAt**: Fecha y hora
- **Reference**: Referencia/motivo
- **Amount**: Monto en $
- **Status**: PENDING, COMPLETED, FAILED
- **Type**: RECHARGE, USAGE
- **IntegrationName**: Nequi, Bold, Débito manual
- **BusinessID**: ID del negocio (ahora incluido en respuestas)
- **WalletID**: ID de la billetera

### Por Wallet:
- **ID**: UUID de la billetera
- **BusinessID**: ID del negocio
- **Balance**: Saldo disponible en $
- **CreatedAt/UpdatedAt**: Fechas

---

## 🎯 Funcionalidades Clave

✅ Vista diferenciada Admin/Business  
✅ Filtro de fechas en Aprobados/Rechazados  
✅ Paginación en todas las tablas  
✅ Métodos de pago múltiples (Nequi, Bold)  
✅ Estados de transacción visuales  
✅ Historial con búsqueda por estado  
✅ Tarjeta virtual con saldo en tiempo real  
✅ Montos rápidos para recarga  
✅ Avisos y confirmaciones claras  

---

## 🚀 Oportunidades de Rediseño

1. **Visualización de Saldos**: Dashboard con gráficos de tendencia
2. **Historial Mejorado**: Timeline vs tablas tradicionales
3. **Métodos de Pago**: Flujo más visual y directo
4. **Estados**: Animaciones y feedback más inmersivo
5. **Filtros**: Más opciones (por método, estado, negocio)
6. **Alerts Inteligentes**: Notificaciones contextuales
7. **Mobile First**: Adaptación para dispositivos móviles
8. **Accesibilidad**: Mejor navegación por keyboard

---

## 🛠️ Stack Técnico

- **Framework**: Next.js 16 + React 19
- **Styling**: TailwindCSS 4
- **Componentes**: Custom UI components (Spinner, Button, Modal, Alert, Table, Input)
- **Estado**: React hooks (useState, useCallback, useEffect)
- **Acciones**: Server Actions (loginServerAction, getWalletBalanceAction, etc.)
- **Real-time**: SSE para procesamiento de Bold, Polling cada 5s

---

## 📱 Responsive

- Desktop first approach
- Grid 1-2 columnas (tarjeta + recarga en desktop)
- Tablas scroll-friendly
- Inputs full-width en mobile
- Modales responsivos

---

## ✨ Propuesta para Claude Design

Este módulo es ideal para un **rediseño innovador** porque:
- Tiene múltiples vistas y flujos
- Maneja datos financieros (necesita claridad visual)
- Tiene componentes reutilizables
- Es crítico para la UX del usuario
- Puede beneficiarse de animaciones y micro-interacciones
- Requiere accesibilidad y claridad de información

**Próximos pasos**: Comparte fotos/mockups con Claude Design para colaborar en el rediseño manteniendo funcionalidad, mejorando estética y UX.
