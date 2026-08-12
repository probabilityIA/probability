# Alerta: Siigo - productos nuevos quedan sin nivel de inventario

Fecha: 2026-08-12
Modulo: `back/central/services/modules/inventory`
Detectado en: business 46 (Viga), integracion 228

## Contexto

Comparando los 2.166 productos de Siigo contra Probability, 2.120 tienen stock
identico (98%). El mapeo funciona. Pero 38 productos existen aqui **sin ninguna
fila en `inventory_levels`**, teniendo 285 unidades en Siigo.

Ejemplos: `BT150-S` (40 en Siigo, 0 aqui), `BT150-XL` (35), `BT150-L` (23),
`BN150-S` (20). Todos creados el 2026-08-05 17:03:52 y el 2026-07-28.

## Causa: carrera entre creacion de producto y sync de stock

`internal/app/sync_provider_stock.go:66`, en `applyProviderStockItem`:

```go
productID, _, track, err := uc.repo.GetProductBySKU(ctx, item.SKU, businessID)
if err != nil || productID == "" || !track {
    result.Skipped++
    return
}
```

El producto se crea por la cola `products.provider_upsert.requests` (async) y su
stock se procesa en la misma corrida del sync. Si el producto todavia no existe,
el item se cuenta como `Skipped` y **no se reintenta nunca**.

Los 38 se crearon en el mismo segundo en que corria ese sync. Llegaron tarde.

No es que falte crear el nivel a mano: `AdjustStockTx` ya llama a
`getOrCreateLevelKeyTx`, que crea la fila si no existe. El item simplemente
nunca llego ahi.

## Items

### IMPORTANTE

1. **Que los `Skipped` por producto inexistente se reintenten.** Opciones:
   reencolar el item, o hacer una segunda pasada al final de la corrida con los
   saltados. Hoy se pierden en silencio: el run termina en `completed` y el
   contador de `skipped` no distingue "producto no existe todavia" de "producto
   sin `track_inventory`".
2. **Separar los motivos de `Skipped`** en el resultado del sync, para que el
   negocio vea en la UI cuales quedaron sin aplicar y por que.

### DESEABLE

3. Correr el sync de inventario de Siigo del business 46 para destapar las 285
   unidades. Los productos ya existen, esta vez si los encuentra. No requiere
   cambio de codigo.

## Dato aparte, mismo business

16 productos existen aqui y no en Siigo: `PROD-001` y 15 `SDH*` creados por el
product sync de MercadoLibre el 2026-08-12 17:09 (SKU con las letras invertidas,
`SDH` en vez de `SHD`). Ver `.claude/alerts/meli-variantes-sin-sync.md`.

## Criterio para cerrar

- Un producto creado en la misma corrida en que llega su stock termina con su
  nivel de inventario aplicado.
- El resultado del sync distingue los motivos de `Skipped`.
