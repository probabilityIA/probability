# Siigo no tiene endpoint de remisiones en su API

Fecha: 2026-08-19
Modulo: `back/central/services/integrations/invoicing/siigo`

## Resumen

Se pidio implementar remision de salida de mercancia contra Siigo. La API de
Siigo **no expone remisiones**: la remision existe en Siigo Nube pero solo por
UI. No hay endpoint contra el cual programar.

## Sintoma / punto de partida

El modulo Siigo cubre 11 operaciones (create, retry, check_status, cancel,
cash_receipt, credit_note, create_journal, compare, list_items,
list_bank_accounts, inventory_sync, list_siigo_warehouses) y ninguna es
remision. El cliente solo toca `/v1/{customers,invoices,credit-notes,journals,
vouchers,products,warehouses,payment-types,webhooks}`.

## Diagnostico: la cadena de evidencia

1. **Bloqueo de red.** Desde la sesion de Claude Code en la nube, la politica
   de egress de la organizacion bloquea el dominio de Siigo:
   `developers.siigo.com`, `api.siigo.com`, `*.portaldeclientes.siigo.com` y
   `siigoapi.docs.apiary.io` devuelven `CONNECT tunnel failed, 403`. Tambien
   `grep.app` y `postman.com`. GitHub, PyPI y npm si pasan.

2. **Hipotesis descartada #1: "el SDK oficial es la referencia".**
   `github.com/SiigoSAS/siigo_sdk_javascript` lista 18 clases y ~46
   operaciones, sin remisiones. Pero es de **octubre 2023** y esta
   desactualizado: no incluye `/v1/webhooks`, que nuestro propio cliente ya
   consume. Ausencia ahi no probaba nada.

3. **Hipotesis descartada #2: "buscar en Google la doc del endpoint".** Todas
   las busquedas devuelven paginas del portal de CLIENTES (como elaborar una
   remision desde la UI de Siigo Nube, "documento tipo S - Nota de remision" en
   Siigo Pyme), nunca documentacion de API. Confunde: hace parecer que la
   funcionalidad existe en la API cuando lo que existe es en el producto.

4. **Lo que si funciono.** `github.com/jdlar1/siigo-mcp` (MCP server no
   oficial, 75 tools, v3.2.0, commit 2026-05-05) trae **vendorizado el API
   Blueprint completo de Siigo**: `siigoapi.apib`, 6.029 lineas,
   `HOST: https://api.siigo.com`. Es el mismo documento que sirve Apiary, y se
   puede clonar y grepear porque GitHub no esta bloqueado.

## Causa raiz

Siigo no expone el recurso. Evidencia del blueprint:

```bash
git clone --depth 1 https://github.com/jdlar1/siigo-mcp
grep -ic remis      siigo-mcp/siigoapi.apib   # 0
grep -ic remission  siigo-mcp/siigoapi.apib   # 0
```

Tabla de recursos completa del blueprint: `/products`, `/customers`,
`/quotations`, `/invoices`, `/purchases`, `/credit-notes`, `/vouchers`,
`/payment-receipts`, `/journals`, `/purchase-support-documents`, reportes.

Enum `DocumentType` completo: `FV`, `RC`, `NC`, `FC`, `CC`. (Mas `DS`, `RP` y
`C` como valores de consulta de `/v1/document-types?type=`.) No hay tipo de
remision.

El blueprint esta al dia: trae los cambios del sector salud vigentes desde
2025-07-22, `/v1/invoices/batch`, `/v1/payment-receipts` y
`/v1/purchase-support-documents`, todos posteriores al SDK de 2023. Siigo si
agrega documentos nuevos a la API; remisiones no es uno de ellos.

## Correccion

Ninguna en codigo. Se documento:

- `.claude/plan/siigo-remision-salida.md` - evidencia, plan descartado y
  alternativas (la viable: remision solo en Probability, con PDF propio, sin
  enviarla a Siigo).
- README del modulo Siigo - donde vive la documentacion real y el comando para
  grepearla, mas la nota de que remisiones no existe.

## Verificacion

`grep -ic remis` sobre el blueprint = 0, y el enum de tipos de documento no
tiene remision. Reproducible con las dos lineas de arriba.

## Pendientes

- Residuo de incertidumbre: el blueprint tiene 3,5 meses. Confirmar con
  soporteapi@siigo.com que no haya salido nada entre mayo y hoy, y de paso
  pedir el endpoint como feature.
- Decision de negocio antes de construir la alternativa: la remision se
  necesita como papel que acompana la mercancia, o como descargue de
  inventario en Siigo antes de facturar? Lo segundo no es posible por API.

## Para la proxima

Cuando haya que responder "existe el endpoint X en Siigo?", no buscar en
Google ni confiar en el SDK de 2023: clonar `jdlar1/siigo-mcp` y grepear
`siigoapi.apib`. Diez segundos, y es la doc oficial completa.
