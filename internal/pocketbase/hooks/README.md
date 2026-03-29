# internal/pocketbase/hooks

PocketBase collection and record event hooks.

## Responsibilities

Hooks run custom Go logic in response to PocketBase record lifecycle events. They are registered on the app before `app.Start()`.

## Common Hook Events

| Event | Trigger |
|-------|---------|
| `OnRecordBeforeCreateRequest` | Before a record is created via API |
| `OnRecordAfterCreateRequest` | After a record is created via API |
| `OnRecordBeforeUpdateRequest` | Before a record is updated via API |
| `OnRecordAfterUpdateRequest` | After a record is updated via API |
| `OnRecordBeforeDeleteRequest` | Before a record is deleted via API |
| `OnRecordAfterDeleteRequest` | After a record is deleted via API |

## Expected Files

- One `.go` file per collection or domain (e.g., `users.go`, `posts.go`)
- Each exports a `Register*Hooks(app *pocketbase.PocketBase)` function
- A `register.go` that wires all hook registrations

## Example Pattern

```go
func RegisterPostsHooks(app *pocketbase.PocketBase) {
    app.OnRecordAfterCreateRequest("posts").Add(func(e *core.RecordCreateEvent) error {
        // e.g., send a notification, update a counter
        return nil
    })
}
```
