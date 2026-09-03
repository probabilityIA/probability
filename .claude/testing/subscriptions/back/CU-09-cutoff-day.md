# CU-09: Dia de corte fijo del mes (subscription_cutoff_day)

## Objetivo
Verificar que configurar un dia de corte pospone el bloqueo hasta la proxima
ocurrencia de ese dia, en vez de bloquear apenas vence end_date.

## Precondiciones
- SUPER_TOKEN obtenido (CU-01).
- Negocio de prueba.

## 9.1 Configurar cutoff_day via edit-dates
```
PUT /api/v1/subscriptions/edit-dates
Authorization: Bearer {SUPER_TOKEN}
Content-Type: application/json

{
  "business_id": 26,
  "start_date": "{hoy}",
  "end_date": "{ayer}",
  "cutoff_day": {dia_futuro_cercano}
}
```
- [ ] Status 200
- [ ] `business.subscription_cutoff_day` quedo seteado
- [ ] Auditoria `dates_edited` incluye "dia de corte: N"

## 9.2 cutoffReached respeta el dia de corte (verificado por unit test, no por worker en E2E)
El worker `CheckExpiringSubscriptions` corre por cron, no hay endpoint para
dispararlo manualmente en caliente salvo llamarlo directo en un test. Verificar
la logica via los tests unitarios existentes (`subscription_cutoff_test.go`) y,
si se quiere una prueba E2E real, invocar el usecase manualmente en un entorno
de test o esperar al cron local.

## 9.3 EditSubscriptionDates con cuenta expired y nueva fecha futura: reactiva
- [ ] Marcar negocio como expired
- [ ] Ejecutar edit-dates con end_date futuro
- [ ] `business.subscription_status` vuelve a `active`
- [ ] Auditoria `subscription_reactivated` adicional
