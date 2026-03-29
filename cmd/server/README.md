# cmd/server

Go entry point for the application. Contains `main.go`.

## Responsibilities

This is where all components are wired together before `app.Start()` is called:

1. Create the PocketBase app instance
2. Register collection schemas (`internal/pocketbase/schema`)
3. Register PocketBase hooks (`internal/pocketbase/hooks`)
4. Register custom API routes and WebSocket handler on PocketBase's Echo router (`internal/pocketbase/routes`, `internal/websocket`)
5. Start the Disgo Discord bot in a goroutine (`internal/disgo`)
6. Call `app.Start()` — blocking, runs PocketBase's HTTP server

## Expected Files

- `main.go` — wires and starts the application
- `config.go` (optional) — loads environment variables / config struct
