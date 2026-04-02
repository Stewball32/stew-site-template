<script lang="ts">
	import { Navigation } from '@skeletonlabs/skeleton-svelte';
	import { PanelLeftCloseIcon, PanelLeftOpenIcon } from '@lucide/svelte';
	import { navLinks } from '$lib/config/navigation';

	let { currentPath, expanded = $bindable(true) }: { currentPath: string; expanded: boolean } =
		$props();
</script>

<div class="hidden lg:flex h-full">
	<Navigation layout={expanded ? 'sidebar' : 'rail'}>
		<Navigation.Content>
			<Navigation.Menu>
				{#each navLinks as link}
					<Navigation.TriggerAnchor
						href={link.href}
						aria-current={currentPath === link.href ? 'page' : undefined}
						class="aria-[current=page]:preset-tonal"
					>
						<link.icon class="size-5" />
						{#if expanded}
							<span>{link.label}</span>
						{/if}
					</Navigation.TriggerAnchor>
				{/each}
			</Navigation.Menu>
		</Navigation.Content>
		<Navigation.Footer>
			<button class="btn-icon hover:preset-tonal mx-auto" onclick={() => (expanded = !expanded)}>
				{#if expanded}
					<PanelLeftCloseIcon class="size-5" />
				{:else}
					<PanelLeftOpenIcon class="size-5" />
				{/if}
			</button>
		</Navigation.Footer>
	</Navigation>
</div>
