// Package migrations holds the PocketBase database migrations that define this
// project's schema. It is the SOURCE OF TRUTH for collections — see
// docs/MIGRATIONS.md for the workflow.
//
// Files here self-register via init() + m.Register(up, down) and are applied in
// filename (timestamp) order. main.go blank-imports this package so they load;
// pending ones are applied on boot and recorded in the _migrations table.
//
// Author migrations in dev (Automigrate writes them when you change the schema
// in the admin UI), review the generated file, commit it, then let test and prod
// apply it. Never hand-edit an already-applied migration — add a new one.
package migrations
