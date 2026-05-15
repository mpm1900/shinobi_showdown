package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var Rasengan = MakeRasengan()

func MakeRasengan() game.Action {
	ID := uuid.MustParse("054eb97a-cd6f-4428-8f54-96d9b6b33bfa")
	config := game.ActionConfig{
		Name:        "Rasengan",
		Nature:      game.Ptr(game.NsSage),
		Accuracy:    game.Ptr(90),
		Power:       game.Ptr(90),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(50),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
	})
}
