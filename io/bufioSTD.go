package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func BUFSTDMAIN() {
	fmt.Println("Please enter the contents of the file below: ")
	fmt.Print("Start from here ->>>> ")
	buf := make([]byte, 4069)
	_, err := os.Stdin.Read(buf)
	if err != nil {
		if err == io.EOF {
			fmt.Fprintf(os.Stderr, "\nError occured: %s\n", err)
		}
		fmt.Fprintf(os.Stderr, "Error: %s", err)
	}
	fileName := "bufioTest.txt"
	err = os.WriteFile(fileName, buf, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writiong the file: %s", err)
		return
	}
	defer os.Remove("bufioTest.txt")
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	bufferedReader := bufio.NewReader(file)

	for i := 0; ; i++ {
		line, err := bufferedReader.ReadString('\n')
		if len(line) > 0 {
			fmt.Fprintf(os.Stdout, "Reading line: %d: %s\n", i, strings.TrimSpace(line))
		}
		if strings.HasSuffix(line, "\n") {
			fmt.Fprintf(os.Stdout, "Ended with the line: %d\n", i)
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", "EOF")
		}
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("End of file reached.")
			} else {
				fmt.Fprintf(os.Stderr, "Error reading string: %v\n", err)
			}
			break
		}
	}

}
