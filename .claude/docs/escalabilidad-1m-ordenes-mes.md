# Mejoras a futuro: escalar a 1 millon de ordenes/mes

Fecha del analisis: 2026-08-05
Estado: diagnostico, nada implementado.

Escenario objetivo (ideal): 1.000.000 de ordenes/mes sostenidas.

## Dimensionamiento

- 1M ordenes/mes = ~23 ordenes/min promedio (0,4/s).
- Picos de campana (5-10x el promedio) = 2-4 ordenes/s.
- Cada orden abre fan-out a ~8 colas desde `orders.events`
  (invoicing, score, inventory, customers, shipments, meli, jumpseller, events)
  mas el downstream de cada una.
- Total: 10-20M mensajes/mes, picos de 40-80 msg/s.
- 33.000 ordenes/dia.

Conclusion: RabbitMQ como broker no es el cuello de botella a este volumen.
Los cuellos estan en el despliegue, las garantias de entrega y los terceros.

## Hallazgos (codigo revisado)

Archivos base: `back/central/shared/rabbitmq/rabbitmq.go`,
`back/central/shared/rabbitmq/queues.go`, `back/central/shared/db/db.go`,
`infra/compose-prod/docker-compose.yaml`.

### 1. Mensajes no persistentes (CRITICO, es un bug hoy)

`Publish` y `PublishToExchange` construyen `amqp.Publishing{ContentType, Body}`
sin `DeliveryMode: amqp.Persistent`. Las colas se declaran durable pero los
mensajes viajan transient: si RabbitMQ reinicia se pierde todo lo encolado.
A 33.000 ordenes/dia, cada reinicio del broker son miles de ordenes sin
facturar ni despachar.

### 2. Sin publisher confirms ni patron outbox

El publish ocurre fuera de la transaccion de la orden y sin confirms, asi que
`Publish` retorna `nil` aunque el mensaje nunca llegue al broker. Resultado:
orden persistida en DB, evento perdido, sin senal de error.

### 3. Sin DLQ y requeue infinito

`rabbitmq.go` linea ~308: `msg.Nack(false, true)` reencola indefinidamente.
Un mensaje envenenado gira a maxima velocidad, ocupa el prefetch (50) y quema
el unico CPU disponible. A este volumen aparece en dias, no en meses.
No existe ninguna declaracion de `x-dead-letter-exchange` en el repo.

### 4. Recursos de produccion: 1 vCPU / 512MB

`back-central` en `infra/compose-prod/docker-compose.yaml` tiene
`deploy.resources.limits: cpus "1.00", memory 512M`, en un solo contenedor.
Es un monolito: API HTTP + todos los consumers + todos los tickers
(`reconciliation_worker`, `retry_consumer` de pay e invoicing, `expiry_worker`
de subscriptions, agregado de geozones, bold 24h) en el mismo proceso.
Ya hubo un incidente de OOM en RDS con el volumen actual.

### 5. No escala horizontalmente tal como esta

Levantar replicas duplica los tickers (retries, reconciliacion, expiraciones)
porque no hay leader election ni lock distribuido. Hay que separar un binario
`worker` de un binario `api`, o poner locks en Redis antes de replicar.

### 6. Concurrencia real por cola = 1

Casi todos los modulos usan `Consume(...)`, que internamente es
`ConsumeConcurrent(..., 1)`: un solo worker secuencial por cola.
Con handlers de facturacion (3-8 llamadas HTTP externas, 2-5s por orden)
eso da ~0,3 ordenes/s en esa cola, por debajo del promedio requerido de 0,4/s
y muy lejos de los picos. Este es el cuello mas inmediato del codigo propio.
Excepcion: el consumer de SoftPymes si usa `ConsumeConcurrent` (3 workers).

### 7. Los terceros son el limite duro (ya observado)

Ver `.claude/alerts/softpymes-timeouts-masivo.md`: 3.800 facturas en un dia
tumbaron SoftPymes (406 respuestas 200 con body vacio, 230 timeouts en el POST
de factura, 193 timeouts en el login). 1M/mes son 33.000/dia, 10x eso.
Agravantes: `shared/httpclient` con `RetryCount: 2` sin check de idempotencia
triplica la carga justo cuando el tercero esta lento; no hay rate limiting,
backoff ni circuit breaker por integracion (SoftPymes, Siigo, EnvioClick,
Meli, WhatsApp).

### 8. Idempotencia parcial

Solo Siigo y SoftPymes implementan chequeo de idempotencia, y el de SoftPymes
solo corre cuando `operation == "retry"`. Requeue infinito + retry consumers +
handlers no idempotentes ya produjo el incidente de
`.claude/alerts/guias-duplicadas-doble-cobro.md` (dos guias reales y doble
debito de wallet). A 10x volumen eso escala en dinero perdido.

### 9. SPOF de infraestructura

RabbitMQ y Redis corren en el mismo EC2 que la app, sin cluster ni quorum
queues. El pool de PostgreSQL esta fijo en `MaxOpenConns(25)` / `MaxIdleConns(25)`
(`shared/db/db.go`), compartido entre API y todos los consumers del mismo proceso.

### 10. Canal unico de publicacion

Todos los publish del proceso van por `r.channel` (un solo canal AMQP
compartido). La libreria serializa internamente, asi que es seguro pero es un
punto de serializacion. No es el cuello a este volumen; si lo seria a 10x.

## Plan por fases

### Fase 1 - Evitar perdida de datos (hacer ya, independiente del volumen)

1. Agregar `DeliveryMode: amqp.Persistent` en `Publish` y `PublishToExchange`.
2. Habilitar publisher confirms (`channel.Confirm(false)` + `NotifyPublish`)
   y propagar el error si el broker no confirma.
3. DLQ: declarar colas con `x-dead-letter-exchange` y limite de reintentos;
   eliminar el `Nack(requeue=true)` ciego. Mensaje agotado va a la DLQ,
   no vuelve a la cola principal.

### Fase 2 - Capacidad de proceso (hasta ~300k ordenes/mes)

4. Separar `api` y `worker` en dos servicios; los tickers solo en `worker`
   (o proteger cada ticker con lock en Redis para poder replicar).
5. Subir `back-central` a 2-4 vCPU y 2-4GB en compose-prod.
6. Cambiar los `Consume(...)` de las colas calientes a `ConsumeConcurrent`
   con numero de workers configurable por env (invoicing, shipments,
   inventory, webhooks).
7. Subir/parametrizar el pool de DB por rol (api vs worker) y revisar el
   sizing de RDS despues del OOM.

### Fase 3 - Sostener 1M/mes

8. Rate limiter + backoff exponencial + circuit breaker por integracion
   externa. Ya existe `shared/ratelimit/limiter.go`, no esta aplicado a los
   clientes de facturacion ni de transportadoras.
9. Idempotencia general por `event_id` en todos los consumers
   (tabla de eventos procesados, o SETNX en Redis con TTL).
10. Outbox transaccional para los eventos de orden: escribir el evento en la
    misma transaccion que la orden y publicar desde un relay.
11. Quorum queues y RabbitMQ fuera del EC2 de aplicacion (cluster propio o
    Amazon MQ). Lo mismo para Redis (ElastiCache).
12. Metricas por cola (profundidad, edad del mensaje mas viejo, tasa de DLQ)
    en el Grafana que ya existe en compose-prod, con alertas.

## Criterio de exito

- Reinicio de RabbitMQ sin perdida de mensajes.
- Ningun mensaje reencolado mas de N veces; DLQ observable y vacia en operacion normal.
- Backlog de `orders.events.*` por debajo de 60s de antiguedad en pico.
- Jobs masivos de 30.000+ facturas/dia sin tumbar al proveedor ni generar duplicados.
