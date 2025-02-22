package netPackage

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Server side

func broadcastServer(addr *net.UDPAddr) {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Broadcast server listening on:", addr)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message to broadcast: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		_, err := conn.WriteToUDP([]byte(message), &net.UDPAddr{IP: net.IPv4(192, 168, 4, 255), Port: addr.Port}) // Corrected broadcast address
		if err != nil {
			fmt.Println("Error broadcasting:", err)
		}
	}
}

// Client side

func broadcastClient(addr *net.UDPAddr) {
	//bind to any available local port.
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0}) //bind to port 0 to let OS choose an available port
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Broadcast client listening on:", conn.LocalAddr())

	buffer := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}
		fmt.Println("Received:", string(buffer[:n]))
	}
}

func BroadcastUdp() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go server <port> or go run main.go client <port>")
		return
	}

	port := 3000 //Default port
	if len(os.Args) > 2 {
		fmt.Sscanf(os.Args[2], "%d", &port)
	}

	addr := &net.UDPAddr{Port: port, IP: net.IPv4zero}

	if os.Args[1] == "server" {
		broadcastServer(addr)
	} else if os.Args[1] == "client" {
		broadcastClient(addr)
	} else {
		fmt.Println("Invalid argument. Use 'server' or 'client'.")
	}
}
