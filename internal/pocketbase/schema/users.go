package schema

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func init() {
	register(registerUsersCollection)
}

// registerUsersCollection customizes the built-in "users" auth collection
// that PocketBase creates automatically on first boot.
//
// PocketBase ships the users auth collection with these fields already:
//   - id, password, tokenKey, email, emailVisibility, verified (system/auth)
//   - name (text), avatar (file)
//   - created, updated (timestamps)
//
// Don't redefine those here — only add project-specific fields below.
func registerUsersCollection(app *pocketbase.PocketBase) error {
	// No collectionExists() guard — users always exists, so we reconcile on every boot.
	users, err := requireCollection(app, "users")
	if err != nil {
		return err
	}

	// Guard each field with GetByName so reboots stay idempotent.
	if users.Fields.GetByName("age") == nil {
		users.Fields.Add(&core.NumberField{
			Name:    "age",
			Min:     f64(0),
			Max:     f64(120),
			OnlyInt: true,
		})
	}

	return app.Save(users)
}
