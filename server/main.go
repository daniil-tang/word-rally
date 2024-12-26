package main

import (
	"log"
	"net/http"
	game "server/game"
)

func main() {
	// Instantiate the game server
	gameManager := game.NewGameManager()

	// Handle websocket connections
	http.HandleFunc("/ws", game.HandleWebSocketConnection(gameManager))
	http.HandleFunc("/createlobby", game.HandleCreateLobby(gameManager))
	// http.HandleFunc("/joinlobby", game.HandleJoinLobby(gameManager))
	http.HandleFunc("/createplayer", game.HandleCreatePlayer())

	log.Println("Server is running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
