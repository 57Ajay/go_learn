package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func OSSTD() {
	for {
		fmt.Print("Enter something: ")
		buf := make([]byte, 1024)
		userInput, err := os.Stdin.Read(buf)
		// userInput, err := io.ReadAll(os.Stdin)
		if err != nil {
			if err == io.EOF {
				fmt.Fprintf(os.Stderr, "\nError occured: %s\n", err)
				break
			}
			fmt.Fprintf(os.Stderr, "Error: %s", err)
		}
		input := strings.TrimSpace(string(buf[:userInput]))
		if strings.EqualFold(input, "Exit") {
			fmt.Fprintf(os.Stdout, "%s\n", "Exiting the programme")
			break
		}
		fmt.Fprintf(os.Stdout, "User Entered: %s", string(buf[:userInput]))
	}
}
