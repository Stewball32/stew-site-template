# stew-kit — motion & theme kit

Reference for the theme + motion + fx layer of the Skeleton v5 / Tailwind v4 / Svelte 5 frontend.

## Where files live

| File(s)                          | Purpose                                         |
| -------------------------------- | ----------------------------------------------- |
| `src/lib/themes/custom.css`      | editable cerberus clone (`data-theme="custom"`) |
| `src/lib/themes/stew.css`        | the stew theme (`data-theme="stew"`, default)   |
| `src/lib/styles/motion.css`      | motion/background utility classes               |
| `src/lib/components/fx/*.svelte` | 14 fx components (listed below)                 |

## Wiring

`src/routes/layout.css` — after the skeleton imports:

```css
@import '@skeletonlabs/skeleton/themes/cerberus'; /* keep or remove */
@import '../lib/themes/custom.css';
@import '../lib/themes/stew.css';
@import '../lib/styles/motion.css';
```

`src/app.html` — pick the active theme:

```html
<html lang="en" data-theme="stew"></html>
```

Themes are additive: registering all three costs nothing, and you can switch
at runtime with `document.documentElement.setAttribute('data-theme', 'stew')`.

- **custom** is a pixel-identical cerberus clone — edit its tokens freely;
  the packaged cerberus stays pristine. Generator: <https://themes.skeleton.dev>
- **stew** is a distinct blue-green look (teal / ocean / mint on cool slate).
  Your existing dark-mode gradient in `layout.css` already reads
  `--color-primary-500` / `--color-tertiary-500`, so it re-tints automatically.

## Motion utilities (`motion.css`)

Compose like any Tailwind/Skeleton class; all honor `prefers-reduced-motion`.

```html
<div class="hover-lift card p-6">…</div>
<button class="btn press icon-slide preset-filled">Get Started <ArrowRightIcon /></button>
<a class="btn preset-outlined-primary-500 glow" href="…">Docs</a>

<!-- loading state: shimmer + Skeleton placeholders -->
<div class="shimmer space-y-3 card p-6">
	<div class="placeholder-circle w-10"></div>
	<div class="placeholder w-3/4"></div>
	<div class="placeholder w-1/2"></div>
</div>
```

## Backgrounds

Four token-tinted treatments — mix or match per section:

```html
<!-- CSS-only, from motion.css -->
<section class="bg-aurora">…</section>
<!-- drifting blurred color blobs -->
<section class="bg-dotgrid">…</section>
<!-- dot grid fading out radially -->
```

```svelte
<!-- JS-driven, inside a position:relative container -->
<ParticleField density={40} />
<!-- floating particles -->
<MouseGlow />
<!-- soft glow that follows the pointer -->
```

## Components

```svelte
<script>
	import ParticleField from '$lib/components/fx/ParticleField.svelte';
	import AnimateIn from '$lib/components/fx/AnimateIn.svelte';
	import CountUp from '$lib/components/fx/CountUp.svelte';
</script>

<!-- exciting background: fills nearest position:relative ancestor -->
<section class="relative">
	<ParticleField density={40} />
	<div class="relative">hero content…</div>
</section>

<AnimateIn stagger={60} class="grid gap-4 sm:grid-cols-3">
	{#each features as f}<div class="card p-6">…</div>{/each}
</AnimateIn>

<p class="text-3xl font-bold"><CountUp value={1234} /></p>
<CountUp value={8.3} decimals={1} prefix="+" suffix="%" />
```

`ParticleField` and `MouseGlow` read `--color-primary-500` / `--color-secondary-500` /
`--color-tertiary-500` from the active theme, re-seeds on `data-theme` or
`.dark` changes, pauses offscreen, and goes static under reduced motion.

## App-page polish components

All live in `src/lib/components/fx/`; import like the others. Every one is
theme-token driven and reduced-motion safe.

```svelte
<!-- +layout.svelte — one wrapper, all routes transition -->
<PageTransition>{@render children()}</PageTransition>
<ScrollProgress />
<!-- or target="#docs-pane" for a container -->

<!-- data loading (PocketBase fetches) -->
<SkeletonBlock loaded={!!records} lines={3} avatar>
	{#each records as r}…{/each}
</SkeletonBlock>

<!-- dashboard metrics -->
<Sparkline data={requests} />
<Sparkline data={churn} color="var(--color-error-500)" fill={false} />

<!-- pricing / gallery cards -->
<TiltCard class="card p-6" max={8}>…</TiltCard>

<!-- hero headline -->
<h1>Build <TypeCycle words={['dashboards', 'inboxes', 'wizards']} class="text-primary-400" /></h1>

<!-- success moments (register done, wizard complete) -->
<ConfettiBurst bind:this={confetti} />
<button onclick={(e) => confetti.fire(e.clientX, e.clientY)}>Finish</button>

<!-- zero-data -->
<EmptyState title="No messages yet" description="When someone writes to you, it lands here.">
	{#snippet action()}<button class="btn preset-filled">New message</button>{/snippet}
</EmptyState>
```

Key props (full docs in each file's header comment):

- `PageTransition` — `duration=220`, `y=8`; View Transitions API, silent no-op fallback
- `SkeletonBlock` — `loaded`, `lines=3`, `avatar`, `media`; children render when loaded
- `ScrollProgress` — `target` selector (empty = window), `height=3`, `zIndex=50`
- `TiltCard` — `max=8` deg, `scale=1.02`, `glare`; off on touch + reduced motion
- `Sparkline` — `data: number[]`, `width/height`, `color`, `fill`, `dot`, `animate`
- `ConfettiBurst` — exported `fire(x?, y?, count=90)`; mount once near the page root
- `TypeCycle` — `words: string[]`, `typeSpeed`, `deleteSpeed`, `hold`, `caret`
- `EmptyState` — `title`, `description`, snippets `icon` / `action`
- `BuildStamp` — `label`, `inline`, `corner`, or wrap children for anchor mode; see wiring below

### BuildStamp — stale-bundle spot check

`svelte.config.js` — bake the git commit into the bundle and enable deploy polling:

```js
import { execSync } from 'node:child_process';

const commit = (() => {
	try {
		return execSync('git rev-parse --short HEAD').toString().trim();
	} catch {
		return `${Date.now()}`;
	}
})();

const config = {
	kit: {
		version: { name: commit, pollInterval: 60_000 }
		// …existing adapter/env config
	}
};
```

Then once in `+layout.svelte`:

```svelte
<BuildStamp />
<!-- fixed bottom-right, 35% opacity -->
<BuildStamp inline label="v1.4.2" />
<!-- or in the sidebar footer -->

<!-- Anchor mode: zero footprint. Wrap the app icon in NavToggle.svelte —
     hover/focus the logo to see the commit; an amber corner dot appears
     only when a newer deploy is detected (card gains a Reload button).
     A patched NavToggle.svelte is in integration/NavToggle.svelte. -->
<BuildStamp>
	<a href={resolve('/')} class="flex min-w-0 items-center gap-2" aria-label="{APP_NAME} home">
		<img src={logoUrl} alt="" class="size-8 shrink-0" />
	</a>
</BuildStamp>
```

- **Green dot + commit** — the bundle you are looking at is current (click copies the commit).
- **Amber + “update available”** — SvelteKit's `updated` store found a newer
  `_app/version.json`, i.e. this tab is stale; click hard-reloads.

The Go server's ldflags version (`internal/version`, served at `/api/version`)
can be surfaced too — pass it as `label` to compare server vs. bundle at a glance.
