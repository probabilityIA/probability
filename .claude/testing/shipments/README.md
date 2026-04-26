# E2E Tests - Modulo Shipments / EnvioClick

## Estructura

```
shipments/
  README.md
  shared/
    test_data.md
  back/
    CU-01-setup-test-mode.md
    CU-02-quote.md
    CU-03-generate-guide.md
    CU-04-track.md
    CU-05-cancel.md
    CU-06-webhook.md
    CU-07-cod-flow.md
    RESULTS.md
  front/
    RESULTS.md
```

## Casos de prueba backend

Pruebas E2E del modulo de generacion de guias de envio via EnvioClick.
Se usa el mock server (`back/testing`) en puerto 9091 para simular respuestas de EnvioClick.

### Requisitos

- Backend corriendo en `http://localhost:3050`
- Mock server corriendo en `http://localhost:9091` (`back/testing`)
- PostgreSQL en puerto 5433
- RabbitMQ en puerto 5672
- Redis en puerto 6379
- Integration ID 44 (Demo business) con `is_testing = true`

### Ejecucion

1. Iniciar infra: `./scripts/dev-services.sh start infra`
2. Iniciar mock server: `cd back/testing && go run cmd/main.go`
3. Iniciar backend: `./scripts/dev-services.sh start backend`
4. Activar test mode: ejecutar CU-01 (setup)
5. Ejecutar tests CU-02 a CU-07 en orden
6. Verificar resultados en `back/RESULTS.md`

### Casos de prueba frontend

Pruebas del modal de generacion de guias, filtrado COD, y cancelacion.
Se ejecutan con MCP Playwright contra `http://localhost:3000`.
