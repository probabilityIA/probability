# Alerta: MercadoLibre - la guia PDF no se puede abrir desde la plataforma

Fecha: 2026-07-28
Modulos: `back/central/services/integrations/ecommerce/meli`, `back/central/services/modules/shipments`

## Contexto

Las ordenes de MercadoLibre llegan sin guia y MELI la asigna minutos despues, por
la ruta de actualizacion. Hasta el 2026-07-28 esa ruta solo copiaba tracking y
guia a columnas de `orders` y nunca creaba la fila en `shipments`, asi que las
ordenes del canal quedaban fuera de Envios, manifiestos y recaudo (0 de 2 ordenes
meli de los ultimos 60 dias tenian envio).

Eso ya quedo resuelto (commit `656d7ed2`): el modulo `shipments` escucha
`order.updated` en `orders.events.shipments` y materializa el envio con
`GetOrderExternalGuide`, con guardas para no interferir con la generacion propia
de EnvioClick.

Lo que falta es poder VER esa guia. El `guide_link` que guardamos apunta a
`https://api.mercadolibre.com/shipment_labels?shipment_ids=<id>&response_type=pdf`,
que exige el access token de la integracion. Abierto directo desde el navegador
devuelve 401, asi que el boton "Ver Guia" del modulo de Envios no funciona para
este canal.

Orden de referencia: `2000017632602160` (business 37, integracion 250),
shipment_id de MELI `47626453039`, tracking `MEL47626453039FMDOF01`.

## Items

### IMPORTANTE

1. [RESUELTO 2026-07-28] Endpoint proxy autenticado
   `GET /api/v1/integrations/meli/shipments/:shipment_id/label` en el modulo meli.
   Resuelve business/integracion desde el shipment (nunca del request), exige que
   el envio sea de MercadoLibre, valida pertenencia al negocio y pide
   `business_id` por query a super admin. Refresca el token con `EnsureValidToken`
   y devuelve el PDF.
2. [RESUELTO 2026-07-28] El front pasa por `/internal/meli-label/[id]`
   (`ui/utils/guide-link.ts` detecta las URLs de `api.mercadolibre.com` y
   reescribe el href). Aplicado en `ShipmentList` y `CODShipmentList`.
   `TrackingPanel.tsx` quedo sin tocar: hoy no lo usa ningun componente.
3. Considerar guardar el PDF en S3 la primera vez (como se hace con las guias
   propias) para no depender del token ni de la disponibilidad de la API de MELI
   cada vez que alguien quiere imprimir.

### DESEABLE

4. La orden `2000017632602160` no se materializa porque esta cancelada: la guarda
   de `HasGuide()` excluye ordenes canceladas para no ensuciar Envios y
   manifiestos. Si se quiere el registro historico, hace falta una migracion
   puntual.
5. `orders.carrier` sigue quedando en null para ordenes de canal: la ruta de
   update no copia el carrier del DTO porque la entidad `ProbabilityOrder` no
   tiene ese campo. El envio materializado si trae carrier (derivado de la
   plataforma), pero la columna de la orden queda vacia.

## Criterio para cerrar

Se cierra cuando desde el modulo de Envios se pueda abrir e imprimir la guia de
una orden de MercadoLibre sin errores de autenticacion, verificado con una orden
viva del canal.

Estado 2026-07-28: falta esa verificacion. Hoy NO existe ningun shipment de
MercadoLibre en la base (la unica orden del canal con guia esta cancelada y por
eso no se materializa), asi que el camino contra la API real de MELI todavia no
se ha ejercido. Lo verificado es el contorno: ruta, JWT, regla de super admin,
guarda de canal (un envio que no sea de MELI responde 404) y el proxy del front.
Al llegar la proxima orden viva de MercadoLibre, abrir su guia y cerrar la alerta
si el PDF sale bien.
