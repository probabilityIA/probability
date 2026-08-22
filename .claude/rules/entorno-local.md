# Entorno local primero - probar antes de tocar produccion

Regla de trabajo OBLIGATORIA. Todo cambio se prueba primero contra la copia
local de la base, y solo despues, si hace falta, contra produccion.

## Por que

`back/central/.env` apuntaba historicamente al **RDS de produccion**. Una prueba
local creaba ordenes, cotizaciones y guias REALES, y una guia le programa
recoleccion de verdad a la transportadora y descuenta del wallet. Probar era
caro y daba miedo, asi que se probaba poco.

Con la copia local eso desaparece: se puede romper, borrar y volver a empezar.

## Que hay

| Que | Donde |
|---|---|
| PostgreSQL local | `127.0.0.1:5434`, base `probability`, `postgres`/`postgres` |
| Tunel al RDS de produccion | `127.0.0.1:5433` (`./scripts/aws-tunnel.sh ensure`) |
| Clonar datos a local | `./scripts/local-db-clone.sh [business_id]` |
| Cambiar a que base apunta el backend | `./scripts/dev-db-switch.sh local\|prod\|status` |

**El puerto local es 5434, no 5433.** 5433 es del tunel a produccion, que es a
donde apunta el MCP `postgres-probability`. Antes chocaban y `postgres_local`
quedaba en bucle de reinicio.

## Flujo de trabajo

```bash
./scripts/aws-tunnel.sh ensure          # tunel a produccion (solo para clonar)
./scripts/local-db-clone.sh 26          # estructura completa + datos del business Demo
./scripts/dev-db-switch.sh local        # el backend apunta a la copia local
./scripts/dev-services.sh restart backend
```

1. **Se trabaja siempre en local.** Cambios de codigo, migraciones, pruebas
   manuales, E2E de `.claude/testing/`, mocks: todo contra `5434`.
2. **Se verifica en local** que el caso pasa: endpoint, dato en BD, pantalla.
3. **Solo si el caso lo exige** (integracion con proveedor externo que no se
   puede simular, dato que solo existe en produccion) se pasa a produccion, con
   `dev-db-switch.sh prod`, avisando al usuario y cancelando lo que se cree.
4. **Se vuelve a local** al terminar: `./scripts/dev-db-switch.sh local`.

`dev-db-switch.sh status` dice a que base apunta el backend. **Revisarlo antes
de cualquier prueba que escriba datos.**

## Que clona el script

- **Estructura completa** de produccion: las 157 tablas, indices y constraints.
  Asi las migraciones se prueban contra el esquema real.
- **Datos de un solo business** (por defecto 26, Demo): ~13.500 filas.
  - Tablas con `business_id` -> filtradas por ese negocio.
  - Tablas hijas sin `business_id` (`order_items`, `shipments`,
    `shipment_sync_logs`, `invoice_items`, ...) -> filtradas por su padre.
  - Catalogos globales (estados, tipos de integracion, permisos, geozonas) ->
    completos.
  - `business` y `user` -> solo el negocio y su personal. **No se traen datos de
    otros clientes.**
- Al final reajusta todas las secuencias, para que los IDs nuevos no choquen.

Detalles que ya costaron tiempo y estan resueltos en el script:

- Se copia con **lista explicita de columnas**, excluyendo las **generadas**
  (`shipments.carrier_key`): `COPY FROM` las rechaza.
- Se carga con `session_replication_role = replica`, asi el orden de las tablas
  no importa y no hay que resolver el grafo de FKs.
- `\copy` no se puede mezclar con SQL en un mismo `-c` de psql: van en `-c`
  separados.
- Los clientes de postgres del sistema son viejos (v14) para un servidor 17: el
  script usa `psql`/`pg_dump` dentro del contenedor `postgis/postgis:17-3.4`.

## Repetir el clonado

`local-db-clone.sh` es idempotente y destructivo en local: hace
`DROP SCHEMA public CASCADE` y vuelve a crear todo. Correrlo cuando:

- cambio el esquema en produccion (migracion nueva de otro),
- se ensucio la base local de tanto probar,
- hace falta un dato de produccion que no estaba.

Nunca escribe en produccion: solo hace `SELECT` y `pg_dump --schema-only`.

## Migraciones

Probar SIEMPRE en local antes de correrlas contra produccion:

```bash
./scripts/dev-db-switch.sh local
cd back/migration && go run cmd/main.go
```

Ver `back/migration/MIGRACIONES.md`. Solo correr la migracion puntual, no toda
la cadena del constructor.

## Como se combina con worktrees

Un worktree aisla el **codigo**, no la base. Todos los worktrees comparten el
mismo PostgreSQL local en `5434` y el mismo `back/central/.env` **no**: cada
worktree tiene su propia copia del `.env` (esta gitignored, hay que copiarlo al
crear el worktree).

En la practica:

- **Cambios que no tocan el esquema:** varios worktrees contra la misma base
  local sin problema.
- **Cambios que SI tocan el esquema** (migraciones): trabajar en un worktree a
  la vez, o volver a clonar (`local-db-clone.sh`) al cambiar de rama, porque la
  base queda con el esquema de la ultima migracion corrida.
- Al terminar un worktree con migraciones, correr `local-db-clone.sh` para
  volver al esquema de produccion antes de seguir en otro.

## Prohibido

- Probar contra produccion "porque es mas rapido" cuando el caso se puede
  reproducir en local.
- Dejar el `.env` apuntando a produccion despues de una prueba.
- Escribir en produccion desde el MCP de postgres: es de **solo lectura**
  (`.claude/rules/testing.md`).
- Clonar un business que no sea de prueba sin acordarlo: son datos de un cliente
  real en el disco local.
