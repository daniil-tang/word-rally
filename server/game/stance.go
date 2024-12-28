package game

type StanceType string

const (
	StanceTennis     StanceType = "tennis"
	StanceVolleyball StanceType = "volleyball"
	StanceFootball   StanceType = "football"
)

type Stance interface {
	UseSkill(skill string) string
}

type TennisStance struct{}

// Prolly make more sense to pass the game or the rally, because actions are performed on the rally.
func (s TennisStance) UseSkill(skill string) string {
	return "Tennis " + skill
}

type VolleyballStance struct{}

func (s VolleyballStance) UseSkill(skill string) string {
	return "Volleyball " + skill
}

type FootballStance struct{}

func (s FootballStance) UseSkill(skill string) string {
	return "Football " + skill
}
