# Changelog

All notable changes to this project are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

- CI: require pnpm **v11** (was v10) — `Containerfile` (`corepack prepare pnpm@11`) and GitHub Actions (`pnpm/action-setup` `version: 11`); Node 22+ required.
- CI: migrate the dependency build-script allowlist in `sveltekit/pnpm-workspace.yaml` from v10's `onlyBuiltDependencies` list to v11's `allowBuilds` map (`esbuild: true`, `sqlite3: true`). The `Containerfile` frontend stage still copies the workspace file into the build context. `pnpm-lock.yaml` is unchanged — v11 reads the existing `lockfileVersion: '9.0'` file as-is.

### Deprecated

### Removed

### Fixed

### Security
