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

## Evidencia en shipment_sync_logs (2026-08-21)

Consulta: logs `provider='envioclick'`, `status='success'`,
`request_url LIKE '%/api/v2/shipment%'` agrupados por `myShipmentReference`.
Solo esa URL crea guia; `/api/v2/quotation` tambien se registra con
`operation_type='generate'` y por eso el conteo crudo de "generate" engana
(3 quotations + 1 shipment por guia normal).

Ocho referencias con mas de una creacion real, cada una con `idOrder` distinto
de EnvioClick (guias reales separadas):

| Orden | Negocio | idOrder EnvioClick | Separacion |
|---|---|---|---|
| MYS-0118 | Mystic Rose | 4385504 / 4385509 | 1 min |
| MYS-0157 | Mystic Rose | 4406107 / 4406110 | 1 min |
| MYS-0232 | Mystic Rose | 4424912 / 4424915 | 1 min |
| MYS-0302 | Mystic Rose | 4476735 / 4476740 | 1 min |
| VIG-0010 | Viga | 4491102 / 4491506 | 1 h 38 |
| FIS-0061 | LaPercha | 4501053 / 4576423 / 4593773 | dias |
| FIS-0066 | LaPercha | 4580801 / 4584210 | 2 dias |
| MYS-0718 | Mystic Rose | 4632833 / 4632845 | 1 min 44 |

Las de dias (FIS-0061, FIS-0066) son regeneraciones legitimas tras cancelar.
Las de minutos son el bug.

Doble debito confirmado en `transaction` (USAGE/GUIDE) para MYS-0118 (21.539 x2),
MYS-0157 (19.980 x2), MYS-0232 (15.598 x2), MYS-0302 (16.353 x2) y VIG-0010
(25.051,15 + 25.420,90). Ninguna cancelacion reembolso.

### MYS-0718 - caso nuevo del 18/08, el guard NO lo cubre

shipment **43991**, usuario 20, dos POST `/api/v2/shipment` con 200:

- 19:30:49 -> idOrder 4632833, tracker **2284436982**, PDF emitido
- 19:32:33 -> idOrder 4632845, tracker **2284436983**, PDF emitido

El shipment quedo con `tracking_number = 2284436983` y una sola transaccion de
wallet (19.160 a las 19:32:34). La primera guia **no dejo transaccion ni se
guardo en el shipment**: quedo huerfana en EnvioClick, con recoleccion
programada y costo que Probability paga sin cobrar al negocio.

Diferencia con los casos viejos: aqui no se creo un shipment nuevo, se repitio
sobre el MISMO shipment. `shipmentHasActiveGuide` mira los shipments de la orden
en BD; como la primera creacion no persistio tracking ni guide_url, el guard no
vio nada y dejo pasar la segunda. Emparenta con
`.claude/bitacora/2026-08-17-cotizacion-guia-generada-fantasma.md`.

Falta: idempotencia real contra el proveedor (marcar el shipment "en generacion"
antes de llamar a EnvioClick, o reconciliar por `myShipmentReference` con
`GET /shipment` antes de crear). El guard actual no protege cuando la respuesta
del carrier no alcanza a persistirse.
