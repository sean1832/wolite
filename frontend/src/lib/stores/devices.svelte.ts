import { type Device } from '$lib/types';

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
		}
	}
}

export const deviceStore = new DeviceStore();
