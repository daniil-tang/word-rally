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

	if player.ID == "" {
		player.ID = uuid.NewString()
	}
	players[player.ID] = player
	return player, nil
}
