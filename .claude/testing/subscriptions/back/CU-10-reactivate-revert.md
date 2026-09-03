# CU-10: Reactivar suscripcion, revertir pago, historial y auditoria

## Objetivo
Verificar el boton "Reactivar" (independiente de dias de cortesia/edit-dates),
revertir un pago, y los endpoints de solo lectura de historial/auditoria.

## Precondiciones
- SUPER_TOKEN, DEMO_TOKEN obtenidos (CU-01).

## 10.1 Deshabilitar y reactivar manualmente
```
POST /api/v1/subscriptions/disable?business_id=26
POST /api/v1/subscriptions/reactivate?business_id=26
```
- [ ] disable: status 200, `business.subscription_status = cancelled`
- [ ] reactivate: status 200, `business.subscription_status = active`
- [ ] Auditoria `subscription_suspended` y `subscription_reactivated`

## 10.2 Revertir un pago
Usar el `id` de una fila de `business_subscriptions` reciente (no la mas reciente,
para no romper el estado usado por otras pruebas).
```
POST /api/v1/subscriptions/payments/{id}/revert
Authorization: Bearer {SUPER_TOKEN}
```
- [ ] Status 200
- [ ] La fila queda `status = reverted`
- [ ] Auditoria `payment_reverted`

## 10.3 Historial de pagos
```
GET /api/v1/subscriptions/payments/26
Authorization: Bearer {SUPER_TOKEN}
```
- [ ] Status 200, incluye la fila revertida con status `reverted`

## 10.4 Auditoria
```
GET /api/v1/subscriptions/audit-logs/26
Authorization: Bearer {SUPER_TOKEN}
```
- [ ] Status 200, trae todas las acciones registradas en CU-05, 08, 09 y 10

## 10.5 Seguridad: todos exclusivos de super admin
Repetir 10.1, 10.2, 10.3, 10.4 con DEMO_TOKEN.
- [ ] 403 en las 4
