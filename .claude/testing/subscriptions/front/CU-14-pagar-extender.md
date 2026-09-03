# CU-14: Boton "Pagar / Extender" (frontend)

## Objetivo
Verificar que el boton abre el modal de compra correctamente para un negocio en
plan de catalogo y para uno en plan personalizado (fix de esta sesion), y que el
flujo de compra completa contra la API real.

## Precondiciones
- Frontend en http://localhost:3000, backend local con Demo (business 26) logueable
  con `demo@probability.com` / password de `.env.ai` (hash seteado en local para pruebas).
- Playwright MCP.

## 14.1 Login y navegar a /subscription
- [ ] Login exitoso, redirige al dashboard
- [ ] Navegar a /subscription
- [ ] La tarjeta de "Estado de suscripcion" muestra el plan actual (Basico)

## 14.2 Click en "Pagar / Extender" (plan de catalogo)
- [ ] El modal de compra se abre inmediatamente (sin recargar pagina)
- [ ] Muestra el plan actual preseleccionado y el precio correcto

## 14.3 Plan personalizado: asignar uno a Demo y repetir
- [ ] Via API (super admin), crear un custom plan payable para Demo y asignarselo
- [ ] Recargar /subscription, click en "Pagar / Extender"
- [ ] El modal abre igual (regresion del fix de hoy) mostrando el plan personalizado

## 14.4 Completar la compra desde el modal
- [ ] Con saldo suficiente, confirmar la compra desde el modal
- [ ] Mensaje de exito visible
- [ ] La tarjeta de estado se actualiza con la nueva fecha de vencimiento
