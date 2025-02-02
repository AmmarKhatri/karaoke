package states

import (
	"errors"
	"game-service/models"
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
		//TODO
	}
	return errors.New("invalid event in started state")
}

func (s *StartedState) TransitionToNext(gameRoom *models.GameRoomEntity) GameState {
	return &PausedState{}
}
