<script lang="ts">
	import { fade } from 'svelte/transition';
	import { Navigation } from '@skeletonlabs/skeleton-svelte';
	import NavToggleButton from '$lib/components/NavToggle.svelte';
	import { mainLinks, footerLinks } from '$lib/config/navigation';

	let {
		open = $bindable(false),
		isDesktop,
		isTablet,
		currentPath
	}: {
		open: boolean;
		isDesktop: boolean;
		isTablet: boolean;
		currentPath: string;
	} = $props();

	function close() {
		open = false;
	}

	let navLayout = $derived<'rail' | 'sidebar'>(
		isDesktop || isTablet ? (open ? 'sidebar' : 'rail') : 'sidebar'
	);
</script>

<!-- Mobile backdrop -->
{#if !isDesktop && open}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		class:fixed={!isTablet}
		class:absolute={isTablet}
		class="inset-0 z-40 bg-black/50"
		onclick={close}
		onkeydown={close}
		transition:fade={{ duration: 200 }}
	></div>
{/if}

<!-- Nav panel -->
<div
	class:h-full={isDesktop}
	class:flex={isDesktop}
	class:fixed={!isDesktop && !isTablet}
	class:absolute={isTablet}
	class:inset-y-0={!isDesktop}
	class:left-0={!isDesktop}
	class:z-50={!isDesktop}
	class:transition-transform={!isDesktop && !isTablet}
	class:duration-300={!isDesktop && !isTablet}
	class:-translate-x-full={!isDesktop && !isTablet && !open}
	class:translate-x-0={!isDesktop && !isTablet && open}
>
	<Navigation layout={navLayout} class="h-full">
		{#if !isDesktop && !isTablet}
			<Navigation.Header class="pb-4">
				<NavToggleButton onclick={close} />
			</Navigation.Header>
		{/if}
		<Navigation.Content>
			<Navigation.Menu>
				{#each mainLinks as link}
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
				{#each footerLinks as link}
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
		</Navigation.Footer>
	</Navigation>
</div>
