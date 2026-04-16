# Go Subnet Calculator

A command-line tool written in Go that calculates subnet divisions for IPv4 networks. This tool helps network administrators and students understand subnetting by dividing a network into multiple subnets and displaying all relevant networking information.

## Features

- **Flexible Input**: Accepts IP addresses with or without CIDR notation
- **Automatic Class Detection**: Determines Class A/B/C networks when CIDR is not provided
- **Power-of-2 Subnetting**: Automatically rounds up to the next power of 2 for optimal subnet allocation
- **Comprehensive Output**: Displays all essential subnet information including:
  - Network ID
  - Subnet Mask (both CIDR and dotted decimal notation)
  - First and Last usable IP addresses
  - Broadcast address
  - Total and usable host counts

## Requirements

- Go 1.16 or higher

## Installation

1. Clone or download this repository
2. Navigate to the project directory
3. Build the program:
```bash
go build subnetting.go
```

This will create an executable named `subnetting` (or `subnetting.exe` on Windows).

Alternatively, you can run it directly without building:
```bash
go run subnetting.go
```

## Usage

Run the compiled program:
```bash
./subnetting
```

Or on Windows:
```bash
subnetting.exe
```

The program will prompt you for two inputs:

1. **IP Address**: Enter an IP address with optional CIDR notation
   - With CIDR: `192.168.1.0/24`
   - Without CIDR: `192.168.1.0` (class-based mask will be applied)

2. **Number of Subnets**: Enter how many subnets you need (minimum 2)

### Example Session

```
Go Subnet Calculator - Advanced Tool 

Enter IP address (with or without CIDR, e.g., 192.168.1.0/24 or 192.168.1.0): 192.168.1.0/24
Enter number of subnets required: 4

Calculating subnets...

================================================================================
                    SUBNET CALCULATION RESULTS
================================================================================

Original Network: 192.168.1.0/24
Requested Subnets: 4
Actual Subnets Created: 4 (next power of 2)
New Subnet Mask: 255.255.255.192 (/26)
Hosts per Subnet: 64 total, 62 usable

--------------------------------------------------------------------------------

Subnet #1:
  Network ID:       192.168.1.0/26
  Subnet Mask:      255.255.255.192
  First Usable IP:  192.168.1.1
  Last Usable IP:   192.168.1.62
  Broadcast ID:     192.168.1.63
  Total Hosts:      64
  Usable Hosts:     62
--------------------------------------------------------------------------------

Subnet #2:
  Network ID:       192.168.1.64/26
  Subnet Mask:      255.255.255.192
  First Usable IP:  192.168.1.65
  Last Usable IP:   192.168.1.126
  Broadcast ID:     192.168.1.127
  Total Hosts:      64
  Usable Hosts:     62
--------------------------------------------------------------------------------

Subnet #3:
  Network ID:       192.168.1.128/26
  Subnet Mask:      255.255.255.192
  First Usable IP:  192.168.1.129
  Last Usable IP:   192.168.1.190
  Broadcast ID:     192.168.1.191
  Total Hosts:      64
  Usable Hosts:     62
--------------------------------------------------------------------------------

Subnet #4:
  Network ID:       192.168.1.192/26
  Subnet Mask:      255.255.255.192
  First Usable IP:  192.168.1.193
  Last Usable IP:   192.168.1.254
  Broadcast ID:     192.168.1.255
  Total Hosts:      64
  Usable Hosts:     62
--------------------------------------------------------------------------------

✓ Calculation complete!
```

## How It Works

### 1. **IP Address Parsing**
The program accepts IP addresses in two formats:
- **With CIDR**: `192.168.1.0/24` - Uses the specified subnet mask
- **Without CIDR**: `192.168.1.0` - Applies class-based default masks:
  - Class A (1-126): /8
  - Class B (128-191): /16
  - Class C (192-223): /24

### 2. **Subnet Calculation**
- Calculates the number of bits needed based on the requested subnets
- Uses the formula: `bits_needed = ceil(log2(requested_subnets))`
- Rounds up to the nearest power of 2 (e.g., requesting 5 subnets creates 8)
- Adds the calculated bits to the original CIDR to determine the new subnet mask

### 3. **Address Range Calculation**
For each subnet, the program calculates:
- **Network ID**: The first address in the subnet (not usable for hosts)
- **First Usable IP**: Network ID + 1
- **Last Usable IP**: Broadcast ID - 1
- **Broadcast ID**: The last address in the subnet (not usable for hosts)
- **Usable Hosts**: Total addresses - 2 (excludes network and broadcast)

## Example Use Cases

### Example 1: Creating 4 subnets from a /24 network
```
Input: 192.168.1.0/24, 4 subnets
Result: 4 subnets with /26 masks (64 hosts each, 62 usable)
```

### Example 2: Creating 10 subnets (rounds to 16)
```
Input: 10.0.0.0/8, 10 subnets
Result: 16 subnets with /12 masks (1,048,576 hosts each)
```

### Example 3: Class-based without CIDR
```
Input: 172.16.0.0, 8 subnets
Auto-detects: Class B (/16)
Result: 8 subnets with /19 masks (8,192 hosts each, 8,190 usable)
```

## Understanding the Output

| Field | Description |
|-------|-------------|
| **Network ID** | The network address of the subnet (first address) |
| **Subnet Mask** | The mask in dotted decimal notation |
| **CIDR** | The subnet mask in CIDR notation (e.g., /26) |
| **First Usable IP** | The first IP address that can be assigned to a host |
| **Last Usable IP** | The last IP address that can be assigned to a host |
| **Broadcast ID** | The broadcast address for the subnet (last address) |
| **Total Hosts** | Total number of addresses in the subnet |
| **Usable Hosts** | Number of addresses available for hosts (Total - 2) |

## Limitations

- Maximum supported CIDR is /30 (4 addresses, 2 usable)
- IPv4 only (IPv6 not supported)
- Does not support Variable Length Subnet Masking (VLSM)
- Always creates power-of-2 number of subnets

## Error Handling

The program handles common errors:
- Invalid IP address format
- Empty input
- Non-numeric subnet count
- Too many subnets requested (would exceed /30)
- Minimum 2 subnets required

## File Structure

```
.
├── subnetting.go    # Main program file
└── README.md        # This file
```

## Building for Different Platforms

You can cross-compile for different operating systems:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build subnetting.go

# Windows
GOOS=windows GOARCH=amd64 go build subnetting.go

# macOS
GOOS=darwin GOARCH=amd64 go build subnetting.go
```

## Use Cases

- **Network Planning**: Divide a network into departmental subnets
- **Learning**: Understand subnetting concepts and practice calculations
- **Documentation**: Generate subnet information for network documentation
- **VLSM Planning**: Determine base subnet divisions before applying VLSM
- **Certification Prep**: Practice for networking certifications (CCNA, Network+)

## License

This program is provided as-is for educational and professional use.

## Contributing

Feel free to fork, modify, and improve this tool. Suggestions and pull requests are welcome!