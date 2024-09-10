package helpers

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

func handlePusherEvents(redisClient *redis.Client, roomID, playerID string, ws *websocket.Conn) {
	log.Println("Inside of Push events!!")
	for {
		var event GameRoomEvent

		// Read event from WebSocket
		err := ws.ReadJSON(&event)
		if err != nil {
			log.Printf("Error reading event from WebSocket: %v", err)
			break
		}
		log.Println("Printing event!")
		log.Println(event)
		// Publish event to Redis Stream
		res, err := redisClient.XAdd(context.Background(), &redis.XAddArgs{
			Stream: "stream:" + roomID,
			Values: map[string]interface{}{
				"eventType": event.EventType,
				"playerID":  playerID,
				"data":      event.Data,
			},
		}).Result()
		log.Println("Response: " + res)
		if err != nil {
			log.Printf("Error publishing event to Redis Stream: %v", err)
		}
	}
}
