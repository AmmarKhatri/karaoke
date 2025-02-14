package models

type RoomScores struct {
	SongNote         Note                   `json:"songNote"`
	ConnectedPlayers map[string]PlayerStats `json:"connectedPlayers"`
}
