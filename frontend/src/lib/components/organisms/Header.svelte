<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { ArrowLeft, User, Sun, Moon, LogOut, Settings } from '@lucide/svelte';
	import { toggleMode } from 'mode-watcher';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte';

	let {
		title,
		subtitle,
		backHref = undefined,
		children = undefined,
		showAccount = true
	} = $props();
</script>

<header class="mb-8 flex items-end justify-between py-6 lg:mb-12">
	<div class="flex items-center gap-4">
		{#if backHref}
			<a href={backHref} class="-ml-2 inline-flex p-0">
				<Button
					variant="ghost"
					size="icon"
					class="h-8 w-8 text-muted-foreground/50 hover:text-foreground"
				>
					<ArrowLeft class="h-5 w-5" />
				</Button>
			</a>
		{/if}
		<div class="space-y-0.5">
			<h1 class="text-2xl font-light tracking-tight text-foreground">{title}</h1>
			{#if subtitle}
				<p class="text-xs tracking-widest text-muted-foreground/50 uppercase">{subtitle}</p>
			{/if}
		</div>
	</div>

	<div class="flex items-center gap-4">
		<div class="hidden items-center gap-2 sm:flex">
			{@render children?.()}
		</div>

		{#if showAccount}
			<DropdownMenu.Root>
				<DropdownMenu.Trigger class="focus:outline-none">
					<Avatar.Root class="h-8 w-8 transition-opacity hover:opacity-80">
						<Avatar.Fallback class="bg-muted text-xs text-muted-foreground"
							><User class="h-4 w-4" /></Avatar.Fallback
						>
					</Avatar.Root>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end" class="w-48">
					<DropdownMenu.Label>My Account</DropdownMenu.Label>
					<DropdownMenu.Separator />
					<DropdownMenu.Item
						onclick={async () => {
							await goto('/account');
						}}
					>
						<Settings class="mr-2 h-4 w-4" />
						<span>Settings</span>
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={toggleMode}>
						<Sun
							class="mr-2 h-4 w-4 scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90"
						/>
						<Moon
							class="absolute mr-2 h-4 w-4 scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0"
						/>
						<span>Toggle Theme</span>
					</DropdownMenu.Item>
					<DropdownMenu.Separator />
					<DropdownMenu.Item
						class="text-destructive focus:text-destructive"
						onclick={async () => await authStore.logout(fetch)}
					>
						<LogOut class="mr-2 h-4 w-4" />
						<span>Sign out</span>
					</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		{/if}
	</div>
</header>
