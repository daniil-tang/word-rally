package game

import (
	"fmt"
)

type StanceType string

const (
	StanceTennis     StanceType = "tennis"
	StanceVolleyball StanceType = "volleyball"
	StanceFootball   StanceType = "football"
)

type SkillMetadata struct {
	Cooldown int
}

type Skill interface {
	Execute() string
	GetMetadata() SkillMetadata
}

type Stance interface {
	UseSkill(skill string) string
	GetSkillsMetadata() map[string]SkillMetadata
}

type BaseStance struct {
	StanceType StanceType
	Skills     map[string]Skill
}

func (s BaseStance) UseSkill(skill string) string {
	if skillFunc, exists := s.Skills[skill]; exists {
		return skillFunc.Execute()
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
			"tackle": &TackleSkill{},
		},
	}
}

type SecondServeSkill struct{}

func (s *SecondServeSkill) Execute() string {
	return "Executing a second serve in Tennis"
}

func (s *SecondServeSkill) GetMetadata() SkillMetadata {
	return SkillMetadata{Cooldown: 2}
}

type LiberoSkill struct{}

func (s *LiberoSkill) Execute() string {
	return "Libero player receives the ball in Volleyball"
}

func (s *LiberoSkill) GetMetadata() SkillMetadata {
	return SkillMetadata{Cooldown: 3}
}

type TackleSkill struct{}

func (s *TackleSkill) Execute() string {
	return "Executing a tackle in Football"
}

func (s *TackleSkill) GetMetadata() SkillMetadata {
	return SkillMetadata{Cooldown: 5}
}
