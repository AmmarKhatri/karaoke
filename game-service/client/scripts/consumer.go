package scripts

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

var Addr = "localhost:8081"

// Function to create a listener connection
func StartListener(roomID, playerID string) {
	u := url.URL{Scheme: "ws", Host: Addr, Path: "/ws", RawQuery: "roomID=" + roomID + "&playerID=" + playerID + "&role=listener"}
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
