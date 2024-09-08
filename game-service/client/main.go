package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = "localhost:8081"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [0|1] (0 for listener, 1 for producer)")
	}

	// Handle interrupt signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// RoomID and player IDs
	roomID := "room-68c26dab-4caf-464e-971a-41fba8f9a3cf"
	listenerID := "listener1"
	producerID := "producer1"

	// Check input argument
	mode := os.Args[1]

	if mode == "0" {
		// Start listener
		startListener(roomID, listenerID)
	} else if mode == "1" {
		// Start producer
		startProducer(roomID, producerID)
	} else {
		log.Fatal("Invalid argument. Use 0 for listener and 1 for producer.")
	}

	// Wait for interrupt signal to terminate all connections
	<-interrupt
	log.Println("Received interrupt signal. Closing all connections.")
}

// Function to create a listener connection
func startListener(roomID, playerID string) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws", RawQuery: "roomID=" + roomID + "&playerID=" + playerID + "&role=listener"}
	log.Printf("Listener %s connecting to %s", playerID, u.String())

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Listener connection error:", err)
	}
	defer conn.Close()

	// Read messages from WebSocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Listener read error:", err)
			break
		}
		log.Printf("Listener %s received: %s", playerID, message)
	}
}

// Function to create a producer connection
func startProducer(roomID, playerID string) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws", RawQuery: "roomID=" + roomID + "&playerID=" + playerID + "&role=pusher"}
	log.Printf("Producer %s connecting to %s", playerID, u.String())

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Producer connection error:", err)
	}
	defer conn.Close()

	// Send data every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			// Send a JSON message
			err := conn.WriteJSON(map[string]string{
				"eventType": "sendData",
				"playerID":  playerID,
				"data":      "Data from " + playerID + " at " + t.String(),
			})
			if err != nil {
				log.Println("Producer write error:", err)
				return
			}
			log.Printf("Producer %s sent data", playerID)
		}
	}
}
