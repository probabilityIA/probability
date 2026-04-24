# Testing E2E - Módulo Inventario

Casos de prueba que validan el flujo completo del módulo WMS afectando múltiples submódulos.

## Orden sugerido

1. `front/CU-01-adjust-with-full-hierarchy.md` — Crea bodega con jerarquía completa + productos + ajustes por ubicación y lote

## Submódulos que se cruzan

- **Bodegas**: jerarquía física (zonas, pasillos, racks, niveles, posiciones)
- **Productos**: con `track_inventory=true`
- **Trazabilidad**: lotes con fecha de vencimiento
- **Stock**: ajustes que afectan niveles específicos
- **Movimientos**: kardex histórico

## Credenciales

Ver `/home/cam/Desktop/probability/.claude/rules/test-credentials.md`.

Para estos casos se usa el usuario `demo@probability.com` (business_id=26).
