# CU-04: Adjuntos en el chat y retencion de 1 ano

Recorre el envio de archivos desde la bandeja, la recepcion de archivos del
cliente y el job que borra el historial de mas de 365 dias. Todo contra el mock
de WhatsApp: nada llega a Meta.

Los archivos SI se suben al bucket real `probability-chat-attachments` (en local
no hay MinIO). Quedan con la regla de 365 dias del bucket; al terminar se borran
los objetos de prueba.

## Precondiciones

1. Base local, negocio Demo (26), migracion `migrateWhatsappMessageMedia`
   aplicada (columnas `media_*` en `whatsapp_message_logs`).
2. Backend en 3050 con credenciales S3 en `back/central/.env`.
3. Mock en 9103 con las rutas de archivos (`GET /_mock/media` responde) y
   respuestas automaticas apagadas.
4. `whatsapp_url` de `integration:platform_creds:2` apuntando al mock, escrito
   despues de arrancar el backend.
5. Una conversacion con envio previo en el mock (por ejemplo la de CU-03, A =
   573102020202), para que el mock sepa a que `phone_number_id` responder.

Endpoint de envio: `POST /api/v1/integrations/whatsapp/conversations/{id}/reply-media`
(multipart: `file`, `caption`, `phone_number`, `business_id`).

## Paso 1: enviar una imagen con texto

PNG pequeno con `caption`.

**Esperado:**
- `200 {"status":"sent"}`.
- El mock registra `POST /{phone}/media` y un envio `type=image` con `media_id`
  y el caption (`GET /_mock/sent`).
- En el detalle de la conversacion aparece el mensaje con `media.type=image`,
  `media.available=true` y una `url` firmada que responde 200 con `image/png`.
- El objeto existe en `s3://probability-chat-attachments/whatsapp/26/<yyyy>/<mm>/`.

## Paso 2: enviar un PDF sin texto

**Esperado:** envio `type=document` con `filename`; en el detalle
`media.type=document`, `filename` correcto y `content` vacio.

## Paso 3: rechazos

- Archivo `.exe` -> 400 `invalid_media`, no llega nada al mock ni a S3.
- Imagen de mas de 5 MB -> 400.
- Archivo de mas de 10 MB -> 413.
- Negocio normal enviando `business_id` de otro negocio en el form: se ignora y
  se usa el del token.

## Paso 4: el cliente envia una imagen

`POST /_mock/inbound-media {"phone": A, "type": "image", "caption": "Asi llego"}`.

**Esperado:** mensaje entrante con `media.type=image`, `media.available=true`,
URL firmada 200, contenido "Asi llego", la conversacion queda sin leer y la
vista previa del listado dice "Asi llego".

## Paso 5: el cliente envia un documento sin texto

`POST /_mock/inbound-media {"phone": A, "type": "document"}`.

**Esperado:** `media.type=document`, `filename=comprobante.pdf`, y la vista
previa del listado dice "[Documento]".

## Paso 6: retencion

El job corre 2 minutos despues de arrancar y luego cada 24 h, con corte en
365 dias. Las pruebas no escriben la base, asi que se valida sin borrar:

- `EXPLAIN` de los tres `DELETE` (mensajes por lotes, conversaciones vacias,
  marcas de leido huerfanas) contra la base local: planifican sin error.
- Conteo de solo lectura con el mismo filtro: cuantas filas borraria hoy.
- Test unitario de `chatretention`: el corte es exactamente `now - 365 dias`.
- Log del backend: una linea "Retencion de chats" por corrida.

## Limpieza

Borrar los objetos de prueba de `whatsapp/26/` en el bucket y dejar el
`whatsapp_url` como quedo (vuelve solo al reiniciar el backend).
