<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { deviceStore } from '$lib/stores/devices.svelte';
	import type { Device } from '$lib/types';
	import { untrack } from 'svelte';
	import { Textarea } from '../ui/textarea';

	let { open = $bindable(false), device }: { open: boolean; device: Device } = $props();

	let name = $state(untrack(() => device.name));
	let description = $state(untrack(() => device.description || ''));
	let ip_address = $state(untrack(() => device.ip_address));
	let broadcast_ip = $state(untrack(() => device.broadcast_ip));

	// Update local state when device prop changes
	$effect(() => {
		if (device) {
			name = device.name;
			description = device.description || '';
			ip_address = device.ip_address;
			broadcast_ip = device.broadcast_ip;
		}
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		try {
			// Default port to 9 if not specified
            let finalBroadcastIp = broadcast_ip;
            if (finalBroadcastIp && !finalBroadcastIp.includes(':')) {
                finalBroadcastIp += ':9';
            }

			await deviceStore.updateDevice(fetch, device.mac_address, { 
				name, 
				description,
				ip_address, 
				broadcast_ip: finalBroadcastIp 
			});
			open = false;
		} catch (err) {
			// Error is already logged in store
		}
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Edit Device</Dialog.Title>
			<Dialog.Description>
				Make changes to your device here. Click save when you're done.
			</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="grid gap-6 py-4">
			<div class="grid gap-2">
				<Label for="device_name">Device Name</Label>
				<Input
					id="device_name"
					bind:value={name}
					placeholder="e.g. Workstation"
					required
					class="col-span-3"
				/>
			</div>
			<div class="grid gap-2">
				<Label for="description">Description</Label>
				<Textarea id="description" placeholder="e.g. Living Room PC" bind:value={description} class="col-span-3" />
			</div>
			<div class="grid gap-2">
				<Label for="ip_address">IP Address</Label>
				<Input
					id="ip_address"
					bind:value={ip_address}
					placeholder="192.168.1.10"
					required
					class="col-span-3"
				/>
			</div>
			<div class="grid gap-2">
				<Label for="broadcast_ip">Broadcast IP</Label>
				<Input
					id="broadcast_ip"
					bind:value={broadcast_ip}
					placeholder="192.168.1.255:9"
					required
					class="col-span-3"
				/>
			</div>

			<Dialog.Footer>
				<Button type="submit">Save changes</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
