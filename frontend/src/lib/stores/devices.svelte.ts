import { type Device } from '$lib/types';
<<<<<<< HEAD

class DeviceStore {
	devices = $state<Device[]>([
		{ id: '1', name: 'Gaming PC', ip: '192.168.1.10', mac: '00:11:22:33:44:55', status: 'offline' },
		{ id: '2', name: 'Home Server', ip: '192.168.1.20', mac: 'AA:BB:CC:DD:EE:FF', status: 'online' },
		{ id: '3', name: 'Living Room TV', ip: '192.168.1.30', mac: '11:22:33:44:55:66', status: 'offline' }
	]);

	addDevice(device: Omit<Device, 'id' | 'status'>) {
		const newDevice: Device = {
			...device,
			id: crypto.randomUUID(),
			status: 'unknown'
		};
		this.devices.push(newDevice);
	}

	removeDevice(id: string) {
		this.devices = this.devices.filter((d) => d.id !== id);
	}

	updateDevice(id: string, data: Partial<Omit<Device, 'id'>>) {
		const index = this.devices.findIndex((d) => d.id === id);
		if (index !== -1) {
			this.devices[index] = { ...this.devices[index], ...data };
		}
	}

	wakeDevice(id: string) {
		const index = this.devices.findIndex((d) => d.id === id);
		if (index !== -1) {
			this.devices[index].status = 'waking';
			// Mock waking process
			setTimeout(() => {
				if (this.devices[index]) { // Check if still exists
					this.devices[index].status = 'online';
				}
			}, 3000);
=======
import { http } from '$lib/api';

class DeviceStore {
	devices = $state<Device[]>([]);
	loading = $state(false);
	error = $state<string | null>(null);

	/**
	 * Initialize the store by fetching devices from the API
	 * Call this when the app loads
	 */
	async init(fetch: typeof window.fetch) {
		this.loading = true;
		this.error = null;
		try {
			const devices = await http.get<Device[]>(fetch, '/devices');
			this.devices = devices;
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to load devices';
			console.error('Failed to load devices:', err);
		} finally {
			this.loading = false;
		}
	}

	async addDevice(fetch: typeof window.fetch, device: Omit<Device, 'status'>) {
		this.loading = true;
		this.error = null;
		try {
			const newDevice = await http.post<Device>(fetch, '/devices', device);
			this.devices.push(newDevice);
			return newDevice;
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to add device';
			console.error('Failed to add device:', err);
			throw err;
		} finally {
			this.loading = false;
		}
	}

	async removeDevice(fetch: typeof window.fetch, macAddress: string) {
		this.loading = true;
		this.error = null;
		try {
			await http.delete(fetch, `/devices/${macAddress}`);
			this.devices = this.devices.filter((d) => d.mac_address !== macAddress);
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to remove device';
			console.error('Failed to remove device:', err);
			throw err;
		} finally {
			this.loading = false;
		}
	}

	async updateDevice(fetch: typeof window.fetch, macAddress: string, data: Partial<Omit<Device, 'mac_address' | 'status'>>) {
		this.loading = true;
		this.error = null;
		try {
			const updatedDevice = await http.put<Device>(fetch, `/devices/${macAddress}`, data);
			const index = this.devices.findIndex((d) => d.mac_address === macAddress);
			if (index !== -1) {
				this.devices[index] = updatedDevice;
			}
			return updatedDevice;
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to update device';
			console.error('Failed to update device:', err);
			throw err;
		} finally {
			this.loading = false;
		}
	}

	async wakeDevice(fetch: typeof window.fetch, macAddress: string) {
		this.loading = true;
		this.error = null;
		try {
			await http.post(fetch, `/devices/${macAddress}/wake`, {});
			// Optionally update local state to show "waking" status
			const index = this.devices.findIndex((d) => d.mac_address === macAddress);
			if (index !== -1) {
				// Note: Backend doesn't return updated device, so we manually update
				// In a production app, you might poll for status updates
				this.devices[index] = { ...this.devices[index] };
			}
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Failed to wake device';
			console.error('Failed to wake device:', err);
			throw err;
		} finally {
			this.loading = false;
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
		}
	}
}

export const deviceStore = new DeviceStore();
