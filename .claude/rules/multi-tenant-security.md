# Seguridad Multi-Tenant - Aislamiento por Business

Regla de seguridad OBLIGATORIA para TODOS los endpoints (lectura y escritura).
Objetivo: que un usuario solo pueda ver y ejecutar acciones sobre SU propio
business, nunca sobre otro. Un business que conozca el ID de otro NO debe poder
usar sus endpoints para afectar datos ajenos.

## Modelo de usuarios

- **Usuario business normal:** su `business_id` viene SIEMPRE del token (JWT).
  - Prohibido tomar `business_id` de body o query. Si el cliente lo envia, se IGNORA.
  - Solo puede ver y operar sobre su propio business.
- **Super admin:** `business_id = 0` en el token. NO tiene business propio.
  - Ve todo, pero para CUALQUIER accion (incluso una simple consulta) debe
    seleccionar/enviar el `business_id` con el que quiere operar.
  - `business_id` de super admin va por query param (GET/DELETE) o body/query
    segun el endpoint. Sin el = 400.

## Como resolver el business (patron obligatorio)

```go
func resolveBusinessID(c *gin.Context, bodyBusinessID *uint) (uint, bool) {
    businessID, ok := middleware.GetBusinessIDFromContext(c) // viene del token
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
        return 0, false
    }
    if businessID == 0 { // super admin: no tiene business en el token
        if bodyBusinessID == nil || *bodyBusinessID == 0 {
            c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
            return 0, false
        }
        businessID = *bodyBusinessID
    }
    return businessID, true // normal: SIEMPRE el del token, se ignora el del body
}
```

Clave: para el usuario normal, `businessID` del token gana y el enviado por el
cliente NUNCA se usa. Solo el super admin (token = 0) puede aportar el valor.

## Validar propiedad del recurso (no confiar en el ID enviado)

Cuando la accion referencia un recurso con dueno (integracion, orden, producto,
webhook, etc.), verificar que ese recurso pertenece al `business_id` ya resuelto.
No basta con filtrar la consulta principal por business: si el request trae un
`integration_id`/`order_id` de otro business, hay que rechazarlo.

```go
integration, err := core.GetIntegrationByID(ctx, integrationID)
if err != nil || integration == nil {
    return fmt.Errorf("recurso no encontrado")
}
if integration.BusinessID == nil || *integration.BusinessID != businessID {
    return fmt.Errorf("el recurso no pertenece al negocio")
}
```

## Checklist por endpoint

1. Resolver `business_id` con el patron de arriba (token para normal, param para super admin).
2. Toda consulta filtra por `business_id`.
3. Todo recurso referenciado por ID en el request se valida que pertenezca a ese `business_id`.
4. En POST/PUT/DELETE, el `business_id` de super admin va por query param, no en el body de dominio.

Ver tambien: `backend-conventions.md` seccion 5 (Super Admin - Business ID) y las
notas de frontend (`isSuperAdmin` -> selector obligatorio, pasar `business_id` a
todas las operaciones, resetear al cambiar de negocio).
