# WhatsApp: numero propio por negocio (fases 1 a 5)

**Ticket:** sin ticket. La regla `.claude/rules/tickets.md` pide registrarlo,
pero el modulo de tickets exige JWT de super admin contra un backend corriendo y
el trabajo se hizo en un contenedor remoto sin base ni backend. **Crear el ticket
y pegar ahi este resumen** es el primer pendiente.

**Rama:** `claude/whatsapp-numero-por-cliente-s22bjf` (pusheada)
**Plan de origen:** `.claude/project/whatsapp-numero-por-cliente.md`
**Detalle tecnico:** `back/central/services/integrations/messaging/whatsapp/NUMERO-PROPIO.md`
**Alerta de pendientes:** `.claude/alerts/whatsapp-numero-propio-pendientes.md`

## Resumen

Se implementaron las 5 fases del plan para que cada negocio pueda enviar y
recibir WhatsApp desde SU propio numero, con Probability como administrador de su
cuenta de Meta. Todo compila y los tests pasan, pero **nada se ejecuto contra una
base, un Redis ni la Graph API**: falta correr una migracion y probar en local
antes de desplegar.

## Commits

| Hash | Que |
|---|---|
| `ac8f96a` | Fases 1 a 4 (backend): credenciales por negocio, ruteo del webhook por numero, plantillas por WABA, aislamiento de lo que no debe migrar |
| `3f78a7d` | Fase 5 (UI) + invalidacion del indice viejo al cambiar el `config` |
| `486b343` | Fix: activar la `HumanSession` al enrutar por numero propio |

## Que cambio, por fase

### Fase 1 - credenciales por negocio

`credentials_cache.GetWhatsAppConfig` ya recibia el `businessID` y buscaba la
integracion del negocio, pero **descartaba el resultado** y caia siempre a
`GetWhatsAppDefaultConfig()`. Ahora lee de verdad su `config` + `credentials`.

Como se decide de que numero sale cada mensaje:

| `integrations.config` (tipo 2) | Efecto |
|---|---|
| sin `phone_number_id`, o `use_platform_token: true` | numero de Probability (comportamiento historico) |
| `phone_number_id` + `waba_id` + credencial `access_token` | numero del negocio |

El fallback no es opcional: **hoy todos los negocios estan en la primera fila de
esa tabla**, y ninguno debe notar el cambio.

Si un negocio declara `phone_number_id` pero le falta el token, el envio falla
con error **permanente** (no reintentable): es configuracion rota, no un fallo
transitorio, y reintentar no la arregla.

### Fase 2 - ruteo del webhook por numero

El webhook es uno solo para todos; la firma HMAC usa el `webhook_secret` de la
app, igual para todos, asi que eso no cambia.

Lo nuevo: `phone_number_id -> (integration_id, business_id)`, un SELECT sobre
`integrations` filtrando `config->>'phone_number_id'`, cacheado en Redis en
`integration:idx:cfg:2:phone_number_id:<id>`. Corre por cada mensaje entrante.

Cuando llega un mensaje sin conversacion activa ni sesion humana:

- si el numero es de un negocio: se abre una conversacion `inbound`, se persiste
  el mensaje y se publica por SSE a su dashboard;
- si no es de nadie (o es el numero de la plataforma): sigue el camino de
  siempre, el agente AI Sales.

**Esta fase es la que necesita la migracion** (ver abajo).

### Fase 3 - plantillas en el WABA del cliente

Se implemento **replicando desde nuestro WABA** (`GET /{waba_id}/message_templates`)
hacia el del cliente (`POST /{waba_id}/message_templates`), no reconstruyendo el
cuerpo de cada plantilla en Go.

Razon: varias definiciones de `entities/template.go` (`menu_no_confirmacion`,
`tipo_novedad_pedido`, `handoff_asesor`, ...) **no tienen `Body`**. Una provision
armada desde el codigo habria creado plantillas incompletas o directamente
rechazadas por Meta. Con Meta como origen, el catalogo no se desincroniza.

El estado vive en Redis (`whatsapp:templates:{integration_id}`, TTL 6 h), no en
base: la fuente de verdad es Meta y el modulo es cache-first por diseno. El
webhook `message_template_status_update` actualiza ese cache usando el WABA del
`entry.id`.

Endpoints nuevos:

```
GET  /api/v1/integrations/whatsapp/templates/status      (?refresh=true consulta a Meta)
POST /api/v1/integrations/whatsapp/templates/provision
```

### Fase 4 - lo que NO debe migrar

El plan listaba tres consumidores a aislar: `consumeralert` (alertas de
monitoreo), `consumerauthotp` (OTP de login) y `consumerai`. **Los tres ya usaban
`GetWhatsAppDefaultConfig`**, asi que no habia nada que hacer.

Lo que el plan NO listaba y si era un problema real: `consumerwalletalert` y
`consumersubscriptionalert` mandaban por `SendTemplate(businessID)`. Con el
cambio de la fase 1, el aviso de saldo bajo de billetera y el de ventana de pago
de la suscripcion habrian salido **del numero del propio negocio hacia el propio
negocio**, y con una plantilla (`reporte_saldo_billetera_v2`,
`resumen_pago_suscripcion`) que solo existe en nuestro WABA. Se agrego
`SendPlatformTemplate` y los dos consumidores lo usan.

Criterio: si el mensaje es de Probability HACIA el negocio, sale de nuestro
numero. Si es del negocio HACIA su cliente final (`consumerorder`,
`consumershipment`), sale del numero del negocio.

### Fase 5 - UI

Modulo front nuevo `services/integrations/messages/whatsapp` con `domain`, `app`
e `infra` (repositorio + server actions), mas dos componentes montados en
`WhatsAppIntegrationView`:

- `WhatsAppConnectionForm`: toggle "usar mi propio numero" + `waba_id`,
  `phone_number_id` y token. El token solo se manda si se escribe, asi guardar la
  config no borra el guardado (el backend REEMPLAZA credenciales, no las mezcla).
- `WhatsAppTemplatesPanel`: estado por plantilla, refresco contra Meta y boton de
  aprovisionamiento. Con el numero de la plataforma explica que no hay nada que
  gestionar.

## El bug que aparecio revisando el propio codigo

La conversacion `inbound` se creaba en `StateHandoffToHuman`, que
`conversation_cache.go:62` considera **estado terminal**: `Save` le borra el
indice `active`. Resultado: el segundo mensaje del mismo cliente no encontraba
conversacion activa, volvia a caer en el ruteo por `phone_number_id` y abria otra
conversacion. **Una conversacion por mensaje.**

Corregido en `486b343` activando la `HumanSession` que apunta a esa conversacion:
es la rama que ya existia para el chat del dashboard, y los mensajes siguientes
entran por ahi sobre la misma conversacion, sin pasar por el flujo del bot.

## Hipotesis y alternativas descartadas

- **Repo propio en el modulo WhatsApp para la consulta por `phone_number_id`.**
  El plan lo pedia asi por la regla de aislamiento de repositorios. No se hizo:
  el modulo es cache-first y **no recibe `db.IDatabase`** (`messaging.New` no lo
  tiene, a proposito). Se resolvio extendiendo `IPlatformCredentialsGetter`, que
  es el patron que el modulo ya usaba para hablar con `integrations/core`. No
  viola el aislamiento: no importa el repo de otro modulo, pide por un puerto.
- **Guardar el estado de plantillas en base.** Habria exigido tabla + migracion
  para un dato del que Meta ya es fuente de verdad y que se puede re-consultar.
  Se descarto por Redis con TTL.
- **Reconstruir las plantillas desde `entities/template.go`.** Descartado por lo
  del `Body` faltante (ver fase 3).
- **Caer al numero de la plataforma cuando el negocio tiene `phone_number_id`
  pero le falta el token.** Descartado: enmascara una configuracion rota y hace
  que el negocio crea que envia desde su numero cuando no. Mejor error explicito
  y permanente.

## Migracion: SI hace falta, y es una sola

`migrateWhatsappInboundConversationType`
(`back/migration/internal/infra/repository/migrate_whatsapp_inbound_conversation_type.go`)

Amplia el CHECK de `whatsapp_conversations.conversation_type` para aceptar
`inbound`, ademas de `order` y `system_alert`. No crea tablas, no crea columnas,
no toca datos.

**Sin ella, la fase 2 falla con SQLSTATE 23514 y el mensaje entrante se pierde.**

El resto no necesita migracion: `waba_id`, `phone_number_id` y `access_token` van
en las columnas JSONB `integrations.config` e `integrations.credentials`, que ya
existen; y el estado de plantillas vive en Redis.

## Como ejecutarlo desde tu PC

### 1. Traer la rama

```bash
cd ~/probability
git fetch origin claude/whatsapp-numero-por-cliente-s22bjf
git checkout claude/whatsapp-numero-por-cliente-s22bjf
```

### 2. Apuntar a la base LOCAL (obligatorio antes de migrar)

```bash
./scripts/dev-db-switch.sh status     # confirma a que base apunta hoy
./scripts/dev-db-switch.sh local      # 127.0.0.1:5434, NO 5433 (5433 es prod)
./scripts/dev-services.sh status      # que postgres local este arriba
```

Si la base local no existe o quedo sucia de tanto probar:

```bash
./scripts/aws-tunnel.sh ensure        # solo para clonar estructura desde prod
./scripts/local-db-clone.sh 26        # business Demo
```

### 3. Correr la migracion

```bash
cd back/migration && go run cmd/main.go
```

`Migrate()` ya quedo apuntando SOLO a esta migracion, no arrastra la cadena
historica.

### 4. Verificar que quedo

```sql
SELECT pg_get_constraintdef(oid)
FROM pg_constraint
WHERE conname = 'chk_whatsapp_conversations_conversation_type';
```

Debe decir:

```
CHECK (conversation_type IN ('order', 'system_alert', 'inbound'))
```

### 5. Dejar `Migrate()` en cero otra vez

`back/migration/internal/infra/repository/constructor.go`:

```go
func (r *Repository) Migrate(ctx context.Context) error {
    return nil
}
```

Y marcar como corrida la fila del 2026-09-05 en `back/migration/MIGRACIONES.md`
(hoy dice `pendiente`).

### 6. Cargar el `waba_id` de plataforma

Sin este campo, `POST /whatsapp/templates/provision` no tiene de donde copiar las
plantillas y devuelve error.

En el front, como super admin: Integraciones -> tipo WhatsApp -> credenciales de
plataforma. Agregar `waba_id` = `1302830408357767` (WABA de produccion). En
pruebas, el WABA de test es `946521194991666`.

Verificar que quedo en cache:

```bash
redis-cli GET integration:platform_creds:2
```

### 7. Levantar y probar la fase 1 (la que toca a todos)

```bash
./scripts/dev-services.sh restart backend
./scripts/dev-services.sh logs backend 100
```

Caso A - negocio SIN numero propio (todos los de hoy). Debe seguir saliendo por
el numero de Probability, sin tocarle nada:

```sql
SELECT id, business_id, config FROM integrations WHERE integration_type_id = 2;
```

Manda una plantilla al numero de pruebas desde el panel y confirma en el log que
el `phone_number_id` es el de la plataforma.

Caso B - negocio CON numero propio. Con el numero de test:

1. En el front, integracion WhatsApp del negocio Demo -> "Usar mi propio numero".
2. `waba_id` = `946521194991666`, `phone_number_id` = `921919641007826`, token
   del system user.
3. Guardar, mandar una plantilla y confirmar en el log que el `phone_number_id`
   ahora es el del negocio.

Caso C - ruteo del webhook. Responde desde el celular al numero de test y
confirma:

- en el log, `mensaje enrutado al negocio dueno del numero`;
- en base, una fila nueva en `whatsapp_conversations` con
  `conversation_type = 'inbound'` y el `business_id` correcto;
- que el SEGUNDO mensaje **no** cree otra conversacion (ese era el bug de
  `486b343`).

```sql
SELECT id, phone_number, conversation_type, business_id, current_state, created_at
FROM whatsapp_conversations
WHERE conversation_type = 'inbound'
ORDER BY created_at DESC LIMIT 5;
```

### 8. Volver a dejar el entorno como estaba

```bash
./scripts/dev-db-switch.sh local   # nunca dejar el .env apuntando a prod
```

## Verificacion hecha (y la que no)

Hecho en la sesion:

- `go build ./...` y `go test ./...` de `back/central`: pasan.
- `go build ./...` de `back/migration`: pasa.
- `npx tsc --noEmit` del front: sin errores nuevos (los que quedan son
  preexistentes: `@testing-library/react` y los imports de assets de
  `MapComponent.tsx`).
- `revisar.py` del skill `ortografia-front` sobre el modulo nuevo: 0 hallazgos.

NO hecho, y por eso esta la alerta:

- La migracion no se corrio.
- Ningun envio real, ni por el numero de la plataforma ni por uno propio.
- La fase 3 nunca toco la Graph API: el aprovisionamiento de plantillas y el
  consumo de `message_template_status_update` estan escritos pero sin una sola
  corrida contra Meta.
- El webhook nunca recibio un evento real.

## Pendientes

1. **Correr la migracion** (pasos 2 a 5 de arriba). Bloqueante.
2. **Probar la fase 1 en local** con los casos A y B. Bloqueante: toca a TODOS
   los negocios que hoy envian por nuestro numero.
3. **Cargar el `waba_id` de plataforma** (paso 6). Sin el, la fase 3 no arranca.
4. **Verificar la fase 3 contra un WABA real**, con un cliente conectado o un
   WABA de pruebas compartido con Probabilityapp.
5. **Crear el ticket** y pegar ahi este resumen.
6. **Quality rating por numero.** El riesgo que el plan senala sigue abierto: hoy
   se vigila un `quality_rating` y con N clientes son N. Falta leer
   `GET /{phone_number_id}?fields=quality_rating` y exponerlo en la UI.
7. **Embedded Signup** (camino largo del plan). El codigo de las fases 1 a 5
   sirve igual; lo unico que cambia es que la pantalla de conexion pasa a ser un
   boton. Depende de sacar la app de dev mode y pasar App Review.
