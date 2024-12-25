package helpers

import (
	"encoding/json"
	"game-service/models"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	tvConnections    = make(map[string]*websocket.Conn) // Map to hold TV connections
	phoneConnections = make(map[string]*websocket.Conn) // Map to hold Phone connections
	connectionsLock  sync.Mutex                         // Mutex to ensure thread safety
)

// ConnectTV handles WebSocket connections for TV and Phone
func ConnectTV(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tvID := vars["id"]

	// Extract role from query parameters
	role := r.URL.Query().Get("role")
	if role != "tv" && role != "phone" {
		http.Error(w, "Invalid role. Role must be 'tv' or 'phone'", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	connectionsLock.Lock()
	if role == "tv" {
		// Check if a TV is already connected
		if _, exists := tvConnections[tvID]; exists {
			connectionsLock.Unlock()
			http.Error(w, "TV is already connected to this ID", http.StatusConflict)
			return
		}
		tvConnections[tvID] = ws
	} else if role == "phone" {
		// Check if a Phone is already connected
		if _, exists := phoneConnections[tvID]; exists {
			connectionsLock.Unlock()
			http.Error(w, "Phone is already connected to this ID", http.StatusConflict)
			return
		}
		phoneConnections[tvID] = ws
	}
	connectionsLock.Unlock()

	log.Printf("%s connected to TV ID %s", role, tvID)

	// Handle communication
	handleCommunication(tvID, role, ws)

	// Clean up on disconnect
	connectionsLock.Lock()
	if role == "tv" {
		delete(tvConnections, tvID)
	} else if role == "phone" {
		delete(phoneConnections, tvID)
	}
	connectionsLock.Unlock()
	log.Printf("%s disconnected from TV ID %s", role, tvID)
}

func handleCommunication(tvID string, role string, ws *websocket.Conn) {
	for {
		// Read message from the WebSocket
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Read error for %s: %v", role, err)
			break
		}

		// Parse the incoming message as Instruction
		var instruction models.Instruction
		err = json.Unmarshal(message, &instruction)
		if err != nil {
			log.Printf("Invalid message format from %s: %v", role, err)
			continue
		}

		log.Printf("Received from %s: %+v", role, instruction)

		// Handle "exit" command
		if instruction.Command == "exit" {
			log.Printf("Exit command received from %s. Notifying and closing connections for TV ID %s.", role, tvID)

			// Notify the other device
			connectionsLock.Lock()
			var targetConn *websocket.Conn
			if role == "tv" {
				targetConn = phoneConnections[tvID]
			} else if role == "phone" {
				targetConn = tvConnections[tvID]
			}
			if targetConn != nil {
				targetConn.WriteJSON(models.Instruction{
					Role:    role,
					Command: "exit",
				})
			}
			// Close both connections
			if conn, exists := tvConnections[tvID]; exists {
				conn.Close()
				delete(tvConnections, tvID)
			}
			if conn, exists := phoneConnections[tvID]; exists {
				conn.Close()
				delete(phoneConnections, tvID)
			}
			connectionsLock.Unlock()
			break
		}

		// Handle "disconnect" command
		if instruction.Command == "disconnect" && role == "phone" {
			log.Printf("Disconnect command received from phone for TV ID %s. Notifying TV and closing phone connection.", tvID)

			// Notify the TV
			connectionsLock.Lock()
			if conn, exists := tvConnections[tvID]; exists {
				conn.WriteJSON(models.Instruction{
					Role:    role,
					Command: "disconnect",
				})
			}
			// Close the phone connection
			if conn, exists := phoneConnections[tvID]; exists {
				conn.Close()
				delete(phoneConnections, tvID)
			}
			connectionsLock.Unlock()
			break
		}

		// Forward message to the other device
		connectionsLock.Lock()
		if role == "tv" {
			// TV messages are sent to the phone
			targetConn := phoneConnections[tvID]
			if targetConn != nil && targetConn != ws {
				targetConn.WriteJSON(instruction)
			}
		} else if role == "phone" {
			// Phone messages are sent to the TV
			targetConn := tvConnections[tvID]
			if targetConn != nil && targetConn != ws {
				targetConn.WriteJSON(instruction)
			}
		}
		connectionsLock.Unlock()
	}
}
