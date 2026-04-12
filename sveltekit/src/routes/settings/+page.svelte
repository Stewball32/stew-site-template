<script lang="ts">
	import {
		Tabs,
		Switch,
		Slider,
		FileUpload,
		Avatar,
		ToggleGroup,
		Progress
	} from '@skeletonlabs/skeleton-svelte';
	import {
		UserIcon,
		MailIcon,
		GlobeIcon,
		BellIcon,
		BellOffIcon,
		MailCheckIcon,
		MessageSquareIcon,
		PaletteIcon,
		Trash2Icon,
		UploadIcon,
		XIcon
	} from '@lucide/svelte';

	let activeTab = $state('general');

	// General tab state
	let displayName = $state('Jane Doe');
	let email = $state('jane.doe@example.com');
	let timezone = $state('America/New_York');

	// Notifications tab state
	let pushNotifs = $state(true);
	let emailNotifs = $state(true);
	let inAppNotifs = $state(false);
	let notifFrequency = $state([50]);

	// Appearance tab state
	let density = $state<string[]>(['comfortable']);
	let reduceAnimations = $state(false);

	const themeColors = [
		{ name: 'Primary', class: 'bg-primary-500' },
		{ name: 'Secondary', class: 'bg-secondary-500' },
		{ name: 'Tertiary', class: 'bg-tertiary-500' },
		{ name: 'Success', class: 'bg-success-500' },
		{ name: 'Warning', class: 'bg-warning-500' },
		{ name: 'Error', class: 'bg-error-500' }
	];
</script>

<h1 class="mb-6 h2">Settings</h1>

<div class="max-w-2xl">
	<Tabs value={activeTab} onValueChange={(e) => (activeTab = e.value)}>
		<Tabs.List class="mb-6">
			<Tabs.Trigger value="general">General</Tabs.Trigger>
			<Tabs.Trigger value="notifications">Notifications</Tabs.Trigger>
			<Tabs.Trigger value="appearance">Appearance</Tabs.Trigger>
			<Tabs.Indicator />
		</Tabs.List>

		<!-- General Tab -->
		<Tabs.Content value="general">
			<div class="space-y-6">
				<!-- Account Info -->
				<div class="space-y-6 card p-6">
					<h2 class="h4">Account</h2>

					<!-- Avatar Upload -->
					<div class="flex items-center gap-6">
						<Avatar class="size-20">
							<Avatar.Image src="https://i.pravatar.cc/160?img=5" />
							<Avatar.Fallback>JD</Avatar.Fallback>
						</Avatar>
						<FileUpload maxFiles={1} accept="image/*">
							<FileUpload.Dropzone class="card preset-outlined-surface-200-800 p-4">
								<div class="flex flex-col items-center gap-2 text-center">
									<UploadIcon class="size-8 text-surface-400-600" />
									<p class="text-sm">
										<span class="font-semibold text-primary-500">Click to upload</span> or drag and drop
									</p>
									<p class="text-xs opacity-50">PNG, JPG up to 2MB</p>
								</div>
							</FileUpload.Dropzone>
							<FileUpload.HiddenInput />
						</FileUpload>
					</div>

					<hr class="hr" />

					<!-- Input Groups -->
					<label class="label">
						<span>Display Name</span>
						<div class="input-group grid-cols-[auto_1fr_auto]">
							<div class="ig-cell"><UserIcon class="size-4" /></div>
							<input type="text" class="ig-input" bind:value={displayName} />
						</div>
					</label>

					<label class="label">
						<span>Email</span>
						<div class="input-group grid-cols-[auto_1fr_auto]">
							<div class="ig-cell"><MailIcon class="size-4" /></div>
							<input type="email" class="ig-input" bind:value={email} />
						</div>
					</label>

					<label class="label">
						<span>Timezone</span>
						<div class="input-group grid-cols-[auto_1fr_auto]">
							<div class="ig-cell"><GlobeIcon class="size-4" /></div>
							<select class="ig-select" bind:value={timezone}>
								<option value="America/New_York">Eastern Time (ET)</option>
								<option value="America/Chicago">Central Time (CT)</option>
								<option value="America/Denver">Mountain Time (MT)</option>
								<option value="America/Los_Angeles">Pacific Time (PT)</option>
								<option value="Europe/London">Greenwich Mean Time (GMT)</option>
								<option value="Europe/Berlin">Central European Time (CET)</option>
								<option value="Asia/Tokyo">Japan Standard Time (JST)</option>
							</select>
						</div>
					</label>
				</div>

				<!-- Danger Zone -->
				<div class="space-y-3 card preset-outlined-error-500 p-6">
					<h2 class="h4 text-error-500">Danger Zone</h2>
					<p class="text-sm opacity-70">
						Permanently delete your account and all associated data. This action cannot be undone.
					</p>
					<button class="btn preset-filled-error-500">
						<Trash2Icon class="size-4" />
						<span>Delete Account</span>
					</button>
				</div>

				<!-- Save -->
				<div class="flex gap-3">
					<button class="btn preset-filled">Save Changes</button>
					<button class="btn preset-tonal">Reset</button>
				</div>
			</div>
		</Tabs.Content>

		<!-- Notifications Tab -->
		<Tabs.Content value="notifications">
			<div class="space-y-6">
				<div class="space-y-6 card p-6">
					<h2 class="h4">Notification Preferences</h2>

					<div class="flex items-center justify-between">
						<div class="flex items-center gap-3">
							<BellIcon class="size-5 text-primary-500" />
							<div>
								<p class="font-semibold">Push Notifications</p>
								<p class="text-sm opacity-70">Receive alerts in your browser</p>
							</div>
						</div>
						<Switch checked={pushNotifs} onCheckedChange={(e) => (pushNotifs = e.checked)} />
					</div>

					<hr class="hr" />

					<div class="flex items-center justify-between">
						<div class="flex items-center gap-3">
							<MailCheckIcon class="size-5 text-secondary-500" />
							<div>
								<p class="font-semibold">Email Notifications</p>
								<p class="text-sm opacity-70">Get updates delivered to your inbox</p>
							</div>
						</div>
						<Switch checked={emailNotifs} onCheckedChange={(e) => (emailNotifs = e.checked)} />
					</div>

					<hr class="hr" />

					<div class="flex items-center justify-between">
						<div class="flex items-center gap-3">
							<MessageSquareIcon class="size-5 text-tertiary-500" />
							<div>
								<p class="font-semibold">In-App Messages</p>
								<p class="text-sm opacity-70">Show notification popups within the app</p>
							</div>
						</div>
						<Switch checked={inAppNotifs} onCheckedChange={(e) => (inAppNotifs = e.checked)} />
					</div>
				</div>

				<div class="space-y-4 card p-6">
					<h2 class="h4">Digest Frequency</h2>
					<p class="text-sm opacity-70">Control how often you receive notification summaries.</p>
					<Slider
						value={notifFrequency}
						onValueChange={(e) => (notifFrequency = e.value)}
						min={0}
						max={100}
						step={25}
					>
						<div class="flex items-center justify-between">
							<Slider.Label class="text-sm font-semibold">Frequency</Slider.Label>
							<Slider.ValueText class="text-sm opacity-70" />
						</div>
						<Slider.Control>
							<Slider.Track class="preset-outlined-surface-200-800">
								<Slider.Range />
							</Slider.Track>
							<Slider.Thumb index={0} />
						</Slider.Control>
						<Slider.MarkerGroup>
							<Slider.Marker value={0} class="text-xs opacity-50">Instant</Slider.Marker>
							<Slider.Marker value={25} class="text-xs opacity-50">Hourly</Slider.Marker>
							<Slider.Marker value={50} class="text-xs opacity-50">Daily</Slider.Marker>
							<Slider.Marker value={75} class="text-xs opacity-50">Weekly</Slider.Marker>
							<Slider.Marker value={100} class="text-xs opacity-50">Monthly</Slider.Marker>
						</Slider.MarkerGroup>
					</Slider>
				</div>

				<div class="flex gap-3">
					<button class="btn preset-filled">Save Changes</button>
					<button class="btn preset-tonal">Reset</button>
				</div>
			</div>
		</Tabs.Content>

		<!-- Appearance Tab -->
		<Tabs.Content value="appearance">
			<div class="space-y-6">
				<div class="space-y-6 card p-6">
					<h2 class="h4">Display</h2>

					<div class="space-y-3">
						<p class="font-semibold">Layout Density</p>
						<p class="text-sm opacity-70">Choose how compact the interface appears.</p>
						<ToggleGroup value={density} onValueChange={(e) => (density = e.value)}>
							<ToggleGroup.Item value="compact" class="btn preset-tonal btn-sm">
								Compact
							</ToggleGroup.Item>
							<ToggleGroup.Item value="comfortable" class="btn preset-tonal btn-sm">
								Comfortable
							</ToggleGroup.Item>
							<ToggleGroup.Item value="spacious" class="btn preset-tonal btn-sm">
								Spacious
							</ToggleGroup.Item>
						</ToggleGroup>
					</div>

					<hr class="hr" />

					<div class="flex items-center justify-between">
						<div>
							<p class="font-semibold">Reduce Animations</p>
							<p class="text-sm opacity-70">Minimize motion effects throughout the UI</p>
						</div>
						<Switch
							checked={reduceAnimations}
							onCheckedChange={(e) => (reduceAnimations = e.checked)}
						/>
					</div>
				</div>

				<div class="space-y-4 card p-6">
					<h2 class="h4">Theme Colors</h2>
					<p class="text-sm opacity-70">Current color palette from the active theme.</p>
					<div class="grid grid-cols-6 gap-3">
						{#each themeColors as color}
							<div class="space-y-1 text-center">
								<div class="mx-auto size-10 rounded-full {color.class}"></div>
								<p class="text-xs opacity-70">{color.name}</p>
							</div>
						{/each}
					</div>
				</div>

				<div class="space-y-4 card p-6">
					<h2 class="h4">Component Preview</h2>
					<p class="text-sm opacity-70">See how the current theme looks on common components.</p>
					<div class="flex flex-wrap items-center gap-3">
						<button class="btn preset-filled-primary-500 btn-sm">Primary</button>
						<button class="btn preset-filled-secondary-500 btn-sm">Secondary</button>
						<button class="btn preset-tonal-tertiary btn-sm">Tertiary</button>
						<button class="preset-outlined-surface btn btn-sm">Outlined</button>
					</div>
					<div class="flex flex-wrap items-center gap-2">
						<span class="badge preset-filled-success-500">Success</span>
						<span class="badge preset-filled-warning-500">Warning</span>
						<span class="badge preset-filled-error-500">Error</span>
						<span class="badge preset-tonal-primary">Tonal</span>
					</div>
					<Progress value={65} max={100}>
						<Progress.Track class="h-2 preset-outlined-surface-200-800">
							<Progress.Range class="bg-primary-500" />
						</Progress.Track>
					</Progress>
				</div>

				<div class="flex gap-3">
					<button class="btn preset-filled">Save Changes</button>
					<button class="btn preset-tonal">Reset</button>
				</div>
			</div>
		</Tabs.Content>
	</Tabs>
</div>
