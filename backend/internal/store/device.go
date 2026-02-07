package store

type Status string

const (
	StatusOnline  Status = "online"
	StatusOffline Status = "offline"
	StatusUnknown Status = "unknown"
	StatusError   Status = "error"
)

type Device struct {
	MACAddress  string `json:"mac_address"`           // unique identifier for the device
	Name        string `json:"name"`                  // human-readable name for the device
	Description string `json:"description,omitempty"` // optional description of the device
	IPAddress   string `json:"ip_address,omitempty"`  // optional IP address of the device

	Status Status `json:"status"` // current status of the device
}

func NewDevice(macAddress, name, description, ipAddress string, status Status) *Device {
	return &Device{
		MACAddress:  macAddress,
		Name:        name,
		Description: description,
		IPAddress:   ipAddress,
		Status:      status,
	}
}
