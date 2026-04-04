<script lang="ts">
	import { Navigation } from '@skeletonlabs/skeleton-svelte';
	import { LogInIcon, LogOutIcon } from '@lucide/svelte';
	import { navLinks } from '$lib/config/navigation';
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	let { open = $bindable(false), currentPath }: { open: boolean; currentPath: string } = $props();

	function close() {
		open = false;
	}
</script>

{#if open}
	<!-- Backdrop -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div class="fixed inset-0 z-40 bg-black/50" onkeydown={close} onclick={close}></div>

	<!-- Drawer panel -->
	<div class="fixed inset-y-0 left-0 z-50 w-64">
		<Navigation layout="sidebar" class="h-full">
			<Navigation.Header>
				<span class="p-4 text-xl font-bold">App Template</span>
			</Navigation.Header>
			<Navigation.Content>
				<Navigation.Menu>
					{#each navLinks as link}
						<Navigation.TriggerAnchor
							href={link.href}
							aria-current={currentPath === link.href ? 'page' : undefined}
							class="aria-[current=page]:preset-tonal"
							onclick={close}
						>
							<link.icon class="size-5" />
							<span>{link.label}</span>
						</Navigation.TriggerAnchor>
					{/each}
				</Navigation.Menu>
			</Navigation.Content>
			<Navigation.Footer>
				{#if auth.isLoggedIn}
					<button
						class="btn w-full preset-tonal"
						onclick={() => {
							auth.logout();
							close();
							goto('/login/');
						}}
					>
						<LogOutIcon class="size-5" />
						<span>Logout</span>
					</button>
				{:else}
					<a href="/login/" class="btn w-full preset-tonal" onclick={close}>
						<LogInIcon class="size-5" />
						<span>Sign In</span>
					</a>
				{/if}
			</Navigation.Footer>
		</Navigation>
	</div>
{/if}
