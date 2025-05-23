package scanner

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/svgreg/wifi-bander/internal/analyzer"
)

// MacOSScanner implements WiFi scanning for macOS systems
type MacOSScanner struct{}

// Scan performs WiFi network scanning on macOS
func (m *MacOSScanner) Scan() ([]WiFiNetwork, error) {
	// Try using the airport command if available
	cmd := exec.Command("/usr/sbin/airport", "-s")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to system_profiler
		return m.scanWithSystemProfiler()
	}

	networks, err := m.parseAirportOutput(string(output))
	if err != nil {
		return m.scanWithSystemProfiler()
	}

	return networks, nil
}

// scanWithSystemProfiler uses system_profiler to scan WiFi networks
func (m *MacOSScanner) scanWithSystemProfiler() ([]WiFiNetwork, error) {
	cmd := exec.Command("system_profiler", "SPAirPortDataType")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("system_profiler command failed: %v", err)
	}

	return m.parseSystemProfilerOutput(string(output))
}

// parseAirportOutput parses the output from the airport command
func (m *MacOSScanner) parseAirportOutput(output string) ([]WiFiNetwork, error) {
	var networks []WiFiNetwork
	channelMap := make(map[int]*ChannelInfo)

	lines := strings.Split(output, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if i == 0 || line == "" {
			continue // Skip header line or empty lines
		}

		// Airport output format varies, let's be more flexible
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		var ssid string
		var signal, channel int
		var err error

		// Work backwards to find signal and channel
		if len(parts) >= 3 {
			channelIdx := len(parts) - 2
			signalIdx := len(parts) - 3

			if channelIdx >= 0 && signalIdx >= 0 {
				channel, err = strconv.Atoi(parts[channelIdx])
				if err != nil {
					continue
				}

				signal, err = strconv.Atoi(parts[signalIdx])
				if err != nil {
					continue
				}

				// Find MAC address and extract SSID
				macIdx := -1
				for j := 1; j < len(parts)-2; j++ {
					if strings.Count(parts[j], ":") == 5 && len(parts[j]) == 17 {
						macIdx = j
						break
					}
				}

				if macIdx > 0 {
					ssid = strings.Join(parts[0:macIdx], " ")
				} else {
					ssid = parts[0]
				}
			}
		}

		if ssid == "" || channel == 0 {
			continue
		}

		network := m.createWiFiNetwork(ssid, channel, signal)
		networks = append(networks, network)
		updateChannelMap(channelMap, channel, signal)
	}

	// Calculate congestion scores
	for i := range networks {
		networks[i].CongestionScore = analyzer.CalculateCongestionScore(networks[i], channelMap)
	}

	return networks, nil
}

// parseSystemProfilerOutput parses the output from system_profiler
func (m *MacOSScanner) parseSystemProfilerOutput(output string) ([]WiFiNetwork, error) {
	var networks []WiFiNetwork
	channelMap := make(map[int]*ChannelInfo)

	lines := strings.Split(output, "\n")
	var currentNetwork *WiFiNetwork
	inOtherNetworks := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check if we're entering the "Other Local Wi-Fi Networks" section
		if strings.Contains(line, "Other Local Wi-Fi Networks:") {
			inOtherNetworks = true
			continue
		}

		if !inOtherNetworks {
			continue
		}

		// Network names are indented and end with ":"
		if strings.HasSuffix(line, ":") && !strings.Contains(line, "PHY Mode") &&
			!strings.Contains(line, "Channel") && !strings.Contains(line, "Network Type") &&
			!strings.Contains(line, "Security") && !strings.Contains(line, "Signal") {

			// Save previous network
			if currentNetwork != nil && currentNetwork.SSID != "" && currentNetwork.Channel != 0 {
				networks = append(networks, *currentNetwork)
				updateChannelMap(channelMap, currentNetwork.Channel, currentNetwork.Signal)
			}

			ssid := strings.TrimSuffix(line, ":")
			currentNetwork = &WiFiNetwork{SSID: ssid}
		} else if currentNetwork != nil && strings.Contains(line, ":") {
			m.parseNetworkProperty(currentNetwork, line)
		}
	}

	// Add the last network
	if currentNetwork != nil && currentNetwork.SSID != "" && currentNetwork.Channel != 0 {
		networks = append(networks, *currentNetwork)
		updateChannelMap(channelMap, currentNetwork.Channel, currentNetwork.Signal)
	}

	// Calculate congestion scores
	for i := range networks {
		networks[i].CongestionScore = analyzer.CalculateCongestionScore(networks[i], channelMap)
	}

	return networks, nil
}

// parseNetworkProperty parses a property line for the current network
func (m *MacOSScanner) parseNetworkProperty(network *WiFiNetwork, line string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "Channel":
		// Parse channel like "36 (5GHz, 80MHz)" or "6 (2GHz, 20MHz)"
		channelParts := strings.Fields(value)
		if len(channelParts) > 0 {
			if ch, err := strconv.Atoi(channelParts[0]); err == nil {
				network.Channel = ch
				network.Frequency = channelToFrequency(ch)
				if network.Frequency > 5000 {
					network.Band = "5G"
				} else {
					network.Band = "2.4G"
				}
			}
		}

		// Extract channel width from parentheses like "(5GHz, 80MHz)"
		if strings.Contains(value, "(") && strings.Contains(value, ")") {
			parenContent := strings.Split(strings.Split(value, "(")[1], ")")[0]
			if strings.Contains(parenContent, "MHz") {
				widthParts := strings.Split(parenContent, ",")
				if len(widthParts) > 1 {
					network.ChannelWidth = strings.TrimSpace(widthParts[1])
				}
			}
		}

	case "Signal / Noise":
		// Parse signal like "-77 dBm / -94 dBm"
		signalParts := strings.Fields(value)
		if len(signalParts) >= 4 {
			// Signal strength
			signalStr := strings.TrimSuffix(signalParts[0], " dBm")
			if sig, err := strconv.Atoi(signalStr); err == nil {
				network.Signal = sig
				network.StationCount = estimateStationCount(sig, network.Channel)
			}

			// Noise level
			if len(signalParts) >= 4 {
				noiseStr := strings.TrimSuffix(signalParts[3], " dBm")
				if noise, err := strconv.Atoi(noiseStr); err == nil {
					network.Noise = noise
					// Calculate SNR
					if network.Signal != 0 && network.Noise != 0 {
						network.SNR = network.Signal - network.Noise
					}
					// Calculate quality percentage (rough approximation)
					network.Quality = calculateQuality(network.Signal)
				}
			}
		}

	case "Security":
		network.Security = value
		if network.Security == "" {
			network.Security = "Open"
		}

	case "PHY Mode":
		network.PHYMode = value

	case "Network Type":
		network.NetworkType = value

	case "BSSID":
		network.BSSID = value
		// Extract vendor from MAC address
		network.Vendor = getVendorFromMAC(value)
	}
}

// createWiFiNetwork creates a WiFiNetwork struct from basic parameters
func (m *MacOSScanner) createWiFiNetwork(ssid string, channel, signal int) WiFiNetwork {
	frequency := channelToFrequency(channel)
	band := "2.4G"
	if frequency > 5000 {
		band = "5G"
	}

	return WiFiNetwork{
		SSID:         ssid,
		Channel:      channel,
		Signal:       signal,
		Band:         band,
		Frequency:    frequency,
		StationCount: estimateStationCount(signal, channel),
		Quality:      calculateQuality(signal),
		Security:     "Unknown",
		PHYMode:      "Unknown",
		ChannelWidth: "Unknown",
		NetworkType:  "Infrastructure",
		BSSID:        "Unknown",
		Vendor:       "Unknown",
		Noise:        0,
		SNR:          0,
	}
}
