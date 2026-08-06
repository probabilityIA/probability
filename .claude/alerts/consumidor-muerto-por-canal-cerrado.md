# Consumidores que mueren sin reconectar (canal cerrado con conexion viva)

Fecha: 2026-08-06
Incidente: el consumidor de `invoicing.softpymes.requests` se cayo el 2026-08-05
a las 17:47 UTC y no volvio hasta que reiniciaron el contenedor a las 03:01 UTC
del 2026-08-06. Nueve horas sin procesar facturas.

## Causa raiz

`shared/rabbitmq` solo vigilaba la CONEXION (`r.conn.NotifyClose`). Cada consumidor
corre sobre su propio canal (`r.conn.Channel()`) y ese canal no se vigilaba. Cuando
el broker cierra el canal pero deja viva la conexion, la goroutine del worker
terminaba con un log que decia "will be restored on reconnection", pero como la
conexion seguia arriba `NotifyClose` nunca disparaba y `reregisterConsumers` nunca
corria. El consumidor quedaba muerto hasta el siguiente reinicio.

Evidencia del broker (`docker logs rabbitmq_prod`):

```
2026-08-05 17:47:10 [warning] Consumer 'ctag-...-26' on channel 65 and queue
  'invoicing.softpymes.requests' has timed out waiting for a consumer
  acknowledgement of a delivery with delivery tag = 39.
  Timeout used: 1800000 ms
2026-08-05 17:47:10 [error] channel exception precondition_failed:
  delivery acknowledgement on channel 65 timed out
```

La conexion `<0.668963.0>` siguio viva hasta las 20:48, tres horas despues de que
murio el canal 65.

## Por que le toco a softpymes

`consumerPrefetchCount` era 50 fijo para todos los consumidores. Softpymes quedo
en `workers=1` tras el incidente de julio (ver `softpymes-timeouts-masivo.md`) y
cada factura son 3-8 llamadas HTTP con timeout de 90s.

Con prefetch 50 y 1 worker el broker entrega 50 mensajes de una sola vez y 49
quedan sin ACK mientras el worker procesa el primero. El `consumer_timeout` del
broker (1800000 ms = 30 min) se mide desde la ENTREGA hasta el ACK, no desde que
el worker empieza a trabajar el mensaje. El `delivery tag = 39` del log confirma
que el mensaje que reviento era el 39 de esa tanda.

## Arreglado 2026-08-06 (sin desplegar aun)

En `shared/rabbitmq/rabbitmq.go`:

1. `watchConsumerChannel`: cada canal de consumidor registra `NotifyClose` y, si
   se cierra con la conexion viva, reinicia ese consumidor con backoff.
2. Contador de epoca de conexion (`connEpoch`): evita duplicar consumidores
   cuando la conexion tambien murio y `reregisterConsumers` ya los recreo.
3. `consumerPrefetch(workers)` reemplaza la constante 50: ahora es un mensaje sin
   ACK por trabajador, asi la espera de un mensaje no depende de cuantos tenga
   delante.
4. El log enganoso ("will be restored on reconnection") dice ahora lo que pasa
   de verdad.

`decidirRecuperacion` quedo como funcion pura para poder testear las cinco ramas
sin broker (`shared/rabbitmq/rabbitmq_test.go`).

## Items

### Urgente

- [ ] Desplegar. Mientras no se despliegue, cualquier consumidor puede volver a
      morir igual; softpymes es solo el mas expuesto.
- [ ] Revisar si quedaron facturas sin procesar de la ventana 17:47-03:01 del
      2026-08-05/06 (la cola quedo en 0 mensajes, asi que probablemente se
      consumieron al reiniciar, pero hay que confirmar contra `invoice_sync_logs`).

### Importante

- [ ] Bajar el prefetch reduce el throughput de los consumidores rapidos. Medir
      despues del deploy; si algun consumidor de alto volumen se queda corto, la
      solucion NO es volver a 50 sino agregarle workers.
- [ ] Considerar subir `consumer_timeout` del broker como defensa en profundidad.
      Con 8 llamadas de 90s una sola factura puede tardar ~12 min; sigue por
      debajo de los 30 min pero el margen no es grande.
- [ ] No hay alerta de "cola con mensajes y cero consumidores". Es la deteccion
      que hubiera avisado a las 17:47 en vez de a las 03:00.

## Criterio para cerrar

Desplegado, y verificado en los logs del broker que tras un cierre de canal
aparece "Consumer restarted successfully after channel failure" sin reinicio del
contenedor.
