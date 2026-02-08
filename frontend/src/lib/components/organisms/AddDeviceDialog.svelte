<script lang="ts">
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { deviceStore } from "$lib/stores/devices.svelte";
    import { Plus } from '@lucide/svelte';
    
    let open = $state(false);
    let name = $state('');
    let ip_address = $state('');
    let mac_address = $state('');
    let broadcast_ip = $state('');

    async function handleSubmit(e: Event) {
        e.preventDefault();
        try {
            await deviceStore.addDevice(fetch, { 
                name, 
                mac_address, 
                ip_address, 
                broadcast_ip 
            });
            open = false;
            name = '';
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
            <Button {...props} size="sm" class="h-9 gap-2 px-4 shadow-sm">
                <Plus class="w-4 h-4" />
                <span>Add Device</span>
            </Button>
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
                <Label for="name">Name</Label>
                <Input id="name" bind:value={name} placeholder="e.g. Workstation" required class="col-span-3" />
            </div>
            <div class="grid gap-2">
                <Label for="ip_address">IP Address</Label>
                <Input id="ip_address" bind:value={ip_address} placeholder="192.168.1.10" required class="col-span-3" />
            </div>
            <div class="grid gap-2">
                <Label for="mac_address">MAC Address</Label>
                <Input id="mac_address" bind:value={mac_address} placeholder="AA:BB:CC:DD:EE:FF" required class="col-span-3" />
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
