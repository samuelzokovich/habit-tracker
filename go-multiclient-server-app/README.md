# Go Multi-Client Server Application

This project implements a multi-client server application in Go. It consists of a server that can handle multiple client connections and a client that connects to the server to perform health checks.

## Project Structure

```
go-multiclient-server-app
├── cmd
│   ├── client
│   │   └── main.go        # Entry point for the client application
│   └── server
│       └── main.go        # Entry point for the server application
├── internal
│   ├── health
│   │   └── health.go      # Logic for health check functionality
│   └── network
│       ├── client.go      # Client-side networking logic
│       └── server.go      # Server-side networking logic
├── go.mod                  # Module definition for the Go application
└── README.md               # Documentation for the project
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd go-multiclient-server-app
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Run the server:**
   ```
   go run cmd/server/main.go
   ```

4. **Run the client:**
   ```
   go run cmd/client/main.go
   ```

## Usage

- The server listens for incoming client connections and handles health check requests.
- The client connects to the server, performs a health check, and prints the latency of the connection.

## Contributing

Feel free to submit issues or pull requests for improvements or bug fixes.