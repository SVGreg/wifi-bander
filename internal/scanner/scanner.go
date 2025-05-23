package scanner

import (
	"fmt"
	"runtime"
	"strings"
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

// calculateQuality estimates signal quality as a percentage
func calculateQuality(signal int) int {
	// Convert dBm to quality percentage (rough approximation)
	// -30 dBm = 100%, -90 dBm = 0%
	if signal >= -30 {
		return 100
	}
	if signal <= -90 {
		return 0
	}
	// Linear approximation
	quality := ((signal + 90) * 100) / 60
	if quality < 0 {
		quality = 0
	}
	if quality > 100 {
		quality = 100
	}
	return quality
}

// getVendorFromMAC extracts vendor information from MAC address (simplified)
func getVendorFromMAC(mac string) string {
	if len(mac) < 8 {
		return "Unknown"
	}

	// Extract first 3 octets (OUI)
	oui := strings.ToUpper(strings.ReplaceAll(mac[:8], ":", ""))

	// Common vendor prefixes (simplified lookup)
	vendors := map[string]string{
		"00:1B:63": "Apple",
		"00:23:DF": "Apple",
		"00:26:BB": "Apple",
		"04:0C:CE": "Apple",
		"04:D3:CF": "Apple",
		"08:00:07": "Apple",
		"0C:74:C2": "Apple",
		"10:9A:DD": "Apple",
		"14:7D:DA": "Apple",
		"18:AF:61": "Apple",
		"1C:AB:A7": "Apple",
		"20:C9:D0": "Apple",
		"28:E0:2C": "Apple",
		"2C:B4:3A": "Apple",
		"30:90:AB": "Apple",
		"34:A3:95": "Apple",
		"38:2D:E8": "Apple",
		"3C:2E:FF": "Apple",
		"40:B3:95": "Apple",
		"44:D8:84": "Apple",
		"48:74:6E": "Apple",
		"4C:57:CA": "Apple",
		"50:05:DA": "Apple",
		"54:72:4F": "Apple",
		"58:66:BA": "Apple",
		"5C:95:AE": "Apple",
		"60:33:4B": "Apple",
		"64:76:BA": "Apple",
		"68:AB:BC": "Apple",
		"6C:72:20": "Apple",
		"70:48:0F": "Apple",
		"74:E2:F5": "Apple",
		"78:31:C1": "Apple",
		"7C:6D:62": "Apple",
		"80:92:9F": "Apple",
		"84:38:35": "Apple",
		"88:1F:A1": "Apple",
		"8C:85:90": "Apple",
		"90:72:40": "Apple",
		"94:9A:A8": "Apple",
		"98:F0:AB": "Apple",
		"9C:04:EB": "Apple",
		"A0:99:9B": "Apple",
		"A4:5E:60": "Apple",
		"A8:66:7F": "Apple",
		"AC:87:A3": "Apple",
		"B0:65:BD": "Apple",
		"B4:F0:AB": "Apple",
		"B8:78:2E": "Apple",
		"BC:52:B7": "Apple",
		"C0:B6:58": "Apple",
		"C4:B3:01": "Apple",
		"C8:BC:C8": "Apple",
		"CC:25:EF": "Apple",
		"D0:23:DB": "Apple",
		"D4:90:9C": "Apple",
		"D8:30:62": "Apple",
		"DC:2B:2A": "Apple",
		"E0:C9:7A": "Apple",
		"E4:CE:8F": "Apple",
		"E8:80:2E": "Apple",
		"EC:35:86": "Apple",
		"F0:18:98": "Apple",
		"F4:0F:24": "Apple",
		"F8:1E:DF": "Apple",
		"FC:25:3F": "Apple",

		// TP-Link
		"EC:08:6B": "TP-Link",
		"F4:F2:6D": "TP-Link",
		"A4:2B:B0": "TP-Link",
		"C4:E9:84": "TP-Link",
		"50:C7:BF": "TP-Link",
		"AC:84:C6": "TP-Link",
		"B0:4E:26": "TP-Link",
		"98:DA:C4": "TP-Link",
		"14:CF:92": "TP-Link",
		"E8:DE:27": "TP-Link",

		// ASUS
		"2C:56:DC": "ASUS",
		"1C:87:2C": "ASUS",
		"AC:9E:17": "ASUS",
		"04:D4:C4": "ASUS",
		"30:5A:3A": "ASUS",
		"50:46:5D": "ASUS",
		"70:4D:7B": "ASUS",
		"9C:5C:8E": "ASUS",
		"B0:6E:BF": "ASUS",

		// Netgear
		"84:1B:5E": "Netgear",
		"A0:04:60": "Netgear",
		"C4:04:15": "Netgear",
		"E0:91:F5": "Netgear",
		"9C:3D:CF": "Netgear",
		"CC:40:D0": "Netgear",

		// Linksys
		"C8:D7:19": "Linksys",
		"48:F8:B3": "Linksys",
		"94:10:3E": "Linksys",
		"20:AA:4B": "Linksys",
	}

	// Check exact match first
	for prefix, vendor := range vendors {
		if strings.HasPrefix(mac, prefix) {
			return vendor
		}
	}

	// Check OUI prefix
	if len(oui) >= 6 {
		ouiPrefix := oui[:6]
		for prefix, vendor := range vendors {
			vendorOUI := strings.ReplaceAll(prefix, ":", "")
			if ouiPrefix == vendorOUI {
				return vendor
			}
		}
	}

	return "Unknown"
}
