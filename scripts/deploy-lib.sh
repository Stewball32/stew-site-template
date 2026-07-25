#!/usr/bin/env bash
# Shared deploy helpers — sourced by deploy-prod.sh / deploy-pre.sh so every
# tier behaves identically: guard → build → deploy → health-check.
#
# Not executable on its own.

set -euo pipefail

# --- output ------------------------------------------------------------------
_c_red=$'\033[31m'; _c_grn=$'\033[32m'; _c_ylw=$'\033[33m'; _c_off=$'\033[0m'
say()  { printf '%s\n' "$*"; }
info() { printf '  %s\n' "$*"; }
ok()   { printf '%s✓%s %s\n' "$_c_grn" "$_c_off" "$*"; }
warn() { printf '%s!%s %s\n' "$_c_ylw" "$_c_off" "$*" >&2; }
die()  { printf '%s✗%s %s\n' "$_c_red" "$_c_off" "$*" >&2; exit 1; }

# --- guards ------------------------------------------------------------------

# require_cmd <cmd>...
require_cmd() {
  local c
  for c in "$@"; do
    command -v "$c" >/dev/null 2>&1 || die "required command not found: $c"
  done
}

# require_file <path> <hint>
require_file() {
  [ -f "$1" ] || die "missing $1 — $2"
}

# require_branch <expected>  (skip with ALLOW_ANY_BRANCH=1)
require_branch() {
  local want="$1" have
  have="$(git rev-parse --abbrev-ref HEAD)"
  if [ "$have" != "$want" ]; then
    [ "${ALLOW_ANY_BRANCH:-0}" = "1" ] \
      && { warn "on '$have', expected '$want' (ALLOW_ANY_BRANCH=1)"; return 0; }
    die "on branch '$have', expected '$want'. Checkout $want, or re-run with ALLOW_ANY_BRANCH=1."
  fi
  ok "branch: $have"
}

# require_clean_tree  (skip with ALLOW_DIRTY=1)
# Deploys build from the working tree, so a dirty tree means shipping something
# that isn't committed — you couldn't reproduce or roll back to it.
require_clean_tree() {
  if [ -n "$(git status --porcelain)" ]; then
    [ "${ALLOW_DIRTY:-0}" = "1" ] \
      && { warn "working tree is dirty (ALLOW_DIRTY=1)"; return 0; }
    git status --short >&2
    die "working tree is dirty. Commit or stash first, or re-run with ALLOW_DIRTY=1."
  fi
  ok "working tree clean"
}

# --- deploy ------------------------------------------------------------------

# compose_deploy <tier-label> <host-port> [-f composefile ...]
# Builds the image and recreates the stack, then waits for health.
compose_deploy() {
  local label="$1" port="$2"; shift 2
  local -a compose=(docker compose "$@")

  say "── building ${label} ─────────────────────────────────────"
  "${compose[@]}" build

  say "── deploying ${label} ────────────────────────────────────"
  # `up -d --force-recreate`: always bounce the container, even when the image
  # and config are byte-identical (unchanged code). Without --force-recreate an
  # unchanged rebuild is a no-op — the old process keeps running and the on-boot
  # startup reconcile (Discord scheduled-events sync, migrations) never re-fires.
  # Recreate also guarantees a changed env_file is re-read. Migrations apply on boot.
  "${compose[@]}" up -d --force-recreate

  "${compose[@]}" ps
  wait_healthy "$label" "$port" "${compose[@]}"
}

# wait_healthy <label> <host-port> <compose cmd...>
wait_healthy() {
  local label="$1" port="$2"; shift 2
  local -a compose=("$@")
  local url="http://127.0.0.1:${port}/api/health"
  local i

  say "── health check (${url}) ─────────────────"
  for i in $(seq 1 30); do
    if curl -fsS -o /dev/null "$url" 2>/dev/null; then
      ok "${label} healthy on 127.0.0.1:${port}"
      return 0
    fi
    sleep 2
  done

  warn "${label} did not become healthy in ~60s — last logs:"
  "${compose[@]}" logs --tail=40 >&2 || true
  die "deploy failed health check"
}
