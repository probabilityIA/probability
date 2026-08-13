# 2026-08-12 - El detector de SKU sugeria cambiar de talla

## Contexto

Negocio Viga (`business_id` 46). El comparativo de Mercado Libre mostraba 19
"posibles errores de digitacion". Al revisarlos contra la base, **17 de 19 eran
falsos** y varios eran consejos daninos.

## Que estaba mal

El detector comparaba el SKU completo con distancia de edicion 1 y cantidades
parecidas. En un catalogo de ropa con tallas eso rompe de dos formas:

1. **Proponia cambiar la talla.** `TD22 - 5XL -> TD22-6XL` no es una correccion
   de digitacion: son dos productos distintos. Igual `TD250 - 6XL -> TD250-5XL`.
2. **Elegia al azar entre vecinos.** Con SKU correlativos (TD21, TD212, TD213,
   TD214, TD22, TD25, TD250, TD26) casi todo queda a un caracter de algo. Para
   `TD212 - 4XL` existian al menos dos candidatos validos y el detector se
   quedaba con el primero que encontraba.

## La causa real que tapaba

Lo importante: **la mayoria de esos SKU no eran errores de digitacion, sino
espacios de sobra.** ML tenia `BN14 -5XL` y en Probability existe `BN14-5XL` con
10 unidades. El emparejamiento no normaliza espacios interiores, el producto
quedaba sin pareja, y el detector le inventaba un "parecido" cualquiera.

Consulta que lo dejo en evidencia:

```sql
-- SKU del canal con espacio que SI existen en Probability al quitarlo
WITH canal AS (
  SELECT DISTINCT it.sku, it.group_code
  FROM integration_sync_run_items it
  JOIN integration_sync_runs r ON r.id = it.run_id
  WHERE r.business_id = 46 AND it.group_code IN ('only_channel','sku_typo')
    AND it.sku LIKE '% %'
), prob AS (
  SELECT REPLACE(UPPER(TRIM(p.sku)), ' ', '') AS k, COALESCE(SUM(il.quantity),0) AS qty
  FROM products p LEFT JOIN inventory_levels il
    ON il.product_id = p.id AND il.deleted_at IS NULL
  WHERE p.business_id = 46 AND p.deleted_at IS NULL GROUP BY 1
)
SELECT canal.group_code, COUNT(*) AS con_espacio, COUNT(prob.k) AS resuelve_sin_espacio,
       SUM(COALESCE(prob.qty,0)) AS unidades
FROM canal LEFT JOIN prob ON prob.k = REPLACE(UPPER(TRIM(canal.sku)), ' ', '')
GROUP BY canal.group_code;
```

Resultado: **53 de 58** SKU con espacio resolvian quitandolo, 35 con stock real,
**457 unidades** que no estaban sincronizando.

## Que se corrigio

- El SKU se parte en codigo y talla por el ultimo guion. **La talla tiene que ser
  identica**; nunca se sugiere cambiarla.
- Si mas de un candidato califica, no se sugiere nada. Callarse es mejor que
  acertar por casualidad.
- Los espacios pasaron a ser su propia categoria (`sku_spacing`), con confianza
  alta e instruccion exacta, separada de las sospechas (`sku_typo`).
- Cada hallazgo guarda la evidencia: SKU y nombre de los dos lados, stock de cada
  lado, patron y de que lado conviene corregir.

Resultado en la misma corrida: **19 sospechas -> 2**, mas 38 espacios ciertos.

## Autoridad: quien manda cuando los dos lados diferen

Si el negocio tiene Siigo activo con `inventory_sync_enabled`, el catalogo de
Probability viene del ERP y la correccion va **del lado del canal**. En el
comparativo del propio Siigo la autoridad es Siigo, asi que la correccion va del
lado de **Probability**. Sin ERP se mira que lado tiene existencias reales.

Los espacios son la excepcion: se corrigen donde esta el espacio, sin importar la
autoridad. Es un hecho objetivo, no una inferencia.

## Trampa al sumar Siigo al cruce entre canales

Al darle comparativo a Siigo, los hallazgos que cruzan canales lo contaron como
canal de venta y **`not_published` cayo de 443 a 1**: un producto que solo existe
en Siigo pasaba a estar "publicado". El cruce ahora filtra por
`integration_categories.code = 'ecommerce'`, asi Siigo participa como integracion
pero no como vitrina.

## Hipotesis descartadas

- *"Los SDH/SHD son el patron principal"*: eran 15 casos reales, pero en la
  corrida siguiente ya no aparecian y el grueso del ruido era otra cosa.
- *"Subir el umbral de distancia de edicion arregla el ruido"*: no. El problema
  no era la distancia sino comparar la talla y no resolver ambiguedades.
