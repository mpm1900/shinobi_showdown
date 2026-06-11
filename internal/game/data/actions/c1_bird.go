package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var C1Bird = MakeC1Bird()

func MakeC1Bird() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("9e8ecd72-8df3-5551-9672-0040d622beb1"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "C1: Bird",
			Description: "+1 priority.",
			Nature:      game.Ptr(game.NsExplosion),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(70),
			Stat:        game.Ptr(game.StatChakraAttack),
			Cost:        game.Ptr(30),
			Cooldown:    game.Ptr(0),
			Jutsu:       game.Ninjutsu,
		}),
		Priority: game.Ptr(game.ActionPriorityP1),
	})
}
