<script lang="ts">
  import { type Device } from '$lib/types';
  import { deviceStore } from '$lib/stores/devices.svelte';
  import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
  import { Button } from '$lib/components/ui/button';
<<<<<<< HEAD
  import { MoreHorizontal, Power } from '@lucide/svelte';
=======
  import { MoreVertical, Power } from '@lucide/svelte';
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
  import { cn } from '$lib/utils';
  
  import EditDeviceDialog from '$lib/components/organisms/EditDeviceDialog.svelte';
  
  let { device }: { device: Device } = $props();

  let isEditDialogOpen = $state(false);

<<<<<<< HEAD
  let statusColor = $derived(
    device.status === 'online' ? "bg-emerald-500" : 
    device.status === 'offline' ? "bg-zinc-200 dark:bg-zinc-800" :
=======
  // Status color - minimal dot
  let statusColor = $derived(
    device.status === 'online' ? "bg-emerald-500 shadow-[0_0_8px_-2px_rgba(16,185,129,0.5)]" : 
    device.status === 'offline' ? "bg-zinc-300 dark:bg-zinc-700" :
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
    "bg-amber-500 animate-pulse"
  );

  let isOnline = $derived(device.status === 'online');

<<<<<<< HEAD
  function handleWake() {
      if (isOnline) return;
      deviceStore.wakeDevice(device.id);
  }
</script>

<div class="group relative bg-background border border-border/40 p-5 rounded-xl transition-all duration-300 hover:shadow-sm hover:border-border/80">
  <div class="flex justify-between items-start mb-4">
    <div class="flex items-center gap-2.5">
        <div class={cn("w-1.5 h-1.5 rounded-full transition-colors duration-500", statusColor)}></div>
        <span class="text-[10px] uppercase tracking-widest text-muted-foreground/70 font-medium">{device.status}</span>
=======
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
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
    </div>
    
    <DropdownMenu.Root>
      <DropdownMenu.Trigger>
        {#snippet child({ props })}
<<<<<<< HEAD
            <Button variant="ghost" size="icon" class="h-6 w-6 -mr-2 text-muted-foreground/50 hover:text-foreground" {...props}>
                <MoreHorizontal class="w-4 h-4" />
                <span class="sr-only">Menu</span>
            </Button>
        {/snippet}
      </DropdownMenu.Trigger>
      <DropdownMenu.Content align="end" class="w-32">
        <DropdownMenu.Item onclick={() => isEditDialogOpen = true}>Edit</DropdownMenu.Item>
        <DropdownMenu.Separator />
        <DropdownMenu.Item class="text-destructive focus:text-destructive" onclick={() => deviceStore.removeDevice(device.id)}>
=======
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
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
            Delete
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  </div>

  <EditDeviceDialog bind:open={isEditDialogOpen} {device} />

<<<<<<< HEAD
  <div class="space-y-1 mb-6">
    <h3 class="font-medium text-sm tracking-tight text-foreground">{device.name}</h3>
    <div class="flex flex-col gap-0.5">
        <code class="text-[10px] text-muted-foreground/60 font-mono">{device.ip}</code>
    </div>
  </div>

  <Button 
      variant="outline" 
      class={cn(
          "w-full h-9 justify-between group/btn transition-all duration-300 border-input bg-background hover:bg-accent hover:text-accent-foreground hover:border-primary/30 dark:hover:border-primary/50", 
          isOnline && "opacity-50 cursor-not-allowed pointer-events-none"
      )}
      onclick={handleWake}
      disabled={isOnline}
  >
      <span class="text-xs font-medium">{isOnline ? 'Online' : 'Wake'}</span>
      <Power class={cn("w-3.5 h-3.5 transition-transform duration-300 group-hover/btn:scale-110", isOnline ? "text-emerald-500" : "")} />
  </Button>
=======
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
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
</div>
