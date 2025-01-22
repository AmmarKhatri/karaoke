package models

type GameRoomStatus string

const (
	Waiting  GameRoomStatus = "waiting"
	Started  GameRoomStatus = "started"
	Paused   GameRoomStatus = "paused"
	Finished GameRoomStatus = "finished"
)

// GameRoom represents the structure of a game room in Redis
type GameRoomEntity struct {
	ID               string                 `json:"id"`
	MinPlayers       int                    `json:"minPlayers"`
	MaxPlayers       int                    `json:"maxPlayers"`
	CreatedBy        string                 `json:"createdBy"`
	Type             string                 `json:"type"`
	Status           GameRoomStatus         `json:"status"`
	ConnectedPlayers map[string]PlayerStats `json:"connectedPlayers"`
	JoinedPlayers    int                    `json:"joinedPlayers"`
}

type PlayerStats struct {
	SkillLevel     string `json:"skillLevel"`
	Points         int    `json:"points"`
	PhoneConnected bool   `json:"phoneConnected"`
	UnityConnected bool   `json:"unityConnected"`
}
