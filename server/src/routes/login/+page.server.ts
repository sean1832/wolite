import { auth } from '$lib/server/auth';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
    // If no users, redirect to setup
    if (!auth.hasUsers()) {
        throw redirect(303, '/setup');
    }
    return {};
};

export const actions: Actions = {
    default: async ({ request, cookies }) => {
        const data = await request.formData();
        const username = data.get('username') as string;
        const password = data.get('password') as string;
        const otpToken = data.get('otp-token') as string;

        if (!username || !password) {
            return fail(400, { missing: true });
        }

        const user = await auth.verifyUser(username, password);

        if (!user) {
            return fail(400, { invalid: true });
        }

        if (user.hasOTP) {
            if (!otpToken) {
                // Return flag to show OTP input, preserve username
                return fail(400, { otpRequired: true, username });
            }
            
            const secret = auth.getUserOTPSecret(username);
            if (!secret || !auth.verifyOTP(otpToken, secret)) {
                return fail(400, { otpInvalid: true, username, otpRequired: true });
            }
        }

        // Login successful
        const token = auth.createSessionToken(username);
        cookies.set('session', token, {
            path: '/',
            httpOnly: true,
            sameSite: 'strict',
            secure: process.env.NODE_ENV === 'production',
            maxAge: 60 * 60 * 24 * 7 // 1 week
        });

        throw redirect(303, '/');
    }
};
