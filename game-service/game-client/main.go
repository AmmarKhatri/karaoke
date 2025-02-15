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
	roomID := "room-f8108778-e3a2-4bfc-be06-3c2164c823ac"
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
