package main

import (
	"log"
	"net/http"
	game "server/game"
)

func main() {
	// Instantiate the game server

	// Handle websocket connections
	http.HandleFunc("/ws", game.HandleWebSocketConnection())
	http.HandleFunc("/createlobby", game.HandleLobbyCreation())
	// http.HandleFunc("/joinlobby", game.HandleLobbyJoin())
	http.HandleFunc("/createplayer", game.HandlePlayerCreation())

	log.Println("Server is running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
