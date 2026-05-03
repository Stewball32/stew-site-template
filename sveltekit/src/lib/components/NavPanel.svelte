<script lang="ts">
	import { fade } from 'svelte/transition';
	import { Navigation } from '@skeletonlabs/skeleton-svelte';
	import NavToggleButton from '$lib/components/NavToggle.svelte';
	import { mainGroups, footerLinks, type NavLink } from '$lib/config/navigation';

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

	$effect(() => {
		if (!open || isDesktop) return;
		const onKey = (e: KeyboardEvent) => {
			if (e.key === 'Escape') close();
		};
		window.addEventListener('keydown', onKey);
		return () => window.removeEventListener('keydown', onKey);
	});

	let navLayout = $derived<'rail' | 'sidebar'>(
		isDesktop || isTablet ? (open ? 'sidebar' : 'rail') : 'sidebar'
	);

	function isVisible(link: NavLink, layout: 'rail' | 'sidebar'): boolean {
		return layout === 'sidebar' ? (link.showInDrawer ?? true) : (link.showInRail ?? true);
	}

	let visibleGroups = $derived(
		mainGroups
			.map((group) => ({ ...group, links: group.links.filter((l) => isVisible(l, navLayout)) }))
			.filter((group) => group.links.length > 0)
	);

	let visibleFooterLinks = $derived(footerLinks.filter((l) => isVisible(l, navLayout)));
</script>

<!-- Mobile backdrop -->
{#if !isDesktop && open}
	<button
		type="button"
		aria-label="Close navigation"
		class:fixed={!isTablet}
		class:absolute={isTablet}
		class="inset-0 z-40 cursor-default bg-black/50"
		onclick={close}
		transition:fade={{ duration: 200 }}
	></button>
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
	<Navigation layout={navLayout} class="flex h-full min-h-0 flex-col overflow-hidden">
		{#if !isDesktop && !isTablet}
			<Navigation.Header class="pb-4">
				<NavToggleButton onclick={close} />
			</Navigation.Header>
		{/if}
		<Navigation.Content class="flex min-h-0 flex-1 flex-col overflow-y-auto">
			{#each visibleGroups as group (group.label)}
				<Navigation.Group>
					{#if navLayout === 'sidebar'}
						<Navigation.Label>
							{#if group.href}
								<!-- group.href is config-driven, not a static route literal -->
								<!-- eslint-disable svelte/no-navigation-without-resolve -->
								<a
									href={group.href}
									aria-current={currentPath === group.href ? 'page' : undefined}
									class="hover:underline aria-[current=page]:underline"
									onclick={!isDesktop ? close : undefined}
								>
									{group.label}
								</a>
								<!-- eslint-enable svelte/no-navigation-without-resolve -->
							{:else}
								{group.label}
							{/if}
						</Navigation.Label>
					{/if}
					<Navigation.Menu>
						{#each group.links as link (link.href)}
							<Navigation.TriggerAnchor
								href={link.href}
								aria-current={currentPath === link.href ? 'page' : undefined}
								class="aria-[current=page]:preset-tonal"
								onclick={!isDesktop ? close : undefined}
							>
								<link.icon class="size-5" />
								<Navigation.TriggerText>{link.label}</Navigation.TriggerText>
							</Navigation.TriggerAnchor>
						{/each}
					</Navigation.Menu>
				</Navigation.Group>
			{/each}
		</Navigation.Content>
		<Navigation.Footer>
			<Navigation.Menu>
				{#each visibleFooterLinks as link (link.href)}
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
