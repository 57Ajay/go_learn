package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Password string `json:"-"` // Always omitted
	IsActive bool   `json:"isActive"`
}

type Product struct {
	Name     string         `json:"productName"`
	Price    float64        `json:"price"`
	Tags     []string       `json:"tags"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func JsonMarshalMain() {
	user1 := User{
		ID:       1,
		Username: "johndoe",
		Password: "supersecret", // This will be ignored
		IsActive: true,
	}
	userJSON, err := json.Marshal(user1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling user: %v\n", err)
		return
	}
	fmt.Println("User JSON:", string(userJSON))

	// Marshalling a slice of structs
	users := []User{
		{ID: 2, Username: "janedoe", Email: "jane@example.com", IsActive: true},
		{ID: 3, Username: "guest", IsActive: false},
	}
	usersJSON, err := json.Marshal(users)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling users slice: %v\n", err)
		return
	}
	fmt.Println("Users Slice JSON:", string(usersJSON))

	// Marshalling a map
	productMeta := map[string]any{
		"sku":    "SKU12345",
		"weight": 2.5,
		"notes":  "Handle with care",
	}
	product1 := Product{
		Name:     "Awesome Gadget",
		Price:    99.99,
		Tags:     []string{"tech", "gadget", "cool"},
		Metadata: productMeta,
	}
	productJSON, err := json.Marshal(product1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling product: %v\n", err)
		return
	}
	fmt.Println("Product JSON:", string(productJSON))

	prettyUserJSON, err := json.MarshalIndent(user1, "", "  ") // "" prefix, "  " (2 spaces) for indentation
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling with indent: %v\n", err)
		return
	}
	fmt.Println("\nPretty User JSON:")
	fmt.Println(string(prettyUserJSON))
}

func JsonUnMarshalMain() {
	// Unmarshalling into a struct
	userJSONData := []byte(`{"id":101, "username":"testuser", "email":"test@example.com", "isActive":false, "extraField":"ignored"}`)
	var userFromFile User
	err := json.Unmarshal(userJSONData, &userFromFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling user: %v\n", err)
		return
	}
	fmt.Printf("Unmarshalled User: %+v\n", userFromFile) // City will be empty string (zero value)

	// Unmarshalling into a slice of structs
	usersJSONData := []byte(`[
		{"id":201, "username":"alice", "isActive":true},
		{"id":202, "username":"bob", "email":"bob@work.com", "isActive":false}
	]`)
	var usersFromFile []User
	err = json.Unmarshal(usersJSONData, &usersFromFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling users slice: %v\n", err)
		return
	}
	fmt.Println("Unmarshalled Users Slice:")
	for _, u := range usersFromFile {
		fmt.Printf("  %+v\n", u)
	}

	// Unmarshalling into a map[string]any for arbitrary JSON structure
	productJSONData := []byte(`{
		"productName": "Super Widget",
		"price": 49.95,
		"tags": ["utility", "new"],
		"dimensions": {"height": 10, "width": 5, "depth": 2.5},
		"manufacturer": null
	}`)
	var productData map[string]any
	err = json.Unmarshal(productJSONData, &productData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling product data into map: %v\n", err)
		return
	}
	fmt.Println("Unmarshalled Product Data (map):")
	for key, value := range productData {
		fmt.Printf("  %s: %v (Type: %T)\n", key, value, value)
	}
	// Accessing a nested map (dimensions)
	if dims, ok := productData["dimensions"].(map[string]any); ok {
		if height, ok := dims["height"].(float64); ok { // JSON numbers are float64 by default into any
			fmt.Printf("Product height: %.2f\n", height)
		}
	}
	if manu, ok := productData["manufacturer"]; ok && manu == nil {
		fmt.Println("Manufacturer is JSON null.")
	}
}

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

func EncoderDecoderMain() {
	// --- Using json.Encoder ---
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	msg1 := Message{Sender: "Alice", Receiver: "Bob", Content: "Hello Bob!"}
	msg2 := Message{Sender: "Bob", Receiver: "Alice", Content: "Hi Alice!"}

	fmt.Println("Encoding messages to buffer...")
	err := encoder.Encode(msg1) // Encode writes a newline after each JSON object
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding msg1: %v\n", err)
	}
	err = encoder.Encode(msg2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding msg2: %v\n", err)
	}
	fmt.Println("Buffer content after encoding (note newlines):")
	fmt.Print(buf.String())

	// --- Using json.Decoder ---
	// The buffer now contains two JSON objects, one per line
	jsonStream := buf.String()
	decoder := json.NewDecoder(strings.NewReader(jsonStream))

	fmt.Println("\nDecoding messages from stream:")
	var decodedMsg Message
	// Decode multiple JSON objects from the stream
	for decoder.More() { // Check if there is more data to decode
		err := decoder.Decode(&decodedMsg) // Pass pointer
		if err != nil {
			// io.EOF is expected at the end of a valid stream of JSON objects
			// if decoder.More() was true but Decode fails for other reasons, it's an error.
			if err.Error() == "EOF" { // Simplistic EOF check
				fmt.Println("Reached end of JSON stream.")
				break
			}
			fmt.Fprintf(os.Stderr, "Error decoding message: %v\n", err)
			break // Exit on other errors
		}
		fmt.Printf("  Decoded: %+v\n", decodedMsg)
	}

	// Example with a file (conceptual)
	/*
	   // Writing to a file
	   file, _ := os.Create("messages.json")
	   defer file.Close()
	   fileEncoder := json.NewEncoder(file)
	   fileEncoder.Encode(msg1)
	   fileEncoder.Encode(msg2)

	   // Reading from a file
	   readFile, _ := os.Open("messages.json")
	   defer readFile.Close()
	   fileDecoder := json.NewDecoder(readFile)
	   var m Message
	   for fileDecoder.More() {
	       fileDecoder.Decode(&m)
	       fmt.Printf("From file: %+v\n", m)
	   }
	   os.Remove("messages.json") // cleanup
	*/
}
