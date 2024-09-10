package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"socket-service/utils"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// Create a new Socket.IO server with polling-only transport
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
		PingTimeout:  120 * 1000, // Ping timeout: 120 seconds
		PingInterval: 50 * 1000,  // Ping interval: 50 seconds
	})

	// Initialize Redis connection
	utils.ConnectToRedis()

	// Handle connection events
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("Connected:", s.ID())
		s.Emit("connect", "Connected to server")
		return nil
	})

	// Handle "joinRoom" event
	server.OnEvent("/", "joinRoom", func(s socketio.Conn, roomID string) {
		fmt.Println("Checking room existence:", roomID)
		streamExists, err := checkRedisStreamExists(roomID)
		if err != nil || !streamExists {
			s.Emit("error", fmt.Sprintf("Room %s does not exist", roomID))
			return
		}
		s.Join(roomID)
		fmt.Printf("User %s joined room %s\n", s.ID(), roomID)
		go listenToRedisStream(s, roomID)
	})

	// Handle "produceEvent" event
	server.OnEvent("/", "produceEvent", func(s socketio.Conn, roomID, data string) {
		err := publishToRedisStream(roomID, data)
		if err != nil {
			s.Emit("error", "Failed to publish event to Redis stream")
			return
		}
		fmt.Printf("Published event to room %s: %s\n", roomID, data)
	})

	// Handle disconnection
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Printf("User %s disconnected: %s\n", s.ID(), reason)
	})

	// Handle errors
	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("Error occurred:", e)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO listen error: %s\n", err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("../asset")))

	log.Println("Serving at localhost:3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// Check if the Redis stream exists for the game room
func checkRedisStreamExists(roomID string) (bool, error) {
	streamInfo, err := utils.Redis.XInfoStream(ctx, "stream:"+roomID).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return streamInfo.Length > 0, nil
}

// Publish an event to the Redis stream
func publishToRedisStream(roomID, data string) error {
	_, err := utils.Redis.XAdd(ctx, &redis.XAddArgs{
		Stream: "stream:" + roomID,
		Values: map[string]interface{}{
			"data": data,
		},
	}).Result()
	return err
}

// Listen to the Redis stream and send events to the Socket.IO client
func listenToRedisStream(s socketio.Conn, roomID string) {
	lastID := "$" // Start reading new messages from the latest entry

	for {
		entries, err := utils.Redis.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"stream:" + roomID, lastID},
			Block:   0,
			Count:   1,
		}).Result()

		if err != nil {
			log.Printf("Error reading from Redis stream: %v", err)
			break
		}

		for _, entry := range entries[0].Messages {
			data := entry.Values["data"].(string)
			s.Emit("event", data)
			lastID = entry.ID
		}
	}
}
