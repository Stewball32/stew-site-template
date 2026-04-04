<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/Header.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import MobileDrawer from '$lib/components/MobileDrawer.svelte';
	import MobileNav from '$lib/components/MobileNav.svelte';
	import { page } from '$app/state';

	let { children } = $props();
	let sidebarExpanded = $state(true);
	let drawerOpen = $state(false);

	function handleToggle() {
		if (window.matchMedia('(max-width: 1023px)').matches) {
			drawerOpen = !drawerOpen;
		} else {
			sidebarExpanded = !sidebarExpanded;
		}
	}
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="flex h-screen flex-col">
	<Header onToggle={handleToggle} />
	<div class="flex flex-1 overflow-hidden">
		<Sidebar currentPath={page.url.pathname} bind:expanded={sidebarExpanded} />
		<main class="flex-1 overflow-y-auto p-4 pb-20 lg:p-8 lg:pb-8">
			{@render children()}
		</main>
	</div>
	<MobileNav currentPath={page.url.pathname} />
</div>

<MobileDrawer bind:open={drawerOpen} currentPath={page.url.pathname} />
