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

## Pendiente de fondo

Esto es una mitigacion en el handler. La solucion completa es limite de
reintentos + backoff + DLQ en `shared/rabbitmq`, que hoy no existe. Ver
`.claude/alerts/consumidor-muerto-por-canal-cerrado.md` y la Fase 1 de
`.claude/docs/escalabilidad-1m-ordenes-mes.md`.
