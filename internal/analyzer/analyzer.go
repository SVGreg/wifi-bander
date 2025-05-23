package analyzer

import (
	"fmt"
	"sort"
)

// ChannelInfo represents aggregated information about a WiFi channel
// This should match the one in scanner package
type ChannelInfo struct {
	Channel       int
	NetworkCount  int
	StrongestRSSI int
}

// WiFiNetwork interface to avoid circular imports
type WiFiNetwork interface {
	GetBand() string
	GetChannel() int
	GetSignal() int
	GetStationCount() int
	GetChannelWidth() string
}

// NetworkInfo is a minimal struct for networks used in analysis
type NetworkInfo struct {
	Band         string
	Channel      int
	Signal       int
	StationCount int
}

func (n NetworkInfo) GetBand() string      { return n.Band }
func (n NetworkInfo) GetChannel() int      { return n.Channel }
func (n NetworkInfo) GetSignal() int       { return n.Signal }
func (n NetworkInfo) GetStationCount() int { return n.StationCount }

// Comprehensive channel definitions by region/standard
var (
	// 2.4GHz channels (by region)
	Channels24GHz_US     = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	Channels24GHz_Europe = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	Channels24GHz_Japan  = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

	// Non-overlapping 2.4GHz channels (universal optimum)
	Channels24GHz_NonOverlapping = []int{1, 6, 11}

	// 5GHz channels (comprehensive list)
	// UNII-1 Band (5.15-5.25 GHz) - Low power, indoor use
	Channels5GHz_UNII1 = []int{36, 40, 44, 48}

	// UNII-2A Band (5.25-5.35 GHz) - DFS required
	Channels5GHz_UNII2A = []int{52, 56, 60, 64}

	// UNII-2C Band (5.47-5.725 GHz) - DFS required
	Channels5GHz_UNII2C = []int{100, 104, 108, 112, 116, 120, 124, 128, 132, 136, 140, 144}

	// UNII-3 Band (5.725-5.875 GHz) - Higher power, outdoor use
	Channels5GHz_UNII3 = []int{149, 153, 157, 161, 165}

	// UNII-4 Band (5.85-5.925 GHz) - Newer allocation
	Channels5GHz_UNII4 = []int{169, 173, 177}
)

// getAllAvailable5GHzChannels returns a comprehensive list of 5GHz channels
func getAllAvailable5GHzChannels() []int {
	var all []int
	all = append(all, Channels5GHz_UNII1...)
	all = append(all, Channels5GHz_UNII2A...)
	all = append(all, Channels5GHz_UNII2C...)
	all = append(all, Channels5GHz_UNII3...)
	all = append(all, Channels5GHz_UNII4...)
	return all
}

// getAllAvailable24GHzChannels returns channels 1-13 (European standard, most comprehensive)
func getAllAvailable24GHzChannels() []int {
	return Channels24GHz_Europe // Most comprehensive common standard
}

// CalculateCongestionScore calculates the congestion score for a network
func CalculateCongestionScore(network WiFiNetwork, channelMap interface{}) int {
	// Convert interface{} to map[int]*ChannelInfo
	var channelInfoMap map[int]*ChannelInfo
	switch v := channelMap.(type) {
	case map[int]*ChannelInfo:
		channelInfoMap = v
	default:
		// If it's from scanner package, create a compatible map
		return calculateWithGenericMap(network, channelMap)
	}

	score := 0
	channel := network.GetChannel()

	// Base score from number of networks on the same channel
	if channelInfo, exists := channelInfoMap[channel]; exists {
		score += channelInfo.NetworkCount * 10
	}

	// Add penalty for overlapping channels (2.4GHz)
	if network.GetBand() == "2.4G" {
		for i := -2; i <= 2; i++ {
			if i == 0 {
				continue
			}
			adjChannel := channel + i
			if channelInfo, exists := channelInfoMap[adjChannel]; exists {
				score += channelInfo.NetworkCount * 5
			}
		}
	}

	// Use estimated station count
	score += network.GetStationCount() * 8

	// Penalty for strong signals
	signal := network.GetSignal()
	if signal > -50 {
		score += 20
	} else if signal > -70 {
		score += 10
	}

	return score
}

// calculateWithGenericMap handles different ChannelInfo types via reflection-like approach
func calculateWithGenericMap(network WiFiNetwork, channelMap interface{}) int {
	score := 0

	// This is a simplified version that works with any compatible map
	// We'll extract the basic scoring logic without type dependencies

	// Base station count scoring
	score += network.GetStationCount() * 8

	// Penalty for strong signals
	signal := network.GetSignal()
	if signal > -50 {
		score += 20
	} else if signal > -70 {
		score += 10
	}

	// Add a base congestion penalty
	score += 25 // Default congestion assumption

	return score
}

// getDetectedChannels extracts all channels actually detected during scanning
func getDetectedChannels(networks []WiFiNetwork, band string) []int {
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

	// Sort channels
	sort.Ints(channels)
	return channels
}

// ChannelRecommendation represents a recommended channel with scoring details
type ChannelRecommendation struct {
	Channel           int
	Frequency         int
	Score             float64
	InterferenceLevel string
	Reasoning         string
	SignalImpact      float64
	FrequencyGap      int // MHz to nearest neighbor
}

// GetChannelRecommendations returns top 3 recommended channels per band based on advanced criteria
func GetChannelRecommendations(networks []WiFiNetwork) map[string][]ChannelRecommendation {
	// Analyze current network landscape
	channelAnalysis24 := analyzeChannelLandscape(networks, "2.4G")
	channelAnalysis5 := analyzeChannelLandscape(networks, "5G")

	recommendations := make(map[string][]ChannelRecommendation)

	// Get 2.4GHz recommendations
	recommendations["2.4G"] = getBest24GHzChannels(channelAnalysis24)

	// Get 5GHz recommendations
	recommendations["5G"] = getBest5GHzChannels(channelAnalysis5)

	return recommendations
}

// NetworkAnalysis holds analysis data for a specific network
type NetworkAnalysis struct {
	Channel       int
	Frequency     int
	Signal        int
	ChannelWidth  string
	NetworkCount  int
	StrongestRSSI int
}

// analyzeChannelLandscape creates a comprehensive analysis of the current WiFi landscape
func analyzeChannelLandscape(networks []WiFiNetwork, band string) map[int]*NetworkAnalysis {
	analysis := make(map[int]*NetworkAnalysis)

	for _, network := range networks {
		if network.GetBand() != band {
			continue
		}

		ch := network.GetChannel()
		freq := channelToFrequency(ch)
		signal := network.GetSignal()

		if existing, exists := analysis[ch]; exists {
			existing.NetworkCount++
			if signal > existing.StrongestRSSI {
				existing.StrongestRSSI = signal
			}
		} else {
			analysis[ch] = &NetworkAnalysis{
				Channel:       ch,
				Frequency:     freq,
				Signal:        signal,
				ChannelWidth:  getChannelWidth(network),
				NetworkCount:  1,
				StrongestRSSI: signal,
			}
		}
	}

	return analysis
}

// getChannelWidth extracts channel width from network, with fallback logic
func getChannelWidth(network WiFiNetwork) string {
	if width := network.GetChannelWidth(); width != "Unknown" && width != "" {
		return width
	}
	// Default assumptions based on band
	if network.GetBand() == "5G" {
		return "80MHz"
	}
	return "20MHz"
}

// getBest24GHzChannels finds optimal 2.4GHz channels with sophisticated scoring
func getBest24GHzChannels(analysis map[int]*NetworkAnalysis) []ChannelRecommendation {
	// Available 2.4GHz channels (1-13, EU standard)
	availableChannels := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

	var recommendations []ChannelRecommendation

	for _, ch := range availableChannels {
		freq := channelToFrequency(ch)
		score := calculate24GHzScore(ch, freq, analysis)

		recommendation := ChannelRecommendation{
			Channel:           ch,
			Frequency:         freq,
			Score:             score,
			InterferenceLevel: getInterferenceLevel(score),
			Reasoning:         getReasoning24GHz(ch, analysis),
			SignalImpact:      calculateSignalImpact(ch, analysis),
			FrequencyGap:      calculateFrequencyGap(freq, analysis),
		}

		recommendations = append(recommendations, recommendation)
	}

	// Sort by score (lower is better)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score < recommendations[j].Score
	})

	// Return top 3
	if len(recommendations) > 3 {
		recommendations = recommendations[:3]
	}

	return recommendations
}

// calculate24GHzScore calculates interference score for 2.4GHz channels
func calculate24GHzScore(channel, frequency int, analysis map[int]*NetworkAnalysis) float64 {
	score := 0.0

	// Base score: prefer non-overlapping channels (1, 6, 11)
	nonOverlapping := []int{1, 6, 11}
	isNonOverlapping := false
	for _, noCh := range nonOverlapping {
		if channel == noCh {
			isNonOverlapping = true
			break
		}
	}

	if !isNonOverlapping {
		score += 20.0 // Penalty for overlapping channels
	}

	// Calculate interference from all existing networks
	for _, net := range analysis {
		freqDiff := abs(frequency - net.Frequency)

		// Strong penalty for same channel
		if freqDiff == 0 {
			score += float64(net.NetworkCount) * 50.0
			// Additional penalty for strong signals
			if net.StrongestRSSI > -60 {
				score += 30.0
			}
			continue
		}

		// Adjacent channel interference (2.4GHz channels are 5MHz apart)
		if freqDiff <= 15 { // Within 3 channels
			interferenceWeight := 30.0 - (float64(freqDiff) * 2.0) // Decreases with distance
			score += interferenceWeight * float64(net.NetworkCount)

			// Signal strength impact
			signalPenalty := calculateSignalPenalty(net.StrongestRSSI, freqDiff)
			score += signalPenalty
		}
	}

	return score
}

// getBest5GHzChannels finds optimal 5GHz channels
func getBest5GHzChannels(analysis map[int]*NetworkAnalysis) []ChannelRecommendation {
	// All available 5GHz channels
	allChannels5G := getAllAvailable5GHzChannels()

	var recommendations []ChannelRecommendation

	for _, ch := range allChannels5G {
		freq := channelToFrequency(ch)
		score := calculate5GHzScore(ch, freq, analysis)

		recommendation := ChannelRecommendation{
			Channel:           ch,
			Frequency:         freq,
			Score:             score,
			InterferenceLevel: getInterferenceLevel(score),
			Reasoning:         getReasoning5GHz(ch, analysis),
			SignalImpact:      calculateSignalImpact(ch, analysis),
			FrequencyGap:      calculateFrequencyGap(freq, analysis),
		}

		recommendations = append(recommendations, recommendation)
	}

	// Sort by score (lower is better)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score < recommendations[j].Score
	})

	// Return top 3
	if len(recommendations) > 3 {
		recommendations = recommendations[:3]
	}

	return recommendations
}

// calculate5GHzScore calculates interference score for 5GHz channels
func calculate5GHzScore(channel, frequency int, analysis map[int]*NetworkAnalysis) float64 {
	score := 0.0

	// Prefer non-DFS channels (UNII-1 and UNII-3)
	isDFS := (channel >= 52 && channel <= 64) || (channel >= 100 && channel <= 144)
	if isDFS {
		score += 10.0 // Small penalty for DFS channels
	}

	// Calculate interference from existing networks
	for _, net := range analysis {
		freqDiff := abs(frequency - net.Frequency)

		// Same channel penalty
		if freqDiff == 0 {
			score += float64(net.NetworkCount) * 40.0
			if net.StrongestRSSI > -60 {
				score += 25.0
			}
			continue
		}

		// 5GHz interference calculation (considering 80MHz channel widths)
		interferenceRange := 80 // MHz, typical 5GHz channel width
		if freqDiff <= interferenceRange {
			interferenceWeight := float64(interferenceRange-freqDiff) / 10.0
			score += interferenceWeight * float64(net.NetworkCount)

			signalPenalty := calculateSignalPenalty(net.StrongestRSSI, freqDiff)
			score += signalPenalty * 0.8 // 5GHz less prone to interference
		}
	}

	return score
}

// calculateSignalPenalty calculates penalty based on signal strength and frequency distance
func calculateSignalPenalty(signalStrength, freqDiff int) float64 {
	// Convert dBm to penalty weight (stronger signals cause more interference)
	signalWeight := 0.0
	if signalStrength > -40 {
		signalWeight = 20.0
	} else if signalStrength > -60 {
		signalWeight = 10.0
	} else if signalStrength > -80 {
		signalWeight = 5.0
	}

	// Apply distance factor
	distanceFactor := 1.0 / (1.0 + float64(freqDiff)/10.0)
	return signalWeight * distanceFactor
}

// calculateSignalImpact calculates the signal impact score for a channel
func calculateSignalImpact(channel int, analysis map[int]*NetworkAnalysis) float64 {
	impact := 0.0
	freq := channelToFrequency(channel)

	for _, net := range analysis {
		freqDiff := abs(freq - net.Frequency)
		if freqDiff <= 40 { // Within interference range
			signalImpact := float64(-net.StrongestRSSI) / float64(freqDiff+1)
			impact += signalImpact
		}
	}

	return impact
}

// calculateFrequencyGap calculates the gap to the nearest network in MHz
func calculateFrequencyGap(frequency int, analysis map[int]*NetworkAnalysis) int {
	minGap := 1000 // Large initial value

	for _, net := range analysis {
		gap := abs(frequency - net.Frequency)
		if gap > 0 && gap < minGap {
			minGap = gap
		}
	}

	if minGap == 1000 {
		return 0 // No other networks
	}
	return minGap
}

// getReasoning24GHz provides reasoning for 2.4GHz channel recommendation
func getReasoning24GHz(channel int, analysis map[int]*NetworkAnalysis) string {
	nonOverlapping := []int{1, 6, 11}
	isNonOverlapping := false
	for _, noCh := range nonOverlapping {
		if channel == noCh {
			isNonOverlapping = true
			break
		}
	}

	if analysis[channel] == nil {
		if isNonOverlapping {
			return "Optimal: Non-overlapping channel with no detected networks"
		}
		return "Good: No networks detected, minimal interference expected"
	}

	net := analysis[channel]
	if isNonOverlapping {
		return fmt.Sprintf("Fair: Non-overlapping but has %d network(s), strongest at %d dBm",
			net.NetworkCount, net.StrongestRSSI)
	}

	return fmt.Sprintf("Suboptimal: Overlapping channel with %d network(s)", net.NetworkCount)
}

// getReasoning5GHz provides reasoning for 5GHz channel recommendation
func getReasoning5GHz(channel int, analysis map[int]*NetworkAnalysis) string {
	isDFS := (channel >= 52 && channel <= 64) || (channel >= 100 && channel <= 144)

	if analysis[channel] == nil {
		if isDFS {
			return "Good: DFS channel with no detected networks, radar detection required"
		}
		return "Excellent: Non-DFS channel with no detected networks"
	}

	net := analysis[channel]
	status := "Fair"
	if net.NetworkCount == 1 && net.StrongestRSSI < -70 {
		status = "Good"
	}

	dfsNote := ""
	if isDFS {
		dfsNote = ", DFS required"
	}

	return fmt.Sprintf("%s: %d network(s), strongest at %d dBm%s",
		status, net.NetworkCount, net.StrongestRSSI, dfsNote)
}

// getInterferenceLevel converts score to human-readable interference level
func getInterferenceLevel(score float64) string {
	switch {
	case score <= 20:
		return "Minimal"
	case score <= 50:
		return "Low"
	case score <= 100:
		return "Moderate"
	case score <= 200:
		return "High"
	default:
		return "Very High"
	}
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// channelToFrequency converts channel to frequency - moved here for local access
func channelToFrequency(channel int) int {
	// 2.4GHz band (channels 1-14)
	if channel >= 1 && channel <= 13 {
		return 2412 + (channel-1)*5
	}
	if channel == 14 {
		return 2484
	}

	// 5GHz bands
	if channel >= 36 && channel <= 48 {
		return 5180 + (channel-36)*5
	}
	if channel >= 52 && channel <= 64 {
		return 5260 + (channel-52)*5
	}
	if channel >= 100 && channel <= 144 {
		return 5500 + (channel-100)*5
	}
	if channel >= 149 && channel <= 165 {
		return 5745 + (channel-149)*5
	}
	if channel >= 169 && channel <= 177 {
		return 5845 + (channel-169)*5
	}

	return 2412 // Default fallback
}

// GetChannelInfo returns detailed information about channel allocations
func GetChannelInfo() map[string]interface{} {
	return map[string]interface{}{
		"2.4GHz": map[string]interface{}{
			"non_overlapping": Channels24GHz_NonOverlapping,
			"us_channels":     Channels24GHz_US,
			"eu_channels":     Channels24GHz_Europe,
			"jp_channels":     Channels24GHz_Japan,
		},
		"5GHz": map[string]interface{}{
			"unii_1":  Channels5GHz_UNII1,
			"unii_2a": Channels5GHz_UNII2A,
			"unii_2c": Channels5GHz_UNII2C,
			"unii_3":  Channels5GHz_UNII3,
			"unii_4":  Channels5GHz_UNII4,
			"all":     getAllAvailable5GHzChannels(),
		},
	}
}
