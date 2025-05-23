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
	// 2.4GHz band
	if channel >= 1 && channel <= 13 {
		return 2412 + (channel-1)*5
	}
	if channel == 14 {
		return 2484
	}

	// 5GHz band (simplified mapping)
	if channel >= 36 && channel <= 64 {
		return 5180 + (channel-36)*5
	}
	if channel >= 100 && channel <= 144 {
		return 5500 + (channel-100)*5
	}
	if channel >= 149 && channel <= 165 {
		return 5745 + (channel-149)*5
	}

	// Default fallback
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
