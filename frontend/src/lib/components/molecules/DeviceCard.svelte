<script lang="ts">
  import { type Device } from '$lib/types';
  import { deviceStore } from '$lib/stores/devices.svelte';
  import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
  import { Button } from '$lib/components/ui/button';
  import { MoreVertical, Power } from '@lucide/svelte';
  import { cn } from '$lib/utils';
  
  import EditDeviceDialog from '$lib/components/organisms/EditDeviceDialog.svelte';
  
  let { device }: { device: Device } = $props();

  let isEditDialogOpen = $state(false);

  // Status color - minimal dot
  let statusColor = $derived(
    device.status === 'online' ? "bg-emerald-500 shadow-[0_0_8px_-2px_rgba(16,185,129,0.5)]" : 
    device.status === 'offline' ? "bg-zinc-300 dark:bg-zinc-700" :
    "bg-amber-500 animate-pulse"
  );

  let isOnline = $derived(device.status === 'online');

  async function handleWake() {
      if (isOnline) return;
      try {
        await deviceStore.wakeDevice(fetch, device.mac_address);
      } catch (err) {
        // Error is already logged in store
      }
  }

  async function handleDelete() {
    try {
        await deviceStore.removeDevice(fetch, device.mac_address);
    } catch (err) {
        // Error is already logged in store
    }
  }
</script>

<div class="group relative bg-card text-card-foreground p-5 rounded-sm transition-all duration-300 hover:bg-accent/30 border border-border/40 dark:border-transparent shadow-sm hover:shadow-md dark:shadow-none">
  <div class="flex justify-between items-start mb-6">
    <div class="flex items-center gap-2">
        <div class={cn("w-1.5 h-1.5 rounded-full transition-all duration-500", statusColor)}></div>
        <span class="text-[10px] uppercase tracking-widest text-muted-foreground/60 font-medium">{device.status}</span>
    </div>
    
    <DropdownMenu.Root>
      <DropdownMenu.Trigger>
        {#snippet child({ props })}
        <Button 
            variant="ghost" 
            size="icon" 
            class="h-6 w-6 -mr-2 text-muted-foreground/30 hover:text-foreground transition-colors" 
            {...props}
        >
            <MoreVertical class="w-3.5 h-3.5" />
            <span class="sr-only">Menu</span>
        </Button>
        {/snippet}
      </DropdownMenu.Trigger>
      <DropdownMenu.Content align="end" class="w-32 border-border/40 font-light">
        <DropdownMenu.Item onclick={() => isEditDialogOpen = true} class="text-xs">Edit</DropdownMenu.Item>
        <DropdownMenu.Separator class="bg-border/30" />
        <DropdownMenu.Item class="text-destructive focus:text-destructive text-xs" onclick={handleDelete}>
            Delete
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  </div>

  <EditDeviceDialog bind:open={isEditDialogOpen} {device} />

  <div class="space-y-3 mb-6">
    <div>
        <h3 class="font-normal text-base tracking-tight text-foreground">{device.name}</h3>
        {#if device.description}
            <p class="text-[11px] text-muted-foreground/70 truncate mt-0.5 max-w-[90%]">{device.description}</p>
        {/if}
    </div>
    <div class="flex flex-col gap-0.5">
        <code class="text-[10px] text-muted-foreground/40 font-mono tracking-wide">{device.ip_address}</code>
    </div>
  </div>

  <div class="flex items-center justify-end">
    {#if !isOnline}
        <Button 
            variant="ghost" 
            size="sm"
            class="h-8 px-3 text-xs font-medium text-foreground/80 hover:text-primary hover:bg-primary/5 transition-all duration-300 gap-1.5 group/wake w-full md:w-auto border"
            onclick={handleWake}
        >
            <span class="hidden md:block">Wake</span>
            <Power class="w-3 h-3 group-hover/wake:text-primary transition-colors" />
        </Button>
    {/if}
  </div>
</div>
