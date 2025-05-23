package analyzer

import (
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

// GetChannelRecommendations returns recommended channels based on congestion analysis
func GetChannelRecommendations(networks []WiFiNetwork) map[string][]int {
	channelUsage24 := make(map[int]int)
	channelUsage5 := make(map[int]int)

	// Count networks per channel
	for _, network := range networks {
		if network.GetBand() == "2.4G" {
			channelUsage24[network.GetChannel()]++
		} else {
			channelUsage5[network.GetChannel()]++
		}
	}

	recommendations := make(map[string][]int)

	// For 2.4GHz: Use detected channels, but prefer non-overlapping channels in recommendations
	detected24 := getDetectedChannels(networks, "2.4G")
	available24 := detected24

	// If no channels detected, fall back to comprehensive list
	if len(available24) == 0 {
		available24 = getAllAvailable24GHzChannels()
	}

	// For recommendations, prefer non-overlapping channels if they exist in available channels
	recommend24 := []int{}
	for _, ch := range Channels24GHz_NonOverlapping {
		// Check if this non-overlapping channel exists in our available channels
		for _, avail := range available24 {
			if ch == avail {
				recommend24 = append(recommend24, ch)
				break
			}
		}
	}

	// If we don't have enough non-overlapping channels, add others
	if len(recommend24) < 3 {
		for _, ch := range available24 {
			// Add channels not already in recommend24
			found := false
			for _, existing := range recommend24 {
				if ch == existing {
					found = true
					break
				}
			}
			if !found {
				recommend24 = append(recommend24, ch)
			}
		}
	}

	bestChannels24 := findLeastCongestedChannels(channelUsage24, recommend24)
	recommendations["2.4G"] = bestChannels24

	// For 5GHz: Use detected channels with comprehensive fallback
	detected5 := getDetectedChannels(networks, "5G")
	available5 := detected5

	// If no channels detected, fall back to comprehensive list
	if len(available5) == 0 {
		available5 = getAllAvailable5GHzChannels()
	}

	bestChannels5 := findLeastCongestedChannels(channelUsage5, available5)
	recommendations["5G"] = bestChannels5

	return recommendations
}

// findLeastCongestedChannels finds the least congested channels from available options
func findLeastCongestedChannels(usage map[int]int, availableChannels []int) []int {
	type channelScore struct {
		channel int
		count   int
	}

	var scores []channelScore
	for _, ch := range availableChannels {
		scores = append(scores, channelScore{
			channel: ch,
			count:   usage[ch],
		})
	}

	// Sort by usage (ascending)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].count < scores[j].count
	})

	// Return top 3 least congested channels
	result := make([]int, 0, 3)
	for i := 0; i < len(scores) && i < 3; i++ {
		result = append(result, scores[i].channel)
	}

	return result
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
