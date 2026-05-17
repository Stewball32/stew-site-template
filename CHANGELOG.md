# Changelog

All notable changes to this project are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

- CI: pin pnpm to v10 in Containerfile and GitHub Actions. pnpm v11 (pulled by `latest`) raises `ERR_PNPM_IGNORED_BUILDS` during `--frozen-lockfile` installs even when packages are listed in `onlyBuiltDependencies`.
- CI: consolidate dependency build-script approvals in `sveltekit/pnpm-workspace.yaml` (esbuild + sqlite3) and remove the now-redundant `pnpm.onlyBuiltDependencies` block from `package.json`. pnpm v10+ treats the workspace file as authoritative.
- CI: copy `pnpm-workspace.yaml` into the container build context so the approval list is visible during the frontend stage.

### Security
