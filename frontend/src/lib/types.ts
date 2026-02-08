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
}
