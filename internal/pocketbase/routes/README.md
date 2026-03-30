# internal/pocketbase/routes

Custom API routes registered on PocketBase's ServeMux router.

## Responsibilities

PocketBase exposes its underlying `net/http.ServeMux` router via `app.OnServe()`. This package registers additional HTTP endpoints alongside PocketBase's built-in REST API.

This is also where the WebSocket upgrade endpoint is mounted (delegating to `internal/websocket`).

## Expected Files

- `routes.go` — registers all custom routes via `app.OnServe()`
- Additional handler files per route group (e.g., `ws.go` for WebSocket mount)

## Example Pattern

```go
func RegisterRoutes(app *pocketbase.PocketBase, wsHandler http.Handler) {
    app.OnServe().BindFunc(func(se *core.ServeEvent) error {
        se.Router.GET("/api/custom/hello", func(e *core.RequestEvent) error {
            return e.JSON(200, map[string]string{"message": "hello"})
        })
        // Mount WebSocket handler
        se.Router.GET("/api/ws", wsHandler)
        return se.Next()
    })
}
```
