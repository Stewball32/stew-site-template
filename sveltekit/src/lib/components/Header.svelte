<script lang="ts">
	import { AppBar } from '@skeletonlabs/skeleton-svelte';
	import NavToggleButton from '$lib/components/NavToggle.svelte';

	let { onToggle }: { onToggle: () => void } = $props();
</script>

<AppBar class="h-16 p-4 ">
	<AppBar.Toolbar class="">
		<AppBar.Lead>
			<NavToggleButton onclick={onToggle} />
		</AppBar.Lead>
	</AppBar.Toolbar>
</AppBar>
