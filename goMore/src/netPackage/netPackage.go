package netPackage

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Error Reading: ", err)
			}
			return
		}
		fmt.Println("Received: ", string(buffer[:n]))

		if num, err := strconv.Atoi(strings.TrimSpace(string(buffer[:n]))); err == nil {
			fmt.Println("Received number: ", num)
			sum := num * (num + 1) / 2
			fmt.Println("Sum is: ", sum)
			conn.Write([]byte(strconv.Itoa(sum)))
		} else {
			_, err = conn.Write(buffer[:n])
			if err != nil {
				fmt.Println("Error writing:", err)
				return
			}
		}

	}
}

func startConnection() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn)
	}

}

func Netmain() {
	// startConnection()
	// IPMain()
	// UdpMain()
	BroadcastUdp()
}
