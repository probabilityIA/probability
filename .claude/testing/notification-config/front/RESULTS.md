# Resultados front - notification-config

## 2026-09-13 - CU-05 crear plantillas por las tres entradas con reglas de Meta

Base local (Demo, 26), backend y front locales, mock de WhatsApp en 9103.
Playwright contra `http://localhost:3000`.

| Paso | Resultado |
|---|---|
| 1. Pestana Plantillas WhatsApp: boton "+ Nueva plantilla" | OK: visible, abre el modal con selector de uso |
| 1. Encabezado con emoji | OK: aviso "El encabezado no admite emojis: Meta lo rechaza." |
| 1. Variable al final del cuerpo | OK: aviso en vivo y "Guardar borrador" no guarda (0 filas) |
| 1. Cuerpo valido | OK: id 35, `scope=campaign`, `components` con HEADER y BODY; enviada a revision, el mock la aprueba |
| 2. Flujos > + Crear respuesta con variable en el pie | OK: aviso "El pie no admite variables." |
| 2. Respuesta valida | OK: id 36 enlazada al boton "genial", con `components`; enviada a revision -> `approved` |
| 3. Reglas > Programadas > Nueva plantilla con `*Oferta*` | OK: aviso "no admite formato (*, _, ~)" |
| 3. Encabezado valido | OK: id 37, `scope=scheduled`, con `components` |
| 4. API con emoji, variable al inicio, variables seguidas, `{{nombre}}` | OK: los 4 responden `success:false` con el motivo; 0 filas creadas |

**Encontrado de paso:** el boton "+ variable" pegaba `{{n}}` sin espacio al
texto anterior (`contactamos.{{1}}`). Ahora agrega un espacio si hace falta.

**Contexto:** en produccion las plantillas 30-32 de LaPerchaDel10 fallaban con
code 100 porque el repositorio no leia `components` y Meta recibia la plantilla
vacia. Corregido en `e1a0afed`; este caso cubre ademas que el mock rechace una
plantilla sin BODY para que no vuelva a pasar desapercibido.
