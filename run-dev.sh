#!/usr/bin/env bash
#
# run-dev.sh — DEV tier: live-reload straight from the working tree.
#
#   backend  Air rebuilds the Go/PocketBase server on every .go save
#            (`-tags dev` → Automigrate ON + dev seeder), with an EPHEMERAL
#            database under tmp/dev_pb_data (wiped by Air on exit).
#   frontend Vite dev server with HMR, proxying /api to the backend.
#
# Exposed publicly (optional): set DEV_ALLOWED_HOST to the dev hostname and the
# Vite config adds it to allowedHosts and points the HMR websocket at
# wss://<host>:443, so HMR works through the cloudflared tunnel.
#
# There is NO Discord bot here — by design, not a toggle. The dev build
# (`-tags dev`) compiles the bot subsystem out entirely
# (internal/disgo.SubsystemEnabled == false), so a dev-tier server never opens a
# gateway and needs no token. Dev is the frontend/HMR tier; `pre` (test) is the
# lowest tier that runs a bot. See docs/DEPLOYMENTS.md.
#
# Config: ./.env.dev (optional, see .env.dev.example) or env vars.
#   DEV_PB_PORT       backend/PocketBase  (default 8092, internal)
#   DEV_VITE_PORT     Vite HMR server     (default 8093, tunnel target)
#   DEV_ALLOWED_HOST  public dev hostname (default empty = localhost only)
#
# Usage:  ./run-dev.sh        Ctrl-C stops both processes
set -euo pipefail
cd "$(dirname "$0")"

# Optional per-machine dev overrides.
if [ -f .env.dev ]; then
  set -a
  # shellcheck disable=SC1091
  . ./.env.dev
  set +a
fi

export DEV_PB_PORT="${DEV_PB_PORT:-8092}"
export DEV_VITE_PORT="${DEV_VITE_PORT:-8093}"
export DEV_ALLOWED_HOST="${DEV_ALLOWED_HOST:-}"

# cmd/server self-binds to PUBLIC_PB_PORT; point it at the dev backend port.
export PUBLIC_PB_PORT="$DEV_PB_PORT"

# No DISCORD_BOT_TOKEN is set here on purpose — the dev build has no bot
# subsystem, so there is nothing to configure and nothing opens a gateway.

command -v air >/dev/null 2>&1 || {
  echo "run-dev.sh: 'air' not found — go install github.com/air-verse/air@latest" >&2
  exit 1
}

echo "── dev ───────────────────────────────────────────────────"
echo "  backend (Air, PocketBase) : http://127.0.0.1:${DEV_PB_PORT}   (ephemeral DB)"
echo "  frontend (Vite HMR)       : http://127.0.0.1:${DEV_VITE_PORT}"
[ -n "$DEV_ALLOWED_HOST" ] && echo "  public dev host           : https://${DEV_ALLOWED_HOST}"
echo "  Discord bot               : none (compiled out of the dev build)"
echo "  Automigrate               : ON (dev build tag)"
echo "──────────────────────────────────────────────────────────"

air -c .air.dev.toml &
back=$!
( cd sveltekit && pnpm dev ) &
front=$!

trap 'kill "$back" "$front" 2>/dev/null || true' EXIT INT TERM
wait
