<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { ArrowLeft } from "@lucide/svelte";
    
    let { title, subtitle, backHref = undefined, children } = $props();
</script>

<header class="flex items-end justify-between border-b border-border/10 pb-6 mb-8">
    <div class="flex items-center gap-4">
        {#if backHref}
            <a href={backHref} class="-ml-2">
                <Button variant="ghost" size="icon" class="h-8 w-8 text-muted-foreground/50 hover:text-foreground">
                    <ArrowLeft class="h-5 w-5" />
                </Button>
            </a>
        {/if}
        <div class="space-y-1">
            <h1 class="text-xl font-medium tracking-tight text-foreground">{title}</h1>
            {#if subtitle}
                <p class="text-sm text-muted-foreground/60">{subtitle}</p>
            {/if}
        </div>
    </div>
    
    <div class="flex items-center gap-2">
        {@render children?.()}
    </div>
</header>
