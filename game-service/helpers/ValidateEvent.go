package helpers

import (
	"errors"
	"fmt"
	"game-service/helpers/states"
	"game-service/models"
	"game-service/utils"
)

func ValidateEvent(event models.GameRoomEvent, roomID string) (error, bool) {
	if event.EventType == "" {
		return errors.New("event type is required"), false
	}
	// if event.Data == "" {
	// 	return errors.New("event data is required"), false
	// }

	// Create a variable to store the game room data
	var gameRoom models.GameRoomEntity

	// Retrieve the game room from Redis
	err := utils.Get(utils.Redis, roomID, &gameRoom)
	if err != nil {
		return errors.New("game room not found"), true
	}

	// Determine the current state based on the game room's status
	var currentState states.GameState
	fmt.Println("Current Status: ", gameRoom.Status)
	switch gameRoom.Status {
	case models.Waiting:
		currentState = &states.WaitingState{}
	case models.Started:
		currentState = &states.StartedState{}
	case models.Paused:
		currentState = &states.PausedState{}
	case models.Finished:
		currentState = &states.FinishedState{}
	default:
		return errors.New("invalid game room status"), true
	}

	// Use the current state to handle the event
	err = currentState.HandleEvent(event, &gameRoom)
	if err != nil {
		return err, false
	}

	// Save the updated game room state back to Redis
	err = utils.Set(utils.Redis, roomID, &gameRoom, 0)
	if err != nil {
		return errors.New("failed to update game room"), true
	}

	return nil, false
}
