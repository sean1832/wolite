<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Checkbox } from "$lib/components/ui/checkbox";
    import { enhance } from '$app/forms';
    import { ShieldCheck, Loader2 } from '@lucide/svelte';

    let { form } = $props();
    
    let setupOTP = $state(false);
    let username = $state('');
    let otpData = $state<{ secret: string, imageUrl: string } | null>(null);

    $effect(() => {
        // Safe access to form properties
        const f = form as any;
        if (f?.step === 'otp-generated') {
             otpData = { secret: f.secret, imageUrl: f.imageUrl };
        }
    });
</script>

<div class="h-screen w-full flex items-center justify-center bg-zinc-50 dark:bg-zinc-950">
    <div class="w-full max-w-md p-8 space-y-8 bg-background border border-border/40 rounded-xl shadow-sm">
        <div class="flex flex-col items-center text-center space-y-2">
            <div class="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-2">
                <ShieldCheck class="w-6 h-6 text-primary" />
            </div>
            <h1 class="text-2xl font-semibold tracking-tight">Welcome to Wolite</h1>
            <p class="text-sm text-muted-foreground">Setup your admin account to get started.</p>
        </div>

        <!-- Main Setup Form -->
        <form method="POST" action="?/register" class="space-y-6" use:enhance={({ action }) => {
            return async ({ update }) => {
                if (action.search.includes('?/generateOTP')) {
                    await update({ reset: false });
                } else {
                    await update();
                }
            };
        }}>
            <div class="space-y-2">
                <Label for="username">Username</Label>
                <Input id="username" name="username" bind:value={username} placeholder="admin" required />
            </div>
            
            <div class="space-y-2">
                <Label for="password">Password</Label>
                <Input id="password" name="password" type="password" required />
            </div>

            <div class="space-y-2">
                <Label for="confirm-password">Confirm Password</Label>
                <Input id="confirm-password" name="confirm-password" type="password" required />
            </div>

            <div class="flex items-center space-x-2 pt-2">
                <Checkbox id="setup-otp" bind:checked={setupOTP} />
                <input type="hidden" name="setup-otp" value={setupOTP ? 'on' : 'off'} />
                <Label for="setup-otp" class="font-normal text-muted-foreground">Enable Two-Factor Authentication (Recommended)</Label>
            </div>

            {#if setupOTP}
                <div class="p-4 border border-border/50 rounded-lg bg-muted/30 space-y-4">
                    {#if !otpData}
                        <div class="text-center py-4">
                             <p class="text-sm text-muted-foreground mb-4">Click below to generate your unique QR code.</p>
                             <Button type="submit" formaction="?/generateOTP" variant="secondary" size="sm" class="w-full">
                                Generate QR Code
                             </Button>
                        </div>
                    {:else}
                        <div class="flex flex-col items-center space-y-4">
                             <img src={otpData.imageUrl} alt="QR Code" class="w-40 h-40 bg-white rounded-lg p-2" />
                             <p class="text-xs text-muted-foreground text-center">Scan with your authenticator app<br/>or enter code manually:</p>
                             <code class="text-xs font-mono bg-muted px-2 py-1 rounded">{otpData.secret}</code>
                             <input type="hidden" name="otp-secret" value={otpData.secret} />
                             
                             <div class="w-full space-y-2">
                                <Label for="otp-token">Verification Code</Label>
                                <Input id="otp-token" name="otp-token" placeholder="000000" class="text-center tracking-[0.5em] font-mono" maxlength={6} />
                             </div>
                        </div>
                    {/if}
                </div>
            {/if}
            
            {#if form?.error}
                <p class="text-sm text-destructive text-center">{form.error}</p>
            {/if}
            {#if form?.passwordMismatch}
                <p class="text-sm text-destructive text-center">Passwords do not match.</p>
            {/if}
            {#if form?.otpInvalid}
               <p class="text-sm text-destructive text-center">Invalid OTP code.</p>
            {/if}

            <Button type="submit" class="w-full">Create Account</Button>
        </form>
    </div>
</div>
