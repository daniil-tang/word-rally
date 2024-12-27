package game

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type GameState string

const (
	StateWaiting    GameState = "waiting"
	StateInProgress GameState = "inprogress"
	StateFinished   GameState = "finished"
)

type Game struct {
	ID    string
	State GameState //waiting, inprogress, finished
	// Turn          int    //Turns within a rally, starts with current server
	Score map[string]int
	// Guesses       map[string][]rune
	// Word          string
	CurrentServer int //Randomize first server. Changes after every rally. First server = starting player
	Rally         *Rally
}

type Rally struct {
	Turn    int
	Guesses map[string][]rune //Make sure to initialize it with make(map[string][]rune)
	Word    string
	// CurrentServer int
}

func (lobby *Lobby) CreateNewGame() *Lobby {
	lobby.Game = &Game{
		ID:    uuid.NewString(),
		State: StateWaiting,
		Score: make(map[string]int),
	}

	for _, player := range lobby.Players {
		lobby.Game.Score[player.ID] = 0
	}
	return lobby
}

func (lobby *Lobby) StartGame() (*Lobby, error) {
	if lobby.Game == nil {
		return nil, fmt.Errorf("Game instance not found")
	}

	// Check if game has already started
	if lobby.Game.State == StateInProgress {
		return nil, fmt.Errorf("Game already in progress")
	}

	//Get server
	serverIndex := rand.Intn(len(lobby.Players))
	lobby.Game.CurrentServer = serverIndex

	// Initialize rally
	lobby.initializeRally()

	// Update game state, the UI will swap once it receives the new game state and data
	lobby.Game.State = StateInProgress
	return lobby, nil
}

func (lobby *Lobby) initializeRally() *Lobby {
	lobby.Game.Rally = &Rally{
		Word:    "Hello",
		Guesses: make(map[string][]rune),
		Turn:    lobby.Game.CurrentServer,
	}

	for _, player := range lobby.Players {
		lobby.Game.Rally.Guesses[player.ID] = make([]rune, len(lobby.Game.Rally.Word))
	}

	return lobby
}

func (lobby *Lobby) Guess(player Player, actionDetails ActionDetails) (*Lobby, error) {
	// lobby.
	// Check if game is in progress
	if lobby.Game.State != StateInProgress {
		return nil, fmt.Errorf("Game not in progress")
	}

	for _, guessedLetter := range actionDetails.GuessedLetters {
		// Get unguessed letter index and put them in an array
		unguessedIndexes := []int{}
		for i, _ := range lobby.Game.Rally.Guesses[player.ID] {
			if lobby.Game.Rally.Guesses[player.ID][i] == '\x00' {
				unguessedIndexes = append(unguessedIndexes, i)
			}

			for _, i := range unguessedIndexes {
				if lobby.Game.Rally.Word[i] == byte(guessedLetter) {
					lobby.Game.Rally.Guesses[player.ID][i] = guessedLetter
					break
				}
			}

			if isRuneArrayFilled(lobby.Game.Rally.Guesses[player.ID]) {
				// Player wins the rally
				lobby.incrementScore(player.ID)
			} else {
				// If the guess failed....return something else?
				lobby.Game.Rally.Turn = (lobby.Game.Rally.Turn + 1) % len(lobby.Players)
			}
		}

		// If there's duplicate letters in a word, only reveal one
	}
	return lobby, nil
}

func isRuneArrayFilled(runes []rune) bool {
	for _, r := range runes {
		if r == '\x00' {
			return false
		}
	}
	return true
}

func (lobby *Lobby) incrementScore(playerID string) {
	lobby.Game.Score[playerID] += 1

	//Reset rally or declare game winner if score >= 3
	if lobby.Game.Score[playerID] >= 3 {
		lobby.Game.State = StateFinished
	} else {
		lobby.Game.CurrentServer = (lobby.Game.CurrentServer + 1) % len(lobby.Players)
		lobby.initializeRally()
	}
}

// func CreateNewGame(lobbyID string) {
// 	// Instantiate new game with whatever params
// 	lobbies[lobbyID].Game = &Game{
// 		ID:    uuid.NewString(),
// 		State: "waiting",
// 		Score: make(map[string]int),
// 	}

// 	for _, player := range lobbies[lobbyID].Players {
// 		lobbies[lobbyID].Game.Score[player.ID] = 0
// 	}
// }

// func StartGame(lobbyID string) error {
// 	if lobbies[lobbyID].Game == nil {
// 		return fmt.Errorf("Game not initialized")
// 	} else {
// 		Start(lobbyID)
// 		return nil
// 	}
// }

// func Start(lobbyID string) {
// 	game := lobbies[lobbyID].Game
// 	game.State = "inprogress"
// 	NewRally(lobbyID)

// 	// game.CurrentServer
// 	// Assign first server here?
// }

// func (g *Game) Guess() {
// 	// needs player ID
// 	// Breakdown word and compare letter by letter..?
// 	// Check if someone wins a rally
// 	// Update rally?

// }

// func NewRally(lobbyID string) {
// 	// Randomly select word, hardcoded to "Hello" for now
// 	game := lobbies[lobbyID].Game
// 	game.Rally = &Rally{
// 		Turn:    0,
// 		Guesses: make(map[string][]rune),
// 		Word:    "Hello",
// 	}
// }

// Select Stance

/*
Gamelogic
Where do I start the game? I'll assume everything will be done via lobbies. Because the "session" is essentially lobbies
So every request will include the lobby ID and player ID
The server knows which player's turn it is, how will the client know?

Does lobbies need to be aware of this? yes lobbies need to be aware of the Game
*/
