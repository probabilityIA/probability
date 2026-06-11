# Alertas - Pendientes Urgentes

La carpeta `.claude/alerts/` contiene pendientes importantes y urgentes del proyecto.

## Reglas

1. **Al iniciar cualquier sesion de trabajo**: revisar `ls .claude/alerts/` y leer
   las alertas existentes. Son pendientes de alta prioridad, no documentacion.
2. **Antes de tocar un modulo**: si existe una alerta relacionada con ese modulo,
   leerla completa y tenerla en cuenta (puede haber riesgos conocidos, p.ej.
   operaciones sin idempotencia o endpoints sin validar).
3. **Mencionar al usuario** las alertas relevantes al trabajo en curso, sobre todo
   si el cambio que pide toca un punto listado como urgente.
4. **Crear alerta**: cuando quede trabajo critico inconcluso (riesgo de datos,
   bloqueante de produccion, deuda peligrosa), crear `.claude/alerts/<tema>.md`
   con: fecha, contexto, items clasificados (urgente / importante / deseable) y
   criterio para cerrarla.
5. **Cerrar alerta**: solo eliminar el archivo cuando los items urgentes e
   importantes esten resueltos y verificados. Actualizar el archivo si se resuelve
   un item parcial (marcarlo como resuelto con fecha, no borrarlo de la lista).
