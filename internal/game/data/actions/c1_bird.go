package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var C1Bird = MakeC1Bird()

func MakeC1Bird() game.Action {
	ID := uuid.MustParse("9e8ecd72-8df3-5551-9672-0040d622beb1")
	config := game.ActionConfig{
		Name:        "C1: Bird",
		Description: "+1 priority.",
		Nature:      game.Ptr(game.NsExplosion),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(70),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(30),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:       ID,
		Config:   config,
		Priority: game.Ptr(game.ActionPriorityP1),
	})
}
