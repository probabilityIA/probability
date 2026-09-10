**Ticket:** TKT-000085
**Negocio:** todos (produccion completa) - reportado por LaPerchaDel10
**Canal afectado:** creacion manual de ordenes + toda orden entrante de integraciones (Shopify, WooCommerce, MercadoLibre, etc.) via la cola `probability.orders.canonical`

## Resumen

Una migracion de base de datos se probo y corrio en local pero nunca se corrio
en produccion. El codigo desplegado ya esperaba columnas que la tabla real de
produccion no tenia, tumbando toda consulta a la tabla `client` durante ~11
horas y generando un bucle caliente de 112,356 reintentos.

## Sintoma

Usuario reporta "Error al crear orden" (mensaje generico, sin detalle) al
intentar crear una orden manual en el negocio LaPerchaDel10, cliente
`david rodriguez` (DNI `1025460955`), con el modal completamente lleno
(cliente, direccion, financiera, medio de pago = Efectivo).

## Diagnostico

1. Se reprodujo el mismo formulario, mismos datos, mismo negocio -> mismo
   error "Error al crear orden" en el modal.
2. Hipotesis descartada: medio de pago vacio. Se probo dejandolo vacio a
   proposito: produce un toast especifico ("Selecciona el medio de pago"),
   no el banner generico. No era esto.
3. `read_network_requests` del navegador (Chrome DevTools) mostro 4 peticiones
   `_rsc=...` con status 503 justo despues del submit. Hipotesis descartada:
   se penso que el refresco de RSC de Next.js estaba fallando y que la UI
   interpretaba ese fallo como error de creacion. El `access.log` de nginx en
   produccion (fuente de verdad) muestra las mismas peticiones exactas con
   status **200** — el 503 fue una lectura a mitad de vuelo del tool del
   navegador, no un error real.
4. Se leyeron los logs del backend en produccion
   (`docker logs central_reserve_prod_green`) en la ventana del intento:
   ```
   ERR Error processing message error="failed to map and save order:
   error processing customer: error searching client by email:
   ERROR: column client.address does not exist (SQLSTATE 42703)"
   function=func1 queue=probability.orders.canonical
   ```
   Causa encontrada.

## Causa raiz

El commit `feat(clientes): carga masiva por Excel/CSV y rediseno del modal de
alta` (2026-09-09) agrego las columnas `address`, `city`, `notes` a
`back/migration/shared/models/models.go` (`Client`) y actualizo el
`mapClientToEntity`/`modelToEntity`/`Create`/`Update` del repositorio de
`customers` para que el `SELECT` de **cualquier** consulta de cliente incluya
esas 3 columnas.

La migracion (`migrateClientAddressFields`, AutoMigrate) se corrio y se
verifico contra la base local (`127.0.0.1:5434`). El codigo se pusheo a
`main` y el deploy automatico lo llevo a produccion
(`central_reserve_prod_green`, activo desde `2026-09-10T02:03:35Z`). **Nadie
corrio la migracion contra el RDS de produccion.** El codigo desplegado
esperaba columnas que la tabla real no tenia.

## Impacto medido

- **112,356** ocurrencias del error en los logs del backend en la ventana de
  deteccion (ultimas 6 horas antes del fix).
- El consumidor de RabbitMQ trata cualquier error no clasificado como
  transitorio (`Nack(false, true)`, sin limite de reintentos ni backoff — ver
  `.claude/rules/colas-errores-permanentes.md`). Como el error era permanente
  (la columna nunca iba a aparecer sola), el mensaje se reencolaba y fallaba
  de inmediato, en bucle, durante ~11 horas.
- La cola `probability.orders.canonical` (procesa **todas** las ordenes
  entrantes de integraciones, no solo la creacion manual) acumulo **89
  mensajes** sin procesar en el momento del diagnostico.
- Afecto a **todos los negocios** de produccion, no solo LaPerchaDel10:
  cualquier creacion/busqueda de cliente (por email o DNI) tocaba el `SELECT`
  roto.

## Correccion

```bash
./scripts/aws-tunnel.sh ensure          # tunel SSM al RDS de produccion
# back/migration/.env apuntado temporalmente a 127.0.0.1:5433 (RDS real)
cd back/migration && go run cmd/main.go # ~4 min por la latencia del tunel
# back/migration/.env restaurado a la configuracion local inmediatamente despues
```

No se toco codigo: el codigo desplegado ya era correcto, solo faltaba el
`AutoMigrate` en la base real.

## Verificacion

1. `psql` directo contra el RDS de produccion: `client.address`,
   `client.city`, `client.notes` existen.
2. `docker logs central_reserve_prod_green --since 3m` post-fix: 0
   ocurrencias del error.
3. `rabbitmqctl list_queues`: `probability.orders.canonical` paso de 89 a 0
   mensajes sin intervencion manual — el backlog se autoproceso al quedar
   corregido el schema. No se perdieron ordenes, solo llegaron con horas de
   retraso.

## Pendientes

- Evaluar clasificar los errores de SQL de schema (columna/tabla inexistente)
  como no-retryables en el consumidor de ordenes, para que un futuro desajuste
  de migracion no vuelva a generar un bucle caliente. Ver
  `.claude/alerts/consumidor-muerto-por-canal-cerrado.md` y
  `.claude/rules/colas-errores-permanentes.md`.
- Agregar un paso explicito de "correr migracion contra produccion" al
  checklist de deploy (`.claude/rules/deploy.md`) cuando un cambio de backend
  incluya migracion de base de datos. Hoy el deploy de codigo es automatico
  (push a `main`) pero las migraciones son 100% manuales
  (`back/migration/MIGRACIONES.md`), y nada avisa si quedaron desincronizadas.
