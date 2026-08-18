# Backend en crash loop por orden de negocio soft-deleted (panic nil en mapper)

2026-08-18. El backend de prod (`central_reserve_prod`) estuvo caido en crash
loop ~1.5 h (33+ reinicios). El login y toda la API respondian el 502 HTML de
nginx; el front lo mostraba como `Unexpected token '<' ... is not valid JSON`.

## Sintoma

- Login roto en www.probabilityia.com.co con error de parseo JSON.
- `/health` respondia `502 Bad Gateway` de nginx.
- Contenedor `central_reserve_prod`: `Restarting (2)` en loop, exit code 2.

## Causa raiz

1. 10:32 -0500: se hizo soft delete del negocio 58 "Gamers Paradise"
   (`business.deleted_at` quedo con fecha).
2. 11:35:45: webhook de Shopify (`zrmr0x-af.myshopify.com`) creo la orden
   WEB-1018 (`f695d686-7dc1-4c90-a9f8-66613eee11dc`) con `business_id = 58`.
3. 11:35:48: webhook `orders/updated` de la misma orden se publico a
   `probability.orders.canonical`.
4. El consumer llamo `MapAndSaveOrder` -> `GetOrderByExternalID`, que hace
   `Preload("Business")`. GORM excluye registros soft-deleted, asi que
   `o.Business` quedo `nil`.
5. `mappers/mapper.go:152` (`BusinessName: o.Business.Name`) -> panic nil
   pointer -> en Go un panic sin `recover` mata el proceso entero.
6. El mensaje nunca se ACKeo: se reencolo, Docker reinicio el contenedor, el
   consumer tomo el mismo mensaje y panic otra vez. Loop infinito.

Un solo mensaje venenoso tumbo toda la API para todos los negocios.

## Hipotesis descartadas

- Nginx: estaba sano, solo mostraba el 502 del backend caido. Reiniciarlo no
  sirve de nada en este escenario.
- Recursos (memoria/CPU/RDS): no fue. Exit code 2 = panic de Go, no OOM.

## Solucion aplicada (inmediata)

- Descartar el mensaje venenoso de `probability.orders.canonical`
  (`rabbitmqadmin get ... ackmode=ack_requeue_false`) y
  `docker compose restart back-central`. La orden ya estaba persistida, no se
  perdio informacion.
- Ocurrio 2 veces mas: Shopify siguio mandando webhooks de la misma orden
  (la integracion 230 sigue activa aunque el negocio este eliminado). Hubo
  que descartar 3 mensajes en total, todos de WEB-1018.

## Fix de codigo (pendiente de deploy)

- `orders/internal/app/usecasecreateorder/map_and_save_order.go`: al inicio de
  `MapAndSaveOrder` se valida el negocio con `GetBusinessNameByID` (filtra
  `deleted_at IS NULL`); si no existe o esta eliminado retorna
  `domainerrors.ErrOrderBusinessDeleted`.
- `orders/internal/domain/errors/errors.go`: nuevo sentinel
  `ErrOrderBusinessDeleted`.
- `orders/internal/infra/primary/queue/consumer.go`: el error se clasifica
  como permanente -> `Warn` + `publishRejected` + `return nil` (ACK, no
  reencola), siguiendo `.claude/rules/colas-errores-permanentes.md`.
- `orders/internal/infra/secondary/repository/mappers/mapper.go`: guard de nil
  para `o.Business` en `ToDomainOrder` (defensa para todos los demas callers).

Sin este deploy, cualquier webhook de una orden existente de un negocio
eliminado vuelve a tirar el backend.

## No se perdio data

- Los mensajes de RabbitMQ persisten mientras el backend esta caido; se
  procesaron al volver (147 facturas Softpymes quedaron encoladas y se
  drenaron despues).
- Shopify reintenta webhooks fallidos automaticamente (~48 h).

## Lecciones

- Un panic en cualquier goroutine de un consumer mata todo el proceso. Sigue
  pendiente el `recover` + retry limit + DLQ en `shared/rabbitmq` (ver alerta
  `consumidor-muerto-por-canal-cerrado.md`).
- Eliminar un negocio no detiene sus integraciones: Shopify sigue mandando
  webhooks y el flujo de ordenes seguia creando registros con ese
  `business_id`.
