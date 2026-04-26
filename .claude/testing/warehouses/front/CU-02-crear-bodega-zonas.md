# CU-02 — Crear bodega tipo Con Zonas

## Objetivo
Crear una bodega con estructura "Con Zonas" y datos de contacto.

## Precondiciones
- Usuario autenticado
- Frontend corriendo en http://localhost:3000

## Pasos

1. Navegar a /warehouses
2. Clic en "Nueva bodega"
3. Ingresar Nombre: "Bodega Bogota Zonas"
4. Seleccionar card "Con Zonas"
5. Ingresar Direccion: "Calle 80 #68-50"
6. Ingresar Ciudad: "Bogota", Departamento: "Cundinamarca", Pais: "Colombia"
7. En seccion Contacto: Telefono "3201234567", Nombre "Carlos Lopez"
8. Dejar "Bodega principal" en OFF
9. Clic en "Crear bodega"
10. Verificar mensaje de exito
11. Verificar que aparece en lista

## Resultado esperado
- Bodega creada correctamente
- Estructura "zones" guardada en localStorage con key wh_struct_{id}
- Datos de contacto persistidos
