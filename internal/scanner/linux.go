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
	cmd := exec.Command("nmcli", "-t", "-f", "SSID,CHAN,SIGNAL,FREQ", "dev", "wifi")
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

		network := l.createWiFiNetwork(ssid, channel, signal, frequency)
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

	network.StationCount = estimateStationCount(network.Signal, network.Channel)

	if network.SSID == "" || network.Channel == 0 {
		return network, fmt.Errorf("incomplete network data")
	}

	return network, nil
}

// createWiFiNetwork creates a WiFiNetwork struct from parsed parameters
func (l *LinuxScanner) createWiFiNetwork(ssid string, channel, signal, frequency int) WiFiNetwork {
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
	}
}
