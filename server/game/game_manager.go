package game

import (
	"fmt"
	"math/rand"
	"sync"
)

func generateLobbyID() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 4)
	for i := 0; i < 4; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

type GameManager struct {
	lobbies map[string]*Lobby
	mutex   sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		lobbies: make(map[string]*Lobby),
	}
}

func (gm *GameManager) CreateLobby(hostPlayer Player) (*Lobby, error) {
	if _, exists := players[hostPlayer.ID]; !exists {
		return nil, fmt.Errorf("Player not found.")
	}

	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Keep generating lobby ID until a unique ID is found
	lobbyID := generateLobbyID()
	for _, exists := gm.lobbies[lobbyID]; exists; {
		lobbyID = generateLobbyID()
	}

	newLobby := NewLobby(lobbyID, hostPlayer)
	gm.lobbies[lobbyID] = newLobby

	return newLobby, nil
}

// Should the arguments be pointers at all? Prolly not
func (gm *GameManager) JoinLobby(lobbyID string, player Player) (*Lobby, error) {
	if gm.lobbies[lobbyID] == nil {
		return nil, fmt.Errorf("Lobby with ID %s doesn't exist.", lobbyID)
	}

	lobby := gm.lobbies[lobbyID]

	if len(lobby.Players) >= lobby.MaxPlayers {
		return nil, fmt.Errorf("Lobby is full.")
	}

	lobby.Players = append(lobby.Players, &player)

	return lobby, nil
}

func (gm *GameManager) CreateGame(lobbyID string, player Player) (*Lobby, error) {
	if gm.lobbies[lobbyID] == nil {
		return nil, fmt.Errorf("Lobby with ID %s doesn't exist", lobbyID)
	}

	lobby := gm.lobbies[lobbyID]

	if player.ID != lobby.Host {
		return nil, fmt.Errorf("Player is not the host of the lobby")
	}

	lobby.CreateNewGame()
	return lobby, nil
}

func (gm *GameManager) StartGame(lobbyID string, player Player) (*Lobby, error) {
	if gm.lobbies[lobbyID] == nil {
		return nil, fmt.Errorf("Lobby with ID %s doesn't exist", lobbyID)
	}

	lobby := gm.lobbies[lobbyID]

	if player.ID != lobby.Host {
		return nil, fmt.Errorf("Player is not the host of the lobby")
	}

	if len(lobby.Players) <= 1 {
		return nil, fmt.Errorf("Not enough players to start game")
	}

	lobby.StartGame()
	return lobby, nil
}
