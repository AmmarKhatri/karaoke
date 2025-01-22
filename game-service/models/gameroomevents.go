package models

type EventType string

const (
	StartGame          EventType = "startGame"
	PauseGame          EventType = "pauseGame"
	ResumeGame         EventType = "resumeGame"
	PlayerConnected    EventType = "playerConnected"
	playerDisconnected EventType = "playerDisconnected"
	EndGame            EventType = "endGame"
)

// Define a struct for game room events
type GameRoomEvent struct {
	EventType EventType `json:"eventType"`
	PlayerID  string    `json:"playerID"`
	Data      any       `json:"data"`
}
