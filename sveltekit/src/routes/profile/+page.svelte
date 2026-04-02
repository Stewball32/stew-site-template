<script lang="ts">
	import { Avatar } from '@skeletonlabs/skeleton-svelte';
</script>

<h1 class="h2 mb-6">Profile</h1>

<div class="card p-6 max-w-lg space-y-6">
	<div class="flex items-center gap-4">
		<Avatar src="https://i.pravatar.cc/160?img=5" name="Jane Doe" class="size-20" />
		<div>
			<h2 class="h4">Jane Doe</h2>
			<p class="text-sm opacity-70">jane.doe@example.com</p>
		</div>
	</div>

	<hr class="hr" />

	<div class="space-y-4">
		<div>
			<span class="text-sm font-semibold opacity-70">Display Name</span>
			<p>Jane Doe</p>
		</div>
		<div>
			<span class="text-sm font-semibold opacity-70">Email</span>
			<p>jane.doe@example.com</p>
		</div>
		<div>
			<span class="text-sm font-semibold opacity-70">Bio</span>
			<p>Full-stack developer who enjoys building things with SvelteKit and Go.</p>
		</div>
	</div>

	<button class="btn preset-filled">Edit Profile</button>
</div>
