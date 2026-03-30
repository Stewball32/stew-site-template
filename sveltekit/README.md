# SvelteKit Frontend

This directory will contain the SvelteKit project. Scaffold it with:

```bash
cd sveltekit
pnpm create svelte@latest . --template skeleton --types ts
pnpm install
pnpm add -D @skeletonlabs/skeleton @skeletonlabs/skeleton-svelte
pnpm add pocketbase
```

## Post-scaffold setup

- **svelte.config.js**: Switch to `@sveltejs/adapter-static` with `fallback: 'index.html'`
- **vite.config.ts**: Add proxy for `/api` and `/_` to `http://localhost:8090`
- **src/app.css**: Import Tailwind, Skeleton, and a theme:
  ```css
  @import 'tailwindcss';
  @import '@skeletonlabs/skeleton';
  @import '@skeletonlabs/skeleton-svelte';
  @import '@skeletonlabs/skeleton/themes/cerberus';
  ```
- **src/app.html**: Add `data-theme="cerberus"` to the `<html>` tag
- **src/routes/+layout.ts**: Set `ssr = false`, `prerender = true`, `trailingSlash = 'always'`
- **src/lib/pocketbase.ts**: Create PB client singleton using `PUBLIC_POCKETBASE_URL`
- **sveltekit/.env**: Set `PUBLIC_POCKETBASE_URL=http://localhost:8090`
