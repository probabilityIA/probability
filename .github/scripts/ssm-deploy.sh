#!/usr/bin/env bash
# Ejecuta un script de deploy en el EC2 de produccion via SSM (sin SSH ni llave .pem).
#
# Uso:
#   ssm-deploy.sh <script-local> [VAR=valor ...]
#
# Sube el script (y lo que haya en el directorio artifacts/ si existe) al bucket
# de artefactos, lo ejecuta con AWS-RunShellScript como usuario ubuntu, y
# transmite la salida. Devuelve el mismo codigo de salida del script remoto.
set -euo pipefail

INSTANCE_ID="${SSM_INSTANCE_ID:-i-0f3284d2a87127e57}"
BUCKET="${DEPLOY_BUCKET:-probability-deploy-artifacts}"
REGION="${AWS_REGION:-us-east-1}"
RUN_ID="${GITHUB_RUN_ID:-manual-$(date +%s)}"
PREFIX="deploy/${RUN_ID}"

SCRIPT_LOCAL="${1:?falta el script de deploy}"
shift || true

ENV_EXPORTS=""
for kv in "$@"; do
  ENV_EXPORTS="${ENV_EXPORTS}export ${kv%%=*}='${kv#*=}'; "
done

echo "📤 Subiendo artefactos a s3://${BUCKET}/${PREFIX}/"
aws s3 cp "$SCRIPT_LOCAL" "s3://${BUCKET}/${PREFIX}/deploy.sh" --only-show-errors
if [ -d artifacts ]; then
  aws s3 cp artifacts "s3://${BUCKET}/${PREFIX}/artifacts" --recursive --only-show-errors
fi

REMOTE_CMD="set -e
rm -rf /tmp/${RUN_ID} && mkdir -p /tmp/${RUN_ID}
aws s3 cp s3://${BUCKET}/${PREFIX} /tmp/${RUN_ID} --recursive --only-show-errors
chmod +x /tmp/${RUN_ID}/deploy.sh
chown -R ubuntu:ubuntu /tmp/${RUN_ID}
runuser -l ubuntu -c \"${ENV_EXPORTS}export DEPLOY_DIR=/tmp/${RUN_ID}; bash /tmp/${RUN_ID}/deploy.sh\"
rm -rf /tmp/${RUN_ID}"

echo "🚀 Enviando comando SSM a ${INSTANCE_ID}"
CMD_ID=$(aws ssm send-command \
  --region "$REGION" \
  --cli-input-json "$(jq -n --arg iid "$INSTANCE_ID" --arg cmd "$REMOTE_CMD" --arg c "deploy ${RUN_ID}" \
      '{InstanceIds:[$iid], DocumentName:"AWS-RunShellScript", Comment:$c, TimeoutSeconds:600,
        Parameters:{commands:($cmd|split("\n")), executionTimeout:["1800"]}}')" \
  --query 'Command.CommandId' --output text)

echo "🆔 CommandId: ${CMD_ID}"

STATUS=Pending
for _ in $(seq 1 120); do
  sleep 10
  STATUS=$(aws ssm get-command-invocation --region "$REGION" \
    --command-id "$CMD_ID" --instance-id "$INSTANCE_ID" \
    --query 'Status' --output text 2>/dev/null || echo Pending)
  case "$STATUS" in
    Success|Failed|Cancelled|TimedOut) break ;;
  esac
  echo "⏳ ${STATUS}..."
done

echo "──────── salida remota ────────"
aws ssm get-command-invocation --region "$REGION" \
  --command-id "$CMD_ID" --instance-id "$INSTANCE_ID" \
  --query 'StandardOutputContent' --output text
ERR=$(aws ssm get-command-invocation --region "$REGION" \
  --command-id "$CMD_ID" --instance-id "$INSTANCE_ID" \
  --query 'StandardErrorContent' --output text)
[ -n "$ERR" ] && [ "$ERR" != "None" ] && { echo "──────── stderr ────────"; echo "$ERR"; }
echo "───────────────────────────────"

aws s3 rm "s3://${BUCKET}/${PREFIX}" --recursive --only-show-errors || true

if [ "$STATUS" != "Success" ]; then
  echo "❌ Deploy fallo con estado: ${STATUS}"
  exit 1
fi
echo "✅ Deploy completado via SSM"
