<script lang="ts">
    import { Button } from "$lib/components/ui/button";
<<<<<<< HEAD
    import { ArrowLeft } from "@lucide/svelte";
    
    let { title, subtitle, backHref = undefined, children } = $props();
</script>

<header class="flex items-end justify-between border-b border-border/10 pb-6 mb-8">
=======
    import * as Avatar from "$lib/components/ui/avatar";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
    import { ArrowLeft, User, Sun, Moon, LogOut, Settings } from "@lucide/svelte";
    import { toggleMode } from "mode-watcher";
    import { goto } from "$app/navigation";
    import { authStore } from "$lib/stores/auth.svelte";
    
    let { title, subtitle, backHref = undefined, children = undefined, showAccount = true } = $props();
</script>

<header class="flex items-end justify-between py-6 mb-8 lg:mb-12">
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
    <div class="flex items-center gap-4">
        {#if backHref}
            <a href={backHref} class="-ml-2">
                <Button variant="ghost" size="icon" class="h-8 w-8 text-muted-foreground/50 hover:text-foreground">
                    <ArrowLeft class="h-5 w-5" />
                </Button>
            </a>
        {/if}
<<<<<<< HEAD
        <div class="space-y-1">
            <h1 class="text-xl font-medium tracking-tight text-foreground">{title}</h1>
            {#if subtitle}
                <p class="text-sm text-muted-foreground/60">{subtitle}</p>
=======
        <div class="space-y-0.5">
            <h1 class="text-2xl font-light tracking-tight text-foreground">{title}</h1>
            {#if subtitle}
                <p class="text-xs text-muted-foreground/50 uppercase tracking-widest">{subtitle}</p>
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
            {/if}
        </div>
    </div>
    
<<<<<<< HEAD
    <div class="flex items-center gap-2">
        {@render children?.()}
=======
    <div class="flex items-center gap-4">
        <div class="hidden sm:flex items-center gap-2">
            {@render children?.()}
        </div>
        
        {#if showAccount}
        <DropdownMenu.Root>
            <DropdownMenu.Trigger class="focus:outline-none">
                <Avatar.Root class="h-8 w-8 transition-opacity hover:opacity-80">
                    <Avatar.Fallback class="bg-muted text-muted-foreground text-xs"><User class="w-4 h-4" /></Avatar.Fallback>
                </Avatar.Root>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content align="end" class="w-48">
                <DropdownMenu.Label>My Account</DropdownMenu.Label>
                <DropdownMenu.Separator />
                <DropdownMenu.Item onclick={() => goto('/account')}>
                    <Settings class="mr-2 h-4 w-4" />
                    <span>Settings</span>
                </DropdownMenu.Item>
                <DropdownMenu.Item onclick={toggleMode}>
                    <Sun class="mr-2 h-4 w-4 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
                    <Moon class="absolute mr-2 h-4 w-4 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
                    <span>Toggle Theme</span>
                </DropdownMenu.Item>
                <DropdownMenu.Separator />
                <DropdownMenu.Item class="text-destructive focus:text-destructive" onclick={() => authStore.logout(fetch)}>
                    <LogOut class="mr-2 h-4 w-4" />
                    <span>Sign out</span>
                </DropdownMenu.Item>
            </DropdownMenu.Content>
        </DropdownMenu.Root>
        {/if}
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
    </div>
</header>
