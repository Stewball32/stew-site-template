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

## The model: granular in dev, squashed on release

Two different needs, so two different shapes of history:

| | dev / test branches | `main` → prod |
| --- | --- | --- |
| Migrations | **many small files**, one per change | **one file per approved release** |
| History | messy and iterative — expected | a clean version log |
| Source | Automigrate + `migrate create` | `migrate:squash` |

**Branch discipline:** the granular churn lives on dev/test branches only. `main`
carries nothing but squashed release migrations. **The squash happens on the way
to main** — never merge granular files into it.

So prod's `_migrations` table reads like a release history:

```
1784482228_collections_snapshot.go   ← baseline
1784483479_release_v0_3_0.go         ← everything approved in v0.3.0, as one file
1784490000_release_v0_4_0.go         ← …v0.4.0
```

## Workflow

```
dev (Automigrate, many small files)
  └─► test applies them on boot          ← iterate freely here
        └─► APPROVED → task migrate:squash   (granular ──► one release migration)
              └─► task migrate:verify        (proves squashed ≡ granular)
                    └─► merge to main ──► prod applies ONE clean migration
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

### Squashing for a release

When the branch is approved, collapse its granular migrations into one release
migration and prove the collapse is faithful:

```sh
# 1. squash — applies every current migration to a throwaway DB, snapshots the
#    resulting schema as ONE release file, archives the granular ones
task migrate:squash -- --version v0.3.0

# 2. verify — THIS IS THE PROVE-BEFORE-PROD CHECK. Run it against a COPY of
#    prod's pb_data so you verify the exact upgrade path prod will take.
task migrate:verify -- --from /tmp/prod-pb_data-copy

# 3. only if verify passes:
git add -A && git commit -m "chore(release): squash migrations for v0.3.0"
```

**What verify actually does.** From the *same* starting database it builds two
independent trees and compares the end schemas byte-for-byte:

- **A (granular)** — starting data + release-line files + the archived granular files
- **B (squashed)** — starting data + release-line files + the new release migration

Identical → the squash is equivalent and safe to merge. Different → it prints the
diff and refuses (exit 1).

> **The one case that legitimately fails:** if the release **deletes** a
> collection. A generated snapshot imports with `deleteMissing=false`, so it
> creates and updates but never removes — the squashed schema keeps the deleted
> collection and verify catches the divergence. Re-run the squash with
> `--delete-missing` to make the release migration authoritative (it will drop
> collections absent from the snapshot — destructive, so it is opt-in), then
> verify again.

### Getting a copy of prod's data

Never let prod be the first place a migration meets real data.

```sh
# copy prod's volume into the test tier's volume, then deploy test
docker run --rm -v <proj>_pb_data:/from -v <proj>-pre_pb_data_pre:/to alpine \
  sh -c 'rm -rf /to/* && cp -a /from/. /to/'
./deploy-pre.sh          # boots test on the copy → migrations apply there first

# or extract it to a directory for `migrate:verify --from`
docker run --rm -v <proj>_pb_data:/from -v /tmp:/out alpine \
  sh -c 'rm -rf /out/prod-pb_data-copy && cp -a /from /out/prod-pb_data-copy'
```

Check the logs for migration errors and hit `/api/health`. If it goes wrong you
throw the copy away; prod was never touched.

### Why a snapshot, not a hand-written diff

The release migration is a full collections snapshot rather than a computed
"diff since last release". PocketBase generates snapshots reliably, and applying
one reproduces the exact end state whatever path got you there — which is the
property that matters. `migrate:verify` is what turns that into a guarantee
rather than an assumption.

## Commands

| Task | What it does |
| --- | --- |
| `task migrate:up` | Apply pending migrations (test/prod also do this on boot) |
| `task migrate:create -- <name>` | New blank migration |
| `task migrate:collections` | Snapshot current collections into a migration (baseline / re-sync) |
| `task migrate:down` | Revert the last migration — **dev only, destructive** |
| `task migrate:squash -- --version vX.Y.Z` | Collapse this branch's granular migrations into one release migration |
| `task migrate:verify -- --from <pb_data>` | Prove the squash produces the identical schema — **run before merging** |

Migration files are classified by filename:

- **release-line** — `*_collections_snapshot.go` (the baseline) and `*_release_*.go`.
  These are what `main`/prod carry.
- **granular** — everything else in `migrations/`. Dev/test only; squashed away
  before merge.

`migrate:squash` archives the granular files to `.migrate-archive/<version>/`
(gitignored) so `migrate:verify` can rebuild the "before" side to compare against.

## Rules

- **Never edit an applied migration.** It's already recorded in `_migrations` on
  some tier. Add a new migration instead.
- **Never hand-write into `_migrations`.**
- **Never merge granular migrations to `main`.** Squash first — that's the whole
  point of the release-migration history.
- **Never merge a squash that hasn't passed `migrate:verify`.**
- **Squashing is only safe because prod never applied the granular files.** They
  exist solely on dev/test branches, so replacing them with one release migration
  loses nothing prod knew about. (Test tiers that *did* apply them are rebuilt
  from a prod copy anyway.)
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
