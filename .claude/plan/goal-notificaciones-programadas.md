# TKT-000065 - Notificaciones programadas + plantillas WhatsApp propias

Estado del trabajo. Leer al inicio de cada vuelta del loop; actualizar al final.

## Decisiones cerradas con el usuario

- Plantillas van DIRECTO a Meta, sin aprobacion previa de super admin, incluso
  en el numero compartido de Probability.
- Mitigaciones que se construyen igual: tope de envios por negocio/dia,
  `accepts_marketing` con opt-out, alarma si el quality rating baja de GREEN.
- Alcance de la entrega: segmento `customers_inactive` + modulo de plantillas.

## Hechos del terreno (ya verificados)

- El webhook `message_template_status_update` YA se consume:
  `whatsapp/internal/infra/primary/queue/consumerwebhook/consumer.go:76`.
- `CreateTemplate` y `ListTemplates` YA existen:
  `whatsapp/internal/infra/secondary/client/templates_client.go:143` y `:60`.
- Snapshot de plantillas en Redis con TTL 6h: `cache/templates_cache.go:17`.
- Segmentacion sale de `customer_summary.last_order_at` + `client.phone`.
- `client` NO tiene `accepts_marketing`: lo agrega esta migracion.
- La base local estaba atras de prod; se le aplico a mano el DDL de `sprints`
  y `tickets.sprint_id` (prod ya los tenia, no se toco).
- `migrateSprints` sin registrar en `constructor.go` NO es un bug: hay ~55
  migraciones igual. La convencion del repo es sacarlas de la cadena una vez
  aplicadas en produccion. No tocar.
- TRAMPA: `back/migration/.env` apuntaba al RDS de **produccion**, y
  `dev-db-switch.sh` solo cambia `back/central/.env`. El flujo documentado
  `cd back/migration && go run cmd/main.go` corria contra prod. Se apunto a
  local (5434) con respaldo en `back/migration/.env.dbprod`.
- Correr la migracion puntual: wrapper exportado temporal en el paquete
  `repository` + `cmd/tmponly`, y borrar ambos al terminar. `go run cmd/main.go`
  corre TODA la cadena, que la regla prohibe.
- Segmento verificado en local: business 26 tiene 229 clientes con
  `last_order_at` de mas de 30 dias. La consulta une `client` con
  `customer_summary` por `(customer_id, business_id)`.
- Los telefonos vienen sin normalizar (`3026600431` y `+573026600431` para la
  misma persona). MEDIDO en business 26: el segmento da 204 filas sin dedupe y
  111 con dedupe por telefono. Casi la mitad son la misma persona repetida.
  El dedupe va en la consulta (`DISTINCT ON (phone_key)`), y el cooldown se
  evalua por `client_id` O por telefono normalizado, no solo por `client_id`.
- La persistencia de WhatsApp NO vive en el modulo de integracion: este publica
  a una cola y `notification_config` persiste
  (`whatsapp_persistence_consumer`). Las tablas nuevas siguen ese patron y
  viven en `notification_config`.
- `GetStatus` de usecasetemplates devuelve vacio si el negocio NO tiene numero
  propio. Para plantillas propias en el numero compartido hace falta otra ruta
  que use el WABA de plataforma.

## Pasos

- [x] 1. Ticket TKT-000065 creado
- [x] 2. Modelos GORM + migracion de las 4 tablas (corrida en local, idempotente)
- [~] 3. Modulo whatsapp_templates: hecho dominio, ports, 5 repositorios,
      DTOs, usecase completo (`internal/app/templates/`) con validacion y
      9 tests en verde. Colas nuevas declaradas en `shared/rabbitmq/queues.go`.
      Hechos tambien: handlers HTTP (7 rutas en `/whatsapp-templates`),
      publisher a `whatsapp.templates.submit.requests`, consumer de
      `whatsapp.templates.submit.results`, y todo cableado en `bundle.go`.
      Circuito con Meta cerrado: `consumertemplates` toma la cola de requests
      y llama `CreateTemplate`; `SubmitCustom` resuelve el WABA (propio si
      `OwnNumber`, si no el de plataforma); el webhook publica el estado a la
      cola de results y `template_result_consumer` enruta por `template_id`
      (>0 = resultado de envio, 0 = estado del webhook por waba+nombre).
      PENDIENTE DE VERIFICAR: las rutas no se probaron con el backend
      levantado porque el loop tiene prohibido reiniciar servicios.
- [~] 4. Segmento customers_inactive: hecho el usecase `internal/app/scheduled/`
      (CRUD de reglas + motor `RunDueRules`/`RunRuleNow`, ventana horaria por
      timezone, tope diario, tope por lote, cooldown, idempotencia por
      `ReserveSend`). Hechos tambien: publisher a `notification.scheduled.sends`,
      worker con ticker de 5 min, 7 handlers en `/scheduled-notification-rules`,
      `consumerscheduled` del lado whatsapp que envia, y
      `scheduled_result_consumer` que guarda el resultado. Todo en bundles.

  HALLAZGO: `SendTemplate` de usecasemessaging solo sirve para plantillas del
  catalogo hardcodeado (`entities.Templates`); una plantilla propia del negocio
  daba `ErrTemplateNotFound`. Se agrego `SendCustomTemplate`, que arma el
  mensaje por posicion sin pasar por el catalogo.

  BUG DE DATOS ENCONTRADO Y CORREGIDO: 48 de 229 filas de `customer_summary`
  del business 26 tienen `last_order_at` en fecha cero (año 1), no NULL. Sin
  filtro daban `days_inactive = 739866` y esa gente recibiria "hace 739866
  dias que no compras". La consulta ahora exige
  `cs.last_order_at > DATE '2000-01-01'`. Con el filtro: 102 candidatos
  reales, maximo 218 dias. Vale revisar por que se escriben esas filas.
- [~] 5. UI: hechos `domain/scheduled-types.ts` y las server actions
      (`infra/actions/whatsapp-templates.ts` y `scheduled-rules.ts`).
      `npx tsc --noEmit` sin errores en los archivos nuevos (los que salen son
      preexistentes y todos en `.test.ts`).
      Hechos los 3 componentes: `TemplateBuilder.tsx`, `ScheduledRuleForm.tsx`
      y `ScheduledRulesSection.tsx`. tsc limpio y la skill `ortografia-front`
      no reporta hallazgos.
      Colgado del modal `IntegrationRulesForm` como bloque desplegable
      "Envios programados por segmento". tsc limpio.
      OJO: no correr `next build` (satura la RAM); usar `npx tsc --noEmit`.
- [~] 6. E2E corrido en local sobre Demo (26). Verificado: catalogo de
      variables, las 3 validaciones, creacion real en Meta
      (meta_template_id=1564011728385611, WABA 1302830408357767), creacion de
      regla, guard de plantilla no aprobada, reprogramacion, historial de
      corridas y aislamiento multi-tenant. Ticket comentado, estado `testing`.

  BUG CORREGIDO EN LA PRUEBA: las 4 colas nuevas no existian en RabbitMQ
  ("NOT_FOUND - no queue"). En este repo cada consumer declara su cola antes
  de consumir; faltaba. Se agrego `DeclareQueue` en los 4.

  MISMO BUG PREEXISTENTE, NO TOCADO: `notification.delivery.results` (consumer
  de email) falla igual al arrancar. Fuera del alcance del ticket.

  GAP NO CORREGIDO: `DELETE /whatsapp-templates/:id` solo borra la fila local,
  NO borra la plantilla en Meta. Sigue ocupando cupo en el WABA.

  PENDIENTE: el envio real depende de que Meta apruebe la plantilla (sigue en
  PENDING). No se forzo `approved` en la base a proposito. Borrar
  `prueba_e2e_probability` del WABA cuando se cierre ese tramo.

## Reglas de esta corrida

- Local siempre (`127.0.0.1:5434`). Verificar con `dev-db-switch.sh status`.
- Nada de push, nada de escribir en prod, nada de reiniciar servicios (loop).
- Sin comentarios en Go/TS. Archivos 500+ lineas: cero non-ASCII.


## Cola de pendientes (fuera del alcance original del ticket)

- [ ] Ocultar las pestanas "Canales" y "Tipos de Eventos" del Centro de
      Notificaciones a los usuarios de negocio; solo visibles para super admin.
      Son catalogos de plataforma, un negocio no deberia verlos ni tocarlos.
      Archivo: `front/central/src/app/(auth)/notification-config/` (barra de
      subnavegacion) o el layout que la pinta.
- [ ] Aviso de ventana de 24h cerrada en el chat con selector de plantilla:
      hoy avisa pero igual deja escribir algo que Meta va a rechazar.
- [ ] Alarma si el quality rating del numero compartido baja de GREEN.
- [ ] Investigar las 48 filas de `customer_summary` con `last_order_at` en
      fecha cero (ano 1) en el business 26.
- [ ] `embedded_signup_config_id` real en las platform credentials de
      produccion (el valor de la copia local era relleno: 1234567890123456).
- [ ] `dev-db-switch.sh prod` deja `DB_HOST` en el RDS directo, que no es
      alcanzable desde 2026-08-21. Deberia apuntar al tunel 127.0.0.1:5433.
      Ademas solo cambia `back/central/.env`, no el de `back/migration`.
- [ ] Selector de negocios del header trae solo los primeros 20: en produccion
      un negocio fuera de esa pagina se muestra como "Negocio #26".
- [ ] Centralizar el mapeo evento -> plantilla: hoy esta repetido en
      `whatsapp_publisher.go`, `order_consumer.go`, `shipment_consumer.go` y
      `preview.go`. El campo `template_config` de `notification_event_types`
      seria el lugar natural.

## IA de ventas: DESCONECTADA (2026-09-07)

Por decision del usuario, no se usa por ahora.

- `ai_sales_enabled = false` en las platform credentials de WhatsApp
  (integration_types id 2) en PRODUCCION. Es el switch: `ai_forwarder.go` lo
  lee de Redis y no reenvia nada si esta en false.
- En el chat se quito el bloqueo del compositor y el toggle "IA Activa /
  Pausada". Ahora siempre se puede escribir; el unico limite real es la
  ventana de 24h de Meta.
- Para reactivarla: poner el flag en true otra vez. El codigo del modulo
  `ai_sales` y las colas siguen intactos.
