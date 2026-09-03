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
| 2026-08-25 | [Guias duplicadas en EnvioClick: la ventana entre el envio y la respuesta](2026-08-25-guias-duplicadas-idempotencia.md) | Estado `generating` + lock por orden (advisory lock) + `needs_verification` cuando el carrier no confirma. Modal ya no rehabilita el boton a los 45 s. Probado en local con mock. Sin desplegar |
| 2026-08-25 | [Configuracion de envios por negocio y el paquete que se declara al carrier](2026-08-25-configuracion-envios-por-business.md) | Modulo nuevo `shippingconfig` (cajas, transportadoras, bodega predeterminada) + resolvedor compartido usado por checkout Woo/Shopify, cotizacion del panel y generacion de guia. Corrige el 1 kg 10x10x10 de Viga. Sin desplegar |
| 2026-08-27 | [Siigo: documento inactivo, documento sin electronica y errores crudos](2026-08-27-siigo-documento-inactivo-y-errores-crudos.md) | Causa en la cuenta Siigo del cliente (activaron el documento 30606 y ya emite). Errores traducidos a codigo canonico, vista previa de factura para Siigo y sin reintentos en errores de configuracion. Sin desplegar |
| 2026-08-27 | [Siigo: el IVA no se aplica automatico por producto](2026-08-27-siigo-iva-por-producto.md) | Se agrego lectura de `taxes[]` por producto (sin desplegar, sin consultar aun datos reales de Viga). Pendiente decidir si vale la pena IVA por SKU |
| 2026-08-31 | [WooCommerce COD: la comision del carrier se declaraba dos veces](2026-08-31-woocommerce-cod-comision-doble.md) | Corregido y desplegado, verificado por el usuario en Woo de pruebas. Pendiente: reintentar guias 14666/14668 de Viga |
| 2026-08-31 | [Viga: la comision COD del carrier no se factura en Siigo](2026-08-31-siigo-viga-comision-cod-no-facturada.md) | Confirmado con los PDF reales de Siigo: 7 facturas COD quedaron cortas $38.702 en total porque la linea "Envio" usa el costo de la guia y no lo cobrado al cliente. El "doble envio" no se reproduce (era el modal de Editar Orden, corregido en b3611857; quedan 14 ordenes viejas con el dato). Corregido leyendo el meta `cod_carrier_fee` del canal (no se deduce de cod_total: eso rompia las 350 facturas COD de Shopify/Softpymes del negocio 34). Solo afecta Viga y Demo. Faltan las notas credito de las 7 facturas timbradas. Complementa a `2026-08-31-woocommerce-cod-comision-doble.md`: aquel corrige el codValue de la guia, este la factura |
| 2026-09-01 | [Suscripciones: el excedente de uso nunca se cobraba al renovar](2026-09-01-suscripcion-excedente-no-se-cobraba.md) | `PurchaseSubscription` cobraba solo `precio_plan x meses`, ignorando el excedente de envios/facturas/ordenes que si mostraba el "pago pronosticado". Mystic (173/100 envios) pagaba $99.000 en vez de $142.800. Corregido reutilizando el mismo calculo del pronostico; recreado y verificado en local (misma cifra pronosticada y cobrada). Pendiente: no se cobro retroactivamente el excedente historico ya perdonado |
| 2026-09-02 | [WooCommerce: las fechas de orden se guardaban 5 horas corridas por timezone](2026-09-02-woocommerce-desfase-timezone-fechas.md) | `date_created`/`date_paid` de WooCommerce (hora local Bogota, sin offset) se interpretaban como UTC. Orden 14684: creada y pagada con 43s de diferencia real, no 5h -- el gap visible era el mismo timestamp mostrado con dos conversiones distintas. Corregido leyendo `date_created_gmt`/`date_paid_gmt` (UTC real) con fallback a `America/Bogota`. Compilado, sin desplegar, sin corregir datos historicos |
| 2026-09-02 | [Cancelacion revertida por un webhook de tracking atrasado](2026-09-02-cancelacion-revertida-por-webhook.md) | La guia 034058470201 (MYS-0852) si se cancelo en EnvioClick y se marco `cancelled`, pero un `webhook_update` con la foto anterior (evento del carrier a medianoche) la devolvio a `pending` un segundo despues. `cancelled` pasa a ser terminal: ningun tracking ni webhook lo saca de ahi, solo deja un `Warn`. El popup de "Server Action not found" que vio el usuario era la pestana con un bundle viejo, no el backend. Sin desplegar |
| 2026-09-02 | [Una orden con direccion en Zarzal se geocodifico en Guadalajara de Buga](2026-09-02-orden-geocodificada-en-buga-siendo-zarzal.md) | El texto de la direccion siempre decia Zarzal; las coordenadas quedaron pegadas de una sugerencia de Google que resolvio en Buga porque corregir la ciudad a mano no vuelve a geocodificar. Guia devuelta. Corregido: `handleCitySelect` re-geocodifica al corregir la ciudad. Sin desplegar |
| 2026-09-02 | [Viga: recarga fantasma por acceso equivocado, y por que "editar fechas" no desbloqueaba la cuenta](2026-09-02-viga-recarga-fantasma-y-fecha-de-corte.md) | Recarga PENDING de $323.907 nunca llego a Bold (nunca se abrio el checkout), sin riesgo ni movimiento de saldo: eliminada en produccion con autorizacion. Aparte, editar fechas/dias de cortesia nunca reactivaba `subscription_status`, dejando la cuenta bloqueada aunque la fecha ya fuera futura. Corregido con reactivacion automatica. Sin desplegar |
| 2026-09-02 | [Suscripciones: cuentas deshabilitadas con mas acceso que un plan pagado, y planes personalizados que nadie podia pagar](2026-09-02-suscripciones-control-de-acceso-y-planes.md) | Testing E2E completo (CU-01 a CU-15) encontro: `HasModuleAccess` daba mas modulos a una cuenta deshabilitada que a un plan pagado (mismo fallback de "sin plan" mal aplicado, 3 veces); todo plan personalizado nuevo nacia no-comprable (`payable` nunca se seteaba ni se persistia). Ambos corregidos con tests nuevos. Sin desplegar |
