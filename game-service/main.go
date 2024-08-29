package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID   string
	Conn *websocket.Conn
}

type GameRoom struct {
	ID      string
	Players []*Player
}

var (
	upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	gameRooms = make(map[string]*GameRoom) // All active game rooms
)

// Handle WebSocket connections.
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	player := &Player{
		ID:   r.RemoteAddr,
		Conn: conn,
	}

	// Listen for incoming messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		var request map[string]interface{}
		if err := json.Unmarshal(msg, &request); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		switch request["type"] {
		case "create_room":
			createRoom(player)
		case "join_room":
			roomID := request["room_id"].(string)
			joinRoom(player, roomID)
		default:
			log.Println("Unknown request type:", request["type"])
		}
	}
}

// Create a new game room.
func createRoom(player *Player) {
	roomID := fmt.Sprintf("room-%d", len(gameRooms)+1)
	room := &GameRoom{
		ID:      roomID,
		Players: []*Player{player},
	}
	gameRooms[roomID] = room

	response := map[string]string{"type": "room_created", "room_id": roomID}
	player.Conn.WriteJSON(response)
	log.Printf("Room %s created by %s", roomID, player.ID)
}

// Join an existing game room.
func joinRoom(player *Player, roomID string) {
	room, exists := gameRooms[roomID]
	if !exists {
		player.Conn.WriteJSON(map[string]string{"type": "error", "message": "Room does not exist"})
		return
	}

	if len(room.Players) >= 2 {
		player.Conn.WriteJSON(map[string]string{"type": "error", "message": "Room is full"})
		return
	}

	room.Players = append(room.Players, player)
	player.Conn.WriteJSON(map[string]string{"type": "room_joined", "room_id": roomID})
	log.Printf("Player %s joined room %s", player.ID, roomID)

	// Notify all players in the room that the game can start
	for _, p := range room.Players {
		p.Conn.WriteJSON(map[string]string{"type": "start_game"})
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	log.Println("Game Service WebSocket server started on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
