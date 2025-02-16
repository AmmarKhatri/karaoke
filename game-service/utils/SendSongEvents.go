package utils

import (
	"context"
	"encoding/json"
	"game-service/models"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func SendSongEvents(gameRoom *models.GameRoomEntity, song *models.UltraStarSong) {
	stream := "stream:" + gameRoom.ID

	log.Printf("Starting song playback: %s by %s", song.Title, song.Artist)

	// Store the initial song note separately
	err := Set(Redis, "room_scores:"+gameRoom.ID+":note", models.Note{}, 0)
	if err != nil {
		log.Printf("Failed to store initial song note in Redis: %v", err)
		return
	}

	// Initialize player scores in Redis Hash
	for playerID := range gameRoom.ConnectedPlayers {
		_, err := Redis.HSet(context.Background(), "room_scores:"+gameRoom.ID+":scores", playerID, 0).Result()
		if err != nil {
			log.Printf("Failed to initialize score for player %s: %v", playerID, err)
			return
		}
	}

	var lastTimestamp int

	for i, note := range song.Notes {
		var delay int
		if i == 0 {
			delay = song.GAP
		} else {
			delay = note.Timestamp - lastTimestamp
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)

		// Note details struct
		noteDetails := models.NoteDetails{
			Note:          note,
			Offset:        song.Offset,
			TotalDuration: song.TotalDuration,
		}

		// Update only the song note in Redis
		err = Set(Redis, "room_scores:"+gameRoom.ID+":note", noteDetails, 0)
		if err != nil {
			log.Printf("Failed to update song note in Redis: %v", err)
			return
		}

		// Send stringified JSON
		jsonNote, err := json.Marshal(noteDetails)
		if err != nil {
			log.Printf("Failed to marshal song details: %v", err)
			return
		}

		// Send the note event to Redis stream
		noteEvent := map[string]interface{}{
			"eventType": string(models.SongNote),
			"createdBy": "system",
			"data":      string(jsonNote),
		}

		_, err = Redis.XAdd(context.Background(), &redis.XAddArgs{
			Stream: stream,
			Values: noteEvent,
		}).Result()
		if err != nil {
			log.Printf("Failed to send note to Redis: %v", err)
			return
		}

		lastTimestamp = note.Timestamp
	}

	// Pull the latest game room entity
	latestGameRoom := models.GameRoomEntity{}
	err = Get(Redis, gameRoom.ID, &latestGameRoom)
	if err != nil {
		log.Printf("Failed to retrieve game room from Redis: %v", err)
		return
	}

	// Retrieve final scores from Redis
	scores, err := Redis.HGetAll(context.Background(), "room_scores:"+gameRoom.ID+":scores").Result()
	if err != nil {
		log.Printf("Failed to retrieve final scores: %v", err)
		return
	}

	// Log final scores
	log.Printf("Final scores for game room %s: %+v", gameRoom.ID, scores)

	// Update player scores in the latestGameRoom object
	for playerID, scoreStr := range scores {
		score, err := strconv.Atoi(scoreStr) // Convert score from string to int
		if err != nil {
			log.Printf("Failed to convert score for player %s: %v", playerID, err)
			continue
		}

		// Fetch the player object from the map
		player, exists := latestGameRoom.ConnectedPlayers[playerID]
		if !exists {
			log.Printf("Player %s not found in ConnectedPlayers map", playerID)
			continue
		}

		// Update player points
		player.Points = score

		// Reassign updated player back to the map
		latestGameRoom.ConnectedPlayers[playerID] = player
	}

	// Mark the game as finished
	latestGameRoom.Status = models.Finished

	// Save updated game room with scores back to Redis
	err = Set(Redis, gameRoom.ID, latestGameRoom, 0)
	if err != nil {
		log.Printf("Failed to update game state to finished: %v", err)
		return
	}

	// Convert final scores map to JSON for Redis stream
	scoreJSON, err := json.Marshal(scores)
	if err != nil {
		log.Printf("Failed to marshal final scores: %v", err)
		return
	}

	// Send scores as a message to the Redis stream
	scoreEvent := map[string]interface{}{
		"eventType": "gameScores",
		"createdBy": "system",
		"data":      string(scoreJSON),
	}

	_, err = Redis.XAdd(context.Background(), &redis.XAddArgs{
		Stream: stream,
		Values: scoreEvent,
	}).Result()
	if err != nil {
		log.Printf("Failed to send final scores to Redis: %v", err)
		return
	}

	// Send system event for game completion
	gameFinishedEvent := map[string]interface{}{
		"eventType": "endGame",
		"createdBy": "system",
		"data":      "Game has finished.",
	}

	_, err = Redis.XAdd(context.Background(), &redis.XAddArgs{
		Stream: stream,
		Values: gameFinishedEvent,
	}).Result()
	if err != nil {
		log.Printf("Failed to send game finished event: %v", err)
		return
	}

	log.Printf("Song playback completed: %s by %s", song.Title, song.Artist)
}
