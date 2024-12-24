package clientutils

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Instruction represents the message structure for communication
type Instruction struct {
	Role    string `json:"role"`
	Command string `json:"command"`
}

// ConnectWebSocket connects to the WebSocket server
func ConnectWebSocket(baseURL, id, role string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: baseURL, Path: "/connect-tv/" + id, RawQuery: "role=" + role}
	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected as %s to TV ID %s", role, id)
	return conn, nil
}

// SendInstruction sends an instruction to the WebSocket server
func SendInstruction(conn *websocket.Conn, instruction Instruction) error {
	message, err := json.Marshal(instruction)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	log.Printf("Sent: %+v", instruction)
	return nil
}

// ReceiveMessages listens for messages from the WebSocket server
func ReceiveMessages(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var instruction Instruction
		err = json.Unmarshal(message, &instruction)
		if err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		log.Printf("Received: %+v", instruction)
	}
}
