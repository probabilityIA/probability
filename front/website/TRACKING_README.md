# Módulo de Rastreo (Tracking) - Website

## Descripción
Módulo de rastreo público para el sitio web Astro. Permite a clientes buscar y consultar el estado de sus envíos sin autenticación.

## Archivos Creados

```
src/pages/rastreo.astro                              # Página pública
src/components/TrackingClient.tsx                    # Componente cliente (Preact)
src/components/tracking/
├── TrackingSearchInput.tsx                          # Input de búsqueda
├── TrackingProgressBar.tsx                          # Barra de progreso (5 pasos)
├── TrackingDetails.tsx                              # Tarjetas de información
└── TrackingTimeline.tsx                             # Timeline de eventos
src/types/tracking.ts                                # Tipos de datos
```

## Ruta Pública
- **URL:** `/rastreo`
- **Acceso:** Público (sin autenticación)
- **Framework:** Astro + Preact (client:load)

## Características

✅ Búsqueda por número de tracking
✅ Búsqueda por número de orden
✅ Barra de progreso visual (5 pasos)
✅ Timeline de eventos con fecha y ubicación
✅ Información detallada del envío
✅ Links a transportista y guía
✅ Manejo de errores con mensajes amigables
✅ Carga de historial automático
✅ Responsive design
✅ Optimizado para Preact (ligero)

## Componentes Astro

### rastreo.astro (Página Pública)
```astro
---
import Layout from '../layouts/Layout.astro';
import TrackingClient from '../components/TrackingClient';
---

<Layout title="Rastreo de Envíos | Probability">
  <!-- Hero header -->
  <TrackingClient client:load />
  <!-- Footer -->
</Layout>
```

- **client:load**: Hidrata componente Preact inmediatamente
- **Layout**: Usa el layout default del website
- **Responsive**: Funciona en mobile/tablet/desktop

## Componentes Preact

### TrackingClient (Principal)
- Estado: shipment, history, error, isLoading
- Llamadas HTTP al backend
- Coordinación entre componentes hijo
- Manejo de búsqueda y reset

### TrackingSearchInput
- Input con ícono de búsqueda
- Botón submit con estado loading
- Validación de entrada
- Disabled durante búsqueda

### TrackingProgressBar
- 5 steps con emojis
- Barra de progreso animada
- Colores según estado
- Mensaje personalizado

### TrackingDetails
- Grid 1-2 columnas (responsive)
- 6 tarjetas con información
- Formato de moneda (COP)
- Links a transportista/guía

### TrackingTimeline
- Timeline vertical
- Línea gradiente
- Círculos para eventos
- Datos: fecha, status, descripción, ubicación

## Tipos de Datos

```typescript
type TrackingStatus = 'pending' | 'picked_up' | 'in_transit' | 'out_for_delivery' | 'delivered' | 'failed';

interface TrackingSearchResult {
  tracking_number: string;
  carrier: string;
  status: TrackingStatus;
  client_name?: string;
  destination_address?: string;
  shipping_cost?: number;
  // ... más campos
}

interface TrackingHistory {
  date: string;
  status: string;
  description: string;
  location: string;
}
```

## Endpoints Requeridos

### Búsqueda
```http
GET /api/v1/tracking/search?tracking_number=ABC123
GET /api/v1/tracking/search?order_number=UUID-123
```

### Historial
```http
GET /api/v1/tracking/{trackingNumber}/history
```

## Configuración de Entorno

En `astro.config.mjs` o `.env`:
```
PUBLIC_API_URL=http://localhost:3050/api/v1
```

El componente accede a: `import.meta.env.PUBLIC_API_URL`

## Estados del Envío

| Estado | Emoji | Color | Significado |
|--------|-------|-------|-------------|
| pending | 📦 | Ámbar | Creado/Pendiente |
| picked_up | 🚚 | Azul | Recogido |
| in_transit | 📍 | Azul | En Tránsito |
| out_for_delivery | 🏠 | Índigo | En Reparto |
| delivered | ✅ | Esmeralda | Entregado |
| failed | ⚠️ | Rojo | Fallido |

## Uso en Astro

### Acceder a la página
```
https://tudominio.com/rastreo
```

### Integrar en menú
Agregar link en layout o navbar:
```astro
<a href="/rastreo">Rastreo de Envíos</a>
```

## Flujo de Usuario

1. Accede a `/rastreo`
2. Ve página limpia con input
3. Ingresa tracking o orden
4. Click en "Rastrear"
5. Sistema busca en backend
6. Muestra: progreso + detalles + timeline
7. Puede rastrear otro o ir atrás

## Estilos

- **Framework:** TailwindCSS (mismo que website)
- **Colores:** Gradientes azul/índigo
- **Responsive:** Mobile-first
- **Animaciones:** Fade-in, pulse, transitions
- **Emojis:** Para iconografía (ligero)

## Performance

- **Framework:** Preact (10KB vs 40KB React)
- **Loading:** client:load (carga con página)
- **API Calls:** Directas al backend
- **Caché:** Browser cache HTTP standard
- **Assets:** Mínimos (sin librerías externas)

## Ventajas sobre Central

✅ Público - No requiere login
✅ Ligero - Preact (10KB)
✅ Rápido - Astro SSG
✅ Simple - Solo rastreo
✅ SEO friendly - Astro meta tags
✅ Mejor para usuarios finales
✅ No contamina dashboard admin

## Próximas Mejoras

- Notificaciones por email
- Suscripción a actualizaciones
- Descarga de prueba de entrega
- Chat con soporte
- Historial de búsquedas previas
- Integración con más transportistas

## Troubleshooting

### "No se encontró información del envío"
- Verifica que el tracking es correcto
- Asegúrate que existe en el backend
- Comprueba conectividad API

### Timeline vacío
- Es normal si el transportista no ha reportado
- Los eventos se cargan conforme se actualizan
- Recarga para ver cambios

### Errores CORS
- Asegúrate que `PUBLIC_API_URL` es correcto
- Backend debe permitir CORS desde website
- Verifica configuración de backend

## API Response Example

```json
{
  "success": true,
  "data": {
    "shipment": {
      "tracking_number": "PROB123XYZ",
      "carrier": "Envioclik",
      "status": "in_transit",
      "client_name": "Juan Pérez",
      "destination_address": "Cra 5 #10-20, Bogotá",
      "estimated_delivery": "2026-03-25T18:00:00Z",
      "shipped_at": "2026-03-20T10:00:00Z",
      "shipping_cost": 15000,
      "tracking_url": "https://...",
      "guide_url": "https://..."
    },
    "history": [
      {
        "date": "2026-03-20T15:30:00Z",
        "status": "Paquete recibido",
        "description": "Recibido en plataforma",
        "location": "Bogotá"
      }
    ]
  }
}
```

## Build & Deploy

```bash
# Development
npm run dev

# Build
npm run build

# Preview
npm run preview
```

La página se genera estáticamente durante build y se sirve como SSG.
