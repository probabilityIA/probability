# 2026-08-06 - COD EnvioClick: valor a recaudar mal calibrado

## Resumen

Las 3 guias contra entrega generadas hoy para Viga (business 46) declararon un
valor menor al que corresponde y el negocio recibio 6.104 menos de lo debido.
Causa: EnvioClick deja de devolver la tarifa COD de Interrapidisimo cuando
`codValue` y `contentValue` no son iguales, y el codigo tapaba esa ausencia
usando la comision de otra transportadora.

## Sintoma

En Probability la orden VIG-0071 decia "Se cobrara contra entrega: COP 111.113"
(productos 95.500 + envio 15.613). En el panel de EnvioClick esa guia aparecia
con Monto 109.255. Diferencia: 1.858.

Las tres de hoy:

| Orden | Guia | Objetivo | Declarado (PDF) | Comision | Neto | Falta |
|-------|------|----------|-----------------|----------|------|-------|
| VIG-0071 | 240058568295 | 111.113 | 117.229 | 7.974 | 109.255 | 1.858 |
| VIG-0072 | 240058578276 | 90.254 | 96.057 | 6.714 | 89.343 | 911 |
| VIG-0069 | 240058579495 | 135.929 | 142.045 | 9.451 | 132.594 | 3.335 |

Las de ayer (VIG-0066, 0067, 0068) y las de julio cerraron en 0. El problema
aparece solo despues del deploy del commit `aab3a08f` (2026-08-05 17:35).

## Como funciona el mecanismo

EnvioClick descuenta su comision de recaudo del monto recaudado. Para que al
negocio le llegue `productos + envio`, hay que declarar de mas: se cobra al
comprador `objetivo + comision`. `ResolveCODValue`
(`envioclick/internal/app/cod_value.go`) cotiza dos veces para medir la comision
real, despeja el valor a declarar y verifica con una tercera cotizacion que
`declarado - comision >= objetivo`.

## Diagnostico

Evidencia principal: `shipment_sync_logs`, que guarda request y response de cada
llamada a EnvioClick.

En las 9 cotizaciones de sondeo (3 por guia) EnvioClick devolvio la tarifa
normal de INTERRAPIDISIMO pero **sin `codDetails`**. `FindCODFee` no fallaba:
caia a un fallback que tomaba la primera tarifa con `codDetails` de la lista, o
sea COORDINADORA (6.116, plana) o Servientrega (5.512). Calibro con la tarifa de
otra transportadora y declaro de menos.

### Hipotesis descartadas (importante, no repetirlas)

1. **"El payload enriquecido de generacion suprime la tarifa"** (`suburb`,
   nombres, telefono, `pickupDate`). Falso: se cotizo con el payload exacto del
   front, con `suburb` vacio, y tampoco aparecia.
2. **"La disponibilidad va y viene en su API"**. Falso: es determinista.
3. **"Se demora entre cotizar y generar y la tarifa expira"**. Falso: los gaps
   reales fueron 84s, 37s y 40s.
4. **"La tarifa de comision cambia"**. Falso: la formula de Interrapidisimo
   ajusta exacto desde el 27/07 hasta hoy en 8 puntos, mezclando cotizaciones y
   cobros reales.

## Causa raiz

**EnvioClick solo devuelve el producto COD de Interrapidisimo cuando `codValue`
es exactamente igual a `contentValue`.** Un peso de diferencia basta para que
desaparezca. A Coordinadora, Envia y Servientrega no les afecta.

Comprobado con 15 cotizaciones contra la API real:

| codValue | contentValue | Interrapidisimo COD |
|----------|--------------|---------------------|
| 117.229 | 117.229 | 7.975 |
| 117.229 | 130.000 | NO |
| 109.471 | 109.470 | NO |
| 105.000 | 109.470 | NO |
| 88.920 | 88.920 | 6.290 |
| 96.057 | 88.920 | NO |

`probeCODFee` subia `CODValue` para sondear y dejaba `ContentValue` en el
original. Por eso los 9 sondeos perdian la tarifa. La cotizacion del front nunca
fallo porque manda los dos iguales.

Encaja con todos los registros historicos, sin excepcion.

## Tarifas medidas

Formula que ajusta exacto en todos los puntos observados:

- **Interrapidisimo**: `floor((840 + 0,05 * declarado) * 1,19)`
- **Coordinadora**: `floor((840 + 0,02 * declarado) * 1,19)`, minimo 6.116

No estan quemadas en el codigo (a proposito, ver `aab3a08f`); quedan aca solo
como referencia para verificar calculos a mano.

## Correccion de codigo

- `envioclick/internal/app/cod_value.go`: `probe.ContentValue = declared` en
  `probeCODFee`. Esta es la causa raiz.
- `envioclick/internal/domain/cod_calibration.go`: `FindCODFee` recibe el
  `rateID` de la guia, busca primero por ese id, luego por nombre de
  transportadora, y **no cae nunca a la tarifa de otra**. Si no encuentra la del
  carrier del envio devuelve `false`, la calibracion se aborta con un Warn y se
  conserva el valor sin calibrar.
- Tests en `cod_calibration_test.go`: prioridad del `idRate` y el caso exacto de
  hoy (Coordinadora y Servientrega presentes, Interrapidisimo sin `codDetails`
  -> debe fallar, no calibrar).

## Verificacion

Dos guias reales generadas en el negocio Demo (business 26) desde el backend
local con el arreglo, y canceladas despues:

| Guia | Transportadora | Objetivo | Declarado (PDF) | Comision | Neto |
|------|----------------|----------|-----------------|----------|------|
| 240058587632 | Interrapidisimo | 97.250 | 104.466 | 7.215 | 97.251 |
| 84151705886 | Coordinadora | 97.250 | 103.366 | 6.116 | 97.250 |

En ambas los 3 sondeos devolvieron la tarifa del carrier correcto. Tarifa
proporcional en una y plana en la otra: el algoritmo resuelve las dos.

## Correccion de data

Migracion `fixVigaCodCalibracionFallida`, corrida en produccion el 2026-08-06:

| Orden | cod_total | cod_carrier_fee |
|-------|-----------|-----------------|
| VIG-0071 | 111.113 -> 109.255 | 7.513 -> 7.974 |
| VIG-0072 | 90.254 -> 89.343 | 6.290 -> 6.714 |
| VIG-0069 | 135.929 -> 132.594 | 8.968 -> 9.451 |

Los valores se validaron por dos vias independientes: el panel de EnvioClick y
una cotizacion por el valor declarado exacto de cada guia. Coinciden al peso.

## Hallazgos secundarios

- **11 ordenes de Viga de julio siguen descuadradas** (VIG-0029 a VIG-0050,
  aprox. 22.600 en total). Son del problema original que motivo `aab3a08f`, que
  solo corrigio las entregadas y sin corte. Estas ya estaban dentro de un corte
  cerrado. **Decision del 06/08: no se corrigen**, lo pagado ya se pago.
- La respuesta de generacion de guia **no trae la comision**, solo hace eco del
  `codValue`. El tracking tampoco. La cotizacion es el unico punto donde
  EnvioClick informa cuanto cobra de recaudo.
- Hay negocios con transportadoras sin margen configurado (visto en Demo con
  Coordinadora: `applied_margin = 0`). Vale revisar si pasa en negocios reales.
- La integracion de LaPerchaDel10 (#49) tiene llave propia de sandbox y
  `use_platform_token: false`, a diferencia del resto que usa la llave de
  plataforma. Durante la sesion se le apago el modo test y quedo inoperante.

## Pendientes

- [ ] Deploy del arreglo (sin el, cada guia COD de Interrapidisimo sigue
      perdiendo entre 900 y 3.500).
- [ ] Devolver la integracion #49 de LaPerchaDel10 a su estado, o pasarla a
      `use_platform_token: true`.
- [ ] Evaluar mover la calibracion al momento de la cotizacion en vez de al
      generar la guia.

## Donde mirar si vuelve a pasar

```sql
SELECT id, operation_type, created_at,
       (request_payload->>'codValue')::numeric  cod,
       (request_payload->>'contentValue')::numeric content,
       (request_payload->>'idRate') idrate
FROM shipment_sync_logs
WHERE shipment_id = <id>
ORDER BY created_at;
```

En los 3 primeros registros (los sondeos) `cod` y `content` deben ser iguales, y
la respuesta debe traer `codDetails` para la transportadora de la guia. Si no,
la calibracion no midio y hay que revisar el log del Warn
"La cotizacion no devolvio la comision COD de esta transportadora".
