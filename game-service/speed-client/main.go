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
	roomID := "room-8dacefa4-1d2c-4bfb-bc25-641b2fc0f509"
	producerID := "producer1"
	listenerID := "listener1"
	// Start the high-speed producer
	log.Println("Starting high-speed event producer...")
	go StartFastListener(roomID, listenerID, interrupt)
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

// StartFastListener counts the number of events received via WebSocket and logs every 2 seconds.
func StartFastListener(roomID, listenerID string, interrupt chan os.Signal) {
	u := url.URL{Scheme: "ws", Host: "localhost:8081", Path: "/ws", RawQuery: "roomID=" + roomID + "&playerID=" + listenerID + "&role=tv"}
	log.Printf("Listener %s connecting to %s", listenerID, u.String())

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Listener connection error:", err)
	}
	defer conn.Close()

	eventCount := 0
	done := make(chan bool)
	go func() {
		<-interrupt
		done <- true
	}()

	// Log the event count every 2 seconds
	go func() {
		for {
			time.Sleep(2 * time.Second)
			log.Printf("Listener %s has received %d events so far.", listenerID, eventCount)
		}
	}()

	for {
		select {
		case <-done:
			log.Printf("Listener %s received %d events.", listenerID, eventCount)
			return
		default:
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Listener read error:", err)
				return
			}
			eventCount++
			//log.Printf("Listener %s received event: %s", listenerID, message)
		}
	}
}
