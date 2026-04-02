<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth.svelte';
	import { LogInIcon } from '@lucide/svelte';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleLogin(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		loading = true;
		try {
			await auth.login(email, password);
			goto('/');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Login failed. Please try again.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex items-center justify-center min-h-[70vh]">
	<div class="card p-8 w-full max-w-md space-y-6">
		<div class="text-center space-y-2">
			<LogInIcon class="size-10 mx-auto text-primary-500" />
			<h1 class="h3">Sign In</h1>
			<p class="text-sm opacity-70">Enter your credentials to continue</p>
		</div>

		{#if error}
			<aside class="alert preset-filled-error-500">
				<p>{error}</p>
			</aside>
		{/if}

		<form class="space-y-4" onsubmit={handleLogin}>
			<label class="label">
				<span>Email</span>
				<input type="email" class="input" placeholder="you@example.com" bind:value={email} required />
			</label>

			<label class="label">
				<span>Password</span>
				<input type="password" class="input" placeholder="••••••••" bind:value={password} required />
			</label>

			<button type="submit" class="btn preset-filled w-full" disabled={loading}>
				{loading ? 'Signing in...' : 'Sign In'}
			</button>
		</form>
	</div>
</div>
