<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { ShieldCheck, Loader2 } from '@lucide/svelte';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { http } from '$lib/api';
	import { authStore } from '$lib/stores/auth.svelte';
	import {
		InputOTP,
		InputOTPGroup,
		InputOTPSeparator,
		InputOTPSlot
	} from "$lib/components/ui/input-otp";

	let setupOTP = $state(false);
	let username = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let otpData = $state<{ secret: string; imageUrl: string } | null>(null);
	let loading = $state(false);
	let error = $state('');

	async function handleSetup(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			loading = false;
			return;
		}

		try {
			// Setup request
			const response = await http.post<any>(fetch, '/users', {
				username,
				password,
				use_otp: setupOTP
			});

			// Check if OTP url was returned
			if (setupOTP && response && response.otp_url) {
				// Generate QR code for the returned URL
				const secretMatch = response.otp_url.match(/secret=([A-Za-z0-9]+)/);
				const secret = secretMatch ? secretMatch[1] : 'Unknown Secret';
				
				otpData = {
					secret: secret,
					imageUrl: `https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${encodeURIComponent(response.otp_url)}`
				};

				// Don't redirect yet, let user scan
				toast.success('Account created! Scan the QR code to set up 2FA.');
				return;
			}

			// Auto-login successful
			await authStore.init(fetch);
			toast.success('Admin account created successfully');
			goto('/');
		} catch (err: any) {
			error = err.message || 'Setup failed';
		} finally {
			loading = false;
		}
	}

	let otpCode = $state('');
	
	async function verifyAndFinish() {
		if (otpCode.length !== 6) {
			error = 'Please enter a valid 6-digit code';
			return;
		}
		
		loading = true;
		error = '';

		try {
			await http.post(fetch, '/users/otp/verify', { code: otpCode });
			await authStore.init(fetch);
			toast.success('2FA Verified! Setup complete.');
			goto('/');
		} catch (err: any) {
			error = err.message || 'Verification failed';
			loading = false;
		}
	}
</script>

<div class="flex h-screen w-full items-center justify-center bg-zinc-50 dark:bg-zinc-950">
	<div
		class="w-full max-w-md space-y-8 rounded-xl border border-border/40 bg-background p-8 shadow-sm"
	>
		<div class="flex flex-col items-center space-y-2 text-center">
			<div class="mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
				<ShieldCheck class="h-6 w-6 text-primary" />
			</div>
			<h1 class="text-2xl font-semibold tracking-tight">Create Admin Account</h1>
			<p class="text-sm text-muted-foreground">Set up your administrator credentials.</p>
		</div>

		{#if otpData}
			<div class="animate-in space-y-6 text-center fade-in slide-in-from-bottom-4">
				<div class="inline-block rounded-xl border bg-white p-4 shadow-sm">
					<img src={otpData.imageUrl} alt="QR Code" class="h-48 w-48" />
				</div>

				<div class="space-y-2">
					<h3 class="font-medium">Scan this code</h3>
					<p class="text-sm text-muted-foreground">
						Open your authenticator app and scan the QR code above.
					</p>
				</div>

				<div class="space-y-4">
					<div class="rounded bg-muted/30 p-3 font-mono text-xs select-all">
						{otpData.secret}
					</div>
					
					<div class="flex justify-center">
						<InputOTP maxlength={6} bind:value={otpCode}>
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

					{#if error}
						<p class="text-center text-sm text-destructive">{error}</p>
					{/if}

					<Button class="w-full" onclick={verifyAndFinish} disabled={loading || otpCode.length !== 6}>
						{#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin" />{/if}
						Verify & Complete Setup
					</Button>
				</div>
			</div>
		{:else}
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
					<Input
						id="confirm-password"
						name="confirm-password"
						type="password"
						bind:value={confirmPassword}
						required
					/>
				</div>

				<div class="flex items-center space-x-2 pt-2">
					<Checkbox id="setup-otp" bind:checked={setupOTP} />
					<Label for="setup-otp" class="font-normal text-muted-foreground"
						>Enable Two-Factor Authentication</Label
					>
				</div>

				{#if error}
					<p class="text-center text-sm text-destructive">{error}</p>
				{/if}

				<Button type="submit" class="w-full" disabled={loading}>
					{#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin" />{/if}
					Create Account
				</Button>
			</form>
		{/if}
	</div>
</div>
