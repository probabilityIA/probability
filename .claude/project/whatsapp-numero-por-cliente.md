# WhatsApp: numero propio por cliente

Fecha: 2026-09-04

Objetivo: que cada negocio pueda enviar y recibir WhatsApp desde SU propio
numero, con Probability como administrador de su cuenta, sin que ningun cliente
tenga que hacer tramites con Meta.

## Estado actual

Todos los negocios envian desde el numero de Probability
(+57 300 5636160, phone_number_id 1077369948787698).

```
integrations (integration_type_id = 2)
 id=2   business_id=NULL  config: {"phone_number_id": "921919641007826"}
 id=61  business_id=26    config: {"use_platform_token": true}
```

Verificado en Graph API el 2026-09-04:

| Campo | Valor |
|---|---|
| business_verification_status | verified |
| account_review_status | APPROVED |
| numero prod | CONNECTED, quality GREEN |
| scopes del token | whatsapp_business_management, whatsapp_business_messaging |
| privacy_policy_url / terms_of_service_url | vacios |
| GET /me/businesses | (#100) Missing Permission |

La cuenta esta verificada. Lo que NO esta es la app: falta `business_management`,
faltan las URLs legales y los permisos estan en Standard Access. Eso solo importa
para Embedded Signup, no para el camino que se describe aqui.

## Dos caminos

### Camino corto (este documento)

El cliente comparte su WABA con nuestro Business desde su Business Manager.
Nuestro system user `cam-adm` queda como administrador de su cuenta.

- Tramites con Meta de nuestro lado: ninguno.
- Trabajo: solo codigo.
- Limite: es manual, el cliente hace unos clics. No sirve para autoservicio.

### Camino largo (Embedded Signup, mas adelante)

Boton en el dashboard, el cliente se conecta solo. Exige sacar la app de dev
mode, publicar politica de privacidad / terminos / eliminacion de datos, y pasar
App Review con Advanced Access mas `business_management`. Meta tarda 1 a 2
semanas si aprueba a la primera, 3 a 4 si rechaza una vez.

**El codigo de las fases 1 a 5 es el mismo para los dos caminos.** Por eso se
arranca por el corto sin perder nada.

## Hallazgos que reducen el trabajo

1. `credentials_cache.go:39` ya recibe `businessID` y busca la integracion del
   negocio, pero descarta el resultado y cae a `GetWhatsAppDefaultConfig()`. El
   hueco esta concentrado en un solo punto.
2. `request/webhook_payload.go:33` ya parsea `phone_number_id`. Nunca se usa para
   rutear, pero no hay que tocar el parser.
3. El core de integraciones ya expone `DecryptCredential`, `GetIntegrationConfig`
   y `UpdateIntegrationCredentials`. No hay que construir cifrado ni storage.

## Fases

### Fase 0 - Onboarding manual (sin codigo)

El cliente: Business Manager -> Configuracion del negocio -> Cuentas de WhatsApp
-> Socios -> agregar el Business ID de Probabilityapp con control total.

Nosotros anotamos `waba_id` y `phone_number_id`, y suscribimos la app:

```
POST /{waba_id}/subscribed_apps
```

### Fase 1 - Credenciales por negocio (1-2 dias)

- `credentials_cache.go:39-56`: que `GetWhatsAppConfig` lea de verdad el `config`
  y las `credentials` de la fila `Integration` del negocio. Con
  `phone_number_id` propio se usa el del cliente; con `use_platform_token: true`
  cae al numero de Probability como hoy.
- `ports.go:127` `IPlatformCredentialsGetter`: agregar lectura de config y
  credenciales por `integrationID`.
- `ports.go:109` `WhatsAppConfig`: agregar `WABAID`.
- Cachear por business en Redis, con la invalidacion que ya usa
  `integration:creds:*`.

El fallback no es opcional: los negocios que hoy usan nuestro numero deben
seguir funcionando sin tocarles nada.

### Fase 2 - Ruteo del webhook por numero (1-2 dias)

- Consulta nueva `phone_number_id -> (business_id, integration_id)`: SELECT sobre
  `integrations` filtrando `config->>'phone_number_id'`. Replicada en el repo del
  modulo, no compartida (regla de aislamiento de repositorios).
- Cachearla en Redis: corre por cada mensaje entrante.
- `handle-webhook.go`: cuando no hay conversacion ni sesion humana en Redis,
  resolver el negocio por `metadata.phone_number_id`.
- La firma HMAC no cambia: el `webhook_secret` es de la app, igual para todos.

### Fase 3 - Plantillas en el WABA del cliente (2-3 dias)

Las 13 plantillas viven en nuestro WABA. Cada cliente necesita las suyas.

- Provision al conectar: `POST /{waba_id}/message_templates` con las 13.
- Guardar estado por plantilla y por WABA. Meta aprueba en horas o dias, y
  mientras tanto no se pueden usar.
- Consumir `message_template_status_update` en el webhook para actualizar ese
  estado.

Es la fase que mas se subestima: introduce el estado "cliente conectado pero sin
plantillas aprobadas", que la UI tiene que saber mostrar.

### Fase 4 - Lo que NO debe migrar (medio dia)

Estos consumidores deben seguir en el numero de Probability:

- `consumeralert/alert_consumer.go:72` - alertas de monitoreo, son nuestras
- `consumerauthotp/otp_consumer.go:57` - OTP de login de la plataforma
- `consumerai/response_consumer.go:58` - revisar caso por caso

Si se mezcla, las alertas del servidor le salen al cliente por su propio numero.

### Fase 5 - UI (1-2 dias)

Pantalla de conexion: `waba_id`, `phone_number_id`, token, boton de probar
conexion (ya existe `usecasetestconnection`) y estado de plantillas. Con Embedded
Signup esta pantalla se reemplaza por un boton.

## Esfuerzo

| Fase | Esfuerzo | Probable sin cliente conectado |
|---|---|---|
| 1 credenciales por negocio | 1-2 dias | si, con el numero de test |
| 2 ruteo del webhook | 1-2 dias | si |
| 3 plantillas por WABA | 2-3 dias | no, necesita un WABA real |
| 4 aislar alertas/OTP | medio dia | si |
| 5 UI | 1-2 dias | si |

Total aproximado: semana y media. Las fases 1, 2 y 4 se verifican completas
contra la base local y el numero de test.

## Riesgos

- **Calidad y ventana de 24h pasan a ser por numero.** Hoy se vigila un
  `quality_rating`; con N clientes son N. Exponerlo en la UI desde el principio.
- **El cliente puede revocar el acceso** desde su Business Manager sin avisar.
  Los envios fallan con 401 / codigo 190: hay que clasificarlo como error
  permanente y no reintentable, o se cae en el bucle caliente descrito en
  `.claude/rules/colas-errores-permanentes.md`.
- **`business_id` NULL vs no NULL.** La fila global (id=2) convive con las de
  cada negocio. Una consulta mal filtrada hace que un negocio envie por el numero
  de otro.

## Lo que tiene que poner el cliente

- Un Business Manager de Meta.
- Un numero: nuevo, o migrado desde WhatsApp Business App, o en modo
  coexistencia. Verificacion por OTP y PIN de 2FA.
- Metodo de pago propio en su WABA.
- Si necesita mas de 250 conversaciones diarias, verificacion de su negocio ante
  Meta (camara de comercio y RUT, 2 a 10 dias habiles). Es tramite del cliente,
  no nuestro.

## Referencias

- `back/central/services/integrations/messaging/whatsapp/README.md`
- `.claude/rules/meta-devtools-mcp.md`
- `.claude/rules/colas-errores-permanentes.md`
