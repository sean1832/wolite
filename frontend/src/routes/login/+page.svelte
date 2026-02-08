<script lang="ts">
<<<<<<< HEAD
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { LogIn, KeyRound } from '@lucide/svelte';
    import { goto } from '$app/navigation';

    // Mock state for UI demonstration
    let usernameInput = $state('admin');
    let passwordInput = $state('');
    let otpInput = $state('');
    let showOTP = $state(false);
    let error = $state('');

    async function handleLogin(e: Event) {
        e.preventDefault();
        error = '';
        
        // Simulating backend logic for UI testing
        if (!showOTP) {
            if (usernameInput === 'admin' && passwordInput === 'password') {
                showOTP = true;
            } else {
                error = 'Invalid username or password (try admin/password)';
            }
        } else {
            if (otpInput === '123456') {
                goto('/');
            } else {
                error = 'Invalid authenticator code (try 123456)';
            }
        }
    }
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

        <form class="space-y-6" onsubmit={handleLogin}>
            <div class="space-y-4">
                <div class="space-y-2">
                    <Label for="username">Username</Label>
                    <Input id="username" name="username" bind:value={usernameInput} placeholder="admin" required readonly={showOTP} class={showOTP ? "text-muted-foreground bg-muted/40" : ""} />
                </div>
                
                <div class="space-y-2">
                    <Label for="password">Password</Label>
                    <Input id="password" name="password" type="password" bind:value={passwordInput} required readonly={showOTP} class={showOTP ? "text-muted-foreground bg-muted/40" : ""} />
                </div>

                {#if showOTP}
                    <div class="space-y-2 animate-in fade-in slide-in-from-top-2">
                        <Label for="otp-token" class="flex items-center gap-2">
                            <KeyRound class="w-3.5 h-3.5" />
                            Authenticator Code
                        </Label>
                        <Input id="otp-token" name="otp-token" bind:value={otpInput} placeholder="000000" class="text-center tracking-[0.5em] font-mono" maxlength={6} autofocus />
                    </div>
                {/if}
            </div>

            {#if error}
                <p class="text-sm text-destructive text-center">{error}</p>
            {/if}

            <Button type="submit" class="w-full">
                {showOTP ? 'Verify & Sign In' : 'Sign In'}
            </Button>
        </form>
    </div>
=======
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { LogIn, KeyRound, Loader2 } from '@lucide/svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte';
	import { toast } from 'svelte-sonner';
	import {
		InputOTP,
		InputOTPGroup,
		InputOTPSeparator,
		InputOTPSlot
	} from "$lib/components/ui/input-otp";

	let usernameInput = $state('');
	let passwordInput = $state('');
	let otpInput = $state('');
	let showOTP = $state(false);
	let error = $state('');
	let loading = $state(false);

	async function handleLogin(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;

		try {
			// Attempt login
			await authStore.login(fetch, usernameInput, passwordInput, otpInput);
			toast.success('Welcome back!');
			goto('/');
		} catch (err: any) {
			// Handle specific errors
			const msg = err.message || 'Login failed';

			if (msg.includes('OTP required')) {
				showOTP = true;
				error = ''; // Clear error if it was just asking for OTP
			} else if (msg.includes('Invalid OTP')) {
				error = 'Invalid authenticator code';
			} else if (msg.includes('Invalid credentials')) {
				error = 'Invalid username or password';
			} else {
				error = msg;
			}
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex h-screen w-full items-center justify-center bg-zinc-50 dark:bg-zinc-950">
	<div
		class="w-full max-w-sm space-y-8 rounded-xl border border-border/40 bg-background p-8 shadow-sm"
	>
		<div class="flex flex-col items-center space-y-2 text-center">
			<div class="mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
				<LogIn class="h-6 w-6 text-primary" />
			</div>
			<h1 class="text-2xl font-semibold tracking-tight">Sign in</h1>
			<p class="text-sm text-muted-foreground">Enter your credentials to access the dashboard.</p>
		</div>

		<form class="space-y-6" onsubmit={handleLogin}>
			<div class="space-y-4">
				<div class="space-y-2">
					<Label for="username">Username</Label>
					<Input
						id="username"
						name="username"
						bind:value={usernameInput}
						placeholder="admin"
						required
						readonly={showOTP}
						class={showOTP ? 'bg-muted/40 text-muted-foreground' : ''}
					/>
				</div>

				<div class="space-y-2">
					<Label for="password">Password</Label>
					<Input
						id="password"
						name="password"
						type="password"
						bind:value={passwordInput}
						required
						readonly={showOTP}
						class={showOTP ? 'bg-muted/40 text-muted-foreground' : ''}
					/>
				</div>

				{#if showOTP}
					<div class="animate-in space-y-2 fade-in slide-in-from-top-2">
						<Label for="otp-token" class="flex items-center gap-2">
							<KeyRound class="h-3.5 w-3.5" />
							Authenticator Code
						</Label>
						<div class="flex justify-center py-2">
							<InputOTP maxlength={6} bind:value={otpInput}>
								{#snippet children({ cells })}
									<InputOTPGroup>
										{#each cells.slice(0, 3) as cell}
											<InputOTPSlot {cell} />
										{/each}
									</InputOTPGroup>
									<InputOTPSeparator />
									<InputOTPGroup>
										{#each cells.slice(3, 6) as cell}
											<InputOTPSlot {cell} />
										{/each}
									</InputOTPGroup>
								{/snippet}
							</InputOTP>
						</div>
					</div>
				{/if}
			</div>

			{#if error}
				<p class="text-center text-sm text-destructive">{error}</p>
			{/if}

			<Button type="submit" class="w-full" disabled={loading}>
				{#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin" />{/if}
				{showOTP ? 'Verify & Sign In' : 'Sign In'}
			</Button>
		</form>
	</div>
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
</div>
