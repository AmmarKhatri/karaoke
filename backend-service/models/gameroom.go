package models

// GameRoomRequest represents the request payload for creating a game room
type GameRoomRequest struct {
	MinPlayers int    `json:"minPlayers" binding:"required"`
	MaxPlayers int    `json:"maxPlayers" binding:"required"`
	CreatedBy  string `json:"createdBy" binding:"required"`
}

// GameRoom represents the structure of a game room in Redis
type GameRoomEntity struct {
	ID                    string   `json:"id"`
	MinPlayers            int      `json:"minPlayers"`
	MaxPlayers            int      `json:"maxPlayers"`
	CreatedBy             string   `json:"createdBy"`
	Type                  string   `json:"type"`
	Status                string   `json:"status"`
	ConnectedPlayers      []string `json:"connectedPlayers"`
	UnityConnectedPlayers []string `json:"unityConnectedPlayers"`
}
type LiveGameRoomRequest struct {
	SkillLevel string `json:"skillLevel" binding:"required"`
	GameType   string `json:"gameType" binding:"required"`
	CreatedBy  string `json:"createdBy" binding:"required"`
}