import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

// DEV tier wiring — active only when DEV_ALLOWED_HOST is set (by ../run-dev.sh).
//
// It exposes the Vite dev server on the public dev hostname behind the
// cloudflared tunnel and proxies /api to the Go/PocketBase backend, so the dev
// frontend talks to a real backend same-origin.
//
// When the var is unset (plain `pnpm dev`, `pnpm build`, and every container
// build) none of this applies — `server` config only affects `vite dev`, never
// `vite build` — so the prod/test images are unaffected.
const devHost = process.env.DEV_ALLOWED_HOST;
const devVitePort = process.env.DEV_VITE_PORT ? Number(process.env.DEV_VITE_PORT) : undefined;
// cmd/server binds PUBLIC_PB_PORT; run-dev.sh sets it to DEV_PB_PORT.
const backendPort = process.env.PUBLIC_PB_PORT || '8090';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	envDir: '..',
	...(devHost
		? {
				server: {
					host: '127.0.0.1',
					port: devVitePort,
					strictPort: true,
					// Vite rejects unknown Host headers; allow the tunnel's hostname.
					allowedHosts: [devHost],
					// HMR rides back out through the tunnel over TLS on 443.
					hmr: { protocol: 'wss', host: devHost, clientPort: 443 },
					// Same-origin /api → the local backend (ws:true covers /api/ws).
					proxy: {
						'/api': {
							target: `http://127.0.0.1:${backendPort}`,
							changeOrigin: true,
							ws: true
						}
					}
				}
			}
		: {})
});
