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

	// Role-specific logic
	if role == "tv" {
		// TV listens for instructions from the phone
		log.Println("TV is ready to receive instructions...")
		clientutils.ReceiveMessages(conn)
	} else if role == "phone" {
		// Phone sends instructions to the TV
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter command for TV (or 'exit' to quit): ")
			command, _ := reader.ReadString('\n')
			command = strings.TrimSpace(command)

			if command == "exit" {
				break
			}

			// Send the command as an Instruction
			instruction := clientutils.Instruction{
				Role:    role,
				Command: command,
			}

			err = clientutils.SendInstruction(conn, instruction)
			if err != nil {
				log.Printf("Failed to send instruction: %v", err)
				break
			}
		}

		log.Println("Phone client exiting...")
	}
}
