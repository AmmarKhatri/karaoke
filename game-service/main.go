package main

import (
	"fmt"
	"log"
	"net/http"

	"game-service/helpers"
	"game-service/utils"
)

func main() {
	//load it
	filename := "/app/song.txt"

	song, err := helpers.LoadUltraStarFile(filename)
	if err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		return
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
	http.HandleFunc("/ws", helpers.HandleConnections)
	http.HandleFunc("/connect-tv/", helpers.ConnectTV)
	log.Println("WebSocket server started on :8081")

	// connect to Redis
	utils.ConnectToRedis()

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
