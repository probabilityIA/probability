# 2026-08-22 - Tiendanube: 5 ciclos de borrar y reconectar por OAuth en produccion

## Contexto

Se pidio validar que el OAuth de Tiendanube funciona de forma repetible: eliminar
la integracion del business Demo (26) en produccion y volver a conectarla cinco
veces seguidas, con Probability y Tiendanube abiertos en el mismo navegador.

Tienda usada: `sebas y corotos` (store_id `8126740`, `sebasycorotos2.mitiendanube.com`),
app de Tiendanube `39928`, cuenta partner `probabilitysas@gmail.com` (credenciales
en `.env.ai` como `TIENDANUBE_PARTNER_*`).

## Resultado: OAuth OK, 5 de 5

| Ciclo | Integracion eliminada | Integracion creada | Hora |
|---|---|---|---|
| 1 | 259 | 260 | 20:34 -> 20:36 |
| 2 | 260 | 261 | 20:38 -> 20:39 |
| 3 | 261 | 262 | 20:40 -> 20:40 |
| 4 | 262 | 263 | 20:41 -> 20:41 |
| 5 | 263 | 264 | 20:41 -> 20:42 |

En cada ciclo el `access_token` cifrado quedo distinto y los 7 webhooks se
crearon solos contra el nuevo `integration_id`, lo que prueba que el token nuevo
sirve (crear webhooks es una llamada autenticada a la API de Tiendanube).

## Como se autentica el navegador contra la tienda

El boton "Conectar con Tiendanube" manda a `www.tiendanube.com/apps/<client_id>/authorize`,
que exige sesion de **admin de la tienda**, no de partner. La sesion de partner no
sirve por si sola. El camino es:

`partners.tiendanube.com/stores/details/8126740` -> enlace "Administrar tienda"
(`tiendanube.com/auth/partners/login?token=...`) -> abre el admin de la tienda.

Con esa sesion abierta, el authorize ya no pide login y, como la app sigue
instalada, ni siquiera muestra la pantalla de consentimiento: redirige directo al
callback. Por eso los 5 ciclos corrieron sin intervencion manual.

## Hallazgo 1: eliminar la integracion NO borra los webhooks en Tiendanube

Tras los 5 ciclos habia **35 webhooks** en la tienda (5 juegos de 7), 28 de ellos
apuntando a integraciones ya eliminadas:

```
{'260': 7, '261': 7, '262': 7, '263': 7, '264': 7}
```

Cada webhook huerfano sigue disparando `POST /api/v1/tiendanube/webhook?integration_id=<id borrado>`
en cada evento de la tienda. Se limpiaron a mano con
`DELETE /api/v1/integrations/264/webhooks/<webhook_id>`; quedaron solo los 7 de la 264.

Nota: la integracion 259 no tenia webhooks propios en la tienda, por eso el conteo
despues del ciclo 1 fue 7 y no 14.

Pendiente de fondo en `.claude/alerts/tiendanube-webhooks-huerfanos.md`.

## Hallazgo 2: el modal "Sincronizacion de Productos" decia "Todo sincronizado" con 0

`POST /tiendanube/products/reconcile` es **asincrono**: responde 202 con
`{success, correlation_id, message}` y publica el resultado por SSE
(`tiendanube.product.reconcile.completed`) y en `integration_sync_runs`.

`TiendanubeProductSyncModal.analyze()` leia `res.matched`, `res.only_in_probability`
y `res.only_in_tiendanube` de esa respuesta 202: todos `undefined` -> 0 -> pintaba
"Todo sincronizado". Ademas el modal solo escuchaba `tiendanube.product.sync.*`,
nunca los eventos `reconcile.*`, asi que el resultado real jamas llegaba a la UI.

El backend si calculaba bien. La corrida disparada a mano contra la 264 quedo en
`integration_sync_runs` id 31: `only_in_probability = 57`, `matched = 0`.

Arreglado: el modal ahora guarda el `correlation_id`, escucha
`tiendanube.product.reconcile.completed` y rearma las listas leyendo el detalle
persistido (`/internal/sync-run-items`, grupos `not_associated`, `only_probability`,
`only_channel`). Se agregaron `probability_no_sku` y `tiendanube_no_sku` al evento.

**El mismo defecto vive en Shopify y Jumpseller** (sus handlers de reconcile
tambien son asincronos y sus modales leen la respuesta directa). WooCommerce y
VTEX no lo tienen porque su reconcile sigue siendo sincrono.

## Hipotesis descartadas

- "El listado de integraciones del usuario demo esta roto": el modulo mostraba
  `1 activas - 1 en total` en la pestana "Todos" y la integracion de Tiendanube no
  aparecia. El API devolvia las 19. Era la pestana: la de Tiendanube esta bajo
  **E-commerce**, no bajo "Todos".
- "Probar conexion fallo por el OAuth": el boton "Probar" devuelve
  `permisos insuficientes` porque `TestIntegrationHandler` es **solo super admin**
  (`middleware.IsSuperAdmin`), no porque el token estuviera mal.
- "No hay productos porque el catalogo esta vacio": Demo tiene 57 productos; lo
  que no existe es ninguna fila en `product_business_integrations` para las
  integraciones de Tiendanube, por eso `matched = 0`.
