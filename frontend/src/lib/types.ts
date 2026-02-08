<<<<<<< HEAD
export interface Device {
	id: string;
	name: string;
	ip: string;
	mac: string;
	status: 'online' | 'offline' | 'waking' | 'unknown';
=======
// Device represents a network device that can be woken
export interface Device {
  mac_address: string; // Unique identifier
  name: string;
  description?: string;
  ip_address: string;
  broadcast_ip: string; // For Wake-on-LAN
  status: "online" | "offline" | "unknown" | "error";
}

// API Response wrapper from backend
export interface ApiResponse<T> {
	code: number;
	message?: string;
	data?: T;
}

// User represents an authenticated user
export interface User {
	username: string;
	has_otp: boolean;
}

// Auth response for login/setup
export interface AuthResponse {
	qr_code?: string; // base64 encoded QR code image for OTP setup
	secret?: string; // OTP secret (only during setup)
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
}
