package game

import (
	"sync"

	"github.com/google/uuid"
)

type Player struct {
	ID   string
	Name string
	// Needs a way to clear old players, expire in 12 hours maybe. The frontend can store in local storage and just update this list again.
}

var (
	players  = make(map[string]*Player)
	playerMu sync.RWMutex
)

func CreatePlayer(player *Player) (*Player, error) {
	if existingPlayer, exists := players[player.ID]; exists {
		return existingPlayer, nil
	}

	playerMu.Lock()
	defer playerMu.Unlock()

	player.ID = uuid.NewString()
	players[player.ID] = player
	return player, nil
}

// func HandlePlayerCreation() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		var player Player
// 		decoder := json.NewDecoder(r.Body)
// 		err := decoder.Decode(&player)
// 		if err != nil {
// 			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
// 		}

// 		if existingPlayer, exists := players[player.ID]; exists {
// 			response, err := json.Marshal(existingPlayer)
// 			if err != nil {
// 				http.Error(w, "Failed to marshal player data", http.StatusInternalServerError)
// 				return
// 			}

// 			w.WriteHeader(http.StatusOK)
// 			w.Write(response)
// 			return
// 		}

// 		playerMu.Lock()
// 		defer playerMu.Unlock()

// 		player.ID = uuid.NewString()

// 		players[player.ID] = &player

// 		response, err := json.Marshal(player)
// 		if err != nil {
// 			http.Error(w, "Failed to marshal palyer data", http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusCreated)
// 		w.Write([]byte(response))
// 	}
// }
