# Resultados COD (2026-04-25 21:17:24)

PASS: 32  FAIL: 0

- OK   - demo token obtenido (325 chars)
- OK   - super admin token obtenido
- OK   - 8.1 listado base 200 con data array
- OK   - 8.2 filtro status=delivered
- OK   - 8.3 filtro is_paid=false
- OK   - 8.5 super admin sin business_id -> 400
- OK   - 8.5 super admin con business_id=26 -> 200
- OK   - 9.1 creada orden N=1 cod=50000 -> ea2330f2-4a0b-44a1-b047-f408fa765162 (prob-0057)
- OK   - 9.2 creada orden N=2 cod=1500000 -> cbad4269-ee2d-4aa3-bc5d-91dfd59a86fc (prob-0058)
- OK   - 9.3 creada orden N=3 cod=250000 -> 53a4db6c-5be7-4b7c-9f08-6bb06a5f612c (prob-0059)
- OK   - 9.4 creada orden N=4 cod=320000 -> 47fb6221-0d0e-4185-8856-8aac25f16029 (prob-0060)
- OK   - 9.5 creada orden N=5 cod=0 -> ea73f28c-045e-4ed1-99d6-627b1f3491ff (prob-0061)
- OK   - 10.1 shipment_id=34217 generado
- OK   - 10.2 shipment_id=34218 generado
- OK   - 10.3 shipment_id=34219 generado
- OK   - 10.4 shipment_id=34220 generado
- OK   - 10.1 shipment 34217 tracking=SRV-128489431
- OK   - 10.2 shipment 34218 tracking=IRP-844334889
- OK   - 10.3 shipment 34219 tracking=TCC-332406552
- OK   - 10.4 shipment 34220 tracking=797468111132
- OK   - 11.1 shipment 34217 -> delivered
- OK   - 11.2 shipment 34218 -> delivered
- OK   - 11.4 shipment 34220 -> delivered
- OK   - 12 cobro 50k OK is_paid=true
- OK   - 13 cobro 1.5M cod_total=1500000
- OK   - 14 rechazo 400 con mensaje delivered
- OK   - 15.1 primer cobro OK
- OK   - 15.2 segundo cobro 400 already paid
- OK   - 16.1 shipment manual creado id=34221
- OK   - 16.2 rechazo 400 not COD
- OK   - 17.1 lista is_paid=true incluye SHIPMENTS 1,2,4
- OK   - 17.2 lista is_paid=false NO incluye SHIPMENTS pagados

## IDs generados
- N=1  order=ea2330f2-4a0b-44a1-b047-f408fa765162  number=prob-0057  shipment=34217  tracking=SRV-128489431
- N=2  order=cbad4269-ee2d-4aa3-bc5d-91dfd59a86fc  number=prob-0058  shipment=34218  tracking=IRP-844334889
- N=3  order=53a4db6c-5be7-4b7c-9f08-6bb06a5f612c  number=prob-0059  shipment=34219  tracking=TCC-332406552
- N=4  order=47fb6221-0d0e-4185-8856-8aac25f16029  number=prob-0060  shipment=34220  tracking=797468111132
- N=5  order=ea73f28c-045e-4ed1-99d6-627b1f3491ff  number=prob-0061  shipment=  tracking=
