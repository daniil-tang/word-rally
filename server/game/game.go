package game

type Game struct {
	ID            string
	State         string //waiting, inprogress, finished
	Turn          int    //Turns within a rally, starts with current server
	Score         map[string]int
	Guesses       map[string][]rune
	Word          string
	CurrentServer int //Randomize first server. Changes after every rally. First server = starting player
}
