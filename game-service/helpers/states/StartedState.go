package states

import (
	"errors"
	"game-service/models"
	"game-service/utils"
	"math"
	"strconv"
)

// StartedState represents the started state of the game room
type StartedState struct{}

func (s *StartedState) HandleEvent(event models.GameRoomEvent, gameRoom *models.GameRoomEntity) error {
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
		// add player notes
		pitchReceived, err := strconv.Atoi(event.Data.(string))
		if err != nil {
			return err
		}
		// Get game room scores
		roomScores := models.RoomScores{}
		err = utils.Get(utils.Redis, "room_scores:"+gameRoom.ID, roomScores)
		if err != nil {
			return err
		}
		// Do calculation against player and add the score
		if math.Abs(float64(roomScores.SongNote.Pitch-pitchReceived)) < 3 {
			playerScore := roomScores.ConnectedPlayers[event.PlayerID]
			playerScore.Points += 10000
			roomScores.ConnectedPlayers[event.PlayerID] = playerScore
		}
		// Save the scores
		err = utils.Set(utils.Redis, "room_scores:"+gameRoom.ID, roomScores, 0)
		if err != nil {
			return err
		}

	}
	return errors.New("invalid event in started state")
}

func (s *StartedState) TransitionToNext(gameRoom *models.GameRoomEntity) GameState {
	return &PausedState{}
}
