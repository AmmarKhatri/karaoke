package utils

import (
	"context"
	"encoding/json"
	"log"
)

// findMatch finds another request in the queue with the same skill level
func FindMatch(ctx context.Context, queueKey, skillLevel string, currentRequest string) (map[string]string, error) {
	// Retrieve the queue data
	queueData, err := Redis.LRange(ctx, queueKey, 0, -1).Result()
	if err != nil {
		log.Printf("Error fetching queue data: %v", err)
		return nil, err
	}

	for _, requestJSON := range queueData {
		var player map[string]string
		if err := json.Unmarshal([]byte(requestJSON), &player); err != nil {
			log.Printf("Error unmarshalling player data: %v", err)
			continue
		}

		// Check if the request matches
		if player["skillLevel"] == skillLevel && requestJSON != currentRequest {
			// Remove the matched request from the queue
			Redis.LRem(ctx, queueKey, 1, requestJSON)
			return player, nil
		}
	}

	return nil, nil
}
