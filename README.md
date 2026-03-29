# stew-site-template

A reusable project template combining a Go backend with a SvelteKit frontend.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend / Auth / DB | [PocketBase](https://pocketbase.io/) (Go framework) |
| Discord Bot | [Disgo](https://github.com/disgoorg/disgo) |
| WebSocket Server | [coder/websocket](https://github.com/coder/websocket) |
| Frontend | [SvelteKit](https://kit.svelte.dev/) + [Skeleton UI v4](https://www.skeleton.dev/) |
| Frontend Build | `@sveltejs/adapter-static` → served by PocketBase |

## Architecture Overview

```
┌─────────────────────────────────────────┐
│  Go Binary (cmd/server)                 │
│                                         │
│  ┌─────────────┐  ┌──────────────────┐  │
│  │  PocketBase │  │  Disgo Bot       │  │
│  │  - REST API │  │  - Commands      │  │
│  │  - Auth     │  │  - Events        │  │
│  │  - SQLite   │  └──────────────────┘  │
│  │  - Static   │                        │
│  │    Files    │  ┌──────────────────┐  │
│  │  - Echo     │  │  WebSocket       │  │
│  │    Router   │  │  (coder/ws)      │  │
│  └─────────────┘  └──────────────────┘  │
└─────────────────────────────────────────┘
         ↑ serves
┌─────────────────┐
│  pb_public/     │  ← SvelteKit static build
│  (SvelteKit)    │
└─────────────────┘
```

## Project Structure

```
.
├── cmd/server/          # Go entry point
├── internal/
│   ├── pocketbase/      # PocketBase app, schema, hooks, routes
│   ├── disgo/           # Discord bot commands and events
│   └── websocket/       # WebSocket handler
├── sveltekit/           # SvelteKit frontend source
├── pb_public/           # SvelteKit build output (served by PocketBase)
└── pb_data/             # PocketBase runtime data (gitignored)
```

## Development

```sh
# Backend
go run ./cmd/server serve

# Frontend (in sveltekit/)
npm install
npm run dev

# Build frontend for production (outputs to pb_public/)
npm run build
```
