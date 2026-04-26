# Datos de prueba compartidos

## Variables de entorno

| Variable | Valor | Descripcion |
|----------|-------|-------------|
| BASE_URL | http://localhost:3050/api/v1 | URL base de la API |
| SUPER_ADMIN_TOKEN | (obtener via login) | Token JWT de super admin |
| BUSINESS_USER_TOKEN | (obtener via login) | Token JWT de usuario con business_id > 0 |

## IDs generados durante las pruebas

Estos IDs se generan en CU-02 y se usan en tests subsiguientes:

| Variable | Generado en | Descripcion |
|----------|-------------|-------------|
| ANNOUNCEMENT_ID_1 | CU-02 Caso 2.1 | Anuncio global modal_text |
| ANNOUNCEMENT_ID_2 | CU-02 Caso 2.2 | Anuncio con links |
| ANNOUNCEMENT_ID_3 | CU-02 Caso 2.3 | Anuncio scheduled |
| ANNOUNCEMENT_ID_4 | CU-02 Caso 2.4 | Anuncio segmentado con targets |
| LINK_ID_1 | CU-02 Caso 2.2 | Primer link del ANNOUNCEMENT_ID_2 |

## Orden de ejecucion

Las pruebas deben ejecutarse en orden numerico:

1. CU-01: Categorias (verificar seed)
2. CU-02: Crear anuncios (genera IDs)
3. CU-03: Listar y obtener (usa IDs de CU-02)
4. CU-04: Actualizar (modifica ANNOUNCEMENT_ID_1)
5. CU-05: Cambiar estado (modifica ANNOUNCEMENT_ID_1)
6. CU-06: Views y stats (registra interacciones)
7. CU-07: Anuncios activos (verifica filtros de frecuencia)
8. CU-08: Force redisplay (limpia views)
9. CU-09: Delete (cleanup final)

## Categorias seedeadas (referencia)

| ID | Code | Name |
|----|------|------|
| 1 | promotion | Promocion |
| 2 | alert | Alerta |
| 3 | informative | Informativo |
| 4 | tutorial | Tutorial |
| 5 | update | Actualizacion |
| 6 | terms | Terminos y Condiciones |
