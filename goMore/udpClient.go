// run this code seperately ( go run udpClient.go) and
// then run make run in the terminal to run the rest of the code

package main

import (
	"fmt"
	"net"
)

func udpClient() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP("localhost"), Port: 3000})
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("Hello from UDP client"))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Received: ", string(buffer[:n]))

}

func main() {
	udpClient()
}
