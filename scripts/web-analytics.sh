#!/usr/bin/env bash
# Analitica del sitio publico desde la terminal: GA4 (clics dentro de la web) y
# Search Console (clics desde los resultados de busqueda de Google).
#
#   ./scripts/web-analytics.sh props            propiedades GA4 visibles
#   ./scripts/web-analytics.sh sites            sitios de Search Console
#   ./scripts/web-analytics.sh pages [dias]     paginas mas vistas (GA4)
#   ./scripts/web-analytics.sh sources [dias]   de donde llega el trafico (GA4)
#   ./scripts/web-analytics.sh events [dias]    eventos, incluye clics (GA4)
#   ./scripts/web-analytics.sh queries [dias]   consultas de busqueda (Search Console)
#   ./scripts/web-analytics.sh clicks [dias]    paginas con mas clics desde Google
#   ./scripts/web-analytics.sh token            imprime el token y sale
#
# Autenticacion: impersonacion de una cuenta de servicio, sin llaves JSON.
# La cuenta de servicio debe estar invitada DENTRO de GA4 y de Search Console;
# el permiso de Google Cloud por si solo no alcanza (devuelve 403).
set -euo pipefail

SA="${ANALYTICS_SA:-probability-analytics@probabilityia.iam.gserviceaccount.com}"
PROPERTY="${GA4_PROPERTY_ID:-548369470}"
SITE="${GSC_SITE_URL:-sc-domain:probabilityia.com.co}"
DAYS="${2:-28}"

export PATH="$HOME/google-cloud-sdk/bin:$PATH"
export CLOUDSDK_ACTIVE_CONFIG_NAME="${CLOUDSDK_ACTIVE_CONFIG_NAME:-probability}"

command -v gcloud >/dev/null || { echo "falta gcloud (~/google-cloud-sdk)"; exit 1; }
command -v jq >/dev/null || { echo "falta jq: sudo apt install jq"; exit 1; }

token() {
  local user_token
  user_token=$(gcloud auth print-access-token 2>/dev/null) || {
    echo "sin credenciales: corre 'gcloud auth login'" >&2; exit 1; }
  curl -s -X POST \
    "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/${SA}:generateAccessToken" \
    -H "Authorization: Bearer ${user_token}" -H "Content-Type: application/json" \
    -d '{"scope":["https://www.googleapis.com/auth/analytics.readonly","https://www.googleapis.com/auth/webmasters.readonly"],"lifetime":"3600s"}' \
    | jq -r '.accessToken // empty'
}

TOKEN=""
need_token() {
  [ -n "$TOKEN" ] && return
  TOKEN=$(token)
  [ -n "$TOKEN" ] || { echo "no se pudo generar el token de la cuenta de servicio" >&2; exit 1; }
}

resolve_property() {
  [ -n "$PROPERTY" ] && return
  PROPERTY=$(curl -s -H "Authorization: Bearer ${TOKEN}" \
    "https://analyticsadmin.googleapis.com/v1beta/accountSummaries" \
    | jq -r '[.accountSummaries[]?.propertySummaries[]?.property] | .[0] // empty' \
    | sed 's|properties/||')
  [ -n "$PROPERTY" ] || {
    echo "sin propiedades GA4 visibles. Invita a ${SA} en GA4 > Administrar > Gestion de accesos." >&2
    exit 1; }
}

ga_report() {
  local dimension="$1" metrics="$2" limit="${3:-20}"
  need_token; resolve_property
  local metrics_json first_metric
  metrics_json=$(echo "$metrics" | tr ',' '\n' | jq -R -s -c 'split("\n") | map(select(length>0)) | map({name: .})')
  first_metric=$(echo "$metrics" | cut -d, -f1)
  curl -s -X POST \
    "https://analyticsdata.googleapis.com/v1beta/properties/${PROPERTY}:runReport" \
    -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" \
    -d "{\"dateRanges\":[{\"startDate\":\"${DAYS}daysAgo\",\"endDate\":\"today\"}],
         \"dimensions\":[{\"name\":\"${dimension}\"}],
         \"metrics\":${metrics_json},
         \"limit\":${limit},
         \"orderBys\":[{\"desc\":true,\"metric\":{\"metricName\":\"${first_metric}\"}}]}"
}

gsc_report() {
  local dimension="$1" limit="${2:-20}"
  need_token
  local start end
  if date -u -d "${DAYS} days ago" +%F >/dev/null 2>&1; then
    start=$(date -u -d "${DAYS} days ago" +%F)
    end=$(date -u +%F)
  else
    start=$(date -u -v-"${DAYS}"d +%F)
    end=$(date -u +%F)
  fi
  curl -s -X POST \
    "https://searchconsole.googleapis.com/webmasters/v3/sites/$(printf %s "$SITE" | jq -sRr @uri)/searchAnalytics/query" \
    -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" \
    -d "{\"startDate\":\"${start}\",\"endDate\":\"${end}\",\"dimensions\":[\"${dimension}\"],\"rowLimit\":${limit}}"
}

show_ga() {
  jq -r '
    if .error then "ERROR: " + .error.message
    else (["DIMENSION","VALOR"] | @tsv),
         (.rows // [] | .[] | [.dimensionValues[0].value, .metricValues[0].value] | @tsv)
    end' | column -t -s "$(printf '\t')"
}

show_gsc() {
  jq -r '
    if .error then "ERROR: " + .error.message
    else (["CLAVE","CLICS","IMPRESIONES","CTR","POSICION"] | @tsv),
         (.rows // [] | .[] | [.keys[0], (.clicks|tostring), (.impressions|tostring),
            ((.ctr*100*100|round)/100|tostring) + "%", ((.position*10|round)/10|tostring)] | @tsv)
    end' | column -t -s "$(printf '\t')"
}

case "${1:-}" in
  token)   token ;;
  props)   need_token
           curl -s -H "Authorization: Bearer ${TOKEN}" \
             "https://analyticsadmin.googleapis.com/v1beta/accountSummaries" \
             | jq -r 'if .error then "ERROR: " + .error.message
                      else .accountSummaries[]?.propertySummaries[]?
                           | .property + "  " + .displayName end' ;;
  sites)   need_token
           curl -s -H "Authorization: Bearer ${TOKEN}" \
             "https://searchconsole.googleapis.com/webmasters/v3/sites" \
             | jq -r 'if .error then "ERROR: " + .error.message
                      else .siteEntry[]? | .siteUrl + "  " + .permissionLevel end' ;;
  pages)   ga_report pagePath screenPageViews,totalUsers | show_ga ;;
  sources) ga_report sessionSource sessions,totalUsers   | show_ga ;;
  events)  ga_report eventName eventCount                | show_ga ;;
  queries) gsc_report query | show_gsc ;;
  clicks)  gsc_report page  | show_gsc ;;
  *)       sed -n '2,17p' "$0" | sed 's/^# \{0,1\}//' ;;
esac
