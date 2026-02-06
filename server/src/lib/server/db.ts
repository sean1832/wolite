import fs from 'node:fs';
import path from 'node:path';

const DATA_DIR = path.join(process.cwd(), 'data');
const DB_FILE = 'users.json';
const dbPath = path.join(DATA_DIR, DB_FILE);

// Ensure data directory exists
function ensureDataDir() {
    if (!fs.existsSync(DATA_DIR)) {
        try {
            fs.mkdirSync(DATA_DIR, { recursive: true });
        } catch (e) {
            console.error("Failed to create data dir:", e);
        }
    }
}

// Run once on load
ensureDataDir();

export interface User {
    username: string;
    passwordHash: string;
    otpSecret?: string;
}

export const db = {
    getUsers: (): User[] => {
        if (!fs.existsSync(dbPath)) {
            return [];
        }
        try {
            const data = fs.readFileSync(dbPath, 'utf-8');
            return JSON.parse(data);
        } catch (e) {
            console.error("Failed to read database:", e);
            return [];
        }
    },

    saveUsers: (users: User[]) => {
        try {
            ensureDataDir();
            fs.writeFileSync(dbPath, JSON.stringify(users, null, 2));
        } catch (e) {
            console.error("Failed to write database:", e);
            throw e; // Re-throw to see error in actions
        }
    },

    addUser: (user: User) => {
        const users = db.getUsers();
        users.push(user);
        db.saveUsers(users);
    },

    findUser: (username: string): User | undefined => {
        const users = db.getUsers();
        return users.find(u => u.username === username);
    },
    
    hasUsers: (): boolean => {
        const users = db.getUsers();
        return users.length > 0;
    }
};
