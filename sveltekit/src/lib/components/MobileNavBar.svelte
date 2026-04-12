<script lang="ts">
	import { Navigation } from '@skeletonlabs/skeleton-svelte';
	import { navLinks } from '$lib/config/navigation';

	let { currentPath }: { currentPath: string } = $props();
</script>

<div class="fixed inset-x-0 bottom-0 z-30 flex sm:hidden">
	<Navigation layout="bar" class="w-full">
		<Navigation.Content>
			<Navigation.Menu>
				{#each navLinks as link}
					<Navigation.TriggerAnchor
						href={link.href}
						aria-current={currentPath === link.href ? 'page' : undefined}
						class="aria-[current=page]:preset-tonal"
					>
						<link.icon class="size-5" />
						<span class="text-xs">{link.label}</span>
					</Navigation.TriggerAnchor>
				{/each}
			</Navigation.Menu>
		</Navigation.Content>
	</Navigation>
</div>
