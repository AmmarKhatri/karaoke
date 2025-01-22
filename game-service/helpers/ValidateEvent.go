package helpers

import (
	"errors"
	"game-service/models"
	"game-service/utils"
)

func ValidateEvent(event models.GameRoomEvent, roomID string) (error, bool) {
	if event.EventType == "" {
		return errors.New("event type is required"), false
	}
	if event.Data == "" {
		return errors.New("event data is required"), false
	}
	// Create a variable to store the game room data
	var gameRoom models.GameRoomEntity

	// Retrieve the game room from Redis
	err := utils.Get(utils.Redis, roomID, &gameRoom)
	if err != nil {
		return errors.New("game room not found"), true
	}
	//start game
	if gameRoom.Status == "waiting" && event.PlayerID == gameRoom.CreatedBy && event.EventType == "startGame" {
		gameRoom.Status = "started"
		err := utils.Set(utils.Redis, roomID, &gameRoom, 0)
		if err != nil {
			return errors.New("game room not found"), true
		}
		return nil, false
	}
	return nil, false
}
