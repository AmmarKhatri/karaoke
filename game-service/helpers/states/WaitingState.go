package states

import (
	"errors"
	"game-service/models"
)

// WaitingState represents the waiting state of the game room
type WaitingState struct{}

func (w *WaitingState) HandleEvent(event models.GameRoomEvent, gameRoom *models.GameRoomEntity) error {
	if event.EventType == "startGame" && event.PlayerID == gameRoom.CreatedBy {
		gameRoom.Status = models.Started
		return nil
	}
	return errors.New("invalid event in waiting state")
}

func (w *WaitingState) TransitionToNext(gameRoom *models.GameRoomEntity) GameState {
	return &StartedState{}
}
