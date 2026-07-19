# Deployments

How this site is deployed. Three tiers, one interface. **Fill in the tables
below when scaffolding the project** — `scripts/scaffold-site.sh` pre-fills the
ports and hostnames for you.

## Tiers

| Tier | Purpose | Runs from | Command | Compose | Env | Bot |
| --- | --- | --- | --- | --- | --- | --- |
| **dev** | live-reload coding | working tree (Air + Vite HMR) | `./run-dev.sh` | — | `.env.dev` | **off** (forced) |
| **test** (`pre`) | the gate before prod | built image, isolated | `./deploy-pre.sh` | `compose.pre.yml` | `.env.pre` | optional (own token) |
| **prod** | the live site | built image | `./deploy-prod.sh` | `compose.yml` | `.env` | optional (own token) |

## This project — FILL IN

| Tier | Host port | Public hostname | Branch |
| --- | --- | --- | --- |
| dev (vite) | `<DEV_VITE_PORT>` | `<DEV_HOST>` | any |
| dev (backend) | `<DEV_PB_PORT>` | _internal — proxied by Vite_ | any |
| test | `<PRE_PB_PORT>` | `<TEST_HOST>` | `pre` |
| prod | `<PROD_PB_PORT>` | `<PROD_HOST>` | `main` |

- Ports come from this project's block in [`../PORTS.md`](../PORTS.md).
- Every tier binds **loopback only**; the cloudflared tunnel is the only way in
  (see [`../deploy/cloudflared-ingress.snippet.yml`](../deploy/cloudflared-ingress.snippet.yml)).
- Data: prod uses the `pb_data` volume, test uses `pb_data_pre` — **separate**.
  Dev's database is ephemeral (`tmp/dev_pb_data`, wiped on exit).

## Promotion path

```
dev (working tree)  ──►  test (pre)  ──►  prod
   run-dev.sh            deploy-pre.sh     deploy-prod.sh
```

Nothing reaches prod without passing through test. For anything touching the
schema, test must run against a **copy of prod `pb_data`** first — see
[MIGRATIONS.md](MIGRATIONS.md).

## Deploying

Both deploy scripts share one interface:

```sh
./deploy-pre.sh           # guard → build → deploy → health-check
./deploy-pre.sh logs      # follow logs
./deploy-pre.sh down      # stop (keeps the database volume)

./deploy-prod.sh          # same, for prod
```

**Guards.** Each refuses to run unless you're on the tier's branch with a clean
tree — a deploy builds from the working tree, so a dirty tree ships something
that isn't committed and can't be reproduced or rolled back to. Override
deliberately:

```sh
ALLOW_DIRTY=1 ./deploy-pre.sh          # uncommitted changes
ALLOW_ANY_BRANCH=1 ./deploy-pre.sh     # different branch
PRE_BRANCH=main ./deploy-pre.sh        # this project promotes off main
```

**What a deploy does:** builds the image → `up -d` (recreate, so a changed
`env_file` is re-read — `restart` would not) → applies pending migrations on
boot → polls `/api/health` until healthy, dumping logs and failing if not.

## First-time setup

1. Claim a port block in [`../PORTS.md`](../PORTS.md) and record it there.
2. `cp .env.example .env` and `cp .env.pre.example .env.pre`; set ports + secrets.
   (`.env.dev` optional — see `.env.dev.example`.)
3. Give each tier its **own** `SEED_SUPERUSER_*`; never share prod's.
4. If the bot runs: one tier only, its own Discord application per tier (a token
   allows a single gateway connection).
5. Merge the ingress rules from `deploy/` and create the DNS records.
6. `./deploy-pre.sh`, verify, then `./deploy-prod.sh`.

## Rollback

```sh
git checkout <last-good-tag-or-sha>
./deploy-prod.sh            # rebuilds and redeploys that revision
```

⚠️ Code rolls back; **migrations do not**. A migration already applied stays
applied — that's why they're proven on test first. To undo a schema change,
write a new forward migration.

## Backups

`pb_data` is the whole database. Before any risky deploy:

```sh
docker run --rm -v <project>_pb_data:/data -v "$PWD":/backup alpine \
  tar czf /backup/pb_data-$(date +%F).tar.gz -C /data .
```

_TODO: record where backups are stored and the retention policy._
