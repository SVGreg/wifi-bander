package display

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/svgreg/wifi-bander/internal/analyzer"
)

// WiFiNetwork interface for display purposes
type WiFiNetwork interface {
	GetSSID() string
	GetBand() string
	GetChannel() int
	GetSignal() int
	GetStationCount() int
	GetCongestionScore() int
	GetFrequency() int
	GetSecurity() string
	GetPHYMode() string
	GetChannelWidth() string
	GetNetworkType() string
	GetBSSID() string
	GetVendor() string
	GetQuality() int
	GetNoise() int
	GetSNR() int
}

// DisplayResults shows the WiFi scan results in a comprehensive formatted table
func DisplayResults(networks []WiFiNetwork) {
	fmt.Printf("\n=== WiFi Network Analysis - %s ===\n", time.Now().Format("15:04:05"))

	if len(networks) == 0 {
		fmt.Println("No networks detected.")
		return
	}

	// Create a new tabwriter with wider spacing for comprehensive data
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	// Comprehensive header
	fmt.Fprintln(w, "SSID\tBand\tCh\tSignal\tQuality\tSecurity\tPHY Mode\tWidth\tVendor\tCongestion\tFreq\t")
	fmt.Fprintln(w, "----\t----\t--\t------\t-------\t--------\t--------\t-----\t------\t----------\t----\t")

	// Print each network's comprehensive information
	for _, net := range networks {
		congestionLevel := GetCongestionLevel(net.GetCongestionScore())

		// Truncate long values for better table formatting, but make Security and PHY Mode wider
		ssid := truncateString(net.GetSSID(), 16)
		security := truncateString(net.GetSecurity(), 18) // Increased from 10 to 18
		phyMode := truncateString(net.GetPHYMode(), 15)   // Increased from 10 to 15
		vendor := truncateString(net.GetVendor(), 8)

		fmt.Fprintf(w, "%s\t%s\t%d\t%d dBm\t%d%%\t%s\t%s\t%s\t%s\t%s\t%d\t\n",
			ssid,
			net.GetBand(),
			net.GetChannel(),
			net.GetSignal(),
			net.GetQuality(),
			security,
			phyMode,
			net.GetChannelWidth(),
			vendor,
			congestionLevel,
			net.GetFrequency(),
		)
	}

	// Flush the tabwriter to display the table
	w.Flush()

	// Show network count summary
	fmt.Printf("\nTotal networks detected: %d\n", len(networks))
}

// truncateString truncates a string to a maximum length for table formatting
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// DisplayCompactResults shows a compact view of the networks
func DisplayCompactResults(networks []WiFiNetwork) {
	fmt.Printf("\n=== WiFi Networks (Compact View) - %s ===\n", time.Now().Format("15:04:05"))

	if len(networks) == 0 {
		fmt.Println("No networks detected.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Compact header
	fmt.Fprintln(w, "SSID\tBand\tChannel\tSignal\tSecurity\tCongestion\t")
	fmt.Fprintln(w, "----\t----\t-------\t------\t--------\t----------\t")

	// Print compact network information
	for _, net := range networks {
		congestionLevel := GetCongestionLevel(net.GetCongestionScore())
		ssid := truncateString(net.GetSSID(), 20)
		security := truncateString(net.GetSecurity(), 12)

		fmt.Fprintf(w, "%s\t%s\t%d\t%d dBm\t%s\t%s\t\n",
			ssid,
			net.GetBand(),
			net.GetChannel(),
			net.GetSignal(),
			security,
			congestionLevel,
		)
	}

	w.Flush()
}

// DisplayRecommendations shows channel recommendations
func DisplayRecommendations(networks []analyzer.WiFiNetwork) {
	recommendations := analyzer.GetChannelRecommendations(networks)

	fmt.Println("\n=== Channel Recommendations (Top 3 Optimal Choices) ===")
	fmt.Println("Advanced analysis considering frequency separation, signal strength, and interference patterns")

	for band, recs := range recommendations {
		fmt.Printf("\n🔸 %s Band Recommendations:\n", band)

		if len(recs) == 0 {
			fmt.Printf("  No recommendations available for %s band\n", band)
			continue
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Rank\tChannel\tFreq(MHz)\tInterference\tGap(MHz)\tReasoning\t")
		fmt.Fprintln(w, "----\t-------\t---------\t-----------\t--------\t---------\t")

		for i, rec := range recs {
			rank := fmt.Sprintf("#%d", i+1)
			gap := "N/A"
			if rec.FrequencyGap > 0 {
				gap = fmt.Sprintf("%d", rec.FrequencyGap)
			}

			fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%s\t%s\t\n",
				rank,
				rec.Channel,
				rec.Frequency,
				rec.InterferenceLevel,
				gap,
				rec.Reasoning,
			)
		}
		w.Flush()

		// Show frequency separation analysis
		if len(recs) >= 2 {
			fmt.Printf("\n  📊 Frequency Separation Analysis:\n")
			for i := 0; i < len(recs)-1; i++ {
				separation := abs(recs[i].Frequency - recs[i+1].Frequency)
				fmt.Printf("     • Channel %d ↔ Channel %d: %d MHz separation\n",
					recs[i].Channel, recs[i+1].Channel, separation)
			}
		}

		// Show band-specific advice
		if band == "2.4G" {
			fmt.Printf("\n  💡 2.4GHz Advice: Prefer channels 1, 6, or 11 (non-overlapping). Avoid channels with strong nearby signals.\n")
		} else {
			fmt.Printf("\n  💡 5GHz Advice: More spectrum available. DFS channels may require radar detection but are often less congested.\n")
		}
	}

	fmt.Println("\n🎯 Configuration Tips:")
	fmt.Println("   • Choose the #1 ranked channel for optimal performance")
	fmt.Println("   • Monitor performance and try #2 or #3 if issues occur")
	fmt.Println("   • Consider channel width: 80MHz for 5GHz, 20MHz for 2.4GHz in crowded areas")
	fmt.Println("   • Update analysis periodically as WiFi landscape changes")
	fmt.Println("\nPress Ctrl+C to exit...")
}

// abs helper function for frequency calculations
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetCongestionLevel returns a human-readable congestion level
func GetCongestionLevel(score int) string {
	switch {
	case score <= 15:
		return "Low"
	case score <= 30:
		return "Medium"
	case score <= 50:
		return "High"
	default:
		return "Very High"
	}
}

// DisplayChannelInfo shows detailed information about detected and available channels
func DisplayChannelInfo(networks []analyzer.WiFiNetwork) {
	fmt.Println("\n=== Channel Analysis ===")

	// Get detected channels
	detected24 := getDetectedChannelsByBand(networks, "2.4G")
	detected5 := getDetectedChannelsByBand(networks, "5G")

	// Get channel info
	channelInfo := analyzer.GetChannelInfo()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// 2.4GHz information
	fmt.Fprintln(w, "\n2.4GHz Band Analysis:")
	fmt.Fprintf(w, "Detected channels:\t%v\t\n", detected24)
	fmt.Fprintf(w, "Non-overlapping (optimal):\t%v\t\n", channelInfo["2.4GHz"].(map[string]interface{})["non_overlapping"])
	fmt.Fprintf(w, "US standard (1-11):\t%v\t\n", channelInfo["2.4GHz"].(map[string]interface{})["us_channels"])
	fmt.Fprintf(w, "EU standard (1-13):\t%v\t\n", channelInfo["2.4GHz"].(map[string]interface{})["eu_channels"])

	// 5GHz information
	fmt.Fprintln(w, "\n5GHz Band Analysis:")
	fmt.Fprintf(w, "Detected channels:\t%v\t\n", detected5)
	fmt.Fprintf(w, "UNII-1 (36-48):\t%v\t\n", channelInfo["5GHz"].(map[string]interface{})["unii_1"])
	fmt.Fprintf(w, "UNII-2A (52-64, DFS):\t%v\t\n", channelInfo["5GHz"].(map[string]interface{})["unii_2a"])
	fmt.Fprintf(w, "UNII-2C (100-144, DFS):\t%v\t\n", channelInfo["5GHz"].(map[string]interface{})["unii_2c"])
	fmt.Fprintf(w, "UNII-3 (149-165):\t%v\t\n", channelInfo["5GHz"].(map[string]interface{})["unii_3"])
	fmt.Fprintf(w, "UNII-4 (169-177):\t%v\t\n", channelInfo["5GHz"].(map[string]interface{})["unii_4"])

	w.Flush()

	// Channel usage statistics
	displayChannelUsageStats(networks)
}

// getDetectedChannelsByBand extracts detected channels for a specific band
func getDetectedChannelsByBand(networks []analyzer.WiFiNetwork, band string) []int {
	channelSet := make(map[int]bool)
	for _, network := range networks {
		if network.GetBand() == band {
			channelSet[network.GetChannel()] = true
		}
	}

	var channels []int
	for ch := range channelSet {
		channels = append(channels, ch)
	}

	sort.Ints(channels)
	return channels
}

// displayChannelUsageStats shows how many networks are on each channel
func displayChannelUsageStats(networks []analyzer.WiFiNetwork) {
	fmt.Println("\n=== Channel Usage Statistics ===")

	usage24 := make(map[int]int)
	usage5 := make(map[int]int)

	// Count networks per channel
	for _, network := range networks {
		if network.GetBand() == "2.4G" {
			usage24[network.GetChannel()]++
		} else {
			usage5[network.GetChannel()]++
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	// 2.4GHz horizontal display - show channels 1-13 (EU standard, most comprehensive)
	fmt.Fprintln(w, "\n2.4GHz Channel Usage:")

	// Channel headers
	fmt.Fprint(w, "Channel\t")
	for ch := 1; ch <= 13; ch++ {
		fmt.Fprintf(w, "%d\t", ch)
	}
	fmt.Fprintln(w)

	// Separator line
	fmt.Fprint(w, "-------\t")
	for ch := 1; ch <= 13; ch++ {
		fmt.Fprint(w, "-\t")
	}
	fmt.Fprintln(w)

	// Network counts
	fmt.Fprint(w, "Networks\t")
	for ch := 1; ch <= 13; ch++ {
		count := usage24[ch] // Will be 0 if channel not found
		fmt.Fprintf(w, "%d\t", count)
	}
	fmt.Fprintln(w)

	// 5GHz horizontal display - show all UNII band channels
	if len(usage5) > 0 {
		fmt.Fprintln(w, "\n5GHz Channel Usage:")

		// Get all 5GHz channels from analyzer
		channelInfo := analyzer.GetChannelInfo()
		allChannels5G := channelInfo["5GHz"].(map[string]interface{})["all"].([]int)

		// Channel headers
		fmt.Fprint(w, "Channel\t")
		for _, ch := range allChannels5G {
			fmt.Fprintf(w, "%d\t", ch)
		}
		fmt.Fprintln(w)

		// Separator line
		fmt.Fprint(w, "-------\t")
		for range allChannels5G {
			fmt.Fprint(w, "--\t")
		}
		fmt.Fprintln(w)

		// Network counts
		fmt.Fprint(w, "Networks\t")
		for _, ch := range allChannels5G {
			count := usage5[ch] // Will be 0 if channel not found
			fmt.Fprintf(w, "%d\t", count)
		}
		fmt.Fprintln(w)
	} else {
		fmt.Fprintln(w, "\n5GHz Channel Usage: No 5GHz networks detected")
	}

	w.Flush()
}
