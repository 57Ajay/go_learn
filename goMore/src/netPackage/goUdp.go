package netPackage

import (
	"fmt"
	"net"
	"os"
	"time"
)

func goUdp() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":3000")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Println("Listening on udp port 3000")

	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		fmt.Println("Received: ", string(buffer[:n]))
		if string(buffer[:n]) == "exit" {
			_, err = conn.WriteToUDP([]byte("Exiting..."), addr)
			if err != nil {
				fmt.Println("Error Exiting : ", err)
				continue
			}
			break
		}
		if string(buffer[:n]) == "time" {
			_, err = conn.WriteToUDP([]byte("Time is: "+fmt.Sprint(time.Now())), addr)
			continue
		}
		// This will echo the message back to the client
		_, err = conn.WriteToUDP(buffer[:n], addr)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func UdpMain() {
	goUdp()
}
