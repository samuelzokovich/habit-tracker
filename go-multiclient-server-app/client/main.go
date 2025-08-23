package main

import (
    "fmt"
    "net"
    "time"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        return
    }
    defer conn.Close()

    start := time.Now()
    _, err = conn.Write([]byte("health"))
    if err != nil {
        fmt.Println("Error sending health check:", err)
        return
    }

    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Error reading response:", err)
        return
    }

    latency := time.Since(start)
    fmt.Printf("Server health status: %s\n", string(buffer[:n]))
    fmt.Printf("Latency: %v\n", latency)
}