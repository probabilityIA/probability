# Sistema de Geozonas - Probability

**Fecha de creacion:** 2026-05-05
**Estado actual:** Funcional en local, pendiente de deploy a prod

## Objetivo del proyecto

Implementar un sistema de geozonas con PostGIS para:

1. **Visualizacion jerarquica** de zonas geograficas en el frontend (drill-down).
2. **Calculos de pertenencia** (point-in-polygon) para asignar entregas a zonas.
3. **Estadisticas por zona** (ej. tasa de entregas exitosas por barrio/municipio para cada carrier).
4. **Geozonas personalizadas por business** (ej. zonas de cobertura propias).

Caso de uso primario: saber en que barrio/municipio se entregan ordenes y calcular tasa de exito por zona y carrier (Servientrega, Interrapidisimo, EnvioClick, etc).

## Arquitectura

### Stack
- **Backend:** Go 1.23 + Gin + GORM + PostgreSQL 17.6 + PostGIS 3.5.1
- **Frontend:** Next.js 16 + React 19 + react-leaflet 5 + OpenStreetMap (gratis)
- **DB:** AWS RDS (`database-1.capmmoe4cw2e.us-east-1.rds.amazonaws.com`)
- **Cache:** Redis (local dev container)
- **Mapas:** OSM via tiles publicos (sin costo)

### Modulo backend
Hexagonal estandar en `back/central/services/modules/geozones/`:
```
internal/
  domain/
    entities/geozone.go
    dtos/geozone.go      (types: country, state, city, admin_district, locality, neighborhood, barrio, custom)
    dtos/display.go      (DisplayParams, Bbox, FeatureCollection)
    ports/ports.go       (IRepository, IDisplayCache, IUseCase)
    errors/errors.go     (en español: "geozona no encontrada", etc)
  app/
    constructor.go
    usecases.go          (Create, BulkImport, Get, List, Lookup, Delete)
    display.go           (GetForDisplay con tolerance + cache)
  infra/
    primary/handlers/    (Display, Lookup, List, Get, Create, Bulk, Delete, FlushDisplayCache)
    secondary/
      repository/        (queries con SQL crudo + GORM, ST_Contains, ST_SimplifyPreserveTopology)
      cache/             (Redis con gzip + base64, TTL 30 dias)
```

### Modulo frontend
`front/central/src/services/modules/geozones/`:
```
domain/types.ts          (GeozoneType, DrillState, etc)
domain/ports.ts
app/use-cases.ts
infra/repository/api-repository.ts
infra/actions/index.ts   (Server Actions con session_token)
ui/components/
  GeozoneManager.tsx     (container, drill-down state)
  GeozoneMap.tsx         (Leaflet OSM, FitBounds, FocusOnSelected)
  GeozoneDrawMap.tsx     (dibujar poligono custom)
  GeozoneList.tsx        (lista lateral)
  GeozoneForm.tsx        (modal crear)
  TypeChip.tsx
```

Pagina: `src/app/(auth)/delivery/geozones/page.tsx`
Subnavbar de Ultima Milla: tab "Geozonas 📍"

## Endpoints

| Endpoint | Funcion |
|----------|---------|
| `POST /api/v1/geozones` | Crear desde GeoJSON |
| `POST /api/v1/geozones/bulk` | Importar FeatureCollection |
| `GET /api/v1/geozones?type=&parent_id=&search=` | Listar paginado |
| `GET /api/v1/geozones/:id?include_geometry=` | Obtener una |
| `GET /api/v1/geozones/lookup?lat=&lng=&type=` | Point-in-polygon |
| `GET /api/v1/geozones/display?type=&parent_id=&zoom=&bbox=` | Display optimizado con cache Redis |
| `POST /api/v1/geozones/display/flush-cache` | Invalidar cache |
| `DELETE /api/v1/geozones/:id` | Soft delete |
| `GET /api/v1/shipments/stats/by-geozone?carrier=&type=&from=&to=` | Stats entregas |

Todos protegidos con JWT.

## Datos cargados en RDS prod

Verificado 2026-05-05:

| Tipo | Cantidad | Fuente | Notas |
|------|----------|--------|-------|
| `country` | 1 | DANE 2025 union | Colombia |
| `state` | 33 | DANE MGN 2025 oficial | Departamentos |
| `city` | 1.121 | DANE MGN 2025 oficial | Municipios |
| `admin_district` | 20 | datosabiertos.bogota.gov.co (SDHT) | Localidades politicas Bogota (Chapinero, Kennedy, etc) |
| `locality` | 7.293 | DANE MGN 2025 (URB_ZONA_URBANA) | Corregimientos rurales nacional |
| `neighborhood` | 112 | catastrobogota arcgis (UPZ) | UPZ Bogota POT 190 |
| `barrio` | 1.216 | catastrobogota arcgis (sectorcatastral) | Sectores catastrales Bogota |
| **TOTAL** | **9.796** | | |

Calidad: avg 5.524 puntos por municipio (DANE 2025 vs 312 de la version "_web" inicial). Bogota: 24.420 puntos.

### Jerarquia drill-down (5 niveles para Bogota)
```
Colombia (country)
  Bogota DC (state)
    Bogota DC (city)
      [20 admin_district] + [10 locality corregimientos rurales Sumapaz]
        Chapinero, Kennedy, Suba, etc → tienen UPZ asignados via PostGIS ST_Contains
          [112 neighborhood UPZ] (asignados a su admin_district)
            [1.216 barrio] (asignados a su UPZ via ST_Intersects + mayor area)
```

Los corregimientos rurales (locality) solo aplican para municipios fuera de Bogota DC; en Bogota tambien hay 10 que cubren la zona del Sumapaz.

## Cache strategy

### Redis (backend)
- **Solo data DANE/global** (`business_id=0`) se cachea.
- Key: `geozones:display:{type}:{parent_id}:{zoom_bucket}:{version}`
- Maximo ~50 keys totales (por nivel × parent unico)
- gzip + base64
- TTL 30 dias
- Memoria total: ~10 MB
- **NO se duplica por business** — todos los businesses comparten el cache de zonas oficiales.

### Browser cache
- `Cache-Control: no-cache` + ETag (revalidacion siempre)
- Bumping version (v1 → v7) invalida automaticamente

### HTTP Route Handler (Next.js)
- `app/api/geozones-display/route.ts`
- **Streamea** el response del backend al browser sin reserializar (clave para payloads grandes)
- Server Actions reservan para CRUD chico solamente

## Tolerancias de simplificacion (PostGIS)

```go
func toleranceForZoom(zoom int) float64 {
    switch {
    case zoom <= 7:  return 0.0005   // ~55m, vista pais
    case zoom <= 9:  return 0.0001   // ~11m, vista departamento
    }
    return 0                          // full precision (zoom >= 10)
}
```

- Country: ~2 MB
- State (Antioquia): ~2.9 MB
- City (Bogota localidades): ~1.5 MB
- UPZ Kennedy: ~300 KB
- Barrios: <100 KB

## UX drill-down

**Flujo:**
1. Vista inicial: 33 deptos pintados.
2. Click en depto → mapa zoom + carga sus municipios.
3. Click en municipio → carga sus zonas (admin_districts + corregimientos).
4. Click en localidad → carga sus UPZ.
5. Click en UPZ → carga sus barrios.
6. Click en barrio → solo selecciona/zoom.

**UI:**
- Breadcrumb clickeable: `🏠 Colombia > BOGOTA DC > BOGOTA DC > KENNEDY > PATIO BONITO > [barrio]`
- Stat cards: cantidad visible, payload, nivel
- Lista lateral con boton "← Atras"
- Click en lista hace zoom directo al item (FocusOnSelected)
- Toggle "Mostrar mis zonas personalizadas"
- Boton "+ Nueva geozona" con modal (dibujar en mapa o pegar GeoJSON)

**Colores por tipo en mapa:**
- country: sky-500
- state: violet-500
- city: emerald-500
- admin_district: indigo-500
- locality: amber-500
- neighborhood (UPZ): red-500
- barrio: red-600
- custom: pink-500

## Decisiones importantes

1. **PostGIS habilitado en RDS** sin downtime via `CREATE EXTENSION postgis`.
2. **Datos DANE 2025** elegidos sobre GADM (mejor calidad: 5.524 vs 192 puntos avg).
3. **Misma tabla `geozones`** para todos los niveles (no tablas separadas).
4. **Tipo `admin_district` separado** de `locality` para evitar colision (Sumapaz politica vs corregimiento).
5. **Drill-down en lugar de auto-LOD por zoom** — selectividad jerarquica por click del usuario.
6. **Bbox filter** disponible pero no se usa en drill-down (cada nivel ya esta acotado por parent).
7. **5° nivel de drill** activado para Bogota (admin_district → neighborhood → barrio).
8. **Stacking context aislado** en el mapa (`isolation: isolate`) para que Leaflet no se cuele sobre modales.

## Pendientes

### Funcionalidad
- [ ] Stats por geozona en el frontend (endpoint backend ya existe `/shipments/stats/by-geozone`).
- [ ] Crear/editar/eliminar geozonas custom (modal ya existe pero falta validar end-to-end).
- [ ] Comunas + barrios de Medellin (Geomedellin).
- [ ] Comunas + barrios de Cali (IDESC).
- [ ] Otras ciudades on-demand.
- [ ] Configuracion `parent_id` automatico para zonas custom (ej. cuando creas dentro de Kennedy, se asocia automaticamente).

### Integracion con shipments
- [x] Migracion: `destination_geozone_id`, `destination_point` agregados a `shipments`.
- [x] Backfill: 9.298 / 9.976 shipments con `shipping_lat/lng` ya tienen geozone resuelto.
- [x] Hook en `CreateShipment` use case para resolver geozone al crear.
- [ ] Re-correr backfill para asignar a UPZ/barrio (no solo city).

### Performance / produccion
- [ ] Validar que en prod (mismo VPC RDS) el cache miss baje de 15s a 3-5s.
- [ ] Considerar `pg_tileserv` o vector tiles si llegamos a millones de barrios.
- [ ] Bbox filtering activable a nivel pais cuando carguemos barrios masivos.

### UI / UX
- [ ] Heatmap por zona segun tasa de exito de entregas (combinar con `/shipments/stats/by-geozone`).
- [ ] Filtros adicionales (por carrier, periodo) en la pagina geozonas.
- [ ] Mostrar count de envios por zona en el listado.

## Comandos utiles

### Verificar BD
```sql
SELECT type, COUNT(*) FROM geozones WHERE business_id=0 AND deleted_at IS NULL GROUP BY type;
```

### Probar lookup
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:3050/api/v1/geozones/lookup?lat=4.7110&lng=-74.0721&business_id=1"
```

### Flush cache Redis
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:3050/api/v1/geozones/display/flush-cache
```

### Stats por ciudad
```sql
SELECT g.name, g.code,
  COUNT(*) AS total,
  COUNT(*) FILTER (WHERE s.status='delivered') AS delivered,
  ROUND(100.0 * COUNT(*) FILTER (WHERE s.status='delivered') / COUNT(*), 2) AS pct
FROM shipments s JOIN geozones g ON g.id = s.destination_geozone_id
WHERE g.type='city'
GROUP BY g.id ORDER BY total DESC LIMIT 10;
```

## Resultados reales medidos

Top ciudades por volumen de envios + tasa de exito:
- Bogota: 3.692 envios, 93.23% exito
- Medellin: 1.233 envios, 93.59% exito
- Cali: 598 envios, 93.48% exito
- Itagui: 96.03% exito

## Archivos clave

### Backend
- `back/central/services/modules/geozones/` (modulo completo)
- `back/migration/internal/infra/repository/sql/create_geozones.sql`
- `back/migration/internal/infra/repository/sql/add_shipments_geozone.sql`
- `back/migration/cmd/seed-geozones/` (DANE deptos+munis basico)
- `back/migration/cmd/upgrade-geozones-dane2025/` (datos oficiales DANE)
- `back/migration/cmd/seed-corregimientos/` (locality DANE)
- `back/migration/cmd/seed-bogota-localidades/` (admin_district)
- `back/migration/cmd/seed-bogota-upz/` (neighborhood UPZ)
- `back/migration/cmd/seed-bogota-barrios/` (barrio sector catastral)

### Frontend
- `front/central/src/services/modules/geozones/` (modulo completo)
- `front/central/src/app/(auth)/delivery/geozones/page.tsx`
- `front/central/src/app/api/geozones-display/route.ts` (proxy streaming)
- `front/central/src/shared/ui/delivery-subnavbar.tsx` (entrada en menu)

### Testing
- `.claude/testing/geozones/geozones-api-test.html` (tester HTML standalone)

## Bugs resueltos en esta sesion

1. **Fuente de datos web simplificada** → bajamos shapefile original DANE 2025 (18x mas detalle).
2. **Coords swapped en localidades Bogota** (lat/lng invertidas pese a declarar EPSG:4326) → fix con flip Python.
3. **FitBounds con array vacio** marcaba lastKey y nunca refiteaba → ahora solo marca si bounds.isValid().
4. **Browser cache HTTP duro** servia respuestas viejas → cambio a `Cache-Control: no-cache` + revalidacion ETag.
5. **Next.js Server Actions reserializaban 15 MB** → Route Handler con streaming directo.
6. **Cache key mismo entre business** → confirmado correcto (DANE es global).
7. **Leaflet z-index salia sobre modales** → wrapper con `isolation: isolate` + `z-[2000]` al modal.

## Para continuar mañana

1. Probar end-to-end el modal de crear geozona custom (drawing + paste GeoJSON).
2. Conectar el endpoint `/shipments/stats/by-geozone` al frontend con visualizacion (heatmap).
3. Decidir orden de carga de Medellin/Cali/Barranquilla.
4. Re-correr backfill de shipments para asignar geozone_id mas profundo (UPZ/barrio cuando aplique).
