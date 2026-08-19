# CU-05: Operaciones de super admin sobre suscripciones

## Objetivo
Verificar registro de pago, edicion de fechas y deshabilitar suscripcion — todas
protegidas con `requireSuperAdmin`.

## Precondiciones
- SUPER_TOKEN y DEMO_TOKEN obtenidos (CU-01).
- SUBSCRIPTION_TYPE_ID de CU-02.

## 5.1 Registrar pago para Demo (super admin)
```
POST /api/v1/subscriptions/register-payment
Authorization: Bearer {SUPER_TOKEN}
Content-Type: application/json

{
  "business_id": 26,
  "subscription_type_id": {SUBSCRIPTION_TYPE_ID},
  "months": 1,
  "payment_reference": "TEST-PAY-{TIMESTAMP}",
  "notes": "Pago de prueba E2E"
}
```
- [ ] Status 200, `data.business_id == 26`

## 5.2 Editar fechas de la suscripcion
```
PUT /api/v1/subscriptions/edit-dates
Authorization: Bearer {SUPER_TOKEN}
Content-Type: application/json

{
  "business_id": 26,
  "start_date": "2026-08-01",
  "end_date": "2026-12-31"
}
```
- [ ] Status 200
- [ ] `data.end_date` refleja "2026-12-31"

## 5.3 Un business normal NO puede llamar estos endpoints (403 esperado)
Repetir 5.1 y 5.2 con `Authorization: Bearer {DEMO_TOKEN}`.
- [ ] Status 403, `error == "super admin access required"` en ambos casos

## 5.4 Deshabilitar suscripcion
```
POST /api/v1/subscriptions/disable?business_id=26
Authorization: Bearer {SUPER_TOKEN}
```
- [ ] Status 200
- [ ] Verificar en DB que la suscripcion de Demo quedo marcada inactiva/deshabilitada
- [ ] **Regresion del fix de acceso no revocado**: inmediatamente despues, `GET /my-modules`
      con DEMO_TOKEN ya NO debe incluir los modulos del plan de Demo (antes del fix,
      `has_module_access.go` ignoraba el status y seguian apareciendo indefinidamente)
- [ ] Despues de este paso, `GET /me` (CU-03.1) para Demo ya no debe mostrar la
      suscripcion como activa -> re-ejecutar CU-03.1 y confirmar

## 5.5 Restaurar estado (dejar Demo funcional para no romper otras pruebas)
Volver a registrar un pago (5.1) para reactivar la suscripcion de Demo al finalizar,
o correr `edit-dates` extendiendo `end_date` a futuro si el disable no borra la fila.
Confirmar con CU-03.1 que Demo vuelve a tener suscripcion activa.
