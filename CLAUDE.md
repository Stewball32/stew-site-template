# CLAUDE.md

> See also: [README.md](README.md) for project overview, tech stack, architecture diagram, and quick-start guide.

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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
4. Register OnTerminate hook to shut down Disgo bot cleanly
5. Call `app.Start()` (blocking)

### Key packages

| Package                          | Purpose                                                                |
| -------------------------------- | ---------------------------------------------------------------------- |
| `internal/pocketbase/schema`     | Programmatic collection definitions — one file per domain              |
| `internal/pocketbase/hooks`      | Record event hooks — fire Discord notifications, push to WS clients    |
| `internal/pocketbase/oauth`      | OAuth2 provider config — env-driven, self-registering, one per file    |
| `internal/pocketbase/routes`     | Custom endpoints + protected page serving via auth-gated routes        |
| `internal/pocketbase/routes/middleware` | Auth middleware, role checks, global middleware registry         |
| `internal/pocketbase/routes/admin`     | Route group for `/api/admin` — auth + admin middleware           |
| `internal/websocket`             | WebSocket Hub, client management, message routing, JWT auth upgrade    |
| `internal/disgo`                 | Bot client setup: NewBot(), OpenGateway(), Close()                     |
| `internal/disgo/commands`        | Slash command definitions and handler functions                        |
| `internal/disgo/events`          | Discord gateway event listeners for non-interaction events             |

## Frontend Structure

- **UI framework:** Skeleton UI v4 (Svelte 5 + Tailwind CSS v4)
- **API client:** PocketBase JS SDK (`pocketbase` npm package) for REST/auth
- **WebSocket:** Browser native `WebSocket` API connecting to `/api/ws?token=PB_JWT`
- **Routing:** SvelteKit file-based routing in `sveltekit/src/routes/`
- **Package manager:** pnpm

## Conventions

- PocketBase v0.36.7 — uses `net/http.ServeMux`, NOT Echo. Hooks use `OnServe` not `OnBeforeServe`.
- PocketBase extensions follow a registration pattern: hooks register before OnServe, schema/routes register inside OnServe via `RegisterAll()`
- One `.go` file per logical domain in `schema/`, `hooks/`, `routes/`, and `commands/`
- PB record hooks use `routine.FireAndForget` for async external calls (Discord API)
- Clone record data into local variables before goroutines — event objects are not concurrent-safe
- WebSocket auth: validate `?token=` query param, attach user if valid, allow anonymous if not
- WebSocket Hub supports: Broadcast (all clients), SendToUser (by user ID), SendToRoom (room members)
- Disgo uses `discord.SlashCommandCreate` for slash commands, raw event listeners for gateway events
- Custom routes registered in OnServe take priority over pb_public/ static file serving
