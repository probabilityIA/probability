#!/usr/bin/env bash
# Clona la estructura completa de produccion y SOLO los datos de un business
# hacia la base local. Lectura de produccion, escritura local. Nunca al reves.
#
#   ./scripts/local-db-clone.sh            # clona el business por defecto (26, Demo)
#   ./scripts/local-db-clone.sh 26         # explicito
#   ./scripts/local-db-clone.sh --check    # solo verifica requisitos
#
# Ademas del personal del business, SIEMPRE trae los usuarios de scope platform
# (super admins) con sus user_roles: sin ellos no se puede entrar a local como
# super admin y no se pueden probar los flujos que piden business_id por query.
#
# Requiere el tunel SSM arriba: ./scripts/aws-tunnel.sh ensure

set -e

BUSINESS_ID="${1:-26}"
[ "$BUSINESS_ID" = "--check" ] && { CHECK_ONLY=1; BUSINESS_ID=26; }

PROD_HOST=127.0.0.1
PROD_PORT=5433
LOCAL_HOST=127.0.0.1
LOCAL_PORT=5434
LOCAL_DB=probability
LOCAL_USER=postgres
LOCAL_PASS=postgres
CLIENT_IMAGE=postgis/postgis:17-3.4

DOCKER_HOST_ADDR=127.0.0.1
[ "$(uname -s)" = "Darwin" ] && DOCKER_HOST_ADDR=host.docker.internal

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="$ROOT_DIR/back/central/.env"

RED=$'\033[0;31m'; GREEN=$'\033[0;32m'; YELLOW=$'\033[1;33m'; BLUE=$'\033[0;34m'; NC=$'\033[0m'
say()  { echo "${BLUE}[clone]${NC} $*"; }
ok()   { echo "${GREEN}[clone]${NC} $*"; }
warn() { echo "${YELLOW}[clone]${NC} $*"; }
die()  { echo "${RED}[clone]${NC} $*" >&2; exit 1; }

[ -f "$ENV_FILE" ] || die "no existe $ENV_FILE"
env_val() { grep -E "^$1=" "$ENV_FILE" | head -1 | cut -d= -f2- | sed -e 's/^"//' -e 's/"$//' -e "s/^'//" -e "s/'$//"; }
PROD_DB=$(env_val DB_NAME)
PROD_USER=$(env_val DB_USER)
PROD_PASS=$(env_val DB_PASS)
[ -n "$PROD_DB" ] && [ -n "$PROD_USER" ] || die "no pude leer DB_NAME/DB_USER de $ENV_FILE"

prod() { docker run --rm --network host -e PGPASSWORD="$PROD_PASS" "$CLIENT_IMAGE" "$@"; }
prod_in() { docker run --rm -i --network host -e PGPASSWORD="$PROD_PASS" "$CLIENT_IMAGE" "$@"; }
locl() { docker run --rm --network host -e PGPASSWORD="$LOCAL_PASS" "$CLIENT_IMAGE" "$@"; }
locl_in() { docker run --rm -i --network host -e PGPASSWORD="$LOCAL_PASS" "$CLIENT_IMAGE" "$@"; }

PSQL_PROD=(psql -h "$DOCKER_HOST_ADDR" -p "$PROD_PORT" -U "$PROD_USER" -d "$PROD_DB" -v ON_ERROR_STOP=1)
PSQL_LOCAL=(psql -h "$DOCKER_HOST_ADDR" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$LOCAL_DB" -v ON_ERROR_STOP=1)

say "verificando requisitos"
docker info >/dev/null 2>&1 || die "docker no responde"
prod "${PSQL_PROD[@]}" -tAc 'SELECT 1' >/dev/null 2>&1 || \
  die "no hay conexion a produccion en $PROD_HOST:$PROD_PORT. Levanta el tunel: ./scripts/aws-tunnel.sh ensure"
locl "${PSQL_LOCAL[@]}" -tAc 'SELECT 1' >/dev/null 2>&1 || \
  die "no hay postgres local en $LOCAL_HOST:$LOCAL_PORT. Levanta: cd infra/compose-local && docker-compose up -d postgres"

BUSINESS_NAME=$(prod "${PSQL_PROD[@]}" -tAc "SELECT name FROM business WHERE id=$BUSINESS_ID" | tr -d '\r')
[ -n "$BUSINESS_NAME" ] || die "el business $BUSINESS_ID no existe en produccion"
ok "produccion OK, local OK, business $BUSINESS_ID = $BUSINESS_NAME"
[ -n "${CHECK_ONLY:-}" ] && exit 0

# ---------------------------------------------------------------------------
# Filtros para tablas que NO tienen business_id y cuelgan de una que si
# ---------------------------------------------------------------------------
ORDERS_SUB="SELECT id FROM orders WHERE business_id=$BUSINESS_ID"

declare -A CHILD_FILTER=(
  [order_items]="order_id IN ($ORDERS_SUB)"
  [order_channel_metadata]="order_id IN ($ORDERS_SUB)"
  [order_history]="order_id IN ($ORDERS_SUB)"
  [addresses]="order_id IN ($ORDERS_SUB)"
  [payments]="order_id IN ($ORDERS_SUB)"
  [route_stop]="order_id IN ($ORDERS_SUB)"
  [shipments]="order_id IN ($ORDERS_SUB)"
  [shipment_sync_logs]="shipment_id IN (SELECT id FROM shipments WHERE order_id IN ($ORDERS_SUB))"
  [invoice_items]="invoice_id IN (SELECT id FROM invoices WHERE business_id=$BUSINESS_ID)"
  [invoice_sync_logs]="invoice_id IN (SELECT id FROM invoices WHERE business_id=$BUSINESS_ID)"
  [bulk_invoice_job_items]="job_id IN (SELECT id FROM bulk_invoice_jobs WHERE business_id=$BUSINESS_ID)"
  [accounting_invoice_items]="invoice_id IN (SELECT id FROM accounting_invoices WHERE business_id=$BUSINESS_ID)"
  [invoicing_config_integrations]="config_id IN (SELECT id FROM invoicing_configs WHERE business_id=$BUSINESS_ID)"
  [integration_notification_configs]="integration_id IN (SELECT id FROM integrations WHERE business_id=$BUSINESS_ID)"
  [integration_sync_run_items]="run_id IN (SELECT id FROM integration_sync_runs WHERE business_id=$BUSINESS_ID)"
  [warehouse_locations]="warehouse_id IN (SELECT id FROM warehouses WHERE business_id=$BUSINESS_ID)"
  [whatsapp_message_logs]="conversation_id IN (SELECT id FROM whatsapp_conversations WHERE business_id=$BUSINESS_ID)"
  [ticket_status_history]="ticket_id IN (SELECT id FROM tickets WHERE business_id=$BUSINESS_ID)"
  [payment_sync_log]="payment_transaction_id IN (SELECT id FROM payment_transaction WHERE business_id=$BUSINESS_ID)"
  [business]="id=$BUSINESS_ID"
  ["user"]="id IN (SELECT user_id FROM business_staff WHERE business_id=$BUSINESS_ID) OR scope_id IN (SELECT id FROM scope WHERE code='platform')"
  [user_roles]="user_id IN (SELECT user_id FROM business_staff WHERE business_id=$BUSINESS_ID) OR user_id IN (SELECT id FROM \"user\" WHERE scope_id IN (SELECT id FROM scope WHERE code='platform'))"
)

# Datos de otros negocios o ruido que no aporta en local
SKIP_TABLES="commercial_prospects payment_webhook_events bold_webhook_events spatial_ref_sys"

# ---------------------------------------------------------------------------
say "1/4 estructura (schema-only) desde produccion"
SCHEMA_SQL=$(mktemp)
trap 'rm -f "$SCHEMA_SQL"' EXIT
prod pg_dump -h "$DOCKER_HOST_ADDR" -p "$PROD_PORT" -U "$PROD_USER" -d "$PROD_DB" \
  --schema-only --no-owner --no-privileges --no-comments \
  --exclude-schema='tiger*' --exclude-schema=topology > "$SCHEMA_SQL"
ok "estructura obtenida ($(wc -l < "$SCHEMA_SQL") lineas)"

say "2/4 recreando schema public en local"
locl "${PSQL_LOCAL[@]}" -q -c 'DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;' >/dev/null
locl "${PSQL_LOCAL[@]}" -q -c 'CREATE EXTENSION IF NOT EXISTS postgis; CREATE EXTENSION IF NOT EXISTS "uuid-ossp";' >/dev/null 2>&1 || true
locl_in "${PSQL_LOCAL[@]}" -q -f - < "$SCHEMA_SQL" > /dev/null 2>&1 || \
  warn "la carga del schema reporto avisos (normal por extensiones), se continua"
TABLAS=$(locl "${PSQL_LOCAL[@]}" -tAc "SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'")
ok "schema local con $TABLAS tablas"

say "3/4 copiando datos del business $BUSINESS_ID"
TABLE_LIST=$(prod "${PSQL_PROD[@]}" -tAc "
  SELECT c.relname
  FROM pg_class c JOIN pg_namespace n ON n.oid=c.relnamespace
  WHERE n.nspname='public' AND c.relkind='r'
  ORDER BY c.relname" | tr -d '\r')

BUSINESS_TABLES=$(prod "${PSQL_PROD[@]}" -tAc "
  SELECT table_name FROM information_schema.columns
  WHERE table_schema='public' AND column_name='business_id'" | tr -d '\r')

total_filas=0; copiadas=0; vacias=0; saltadas=0
for t in $TABLE_LIST; do
  case " $SKIP_TABLES " in *" $t "*) saltadas=$((saltadas+1)); continue;; esac

  if [ -n "${CHILD_FILTER[$t]:-}" ]; then
    where="${CHILD_FILTER[$t]}"
  elif echo "$BUSINESS_TABLES" | grep -qx "$t"; then
    where="business_id=$BUSINESS_ID OR business_id IS NULL"
  else
    where="true"
  fi

  n=$(prod "${PSQL_PROD[@]}" -tAc "SELECT count(*) FROM \"$t\" WHERE $where" 2>/dev/null | tr -d '\r' || echo 0)
  if [ "${n:-0}" = "0" ]; then vacias=$((vacias+1)); continue; fi

  # Lista explicita de columnas, sin las generadas: COPY FROM las rechaza
  cols=$(prod "${PSQL_PROD[@]}" -tAc "
    SELECT string_agg(quote_ident(a.attname), ',' ORDER BY a.attnum)
    FROM pg_attribute a
    WHERE a.attrelid='public.\"$t\"'::regclass
      AND a.attnum > 0 AND NOT a.attisdropped AND a.attgenerated = ''" | tr -d '\r')
  [ -n "$cols" ] || { warn "  $t: sin columnas copiables, se omite"; continue; }

  prod "${PSQL_PROD[@]}" -c "\copy (SELECT $cols FROM \"$t\" WHERE $where) TO STDOUT WITH (FORMAT csv)" 2>/dev/null | \
    locl_in "${PSQL_LOCAL[@]}" -q -c "SET session_replication_role = replica" \
      -c "\copy \"$t\" ($cols) FROM STDIN WITH (FORMAT csv)" >/dev/null 2>&1 || {
      warn "  $t: fallo la copia, se omite"; continue; }

  printf '  %-42s %8s filas\n' "$t" "$n"
  total_filas=$((total_filas + n)); copiadas=$((copiadas+1))
done

say "4/4 reajustando secuencias"
locl "${PSQL_LOCAL[@]}" -q -c "
DO \$\$
DECLARE r record; maxid bigint;
BEGIN
  FOR r IN
    SELECT s.relname AS seq, t.relname AS tbl, a.attname AS col
    FROM pg_class s
    JOIN pg_depend d ON d.objid=s.oid AND d.classid='pg_class'::regclass AND d.deptype='a'
    JOIN pg_class t ON t.oid=d.refobjid
    JOIN pg_attribute a ON a.attrelid=t.oid AND a.attnum=d.refobjsubid
    JOIN pg_namespace n ON n.oid=s.relnamespace
    WHERE s.relkind='S' AND n.nspname='public'
  LOOP
    EXECUTE format('SELECT COALESCE(MAX(%I),0) FROM %I', r.col, r.tbl) INTO maxid;
    EXECUTE format('SELECT setval(%L, GREATEST(%s,1), %L::boolean)', r.seq, maxid, maxid > 0);
  END LOOP;
END \$\$;" >/dev/null

echo
ok "listo: $copiadas tablas con datos, $total_filas filas, $vacias vacias, $saltadas saltadas"
echo
echo "Para que el backend use esta base, en back/central/.env:"
echo "  DB_HOST=127.0.0.1"
echo "  DB_PORT=$LOCAL_PORT"
echo "  DB_USER=$LOCAL_USER"
echo "  DB_PASS=$LOCAL_PASS"
echo "  DB_NAME=$LOCAL_DB"
echo "  REDIS_PASSWORD=localdev"
