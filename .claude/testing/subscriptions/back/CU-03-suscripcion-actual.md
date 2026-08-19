# CU-03: Consultar suscripcion actual del negocio

## Objetivo
Verificar `GET /me` para el negocio del token (business normal) y para super admin
(que debe requerir `business_id` explicito).

## Precondiciones
- DEMO_TOKEN y SUPER_TOKEN obtenidos (CU-01).

## Pasos

### 3.1 Business normal - usa business_id del token
```
GET /api/v1/subscriptions/me
Authorization: Bearer {DEMO_TOKEN}
```
- [ ] Status 200
- [ ] Si Demo tiene suscripcion activa: `data.business_id == 26`
- [ ] Si no tiene: `data: null` (no error)

### 3.2 Super admin SIN business_id -> debe fallar
```
GET /api/v1/subscriptions/me
Authorization: Bearer {SUPER_TOKEN}
```
- [ ] Status 400 (segun patron resolveBusinessID de multi-tenant-security.md)

### 3.3 Super admin CON business_id de Demo
```
GET /api/v1/subscriptions/me?business_id=26
Authorization: Bearer {SUPER_TOKEN}
```
- [ ] Status 200
- [ ] Mismo resultado que 3.1

### 3.4 Business normal intenta ver otro negocio via query (debe ignorarse)
```
GET /api/v1/subscriptions/me?business_id=999
Authorization: Bearer {DEMO_TOKEN}
```
- [ ] Status 200
- [ ] `data.business_id == 26` (el query param se IGNORA, NO 999) -> si devuelve datos de otro
      negocio o error revelando que 999 existe, es una violacion de multi-tenant-security.md
