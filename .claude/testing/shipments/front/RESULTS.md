# Resultados E2E - Shipments / EnvioClick Frontend

**Ultima ejecucion:** pendiente
**Resumen:** - PASS / - FAIL de 12 casos

---

## Casos de prueba frontend

Ejecutar con MCP Playwright contra http://localhost:3000
Login con super admin: ${AI_SUPER_ADMIN_EMAIL} / ${AI_SUPER_ADMIN_PASSWORD} (ver .env.ai)
Seleccionar negocio Demo (business_id=26)

### TC-F01: Navegar a seccion envios
- [ ] Login exitoso
- [ ] Seleccionar negocio Demo
- [ ] Navegar a /shipments
- [ ] Lista de envios carga correctamente

### TC-F02: Abrir modal de generacion de guia (orden normal)
- [ ] Seleccionar una orden SIN cod_total desde /orders
- [ ] Click en "Generar guia"
- [ ] Modal se abre con 4 pasos
- [ ] Campo "Valor contra entrega (COD)" visible y en 0
- [ ] NO se muestra banner "Orden Contra Entrega"

### TC-F03: Abrir modal de generacion de guia (orden COD)
- [ ] Preparar: UPDATE orders SET cod_total = 150000 WHERE id = '{ORDER_UUID}'
- [ ] Seleccionar la orden COD desde /orders
- [ ] Click en "Generar guia"
- [ ] Banner "Orden Contra Entrega - $150,000 COP" visible
- [ ] Campo codValue pre-llenado con 150000
- [ ] Campo codValue es read-only

### TC-F04: Cotizar envio (Step 1 -> Step 2)
- [ ] Llenar datos de origen (DANE Bogota)
- [ ] Llenar datos de destino (DANE Medellin)
- [ ] Dimensiones y peso correctos
- [ ] Click "Cotizar"
- [ ] Spinner de carga aparece
- [ ] Rates se muestran en grid de 4 columnas

### TC-F05: Filtrado COD en cotizaciones
- [ ] Con orden COD (codValue > 0): solo se muestran rates con badge "COD"
- [ ] Badge "Contra Entrega - Solo opciones COD" visible
- [ ] Sin orden COD (codValue = 0): se muestran TODAS las rates
- [ ] Rates sin COD muestran su precio normalmente (sin badge COD)

### TC-F06: Rates sin COD disponible
- [ ] Si ninguna rate tiene cod=true y codValue > 0
- [ ] Mensaje "No hay transportadoras disponibles con opcion contra entrega" visible
- [ ] No se pueden seleccionar rates

### TC-F07: Seleccionar rate y continuar a Step 3
- [ ] Click en una rate card
- [ ] Avanza a Step 3 (detalles de contacto)
- [ ] Campos de origen y destino pre-llenados desde la orden

### TC-F08: Generar guia (Step 4)
- [ ] Completar Step 3 con datos validos
- [ ] Avanzar a Step 4 (pago)
- [ ] Balance de wallet visible
- [ ] Click "Confirmar y generar"
- [ ] Spinner de generacion
- [ ] Tracking number y PDF URL aparecen al completar

### TC-F09: Cancelar envio desde lista
- [ ] Ir a /shipments
- [ ] Seleccionar un envio con status "pending"
- [ ] Click "Cancelar"
- [ ] Modal de confirmacion aparece
- [ ] Confirmar cancelacion
- [ ] Envio cambia a status "cancelled" (via SSE, sin reload)

### TC-F10: Error de cancelacion via SSE
- [ ] Intentar cancelar un envio ya cancelado
- [ ] Alert muestra "Cancelacion fallida: ..."

### TC-F11: Generacion masiva - ordenes COD marcadas
- [ ] Ir a /shipments
- [ ] Abrir modal de generacion masiva
- [ ] Ordenes con cod_total > 0 muestran badge "COD $X"
- [ ] En tabla de confirmacion, columna Orden muestra badge "COD"

### TC-F12: Generacion masiva - COD en payload
- [ ] Seleccionar orden COD para generacion masiva
- [ ] Verificar en network/logs que el quote payload incluye codValue y codPaymentMethod
- [ ] Verificar que el generate payload incluye codValue

---

## Bugs encontrados y corregidos durante testing

(Se documentaran durante la ejecucion)

## Archivos modificados durante testing (bugfixes)

| Archivo | Cambio |
|---------|--------|
| (pendiente) | (pendiente) |

## CU-18 Playwright UI - 2026-04-25

PASS: 9/9

- [OK] Login demo via UI funciona
- [OK] Tab "Contra entrega" visible en orders-subnavbar
- [OK] Click navega a /shipments/cod, tab queda activo (purple)
- [OK] Lista muestra 7 cards con badges Recaudado/Por cobrar
- [OK] Click en card pendiente: detalle correcto, boton "Marcar recaudado" disabled
- [OK] Click en card delivered no pagado (prob-0021): detalle con mapas + InfoCards + boton enabled
- [OK] Click "Marcar recaudado" abre modal con monto readonly y textarea notas
- [OK] Confirmar cobro: toast "Cobro registrado exitosamente" + badge cambia a "Recaudado"
- [OK] Verificacion DB via API: prob-0021.is_paid=true, paid_at=2026-04-25T21:19:22

Screenshot: cu18-after-collect.png
