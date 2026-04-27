# Plan: Integracion Bold (Pasarela de Pagos)

Documento de diagnostico + plan por fases para terminar la integracion.

---

## 1. Diagnostico (estado actual)

### Backend - integraciones/pay/bold/

| Archivo | Estado |
|---|---|
| bundle.go | OK - inicializa todas las capas, arranca consumer en goroutine |
| app/process_payment.go | OK - flujo: config -> link -> publish |
| domain/entities/bold_payment.go | OK |
| domain/ports/ports.go | OK - 3 interfaces (IBoldClient, IIntegrationRepository, IResponsePublisher) |
| domain/errors/errors.go | OK |
| infra/primary/consumer/bold_consumer.go | OK - consume `pay.bold.requests` |
| infra/secondary/client/bold_client.go | OK - POST `/online/link/v1`, soporta CREDIT_CARD/PSE/NEQUI/BOTON_BANCOLOMBIA |
| infra/secondary/repository/integration_repository.go | OK pero sin cache Redis |
| infra/secondary/queue/response_publisher.go | OK |

### Backend - modules/pay/ (Bold-related)

| Archivo | Estado | Problema |
|---|---|---|
| internal/app/bold_recharge.go | ESQUELETO | Lee `BOLD_IDENTITY_KEY` y `BOLD_SECRET_KEY` de env. Debe leer de DB encriptada como el resto. Sin timeout, sin retry, sin manejo de 404. |
| internal/domain/dtos/bold_dtos.go | OK |
| internal/infra/primary/handlers/bold_handlers.go | OK |
| Rutas: `POST /pay/wallet/bold/signature`, `GET /pay/wallet/bold/status/:id` | OK |

### Frontend / Mobile

- `front/central/src/services/integrations/pay/bold/` - existe carpeta (revisar contenido en fase de UI)
- `mobile/mobile_central/lib/services/integrations/pay/bold/` - existe carpeta + tests

### Lo que NO existe

- Webhook de Bold (recepcion de eventos `SALE_APPROVED`, `SALE_REJECTED`, `VOID_APPROVED`, `VOID_REJECTED`).
- Cache Redis de credenciales descifradas.
- Endpoint para revelar credenciales descifradas al frontend (super admin).
- Migracion / seed que cree `integration_types.code = 'bold_pay'` con su `platform_credentials_encrypted`, `config_schema`, `credentials_schema`, `setup_instructions`, `base_url`, `base_url_test`.
- Tests unitarios e integracion.
- Documentacion (cubierto con README en esta tarea).

---

## 2. Decisiones arquitectonicas

1. **Credenciales**: SIEMPRE en `integration_types.platform_credentials_encrypted` (AES-256-GCM con `ENCRYPTION_KEY`). NUNCA en `.env`. Esto incluye sacar `BOLD_IDENTITY_KEY` y `BOLD_SECRET_KEY` de `bold_recharge.go`.
2. **Cache**: Redis con TTL de 10 min para creds descifradas. Invalidar al actualizar `integration_types`.
3. **Reveal en frontend**: endpoint dedicado en `modules/integrations` o `modules/my-integrations` con guard de super admin que devuelve creds descifradas (no el blob bytea).
4. **Webhook**: handler en `modules/pay/internal/infra/primary/handlers/bold_webhook.go` (no en la integracion). La integracion solo emite/consulta. La logica de actualizacion de orden vive en pay/payments.
5. **Idempotencia webhook**: tabla `bold_webhook_events` (id_evento Bold UUID) para deduplicar reintentos.

---

## 3. Plan por fases

### Fase 1 - Saneamiento (medio dia)

**Objetivo**: alinear modulo `pay/bold_recharge.go` con el patron del proyecto.

- [ ] Mover `BoldGenerateSignature` y `GetBoldStatus` a usar credenciales desde DB (via repositorio replicado en `modules/pay`, solo SELECT).
- [ ] Eliminar lecturas de `BOLD_IDENTITY_KEY` y `BOLD_SECRET_KEY` de env.
- [ ] Anadir timeout (10s) y retry (max 3) al cliente HTTP que consulta status.
- [ ] Validar codigos de respuesta (404 -> not found, 401 -> creds invalidas, 5xx -> reintentar).
- [ ] Unificar base URL: usar `base_url` / `base_url_test` de `integration_types`, no hardcoded.

Criterio de aceptacion: backend arranca sin requerir variables `BOLD_*`. Endpoint `GET /pay/wallet/bold/status/:id` responde 200 con un id valido y 404 con uno invalido.

### Fase 2 - Cache Redis de credenciales (medio dia)

**Objetivo**: reducir lecturas a DB y descifrado en cada request.

- [ ] Introducir wrapper en `IIntegrationRepository.GetBoldConfig` que consulta Redis primero (`cache:bold:config`).
- [ ] TTL 10 min. Invalidacion: al actualizar `integration_types` desde admin (publish event `integration.updated`).
- [ ] Logging de hits/misses.

Criterio: 99% de las solicitudes `pay.bold.requests` resuelven config sin tocar DB.

### Fase 3 - Migracion y seed inicial (1-2 horas)

**Objetivo**: dejar la integracion configurable por la admin UI.

- [ ] Migracion en `back/migration` que crea `integration_types` con:
  - `code = 'bold_pay'`
  - `name = 'Bold'`
  - `category = 'pay'`
  - `base_url = 'https://integrations.api.bold.co'`
  - `base_url_test = 'https://integrations.api.bold.co'` (Bold no separa env por URL, sino por API key)
  - `config_schema` JSON: `{required_fields: [], optional_fields: ["webhook_url"]}`
  - `credentials_schema` JSON: `{required_fields: ["api_key", "secret_key"], optional_fields: ["environment"]}`
  - `setup_instructions` (texto: como obtener llaves desde panel Bold)
  - `image_url` con logo
- [ ] Seed de prueba (sandbox) en migracion DML, eliminada despues de prod.

Criterio: el equipo puede configurar Bold desde la UI de integraciones sin tocar codigo.

### Fase 4 - Reveal de credenciales (medio dia)

**Objetivo**: super admin puede ver creds descifradas en formulario edit.

- [ ] Endpoint backend `GET /integrations/types/:code/credentials/reveal` con guard de super admin.
- [ ] Componente UI en `front/central/src/services/modules/my-integrations/` que pide reveal y muestra los campos en claro temporalmente (con boton "ocultar").
- [ ] Audit log: registrar cada reveal con user_id + timestamp.

Criterio: solo super admin ve creds. Auditable.

### Fase 5 - Webhook Bold (1-2 dias) [CRITICO]

**Objetivo**: confirmar pagos en tiempo real.

- [ ] Modelo `BoldWebhookEvent`: `id` (UUID Bold), `type`, `payload` (jsonb), `processed_at`, `payment_id`.
- [ ] Handler `POST /webhooks/bold`:
  1. Leer body raw + header `X-Bold-Signature`.
  2. Verificar firma `HMAC-SHA256(base64(body), secret_key)` en hex.
  3. Si invalida -> 401 (no reintentar).
  4. Parsear CloudEvents payload.
  5. Verificar idempotencia por `id` (insert con unique constraint, on conflict skip).
  6. Publicar a `pay.bold.webhook.events` para procesamiento asincrono.
  7. Responder 200 inmediatamente (menos de 2s).
- [ ] Consumer en `modules/pay`:
  - Mapear evento -> estado de orden (`SALE_APPROVED` -> paid, `SALE_REJECTED` -> failed, `VOID_*` -> refunded).
  - Actualizar `payment` correspondiente.
  - Notificar (notification_config) si aplica.
- [ ] Endpoint admin `GET /admin/bold/webhooks` para ver eventos recientes (debugging).
- [ ] Job de reconciliacion diario: cruzar `pay.bold.webhook.events` con `payments` para detectar perdidas (Bold reintenta 5 veces, podemos perder despues).

Criterio: pago hecho en sandbox aparece como `paid` en menos de 5s sin polling. Reintento de Bold no duplica el pago.

### Fase 6 - Cliente robusto + retry (medio dia)

**Objetivo**: cliente HTTP a Bold tolerante a fallos.

- [ ] Timeout configurable (default 10s).
- [ ] Retry exponencial (max 3 reintentos) en errores 5xx y network.
- [ ] Circuit breaker para Bold API (si falla 10 veces seguidas, abrir 1 min).
- [ ] Logging estructurado: request_id propio + bold_link_id.
- [ ] Metricas: latencia, % errores por tipo, % timeouts.

Criterio: caer Bold por 30s no tumba el backend. Logs muestran reintentos.

### Fase 7 - Tests (1 dia)

**Objetivo**: cobertura minima.

- [ ] Unit tests para:
  - `process_payment` usecase con mocks de IBoldClient/Repo/Publisher.
  - `bold_client.CreatePaymentLink` con httptest.
  - Verificacion firma webhook (`HMAC-SHA256`).
  - `bold_recharge.GetBoldStatus` con httptest (200, 404, 401, 5xx).
- [ ] Integration test que arranca consumer + publica un mensaje en `pay.bold.requests` y verifica respuesta.
- [ ] E2E en `.claude/testing/pay/bold/` con CU-01 (crear link), CU-02 (recibir webhook), CU-03 (consultar status).

Criterio: cobertura 70%+ en archivos Bold. CI verde.

### Fase 8 - Frontend / Mobile (1-2 dias)

**Objetivo**: UX completa.

- [ ] Front central: revisar `services/integrations/pay/bold/` y completar:
  - Server Action que llama al backend para iniciar pago.
  - Componente Checkout que renderiza el iframe / redirect a Bold.
  - Pagina de retorno (success / fail) que consulta estado.
- [ ] Mobile: integrar el SDK Bold mobile (si existe) o WebView con el checkout URL.
- [ ] Notificacion al usuario via WebSocket cuando llegue el webhook (para refrescar UI sin polling).

Criterio: usuario completa pago de prueba en sandbox desde web y mobile.

### Fase 9 - Refunds y suscripciones (futuro, no bloqueante)

- [ ] Refund total/parcial via API Bold.
- [ ] Suscripciones recurrentes (si Bold lo soporta - validar).

---

## 4. Dependencias entre fases

```
Fase 1 (saneamiento)
    -> Fase 2 (cache)
        -> Fase 4 (reveal)

Fase 3 (migracion seed) --> independiente, ejecutar pronto

Fase 5 (webhook) --> CRITICO, depende de Fase 1 + Fase 3
    -> Fase 6 (retry/circuit breaker, refuerza)
    -> Fase 7 (tests)
    -> Fase 8 (frontend)

Fase 9 --> backlog
```

## 5. Estimacion total

- Fases criticas (1-7): ~5-7 dias dev.
- Frontend / Mobile (8): +2 dias.
- Total para MVP productivo: **~9 dias** de un dev.

## 6. Riesgos

- **Bold cambia firma del webhook**: mitigar con tests del verificador HMAC.
- **Limite de tiempo 2s en webhook**: el procesamiento real va a queue, el handler solo valida + persiste evento + 200.
- **Perdida de webhooks despues de 5 reintentos de Bold**: job de reconciliacion (Fase 5) y consulta manual con `GET /payments/webhook/notifications/{payment_id}`.
- **Credenciales en cache Redis**: si Redis cae, fallback a DB. Si key rotation, invalidar cache.
- **Doble cobro por reintento Bold**: idempotencia por `id` del evento (Fase 5).

## 7. Out of scope (no incluido en este plan)

- Auditoria de seguridad PCI-DSS (no aplica si no almacenamos PAN).
- Integracion con datafonos fisicos Bold (es otra API).
- Reportes financieros (modulo aparte).

## 8. Criterios de aceptacion globales del MVP

- [ ] Generar link de pago desde admin sin variables `.env` de Bold.
- [ ] Pago en sandbox completo (link -> 3DS -> webhook -> orden actualizada).
- [ ] Reintento de Bold no genera duplicados.
- [ ] Frontend muestra estado en tiempo real.
- [ ] Cobertura de tests >= 70%.
- [ ] README + plan documentados (este archivo).
