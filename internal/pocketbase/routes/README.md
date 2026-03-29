# internal/pocketbase/routes

Custom API routes registered on PocketBase's Echo router.

## Responsibilities

PocketBase exposes its underlying [Echo](https://echo.labstack.com/) router via `app.OnBeforeServe()`. This package registers additional HTTP endpoints alongside PocketBase's built-in REST API.

This is also where the WebSocket upgrade endpoint is mounted (delegating to `internal/websocket`).

## Expected Files

- `routes.go` — registers all custom routes via `app.OnBeforeServe()`
- Additional handler files per route group (e.g., `ws.go` for WebSocket mount)

## Example Pattern

```go
func RegisterRoutes(app *pocketbase.PocketBase, wsHandler http.Handler) {
    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        e.Router.GET("/api/custom/hello", func(c echo.Context) error {
            return c.JSON(200, map[string]string{"message": "hello"})
        })
        // Mount WebSocket handler
        e.Router.GET("/ws", echo.WrapHandler(wsHandler))
        return nil
    })
}
```
