// run this code seperately ( go run udpClient.go) and
// then run make run in the terminal to run the rest of the code

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func udpClient() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP("localhost"), Port: 3000})
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer conn.Close()
	fmt.Println("Connected to server on port 3000")
clientLoop:
	for {
		time.Sleep(1 * time.Second)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Raed Error: ", err)
		}
		input = strings.TrimSpace(input)
		if input == "exit" {
			fmt.Println("Exiting...")
			_, err = conn.Write([]byte(input))
			if err != nil {
				fmt.Println("Error Exiting: ", err)
				continue
			}
			break clientLoop
		}
		_, err = conn.Write([]byte(input))
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

}

func main() {
	udpClient()
}
