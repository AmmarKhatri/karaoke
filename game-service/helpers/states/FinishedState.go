package states

import (
	"errors"
	"game-service/models"
)

// FinishedState represents the finished state of the game room
type FinishedState struct{}

func (f *FinishedState) HandleEvent(event models.GameRoomEvent, gameRoom *models.GameRoomEntity) error {
	return errors.New("no events allowed in finished state")
}

func (f *FinishedState) TransitionToNext(gameRoom *models.GameRoomEntity) GameState {
	return nil // No transition from finished state
}
