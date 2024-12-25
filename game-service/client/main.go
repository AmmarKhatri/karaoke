package main

import (
	"game-service/client/scripts"
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
	roomID := "room-fcc8b5cd-b879-45a5-9c59-18a2f5058791"
	listenerID := "listener1"
	producerID := "producer1"

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
