package utils

import (
	"bufio"
	"fmt"
	"game-service/models"
	"os"
	"strconv"
	"strings"
)

// LoadUltraStarFile loads and parses an UltraStar Deluxe song file
func LoadUltraStarFile(filename string) (*models.UltraStarSong, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	song := &models.UltraStarSong{}
	scanner := bufio.NewScanner(file)

	var beatDuration float64 // Beat duration in milliseconds

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line[1:], ":", 2)
			if len(parts) != 2 {
				continue
			}

			key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
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
				if err != nil || bpm <= 0 {
					return nil, fmt.Errorf("invalid BPM value: %s", value)
				}
				song.BPM = bpm
				beatDuration = 60000 / bpm // Beat duration in milliseconds
			case "GAP":
				gap, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("invalid GAP value: %s", value)
				}
				song.GAP = gap
			}
		} else if len(line) > 0 && (line[0] == ':' || line[0] == '*' || line[0] == '-' || line[0] == 'F') {
			note, err := parseNoteLine(line, beatDuration)
			if err != nil {
				return nil, fmt.Errorf("error parsing note line '%s': %w", line, err)
			}
			song.Notes = append(song.Notes, note)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Debugging: Calculate and print total song duration
	if len(song.Notes) > 0 {
		lastNote := song.Notes[len(song.Notes)-1]
		totalDuration := (float64(lastNote.Timestamp+lastNote.Duration) * beatDuration) + float64(song.GAP)
		fmt.Printf("Total Song Duration: %.2f seconds\n", totalDuration/1000)
	}

	return song, nil
}
