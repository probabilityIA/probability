# Plantillas: Meta las rechazaba con code 100 y las aprobadas se quedaban en pending

**Ticket:** sin ticket (trabajo hecho en sesion, falta registrarlo)
**Negocio:** LaPerchaDel10 (37), integracion 70, WABA `1302830408357767`
**Canal:** WhatsApp Cloud API

## Resumen

Dos fallas en el camino de las plantillas propias. (1) Toda plantilla enviada a
revision llegaba a Meta sin componentes: Meta respondia code 100 y el error
mostrado hablaba de un "Phone Number ID '0'" que no tenia nada que ver. (2) Si se
perdia el webhook de aprobacion, la plantilla quedaba `pending` para siempre y
bloqueaba el lanzamiento de campanas con flujo.

## Sintoma

- Plantillas 30, 31 y 32 (`37_percha_prueba_*`) en `failed` con
  `Solicitud invalida (code 100). Verifica que el Phone Number ID '0' sea correcto`.
- En produccion **todas** las plantillas de negocio (27 a 32) tenian
  `components = 'null'` (el texto JSON, no SQL NULL).
- Tras corregir (1), Meta aprobo las tres, pero la 32 siguio `pending` en la
  base. La campana 1 respondio `400` al lanzar: "1 respuesta sin aprobar".

## Diagnostico

1. Hipotesis descartada: credenciales viejas por integracion. LaPerchaDel10 usa
   las mismas credenciales de Probability; se probo que tanto el token de
   `.env.ai` como el de Redis de produccion crean la misma plantilla en Meta
   (plantillas de sonda creadas y borradas).
2. Hipotesis descartada: solo las respuestas creadas desde el flujo venian sin
   componentes. La 30 (creada desde Reglas) tambien tenia `null`.
3. `templateToDomain` (repositorio de `notification_config`) nunca deserializaba
   `Components`. Cualquier plantilla leida de la base quedaba sin componentes,
   el mensaje de envio a Meta viajaba vacio y el `UpdateTemplate` posterior
   escribia `null` en la columna.
4. El mensaje enganoso venia de `MetaGraphError.FriendlyMessage()`: el caso
   code 100 generico asume un envio de mensaje e imprime `PhoneNumberID`, que en
   la creacion de plantillas es 0.
5. El mock local no validaba componentes, por eso nunca aparecio en pruebas.
6. La 32: el webhook de aprobacion llego (o no) durante un deploy blue-green y
   no hay nada que vuelva a consultar el estado. Relacionado con
   [2026-09-12](2026-09-12-aprobacion-de-plantilla-se-pierde-con-cache-fria.md),
   que corrigio otra forma de perder el mismo aviso.

## Causa raiz

- Deserializacion faltante de `components` en el repositorio.
- Ausencia total de reconciliacion de estado con Meta: el estado solo cambiaba
  por webhook.

## Correccion

| Commit | Que |
|---|---|
| `e1a0afed` | Leer `components`; reconstruirlos con `BuildMetaComponents` al enviar; rechazar envio sin componentes; error de Meta con mensaje y subcode reales; el mock exige BODY |
| `85054162` | Boton crear en la pestana; reglas de Meta validadas en front y backend (variable al inicio/final, variables seguidas, `{{nombre}}`, emojis/formato en encabezado, variables en pie/botones) |
| `ba3f4344` | Boton "Consultar estado en Meta" (`POST /whatsapp-templates/sync-status`) y job `template_status_sync` cada 10 min que solo actua si hay plantillas pendientes enviadas hace mas de 5 min. El modulo WhatsApp lista el WABA y aplica el estado por el mismo camino del webhook. Mock con `silent` para simular aviso perdido |

No se escribio data a mano en produccion: la 32 quedo `approved` con el boton.

## Verificacion

- Local contra mock: plantilla creada, enviada, aprobada; aviso perdido
  simulado (webhook a puerto muerto + `silent`), el boton y el job la dejaron en
  `approved`. CU-05 en `.claude/testing/notification-config/front/`.
- Produccion (Playwright, LaPerchaDel10): 30/31/32 aceptadas por Meta; boton
  paso la 32 a `approved`; campana 1 lanzada solo al cliente 581219
  (3192611891), inicio enviado, "Opcion A" -> camino_a y "Opcion B" -> camino_b
  en menos de 1 s, campana `completed` 1/1.

## Pendientes

- Registrar el ticket y enlazar esta entrada.
- Advertencia sin impacto: `key not found: whatsapp:templates:70` en la cache al
  aplicar el estado (se persiste igual).
- Datos de prueba en LaPerchaDel10: cliente 581219, campana 1, flujo 1 y
  plantillas 30-32.
- Plantillas 27-29 (negocio 66) quedaron con `components = 'null'`; se
  reconstruyen solas al reenviar o editar.
