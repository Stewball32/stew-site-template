<script lang="ts">
	import { AppBar } from '@skeletonlabs/skeleton-svelte';
	import { MenuIcon, CircleUserIcon } from '@lucide/svelte';

	let { onToggle }: { onToggle: () => void } = $props();
</script>

<AppBar>
	<AppBar.Toolbar>
		<AppBar.Lead>
			<button class="btn-icon hover:preset-tonal" aria-label="Toggle navigation" onclick={onToggle}>
				<MenuIcon class="size-5" />
			</button>
			<span class="text-xl font-bold">App Template</span>
		</AppBar.Lead>
		<AppBar.Trail>
			<button class="btn-icon hover:preset-tonal" aria-label="User menu">
				<CircleUserIcon class="size-6" />
			</button>
		</AppBar.Trail>
	</AppBar.Toolbar>
</AppBar>
