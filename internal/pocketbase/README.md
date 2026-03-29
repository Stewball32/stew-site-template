# internal/pocketbase

PocketBase application setup and customization.

## Responsibilities

Initializes and configures the PocketBase `*pocketbase.PocketBase` app instance. Acts as the integration point between PocketBase's lifecycle and the rest of the application.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `schema/` | Go files that define and register PocketBase collections |
| `hooks/` | Event hooks on PocketBase collections and records |
| `routes/` | Custom API routes registered on PocketBase's Echo router |

## Expected Files

- `app.go` — constructs and returns the configured `*pocketbase.PocketBase` instance, wires schema/hooks/routes
