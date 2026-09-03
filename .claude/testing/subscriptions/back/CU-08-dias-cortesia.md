# CU-08: Dias de cortesia (subscription_courtesy_until)

## Objetivo
Verificar que ExtendCourtesy pospone solo la fecha de bloqueo (courtesy_until),
sin mover el rango de facturacion (business_subscriptions.start_date/end_date ni
business.subscription_end_date), y que reactiva la cuenta si estaba expired.

## Precondiciones
- SUPER_TOKEN obtenido (CU-01).
- Negocio de prueba con subscripcion activa y con end_date conocido.

## 8.1 Extender cortesia con cuenta activa: NO debe tocar el rango de facturacion
```
POST /api/v1/subscriptions/extend-days
Authorization: Bearer {SUPER_TOKEN}
Content-Type: application/json

{
  "business_id": 26,
  "days": 5,
  "reason": "prueba E2E cortesia"
}
```
- [ ] Status 200
- [ ] Verificar en DB: `business.subscription_end_date` y `business_subscriptions.end_date`
      NO cambiaron
- [ ] Verificar en DB: `business.subscription_courtesy_until` quedo en
      `max(fecha_de_corte_actual, ahora) + 5 dias`
- [ ] Se registro auditoria `courtesy_extended`

## 8.2 Extender cortesia con cuenta expired: debe reactivar
Preparar: marcar el negocio como expired manualmente (o via disable) antes de este paso.
- [ ] Antes: `business.subscription_status = expired`
- [ ] Ejecutar 8.1 de nuevo (mismo body)
- [ ] Despues: `business.subscription_status = active`
- [ ] Se registro auditoria adicional `subscription_reactivated`
- [ ] `GET /my-modules` con un token de ese negocio vuelve a mostrar los modulos del plan

## 8.3 Dias invalidos
```
POST /api/v1/subscriptions/extend-days
{ "business_id": 26, "days": 0, "reason": "x" }
```
- [ ] Status 400

## 8.4 Business normal no puede llamar el endpoint
Repetir 8.1 con `Authorization: Bearer {DEMO_TOKEN}`.
- [ ] Status 403
