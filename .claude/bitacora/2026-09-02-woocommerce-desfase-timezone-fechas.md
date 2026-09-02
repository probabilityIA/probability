# WooCommerce: las fechas de orden se guardaban 5 horas corridas por timezone

Resumen: `date_created`/`date_paid` de WooCommerce vienen en hora local de la
tienda (Bogota) sin offset, y el backend los interpretaba como UTC sin
convertir. El instante absoluto guardado en `orders` quedaba 5 horas atrasado
respecto a la realidad, lo que hacia que ordenes pendientes recientes
aparecieran ordenadas por debajo de ordenes mas antiguas ya generadas.

## Sintoma

Cliente reporto: la orden #14684 (Luz Andrea Osorio Salcedo, WooCommerce)
aparecia en el listado principal con estado "Pendiente" por debajo de ordenes
ya generadas, como si fuera mas vieja. En pantalla se veian dos horas para el
mismo pedido: "2 sept 2026, 11:29 a.m." (creacion) y "WooCommerce 2 de sept,
04:30 p.m." (pagado) -- un aparente gap de 5 horas entre cotizar y pagar.

## Diagnostico

Se descarto primero la hipotesis de que el cliente realmente se hubiera
tardado 5 horas en pagar (carrito abandonado + retorno por remarketing de Ig).
Consulta real en produccion (solo lectura, tunel SSM) sobre la orden 14684:

```
order_number 14684
created_at  = 2026-09-02 16:29:05 UTC
occurred_at = 2026-09-02 16:29:05 UTC   (identico, al segundo)
paid_at     = 2026-09-02 16:29:48 UTC   (43 segundos despues)
```

Dos hallazgos:

1. `created_at` y `occurred_at` son identicos al segundo. Eso confirma que
   `orders/internal/infra/secondary/repository/mappers/mapper.go:18-19`
   sobreescribe el `created_at` (que deberia ser la hora real del INSERT) con
   `OccurredAt` (la fecha de WooCommerce, ya corrompida) cuando el mapeo trae
   `createdAt` vacio. El campo "Creado (DB)" que se muestra en el detalle NO
   es confiable para ordenes que entraron por WooCommerce.
2. `created_at` y `paid_at` distan 43 segundos, no 5 horas. La orden se creo y
   se pago casi simultaneamente. El "gap de 5 horas" que se veia en pantalla
   era el MISMO timestamp mostrado dos veces con dos tratamientos de zona
   horaria distintos: una vista le resta 5h (asumiendo que el valor guardado
   es UTC real) y muestra "11:29 a.m."; otra imprime los digitos crudos sin
   convertir y muestra "04:30 p.m.". Como el valor guardado en BD ya estaba
   mal etiquetado (era hora local de Bogota, no UTC), la resta de 5h en la
   primera vista duplica el error. La hora real del evento (creacion y pago)
   fue ~4:29 p.m. Bogota, no 11:29 a.m.

## Causa raiz

`back/central/services/integrations/ecommerce/woocommerce/internal/infra/secondary/client/response/woo_order_response.go`,
funcion `parseWooDate` (antes de la correccion). WooCommerce manda
`date_created`/`date_paid` en formato `"2006-01-02T15:04:05"` (hora local de
la tienda, sin offset). El fallback de parseo usaba
`time.Parse("2006-01-02T15:04:05", s)`, y Go le asigna UTC por defecto cuando
el string no trae zona. El struct tampoco capturaba `date_created_gmt` ni
`date_paid_gmt` (campos que WooCommerce SI manda en UTC real).

No hubo doble conversion -- fue ausencia total de manejo de timezone en todo
el paquete.

## Correccion

Mismo archivo. Se agregaron al struct `date_created_gmt`, `date_modified_gmt`,
`date_paid_gmt`, `date_completed_gmt`. El parseo ahora:

1. Si viene el campo `*_gmt`, se usa directo (ya es UTC real).
2. Si no viene, se interpreta el campo local con
   `time.ParseInLocation(..., America/Bogota)` en vez de asumir UTC.

No se corrigieron datos historicos en produccion (fuera de alcance de esta
sesion). Las ordenes ya guardadas antes del fix mantienen su `occurred_at`/
`paid_at`/`created_at` corridos 5 horas.

## Verificacion

`go build ./...` en `back/central` compila sin errores. No se probo aun contra
un webhook real de WooCommerce (pendiente antes de desplegar: forzar una orden
de prueba en `wordpress/` local o revisar el proximo webhook real en
staging/produccion y confirmar que `date_created_gmt` llega poblado).

## Pendientes

- Desplegar el fix.
- Decidir si vale la pena un script de correccion retroactiva para ordenes
  historicas de WooCommerce (recalcular `occurred_at`/`paid_at`/`created_at`
  sumando 5h a las que se guardaron antes del fix). Bajo impacto: solo afecta
  orden relativo dentro del listado y reportes por rango de fecha, no montos.
- Revisar si `mapper.go:18-19` (orders) deberia dejar de sobreescribir
  `created_at` con `OccurredAt` -- son conceptos distintos (cuando entro a
  Probability vs. cuando WooCommerce dice que se creo) y mezclarlos oculta
  bugs como este.
