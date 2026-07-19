# Database migrations

**Migrations are the source of truth for the schema.** Collections are created
and changed by files in [`migrations/`](../migrations), applied on boot and
tracked in PocketBase's `_migrations` table. There is no on-serve
"schema-as-code" step — that was replaced by the baseline migration.

## Why

Schema-as-code (recreating collections in an `OnServe` hook) works for *creating*
collections but has no answer for *changing* one: no ordering, no record of what
ran, no way to know if test and prod are at the same schema, and edits made in the
admin UI silently diverge. Migrations give an ordered, reviewable, idempotent
history that every tier applies identically.

## How it's wired

- `cmd/server/main.go` registers `migratecmd` with `Dir: "migrations"` and
  blank-imports `migrations/` so the files self-register.
- **Automigrate is ON in dev builds only** — `internal/pocketbase/migrateconf`
  (`//go:build dev` → `true`, `!dev` → `false`). `task dev` runs Air with
  `-tags dev`; container/prod builds have no tag.
  - **dev:** change the schema in the admin UI → a migration file is written to
    `migrations/` automatically. That file is what you commit.
  - **test/prod:** never author migrations. They only *apply* what was committed.
- Pending migrations run **before** `OnServe`, so routes/hooks can assume the
  schema exists.

## Workflow

```
dev (Automigrate) ──► commit migration file ──► test applies on boot ──► prod applies on boot
```

1. **Author in dev.** `task dev`, then edit collections in the admin UI
   (`http://localhost:<dev-backend-port>/_/`). A file appears in `migrations/`.
   - Prefer a blank migration for data/logic changes:
     `task migrate:create -- backfill_widget_slugs`
2. **Review it.** Open the generated file. It is plain Go — confirm it does what
   you expect, and that `down` actually reverses `up`.
3. **Commit it** with the code that depends on it, in the same commit.
4. **Prove it on test.** Deploy to the test tier **against a copy of prod
   `pb_data`** (see below) and confirm the app boots and behaves.
5. **Deploy prod.** It applies the same pending migrations on boot.

### Proving a migration against real data (do this before prod)

Never let prod be the first place a migration meets real data.

```sh
# snapshot prod's volume into the test tier's volume, then deploy test
docker run --rm -v <proj>_pb_data:/from -v <proj>-pre_pb_data_pre:/to alpine \
  sh -c 'rm -rf /to/* && cp -a /from/. /to/'
./deploy-pre.sh          # boots test on the copy → migrations apply there first
```

Check the logs for migration errors and hit `/api/health`. If it goes wrong you
throw the copy away; prod was never touched.

## Commands

| Task | What it does |
| --- | --- |
| `task migrate:up` | Apply pending migrations (test/prod also do this on boot) |
| `task migrate:create -- <name>` | New blank migration |
| `task migrate:collections` | Snapshot current collections into a migration (baseline / re-sync) |
| `task migrate:down` | Revert the last migration — **dev only, destructive** |

## Rules

- **Never edit an applied migration.** It's already recorded in `_migrations` on
  some tier. Add a new migration instead.
- **Never hand-write into `_migrations`.**
- **One logical change per migration**, named for what it does.
- **Migrations must be idempotent-safe to re-run in sequence** on a fresh DB —
  that's the guarantee that a new environment can be built from scratch.
- **Commit the migration with the code that needs it**, so a rollback of one is a
  rollback of both.

## The baseline

`migrations/*_collections_snapshot.go` is the **baseline**: a full snapshot of the
collections, generated with `migrate collections` from the schema this template
previously created in code. A fresh `pb_data` is built entirely from it — verified:
a brand-new database produces `posts` (`id, title, body, author`) and `users` with
no schema-as-code present.

## Migrating an existing project off schema-as-code

For a child repo (e.g. `norcal-halo-site`) that still creates collections in an
`OnServe` hook:

1. Wire `migratecmd` + `migrateconf` + `migrations/` as here.
2. Boot the existing app once so schema-as-code creates/reconciles everything.
3. `task migrate:collections` → baseline that matches the live schema exactly.
4. Remove the `schema.RegisterAll(...)` call and retire `internal/pocketbase/schema/`.
5. Verify a **fresh** `pb_data` boots and produces the same collections.
6. On tiers with existing data, the baseline is recorded as applied without
   re-creating anything (collections already match) — confirm on test first.
