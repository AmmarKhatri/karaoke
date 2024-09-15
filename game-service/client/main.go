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

	// Create a 'done' channel to signal when to exit

	// RoomID and player IDs
	roomID := "room-7174a51b-a936-4663-8238-faae7466f13b"
	listenerID := "listener1"
	producerID := "producer1"

	// Check input argument
	mode := os.Args[1]

	if mode == "0" {
		// Start listener in a goroutine
		go scripts.StartListener(roomID, listenerID, interrupt)
	} else if mode == "1" {
		// Start producer
		scripts.StartProducer(roomID, producerID, interrupt)
	} else {
		log.Fatal("Invalid argument. Use 0 for listener and 1 for producer.")
	}

	// Wait for interrupt signal or done signal to terminate all connections
	select {
	case <-interrupt:
		log.Println("Received interrupt signal. Closing all connections.")
		os.Exit(0) // Exit after the interrupt signal
	}

	// Perform any cleanup if necessary
	os.Exit(0)
}
