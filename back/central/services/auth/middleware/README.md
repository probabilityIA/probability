# Middleware de autenticacion

Este paquete autentica: valida el token (cookie `session_token` o header
`Authorization`) y deja en el contexto de gin `user_id`, `business_id`,
`role_id` e `is_super_admin`.

**No autoriza.** Quien puede usar cada ruta lo decide el modulo
`services/auth/authz` con la tabla `shared/authz/route_policies.go`. Ver
`.claude/rules/autorizacion.md`.

## Middlewares

| Funcion | Uso |
|---|---|
| `JWT()` | exige sesion valida en la ruta |
| `Authenticate(c)` | valida la sesion sin abortar; lo usa el middleware global de autorizacion |
| `RequireSuperAdmin()` | solo super admin (`business_id = 0` en el token) |
| `RequireJWT()` / `RequireAPIKey()` | exige un tipo de autenticacion |

`RequireRole` y `RequireAnyRole` se eliminaron el 2026-09-14: leian
`user_roles`, que el camino JWT nunca escribia, y ninguna ruta los usaba.

## Helpers de contexto

```go
userID, ok := middleware.GetUserID(c)
businessID, ok := middleware.GetBusinessIDFromContext(c)
isSuper := middleware.IsSuperAdmin(c)
```

Para el `business_id` efectivo seguir `.claude/rules/multi-tenant-security.md`:
el del token manda para un usuario normal y solo el super admin puede aportar
uno por query o body.

## Respuestas de error

| Codigo | Cuando |
|---|---|
| 401 | token ausente, invalido o vencido |
| 402 | suscripcion vencida o cancelada (middleware de autorizacion) |
| 403 | la ruta exige un permiso o super admin que el usuario no tiene |
