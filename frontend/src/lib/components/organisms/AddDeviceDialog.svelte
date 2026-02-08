<script lang="ts">
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { deviceStore } from "$lib/stores/devices.svelte";
    import { Plus } from '@lucide/svelte';
	import { Textarea } from "../ui/textarea";
    
    let { open = $bindable(false), trigger = undefined } = $props();

    let name = $state('');
    let description = $state('');
    let ip_address = $state('');
    let mac_address = $state('');
    let broadcast_ip = $state('');

    function handleIpInput(e: Event) {
        const input = e.target as HTMLInputElement;
        const value = input.value;
        
        // Simple heuristic: if it looks like an IP, suggest broadcast
        // Matches roughly 3 segments of digits/dots
        const parts = value.split('.');
        if (parts.length >= 3) {
            // Reconstruct subnet
            const subnet = parts.slice(0, 3).join('.');
            if (subnet.length > 0) {
                 broadcast_ip = `${subnet}.255:9`;
            }
        }
    }

    async function handleSubmit(e: Event) {
        e.preventDefault();
        try {
            // Default port to 9 if not specified
            let finalBroadcastIp = broadcast_ip;
            if (finalBroadcastIp && !finalBroadcastIp.includes(':')) {
                finalBroadcastIp += ':9';
            }

            await deviceStore.addDevice(fetch, { 
                name, 
                description,
                mac_address, 
                ip_address, 
                broadcast_ip: finalBroadcastIp 
            });
            open = false;
            name = '';
            description = '';
            ip_address = '';
            mac_address = '';
            broadcast_ip = '';
        } catch (err) {
            // Error is already logged in store
        }
    }
</script>

<Dialog.Root bind:open>
    <Dialog.Trigger>
        {#snippet child({ props })}
            {#if trigger}
                {@render trigger(props)}
            {:else}
                <Button {...props} size="sm" class="h-9 gap-2 px-4 shadow-sm border border-border/50">
                    <Plus class="w-4 h-4" />
                    <span>Add Device</span>
                </Button>
            {/if}
        {/snippet}
    </Dialog.Trigger>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Add Device</Dialog.Title>
            <Dialog.Description>
                Add a new computer to wake up remotely.
            </Dialog.Description>
        </Dialog.Header>
        
        <form onsubmit={handleSubmit} class="grid gap-6 py-4">
            <div class="grid gap-2">
                <Label for="device_name">Device Name</Label>
                <Input id="device_name" bind:value={name} placeholder="e.g. Workstation" required class="col-span-3" />
            </div>
            <div class="grid gap-2">
                <Label for="description">Description</Label>
                <Textarea placeholder="e.g. Living Room PC" bind:value={description} class="col-span-3" />
            </div>
            <div class="grid gap-2">
                <Label for="ip_address">IP Address</Label>
                <Input 
                    id="ip_address" 
                    bind:value={ip_address} 
                    oninput={handleIpInput}
                    placeholder="192.168.1.10" 
                    required 
                    class="col-span-3" 
                />
            </div>
            <div class="grid gap-2">
                <Label for="mac">MAC</Label>
                <Input id="mac" bind:value={mac_address} placeholder="AA:BB:CC:DD:EE:FF" required class="col-span-3" />
            </div>
            <div class="grid gap-2">
                <Label for="broadcast_ip">Broadcast IP</Label>
                <Input id="broadcast_ip" bind:value={broadcast_ip} placeholder="192.168.1.255:9" required class="col-span-3" />
            </div>

            <Dialog.Footer>
                <Button type="submit">Add Device</Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>
