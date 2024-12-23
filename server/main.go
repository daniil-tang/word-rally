package main

import (
	"log"
	"net/http"
	websocket "server/game"
)

func main() {
	// Instantiate the game server

	// Handle websocket connections
	http.HandleFunc("/ws", websocket.HandleWebSocketConnection())

	log.Println("Server is running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
