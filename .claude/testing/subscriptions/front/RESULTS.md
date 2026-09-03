# Resultados - Modulo Subscriptions (front)

| Fecha | Caso | Resultado | Bug | Commit fix |
|-------|------|-----------|-----|------------|
| 2026-09-02 | CU-14 boton Pagar/Extender | OK | 14.1-14.4 pasan: login, modal abre para plan de catalogo (Basico) y para plan personalizado (regresion del fix de hoy), detecta saldo insuficiente/suficiente correctamente, compra completa via API real y actualiza la UI. Nota de entorno: el dev server (Turbopack) tuvo una compilacion stale que hizo parecer que el modal no abria para el plan personalizado -- un restart limpio del frontend confirmo que si funciona, no era un bug real | - |
| 2026-09-02 | CU-15 SubscriptionGuard | OK | 15.1 cuenta deshabilitada bloquea /orders con el muro "Cuenta suspendida". 15.2 /subscription sigue accesible bajo bloqueo. 15.3 reactivar por API sin relogin NO desbloquea (JWT stale, comportamiento documentado y esperado). 15.4 con sesion realmente limpia (equivalente a `TokenStorage.clearSession()`, que si limpia `permissions` de localStorage) el relogin desbloquea correctamente. Nota de metodologia: usar `page.context().clearCookies()` solo (sin limpiar localStorage) da un falso bloqueo persistente porque `permissions-context.tsx` cachea `subscription_status` en localStorage y no se limpia solo con borrar cookies -- esto no es un bug de la app (el logout real si limpia localStorage via `clearSession()`), fue un artefacto de la prueba | - |

## Notas de entorno para la proxima vez

- El dev server de frontend (Next.js 16 + Turbopack) se puso inestable durante pruebas repetidas con Playwright: paginas que tardaban 60s+ en compilar, procesos que se mataban solos, requiere `./scripts/dev-services.sh kill-zombies` + restart limpio de cuando en cuando.
- `./scripts/dev-services.sh` tiene un bug menor: usa `grep -P` en varios puntos, que no existe en el grep BSD de macOS (imprime "usage: grep..." pero no bloquea la ejecucion). No se corrigio por estar fuera de alcance de este testing.
- Para simular "cerrar sesion y volver a entrar" en pruebas con Playwright, usar `page.context().clearCookies()` Y `page.evaluate(() => localStorage.clear())` juntos, no solo cookies -- de lo contrario el permissions cache de localStorage da un falso positivo de bloqueo persistente.
