package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type StanceType string

const (
	StanceTennis     StanceType = "tennis"
	StanceVolleyball StanceType = "volleyball"
	StanceFootball   StanceType = "football"
)

type SkillMetadata struct {
	Cooldown int
	Duration int
}

type Skill interface {
	Execute(lobby *Lobby) string
	GetMetadata() SkillMetadata
}

type Stance interface {
	UseSkill(lobby *Lobby, skill string) string
	GetSkillsMetadata() map[string]SkillMetadata
}

type BaseStance struct {
	StanceType StanceType
	Skills     map[string]Skill
}

func (s BaseStance) UseSkill(lobby *Lobby, skill string) string {
	if skillFunc, exists := s.Skills[skill]; exists {
		return skillFunc.Execute(lobby)
	}
	return fmt.Sprintf("Skill %s not found!", skill)
}

func (s BaseStance) GetSkillsMetadata() map[string]SkillMetadata {
	metadata := make(map[string]SkillMetadata)
	for name, skill := range s.Skills {
		metadata[name] = skill.GetMetadata()
	}
	return metadata
}

func NewTennisStance() Stance {
	return BaseStance{
		StanceType: StanceTennis,
		Skills: map[string]Skill{
			"secondserve": &SecondServeSkill{},
		},
	}
}

func NewVolleyballStance() Stance {
	return BaseStance{
		StanceType: StanceVolleyball,
		Skills: map[string]Skill{
			"libero": &LiberoSkill{},
		},
	}
}

func NewFootballStance() Stance {
	return BaseStance{
		StanceType: StanceFootball,
		Skills: map[string]Skill{
			"tackle":     &TackleSkill{},
			"goalkeeper": &GoalkeeperSkill{},
		},
	}
}

type SecondServeSkill struct{}

func (s *SecondServeSkill) Execute(lobby *Lobby) string {
	return "Executing a second serve in Tennis"
}

func (s *SecondServeSkill) GetMetadata() SkillMetadata {
	return SkillMetadata{Cooldown: 2}
}

type LiberoSkill struct{}

func (s *LiberoSkill) Execute(lobby *Lobby) string {
	return "Libero player receives the ball in Volleyball"
}

func (s *LiberoSkill) GetMetadata() SkillMetadata {
	return SkillMetadata{Cooldown: 3}
}

type TackleSkill struct{}

func (s *TackleSkill) Execute(lobby *Lobby) string {
	// Get the current player and opponent IDs
	currentPlayerID := lobby.Game.Rally.Turn % 2
	opponentID := (currentPlayerID + 1) % 2

	var playerIDs []string
	for id := range lobby.Game.Rally.Guesses {
		playerIDs = append(playerIDs, id)
	}

	currentPlayer := playerIDs[currentPlayerID]
	opponent := playerIDs[opponentID]

	// Find all correct guesses from opponent that current player hasn't guessed yet
	var availableGuesses []int
	for i, opponentGuess := range lobby.Game.Rally.Guesses[opponent] {
		if opponentGuess != '\x00' && lobby.Game.Rally.Guesses[currentPlayer][i] == '\x00' {
			availableGuesses = append(availableGuesses, i)
		}
	}

	// If there are no available guesses to steal, return early
	if len(availableGuesses) == 0 {
		return "No new correct guesses available to steal"
	}

	// Randomly select one of the available guesses
	rand.Seed(time.Now().UnixNano())
	selectedIndex := availableGuesses[rand.Intn(len(availableGuesses))]

	// Copy the guess to the current player's guesses and remove it from opponent
	stolenGuess := lobby.Game.Rally.Guesses[opponent][selectedIndex]
	lobby.Game.Rally.Guesses[currentPlayer][selectedIndex] = stolenGuess
	lobby.Game.Rally.Guesses[opponent][selectedIndex] = '\x00'

	lobby.Game.PlayerCooldowns[currentPlayer][Tackle] = s.GetMetadata().Cooldown

	return "Successfully stole one correct guess from opponent"
}

func (s *TackleSkill) GetMetadata() SkillMetadata {
	return SkillMetadata{Cooldown: 5, Duration: 0}
}

type GoalkeeperSkill struct{}

func (s *GoalkeeperSkill) Execute(lobby *Lobby) string {
	currentPlayer := lobby.Players[lobby.Game.Rally.Turn%2]
	opponent := lobby.Players[(lobby.Game.Rally.Turn+1)%2]
	log.Printf("OPPONENT NAME " + opponent.Name)
	lobby.Game.Rally.StatusEffects[opponent.ID][Goalkeeper] = &StatusEffect{
		IsActive: true,
		Duration: s.GetMetadata().Duration,
	}

	lobby.Game.PlayerCooldowns[currentPlayer.ID][Goalkeeper] = s.GetMetadata().Cooldown

	return fmt.Sprintf("Player %d activates Goalkeeper skill - next 3 correct guesses by opponent will be blocked!")
}

func (s *GoalkeeperSkill) GetMetadata() SkillMetadata {
	return SkillMetadata{Cooldown: 4, Duration: 1} // Higher cooldown due to powerful effect
}

/*
Football (Soccer) Stance:
Tackle: Steal a correct guess from the opponent.
Goalkeeper: Block the opponent's correct guess.
Offside Trap: Force the opponent to lose their next turn or guess.
Tennis Stance:
Ace: Make an unpredictable guess that can bypass the opponent’s block.
Rally: Get an additional guess point for the next turn.
Slice: Delay the opponent's next action by one turn.
Volleyball Stance:
Spike: Knock out one of the opponent’s guesses.
Serve: Reveal a letter from the word, but only for you.
Block: Prevent the opponent from using their skill for the next turn.
*/
