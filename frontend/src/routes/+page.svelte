<script lang="ts">
	import { deviceStore } from '$lib/stores/devices.svelte';
	import DeviceCard from '$lib/components/molecules/DeviceCard.svelte';
	import DeviceCardSkeleton from '$lib/components/molecules/DeviceCardSkeleton.svelte';
	import AddDeviceDialog from '$lib/components/organisms/AddDeviceDialog.svelte';
	import Header from '$lib/components/organisms/Header.svelte';
	import FloatingActionButton from '$lib/components/atoms/FloatingActionButton.svelte';
    import { Button } from "$lib/components/ui/button";
    import { Plus } from "@lucide/svelte";
	import { onMount } from 'svelte';

	let isAddDialogOpen = $state(false);

	onMount(async () => {
		await deviceStore.init(fetch);
	});
</script>

<div class="container mx-auto max-w-5xl px-6 min-h-[calc(100vh-4rem)] flex flex-col">
	<div class="flex flex-col gap-6 flex-1">
		<Header title="Wolite" subtitle="Control Center">
            <AddDeviceDialog bind:open={isAddDialogOpen}>
                {#snippet trigger(props: any)}
                    <Button variant="default" size="icon" class="hidden sm:flex h-9 w-auto shadow-sm hover:shadow-md transition-all px-2" {...props}>
                        <span>+ Add Device</span>
                    </Button>
                {/snippet}
            </AddDeviceDialog>
        </Header>

		<div class="grid grid-cols-1 gap-4 md:gap-6 sm:grid-cols-2 lg:grid-cols-3">
			{#if deviceStore.loading && deviceStore.devices.length === 0}
				{#each Array(3) as _}
					<DeviceCardSkeleton />
				{/each}
			{:else}
				{#each deviceStore.devices as device (device.mac_address)}
					<DeviceCard {device} />
				{/each}

				{#if deviceStore.devices.length === 0}
                    <!-- Empty state is less relevant now that we have the Add Card, but keeping it for mobile or if list empty -->
                    <div class="col-span-full flex flex-col items-center justify-center space-y-6 py-12 text-center sm:hidden">
						<div class="h-px w-24 bg-border/40"></div>
						<div>
							<p class="text-sm font-medium text-muted-foreground/60 uppercase tracking-widest">No devices</p>
						</div>
						<div class="h-px w-24 bg-border/40"></div>
					</div>
				{/if}
			{/if}
		</div>
	</div>
</div>

<div class="sm:hidden">
	<AddDeviceDialog bind:open={isAddDialogOpen}>
		{#snippet trigger(props: any)}
			<FloatingActionButton {...props} />
		{/snippet}
	</AddDeviceDialog>
</div>
