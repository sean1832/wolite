<script lang="ts">
  import { type Device } from '$lib/types';
  import { deviceStore } from '$lib/stores/devices.svelte';
  import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
  import { Button } from '$lib/components/ui/button';
  import { MoreHorizontal, Power } from '@lucide/svelte';
  import { cn } from '$lib/utils';
  
  import EditDeviceDialog from '$lib/components/organisms/EditDeviceDialog.svelte';
  
  let { device }: { device: Device } = $props();

  let isEditDialogOpen = $state(false);

  let statusColor = $derived(
    device.status === 'online' ? "bg-emerald-500" : 
    device.status === 'offline' ? "bg-zinc-200 dark:bg-zinc-800" :
    "bg-amber-500 animate-pulse"
  );

  let isOnline = $derived(device.status === 'online');

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
    </div>
    
    <DropdownMenu.Root>
      <DropdownMenu.Trigger>
        {#snippet child({ props })}
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
            Delete
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  </div>

  <EditDeviceDialog bind:open={isEditDialogOpen} {device} />

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
</div>
