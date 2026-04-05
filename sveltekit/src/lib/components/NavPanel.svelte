<script lang="ts">
	import { fade } from 'svelte/transition';
	import { Navigation } from '@skeletonlabs/skeleton-svelte';
	import { CircleUserIcon, LogInIcon, LogOutIcon } from '@lucide/svelte';
	import NavToggleButton from '$lib/components/NavToggle.svelte';
	import { navLinks } from '$lib/config/navigation';
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	let {
		open = $bindable(false),
		isDesktop,
		currentPath
	}: {
		open: boolean;
		isDesktop: boolean;
		currentPath: string;
	} = $props();

	function close() {
		open = false;
	}

	function handleLogout() {
		auth.logout();
		close();
		goto('/login/');
	}

	let navLayout = $derived<'rail' | 'sidebar'>(
		isDesktop ? (open ? 'sidebar' : 'rail') : 'sidebar'
	);
</script>

<!-- Mobile backdrop -->
{#if !isDesktop && open}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		class="fixed inset-0 z-40 bg-black/50"
		onclick={close}
		onkeydown={close}
		transition:fade={{ duration: 200 }}
	></div>
{/if}

<!-- Nav panel -->
<div
	class:h-full={isDesktop}
	class:flex={isDesktop}
	class:fixed={!isDesktop}
	class:inset-y-0={!isDesktop}
	class:left-0={!isDesktop}
	class:z-50={!isDesktop}
	class:transition-transform={!isDesktop}
	class:duration-300={!isDesktop}
	class:-translate-x-full={!isDesktop && !open}
	class:translate-x-0={!isDesktop && open}
>
	<Navigation layout={navLayout} class="h-full">
		{#if !isDesktop}
			<Navigation.Header class="pb-4">
				<NavToggleButton onclick={close} />
			</Navigation.Header>
		{/if}
		<Navigation.Content>
			<Navigation.Menu>
				{#each navLinks as link}
					<Navigation.TriggerAnchor
						href={link.href}
						aria-current={currentPath === link.href ? 'page' : undefined}
						class="aria-[current=page]:preset-tonal"
						onclick={!isDesktop ? close : undefined}
					>
						<link.icon class="size-5" />
						{#if navLayout === 'sidebar'}
							<span>{link.label}</span>
						{/if}
					</Navigation.TriggerAnchor>
				{/each}
			</Navigation.Menu>
		</Navigation.Content>
		<Navigation.Footer>
			<Navigation.Menu>
				{#if auth.isLoggedIn}
					<Navigation.TriggerAnchor
						href="/profile/"
						aria-current={currentPath === '/profile/' ? 'page' : undefined}
						class="aria-[current=page]:preset-tonal"
						onclick={!isDesktop ? close : undefined}
					>
						<CircleUserIcon class="size-5" />
						{#if navLayout === 'sidebar'}
							<span class="flex-1 truncate">{auth.user?.email}</span>
						{/if}
					</Navigation.TriggerAnchor>
					{#if navLayout === 'sidebar'}
						<Navigation.Trigger onclick={handleLogout}>
							<LogOutIcon class="size-5" />
							<span>Logout</span>
						</Navigation.Trigger>
					{/if}
				{:else}
					<Navigation.TriggerAnchor href="/login/" onclick={!isDesktop ? close : undefined}>
						<LogInIcon class="size-5" />
						{#if navLayout === 'sidebar'}
							<span>Sign In</span>
						{/if}
					</Navigation.TriggerAnchor>
				{/if}
			</Navigation.Menu>
		</Navigation.Footer>
	</Navigation>
</div>
