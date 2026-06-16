# Estructura Visual - Módulo Saldos

## 👨‍💼 VISTA ADMIN (Super Admin)

```
┌─────────────────────────────────────────────────────────┐
│  SALDOS DE NEGOCIOS                                     │
├─────────────────────────────────────────────────────────┤
│  Negocio     │  Saldo      │  Acciones                 │
├─────────────────────────────────────────────────────────┤
│  Gylshop     │  $13.237    │  [Agregar] [Borrar]       │
│  Demo        │  $337.602   │  [Agregar] [Borrar]       │
│  vitality    │  $0         │  [Agregar] [Borrar]       │
│  ...         │  ...        │  ...                      │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  EN REVISIÓN ▼ (Acordeón - cerrado por defecto)        │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  APROBADOS                       [Desde ___]  [Hasta ___]│
├─────────────────────────────────────────────────────────┤
│  Fecha      │  Negocio  │  Monto                        │
├─────────────────────────────────────────────────────────┤
│  20/04/2026 │  Demo     │  $50.000                      │
│  21/04/2026 │  vitality │  $100.000                     │
│  ...        │  ...      │  ...                          │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  RECHAZADOS                      [Desde ___]  [Hasta ___]│
├─────────────────────────────────────────────────────────┤
│  Fecha      │  Negocio  │  Monto                        │
├─────────────────────────────────────────────────────────┤
│  22/04/2026 │  Gylshop  │  $75.000                      │
│  ...        │  ...      │  ...                          │
└─────────────────────────────────────────────────────────┘
```

---

## 💳 VISTA BUSINESS (Usuario Normal)

```
┌──────────────────────────────────────────────────────────┐
│                                                          │
│  ┌────────────────────┐  ┌─────────────────────────────┐│
│  │  [Tarjeta Virtual] │  │  RECARGAR SALDO             ││
│  │                    │  │                             ││
│  │  $ 41,562          │  │  Monto a recargar:          ││
│  │  SALDO DISPONIBLE  │  │  [______ Ej: 50000]         ││
│  │                    │  │                             ││
│  │  Probability       │  │  [15k] [50k] [100k]...      ││
│  │  FINTECH           │  │                             ││
│  │  1153              │  │  ⚠️  Debe consignar...       ││
│  └────────────────────┘  │  [PROCEDER AL PAGO]         ││
│                          └─────────────────────────────┘│
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│  HISTORIAL DE TRANSACCIONES                             │
├──────────────────────────────────────────────────────────┤
│  Transacciones Recientes / Procesadas                   │
├──────────────────────────────────────────────────────────┤
│  Fecha    │ Ref │ Método      │ Monto  │ Estado       │
├──────────────────────────────────────────────────────────┤
│  20/04    │ -   │ [Nequi]     │ +50k   │ Completado   │
│  21/04    │ -   │ [Bold]      │ +100k  │ Completado   │
│  22/04    │ -   │ [Manual]    │ -5k    │ Completado   │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│  Transacciones Pendientes                               │
├──────────────────────────────────────────────────────────┤
│  Fecha    │ Ref │ Método      │ Monto  │ Estado       │
├──────────────────────────────────────────────────────────┤
│  23/04    │ -   │ [Bold]      │ +200k  │ Pendiente    │
└──────────────────────────────────────────────────────────┘
```

---

## 🔄 MODALES / POP-UPS

### 1️⃣ Modal QR Nequi
```
╔════════════════════════════════════╗
║      Escanea para Pagar            ║
╠════════════════════════════════════╣
║                                    ║
║        [   QR CODE   ]             ║
║                                    ║
║        $50,000                     ║
║   Total a pagar vía Nequi          ║
║                                    ║
║   📋 Siguientes pasos:             ║
║   • Escanea el código QR           ║
║   • Verifica monto y llave         ║
║   • Realiza el pago                ║
║                                    ║
║  [Ya generé el pago] [Regresar]    ║
╚════════════════════════════════════╝
```

### 2️⃣ Modal Selector de Métodos
```
╔════════════════════════════════════╗
║    Selecciona Método de Pago       ║
╠════════════════════════════════════╣
║  ☑️ Nequi/Bancolombia              ║
║     Escanea QR y paga al instante  ║
║                                    ║
║  □ Bold (Tarjeta de Crédito)       ║
║     Abre checkout seguro           ║
║                                    ║
║  □ Otros Métodos 🚀 Próximamente   ║
║     (Wallet digital, etc)          ║
║                                    ║
║        [SIGUIENTE]                 ║
╚════════════════════════════════════╝
```

### 3️⃣ Modal Procesamiento Bold
```
╔════════════════════════════════════╗
║    Procesando pago...              ║
╠════════════════════════════════════╣
║          [SPINNER]                 ║
║                                    ║
║   Esperando confirmación de Bold.  ║
║                                    ║
║   ⏱️  La confirmación puede        ║
║       tardar hasta 5 minutos       ║
║                                    ║
║   Monto: $100,000                  ║
║   Orden: ABC123XYZ                 ║
║   Tiempo: 45s                      ║
║                                    ║
║  [Seguir esperando en segundo plano]
╚════════════════════════════════════╝
```

### 4️⃣ Modal Confirmación Pago
```
╔════════════════════════════════════╗
║    ✅ ¡Pago Confirmado!            ║
╠════════════════════════════════════╣
║                                    ║
║   Tu billetera fue recargada con   ║
║   $50,000                          ║
║                                    ║
║   Nuevo saldo: $91,562             ║
║                                    ║
║         [CERRAR]                   ║
╚════════════════════════════════════╝
```

---

## 🎨 PALETA DE COLORES ACTUAL

| Elemento | Color | Uso |
|----------|-------|-----|
| Primary | Azul claro | Info, borders |
| Secondary | - | - |
| Tertiary | Naranja/Rojo | CTAs, botones |
| Quaternary | Verde agua | Métodos alternativos |
| Success | Verde (#16a34a) | Pagos completados |
| Warning | Amarillo (#fef3c7) | Advertencias |
| Error | Rojo (#dc2626) | Errores, rechazos |
| Background | Gris claro/Oscuro | Fondos, modales |

---

## 📊 DATOS MOSTRADOS

### Por Transacción:
- Fecha (CreatedAt)
- Referencia
- Método (Nequi, Bold, Manual)
- Monto (+ o -)
- Estado (Pending, Completed, Failed)

### Por Wallet:
- Saldo disponible
- Nombre del negocio
- Últimos 4 dígitos
- Estado (Activa)

---

## 🎯 FUNCIONES CLAVE

✅ Ver saldos de todos los negocios (Admin)  
✅ Filtrar por fechas (Admin)  
✅ Recargar billetera (Usuario)  
✅ Ver historial de transacciones  
✅ Métodos de pago múltiples  
✅ Estados visuales de transacciones  
✅ Montos rápidos  
✅ Modales de confirmación  

---

## 💡 IDEAS PARA REDISEÑO

### Dashboard Analytics
- Gráfico de recargas/uso en el tiempo
- Resumen de tendencias
- Alertas de saldo bajo

### Timeline Historial
- En lugar de tabla, timeline visual
- Eventos agrupados por día
- Micro-animaciones al expandir

### Método de Pago Mejorado
- Animación de transición entre métodos
- Preview del QR antes de procesar
- Contador visual de tiempo en procesamiento Bold

### Filtros Avanzados
- Por método, estado, rango de monto
- Búsqueda por referencia
- Exportar reporte

### Cards Mejoradas
- Glassmorphism para tarjeta virtual
- Datos en cards colapsables
- Indicadores visuales de estado

### Notificaciones
- Toast con sonido para confirmaciones
- Badge contador de pendientes
- Sistema de alertas inteligentes

