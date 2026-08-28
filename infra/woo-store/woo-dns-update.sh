#!/bin/bash
set -u
export AWS_DEFAULT_REGION=us-east-1
ZONE_ID="Z06743392M4CV4HQ9DTIT"
RECORD="woo.probabilityia.com.co"
LOG="/var/log/woo-dns.log"

log() { echo "$(date -Is) $*" | tee -a "$LOG"; }

get_ip() {
  local token ip
  token=$(curl -s -m 5 -X PUT "http://169.254.169.254/latest/api/token" \
    -H "X-aws-ec2-metadata-token-ttl-seconds: 300") || return 1
  ip=$(curl -s -m 5 -H "X-aws-ec2-metadata-token: $token" \
    "http://169.254.169.254/latest/meta-data/public-ipv4") || return 1
  [ -n "$ip" ] || return 1
  echo "$ip"
}

IP=""
for i in 1 2 3 4 5 6 7 8 9 10; do
  IP=$(get_ip) && [ -n "$IP" ] && break
  log "intento $i: sin IP publica todavia"
  sleep 5
done

if [ -z "$IP" ]; then
  log "ERROR: no se obtuvo IP publica, se aborta"
  exit 1
fi

CUR=$(aws route53 list-resource-record-sets --hosted-zone-id "$ZONE_ID" \
  --start-record-name "$RECORD" --start-record-type A --max-items 1 \
  --query 'ResourceRecordSets[0].ResourceRecords[0].Value' --output text 2>/dev/null || echo "")

if [ "$CUR" = "$IP" ]; then
  log "DNS ya apunta a $IP, sin cambios"
  exit 0
fi

cat > /tmp/woo-dns-change.json <<EOF
{"Comment":"woo-store autoregistro en arranque",
 "Changes":[{"Action":"UPSERT","ResourceRecordSet":{
   "Name":"$RECORD","Type":"A","TTL":60,
   "ResourceRecords":[{"Value":"$IP"}]}}]}
EOF

for i in 1 2 3; do
  if aws route53 change-resource-record-sets --hosted-zone-id "$ZONE_ID" \
      --change-batch file:///tmp/woo-dns-change.json >>"$LOG" 2>&1; then
    log "DNS actualizado: ${CUR:-vacio} -> $IP"
    exit 0
  fi
  log "intento $i: fallo el UPSERT, reintentando"
  sleep 10
done

log "ERROR: no se pudo actualizar el DNS a $IP"
exit 1
