# Actualizar datos desde un canal fallaba con ON CONFLICT (SQLSTATE 21000)

Un producto de Probability puede estar mapeado a mas de una publicacion del mismo
canal. El motor de actualizacion de datos asumia una fila por producto, asi que
el upsert de origen tocaba la misma fila dos veces y Postgres abortaba todo.

## Sintoma

Al aplicar "Actualizar categoria en Probability" desde Mercadolibre para el
negocio 46 (Viga), la modal mostraba en rojo:

```
ERROR: ON CONFLICT DO UPDATE command cannot affect row a second time (SQLSTATE 21000)
```

No se aplicaba nada: la transaccion completa hacia rollback.

La modal decia "897 se llenan en Probability".

## Diagnostico

Conteo de mappings con snapshot por canal:

```sql
SELECT m.integration_id, COUNT(*) AS filas, COUNT(DISTINCT m.product_id) AS productos
FROM product_business_integrations m
WHERE m.business_id = 46 AND m.deleted_at IS NULL AND m.channel_snapshot IS NOT NULL
GROUP BY m.integration_id;
```

| integracion | filas | productos | duplicadas |
|---|---|---|---|
| 254 Mercadolibre | 897 | 860 | **37** |
| 221 WooCommerce | 1687 | 1687 | 0 |
| 228 Siigo | 2166 | 2166 | 0 |

37 productos estan publicados en dos listings distintos de Mercadolibre. Por eso
el problema solo aparecia con ese canal: Woo y Siigo tienen relacion 1 a 1.

Hipotesis descartadas:

- **Indice unico mal definido en `product_field_origins`.** El indice
  `(product_id, field)` es correcto; el problema era la fuente de filas, no el
  destino.
- **Concurrencia entre dos aplicaciones simultaneas.** No: el error 21000 es
  dentro de un mismo comando, no entre transacciones.
- **Datos corruptos en `channel_snapshot`.** Los snapshots estaban bien; lo que
  se repetia era el `product_id`.

## Causa raiz

La CTE `objetivo` de `aplicarPlantilla` (products/.../repository/dataapply.go)
partia de `product_business_integrations` y no deduplicaba por producto:

```sql
SELECT p.id AS product_id, ...
FROM product_business_integrations m
JOIN products p ON p.id = m.product_id
```

Con dos mappings para el mismo producto, `objetivo` traia dos filas iguales y el
`INSERT ... ON CONFLICT (product_id, field) DO UPDATE` intentaba actualizar la
misma fila dos veces en el mismo comando. Postgres lo prohibe.

El mismo defecto inflaba los conteos: el resumen y el preview contaban filas de
mapping (`COUNT(*)`), no productos, y por eso decian 897 en vez de 860.

## Correccion

Solo codigo, sin tocar datos:

- `dataapply.go`: `SELECT DISTINCT ON (p.id) ... ORDER BY p.id, m.snapshot_at DESC NULLS LAST, m.id DESC`.
  Ante dos publicaciones gana el snapshot mas reciente.
- `variant_apply.go`: mismo `DISTINCT ON` en `candidatosVarianteQuery`.
- `syncruns/.../datadiff.go`: `COUNT(*)` -> `COUNT(DISTINCT p.id)` en resumen y
  preview, y `DISTINCT ON (p.id)` en la consulta de muestras.

## Verificacion

- Preview de categoria desde Mercadolibre: 860 (antes 897).
- Apply real, batch `b49cf780-5ec9-472f-8d4c-7b184eff1057`: 860 productos, 860
  filas de historial, cero duplicados, sin error.
- Categorias resultantes: Buzos 166, Shorts 161, Camisetas y Polos 122,
  Pantalones y Sudaderas 88, Leggings 77, Camisetas 68, y 5 mas.

## Pendientes

- Los 37 productos con doble publicacion en Mercadolibre siguen ahi. Hoy gana el
  snapshot mas reciente, que es una decision razonable pero silenciosa: no se le
  avisa al usuario que ese producto tiene dos listings con datos distintos.
- No hay restriccion unica en `product_business_integrations` por
  `(business_id, integration_id, product_id)`, y no deberia haberla: el mapping
  multiple es legitimo. La deduplicacion tiene que vivir en cada consulta que
  asuma un producto por fila. Revisar ese supuesto antes de escribir una nueva.
