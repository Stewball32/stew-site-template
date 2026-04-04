<script lang="ts">
	import { AppBar } from '@skeletonlabs/skeleton-svelte';
	import { MenuIcon, LogInIcon, LogOutIcon, CircleUserIcon } from '@lucide/svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	let { onToggle }: { onToggle: () => void } = $props();

	function handleLogout() {
		auth.logout();
		goto('/login/');
	}
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
			{#if auth.isLoggedIn}
				<span class="hidden text-sm opacity-70 sm:inline">{auth.user?.email}</span>
				<button class="btn-icon hover:preset-tonal" aria-label="Logout" onclick={handleLogout}>
					<LogOutIcon class="size-5" />
				</button>
			{:else}
				<a href="/login/" class="btn-icon hover:preset-tonal" aria-label="Login">
					<LogInIcon class="size-5" />
				</a>
			{/if}
		</AppBar.Trail>
	</AppBar.Toolbar>
</AppBar>
