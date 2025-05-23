package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/svgreg/wifi-bander/internal/analyzer"
	"github.com/svgreg/wifi-bander/internal/display"
	"github.com/svgreg/wifi-bander/internal/scanner"
)

func main() {
	fmt.Println("WiFi Bander - Cross-Platform WiFi Network Analyzer")
	fmt.Println("Initializing scanner...")

	// Test the scanner once before starting the loop
	networks, err := scanner.ScanWiFiNetworks()
	if err != nil {
		log.Fatalf("Initial scan failed: %v\nPlease ensure you have the correct permissions to scan WiFi networks.", err)
	}

	fmt.Println("Scanner initialized successfully.")

	// Show detailed channel information on first run
	if len(networks) > 0 {
		// Convert to analyzer interface for channel analysis
		analyzerNetworks := make([]analyzer.WiFiNetwork, len(networks))
		for i, net := range networks {
			analyzerNetworks[i] = net
		}

		display.DisplayChannelInfo(analyzerNetworks)
		fmt.Println("\nStarting continuous scan...")
		time.Sleep(3 * time.Second) // Give user time to read
	} else {
		fmt.Println("No networks detected in initial scan. Starting continuous scan...")
	}

	for {
		networks, err := scanner.ScanWiFiNetworks()
		if err != nil {
			log.Printf("Error scanning networks: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Sort networks by congestion score (ascending - least congested first)
		sort.Slice(networks, func(i, j int) bool {
			return networks[i].CongestionScore < networks[j].CongestionScore
		})

		// Convert to display interface
		displayNetworks := make([]display.WiFiNetwork, len(networks))
		for i, net := range networks {
			displayNetworks[i] = net
		}

		// Convert to analyzer interface
		analyzerNetworks := make([]analyzer.WiFiNetwork, len(networks))
		for i, net := range networks {
			analyzerNetworks[i] = net
		}

		display.DisplayResults(displayNetworks)
		display.DisplayRecommendations(analyzerNetworks)

		time.Sleep(10 * time.Second)
	}
}
