package controllers

import (
	"backend-service/models"
	"backend-service/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func CreateLocalGameRoom(c *gin.Context) {
	var request models.GameRoomRequest

	// Bind JSON payload to struct and check for errors
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Generate game room ID
	roomID := "room-" + uuid.New().String()

	// Create a new GameRoom instance
	gameRoom := models.GameRoomEntity{
		ID:                    roomID,
		MinPlayers:            request.MinPlayers,
		MaxPlayers:            request.MaxPlayers,
		CreatedBy:             request.CreatedBy,
		Type:                  "local",
		Status:                "waiting",  // Initial status set to "waiting"
		ConnectedPlayers:      []string{}, // Initialize empty list for connected players
		UnityConnectedPlayers: []string{}, // Initialize empty list for Unity connected players
	}

	// Save the game room to Redis (key-value store for room info)
	err := utils.Set(utils.Redis, gameRoom.ID, gameRoom, 24*time.Hour) // Set expiration as needed
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create game room"})
		return
	}

	// Create a Redis Stream for the game room (stream to handle events)
	streamKey := "stream:" + roomID
	res, err := utils.Redis.XAdd(c, &redis.XAddArgs{
		Stream: streamKey,
		Values: map[string]interface{}{
			"eventType": "gameRoomCreated",
			"message":   "Game room has been created",
		},
	}).Result()
	log.Println(res)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create game room stream"})
		return
	}

	// Respond with success
	c.JSON(200, gin.H{"message": "Game room and stream created successfully", "roomID": gameRoom.ID})
}
