<script lang="ts">
	import { type Device } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { deviceStore } from '$lib/stores/devices.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Button } from '$lib/components/ui/button';
	import { MoreVertical, Power } from '@lucide/svelte';
	import { cn } from '$lib/utils';

	import EditDeviceDialog from '$lib/components/organisms/EditDeviceDialog.svelte';
	import PairCompanionDialog from '$lib/components/organisms/PairCompanionDialog.svelte';

	let { device }: { device: Device } = $props();

	let isEditDialogOpen = $state(false);
	let isPairDialogOpen = $state(false);

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

	async function handleWake(e: Event) {
		e.stopPropagation();
		if (isOnline) return;

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

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	class="group relative cursor-pointer rounded-sm border border-border/40 bg-card p-5 text-card-foreground shadow-sm transition-all duration-300 hover:bg-accent/30 hover:shadow-md dark:border-transparent dark:shadow-none"
	onclick={handleCardConfig}
>
	<div class="mb-6 flex items-start justify-between">
		<div class="flex items-center gap-2">
			<div class={cn('h-1.5 w-1.5 rounded-full transition-all duration-500', statusColor)}></div>
			<span class="text-[10px] font-medium tracking-widest text-muted-foreground/60 uppercase"
				>{device.status}</span
			>
		</div>

		<DropdownMenu.Root>
			<div onclick={(e) => e.stopPropagation()}>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Button
							variant="ghost"
							size="icon"
							class="-mr-2 h-6 w-6 text-muted-foreground/30 transition-colors hover:text-foreground"
							{...props}
						>
							<MoreVertical class="h-3.5 w-3.5" />
							<span class="sr-only">Menu</span>
						</Button>
					{/snippet}
				</DropdownMenu.Trigger>
			</div>
			<DropdownMenu.Content align="end" class="w-48 border-border/40 font-light">
				<DropdownMenu.Item onclick={() => (isEditDialogOpen = true)} class="text-xs">
					Edit
				</DropdownMenu.Item>

				{#if device.companion_url}
					<DropdownMenu.Sub>
						<DropdownMenu.SubTrigger class="text-xs">Power Actions</DropdownMenu.SubTrigger>
						<DropdownMenu.SubContent class="w-32 border-border/40">
							<DropdownMenu.Item onclick={() => handleCompanionAction('sleep')} class="text-xs">
								Sleep
							</DropdownMenu.Item>
							<DropdownMenu.Item onclick={() => handleCompanionAction('restart')} class="text-xs">
								Restart
							</DropdownMenu.Item>
							<DropdownMenu.Item onclick={() => handleCompanionAction('hibernate')} class="text-xs">
								Hibernate
							</DropdownMenu.Item>
							<DropdownMenu.Separator class="bg-border/30" />
							<DropdownMenu.Item
								onclick={() => handleCompanionAction('shutdown')}
								class="text-xs text-destructive focus:text-destructive"
							>
								Shutdown
							</DropdownMenu.Item>
						</DropdownMenu.SubContent>
					</DropdownMenu.Sub>
					<DropdownMenu.Item onclick={handleUnpair} class="text-xs text-muted-foreground">
						Unpair Companion
					</DropdownMenu.Item>
				{:else}
					<DropdownMenu.Item onclick={() => (isPairDialogOpen = true)} class="text-xs">
						Pair Companion
					</DropdownMenu.Item>
				{/if}

				<DropdownMenu.Separator class="bg-border/30" />
				<DropdownMenu.Item
					class="text-xs text-destructive focus:text-destructive"
					onclick={handleDelete}
				>
					Delete
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</div>

	<EditDeviceDialog bind:open={isEditDialogOpen} {device} />
	<PairCompanionDialog bind:open={isPairDialogOpen} {device} />

	<div class="mb-6 space-y-3">
		<div>
			<h3 class="text-base font-normal tracking-tight text-foreground">{device.name}</h3>
			{#if device.description}
				<p class="mt-0.5 max-w-[90%] truncate text-[11px] text-muted-foreground/70">
					{device.description}
				</p>
			{/if}
		</div>
		<div class="flex flex-col gap-0.5">
			<code class="font-mono text-[10px] tracking-wide text-muted-foreground/40"
				>{device.ip_address}</code
			>
			<code class="font-mono text-[10px] tracking-wide text-muted-foreground/30"
				>{device.mac_address}</code
			>
		</div>
	</div>

	<div class="flex items-center justify-end">
		{#if !isOnline}
			<Button
				variant="ghost"
				size="sm"
				class="group/wake h-8 w-full gap-1.5 border px-3 text-xs font-medium text-foreground/80 transition-all duration-300 hover:bg-primary/5 hover:text-primary md:w-auto"
				onclick={handleWake}
			>
				<span class="hidden md:block">Wake</span>
				<Power class="h-3 w-3 transition-colors group-hover/wake:text-primary" />
			</Button>
		{/if}
	</div>
</div>
