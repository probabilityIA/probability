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

## Notas

- Login super admin de .env.ai fallo (password desactualizado en DB local). Se genero JWT
  firmado con JWT_SECRET local (user_id=1, business_id=0). Credenciales demo OK.
- Bug encontrado y corregido durante desarrollo (antes de commit): nombres de tabla
  singulares (cod_payment_cut, business) por SingularTable=true en GORM.
