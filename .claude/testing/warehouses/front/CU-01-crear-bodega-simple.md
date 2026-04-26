# CU-01 — Crear bodega tipo Simple

## Objetivo
Crear una bodega con estructura Simple (sin zonas ni ubicaciones), con direccion autocompletada y mapa visible.

## Precondiciones
- Usuario autenticado como demo@probability.com
- Frontend corriendo en http://localhost:3000

## Pasos

1. Navegar a /warehouses
2. Clic en "Nueva bodega"
3. Verificar que el modal se abre con el formulario en 2 columnas
4. Ingresar Nombre: "Bodega Medellin Test"
5. Verificar que la card "Simple" esta seleccionada por defecto (borde morado)
6. En campo Direccion digitar "Calle 10 Sur" y esperar sugerencias (min 8 chars)
7. Ingresar Ciudad: "Medellin" manualmente (si no autocompleta)
8. Ingresar Departamento: "Antioquia"
9. Ingresar Pais: "Colombia"
10. Activar el switch "Bodega principal"
11. Clic en "Crear bodega"
12. Verificar mensaje de exito
13. Verificar que la bodega aparece en la lista

## Resultado esperado
- Bodega creada exitosamente
- Aparece en lista con nombre "Bodega Medellin Test"
- Switch de bodega principal activo (se guarda is_default=true)
- Estructura guardada en localStorage como "simple"
