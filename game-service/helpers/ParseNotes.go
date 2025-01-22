package helpers

import (
	"fmt"
	"game-service/models"
	"strconv"
	"strings"
)

// Function to parse a note line
func parseNoteLine(line string, beatDuration float64) (models.Note, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 && line[0] != '-' { // Skip invalid note lines
		return models.Note{}, fmt.Errorf("invalid note line: %s", line)
	}

	// Handle pauses (lines starting with '-')
	if line[0] == '-' {
		timestamp, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "-")))
		if err != nil {
			return models.Note{}, err
		}
		return models.Note{
			Type:      "-",
			Timestamp: int(float64(timestamp) * beatDuration),
			Duration:  0,
			Pitch:     0,
			Text:      "",
		}, nil
	}

	// Parse regular or special notes
	timestamp, err := strconv.Atoi(parts[1])
	if err != nil {
		return models.Note{}, err
	}

	duration, err := strconv.Atoi(parts[2])
	if err != nil {
		return models.Note{}, err
	}

	pitch, err := strconv.Atoi(parts[3])
	if err != nil {
		return models.Note{}, err
	}

	// Combine the remaining parts as text
	text := strings.Join(parts[4:], " ")

	return models.Note{
		Type:      string(line[0]),
		Timestamp: int(float64(timestamp) * beatDuration),
		Duration:  int(float64(duration) * beatDuration),
		Pitch:     pitch,
		Text:      text,
	}, nil
}
