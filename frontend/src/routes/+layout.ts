import { redirect } from '@sveltejs/kit';
import { browser } from '$app/environment';
import { authStore } from '$lib/stores/auth.svelte';

export const load = async ({ url }) => {
    // Client-side route guard
    // Note: In SvelteKit + SPA mode with client-side auth, we rely on the store.
    // However, store is initialized asynchronously.
    // A more robust way is to check in +layout.svelte or a wrapper component.
    // But let's add a basic check here for initial load if possible, 
    // or rely on the layout effect.
    
    // Since we are using client-side only auth (cookie not readable by JS),
    // we must wait for authStore.init() to complete.
    // We'll handle redirection in +layout.svelte effect.
    return {
        url: url.pathname
    };
};
