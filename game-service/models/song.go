package models

// Struct to represent a note in the UltraStar file
type Note struct {
	Type      string // Note type (:, *, -, etc.)
	Timestamp int    // Start time in milliseconds
	Duration  int    // Duration in milliseconds
	Pitch     int    // Pitch of the note
	Text      string // Lyric or note text
}

// Struct to represent an UltraStar song
type UltraStarSong struct {
	Artist   string
	Title    string
	MP3      string
	Creator  string
	Genre    string
	Year     string
	Language string
	BPM      float64
	GAP      int
	Notes    []Note
}
