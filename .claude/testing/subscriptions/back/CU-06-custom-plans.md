# CU-06: Planes personalizados por negocio

## Objetivo
Verificar CRUD de `custom-plans`, exclusivo de super admin.

## Precondiciones
- SUPER_TOKEN y DEMO_TOKEN obtenidos (CU-01).

## 6.1 Crear plan personalizado para Demo
```
POST /api/v1/subscriptions/custom-plans
Authorization: Bearer {SUPER_TOKEN}
Content-Type: application/json

{
  "name": "Plan E2E Test",
  "code": "e2e-test-{TIMESTAMP}",
  "description": "Plan de prueba",
  "price": 10000,
  "billing_period": "monthly",
  "module_codes": ["orders", "shipments"],
  "max_ecommerce_channels": 1,
  "business_id": 26,
  "months": 1,
  "payment_reference": "E2E-REF",
  "notes": "Creado por prueba E2E"
}
```
- [ ] Status 200/201, `data.id` generado -> guardar CUSTOM_PLAN_ID
- [ ] `data.business_id == 26`

## 6.2 Listar planes personalizados
```
GET /api/v1/subscriptions/custom-plans
Authorization: Bearer {SUPER_TOKEN}
```
- [ ] Status 200, el plan creado en 6.1 aparece en el listado
- [ ] Si el listado no tiene filtro por business_id ni paginacion, evaluar si califica
      como catalogo <50 (probablemente si, dado que son pocos por negocio) o si es un
      hallazgo para el reporte de arquitectura

## 6.3 Actualizar plan personalizado
```
PUT /api/v1/subscriptions/custom-plans/{CUSTOM_PLAN_ID}
Authorization: Bearer {SUPER_TOKEN}
Content-Type: application/json

{
  "name": "Plan E2E Test Actualizado",
  "code": "e2e-test-updated",
  "price": 15000,
  "billing_period": "monthly",
  "module_codes": ["orders"],
  "max_ecommerce_channels": 1,
  "business_id": 26,
  "months": 1
}
```
- [ ] Status 200, `data.price == 15000`

## 6.4 Business normal NO puede crear/listar/editar custom-plans
Repetir 6.1 y 6.2 con `Authorization: Bearer {DEMO_TOKEN}`.
- [ ] Status 403 en ambos (si custom-plans no tiene `requireSuperAdmin`, es un
      hallazgo critico de seguridad: cualquier negocio podria crear planes con
      precio arbitrario para si mismo pasando su propio business_id)

## 6.4b Regresion del fix de ownership: otro negocio no puede comprar el custom plan de Demo
Requiere un segundo token de negocio distinto a Demo (o crear uno de prueba). Con ese
token, intentar:
```
POST /api/v1/subscriptions/purchase
Authorization: Bearer {OTRO_NEGOCIO_TOKEN}
Content-Type: application/json

{
  "subscription_type_id": {CUSTOM_PLAN_ID},
  "months": 1
}
```
- [ ] Status 400/404 (`ErrSubscriptionTypeNotFound`) -> antes del fix, cualquier negocio
      que conociera el ID del custom plan de Demo podia comprarlo con sus condiciones
      de precio pensadas para otro negocio

## 6.5 Eliminar plan personalizado (limpieza)
```
DELETE /api/v1/subscriptions/custom-plans/{CUSTOM_PLAN_ID}
Authorization: Bearer {SUPER_TOKEN}
```
- [ ] Status 200
- [ ] Verificar en DB que ya no aparece activo (soft delete esperado, no fisico)
