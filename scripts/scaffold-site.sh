#!/usr/bin/env bash
#
# scaffold-site.sh — turn a fresh clone of stew-site-template into a
# standard-compliant site: renames the Go module + image, assigns the project's
# port block, and writes the per-tier env/compose/deploy/ingress config.
#
# Run ONCE, right after cloning. Supersedes rename-module.sh (which only did the
# module rename); this calls it, then does everything else.
#
# Usage:
#   ./scripts/scaffold-site.sh \
#       --module github.com/Stewball32/align-the-day \
#       --port-base 8110 \
#       --prod-host alignthe.day \
#       [--name align-the-day] \
#       [--test-host pre.alignthe.day] [--dev-host lab.alignthe.day] \
#       [--tunnel xemu-cartographer]
#
# Required:
#   --module      Go module path (also sets the image/project name)
#   --port-base   base of this project's block from PORTS.md (prod=base+0,
#                 test=base+1, dev-backend=base+2, dev-vite=base+3)
#   --prod-host   public hostname for prod (apex or subdomain — both supported)
#
# Optional:
#   --name        project/image name        (default: basename of --module)
#   --test-host   test hostname             (default: pre.<prod-host>)
#   --dev-host    dev hostname              (default: lab.<prod-host>)
#   --tunnel      cloudflared tunnel name   (default: <TUNNEL> placeholder)
#   --dry-run     print what would change, write nothing
#
# ⚠️ Claim the port block in PORTS.md FIRST — this script does not reserve it.
set -euo pipefail

MODULE=""; PORT_BASE=""; PROD_HOST=""; NAME=""; TEST_HOST=""; DEV_HOST=""
TUNNEL="<TUNNEL>"; DRY_RUN=0

die() { printf '✗ %s\n' "$*" >&2; exit 1; }
say() { printf '%s\n' "$*"; }

while [ $# -gt 0 ]; do
  case "$1" in
    --module)     MODULE="${2:-}"; shift 2 ;;
    --port-base)  PORT_BASE="${2:-}"; shift 2 ;;
    --prod-host)  PROD_HOST="${2:-}"; shift 2 ;;
    --name)       NAME="${2:-}"; shift 2 ;;
    --test-host)  TEST_HOST="${2:-}"; shift 2 ;;
    --dev-host)   DEV_HOST="${2:-}"; shift 2 ;;
    --tunnel)     TUNNEL="${2:-}"; shift 2 ;;
    --dry-run)    DRY_RUN=1; shift ;;
    -h|--help)    sed -n '2,40p' "$0"; exit 0 ;;
    *)            die "unknown argument: $1" ;;
  esac
done

[ -n "$MODULE" ]    || die "--module is required"
[ -n "$PORT_BASE" ] || die "--port-base is required (see PORTS.md)"
[ -n "$PROD_HOST" ] || die "--prod-host is required"
case "$PORT_BASE" in ''|*[!0-9]*) die "--port-base must be a number" ;; esac

NAME="${NAME:-$(basename "$MODULE")}"
TEST_HOST="${TEST_HOST:-pre.${PROD_HOST}}"
DEV_HOST="${DEV_HOST:-lab.${PROD_HOST}}"

PROD_PORT="$PORT_BASE"
PRE_PORT="$((PORT_BASE + 1))"
DEV_PB_PORT="$((PORT_BASE + 2))"
DEV_VITE_PORT="$((PORT_BASE + 3))"

cd "$(dirname "$0")/.."

say "── scaffold ──────────────────────────────────────────────"
say "  module      : $MODULE"
say "  name/image  : $NAME"
say "  port block  : ${PORT_BASE}-$((PORT_BASE + 9))"
say "    prod      : $PROD_PORT   → $PROD_HOST"
say "    test      : $PRE_PORT   → $TEST_HOST"
say "    dev (be)  : $DEV_PB_PORT   (internal)"
say "    dev (vite): $DEV_VITE_PORT   → $DEV_HOST"
say "  tunnel      : $TUNNEL"
say "──────────────────────────────────────────────────────────"

if [ "$DRY_RUN" = "1" ]; then
  say "(--dry-run: nothing written)"
  exit 0
fi

# --- 1. module + image rename ------------------------------------------------
if grep -q '^module github.com/youruser/yourproject$' go.mod 2>/dev/null; then
  ./scripts/rename-module.sh "$MODULE" "$NAME"
else
  say "! module already renamed — skipping rename step"
fi

# --- 2. port defaults in the committed env templates -------------------------
sed -i \
  -e "s/^PROD_PB_PORT=.*/PROD_PB_PORT=${PROD_PORT}/" \
  .env.example
sed -i \
  -e "s/^PRE_PB_PORT=.*/PRE_PB_PORT=${PRE_PORT}/" \
  .env.pre.example
sed -i \
  -e "s/^DEV_PB_PORT=.*/DEV_PB_PORT=${DEV_PB_PORT}/" \
  -e "s/^DEV_VITE_PORT=.*/DEV_VITE_PORT=${DEV_VITE_PORT}/" \
  -e "s/^DEV_ALLOWED_HOST=.*/DEV_ALLOWED_HOST=${DEV_HOST}/" \
  .env.dev.example

# --- 3. compose defaults (fallbacks match this project's block) --------------
sed -i -e "s/\${PROD_PB_PORT:-[0-9]*}/\${PROD_PB_PORT:-${PROD_PORT}}/" compose.yml
sed -i -e "s/\${PRE_PB_PORT:-[0-9]*}/\${PRE_PB_PORT:-${PRE_PORT}}/" compose.pre.yml

# --- 4. run-dev.sh defaults --------------------------------------------------
sed -i \
  -e "s/DEV_PB_PORT:-[0-9]*/DEV_PB_PORT:-${DEV_PB_PORT}/" \
  -e "s/DEV_VITE_PORT:-[0-9]*/DEV_VITE_PORT:-${DEV_VITE_PORT}/" \
  run-dev.sh

# --- 5. cloudflared ingress snippet — fill the placeholders ------------------
sed -i \
  -e "s|<PROD_HOST>|${PROD_HOST}|g" \
  -e "s|<TEST_HOST>|${TEST_HOST}|g" \
  -e "s|<DEV_HOST>|${DEV_HOST}|g" \
  -e "s|<PROD_PORT>|${PROD_PORT}|g" \
  -e "s|<TEST_PORT>|${PRE_PORT}|g" \
  -e "s|<DEV_PORT>|${DEV_VITE_PORT}|g" \
  -e "s|<PROJECT>|${NAME}|g" \
  -e "s|<TUNNEL>|${TUNNEL}|g" \
  deploy/cloudflared-ingress.snippet.yml

# --- 6. DEPLOYMENTS.md — fill the project table ------------------------------
sed -i \
  -e "s|\`<DEV_VITE_PORT>\`|\`${DEV_VITE_PORT}\`|" \
  -e "s|\`<DEV_PB_PORT>\`|\`${DEV_PB_PORT}\`|" \
  -e "s|\`<PRE_PB_PORT>\`|\`${PRE_PORT}\`|" \
  -e "s|\`<PROD_PB_PORT>\`|\`${PROD_PORT}\`|" \
  -e "s|\`<DEV_HOST>\`|\`${DEV_HOST}\`|" \
  -e "s|\`<TEST_HOST>\`|\`${TEST_HOST}\`|" \
  -e "s|\`<PROD_HOST>\`|\`${PROD_HOST}\`|" \
  docs/DEPLOYMENTS.md

# --- 7. create the local (gitignored) tier env files -------------------------
# The frontend imports PUBLIC_PB_PORT from $env/static/public, so `pnpm build`
# fails without a .env present. Create them now (blank secrets) so the scaffolded
# site builds and runs immediately. Never overwrite an existing file.
for pair in ".env.example:.env" ".env.pre.example:.env.pre" ".env.dev.example:.env.dev"; do
  src="${pair%%:*}"; dst="${pair##*:}"
  if [ -f "$dst" ]; then
    say "! $dst already exists — left alone"
  elif [ -f "$src" ]; then
    cp "$src" "$dst"
    say "  created $dst (from $src)"
  fi
done

say
say "✓ scaffolded."
say
say "Next steps:"
say "  1. Record the ${PORT_BASE}-$((PORT_BASE + 9)) block for '${NAME}' in PORTS.md"
say "  2. go mod tidy && (cd sveltekit && pnpm install)"
say "  3. Fill in secrets in .env / .env.pre (created for you, values blank)"
say "  4. ./run-dev.sh                                            # verify dev boots"
say "  5. Merge deploy/cloudflared-ingress.snippet.yml into the tunnel config"
say "  6. git add -A && git commit -m 'chore: scaffold ${NAME} from template'"
