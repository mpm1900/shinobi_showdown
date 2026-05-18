package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var HumanBoulder = MakeHumanBoulder()

func MakeHumanBoulder() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Human Boulder",
		Description: "Damage is based on the user's Defense rather than Attack.",
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(70),
		Stat:        game.Ptr(game.StatDefense),
		Nature:      game.Ptr(game.NsEarth),
		Cost:        game.Ptr(0),
		Jutsu:       game.Taijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("05b5376a-5c76-4f72-bc2c-c148ad068e40"),
		Config: config,
	})
}

// proxies
var SandCoffin = MakeSandCoffin()

func MakeSandCoffin() game.Action {
	action := MakeHumanBoulder()
	action.ID = uuid.MustParse("a351e716-8dc8-4c05-99a2-d0ec4ea87065")
	action.Config.Name = "Sand Coffin"
	action.Config.Power = game.Ptr(55)
	return action
}
