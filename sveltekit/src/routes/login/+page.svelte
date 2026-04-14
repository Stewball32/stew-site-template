<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth.svelte';
	import pb from '$lib/pocketbase';
	import { OAUTH_PROVIDERS } from '$lib/config/oauth';
	import { Switch } from '@skeletonlabs/skeleton-svelte';
	import { LogInIcon, MailIcon, LockIcon } from '@lucide/svelte';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);
	let rememberMe = $state(false);
	let enabledProviders = $state<string[]>([]);

	const visibleProviders = $derived(
		enabledProviders
			.filter((name) => name in OAUTH_PROVIDERS)
			.map((name) => ({ name, meta: OAUTH_PROVIDERS[name] }))
	);

	onMount(async () => {
		try {
			const methods = await pb.collection('users').listAuthMethods();
			enabledProviders = methods.oauth2?.providers?.map((p) => p.name) ?? [];
		} catch {
			enabledProviders = [];
		}
	});

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

	async function handleOAuth(provider: string) {
		error = '';
		loading = true;
		try {
			await auth.loginWithOAuth(provider);
			goto('/');
		} catch (err) {
			error = err instanceof Error ? err.message : 'OAuth sign-in failed. Please try again.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-[70vh] items-center justify-center">
	<div class="w-full max-w-md space-y-6">
		<div class="space-y-6 card p-8">
			<!-- Header -->
			<div class="space-y-2 text-center">
				<LogInIcon class="mx-auto size-10 text-primary-500" />
				<h1 class="h3">Sign In</h1>
				<p class="text-sm opacity-70">Enter your credentials to continue</p>
			</div>

			<!-- Error Alert -->
			{#if error}
				<aside class="alert preset-filled-error-500">
					<p>{error}</p>
				</aside>
			{/if}

			<!-- Form -->
			<form class="space-y-4" onsubmit={handleLogin}>
				<label class="label">
					<span>Email</span>
					<div class="input-group grid-cols-[auto_1fr_auto]">
						<div class="ig-cell"><MailIcon class="size-4" /></div>
						<input
							type="email"
							class="ig-input"
							placeholder="you@example.com"
							bind:value={email}
							required
						/>
					</div>
				</label>

				<label class="label">
					<span>Password</span>
					<div class="input-group grid-cols-[auto_1fr_auto]">
						<div class="ig-cell"><LockIcon class="size-4" /></div>
						<input
							type="password"
							class="ig-input"
							placeholder="••••••••"
							bind:value={password}
							required
						/>
					</div>
				</label>

				<div class="flex items-center justify-between">
					<div class="flex items-center gap-2">
						<Switch checked={rememberMe} onCheckedChange={(e) => (rememberMe = e.checked)} />
						<span class="text-sm">Remember me</span>
					</div>
					<a href="/login/" class="text-sm text-primary-500 hover:underline">Forgot password?</a>
				</div>

				<button type="submit" class="btn w-full preset-filled" disabled={loading}>
					{loading ? 'Signing in...' : 'Sign In'}
				</button>
			</form>

			{#if visibleProviders.length > 0}
				<!-- Divider -->
				<div class="flex items-center gap-4">
					<hr class="hr flex-1" />
					<span class="text-xs opacity-50">or continue with</span>
					<hr class="hr flex-1" />
				</div>

				<!-- Social Login -->
				<div class="grid grid-cols-2 gap-3">
					{#each visibleProviders as { name, meta } (name)}
						<button
							type="button"
							class="btn w-full preset-tonal"
							disabled={loading}
							onclick={() => handleOAuth(name)}
						>
							<span>{meta.label}</span>
						</button>
					{/each}
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<!-- <p class="text-center text-xs opacity-50">
			By continuing, you agree to our Terms of Service and Privacy Policy.
		</p> -->
	</div>
</div>
