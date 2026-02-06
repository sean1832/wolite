import { redirect, type Handle } from '@sveltejs/kit';
import { auth } from '$lib/server/auth';
import { db } from '$lib/server/db';

export const handle: Handle = async ({ event, resolve }) => {
    const sessionToken = event.cookies.get('session');
    let user = null;

    if (sessionToken) {
        const decoded = auth.verifySessionToken(sessionToken);
        if (decoded) {
            const dbUser = db.findUser(decoded.username);
            if (dbUser) {
                user = { username: dbUser.username, hasOTP: !!dbUser.otpSecret };
            }
        }
    }

    if (user) {
        event.locals.user = user;
    }

    const path = event.url.pathname;

    // 1. If no users exist, enforce setup flow
    // But allow assets and API calls needed for setup (if any)
    if (!auth.hasUsers()) {
        if (path !== '/setup' && !path.startsWith('/api')) {
            throw redirect(303, '/setup');
        }
    } else {
        // Users exist. If trying to access /setup, redirect to / (unless we want to allow re-setup? Security risk if unauthenticated).
        // Let's block /setup if users exist for now, or require auth. 
        // For simplicity: if users exist, /setup is forbidden or redirects home.
        if (path === '/setup') {
            throw redirect(303, '/');
        }

        // 2. Protect protected routes
        const isPublic = path === '/login';
        if (!user && !isPublic) {
            throw redirect(303, '/login');
        }

        // 3. Redirect logged-in users away from /login
        if (user && path === '/login') {
            throw redirect(303, '/');
        }
    }

    const response = await resolve(event);
    return response;
};
