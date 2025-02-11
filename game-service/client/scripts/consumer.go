package scripts

import (
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

// StartListener handles the listener role
func StartListener(roomID, playerID string, interrupt chan os.Signal) {

	u := url.URL{Scheme: "ws", Host: Addr, Path: "/ws", RawQuery: "roomID=" + roomID + "&playerID=" + playerID + "&role=tv"}
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
			log.Printf("Listener read error: %v", err)
			interrupt <- os.Interrupt
			break
		}
		log.Printf("Listener %s received: %s", playerID, message)
	}
}
