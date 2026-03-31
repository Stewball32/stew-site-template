package hooks

import "github.com/pocketbase/pocketbase"

var registry []func(app *pocketbase.PocketBase)

// register adds a hook function to the registry.
// Call this from init() in each domain file.
func register(fn func(app *pocketbase.PocketBase)) {
	registry = append(registry, fn)
}

// RegisterAll wires all record lifecycle hooks.
// Called from cmd/server/main.go.
func RegisterAll(app *pocketbase.PocketBase) {
	for _, fn := range registry {
		fn(app)
	}
}
