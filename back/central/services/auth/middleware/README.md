# 🔐 Middleware de Autenticación JWT

Este directorio contiene middlewares para la autenticación y autorización de la API.

## 📋 Middlewares Disponibles

### 1. AuthMiddleware
Middleware principal para validar tokens JWT y extraer información del usuario.

```go
// Aplicar a una ruta específica
router.GET("/protected", middleware.AuthMiddleware(jwtService, logger), handler)

// Aplicar a un grupo de rutas
protectedGroup := router.Group("/api/v1")
protectedGroup.Use(middleware.AuthMiddleware(jwtService, logger))
{
    protectedGroup.GET("/users", handler)
    protectedGroup.POST("/data", handler)
}
```

### 2. RequireRole
Middleware que requiere un rol específico.

```go
// Requerir rol "admin"
router.GET("/admin-only", 
    middleware.AuthMiddleware(jwtService, logger),
    middleware.RequireRole("admin"),
    handler,
)
```

### 3. RequireAnyRole
Middleware que requiere al menos uno de los roles especificados.

```go
// Requerir rol "admin" o "manager"
router.GET("/management", 
    middleware.AuthMiddleware(jwtService, logger),
    middleware.RequireAnyRole("admin", "manager"),
    handler,
)
```

## 🔧 Funciones de Utilidad

### Obtener Información del Usuario

```go
func MyHandler(c *gin.Context) {
    // Obtener ID del usuario
    userID, exists := middleware.GetUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
        return
    }

    // Obtener email del usuario
    email, exists := middleware.GetUserEmail(c)
    
    // Obtener roles del usuario
    roles, exists := middleware.GetUserRoles(c)
    
    // Obtener claims completos del JWT
    claims, exists := middleware.GetJWTClaims(c)
}
```

## 📝 Ejemplo Completo

```go
package main

import (
    "central_reserve/internal/infra/primary/http2/middleware"
    "central_reserve/internal/pkg/jwt"
    "central_reserve/internal/pkg/log"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine, jwtService *jwt.JWTService, logger log.ILogger) {
    // Rutas públicas
    router.POST("/login", loginHandler)
    
    // Rutas protegidas
    protected := router.Group("/api/v1")
    protected.Use(middleware.AuthMiddleware(jwtService, logger))
    {
        // Rutas para cualquier usuario autenticado
        protected.GET("/profile", getProfileHandler)
        
        // Rutas que requieren rol específico
        admin := protected.Group("/admin")
        admin.Use(middleware.RequireRole("admin"))
        {
            admin.GET("/users", getAllUsersHandler)
            admin.POST("/users", createhandlers)
        }
        
        // Rutas que requieren uno de varios roles
        management := protected.Group("/management")
        management.Use(middleware.RequireAnyRole("admin", "manager"))
        {
            management.GET("/reports", getReportsHandler)
        }
    }
}

func getProfileHandler(c *gin.Context) {
    userID, exists := middleware.GetUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
        return
    }
    
    // Usar userID para obtener datos del perfil
    // ...
}
```

## 🛡️ Seguridad

- **Validación de Token**: El middleware valida automáticamente la firma y expiración del token
- **Extracción de Claims**: Los claims del JWT se extraen y almacenan en el contexto
- **Logging**: Se registran eventos de autenticación para auditoría
- **Manejo de Errores**: Respuestas HTTP apropiadas para diferentes tipos de errores

## 🔄 Flujo de Autenticación

1. **Cliente envía request** con header `Authorization: Bearer <token>`
2. **AuthMiddleware intercepta** la request
3. **Valida el token** usando el JWT service
4. **Extrae claims** y los almacena en el contexto
5. **Handler procesa** la request con información del usuario disponible
6. **Funciones de utilidad** permiten acceder a la información del usuario

## 📊 Respuestas de Error

| Código | Descripción |
|--------|-------------|
| 401 | Token requerido o inválido |
| 403 | Acceso denegado (rol insuficiente) |
| 500 | Error interno del servidor | 