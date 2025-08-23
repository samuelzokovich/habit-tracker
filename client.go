package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// measureRequestTime sends a request and measures round-trip latency
func measureRequestTime(url string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect: %v", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("Response from %s: %s\n", url, string(body))
	fmt.Printf("Round-trip time: %v\n", duration)
}

func main() {
	serverBase := "http://localhost:8080"

	fmt.Println("Checking server health...")
	measureRequestTime(serverBase + "/health")

	fmt.Println("\nCalling hello endpoint...")
	measureRequestTime(serverBase + "/hello")
}
