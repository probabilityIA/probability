# Bitacora

Historico de soportes, incidentes y diagnosticos del proyecto. Un archivo por
caso, con el diagnostico completo: que se rompio, como se encontro, que se
descarto en el camino, que se corrigio y que quedo pendiente.

## Para que sirve

Cuando vuelve a aparecer un problema parecido, o cuando alguien (persona o IA)
necesita entender por que el codigo hace algo raro, aca esta el contexto que no
cabe en un commit. Los commits dicen QUE cambio; la bitacora dice POR QUE y con
que evidencia.

## Cuando escribir una entrada

- Se investigo un problema de produccion que costo tiempo o plata.
- Se descubrio un comportamiento no documentado de un proveedor externo.
- Se corrigio data en produccion.
- Una hipotesis razonable resulto falsa (dejarla escrita evita que el siguiente
  la repita).

No hace falta entrada para un fix trivial o un cambio de UI.

## Nomenclatura

`YYYY-MM-DD-tema-corto.md`

## Estructura sugerida

1. Resumen en dos lineas
2. Sintoma (con numeros reales)
3. Diagnostico: la cadena de evidencia, incluidas las hipotesis descartadas
4. Causa raiz
5. Correccion (codigo y/o data)
6. Verificacion
7. Pendientes

## Indice

| Fecha | Tema | Estado |
|-------|------|--------|
| 2026-08-06 | [COD EnvioClick: valor a recaudar mal calibrado](2026-08-06-cod-envioclick-calibracion.md) | Corregido, pendiente deploy |
| 2026-08-11 | [EnvioClick: cancelaciones con falso positivo](2026-08-11-envioclick-cancelacion-falso-positivo.md) | Fix desplegado, pendiente reclamo a EnvioClick |
| 2026-08-12 | [El detector de SKU sugeria cambiar de talla](2026-08-12-detector-sku-sugerencias-falsas.md) | Corregido en local, sin desplegar |
| 2026-08-13 | [Mappings duplicados rompian el apply con ON CONFLICT](2026-08-13-mappings-duplicados-on-conflict.md) | Corregido en local, sin desplegar |
| 2026-08-17 | [Cotizaciones mostraba "Guia generada" para guias que nunca salieron](2026-08-17-cotizacion-guia-generada-fantasma.md) | Data corregida en produccion, codigo sin desplegar, causa raiz abierta |
| 2026-08-17 | [WooCommerce: las ordenes entraban sin direccion de envio](2026-08-17-woocommerce-direccion-en-billing.md) | Fix desplegado, 53 ordenes corregidas en produccion |
| 2026-08-17 | [MercadoLibre: sin JSON crudo, sin direccion y notificaciones de envio rotas](2026-08-17-meli-json-crudo-y-envios.md) | Cerrado: direccion, packs duplicados, estados y etiqueta corregidos y desplegados |
| 2026-08-19 | [Siigo no tiene endpoint de remisiones en su API](2026-08-19-siigo-sin-endpoint-de-remisiones.md) | Cerrado: no es implementable por API, documentado |
| 2026-08-20 | [Tiendanube devuelve 404 cuando el listado esta vacio](2026-08-20-tiendanube-404-last-page-is-0.md) | Cerrado: 404 de paginacion tratado como lista vacia, verificado en la tienda real |
| 2026-08-21 | [Viga: 3 envios COD sin pagar y ordenes duplicadas](2026-08-21-viga-cod-pendiente-y-ordenes-duplicadas.md) | Reclamo valido: nunca entraron a un corte COD. Duplicados por reintento de guia; doble cobro sin reembolso |
| 2026-08-21 | [Las alarmas de CloudWatch no le avisaban a nadie](2026-08-21-monitoreo-sin-destinatario.md) | Cerrado: SNS suscrito y verificado. CloudTrail activado. Faltan metricas de RAM/disco del EC2 |
| 2026-08-22 | [El guard del SSE de cotizacion no filtraba nada](2026-08-22-sse-cotizacion-guard-sin-filtro.md) | Corregido en local, sin desplegar. Pendiente: campo `status` en POST /shipments/quote |
| 2026-08-22 | [Tiendanube: 5 ciclos de OAuth, webhooks huerfanos y reconcile mudo](2026-08-22-tiendanube-oauth-5-ciclos.md) | OAuth verificado 5/5. Modal de productos corregido (sin desplegar). Webhooks huerfanos limpiados a mano, fix de backend pendiente |
| 2026-08-22 | [Tiendanube: las ordenes entran, pero la pagada queda sin estado](2026-08-22-tiendanube-ordenes-sin-estado.md) | Corregido en local (statusmapper, semilla de order_status_mappings 17, firma del webhook, push-back de estado y guia). Direccion de estados configurable por integracion. Sin desplegar |
