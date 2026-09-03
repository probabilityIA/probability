# Bitacora - Historico de soportes y correcciones

La carpeta `.claude/bitacora/` guarda el historico de soportes, incidentes,
correcciones y arreglos de bugs. **Un archivo por caso**, nunca un archivo
acumulativo que crezca sin fin.

Indice y estructura de una entrada: `.claude/bitacora/README.md`.

Toda entrada de bitacora referencia el ticket del que salio, y el comentario de
cierre del ticket referencia la entrada. Reglas: `.claude/rules/tickets.md`.

## Diferencia con `.claude/alerts/`

| | alerts | bitacora |
|---|---|---|
| Que es | pendiente urgente, trabajo sin terminar | historico de algo ya diagnosticado |
| Cuando se borra | al resolver los items urgentes | nunca |
| Se lee | al iniciar sesion, siempre | cuando el problema se parece |

Un caso puede empezar como alerta y terminar como entrada de bitacora.

## Reglas

1. **Antes de investigar un problema**: revisar `.claude/bitacora/` por si el
   caso o uno parecido ya esta diagnosticado. Sobre todo antes de tocar
   integraciones con proveedores externos, donde el comportamiento no
   documentado ya costo tiempo una vez.
2. **Crear entrada** cuando: se investigo un problema de produccion que costo
   tiempo o plata, se descubrio un comportamiento no documentado de un proveedor,
   se corrigio data en produccion, o una hipotesis razonable resulto falsa.
3. **Dejar escritas las hipotesis descartadas.** Es la parte mas util: evita que
   el siguiente repita el mismo callejon sin salida.
4. **Numeros reales, no descripciones vagas.** Ids, montos, trackings, consultas
   SQL que sirvan para reproducir el diagnostico.
5. **Nomenclatura**: `YYYY-MM-DD-tema-corto.md`. Actualizar el indice del README.
6. No hace falta entrada para un fix trivial o un cambio de UI.
