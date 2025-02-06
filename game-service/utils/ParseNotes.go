package utils

import (
	"fmt"
	"game-service/models"
	"strconv"
	"strings"
)

func parseNoteLine(line string, bpm float64) (models.Note, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 && line[0] != '-' { // Skip invalid note lines
		return models.Note{}, fmt.Errorf("invalid note line: %s", line)
	}

	// Handle pauses (lines starting with '-')
	if line[0] == '-' {
		timestampTicks, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "-")))
		if err != nil {
			return models.Note{}, err
		}
		timestampMs := int(float64(timestampTicks) * (15.0 / bpm) * 1000)

		return models.Note{
			Type:      "-",
			Timestamp: timestampMs,
			Duration:  0,
			Pitch:     0,
			Text:      "",
		}, nil
	}

	// Parse regular or special notes
	timestampTicks, err := strconv.Atoi(parts[1])
	if err != nil {
		return models.Note{}, err
	}

	durationTicks, err := strconv.Atoi(parts[2])
	if err != nil {
		return models.Note{}, err
	}

	pitch, err := strconv.Atoi(parts[3])
	if err != nil {
		return models.Note{}, err
	}

	// Remaining parts are the lyrics text
	text := strings.Join(parts[4:], " ")

	// Convert TICKS to milliseconds using correct formula
	timestampMs := int(float64(timestampTicks) * (15.0 / bpm) * 1000)
	durationMs := int(float64(durationTicks) * (15.0 / bpm) * 1000)

	return models.Note{
		Type:      string(line[0]),
		Timestamp: timestampMs,
		Duration:  durationMs,
		Pitch:     pitch,
		Text:      text,
	}, nil
}
