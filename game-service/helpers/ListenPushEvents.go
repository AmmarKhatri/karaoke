package helpers

import (
	"context"
	"game-service/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// TV function watching over the gameroom
func listenToGameRoom(redisClient *redis.Client, roomID string, ws *websocket.Conn) {
	defer ws.Close()

	for {
		entries, err := redisClient.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{"stream:" + roomID, "$"},
			Block:   time.Duration(0 * time.Millisecond), // No blocking
		}).Result()

		if err != nil && err != redis.Nil {
			log.Printf("Error reading from Redis stream: %v", err)
			break
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
					return
				}
			}
		}
	}
}
