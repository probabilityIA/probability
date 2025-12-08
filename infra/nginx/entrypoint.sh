#!/bin/sh
set -e
envsubst '\$DOMAIN \$SSL_CERT_PATH \$SSL_KEY_PATH' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf
exec nginx -g 'daemon off;'
