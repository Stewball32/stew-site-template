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
│ Go Binary (cmd/server) │
│ │
│ ┌──────────────┐ ┌─────────────────────┐ │
│ │ PocketBase │ │ Disgo Bot │ │
│ │ - REST API │ │ - Slash Commands │ │
│ │ - Auth/JWT │ │ - Event Listeners │ │
│ │ - SQLite │ └─────────────────────┘ │
│ │ - ServeMux │ │
│ │ Router │ ┌─────────────────────┐ │
│ └──────┬───────┘ │ WebSocket │ │
│ │ │ (coder/websocket) │ │
│ │ │ - Optional JWT │ │
│ │ │ - Hub / Rooms │ │
│ │ └─────────────────────┘ │
│ │ │
│ guards/Services → cross-system DI │
│ PB Hooks → Discord notifications │
│ PB Hooks → WS Hub broadcasts │
│ PB Routes → Auth-gated page serving │
│ │
└─────────┬────────────────────────────────────┘
│ serves
┌─────────▼───────┐
│ pb_public/ │ ← SvelteKit static build
│ (SvelteKit) │
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
│   │   ├── actions/           # Reusable PB data operations
│   │   └── resolvers/         # PB data lookups (one function per file)
│   ├── disgo/
│   │   ├── bot.go             # Bot client + interface methods + lifecycle
│   │   ├── commands/          # Slash command definitions and handlers
│   │   ├── events/            # Discord gateway event listeners
│   │   ├── actions/           # Reusable Discord API calls
│   │   ├── resolvers/         # Discord data lookups via Services
│   │   ├── components/        # UI builders (buttons, embeds, rows)
│   │   └── guards/            # Bot-side permission checks
│   └── websocket/
│       ├── hub.go             # Client registry, rooms, message routing
│       ├── handler.go         # WS upgrade with optional JWT auth
│       ├── client.go          # Single connection read/write pumps
│       ├── message.go         # Wire format for WS messages
│       ├── handlers/          # Self-registering message type handlers
│       ├── rooms/             # Room type definitions with guard lists
│       └── resolvers/         # WS state lookups via Services
├── sveltekit/                 # SvelteKit frontend (scaffold separately)
├── .env.example               # Backend env template
├── .air.toml                  # Go hot reload config
├── .gitignore
├── Taskfile.yml               # Build orchestration
├── Containerfile              # Multi-stage Podman/Docker build
├── compose.yml                # Container compose config
├── go.mod
└── LICENSE
```

## Prerequisites

- **Go 1.25+**
- **pnpm** — `npm install -g pnpm`
- **Task** — `go install github.com/go-task/task/v3/cmd/task@latest`
- **Air** (for dev) — `go install github.com/air-verse/air@latest`
- **Podman** (optional, for containers)

## Quick Start

1. **Clone and rename the module:**
   ```bash
   git clone https://github.com/Stewball32/stew-site-template.git my-project
   cd my-project
   grep -rl 'github.com/youruser/yourproject' . | xargs sed -i 's|github.com/youruser/yourproject|github.com/you/yourapp|g'
   ```

2. **Configure environment:**
   ```bash
   cp .env.example .env

   # Edit .env with your Discord bot token, guild ID, etc.

   ```

3. **Scaffold the frontend** (see `sveltekit/README.md` for instructions)

4. **Run in development:**
   ```bash
   task dev
   ```

5. **Build for production:**
   ```bash
   task build
   ./bin/server serve
   ```

## Containers

```bash

# Build image

task container:build

# Run (PB_PORT defaults to 8090, set in .env)

task container:run
```

For multiple instances on the same machine, set a unique `PB_PORT` in each project's `.env`. Route traffic with cloudflared or a reverse proxy.

## Notes

- `pb_data/` — PocketBase runtime data (SQLite DB, uploads). Created at runtime, gitignored.
- `pb_public/` — SvelteKit build output. Created by `task build:frontend`, gitignored.
- Schema can be managed via PocketBase admin UI or programmatically in `internal/pocketbase/schema/`.
- Protected pages are served through auth-gated custom routes; public pages are served directly from `pb_public/`.
