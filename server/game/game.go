package game

type GameState string

const (
	StateWaiting    GameState = "waiting"
	StateInProgress GameState = "inprogress"
	StateFinished   GameState = "finished"
)

type Game struct {
	ID    string
	State string //waiting, inprogress, finished
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
