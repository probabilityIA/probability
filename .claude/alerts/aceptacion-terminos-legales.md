# ALERTA: Aceptacion de terminos legales - pasos pendientes para produccion

Creada: 2026-08-27
Contexto: se implemento el gate de aceptacion de Terminos y Condiciones v1.0 y
Politica de Tratamiento de Datos Personales v1.0 al iniciar sesion.

## Urgente (sin esto la funcionalidad no existe en produccion)

1. **Correr la migracion en produccion.** `migrateLegalDocuments` crea
   `legal_documents` + `legal_acceptances` y siembra los dos documentos v1.0.
   Ya se corrio en local (2026-08-27). En produccion el deploy NO corre
   migraciones: hay que llamarla desde `Migrate()` y ejecutar el binario contra
   el RDS, o aplicar el DDL a mano por el tunel.
2. **Verificar que los PDFs quedaron desplegados** en
   `front/central/public/legal/`. El registro de aceptacion guarda el SHA-256 de
   cada archivo:
   - terminos-y-condiciones-v1.0.pdf: `25639ec867d7261189d61f262a80bdbedd36628024c0d8959d8cfbb06730c5ab`
   - politica-datos-personales-v1.0.pdf: `f97b9ff947034a6b8b234bda5a039285b833809ec4f9ec0ccb7ec0e8b05eb7ba`
   Si un PDF se reemplaza sin cambiar de version, las aceptaciones viejas quedan
   apuntando a un hash que ya no corresponde al archivo servido.

## Importante

3. **El bloqueo es de UI, no de API.** El modal impide usar el panel, pero un
   usuario con token valido podria seguir llamando la API sin aceptar. Si se
   quiere blindaje real, agregar un middleware que rechace las rutas de negocio
   con 403 `legal_acceptance_required` mientras haya documentos pendientes.
4. **Prueba E2E pendiente**: login de un usuario de negocio, ver el modal,
   aceptar, y confirmar la fila en `legal_acceptances` (user_id, version, sha256,
   ip, user_agent). Solo hay tests unitarios del caso de uso.

## Como publicar una version nueva de un documento

1. Dejar el PDF nuevo en `front/central/public/legal/<nombre>-vX.Y.pdf`.
2. `sha256sum` del archivo.
3. Insertar la fila en `legal_documents` con la version nueva `is_active = TRUE`
   y poner `is_active = FALSE` en la anterior. **Nunca editar la fila vieja ni
   borrar aceptaciones**: son la evidencia de que el usuario acepto ESA version.
4. El gate vuelve a aparecer solo, porque el usuario no tiene aceptacion para el
   `legal_document_id` nuevo.

## Como cerrar esta alerta

Cuando los puntos 1 y 2 esten hechos en produccion y el punto 4 verificado.
