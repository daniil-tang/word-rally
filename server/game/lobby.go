package game

// Lobby represents a game lobby
type Lobby struct {
	ID             string
	Players        []*Player
	Game           *Game
	Host           string
	MaxPlayers     int
	PlayerSettings map[string]*PlayerSettings
	// Needs a way to clear old lobbies, expire in 12 hours maybe
	// Where can a player select stances? Let's do stances later.
	// Where to store game state? Just use nil?
}

type PlayerSettings struct {
	Stance StanceType
}

// Create lobby
func NewLobby(lobbyID string, hostPlayer Player) *Lobby {

	lobby := &Lobby{
		ID:             lobbyID,
		Players:        nil,           // Host is the first player
		Game:           nil,           // No game started yet
		Host:           hostPlayer.ID, //Use the username of the person who joined(or maybe user id)
		MaxPlayers:     2,             // Always set to 2 as required
		PlayerSettings: make(map[string]*PlayerSettings),
	}

	// Seems like have to append, can't assign directly to the lobby
	lobby.Players = append(lobby.Players, &hostPlayer)

	return lobby
}

func (lobby *Lobby) UpdatePlayerSettings(playerID string, settings PlayerSettings) *Lobby {
	// Is it alright to use address here?
	lobby.PlayerSettings[playerID] = &settings
	return lobby
}
