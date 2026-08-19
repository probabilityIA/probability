# CU-07: Overrides de acceso a modulos por negocio

## Objetivo
Verificar grant/list/revoke de overrides — dan acceso a un modulo a un negocio sin
necesidad de que su plan lo incluya. Exclusivo de super admin por ser sensible
(otorga funcionalidad de pago gratis).

## Precondiciones
- SUPER_TOKEN y DEMO_TOKEN obtenidos (CU-01).
- Elegir un `module_code` que Demo NO tenga activo por su plan actual
  (comparar contra CU-02.5 `my-modules`), ej `whatsapp` o similar.

## 7.1 Otorgar override
```
POST /api/v1/subscriptions/overrides
Authorization: Bearer {SUPER_TOKEN}
Content-Type: application/json

{
  "business_id": 26,
  "module_code": "{MODULE_CODE}",
  "notes": "Override de prueba E2E"
}
```
- [ ] Status 201, mensaje de confirmacion
- [ ] Verificar en DB que la fila quedo con `granted_by_user_id` = el ID del super admin
      (confirmar via `middleware.GetUserID` -> revisar tabla real con SELECT)

## 7.2 Confirmar que my-modules de Demo ahora incluye el modulo
```
GET /api/v1/subscriptions/my-modules
Authorization: Bearer {DEMO_TOKEN}
```
- [ ] `{MODULE_CODE}` aparece en la lista aunque el plan de Demo no lo incluya

## 7.3 Listar overrides del negocio
```
GET /api/v1/subscriptions/overrides/26
Authorization: Bearer {SUPER_TOKEN}
```
- [ ] Status 200, el override de 7.1 aparece

## 7.4 Business normal no puede otorgar/listar/revocar overrides
Repetir 7.1 y 7.3 con `Authorization: Bearer {DEMO_TOKEN}` (y probar tambien
`GET /overrides/26` con el propio token de Demo, aunque sea su propio negocio).
- [ ] Status 403 en los tres casos -> si Demo puede LISTAR sus propios overrides
      con su propio token, no es necesariamente un problema de seguridad grave
      (es su propio negocio) pero violaria el diseno actual "todo por super admin";
      registrar como hallazgo de severidad baja/media si ocurre, no critico

## 7.5 Revocar override
```
DELETE /api/v1/subscriptions/overrides/26/{MODULE_CODE}
Authorization: Bearer {SUPER_TOKEN}
```
- [ ] Status 200
- [ ] Repetir 7.2: `{MODULE_CODE}` ya NO aparece en `my-modules` de Demo
