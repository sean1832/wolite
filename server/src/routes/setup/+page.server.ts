import { auth } from '$lib/server/auth';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
    if (auth.hasUsers()) {
        throw redirect(303, '/login');
    }
    return {};
};

export const actions: Actions = {
    register: async ({ request, cookies }) => {
        const data = await request.formData();
        const username = data.get('username') as string;
        const password = data.get('password') as string;
        const confirmPassword = data.get('confirm-password') as string;
        const setupOTP = data.get('setup-otp') === 'on';
        const otpToken = data.get('otp-token') as string;
        const otpSecret = data.get('otp-secret') as string;

        console.log('Setup Action:', { username, setupOTP, otpToken, otpSecretPayload: !!otpSecret });

        if (!username || !password) {
            console.log('Setup Failed: Missing fields');
            return fail(400, { missing: true });
        }

        if (password !== confirmPassword) {
            console.log('Setup Failed: Password mismatch');
            return fail(400, { passwordMismatch: true });
        }

        let userSecret = undefined;

        if (setupOTP) {
            if (!otpToken || !otpSecret) {
                console.log('Setup Failed: OTP missing');
                return fail(400, { otpMissing: true });
            }
            const isValid = auth.verifyOTP(otpToken, otpSecret);
            if (!isValid) {
                console.log('Setup Failed: OTP invalid');
                return fail(400, { otpInvalid: true });
            }
            userSecret = otpSecret;
        }

        try {
            await auth.createUser(username, password, userSecret);
            console.log('Setup Success: User created');
        } catch (e) {
            console.error('Setup Error:', e);
            return fail(500, { error: 'Failed to create user' });
        }
        
        const token = auth.createSessionToken(username);
        cookies.set('session', token, {
            path: '/',
            httpOnly: true,
            sameSite: 'strict',
            secure: process.env.NODE_ENV === 'production',
            maxAge: 60 * 60 * 24 * 7 // 1 week
        });

        throw redirect(303, '/');
    },
    generateOTP: async ({ request }) => {
        console.log('Generate OTP Action Called');
        const data = await request.formData();
        const username = data.get('username') as string;
        
        console.log('Generate OTP for:', username);

        if (!username) return fail(400, { error: 'Username required' });
        
        const { secret, imageUrl } = await auth.generateOTP(username);
        console.log('OTP Generated');
        return { secret, imageUrl, step: 'otp-generated' };
    }
};
