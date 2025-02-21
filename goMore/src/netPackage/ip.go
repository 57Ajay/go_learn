package netPackage

import (
	"fmt"
	"net"
)

func ip() {

}

func parseIp() {

	// Create an IPv4 address
	ip4 := net.ParseIP("192.168.1.1")
	fmt.Println("IPv4:", ip4) // Output: IPv4: 192.168.1.1

	// Create an IPv6 address
	ip6 := net.ParseIP("2001:db8::1")
	fmt.Println("IPv6:", ip6) // Output: IPv6: 2001:db8::1

	// Check if an IP address is IPv4 or IPv6
	fmt.Println("Is IPv4?", ip4.To4() != nil) // Output: Is IPv4? true
	fmt.Println("Is IPv4?", ip6.To4() != nil) // Output: Is IPv4? false

	// Convert IPv4 to IPv6 representation
	ip4to6 := ip4.To16()
	fmt.Println("IPv4 as IPv6:", ip4to6) // Output: IPv4 as IPv6: ::ffff:192.168.1.1

	// Accessing the bytes of an IP Address
	fmt.Println("IP Address Bytes", ip4)

	// An IPv4 address is just four bytes:
	ip4Manual := net.IP{192, 168, 1, 1}
	fmt.Println("Manual IPv4:", ip4Manual) // Note that To4() will fail.
	fmt.Println("Manual IPv4 - To4():", ip4Manual.To4())
	fmt.Println("Manual IPv4 - To16():", ip4Manual.To16())
	ip4Correct := net.IP(net.CIDRMask(32, 32))
	// copy(ip4Correct[12:16], ip4Manual)
	fmt.Println("Corrected:", ip4Correct, "To 4:", ip4Correct.To4())

	// Join host and port
	addr := net.JoinHostPort("192.168.1.1", "8080")
	fmt.Println("Joined address:", addr) // Output: Joined address: 192.168.1.1:8080

	// Split host and port
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("Error splitting:", err)
		return
	}
	fmt.Println("Host:", host) // Output: Host: 192.168.1.1
	fmt.Println("Port:", port) // Output: Port: 8080

	// Example of incorrect usage, missing port
	_, _, err = net.SplitHostPort("192.168.1.1")
	if err != nil {
		fmt.Println("Error splitting:", err) // Output: Error splitting: address 192.168.1.1: missing port in address
	}

}

func resolveIp() {
	// This looks for ip address for a host name; hostname -> IP address
	ips, err := net.LookupIP("google.com")
	if err != nil {
		fmt.Println("Error looking up IP address:", err)
		return
	}
	for _, ip := range ips {
		fmt.Println("IP address:", ip)
	}

	// IP address -> hostname
	names, err := net.LookupAddr("8.8.8.8")
	if err != nil {
		fmt.Println("Error looking up hostname:", err)
	}
	for _, name := range names {
		fmt.Println("Hostname:", name)
	}

	// LookUp port
	port, err := net.LookupPort("tcp", "https")
	if err != nil {
		fmt.Println("Error looking up port:", err)
	}
	fmt.Println("Port:", port)
}

func IPMain() {
	fmt.Println("This is the IPMain function")
	parseIp()
	fmt.Print("\n -------resolveIp-------\n")
	resolveIp()
}
