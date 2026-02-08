<script lang="ts">
	import { deviceStore } from '$lib/stores/devices.svelte';
	import DeviceCard from '$lib/components/molecules/DeviceCard.svelte';
	import AddDeviceDialog from '$lib/components/organisms/AddDeviceDialog.svelte';
	import ThemeToggle from '$lib/components/atoms/ThemeToggle.svelte';
	import Header from '$lib/components/organisms/Header.svelte';
	import { Button } from '$lib/components/ui/button';
	import { UserIcon } from '@lucide/svelte';
	import { onMount } from 'svelte';

	onMount(async () => {
		await deviceStore.init(fetch);
	});
</script>

<div class="container mx-auto max-w-6xl space-y-8 px-6 py-20">
	<Header title="Wolite" subtitle="Network Power Control">
		<ThemeToggle />
		<a href="/account">
			<Button variant="outline" size="icon" class="h-9 w-9">
				<UserIcon class="h-4 w-4" />
			</Button>
		</a>
		<AddDeviceDialog />
	</Header>

	<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
		{#each deviceStore.devices as device (device.mac_address)}
			<DeviceCard {device} />
		{/each}

		{#if deviceStore.devices.length === 0}
			<div
				class="col-span-full flex flex-col items-center justify-center space-y-4 py-24 text-center text-muted-foreground"
			>
				<div class="flex h-12 w-12 items-center justify-center rounded-full bg-muted/30">
					<span class="text-xl opacity-20">?</span>
				</div>
				<div>
					<p class="font-medium text-foreground">No devices yet</p>
					<p class="text-sm opacity-50">Add your first machine to get started.</p>
				</div>
			</div>
		{/if}
	</div>
</div>
