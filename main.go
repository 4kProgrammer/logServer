package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	HOST      = "94.101.186.207"
	TCP_PORT  = "8080"
	HTTP_PORT = "8081"
	TYPE      = "tcp"
)

func main() {
	// Create the log file
	logFile, err := os.OpenFile("messages.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	// Start a HTTP server on port 8080
	go func() {
		http.HandleFunc("/log", handleLog)
		http.ListenAndServe(HOST+":"+HTTP_PORT, nil)
	}()

	// Listen for incoming connections on port 8081
	listener, err := net.Listen("tcp", HOST+":"+TCP_PORT)
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()
	fmt.Println("Listening on :8081")

	for {
		// Wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn, logFile)
	}
}

func handleConnection(conn net.Conn, logFile *os.File) {
	//defer conn.Close()

	for {
		// Read data from the connection
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error reading:", err)
			return
		}

		// Print the received data to the console
		message := string(buffer[:n])
		fmt.Println("Received data:", message)

		// Write the received message to the log file
		_, err = logFile.WriteString(message)
		if err != nil {
			fmt.Println("Error writing to log file:", err.Error())
			return
		}

		// Send a response back to the client
		response := "Hello, client!"
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Println("Error writing:", err)
			return
		}

		// Wait for a second before sending the next response
		//time.Sleep(time.Second)
	}
}

func handleLog(w http.ResponseWriter, r *http.Request) {
	// Open the log file and write its contents to the HTTP response

	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// open log file for reading
	logFile, err := os.Open("messages.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logFile.Close()

	// read log file line by line and send each line to client
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		message := scanner.Text()
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			fmt.Println(err)
			break
		}
		time.Sleep(time.Millisecond * 500) // wait half a second before sending next message
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
