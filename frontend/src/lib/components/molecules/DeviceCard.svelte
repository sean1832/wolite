<script lang="ts">
	import { type Device } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { deviceStore } from '$lib/stores/devices.svelte';
	import { onMount } from 'svelte';
	import { cn } from '$lib/utils.js';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import {
		Monitor,
		MoreVertical,
		Power,
		RotateCw,
		Moon,
		Trash2,
		Link2,
		Loader2,
		Unlink2,
		ChevronDown,
		Snowflake
	} from '@lucide/svelte';
	import PairCompanionDialog from '../organisms/PairCompanionDialog.svelte';
	import EditDeviceDialog from '../organisms/EditDeviceDialog.svelte';

	let {
		device,
		class: className,
		...restProps
	}: {
		device: Device;
		class?: string;
		[key: string]: unknown;
	} = $props();

	onMount(() => {
		if (device.companion_url) {
			deviceStore.checkDeviceStatus(window.fetch, device.mac_address);
		}
	});

	let isEditDialogOpen = $state(false);
	let isPairDialogOpen = $state(false);
	let loading = $state(false);

	// Status color - minimal dot
	let statusColor = $derived(
		device.status === 'online'
			? 'bg-emerald-500 shadow-[0_0_8px_-2px_rgba(16,185,129,0.5)]'
			: device.status === 'offline'
				? 'bg-zinc-300 dark:bg-zinc-700'
				: 'bg-amber-500 animate-pulse'
	);

	let isOnline = $derived(device.status === 'online');

	function handleCardConfig() {
		// Prevent text selection from triggering edit
		if (window.getSelection()?.toString()) return;
		isEditDialogOpen = true;
	}

	async function handleWake() {
		if (isOnline) return;
		loading = true;

		const promise = deviceStore.wakeDevice(fetch, device.mac_address);

		toast.promise(promise, {
			loading: 'Sending wake command...',
			success: `Wake command sent to ${device.name}`,
			error: 'Failed to send wake command'
		});

		try {
			await promise;
		} catch {
			// Error is already logged in store
		} finally {
			loading = false;
		}
	}

	async function handleDelete() {
		try {
			await deviceStore.removeDevice(fetch, device.mac_address);
		} catch {
			// Error is already logged in store
		}
	}

	async function handleUnpair() {
		const promise = deviceStore.unpairCompanion(fetch, device.mac_address);
		toast.promise(promise, {
			loading: 'Unpairing companion...',
			success: 'Companion unpaired',
			error: 'Failed to unpair companion'
		});
	}

	async function handleCompanionAction(action: string) {
		const promise = deviceStore.companionAction(fetch, device.mac_address, action);
		toast.promise(promise, {
			loading: `Sending ${action} command...`,
			success: `Command sent to ${device.name}`,
			error: `Failed to send ${action} command`
		});
	}
</script>

<div
	role="button"
	tabindex="0"
	class={cn(
		'group relative flex min-h-45 flex-col justify-between rounded-xl border border-border/40 bg-card p-5 text-card-foreground shadow-sm transition-all duration-300 hover:bg-accent/30 hover:shadow-md dark:border-transparent dark:shadow-none',
		className
	)}
	onclick={handleCardConfig}
	draggable={true}
	{...restProps}
>
	<!-- Card Header -->
	<div class="flex items-start justify-between">
		<div class="flex items-center gap-3">
			<div
				class={cn(
					'flex h-10 w-10 items-center justify-center rounded-full bg-muted transition-colors',
					statusColor
				)}
			>
				<Monitor class="h-5 w-5 text-foreground" />
			</div>
			<div>
				<h3 class="leading-none font-semibold tracking-tight text-foreground/90">
					{device.name}
				</h3>
				<p class="mt-1 text-xs text-muted-foreground">{device.ip_address}</p>
			</div>
		</div>

		<div class="flex items-start gap-1">
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props }: { props: Record<string, unknown> })}
						<Button
							{...props}
							variant="ghost"
							size="icon"
							class="h-8 w-8 text-muted-foreground/50 hover:text-foreground"
						>
							<MoreVertical class="h-4 w-4" />
							<span class="sr-only">Open menu</span>
						</Button>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					<DropdownMenu.Item onclick={() => (isEditDialogOpen = true)}>
						<span class="ml-6">Edit</span>
					</DropdownMenu.Item>
					<DropdownMenu.Separator />
					<DropdownMenu.Item class="text-destructive focus:text-destructive" onclick={handleDelete}>
						<Trash2 class="mr-2 h-4 w-4" />
						<span>Remove</span>
					</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
	</div>

	<!-- Status Indicator -->
	<div class="mt-4 flex items-end justify-between">
		<div class="flex flex-col gap-1">
			<span class="text-[10px] font-medium text-muted-foreground uppercase">Status</span>
			<div class="flex items-center gap-2">
				<div class={cn('h-2 w-2 rounded-full transition-all duration-500', statusColor)}></div>
				<span class="text-[10px] font-medium tracking-widest text-muted-foreground/60 uppercase"
					>{device.status}</span
				>
			</div>
		</div>

		<div class="flex items-center gap-2">
			<!-- Secondary Action (Pair/Unpair) -->
			{#if device.status !== 'offline'}
				{#if device.companion_auth_fingerprint}
					<Button
						size="sm"
						variant="ghost"
						class="h-8 gap-2 text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
						onclick={(e) => {
							e.stopPropagation();
							handleUnpair();
						}}
					>
						<Unlink2 class="h-3.5 w-3.5" />
						<span class="hidden sm:inline">Unpair</span>
					</Button>
				{:else}
					<Button
						size="sm"
						variant="secondary"
						class="h-8 gap-2 hover:bg-secondary/80"
						onclick={(e) => {
							e.stopPropagation();
							isPairDialogOpen = true;
						}}
					>
						<Link2 class="h-3.5 w-3.5" />
						<span class="hidden sm:inline">Pair</span>
					</Button>
				{/if}
			{/if}

			<!-- Primary Action (Wake/Sleep/Status) -->
			{#if device.status === 'online'}
				{#if device.companion_auth_fingerprint}
					<div class="flex items-center rounded-md border border-primary/20 bg-primary/5 shadow-sm">
						<Button
							size="sm"
							variant="ghost"
							class="h-8 gap-2 rounded-r-none border-r border-primary/10 px-3 text-primary hover:bg-primary/10 hover:text-primary"
							onclick={(e) => {
								e.stopPropagation();
								handleCompanionAction('sleep');
							}}
						>
							<Moon class="h-3.5 w-3.5" />
							<span>Sleep</span>
						</Button>
						<DropdownMenu.Root>
							<DropdownMenu.Trigger>
								{#snippet child({ props }: { props: Record<string, unknown> })}
									<Button
										{...props}
										variant="ghost"
										size="icon"
										class="h-8 w-6 rounded-l-none text-primary hover:bg-primary/10 hover:text-primary"
									>
										<ChevronDown class="h-3.5 w-3.5" />
										<span class="sr-only">More options</span>
									</Button>
								{/snippet}
							</DropdownMenu.Trigger>
							<DropdownMenu.Content align="end">
								<DropdownMenu.Item onclick={() => handleCompanionAction('hibernate')}>
									<Snowflake class="mr-2 h-4 w-4" />
									<span>Hibernate</span>
								</DropdownMenu.Item>
								<DropdownMenu.Item onclick={() => handleCompanionAction('restart')}>
									<RotateCw class="mr-2 h-4 w-4" />
									<span>Restart</span>
								</DropdownMenu.Item>
								<DropdownMenu.Separator />
								<DropdownMenu.Item
									class="text-destructive focus:text-destructive"
									onclick={() => handleCompanionAction('shutdown')}
								>
									<Power class="mr-2 h-4 w-4" />
									<span>Shutdown</span>
								</DropdownMenu.Item>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</div>
				{:else}
					<Button
						size="sm"
						variant="outline"
						class="h-8 gap-2 border-primary/20 bg-primary/5 text-primary hover:bg-primary/10"
						disabled
					>
						<Monitor class="h-3.5 w-3.5" />
						<span>Online</span>
					</Button>
				{/if}
			{:else}
				<Button
					size="sm"
					class={cn(
						'h-8 gap-2 shadow-sm transition-all duration-300 hover:shadow-md',
						'bg-linear-to-r from-primary to-primary/90 hover:to-primary'
					)}
					onclick={(e) => {
						e.stopPropagation();
						handleWake();
					}}
					disabled={loading}
				>
					{#if loading}
						<Loader2 class="h-3.5 w-3.5 animate-spin" />
						<span>Waking...</span>
					{:else}
						<Power class="h-3.5 w-3.5" />
						<span>Wake Up</span>
					{/if}
				</Button>
			{/if}
		</div>
	</div>
</div>

<EditDeviceDialog bind:open={isEditDialogOpen} {device} />
<PairCompanionDialog bind:open={isPairDialogOpen} {device} />
