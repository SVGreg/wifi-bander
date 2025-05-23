package scanner

// WiFiNetwork represents a detected WiFi network with all its properties
type WiFiNetwork struct {
	SSID            string // Network name
	Channel         int    // WiFi channel (1-13 for 2.4GHz, 36+ for 5GHz)
	Signal          int    // Signal strength in dBm
	Band            string // "2.4G" or "5G"
	CongestionScore int    // Calculated congestion level
	Frequency       int    // Frequency in MHz
	StationCount    int    // Estimated connected devices
}

// Interface methods for analyzer package compatibility
func (w WiFiNetwork) GetBand() string      { return w.Band }
func (w WiFiNetwork) GetChannel() int      { return w.Channel }
func (w WiFiNetwork) GetSignal() int       { return w.Signal }
func (w WiFiNetwork) GetStationCount() int { return w.StationCount }

// Interface methods for display package compatibility
func (w WiFiNetwork) GetSSID() string         { return w.SSID }
func (w WiFiNetwork) GetCongestionScore() int { return w.CongestionScore }
func (w WiFiNetwork) GetFrequency() int       { return w.Frequency }

// ChannelInfo holds aggregated information about a specific channel
type ChannelInfo struct {
	Channel       int // Channel number
	NetworkCount  int // Number of networks on this channel
	StrongestRSSI int // Strongest signal strength seen on this channel
}

// Scanner defines the interface for WiFi network scanning
type Scanner interface {
	Scan() ([]WiFiNetwork, error)
}
