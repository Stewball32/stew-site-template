# internal/pocketbase/schema

Go files that programmatically define PocketBase collections.

## Responsibilities

Each file in this directory defines one or more PocketBase collections — their fields, types, validation rules, indexes, and API access rules. These are registered with the app before `app.Start()` is called.

PocketBase will create or update collections to match these definitions on startup (when using `automigrate` or explicit collection upserts).

## Expected Files

- One `.go` file per logical domain (e.g., `users.go`, `posts.go`)
- Each file exports a function like `RegisterUsersCollection(app *pocketbase.PocketBase) error`
- A `register.go` that calls all individual registration functions

## Example Pattern

```go
func RegisterPostsCollection(app *pocketbase.PocketBase) error {
    collection := &models.Collection{
        Name:   "posts",
        Type:   models.CollectionTypeBase,
        Schema: schema.NewSchema(
            &schema.SchemaField{Name: "title", Type: schema.FieldTypeText, Required: true},
            &schema.SchemaField{Name: "body",  Type: schema.FieldTypeEditor},
        ),
    }
    return app.Dao().SaveCollection(collection)
}
```
