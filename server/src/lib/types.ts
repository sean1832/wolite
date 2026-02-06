export interface Device {
	id: string;
	name: string;
	ip: string;
	mac: string;
	status: 'online' | 'offline' | 'waking' | 'unknown';
}
