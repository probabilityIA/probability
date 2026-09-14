# CU-05: Crear plantillas por las tres entradas con las reglas de Meta

Recorre las tres formas de crear una plantilla desde el front y verifica que las
reglas que Meta usa para rechazar se atajan antes de guardar, en el front y en el
backend. Todo contra el mock de WhatsApp: nada llega a Meta.

## Precondiciones

1. Base local, negocio Demo (26), backend en 3050 y front en 3000.
2. Mock en 9103 y `whatsapp_url` de `integration:platform_creds:2` apuntando a
   el, escrito despues de arrancar el backend.
3. El mock rechaza con code 100 una plantilla sin componente BODY, como Meta.

## Paso 1: pestana Plantillas WhatsApp

1. Boton "+ Nueva plantilla" visible junto a "Propias del negocio".
2. Selector de uso: "Campanas y flujos" (`campaign`) o "Reglas programadas"
   (`scheduled`).
3. Encabezado con emoji -> aviso "El encabezado no admite emojis".
4. Boton "+ Nombre del cliente" con el cuerpo terminado en punto -> aviso "no
   puede terminar con una variable"; "Guardar borrador" no guarda nada.
5. Cuerpo valido -> se guarda en borrador con `scope=campaign` y `components`
   completos; "Enviar a revision" -> el mock la aprueba.

## Paso 2: Flujos > + Crear respuesta

1. Pie con variable -> aviso "El pie no admite variables".
2. Pie valido -> se guarda, queda enlazada al boton del flujo y tiene
   `components`.
3. Enviarla a revision -> `approved` con `meta_template_id` del mock.

## Paso 3: Reglas > Programadas por segmento > Nueva plantilla

1. Encabezado con formato (`*Oferta*`) -> aviso "no admite formato".
2. Encabezado valido -> se guarda con `scope=scheduled`.

## Paso 4: el backend aplica las reglas aunque se salte el front

`POST /api/v1/whatsapp-templates?business_id=26` debe responder error para:

- encabezado con emoji
- cuerpo que empieza con variable
- dos variables seguidas
- variable mal escrita (`{{nombre}}`)

Y no debe quedar ninguna fila creada.
