package helpers

import (
	"context"
	"game-service/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

func listenToGameRoom(redisClient *redis.Client, roomID string, ws *websocket.Conn) {
	defer ws.Close()

	// Create a channel to detect WebSocket disconnects
	done := make(chan struct{})

	// Start a goroutine to listen for WebSocket disconnects
	go func() {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		ws.SetPongHandler(func(string) error {
			ws.SetReadDeadline(time.Now().Add(60 * time.Second))
			return nil
		})

		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				log.Printf("WebSocket disconnected: %v", err)
				close(done) // Signal the main function to stop
				return
			}
		}
	}()

	for {
		select {
		case <-done: // WebSocket has disconnected, break the loop
			log.Printf("Stopping listenToGameRoom for room %s due to WebSocket disconnect", roomID)
			updatePlayerConnectionStatus(roomID, "Unknown", "tv", false) // Mark player as disconnected
			return

		default:
			// Read from Redis stream
			entries, err := redisClient.XRead(context.Background(), &redis.XReadArgs{
				Streams: []string{"stream:" + roomID, "$"},
				Block:   time.Duration(500 * time.Millisecond), // Small timeout to allow checks
			}).Result()

			if err != nil && err != redis.Nil {
				log.Printf("Error reading from Redis stream: %v", err)
				return
			}

			// Forward events to WebSocket
			for _, entry := range entries {
				for _, message := range entry.Messages {
					event := models.GameRoomEvent{
						EventType: message.Values["eventType"].(models.EventType),
						PlayerID:  message.Values["createdBy"].(string),
						Data:      message.Values["data"],
					}
					if err := ws.WriteJSON(event); err != nil {
						log.Printf("Error sending event to WebSocket: %v", err)
						close(done) // Signal disconnection
						return
					}
				}
			}
		}
	}
}
