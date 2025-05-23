package scanner

import (
	"fmt"
	"runtime"
)

// ScanWiFiNetworks detects the operating system and calls the appropriate scanner
func ScanWiFiNetworks() ([]WiFiNetwork, error) {
	switch runtime.GOOS {
	case "linux":
		scanner := &LinuxScanner{}
		return scanner.Scan()
	case "darwin":
		scanner := &MacOSScanner{}
		return scanner.Scan()
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// channelToFrequency converts a WiFi channel number to frequency in MHz
func channelToFrequency(channel int) int {
	// 2.4GHz band (channels 1-14)
	if channel >= 1 && channel <= 13 {
		return 2412 + (channel-1)*5
	}
	if channel == 14 {
		return 2484 // Special case for Japan
	}

	// 5GHz band - comprehensive mapping
	// UNII-1 Band (5.15-5.25 GHz) - channels 36, 40, 44, 48
	if channel >= 36 && channel <= 48 {
		return 5180 + (channel-36)*5
	}

	// UNII-2A Band (5.25-5.35 GHz) - channels 52, 56, 60, 64 (DFS required)
	if channel >= 52 && channel <= 64 {
		return 5260 + (channel-52)*5
	}

	// UNII-2C Band (5.47-5.725 GHz) - channels 100-144 (DFS required)
	if channel >= 100 && channel <= 144 {
		return 5500 + (channel-100)*5
	}

	// UNII-3 Band (5.725-5.875 GHz) - channels 149, 153, 157, 161, 165
	if channel >= 149 && channel <= 165 {
		return 5745 + (channel-149)*5
	}

	// UNII-4 Band (5.85-5.925 GHz) - channels 169, 173, 177 (newer allocation)
	if channel >= 169 && channel <= 177 {
		return 5845 + (channel-169)*5
	}

	// Default fallback for unknown channels
	return 2412
}

// updateChannelMap updates the channel tracking map with network information
func updateChannelMap(channelMap map[int]*ChannelInfo, channel, signal int) {
	if channelInfo, exists := channelMap[channel]; exists {
		channelInfo.NetworkCount++
		if signal > channelInfo.StrongestRSSI {
			channelInfo.StrongestRSSI = signal
		}
	} else {
		channelMap[channel] = &ChannelInfo{
			Channel:       channel,
			NetworkCount:  1,
			StrongestRSSI: signal,
		}
	}
}

// estimateStationCount estimates the number of stations based on signal patterns
func estimateStationCount(signal, channel int) int {
	baseCount := 1

	// Stronger signals might indicate more active networks
	if signal > -40 {
		baseCount += 3
	} else if signal > -60 {
		baseCount += 2
	} else if signal > -80 {
		baseCount += 1
	}

	// 2.4GHz tends to be more congested
	if channel <= 14 {
		baseCount += 1
	}

	return baseCount
}
