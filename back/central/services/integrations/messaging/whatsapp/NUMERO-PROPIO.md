# WhatsApp: numero propio por negocio

Como un negocio envia y recibe WhatsApp desde SU numero, con Probability como
administrador de su cuenta, sin que el cliente haga tramites con Meta.

Contexto y fases: `.claude/project/whatsapp-numero-por-cliente.md`.

## Como se decide de que numero sale cada mensaje

Desde 2026-09-05 un negocio puede enviar y recibir desde **su propio numero**, con
Probability como administrador de su cuenta. Lo decide su fila de `integrations`
(tipo 2):

| `config` | Efecto |
|---|---|
| sin `phone_number_id`, o `use_platform_token: true` | sale por el numero de Probability (comportamiento historico) |
| `phone_number_id` + `hosted_by_platform: true` | numero del negocio **dentro del WABA de Probability**: modo por defecto |
| `phone_number_id` + `waba_id` del cliente, sin credencial propia | WABA del cliente, con el token de Probability |
| `phone_number_id` + `waba_id` + credencial `access_token` | WABA del cliente, con el token del cliente |

## Los dos modelos de numero propio

Decidido el 2026-09-06: **por defecto el numero del cliente se agrega al WABA de
Probability** (`hosted_by_platform: true`, el `waba_id` se resuelve solo).

|  | Numero alojado en nuestro WABA (defecto) | WABA del cliente |
|---|---|---|
| Quien paga las conversaciones | Probability, y se le cobra al negocio | el cliente, a Meta |
| Plantillas | las 26 ya aprobadas, no hay que replicar nada | hay que crearlas y esperar aprobacion |
| Tramite con Meta | ninguno | compartir el WABA y asignarlo a `cam-adm` |
| Riesgo | una sancion de politica es a nivel de WABA: cae para todos | aislado por cliente |
| Si el cliente se va | migrar el numero a otro WABA | se lleva su cuenta |

El segundo queda para el cliente grande que exija ser dueno de su cuenta o que
ya la tenga con historia.

`credentials_cache.GetWhatsAppConfig` lee esa fila (config + credenciales
desencriptadas, cacheadas por `integrations/core` en `integration:meta:*` y
`integration:creds:*`) y solo cae a la plataforma cuando falta configuracion
propia.

El caso del medio es el del **camino corto**: el cliente comparte su WABA con
Probability desde su Business Manager y nuestro system user `cam-adm` queda de
administrador, asi que enviamos con nuestro token sobre SU numero. El cliente no
tiene que crear una app en Meta ni generar un token permanente. Solo hace falta
un `access_token` propio si el cliente prefiere administrar su cuenta el mismo.

**Lo que NO migra al numero del negocio** (siempre sale del de Probability):

- alertas de monitoreo (`consumeralert`)
- OTP de login (`consumerauthotp`)
- respuestas del agente AI de la plataforma (`consumerai`)
- aviso de saldo bajo de billetera y de ventana de pago de suscripcion
  (`consumerwalletalert`, `consumersubscriptionalert`): son mensajes de
  Probability al negocio, por eso usan `SendPlatformTemplate`.

### Quien puede escribir `phone_number_id` y `waba_id`

Son campos **protegidos**: el `PUT /integrations/:id` generico los ignora y
conserva los guardados (`handlerintegrations/protected_config.go`), igual que
`use_platform_token`. El unico camino para escribirlos es
`PUT /whatsapp/connection`, que antes de guardar:

1. rechaza el `phone_number_id` de la plataforma,
2. rechaza un `phone_number_id` que ya pertenezca a otra integracion,
3. le pregunta a Meta (`GET /{waba_id}/phone_numbers` con el token que se va a
   usar) si ese numero de verdad esta en ese WABA y si tenemos acceso.

Ademas hay un indice unico parcial en base
(`uq_integrations_whatsapp_phone_number_id`) por si alguien escribe por otro
camino. Sin esto, un negocio podia declarar el `phone_number_id` de otro (o el
de Probability) y quedarse con los mensajes entrantes de ese numero: el indice
de ruteo en Redis y la consulta del webhook apuntan al dueno declarado.

### Ruteo del webhook por numero

El webhook es uno solo para todos (la firma HMAC usa el `webhook_secret` de la
app, igual para todos). Cuando llega un mensaje sin conversacion activa ni sesion
humana, se resuelve el dueno con
`phone_number_id -> (integration_id, business_id)`: un SELECT sobre `integrations`
filtrando `config->>'phone_number_id'`, con indice en Redis
(`integration:idx:cfg:2:phone_number_id:<id>`) que se mantiene al cachear la
integracion y se invalida al cambiarle el `config`.

- Si el numero es de un negocio: se abre una conversacion `inbound` para ese
  negocio, se persiste el mensaje, se publica por SSE a su dashboard y se
  activa la `HumanSession` de ese telefono apuntando a esa conversacion. Esa
  activacion no es un detalle: el estado `HANDOFF_TO_HUMAN` es terminal, asi que
  `Save` borra el indice `active` y sin la sesion humana cada mensaje siguiente
  abriria una conversacion nueva. Con ella, los mensajes posteriores entran por
  la rama de sesion humana que ya existia y van al chat del dashboard.
- Si no es de nadie (o es el numero de la plataforma): sigue el camino de
  siempre, el agente AI Sales.

### Alta del numero desde nuestra propia pantalla (modo por defecto)

El negocio digita su numero en Probability y no ve nunca una pantalla de Meta.
Los cuatro pasos son API, con el token del system user (`whatsapp_business_management`
+ `whatsapp_business_messaging`); **no requiere App Review ni Embedded Signup**.

| Paso | Endpoint nuestro | Llamada a Meta | `number_status` |
|---|---|---|---|
| 1. Agregar el numero | `POST /whatsapp/numbers` | `POST /{waba_id}/phone_numbers` (`cc`, `phone_number`, `verified_name`) | `esperando_codigo` |
| 2. Pedir el codigo | `POST /whatsapp/numbers/code` | `POST /{phone_number_id}/request_code` (SMS o VOICE) | `esperando_codigo` |
| 3. Verificar | `POST /whatsapp/numbers/verify` | `POST /{phone_number_id}/verify_code` | `verificado` |
| 4. Activar | `POST /whatsapp/numbers/register` | `POST /{phone_number_id}/register` con PIN | `registrado` |
| Deshacer | `DELETE /whatsapp/numbers` | `POST /{phone_number_id}/deregister` (solo si estaba registrado) + `DELETE /{phone_number_id}` | `sin_numero` |

`GET /whatsapp/numbers` devuelve el estado, y lo completa con
`GET /{phone_number_id}` de Meta (`display_phone_number`, `quality_rating`,
`name_status`, `code_verification_status`).

- El `use_platform_token` recien pasa a `false` en el paso 4: hasta que el numero
  no esta registrado, el negocio sigue enviando por el numero de Probability.
- El **PIN de 2FA lo genera el backend**, se guarda cifrado en las credenciales
  de la integracion (`two_factor_pin`) y se le muestra al usuario UNA vez.
- `number_status`, `verified_name` y `hosted_by_platform` son campos protegidos:
  el `PUT /integrations/:id` generico no los puede tocar. Sin eso, un negocio
  podria marcarse `registrado` sin haber verificado nada.
- No hay paso de plantillas: el numero usa las que ya estan aprobadas.

## Los digitos del numero no se pueden editar

Meta no tiene ninguna operacion para cambiarle `cc` o `phone_number` a un numero
ya creado: solo existe crear y borrar. Lo unico editable de un numero vivo es su
perfil (`POST /{phone_number_id}/whatsapp_business_profile`), el nombre visible
y el PIN de dos pasos.

Por eso un numero mal escrito se corrige con `DELETE /whatsapp/numbers` y volver
a agregarlo. Si el numero ya estaba `registrado`, primero se hace `deregister`
(un fallo ahi no aborta: se registra en el log y se sigue con el borrado), y
`use_platform_token` vuelve a `true`, asi que mientras tanto el negocio envia por
el numero de Probability en vez de quedarse mudo.

Si Meta responde que el objeto ya no existe (`803`, o `100/33`, o 404) se limpia
la configuracion local igual: de otro modo un `phone_number_id` fantasma dejaria
al negocio trabado, sin poder agregar otro numero.

## La validacion antes de llamar a Meta

`validarNumero` rechaza el numero **antes** de crearlo en Meta. Nacio del caso de
`LaPerchaDel10` (2026-09-11): se agrego `+57 30279400815`, once digitos. Meta lo
acepto sin chistar y devolvio un `phone_number_id`, pero el codigo de
verificacion no podia llegar a ningun lado, asi que `request_code` y
`verify_code` fallaban con 400 y el negocio quedaba atascado en
`esperando_codigo` sin forma de deshacerlo.

Reglas: indicativo de maximo 3 digitos, indicativo mas numero de maximo 15
(E.164), y para Colombia (`57`) exactamente 10 digitos empezando por 3. La regla
de Colombia es la unica especifica de un pais, y es la que hubiera atajado este
caso: la de E.164 no, porque 2 + 11 = 13 digitos es un numero internacional
perfectamente valido.

**Como se detecta a ojo un numero que Meta no reconoce:** en
`GET /{waba_id}/phone_numbers`, un numero valido vuelve formateado por pais
(`+57 350 6886697`) y uno invalido vuelve pegado (`+57 30279400815`).

Limites de Meta que se van a ver en la practica:

- **Cupo de numeros por negocio.** Hoy estamos en el tope: la API responde
  `Phone Numbers Count Exceeded Limit Per Business` (subcodigo 2388386). Se pide
  ampliacion en Meta Business Suite; no es App Review.
- El numero no puede estar activo en WhatsApp normal.
- El `verified_name` pasa por revision de Meta. El webhook
  `phone_number_name_update` avisa cuando lo aprueban y **hay que volver a
  llamar `register`**. Ese campo todavia NO esta suscrito en la app: hay que
  agregarlo a la suscripcion del webhook.
- Maximo 10 intentos de `register` por numero cada 72 horas.

**Pendiente para este modelo: medir el consumo por numero y cobrarlo.** Hoy
nadie cuenta cuantos mensajes manda cada negocio, asi que Probability pone la
plata y no la recupera. El punto unico donde registrarlo es `sendTemplate`.

### Onboarding con WABA del cliente (sin Embedded Signup)

1. El cliente entra a su Business Manager -> Configuracion del negocio ->
   Cuentas de WhatsApp -> Socios -> agrega el Business ID de Probabilityapp con
   control total.
2. Se le asigna ese WABA al system user `cam-adm` en Business Manager. Sin eso
   el token no lo ve: `GET /me/assigned_whatsapp_business_accounts` devuelve
   vacio y el paso 4 falla con permisos.
3. Se anota su `waba_id` y su `phone_number_id` y se suscribe la app:
   `POST /{waba_id}/subscribed_apps`.
4. En el front, en la integracion de WhatsApp del negocio, se activa "Usar mi
   propio numero" y se pegan `waba_id` y `phone_number_id`. **El token se deja
   vacio**: se usa el de Probability. Al guardar, el backend le pregunta a Meta
   si ese numero pertenece a ese WABA; si no, no guarda nada.
5. Se pulsa "Crear las plantillas que faltan" y se espera la aprobacion de Meta.

Si el cliente revoca el acceso desde su Business Manager, Meta responde 190 / 401
y el clasificador lo trata como error permanente: el mensaje se descarta con
`Warn` en vez de reencolarse en bucle (ver
`.claude/rules/colas-errores-permanentes.md`).

## Plantillas por WABA

Las 13+ plantillas viven en el WABA de Probability. Un negocio con numero propio
necesita las suyas, aprobadas por Meta en SU cuenta.

- `POST /whatsapp/templates/provision` lee las plantillas del WABA de la
  plataforma (`GET /{waba_id}/message_templates`) y crea en el WABA del cliente
  las que le falten (`POST /{waba_id}/message_templates`). El catalogo origen es
  Meta, no una copia en Go: asi no se desincronizan.
- **No se copian todas.** Las plantillas que son de Probability hacia el negocio
  (`alerta_servidor`, `reporte_saldo_billetera`, `reporte_saldo_billetera_v2`,
  `resumen_pago_suscripcion`, `recuperacion_codigo`) y `hello_world` quedan
  fuera; salen en `skipped`. La lista vive en `usecasetemplates/catalog.go`.
- Las plantillas de categoria `AUTHENTICATION` **no se pueden copiar tal cual**:
  Meta devuelve en el GET el texto ya renderizado, pero en el POST solo acepta
  `BODY {add_security_recommendation}`, `FOOTER {code_expiration_minutes}` y
  botones `OTP`. `componentsForCreate` las reconstruye con esa forma.
- `GET /whatsapp/templates/status` devuelve el estado por plantilla. Se cachea en
  Redis (`whatsapp:templates:{integration_id}`, TTL 6 h) y `?refresh=true` fuerza
  la consulta a Meta.
- El webhook `message_template_status_update` actualiza ese cache usando el WABA
  del `entry.id` (`whatsapp:waba:{waba_id}` -> integration_id).

Mientras una plantilla no este `APPROVED`, los mensajes que la usan no se pueden
enviar: es el estado "cliente conectado pero sin plantillas aprobadas" que la UI
muestra explicitamente.

## Registro insertado (Embedded Signup) - detras de flag

El cliente pulsa "Conectar cuenta de WhatsApp", inicia sesion con su Facebook en
la ventana de Meta, elige o crea su WABA y su numero, y vuelve conectado. Es el
mismo destino que el asistente manual, pero sin que nadie copie IDs.

**Esta apagado por defecto.** Se prende desde las platform credentials del
integration_type 2 (nunca desde `.env`, ver la regla de credenciales OAuth):

| Campo | Para que |
|---|---|
| `embedded_signup_enabled` | el flag; sin el, `GET /whatsapp/embedded-signup/config` responde `enabled: false` y el boton no se pinta |
| `app_id` | el de la app de Meta (`2812884712240202`) |
| `app_secret` | necesario para canjear el codigo; **no esta en el repo**, hay que cargarlo |
| `embedded_signup_config_id` | el id de la configuracion de registro insertado que se crea en el panel de la app |

Flujo:

1. El front carga el SDK de Meta y llama `FB.login` con `config_id`,
   `response_type: code` y `override_default_response_type`.
2. Del evento `WA_EMBEDDED_SIGNUP` (postMessage) salen `waba_id` y
   `phone_number_id`; del callback sale el `code`.
3. `POST /whatsapp/embedded-signup` con los tres. El backend:
   - canjea el `code` por el token del negocio (`GET /oauth/access_token`),
   - comprueba con ese token que el numero esta en ese WABA,
   - rechaza el numero si ya pertenece a otra integracion,
   - suscribe el webhook (`POST /{waba_id}/subscribed_apps`),
   - registra el numero (`POST /{phone_number_id}/register`) con un PIN generado,
   - guarda el token del negocio y el PIN cifrados, y deja
     `use_platform_token: false` solo si el registro salio bien.
4. Si el registro falla, la conexion queda guardada con
   `number_status: verificado` y el negocio sigue enviando por el numero de
   Probability hasta que se resuelva. La respuesta trae `warning`.

**Limite que no desaparece:** con la app en `dev_mode` el flujo solo funciona
para cuentas con rol en la app. Para clientes reales hace falta Advanced Access
en `whatsapp_business_management` y `whatsapp_business_messaging`, que hoy estan
en `REJECTED`. El codigo no depende de eso; el permiso si.
