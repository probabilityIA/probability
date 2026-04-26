# Plan: Refactor del webhook de Envioclick y mapeo de estados

## Objetivo

Desacoplar el webhook de Envioclick del modulo `shipments`. Que cada carrier sea responsable
de su parseo y del mapeo a estados canonicos de Probability. Shipments debe recibir datos
ya normalizados sin saber quien es el carrier. Ampliar el set de estados canonicos y auditar
los webhooks en una tabla reutilizable.

## Decisiones alineadas con el usuario

1. **Constantes de estados canonicos de Probability** viven en
   `back/central/services/integrations/transport/envioclick/internal/domain`
   (cada carrier replica lo que necesita, respetando aislamiento hexagonal).
2. **Estados actuales en DB** (columna `shipments.status VARCHAR(64)`):
   `pending, in_transit, delivered, failed, cancelled`.
   Propuesta de mapeo ampliado (abajo).
3. **Ruta del webhook** se registra en
   `back/central/services/integrations/transport/envioclick/internal/infra/primary/handlers/`
   respetando hexagonal. Shipments pierde la ruta `/webhooks/envioclick`.
4. **Tabla generica de webhooks** (reutilizable) con rolling window de 50 por provider.
5. **Enviame / MiPaquete tendran webhook** a futuro: el patron se disena generico desde ya.

---

## Estados canonicos propuestos para Probability

Hoy hay 5 estados en DB. Envioclick emite 4 familias en `statusStep` + flag `incidence`.
Propongo ampliar a 9 canonicos para aprovechar granularidad de carriers modernos:

| Canonico Probability | Semantica | Equivalente hoy |
|---------------------|-----------|-----------------|
| `pending` | Creada, no recolectada aun | existente |
| `picked_up` | Recolectada en origen (paso inicial) | nuevo, antes era `in_transit` |
| `in_transit` | En viaje entre centros | existente |
| `out_for_delivery` | En reparto final (en distribucion) | nuevo, antes era `in_transit` |
| `delivered` | Entregada exitosamente | existente |
| `on_hold` | Incidencia / novedad reintentable | nuevo, antes era `failed` |
| `failed` | Fallo permanente, no entregable | existente |
| `returned` | Devuelta al remitente | nuevo |
| `cancelled` | Cancelada manualmente por el negocio | existente |

### Mapeo Envioclick -> Probability (case-insensitive + normalizacion sin acentos)

| `statusStep` de Envioclick | Probability |
|--------------------------|-------------|
| `Pendiente`, `Pendiente de Recoleccion` | `pending` |
| `Envio Recolectado`, `Recolectado` | `picked_up` |
| `En Transito` (cualquier variante) | `in_transit` |
| `En Distribucion` | `out_for_delivery` |
| `Entregado`, `Entregada` | `delivered` |
| `Novedad`, `Incidencia` (con `incidence:true` reintentable) | `on_hold` |
| `No Entregado`, `No entregado` | `failed` |
| `Devuelto`, `Devuelta`, `Regresado` | `returned` |
| _cualquier otro_ | `in_transit` (fallback + log.Warn + registro en `webhook_logs`) |

Reglas especiales:
- Si `incidence: true` con `statusStep != "Entregado"` -> `on_hold` (antes: `failed`).
- Si `incidence: true` con `statusStep == "Entregado"` -> `delivered` (incidencia resuelta).
- Normalizar `statusStep` con `strings.ToLower` + remover acentos antes de comparar.

### NO hay tabla de estados en DB

La columna `shipments.status` es `VARCHAR(64)` sin FK. Los valores nuevos no requieren
migracion de schema, solo poblacion natural conforme lleguen webhooks. Se puede ejecutar
una migracion opcional para backfill si se quiere (no incluido en este plan).

---

## Arquitectura objetivo

### Flujo nuevo del webhook

```
Envioclick -> POST /api/v1/integrations/transport/envioclick/webhook
   |
   v
envioclick/internal/infra/primary/handlers/webhook_handler.go
   1. Persiste body + metadata en webhook_logs (async, no-block)
   2. Parsea payload EnvioClickWebhookPayload (DTO en envioclick/domain)
   3. Mapea a ProbabilityShipmentStatus via envioclick/domain mapper
   4. Construye TransportResponseMessage (operation="webhook_update", provider="envioclick")
   5. Publica a transport.responses via queue publisher existente
   6. Responde 200 OK a Envioclick (rapido, <30ms)
   |
   v (RabbitMQ transport.responses)
shipments/internal/infra/primary/queue/consumer/response_consumer.go
   - Nuevo case "webhook_update"
   - Busca shipment por tracking_number o shipment_id
   - Actualiza status, delivered_at, shipped_at con datos ya normalizados
   - SSE publishing como ya hace hoy
```

### Por que async (no HTTP directo)

- Consistente con el patron existente (quote, generate, track, cancel).
- Shipments mantiene solo UN punto de actualizacion via `transport.responses`.
- Si shipments esta caido, el webhook_logs preserva el evento para reintento.
- Permite que futuros carriers (Enviame, MiPaquete) solo publiquen al mismo canal.

---

## Tabla nueva: `webhook_logs` (generica, reutilizable)

Modelo GORM en `back/migration/shared/models/webhook_log.go`:

```go
type WebhookLog struct {
    ID            uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    CreatedAt     time.Time       `gorm:"not null;index"`
    UpdatedAt     time.Time       `gorm:"not null"`

    // Source identification (reutilizable para cualquier webhook futuro)
    Source        string          `gorm:"size:64;not null;index"`    // "envioclick", "enviame", "mipaquete", "shopify", etc.
    EventType     string          `gorm:"size:128;not null;index"`   // "tracking_update", "order_paid", etc.

    // HTTP metadata
    URL           string          `gorm:"size:512;not null"`
    Method        string          `gorm:"size:8;not null;default:POST"`
    Headers       datatypes.JSON  `gorm:"type:jsonb"`
    RequestBody   datatypes.JSON  `gorm:"type:jsonb;not null"`
    RemoteIP      string          `gorm:"size:64"`

    // Procesamiento
    Status        string          `gorm:"size:32;not null;index"`    // "received", "processed", "failed", "ignored"
    ResponseCode  int             `gorm:"not null;default:200"`
    ProcessedAt   *time.Time
    ErrorMessage  *string         `gorm:"type:text"`

    // Relaciones opcionales (no FK dura; solo referencia)
    ShipmentID    *uint           `gorm:"index"`
    BusinessID    *uint           `gorm:"index"`
    CorrelationID *string         `gorm:"size:128;index"`

    // Campos extraidos clave para queries rapidas sin parsear JSON
    TrackingNumber *string        `gorm:"size:128;index"`
    MappedStatus   *string        `gorm:"size:64"`   // estado canonico de Probability resultante
    RawStatus      *string        `gorm:"size:128"`  // estado raw del carrier
}

func (WebhookLog) TableName() string { return "webhook_logs" }
```

### Rolling window de 50 por source

Despues de cada insert, ejecutar en background (goroutine, no bloquea la respuesta):

```sql
DELETE FROM webhook_logs
WHERE id IN (
    SELECT id FROM webhook_logs
    WHERE source = $1
    ORDER BY created_at DESC
    OFFSET 50
);
```

- Mantiene los 50 ultimos **por source** (envioclick, enviame, etc.).
- Se ejecuta async para no afectar el SLA del webhook (<30ms a Envioclick).
- Alternativa: trigger en PostgreSQL (mas opaco, preferimos codigo explicito).

---

## Estructura de archivos resultante

### Envioclick (dueno del mapeo y el webhook)

```
back/central/services/integrations/transport/envioclick/
├── bundle.go                                       # wire webhook handler + routes
└── internal/
    ├── domain/
    │   ├── entities.go                             # existente (API models)
    │   ├── ports.go                                # existente + IWebhookLogRepository
    │   ├── webhook_payload.go                      # NUEVO — EnvioClickWebhookPayload (movido desde shipments)
    │   ├── probability_status.go                   # NUEVO — constantes canonicas Probability
    │   └── status_mapper.go                        # NUEVO — MapToProbabilityStatus (con normalizacion)
    ├── app/
    │   ├── constructor.go                          # existente
    │   ├── operations.go                           # existente
    │   └── webhook_usecase.go                      # NUEVO — parsea, mapea, persiste log, publica msg
    └── infra/
        ├── primary/
        │   ├── consumer/                           # existente
        │   └── handlers/                           # NUEVA carpeta
        │       ├── constructor.go                  # wire del handler
        │       ├── routes.go                       # RegisterRoutes(router)
        │       └── webhook_handler.go              # HTTP handler
        └── secondary/
            ├── client/                             # existente
            ├── queue/
            │   ├── response_publisher.go           # existente
            │   └── webhook_response_publisher.go   # NUEVO (o extender el existente)
            └── repository/                         # NUEVA carpeta
                ├── constructor.go
                └── webhook_log_repository.go       # GORM repo para webhook_logs
```

### Shipments (solo consume el update normalizado)

```
back/central/services/modules/shipments/
└── internal/
    ├── domain/
    │   ├── envioclick_webhook.go                   # ELIMINAR (queda solo tests de regresion movidos)
    │   └── envioclick_webhook_test.go              # ELIMINAR
    └── infra/
        └── primary/
            ├── handlers/
            │   ├── envioclick-webhook.go           # ELIMINAR
            │   └── router.go                       # QUITAR linea webhooks/envioclick
            └── queue/
                └── consumer/
                    └── response_consumer.go        # AGREGAR case "webhook_update"
```

### Migration

```
back/migration/
├── shared/models/
│   └── webhook_log.go                              # NUEVO modelo GORM
└── internal/infra/repository/
    ├── constructor.go                              # AGREGAR migrateWebhookLogs()
    └── migrate_webhook_logs.go                     # NUEVO (AutoMigrate + indices)
```

---

## Contrato del mensaje `webhook_update`

`TransportResponseMessage` ya existe. Solo agregamos una nueva `Operation`:

```go
// Publicado por envioclick, consumido por shipments
{
    "shipment_id": 34081,              // resuelto por envioclick al buscar tracking
    "business_id": 17,                  // resuelto por envioclick desde el shipment
    "provider": "envioclick",
    "operation": "webhook_update",     // NUEVO
    "status": "success",
    "correlation_id": "wh-<uuid>",
    "timestamp": "2026-04-17T01:31:50Z",
    "data": {
        "tracking_number": "034056642049",
        "probability_status": "in_transit",    // estado canonico YA mapeado
        "raw_status": "En Transito",           // para debugging
        "shipped_at": "2026-04-16T15:21:00Z",  // si aplica
        "delivered_at": null,
        "has_incidence": false,
        "event_description": "Paquete en tr nsito a destino",
        "event_timestamp": "2026-04-17T01:31:45Z"
    }
}
```

El `ResponseConsumer` de shipments **solo lee campos normalizados**. Nunca parsea
`raw_status` ni conoce el carrier.

---

## Fases de implementacion

### Fase 0 — Planning (este documento)
- [x] Diagnostico
- [x] Plan escrito

### Fase 1 — Tabla de webhook_logs y modelo
1. Crear `migration/shared/models/webhook_log.go`.
2. Crear `migration/internal/infra/repository/migrate_webhook_logs.go` con AutoMigrate + indices.
3. Registrar en `constructor.go`.
4. Ejecutar `cd back/migration && go run cmd/main.go`.
5. Verificar con `mcp__postgres-probability__query` que la tabla existe con todos los indices.

### Fase 2 — Dominio envioclick: constantes + mapper
1. Crear `envioclick/internal/domain/probability_status.go` con constantes.
2. Crear `envioclick/internal/domain/webhook_payload.go` moviendo `EnvioClickWebhookPayload` desde shipments.
3. Crear `envioclick/internal/domain/status_mapper.go` con:
   - Normalizacion (lowercase + remove accents via `transform.Chain(norm.NFD, ...)`).
   - Switch ampliado con los 9 estados canonicos.
   - Fallback `in_transit` con retorno de flag `isUnknown bool` para logging.
4. Crear `envioclick/internal/domain/status_mapper_test.go` con casos existentes + nuevos.

### Fase 3 — Repositorio webhook_logs en envioclick
1. Crear `envioclick/internal/domain/ports.go` -> agregar `IWebhookLogRepository`.
2. Crear `envioclick/internal/infra/secondary/repository/webhook_log_repository.go`:
   - Metodo `Save(ctx, log)`.
   - Metodo `TrimOldBySource(ctx, source, keepCount=50)` (async post-insert).
3. Si envioclick ya no tiene repo (hoy no tiene), crear carpeta + constructor.

### Fase 4 — UseCase de webhook en envioclick
1. Crear `envioclick/internal/app/webhook_usecase.go`:
   - Input: `EnvioClickWebhookPayload` + request metadata (URL, IP, headers).
   - Output: `TransportResponseMessage` listo para publicar.
   - Acciones: persistir raw log, mapear status, construir mensaje.
2. Actualizar `envioclick/internal/app/constructor.go` para incluir el nuevo usecase.

### Fase 5 — Handler HTTP en envioclick
1. Crear carpeta `envioclick/internal/infra/primary/handlers/`.
2. Crear `constructor.go`, `routes.go`, `webhook_handler.go`.
3. Handler: leer body raw, parsear, llamar usecase, publicar a `transport.responses`, responder 200 OK rapido.
4. Ruta: `POST /api/v1/integrations/transport/envioclick/webhook`.

### Fase 6 — Wire en bundle de envioclick
1. Actualizar `envioclick/bundle.go` para:
   - Recibir `db.IDatabase` y `*gin.RouterGroup` como parametros.
   - Instanciar repo + usecase + handler.
   - Llamar `handler.RegisterRoutes(router)`.
2. Actualizar `integrations/transport/bundle.go` para pasar db + router.
3. Actualizar llamada desde `cmd/internal/server/init.go` o donde se arma transport.

### Fase 7 — Consumer de shipments: nuevo case webhook_update
1. En `shipments/.../response_consumer.go`:
   - Agregar `case "webhook_update"` en el switch.
   - Implementar `handleWebhookUpdate(ctx, response)`:
     - Resolver `ShipmentID` (si viene, usar; si no, buscar por `tracking_number` del Data).
     - Leer `data.probability_status`, `data.delivered_at`, `data.shipped_at`.
     - Update shipment en DB.
     - Publicar SSE de status change.

### Fase 8 — Eliminar legacy de shipments
1. Borrar `shipments/internal/domain/envioclick_webhook.go` y su test.
2. Borrar `shipments/.../handlers/envioclick-webhook.go`.
3. Quitar linea `router.POST("/webhooks/envioclick", ...)` de `handlers/router.go`.
4. Verificar que `shipments/internal/domain/ports.go` no tenga referencias muertas a envioclick.

### Fase 9 — Verificacion
1. `go build ./...` desde `back/central/` — debe compilar.
2. Tests: `go test ./services/integrations/transport/envioclick/... ./services/modules/shipments/...`.
3. Arrancar backend local, simular webhook con `curl`:
   ```bash
   curl -X POST http://localhost:3050/api/v1/integrations/transport/envioclick/webhook \
     -H 'Content-Type: application/json' \
     -d @testdata/envioclick_webhook_sample.json
   ```
4. Verificar en DB:
   - `SELECT * FROM webhook_logs ORDER BY created_at DESC LIMIT 5` (row nueva).
   - `SELECT status FROM shipments WHERE tracking_number = '...'` (actualizado).
5. Verificar en logs:
   - Envioclick: `✅ Webhook processed successfully`.
   - Shipments: `📨 Processing transport response operation=webhook_update`.
6. Rolling window: insertar 51 webhooks y verificar que hay 50.

### Fase 10 — Deploy
1. Commit atomico por fase (preferible) o un commit grande con mensaje descriptivo.
2. Merge a main -> CI/CD despliega via GitHub Actions (~3-4 min).
3. **Actualizar URL del webhook en panel de Envioclick** de `/api/v1/webhooks/envioclick` a `/api/v1/integrations/transport/envioclick/webhook`.
   - Durante la transicion, mantener la ruta vieja funcional o hacer redirect 301 desde shipments.
   - **Riesgo:** si se cambia la URL antes del deploy, webhooks caen. Plan: deploy primero con AMBAS rutas activas, verificar nueva funciona, luego quitar vieja en PR siguiente.

---

## Consideraciones de compatibilidad hacia atras

**Estrategia dual-route durante migracion:**

Durante fases 5-8, mantener ambas rutas activas:
- Vieja: `/api/v1/webhooks/envioclick` (shipments) — sigue funcionando
- Nueva: `/api/v1/integrations/transport/envioclick/webhook` (envioclick) — nueva logica

Esto permite:
1. Desplegar sin romper webhooks en vuelo.
2. Probar la nueva ruta con un subconjunto antes de cambiar config en Envioclick.
3. Rollback simple: solo revertir el cambio de URL en Envioclick.

Un PR posterior (fase 11, no incluida aqui) elimina la ruta vieja tras ~1 semana.

---

## Riesgos identificados

| Riesgo | Mitigacion |
|--------|------------|
| Nuevos estados (`picked_up`, `out_for_delivery`, etc.) rompen UI/dashboard | Auditar frontend antes de poblar. Si la UI solo espera 5 estados, mapear conservadoramente y feature-flag los nuevos |
| `webhook_logs` crece sin control | Rolling window 50 por source + indice por `(source, created_at)` |
| Shipments consumer no encuentra shipment por tracking | Mismo fallback que hoy: responder 200 y log.Warn; webhook_logs queda con `status=ignored` |
| Mensajes perdidos si RabbitMQ cae durante el webhook | `webhook_logs` conserva el body; plan de reprocesamiento manual via script |
| Orden de eventos invertido (carrier manda evento antiguo despues del nuevo) | Guardar `event_timestamp` del carrier y descartar si es menor al `updated_at` del shipment |

---

## Pregunta abierta (no bloqueante)

¿Exponemos en la UI de shipments los nuevos estados (`picked_up`, `out_for_delivery`, `on_hold`, `returned`)? Si si, se requiere un PR adicional en frontend para:
- Badge/color por estado.
- Filtros en listado.
- Traducciones es/en.

Sugerencia: hacer un PR de backend primero (este plan), poblar DB con estados reales 1-2 dias, luego PR de frontend basado en datos reales.

---

## Checklist de entrega

- [ ] Migracion de `webhook_logs` aplicada en DB local y verificada via MCP
- [ ] Constantes canonicas en envioclick/domain con godoc
- [ ] Mapper ampliado con tests (casos existentes + 4 nuevos estados)
- [ ] Handler HTTP en envioclick responde <30ms
- [ ] webhook_logs recibe y retiene max 50 por source
- [ ] ResponseConsumer de shipments procesa operation=webhook_update
- [ ] Ruta vieja y nueva coexisten hasta deploy
- [ ] `go build ./...` limpio
- [ ] Tests pasan (`go test ./...`)
- [ ] Prueba end-to-end local con curl exitosa
- [ ] Plan compartido con el usuario para validacion antes de ejecutar
