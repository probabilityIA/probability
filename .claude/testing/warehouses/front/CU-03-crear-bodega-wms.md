# CU-03 — Crear bodega tipo WMS Completo

## Objetivo
Crear una bodega con estructura WMS Completo y verificar que el mapa aparece al seleccionar direccion con autocomplete.

## Precondiciones
- Usuario autenticado
- Frontend corriendo en http://localhost:3000

## Pasos

1. Navegar a /warehouses
2. Clic en "Nueva bodega"
3. Ingresar Nombre: "Centro Logistico WMS"
4. Seleccionar card "WMS Completo"
5. En campo Direccion digitar "Avenida El Dorado" y esperar sugerencias del autocomplete
6. Si aparecen sugerencias, seleccionar una — verificar que Ciudad/Departamento se autocompletan
7. Verificar que el mapa aparece debajo de los campos de direccion
8. Ingresar Contacto: Email "logistica@empresa.com"
9. Dejar "Bodega principal" en OFF
10. Clic en "Crear bodega"
11. Verificar mensaje de exito
12. Verificar que aparece en lista

## Resultado esperado
- Bodega creada correctamente
- Si se selecciono sugerencia: ciudad/dpto/codpostal autocompletados, mapa visible
- Estructura "wms" guardada en localStorage con key wh_struct_{id}
