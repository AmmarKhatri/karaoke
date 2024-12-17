package helpers

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

func pushToGameRoom(redisClient *redis.Client, roomID, playerID string, ws *websocket.Conn) {
	defer ws.Close()

	for {
		var event GameRoomEvent
		if err := ws.ReadJSON(&event); err != nil {
			log.Printf("Error reading event: %v", err)
			break
		}

		err, kill := ValidateEvent(event, roomID)
		if err != nil {
			log.Printf("Invalid event: %v", err)
			if kill {
				sendCloseMessage(ws, err.Error())
				return
			}
			continue
		}

		_, err = redisClient.XAdd(context.Background(), &redis.XAddArgs{
			Stream: "stream:" + roomID,
			Values: map[string]interface{}{
				"eventType": event.EventType,
				"createdBy": playerID,
				"data":      event.Data,
			},
		}).Result()
		if err != nil {
			log.Printf("Error pushing event to Redis stream: %v", err)
		}
	}
}
