package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func BufioEx1() {
	fileName := "Bufio.txt"
	content := "Line 1: Hello from Go.\nLine 2: This is a test of bufio.Reader.\nLine 3: Enjoy!"
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writiong the file: %s", err)
		return
	}
	defer os.Remove(fileName)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return
	}
	defer file.Close()
	bufferedReader := bufio.NewReader(file)
	fmt.Println("Reading line by line using ReadString method: fileName: ", fileName)

	for i := 1; ; i++ {
		line, err := bufferedReader.ReadString('\n')
		if len(line) > 0 {
			fmt.Printf("Line-> %d: %s\n", i, strings.TrimSpace(line))
			if strings.HasSuffix(line, "\n") {
				fmt.Println("Ended with line: ", i)
			} else {
				fmt.Println(" (ended without newline - likely EOF)")
			}
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
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error seeking file: %v\n", err)
		return
	}
	bufferedReader.Reset(file)
	fmt.Println("\nReading first 10 bytes one by one using ReadByte():")
	for i := range 10 {
		b, err := bufferedReader.ReadByte()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading byte: %v\n", err)
			break
		}
		fmt.Printf("Byte %d: %c (%d)\n", i+1, b, b)
	}

}

func BufioEx2() {
	fileName := "bufioTest.txt"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error creating the file: %s", err)
	}
	defer os.Remove(fileName)
	bufferedWriter := bufio.NewWriter(file)
	defer bufferedWriter.Flush()
	defer file.Close()

	linesToWrite := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the contents of the file (type 'END' on a new line to finish):")
	fmt.Print("Start from here ->>>> \n")
	for {
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		if strings.TrimSpace(text) == "END" {
			break
		}
		linesToWrite = append(linesToWrite, text)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return
	}

	for i, lines := range linesToWrite {
		n, err := bufferedWriter.WriteString(lines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing string to buffer: %v\n", err)
			return
		}
		fmt.Printf("Wrote %d bytes for line %d. Buffered: %d bytes.\n", n, i+1, bufferedWriter.Buffered())
		if bufferedWriter.Buffered() > 1024 { // Arbitrary threshold for testing
			fmt.Println("Buffer threshold reached, flushing...")
			err = bufferedWriter.Flush()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error flushing buffer: %v\n", err)
				return
			}
			fmt.Println("Flushed. Buffer now empty.")
		}
	}
	fmt.Println("All lines queued to writer.")
	fmt.Println("Final Flush will happen on defer (or explicitly if called before defer).")
}

func BufioMain() {
	// bufioEx1()
	BufioEx2()
}
