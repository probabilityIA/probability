# CU-08: Listado COD y filtros

## Precondiciones
- Backend en `http://localhost:3050` con endpoint `GET /api/v1/shipments/cod`
- Demo business id=26 con al menos una orden con `cod_total > 0` y un shipment asociado.

## Caso 8.1: Listado base (usuario business demo)

```
POST /api/v1/auth/login        (X-Client-Type: api)  -> DEMO_TOKEN
GET  /api/v1/shipments/cod?page=1&page_size=10
Authorization: Bearer {DEMO_TOKEN}
```

Verificaciones:
- [ ] Status 200, `success=true`
- [ ] `data` es array (puede estar vacio)
- [ ] Cada item incluye `cod_total > 0`, `is_paid` (bool), `payment_method_code`, `order_total_amount`, `order_currency`
- [ ] `total`, `page`, `page_size`, `total_pages` presentes

## Caso 8.2: Filtro por estado

```
GET /api/v1/shipments/cod?status=delivered&page_size=20
```

- [ ] Todos los items tienen `status="delivered"`

## Caso 8.3: Filtro por is_paid

```
GET /api/v1/shipments/cod?is_paid=false
GET /api/v1/shipments/cod?is_paid=true
```

- [ ] Para `is_paid=false`: ningun item con `is_paid=true`
- [ ] Para `is_paid=true`: todos los items con `is_paid=true`

## Caso 8.4: Paginacion

```
GET /api/v1/shipments/cod?page=1&page_size=2
GET /api/v1/shipments/cod?page=2&page_size=2
```

- [ ] Diferente set de IDs entre pagina 1 y pagina 2
- [ ] `page_size=2` respetado
- [ ] `total_pages = ceil(total/2)`

## Caso 8.5: Super admin debe pasar business_id

```
POST /api/v1/auth/login (super admin) -> SUPER_TOKEN
GET  /api/v1/shipments/cod                              -> 400 (super admin: business_id requerido)
GET  /api/v1/shipments/cod?business_id=26               -> 200
```

- [ ] Sin `business_id`: status 400 con mensaje de business_id requerido
- [ ] Con `business_id=26`: status 200, items pertenecen a business 26
