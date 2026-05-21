#!/usr/bin/env bash
# ═══════════════════════════════════════════════════════════════════════
#  gh-env.sh — Exporta GH_TOKEN scoped al repo `secamc93/probability`
#  sin tocar la configuración global (~/.config/gh).
#
#  Uso:
#      eval "$(./scripts/gh-env.sh)"
#      gh pr list
#
#  Búsqueda del token (primer hit gana):
#    1. archivo `.gh-token`             ← una línea con el PAT
#    2. `GITHUB_PERSONAL_ACCESS_TOKEN`  en `.mcp.json`
#
#  Ambos archivos deberían estar en `.gitignore`.
# ═══════════════════════════════════════════════════════════════════════

set -euo pipefail
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

REPO_OWNER="probabilityIA"
REPO_NAME="probability"
REPO_FULL="${REPO_OWNER}/${REPO_NAME}"
# Cuenta de usuario que dueña del PAT (para el banner; el PAT pertenece a
# este usuario pero opera sobre el repo del org de arriba).
GH_USER="secamc93"

token=""
if [ -f "$PROJECT_ROOT/.gh-token" ]; then
    token="$(tr -d '[:space:]' < "$PROJECT_ROOT/.gh-token")"
elif [ -f "$PROJECT_ROOT/.mcp.json" ]; then
    token="$(grep -oE '"GITHUB_PERSONAL_ACCESS_TOKEN"\s*:\s*"[^"]+"' "$PROJECT_ROOT/.mcp.json" \
        | head -1 \
        | sed -E 's/.*"([^"]+)"$/\1/')"
fi

if [ -z "$token" ]; then
    echo "echo '❌ No se encontró GH_TOKEN. Crea .gh-token o ajusta .mcp.json.' >&2; return 1 2>/dev/null || exit 1"
    exit 1
fi

echo "export GH_TOKEN=$token"
echo "export GITHUB_TOKEN=$token"
echo "export GH_REPO=$REPO_FULL"
echo "export GH_HOST=github.com"
echo "export GH_PROMPT_INFO='[gh:$REPO_FULL]'"
echo "echo '✅ gh CLI scoped → $REPO_FULL (cuenta $GH_USER)'"
echo "echo '   GH_TOKEN     = ${token:0:18}…'"
echo "echo '   GH_REPO      = $REPO_FULL  (gh usa este repo por defecto)'"
