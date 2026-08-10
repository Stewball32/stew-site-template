#!/usr/bin/env bash
#
# deploy-prod.sh — deploy the PROD tier (compose.yml + .env).
#
# Uniform tier interface (same flags/behaviour as deploy-pre.sh):
#   ./deploy-prod.sh          guard → build → deploy → health-check
#   ./deploy-prod.sh down     stop the stack (keeps the database volume)
#   ./deploy-prod.sh logs     follow logs
#
# Guards (override only when you mean it):
#   ALLOW_ANY_BRANCH=1   deploy from a branch other than PROD_BRANCH
#   ALLOW_DIRTY=1        deploy with uncommitted changes
#
# Pending database migrations apply automatically on boot (docs/MIGRATIONS.md).
set -euo pipefail
cd "$(dirname "$0")"
# shellcheck source=scripts/deploy-lib.sh
. scripts/deploy-lib.sh

TIER="prod"
PROD_BRANCH="${PROD_BRANCH:-main}"
COMPOSE_FILE="compose.yml"

case "${1:-up}" in
  down) exec docker compose -f "$COMPOSE_FILE" down ;;
  logs) exec docker compose -f "$COMPOSE_FILE" logs -f ;;
esac

require_cmd docker git curl
require_file .env "copy .env.example to .env and fill it in"

# Host port must match what the cloudflared ingress points at (see deploy/).
PORT="$(grep -E '^PROD_PB_PORT=' .env | cut -d= -f2- | tr -d '"' || true)"
PORT="${PORT:-8090}"

say "── ${TIER} ───────────────────────────────────────────────"
info "branch expected : ${PROD_BRANCH}"
info "compose file    : ${COMPOSE_FILE}"
info "host port       : 127.0.0.1:${PORT}"
say

require_branch "$PROD_BRANCH"
require_clean_tree

compose_deploy "$TIER" "$PORT" -f "$COMPOSE_FILE"

say
ok "prod deployed. Migrations (if any) applied on boot — check logs:"
info "docker compose -f ${COMPOSE_FILE} logs -f"
