package game

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

// Lobby represents a game lobby
type Lobby struct {
	ID         string
	Players    string
	Game       string
	Host       string
	MaxPlayers int
}

var (
	lobbies = make(map[string]*Lobby)
	lobbyMu sync.RWMutex
)

// generateID creates a random 4 character alphanumeric ID
func generateID() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 4)
	for i := 0; i < 4; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// Create lobby
func HandleLobbyCreation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lobbyMu.Lock()
		defer lobbyMu.Unlock()

		lobbyID := generateID()
		
		// Keep generating new ID if collision occurs
		for _, exists := lobbies[lobbyID]; exists; {
			lobbyID = generateID()
		}

		lobby := &Lobby{
			ID:         lobbyID,
			Players:    "1",  // Host is the first player
			Game:       "0",  // No game started yet
			Host:       lobbyID,
			MaxPlayers: 2,    // Always set to 2 as required
		}
		
		lobbies[lobby.ID] = lobby

		// Return lobby ID to client
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(lobby.ID))
	}
}

// Join lobby
func HandleLobbyJoin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lobbyID := r.URL.Query().Get("id")
		
		lobbyMu.Lock()
		defer lobbyMu.Unlock()

		lobby, exists := lobbies[lobbyID]
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		currentPlayers := 0
		fmt.Sscanf(lobby.Players, "%d", &currentPlayers)
		
		if currentPlayers >= lobby.MaxPlayers {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		lobby.Players = fmt.Sprintf("%d", currentPlayers + 1)
		w.WriteHeader(http.StatusOK)
	}
}
