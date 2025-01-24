package utils

import (
	"bufio"
	"fmt"
	"game-service/models"
	"os"
	"strconv"
	"strings"
)

// Function to load and parse an UltraStar Deluxe song file
func LoadUltraStarFile(filename string) (*models.UltraStarSong, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	song := &models.UltraStarSong{}
	scanner := bufio.NewScanner(file)

	var beatDuration float64 // Calculate beat duration in milliseconds

	for scanner.Scan() {
		line := scanner.Text()

		// Check for key-value pairs
		if strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line[1:], ":", 2)
			if len(parts) == 2 {
				key, value := parts[0], strings.TrimSpace(parts[1])
				switch key {
				case "ARTIST":
					song.Artist = value
				case "TITLE":
					song.Title = value
				case "MP3":
					song.MP3 = value
				case "CREATOR":
					song.Creator = value
				case "GENRE":
					song.Genre = value
				case "YEAR":
					song.Year = value
				case "LANGUAGE":
					song.Language = value
				case "BPM":
					bpm, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return nil, fmt.Errorf("invalid BPM value: %s", value)
					}
					song.BPM = bpm
					beatDuration = 60000 / bpm
				case "GAP":
					gap, err := strconv.Atoi(value)
					if err != nil {
						return nil, fmt.Errorf("invalid GAP value: %s", value)
					}
					song.GAP = gap
				}
			}
		} else if len(line) > 0 && (line[0] == ':' || line[0] == '*' || line[0] == '-' || line[0] == 'F') {
			// Parse the note line into the Note struct
			note, err := parseNoteLine(line, beatDuration)
			if err != nil {
				return nil, err
			}
			song.Notes = append(song.Notes, note)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return song, nil
}
