package states

import (
	"context"
	"errors"
	"fmt"
	"game-service/models"
	"game-service/utils"
	"math"
	"strconv"
)

// StartedState represents the started state of the game room
type StartedState struct{}

func (s *StartedState) HandleEvent(event models.GameRoomEvent, gameRoom *models.GameRoomEntity) error {
	fmt.Println(event)

	if event.EventType == "pauseGame" && event.PlayerID == gameRoom.CreatedBy {
		// VALIDATION OF CREATOR ONLY
		gameRoom.Status = models.Paused
		return nil
	}

	if event.EventType == "endGame" && event.PlayerID == gameRoom.CreatedBy {
		// VALIDATION OF CREATOR ONLY
		gameRoom.Status = models.Finished
		return nil
	}

	if event.EventType == "playerNote" {
		// Convert received pitch to integer
		pitchReceived, err := strconv.Atoi(event.Data.(string))
		if err != nil {
			return err
		}

		// Fetch only the SongNote from Redis
		var songNote models.Note
		err = utils.Get(utils.Redis, "room_scores:"+gameRoom.ID+":note", &songNote)
		if err != nil {
			return err
		}

		// Check if pitch is within range and update the score atomically
		if math.Abs(float64(songNote.Pitch-pitchReceived)) < 3 {
			err := utils.Redis.HIncrBy(context.Background(), "room_scores:"+gameRoom.ID+":scores", event.PlayerID, 10000).Err()
			if err != nil {
				return err
			}
			fmt.Println("Updated score for", event.PlayerID)
		}

		return nil
	}

	return errors.New("invalid event in started state")
}

func (s *StartedState) TransitionToNext(gameRoom *models.GameRoomEntity) GameState {
	return &PausedState{}
}
