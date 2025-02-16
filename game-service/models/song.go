package models

// Struct to represent a note in the UltraStar file
type Note struct {
	Type      string `json:"type"`      // Note type (:, *, -, etc.)
	Timestamp int    `json:"timestamp"` // Start time in milliseconds
	Duration  int    `json:"duration"`  // Duration in milliseconds
	Pitch     int    `json:"pitch"`     // Pitch of the note
	Text      string `json:"text"`      // Lyric or note text
}

// Struct to represent an UltraStar song
type UltraStarSong struct {
	Artist        string  `json:"artist"`
	Title         string  `json:"title"`
	MP3           string  `json:"mp3"`
	Creator       string  `json:"creator"`
	Genre         string  `json:"genre"`
	Year          string  `json:"year"`
	Language      string  `json:"language"`
	BPM           float64 `json:"bpm"`
	GAP           int     `json:"gap"`
	TotalDuration int     `json:"total_duration"`
	Offset        int     `json:"offset"`
	Notes         []Note  `json:"notes"`
}

type NoteDetails struct {
	Note          Note `json:"note"`
	Offset        int  `json:"offset"`
	TotalDuration int  `json:"total_duration"`
}
