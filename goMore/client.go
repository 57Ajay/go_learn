package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error dialing:", err)
		os.Exit(1)
	}
	defer conn.Close()

	input := bufio.NewReader(os.Stdin)

client:
	for {
		fmt.Print("Enter text to send (or type 'exit' to quit): ")
		text, _ := input.ReadString('\n')

		// Send the input to the server
		_, err = conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error sending:", err)
			return
		}

		if text == "exit" {
			fmt.Println("Exiting.")
			break client
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving:", err)
			return
		}
		fmt.Printf("Server replied: %s\n", buffer[:n])
	}
}
