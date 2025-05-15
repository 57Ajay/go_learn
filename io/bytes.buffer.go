package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func processData(r io.Reader) {
	fmt.Println("\n--- Processing data (expecting an io.Reader) ---")
	data, err := io.ReadAll(r) // io.ReadAll reads everything from r until EOF
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in processData: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %s\n", string(data))
}

func ProcessData() {
	s := "This is a string that we want to treat as an io.Reader."

	stringReader := strings.NewReader(s)

	fmt.Println("Copying from strings.Reader to os.Stdout:")
	_, err := io.Copy(os.Stdout, stringReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying stringReader to stdout: %v\n", err)
	}
	fmt.Println("\n--- End of copy ---")

	offset, err := stringReader.Seek(0, io.SeekStart) // 0 means from the beginning
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error seeking: %v\n", err)
		return
	}
	fmt.Printf("Seeked to offset: %d\n", offset)

	// Now pass it to our generic processData function
	processData(stringReader)

	// Demonstrate reading with a buffer
	stringReader.Seek(0, io.SeekStart) // Reset for another read
	smallBuffer := make([]byte, 10)
	n, err := stringReader.Read(smallBuffer)
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Error reading into smallBuffer: %v\n", err)
	}
	fmt.Printf("\nRead %d bytes into smallBuffer: '%s'\n", n, string(smallBuffer[:n]))
	fmt.Printf("Remaining in stringReader: %d bytes\n", stringReader.Len())
}

func BytesBufferMain() {
	var buf bytes.Buffer

	// Writing to the buffer (bytes.Buffer implements io.Writer)
	s1 := "Hello, "
	buf.WriteString(s1)                                    // Write a string
	fmt.Fprintf(&buf, "World! From %s.\n", "bytes.Buffer") // Use fmt.Fprintf

	fmt.Printf("Buffer length: %d\n", buf.Len())
	fmt.Printf("Buffer content: %s\n", buf.String())

	// Reading from the buffer (bytes.Buffer implements io.Reader)
	fmt.Println("\nReading from buffer:")
	// We can copy its content to os.Stdout (which is an io.Writer)
	// io.Copy will read from buf until EOF (or error) and write to os.Stdout
	cont, err := io.Copy(os.Stdout, &buf) // Note: After this, buf will be empty
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying buffer to stdout: %v\n", err)
	}

	fmt.Println("Cont: ", cont)

	// The buffer is now empty because io.Copy read all its content
	fmt.Printf("\nBuffer length after reading with io.Copy: %d\n", buf.Len())

	// Let's write to it again
	buf.WriteString("Another piece of data.")

	// Read using the Read method directly
	outputBuffer := make([]byte, 10)
	n, err := buf.Read(outputBuffer)
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Error reading from buffer: %v\n", err)
	}
	fmt.Printf("Read %d bytes: %s\n", n, string(outputBuffer[:n]))

	n, err = buf.Read(outputBuffer) // Try reading again
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Error reading from buffer: %v\n", err)
	}
	fmt.Printf("Read %d bytes: %s\n", n, string(outputBuffer[:n]))
}
