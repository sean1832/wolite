<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import { Toaster } from '$lib/components/ui/sonner';
	import { authStore } from '$lib/stores/auth.svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import type { ToasterProps } from 'svelte-sonner';

	let { children } = $props();

	let toastPosition = $state<ToasterProps['position']>('bottom-left');

	// Protected routes that require authentication
	// We allow public access only to /login and /setup
	const publicRoutes = ['/login', '/setup'];

	// Effect to handle auth checks and redirects
	$effect(() => {
		// Wait for both auth check and initialization check to complete
		if (authStore.loading || authStore.initialized === null) return;

		const currentPath = $page.url.pathname;
		const isPublic = publicRoutes.some((route) => currentPath.startsWith(route));

		if (!authStore.isAuthenticated && !isPublic) {
			// Not logged in and trying to access protected route
			// Redirect to setup if no users exist, otherwise login
			if (authStore.initialized) {
				goto('/login');
			} else {
				goto('/setup');
			}
		} else if (authStore.isAuthenticated && isPublic) {
			// Logged in and trying to access public route (login/setup) -> Dashboard
			goto('/');
		} else if (!authStore.isAuthenticated && currentPath === '/login' && !authStore.initialized) {
			// On login page but no users exist -> redirect to setup
			goto('/setup');
		}
	});

	onMount(() => {
		authStore.init(fetch);

		const mediaQuery = window.matchMedia('(max-width: 640px)');
		const updatePosition = () => {
			toastPosition = mediaQuery.matches ? 'top-center' : 'bottom-right';
		};

		updatePosition();
		mediaQuery.addEventListener('change', updatePosition);

		return () => mediaQuery.removeEventListener('change', updatePosition);
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<ModeWatcher />
<Toaster position={toastPosition} />
{@render children()}
