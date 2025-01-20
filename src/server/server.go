package server

import (
	"fmt"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Extract the query parameter
	queryValue := r.URL.Query().Get("number")

	// Convert query parameter to an integer
	number, err := strconv.Atoi(queryValue)
	if err != nil {
		// Handle the error if conversion fails
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	// Use the converted integer
	sum := number + 100
	fmt.Fprintf(w, "The number is: %d, and after adding 100, it becomes: %d", number, sum)
}

func StartServer() {
	fmt.Println("Server starting on localhost:3000")
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:3000", nil)
}
