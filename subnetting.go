package main

import (
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
)

// Subnet represents a single subnet with all its properties
type Subnet struct {
	SubnetNumber     int
	NetworkID        net.IP
	FirstUsableIP    net.IP
	LastUsableIP     net.IP
	BroadcastID      net.IP
	SubnetMask       net.IPMask
	SubnetMaskDotted string
	CIDR             int
	TotalHosts       int
	UsableHosts      int
}

// ipToInt converts an IP address to uint32
func ipToInt(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// intToIP converts uint32 to IP address
func intToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

// calculateSubnets performs the subnetting calculation
func calculateSubnets(ipAddr string, numSubnets int) ([]Subnet, error) {
	// Parse the IP address and get the original CIDR
	ip, ipNet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		// If no CIDR notation, try to determine the class
		ip = net.ParseIP(ipAddr)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address: %s", ipAddr)
		}

		// Determine default class-based mask
		firstOctet := ip.To4()[0]
		var cidr int
		if firstOctet < 128 {
			cidr = 8 // Class A
		} else if firstOctet < 192 {
			cidr = 16 // Class B
		} else {
			cidr = 24 // Class C
		}

		ipNet = &net.IPNet{
			IP:   ip.Mask(net.CIDRMask(cidr, 32)),
			Mask: net.CIDRMask(cidr, 32),
		}
	}

	// Get the original network
	originalCIDR, _ := ipNet.Mask.Size()

	// Calculate bits needed for subnets
	bitsNeeded := int(math.Ceil(math.Log2(float64(numSubnets))))

	// Calculate new CIDR
	newCIDR := originalCIDR + bitsNeeded
	if newCIDR > 30 {
		return nil, fmt.Errorf("too many subnets requested - would result in CIDR /%d (max is /30)", newCIDR)
	}

	// Calculate subnet size
	subnetSize := uint32(1 << (32 - newCIDR))

	// Create subnet mask
	subnetMask := net.CIDRMask(newCIDR, 32)

	// Get network address as integer
	networkInt := ipToInt(ipNet.IP.Mask(ipNet.Mask))

	// Calculate actual number of subnets created (power of 2)
	actualSubnets := 1 << bitsNeeded

	subnets := make([]Subnet, actualSubnets)

	for i := 0; i < actualSubnets; i++ {
		// Calculate network ID for this subnet
		subnetNetworkInt := networkInt + (uint32(i) * subnetSize)
		networkID := intToIP(subnetNetworkInt)

		// Calculate broadcast address
		broadcastInt := subnetNetworkInt + subnetSize - 1
		broadcastID := intToIP(broadcastInt)

		// Calculate first and last usable IPs
		firstUsableInt := subnetNetworkInt + 1
		lastUsableInt := broadcastInt - 1

		firstUsableIP := intToIP(firstUsableInt)
		lastUsableIP := intToIP(lastUsableInt)

		// Calculate host counts
		totalHosts := int(subnetSize)
		usableHosts := totalHosts - 2 // Subtract network and broadcast

		if usableHosts < 0 {
			usableHosts = 0
		}

		subnets[i] = Subnet{
			SubnetNumber:     i + 1,
			NetworkID:        networkID,
			FirstUsableIP:    firstUsableIP,
			LastUsableIP:     lastUsableIP,
			BroadcastID:      broadcastID,
			SubnetMask:       subnetMask,
			SubnetMaskDotted: fmt.Sprintf("%d.%d.%d.%d", subnetMask[0], subnetMask[1], subnetMask[2], subnetMask[3]),
			CIDR:             newCIDR,
			TotalHosts:       totalHosts,
			UsableHosts:      usableHosts,
		}
	}

	return subnets, nil
}

// printSubnetInfo prints detailed information about all subnets
func printSubnetInfo(subnets []Subnet, originalIP string, requestedSubnets int) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("                    SUBNET CALCULATION RESULTS")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\nOriginal Network: %s\n", originalIP)
	fmt.Printf("Requested Subnets: %d\n", requestedSubnets)
	fmt.Printf("Actual Subnets Created: %d (next power of 2)\n", len(subnets))
	fmt.Printf("New Subnet Mask: %s (/%d)\n", subnets[0].SubnetMaskDotted, subnets[0].CIDR)
	fmt.Printf("Hosts per Subnet: %d total, %d usable\n", subnets[0].TotalHosts, subnets[0].UsableHosts)

	fmt.Println("\n" + strings.Repeat("-", 80))

	for _, subnet := range subnets {
		fmt.Printf("\nSubnet #%d:\n", subnet.SubnetNumber)
		fmt.Printf("  Network ID:       %s/%d\n", subnet.NetworkID, subnet.CIDR)
		fmt.Printf("  Subnet Mask:      %s\n", subnet.SubnetMaskDotted)
		fmt.Printf("  First Usable IP:  %s\n", subnet.FirstUsableIP)
		fmt.Printf("  Last Usable IP:   %s\n", subnet.LastUsableIP)
		fmt.Printf("  Broadcast ID:     %s\n", subnet.BroadcastID)
		fmt.Printf("  Total Hosts:      %d\n", subnet.TotalHosts)
		fmt.Printf("  Usable Hosts:     %d\n", subnet.UsableHosts)
		fmt.Println(strings.Repeat("-", 80))
	}
}

func main() {
	fmt.Println("Go Subnet Calculator - Advanced Tool ")

	// Get IP address from user
	var ipAddr string
	fmt.Print("\nEnter IP address (with or without CIDR, e.g., 192.168.1.0/24 or 192.168.1.0): ")
	fmt.Scanln(&ipAddr)

	if strings.TrimSpace(ipAddr) == "" {
		fmt.Println("Error: IP address cannot be empty")
		os.Exit(1)
	}

	// Get number of subnets
	var numSubnetsStr string
	fmt.Print("Enter number of subnets required: ")
	fmt.Scanln(&numSubnetsStr)

	numSubnets, err := strconv.Atoi(numSubnetsStr)
	if err != nil || numSubnets < 2 {
		fmt.Println("Error: Please enter a valid number of subnets (minimum 2)")
		os.Exit(1)
	}

	// Calculate subnets
	fmt.Println("\nCalculating subnets...")
	subnets, err := calculateSubnets(ipAddr, numSubnets)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Print results
	printSubnetInfo(subnets, ipAddr, numSubnets)

	fmt.Println("\n✓ Calculation complete!")
}
