# internal/pocketbase/middleware

Custom middleware for PocketBase routes.

## Responsibilities

Provides reusable middleware functions for custom routes registered in `internal/pocketbase/routes`. Middleware runs before route handlers to enforce auth, roles, or other request-level checks.

## Expected Files

- `auth.go` — JWT validation middleware. Checks Authorization header or `pb_auth` cookie. Can use `apis.RequireAuth()` for simple gating or custom logic for finer control.
- `roles.go` — Role-based access. Checks user record fields (e.g., admin, member) after auth validation.

## Usage

Middleware is bound to routes in `internal/pocketbase/routes` using PocketBase's `Bind()` method:

```go
se.Router.GET("/api/private", handler).Bind(customAuthMiddleware)
```

PocketBase also provides built-in middleware:

- `apis.RequireAuth()` — rejects unauthenticated requests
- `apis.RequireGuestOnly()` — rejects authenticated requests
- `apis.RequireSuperuserAuth()` — restricts to superusers
