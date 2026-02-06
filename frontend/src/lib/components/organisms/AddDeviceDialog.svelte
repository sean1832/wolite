<script lang="ts">
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { deviceStore } from "$lib/stores/devices.svelte";
    import { Plus } from '@lucide/svelte';
    
    let open = $state(false);
    let name = $state('');
    let ip = $state('');
    let mac = $state('');

    function handleSubmit(e: Event) {
        e.preventDefault();
        deviceStore.addDevice({ name, ip, mac });
        open = false;
        name = '';
        ip = '';
        mac = '';
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
                <Label for="ip">IP Address</Label>
                <Input id="ip" bind:value={ip} placeholder="192.168.1.10" required class="col-span-3" />
            </div>
            <div class="grid gap-2">
                <Label for="mac">MAC Address</Label>
                <Input id="mac" bind:value={mac} placeholder="AA:BB:CC:DD:EE:FF" required class="col-span-3" />
            </div>

            <Dialog.Footer>
                <Button type="submit">Add Device</Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>
