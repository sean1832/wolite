<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Plus } from "@lucide/svelte";
    import { cn } from "$lib/utils";

    let { onclick, class: className } = $props();
</script>

<Button 
    onclick={onclick}
    size="icon"
    class={cn(
        "fixed bottom-6 right-6 h-14 w-14 rounded-full shadow-lg z-50 sm:hidden",
        "bg-primary text-primary-foreground hover:bg-primary/90",
        "transition-all duration-300 hover:scale-105 active:scale-95",
        className
    )}
>
    <Plus class="h-6 w-6" />
</Button>
