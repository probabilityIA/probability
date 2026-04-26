#!/usr/bin/env bash
# Runner para CU-08..CU-17 del modulo COD.
# Uso: ./run-cod-tests.sh
# Requiere: backend en :3050, mock en :9091, jq, curl, python3.

set -uo pipefail

BASE="http://localhost:3050/api/v1"
DEMO_EMAIL="demo@probability.com"
DEMO_PASS="ProbabilityDemo"
SUPER_EMAIL="${AI_SUPER_ADMIN_EMAIL}"
SUPER_PASS="Seb51923662#"

PASS=0
FAIL=0
RESULTS=()

c_red()   { printf "\033[0;31m%s\033[0m" "$1"; }
c_green() { printf "\033[0;32m%s\033[0m" "$1"; }
c_blue()  { printf "\033[0;34m%s\033[0m" "$1"; }

ok()   { PASS=$((PASS+1)); RESULTS+=("OK   - $1"); echo "  $(c_green '[OK]') $1"; }
fail() { FAIL=$((FAIL+1)); RESULTS+=("FAIL - $1"); echo "  $(c_red '[FAIL]') $1"; }
section() { echo; echo "$(c_blue "=== $1 ===")"; }

# JSON helpers
jq_or() { python3 -c "import json,sys;d=json.load(sys.stdin);
keys='$1'.split('.')
for k in keys:
    if isinstance(d,dict): d=d.get(k)
    else: d=None
    if d is None: break
print(d if d is not None else '$2')"
}

login() {
  local email="$1" pass="$2"
  curl -s -X POST "$BASE/auth/login" \
    -H 'Content-Type: application/json' -H 'X-Client-Type: api' \
    -d "{\"email\":\"$email\",\"password\":\"$pass\"}" \
    | python3 -c "import json,sys;print(json.load(sys.stdin)['data']['token'])"
}

section "Setup: login demo + super admin"
DEMO_TOKEN=$(login "$DEMO_EMAIL" "$DEMO_PASS")
SUPER_TOKEN=$(login "$SUPER_EMAIL" "$SUPER_PASS")
[ -n "$DEMO_TOKEN" ] && [ "${#DEMO_TOKEN}" -gt 50 ] && ok "demo token obtenido (${#DEMO_TOKEN} chars)" || { fail "demo token vacio"; exit 1; }
[ -n "$SUPER_TOKEN" ] && [ "${#SUPER_TOKEN}" -gt 50 ] && ok "super admin token obtenido" || fail "super admin token vacio"

# ============================================================
section "CU-08: Listado y filtros"
RESP=$(curl -s -H "Authorization: Bearer $DEMO_TOKEN" "$BASE/shipments/cod?page=1&page_size=10")
echo "$RESP" | python3 -c "import json,sys;d=json.load(sys.stdin);assert d.get('success') and isinstance(d.get('data'),list), d" 2>/dev/null \
  && ok "8.1 listado base 200 con data array" || fail "8.1 listado base"

RESP=$(curl -s -H "Authorization: Bearer $DEMO_TOKEN" "$BASE/shipments/cod?status=delivered&page_size=20")
ALL_DELIVERED=$(echo "$RESP" | python3 -c "import json,sys;d=json.load(sys.stdin)['data'];print(all(x.get('status')=='delivered' for x in d))" 2>/dev/null)
[ "$ALL_DELIVERED" = "True" ] && ok "8.2 filtro status=delivered" || fail "8.2 filtro status (got $ALL_DELIVERED)"

RESP=$(curl -s -H "Authorization: Bearer $DEMO_TOKEN" "$BASE/shipments/cod?is_paid=false&page_size=50")
ALL_UNPAID=$(echo "$RESP" | python3 -c "import json,sys;d=json.load(sys.stdin)['data'];print(all(not x.get('is_paid',False) for x in d))" 2>/dev/null)
[ "$ALL_UNPAID" = "True" ] && ok "8.3 filtro is_paid=false" || fail "8.3 filtro is_paid=false"

# Super admin sin business_id -> 400
HTTP=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $SUPER_TOKEN" "$BASE/shipments/cod")
[ "$HTTP" = "400" ] && ok "8.5 super admin sin business_id -> 400" || fail "8.5 super admin sin business_id (HTTP=$HTTP)"

HTTP=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $SUPER_TOKEN" "$BASE/shipments/cod?business_id=26&page_size=5")
[ "$HTTP" = "200" ] && ok "8.5 super admin con business_id=26 -> 200" || fail "8.5 super admin con business_id (HTTP=$HTTP)"

# ============================================================
section "CU-09: Crear 5 ordenes (4 COD + 1 sin COD)"
declare -A ORDER_ID ORDER_NUMBER
TS_NS=$(date +%s%N)

create_order() {
  local n="$1" amount="$2" cod="$3"
  local body
  if [ "$cod" = "0" ]; then
    body=$(printf '{"business_id":26,"integration_id":35,"integration_type":"platform","platform":"manual","external_id":"cod-test-%s-%s","subtotal":%s,"total_amount":%s,"currency":"COP","customer_name":"COD Test %s","customer_email":"codtest%s@test.com","customer_phone":"30010000%s","customer_dni":"100000000%s","status":"pending","payment_method_id":1,"occurred_at":"2026-04-25T12:00:00Z","imported_at":"2026-04-25T12:00:00Z","order_items":[{"product_sku":"SKU-COD-%s","product_name":"P %s","quantity":1,"unit_price":%s,"total_price":%s,"currency":"COP"}]}' "$TS_NS" "$n" "$amount" "$amount" "$n" "$n" "$n" "$n" "$n" "$n" "$amount" "$amount")
  else
    body=$(printf '{"business_id":26,"integration_id":35,"integration_type":"platform","platform":"manual","external_id":"cod-test-%s-%s","subtotal":%s,"total_amount":%s,"cod_total":%s,"currency":"COP","customer_name":"COD Test %s","customer_email":"codtest%s@test.com","customer_phone":"30010000%s","customer_dni":"100000000%s","status":"pending","payment_method_id":1,"occurred_at":"2026-04-25T12:00:00Z","imported_at":"2026-04-25T12:00:00Z","order_items":[{"product_sku":"SKU-COD-%s","product_name":"P %s","quantity":1,"unit_price":%s,"total_price":%s,"currency":"COP"}]}' "$TS_NS" "$n" "$amount" "$amount" "$cod" "$n" "$n" "$n" "$n" "$n" "$n" "$amount" "$amount")
  fi
  local resp=$(curl -s -X POST -H "Authorization: Bearer $DEMO_TOKEN" -H 'Content-Type: application/json' -d "$body" "$BASE/orders")
  local oid=$(echo "$resp" | python3 -c "import json,sys;print(json.load(sys.stdin).get('data',{}).get('ID',''))" 2>/dev/null)
  local onum=$(echo "$resp" | python3 -c "import json,sys;print(json.load(sys.stdin).get('data',{}).get('OrderNumber',''))" 2>/dev/null)
  if [ -n "$oid" ]; then
    ORDER_ID[$n]="$oid"
    ORDER_NUMBER[$n]="$onum"
    ok "9.$n creada orden N=$n cod=$cod -> $oid ($onum)"
  else
    fail "9.$n no se pudo crear orden N=$n: $(echo "$resp" | head -c 200)"
  fi
}

create_order 1 50000   50000
create_order 2 1500000 1500000
create_order 3 250000  250000
create_order 4 320000  320000
create_order 5 99000   0

# ============================================================
section "CU-10: Generar guias COD via mock (N=1..4)"
declare -A SHIPMENT_ID TRACKING

generate_guide() {
  local n="$1" cod="${2}"
  local oid="${ORDER_ID[$n]}" onum="${ORDER_NUMBER[$n]}"
  [ -z "$oid" ] && { fail "10.$n no order_id"; return; }

  local body=$(python3 -c "
import json
print(json.dumps({
  'idRate': 1, 'carrier': 'COORDINADORA',
  'order_uuid': '$oid', 'myShipmentReference': 'COD $onum',
  'external_order_id': '$onum',
  'requestPickup': False, 'pickupDate': '2026-04-26',
  'insurance': False, 'description': 'COD shipment',
  'contentValue': $cod, 'codValue': $cod, 'codPaymentMethod': 'cash',
  'includeGuideCost': False, 'totalCost': 15000,
  'packages': [{'weight':1,'height':10,'width':15,'length':20}],
  'origin': {'daneCode':'11001000','address':'Calle 100 #15-20','company':'Demo','firstName':'Demo','lastName':'Test','email':'demo@probability.com','phone':'3001234567','suburb':'Usaquen','crossStreet':'Calle 100','reference':'Of 301'},
  'destination': {'daneCode':'08001000','address':'Carrera 54 #72-80','company':'Cliente $n','firstName':'Cli','lastName':'$n','email':'c$n@test.com','phone':'30010000$n','suburb':'Norte','crossStreet':'Carrera 54','reference':'Apto 1'}
}))")
  local resp=$(curl -s -X POST -H "Authorization: Bearer $DEMO_TOKEN" -H 'Content-Type: application/json' -d "$body" "$BASE/shipments/generate")
  local sid=$(echo "$resp" | python3 -c "import json,sys;d=json.load(sys.stdin);print(d.get('shipment_id') or d.get('data',{}).get('shipment_id') or '')" 2>/dev/null)
  if [ -n "$sid" ] && [ "$sid" != "0" ] && [ "$sid" != "None" ]; then
    SHIPMENT_ID[$n]="$sid"
    ok "10.$n shipment_id=$sid generado"
  else
    fail "10.$n generar guia fallo: $(echo "$resp" | head -c 250)"
  fi
}

generate_guide 1 50000
generate_guide 2 1500000
generate_guide 3 250000
generate_guide 4 320000

echo "  ... esperando 6s para que el mock procese y el consumer cree los shipments en DB"
sleep 6

# Verificar que los shipments tienen tracking_number
for n in 1 2 3 4; do
  sid="${SHIPMENT_ID[$n]:-}"
  [ -z "$sid" ] && continue
  resp=$(curl -s -H "Authorization: Bearer $DEMO_TOKEN" "$BASE/shipments/$sid")
  trk=$(echo "$resp" | python3 -c "import json,sys;print(json.load(sys.stdin).get('data',{}).get('tracking_number') or '')" 2>/dev/null)
  if [ -n "$trk" ]; then
    TRACKING[$n]="$trk"
    ok "10.$n shipment $sid tracking=$trk"
  else
    fail "10.$n shipment $sid sin tracking_number todavia"
  fi
done

# ============================================================
section "CU-11: Marcar shipments 1, 2, 4 como delivered"

mark_delivered() {
  local n="$1"
  local sid="${SHIPMENT_ID[$n]:-}"
  [ -z "$sid" ] && { fail "11.$n no shipment id"; return; }
  resp=$(curl -s -X PUT -H "Authorization: Bearer $DEMO_TOKEN" -H 'Content-Type: application/json' \
    -d '{"status":"delivered","delivered_at":"2026-04-25T18:00:00Z"}' \
    "$BASE/shipments/$sid")
  st=$(echo "$resp" | python3 -c "import json,sys;print(json.load(sys.stdin).get('data',{}).get('status',''))" 2>/dev/null)
  [ "$st" = "delivered" ] && ok "11.$n shipment $sid -> delivered" || fail "11.$n PUT status=$st"
}

mark_delivered 1
mark_delivered 2
mark_delivered 4

# ============================================================
section "CU-12: Cobro 50.000 (happy path)"
sid="${SHIPMENT_ID[1]:-}"
oid="${ORDER_ID[1]:-}"
resp=$(curl -s -X POST -H "Authorization: Bearer $DEMO_TOKEN" -H 'Content-Type: application/json' \
  -d '{"notes":"CU-12 efectivo"}' "$BASE/shipments/$sid/collect-cod")
ip=$(echo "$resp" | python3 -c "import json,sys;d=json.load(sys.stdin);print(d.get('data',{}).get('is_paid'))" 2>/dev/null)
[ "$ip" = "True" ] && ok "12 cobro 50k OK is_paid=true" || fail "12 cobro 50k: $(echo "$resp" | head -c 200)"

# ============================================================
section "CU-13: Cobro 1.500.000 (happy path)"
sid="${SHIPMENT_ID[2]:-}"
oid="${ORDER_ID[2]:-}"
resp=$(curl -s -X POST -H "Authorization: Bearer $DEMO_TOKEN" -H 'Content-Type: application/json' \
  -d '{"notes":"CU-13 INV-2026-001"}' "$BASE/shipments/$sid/collect-cod")
ct=$(echo "$resp" | python3 -c "import json,sys;d=json.load(sys.stdin);print(d.get('data',{}).get('cod_total'))" 2>/dev/null)
[ "$ct" = "1500000" ] || [ "$ct" = "1500000.0" ] && ok "13 cobro 1.5M cod_total=$ct" || fail "13 cobro 1.5M: $(echo "$resp" | head -c 200)"

# ============================================================
section "CU-14: Cobrar shipment 3 (pending) -> ErrShipmentNotDelivered"
sid="${SHIPMENT_ID[3]:-}"
HTTP=$(curl -s -o /tmp/cu14.json -w "%{http_code}" -X POST -H "Authorization: Bearer $DEMO_TOKEN" -H 'Content-Type: application/json' \
  -d '{"notes":"no debe permitir"}' "$BASE/shipments/$sid/collect-cod")
MSG=$(python3 -c "import json;print(json.load(open('/tmp/cu14.json')).get('message',''))" 2>/dev/null)
[ "$HTTP" = "400" ] && [[ "$MSG" == *delivered* ]] && ok "14 rechazo 400 con mensaje delivered" || fail "14 rechazo no-delivered (HTTP=$HTTP MSG=$MSG)"

# ============================================================
section "CU-15: Doble cobro shipment 4"
sid="${SHIPMENT_ID[4]:-}"
oid="${ORDER_ID[4]:-}"
resp1=$(curl -s -X POST -H "Authorization: Bearer $DEMO_TOKEN" -H 'Content-Type: application/json' \
  -d '{"notes":"primer CU-15"}' "$BASE/shipments/$sid/collect-cod")
ip=$(echo "$resp1" | python3 -c "import json,sys;d=json.load(sys.stdin);print(d.get('data',{}).get('is_paid'))" 2>/dev/null)
[ "$ip" = "True" ] && ok "15.1 primer cobro OK" || fail "15.1 primer cobro: $(echo "$resp1" | head -c 200)"

HTTP=$(curl -s -o /tmp/cu15.json -w "%{http_code}" -X POST -H "Authorization: Bearer $DEMO_TOKEN" -H 'Content-Type: application/json' \
  -d '{"notes":"doble cobro"}' "$BASE/shipments/$sid/collect-cod")
MSG=$(python3 -c "import json;print(json.load(open('/tmp/cu15.json')).get('message',''))" 2>/dev/null)
[ "$HTTP" = "400" ] && [[ "$MSG" == *already*paid* ]] && ok "15.2 segundo cobro 400 already paid" || fail "15.2 doble cobro (HTTP=$HTTP MSG=$MSG)"

# ============================================================
section "CU-16: Orden no COD (N=5)"
oid5="${ORDER_ID[5]:-}"
[ -n "$oid5" ] || { fail "16 no order_id 5"; }
TS=$(date +%s)
resp=$(curl -s -X POST -H "Authorization: Bearer $DEMO_TOKEN" -H 'Content-Type: application/json' \
  -d "{\"order_id\":\"$oid5\",\"client_name\":\"No COD\",\"destination_address\":\"Calle 99 #1-2\",\"tracking_number\":\"TEST-NOCOD-$TS\",\"carrier\":\"TEST\",\"status\":\"delivered\",\"delivered_at\":\"2026-04-25T18:00:00Z\"}" \
  "$BASE/shipments")
sid5=$(echo "$resp" | python3 -c "import json,sys;d=json.load(sys.stdin);print(d.get('data',{}).get('id') or '')" 2>/dev/null)
if [ -n "$sid5" ]; then
  ok "16.1 shipment manual creado id=$sid5"
  HTTP=$(curl -s -o /tmp/cu16.json -w "%{http_code}" -X POST -H "Authorization: Bearer $DEMO_TOKEN" -H 'Content-Type: application/json' \
    -d '{"notes":"no debe permitir"}' "$BASE/shipments/$sid5/collect-cod")
  MSG=$(python3 -c "import json;print(json.load(open('/tmp/cu16.json')).get('message',''))" 2>/dev/null)
  [ "$HTTP" = "400" ] && [[ "$MSG" == *cash*on*delivery* || "$MSG" == *not*COD* || "$MSG" == *cash* ]] && ok "16.2 rechazo 400 not COD" || fail "16.2 (HTTP=$HTTP MSG=$MSG)"
else
  fail "16.1 crear shipment manual: $(echo "$resp" | head -c 200)"
fi

# ============================================================
section "CU-17: Listado filtrado y suma"
RESP=$(curl -s -H "Authorization: Bearer $DEMO_TOKEN" "$BASE/shipments/cod?is_paid=true&page_size=100")
COUNT_PAID=$(echo "$RESP" | python3 -c "import json,sys;d=json.load(sys.stdin);ids={int(s['id']) for s in d['data']};want={int('${SHIPMENT_ID[1]:-0}'), int('${SHIPMENT_ID[2]:-0}'), int('${SHIPMENT_ID[4]:-0}')};print(len(want.intersection(ids)))" 2>/dev/null)
[ "$COUNT_PAID" = "3" ] && ok "17.1 lista is_paid=true incluye SHIPMENTS 1,2,4" || fail "17.1 incluidos=$COUNT_PAID/3"

RESP=$(curl -s -H "Authorization: Bearer $DEMO_TOKEN" "$BASE/shipments/cod?is_paid=false&page_size=100")
INCLUDE_PAID=$(echo "$RESP" | python3 -c "import json,sys;d=json.load(sys.stdin);ids={int(s['id']) for s in d['data']};paid={int('${SHIPMENT_ID[1]:-0}'), int('${SHIPMENT_ID[2]:-0}'), int('${SHIPMENT_ID[4]:-0}')};print(len(paid.intersection(ids)))" 2>/dev/null)
[ "$INCLUDE_PAID" = "0" ] && ok "17.2 lista is_paid=false NO incluye SHIPMENTS pagados" || fail "17.2 incluye pagados (debio 0): $INCLUDE_PAID"

# ============================================================
section "Resumen"
echo "  PASSED: $(c_green $PASS)"
echo "  FAILED: $(c_red $FAIL)"

# Persist results
{
  echo "# Resultados COD ($(date +'%Y-%m-%d %H:%M:%S'))"
  echo
  echo "PASS: $PASS  FAIL: $FAIL"
  echo
  for r in "${RESULTS[@]}"; do
    echo "- $r"
  done
  echo
  echo "## IDs generados"
  for n in 1 2 3 4 5; do
    echo "- N=$n  order=${ORDER_ID[$n]:-}  number=${ORDER_NUMBER[$n]:-}  shipment=${SHIPMENT_ID[$n]:-}  tracking=${TRACKING[$n]:-}"
  done
} > "$(dirname "$0")/RESULTS.md"

[ "$FAIL" -eq 0 ] && exit 0 || exit 1
