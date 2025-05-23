# WiFi Bander

A cross-platform Go application that analyzes WiFi network congestion and helps identify optimal channels for better connectivity.

## Features

- **Cross-platform support** (Linux and macOS)
- Scans both 2.4GHz and 5GHz WiFi bands
- **Channel congestion analysis** instead of user counting (more reliable)
- Displays signal strength and frequency information
- **Smart channel recommendations** based on interference analysis
- Real-time updates every 10 seconds
- Clean table-formatted output
- **Production-ready** - no debug output, clean interface
- **Modular architecture** - well-organized, maintainable codebase

## How It Works

Instead of trying to count actual users (which requires specialized hardware), this application focuses on **channel congestion analysis**:

1. **Scans available WiFi networks** using OS-specific system commands
2. **Analyzes channel overlap** (especially important for 2.4GHz)
3. **Calculates congestion scores** based on:
   - Number of networks on the same channel
   - Adjacent channel interference (2.4GHz)
   - Signal strength indicators
   - Estimated station counts
4. **Provides channel recommendations** for optimal performance

## Architecture & Implementation

### **Project Structure**
```
wifi-bander/
├── main.go                          # Application entry point
├── go.mod                           # Module definition (github.com/svgreg/wifi-bander)
├── go.sum                           # Dependency checksums
├── Makefile                         # Build and development tasks
├── README.md                        # This documentation
└── internal/                        # Internal packages (not importable)
    ├── scanner/                     # WiFi scanning functionality
    │   ├── types.go                # Data structures and interfaces
    │   ├── scanner.go              # Main scanner with OS detection
    │   ├── macos.go                # macOS-specific implementation
    │   └── linux.go                # Linux-specific implementation
    ├── analyzer/                   # Congestion analysis algorithms
    │   └── analyzer.go             # Scoring & recommendations
    └── display/                    # Output formatting
        └── display.go              # Table formatting & display
```

### **Package Responsibilities**

#### **🔍 `internal/scanner`**
- **OS Detection**: Automatically detects Linux vs macOS
- **Platform Abstraction**: Common `Scanner` interface for all platforms
- **Data Structures**: `WiFiNetwork` and `ChannelInfo` types
- **Utility Functions**: Channel-to-frequency conversion, station estimation

#### **📊 `internal/analyzer`** 
- **Congestion Scoring**: Advanced algorithm considering channel overlap
- **Channel Recommendations**: Identifies optimal channels per band
- **Interface-Based**: Works with any `WiFiNetwork` implementation

#### **🖥️ `internal/display`**
- **Table Formatting**: Clean tabwriter-based output
- **Recommendations Display**: Channel suggestion formatting
- **Congestion Levels**: Human-readable scoring (Low/Medium/High/Very High)

### **Platform Implementations**

#### **macOS Implementation**
- **Primary method**: `system_profiler SPAirPortDataType` 
- **Fallback**: Previously used `airport` command (now deprecated on newer macOS)
- **Why system_profiler**: The traditional `/usr/sbin/airport` command is no longer available on modern macOS versions
- **Data parsing**: Extracts network information from the "Other Local Wi-Fi Networks" section
- **No root required**: Works with standard user permissions

#### **Linux Implementation**  
- **Primary method**: `nmcli` (NetworkManager)
- **Fallback**: `iwlist` scanning
- **Requirements**: WiFi interface with nl80211 support
- **Root access**: May require `sudo` for some operations

### **Interface-Based Design**
```go
// Scanner interface enables platform extensibility
type Scanner interface {
    Scan() ([]WiFiNetwork, error)
}

// WiFiNetwork interface enables loose coupling
type WiFiNetwork interface {
    GetBand() string
    GetChannel() int
    GetSignal() int
    GetStationCount() int
    GetSSID() string
    GetCongestionScore() int
    GetFrequency() int
}
```

### **Congestion Calculation Algorithm**
1. **Base score**: Number of networks on same channel × 10
2. **Adjacent interference** (2.4GHz only): Nearby channels × 5  
3. **Station estimation**: Based on signal strength patterns × 8
4. **Signal strength penalty**: Strong signals indicate busy networks
5. **Final score categories**:
   - Low (≤15): Minimal interference
   - Medium (16-30): Some interference  
   - High (31-50): Significant interference
   - Very High (>50): Heavy interference

## Requirements

- **Linux**: WiFi interface with nl80211 support, `nmcli` or `iwlist`
- **macOS**: Built-in WiFi support (macOS 10.13+)
- **Go**: Version 1.21 or later
- **Permissions**: Standard user (macOS), may need `sudo` (Linux)

## Installation

### **Option 1: From Source**
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

### On Linux:
```bash
# Try without sudo first
./wifi-bander

# If permission errors occur:
sudo ./wifi-bander
```

### On macOS:
```bash
./wifi-bander
```

## Real Output Example

```
WiFi Bander - Cross-Platform WiFi Network Analyzer
Initializing scanner...
Scanner initialized successfully. Starting continuous scan...

=== WiFi Network Analysis - 13:18:14 ===
SSID             Band  Channel  Signal (dBm)  Stations  Congestion  Freq (MHz)  
----             ----  -------  -----------   --------  ----------  ---------   
TP-Link_B618     5G    48       -81           1         High        5240        
BASH-2_5G        5G    36       -91           1         High        5180        
kv140_5G         5G    48       -91           1         High        5240        
TP-Link_91E8_5G  5G    40       -88           1         High        5200        
Vasyl Ilba       5G    100      -88           1         High        5500        
NeroDiablo 5G    5G    36       -78           2         High        5180        
SuperMax         2.4G  4        -66           3         Very High   2427        
TP-Link_91E8     2.4G  7        -63           3         Very High   2442        
NeroDiablo 2G    2.4G  9        -60           3         Very High   2452        
TP-Link_B618     2.4G  10       -58           4         Very High   2457        

=== Channel Recommendations ===
Recommended channels for optimal performance:
Band  Best Channels  
----  -------------  
2.4G  11, 6, 1       
5G    44, 149, 153   

Press Ctrl+C to exit...
```

## Understanding the Output

### Signal Strength (dBm)
- **-30 to -50**: Excellent signal
- **-50 to -70**: Good signal  
- **-70 to -80**: Fair signal
- **-80+**: Weak signal

### Congestion Levels
- **Low** (≤15): Minimal interference, optimal performance
- **Medium** (16-30): Some interference, good performance
- **High** (31-50): Significant interference, consider switching
- **Very High** (>50): Heavy interference, definitely switch channels

### Channel Information
- **2.4GHz optimal channels**: 1, 6, 11 (non-overlapping)
- **5GHz channels**: Much more available, less congested
- **Frequency**: Shows exact MHz for technical analysis

## Why Channel Analysis Instead of User Counting?

1. **User counting requires specialized hardware** or monitor mode
2. **Channel congestion is the real performance factor**
3. **Cross-platform compatibility** - works with standard APIs
4. **More actionable insights** - you can actually change channels
5. **Privacy-friendly** - doesn't attempt to sniff user traffic
6. **Reliable detection** - doesn't depend on device cooperation

## Development

### **Building & Testing**
```bash
# Build the application
make build

# Run tests
make test

# Run linting
make lint

# Clean build artifacts
make clean

# Show all available commands
make help
```

### **Adding New Platforms**
1. Create new scanner file: `internal/scanner/{platform}.go`
2. Implement the `Scanner` interface:
   ```go
   type PlatformScanner struct{}
   func (p *PlatformScanner) Scan() ([]WiFiNetwork, error) { ... }
   ```
3. Add platform detection in `internal/scanner/scanner.go`:
   ```go
   case "windows":
       scanner := &WindowsScanner{}
       return scanner.Scan()
   ```
4. Implement platform-specific parsing logic
5. Test with various network configurations

### **Extending Features**
- **Historical tracking**: Store scan results over time
- **Network security analysis**: Add security protocol detection  
- **Configuration files**: Add YAML/JSON config support
- **Web interface**: Create HTTP API endpoints
- **Database storage**: Persist network data

### **Code Organization Principles**
- **Single Responsibility**: Each package has one clear purpose
- **Interface Segregation**: Small, focused interfaces
- **Dependency Injection**: Loose coupling through interfaces
- **Platform Abstraction**: OS-specific code isolated in separate files

## Troubleshooting

### macOS Issues
- **"system_profiler not found"**: Ensure you're on macOS 10.13+
- **No networks detected**: Check WiFi is enabled and scanning is allowed
- **Permission denied**: Try running with `sudo` (usually not needed)

### Linux Issues
- **"nmcli not found"**: Install NetworkManager: `sudo apt install network-manager`
- **"iwlist not found"**: Install wireless tools: `sudo apt install wireless-tools`
- **Permission denied**: Run with `sudo` or add user to appropriate groups
- **No WiFi interface**: Check `ip link` or `iwconfig` for available interfaces

### Common Issues
- **Empty table**: Wait 10-15 seconds for first scan to complete
- **Inconsistent results**: WiFi environment changes frequently, this is normal
- **High congestion everywhere**: Consider 5GHz band or different location

## Performance Notes

- **Scan frequency**: Every 10 seconds (configurable)
- **Memory usage**: Minimal (< 10MB typical)
- **CPU usage**: Low impact, brief spikes during scans
- **Network impact**: Read-only scanning, no interference

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes following the established patterns
4. Add tests if applicable: `make test`
5. Run linting: `make lint`
6. Submit a pull request

### **Development Guidelines**
- Follow Go conventions and `gofmt` styling
- Keep packages focused and interfaces small
- Add comprehensive error handling
- Include documentation for public APIs
- Test on multiple platforms when possible

## License

MIT License

## Changelog

### v2.0.0 (Current - Refactored)
- ✅ **Modular architecture** with `internal/` packages
- ✅ **Interface-based design** for better extensibility
- ✅ **Updated module name** to `github.com/svgreg/wifi-bander`
- ✅ **Separated concerns** into scanner/analyzer/display packages
- ✅ **Enhanced maintainability** and code organization
- ✅ **Production Makefile** with common development tasks

### v1.1.0 (Previous)
- ✅ Fixed macOS compatibility using system_profiler
- ✅ Improved parsing reliability
- ✅ Added comprehensive congestion analysis
- ✅ Clean production output
- ✅ Enhanced error handling

### v1.0.0 (Initial)
- Initial release with basic WiFi scanning
- Cross-platform support (Linux/macOS)
- Channel recommendation system 