<script lang="ts">
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { deviceStore } from "$lib/stores/devices.svelte";
    import type { Device } from "$lib/types";

    let { open = $bindable(false), device }: { open: boolean, device: Device } = $props();
    
    let name = $state(device.name);
    let ip = $state(device.ip);
    let mac = $state(device.mac);

    // Update local state when device prop changes
    $effect(() => {
        if (device) {
            name = device.name;
            ip = device.ip;
            mac = device.mac;
        }
    });

    function handleSubmit(e: Event) {
        e.preventDefault();
        deviceStore.updateDevice(device.id, { name, ip, mac });
        open = false;
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
                <Button type="submit">Save changes</Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>
