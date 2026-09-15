# MYS-0853: la guia salio sin conjunto ni casa y la devolvieron

**Ticket:** sin ticket propio. La causa ya estaba corregida por **TKT-000086**
(`37ff2a03`, `3adb2b8c`, 2026-09-10/11); este caso es la verificacion de ese fix.
**Negocio:** Mystic Rose - Official (business 36). **Canal:** orden manual.
**Transportadora:** Interrapidisimo via EnvioClick.

## Resumen

La guia 240060420582 (2026-09-02) imprimio solo "Entrada 16, El Tigre Cerritos":
el conjunto y la casa viajaron en el campo `suburb`, que Interrapidisimo no
imprime. Paso 4 dias en "Confirmacion Telefonica" y se devolvio. Es el unico
caso de 1.321 guias de Mystic donde este error termino en devolucion. Con el
formulario nuevo de direccion la guia sale completa (verificado con guia real).

## Sintoma

| Dato | Valor |
|---|---|
| Orden | MYS-0853, `48f66620-0b27-4916-98de-4cc7a25cfc03` |
| Shipment | 48331, guia `240060420582`, $17.165, no COD |
| Destino | Pereira (Risaralda) |
| Estado | `returned`, "Devuelto Al Remitente" el 2026-09-08 |
| Tracking | 4-7 sep "En Confirmacion Telefonica", 7 sep "Para Nuevo Intento Entrega", 8 sep devuelto |

`shipping_street` guardado:

```
entrada 16, el tigre Cerritos | entrada 16, el tigre Cerritos conjunto | bosques de Yarima casa 41
```

Payload enviado a EnvioClick (`shipment_sync_logs` 27251, `generate`):

```
address:     entrada 16, el tigre Cerritos
crossStreet: entrada 16, el tigre Cerritos
reference:   conjunto residencial
suburb:      bosques de Yarima casa 41
```

PDF impreso:

```
PARA: ENTRADA 16, EL TIGRE CERRITOS
      PEREIRA\RISA\COL
OBS:  entrada 16 el tigre Cerritos, conjunto residencial, ...
```

Faltaban "Bosques de Yarima" y "Casa 41". El mensajero llego al sector sin
saber a que conjunto ni a que casa ir.

## Diagnostico

1. **Interrapidisimo no imprime `suburb`.** En PARA va `address`; en OBS van
   `crossStreet`, `reference` y lo que quepa, cortado con "...".
2. **La direccion venia mal repartida.** La segunda parte repetia la calle y el
   dato clave estaba en la tercera parte, la que el codigo de entonces mandaba
   como barrio.
3. **El fix del 2026-09-10 no cubria este formato de texto libre.** Pasando el
   mismo `shipping_street` por el `buildGuideDestination` actual:
   `crossStreet = "entrada 16, el tigre Cerritos conjunto bosques de"` (cortado
   en 50), `reference = "entrada 16, el tigre"`, `suburb = "bosques de Yarima casa 41"`.
   La casa se seguia perdiendo. Lo que si lo resuelve es el formulario con
   campos separados (`3adb2b8c`).

### Barrido de todas las guias de Mystic

Se descargaron los 1.321 PDFs de `shipments.guide_url` del business 36 y se
compararon con la direccion de su orden (palabras del complemento ausentes en
el PDF).

| Estado | Guias |
|---|---|
| Entregadas | 1.190 |
| Devueltas | 3 (todas Interrapidisimo) |

Las 3 devoluciones:

| Orden | Direccion | PDF | Veredicto |
|---|---|---|---|
| **MYS-0853** | ... \| bosques de Yarima casa 41 | sin conjunto ni casa | **Error nuestro, evitable** |
| MYS-0655 | El jardin \| etapa 1 manzana 24 casa 2 segundo piso \| El jardin | "El jardin, etapa 1 manzana 24 casa 2" | No: salio casi completa, la direccion no tenia calle |
| MYS-0021 | Cimitarra departamento de Santander | igual | No: la orden no tenia direccion |

Caso casi igual que se entrego: **MYS-0703** (2026-08-18). "Villas de la
Candelaria casa 41, km 14" quedo en `suburb` y la OBS cortada; tuvo "Para Nuevo
Intento Entrega" pero llego.

En las guias de Interrapidisimo desde 2026-08-17 (desde cuando hay payload en
`shipment_sync_logs`): 19 mandaron casa/apto/torre/conjunto en `suburb`, y solo
MYS-0853 se devolvio. El error existia desde antes del fix; casi siempre el
mensajero lo resolvia llamando al cliente.

### Hipotesis descartadas

- **"Lo rompio el fix de TKT-000086"**: falso. La guia es del 2026-09-02 y el fix
  del 2026-09-10.
- **Envia con el mismo problema**: MYS-0156, MYS-0500 y MYS-0796 tuvieron novedad
  "Direccion incorrecta/insuficiente", pero sus PDFs traen la direccion completa.
  Envia imprime todo en una linea; el error es propio de Interrapidisimo.

## Causa raiz

Antes del 2026-09-10 la direccion era una sola cadena con `|` y el reparto a los
campos de EnvioClick dependia de en que parte escribiera el operador el
complemento. Cuando quedaba en la tercera parte iba a `suburb`, que
Interrapidisimo no imprime.

## Correccion

Ya estaba en produccion (TKT-000086):

- `37ff2a03` - `buildGuideDestination` (`front/central/src/shared/utils/guide-destination.ts`)
  como unica construccion del destino.
- `3adb2b8c` - campos separados en el formulario de orden: Tipo + Numero ->
  `reference`, Torre y Edificio/Conjunto -> `crossStreet`, Barrio -> `suburb`.

No hubo cambio de codigo en este caso.

## Verificacion (2026-09-14)

Contra la base local (business Demo) y EnvioClick real, por la interfaz con
Playwright:

1. Orden **DEM-0048** con los datos de la clienta y telefono de prueba
   (3023406789): Direccion `Entrada 16, El Tigre Cerritos`, Tipo `Casa`, Numero
   `41`, Conjunto `Conjunto Bosques de Yarima`, Barrio `Cerritos`, Pereira.
2. Guia Interrapidisimo **240061250779** ($18.275, monedero). Payload:
   `address: Entrada 16, El Tigre Cerritos`, `crossStreet: Conjunto Bosques de Yarima`,
   `reference: Casa 41`, `suburb: Cerritos`.
3. PDF:
   ```
   PARA: ENTRADA 16, EL TIGRE CERRITOS
         PEREIRA\RISA\COL
   OBS:  Conjunto Bosques de Yarima, Casa 41, Cerritos
   ```
4. Guia cancelada desde Envios: EnvioClick `to_refund_orders: [4698706]`, estado
   final "Cancelado". Shipment local 44535 en `cancelled`.

## Pendientes

- **Uso del formulario en Mystic.** De 6 ordenes manuales creadas desde el
  2026-09-11, ninguna usa Tipo/Numero y dos tienen ciudad y departamento en
  Torre/Edificio (MYS-0930, MYS-0931). Parece que algo reparte la cadena con
  `|` en esos campos; no se investigo.
- **MYS-0929**: el payload enviado no coincide con lo que produce
  `buildGuideDestination` con sus campos (editado en el modal o via masiva); no
  se investigo.
- Detalles de UI vistos en la prueba, sin corregir:
  - `OrderForm.tsx` `handleCityBlur`: muestra "Selecciona una opcion del
    listado" despues de elegir la ciudad con clic (closure viejo); la orden se
    crea bien.
  - Paso 3 de la guia parte el nombre como "Angela" / "Maria Arango Marin" y el
    apellido supera el limite de 14 caracteres.
  - El detalle del envio muestra la bodega de la orden, no la usada en la guia.
  - En el modal de cancelar envio el mapa se superpone a los detalles.
