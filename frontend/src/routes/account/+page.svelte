<script lang="ts">
    import * as Card from "$lib/components/ui/card";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Button } from "$lib/components/ui/button";
    import { User, Lock, ShieldCheck, LogOut, ChevronRight, Loader2 } from '@lucide/svelte';
    import { toast } from "svelte-sonner";
    import { goto } from '$app/navigation';

    import Header from "$lib/components/organisms/Header.svelte";

    // Mock initial data
    let user = $state({
        username: 'admin',
        hasOTP: true
    });
    
    let isUsernameOpen = $state(false);
    let isPasswordOpen = $state(false);
    let isOtpOpen = $state(false);
    let loading = $state(false);
    let newOtpQr = $state<string | null>(null);

    async function simulateAction(type: string, data: any) {
        loading = true;
        // Simulate network delay
        await new Promise(r => setTimeout(r, 800));
        loading = false;

        if (type === 'logout') {
            goto('/login');
            return;
        }

        if (type === 'username') {
            user.username = data.get('newUsername');
            isUsernameOpen = false;
            toast.success('Username updated (simulated)');
        } else if (type === 'password') {
            isPasswordOpen = false;
            toast.success('Password updated (simulated)');
        } else if (type === 'toggleOtp') {
            user.hasOTP = !user.hasOTP;
            isOtpOpen = false;
            toast.success(`2FA ${user.hasOTP ? 'enabled' : 'disabled'} (simulated)`);
        } else if (type === 'regenerateOtp') {
            newOtpQr = 'https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=otpauth://totp/Wolite:admin?secret=JBSWY3DPEHPK3PXP&issuer=Wolite';
            toast.success('New OTP secret generated (simulated)');
        }
    }
</script>

<div class="container max-w-lg mx-auto py-20 px-6">
    <Header title="Account" subtitle="Manage your credentials and security." backHref="/">
         <form onsubmit={(e) => { e.preventDefault(); simulateAction('logout', null); }}>
             <Button variant="ghost" size="icon" type="submit" aria-label="Logout">
                <LogOut class="h-5 w-5 opacity-70 hover:opacity-100 transition-opacity" />
            </Button>
        </form>
    </Header>

    <div class="grid gap-4">
        <!-- Profile Card -->
        <button class="text-left w-full group" onclick={() => isUsernameOpen = true}>
            <Card.Root class="hover:bg-muted/40 transition-colors cursor-pointer border-muted/60 shadow-sm hover:shadow-md">
                <Card.Content class="p-4 flex items-center gap-4">
                    <div class="h-10 w-10 rounded-full bg-muted flex items-center justify-center shrink-0">
                        <User class="h-5 w-5 opacity-60" />
                    </div>
                    <div class="flex-1 min-w-0">
                        <p class="text-sm font-medium leading-none">Username</p>
                        <p class="text-xs text-muted-foreground mt-1 truncate">{user.username}</p>
                    </div>
                    <ChevronRight class="h-4 w-4 opacity-30 group-hover:opacity-100 transition-opacity" />
                </Card.Content>
            </Card.Root>
        </button>

        <!-- Password Card -->
        <button class="text-left w-full group" onclick={() => isPasswordOpen = true}>
             <Card.Root class="hover:bg-muted/40 transition-colors cursor-pointer border-muted/60 shadow-sm hover:shadow-md">
                <Card.Content class="p-4 flex items-center gap-4">
                    <div class="h-10 w-10 rounded-full bg-muted flex items-center justify-center shrink-0">
                        <Lock class="h-5 w-5 opacity-60" />
                    </div>
                    <div class="flex-1 min-w-0">
                        <p class="text-sm font-medium leading-none">Password</p>
                        <p class="text-xs text-muted-foreground mt-1">••••••••</p>
                    </div>
                    <ChevronRight class="h-4 w-4 opacity-30 group-hover:opacity-100 transition-opacity" />
                </Card.Content>
            </Card.Root>
        </button>

        <!-- 2FA Card -->
        <button class="text-left w-full group" onclick={() => isOtpOpen = true}>
             <Card.Root class="hover:bg-muted/40 transition-colors cursor-pointer border-muted/60 shadow-sm hover:shadow-md">
                <Card.Content class="p-4 flex items-center gap-4">
                    <div class="h-10 w-10 rounded-full bg-muted flex items-center justify-center shrink-0">
                        <ShieldCheck class="h-5 w-5 opacity-60" />
                    </div>
                    <div class="flex-1 min-w-0">
                        <p class="text-sm font-medium leading-none">Two-Factor Authentication</p>
                         <p class="text-xs text-muted-foreground mt-1 flex items-center gap-2">
                             {#if user.hasOTP}
                                <span class="h-1.5 w-1.5 rounded-full bg-green-500 inline-block"></span> Enabled
                             {:else}
                                <span class="h-1.5 w-1.5 rounded-full bg-yellow-500 inline-block"></span> Disabled
                             {/if}
                         </p>
                    </div>
                    <ChevronRight class="h-4 w-4 opacity-30 group-hover:opacity-100 transition-opacity" />
                </Card.Content>
            </Card.Root>
        </button>
    </div>
</div>

<!-- Username Dialog -->
<Dialog.Root bind:open={isUsernameOpen}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Change Username</Dialog.Title>
            <Dialog.Description>
                Enter your new username. You will be logged out after this change.
            </Dialog.Description>
        </Dialog.Header>
         <form onsubmit={(e) => { e.preventDefault(); simulateAction('username', new FormData(e.currentTarget)); }} class="grid gap-4 py-4">
            <div class="grid gap-2">
                <Label for="newUsername">New Username</Label>
                <Input id="newUsername" name="newUsername" value={user.username} required />
            </div>
            <div class="grid gap-2">
                <Label for="currentPasswordUsername">Current Password</Label>
                <Input id="currentPasswordUsername" name="currentPassword" type="password" required />
            </div>
            
            <Dialog.Footer>
                <Button type="submit" disabled={loading}>
                    {#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin"/>{/if}
                    Update Username
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>

<!-- Password Dialog -->
<Dialog.Root bind:open={isPasswordOpen}>
    <Dialog.Content class="sm:max-w-[425px]">
         <Dialog.Header>
            <Dialog.Title>Change Password</Dialog.Title>
            <Dialog.Description>
                Ensure your account is using a strong password.
            </Dialog.Description>
        </Dialog.Header>
         <form onsubmit={(e) => { e.preventDefault(); simulateAction('password', null); }} class="grid gap-4 py-4">
            <div class="grid gap-2">
                <Label for="newPassword">New Password</Label>
                <Input id="newPassword" name="newPassword" type="password" required />
            </div>
            <div class="grid gap-2">
                <Label for="currentPasswordPassword">Current Password</Label>
                <Input id="currentPasswordPassword" name="currentPassword" type="password" required />
            </div>
            
            <Dialog.Footer>
                <Button type="submit" disabled={loading}>
                    {#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin"/>{/if}
                    Update Password
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>

<!-- OTP Dialog -->
<Dialog.Root bind:open={isOtpOpen} onOpenChange={(open) => { if (!open) newOtpQr = null; }}>
    <Dialog.Content class="sm:max-w-[425px]">
        <Dialog.Header>
            <Dialog.Title>Two-Factor Authentication</Dialog.Title>
            <Dialog.Description>
                Protect your account with an extra layer of security.
            </Dialog.Description>
        </Dialog.Header>

        {#if newOtpQr}
             <div class="space-y-4 py-4 text-center">
                <div class="bg-white p-2 inline-block rounded-lg shadow-sm">
                    <img src={newOtpQr} alt="QR Code" class="h-40 w-40" />
                </div>
                <div class="text-sm text-center">
                    <p class="font-medium">Scan this QR code</p>
                    <p class="text-xs text-muted-foreground mt-1">Use your authenticator app</p>
                    <p class="font-mono text-xs bg-muted mt-2 p-2 rounded select-all">JBSWY3DPEHPK3PXP</p>
                </div>
                <Button variant="outline" class="w-full" onclick={() => isOtpOpen = false}>Done</Button>
            </div>
        {:else}
            <div class="grid gap-6 py-4">
                {#if user.hasOTP}
                     <div class="space-y-4">
                        <form onsubmit={(e) => { e.preventDefault(); simulateAction('toggleOtp', null); }} class="grid gap-4 border p-4 rounded-lg bg-destructive/5 border-destructive/10">
                            <div class="grid gap-2">
                                <Label for="currentPasswordDisable">Disable 2FA</Label>
                                <Input id="currentPasswordDisable" name="currentPassword" type="password" placeholder="Verify password to disable" required class="bg-white/50" />
                            </div>
                            <Button variant="destructive" size="sm" type="submit" disabled={loading}>
                                {#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin"/>{/if}
                                Disable 2FA
                            </Button>
                        </form>
                    </div>
                {/if}

                <div class="space-y-4">
                     <form onsubmit={(e) => { e.preventDefault(); simulateAction('regenerateOtp', null); }} class="grid gap-4">
                         <div class="grid gap-2">
                             <Label for="currentPasswordRegen">{user.hasOTP ? 'Regenerate 2FA Secret' : 'Enable 2FA'}</Label>
                             <Input id="currentPasswordRegen" name="currentPassword" type="password" placeholder="Verify password to continue" required />
                        </div>
                         <Button type="submit" disabled={loading} variant={user.hasOTP ? "secondary" : "default"}>
                             {#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin"/>{/if}
                             {user.hasOTP ? 'Regenerate Secret' : 'Set up 2FA'}
                        </Button>
                    </form>
                     {#if user.hasOTP}
                        <p class="text-xs text-muted-foreground text-center">Regenerating creates a new secret key. You will need to scan a new QR code.</p>
                    {/if}
                </div>
            </div>
        {/if}
    </Dialog.Content>
</Dialog.Root>

<style>
</style>
