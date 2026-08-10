# stew-site-template

> **For AI assistants:** See [CLAUDE.md](CLAUDE.md) for development commands, conventions, and implementation details.

A reusable project template combining a Go backend with a SvelteKit frontend.

## Tech Stack

| Layer               | Technology                                                                           |
| ------------------- | ------------------------------------------------------------------------------------ |
| Backend / Auth / DB | [PocketBase](https://pocketbase.io/) v0.36.7 (Go framework, ServeMux router)         |
| Discord Bot         | [Disgo](https://github.com/disgoorg/disgo) v0.19.3                                   |
| WebSocket Server    | [coder/websocket](https://github.com/coder/websocket)                                |
| Frontend            | [SvelteKit](https://kit.svelte.dev/) 2 + [Skeleton UI v4](https://www.skeleton.dev/) |
| Frontend Build      | `@sveltejs/adapter-static` → served by PocketBase                                    |
| Build Orchestration | [Taskfile](https://taskfile.dev/)                                                    |
| Container           | [Podman](https://podman.io/)                                                         |

## Architecture Overview

```
┌──────────────────────────────────────────────┐
│            Go Binary (cmd/server)            │
│                                              │
│  ┌──────────────┐   ┌─────────────────────┐  │
│  │  PocketBase  │   │      Disgo Bot      │  │
│  │  - REST API  │   │  - Slash Commands   │  │
│  │  - Auth/JWT  │   │  - Event Listeners  │  │
│  │  - SQLite    │   └─────────────────────┘  │
│  │  - ServeMux  │                            │
│  │    Router    │   ┌─────────────────────┐  │
│  └──────┬───────┘   │      WebSocket      │  │
│         │           │  (coder/websocket)  │  │
│         │           │  - Optional JWT     │  │
│         │           │  - Hub / Rooms      │  │
│         │           └─────────────────────┘  │
│         │                                    │
│  guards/Services → cross-system DI           │
│  PB Hooks → Discord notifications            │
│  PB Hooks → WS Hub broadcasts                │
│  PB Routes → Auth-gated page serving         │
│                                              │
└─────────┬────────────────────────────────────┘
          │ serves
┌─────────▼───────┐
│   pb_public/    │ ← SvelteKit static build
│   (SvelteKit)   │
└─────────────────┘
```

## Project Structure

```
.
├── cmd/server/                # Go entrypoint
│   └── main.go
├── internal/
│   ├── guards/                # Unified cross-system guards + Services DI
│   │   ├── interfaces/
│   │   │   ├── discord/       # Per-method Discord interfaces (one per file)
│   │   │   ├── websocket/     # Per-method WS interfaces (one per file)
│   │   │   └── pocketbase/    # Per-method PB interfaces (one per file)
│   │   ├── services.go        # Services struct (bundles all system interfaces)
│   │   ├── guard.go           # GuardFunc type definition
│   │   └── require_*.go       # Guard implementations
│   ├── pocketbase/
│   │   ├── service.go         # PB service wrapper (implements pbiface.Service)
│   │   ├── schema/            # Programmatic collection definitions
│   │   ├── hooks/             # Record event hooks (PB → Discord, PB → WS)
│   │   ├── routes/            # Custom API routes + protected page serving
│   │   │   └── middleware/    # Auth middleware, role checks
│   │   ├── oauth/             # OAuth2 provider configuration
│   │   ├── seed/              # Dev-only seeder (compiled with -tags dev)
│   │   └── resolvers/         # PB data lookups (one function per file)
│   ├── disgo/
│   │   ├── bot.go             # Bot client + interface methods + lifecycle
│   │   ├── commands/          # Slash command definitions and handlers
│   │   ├── events/            # Discord gateway event listeners
│   │   ├── actions/           # Reusable Discord API calls
│   │   ├── resolvers/         # Discord data lookups via Services
│   │   └── components/        # UI builders (buttons, embeds, rows)
│   └── websocket/
│       ├── hub.go             # Client registry, rooms, message routing
│       ├── handler.go         # WS upgrade with optional JWT auth
│       ├── client.go          # Single connection read/write pumps
│       ├── message.go         # Wire format for WS messages
│       ├── handlers/          # Self-registering message type handlers
│       ├── rooms/             # Room type definitions with guard lists
│       └── resolvers/         # WS state lookups via Services
├── sveltekit/                 # SvelteKit frontend (Skeleton UI v4, adapter-static → pb_public/)
├── docs/                      # Project meta-docs (status, milestones, decisions, runbook) — see docs/README.md
├── reference/                 # (gitignored) local-only legacy/reference material — see docs/README.md
├── CHANGELOG.md               # Keep-a-Changelog, SemVer
├── .env.example               # Env template (shared by backend + frontend via envDir)
├── .air.toml                  # Go hot reload config
├── .gitignore
├── Taskfile.yml               # Build orchestration
├── Containerfile              # Multi-stage Podman/Docker build
├── compose.yml                # Container compose config
├── go.mod
└── LICENSE
```

## Meta-docs

Planning, status, and architectural decisions live in [`docs/`](docs/README.md), with this shape:

- [`docs/STATUS.md`](docs/STATUS.md) — current state (Goals / Now / Next / Maybe / Out of scope)
- [`docs/milestones/`](docs/milestones/README.md) — one markdown per milestone (`M??-kebab-name.md`)
- [`docs/decisions/`](docs/decisions/README.md) — one ADR per decision (`????-kebab-name.md`)
- [`docs/RUNBOOK.md`](docs/RUNBOOK.md) — operational procedures (deploy, backup, restore, secret rotation)
- [`CHANGELOG.md`](CHANGELOG.md) — Keep-a-Changelog, SemVer
- `reference/` — gitignored, local-only legacy/reference material

See [`docs/README.md`](docs/README.md) for the full convention.

## Prerequisites

- **Go 1.25+** — runs the backend; [go.dev/dl](https://go.dev/dl)
- **pnpm** _(preferred)_ — package manager for the frontend; `npm install -g pnpm`. npm and yarn work but the project is developed with pnpm.
- **Podman** _(optional)_ — for building and running containers; Docker works as a drop-in alternative.
- **Task** _(optional)_ — task runner for dev commands like `task dev` and `task build`; `go install github.com/go-task/task/v3/cmd/task@latest`. Without it, run `air` and `pnpm dev` in separate terminals.
- **Air** _(optional)_ — Go hot-reload; auto-rebuilds the server on `.go` file saves; `go install github.com/air-verse/air@latest`. Without it, use `go run ./cmd/server serve` and restart manually.
- **Lefthook** _(optional)_ — runs format/lint checks on staged files before each commit; `go install github.com/evilmartians/lefthook@latest`, then `task install:hooks` once.

## Quick Start

1. **Clone and scaffold.** `scaffold-site.sh` renames the module/image _and_
   assigns the project's port block, hostnames, and per-tier config, so the new
   site is standard-compliant from creation. **Claim a port block in
   [PORTS.md](PORTS.md) first.**

   ```bash
   git clone https://github.com/Stewball32/stew-site-template.git my-project
   cd my-project
   ./scripts/scaffold-site.sh \
     --module github.com/you/yourapp \
     --port-base 8130 \
     --prod-host yourapp.example \
     --tunnel your-tunnel
   go mod tidy
   ```

   It rewrites `go.mod` + imports + the image/compose project names, sets the
   tier ports (prod `base+0`, test `base+1`, dev-backend `base+2`, dev-vite
   `base+3`), fills the cloudflared ingress snippet and `docs/DEPLOYMENTS.md`,
   and creates `.env` / `.env.pre` / `.env.dev` with blank secrets so the project
   builds immediately. Add `--dry-run` to preview.

   (`scripts/rename-module.sh` remains for a rename-only change.)

2. **Configure environment:**

   ```bash
   cp .env.example .env
   # Only PUBLIC_PB_PORT is required. Discord bot and OAuth vars are optional.
   ```

3. **Install frontend dependencies:**

   ```bash
   cd sveltekit && pnpm install && cd ..
   ```

4. **Run in development:**

   ```bash
   task dev
   ```

5. **Build for production:**
   ```bash
   task build
   ./bin/server serve
   ```

## Deployment — three tiers

One interface per tier; ports come from this project's block in [PORTS.md](PORTS.md).
Full runbook: **[docs/DEPLOYMENTS.md](docs/DEPLOYMENTS.md)**.

| Tier | Command | Compose | Env |
| --- | --- | --- | --- |
| **dev** | `./run-dev.sh` | — (working tree, Air + Vite HMR) | `.env.dev` |
| **test** | `./deploy-pre.sh` | `compose.pre.yml` | `.env.pre` |
| **prod** | `./deploy-prod.sh` | `compose.yml` | `.env` |

Both deploy scripts do the same thing: refuse to run on the wrong branch or a
dirty tree → build → recreate the stack → apply pending migrations on boot →
poll `/api/health` and fail loudly if it never comes up. `down` and `logs`
subcommands on each; `ALLOW_DIRTY=1` / `ALLOW_ANY_BRANCH=1` to override a guard.

Every tier binds **loopback only** and is fronted by cloudflared — fill in
[`deploy/cloudflared-ingress.snippet.yml`](deploy/cloudflared-ingress.snippet.yml)
(the scaffold does this for you) and merge it into the tunnel config.

## Database migrations

Collections are defined by **migrations** (`migrations/`), not by on-serve schema
code — see **[docs/MIGRATIONS.md](docs/MIGRATIONS.md)**. Automigrate is ON in dev
builds only: change the schema in the dev admin UI, a migration file is written,
you review and commit it, and test/prod apply it on boot (tracked in
`_migrations`). Prove schema changes on the test tier against a copy of prod
`pb_data` before prod.

The history is **granular in dev, squashed on release**: dev/test branches carry
many small migrations, and on approval they're collapsed into a single release
migration on the way to `main` — so prod's `_migrations` reads like a version log.

```bash
task migrate:up                        # apply pending
task migrate:create -- add_widgets     # new blank migration

# on approval, heading to main:
task migrate:squash -- --version v0.3.0            # granular ──► one release file
task migrate:verify -- --from /tmp/prod-pb_data    # prove it matches, THEN merge
```

## Production deployment

The container runs as non-root user `app` (UID 1000) and exposes a healthcheck on `/api/health`. A few things to know:

- **Never commit a populated `.env`.** The committed `.env.example` is the public template; real secrets stay out of the repo. Inject them via your orchestrator:
  - **Compose**: pass `--env-file` to `docker compose` / `podman-compose`, or use a `compose.override.yml` that's gitignored.
  - **Kubernetes**: mount a `Secret` and reference it in the pod spec.
  - **Bare `podman run` / `docker run`**: pass each variable with `-e VAR=value`.
- **Persistent data**: the `pb_data` volume holds the SQLite DB and uploads — back it up regularly. On first run, the volume is initialized owned by UID 1000; if you bind-mount a host directory instead, `chown 1000:1000` it first.
- **TLS**: PocketBase can terminate TLS itself with `--https`, but most deployments put it behind a reverse proxy (Caddy, Cloudflare Tunnel, nginx) that handles certificates.

See the [PocketBase deployment docs](https://pocketbase.io/docs/going-to-production/) for backup, migration, and admin-creation patterns.

## Notes

- `pb_data/` — PocketBase runtime data (SQLite DB, uploads). Created at runtime, gitignored. Nothing wipes it automatically — if data disappears, check for `git clean -fdx` in your workflow.
- `pb_public/` — SvelteKit build output. Created by `task build:frontend`, gitignored.
- Schema can be managed via PocketBase admin UI or programmatically in `internal/pocketbase/schema/`.
- Protected pages are served through auth-gated custom routes; public pages are served directly from `pb_public/`.
