package helpers

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Define a struct for game room events
type GameRoomEvent struct {
	EventType string `json:"eventType"`
	PlayerID  string `json:"playerID"`
	Data      string `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Store WebSocket clients and game rooms
var clients = make(map[*websocket.Conn]bool)
var gameRooms = make(map[string]*redis.Client)
var mu sync.RWMutex // Use RWMutex for safe concurrent access to gameRooms
