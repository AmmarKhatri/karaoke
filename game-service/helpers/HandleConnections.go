package helpers

import (
	"context"
	"game-service/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	roomID := r.URL.Query().Get("roomID")
	playerID := r.URL.Query().Get("playerID")
	role := r.URL.Query().Get("role") // "tv" or "phone"

	exists, err := utils.Redis.Exists(context.Background(), roomID).Result()
	if err != nil || exists == 0 {
		log.Printf("Game room %s does not exist", roomID)
		sendCloseMessage(ws, "Game room does not exist")
		return
	}

	//Check if player can join or already exists
	CheckAndAddPlayer(playerID, roomID, role)

	done := make(chan struct{}) // Signal to stop heartbeat
	setupHeartbeat(ws, done)

	if role == "tv" {
		log.Printf("TV connected: PlayerID: %s, RoomID: %s", playerID, roomID)
		listenToGameRoom(utils.Redis, roomID, ws)
	} else if role == "phone" {
		log.Printf("Phone connected: PlayerID: %s, RoomID: %s", playerID, roomID)
		pushToGameRoom(utils.Redis, roomID, playerID, ws)
	} else {
		log.Printf("Invalid role: %s", role)
		sendCloseMessage(ws, "Invalid role")
	}

	close(done) // Stop heartbeat when connection handler exits
}

// sendCloseMessage sends a close frame with a custom reason
func sendCloseMessage(ws *websocket.Conn, reason string) {
	closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason)
	ws.WriteMessage(websocket.CloseMessage, closeMessage)
}

// setupHeartbeat ensures the connection stays alive with pings/pongs
func setupHeartbeat(ws *websocket.Conn, done chan struct{}) {
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	go func() {
		ticker := time.NewTicker(50 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				log.Println("Stopping heartbeat: Connection closed")
				return // Stop the heartbeat goroutine
			case <-ticker.C:
				if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Heartbeat ping failed: %v", err)
					ws.Close()
					return
				}
			}
		}
	}()
}
