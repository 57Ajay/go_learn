package main

import (
	"fmt"
	"os"
)

func Fmain() {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_WRONLY, 0644) //os.OpenFile gives more control
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Println("Successfully created/opened output.txt for writing.")

	data1 := []byte("Hello from Go!\n")
	data2 := []byte("Writing to a file is easy.\n")

	n1, err := file.Write(data1) // file implements io.Writer
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing data1: %v\n", err)
		return
	}
	fmt.Printf("Wrote %d bytes for data1\n", n1)

	// Write data2
	n2, err := file.Write(data2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing data2: %v\n", err)
		return
	}
	fmt.Printf("Wrote %d bytes for data2\n", n2)

	_, err = fmt.Fprintln(file, "This line was written using fmt.Fprintln!")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error using fmt.Fprintln: %v\n", err)
	}

	fmt.Println("Data written to output.txt")
}
