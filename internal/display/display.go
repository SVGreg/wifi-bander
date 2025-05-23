package display

import (
	"fmt"
	"os"
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
