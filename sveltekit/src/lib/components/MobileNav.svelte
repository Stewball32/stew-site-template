<script lang="ts">
	import { Navigation } from '@skeletonlabs/skeleton-svelte';
	import { navLinks } from '$lib/config/navigation';

	let { currentPath }: { currentPath: string } = $props();
</script>

<div class="flex lg:hidden fixed bottom-0 inset-x-0 z-30">
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
