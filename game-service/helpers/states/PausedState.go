package states

import (
	"errors"
	"game-service/models"
)

// PausedState represents the paused state of the game room
type PausedState struct{}

func (p *PausedState) HandleEvent(event models.GameRoomEvent, gameRoom *models.GameRoomEntity) error {
	if event.EventType == "resumeGame" {
		gameRoom.Status = models.Started
		return nil
	}
	return errors.New("invalid event in paused state")
}

func (p *PausedState) TransitionToNext(gameRoom *models.GameRoomEntity) GameState {
	return &StartedState{}
}
