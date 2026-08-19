# CU-04: Comprar suscripcion (flujo business normal)

## Objetivo
Verificar el flujo de auto-compra de suscripcion via wallet del negocio (`PurchaseSubscription`).

## Precondiciones
- DEMO_TOKEN obtenido (CU-01).
- SUBSCRIPTION_TYPE_ID de un plan activo con precio bajo (CU-02.1).
- Verificar previamente el saldo de wallet del negocio 26 via
  `mcp__postgres-probability__query` (SELECT balance FROM wallets WHERE business_id = 26)
  para saber si alcanza o se espera `ErrInsufficientBalance`.

## Endpoint
```
POST /api/v1/subscriptions/purchase
Authorization: Bearer {DEMO_TOKEN}
Content-Type: application/json
```

## Body
```json
{
  "subscription_type_id": {SUBSCRIPTION_TYPE_ID},
  "months": 1
}
```

## Verificaciones (caso saldo suficiente)
- [ ] Status 200, `success` implicito por `data` presente
- [ ] `data.business_id == 26`
- [ ] `data.subscription_type_id == SUBSCRIPTION_TYPE_ID`
- [ ] Verificar en DB (SELECT) que se desconto el saldo del wallet exactamente
      `price * months`
- [ ] Verificar en DB que la fila de `business_subscriptions` quedo con
      `start_date`/`end_date` coherentes con `months`

## Verificaciones (caso saldo insuficiente)
- [ ] Status 402 (Payment Required)
- [ ] `error == "insufficient wallet balance"`
- [ ] Verificar que NO se creo/modifico ninguna fila de suscripcion ni se desconto saldo
      (si se descuenta saldo sin crear suscripcion o viceversa, es un bug de transaccion)

## Caso de seguridad: business_id ajeno ignorado
El `PurchaseSubscriptionRequest` no acepta `business_id` en el body (solo
`subscription_type_id` y `months`), asi que no aplica el vector de inyeccion por body.
Confirmar leyendo la respuesta que `data.business_id` siempre es 26 (el del token),
nunca otro valor aunque se intente mandar `"business_id": 999` extra en el JSON (debe
ser ignorado silenciosamente por el binding de Gin).
