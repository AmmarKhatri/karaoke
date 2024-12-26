package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Handle interrupt signals for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// RoomID and player IDs
	roomID := "room-e0b27af1-b018-4cef-a0d4-0d9855f0f5b9"
	producerID := "producer1"

	// Start the high-speed producer
	log.Println("Starting high-speed event producer...")
	go StartFastProducer(roomID, producerID, interrupt)

	// Wait for interrupt signal
	<-interrupt
	log.Println("Producer shutting down gracefully...")
}

var Addr = "localhost:8081"

// StartFastProducer sends events at a high rate to the specified room for 2 seconds.
func StartFastProducer(roomID, producerID string, interrupt chan os.Signal) {
	u := url.URL{Scheme: "ws", Host: Addr, Path: "/ws", RawQuery: "roomID=" + roomID + "&playerID=" + producerID + "&role=phone"}
	log.Printf("Listener %s connecting to %s", producerID, u.String())

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Listener connection error:", err)
	}
	defer conn.Close()
	startTime := time.Now()
	done := make(chan bool)
	go func() {
		<-interrupt
		done <- true
	}()

	// Establish a 2-second time limit
	for {
		select {
		case <-done:
			return
		default:
			// Send JSON event every second
			err := conn.WriteJSON(map[string]string{
				"eventType": "sendData",
				"playerID":  producerID,
				"data":      "Data from " + producerID + " at " + time.Now().String(),
			})
			if err != nil {
				log.Println("Producer write error:", err)
				close(done)
				return
			}
			//log.Printf("Producer %s sent data", producerID)

			if time.Since(startTime) > 1*time.Second {
				log.Println("Stopping event transmission after 1 seconds.")
				close(done)
			}

		}
	}
}
