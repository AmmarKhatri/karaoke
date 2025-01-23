package scripts

import (
	"bufio"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

// StartProducer handles the producer role with interactive CLI inputs
func StartProducer(roomID, playerID string, interrupt chan os.Signal) {
	u := url.URL{Scheme: "ws", Host: Addr, Path: "/ws", RawQuery: "roomID=" + roomID + "&playerID=" + playerID + "&role=phone"}
	log.Printf("Producer %s connecting to %s", playerID, u.String())

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Producer connection error:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	// Goroutine to listen for messages from the server
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Producer read error: %v", err)
				close(done)
				return
			}
			log.Printf("Producer %s received: %s", playerID, message)
		}
	}()

	// Input reader from CLI
	inputReader := bufio.NewReader(os.Stdin)

	log.Println("Producer ready. Type eventType to send to the server or type 'exit' to quit.")

	// Main loop
	for {
		select {
		case <-done:
			// Handle connection close
			log.Println("Exiting producer due to connection close.")
			interrupt <- os.Interrupt
			return

		case <-interrupt:
			// Handle interrupt signal for graceful shutdown
			log.Println("Interrupt received. Sending CloseMessage and shutting down producer.")
			closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Producer shutting down")
			if err := conn.WriteMessage(websocket.CloseMessage, closeMessage); err != nil {
				log.Printf("Error sending close message: %v", err)
			}
			return

		default:
			// Read input from CLI
			log.Print("Enter eventType: ")
			input, err := inputReader.ReadString('\n')
			if err != nil {
				log.Println("Error reading input:", err)
				continue
			}

			input = input[:len(input)-1] // Remove newline character
			if input == "exit" {
				log.Println("Exiting producer.")
				interrupt <- os.Interrupt
				return
			}

			// Send the input as an event to the server
			err = conn.WriteJSON(map[string]string{
				"eventType": input,
				"playerID":  playerID,
			})
			if err != nil {
				log.Println("Producer write error:", err)
				close(done)
				return
			}
			log.Printf("Producer %s sent eventType: %s", playerID, input)
		}
	}
}
