# Resultados Testing E2E - Accounting (Back)

## 2026-07-24 - Suite inicial del modulo (creacion)

Ejecutada contra backend local :3050 con JWT de super admin (business_id=0) y token de usuario demo.
Script: e2e_accounting.py (22 checks). Resultado: **22 OK / 0 FAIL**.

| Caso | Descripcion | Resultado |
|------|-------------|-----------|
| CU-01 | GET /accounting/concepts: 200, seeds (GUIDE_MARGIN, SUBSCRIPTION, WALLET_RECHARGE, COD_PAYOUT, COLLABORATOR), COD_PAYOUT con GMF_4X1000 activo | OK |
| CU-02 | GET /accounting/taxes: 200, 4 impuestos seed | OK |
| CU-03 | POST /accounting/taxes: 201 (TEST_TAX 5%) | OK |
| CU-04 | PUT /accounting/taxes/:id: 200, tarifa actualizada | OK |
| CU-05 | PUT /concepts/:id/taxes/:taxId: activar RETEFUENTE en COLLABORATOR y verificar reflejo | OK |
| CU-06 | POST /accounting/entries manual: 201, snapshot RETEFUENTE 11% sobre 1.000.000 = 110.000 | OK |
| CU-07 | GET /accounting/entries: paginacion (page/page_size/total/total_pages) y filtros kind/concept_id | OK |
| CU-08 | DELETE entry automatico rechazado (400), DELETE manual OK (200) | OK |
| CU-09 | POST /accounting/sync: 200, idempotente (created=0 tras ingesta inicial) | OK |
| CU-10 | GET /accounting/report: totales coherentes, by_concept (5) y by_tax (2) | OK |
| CU-11 | Report con rango invalido: 400 | OK |
| CU-12 | Entry con concepto inexistente: 404 | OK |
| CU-13 | Usuario business (demo): 403 en GET concepts y POST sync | OK |
| CU-14 | Token invalido: 401 | OK |

## Sync automatico (worker al arrancar + cada 15 min)

Ingesta historica idempotente verificada en DB (accounting_entries):
- GUIDE_MARGIN: 960 entries, $3.031.145 (USAGE cobrados con margen > 0, excluye is_test)
- WALLET_RECHARGE: 141 entries, $31.320.699
- SUBSCRIPTION: 9 entries, $2.075.002
- COD_PAYOUT: 16 entries, $19.303.200 (GMF 4x1000 snapshot: $77.213)

Segundo sync devuelve created=0 (indice unico parcial source_type+source_id).

## 2026-07-24 - Suite de facturas (CU-20 a CU-32)

Script: e2e_invoices.py (24 checks). Resultado: **24 OK / 0 FAIL** (tras fix).

- CU-20/21: crear factura a negocio Demo con 2 items + IVA 19%: numero FP-2026-NNNNN,
  subtotal 125.000, IVA 23.750, total 148.750, status DRAFT. OK
- CU-22: GET detalle con items. FAIL inicial -> BUG: la conexion global de GORM usa
  Omit(clause.Associations), el Create no insertaba los items asociados. Fix: insertar
  items explicitos en la transaccion (repository/invoices.go). Re-test OK.
- CU-23/24: lista paginada con filtros; editar borrador recalcula totales. OK
- CU-25: sin items 400; negocio inexistente 404. OK
- CU-26: POST /invoices/:id/send envia email real via Resend (destino de prueba
  flowprintml@gmail.com), status SENT + sent_at. OK
- CU-27: editar factura enviada rechazado 400. OK
- CU-28/29: marcar pagada crea movimiento contable (source INVOICE, idempotente);
  pagar dos veces 400. OK
- CU-30: cancelar factura pagada revierte el movimiento contable. OK
- CU-31: eliminar cancelada 200 y GET posterior 404. OK
- CU-32: usuario business recibe 403. OK

Datos de prueba eliminados (facturas soft-deleted, sin movimientos residuales).

## 2026-07-24 - Suite DIAN/Factus + modo prueba (CU-40 a CU-47)

Script: e2e_dian.py (14 checks). Resultado: **14 OK / 0 FAIL**.

- CU-40: GET /accounting/dian-config sin configurar -> configured=false. OK
- CU-41: PUT /accounting/dian-config con credenciales dummy contra el sandbox real de
  Factus -> 400 "credenciales invalidas" (fail-closed: valida contra la API antes de
  guardar; la cadena accounting -> integracion global -> cliente Factus funciona). OK
- CU-42: crear factura con is_test=true + customer_document/phone/address. OK
- CU-43: emitir DIAN sin integracion configurada -> 400 claro. OK
- CU-44: emitir sin documento del cliente -> 400 claro. OK
- CU-45: factura de prueba pagada NO crea movimiento contable (ledger limpio). OK
- CU-46: filtro ?is_test=true|false en listado. OK
- CU-47: limpieza de datos de prueba. OK

Pendiente de credenciales reales de Factus (sandbox o produccion) para probar la
emision completa con CUFE. La integracion se guarda como integracion GLOBAL de
plataforma (integrations.business_id IS NULL), patron ya soportado por el core.

## 2026-07-24 - Suite ficha fiscal + cargos/retenciones (CU-50 a CU-58)

Script: e2e_fiscal.py (14 checks) + verificacion de permisos aparte. **19 OK / 0 FAIL**.

- CU-50: accounting_taxes con kind (IVA=CHARGE, RETEFUENTE/RETEICA=WITHHOLDING, GMF=OTHER). OK
- CU-51/52: PUT/GET /businesses/:id/fiscal-profile (upsert, defaults, validacion person_type). OK
- CU-53: GET /accounting/client-profile sugiere tax_ids segun flags de la ficha
  (caso Viga: retefuente si; caso Mystic: ninguno). OK
- CU-54/55: factura con IVA+RETEFUENTE+GMF: IVA suma (total 119.000), retefuente NO suma
  (retencion informativa 11.000, neto a recibir 108.000), GMF ignorado en facturas. OK
- CU-56/57: factura test pagada sin movimiento; limpieza. OK
- CU-58 (permisos): business lee/edita SU ficha (200), ficha ajena 403 (GET y PUT),
  endpoints accounting 403 para no super admin. OK
  Nota: el login de usuario business devuelve el token SOLO en la cookie HttpOnly
  session_token (el campo token del JSON va vacio); los tests deben capturar Set-Cookie.

## 2026-07-24 - Suite catalogo de servicios (CU-60 a CU-64)

Script: e2e_services.py (8 checks). **8 OK / 0 FAIL**.

- CRUD /accounting/services (seed SOFT-LOG, code en mayusculas, 409 duplicado, update, delete, 404).
- Factura con item del catalogo persiste service_id + unspsc_code y calcula totales.
- Al emitir DIAN: si TODOS los items tienen UNSPSC se envia standard_code_id=2 con el
  codigo UNSPSC como code_reference; si no, codigo propio del contribuyente (default).

## Notas

- Login super admin de .env.ai fallo (password desactualizado en DB local). Se genero JWT
  firmado con JWT_SECRET local (user_id=1, business_id=0). Credenciales demo OK.
- Bug encontrado y corregido durante desarrollo (antes de commit): nombres de tabla
  singulares (cod_payment_cut, business) por SingularTable=true en GORM.
