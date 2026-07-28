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

1. Exponer un endpoint en el backend que descargue la etiqueta de MELI usando las
   credenciales de la integracion y la devuelva como PDF (proxy autenticado).
   Debe resolver el business/integracion desde el shipment, no desde el request,
   y validar pertenencia al negocio (ver `.claude/rules/multi-tenant-security.md`).
2. Que el front use ese endpoint para las guias de canal en vez del `guide_url`
   crudo. Hoy el boton "Ver Guia" abre la URL tal cual.
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
