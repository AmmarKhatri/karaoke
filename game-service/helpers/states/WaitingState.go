package states

import (
	"context"
	"errors"
	"fmt"
	"game-service/models"
	"game-service/utils"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// WaitingState represents the waiting state of the game room
type WaitingState struct{}

func (w *WaitingState) HandleEvent(event models.GameRoomEvent, gameRoom *models.GameRoomEntity) error {
	if event.EventType == "startGame" && event.PlayerID == gameRoom.CreatedBy {
		// Update the game room status
		gameRoom.Status = models.Started
		//load it
		filename := "/app/song.txt"

		song, err := utils.LoadUltraStarFile(filename)
		if err != nil {
			fmt.Printf("Error loading file: %v\n", err)
			return err
		}

		// Print song details
		fmt.Printf("Artist: %s\n", song.Artist)
		fmt.Printf("Title: %s\n", song.Title)
		fmt.Printf("BPM: %.2f\n", song.BPM)
		fmt.Printf("Notes:\n")
		for _, note := range song.Notes {
			fmt.Printf("Type: %s, Timestamp: %d ms, Duration: %d ms, Pitch: %d, Text: %s\n",
				note.Type, note.Timestamp, note.Duration, note.Pitch, note.Text)
		}
		// Send the game specifications to the Redis stream
		gameSpecsEvent := map[string]interface{}{
			"eventType": "gameSpecifications",
			"createdBy": "system",
			"data": fmt.Sprintf("Game started with artist: %s, title: %s, bpm: %.2f",
				song.Artist, song.Title, song.BPM),
		}

		_, err = utils.Redis.XAdd(context.Background(), &redis.XAddArgs{
			Stream: "stream:" + gameRoom.ID,
			Values: gameSpecsEvent,
		}).Result()
		if err != nil {
			log.Printf("Failed to send game specifications to Redis: %v", err)
			return errors.New("failed to send game specifications to Redis")
		}

		// Start a goroutine to send song events
		go sendSongEvents(gameRoom, song)

		return nil
	}
	return errors.New("invalid event in waiting state")
}

func (w *WaitingState) TransitionToNext(gameRoom *models.GameRoomEntity) GameState {
	return &StartedState{}
}
func sendSongEvents(gameRoom *models.GameRoomEntity, song *models.UltraStarSong) {
	stream := "stream:" + gameRoom.ID

	log.Printf("Starting song playback: %s by %s", song.Title, song.Artist)

	for _, note := range song.Notes {
		// Simulate playback by waiting for the note's timestamp
		time.Sleep(time.Duration(note.Timestamp) * time.Millisecond)

		// Send the note as an event to Redis
		noteEvent := map[string]interface{}{
			"eventType": "songNote",
			"createdBy": "system",
			"data": fmt.Sprintf("Type: %s, Timestamp: %d, Duration: %d, Pitch: %d, Text: %s",
				note.Type, note.Timestamp, note.Duration, note.Pitch, note.Text),
		}

		_, err := utils.Redis.XAdd(context.Background(), &redis.XAddArgs{
			Stream: stream,
			Values: noteEvent,
		}).Result()
		if err != nil {
			log.Printf("Failed to send note to Redis: %v", err)
			return
		}
	}

	log.Printf("Song playback completed: %s by %s", song.Title, song.Artist)
}
