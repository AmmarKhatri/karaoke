package scripts

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// // Listen for messages from the server
//
//	go func() {
//		for {
//			_, message, err := conn.ReadMessage()
//			if err != nil {
//				log.Printf("Producer read error: %v", err)
//				close(done)
//				return
//			}
//			log.Printf("Producer %s received: %s", playerID, message)
//		}
//	}()
//
// StartProducer handles the producer role
// StartProducer handles the producer role
func StartProducer(roomID, playerID string, interrupt chan os.Signal) {
	u := url.URL{Scheme: "ws", Host: Addr, Path: "/ws", RawQuery: "roomID=" + roomID + "&playerID=" + playerID + "&role=phone"}
	log.Printf("Producer %s connecting to %s", playerID, u.String())

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Producer connection error:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	// Ticker to send events every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Timers
	stopSendingEvents := time.After(22 * time.Second) // Stop sending events
	gracefulShutdown := time.After(40 * time.Second)  // Send CloseMessage and shut down

	// Main loop
	for {
		select {
		case t := <-ticker.C:
			// Send JSON event every second
			err := conn.WriteJSON(map[string]string{
				"eventType": "sendData",
				"playerID":  playerID,
				"data":      "Data from " + playerID + " at " + t.String(),
			})
			if err != nil {
				log.Println("Producer write error:", err)
				close(done)
				return
			}
			log.Printf("Producer %s sent data", playerID)

		case <-stopSendingEvents:
			// Stop the ticker after 22 seconds
			log.Println("Stopping event transmission after 22 seconds.")
			ticker.Stop()

		case <-gracefulShutdown:
			// Send CloseMessage after 40 seconds
			log.Println("Sending CloseMessage and shutting down producer.")
			closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Producer finished work")
			if err := conn.WriteMessage(websocket.CloseMessage, closeMessage); err != nil {
				log.Printf("Error sending close message: %v", err)
			}
			interrupt <- os.Interrupt
			return

		case <-done:
			// Handle connection close
			log.Println("Exiting producer due to connection close.")
			interrupt <- os.Interrupt
			return
		}
	}
}
