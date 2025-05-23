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
}

// DisplayResults shows the WiFi scan results in a formatted table
func DisplayResults(networks []WiFiNetwork) {
	fmt.Printf("\n=== WiFi Network Analysis - %s ===\n", time.Now().Format("15:04:05"))

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header
	fmt.Fprintln(w, "SSID\tBand\tChannel\tSignal (dBm)\tStations\tCongestion\tFreq (MHz)\t")
	fmt.Fprintln(w, "----\t----\t-------\t-----------\t--------\t----------\t---------\t")

	// Print each network's information
	for _, net := range networks {
		congestionLevel := GetCongestionLevel(net.GetCongestionScore())
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%s\t%d\t\n",
			net.GetSSID(),
			net.GetBand(),
			net.GetChannel(),
			net.GetSignal(),
			net.GetStationCount(),
			congestionLevel,
			net.GetFrequency(),
		)
	}

	// Flush the tabwriter to display the table
	w.Flush()
}

// DisplayRecommendations shows channel recommendations
func DisplayRecommendations(networks []analyzer.WiFiNetwork) {
	recommendations := analyzer.GetChannelRecommendations(networks)

	fmt.Println("\n=== Channel Recommendations ===")
	fmt.Println("Recommended channels for optimal performance:")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Band\tBest Channels\t")
	fmt.Fprintln(w, "----\t-------------\t")

	for band, channels := range recommendations {
		channelStr := ""
		for i, ch := range channels {
			if i > 0 {
				channelStr += ", "
			}
			channelStr += fmt.Sprintf("%d", ch)
		}
		fmt.Fprintf(w, "%s\t%s\t\n", band, channelStr)
	}

	w.Flush()
	fmt.Println("\nPress Ctrl+C to exit...")
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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// 2.4GHz usage
	if len(usage24) > 0 {
		fmt.Fprintln(w, "\n2.4GHz Channel Usage:")
		fmt.Fprintln(w, "Channel\tNetworks\tCongestion\t")
		fmt.Fprintln(w, "-------\t--------\t----------\t")

		// Sort channels
		var channels24 []int
		for ch := range usage24 {
			channels24 = append(channels24, ch)
		}
		sort.Ints(channels24)

		for _, ch := range channels24 {
			count := usage24[ch]
			congestion := "Low"
			if count >= 4 {
				congestion = "Very High"
			} else if count >= 3 {
				congestion = "High"
			} else if count >= 2 {
				congestion = "Medium"
			}
			fmt.Fprintf(w, "%d\t%d\t%s\t\n", ch, count, congestion)
		}
	}

	// 5GHz usage
	if len(usage5) > 0 {
		fmt.Fprintln(w, "\n5GHz Channel Usage:")
		fmt.Fprintln(w, "Channel\tNetworks\tCongestion\t")
		fmt.Fprintln(w, "-------\t--------\t----------\t")

		// Sort channels
		var channels5 []int
		for ch := range usage5 {
			channels5 = append(channels5, ch)
		}
		sort.Ints(channels5)

		for _, ch := range channels5 {
			count := usage5[ch]
			congestion := "Low"
			if count >= 3 {
				congestion = "High"
			} else if count >= 2 {
				congestion = "Medium"
			}
			fmt.Fprintf(w, "%d\t%d\t%s\t\n", ch, count, congestion)
		}
	}

	w.Flush()
}
