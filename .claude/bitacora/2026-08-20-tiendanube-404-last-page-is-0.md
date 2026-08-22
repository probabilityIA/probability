# Tiendanube devuelve 404 cuando el listado esta vacio

Tiendanube responde `404 Not Found` con `"description": "Last page is 0"` al
listar recursos cuando no hay resultados, en vez de devolver `[]`. Tratarlo como
error rompe la sincronizacion de una tienda sin ordenes.

## Contexto

Fecha: 2026-08-20. Integracion 259 (`sebas y corotos`), business 26 (Demo),
store_id `8126740`, conectada por OAuth ese mismo dia.

Primera ejecucion de `POST /api/v1/integrations/259/sync` sobre una tienda recien
creada y sin ordenes.

## Sintoma

```
Tiendanube: tiendanube client: GET /orders returned 404: {
    "code": 404,
    "message": "Not Found",
    "description": "Last page is 0"
}
```

La sincronizacion abortaba con error, sin publicar nada y sin dejar claro que el
problema era simplemente que no habia ordenes.

## Hipotesis descartadas

1. **Credenciales o token invalido.** Descartada: `GET /integrations/259/webhooks`
   contra la misma tienda respondia `200` con `[]`. El token funcionaba.
2. **URL mal construida (falta el store_id).** Descartada: el cliente arma
   `https://api.tiendanube.com/v1/8126740/orders` y ese es el formato correcto;
   ademas el endpoint de webhooks usa el mismo `do()` y funciona.
3. **Scope insuficiente.** Descartada: el scope guardado en el config de la
   integracion incluye `read_orders` y `write_orders`. Un scope faltante en
   Tiendanube devuelve 401/403, no 404.
4. **Dominio equivocado (`nuvemshop.com.br` en vez de `tiendanube.com`).**
   Descartada por el mismo motivo que 2: otros endpoints responden bien.

## Causa real

Es el comportamiento de paginacion de Tiendanube: al pedir `page=1` sobre una
coleccion vacia, la ultima pagina es 0, la pagina 1 queda fuera de rango y la API
responde 404. **No es un error, es una coleccion vacia.**

La documentacion de paginacion no lo menciona.

## Correccion

`internal/infra/secondary/client/constructor.go`: el `do()` ahora devuelve
`domain.ErrResourceNotFound` envuelto para cualquier 404, en vez de un error
generico de texto.

`internal/infra/secondary/client/orders.go`: `GetOrders` corta el bucle de
paginacion con `errors.Is(err, domain.ErrResourceNotFound)` y devuelve lo
acumulado. Una tienda sin ordenes ahora termina en `fetched=0 published=0`.

Verificado tras el fix:

```
INF Sincronizacion de ordenes de Tiendanube completada fetched=0 published=0
    function=SyncOrders integration_id=259 module=tiendanube
```

## Cuidado al extender

`GetOrder` (una sola orden) **sigue propagando** el 404: ahi si es un error real,
la orden no existe. La tolerancia al 404 aplica solo a los listados paginados.

Si se agregan otros listados paginados contra Tiendanube (productos, clientes,
cupones), hay que aplicar el mismo corte por `ErrResourceNotFound`, o volveran a
fallar contra una coleccion vacia.
