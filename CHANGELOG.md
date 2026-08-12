# Changelog

All notable changes to this project are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **stew-kit** motion & theme kit for the frontend (`sveltekit/docs/stew-kit.md`):
  - Two new Skeleton v5 themes in `src/lib/themes/` — **stew** (teal/ocean/mint over cool slate, now the default) and **custom** (an editable cerberus clone). All three themes (incl. cerberus) stay registered and switchable at runtime via `data-theme`.
  - `src/lib/styles/motion.css` — Tailwind v4 motion/background utilities (`hover-lift`, `press`, `icon-slide`, `glow`, `shimmer`, `bg-aurora`, `bg-dotgrid`), all `prefers-reduced-motion` safe.
  - `src/lib/components/fx/` — 14 theme-token-driven components: ParticleField, MouseGlow, AnimateIn, CountUp, TypeCycle, TiltCard, Sparkline, SkeletonBlock, EmptyState, ConfettiBurst, PageTransition, ScrollProgress, and BuildStamp.
  - Build stamp + stale-tab detection: `kit.version` is set to the short git commit (60s poll); hovering the app logo shows the running commit, and an amber dot + Reload appears when a newer deploy is live.

### Changed

- Default theme is now **stew** (`data-theme="stew"` in `app.html`); previously cerberus.
- Layout: route navigations cross-fade via the View Transitions API (`PageTransition`), and a top progress bar tracks the main scroll container (`ScrollProgress`).
- Home page: animated hero (particles, pointer glow, `TypeCycle` headline), staggered feature cards, and corrected "Skeleton UI v4" copy to v5 (also in the docs example, READMEs, and CLAUDE.md).
- Examples polish: dashboard stats animate with `CountUp` + `Sparkline`; pricing cards tilt in 3D and re-count on the billing-period toggle; billing tiles count up; team/inbox/data-table/files/calendar zero-data states use `EmptyState`; wizard submit and inbox "mark all read" celebrate with `ConfettiBurst`; cards site-wide gained `hover-lift`.
- Auth pages: OAuth provider lists and connected-accounts panels show `SkeletonBlock` shimmer placeholders while loading (fixes `(user)/profile` showing "No connected accounts" during the fetch).
- CI: require pnpm **v11** (was v10) — `Containerfile` (`corepack prepare pnpm@11`) and GitHub Actions (`pnpm/action-setup` `version: 11`); Node 22+ required.
- CI: migrate the dependency build-script allowlist in `sveltekit/pnpm-workspace.yaml` from v10's `onlyBuiltDependencies` list to v11's `allowBuilds` map (`esbuild: true`, `sqlite3: true`). The `Containerfile` frontend stage still copies the workspace file into the build context. `pnpm-lock.yaml` is unchanged — v11 reads the existing `lockfileVersion: '9.0'` file as-is.
- The dev tier runs no Discord bot by construction rather than by runtime override. The `dev` build tag sets `internal/disgo.SubsystemEnabled = false`, and `cmd/server` skips bot startup entirely when it is false — replacing `run-dev.sh`'s `export DISCORD_BOT_TOKEN=""`. `.env.dev` needs no bot token, and `test` (pre) is the lowest tier that runs a bot. See `docs/DEPLOYMENTS.md`.

### Deprecated

### Removed

### Fixed

- Wizard example crashed on mount (blank page) since the Skeleton v5 migration: `Steps.Separator` consumes the item context in v5 and must be rendered inside `Steps.Item`, not as a sibling.

### Security
