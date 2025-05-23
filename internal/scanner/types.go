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

	// Enhanced network information
	Security     string // Security type (Open, WPA, WPA2, WPA3, etc.)
	PHYMode      string // PHY mode (802.11n, 802.11ac, 802.11ax, etc.)
	ChannelWidth string // Channel width (20MHz, 40MHz, 80MHz, 160MHz)
	NetworkType  string // Network type (Infrastructure, Ad-hoc)
	BSSID        string // MAC address of access point
	Vendor       string // Vendor name (from MAC OUI lookup)
	Quality      int    // Signal quality percentage (0-100)
	Noise        int    // Noise level in dBm
	SNR          int    // Signal-to-Noise Ratio
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
func (w WiFiNetwork) GetSecurity() string     { return w.Security }
func (w WiFiNetwork) GetPHYMode() string      { return w.PHYMode }
func (w WiFiNetwork) GetChannelWidth() string { return w.ChannelWidth }
func (w WiFiNetwork) GetNetworkType() string  { return w.NetworkType }
func (w WiFiNetwork) GetBSSID() string        { return w.BSSID }
func (w WiFiNetwork) GetVendor() string       { return w.Vendor }
func (w WiFiNetwork) GetQuality() int         { return w.Quality }
func (w WiFiNetwork) GetNoise() int           { return w.Noise }
func (w WiFiNetwork) GetSNR() int             { return w.SNR }

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
