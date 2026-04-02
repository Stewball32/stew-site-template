<script lang="ts">
	import { Switch } from '@skeletonlabs/skeleton-svelte';

	let notifications = $state(true);
	let emailUpdates = $state(false);
</script>

<h1 class="h2 mb-6">Settings</h1>

<div class="card p-6 max-w-lg space-y-6">
	<h2 class="h4">General</h2>

	<label class="label">
		<span>Display Name</span>
		<input type="text" class="input" placeholder="Jane Doe" />
	</label>

	<label class="label">
		<span>Email</span>
		<input type="email" class="input" placeholder="jane.doe@example.com" />
	</label>

	<hr class="hr" />

	<h2 class="h4">Notifications</h2>

	<div class="flex items-center justify-between">
		<div>
			<p class="font-semibold">Push Notifications</p>
			<p class="text-sm opacity-70">Receive push notifications for important updates</p>
		</div>
		<Switch bind:checked={notifications} />
	</div>

	<div class="flex items-center justify-between">
		<div>
			<p class="font-semibold">Email Updates</p>
			<p class="text-sm opacity-70">Receive weekly email digests</p>
		</div>
		<Switch bind:checked={emailUpdates} />
	</div>

	<hr class="hr" />

	<button class="btn preset-filled">Save Changes</button>
</div>
