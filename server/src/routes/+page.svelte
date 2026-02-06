<script lang="ts">
  import { deviceStore } from '$lib/stores/devices.svelte';
  import DeviceCard from '$lib/components/molecules/DeviceCard.svelte';
  import AddDeviceDialog from '$lib/components/organisms/AddDeviceDialog.svelte';
  import ThemeToggle from '$lib/components/atoms/ThemeToggle.svelte';
</script>

<div class="container max-w-6xl mx-auto py-20 px-6 space-y-16">
    <header class="flex items-end justify-between border-b border-border/10 pb-6">
        <div class="space-y-1">
            <h1 class="text-xl font-medium tracking-tight text-foreground">Wolite</h1>
            <p class="text-sm text-muted-foreground/60">Network Power Control</p>
        </div>
        <div class="flex items-center gap-2">
            <ThemeToggle />
            <AddDeviceDialog />
        </div>
    </header>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each deviceStore.devices as device (device.id)}
            <DeviceCard {device} />
        {/each}
        
        {#if deviceStore.devices.length === 0}
            <div class="col-span-full flex flex-col items-center justify-center py-24 text-center text-muted-foreground space-y-4">
                <div class="w-12 h-12 rounded-full bg-muted/30 flex items-center justify-center">
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
