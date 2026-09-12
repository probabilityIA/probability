# Bucle caliente en orders.whatsapp.confirmed

**Ticket:** sin ticket (la API de tickets exige JWT de super admin; queda como
pendiente crearlo y pegar este resumen).
**Negocios afectados:** 66 (Sigo-Isabel_Rojas) como origen del mensaje muerto,
46 (VIG) como damnificado. El desgaste lo sufrio todo el backend.
**Canal:** WhatsApp / colas.

## Resumen

Un mensaje que nunca iba a poder procesarse giraba en la cola
`orders.whatsapp.confirmed` a **675 iteraciones por segundo**, quemando 61% de
CPU del backend y borrando por rotacion todos los logs del contenedor. Es la
**tercera vez** que pasa lo mismo por la misma causa de fondo, y la segunda en
dos dias.

| Fecha | Cola | Volumen |
|---|---|---|
| 2026-08-06 | `wallet.balance_alert.requested` | ~1.542 iteraciones/s |
| 2026-09-10 | `probability.orders.canonical` | 112.356 reintentos en 6 h |
| 2026-09-11 | `orders.whatsapp.confirmed` | 675 iteraciones/s |

Las dos primeras se cerraron parcheando el caso puntual (clasificar el handler,
correr la migracion faltante). Esta vez se ataco la causa comun.

## Sintoma

Se encontro de casualidad, buscando OTRA cosa: el error de conexion del numero
de WhatsApp de LaPerchaDel10 no aparecia por ningun lado en los logs.

Y no aparecia porque los logs solo cubrian **12 minutos**:

```
docker logs --since 24h central_reserve_prod_green | wc -l   -> 1.413.725
primera linea disponible: 09-11 19:51:59
ultima linea:             09-11 20:03:41
tamano del directorio de logs del contenedor: 444 MB
```

Tres lineas por iteracion, repetidas sin parar:

```
INF Processing order confirmation from WhatsApp  order_number=PRUEBA-ISABEL phone_number=573192611891
ERR Error getting order for confirmation         error="order not found" business_id=66
ERR Error processing message                     error="order not found" queue=orders.whatsapp.confirmed
```

Medido en vivo:

| Medida | Valor |
|---|---|
| Iteraciones por segundo | 675 (563 en una segunda medicion) |
| CPU del contenedor del backend | 58-61% |
| Load average del `t4g.small` (2 vCPU) | 2,91 |
| Mensajes en la cola | 2 |

## Diagnostico

`handleConfirmed`
(`services/modules/orders/internal/infra/primary/queue/whatsapp_consumer.go`)
devolvia el error crudo cuando la orden no existia. `shared/rabbitmq` hacia
`msg.Nack(false, true)`: reencolar **sin limite de reintentos, sin backoff y sin
DLQ**. El mensaje volvia al instante y fallaba otra vez.

La orden `PRUEBA-ISABEL` no existe en ninguna parte:

```sql
SELECT id, order_number, business_id FROM orders WHERE order_number ILIKE '%ISABEL%';
-- 0 filas
SELECT id, name FROM business WHERE id = 66;
-- 66 | Sigo-Isabel_Rojas
```

El negocio si existe; la orden no. Fue un mensaje de prueba con un
`order_number` inventado.

### Lo que casi sale mal

La reaccion obvia era purgar la cola. **Habria sido un error.** Los 2 mensajes
NO eran el mismo:

| Mensaje | Orden | Negocio | Existe |
|---|---|---|---|
| El que giraba | `PRUEBA-ISABEL` | 66 | no |
| El de atras | `VIG-0158` | 46 | **si**, con `is_confirmed = false` |

`VIG-0158` es la confirmacion real de un cliente, atascada detras del mensaje
envenenado. Purgar la cola la habria borrado en silencio.

Tampoco se podia sacar solo el envenenado: esta permanentemente **unacked** en
manos del consumidor, y `rabbitmqadmin get` solo devuelve mensajes en estado
*ready*, por eso devolvia el otro. La unica salida limpia era desplegar el fix,
que al reiniciar el consumidor le redelivera el mensaje al codigo nuevo.

### Hipotesis descartadas

- **"Los logs no llegan porque el filtro de logtail no matchea".** Falso: el
  filtro deja pasar `ERROR`, y estas lineas son `ERR`. El problema era la
  rotacion por volumen.
- **"El envio a CloudWatch esta bien y la alarma no salto porque el volumen no
  llego al umbral".** Falso, y es un hallazgo aparte: el grupo
  `/probability/back-central` **no recibe un evento desde 2026-09-02 17:12**
  (9 dias). La alarma `logs-ingesta-alta` nunca tuvo que medir.
- **"Purgar la cola es seguro porque son mensajes de prueba".** Falso, ver
  arriba.

## Causa raiz

Dos niveles, y el segundo es el que importa:

1. **Inmediato:** el handler no distinguia error permanente de transitorio, lo
   que viola `.claude/rules/colas-errores-permanentes.md`.
2. **De fondo:** `shared/rabbitmq` reencolaba sin limite. **De 79 archivos
   consumidores, solo 10 clasificaban.** Los otros 69 podian producir el mismo
   bucle en cualquier momento. Parchear el handler era la respuesta que ya se
   habia dado el 2026-08-06 y el 2026-09-10 (ver la tabla del resumen), y por
   eso volvio a pasar. Mientras `shared/rabbitmq` reencolara sin limite, arreglar
   el handler de turno solo compraba tiempo hasta el siguiente.

## Correccion

`863c497e`:

- **`shared/rabbitmq/retry.go` (nuevo).** Al fallar un handler, republica a una
  cola de espera con TTL creciente (10 s, 60 s, 300 s) y ACKea el original. Al
  quinto intento va a `probability.dlq` y deja de circular. Cabeceras
  `x-origin-queue`, `x-retry-count`, `x-last-error`, `x-first-failed-at`.
- 3 exchanges fanout + 3 colas de espera + 1 DLQ **compartidas** por las 83
  colas. Un juego por cola serian 332 colas en el `t4g.small`.
- Las colas principales **no cambian de argumentos**: declararlas con
  `x-dead-letter-exchange` habria roto el arranque con
  `PRECONDITION_FAILED - inequivalent arg`, porque ya existen en produccion
  declaradas con `nil`.
- Si la republicacion falla, cae al `Nack(requeue)` de siempre: mejor un mensaje
  que gira que un mensaje perdido.
- Los tres handlers del consumidor de orders clasifican, y el repositorio
  devuelve el centinela `ErrOrderNotFound` en vez de un `fmt.Errorf` suelto.

### El detalle que hace que 4 colas sirvan para 83

Las colas de espera se declaran con `x-message-ttl` y
`x-dead-letter-exchange: ""`, **sin** `x-dead-letter-routing-key`. Al vencer el
TTL, RabbitMQ usa la routing key propia del mensaje, que es el nombre de la cola
original, y el exchange por defecto lo devuelve ahi solo.

Para que esa routing key sobreviva, la republicacion va a un exchange **fanout**
(ignora la routing key al enrutar, pero la conserva en el mensaje). Publicar al
exchange por defecto la habria pisado con el nombre de la cola de espera y el
mensaje no habria vuelto nunca.

## Verificacion

Contra un RabbitMQ real (local), no solo tests unitarios
(`shared/rabbitmq/retry_integration_test.go`, se salta solo si no hay broker):

| Prueba | Resultado |
|---|---|
| Mensaje venenoso | **2 entregas en 20 s**, contra 675 por segundo |
| Mensaje bueno detras del venenoso | se procesa, no queda atrapado |
| Al agotar los intentos | llega a la DLQ con el cuerpo intacto, su cola de origen y el ultimo error; **no hay una entrega mas** |

`go build ./...`, `go vet` y la suite completa del backend: 0 fallos.

## Pendientes

1. **Alarma sobre la profundidad de `probability.dlq`.** Sin ella los mensajes
   muertos se acumulan sin que nadie mire.
2. **Arreglar el envio de logs a CloudWatch**, caido desde 2026-09-02. Mientras
   siga asi, el proximo incidente tampoco avisa.
3. **Clasificar los 69 consumidores restantes.** Ya no tumban el servidor, pero
   un error permanente sin clasificar gasta 5 vueltas y 11 minutos, y ensucia la
   DLQ con ruido.
4. **Que hacer con lo que cae en la DLQ**: no hay reproceso ni purga.
5. **Crear el ticket** y pegar este resumen.
