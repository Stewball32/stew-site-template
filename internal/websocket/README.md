# internal/websocket

WebSocket server using [coder/websocket](https://github.com/coder/websocket).

## Responsibilities

Provides a WebSocket endpoint mounted on PocketBase's ServeMux router via `internal/pocketbase/routes`. Handles connection upgrades with optional JWT authentication, and provides a Hub for managing connected clients, rooms, and message routing.

## Why coder/websocket

`coder/websocket` is a lightweight, stdlib-compatible WebSocket library. It has native `context.Context` support, safe concurrent writes without external mutexes, zero external dependencies, and works directly with `net/http` handlers. It is the actively maintained successor to `nhooyr/websocket`.

## Auth Flow

1. Browser connects: `new WebSocket("ws://host/api/ws?token=PB_JWT")`
2. Handler checks for `?token=` query parameter
3. If present, validates with `app.FindAuthRecordByToken(token, core.TokenTypeAuth)`
4. Valid token → Client is tagged with the authenticated user record
5. Invalid or missing token → connection stays open as anonymous
6. Connection is upgraded with `websocket.Accept()`
7. Client is registered with the Hub

## Expected Files

- `handler.go` — HTTP upgrade endpoint. Validates optional JWT auth via `?token=` query param using PocketBase's `app.FindAuthRecordByToken()`. Creates a Client and registers it with the Hub.
- `hub.go` — Central Hub struct managing connected clients and rooms. Processes register/unregister channels and routes messages. Methods: `Broadcast()`, `SendToUser()`, `SendToRoom()`, `JoinRoom()`, `LeaveRoom()`.
- `client.go` — Represents a single WebSocket connection. `readPump()` reads messages from the browser and forwards to the Hub. `writePump()` sends Hub messages to the browser. Tracks authenticated user ID (or nil for anonymous).
- `message.go` — Wire format for WebSocket messages. Type field for routing (e.g., `"broadcast"`, `"room"`, `"direct"`, `"join_room"`, `"leave_room"`). Optional Room and Target (user ID) fields. Payload as `json.RawMessage` for project-specific data.

## Message Routing

```
Browser → readPump → Hub routes by message.Type
├── Broadcast() → all connected clients
├── SendToRoom() → clients in a specific room
├── SendToUser() → specific user's connections
└── custom → project-specific logic
```

PocketBase hooks and Disgo event handlers can also call Hub methods directly since they share the same Go process.

## Example Usage

```go
// In internal/pocketbase/routes — mount the WS endpoint
hub := websocket.NewHub()
go hub.Run()

se.Router.GET("/api/ws", websocket.NewHandler(hub, e.App))
```

```go
// In internal/pocketbase/hooks — push to WS clients from a record hook
app.OnRecordAfterCreateSuccess("posts").BindFunc(func(e *core.RecordEvent) error {
hub.Broadcast(websocket.Message{
Type: "new_post",
Payload: json.RawMessage(`{"id":"` + e.Record.Id + `"}`),
})
return e.Next()
})
```
