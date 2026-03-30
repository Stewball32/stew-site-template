package schema

// TODO: Schema registration entry point.
// RegisterAll(app) calls each individual collection registration function.
// Called from cmd/server/main.go inside the OnServe hook.
// Each domain gets its own file (e.g., users.go, posts.go) with a Register*Collection(app) function.
