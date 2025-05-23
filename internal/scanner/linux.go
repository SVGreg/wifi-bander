package scanner

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/svgreg/wifi-bander/internal/analyzer"
)

// LinuxScanner implements WiFi scanning for Linux systems
type LinuxScanner struct{}

// Scan performs WiFi network scanning on Linux
func (l *LinuxScanner) Scan() ([]WiFiNetwork, error) {
	// Try nmcli first (NetworkManager)
	networks, err := l.scanWithNmcli()
	if err == nil {
		return networks, nil
	}

	// Fallback to iwlist
	return l.scanWithIwlist()
}

// scanWithNmcli uses NetworkManager's nmcli to scan for networks
func (l *LinuxScanner) scanWithNmcli() ([]WiFiNetwork, error) {
	// Enhanced nmcli command to get more fields
	cmd := exec.Command("nmcli", "-t", "-f", "SSID,CHAN,SIGNAL,FREQ,SECURITY,MODE,BSSID", "dev", "wifi")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("nmcli command failed: %v", err)
	}

	return l.parseNmcliOutput(string(output))
}

// parseNmcliOutput parses the output from nmcli command
func (l *LinuxScanner) parseNmcliOutput(output string) ([]WiFiNetwork, error) {
	var networks []WiFiNetwork
	channelMap := make(map[int]*ChannelInfo)

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			continue
		}

		ssid := parts[0]
		if ssid == "" || ssid == "--" {
			continue
		}

		channel, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		signal, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}

		frequency, err := strconv.Atoi(parts[3])
		if err != nil {
			continue
		}

		// Extract additional fields if available
		security := "Unknown"
		mode := "Infrastructure"
		bssid := "Unknown"

		if len(parts) >= 5 && parts[4] != "" {
			security = parts[4]
		}
		if len(parts) >= 6 && parts[5] != "" {
			mode = parts[5]
		}
		if len(parts) >= 7 && parts[6] != "" {
			bssid = parts[6]
		}

		network := l.createWiFiNetworkEnhanced(ssid, channel, signal, frequency, security, mode, bssid)
		networks = append(networks, network)
		updateChannelMap(channelMap, channel, signal)
	}

	// Calculate congestion scores
	for i := range networks {
		networks[i].CongestionScore = analyzer.CalculateCongestionScore(networks[i], channelMap)
	}

	return networks, nil
}

// scanWithIwlist uses iwlist as a fallback scanning method
func (l *LinuxScanner) scanWithIwlist() ([]WiFiNetwork, error) {
	// Find WiFi interface
	iface, err := l.findWiFiInterface()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("sudo", "iwlist", iface, "scan")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("iwlist command failed: %v", err)
	}

	return l.parseIwlistOutput(string(output))
}

// findWiFiInterface finds an available WiFi interface
func (l *LinuxScanner) findWiFiInterface() (string, error) {
	cmd := exec.Command("iw", "dev")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to find WiFi interface: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Interface") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", fmt.Errorf("no WiFi interface found")
}

// parseIwlistOutput parses the output from iwlist command
func (l *LinuxScanner) parseIwlistOutput(output string) ([]WiFiNetwork, error) {
	var networks []WiFiNetwork
	channelMap := make(map[int]*ChannelInfo)

	// Split by Cell entries
	cells := strings.Split(output, "Cell ")
	for _, cell := range cells[1:] { // Skip first empty entry
		network, err := l.parseIwlistCell(cell)
		if err != nil {
			continue
		}

		networks = append(networks, network)
		updateChannelMap(channelMap, network.Channel, network.Signal)
	}

	// Calculate congestion scores
	for i := range networks {
		networks[i].CongestionScore = analyzer.CalculateCongestionScore(networks[i], channelMap)
	}

	return networks, nil
}

// parseIwlistCell parses a single cell from iwlist output
func (l *LinuxScanner) parseIwlistCell(cell string) (WiFiNetwork, error) {
	var network WiFiNetwork
	lines := strings.Split(cell, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "ESSID:") {
			ssid := strings.Split(line, "ESSID:")[1]
			ssid = strings.Trim(ssid, "\"")
			network.SSID = ssid
		} else if strings.Contains(line, "Address:") {
			parts := strings.Fields(line)
			if len(parts) >= 5 {
				network.BSSID = parts[4]
				network.Vendor = getVendorFromMAC(parts[4])
			}
		} else if strings.Contains(line, "Channel:") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "Channel:" && i+1 < len(parts) {
					if ch, err := strconv.Atoi(parts[i+1]); err == nil {
						network.Channel = ch
					}
				}
			}
		} else if strings.Contains(line, "Signal level=") {
			parts := strings.Split(line, "Signal level=")
			if len(parts) > 1 {
				signalStr := strings.Fields(parts[1])[0]
				if strings.Contains(signalStr, "/") {
					// Format like "70/70"
					parts := strings.Split(signalStr, "/")
					if sig, err := strconv.Atoi(parts[0]); err == nil {
						network.Signal = sig - 100 // Convert to dBm-like
					}
				} else {
					// Direct dBm value
					if sig, err := strconv.Atoi(strings.TrimSuffix(signalStr, " dBm")); err == nil {
						network.Signal = sig
					}
				}
			}
		} else if strings.Contains(line, "Frequency:") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "Frequency:" && i+1 < len(parts) {
					freqStr := strings.TrimSuffix(parts[i+1], " GHz")
					if freq, err := strconv.ParseFloat(freqStr, 64); err == nil {
						network.Frequency = int(freq * 1000) // Convert to MHz
					}
				}
			}
		} else if strings.Contains(line, "Encryption key:") {
			if strings.Contains(line, "on") {
				network.Security = "WPA/WPA2" // Default assumption for encrypted
			} else {
				network.Security = "Open"
			}
		} else if strings.Contains(line, "IE: IEEE 802.11i/WPA2") {
			network.Security = "WPA2"
		} else if strings.Contains(line, "IE: WPA Version 1") {
			network.Security = "WPA"
		} else if strings.Contains(line, "Extra:") && strings.Contains(line, "wpa_ie") {
			network.Security = "WPA"
		}
	}

	// Determine band and complete network info
	if network.Frequency > 5000 {
		network.Band = "5G"
	} else {
		network.Band = "2.4G"
	}

	// Set default frequency if not parsed
	if network.Frequency == 0 && network.Channel != 0 {
		network.Frequency = channelToFrequency(network.Channel)
	}

	// Fill in additional properties
	network.StationCount = estimateStationCount(network.Signal, network.Channel)
	network.Quality = calculateQuality(network.Signal)

	if network.Security == "" {
		network.Security = "Unknown"
	}
	if network.PHYMode == "" {
		network.PHYMode = "Unknown"
	}
	if network.ChannelWidth == "" {
		network.ChannelWidth = "Unknown"
	}
	if network.NetworkType == "" {
		network.NetworkType = "Infrastructure"
	}
	if network.BSSID == "" {
		network.BSSID = "Unknown"
	}
	if network.Vendor == "" {
		network.Vendor = "Unknown"
	}

	if network.SSID == "" || network.Channel == 0 {
		return network, fmt.Errorf("incomplete network data")
	}

	return network, nil
}

// createWiFiNetworkEnhanced creates a WiFiNetwork struct with enhanced information
func (l *LinuxScanner) createWiFiNetworkEnhanced(ssid string, channel, signal, frequency int, security, mode, bssid string) WiFiNetwork {
	band := "2.4G"
	if frequency > 5000 {
		band = "5G"
	}

	// Estimate channel width based on frequency and band
	channelWidth := "20MHz"
	if band == "5G" {
		channelWidth = "80MHz" // Common default for 5GHz
	}

	return WiFiNetwork{
		SSID:         ssid,
		Channel:      channel,
		Signal:       signal,
		Band:         band,
		Frequency:    frequency,
		StationCount: estimateStationCount(signal, channel),
		Quality:      calculateQuality(signal),
		Security:     security,
		PHYMode:      "Unknown", // Would need more detailed parsing
		ChannelWidth: channelWidth,
		NetworkType:  mode,
		BSSID:        bssid,
		Vendor:       getVendorFromMAC(bssid),
		Noise:        0, // Not easily available from nmcli
		SNR:          0, // Not easily available from nmcli
	}
}

// createWiFiNetwork creates a WiFiNetwork struct from parsed parameters (legacy method)
func (l *LinuxScanner) createWiFiNetwork(ssid string, channel, signal, frequency int) WiFiNetwork {
	return l.createWiFiNetworkEnhanced(ssid, channel, signal, frequency, "Unknown", "Infrastructure", "Unknown")
}
