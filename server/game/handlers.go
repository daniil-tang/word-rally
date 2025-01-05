package game

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type PlayerAction string

const (
	ActionGuess    PlayerAction = "guess"
	ActionUseSkill PlayerAction = "useskill"
	ActionEndTurn  PlayerAction = "endturn"
)

const (
	GameResponseLobby string = "lobby"
)

type GameMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type GameResponse struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type RegisterConnectionRequest struct {
	Player Player `json:"player"`
}

type LobbyJoinRequest struct {
	LobbyID string `json:"lobbyId"` // The ID of the lobby
	Player  Player `json:"player"`  // The player data
}

type GameCreateRequest struct {
	LobbyID string `json:"lobbyId"` // The ID of the lobby
	Player  Player `json:"player"`  // The player data
}

type GameStartRequest struct {
	LobbyID string `json:"lobbyId"` // The ID of the lobby
	Player  Player `json:"player"`  // The player data
}

type PlayerSettingsRequest struct {
	LobbyID        string         `json:"lobbyId"` // The ID of the lobby
	Player         Player         `json:"player"`  // The player data
	PlayerSettings PlayerSettings `json:"playerSettings"`
}

type PlayerActionRequest struct {
	LobbyID       string        `json:"lobbyId"` // The ID of the lobby
	Player        Player        `json:"player"`  // The player data. Use ID to check if player is the host so that they can start the game
	Action        PlayerAction  `json:"action"`  // The action to be performed
	ActionDetails ActionDetails `json:"actionDetails"`
}

type ActionDetails struct {
	GuessedLetters []rune `json:"guessedLetters,omitempty"` // List of guessed letters (if guessing)
	SkillUsed      *Skill `json:"powerUsed,omitempty"`      // Power used (if activating an ability)
}

// type GuessActionDetails struct {
// 	GuessedLetters []rune   `json:"guessedLetters,omitempty"` // List of guessed letters (if guessing)
// 	SkillUsed    *Skill `json:"powerUsed,omitempty"`      // Power used (if activating an ability)
// }

type Skill interface {
	Activate() string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (this can be adjusted as necessary)
	},
}

func HandleCreateLobby(gm *GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
			return
		}

		var player Player
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&player)
		if err != nil {
			http.Error(w, "Error decoding player data.", http.StatusBadRequest)
			return
		}

		lobby, err := gm.CreateLobby(player)
		if err != nil {
			http.Error(w, "Error creating lobby.", http.StatusInternalServerError)
			return
		}

		lobbyJson, err := json.Marshal(lobby)
		if err != nil {
			http.Error(w, "Falied to marshal lobby", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(lobbyJson)
	}
}

func HandleGetLobby(gm *GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
			return
		}

		playerID := r.URL.Query().Get("playerID")
		if playerID == "" {
			http.Error(w, "Player ID is empty.", http.StatusBadRequest)
		}

		lobby, err := gm.GetLobbyByPlayer(playerID)
		if err != nil {
			http.Error(w, "Error retrieving lobby", http.StatusInternalServerError)
		}
		lobbyJson, err := json.Marshal(lobby)
		if err != nil {
			http.Error(w, "Falied to marshal lobby", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(lobbyJson)
	}
}

func HandleCreatePlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
			return
		}

		var player Player
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&player)
		if err != nil {
			http.Error(w, "Error decoding player data.", http.StatusBadRequest)
			return
		}

		_, err = CreatePlayer(&player)
		if err != nil {
			http.Error(w, "Error creating player.", http.StatusInternalServerError)
			return
		}

		playerJson, err := json.Marshal(player)
		if err != nil {
			http.Error(w, "Falied to marshal player", http.StatusInternalServerError)
			return
		}

		//Consider using json.NewEncoder(w).Encode()
		w.WriteHeader(http.StatusCreated)
		w.Write(playerJson)
	}
}

func HandleWebSocketConnection(gm *GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		conn, err := upgrader.Upgrade(w, r, nil)
		// log.Println("CONNECTION! ", conn)
		if err != nil {
			log.Println("Could not open websocket connection", err)
		}

		for {
			var message GameMessage
			err := conn.ReadJSON(&message)
			if err != nil {
				log.Println("Error reading message", err)
				break
			}

			switch message.Event {
			case "registerconnection":
				var registerConnectionRequest RegisterConnectionRequest
				err := json.Unmarshal([]byte(message.Data), &registerConnectionRequest)
				if err != nil {
					log.Println("Error unmarshalling register connection request", err)
					continue
				}
				gm.AddConnection(registerConnectionRequest.Player.ID, conn)
			case "joinlobby":
				var lobbyJoinRequest LobbyJoinRequest
				err := json.Unmarshal([]byte(message.Data), &lobbyJoinRequest)
				if err != nil {
					log.Println("Error unmarshalling lobby join request", err)
					continue
				}
				lobby, err := gm.JoinLobby(lobbyJoinRequest.LobbyID, lobbyJoinRequest.Player)
				// Needs error handling
				if err != nil {
					log.Println("Error joining lobby", err)
					continue
				}
				// Needs propagating the updated lobby to other players
				// err = conn.WriteJSON(lobby)
				encodedMsg, err := lobby.getEncodedLobbyWSResponse()
				if err != nil {
					log.Println("Error encoding message", err)
					continue // Exit if writing fails
				}

				// err = conn.WriteMessage(websocket.TextMessage, encodedMsg)
				// if err != nil {
				// 	log.Println("Error sending lobby back to client", err)
				// 	continue // Exit if writing fails
				// }
				gm.broadcastToLobbyPlayers(lobby.ID, encodedMsg)
			case "leavelobby":
				//TODO
				continue
			case "creategame":
				var gameCreateRequest GameCreateRequest
				err := json.Unmarshal([]byte(message.Data), &gameCreateRequest)
				if err != nil {
					log.Println("Error unmarshalling lobby game create request", err)
					continue
				}
				lobby, err := gm.CreateGame(gameCreateRequest.LobbyID, gameCreateRequest.Player)
				if err != nil {
					log.Println("Error creating game", err)
					continue
				}
				encodedMsg, err := lobby.getEncodedLobbyWSResponse()
				if err != nil {
					log.Println("Error encoding message", err)
					continue // Exit if writing fails
				}
				// err = conn.WriteMessage(websocket.TextMessage, encodedMsg)
				// if err != nil {
				// 	log.Println("Error sending lobby back to client", err)
				// 	continue
				// }
				gm.broadcastToLobbyPlayers(lobby.ID, encodedMsg)
			case "startgame":
				var gameStartRequest GameStartRequest
				err := json.Unmarshal([]byte(message.Data), &gameStartRequest)
				if err != nil {
					log.Println("Error unmarshalling lobby game start request", err)
					continue
				}
				lobby, err := gm.StartGame(gameStartRequest.LobbyID, gameStartRequest.Player)
				if err != nil {
					log.Println("Error starting game", err)
					continue
				}
				encodedMsg, err := lobby.getEncodedLobbyWSResponse()
				if err != nil {
					log.Println("Error encoding message", err)
					continue // Exit if writing fails
				}
				// err = conn.WriteMessage(websocket.TextMessage, encodedMsg)
				// if err != nil {
				// 	log.Println("Error sending lobby back to client", err)
				// 	continue
				// }
				gm.broadcastToLobbyPlayers(lobby.ID, encodedMsg)
			case "playeraction":
				log.Println("Player Action")
				var playerActionRequest PlayerActionRequest
				err := json.Unmarshal([]byte(message.Data), &playerActionRequest)
				if err != nil {
					log.Println("Error unmarshalling player action request", err)
					continue
				}
				lobby, err := gm.HandlePlayerAction(playerActionRequest.LobbyID, playerActionRequest.Player, playerActionRequest.Action, playerActionRequest.ActionDetails)

				if err != nil {
					log.Println("Error handling player action", err)
					continue
				}

				encodedMsg, err := lobby.getEncodedLobbyWSResponse()
				if err != nil {
					log.Println("Error encoding message", err)
					continue // Exit if writing fails
				}
				// err = conn.WriteMessage(websocket.TextMessage, encodedMsg)
				// if err != nil {
				// 	log.Println("Error sending lobby back to client", err)
				// 	continue
				// }
				gm.broadcastToLobbyPlayers(lobby.ID, encodedMsg)
			case "updateplayersettings":
				var PlayerSettingsRequest PlayerSettingsRequest
				err := json.Unmarshal([]byte(message.Data), &PlayerSettingsRequest)
				if err != nil {
					log.Println("Error unmarshalling player settings request", err)
					continue
				}
				lobby, err := gm.UpdatePlayerSettings(PlayerSettingsRequest.LobbyID, PlayerSettingsRequest.Player, PlayerSettingsRequest.PlayerSettings)
				if err != nil {
					log.Println("Error updating player settings", err)
					continue
				}

				encodedMsg, err := lobby.getEncodedLobbyWSResponse()
				if err != nil {
					log.Println("Error encoding message", err)
					continue // Exit if writing fails
				}
				// err = conn.WriteMessage(websocket.TextMessage, encodedMsg)
				// if err != nil {
				// 	log.Println("Error sending lobby back to client", err)
				// 	continue
				// }
				gm.broadcastToLobbyPlayers(lobby.ID, encodedMsg)
			}
		}
	}
}

func getEncodedWSResponse(event string, data string) ([]byte, error) {
	resp := &GameResponse{
		Event: event,
		Data:  data,
	}
	encodedResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error encoding message:", err)
		return nil, err
	}
	return encodedResp, nil
}

func (lobby *Lobby) getEncodedLobbyWSResponse() ([]byte, error) {
	lobbyMarshaled, err := json.Marshal(lobby)
	if err != nil {
		log.Println("Error marshaling lobby:", err)
		return nil, err
	}
	encodedResp, err := getEncodedWSResponse(GameResponseLobby, string(lobbyMarshaled))
	if err != nil {
		log.Println("Error encoding lobby message:", err)
		return nil, err
	}

	return encodedResp, nil
}

func (gm *GameManager) broadcastToLobbyPlayers(lobbyID string, msg []byte) {
	for _, connToBroadcast := range gm.GetLobbyConnections(lobbyID) {
		err := connToBroadcast.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error sending lobby back to client", err)
			continue
		}
	}
}
