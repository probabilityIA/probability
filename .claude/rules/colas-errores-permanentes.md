# Colas: distinguir error permanente de error transitorio

Regla OBLIGATORIA para todo handler de consumidor de RabbitMQ.

## El problema que evita

`shared/rabbitmq` hace `msg.Nack(false, true)` cuando el handler devuelve error:
reencola **sin limite de reintentos, sin backoff y sin DLQ**. Si el error es
permanente, el mensaje vuelve al instante y el handler falla otra vez. El
resultado es un bucle caliente.

Incidente real (2026-08-06): la cola `wallet.balance_alert.requested` tenia 2
mensajes para un negocio sin integracion de WhatsApp configurada. El handler
devolvia error, se reencolaba, y el bucle corria a **~1.542 iteraciones por
segundo** de forma indefinida: quemaba CPU, le pegaba a RDS ~1.500 veces por
segundo, y genero el **100% del volumen de logs** del backend (4,2 MB/min), lo
que ademas enterraba cualquier otro evento util en el log.

## La regla

Antes de devolver un error desde un handler de cola, clasificarlo:

| Tipo | Ejemplos | Que hacer | Efecto |
|------|----------|-----------|--------|
| **Permanente** | integracion no configurada, credenciales ausentes, plantilla inexistente, telefono invalido, payload malformado, codigos definitivos del proveedor | `log.Warn()` + `return nil` | ACK, se descarta |
| **Transitorio** | proveedor caido, timeout, 5xx, rate limit, base de datos inaccesible | `log.Error()` + `return err` | Nack, se reencola |

Criterio: **si repetir la operacion identica no puede cambiar el resultado, es
permanente.** Un negocio sin WhatsApp configurado no va a tener WhatsApp
configurado por reintentar; Meta caida si se puede recuperar sola.

## Ruido en el log

Un descarte permanente va con `Warn`, una sola linea, y no se repite porque el
mensaje se ACKea. Nunca dejar que un caso permanente escriba en `Error` en bucle:
ademas de inutil, tapa los eventos que si importan.

## Como aplicarla en WhatsApp

Ya existe el clasificador compartido, usarlo en vez de duplicar listas:

```go
import whaErrors ".../whatsapp/internal/domain/errors"

if err != nil {
    if whaErrors.IsNonRetryable(err) {
        c.log.Warn().Err(err).Uint("business_id", ev.BusinessID).
            Msg("... skipped - non-retryable error (ACK)")
        return nil
    }
    c.log.Error().Err(err).Msg("... - will be retried")
    return err
}
```

`IsNonRetryable` vive en `whatsapp/internal/domain/errors/retryable.go` y lo usan
`consumerorder`, `consumershipment` y `consumerwalletalert`. Si aparece una frase
o codigo permanente nuevo, se agrega ahi, en un solo lugar.

## Para otros modulos

No hay clasificador compartido fuera de WhatsApp todavia. Mientras no exista,
cada handler debe hacer la distincion explicitamente. Prohibido devolver el error
crudo del proveedor sin clasificarlo: es exactamente lo que produjo el incidente.

## Violaciones criticas

- Handler que devuelve `err` ante una configuracion ausente o un payload que
  nunca va a poder procesarse.
- Descartar un mensaje (`return nil`) ante un fallo transitorio: se pierde el
  evento en silencio. Ante la duda, reintentar es lo correcto.
- Duplicar la lista de frases no reintentables en vez de usar el clasificador.

## La red de seguridad: backoff + DLQ (desde 2026-09-11)

`shared/rabbitmq` ya NO hace `Nack(requeue: true)` a secas. Al fallar un handler,
`handleFailedMessage` (`shared/rabbitmq/retry.go`) republica el mensaje a una cola
de espera y ACKea el original:

| Intento | Espera | A donde va |
|---|---|---|
| 1 | 10 s | `probability.retry.10s` |
| 2 | 60 s | `probability.retry.60s` |
| 3 | 300 s | `probability.retry.300s` |
| 4 | 300 s | `probability.retry.300s` |
| 5 | - | `probability.dlq`, y se ACKea |

Un mensaje muerto llega a la DLQ en ~11 minutos y deja de circular. **El bucle
caliente es ahora imposible por diseno**, clasifique bien el handler o no.

Como funciona: las colas de espera se declaran con `x-message-ttl` y
`x-dead-letter-exchange: ""`, **sin** `x-dead-letter-routing-key`. Al vencer el
TTL, RabbitMQ usa la routing key propia del mensaje, que es el nombre de la cola
original, y el exchange por defecto lo devuelve ahi solo. Para que esa routing
key sobreviva, la republicacion va a un exchange **fanout** (que la ignora al
enrutar pero la conserva en el mensaje), no al exchange por defecto.

Son 3 colas de espera + 1 DLQ compartidas por las 83 colas del sistema, no un
juego por cola: en el `t4g.small` serian 332 colas.

**Las colas principales NO cambiaron de argumentos.** Declararlas con `x-dead-
letter-exchange` habria roto el arranque con `PRECONDITION_FAILED - inequivalent
arg`, porque ya existen en produccion declaradas con `nil`.

Si la republicacion falla (broker caido), se cae al `Nack(requeue: true)` de
siempre: es preferible un mensaje que gira a un mensaje perdido.

En la DLQ cada mensaje lleva `x-origin-queue`, `x-retry-count`, `x-last-error` y
`x-first-failed-at`, asi que se sabe de donde vino y por que murio.

## Clasificar sigue importando

La DLQ evita el incendio, no hace bien el trabajo. Un error permanente sin
clasificar da cinco vueltas y 11 minutos de espera antes de morir, y ensucia la
DLQ con ruido que tapa los mensajes que si hay que revisar. La tabla de arriba
sigue siendo obligatoria; lo que cambia es que equivocarse ya no tumba el
servidor.

## Pendiente

- **Alarma sobre la profundidad de `probability.dlq`.** Sin ella los mensajes
  muertos se acumulan sin que nadie mire. Hoy ademas el envio de logs a
  CloudWatch esta caido desde 2026-09-02, asi que no hay por donde enterarse.
- **Que hacer con lo que cae en la DLQ**: no hay reproceso ni purga; se revisa a
  mano con `rabbitmqctl`.
