# El guard del SSE de cotizacion no filtraba nada

Fecha: 2026-08-22
Modulo: front/central - `src/shared/ui/modals/shipment-guide-modal.tsx`

## Resumen

El modal de generacion de guia tenia codigo inalcanzable desde mayo que dejaba
`pendingCorrelationId` siempre en `null`. Como el guard de los handlers SSE de
cotizacion era `if (pendingCorrelationId && ...)`, con `null` la condicion daba
falso y **el guard no descartaba ningun evento**: cualquier cotizacion del mismo
negocio podia empujar el modal al paso 2 con tarifas ajenas.

## Sintoma

No hubo reporte de usuario. Salio de una corrida de `oxlint` sobre `src`
(`eslint(no-unreachable)` en las lineas 722-723 del archivo).

## Diagnostico

Cadena de evidencia:

1. `handleStep1Submit` hacia `setError("No hay transportadoras disponibles...")`
   + `return` y despues, sin poder ejecutarse nunca, asignaba
   `pendingStep1DataRef.current = data` y `setPendingCorrelationId(...)`.
2. `git log -S "No rates available"` -> commit `b71c1e31` (2026-05-11, SanCam04).
   El bloque se agrego a proposito; lo que quedo mal es que no borraron el codigo
   de abajo.
3. `grep setPendingCorrelationId` mostraba 4 llamadas, **todas con `null`**. La
   unica que asignaba un valor real era la linea muerta.
4. Los handlers SSE (`useShipmentSSE`, `onQuoteReceived` / `onQuoteFailed`) y el
   timeout de 30 s seguian montados y suscritos.

### Hipotesis descartadas

- **"La cotizacion esta rota porque el SSE no se ejecuta".** Falsa. El backend
  responde sincrono: `POST /shipments/quote` (`quote-shipment.go`) publica a
  RabbitMQ y hace polling a Redis con ticker de 200 ms y techo de 30 s. Las
  tarifas vuelven en la misma respuesta HTTP; el SSE nunca hizo falta.
- **"El corte de mayo fue un error".** Tambien falsa, y es lo mas util de esta
  entrada: el front NO PUEDE distinguir "202 aceptado, espera el SSE" de
  "cotizacion exitosa pero el carrier no dio tarifas". Las dos llegan como
  `success:true`, `rates` vacio y `correlation_id` presente. La ambiguedad esta
  en el contrato del backend, no en el componente.

## Causa raiz

Guard escrito con `&&` sobre un estado que en la practica siempre vale `null`.
El operador correcto es `||`: sin cotizacion pendiente hay que descartar el
evento, no aceptarlo.

## Correccion

`shipment-guide-modal.tsx`:

```
- if (pendingCorrelationId && data.correlation_id !== pendingCorrelationId) return;
+ if (!pendingCorrelationId || data.correlation_id !== pendingCorrelationId) return;
```

en `onQuoteReceived` y en `onQuoteFailed`. Ademas se borraron las dos lineas
inalcanzables. El SSE de generacion de guia (`pendingGuideCorrelationId`, que si
se asigna de verdad) no se toco.

Efecto: la ruta SSE de **cotizacion** queda inerte a proposito. La de **guia**
sigue operativa.

## Verificacion

- `tsc --noEmit`: sin errores en el archivo.
- `vitest run`: 739 tests en 53 archivos, todos en verde.
- `eslint` sobre el archivo: limpio.
- E2E manual en local contra el RDS de produccion por el tunel SSM, orden
  DEM-0044 (business 26), Bogota -> Neiva, 1 kg, contra entrega $65.000:
  `POST /api/v1/shipments/quote` -> **200 en 4.287 ms**, el modal avanzo al paso
  2 con 4 transportadoras (ENVIA $17.853, COORDINADORA $18.581, TCC $22.057,
  Servientrega $23.819). No se genero guia.

## Pendientes

1. **El backend no distingue "aceptado" de "sin tarifas".** Solo importa si
   Redis falla o esta ausente: ahi `runQuote` devuelve `quoteStatusAccepted`
   (`quote-shipment.go:144`, condicionado a `h.redisClient == nil`) y el usuario
   veria "No hay transportadoras disponibles" en lugar de esperar el SSE. La
   solucion es un campo `status` en la respuesta de `POST /shipments/quote`.
2. Reactivar la ruta asincrona del modal solo tiene sentido despues de (1).

## Notas de entorno encontradas de paso

Al probar en local, dos cosas ajenas a este caso:

- `back/central/.env` apunta al RDS por DNS publico, inalcanzable desde el
  2026-08-21. Hay que usar `./scripts/aws-tunnel.sh ensure` y `127.0.0.1:5433`.
- El contenedor `redis_local` corre con `--requirepass localdev` mientras
  `.env` trae la password de produccion: `WRONGPASS`, y la cotizacion moria en
  **408 tras 30 s** de polling fallido. El `docker-compose.yaml` de
  `infra/compose-local` ni siquiera define password: el contenedor esta
  desalineado del compose.
