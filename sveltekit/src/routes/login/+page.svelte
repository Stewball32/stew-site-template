<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth.svelte';
	import { Switch } from '@skeletonlabs/skeleton-svelte';
	import { LogInIcon, MailIcon, LockIcon } from '@lucide/svelte';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);
	let rememberMe = $state(false);

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

			<!-- Divider -->
			<div class="flex items-center gap-4">
				<hr class="hr flex-1" />
				<span class="text-xs opacity-50">or continue with</span>
				<hr class="hr flex-1" />
			</div>

			<!-- Social Login -->
			<div class="grid grid-cols-2 gap-3">
				<button class="btn w-full preset-tonal">
					<svg class="size-5" viewBox="0 0 24 24" fill="currentColor">
						<path
							d="M20.317 4.37a19.791 19.791 0 00-4.885-1.515.074.074 0 00-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 00-5.487 0 12.64 12.64 0 00-.617-1.25.077.077 0 00-.079-.037A19.736 19.736 0 003.677 4.37a.07.07 0 00-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 00.031.057 19.9 19.9 0 005.993 3.03.078.078 0 00.084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 00-.041-.106 13.107 13.107 0 01-1.872-.892.077.077 0 01-.008-.128 10.2 10.2 0 00.372-.292.074.074 0 01.077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 01.078.01c.12.098.246.198.373.292a.077.077 0 01-.006.127 12.299 12.299 0 01-1.873.892.077.077 0 00-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 00.084.028 19.839 19.839 0 006.002-3.03.077.077 0 00.032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 00-.031-.03z"
						/>
					</svg>
					<span>Discord</span>
				</button>
				<button class="btn w-full preset-tonal">
					<svg class="size-5" viewBox="0 0 24 24" fill="currentColor">
						<path
							d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"
						/>
					</svg>
					<span>GitHub</span>
				</button>
			</div>
		</div>

		<!-- Footer -->
		<p class="text-center text-xs opacity-50">
			By continuing, you agree to our Terms of Service and Privacy Policy.
		</p>
	</div>
</div>
