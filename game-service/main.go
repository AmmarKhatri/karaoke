package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"game-service/utils"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Define a struct for game room events
type GameRoomEvent struct {
	EventType string `json:"eventType"`
	PlayerID  string `json:"playerID"`
	Data      string `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Store WebSocket clients and game rooms
var clients = make(map[*websocket.Conn]bool)
var gameRooms = make(map[string]*redis.Client)
var mu sync.RWMutex // Use RWMutex for safe concurrent access to gameRooms

func main() {
	http.HandleFunc("/ws", handleConnections)
	log.Println("WebSocket server started on :8081")
	utils.ConnectToRedis() // connect to Redis
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func handleConnections(w http.ResponseWriter, r *http.Request) {
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
		ws.WriteMessage(websocket.TextMessage, []byte("Game room does not exist"))
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

func listenToGameRoom(redisClient *redis.Client, roomID string, ws *websocket.Conn) {
	// Send regular ping messages to keep the connection alive
	ticker := time.NewTicker(50 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping: %v", err)
				return
			}
		default:
			// Set a blocking read timeout for Redis Streams (10 seconds)
			entries, err := redisClient.XRead(context.Background(), &redis.XReadArgs{
				Streams: []string{"stream:" + roomID, "$"},
				Block:   10000, // Block for 10 seconds, allowing reconnection if no data
				Count:   1,
			}).Result()

			if err != nil && err != redis.Nil {
				log.Printf("Error reading from Redis Stream: %v", err)
				return
			}

			// Only process if there is data
			if len(entries) > 0 {
				for _, entry := range entries[0].Messages {
					eventType := entry.Values["eventType"].(string)
					playerID := entry.Values["playerID"].(string)
					data := entry.Values["data"].(string)

					// Send event to the WebSocket client
					event := GameRoomEvent{
						EventType: eventType,
						PlayerID:  playerID,
						Data:      data,
					}
					err := ws.WriteJSON(event)
					if err != nil {
						log.Printf("Error sending event to client: %v", err)
						return
					}
				}
			}
		}
	}
}

func handlePusherEvents(redisClient *redis.Client, roomID, playerID string, ws *websocket.Conn) {
	for {
		var event GameRoomEvent

		// Read event from WebSocket
		err := ws.ReadJSON(&event)
		if err != nil {
			log.Printf("Error reading event from WebSocket: %v", err)
			break
		}

		// Publish event to Redis Stream
		_, err = redisClient.XAdd(context.Background(), &redis.XAddArgs{
			Stream: "stream:" + roomID,
			Values: map[string]interface{}{
				"eventType": event.EventType,
				"playerID":  playerID,
				"data":      event.Data,
			},
		}).Result()

		if err != nil {
			log.Printf("Error publishing event to Redis Stream: %v", err)
		}
	}
}
