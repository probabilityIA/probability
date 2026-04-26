# Testing E2E

## Definicion

"Testing", "pruebas", "testiemos", "lanza prueba back/front X", "prueba CU-XX de Y" = ejecutar casos de uso E2E en `.claude/testing/<modulo>/back/` o `front/`. NO son tests unitarios Go/Jest.

## Estructura

```
.claude/testing/<modulo>/
  back/     CU-NN-kebab-descripcion.md, RESULTS.md
  front/    CU-NN-kebab-descripcion.md, RESULTS.md
  shared/   test_data.md
```

## Reglas BACK

1. Crear/modificar SOLO via API. `mcp__postgres-probability__query` = SOLO lectura (SELECT para verificar efectos y construir bodies).
2. Ejecutar endpoints en el orden del caso. Encadenar outputs entre pasos.
3. HTTP con `mcp__fetch__fetch`. Fallback: `curl`. Nunca inventar respuestas.
4. Analizar status code, campo `success`, forma del payload. Desviaciones = bug a registrar.
5. JWT via `POST /api/v1/auth/login` con credenciales de `test-credentials.md`. Base URL: `http://localhost:3050/api/v1`.
6. Verificar que backend corre antes: `curl -I http://localhost:3050/health`. Si no responde, pedir al usuario.
7. Registrar resultado en `RESULTS.md`: fecha, caso, OK/FAIL, bug, commit del fix.

## Reglas FRONT

1. Playwright MCP (`mcp__playwright__browser_*`) para navegacion/clicks/forms. Chrome DevTools para network/console.
2. No delegar a agentes; usar los MCP directamente.
3. Base URL: `http://localhost:3000`. Validar texto visible, XHR esperado, consola sin errores.
4. Mutaciones: verificar con MCP postgres que quedaron reflejadas.

## Permisos durante testing

Permitido sin preguntar: modificar codigo Go/TS, compilar, commit + push del fix.
Prohibido: INSERT/UPDATE/DELETE en DB, reiniciar servicios sin avisar, desactivar validaciones, editar el caso para que coincida con un bug.

## Logs para debug

- Backend: `./scripts/dev-services.sh logs backend 200`
- Debug file: `/back/central/log/app-YYYY-MM-DD.log` (si `ENABLE_DEBUG_FILE_LOGGING=true`)
- Frontend consola/red: `mcp__playwright__browser_console_messages` / `browser_network_requests`
