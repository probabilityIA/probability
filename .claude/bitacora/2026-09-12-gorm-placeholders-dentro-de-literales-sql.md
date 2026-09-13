# Un `?` dentro de un literal SQL rompe el bindeo de parametros en GORM

**Ticket:** sin ticket, pendiente de crear
**Negocio:** 26 (Demo), local
**Modulo:** notificaciones, audiencia de campanas

## Resumen

Al normalizar nombres de ciudad con una expresion regular, los `?` del patron
quedaron dentro de un literal de texto del SQL. GORM los conto como placeholders
al armar los `$N` de Postgres y rompio las tres consultas del querier de
audiencia.

## Sintoma

Todo endpoint de audiencia devolvia el mismo error:

```json
{"error":"ERROR: could not determine data type of parameter $1 (SQLSTATE 42P18)","success":false}
```

Afectaba a `GET /whatsapp-campaigns/audience-locations` y tambien a
`POST /whatsapp-campaigns/audience-preview`, incluso con `audience_type:
all_clients`, que no manda ningun filtro.

## Diagnostico

**Hipotesis descartada (la que costo tiempo):** que Postgres no pudiera inferir
el tipo del parametro. Es lo que el mensaje sugiere, y habia un candidato
plausible: el parametro de ciudad iba envuelto en `trim(unaccent(lower(...)))`,
funciones sobrecargadas donde la inferencia si puede fallar.

Se probo en psql con la consulta exacta:

```sql
PREPARE t AS <consulta completa con $1 sin cast>;
EXECUTE t(26);
```

**Funciono perfecto**, devolviendo los 43 municipios con sus conteos. Y
`PREPARE p1 AS SELECT regexp_replace(lower(unaccent(trim($1))), 'x', '', 'g');`
tambien preparo sin quejarse. O sea: **psql no reproduce el fallo**. El problema
no estaba en el SQL sino en lo que GORM enviaba.

## Causa raiz

La normalizacion de ciudad usaba:

```
regexp_replace(lower(unaccent(trim(x))), '\s*[,(]?\s*d\.?\s*c\.?\s*\)?\s*$', '', 'g')
```

Ese literal tiene **cuatro `?`**. GORM arma los placeholders numerados barriendo
el SQL, y los `?` de dentro del literal entraron en la cuenta, corriendo la
numeracion de los parametros reales.

El SQL anterior nunca habia tenido un `?` dentro de un literal, por eso el
problema aparecio recien con este cambio.

## Correccion

Se elimino la expresion regular. El unico sufijo que habia que limpiar es el de
"BOGOTA, D.C." (verificado: es el unico nombre del catalogo de geozonas con esa
forma), asi que alcanza con `replace()` encadenado, que no lleva `?` ni `$`:

```
btrim(replace(replace(lower(unaccent(trim(x))), ', d.c.', ''), ' d.c.', ''))
```

Commit `432cfd1a`.

## Verificacion

Contra la base local, negocio 26, por endpoint y no solo por psql:

- `audience-locations` devuelve 43 municipios con su departamento y conteo.
- `audience-preview` con `all_clients`: 610 total, 403 alcanzables.
- filtrando por departamento Choco: 40. Por ciudad "Bogota, D.C.": 20
  alcanzables, que coincide con el conteo del desplegable.

## Leccion

**psql no sirve para descartar este tipo de fallo.** Si una consulta corre bien
a mano y falla desde la aplicacion, sospechar del armado de placeholders antes
que del SQL. La prueba util es llamar al endpoint real.

## Pendientes

- Revisar si hay otro SQL crudo en el repo con `?` dentro de literales. El caso
  tipico es una expresion regular con cuantificadores opcionales.
- Crear el ticket correspondiente.
