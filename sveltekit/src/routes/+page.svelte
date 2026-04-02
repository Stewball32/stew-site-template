<script lang="ts">
	import { UsersIcon, ActivityIcon, FolderIcon, TrendingUpIcon } from '@lucide/svelte';
</script>

<h1 class="h2 mb-6">Dashboard</h1>

<div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
	<div class="card p-6 space-y-2">
		<div class="flex items-center gap-3">
			<UsersIcon class="size-8 text-primary-500" />
			<h2 class="h4">Users</h2>
		</div>
		<p class="text-3xl font-bold">1,234</p>
		<p class="text-sm opacity-70">Active accounts this month</p>
	</div>

	<div class="card p-6 space-y-2">
		<div class="flex items-center gap-3">
			<ActivityIcon class="size-8 text-secondary-500" />
			<h2 class="h4">Activity</h2>
		</div>
		<p class="text-3xl font-bold">5,678</p>
		<p class="text-sm opacity-70">Events in the last 24 hours</p>
	</div>

	<div class="card p-6 space-y-2">
		<div class="flex items-center gap-3">
			<FolderIcon class="size-8 text-tertiary-500" />
			<h2 class="h4">Projects</h2>
		</div>
		<p class="text-3xl font-bold">42</p>
		<p class="text-sm opacity-70">Active projects</p>
	</div>

	<div class="card p-6 space-y-2">
		<div class="flex items-center gap-3">
			<TrendingUpIcon class="size-8 text-success-500" />
			<h2 class="h4">Growth</h2>
		</div>
		<p class="text-3xl font-bold">+12%</p>
		<p class="text-sm opacity-70">Compared to last month</p>
	</div>
</div>
