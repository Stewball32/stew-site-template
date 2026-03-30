package hooks

// TODO: Hook registration entry point.
// RegisterAll(app) wires all record lifecycle hooks.
// Called from cmd/server/main.go inside the OnServe hook.
// Each domain gets its own file (e.g., users.go, posts.go) with a Register*Hooks(app) function.
