package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var Raikiri = MakeRaikiri()

func MakeRaikiri() game.Action {
	ID := uuid.MustParse("fac49e63-a734-41b2-9c50-463c5af92233")
	config := game.ActionConfig{
		Name:        "Raikiri",
		Description: "Deals 30% recoil damage.",
		Nature:      game.Ptr(game.NsLightning),
		Accuracy:    game.Ptr(95),
		Power:       game.Ptr(120),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(70),
		TargetCount: game.Ptr(1),
		Recoil:      game.Ptr(0.3),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
	})
}
