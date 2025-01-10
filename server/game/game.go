package game

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type GameState string

const (
	StateWaiting    GameState = "waiting"
	StateInProgress GameState = "inprogress"
	StateFinished   GameState = "finished"
)

type Game struct {
	ID              string
	State           GameState //waiting, inprogress, finished
	Score           map[string]int
	CurrentServer   int //Randomize first server. Changes after every rally. First server = starting player
	Rally           *Rally
	Settings        *GameSettings
	PlayerCooldowns map[string]map[SkillType]int
}

// StatusEffect represents various effects that can be applied to a player
type StatusEffect struct {
	IsActive bool
	Duration int
}

type SkillType string

const (
	Goalkeeper SkillType = "goalkeeper"
	Tackle     SkillType = "tackle"
	Fault      SkillType = "fault"
)

// type StatusEffectType string

// const (
// 	Goalkeeper StatusEffectType = "goalkeeper"
// )

// NewStatusEffect creates a new StatusEffect instance
// func NewStatusEffect() *StatusEffect {
// 	return &StatusEffect{
// 		BlockNextCorrectGuess: false,
// 		Duration:              0,
// 	}
// }

type Rally struct {
	Turn             int
	TurnActionPoints map[string]*TurnActionPoints
	StatusEffects    map[string]map[SkillType]*StatusEffect // playerID -> StatusEffect
	Guesses          map[string][]rune                      //Make sure to initialize it with make(map[string][]rune)
	Word             string
	// CurrentServer int
}

type TurnActionPoints struct {
	Guess int
	Skill int
}

type GameSettings struct {
	// Points and timer maybe
}

var WordList = []string{"HELLO", "GOODBYE"}

func (lobby *Lobby) CreateNewGame() *Lobby {
	lobby.Game = &Game{
		ID:              uuid.NewString(),
		State:           StateWaiting,
		Score:           make(map[string]int),
		PlayerCooldowns: make(map[string]map[SkillType]int), //Cooldowns should carry over between rallies
	}

	for _, player := range lobby.Players {
		lobby.Game.Score[player.ID] = 0
		lobby.Game.PlayerCooldowns[player.ID] = make(map[SkillType]int)
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
	serverIndex := 0 //rand.Intn(len(lobby.Players))
	lobby.Game.CurrentServer = serverIndex

	// Initialize rally
	lobby.initializeRally()

	// Update game state, the UI will swap once it receives the new game state and data
	lobby.Game.State = StateInProgress
	return lobby, nil
}

func (lobby *Lobby) initializeRally() *Lobby {
	lobby.Game.Rally = &Rally{
		Word:             "HELLO",
		StatusEffects:    make(map[string]map[SkillType]*StatusEffect),
		Guesses:          make(map[string][]rune),
		Turn:             lobby.Game.CurrentServer,
		TurnActionPoints: make(map[string]*TurnActionPoints),
	}

	for _, player := range lobby.Players {
		lobby.Game.Rally.Guesses[player.ID] = make([]rune, len(lobby.Game.Rally.Word))
		lobby.updatePlayerTurnActionPoints(player.ID, 0, 0)
		lobby.Game.Rally.StatusEffects[player.ID] = make(map[SkillType]*StatusEffect)
	}

	lobby.initializeNextPlayerTurn(lobby.Players[lobby.Game.CurrentServer].ID)

	return lobby
}

func (lobby *Lobby) Guess(player Player, actionDetails ActionDetails) (*Lobby, error) {
	// Check if game is in progress
	if lobby.Game.State != StateInProgress {
		return nil, fmt.Errorf("Game not in progress")
	}

	goalkeeperTriggered := false
	for _, guessedLetter := range actionDetails.GuessedLetters {
		// Check if player has guess actions
		if lobby.Game.Rally.TurnActionPoints[player.ID].Guess <= 0 {
			return nil, fmt.Errorf(player.Name + " has no guess actions")
		}

		// Get unguessed letter index and put them in an array
		unguessedIndexes := []int{}
		for i, _ := range lobby.Game.Rally.Guesses[player.ID] {
			if lobby.Game.Rally.Guesses[player.ID][i] == '\x00' {
				unguessedIndexes = append(unguessedIndexes, i)
			}
		}
		// Detect if failed guess? Does it matter?
	GuessLoop:
		for _, i := range unguessedIndexes {
			if lobby.Game.Rally.Word[i] == byte(guessedLetter) {

				// Goalkeeper effect
				goalkeeperEffect, exists := lobby.Game.Rally.StatusEffects[player.ID][Goalkeeper]
				// log.Printf("GOALKEEPER OR NOT" + st)
				if exists && goalkeeperEffect.IsActive && goalkeeperEffect.Duration > 0 {

					goalkeeperEffect.Duration--
					if goalkeeperEffect.Duration == 0 {
						goalkeeperEffect.IsActive = false
					}
					goalkeeperTriggered = true
					log.Printf("PLAYER GUESSED SKILLLLL TEST %s Index: %d, Letter: %c", player.Name, i, lobby.Game.Rally.Word[i])
					break GuessLoop
				} else {
					//Goalkeeper effect end
					lobby.Game.Rally.Guesses[player.ID][i] = guessedLetter
					log.Printf("PLAYER GUESSED TEST %s Index: %d, Letter: %c", player.Name, i, lobby.Game.Rally.Word[i])
					break GuessLoop
				}
			}
		}

		// Moved this to EndTurn

		// Add cooldown to the skill(+1) because the CD will be decremented this turn. Move Ability use to separate function...?
		// Players can guess + they can select a skill.
		// Once ready players hit "Run" or "Initiate" to trigger the guess + skill activation...?
		// Skills should be processed AFTER the guess returns?

		// Check if the correct guess should be blocked by a Goalkeeper skill
		// if effect, exists := lobby.Game.Rally.StatusEffects[player.ID]; exists && effect[StatusEffectBlockNextCorrectGuess] && effect.Duration > 0 {
		// 	// Decrement duration and block the guess
		// 	effect.Duration--
		// 	if effect.Duration == 0 {
		// 		delete(lobby.Game.Rally.StatusEffects, player.ID)
		// 	}
		// 	break // Skip revealing the letter
		// }

		// What if the player has "Action Points". 1 GuessAction point and 1 SkillAction point. Player is free to trigger these in whichever order?
		// Prolly there's a need for "Buffered" or "Queued" effects that will trigger during the opponent's turn.
	}

	// If there's duplicate letters in a word, only reveal one

	// Reduce guess action points
	lobby.Game.Rally.TurnActionPoints[player.ID].Guess -= 1

	if goalkeeperTriggered {
		return lobby, fmt.Errorf(player.Name + "'s correct guess was blocked by the Goalkeeper skill!")
	} else {
		return lobby, nil
	}
}

func (lobby *Lobby) UseSkill(player Player, actionDetails ActionDetails) (*Lobby, error) {
	if lobby.Game.State != StateInProgress {
		return nil, fmt.Errorf("Game not in progress")
	}
	if lobby.Game.Rally.TurnActionPoints[player.ID].Skill <= 0 {
		return nil, fmt.Errorf("Player has no skill actions")
	}

	// Switch and handle skills
	var stance Stance
	switch lobby.PlayerSettings[player.ID].Stance {
	case StanceTennis:
		stance = NewTennisStance()
	case StanceVolleyball:
		stance = NewVolleyballStance()
	case StanceFootball:
		stance = NewFootballStance()
	}

	res := stance.UseSkill(lobby, actionDetails.SkillUsed)
	fmt.Printf("USE LE SKILL %s", res)

	lobby.Game.Rally.TurnActionPoints[player.ID].Skill -= 1

	return lobby, nil
}

func (lobby *Lobby) EndTurn(player Player) (*Lobby, error) {
	log.Printf("END TURN")

	if isRuneArrayFilled(lobby.Game.Rally.Guesses[player.ID]) {
		// Player wins the rally
		lobby.incrementScore(player.ID)

		// Check if player wins the gamee
	} else {
		// If the guess failed....return something else?

		// Reduce cooldowns of all skills of current player
		for skill, cd := range lobby.Game.PlayerCooldowns[player.ID] {
			if cd > 0 {
				lobby.Game.PlayerCooldowns[player.ID][skill] -= 1
			}
		}

		lobby.Game.Rally.Turn = (lobby.Game.Rally.Turn + 1) % len(lobby.Players)

		//Initialize next player
		lobby.initializeNextPlayerTurn(lobby.Players[lobby.Game.Rally.Turn].ID)
	}
	return lobby, nil
}

func (lobby *Lobby) initializeNextPlayerTurn(nextTurnPlayerID string) *Lobby {
	lobby.updatePlayerTurnActionPoints(nextTurnPlayerID, 1, 1)
	return lobby
}

func (lobby *Lobby) updatePlayerTurnActionPoints(playerID string, guessActionPoints int, skillActionPoints int) *Lobby {
	lobby.Game.Rally.TurnActionPoints[playerID] = &TurnActionPoints{
		Guess: guessActionPoints,
		Skill: skillActionPoints,
	}
	return lobby
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

// Select Stance

/*
Gamelogic
Where do I start the game? I'll assume everything will be done via lobbies. Because the "session" is essentially lobbies
So every request will include the lobby ID and player ID
The server knows which player's turn it is, how will the client know?

Does lobbies need to be aware of this? yes lobbies need to be aware of the Game
*/
