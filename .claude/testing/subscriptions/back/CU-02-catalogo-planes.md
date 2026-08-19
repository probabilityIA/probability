# CU-02: Catalogo de planes y modulos (lectura)

## Objetivo
Verificar los endpoints de solo lectura del catalogo, accesibles a cualquier usuario autenticado.

## Precondiciones
- DEMO_TOKEN obtenido (CU-01).

## Pasos

### 2.1 Listar tipos de suscripcion
```
GET /api/v1/subscriptions/types
Authorization: Bearer {DEMO_TOKEN}
```
- [ ] Status 200, `success=true`
- [ ] `data` es un array de planes (esperado <50 registros -> catalogo, sin paginacion es aceptable segun architecture.md)
- [ ] Guardar un `SUBSCRIPTION_TYPE_ID` activo con precio > 0 para CU-03

### 2.2 Obtener un tipo de suscripcion por ID
```
GET /api/v1/subscriptions/types/{SUBSCRIPTION_TYPE_ID}
Authorization: Bearer {DEMO_TOKEN}
```
- [ ] Status 200, `data.id == SUBSCRIPTION_TYPE_ID`

### 2.3 Codigos de modulo
```
GET /api/v1/subscriptions/module-codes
Authorization: Bearer {DEMO_TOKEN}
```
- [ ] Status 200, `data` es un array de codigos/strings

### 2.4 Catalogo completo de modulos
```
GET /api/v1/subscriptions/module-catalog
Authorization: Bearer {DEMO_TOKEN}
```
- [ ] Status 200, `data` no vacio

### 2.5 Mis modulos (segun suscripcion activa del negocio del token)
```
GET /api/v1/subscriptions/my-modules
Authorization: Bearer {DEMO_TOKEN}
```
- [ ] Status 200
- [ ] Los modulos devueltos son coherentes con la suscripcion activa de Demo (comparar con CU-04)
