import bcrypt from 'bcryptjs';
import { authenticator } from '@otplib/preset-default';
import qrcode from 'qrcode';
import jwt from 'jsonwebtoken';
import { env } from '$env/dynamic/private';

// In-memory user store for simplicity (as requested MVP behavior)
// In a real app, use a database (SQLite/Postgres)
// We persist this to a JSON file or just keep in memory since this is a "service" that might stay up.
// For now, let's use a simple global variable, but note it resets on server restart. 
// Given the requirements, we might want to persist to disk if possible, but let's start with in-memory.



import { db } from './db';

const JWT_SECRET = env.JWT_SECRET || 'super-secret-dev-key-change-me';

export const auth = {
    // Check if any user exists (for setup flow)
    hasUsers: () => db.hasUsers(),

    createUser: async (username: string, password: string, otpSecret?: string) => {
        if (db.findUser(username)) throw new Error("User already exists");
        const passwordHash = await bcrypt.hash(password, 10);
        db.addUser({ username, passwordHash, otpSecret });
        return { username };
    },

    verifyUser: async (username: string, password: string) => {
        const user = db.findUser(username);
        if (!user) return null;
        
        const validPass = await bcrypt.compare(password, user.passwordHash);
        if (!validPass) return null;

        return { username: user.username, hasOTP: !!user.otpSecret };
    },

    // OTP Utils
    generateOTP: async (username: string) => {
        const secret = authenticator.generateSecret();
        const otpauth = authenticator.keyuri(username, 'Wolite', secret);
        const imageUrl = await qrcode.toDataURL(otpauth);
        return { secret, imageUrl };
    },

    verifyOTP: (token: string, secret: string) => {
        return authenticator.check(token, secret);
    },

    getUserOTPSecret: (username: string) => {
        const user = db.findUser(username);
        return user?.otpSecret;
    },

    // Session Utils
    createSessionToken: (username: string) => {
        return jwt.sign({ username }, JWT_SECRET, { expiresIn: '7d' });
    },

    verifySessionToken: (token: string) => {
        try {
            return jwt.verify(token, JWT_SECRET) as { username: string };
        } catch {
            return null;
        }
    }
};
