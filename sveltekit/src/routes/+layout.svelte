<script lang="ts">
	import { onMount } from 'svelte';
	import { Toast } from '@skeletonlabs/skeleton-svelte';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/Header.svelte';
	import NavPanel from '$lib/components/NavPanel.svelte';
	import MobileNavBar from '$lib/components/MobileNavBar.svelte';
	import { toaster } from '$lib/stores/toaster';
	import { page } from '$app/state';

	let { children } = $props();
	let navOpen = $state(false);
	let isTablet = $state(false);
	let isDesktop = $state(false);

	onMount(() => {
		const mqlSm = window.matchMedia('(min-width: 640px)');
		const mqlLg = window.matchMedia('(min-width: 1024px)');

		isDesktop = mqlLg.matches;
		isTablet = mqlSm.matches && !mqlLg.matches;
		if (mqlLg.matches) navOpen = true;

		const handlerSm = (e: MediaQueryListEvent) => {
			isTablet = e.matches && !mqlLg.matches;
			if (!e.matches) navOpen = false;
		};
		const handlerLg = (e: MediaQueryListEvent) => {
			isDesktop = e.matches;
			isTablet = mqlSm.matches && !e.matches;
			navOpen = e.matches;
		};

		mqlSm.addEventListener('change', handlerSm);
		mqlLg.addEventListener('change', handlerLg);
		return () => {
			mqlSm.removeEventListener('change', handlerSm);
			mqlLg.removeEventListener('change', handlerLg);
		};
	});

	function handleToggle() {
		navOpen = !navOpen;
	}
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="flex h-screen flex-col">
	<Header onToggle={handleToggle} />
	<div class="relative flex flex-1 overflow-hidden">
		<div class="hidden w-25 shrink-0 sm:flex lg:hidden" aria-hidden="true"></div>
		<NavPanel bind:open={navOpen} {isDesktop} {isTablet} currentPath={page.url.pathname} />
		<main class="flex-1 overflow-y-auto p-4 pb-20 sm:pb-4 lg:p-8 lg:pb-8">
			{@render children()}
		</main>
	</div>
	<MobileNavBar currentPath={page.url.pathname} />
</div>

<Toast.Group {toaster}>
	{#snippet children(toast)}
		<Toast {toast}>
			<Toast.Message>
				<Toast.Title>{toast.title}</Toast.Title>
				<Toast.Description>{toast.description}</Toast.Description>
			</Toast.Message>
			<Toast.CloseTrigger />
		</Toast>
	{/snippet}
</Toast.Group>
