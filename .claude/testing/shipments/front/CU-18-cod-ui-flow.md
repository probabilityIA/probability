# CU-18 (front): Flujo UI completo COD via Playwright

## Objetivo
Validar la pestana "Contra entrega" del navbar de Ordenes, el listado, el detalle con
mapas y el modal de "Marcar recaudado" usando Playwright MCP.

## Precondiciones
- Frontend en `http://localhost:3000`.
- Backend en `http://localhost:3050` con datos COD generados por CU-09..CU-13.
- Usuario demo: `demo@probability.com` / `ProbabilityDemo`.

## Pasos

1. `browser_navigate http://localhost:3000/login`
2. Llenar email y password (`browser_fill_form` o `browser_type`), submit.
3. Esperar redirect a `/home`. `browser_snapshot` para verificar layout.
4. `browser_navigate http://localhost:3000/orders`
   - [ ] Tab "Contra entrega" visible en el navbar (icono $).
5. Click en "Contra entrega".
   - [ ] URL cambia a `/shipments/cod`.
   - [ ] Tab "Contra entrega" queda activo (clase purple).
   - [ ] Aparecen cards en la columna izquierda (>=1 item).
6. Click en una card con badge "Por cobrar".
   - [ ] Panel derecho muestra detalle con: cliente, orden #, monto COD, transportadora, badges.
   - [ ] Aparecen 2 mini-mapas (Origen, Destino).
   - [ ] Aparecen InfoCards (Monto a cobrar, Total orden, Transportadora, Entregado).
   - [ ] Si esta delivered+no pagado, boton "Marcar recaudado" habilitado.
7. Click "Marcar recaudado".
   - [ ] Modal con monto readonly y textarea de notas aparece.
8. Escribir nota y confirmar.
   - [ ] Toast / mensaje "Cobro registrado exitosamente".
   - [ ] Badge cambia a "Recaudado" (verde).
   - [ ] Card en lista actualiza badge a "Recaudado".
9. `browser_console_messages`
   - [ ] Sin errores rojos relacionados con COD/shipments.

## Reporte
Capturar `browser_take_screenshot` antes/despues del cobro para validar visualmente.
