# Meta: MCP DevTools vs Graph API

Hay dos caminos para hablar con Meta y hacen cosas distintas. Elegir mal cuesta
tiempo buscando un tool que no existe.

| | MCP `meta_developer_tools` | Graph API (`graph.facebook.com`) |
|---|---|---|
| Sobre que opera | la **app** de Meta (`2812884712240202`) | el **WABA** (plantillas, mensajes, numeros) |
| Autenticacion | usuario humano, OAuth interactivo (`/mcp`) | system user `cam-adm`, token en Redis / `.env` |
| Uso | inspeccion manual en sesion | codigo, scripts, cron |

## Que hace cada uno

**MCP** (reemplaza entrar a developers.facebook.com):

- `devtools_app_list` - primero siempre, da los App IDs accesibles
- `devtools_app` - settings, dominios, seguridad, restricciones
- `devtools_api_usage` - `rate_limits`, `call_volume`, `deprecations`
- `devtools_api_changelog` - que cambia entre versiones de la Graph API
- `devtools_webhook_list` / `webhook_manage` / `webhook_test`
- `devtools_app_review` - estado, permisos, requisitos
- `devtools_discovery` - buscar documentacion oficial
- `devtools_compliance`

**Graph API** (todo lo de WhatsApp propiamente dicho):

- plantillas: listar, crear, editar, consultar `status`
- enviar mensajes (texto libre y template)
- phone numbers, quality rating, messaging limit tier
- `debug_token`

El MCP **no tiene ningun tool de plantillas ni de envio de mensajes**. Si la
pregunta es "cuantas plantillas hay" o "manda un mensaje", es Graph API.

## Reglas

1. **Nunca meter el MCP en codigo, scripts ni cron.** Esta autenticado con el
   usuario, no con el system user; en headless no existe. Automatizacion = Graph
   API con el token del system user.
2. **No confundir los dos rate limits.** El MCP reporta el limite de la **app**
   (Graph API general). El limite de **mensajeria** de WhatsApp (tier
   1K/10K/100K + quality rating) es otro contador y solo se ve por Graph API en
   el phone number. Al diagnosticar "no salen mensajes", revisar el segundo.
3. **Antes de subir de version de la Graph API** (`v22.0` hoy, hardcodeado en
   `whatsapp_url` de las platform creds): correr `devtools_api_usage` accion
   `deprecations` y `devtools_api_changelog`.
4. **Antes de tocar el handler del webhook**: `devtools_webhook_list` para ver
   los campos suscritos reales, y `webhook_test` para disparar un evento de
   prueba sin mandar un WhatsApp real.
5. **La fuente de verdad de las credenciales es Redis**
   (`integration:platform_creds:2`), poblado desde BD. Las variables `META_*` de
   `back/central/.env` son copias de conveniencia para consultar a mano y pueden
   quedar desactualizadas. Nunca hacer que el backend las lea.
6. Toda llamada a Graph API con el system user puede requerir
   `appsecret_proof = HMAC-SHA256(token, app_secret).hex()`. Sin proof varios
   endpoints fallan con `(#100) Missing Permission`.

## Contexto conocido

- La app esta en `dev_mode` (`is_live: false`), sin privacy policy, terms ni data
  deletion URL, y con el email de contacto sin verificar. No estorba mientras
  WhatsApp opere con token de system user, pero bloquea App Review y cualquier
  permiso nuevo.
- `cam-adm` NO tiene WABAs asignados: `GET /me/assigned_whatsapp_business_accounts`
  devuelve vacio. Usar los WABA IDs directos, no intentar descubrirlos.

Detalle, IDs y comandos: `back/central/services/integrations/messaging/whatsapp/README.md`.
