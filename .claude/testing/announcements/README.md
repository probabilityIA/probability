# E2E Tests - Modulo Announcements

## Estructura

```
announcements/
+-- back/           # Pruebas E2E de API backend
+-- front/          # Pruebas E2E de frontend (pendiente)
+-- shared/         # Utilidades compartidas
+-- README.md
```

## Casos de prueba backend

Todas las pruebas se ejecutan contra la API REST (`http://localhost:3050/api/v1`).
Se usa MCP postgres solo para verificar datos en base de datos.

### Requisitos

- Backend corriendo en puerto 3050
- PostgreSQL corriendo en puerto 5433
- Token JWT valido de super admin

### Ejecucion

Las pruebas son documentos de casos de prueba que se ejecutan manualmente
o con el agente debug-tester.
