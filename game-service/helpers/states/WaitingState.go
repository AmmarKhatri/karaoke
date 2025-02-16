package states

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"game-service/models"
	"game-service/utils"
	"log"

	"github.com/redis/go-redis/v9"
)

// WaitingState represents the waiting state of the game room
type WaitingState struct{}

func (w *WaitingState) HandleEvent(event models.GameRoomEvent, gameRoom *models.GameRoomEntity) error {
	if event.EventType == "startGame" && event.PlayerID == gameRoom.CreatedBy {
		// Update the game room status
		gameRoom.Status = models.Started
		//load it
		filename := "/app/song.txt"

		song, err := utils.LoadUltraStarFile(filename)
		if err != nil {
			fmt.Printf("Error loading file: %v\n", err)
			return err
		}
		jsonSong, err := json.Marshal(song)
		if err != nil {
			fmt.Printf("Error marshalling the song: %v\n", err)
			return err
		}
		// Send the game specifications to the Redis stream
		gameSpecsEvent := map[string]interface{}{
			"eventType": string(models.GameSpecsEvent),
			"createdBy": "system",
			"data":      string(jsonSong),
			//sending full song data as stringified json before starting song streaming
		}

		_, err = utils.Redis.XAdd(context.Background(), &redis.XAddArgs{
			Stream: "stream:" + gameRoom.ID,
			Values: gameSpecsEvent,
		}).Result()
		if err != nil {
			log.Printf("Failed to send game specifications to Redis: %v", err)
			return errors.New("failed to send game specifications to Redis")
		}

		// Start a goroutine to send song events
		go utils.SendSongEvents(gameRoom, song)

		return nil
	}
	return errors.New("invalid event in waiting state")
}

func (w *WaitingState) TransitionToNext(gameRoom *models.GameRoomEntity) GameState {
	return &StartedState{}
}
