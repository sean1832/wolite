<script lang="ts">
	import { deviceStore } from '$lib/stores/devices.svelte';
	import DeviceCard from '$lib/components/molecules/DeviceCard.svelte';
	import DeviceCardSkeleton from '$lib/components/molecules/DeviceCardSkeleton.svelte';
	import AddDeviceDialog from '$lib/components/organisms/AddDeviceDialog.svelte';
	import Header from '$lib/components/organisms/Header.svelte';
	import FloatingActionButton from '$lib/components/atoms/FloatingActionButton.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { onMount } from 'svelte';
	import { cn } from '$lib/utils'; // Import cn

	let isAddDialogOpen = $state(false);

	onMount(async () => {
		await deviceStore.init(fetch);
	});
</script>

<div class="container mx-auto flex min-h-[calc(100vh-4rem)] max-w-5xl flex-col px-6">
	<div class="flex flex-1 flex-col gap-6">
		<Header title="Wolite" subtitle="Control Center">
			<AddDeviceDialog bind:open={isAddDialogOpen}>
				<!-- eslint-disable-next-line @typescript-eslint/no-explicit-any -->
				{#snippet trigger(props: any)}
					<Button
						variant="default"
						size="icon"
						class="hidden h-9 w-auto px-2 shadow-sm transition-all hover:shadow-md sm:flex"
						{...props}
					>
						<span>+ Add Device</span>
					</Button>
				{/snippet}
			</AddDeviceDialog>
		</Header>

		<Tooltip.Provider>
			<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:gap-6 lg:grid-cols-3">
				{#if deviceStore.loading && deviceStore.devices.length === 0}
					<!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
					{#each Array(3) as _, i (i)}
						<DeviceCardSkeleton />
					{/each}
				{:else}
					{#each deviceStore.devices as device, index (device.mac_address)}
						<DeviceCard
							{device}
							ondragstart={(e: DragEvent) => {
								e.dataTransfer?.setData('text/plain', index.toString());
								e.dataTransfer!.effectAllowed = 'move';
							}}
							ondragover={(e: DragEvent) => {
								e.preventDefault(); // Allow drop
								e.dataTransfer!.dropEffect = 'move';
							}}
							ondrop={(e: DragEvent) => {
								e.preventDefault();
								const fromIndexStr = e.dataTransfer!.getData('text/plain');
								if (!fromIndexStr) return;
								const fromIndex = parseInt(fromIndexStr);
								const toIndex = index;

								if (fromIndex === toIndex) return;

								const newOrder = [...deviceStore.devices];
								const [moved] = newOrder.splice(fromIndex, 1);
								newOrder.splice(toIndex, 0, moved);

								// Extract MACs for API
								const newOrderMacs = newOrder.map((d) => d.mac_address);
								deviceStore.reorderDevices(fetch, newOrderMacs);
							}}
							class={cn(
								'cursor-default',
								'transition-opacity',
								deviceStore.loading ? 'opacity-50' : ''
							)}
						/>
					{/each}

					{#if deviceStore.devices.length === 0}
						<!-- Empty state is less relevant now that we have the Add Card, but keeping it for mobile or if list empty -->
						<div
							class="col-span-full flex flex-col items-center justify-center space-y-6 py-12 text-center sm:hidden"
						>
							<div class="h-px w-24 bg-border/40"></div>
							<div>
								<p class="text-sm font-medium tracking-widest text-muted-foreground/60 uppercase">
									No devices
								</p>
							</div>
							<div class="h-px w-24 bg-border/40"></div>
						</div>
					{/if}
				{/if}
			</div>
		</Tooltip.Provider>
	</div>
</div>

<div class="sm:hidden">
	<AddDeviceDialog bind:open={isAddDialogOpen}>
		<!-- eslint-disable-next-line @typescript-eslint/no-explicit-any -->
		{#snippet trigger(props: any)}
			<FloatingActionButton {...props} />
		{/snippet}
	</AddDeviceDialog>
</div>
