package main

import (
	"log"
	"net/http"

	"game-service/helpers"
	"game-service/utils"
)

func main() {
	http.HandleFunc("/ws", helpers.HandleConnections)
	log.Println("WebSocket server started on :8081")

	// connect to Redis
	utils.ConnectToRedis()

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
