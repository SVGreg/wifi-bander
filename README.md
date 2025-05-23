# WiFi Bander

A professional-grade cross-platform Go application that provides comprehensive WiFi network analysis, intelligent channel recommendations, and spectrum optimization for better connectivity.

## 🚀 Features

### **Core Capabilities**
- **Cross-platform support** (Linux and macOS)
- **Comprehensive WiFi analysis** - Both 2.4GHz and 5GHz bands
- **Dynamic channel detection** - Automatically discovers all available channels in your region
- **Professional-grade recommendations** - AI-powered channel optimization
- **Real-time monitoring** - Continuous updates every 10 seconds
- **Production-ready** - Clean interface, no debug output

### **Advanced Analysis**
- **Frequency interference analysis** - Considers actual MHz separation between networks
- **Signal strength correlation** - Weighs network proximity and signal power
- **Channel width awareness** - Accounts for 20MHz, 40MHz, 80MHz channel usage
- **Vendor identification** - Recognizes equipment manufacturers from MAC addresses
- **Security protocol detection** - Identifies WPA, WPA2, WPA3, Open networks
- **PHY mode analysis** - Shows WiFi standards (802.11a/n/ac/ax, 802.11b/g/n)

### **Intelligent Recommendations**
- **Top 3 ranked channel suggestions** per band with detailed reasoning
- **Frequency separation optimization** - Maximizes distance from interfering signals
- **Non-overlapping channel preference** - Prioritizes optimal 2.4GHz channels (1, 6, 11)
- **DFS channel considerations** - Balances availability vs. radar detection requirements
- **Configuration guidance** - Actionable tips for optimal router setup

## How It Works

### **Sophisticated Analysis Engine**

Instead of simple user counting, WiFi Bander employs advanced RF analysis:

1. **🔍 Comprehensive Spectrum Scanning**
   - Detects all networks across 2.4GHz (1-13) and 5GHz (36-177) channels
   - Identifies security protocols, PHY modes, and channel widths
   - Maps vendor equipment using MAC address OUI lookup

2. **📊 Interference Correlation Analysis**
   - Calculates frequency separation between networks (MHz-level precision)
   - Weighs signal strength impact on neighboring channels
   - Considers channel width overlap (especially important for 80MHz 5GHz)

3. **🎯 AI-Powered Channel Scoring**
   - **2.4GHz**: Prioritizes non-overlapping channels, penalizes interference
   - **5GHz**: Maximizes frequency gaps, considers DFS vs. non-DFS availability
   - Provides ranked recommendations with confidence scoring

4. **💡 Actionable Intelligence**
   - Explains reasoning behind each recommendation
   - Shows frequency separation between optimal choices
   - Provides configuration tips and monitoring advice

## Comprehensive Network Information

### **Enhanced Network Analysis Table**
```
SSID             Band Ch  Signal  Quality Security           PHY Mode        Width Vendor  Congestion Freq 
----             ---- --  ------  ------- --------           --------        ----- ------  ---------- ---- 
NeroDiablo 5G    5G   36  -80 dBm 16%     WPA2 Personal      802.11a/n/ac    80MHz Apple   High       5180 
GUASH-2_5G       5G   36  -88 dBm 3%      WPA2 Personal      802.11a/n/ac/ax 80MHz         High       5180 
Sonja_guar       5G   100 -87 dBm 5%      WPA2/WPA3 Personal 802.11a/n/ac/ax 80MHz         High       5500 
SmartCar         2.4G 7   -70 dBm 33%     None               802.11b/g/n     20MHz         High       2442 
SuperMario       2.4G 4   -60 dBm 50%     WPA/WPA2 Personal  802.11b/g/n/ac  20MHz TP-Link Very High  2427 

Total networks detected: 19
```

### **Horizontal Channel Usage Matrix**
```
=== Channel Usage Statistics ===

2.4GHz Channel Usage:
Channel  1 2 3 4 5 6 7 8 9 10 11 12 13 
-------  - - - - - - - - - -  -  -  -  
Networks 2 0 0 1 0 1 3 0 1 0  0  1  0  

5GHz Channel Usage:
Channel  36 40 44 48 52 56 60 64 100 104 108 112 116 120 124 128 132 136 140 144 149 153 157 161 165 169 173 177
-------  -- -- -- -- -- -- -- -- --  --  --  --  --  --  --  --  --  --  --  --  --  --  --  --  --  --  --  --
Networks 1  1  0  4  0  0  0  0  2   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
```

### **Advanced Channel Recommendations**
```
=== Channel Recommendations (Top 3 Optimal Choices) ===
Advanced analysis considering frequency separation, signal strength, and interference patterns

🔸 2.4G Band Recommendations:
Rank  Channel  Freq(MHz)  Interference  Gap(MHz)  Reasoning                                             
----  -------  ---------  -----------   --------  ---------                                             
#1    11       2462       Low           5         Optimal: Non-overlapping channel with no detected networks
#2    13       2472       Low           5         Good: No networks detected, minimal interference expected  
#3    2        2417       Moderate      15        Good: No networks detected, minimal interference expected

  📊 Frequency Separation Analysis:
     • Channel 11 ↔ Channel 13: 10 MHz separation
     • Channel 13 ↔ Channel 2: 55 MHz separation

  💡 2.4GHz Advice: Prefer channels 1, 6, or 11 (non-overlapping). Avoid channels with strong nearby signals.

🔸 5G Band Recommendations:
Rank  Channel  Freq(MHz)  Interference  Gap(MHz)  Reasoning                                             
----  -------  ---------  -----------   --------  ---------                                             
#1    149      5745       Minimal       245       Excellent: Non-DFS channel with no detected networks  
#2    177      5885       Minimal       385       Excellent: Non-DFS channel with no detected networks  
#3    173      5865       Minimal       365       Excellent: Non-DFS channel with no detected networks  

  📊 Frequency Separation Analysis:
     • Channel 149 ↔ Channel 177: 140 MHz separation
     • Channel 177 ↔ Channel 173: 20 MHz separation

  💡 5GHz Advice: More spectrum available. DFS channels may require radar detection but are often less congested.

🎯 Configuration Tips:
   • Choose the #1 ranked channel for optimal performance
   • Monitor performance and try #2 or #3 if issues occur
   • Consider channel width: 80MHz for 5GHz, 20MHz for 2.4GHz in crowded areas
   • Update analysis periodically as WiFi landscape changes
```

## Supported Channel Ranges & Standards

### **2.4GHz Band (Dynamic Detection)**
- **Channels 1-11**: US standard
- **Channels 1-13**: European standard  
- **Channels 1-14**: Japanese standard (includes channel 14)
- **Optimal non-overlapping**: 1, 6, 11
- **Auto-detects regional availability**

### **5GHz Band (Comprehensive UNII Support)**
- **UNII-1 (36-48)**: 5.15-5.25 GHz, indoor use, no DFS
- **UNII-2A (52-64)**: 5.25-5.35 GHz, DFS required
- **UNII-2C (100-144)**: 5.47-5.725 GHz, DFS required  
- **UNII-3 (149-165)**: 5.725-5.875 GHz, outdoor use, no DFS
- **UNII-4 (169-177)**: 5.85-5.925 GHz, newer allocation

## Installation

### **Option 1: From Source (Recommended)**
```bash
git clone https://github.com/svgreg/wifi-bander.git
cd wifi-bander
make build
```

### **Option 2: Direct Go Install**
```bash
go install github.com/svgreg/wifi-bander@latest
```

### **Option 3: Manual Build**
```bash
git clone https://github.com/svgreg/wifi-bander.git
cd wifi-bander
go mod tidy
go build
```

## Usage

### **macOS (Recommended)**
```bash
./wifi-bander
```
*Uses system_profiler - no special permissions required*

### **Linux**
```bash
# Try without sudo first
./wifi-bander

# If permission errors occur:
sudo ./wifi-bander
```
*Requires nmcli (NetworkManager) or iwlist (wireless-tools)*

## Requirements

### **System Requirements**
- **macOS**: 10.13+ with built-in WiFi support
- **Linux**: WiFi interface with nl80211 support
- **Go**: Version 1.21 or later for building from source

### **Dependencies**
- **macOS**: `system_profiler` (built-in)
- **Linux**: `nmcli` (NetworkManager) or `iwlist` (wireless-tools)

## Architecture & Implementation

### **Modular Design**
```
wifi-bander/
├── main.go                           # Application entry point
├── go.mod                            # Module: github.com/svgreg/wifi-bander
├── Makefile                          # Professional build system
├── README.md                         # This documentation
└── internal/                         # Internal packages
    ├── scanner/                      # Platform-specific WiFi scanning
    │   ├── types.go                 # Enhanced network data structures
    │   ├── scanner.go               # Cross-platform scanner interface
    │   ├── macos.go                 # macOS system_profiler implementation  
    │   └── linux.go                 # Linux nmcli/iwlist implementation
    ├── analyzer/                     # Advanced analysis algorithms
    │   └── analyzer.go              # AI-powered channel optimization
    └── display/                      # Professional output formatting
        └── display.go               # Tables, recommendations, statistics
```

### **Platform-Specific Implementation**

#### **macOS: system_profiler Integration**
- **Primary method**: `system_profiler SPAirPortDataType`
- **Extracts**: SSID, Channel, Signal/Noise, Security, PHY Mode, BSSID
- **Advantages**: No root required, comprehensive data, built-in tool
- **Parsing**: Advanced property extraction from structured output

#### **Linux: Multi-tool Approach**  
- **Primary**: `nmcli` (NetworkManager) for modern systems
- **Fallback**: `iwlist` for legacy/minimal installations
- **Extracts**: SSID, Channel, Signal, Frequency, Security, BSSID
- **Requirements**: Standard Linux wireless tools

### **Advanced Scoring Algorithm**

#### **2.4GHz Optimization**
```
Score = Base_Penalty + Same_Channel_Penalty + Adjacent_Interference + Signal_Impact

- Non-overlapping channels (1,6,11): Base penalty = 0
- Overlapping channels (2-5,7-10,12-13): Base penalty = +20
- Same channel networks: +50 per network
- Adjacent interference: Decreases with frequency distance
- Strong signals (-40 to -60 dBm): Additional +10 to +20 penalty
```

#### **5GHz Optimization**
```
Score = DFS_Penalty + Same_Channel_Penalty + Bandwidth_Interference + Signal_Impact

- Non-DFS channels: Base penalty = 0  
- DFS channels: Base penalty = +10
- 80MHz channel width consideration
- Frequency separation optimization (up to 700+ MHz available)
- Lower interference weighting than 2.4GHz (less prone to interference)
```

## Understanding the Output

### **Network Analysis Fields**
- **SSID**: Network name (truncated to 16 chars for display)
- **Band**: 2.4G or 5G frequency band
- **Ch**: Channel number
- **Signal**: Signal strength in dBm (-30 excellent, -90 very weak)
- **Quality**: Signal quality percentage (0-100%)
- **Security**: Full security protocol (WPA2 Personal, WPA/WPA2, None, etc.)
- **PHY Mode**: Complete WiFi standard (802.11a/n/ac/ax, 802.11b/g/n/ac)
- **Width**: Channel width (20MHz, 40MHz, 80MHz, 160MHz)
- **Vendor**: Equipment manufacturer (Apple, TP-Link, ASUS, etc.)
- **Congestion**: Interference level (Low/Medium/High/Very High)
- **Freq**: Exact frequency in MHz

### **Channel Usage Statistics**
- **Horizontal layout**: Quick visual spectrum overview
- **All channels shown**: Including empty channels (0 networks)
- **Complete coverage**: 2.4GHz (1-13), 5GHz (36-177)
- **Pattern recognition**: Spot clustering and gaps instantly

### **Recommendation Confidence Levels**
- **Minimal** (0-20): Optimal choice, minimal interference expected
- **Low** (21-50): Good choice, slight interference possible  
- **Moderate** (51-100): Acceptable choice, some interference likely
- **High** (101-200): Suboptimal choice, significant interference
- **Very High** (>200): Poor choice, heavy interference expected

## Professional Use Cases

### **Network Deployment Planning**
- **Site surveys**: Identify optimal channels before router installation
- **Interference analysis**: Understand existing network landscape
- **Capacity planning**: Find channels with room for additional networks

### **Performance Troubleshooting**
- **Congestion identification**: Pinpoint overcrowded channels
- **Interference source mapping**: Locate problematic networks
- **Optimization recommendations**: Get specific channel suggestions

### **Ongoing Network Management**  
- **Regular monitoring**: Track changes in WiFi environment
- **Proactive optimization**: Adjust channels before performance degrades
- **Documentation**: Maintain records of network landscape changes

## Development & Extension

### **Build Commands**
```bash
make build          # Build application
make test           # Run tests  
make lint           # Run linting
make clean          # Clean build artifacts
make cross-compile  # Build for multiple platforms
make help           # Show all available commands
```

### **Adding New Platforms**
1. Create `internal/scanner/{platform}.go`
2. Implement `Scanner` interface
3. Add platform detection in `scanner.go`
4. Test with various network configurations

### **Extending Analysis**
- **Historical tracking**: Add time-series analysis
- **Machine learning**: Implement predictive interference modeling
- **Web interface**: Create REST API endpoints
- **Database integration**: Store long-term scanning data

## Troubleshooting

### **macOS Issues**
- **Empty results**: Ensure WiFi is enabled and scanning allowed
- **Partial data**: Some fields may be "Unknown" depending on system version
- **Performance**: macOS scanning is generally faster and more reliable

### **Linux Issues**
- **nmcli not found**: `sudo apt install network-manager`
- **iwlist not found**: `sudo apt install wireless-tools`  
- **Permission denied**: Run with `sudo` or add user to appropriate groups
- **No interface**: Check `ip link` or `iwconfig` for WiFi adapters

### **General Issues**
- **No networks**: Wait 15-20 seconds for initial scan completion
- **Inconsistent results**: WiFi landscape changes frequently (normal behavior)
- **High congestion everywhere**: Consider 5GHz migration or location change

## Performance & Resource Usage

- **Scan frequency**: Every 10 seconds (configurable)
- **Memory footprint**: < 15MB typical usage
- **CPU impact**: Minimal, brief spikes during scans only
- **Network impact**: Read-only passive scanning
- **Platform optimization**: Tailored commands for each OS

## Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature-name`
3. Follow established patterns and interfaces
4. Add comprehensive error handling  
5. Test on multiple platforms when possible
6. Run linting: `make lint`
7. Submit pull request

### **Code Guidelines**
- **Go conventions**: Follow `gofmt` and standard practices
- **Interface design**: Keep interfaces small and focused
- **Error handling**: Comprehensive error reporting
- **Documentation**: Document all public APIs
- **Testing**: Add tests for new functionality

## License

MIT License - see LICENSE file for details

## Changelog

### **v3.0.0 (Current - Professional-Grade Analysis)**
- ✅ **Comprehensive network analysis** - 11-field detailed table
- ✅ **Advanced channel recommendations** - AI-powered top 3 suggestions with reasoning
- ✅ **Horizontal channel usage** - Matrix view of entire spectrum
- ✅ **Enhanced network information** - Security, PHY mode, vendor, channel width
- ✅ **Frequency-based interference** - MHz-level precision analysis
- ✅ **Professional output** - Production-ready formatting and guidance
- ✅ **Signal quality calculation** - Percentage-based quality metrics
- ✅ **Vendor identification** - MAC address OUI lookup for major manufacturers

### **v2.1.0 (Previous - Dynamic Channels)**
- ✅ **Dynamic channel detection** - No hardcoded limitations
- ✅ **Comprehensive 5GHz support** - All UNII bands including DFS
- ✅ **Regional awareness** - US/EU/Japan standards
- ✅ **Smart fallback recommendations** - Adaptive to local regulations

### **v2.0.0 (Previous - Refactored)**
- ✅ **Modular architecture** - Clean package separation
- ✅ **Interface-based design** - Extensible and maintainable
- ✅ **Enhanced build system** - Professional Makefile

### **v1.1.0 (Previous - macOS Fix)**
- ✅ **macOS compatibility** - system_profiler integration
- ✅ **Improved parsing** - Reliable data extraction
- ✅ **Production output** - Clean user interface

### **v1.0.0 (Initial)**
- ✅ **Cross-platform scanning** - Linux and macOS support
- ✅ **Basic recommendations** - Simple congestion-based suggestions
- ✅ **Channel analysis** - Fundamental interference detection 