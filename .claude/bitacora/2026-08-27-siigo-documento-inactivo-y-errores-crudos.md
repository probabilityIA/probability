# Siigo: documento inactivo, documento sin electronica, y errores crudos en pantalla

Fecha: 2026-08-27
Business: 46 (Viga ropa deportiva)
Modulo: `back/central/services/integrations/invoicing/siigo`, `front/central` modulo invoicing

## Resumen

Viga no podia facturar. Los 400 de Siigo llegaban a la pantalla como JSON crudo
del proveedor, asi que nadie sabia que hacer con ellos. El problema de fondo
NO era del codigo: eran dos configuraciones del lado de la cuenta Siigo del
cliente. Ademas cada error permanente consumia los 3 reintentos automaticos.

## Evidencia (invoice_sync_logs, business 46)

| document.id | Code de Siigo | Mensaje | Cuando | # |
|---|---|---|---|---|
| 2801 | `document_settings` | "The send cannot be used, you must verify the document settings" (`stamp.send`) | 20 al 27 ago | 35 |
| 2801 | `parameter_inactive` | "The code is inactive: BN116-2XL" (`items[1].code`) | 27 ago 13:58-15:00 | 4 |
| 2801 | `documents_service` | "The Documents service is currently unavailable" | 26 ago | 2 |
| 30606 | `parameter_inactive` | "The id is inactive: 30606" (`document.id`) | 27 ago 18:54-18:57 | 2 |

Consulta para reproducir:

```sql
SELECT l.request_payload->'document'->>'id' AS doc_id,
       l.response_body->'Errors'->0->>'Code' AS code,
       l.response_body->'Errors'->0->>'Message' AS msg,
       count(*), min(l.created_at), max(l.created_at)
FROM invoice_sync_logs l JOIN invoices i ON i.id = l.invoice_id
WHERE i.business_id = 46 AND l.operation_type = 'create'
GROUP BY 1,2,3 ORDER BY 6 DESC;
```

## Causa

1. `2801` esta activo pero **no habilitado para facturacion electronica** en
   Siigo (resolucion DIAN / envio). Por eso rechaza `stamp.send`.
2. `30606` es el documento electronico nuevo del cliente y estaba **inactivo**
   en Siigo cuando lo configuraron (17:49). Lo activaron mas tarde: a las 20:52
   la factura FV-2-1 salio con 201.

Los dos se arreglan en Siigo Nube. **La API de Siigo no permite activar un tipo
de documento**: `/v1/document-types` es solo GET. Confirmado en la doc oficial.

## Hipotesis descartadas

- "Es un bug del payload nuestro": no. El mismo payload con el documento ya
  activo devolvio 201 (logs 87884, 88184, 88190).
- "Es el mapeo de productos": el `parameter_inactive` de `items[1].code` es otro
  caso distinto (SKU BN116-2XL inactivo en Siigo), no la causa del bloqueo.

## Que se cambio

- `client/siigo_errors.go` (nuevo): traduce el sobre de error de Siigo a un
  mensaje en espanol accionable y a un **codigo canonico** (`document_inactive`,
  `document_not_electronic`, `product_inactive`, `provider_unavailable`,
  `invalid_credentials`, ...). Enganchado en create, credit note y cash receipt.
- `domain/errors/provider_error.go` (nuevo): error tipado con codigo, para que
  el consumer lo publique en `ErrorCode` en vez de `api_error` para todo.
- `client/mappers/document_preview.go` (nuevo): mapea la respuesta de Siigo al
  documento canonico que consume el modulo de facturacion (`document_json`), con
  items, pagos, totales, `public_url`, `stamp_status` y el crudo en `raw`.
  Antes Siigo no mandaba `document_json`, por eso sus facturas salian sin vista
  previa mientras las de Softpymes si.
- `modules/invoicing` response consumer: si el codigo es de configuracion, la
  factura falla sin programar reintentos. Antes cada error permanente gastaba
  los 3 reintentos y llenaba el log.
- Front: `domain/invoice-preview.ts` normaliza canonico y el formato viejo de
  Softpymes; `ui/components/InvoicePreview.tsx` pinta la vista previa igual para
  todos los proveedores.

## Que devuelve Siigo al crear (para futuras vistas)

`id`, `name` (FV-2-3), `number`, `prefix`, `date`, `document.id`, `customer`,
`items[]` (code, description, quantity, price, total), `payments[]`, `total`,
`balance`, `observations`, `stamp.status`, `mail.status`, `metadata.created`,
`public_url` (visor web del documento).

**El CUFE no viene en la respuesta de creacion**: `stamp.status` llega en
`Sending`. Se obtiene despues por `GET /v1/invoices/{id}` (check_status).
Para el PDF existe `GET /v1/invoices/{id}/pdf`, que devuelve el archivo en
base64. Todavia no lo consumimos.
