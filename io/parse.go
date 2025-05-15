package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func BufioScanner() {
	// --- Example 1: Reading lines from os.Stdin ---
	fmt.Println("Example 1: bufio.Scanner reading lines from os.Stdin")
	fmt.Println("Enter a few lines of text (Ctrl+D or Ctrl+Z to end):")

	stdinScanner := bufio.NewScanner(os.Stdin)
	lineCount := 0
	for stdinScanner.Scan() {
		lineCount++
		text := stdinScanner.Text()
		fmt.Printf("Line %d: %s\n", lineCount, text)
	}
	// EOF is not reported as an error by Err() if Scan() simply returned false due to it.
	if err := stdinScanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
	}

	// --- Example 2: Reading words from a string and parsing them ---
	fmt.Println("\nExample 2: bufio.Scanner reading words from a string")
	data := "item1 100 item2 250 item3 30"
	stringReader := strings.NewReader(data)
	wordScanner := bufio.NewScanner(stringReader)
	wordScanner.Split(bufio.ScanWords)

	var itemName string
	var itemValue int
	isExpectingName := true

	fmt.Println("Parsing items and values:")
	for wordScanner.Scan() {
		word := wordScanner.Text()
		if isExpectingName {
			itemName = word
			isExpectingName = false
		} else {
			var err error
			itemValue, err = strconv.Atoi(word) // Convert string to int
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error converting '%s' to int: %v\n", word, err)
				continue
			}
			fmt.Printf("Item: %s, Value: %d\n", itemName, itemValue)
			isExpectingName = true
		}
	}
	if err := wordScanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning words: %v\n", err)
	}
}

func SscanFscan() {
	fmt.Println("Example 1: fmt.Scan")
	fmt.Println("Enter your name and age (e.g., Alice 30):")

	var name string
	var age int

	n, err := fmt.Scan(&name, &age)
	if err != nil {
		if err == io.EOF {
			fmt.Fprintln(os.Stderr, "Reached EOF before scanning anything.")
		} else {
			fmt.Fprintf(os.Stderr, "Error scanning input: %v (scanned %d items)\n", err, n)
		}
	} else {
		fmt.Printf("Hello, %s! You are %d years old. (Scanned %d items)\n", name, age, n)
	}

	// --- Example 2: fmt.Fscan (reading from a strings.Reader) ---
	fmt.Println("\nExample 2: fmt.Fscan")
	inputString := "Bob 25 true 3.14"
	reader := strings.NewReader(inputString)

	var fName string
	var fAge int
	var fIsRegistered bool
	var fScore float64

	// Fscan will read space-separated values from the reader
	n, err = fmt.Fscan(reader, &fName, &fAge, &fIsRegistered, &fScore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error using Fscan: %v (scanned %d items)\n", err, n)
	} else {
		fmt.Printf("Fscan result: Name=%s, Age=%d, Registered=%t, Score=%.2f (Scanned %d items)\n",
			fName, fAge, fIsRegistered, fScore, n)
	}

	// Check if there's anything left in the reader
	remainingBytes, _ := io.ReadAll(reader)
	if len(remainingBytes) > 0 {
		fmt.Printf("Remaining in reader: '%s'\n", string(remainingBytes))
	}

	// --- Example 3: fmt.Scanf (reading formatted input from os.Stdin) ---
	fmt.Println("\nExample 3: fmt.Scanf")
	fmt.Println("Enter a date (MM/DD/YYYY):")
	var month, day, year int
	// If the user enters "12-25-2023", this Scanf will fail.
	n, err = fmt.Scanf("%d/%d/%d", &month, &day, &year)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error using Scanf: %v (scanned %d items)\n", err, n)
	} else {
		fmt.Printf("Parsed date: Month=%02d, Day=%02d, Year=%d (Scanned %d items)\n", month, day, year, n)
	}
}

func ParseMain() {
	// SscanFscan()
	BufioScanner()
}
