package scripts

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Function to create a producer connection
func StartProducer(roomID, playerID string) {
	u := url.URL{Scheme: "ws", Host: Addr, Path: "/ws", RawQuery: "roomID=" + roomID + "&playerID=" + playerID + "&role=pusher"}
	log.Printf("Producer %s connecting to %s", playerID, u.String())

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Producer connection error:", err)
	}
	defer conn.Close()

	// Channel to signal when the connection should be closed
	done := make(chan struct{})

	// Goroutine to listen for messages from the server
	go func() {
		for {
			// Read messages from the WebSocket server
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Printf("Server closed connection: %v", err)
					close(done) // Signal to close the ticker and exit
					return
				}
				log.Println("Error reading WebSocket message:", err)
				close(done)
				return
			}
			log.Printf("Producer %s received: %s", playerID, message)
		}
	}()

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
				close(done)
				return
			}
			log.Printf("Producer %s sent data", playerID)
		case <-done: // Exit the loop if the connection is closed
			log.Println("Exiting producer due to server close event")
			return
		}
	}
}
