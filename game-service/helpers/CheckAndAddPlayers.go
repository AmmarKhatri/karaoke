package helpers

import (
	"fmt"
	"game-service/models"
	"game-service/utils"
	"log"
)

func CheckAndAddPlayer(playerID, roomID, role string) {
	// Create a variable to store the game room data
	var gameRoom models.GameRoomEntity

	// Retrieve the game room from Redis
	err := utils.Get(utils.Redis, roomID, &gameRoom)
	if err != nil {
		return
	}
	fmt.Println(gameRoom)
	// Check if the player is already connected
	if player, exists := gameRoom.ConnectedPlayers[playerID]; exists {
		// If the player exists, persist their points
		points := player.Points

		// Update the connection status based on the role
		if role == "phone" {
			player.PhoneConnected = true
		} else if role == "tv" {
			player.UnityConnected = true
		}

		// Ensure points and other data are persisted
		gameRoom.ConnectedPlayers[playerID] = models.PlayerStats{
			PlayerName:     player.PlayerName, // Persist existing playerName
			SkillLevel:     player.SkillLevel, // Persist existing skillLevel
			Points:         points,            // Persist points
			PhoneConnected: player.PhoneConnected,
			UnityConnected: player.UnityConnected,
		}
		log.Printf("Player %s updated in room %s: %+v", playerID, gameRoom.ID, gameRoom.ConnectedPlayers[playerID])
		return
	}

	// If the player does not exist, add them with default values
	defaultPlayerName := "default" // Replace with actual logic in the future
	defaultSkillLevel := "default" // Replace with actual logic in the future
	gameRoom.ConnectedPlayers[playerID] = models.PlayerStats{
		PlayerName:     defaultPlayerName,
		SkillLevel:     defaultSkillLevel,
		Points:         0,               // Default to 0 points for a new player
		PhoneConnected: role == "phone", // Set PhoneConnected based on role
		UnityConnected: role == "tv",    // Set UnityConnected based on role
	}
	gameRoom.JoinedPlayers++ // Increment the count of joined players

	log.Printf("Player %s added to room %s: %+v", playerID, gameRoom.ID, gameRoom.ConnectedPlayers[playerID])
	utils.Set(utils.Redis, roomID, &gameRoom, 0) // save it
}
