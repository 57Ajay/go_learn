package netPackage

import (
	"fmt"
	"net"
	"os"
)

func goUdp() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":3000")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer ln.Close()

	fmt.Println("Listening on udp port 3000")

	buffer := make([]byte, 1024)

	for {
		n, addr, err := ln.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		fmt.Println("Received: ", string(buffer[:n]))

		// This will echo the message back to the client
		_, err = ln.WriteToUDP(buffer[:n], addr)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func UdpMain() {
	goUdp()
}
