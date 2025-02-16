package controllers

import (
	"backend-service/models"
	"backend-service/utils"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FetchGameRoomScores retrieves and returns the scores of a game room
func FetchGameRoomScores(c *gin.Context) {
	gameRoomID := c.Param("id") // Get game room ID from URL parameters

	// Fetch the game room from Redis
	var gameRoom models.GameRoomEntity
	err := utils.Get(utils.Redis, gameRoomID, &gameRoom)
	if err != nil {
		log.Printf("Game room %s not found: %v", gameRoomID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Game room not found"})
		return
	}

	// Retrieve scores from Redis
	scores, err := utils.Redis.HGetAll(context.Background(), "room_scores:"+gameRoomID+":scores").Result()
	if err != nil {
		log.Printf("Failed to retrieve scores for game room %s: %v", gameRoomID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scores"})
		return
	}

	// Return scores in JSON response
	c.JSON(http.StatusOK, gin.H{
		"gameRoomID": gameRoomID,
		"scores":     scores,
	})
}
