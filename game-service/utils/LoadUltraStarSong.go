package utils

import (
	"bufio"
	"fmt"
	"game-service/models"
	"os"
	"strconv"
	"strings"
)

func LoadUltraStarFile(filename string) (*models.UltraStarSong, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	song := &models.UltraStarSong{}
	scanner := bufio.NewScanner(file)

	var bpm float64 // BPM to be used in conversion

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
				parsedBPM, err := strconv.ParseFloat(value, 64)
				if err != nil || parsedBPM <= 0 {
					return nil, fmt.Errorf("invalid BPM value: %s", value)
				}
				song.BPM = parsedBPM
				bpm = parsedBPM // Use for tick conversion
			case "GAP":
				gap, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("invalid GAP value: %s", value)
				}
				song.GAP = gap
			}
		} else if len(line) > 0 && (line[0] == ':' || line[0] == '*' || line[0] == '-' || line[0] == 'F') {
			note, err := parseNoteLine(line, bpm)
			if err != nil {
				return nil, fmt.Errorf("error parsing note line '%s': %w", line, err)
			}
			song.Notes = append(song.Notes, note)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Debugging: Print total song duration
	if len(song.Notes) > 0 {
		lastNote := song.Notes[len(song.Notes)-1]
		totalDuration := float64(lastNote.Timestamp+lastNote.Duration) + float64(song.GAP)
		fmt.Printf("Total Song Duration: %.2f seconds\n", totalDuration/1000)
	}

	return song, nil
}
