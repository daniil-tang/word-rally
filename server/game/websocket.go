package game

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type GameMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type LobbyJoinRequest struct {
	LobbyID string `json:"lobbyId"` // The ID of the lobby
	Player  Player `json:"player"`  // The player data
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (this can be adjusted as necessary)
	},
}

func HandleWebSocketConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Could not open websocket connection", err)
		}

		// Connection succeeded

		for {
			var message GameMessage
			err := conn.ReadJSON(&message)
			if err != nil {
				log.Println("Error reading message", err)
				break
			}

			switch message.Type {
			case "joinlobby":
				var lobbyJoinRequest LobbyJoinRequest
				err := json.Unmarshal([]byte(message.Content), &lobbyJoinRequest)
				if err != nil {
					log.Println("Error unmarshalling lobby join request", err)
					continue
				}
				lobby, err := HandleLobbyJoin(lobbyJoinRequest.LobbyID, &lobbyJoinRequest.Player)
				// Needs error handling
				if err != nil {
					log.Println("Error joining lobby", err)
					break
				}
				// Needs propagating the updated lobby to other players
				err = conn.WriteJSON(lobby)
				if err != nil {
					log.Println("Error sending lobby back to client", err)
					return // Exit if writing fails
				}
			case "playeraction":
			}

			log.Println(message.Content)
		}
	}
}
