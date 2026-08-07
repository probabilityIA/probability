# Bitacora

Historico de soportes, incidentes y diagnosticos del proyecto. Un archivo por
caso, con el diagnostico completo: que se rompio, como se encontro, que se
descarto en el camino, que se corrigio y que quedo pendiente.

## Para que sirve

Cuando vuelve a aparecer un problema parecido, o cuando alguien (persona o IA)
necesita entender por que el codigo hace algo raro, aca esta el contexto que no
cabe en un commit. Los commits dicen QUE cambio; la bitacora dice POR QUE y con
que evidencia.

## Cuando escribir una entrada

- Se investigo un problema de produccion que costo tiempo o plata.
- Se descubrio un comportamiento no documentado de un proveedor externo.
- Se corrigio data en produccion.
- Una hipotesis razonable resulto falsa (dejarla escrita evita que el siguiente
  la repita).

No hace falta entrada para un fix trivial o un cambio de UI.

## Nomenclatura

`YYYY-MM-DD-tema-corto.md`

## Estructura sugerida

1. Resumen en dos lineas
2. Sintoma (con numeros reales)
3. Diagnostico: la cadena de evidencia, incluidas las hipotesis descartadas
4. Causa raiz
5. Correccion (codigo y/o data)
6. Verificacion
7. Pendientes

## Indice

| Fecha | Tema | Estado |
|-------|------|--------|
| 2026-08-06 | [COD EnvioClick: valor a recaudar mal calibrado](2026-08-06-cod-envioclick-calibracion.md) | Corregido, pendiente deploy |
