<script lang="ts">
	import { type Device } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { deviceStore } from '$lib/stores/devices.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { Button } from '$lib/components/ui/button';
	import { MoreVertical, Power, Moon, RefreshCw, Snowflake, Link } from '@lucide/svelte';
	import { cn } from '$lib/utils';
	import { fade } from 'svelte/transition';

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
	class="group relative flex min-h-45 cursor-pointer flex-col justify-between rounded-xl border border-border/40 bg-card p-5 text-card-foreground shadow-sm transition-all duration-300 hover:bg-accent/30 hover:shadow-md dark:border-transparent dark:shadow-none"
	onclick={handleCardConfig}
>
	<div>
		<div class="mb-4 flex items-start justify-between">
			<div class="flex items-center gap-2">
				<div class={cn('h-2 w-2 rounded-full transition-all duration-500', statusColor)}></div>
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
								class="-mr-2 h-8 w-8 text-muted-foreground/40 transition-colors hover:bg-transparent hover:text-foreground"
								{...props}
							>
								<MoreVertical class="h-4 w-4" />
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
						<DropdownMenu.Item onclick={handleUnpair} class="text-xs text-muted-foreground">
							Unpair Companion
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

		<div class="mb-6 space-y-2">
			<div>
				<h3
					class="truncate pr-2 text-lg font-medium tracking-tight text-foreground"
					title={device.name}
				>
					{device.name}
				</h3>
				{#if device.description}
					<p class="mt-1 line-clamp-2 h-8 text-xs text-muted-foreground/70">
						{device.description}
					</p>
				{:else}
					<div class="h-8"></div>
				{/if}
			</div>
			<div class="flex flex-col gap-0.5 pt-1">
				<code class="font-mono text-[10px] tracking-wide text-muted-foreground/50"
					>{device.ip_address}</code
				>
				<code class="font-mono text-[10px] tracking-wide text-muted-foreground/40"
					>{device.mac_address}</code
				>
			</div>
		</div>
	</div>

	<div
		class="mt-auto flex flex-col gap-2 border-t border-border/20 pt-4"
		onclick={(e) => e.stopPropagation()}
	>
		{#if !isOnline}
			<Button
				variant="outline"
				size="default"
				class="w-full gap-2 border-primary/20 bg-primary/5 transition-all duration-300 hover:bg-primary/10 hover:text-primary"
				onclick={handleWake}
			>
				<Power class="h-3.5 w-3.5" />
				<span>Wake</span>
			</Button>
		{/if}

		{#if isOnline && device.companion_url}
			<!-- Mobile: Dropdown Menu -->
			<div class="flex w-full md:hidden" transition:fade onclick={(e) => e.stopPropagation()}>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Button
								variant="secondary"
								size="default"
								class="w-full gap-2 text-muted-foreground transition-all duration-300 hover:text-foreground"
								{...props}
							>
								<Power class="h-4 w-4" />
								<span>Power Actions</span>
							</Button>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content align="center" class="w-48 border-border/40 font-light">
						<DropdownMenu.Item onclick={() => handleCompanionAction('sleep')} class="text-sm">
							<Moon class="mr-2 h-4 w-4" />
							Sleep
						</DropdownMenu.Item>
						<DropdownMenu.Item onclick={() => handleCompanionAction('restart')} class="text-sm">
							<RefreshCw class="mr-2 h-4 w-4" />
							Restart
						</DropdownMenu.Item>
						<DropdownMenu.Item onclick={() => handleCompanionAction('hibernate')} class="text-sm">
							<Snowflake class="mr-2 h-4 w-4" />
							Hibernate
						</DropdownMenu.Item>
						<DropdownMenu.Separator class="bg-border/30" />
						<DropdownMenu.Item
							onclick={() => handleCompanionAction('shutdown')}
							class="text-sm text-destructive focus:text-destructive"
						>
							<Power class="mr-2 h-4 w-4" />
							Shutdown
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</div>

			<!-- Desktop: Icon Buttons with Tooltips -->
			<div class="hidden w-full items-center justify-between gap-1 md:flex" transition:fade>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								variant="ghost"
								size="icon"
								class="h-10 w-10 text-muted-foreground hover:bg-accent/50 hover:text-foreground"
								{...props}
								onclick={() => handleCompanionAction('sleep')}
							>
								<Moon class="h-5 w-5" />
								<span class="sr-only">Sleep</span>
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Sleep</Tooltip.Content>
				</Tooltip.Root>

				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								variant="ghost"
								size="icon"
								class="h-10 w-10 text-muted-foreground hover:bg-accent/50 hover:text-foreground"
								{...props}
								onclick={() => handleCompanionAction('restart')}
							>
								<RefreshCw class="h-5 w-5" />
								<span class="sr-only">Restart</span>
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Restart</Tooltip.Content>
				</Tooltip.Root>

				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								variant="ghost"
								size="icon"
								class="h-10 w-10 text-muted-foreground hover:bg-accent/50 hover:text-foreground"
								{...props}
								onclick={() => handleCompanionAction('hibernate')}
							>
								<Snowflake class="h-5 w-5" />
								<span class="sr-only">Hibernate</span>
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Hibernate</Tooltip.Content>
				</Tooltip.Root>

				<div class="mx-1 h-6 w-px bg-border/40"></div>

				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								variant="ghost"
								size="icon"
								class="h-10 w-10 text-destructive/70 hover:bg-destructive/10 hover:text-destructive"
								{...props}
								onclick={() => handleCompanionAction('shutdown')}
							>
								<Power class="h-5 w-5" />
								<span class="sr-only">Shutdown</span>
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Shutdown</Tooltip.Content>
				</Tooltip.Root>
			</div>
		{/if}

		{#if !device.companion_url}
			<Button
				variant="secondary"
				size="default"
				class="w-full gap-2 text-muted-foreground transition-all duration-300 hover:text-foreground"
				onclick={() => (isPairDialogOpen = true)}
			>
				<Link class="h-3.5 w-3.5" />
				<span>Pair Companion</span>
			</Button>
		{/if}
	</div>
</div>
