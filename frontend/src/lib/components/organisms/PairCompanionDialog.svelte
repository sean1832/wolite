<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { deviceStore } from '$lib/stores/devices.svelte';
	import type { Device } from '$lib/types';
	import { toast } from 'svelte-sonner';

	let { open = $bindable(false), device }: { open: boolean; device: Device } = $props();

	let url = $state('');
	let token = $state('');
	let isLoading = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		isLoading = true;

		const promise = deviceStore.pairCompanion(fetch, device.mac_address, url, token);

		toast.promise(promise, {
			loading: 'Pairing with companion...',
			success: 'Companion paired successfully',
			error: (err: unknown) =>
				err instanceof Error ? `Failed to pair: ${err.message}` : 'Failed to pair'
		});

		try {
			await promise;
			open = false;
			// Reset form
			url = '';
			token = '';
		} catch {
			// Error handled by toast
		} finally {
			isLoading = false;
		}
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Pair Companion App</Dialog.Title>
			<Dialog.Description>
				Connect your device's companion app to enable advanced control.
			</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="grid gap-6 py-4">
			<div class="grid gap-2">
				<Label for="url">Companion URL</Label>
				<Input
					id="url"
					bind:value={url}
					placeholder="https://192.168.1.50:8443"
					required
					class="col-span-3"
				/>
				<p class="text-[10px] text-muted-foreground">Must be reachable from the Wolite server.</p>
			</div>
			<div class="grid gap-2">
				<Label for="token">Access Token</Label>
				<Input
					id="token"
					type="password"
					bind:value={token}
					placeholder="eyJhbGciOiJIUzI1Ni..."
					required
					class="col-span-3"
				/>
			</div>

			<Dialog.Footer>
				<Button type="submit" disabled={isLoading}>
					{isLoading ? 'Pairing...' : 'Pair Device'}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
