# ALERTA: Guias duplicadas con doble cobro de wallet

Fecha: 2026-06-18

## Contexto

Una orden con auto-generacion de guia activa terminaba con DOS guias reales en
el carrier (dos tracking + dos PDF) y DOS debitos de wallet por el mismo envio.

Causa raiz: el handler manual `GenerateGuide`
(`shipments/internal/infra/primary/handlers/generate-guide.go`) solo reutilizaba
un shipment "pending sin tracking/sin guide_url". Cuando la auto-generacion
(`order_created_autogen.go`) ya habia creado la guia (con tracking + guide_url),
el loop de reuso no encontraba nada reutilizable, creaba un shipment NUEVO y
disparaba una SEGUNDA generacion en EnvioClick -> segundo cobro.

Cancelar una guia (`handleCancelResponse`) NO reembolsa el wallet, asi que el
doble cobro queda.

## Fix aplicado (codigo)

- `generate-guide.go`: si la orden ya tiene una guia activa (tracking o
  guide_url y status != cancelled/failed) responde 409 y NO genera otra.
  Helper `shipmentHasActiveGuide`.
- `order_created_autogen.go` `createOrReusePendingShipment`: mismo guard,
  aborta si ya existe guia activa.

Pendiente de commit + push (esperar autorizacion). Una vez en prod, este bug de
codigo queda cerrado.

## Items

### Urgente (datos / dinero) - reembolsos manuales

Ordenes con 2 guias activas y doble debito (verificar y reembolsar la guia
sobrante / cancelar la duplicada en el carrier):

- order `4ad65800-f74e-4f0b-abf2-f55007dfe9fd` (MYS-0302, business 36): shipments
  34678 (cancelada el 18/06, tracking 034057376067) y 34679 (034057376085).
  Ambos debitos de 16.353 siguen en `transaction` (USAGE/GUIDE). 34678 cancelada
  pero SIN reembolso -> falta acreditar 16.353 al wallet fee6cd0c.
- order `96679526-552f-42bc-b15c-e159c17172d5`: shipments 34631 / 34632
  (trackings 84151641204 / 84151641209). Revisar doble cobro.
- order `79edc7eb-165a-45dd-b3c9-d82638dc17bd`: shipments 34247 / 34248
  (trackings 888006801990 / 84151592680). Revisar doble cobro.

### Importante

- El handler de cancelacion deberia reembolsar (credit) el wallet al cancelar
  una guia ya cobrada, hoy no lo hace.

## Criterio para cerrar

Fix en prod + reembolsos/cancelaciones de las 3 ordenes resueltos y verificados
en `transaction`.
