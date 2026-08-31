# Reglas de espanol para los textos del front

Referencia del skill `ortografia-front`. Lo que `revisar.py` no puede decidir
solo se decide aca.

## 1. Tilde diacritica: la decide la frase, no el diccionario

Estas palabras existen con y sin tilde y significan cosas distintas. El script
las marca con `--ambiguos` pero **nunca las corrige**.

| Sin tilde | Con tilde | Como distinguirlas |
|---|---|---|
| esta, este (este pedido) | está, esté (verbo estar) | si se puede reemplazar por "se encuentra", lleva tilde |
| el (el pedido) | él (persona) | si se refiere a alguien, lleva tilde |
| tu (tu cuenta) | tú (tú decides) | posesivo sin tilde, pronombre con tilde |
| mi (mi cuenta) | mí (para mí) | despues de preposicion, con tilde |
| si (si aplica) | sí (afirmacion, sí mismo) | condicional sin tilde |
| se (se genero) | sé (yo sé) | verbo saber, con tilde |
| de (de la orden) | dé (que él dé) | verbo dar, con tilde |
| mas (arcaico: pero) | más (cantidad) | en UI **siempre** `más` |
| aun (incluso) | aún (todavía) | si se puede cambiar por "todavía", con tilde |
| solo (unicamente / sin compania) | — | la RAE **ya no** pone tilde. `solo` siempre |

Interrogativos y exclamativos llevan tilde **siempre**, tanto en pregunta
directa como indirecta:

```
¿Qué pasó?          No sé qué pasó.
¿Cómo se genera?    Revisa cómo se genera.
¿Cuándo llega?      Depende de cuándo llegue.
¿Dónde? ¿Cuál? ¿Cuánto? ¿Quién?
```

Sin tilde cuando son relativos: `la orden que creaste`, `el modo como se hace`.

## 2. Pasado vs presente

Terminaciones `-o` de primera persona (sin tilde) contra `-ó` de tercera
persona del pasado (con tilde). En mensajes de UI casi siempre es la tercera:

| Presente (yo) | Pasado (el/ella/ello) |
|---|---|
| genero | La guía se generó |
| cambio | El estado cambió |
| proceso | Se procesó el archivo |
| sincronizo | La integración sincronizó |
| elimino | Se eliminó el registro |
| creo | Se creó la orden |
| actualizo | Se actualizó el inventario |
| envio | Se envió la notificación |

## 3. Signos de apertura: van en pareja

Toda pregunta abre con `¿` y toda exclamacion con `¡`. No es opcional en
espanol, aunque el ingles no lo use.

```
mal:  confirm('Eliminar este adjunto?')
bien: confirm('¿Eliminar este adjunto?')

mal:  'Guardado con exito!'
bien: '¡Guardado con éxito!'
```

Si la pregunta arranca a mitad de la frase, el `¿` va donde empieza la
pregunta, no al principio de la oracion:

```
Ya revisaste el pedido, ¿lo confirmamos?
```

## 4. Plurales: el singular lleva tilde, el plural no (y al reves)

Trampa constante. La tilde depende de donde cae el acento, y el plural mueve la
silaba.

| Singular | Plural | Nota |
|---|---|---|
| integración | integraciones | el plural **pierde** la tilde |
| conexión | conexiones | igual |
| comisión | comisiones | igual |
| versión | versiones | igual |
| camión | camiones | igual |
| razón | razones | igual |
| imagen | imágenes | el plural **gana** tilde |
| orden | órdenes | igual |
| margen | márgenes | igual |
| volumen | volúmenes | igual |
| resumen | resúmenes | igual |
| origen | orígenes | igual |
| código | códigos | ambos con tilde |
| página | páginas | ambos con tilde |

Regla corta: `-ción/-sión` pierden la tilde al pluralizar; `imagen`, `orden`,
`margen`, `volumen`, `resumen`, `origen` la ganan.

## 5. Mayusculas

- **Las mayusculas llevan tilde.** `Órdenes`, `Última Milla`, `Envíos`,
  `Añadir`. No existe la excepcion que muchos creen.
- En titulos y botones se usa mayuscula solo en la primera palabra
  (capitalizacion de oracion), no en todas como en ingles:
  `Generar guía`, no `Generar Guía`.
- Nombres propios de producto e integraciones se respetan como los escribe el
  proveedor: `WooCommerce`, `MercadoLibre`, `Shopify`, `VTEX`, `EnvioClick`.
- Los dias de la semana y los meses van en minuscula: `lunes`, `enero`.

## 6. Terminologia del producto

Un mismo concepto se escribe siempre igual en toda la UI:

| Usar | No usar |
|---|---|
| orden | pedido, order |
| guía | rótulo, label, etiqueta de envio |
| envío | shipment, embarque |
| transportadora | carrier, courier |
| bodega | almacén, warehouse |
| integración | conector, canal (salvo "canal de venta") |
| contra entrega | COD, cash on delivery (en texto visible) |
| cotización | cotizador (salvo el nombre de la pantalla) |
| negocio | business, tenant |
| sincronización | sync |
| inicio de sesión | login |
| cerrar sesión | logout |
| contraseña | password, clave |
| correo electrónico | email, e-mail |
| ajustes / configuración | settings |

Siglas que se dejan como estan: `SKU`, `API`, `URL`, `PDF`, `CSV`, `NIT`,
`IVA`, `ReteFuente`, `DIAN`, `WhatsApp`.

## 7. Tono y redaccion

- **Tuteo, no usted.** `Revisa tu conexión`, no `Revise su conexión`.
- Mensajes de error: que digan que paso y que hacer. `No se pudo generar la
  guía: la bodega no tiene teléfono de contacto.` En vez de `Error 422`.
- Sin ingles suelto en texto visible: nada de "Loading...", "No data",
  "Are you sure?".
- Sin punto final en labels, botones ni encabezados. Con punto en frases de
  ayuda y mensajes de error.
- Numeros y moneda en formato colombiano: `$ 12.500`, separador de miles con
  punto.

## 8. Errores tipograficos que ya aparecieron en el repo

- `A{'ñ'}adir` - mezcla de literal y escape, ilegible. Va `{'Añadir'}` o
  `Añadir` segun el largo del archivo.
- `Anadir`, `ano`, `anos`, `contrasena`, `disenar` - falta la enie. `ano` no
  significa lo mismo que `año`: es el error mas caro de todos.
- `Interrapidisimo` - la transportadora se llama `Interrapidísimo`.
- Doble espacio entre palabras.
- Espacio antes de `?`, `!`, `,` o `:` (eso es frances, no espanol).
- `...` escrito con tres puntos sueltos donde el diseno espera `…`: dejar los
  tres puntos, es ASCII y se ve igual.
