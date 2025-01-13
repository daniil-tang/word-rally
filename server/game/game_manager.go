package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
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
	lobbies     map[string]*Lobby
	mutex       sync.Mutex
	connections map[string]*websocket.Conn
	connMutex   sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		lobbies:     make(map[string]*Lobby),
		connections: make(map[string]*websocket.Conn),
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
	// lobbyID := "ABCD"

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

	// Needs check whether player is already in a lobby

	lobby.Players = append(lobby.Players, &player)
	lobby.PlayerSettings[player.ID] = &PlayerSettings{
		Stance: StanceTennis, // Default Stance
		Ready:  false,
	}

	return lobby, nil
}

// Call this before join/create lobby
func (gm *GameManager) GetLobbyByPlayer(playerID string) (*Lobby, error) {
	for _, lobby := range gm.lobbies {
		for _, p := range lobby.Players {
			if p.ID == playerID {
				return lobby, nil
			}
		}
	}
	return nil, nil
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

	_, err := lobby.StartGame()
	if err != nil {
		return nil, err
	}

	return lobby, nil
}

func (gm *GameManager) HandlePlayerAction(lobbyID string, player Player, action PlayerAction, actionDetails ActionDetails) (*Lobby, error) {
	if gm.lobbies[lobbyID] == nil {
		return nil, fmt.Errorf("Lobby with ID %s doesn't exist", lobbyID)
	}

	lobby := gm.lobbies[lobbyID]
	log.Printf("ACTION %s", action)
	switch action {
	case ActionGuess:
		playerGuessesBefore := countNonEmptyElements(lobby.Game.Rally.Guesses[player.ID])
		_, err := lobby.Guess(player, actionDetails)
		if err != nil {
			actionLog := NewActionLog(player.ID, err.Error())
			if encodedLog, err := actionLog.getEncodedActionLogResponse(); err == nil {
				gm.broadcastToLobbyPlayers(lobby.ID, encodedLog)
			}
			// return nil, err
		}

		var actionLogMsg string
		if playerGuessesBefore < countNonEmptyElements(lobby.Game.Rally.Guesses[player.ID]) {
			actionLogMsg = player.Name + " made a correct guess"
		} else {
			actionLogMsg = player.Name + " made an incorrect guess"
		}

		actionLog := NewActionLog(player.ID, actionLogMsg)
		if encodedLog, err := actionLog.getEncodedActionLogResponse(); err == nil {
			gm.broadcastToLobbyPlayers(lobby.ID, encodedLog)
		}
	case ActionUseSkill:
		_, err := lobby.UseSkill(player, actionDetails)
		if err != nil {
			return nil, err
		}
		actionLog := NewActionLog(player.ID, fmt.Sprintf(player.Name+" used skill: %s", actionDetails.SkillUsed))
		if encodedLog, err := actionLog.getEncodedActionLogResponse(); err == nil {
			gm.broadcastToLobbyPlayers(lobby.ID, encodedLog)
		}
	case ActionEndTurn:
		_, err := lobby.EndTurn(player)
		if err != nil {
			return nil, err
		}
		actionLog := NewActionLog(player.ID, player.Name+" ended their turn")
		if encodedLog, err := actionLog.getEncodedActionLogResponse(); err == nil {
			gm.broadcastToLobbyPlayers(lobby.ID, encodedLog)
		}
	}

	return lobby, nil
}

func countNonEmptyElements(runes []rune) int {
	count := 0
	for _, r := range runes {
		if r != '\x00' { // Check if the element is an empty rune (zero value)
			count++
		}
	}
	return count
}

func (gm *GameManager) UpdatePlayerSettings(lobbyID string, player Player, settings PlayerSettings) (*Lobby, error) {
	if gm.lobbies[lobbyID] == nil {
		return nil, fmt.Errorf("Lobby with ID %s doesn't exist", lobbyID)
	}

	lobby := gm.lobbies[lobbyID]

	lobby.UpdatePlayerSettings(player.ID, settings)

	return lobby, nil
}

func (gm *GameManager) AddConnection(playerID string, connection *websocket.Conn) {
	gm.connMutex.Lock()
	defer gm.connMutex.Unlock()

	gm.connections[playerID] = connection
}

func (gm *GameManager) GetLobbyConnections(lobbyID string) []*websocket.Conn {
	var lobbyConnections []*websocket.Conn

	for _, p := range gm.lobbies[lobbyID].Players {
		lobbyConnections = append(lobbyConnections, gm.connections[p.ID])
	}

	return lobbyConnections
}
