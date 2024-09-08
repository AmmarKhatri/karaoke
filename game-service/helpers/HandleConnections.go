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
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("WebSocket upgrade failed:", err)
	}
	defer ws.Close()

	// Register the new client
	clients[ws] = true

	// Get roomID and role from query params
	roomID := r.URL.Query().Get("roomID")
	playerID := r.URL.Query().Get("playerID")
	role := r.URL.Query().Get("role") // "listener" or "pusher"

	// Check if the game room exists in Redis (db 0) using utils.Redis
	exists, err := utils.Redis.Exists(context.Background(), roomID).Result()
	if err != nil || exists == 0 {
		log.Printf("Game room %s does not exist", roomID)

		// Send a close frame with a reason to the client
		closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Game room does not exist")
		ws.WriteMessage(websocket.CloseMessage, closeMessage)

		// Close WebSocket connection
		ws.Close()
		return
	}

	// Handle WebSocket pings/pongs
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Lock for reading the map (multiple goroutines can read concurrently)
	mu.RLock()
	clientRedis, existsRoom := gameRooms[roomID]
	mu.RUnlock()

	if !existsRoom {
		// Lock for writing the map (only one goroutine can write at a time)
		mu.Lock()
		clientRedis = utils.Redis
		gameRooms[roomID] = clientRedis
		mu.Unlock()
	}

	// Handle based on role
	if role == "listener" {
		go listenToGameRoom(clientRedis, roomID, ws)
	} else if role == "pusher" {
		go handlePusherEvents(clientRedis, roomID, playerID, ws)
	}
}
