package main

import (
	"bufio"
	"fmt"
	"game-service/connect-client/clientutils"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <role: 0 for phone, 1 for TV>", os.Args[0])
	}

	// Parse role argument
	roleArg := os.Args[1]
	roleNum, err := strconv.Atoi(roleArg)
	if err != nil || (roleNum != 0 && roleNum != 1) {
		log.Fatalf("Invalid role: %s. Use 0 for phone or 1 for TV.", roleArg)
	}

	// Set role based on argument
	role := "phone"
	if roleNum == 1 {
		role = "tv"
	}

	baseURL := "localhost:8081"
	tvID := "tv123" // Replace with the unique TV ID

	// Connect to the WebSocket server
	conn, err := clientutils.ConnectWebSocket(baseURL, tvID, role)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Channel to signal the listener to stop
	stopChan := make(chan struct{})

	// Start a goroutine to listen for messages
	go func() {
		clientutils.ReceiveMessages(conn, stopChan)
	}()

	// Role-specific logic for sending messages
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("[%s] Enter command (or 'exit' to quit): \n", strings.ToUpper(role))
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		// Send other commands as an Instruction
		instruction := clientutils.Instruction{
			Role:    role,
			Command: command,
		}
		err = clientutils.SendInstruction(conn, instruction)
		if err != nil {
			log.Printf("Failed to send exit instruction: %v", err)
		}

		// Wait for a short duration to ensure the server receives the message

		// Handle exit command from the user
		if command == "exit" || (command == "disconnect" && role == "phone") {
			// Signal the listener goroutine to stop
			close(stopChan)
			conn.Close()
			break
		}
	}

	log.Printf("%s client exiting...", strings.ToUpper(role))
}
