package game

// Lobby represents a game lobby
type Lobby struct {
	ID         string
	Players    []*Player
	Game       *Game
	Host       string
	MaxPlayers int
	// Needs a way to clear old lobbies, expire in 12 hours maybe
	// Where can a player select stances? Let's do stances later.
	// Where to store game state? Just use nil?
}

// var (
// 	lobbies = make(map[string]*Lobby)
// 	lobbyMu sync.RWMutex
// )

// generateID creates a random 4 character alphanumeric ID
// func generateID() string {
// 	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	result := make([]byte, 4)
// 	for i := 0; i < 4; i++ {
// 		result[i] = chars[rand.Intn(len(chars))]
// 	}
// 	return string(result)
// }

// Create lobby
func NewLobby(lobbyID string, hostPlayer Player) *Lobby {

	lobby := &Lobby{
		ID:         lobbyID,
		Players:    nil,           // Host is the first player
		Game:       nil,           // No game started yet
		Host:       hostPlayer.ID, //Use the username of the person who joined(or maybe user id)
		MaxPlayers: 2,             // Always set to 2 as required
	}

	// Seems like have to append, can't assign directly to the lobby
	lobby.Players = append(lobby.Players, &hostPlayer)

	return lobby
}

// // Create lobby
// func HandleLobbyCreation() http.HandlerFunc { //Expect username. Do users need a UUID as well?
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		lobbyMu.Lock()
// 		defer lobbyMu.Unlock()

// 		lobbyID := generateID()

// 		// Keep generating new ID if collision occurs
// 		for _, exists := lobbies[lobbyID]; exists; {
// 			lobbyID = generateID()
// 		}

// 		var player Player
// 		decoder := json.NewDecoder(r.Body)
// 		err := decoder.Decode(&player)
// 		if err != nil {
// 			http.Error(w, "Invalid request body", http.StatusBadRequest)
// 			return
// 		}

// 		// Check if player exists, else add to list? No, players should already be added to list whenever they load the game...?
// 		// For now just throw error if player doesn't exist
// 		if _, exists := players[player.ID]; !exists {
// 			http.Error(w, "Player not found", http.StatusNotFound)
// 			return
// 		}

// 		lobby := &Lobby{
// 			ID:         lobbyID,
// 			Players:    nil,       // Host is the first player
// 			Game:       nil,       // No game started yet
// 			Host:       player.ID, //Use the username of the person who joined(or maybe user id)
// 			MaxPlayers: 2,         // Always set to 2 as required
// 		}

// 		// Seems like have to append, can't assign directly to the lobby
// 		lobby.Players = append(lobby.Players, &player)

// 		lobbies[lobby.ID] = lobby

// 		// Return lobby ID to client
// 		lobbyJson, err := json.Marshal(lobby)
// 		if err != nil {
// 			http.Error(w, "Falied to marshal lobby", http.StatusInternalServerError)
// 			return
// 		}
// 		w.WriteHeader(http.StatusCreated)
// 		w.Write([]byte(lobbyJson))
// 	}
// }

// Join lobby
// func HandleLobbyJoin() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// lobbyID := r.URL.Query().Get("id")

// 		// lobbyMu.Lock()
// 		// defer lobbyMu.Unlock()

// 		// lobby, exists := lobbies[lobbyID]
// 		// if !exists {
// 		// 	w.WriteHeader(http.StatusNotFound)
// 		// 	return
// 		// }

// 		// currentPlayers := 0
// 		// fmt.Sscanf(lobby.Players, "%d", &currentPlayers)

// 		// if currentPlayers >= lobby.MaxPlayers {
// 		// 	w.WriteHeader(http.StatusForbidden)
// 		// 	return
// 		// }

// 		// lobby.Players = fmt.Sprintf("%d", currentPlayers+1)
// 		w.WriteHeader(http.StatusOK)
// 	}
// }

// func HandleLobbyJoin(lobbyID string, player *Player) (*Lobby, error) {
// 	if lobby, exists := lobbies[lobbyID]; !exists {
// 		return nil, fmt.Errorf("Lobby not found")
// 	} else {
// 		lobbyMu.Lock()
// 		defer lobbyMu.Unlock()

// 		if len(lobby.Players) >= lobby.MaxPlayers {
// 			return nil, fmt.Errorf("Lobby is full")
// 		}

// 		lobby.Players = append(lobby.Players, player)
// 		return lobby, nil
// 	}
// }
