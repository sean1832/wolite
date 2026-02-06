import { auth } from '$lib/server/auth';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
    if (!locals.user) {
        throw redirect(303, '/login');
    }

    return {
        user: {
            username: locals.user.username,
            hasOTP: locals.user.hasOTP
        }
    };
};

export const actions: Actions = {
    changeUsername: async ({ request, locals }) => {
        const data = await request.formData();
        const currentPassword = data.get('currentPassword') as string;
        const newUsername = data.get('newUsername') as string;

        if (!locals.user) return fail(401);
        if (!currentPassword || !newUsername) return fail(400, { missing: true });

        // Verify password first
        const verified = await auth.verifyUser(locals.user.username, currentPassword);
        if (!verified) {
             return fail(400, { invalidPassword: true });
        }

        try {
            auth.updateUsername(locals.user.username, newUsername);
        } catch (e) {
            return fail(400, { usernameTaken: true });
        }

        // Logout to force re-login with new username (simpler session management)
        // Or update session. Let's redirect to login for simplicity and security.
        throw redirect(303, '/login?message=Username changed. Please login again.');
    },

    changePassword: async ({ request, locals }) => {
        const data = await request.formData();
        const currentPassword = data.get('currentPassword') as string;
        const newPassword = data.get('newPassword') as string;

        if (!locals.user) return fail(401);
        if (!currentPassword || !newPassword) return fail(400, { missing: true });

        const verified = await auth.verifyUser(locals.user.username, currentPassword);
        if (!verified) {
             return fail(400, { invalidPassword: true });
        }

        await auth.updatePassword(locals.user.username, newPassword);
        
        return { success: true, message: 'Password updated successfully' };
    },

    toggleOtp: async ({ request, locals }) => {
        const data = await request.formData();
        const currentPassword = data.get('currentPassword') as string;
        const action = data.get('action') as string; // 'disable'

        if (!locals.user) return fail(401);
        if (!currentPassword) return fail(400, { missing: true });

        const verified = await auth.verifyUser(locals.user.username, currentPassword);
        if (!verified) {
             return fail(400, { invalidPassword: true });
        }

        if (action === 'disable') {
            auth.disableOTP(locals.user.username);
        }

        return { success: true };
    },

    regenerateOtp: async ({ request, locals }) => {
        const data = await request.formData();
        const currentPassword = data.get('currentPassword') as string;

        if (!locals.user) return fail(401);
        if (!currentPassword) return fail(400, { missing: true });

        const verified = await auth.verifyUser(locals.user.username, currentPassword);
        if (!verified) {
             return fail(400, { invalidPassword: true });
        }

        const { imageUrl, secret } = await auth.regenerateOTP(locals.user.username);
        
        return { success: true, newOtpQr: imageUrl, newSecret: secret };
    },

    logout: async ({ cookies }) => {
        cookies.delete('session', { path: '/' });
        throw redirect(303, '/login');
    }
};
