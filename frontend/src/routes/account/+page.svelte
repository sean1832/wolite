<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import { User, Lock, ShieldCheck, LogOut, ChevronRight, Loader2 } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte';
	import { http } from '$lib/api';
	import {
		InputOTP,
		InputOTPGroup,
		InputOTPSeparator,
		InputOTPSlot
	} from "$lib/components/ui/input-otp";

	import Header from '$lib/components/organisms/Header.svelte';

	let isPasswordOpen = $state(false);
	let isOtpOpen = $state(false);
	let loading = $state(false);
	let newOtpQr = $state<string | null>(null);
	let newOtpSecret = $state<string | null>(null);
	let otpCode = $state('');

	// Derived user state from authStore
	let user = $derived(authStore.user || { username: 'Guest', has_otp: false });

	async function handleLogout() {
		loading = true;
		await authStore.logout(fetch);
		loading = false;
	}

	async function verifyOtp() {
		if (otpCode.length !== 6) {
			toast.error('Please enter a valid 6-digit code');
			return;
		}
		
		loading = true;
		try {
			await http.post(fetch, '/users/otp/verify', { code: otpCode });
			toast.success('2FA Verified & Enabled!');
			
			// Refresh user state
			await authStore.init(fetch);
			
			// Close dialog
			isOtpOpen = false;
			newOtpQr = null;
			newOtpSecret = null;
			otpCode = '';
		} catch (err: any) {
			toast.error(err.message || 'Verification failed');
		} finally {
			loading = false;
		}
	}

	async function handleUpdate(
		type: 'username' | 'password' | 'toggleOtp' | 'regenerateOtp',
		formData?: FormData
	) {
		loading = true;
		try {
			const payload: any = {};

			if (type === 'username') {
				payload.username = formData?.get('newUsername');
				// Note: Changing username might require re-login or updating the store
			} else if (type === 'password') {
				payload.password = formData?.get('newPassword');
			} else if (type === 'toggleOtp') {
				// If currently enabled, we are disabling it (use_otp: false)
				// If currently disabled, we are enabling it (use_otp: true)
				payload.use_otp = !user.has_otp;
			} else if (type === 'regenerateOtp') {
				payload.use_otp = true; // Force enable to get new secret
			}

			// Call API
			// Note: Our current backend handleUserUpdate expects 'username', 'password', 'use_otp'
			// AND 'username' in payload must match current user (allow change? backend says forbidden if different)
			// Wait, backend handleUserUpdate checks: if payload.Username != claims.Username -> Forbidden.
			// So we CANNOT change username with current backend logic.
			// I will skip username change for now or we must update backend to allow it (if new username is free).
			// Let's check backend implementation again...
			// "if payload.Username != claims.Username { Forbidden }" -> converting Username is NOT supported.

			// For now, I will disable username editing or just not send it.
			// Actually, the PUT /users payload structure requires us to send what we WANT TO CHANGE?
			// "payload" struct has Username, Password, UseOTP.
			// If I send ONLY Password, Username is empty string -> != claims.Username -> Forbidden.
			// This is a BACKEND BUG/LIMITATION.
			// Payload struct: Username string `json:"username"`
			// Logic:
			// 1. Decode payload.
			// 2. Check payload.Username != claims.Username.
			// This means we MUST send the CURRENT username in the payload for it to work.

			// So we always send:
			payload.username = user.username;

			// Now apply changes:
			// Now apply changes:
			if (type === 'password') {
				// For password change
				payload.password = formData?.get('newPassword');
				payload.old_password = formData?.get('currentPassword');
			}
			// For OTP toggles, we set use_otp

			const response = await http.put<any>(fetch, '/users', payload);

			if (response && response.otp_url) {
				// OTP was enabled or regenerated
				const secretMatch = response.otp_url.match(/secret=([A-Za-z0-9]+)/);
				const secret = secretMatch ? secretMatch[1] : 'Unknown Secret';
				
				newOtpSecret = secret;
				newOtpQr = `https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${encodeURIComponent(response.otp_url)}`;
			} else {
				newOtpQr = null;
				newOtpSecret = null;
				isOtpOpen = false;
				isPasswordOpen = false;
			}

			// Refresh user profile if needed (api doesn't return full user, but we know what changed)
			// Ideally we re-fetch profile. But /auth/status just returns username.
			// We assume it worked.

			// Update local store state for OTP
			if (type === 'toggleOtp') {
				if (authStore.user) authStore.user.has_otp = payload.use_otp;
				toast.success(`2FA ${payload.use_otp ? 'enabled' : 'disabled'}`);
			} else if (type === 'regenerateOtp') {
				toast.success('New 2FA secret generated. Scan the code.');
			} else if (type === 'password') {
				toast.success('Password updated');
			}
		} catch (err: any) {
			toast.error(err.message || 'Update failed');
		} finally {
			loading = false;
		}
	}
</script>

<div class="container mx-auto max-w-lg px-6 py-20">
	<Header title="Account" subtitle="Manage your credentials and security." backHref="/" showAccount={false}>
		<form
			onsubmit={(e) => {
				e.preventDefault();
				handleLogout();
			}}
		>
			<Button variant="ghost" size="icon" type="submit" aria-label="Logout" disabled={loading}>
				{#if loading}
					<Loader2 class="h-4 w-4 animate-spin opacity-70" />
				{:else}
					<LogOut class="h-5 w-5 opacity-70 transition-opacity hover:opacity-100" />
				{/if}
			</Button>
		</form>
	</Header>

	<div class="grid gap-4">

		<!-- Password Card -->
		<button class="group w-full text-left" onclick={() => (isPasswordOpen = true)}>
			<Card.Root
				class="cursor-pointer border-muted/60 shadow-sm transition-colors hover:bg-muted/40 hover:shadow-md"
			>
				<Card.Content class="flex items-center gap-4 p-4">
					<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-muted">
						<Lock class="h-5 w-5 opacity-60" />
					</div>
					<div class="min-w-0 flex-1">
						<p class="text-sm leading-none font-medium">Password</p>
						<p class="mt-1 text-xs text-muted-foreground">••••••••</p>
					</div>
					<ChevronRight class="h-4 w-4 opacity-30 transition-opacity group-hover:opacity-100" />
				</Card.Content>
			</Card.Root>
		</button>

		<!-- 2FA Card -->
		<button class="group w-full text-left" onclick={() => (isOtpOpen = true)}>
			<Card.Root
				class="cursor-pointer border-muted/60 shadow-sm transition-colors hover:bg-muted/40 hover:shadow-md"
			>
				<Card.Content class="flex items-center gap-4 p-4">
					<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-muted">
						<ShieldCheck class="h-5 w-5 opacity-60" />
					</div>
					<div class="min-w-0 flex-1">
						<p class="text-sm leading-none font-medium">Two-Factor Authentication</p>
						<p class="mt-1 flex items-center gap-2 text-xs text-muted-foreground">
							{#if user.has_otp}
								<span class="inline-block h-1.5 w-1.5 rounded-full bg-green-500"></span> Enabled
							{:else}
								<span class="inline-block h-1.5 w-1.5 rounded-full bg-yellow-500"></span> Disabled
							{/if}
						</p>
					</div>
					<ChevronRight class="h-4 w-4 opacity-30 transition-opacity group-hover:opacity-100" />
				</Card.Content>
			</Card.Root>
		</button>
	</div>
</div>


<!-- Password Dialog -->
<Dialog.Root bind:open={isPasswordOpen}>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Change Password</Dialog.Title>
			<Dialog.Description>Ensure your account is using a strong password.</Dialog.Description>
		</Dialog.Header>
		<form
			onsubmit={(e) => {
				e.preventDefault();
				const formData = new FormData(e.currentTarget);
				if (formData.get('newPassword') !== formData.get('confirmPassword')) {
					toast.error('New passwords do not match');
					return;
				}
				handleUpdate('password', formData);
			}}
			class="grid gap-4 py-4"
		>
			<div class="grid gap-2">
				<Label for="currentPassword">Current Password</Label>
				<Input id="currentPassword" name="currentPassword" type="password" required />
			</div>
			<div class="grid gap-2">
				<Label for="newPassword">New Password</Label>
				<Input id="newPassword" name="newPassword" type="password" required />
			</div>
			<div class="grid gap-2">
				<Label for="confirmPassword">Confirm New Password</Label>
				<Input id="confirmPassword" name="confirmPassword" type="password" required />
			</div>

			<Dialog.Footer>
				<Button type="submit" disabled={loading}>
					{#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin" />{/if}
					Update Password
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

<!-- OTP Dialog -->
<Dialog.Root
	bind:open={isOtpOpen}
	onOpenChange={(open) => {
		if (!open) newOtpQr = null;
	}}
>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Two-Factor Authentication</Dialog.Title>
			<Dialog.Description>Protect your account with an extra layer of security.</Dialog.Description>
		</Dialog.Header>

		{#if newOtpQr}
			<div class="space-y-4 py-4 text-center">
				<div class="inline-block rounded-lg bg-white p-2 shadow-sm">
					<img src={newOtpQr} alt="QR Code" class="h-40 w-40" />
				</div>
				<div class="text-center text-sm">
					<p class="font-medium">Scan this QR code</p>
					<p class="mt-1 text-xs text-muted-foreground">Use your authenticator app</p>
				</div>
				
				{#if newOtpSecret}
					<div class="rounded bg-muted/30 p-2 font-mono text-xs select-all">
						{newOtpSecret}
					</div>
				{/if}

				<div class="space-y-2 text-center">
					<Label for="account-otp-verify">Enter Verification Code</Label>
					<div class="flex flex-col items-center gap-4">
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
						
						<Button onclick={verifyOtp} disabled={loading || otpCode.length !== 6} class="w-full">
							{#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin" />{/if}
							Verify & Enable
						</Button>
					</div>
				</div>

				<Button variant="ghost" class="w-full text-muted-foreground" onclick={() => (isOtpOpen = false)}>Cancel</Button>
			</div>
		{:else}
			{#if user.has_otp}
				<div class="space-y-4">
					<form
						onsubmit={(e) => {
							e.preventDefault();
							handleUpdate('toggleOtp');
						}}
						class="grid gap-4 rounded-lg border border-destructive/10 bg-destructive/5 p-4"
					>
						<div class="grid gap-2">
							<Label>Disable 2FA</Label>
							<p class="text-xs text-muted-foreground">
								This will remove the current 2FA requirement.
							</p>
						</div>
						<Button variant="destructive" size="sm" type="submit" disabled={loading}>
							{#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin" />{/if}
							Disable 2FA
						</Button>
					</form>
				</div>
			{/if}

			<div class="space-y-4">
				<form
					onsubmit={(e) => {
						e.preventDefault();
						handleUpdate('regenerateOtp');
					}}
					class="grid gap-4"
				>
					<div class="grid gap-2">
						<Label>{user.has_otp ? 'Regenerate 2FA Secret' : 'Enable 2FA'}</Label>
						<p class="text-xs text-muted-foreground">
							{user.has_otp
								? 'Regenerating creates a new secret key. You will need to scan a new QR code.'
								: 'Enable 2FA to add an extra layer of security.'}
						</p>
					</div>
					<Button type="submit" disabled={loading} variant={user.has_otp ? 'secondary' : 'default'}>
						{#if loading}<Loader2 class="mr-2 h-4 w-4 animate-spin" />{/if}
						{user.has_otp ? 'Regenerate Secret' : 'Set up 2FA'}
					</Button>
				</form>
			</div>
		{/if}
	</Dialog.Content>
</Dialog.Root>
