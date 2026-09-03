# Suscripciones: cuentas deshabilitadas con mas acceso que un plan pagado, y planes personalizados que nadie podia pagar

Sesion de testing E2E completo del modulo de suscripciones (CU-01 a CU-15,
`.claude/testing/subscriptions/`) tras el fix del boton "Pagar / Extender" y los
cambios de dias de cortesia del mismo dia. Aparecieron dos bugs graves que no
tenian sintoma visible hasta que se probaron los caminos correctos.

## Sintoma

Ninguno reportado por un cliente: aparecieron al ejercitar sistematicamente
`HasModuleAccess` y el flujo de creacion de planes personalizados en local.

## Diagnostico

**1. Cuenta deshabilitada = mas acceso, no menos.** `HasModuleAccess`
(`has_module_access.go`) resuelve `GetBusinessCurrentSubscriptionTypeID`, que
devuelve `nil` tanto para "nunca tuvo plan" como para "esta expired/cancelled".
El codigo trataba ese `nil` siempre como "sin restriccion, todo abierto (salvo
storefront/tickets/delivery)". Al deshabilitar una cuenta de prueba, paso de 8
modulos (los de su plan Basico) a 10 de 13 — le sobraban `announcements` y
`notification_config` que ni su plan pagado incluia.

Se confirmo que esto no es enteramente nuevo: hay tests existentes
(`TestHasModuleAccess_SinPlanAsignado_TodoLoNoRestringidoEsVisible`) que
documentan el fallback abierto como intencional para negocios **activos** sin
plan asignado (en produccion hay 5 de esos: 4 de prueba y `business_id 38`,
"Alejandra", un negocio real que nunca tuvo plan y depende de este
comportamiento — el usuario confirmo que no importa si pierde ese acceso al
desplegar el fix). El bug real es que el mismo fallback tambien se aplicaba
a cuentas **no activas**.

**2. Ningun plan personalizado nuevo se podia autopagar.** `CreateCustomPlan`
nunca seteaba `Payable` en la entidad, y ademas el repo `CreateSubscriptionType`
tampoco lo mapeaba al INSERT — doble bug, cualquiera de los dos solo ya bastaba
para reproducirlo. Todo plan personalizado creado desde la API nacia con
`payable=false` y `POST /purchase` lo rechazaba con "este plan no se puede
comprar, se asigna automaticamente", sin importar quien lo intentara, ni
siquiera el negocio dueño. Los 3 planes personalizados que hoy existen en
produccion (ids 4, 5, 8, incluido el de "sin intermediarios", business_id 34)
tienen `payable=true` porque alguien los parcheo a mano en algun momento — el
codigo nunca lo hizo bien.

Esto es directamente relevante al fix del boton "Pagar/Extender" del mismo dia:
sin esta correccion, un negocio con plan personalizado podria abrir el modal
de pago (ya arreglado) pero la compra fallaria igual al confirmar.

**3. Misma falla, tercera variante.** Al purgar un custom plan de prueba
mientras el negocio seguia asignado a el, `business.subscription_type_id` quedo
apuntando a una fila borrada (soft delete). `GetSubscriptionType` devuelve `nil`
para un plan borrado, y `HasModuleAccess` tenia el mismo `if subType == nil {
return true, nil }` — la misma falla de nuevo, esta vez por plan huerfano en vez
de por plan nunca asignado.

## Correccion

- `HasModuleAccess` ahora usa un helper compartido `allowWhenNoPlan` para los
  dos casos de "sin plan" (nunca asignado, y plan borrado/inexistente): si el
  negocio esta `active`, se mantiene el fallback abierto historico; si no,
  cierra el acceso. Tests nuevos:
  `TestHasModuleAccess_SinPlanYCuentaDeshabilitada_TodoQuedaCerrado`,
  `TestHasModuleAccess_PlanBorradoYCuentaDeshabilitada_TodoQuedaCerrado`.
- `CreateCustomPlan` setea `Payable: true` explicitamente, y el repo
  `CreateSubscriptionType` lo mapea al modelo antes del INSERT. Verificado en
  vivo: plan nuevo creado, `payable=true` en BD, comprado con exito por el
  negocio dueño, y el check de ownership (otro negocio intentando comprarlo)
  sigue rechazando con "subscription type not found".

## Hallazgos menores (no bugs, dejados anotados)

- `PUT /custom-plans/{id}` es reemplazo total del objeto: si se omite `active`
  o `description` quedan sobreescritos a `false`/`""`. Semantica PUT estandar,
  pero es una trampa real para cualquier consumidor de la API que mande un
  body parcial.
- `disable` recibe `business_id` por query param; `reactivate` lo recibe en el
  body. Inconsistencia menor, sin romper nada.
- Borrar un custom plan mientras un negocio esta activo en el no lo reasigna a
  ningun otro plan (deja la referencia huerfana descrita arriba). No se corrigio
  por estar fuera de alcance; anotado en `.claude/testing/subscriptions/back/RESULTS.md`.

## Pendiente

- Ninguno de los fixes de esta entrada esta desplegado a produccion.
- Detalle completo de los 15 casos de uso corridos (backend y frontend) en
  `.claude/testing/subscriptions/back/RESULTS.md` y
  `.claude/testing/subscriptions/front/RESULTS.md`.
