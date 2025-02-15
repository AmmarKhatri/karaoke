package main

import (
	"game-service/game-client/scripts"
	"log"
	"os"
	"os/signal"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [0|1] (0 for listener, 1 for producer)")
	}

	// Handle interrupt signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// RoomID and player IDs
	roomID := "room-c534ac0a-c181-45e5-b068-b815afb1b82e"
	listenerID := "Ammar"
	producerID := "Ammar"

	// Check input argument
	mode := os.Args[1]

	if mode == "0" {
		// Start listener
		scripts.StartListener(roomID, listenerID, interrupt)
	} else if mode == "1" {
		// Start producer
		scripts.StartProducer(roomID, producerID, interrupt)
	} else {
		log.Fatal("Invalid argument. Use 0 for listener and 1 for producer.")
	}
}
