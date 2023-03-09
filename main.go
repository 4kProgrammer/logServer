package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// Create a new file for logging
	logFile, err := os.OpenFile("/var/log/tcpserver.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err.Error())
		return
	}
	defer logFile.Close()

	const (
		HOST = "94.101.186.207"
		PORT = "8080"
		TYPE ="tcp"
	)

	// Listen for incoming connections on port 8080
	listener, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Listening on :8080")

	for {
		// Wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn, logFile)
	}
}

func handleConnection(conn net.Conn, logFile *os.File) {
	defer conn.Close()

	// Read data from the connection
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	// Log the received data to the file
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("%s - Received data: %s\n", timestamp, string(buffer))
	_, err = logFile.WriteString(logMsg)
	if err != nil {
		fmt.Println("Error writing to log file:", err.Error())
		return
	}

	// Send a response back to the client
	response := "Hello, client!"
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
		return
	}

	fmt.Println("Response sent")
}