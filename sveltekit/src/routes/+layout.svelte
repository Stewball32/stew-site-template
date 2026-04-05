<script lang="ts">
	import { onMount } from 'svelte';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/Header.svelte';
	import NavPanel from '$lib/components/NavPanel.svelte';
	import MobileNav from '$lib/components/MobileNav.svelte';
	import { page } from '$app/state';

	let { children } = $props();
	let navOpen = $state(false);
	let isDesktop = $state(false);

	onMount(() => {
		const mql = window.matchMedia('(min-width: 1024px)');
		isDesktop = mql.matches;
		if (mql.matches) navOpen = true;
		const handler = (e: MediaQueryListEvent) => {
			isDesktop = e.matches;
			if (e.matches && !navOpen) navOpen = true;
		};
		mql.addEventListener('change', handler);
		return () => mql.removeEventListener('change', handler);
	});

	function handleToggle() {
		navOpen = !navOpen;
	}
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="flex h-screen flex-col">
	<Header onToggle={handleToggle} />
	<div class="flex flex-1 overflow-hidden">
		<NavPanel bind:open={navOpen} {isDesktop} currentPath={page.url.pathname} />
		<main class="flex-1 overflow-y-auto p-4 pb-20 sm:pb-4 lg:p-8 lg:pb-8">
			{@render children()}
		</main>
	</div>
	<MobileNav currentPath={page.url.pathname} />
</div>
