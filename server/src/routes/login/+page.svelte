<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { LogIn, KeyRound } from '@lucide/svelte';
    import { enhance } from '$app/forms';

    let { form }: { form?: { missing?: boolean; invalid?: boolean; otpRequired?: boolean; otpInvalid?: boolean; username?: string } } = $props();
    
    // If server requested OTP, show it immediately
    let showOTP = $derived(form?.otpRequired || false);
    let username = $derived(form?.username || ''); 
    // We bind to a local variable to initiate it, but update if form returns one
    let usernameInput = $state('');

    $effect(() => {
        if (username) usernameInput = username;
    });
</script>

<div class="h-screen w-full flex items-center justify-center bg-zinc-50 dark:bg-zinc-950">
    <div class="w-full max-w-sm p-8 space-y-8 bg-background border border-border/40 rounded-xl shadow-sm">
        <div class="flex flex-col items-center text-center space-y-2">
            <div class="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-2">
                <LogIn class="w-6 h-6 text-primary" />
            </div>
            <h1 class="text-2xl font-semibold tracking-tight">Sign in</h1>
            <p class="text-sm text-muted-foreground">Enter your credentials to access the dashboard.</p>
        </div>

        <form method="POST" class="space-y-6" use:enhance>
            <div class="space-y-4">
                <div class="space-y-2">
                    <Label for="username">Username</Label>
                    <Input id="username" name="username" bind:value={usernameInput} placeholder="admin" required readonly={showOTP} class={showOTP ? "text-muted-foreground bg-muted/40" : ""} />
                </div>
                
                <div class="space-y-2">
                    <Label for="password">Password</Label>
                    <Input id="password" name="password" type="password" required />
                </div>

                {#if showOTP}
                    <div class="space-y-2 animate-in fade-in slide-in-from-top-2">
                        <Label for="otp-token" class="flex items-center gap-2">
                            <KeyRound class="w-3.5 h-3.5" />
                            Authenticator Code
                        </Label>
                        <Input id="otp-token" name="otp-token" placeholder="000000" class="text-center tracking-[0.5em] font-mono" maxlength={6} autofocus />
                    </div>
                {/if}
            </div>

            {#if form?.invalid}
                <p class="text-sm text-destructive text-center">Invalid username or password.</p>
            {/if}
            {#if form?.otpInvalid}
                <p class="text-sm text-destructive text-center">Invalid authenticator code.</p>
            {/if}

            <Button type="submit" class="w-full">
                {showOTP ? 'Verify & Sign In' : 'Sign In'}
            </Button>
        </form>
    </div>
</div>
