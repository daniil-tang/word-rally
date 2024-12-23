package game

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type GameMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
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

			log.Println(message.Content)
		}
	}
}
