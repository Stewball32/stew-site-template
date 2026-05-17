# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

> See also: [README.md](README.md) for project overview, tech stack, architecture diagram, and quick-start guide.

## Reference docs

Before writing or reviewing code that touches a third-party library where the API may have drifted from your training data, consult up-to-date docs rather than guessing.

- **Skeleton UI v4** — [sveltekit/docs/skeleton-llms.txt](sveltekit/docs/skeleton-llms.txt) is a table of contents of Skeleton's official docs (components, theming, Tailwind v4 integration). Read it first to locate the right page, then WebFetch the specific page under `https://www.skeleton.dev/` (e.g. `https://www.skeleton.dev/docs/svelte/framework-components/app-bar.md`, `https://www.skeleton.dev/docs/svelte/tailwind-components/buttons`). Always use the **Svelte** section, not React.
- **SvelteKit, PocketBase JS SDK, Disgo, Tailwind v4** — WebFetch the official docs site (`kit.svelte.dev`, `pocketbase.io/docs`, `disgo.dev`, `tailwindcss.com`) rather than inventing an API.

## Project meta-docs

[`docs/README.md`](docs/README.md) is the source of truth for the project's meta-docs convention. Follow it when adding planning, status, or decision documents.

- New milestone → copy [`docs/milestones/_template.md`](docs/milestones/_template.md) to `docs/milestones/M??-kebab-name.md`, then add a row to `docs/milestones/README.md`.
- New decision → copy [`docs/decisions/_template.md`](docs/decisions/_template.md) to `docs/decisions/????-kebab-name.md`, then add a row to `docs/decisions/README.md`. ADRs are immutable once Accepted — supersede with a new ADR rather than editing.
- Update [`docs/STATUS.md`](docs/STATUS.md) whenever the "Now" set of work changes.
- Add user-visible changes to `CHANGELOG.md` under `[Unreleased]`; cut a versioned section when shipping (SemVer, Keep-a-Changelog format).
- Dates are always absolute (`YYYY-MM-DD`). Milestone `Log` sections and the ADR index are append-only.

[`docs/README.md`](docs/README.md) also defines cross-cutting conventions: Conventional Commit messages, TODO/FIXME tagged with the related `M??`/`ADR-????`, the standard Taskfile target lineup (`dev`/`build`/`test`/`fmt`/`lint`/`clean`/`container:*`), `v`-prefixed SemVer tags, and `.env.example` grouping. Versioning is driven from the git tag: `internal/version` exposes `Version`/`Commit`/`Date` set via ldflags in `Taskfile.yml`, surfaced at `/api/version`. See [ADR-0001](docs/decisions/0001-version-from-git-tag.md).

## Development Commands

```sh

# Install task runner and hot reload

go install github.com/go-task/task/v3/cmd/task@latest
go install github.com/air-verse/air@latest

# Run both backend and frontend dev servers

task dev

# Backend only (hot reload)

task dev:backend

# Frontend only (run from sveltekit/)

task dev:frontend

# Build for production

task build

# Build and run container

task container:build
task container:run

# Clean build artifacts

task clean

# Install Git pre-commit hooks (lefthook)

task install:hooks

# Aggregate test / format / lint across Go + frontend

task test
task fmt
task lint

# Run server directly (no Task/Air)

go run ./cmd/server serve
./bin/server serve

# Run Go tests

go test ./...                              # all packages
go test ./internal/guards/...              # single package
go test -run TestAny ./internal/guards     # single test

# Frontend type-check, lint, format (run from sveltekit/)

cd sveltekit
pnpm check          # svelte-check + TypeScript
pnpm lint           # prettier + eslint
pnpm format         # prettier --write

# Generate PocketBase TypeScript types (requires running dev server).
# Also runs automatically in the background during `task dev:frontend`.

task typegen
```

## Architecture

Single Go binary (`cmd/server`) runs three concurrent systems:

1. **PocketBase** — REST API, auth (JWT), SQLite database, static file server (serves `pb_public/`), uses `net/http.ServeMux` router
2. **Disgo Discord bot** — connects via gateway in PocketBase's OnServe hook, non-blocking
3. **WebSocket handler** (`coder/websocket`) — custom route on PocketBase's router with optional JWT auth, Hub for managing clients/rooms/broadcasting

The SvelteKit frontend is built with `@sveltejs/adapter-static` into `pb_public/`, which PocketBase serves automatically. The `fallback: 'index.html'` config enables SPA-style client-side routing.

Protected pages can be served through custom PocketBase routes that validate JWT auth before serving the static file, while public pages are served directly from `pb_public/`.

## Backend Structure

### Startup sequence (`cmd/server/main.go`)

1. Create PocketBase app instance
2. Register record lifecycle hooks — `hooks.RegisterAll(app)` (callback registration, fires later)
3. In OnServe hook:
   - Register collection schemas (`schema.RegisterAll(app)`)
   - Register OAuth2 providers (`oauth.RegisterAll(app)`) — must run after schema
   - Register custom API routes (`routes.RegisterAll(se)`)
   - Initialize WebSocket Hub, start its Run() goroutine, mount `/api/ws` endpoint
   - Start Disgo bot gateway connection (non-blocking)
   - Wire cross-system `Services` struct — connects all three systems via interfaces
4. Register OnTerminate hook to shut down Disgo bot cleanly
5. Call `app.Start()` (blocking)

Each system has its own README with deeper per-subdirectory detail — read [`internal/guards/README.md`](internal/guards/README.md), [`internal/pocketbase/README.md`](internal/pocketbase/README.md), [`internal/disgo/README.md`](internal/disgo/README.md), [`internal/websocket/README.md`](internal/websocket/README.md) before modifying that system.

### Key packages

| Package                                 | Purpose                                                             |
| --------------------------------------- | ------------------------------------------------------------------- |
| `internal/guards`                       | Unified cross-system guards, `Services` struct, `GuardFunc` type    |
| `internal/guards/interfaces/discord`    | Per-method Discord interfaces (Membership, Roles, Notify, Voice)    |
| `internal/guards/interfaces/websocket`  | Per-method WS interfaces (Connected, Rooms, Broadcast)              |
| `internal/guards/interfaces/pocketbase` | Per-method PB interfaces (Users)                                    |
| `internal/version`                      | Build-time `Version`/`Commit`/`Date` (set via ldflags, surfaced at `/api/version`) — see ADR-0001 |
| `internal/pocketbase`                   | PB service wrapper — implements `pbiface.Service`                   |
| `internal/pocketbase/schema`            | Programmatic collection definitions — one file per domain           |
| `internal/pocketbase/hooks`             | Record event hooks — fire Discord notifications, push to WS clients |
| `internal/pocketbase/oauth`             | OAuth2 provider config — env-driven, self-registering, one per file |
| `internal/pocketbase/routes`            | Custom endpoints + protected page serving via auth-gated routes     |
| `internal/pocketbase/routes/middleware` | Auth middleware, role checks, global middleware registry            |
| `internal/pocketbase/routes/admin`      | Route group for `/api/admin` — auth + admin middleware              |
| `internal/pocketbase/seed`              | In-process dev seeder — `data.go` defines seed vars, compiled only with `-tags dev` |
| `internal/pocketbase/resolvers`         | PB data lookups — one exported function per file                    |
| `internal/websocket`                    | WebSocket Hub, client management, message routing, JWT auth upgrade |
| `internal/websocket/handlers`           | Self-registering WS message handlers — one per file                 |
| `internal/websocket/rooms`              | Room type definitions with guard lists — one per file               |
| `internal/websocket/resolvers`          | WS state lookups via Services — one exported function per file      |
| `internal/websocket/actions`            | Reusable WS Hub operations — one exported function per file         |
| `internal/disgo`                        | Bot client setup: NewBot(), OpenGateway(), Close(), action methods  |
| `internal/disgo/commands`               | Slash command definitions and handler functions                     |
| `internal/disgo/events`                 | Discord gateway event listeners for non-interaction events          |
| `internal/disgo/actions`                | Reusable Discord API calls — one exported function per file         |
| `internal/disgo/resolvers`              | Discord data lookups via Services — one exported function per file  |
| `internal/disgo/components`             | UI builder factories (buttons, embeds, rows, selects, modals)       |

## Frontend Structure

- **UI framework:** Skeleton UI v4 (Svelte 5 + Tailwind CSS v4), cerberus theme
- **API client:** PocketBase JS SDK (`pocketbase` npm package) — singleton in `src/lib/pocketbase.ts`; in dev points to `http://localhost:PORT`, in production passes `undefined` (same-origin relative)
- **Auth store:** `src/lib/stores/auth.svelte.ts` — uses Svelte 5 runes (`$state`/`$derived`), not writable stores
- **Mode store:** `src/lib/stores/mode.svelte.ts` — dark/light mode toggle, persisted in `localStorage`; call `mode.toggle()` or `mode.set('dark'|'light')`
- **Toaster:** `src/lib/stores/toaster.ts` — global Skeleton toast singleton (`toaster`); import and call `toaster.trigger(...)` from any component
- **Navigation:** `src/lib/config/navigation.ts` — central nav link config consumed by all four layout nav components; edit here to add/remove nav links
- **App config:** `src/lib/config/app.ts` — exports `APP_NAME` (displayed app name) and `OAUTH_PROVIDERS` (display labels + icons per provider); actual enabled providers are discovered at runtime from PocketBase's `listAuthMethods()` API
- **WebSocket:** Browser native `WebSocket` API connecting to `/api/ws?token=PB_JWT`
- **Routing:** SvelteKit file-based routing in `sveltekit/src/routes/`; `+layout.ts` sets `ssr = false`, `prerender = true`, `trailingSlash = 'always'` globally
- **Build:** adapter-static outputs directly to `pb_public/` with SPA fallback
- **Env:** `vite.config.ts` uses `envDir: '..'` to read from root `.env` — no separate `sveltekit/.env`
- **Package manager:** pnpm

### Responsive layout

The root layout (`+layout.svelte`) implements a 3-mode navigation system driven by a single `NavPanel` component:

| Breakpoint       | Nav mode                                                             |
| ---------------- | -------------------------------------------------------------------- |
| Mobile (`< sm`)  | Bottom bar (`MobileNav`) + slide-in overlay drawer (`NavPanel`)      |
| Desktop (`< lg`) | Rail sidebar — icons only (`NavPanel layout="rail"`)                 |
| Desktop (`≥ lg`) | Toggle between rail and full sidebar via `NavToggle` in the `Header` |

`NavToggle` toggles `navOpen`, which controls both the desktop rail↔sidebar expansion and the mobile overlay open/close state. `NavPanel` derives its Skeleton `layout` prop (`"rail"` | `"sidebar"`) from `open` and `isDesktop`.

## Cross-System Architecture

The three main systems (PocketBase, Disgo, WebSocket) never import each other. Cross-system communication is mediated through:

1. **Interfaces** (`internal/guards/interfaces/`) — one interface per file, organized in per-system subdirectories (`discord/`, `websocket/`, `pocketbase/`). Small interfaces compose into aggregate `Service` interfaces via embedding.
2. **Services struct** (`internal/guards/services.go`) — bundles all system references. Fields are nil if the system is not running.
3. **Dependency injection** — `main.go` builds the `Services` struct and injects it into all three systems via `SetServices()`.

Handler flow: **Trigger → Resolve → Guard → Action**

- **Resolvers** stay in their own package (`pocketbase/resolvers/`, `disgo/resolvers/`, `websocket/resolvers/`) and only talk to their own system
- **Guards** (`internal/guards/`) take `*Services` and check cross-system permissions
- **Actions** are called through `Services` interfaces (e.g., `svc.Discord.SendNotification()`, `svc.WS.BroadcastRaw()`)

Handlers orchestrate by calling resolvers/guards/actions from multiple systems — no resolver or guard calls sideways into another package's resolvers.

## Conventions

- **Go module path** is `github.com/youruser/yourproject` — rename it when starting a new project (see README Quick Start).
- **Adding new routes/hooks/commands/WS handlers:** create a new file in the relevant package, define a function, and call `register(fn)` from `init()`. No other file needs to change — the package-level `init()` runs automatically on import.
- PocketBase v0.36.7 — uses `net/http.ServeMux`, NOT Echo. Hooks use `OnServe` not `OnBeforeServe`.
- PocketBase extensions follow a registration pattern: hooks register before OnServe, schema/routes register inside OnServe via `RegisterAll()`
- One `.go` file per logical domain in `schema/`, `hooks/`, `routes/`, and `commands/`
- PB record hooks use `routine.FireAndForget` for async external calls (Discord API)
- Clone record data into local variables before goroutines — event objects are not concurrent-safe
- WebSocket auth: validate `?token=` query param, attach user if valid, allow anonymous if not
- WebSocket Hub supports: Broadcast (all clients), SendToUser (by user ID), SendToRoom (room members), plus `*Raw` variants taking `[]byte` for cross-system use via interfaces
- Disgo uses `discord.SlashCommandCreate` for slash commands, raw event listeners for gateway events
- Disgo actions take `*bot.Client` as first param — also exposed as methods on `Bot` for interface compliance
- Disgo components are pure builder functions (no registry, no init) — one file per button/embed/row
- Cross-system guards in `internal/guards/` take `*Services` + `*core.Record`, usable from any system — Discord command handlers and WebSocket room joins both call into this package
- Interface files use one-interface-per-file convention for merge-safe parallel development
- Custom routes registered in OnServe take priority over pb_public/ static file serving
- `PUBLIC_PB_PORT` in root `.env` — single port variable used by Taskfile, compose, Containerfile, and SvelteKit (via `$env/static/public`). The `PUBLIC_` prefix is required by SvelteKit for client-side access. The server binary also self-binds to it on startup: if `PUBLIC_PB_PORT` is set and no `--http` flag was passed, `cmd/server/main.go` injects `--http=0.0.0.0:$PUBLIC_PB_PORT` into `os.Args` — so `./bin/server serve` works without explicit flags
- Use `app.Logger()` (slog-based) inside PocketBase request handlers, hooks, and routes — not the stdlib `log` package — so output stays consistent with PocketBase's structured logs
- SvelteKit `trailingSlash = 'always'` is set globally — all route hrefs must end with `/` (e.g. `/login/`, not `/login`), otherwise navigation breaks with the static adapter
- **Seeding:** Air (`task dev`) builds with `-tags dev`, causing `seed.Run(app)` to execute automatically at server startup using `internal/pocketbase/seed/data.go`. Edit `data.go` to change seed data.
- **Dev vs prod builds:** `air` (dev) compiles with `-tags dev`; `task build:backend` compiles without it. The `//go:build dev` constraint in `internal/pocketbase/seed/` means the seeder is a no-op in production binaries.
- **Dev DB is ephemeral:** Air compiles the server to `tmp/server.exe` and `clean_on_exit = true` wipes `tmp/` on exit — including `tmp/pb_data/` where PocketBase stores its dev database. This is intentional: each `task dev` session starts with a clean slate. TypeScript type generation (`task typegen`) therefore uses `--url` mode against the live server rather than reading the DB file directly.
- **`.go.example` scaffolding:** files like `internal/guards/guard.go.example` and `internal/pocketbase/routes/admin/routes.go.example` are templates — copy and rename to `.go` to add a new guard or admin route. The `.example` suffix keeps them out of the build.
- **pnpm pinned to v10:** both `Containerfile` (`corepack prepare pnpm@10`) and `.github/workflows/ci.yml` (`pnpm/action-setup` `version: 10`). Don't bump to `latest` — pnpm v11 raises `ERR_PNPM_IGNORED_BUILDS` during `--frozen-lockfile` installs even when packages are listed in `onlyBuiltDependencies`. The dep build-script allowlist is authoritative in `sveltekit/pnpm-workspace.yaml` (pnpm v10+ ignores a `pnpm` block in `package.json`), and the `Containerfile`'s frontend stage must copy this file into the build context.
