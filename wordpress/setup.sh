#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

WP_URL="http://localhost:8088"
ADMIN_USER="admin"
ADMIN_PASS="admin"
ADMIN_EMAIL="admin@probability.test"
WEBHOOK_URL="http://host.docker.internal:3050/api/v1/woocommerce/webhook"

if docker compose version >/dev/null 2>&1; then
  DC="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  DC="docker-compose"
else
  echo "ERROR: no se encontro 'docker compose' ni 'docker-compose'." >&2
  exit 1
fi

cli() { $DC run --rm -T wpcli "$@"; }

echo "==> Usando: ${DC}"
echo "==> Levantando contenedores (db + wordpress)..."
$DC up -d db wordpress

echo "==> Esperando a que WordPress responda en ${WP_URL} ..."
for _ in $(seq 1 60); do
  if curl -fsS -o /dev/null "${WP_URL}"; then break; fi
  sleep 2
done

echo "==> Instalando WordPress core..."
if cli wp core is-installed >/dev/null 2>&1; then
  echo "    WordPress ya estaba instalado, omitiendo."
else
  cli wp core install \
    --url="${WP_URL}" \
    --title="Probability WooCommerce Test" \
    --admin_user="${ADMIN_USER}" \
    --admin_password="${ADMIN_PASS}" \
    --admin_email="${ADMIN_EMAIL}" \
    --skip-email
fi

echo "==> Permalinks 'pretty' (requisito de la REST API de WooCommerce)..."
cli wp rewrite structure '/%postname%/' --hard
cli wp rewrite flush --hard

echo "==> Instalando y activando WooCommerce..."
cli wp plugin install woocommerce --activate

echo "==> mu-plugin: aceptar Basic Auth en la REST API sobre HTTP local..."
docker exec woo_test_wp sh -c 'mkdir -p /var/www/html/wp-content/mu-plugins && cat > /var/www/html/wp-content/mu-plugins/probability-rest-ssl.php <<"PHP"
<?php
if (isset($_SERVER["REQUEST_URI"]) && strpos($_SERVER["REQUEST_URI"], "/wp-json/") !== false) {
    $_SERVER["HTTPS"] = "on";
}
PHP
chown 33:33 /var/www/html/wp-content/mu-plugins/probability-rest-ssl.php'

echo "==> Creando API key REST (read_write)..."
PHP_CODE='$ck="ck_".wc_rand_hash(); $cs="cs_".wc_rand_hash(); global $wpdb; $wpdb->insert($wpdb->prefix."woocommerce_api_keys", array("user_id"=>1,"description"=>"probability-test","permissions"=>"read_write","consumer_key"=>wc_api_hash($ck),"consumer_secret"=>$cs,"truncated_key"=>substr($ck,-7))); echo $ck." ".$cs;'
KEYS=$(cli wp eval "$PHP_CODE")
CK=$(echo "$KEYS" | awk "{print \$1}")
CS=$(echo "$KEYS" | awk "{print \$2}")

echo "==> Creando webhook order.created -> ${WEBHOOK_URL}..."
cli wp wc webhook create \
  --name="Probability order.created" \
  --topic="order.created" \
  --delivery_url="${WEBHOOK_URL}" \
  --status="active" \
  --user="${ADMIN_USER}" || echo "    (si falla, crearlo manual en wp-admin -> WooCommerce -> Ajustes -> Avanzado -> Webhooks)"

cat <<EOF

============================================================
  WooCommerce de pruebas LISTO
============================================================
  Tienda:    ${WP_URL}
  wp-admin:  ${WP_URL}/wp-admin   (user: ${ADMIN_USER} / pass: ${ADMIN_PASS})

  REST API (usar al crear la integracion en Probability):
    Store URL:        http://localhost:8088
    Consumer Key:     ${CK}
    Consumer Secret:  ${CS}

  Webhook order.created -> ${WEBHOOK_URL}
  HMAC desactivado mientras WOOCOMMERCE_WEBHOOK_SECRET este vacio en el back.
============================================================
EOF
