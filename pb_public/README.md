# pb_public/

SvelteKit static build output. Served automatically by PocketBase over HTTP.

## How It Gets Here

Running `npm run build` inside `sveltekit/` writes the compiled static site to this directory (configured via `adapter-static` in `svelte.config.js`).

## What PocketBase Does

PocketBase serves the contents of `pb_public/` as a static file server on its HTTP port. Any request that doesn't match a PocketBase API route falls through to this directory, enabling SPA-style client-side routing.

## Git

This directory's build artifacts are gitignored. Only this `README.md` is tracked.
To generate the contents, run `npm run build` from `sveltekit/`.
