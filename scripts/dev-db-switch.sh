#!/usr/bin/env bash
# Cambia a que base apunta el backend local: la copia local del business Demo,
# o el RDS de produccion por el tunel SSM.
#
#   ./scripts/dev-db-switch.sh local    # BD local (por defecto para trabajar)
#   ./scripts/dev-db-switch.sh prod     # RDS de produccion (requiere tunel)
#   ./scripts/dev-db-switch.sh status
#
# La primera vez que se cambia a local, las credenciales de produccion se
# guardan en back/central/.env.dbprod (gitignored) para poder volver.

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="$ROOT_DIR/back/central/.env"
PROD_BAK="$ROOT_DIR/back/central/.env.dbprod"

GREEN=$'\033[0;32m'; YELLOW=$'\033[1;33m'; BLUE=$'\033[0;34m'; RED=$'\033[0;31m'; NC=$'\033[0m'
say()  { echo "${BLUE}[db]${NC} $*"; }
ok()   { echo "${GREEN}[db]${NC} $*"; }
die()  { echo "${RED}[db]${NC} $*" >&2; exit 1; }

[ -f "$ENV_FILE" ] || die "no existe $ENV_FILE"
val() { grep -E "^$1=" "$ENV_FILE" | head -1 | cut -d= -f2-; }
set_val() { local k="$1" v="$2"; if grep -qE "^$k=" "$ENV_FILE"; then
    python3 - "$ENV_FILE" "$k" "$v" <<'PY'
import re,sys
p,k,v=sys.argv[1],sys.argv[2],sys.argv[3]
s=open(p).read()
open(p,'w').write(re.sub(rf'^{re.escape(k)}=.*$', f'{k}={v}', s, count=1, flags=re.M))
PY
  else printf '\n%s=%s\n' "$k" "$v" >> "$ENV_FILE"; fi; }

case "${1:-status}" in
  status)
    h=$(val DB_HOST); p=$(val DB_PORT); n=$(val DB_NAME)
    if [ "$h" = "127.0.0.1" ] && [ "$p" = "5434" ]; then
      ok "LOCAL  -> $h:$p/$n"
    else
      echo "${YELLOW}[db]${NC} PRODUCCION -> $h:$p/$n"
    fi
    ;;

  local)
    if [ ! -f "$PROD_BAK" ]; then
      { echo "DB_HOST=$(val DB_HOST)"; echo "DB_PORT=$(val DB_PORT)"
        echo "DB_USER=$(val DB_USER)"; echo "DB_PASS=$(val DB_PASS)"
        echo "DB_NAME=$(val DB_NAME)"; echo "PGSSLMODE=$(val PGSSLMODE)"
        echo "REDIS_PASSWORD=$(val REDIS_PASSWORD)"; } > "$PROD_BAK"
      chmod 600 "$PROD_BAK"
      say "credenciales de produccion guardadas en back/central/.env.dbprod"
    fi
    set_val DB_HOST 127.0.0.1; set_val DB_PORT 5434
    set_val DB_USER postgres;  set_val DB_PASS postgres
    set_val DB_NAME probability; set_val PGSSLMODE disable
    set_val REDIS_PASSWORD localdev
    ok "apuntando a la BD LOCAL (127.0.0.1:5434/probability)"
    echo "   reinicia el backend: ./scripts/dev-services.sh restart backend"
    ;;

  prod)
    [ -f "$PROD_BAK" ] || die "no hay $PROD_BAK; restaura DB_HOST/DB_PORT/DB_USER/DB_PASS/DB_NAME a mano"
    while IFS='=' read -r k v; do [ -n "$k" ] && set_val "$k" "$v"; done < "$PROD_BAK"
    echo "${YELLOW}[db]${NC} apuntando a PRODUCCION. Requiere el tunel: ./scripts/aws-tunnel.sh ensure"
    echo "   reinicia el backend: ./scripts/dev-services.sh restart backend"
    ;;

  *) die "uso: $0 [local|prod|status]" ;;
esac
