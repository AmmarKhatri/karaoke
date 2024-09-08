package helpers

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

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
