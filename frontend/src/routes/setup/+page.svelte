<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Checkbox } from "$lib/components/ui/checkbox";
    import { ShieldCheck, Loader2 } from '@lucide/svelte';
    import { goto } from '$app/navigation';
    import { toast } from 'svelte-sonner';

    let setupOTP = $state(false);
    let username = $state('');
    let password = $state('');
    let confirmPassword = $state('');
    let otpData = $state<{ secret: string, imageUrl: string } | null>(null);
    let loading = $state(false);
    let error = $state('');

    async function handleSetup(e: Event) {
        e.preventDefault();
        error = '';
        loading = true;

        // Simulate network delay
        await new Promise(r => setTimeout(r, 1000));
        loading = false;

        if (password !== confirmPassword) {
            error = 'Passwords do not match';
            return;
        }

        toast.success('Admin account created (simulated)');
        goto('/login');
    }

    async function generateOTP() {
        loading = true;
        await new Promise(r => setTimeout(r, 500));
        otpData = { 
            secret: 'JBSWY3DPEHPK3PXP', 
            imageUrl: 'https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=otpauth://totp/Wolite:admin?secret=JBSWY3DPEHPK3PXP&issuer=Wolite' 
        };
        loading = false;
    }
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
        <form class="space-y-6" onsubmit={handleSetup}>
            <div class="space-y-2">
                <Label for="username">Username</Label>
                <Input id="username" name="username" bind:value={username} placeholder="admin" required />
            </div>
            
            <div class="space-y-2">
                <Label for="password">Password</Label>
                <Input id="password" name="password" type="password" bind:value={password} required />
            </div>

            <div class="space-y-2">
                <Label for="confirm-password">Confirm Password</Label>
                <Input id="confirm-password" name="confirm-password" type="password" bind:value={confirmPassword} required />
            </div>

            <div class="flex items-center space-x-2 pt-2">
                <Checkbox id="setup-otp" bind:checked={setupOTP} />
                <Label for="setup-otp" class="font-normal text-muted-foreground">Enable Two-Factor Authentication (Recommended)</Label>
            </div>

            {#if setupOTP}
                <div class="p-4 border border-border/50 rounded-lg bg-muted/30 space-y-4">
                    {#if !otpData}
                        <div class="text-center py-4">
                             <p class="text-sm text-muted-foreground mb-4">Click below to generate your unique QR code.</p>
                             <Button type="button" onclick={generateOTP} variant="secondary" size="sm" class="w-full" disabled={loading}>
                                {#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin"/>{/if}
                                Generate QR Code
                             </Button>
                        </div>
                    {:else}
                        <div class="flex flex-col items-center space-y-4">
                             <img src={otpData.imageUrl} alt="QR Code" class="w-40 h-40 bg-white rounded-lg p-2" />
                             <p class="text-xs text-muted-foreground text-center">Scan with your authenticator app<br/>or enter code manually:</p>
                             <code class="text-xs font-mono bg-muted px-2 py-1 rounded">{otpData.secret}</code>
                             
                             <div class="w-full space-y-2">
                                <Label for="otp-token">Verification Code</Label>
                                <Input id="otp-token" name="otp-token" placeholder="000000" class="text-center tracking-[0.5em] font-mono" maxlength={6} />
                             </div>
                        </div>
                    {/if}
                </div>
            {/if}
            
            {#if error}
                <p class="text-sm text-destructive text-center">{error}</p>
            {/if}

            <Button type="submit" class="w-full" disabled={loading}>
                {#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin"/>{/if}
                Create Account
            </Button>
        </form>
    </div>
</div>
