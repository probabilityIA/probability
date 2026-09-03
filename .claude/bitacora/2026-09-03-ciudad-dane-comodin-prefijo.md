# 2026-09-03 - La ciudad del checkout se resolvia por prefijo y despachaba a otro municipio

**Ticket:** TKT-000072 | **Canal:** WooCommerce (todos los negocios) | **Tambien afecta:** Shopify

## Sintoma

La orden 161 del negocio Demo dice **Suba** en la direccion de envio, en la de
facturacion y en el correo al comprador. Por debajo se cotizo y se habria
despachado a **SUBACHOQUE**, un municipio de Cundinamarca a 45 km.

Se noto porque el envio salio caro: 22.752 de flete contra los 18.763 que
mostraba el cotizador del panel para la misma orden.

## Causa raiz

`GetCityDaneByName` (`shipments/.../repository/shopify_quote_queries.go`)
intentaba tres coincidencias contra `geozones`:

```sql
unaccent(lower(g.name)) = unaccent(lower(@city))            -- exacta
OR unaccent(lower(g.name)) LIKE unaccent(lower(@city))||',%' -- "BOGOTA, D.C."
OR unaccent(lower(g.name)) LIKE unaccent(lower(@city))||'%'  -- COMODIN
```

El tercero es el problema: ordena por longitud del nombre y toma el primero, asi
que siempre elige el municipio mas corto que empiece igual. "Suba" no existe como
municipio (es una localidad de Bogota), asi que engancho SUBACHOQUE.

El texto de la ciudad nunca se corrige: solo cambia el `daneCode`, que es lo
unico que EnvioClick usa para tarifar y despachar. Por eso nada en pantalla lo
delata.

## Evidencia

Las dos cotizaciones de la misma orden, con 16 segundos de diferencia:

| Origen | daneCode destino | Flete ENVIA |
|---|---|---|
| Checkout Woo (quote 17000) | 25769000 = Subachoque | 22.752 |
| Panel (quote 17001) | 11001000 = Bogota | 18.763 |

## Hipotesis descartadas

- **"Lo inflo la calibracion COD"**: no. Cotizando en prepago, que no pasa por la
  calibracion, el flete da los mismos 22.752.
- **"Depende del monto declarado"**: no. Cotizando a Bogota con declarado de
  195.000 y de 232.477 el flete da 18.763 en ambos.
- **"Google lo resuelve"**: no. `validate-address` si geocodifica, pero despues
  llama a la MISMA funcion y vuelve a caer en el comodin. Ademas corre despues de
  la cotizacion y no corrige el dane ya usado.

## Alcance medido

WooCommerce + Shopify, ultimos 3 meses, 22.090 ordenes:

| | Ordenes |
|---|---|
| Resuelven exacto o por "ciudad, algo" | 20.628 |
| Resuelven **solo** por el comodin | 220 |
| No resuelven de ninguna forma | 1.242 |

De las 220 del comodin:

```
Cartagena -> CARTAGENA DE INDIAS (13001)   158   correcto
Buga      -> BUGALAGRANDE (76113)           20   INCORRECTO, es 76111
Suba      -> SUBACHOQUE (25769)             19   INCORRECTO, es Bogota 11001
"Ciudad"  -> CIUDAD BOLIVAR                  5   basura
Bogot/Bogo/Bog/Bo -> BOGOTA / BOYACA         5   mezclado
Santander/Villa/Puerto/casa/cauca            5   adivinanza pura
```

Lo de **Buga -> Bugalagrande** es la otra mitad de la causa del incidente del
2026-09-02 (`2026-09-02-orden-geocodificada-en-buga-siendo-zarzal.md`): no fue
solo la coordenada pegada de Google, el dane tambien estaba mal.

## Correccion

Se decidio **no usar tabla de alias**: si el nombre no coincide, es error, pero
se le ofrecen al comprador las opciones validas.

1. Fuera el comodin. Quedan la coincidencia exacta y "ciudad, algo", con
   `unaccent`, `lower` y `trim`.
2. Buscador de ciudad en el checkout (plugin 1.7.0): se consultan las ciudades
   reales del departamento mientras el comprador escribe y el codigo DANE
   elegido viaja por `extensionCartUpdate` hasta el payload de cotizacion.
3. El endpoint de tarifas responde `address_error` con ciudades parecidas en vez
   de una lista vacia sin explicacion.

Commits `a6f8c8fb` y `a2fe9782`.

## Bug latente encontrado al probar

`wp_enqueue_script` registraba los scripts con la cadena `'1.6.7'` quemada, asi
que el navegador servia la copia cacheada aunque los archivos en disco fueran los
nuevos. **Cualquier actualizacion futura del plugin habria tenido el mismo
problema.** La version pasa a la constante `PROBABILITY_SHIPPING_VERSION`.

## Costo para clientes con plugin viejo

El 93% de las ordenes no cambia. Se pierden las que resolvian por comodin: en 60
dias fueron 8 de Cartagena en Moto Mello, 1 de Buga en Moto Mello y 1 de
Cartagena en Viga. Una cada seis dias. Se cambio una perdida silenciosa
(paquetes a otro municipio) por una falla visible y recuperable.

## Contexto util

WooCommerce ya lista pais y departamento (34 opciones) de fabrica, pero **no trae
los municipios de Colombia**. Ese es el hueco que llena este cambio. Ningun
cliente listaba ciudades: verificado en la tienda real de Viga, donde
`shipping-city` es un `input` de texto libre.

## Pendiente

- [ ] Probar el buscador de punta a punta: al cerrar el ticket la tienda todavia
      servia el JS cacheado.
- [ ] Actualizar el plugin en Moto Mello y Viga.
- [ ] Checkout clasico: la ciudad sigue siendo texto libre.
- [ ] El aviso de WooCommerce cuando no hay tarifas sale en ingles porque la
      traduccion de Blocks no esta completa en la tienda. Se puede reemplazar
      desde nuestro JS.
