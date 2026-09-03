# CU-15: SubscriptionGuard (bloqueo de cuenta)

## Objetivo
Verificar que la cuenta se bloquea visualmente cuando `subscription_status` es
expired/cancelled, que /subscription sigue accesible, y que el desbloqueo
requiere volver a loguearse (JWT stale, comportamiento documentado).

## Precondiciones
- Backend local, Demo (business 26) deshabilitado via API antes de loguear.

## 15.1 Cuenta expirada bloquea la app
- [ ] Deshabilitar Demo via API (`POST /subscriptions/disable?business_id=26`)
- [ ] Login como Demo (nuevo JWT con subscription_status=expired/cancelled)
- [ ] Navegar a /orders (o cualquier ruta que no sea /subscription)
- [ ] Se muestra el muro "Cuenta suspendida", no el contenido de la ruta

## 15.2 /subscription sigue accesible bajo bloqueo
- [ ] Navegar a /subscription con la cuenta bloqueada
- [ ] La pagina de suscripcion se ve normal (no el muro)

## 15.3 Reactivar no desbloquea sin relogin
- [ ] Reactivar Demo via API (`POST /subscriptions/reactivate`)
- [ ] Sin cerrar sesion, navegar de nuevo a /orders
- [ ] Sigue mostrando el muro (JWT stale, comportamiento documentado y esperado)

## 15.4 Relogin desbloquea
- [ ] Cerrar sesion y volver a loguearse como Demo
- [ ] Navegar a /orders
- [ ] Ya no aparece el muro
