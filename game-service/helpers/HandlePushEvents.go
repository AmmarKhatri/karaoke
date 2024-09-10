package helpers

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

func handlePusherEvents(redisClient *redis.Client, roomID, playerID string, ws *websocket.Conn) {
	log.Println("Inside of Push events!!")
	//defer ws.Close()
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
		// Validate event before pushing to stream
		err, kill := ValidateEvent(event, roomID)
		if err != nil {
			log.Printf("Invalid event: %v", err)
			continue
		}
		// Kill if breaking event
		if kill {
			closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error())
			ws.WriteMessage(websocket.CloseMessage, closeMessage)
			ws.Close()
			return
		}
		// Publish event to Redis Stream
		res, err := redisClient.XAdd(context.Background(), &redis.XAddArgs{
			Stream: "stream:" + roomID,
			Values: map[string]interface{}{
				"eventType": event.EventType,
				"createdBy": playerID,
				"data":      event.Data,
			},
		}).Result()
		log.Println("Response: " + res)
		if err != nil {
			log.Printf("Error publishing event to Redis Stream: %v", err)
		}
	}
}
