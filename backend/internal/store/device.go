package store

type Status string

const (
	StatusOnline  Status = "online"
	StatusOffline Status = "offline"
	StatusUnknown Status = "unknown"
	StatusError   Status = "error"
)

type Device struct {
	ID          int    `json:"id"` // snowflake ID or auto-incrementing integer
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MACAddress  string `json:"mac_address"`
	IPAddress   string `json:"ip_address,omitempty"`

	Status Status `json:"status"`
}
