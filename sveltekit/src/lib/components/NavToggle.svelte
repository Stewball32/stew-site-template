<script lang="ts">
	import { APP_NAME } from '$lib/config/app';
	import { MenuIcon } from '@lucide/svelte';

	let { onclick }: { onclick: () => void } = $props();
</script>

<div class="flex flex-row gap-3">
	<button class="btn-icon hover:preset-tonal" aria-label="Toggle navigation" {onclick}>
		<MenuIcon class="size-10" />
	</button>
	<span class="text-xl font-bold">{APP_NAME}</span>
</div>
