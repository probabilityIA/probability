# Datos de prueba compartidos - Shipments

## Variables de entorno

| Variable | Valor | Descripcion |
|----------|-------|-------------|
| BASE_URL | `http://localhost:3050/api/v1` | API backend |
| MOCK_URL | `http://localhost:9091/api/v2` | Mock EnvioClick |
| SUPER_ADMIN_EMAIL | `${AI_SUPER_ADMIN_EMAIL}` | Super admin |
| SUPER_ADMIN_PASS | `Seb51923662#` | Super admin |
| DEMO_EMAIL | `demo@probability.com` | Usuario Demo |
| DEMO_PASS | `ProbabilityDemo` | Usuario Demo |
| BUSINESS_ID_DEMO | `26` | ID del negocio Demo |
| INTEGRATION_ID_DEMO | `44` | Integracion EnvioClick Demo |

## IDs generados durante las pruebas

| Variable | Valor | Generado en |
|----------|-------|-------------|
| SUPER_ADMIN_TOKEN | (obtener en CU-01) | CU-01 |
| DEMO_TOKEN | (obtener en CU-01) | CU-01 |
| CORRELATION_ID_QUOTE | (generado por backend) | CU-02 |
| SHIPMENT_ID_1 | (generado por backend) | CU-03 |
| TRACKING_NUMBER_1 | (generado por mock) | CU-03 |
| ENVIOCLICK_ID_ORDER | (generado por mock) | CU-03 |
| ORDER_UUID_TEST | `a9a03251-f01a-409f-b848-9de8598e84de` | Orden existente Demo |
| CORRELATION_ID_QUOTE_COD | (generado por backend) | CU-07 |
| SHIPMENT_ID_COD | (generado por backend) | CU-07 |

## Orden de ejecucion

1. CU-01: Setup test mode (activar is_testing en integracion Demo)
2. CU-02: Quote (cotizar envio)
3. CU-03: Generate guide (generar guia)
4. CU-04: Track (rastrear envio)
5. CU-05: Cancel (cancelar envio)
6. CU-06: Webhook (simular webhook EnvioClick)
7. CU-07: COD flow (flujo contra entrega completo)

## Datos de referencia

### Codigos DANE validos (mock)

| Codigo | Ciudad |
|--------|--------|
| 11001000 | Bogota D.C. |
| 05001000 | Medellin |
| 76001000 | Cali |
| 08001000 | Barranquilla |
| 68001000 | Bucaramanga |

### Orden de prueba existente (Demo business)

| Campo | Valor |
|-------|-------|
| UUID | `a9a03251-f01a-409f-b848-9de8598e84de` |
| order_number | `prob-0016` |
| customer_name | Carlos Arrieta |
| total_amount | 300000 |
| shipping_city | BOGOTA |
| business_id | 26 |
