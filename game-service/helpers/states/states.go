package states

import (
	"game-service/models"
)

type GameState interface {
	HandleEvent(event models.GameRoomEvent, gameRoom *models.GameRoomEntity) error
	TransitionToNext(gameRoom *models.GameRoomEntity) GameState
}
