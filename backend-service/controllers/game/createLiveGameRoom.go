package controllers

import (
	"backend-service/models"
	"backend-service/utils"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func CreateLiveGameRoom(c *gin.Context) {
	var request models.LiveGameRoomRequest

	// Bind the incoming JSON payload to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Defines the Redis queue key according to the game type
	queueKey := "queue:" + request.GameType

	// Player is added to the queue with skill level as metadata
	playerData := map[string]string{
		"skillLevel": request.SkillLevel,
		"createdBy":  request.CreatedBy,
	}
	playerDataJSON, _ := json.Marshal(playerData)

	err := utils.Redis.RPush(context.Background(), queueKey, string(playerDataJSON)).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add request to queue"})
		return
	}
	log.Printf("Request added to queue %s", queueKey)

	// SSE setup
	c.Stream(func(w io.Writer) bool {
		matchFound := false
		var matchedPlayer map[string]string
		var roomID string

		timeout := time.After(50 * time.Second)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				// when timeout reached, remove the request from the queue
				utils.Redis.LRem(context.Background(), queueKey, 1, string(playerDataJSON))
				c.SSEvent("status", gin.H{"status": "timeout", "message": "No match found, please try again later"})
				return false

			case <-ticker.C:
				// Send "waiting" update to the client
				c.SSEvent("status", gin.H{"status": "waiting", "message": "Waiting for game"})

				// find a match
				matchedPlayer, err = utils.FindMatch(context.Background(), queueKey, request.SkillLevel, string(playerDataJSON))
				if err != nil {
					log.Printf("Error finding match: %v", err)
					continue
				}
				if matchedPlayer != nil {
					matchFound = true
					break
				}
			}

			if matchFound {
				break
			}
		}

		if matchFound {
			// Sort the players' IDs for same room ID
			playerIDs := []string{playerData["createdBy"], matchedPlayer["createdBy"]}
			sort.Strings(playerIDs)

			// Generates same room ID using sorted player IDs
			roomIDKey := "room:" + playerIDs[0] + ":" + playerIDs[1]
			roomID, err = utils.Redis.Get(context.Background(), roomIDKey).Result()
			if err == redis.Nil {
				//if Room ID doesn't exist; create
				roomID = "room-" + uuid.New().String()
				// Save the room ID for both players in Redis
				utils.Redis.Set(context.Background(), roomIDKey, roomID, 0)
			} else if err != nil {
				c.SSEvent("status", gin.H{"status": "error", "message": "Failed to create or retrieve room ID"})
				return false
			}

			// Create the game room entity with the same Room ID for both players
			gameRoom := models.GameRoomEntity{
				ID:                    roomID,
				MinPlayers:            2,
				MaxPlayers:            2,
				CreatedBy:             request.CreatedBy,
				Type:                  "live",
				Status:                "waiting",
				ConnectedPlayers:      []string{},
				UnityConnectedPlayers: []string{},
			}

			// Save the game room to Redis
			gameRoomData, _ := json.Marshal(gameRoom)
			err = utils.Redis.Set(context.Background(), roomID, gameRoomData, 0).Err()
			if err != nil {
				c.SSEvent("status", gin.H{"status": "error", "message": "Failed to create game room"})
				return false
			}
			// Create a Redis Stream for the game room
			streamKey := "stream:" + roomID
			_, err = utils.Redis.XAdd(context.Background(), &redis.XAddArgs{
				Stream: streamKey,
				Values: map[string]interface{}{
					"eventType": "gameRoomCreated",
					"createdBy": request.CreatedBy,
					"data":      "Game room has been created",
				},
			}).Result()
			if err != nil {
				c.SSEvent("status", gin.H{"status": "error", "message": "Failed to create game room stream"})
				return false
			}

			// Notify both players about the successful match and the same room ID
			c.SSEvent("status", gin.H{
				"status":  "matched",
				"message": "Game room created successfully",
				"roomID":  roomID,
				"players": []string{playerData["createdBy"], matchedPlayer["createdBy"]},
			})
			return false
		}

		return true
	})
}
