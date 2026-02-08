import { http } from '$lib/api';
import { type User, type AuthResponse } from '$lib/types';
import { goto } from '$app/navigation';

class AuthStore {
	user = $state<User | null>(null);
	loading = $state(true); // Start loading to block rendering until status check
	initialized = $state<boolean | null>(null);
	isAuthenticated = $derived(!!this.user);

	async checkInitialized(fetch: typeof window.fetch) {
		try {
			const response = await http.get<{ initialized: boolean }>(fetch, '/auth/initialized');
			this.initialized = response.initialized;
		} catch (error) {
			console.error('Failed to check initialization status:', error);
			this.initialized = false;
		}
	}

	async init(fetch: typeof window.fetch) {
		try {
			// Check initialization status first or in parallel
			const initPromise = this.checkInitialized(fetch);
			
			const response = await http.get<{ status: string; user: string; has_otp: boolean }>(fetch, '/auth/status');
			// Convert backend response to internal User type
			this.user = { username: response.user || '', has_otp: response.has_otp }; 
			
			await initPromise;
		} catch (error) {
			this.user = null;
			// Ensure initialized check completes even if auth fails
			if (this.initialized === null) await this.checkInitialized(fetch);
		} finally {
			this.loading = false;
		}
	}

	async login(fetch: typeof window.fetch, username: string, password: string, otp?: string) {
		const result = await http.post<AuthResponse & { has_otp: boolean }>(fetch, '/auth/login', { username, password, otp });
		this.user = { username, has_otp: result.has_otp }; // Optimistic update from result
		return result;
	}

	async logout(fetch: typeof window.fetch) {
		try {
			await http.post(fetch, '/auth/logout', {});
		} finally {
			this.user = null;
			goto('/login');
		}
	}
}

export const authStore = new AuthStore();
