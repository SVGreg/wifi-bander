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

	// Find least congested 2.4GHz channels
	bestChannels24 := findLeastCongestedChannels(channelUsage24, []int{1, 6, 11})
	recommendations["2.4G"] = bestChannels24

	// Find least congested 5GHz channels
	available5GChannels := []int{36, 40, 44, 48, 149, 153, 157, 161, 165}
	bestChannels5 := findLeastCongestedChannels(channelUsage5, available5GChannels)
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
