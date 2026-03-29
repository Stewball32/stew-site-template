# internal/websocket

WebSocket server using [coder/websocket](https://github.com/coder/websocket).

## Responsibilities

Provides an `http.Handler` that upgrades HTTP connections to WebSocket and manages the message loop. This handler is mounted on PocketBase's Echo router via `internal/pocketbase/routes`.

## Why coder/websocket

`coder/websocket` is a lightweight, stdlib-compatible WebSocket library. It works with standard `net/http` handlers and supports both server and client usage without requiring a custom server.

## Expected Files

- `handler.go` — exports `NewHandler() http.Handler`; handles upgrade, read/write loop, and connection cleanup
- `hub.go` (optional) — connection registry for broadcasting messages to multiple clients

## Example Pattern

```go
func NewHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        conn, err := websocket.Accept(w, r, nil)
        if err != nil {
            return
        }
        defer conn.CloseNow()
        // read/write loop
    })
}
```
