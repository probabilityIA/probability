# CU-11: Marcar shipments como delivered

## Objetivo
Para los shipments que necesitan ser cobrados (CU-12, CU-13 y CU-15), forzar el estado a
`delivered` via `PUT /shipments/:id` (admitido por el binding del endpoint).

## Precondiciones
- CU-10 ejecutado: SHIPMENT_ID_1, SHIPMENT_ID_2, SHIPMENT_ID_4 disponibles.
- SHIPMENT_ID_3 se deja en `pending` para CU-14.

## Endpoint
```
PUT /api/v1/shipments/{SHIPMENT_ID}
Authorization: Bearer {DEMO_TOKEN}
Content-Type: application/json

{ "status": "delivered", "delivered_at": "2026-04-25T18:00:00Z" }
```

## Verificaciones (por shipment)
- [ ] Status 200
- [ ] `data.status == "delivered"`
- [ ] `data.delivered_at` poblado
- [ ] El shipment aparece en `GET /shipments/cod?status=delivered`
