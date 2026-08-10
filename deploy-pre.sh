#!/usr/bin/env bash
#
# deploy-pre.sh — deploy the TEST tier (compose.pre.yml + .env.pre).
#
# The gate before prod: same image build, fully isolated project/container/port/
# volume. Prove migrations and behaviour here — ideally against a COPY of prod's
# pb_data (see docs/MIGRATIONS.md) — before running deploy-prod.sh.
#
# Uniform tier interface (same flags/behaviour as deploy-prod.sh):
#   ./deploy-pre.sh          guard → build → deploy → health-check
#   ./deploy-pre.sh down     stop the stack (keeps the test database volume)
#   ./deploy-pre.sh logs     follow logs
#
# Guards (override only when you mean it):
#   ALLOW_ANY_BRANCH=1   deploy from a branch other than PRE_BRANCH
#   ALLOW_DIRTY=1        deploy with uncommitted changes
set -euo pipefail
cd "$(dirname "$0")"
# shellcheck source=scripts/deploy-lib.sh
. scripts/deploy-lib.sh

TIER="test (pre)"
# The test tier usually tracks an integration branch. Default to `pre`; set
# PRE_BRANCH=main if this project promotes straight off main.
PRE_BRANCH="${PRE_BRANCH:-pre}"
COMPOSE_FILE="compose.pre.yml"

case "${1:-up}" in
  down) exec docker compose -f "$COMPOSE_FILE" down ;;
  logs) exec docker compose -f "$COMPOSE_FILE" logs -f ;;
esac

require_cmd docker git curl
require_file .env.pre "copy .env.pre.example to .env.pre and fill it in"

PORT="$(grep -E '^PRE_PB_PORT=' .env.pre | cut -d= -f2- | tr -d '"' || true)"
PORT="${PORT:-8091}"

say "── ${TIER} ───────────────────────────────────────────────"
info "branch expected : ${PRE_BRANCH}"
info "compose file    : ${COMPOSE_FILE}"
info "host port       : 127.0.0.1:${PORT}"
say

require_branch "$PRE_BRANCH"
require_clean_tree

compose_deploy "$TIER" "$PORT" -f "$COMPOSE_FILE"

say
ok "test tier deployed. Verify here BEFORE prod:"
info "docker compose -f ${COMPOSE_FILE} logs -f    # migration errors?"
info "curl -sI http://127.0.0.1:${PORT}/"
