# sveltekit/

SvelteKit frontend source. Uses `@sveltejs/adapter-static` to produce a fully static build served by PocketBase.

## Build Output

`adapter-static` is configured to output to `../pb_public/` (one level up from this directory). PocketBase automatically serves anything in `pb_public/` over HTTP.

In `svelte.config.js`:

```js
import adapter from '@sveltejs/adapter-static';

export default {
  kit: {
    adapter: adapter({
      pages: '../pb_public',
      assets: '../pb_public',
      fallback: 'index.html',  // enables SPA-style client-side routing
    }),
  },
};
```

## UI Framework

[Skeleton UI v4](https://www.skeleton.dev/) is used for components and theming on top of Tailwind CSS.

## API Communication

- **PocketBase** — via the [PocketBase JS SDK](https://github.com/pocketbase/js-sdk) (`pocketbase` npm package)
- **WebSocket** — via the browser's native `WebSocket` API or a thin wrapper in `src/lib/api/`

## Development

```sh
npm install
npm run dev      # dev server (proxies API to running Go backend)
npm run build    # outputs static files to ../pb_public/
```
